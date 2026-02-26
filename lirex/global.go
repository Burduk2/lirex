package lirex

type Node interface {
	compile(*CompileContext) (string, error)
}

type Capturable interface {
	Node
	capturableNode()
}

func (ExprNode) capturableNode()          {}
func (AtLeastRepeatNode) capturableNode() {}
func (ExactlyRepeatNode) capturableNode() {}
func (BetweenRepeatNode) capturableNode() {}
func (OptionaRepeatNode) capturableNode() {}

type CharClassable interface {
	charClassableNode()
}

func (LitNode) charClassableNode()      {}
func (MetaCharNode) charClassableNode() {}

type ExprNode struct {
	node Node
}

// EXP TREE ---------------------------------------------------------------------------------
type ExpTree []Node
type ComponentNode struct {
	nodes []Node
}

func Exp(nodes ...Node) ExpTree {
	return nodes
}
func Component(nodes ...Node) ComponentNode {
	return ComponentNode{nodes: nodes}
}

type Options struct {
	// (?i)
	CaseInsensitive bool
	// (?m)
	Multiline bool
	// (?s)
	DotMatchesNewline bool
	ShowWarnings      bool
}
type CompileContext struct {
	groupNames   map[string]struct{}
	showWarnings bool
}

func (tree ExpTree) Build(opts Options) (string, error) {
	mode := ""
	if opts.CaseInsensitive {
		mode += "i"
	}
	if opts.Multiline {
		mode += "m"
	}
	if opts.DotMatchesNewline {
		mode += "s"
	}
	if mode != "" {
		mode = "(?" + mode + ")"
	}

	ctx := &CompileContext{groupNames: make(map[string]struct{}), showWarnings: opts.ShowWarnings}
	result, err := compileNodes(tree, ctx)
	if err != nil {
		return "", err
	}
	return mode + result, nil
}
func (tree ExpTree) MustBuild(opts Options) string {
	result, err := tree.Build(opts)
	if err != nil {
		panic(err)
	}
	return result
}
func (tree ExpTree) Explain() {

}

// TEXT ---------------------------------------------------------------------------------
type RawNode struct {
	value string
}
type MetaCharNode struct {
	value         string
	charClassable bool
}

type LitNode struct {
	value string
}

var (
	// Regex equivalent: \s
	Whitespace = metaChar(`\s`, true)
	// Regex equivalent: \S
	NonWhitespace = metaChar(`\S`, true)
	// Regex equivalent: \t
	Tab = metaChar(`\t`, true)
	// Regex equivalent: \n
	Newline = metaChar(`\n`, true)

	// Regex equivalent: [a-z]
	LowerLatin = metaChar(`[a-z]`, true)
	// Regex equivalent: [A-Z]
	UpperLatin = metaChar(`[A-Z]`, true)
	// Regex equivalent: [a-zA-Z]
	Latin = metaChar(`[a-zA-Z]`, true)
	// Regex equivalent: \w ([a-zA-Z0-9_])
	WordChar = metaChar(`\w`, true)
	// Regex equivalent: \W
	NonWordChar = metaChar(`\W`, true)
	// Regex equivalent: \d
	Digit = metaChar(`\d`, true)
	// Regex equivalent: \D
	NonDigit = metaChar(`\D`, true)
	// Regex equivalent: .
	AnyChar = metaChar(`.`, false)
	// Regex equivalent: $
	LineEnd = metaChar(`$`, false)
	// Regex equivalent: ^
	LineStart = metaChar(`^`, false)

	// Regex equivalent: \r
	Return = metaChar(`\r`, true)
	// Regex equivalent: \b
	WordBoundary = metaChar(`\b`, true)
	// Regex equivalent: \B
	NonWordBoundary = metaChar(`\B`, true)
)

func metaChar(ch string, charClassable bool) ExprNode {
	return ExprNode{node: MetaCharNode{value: ch, charClassable: charClassable}}
}

// Raw string (escaped)
func Lit(s string) ExprNode {
	return ExprNode{node: LitNode{value: s}}
}

// NOT RECOMMENDED. Raw string (unescaped)
func UnsafeRaw(s string) ExprNode {
	return ExprNode{node: RawNode{value: s}}
}

// CAPTURE ---------------------------------------------------------------------------------
type CaptureNode struct {
	name     string
	children []Capturable
}

// Regex equivalent: (?P<name>...)
func Capture(name string, nodes ...Capturable) CaptureNode {
	return CaptureNode{name: name, children: nodes}
}

// GROUP ---------------------------------------------------------------------------------
type GroupNode struct {
	children []Node
}

// Regex equivalent: (?:...)
func Group(nodes ...Node) ExprNode {
	return ExprNode{node: GroupNode{children: nodes}}
}

// OR -------------------------------------------------------------------------------------
type OrNode struct {
	children []Capturable
}

// Regex equivalent: (...|...)
func Or(nodes ...Capturable) ExprNode {
	return ExprNode{node: OrNode{children: nodes}}
}

// CHAR CLASS -----------------------------------------------------------------------------
type CharClassNode struct {
	children []ExprNode
	negate   bool
}

// Regex equivalent: [...]
func CharClass(nodes ...ExprNode) ExprNode {
	return ExprNode{node: CharClassNode{children: nodes}}
}

// Regex equivalent: [^...]
func NotCharClass(nodes ...ExprNode) ExprNode {
	return ExprNode{node: CharClassNode{children: nodes, negate: true}}
}

// REPEAT ---------------------------------------------------------------------------------
type AtLeastRepeatNode struct {
	child ExprNode
	num   uint
}
type ExactlyRepeatNode struct {
	child ExprNode
	num   uint
}
type BetweenRepeatNode struct {
	child ExprNode
	min   uint
	max   uint
}
type OptionaRepeatNode struct {
	child ExprNode
}

// Regex equivalent: ...* || ...+ || ...{n,}
func (node ExprNode) AtLeast(n uint) AtLeastRepeatNode {
	return AtLeastRepeatNode{child: node, num: n}
}

// Regex equivalent: ...*
func (node ExprNode) ZeroOrMore() AtLeastRepeatNode {
	return AtLeastRepeatNode{child: node, num: 0}
}

// Regex equivalent: ...{n}
func (node ExprNode) Exactly(n uint) ExactlyRepeatNode {
	return ExactlyRepeatNode{child: node, num: n}
}

// Regex equivalent: ...{n,m} || ...?
func (node ExprNode) Between(from, to uint) BetweenRepeatNode {
	return BetweenRepeatNode{child: node, min: from, max: to}
}

// Regex equivalent: ...?
func (node ExprNode) Optional() OptionaRepeatNode {
	return OptionaRepeatNode{child: node}
}
