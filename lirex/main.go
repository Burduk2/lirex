package lirex

type Node interface {
	compile() string
}

type Repeatable interface {
	Node
	repeatNode()
}

type Groupable interface {
	Node
	groupNode()
}

type Negatable interface {
	Node
	negNode()
}

type ExprNode struct {
	node Node
}

// TEXT ---------------------------------------------------------------------------------
type TextNode struct {
	value     string
	unescaped bool
	native    bool
}

func (TextNode) repeatNode() {}

func (t TextNode) compile() string {
	return ""
}

func Lit(s string) TextNode {
	return TextNode{value: s}
}

var (
	// \s
	Whitespace = metaChar(`\s`)
	// \S
	NonWhitespace = metaChar(`\S`)
	// \t
	Tab = metaChar(`\t`)
	// \n
	Newline = metaChar(`\n`)

	// [a-z]
	LowerLatin = metaChar(`[a-z]`)
	// [A-Z]
	UpperLatin = metaChar(`[A-Z]`)
	// [a-zA-Z]
	Latin = metaChar(`[a-zA-Z]`)
	// \w
	WordChar = metaChar(`\w`)
	// \W
	NonWordChar = metaChar(`\W`)
	// \d
	Digit = metaChar(`\d`)
	// \D
	NonDigit = metaChar(`\D`)
	// .
	AnyChar = metaChar(`.`)
	// $
	LineEnd = metaChar(`$`)
	// ^
	LineStart = metaChar(`^`)

	// \r
	Return = metaChar(`\r`)
	// \b
	WordBoundary = metaChar(`\b`)
	// \B
	NonWordBoundary = metaChar(`\B`)
)

func metaChar(ch string) TextNode {
	return TextNode{value: ch, unescaped: true, native: true}
}

// EXP TREE ---------------------------------------------------------------------------------
type ExpTree []Node

func Exp(nodes ...Node) ExpTree {
	return nodes
}

type Options struct {
	CaseInsensitive   bool
	Multiline         bool
	DotMatchesNewline bool
}

func (tree ExpTree) Build(opts Options) string {
	return ""
}

// CAPTURE ---------------------------------------------------------------------------------
type CaptureNode struct {
	name     string
	children []Repeatable
}

func (n CaptureNode) compile() string {
	return ""
}
func Capture(name string, nodes ...Repeatable) CaptureNode {
	return CaptureNode{name: name, children: nodes}
}

// GROUP ---------------------------------------------------------------------------------
type GroupNode struct {
	children []Node
}

func (n GroupNode) repeatNode() {}
func (n GroupNode) compile() string {
	return ""
}
func Group(nodes ...Node) GroupNode {
	return GroupNode{children: nodes}
}

// NEG ---------------------------------------------------------------------------------
type negNodeKind int

const (
	notNode negNodeKind = iota
	notAheadNode
	notBehindNode
)

type NegNode struct {
	kind     negNodeKind
	children []Node
}

func (n NegNode) repeatNode() {}
func (n NegNode) compile() string {
	return ""
}

func Not(nodes ...Node) NegNode {
	return NegNode{kind: notNode, children: nodes}
}
func NotAhead(nodes ...Node) NegNode {
	return NegNode{kind: notAheadNode, children: nodes}
}
func NotBehind(nodes ...Node) NegNode {
	return NegNode{kind: notBehindNode, children: nodes}
}

// OR -------------------------------------------------------------------------------------
type OrNode struct {
	children []Node
}

func (n OrNode) repeatNode() {}
func (n OrNode) compile() string {
	return ""
}

func Or(nodes ...Node) OrNode {
	return OrNode{children: nodes}
}

// REPEAT ---------------------------------------------------------------------------------
type RepeatNode struct {
	child     Node
	min       uint
	max       uint
	unlimited bool
}

func (r RepeatNode) repeatNode() {}
func (r RepeatNode) compile() string {
	return ""
}

func (node TextNode) AtLeast(n uint) RepeatNode {
	return RepeatNode{child: node, min: n, unlimited: true}
}
func (node TextNode) Exactly(n uint) RepeatNode {
	return RepeatNode{child: node, min: n, max: n}
}
func (node TextNode) Between(from, to uint) RepeatNode {
	return RepeatNode{child: node, min: from, max: to}
}
func (node TextNode) Optional() RepeatNode {
	return RepeatNode{child: node, min: 0, unlimited: true}
}
