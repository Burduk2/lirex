package lirex

import (
	"fmt"
	"regexp"
	"strings"
)

func (n ComponentNode) compile(ctx *CompileContext) (string, error) {
	return compileNodes(n.nodes, ctx)
}

func compileNodes(nodes []Node, ctx *CompileContext) (string, error) {
	var b strings.Builder
	for _, child := range nodes {
		compiled, err := child.compile(ctx)
		if err != nil {
			return "", err
		}
		b.WriteString(compiled)
	}
	return b.String(), nil
}
func compileCapturableNodes(captNodes []Capturable, ctx *CompileContext) (string, error) {
	nodes := make([]Node, len(captNodes))
	for i, child := range captNodes {
		nodes[i] = child
	}
	return compileNodes(nodes, ctx)
}

func (n ExprNode) compile(ctx *CompileContext) (string, error) {
	result, err := n.node.compile(ctx)
	return result, err
}

func (node LitNode) compile(ctx *CompileContext) (string, error) {
	return compileLit(node.value, false), nil
}
func compileLit(str string, isForCharClass bool) string {
	sensitiveChars := `.*+-!?:#<>()[]{}^$|\`
	if isForCharClass {
		sensitiveChars = `[]-\`
	}

	var builder strings.Builder
	for _, char := range str {
		if strings.ContainsRune(sensitiveChars, char) {
			builder.WriteRune('\\')
		}
		builder.WriteRune(char)
	}

	literal := builder.String()
	if !isForCharClass && len([]rune(literal)) > 1 {
		if !(len([]rune(literal)) == 2 && literal[0] == '\\') {
			literal = "(?:" + literal + ")"
		}
	}
	return literal
}
func (node MetaCharNode) compile(ctx *CompileContext) (string, error) {
	return node.value, nil
}
func (node RawNode) compile(ctx *CompileContext) (string, error) {
	_, err := regexp.Compile(node.value)
	return node.value, err
}

func (node GroupNode) compile(ctx *CompileContext) (string, error) {
	childrenCompiled, err := compileNodes(node.children, ctx)
	return "(?:" + childrenCompiled + ")", err
}

func (node CaptureNode) compile(ctx *CompileContext) (string, error) {
	expr := Exp(
		LineStart,
		Or(LowerLatin, UpperLatin),
		CharClass(WordChar).ZeroOrMore(),
		LineEnd,
	).MustBuild(Options{})

	name := node.name
	if !regexp.MustCompile(expr).MatchString(name) {
		return "", fmt.Errorf("Capture: invalid name for capture group '%s'", name)
	}
	if _, exists := ctx.groupNames[name]; exists {
		return "", fmt.Errorf("Capture: duplicate name for capture group '%s'", name)
	}
	ctx.groupNames[name] = struct{}{}

	childrenCompiled, err := compileCapturableNodes(node.children, ctx)
	return "(?P<" + name + ">" + childrenCompiled + ")", err
}

func (node OrNode) compile(ctx *CompileContext) (string, error) {
	result := []string{}
	for _, child := range node.children {
		compiled, err := child.compile(ctx)
		if err != nil {
			return "", err
		}
		result = append(result, compiled)
	}
	return "(?:" + strings.Join(result, "|") + ")", nil
}

func (node CharClassNode) compile(ctx *CompileContext) (string, error) {
	result := ""
	for _, child := range node.children {
		if _, ok := child.node.(CharClassable); !ok {
			return "", fmt.Errorf("CharClass: child of type %T cannot be put in CharClass", child.node)
		}
		compiled := ""
		var err error = nil
		switch n := child.node.(type) {
		case LitNode:
			compiled = compileLit(n.value, true)
		case MetaCharNode:
			val := n.value
			if !n.charClassable {
				return "", fmt.Errorf("CharClass: child of type %T('%s') cannot be put in CharClass", child.node, val)
			}
			if val[0] == '[' {
				compiled = val[1 : len(val)-1]
			} else {
				compiled = val
			}
		default:
			compiled, err = n.compile(ctx)
		}
		if err != nil {
			return "", err
		}
		result += compiled
	}

	negation := ""
	if node.negate {
		negation = "^"
	} else if result[0] == '^' {
		result = "\\" + result
	}
	return "[" + negation + result + "]", nil
}

func (node AtLeastRepeatNode) compile(ctx *CompileContext) (string, error) {
	q := ""
	switch num := node.num; num {
	case 0:
		q = "*"
	case 1:
		q = "+"
	default:
		q = fmt.Sprintf("{%d,}", num)
	}
	childCompiled, err := node.child.compile(ctx)
	return childCompiled + q, err
}
func (node ExactlyRepeatNode) compile(ctx *CompileContext) (string, error) {
	childCompiled, err := node.child.compile(ctx)
	return childCompiled + fmt.Sprintf("{%d}", node.num), err
}
func (node BetweenRepeatNode) compile(ctx *CompileContext) (string, error) {
	childCompiled, err := node.child.compile(ctx)
	if err != nil {
		return "", err
	}

	q := ""
	min, max := node.min, node.max
	if min == 0 && max == 1 {
		if ctx.showWarnings {
			fmt.Println("WARNING: .Between(0, 1) => Could use .Optional() instead.")
		}
		q = "?"
	} else if min == max {
		num := min
		if ctx.showWarnings {
			fmt.Printf("WARNING: .Between(%d, %d) => Could use .Exactly(%d) instead.\n", num, num, num)
		}
		q = fmt.Sprintf("{%d}", num)
	} else if min > max {
		return "", fmt.Errorf("Between repeat: min > max")
	} else {
		q = fmt.Sprintf("{%d,%d}", min, max)
	}

	return childCompiled + q, nil
}
func (node OptionaRepeatNode) compile(ctx *CompileContext) (string, error) {
	childCompiled, err := node.child.compile(ctx)
	return childCompiled + "?", err
}
