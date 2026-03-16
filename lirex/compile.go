package lirex

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func handleEmptyNode[T Node](node T, ctx *CompileContext) error {
	suffix := fmt.Sprintf("Node has no children or is empty: %s %+v\n", reflect.TypeOf(node), node)
	if ctx.allowRedundant {
		if ctx.showWarnings {
			fmt.Printf("WARNING: %s", suffix)
		}
		return nil
	}
	return fmt.Errorf("Lirex Compile: %s", suffix)
}
func handleOneChildNode[T Node](node T, ctx *CompileContext, child Node, childrenCompiled string) (error, string) {
	if _, yes := child.(SeqNode); yes {
		return nil, "(?:" + childrenCompiled + ")"
	}
	suffix := fmt.Sprintf("%s with only one child: %s %+v\n", reflect.TypeOf(node), reflect.TypeOf(child), child)
	if ctx.allowRedundant {
		if ctx.showWarnings {
			fmt.Printf("WARNING: Group unneeded: %s", suffix)
		}
		return nil, childrenCompiled
	}
	return fmt.Errorf("Lirex Compile: %s", suffix), childrenCompiled
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
func compileLit(str string, isForCharClass bool) string {
	sensitiveChars := `.*+-!?:#<>()[]{}^$|\/`
	if isForCharClass {
		sensitiveChars = `[]-\/`
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
func toRegularNodes[T Node](irrNodes []T) []Node {
	nodes := make([]Node, len(irrNodes))
	for i, child := range irrNodes {
		nodes[i] = child
	}
	return nodes
}

func (node HelperNode) compile(ctx *CompileContext) (string, error) {
	if _, exists := ctx.helpersUsed[node.name]; exists {
		return "", fmt.Errorf("Lirex Compile: Helper '%s' used more than once.", node.name)
	}
	ctx.helpersUsed[node.name] = struct{}{}
	for _, groupName := range node.groupNames {
		if _, exists := ctx.groupNames[groupName]; exists {
			return "", fmt.Errorf("Lirex Compile: Capture group name '%s' is reserved by Lirex. Use another name.", groupName)
		}
		ctx.groupNames[groupName] = struct{}{}
	}

	return node.value, nil
}
func (n MetaCharNode) compile(*CompileContext) (string, error) { return n.value, nil }
func (n RuneCharNode) compile(*CompileContext) (string, error) { return n.value, nil }

func (node SeqNode) compile(ctx *CompileContext) (string, error) {
	childrenCompiled, err := compileNodes(node.nodes, ctx)
	if err != nil {
		return "", err
	}
	if len(node.nodes) == 0 || childrenCompiled == "" {
		return "", handleEmptyNode(node, ctx)
	}
	return childrenCompiled, nil
}

func (node LitNode) compile(ctx *CompileContext) (string, error) {
	val := node.value
	if val == "" {
		return "", handleEmptyNode(node, ctx)
	}
	return compileLit(val, false), nil
}
func (node RawNode) compile(ctx *CompileContext) (string, error) {
	val := node.value
	if val == "" {
		return "", handleEmptyNode(node, ctx)
	}
	_, err := regexp.Compile(val)
	return val, err
}

func (node CaptureNode) compile(ctx *CompileContext) (string, error) {
	children := node.children
	if len(children) == 0 {
		return "", fmt.Errorf("Lirex Compile: Capture() must have have children. Instead Capture (name=%s) has 0 children.", node.name)
	}
	var captureNameExp = Exp(
		LineStart,
		Latin,
		WordChar.ZeroOrMore(),
		LineEnd,
	).MustCompile(Options{})

	name := node.name
	if !captureNameExp.MatchString(name) {
		return "", fmt.Errorf("Lirex Compile: Capture: invalid name for capture group '%s'.", name)
	}
	if _, exists := ctx.groupNames[name]; exists {
		hint := ""
		if _, exists := ReservedGroupNames[name]; exists {
			hint = "\nHint: This name is reserved by Lirex. Use another one."
		}
		return "", fmt.Errorf("Lirex Compile: Capture: duplicate name for capture group '%s'.%s", name, hint)
	}
	ctx.groupNames[name] = struct{}{}

	childrenCompiled, err := compileNodes(toRegularNodes(children), ctx)
	if err != nil {
		return "", err
	}
	if childrenCompiled == "" {
		return "", fmt.Errorf("Lirex Compile: Capture() must have have children. Instead Capture (name=%s) resolved to empty string.", name)
	}
	return "(?P<" + name + ">" + childrenCompiled + ")", nil
}

func (node GroupNode) compile(ctx *CompileContext) (string, error) {
	children := node.children
	childrenCompiled, err := compileNodes(children, ctx)
	if err != nil {
		return "", err
	}
	if childrenCompiled == "" {
		return "", handleEmptyNode(node, ctx)
	}

	if len(children) == 1 {
		err, output := handleOneChildNode(node, ctx, children[0], childrenCompiled)
		if err != nil {
			return "", err
		}
		return output, nil
	}
	return "(?:" + childrenCompiled + ")", nil
}
func (node OrNode) compile(ctx *CompileContext) (string, error) {
	children := node.children
	if len(children) == 0 {
		return "", handleEmptyNode(node, ctx)
	}
	result := []string{}
	for _, child := range children {
		compiled, err := child.compile(ctx)
		if err != nil {
			return "", err
		}
		result = append(result, compiled)
	}

	childrenCompiled := strings.Join(result, "")
	if childrenCompiled == "" {
		return "", handleEmptyNode(node, ctx)
	}
	if len(children) == 1 {
		err, output := handleOneChildNode(node, ctx, children[0], childrenCompiled)
		if err != nil {
			return "", err
		}
		return output, nil
	}
	return "(?:" + strings.Join(result, "|") + ")", nil
}

func (node CharClassNode) compile(ctx *CompileContext) (string, error) {
	children := node.children
	if len(children) == 0 {
		return "", handleEmptyNode(node, ctx)
	}
	var b strings.Builder
	for _, child := range children {
		compiled := ""
		var err error = nil
		switch node := child.(type) {
		case LitNode:
			compiled = compileLit(node.value, true)
		case RuneCharNode:
			val := node.value
			if val[0] == '[' {
				compiled = val[1 : len(val)-1]
			} else {
				compiled = val
			}
		default:
			compiled, err = node.compile(ctx)
		}
		if err != nil {
			return "", err
		}
		b.WriteString(compiled)
	}

	result := b.String()
	if result == "" {
		return "", handleEmptyNode(node, ctx)
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
	childCompiled, err := node.child.compile(ctx)
	if err != nil {
		return "", err
	} else if childCompiled == "" {
		return "", handleEmptyNode(node, ctx)
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
	return childCompiled + q, nil
}
func (node ExactlyRepeatNode) compile(ctx *CompileContext) (string, error) {
	child := node.child
	childCompiled, err := child.compile(ctx)
	if err != nil {
		return "", err
	} else if childCompiled == "" {
		return "", handleEmptyNode(node, ctx)
	}
	num := node.num
	if num == 0 {
		suffix := fmt.Sprintf(".Exactly(0) resolved to empty string on %s %+v\n", reflect.TypeOf(child), child)
		if ctx.allowRedundant {
			if ctx.showWarnings {
				fmt.Printf("WARNING: %s", suffix)
			}
			return "", nil
		}
		return "", fmt.Errorf("Lirex Compile: %s", suffix)
	}
	return childCompiled + fmt.Sprintf("{%d}", num), nil
}
func (node BetweenRepeatNode) compile(ctx *CompileContext) (string, error) {
	child := node.child
	childCompiled, err := child.compile(ctx)
	if err != nil {
		return "", err
	} else if childCompiled == "" {
		return "", handleEmptyNode(node, ctx)
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
			suffix := fmt.Sprintf(".Between(0, 0) resolved to empty string on %s %+v\n", reflect.TypeOf(child), child)
			if ctx.allowRedundant {
				if ctx.showWarnings {
					fmt.Printf("WARNING: %s", suffix)
				}
				return "", nil
			}
			return "", fmt.Errorf("Lirex Compile: %s", suffix)
		}
		q = fmt.Sprintf("{%d}", min)
	} else if min > max {
		return "", fmt.Errorf("Lirex Compile: Between repeat: .Between(%d, %d): min > max", min, max)
	} else {
		q = fmt.Sprintf("{%d,%d}", min, max)
	}

	return childCompiled + q, nil
}
func (node OptionalRepeatNode) compile(ctx *CompileContext) (string, error) {
	childCompiled, err := node.child.compile(ctx)
	if err != nil {
		return "", err
	} else if childCompiled == "" {
		return "", handleEmptyNode(node, ctx)
	}
	return childCompiled + "?", nil
}
