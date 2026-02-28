package lirex

import (
	"fmt"
	"regexp"
	"strings"
)

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
func toRegularNodes[T Node](irrNodes []T) []Node {
	nodes := make([]Node, len(irrNodes))
	for i, child := range irrNodes {
		nodes[i] = child
	}
	return nodes
}

func (n SeqNode) compile(ctx *CompileContext) (string, error) {
	return compileNodes(n.nodes, ctx)
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
func (n MetaCharNode) compile(ctx *CompileContext) (string, error) {
	return n.value, nil
}
func (n RuneCharNode) compile(ctx *CompileContext) (string, error) {
	return n.value, nil
}
func (n RawNode) compile(ctx *CompileContext) (string, error) {
	val := n.value
	_, err := regexp.Compile(val)
	return val, err
}

func (node GroupNode) compile(ctx *CompileContext) (string, error) {
	children := node.children
	if len(children) == 0 {
		return "", nil
	}
	childrenCompiled, err := compileNodes(children, ctx)
	return "(?:" + childrenCompiled + ")", err
}

func (node CaptureNode) compile(ctx *CompileContext) (string, error) {
	children := node.children
	if len(children) == 0 {
		return "", nil
	}
	expr := Exp(
		LineStart,
		Latin,
		WordChar.ZeroOrMore(),
		LineEnd,
	).MustCompile(Options{})

	name := node.name
	if !expr.MatchString(name) {
		return "", fmt.Errorf("Capture: invalid name for capture group '%s'", name)
	}
	if _, exists := ctx.groupNames[name]; exists {
		return "", fmt.Errorf("Capture: duplicate name for capture group '%s'", name)
	}
	ctx.groupNames[name] = struct{}{}

	childrenCompiled, err := compileNodes(toRegularNodes(children), ctx)
	return "(?P<" + name + ">" + childrenCompiled + ")", err
}

func (node OrNode) compile(ctx *CompileContext) (string, error) {
	children := node.children
	if len(children) == 0 {
		return "", nil
	}
	result := []string{}
	for _, child := range children {
		compiled, err := child.compile(ctx)
		if err != nil {
			return "", err
		}
		result = append(result, compiled)
	}
	return "(?:" + strings.Join(result, "|") + ")", nil
}

func (node CharClassNode) compile(ctx *CompileContext) (string, error) {
	children := node.children
	if len(children) == 0 {
		return "", nil
	}
	var b strings.Builder
	for _, child := range children {
		compiled := ""
		var err error = nil
		switch n := child.(type) {
		case LitNode:
			compiled = compileLit(n.value, true)
		case RuneCharNode:
			val := n.value
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
		b.WriteString(compiled)
	}

	result := b.String()
	negation := ""
	if node.negate {
		negation = "^"
	} else if result[0] == '^' {
		result = "\\" + result
	}
	return "[" + negation + result + "]", nil
}

func (node AtLeastRepeatNode) compile(ctx *CompileContext) (string, error) {
	childCompiled, err := node.child.compile(ctx)
	if childCompiled == "" {
		return "", nil
	}
	q := ""
	switch num := node.num; num {
	case 0:
		q = "*"
	case 1:
		q = "+"
	default:
		q = fmt.Sprintf("{%d,}", num)
	}
	return childCompiled + q, err
}
func (node ExactlyRepeatNode) compile(ctx *CompileContext) (string, error) {
	childCompiled, err := node.child.compile(ctx)
	if childCompiled == "" {
		return "", nil
	}
	num := node.num
	if num == 0 {
		return "", nil
	}
	return childCompiled + fmt.Sprintf("{%d}", num), err
}
func (node BetweenRepeatNode) compile(ctx *CompileContext) (string, error) {
	childCompiled, err := node.child.compile(ctx)
	if childCompiled == "" || err != nil {
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
		if ctx.showWarnings {
			fmt.Printf("WARNING: .Between(%d, %d) => Could use .Exactly(%d) instead.\n", min, min, min)
		}
		if min == 0 {
			return "", nil
		}
		q = fmt.Sprintf("{%d}", min)
	} else if min > max {
		return "", fmt.Errorf("Between repeat: min > max")
	} else {
		q = fmt.Sprintf("{%d,%d}", min, max)
	}

	return childCompiled + q, nil
}
func (node OptionalRepeatNode) compile(ctx *CompileContext) (string, error) {
	childCompiled, err := node.child.compile(ctx)
	if childCompiled == "" {
		return "", nil
	}
	return childCompiled + "?", err
}
