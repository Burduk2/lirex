package lirex

type Node interface {
	compile(*CompileContext) (string, error)
	explain() string
}

type Capturable interface {
	Node
	capturableNode()
}
type CharClassable interface {
	Node
	charClassableNode()
}
type Repeatable interface {
	Node
	repeatableNode()
}

func (LitNode) repeatableNode()       {}
func (MetaCharNode) repeatableNode()  {}
func (RuneCharNode) repeatableNode()  {}
func (GroupNode) repeatableNode()     {}
func (OrNode) repeatableNode()        {}
func (CharClassNode) repeatableNode() {}

func (LitNode) capturableNode()            {}
func (MetaCharNode) capturableNode()       {}
func (RuneCharNode) capturableNode()       {}
func (RawNode) capturableNode()            {}
func (GroupNode) capturableNode()          {}
func (OrNode) capturableNode()             {}
func (CharClassNode) capturableNode()      {}
func (AtLeastRepeatNode) capturableNode()  {}
func (ExactlyRepeatNode) capturableNode()  {}
func (BetweenRepeatNode) capturableNode()  {}
func (OptionalRepeatNode) capturableNode() {}
func (CaptureNode) capturableNode()        {}
func (SeqNode) capturableNode()            {}

func (LitNode) charClassableNode()      {}
func (RuneCharNode) charClassableNode() {}

// SEQUENCE ----------------------------------------------------------------------------
type SeqNode struct {
	nodes []Node
}

func Seq(nodes ...Node) SeqNode {
	return SeqNode{nodes: nodes}
}

// TEXT ---------------------------------------------------------------------------------
type RawNode struct {
	value string
}
type LitNode struct {
	value string
}

// Raw string (escaped)
func Lit(s string) LitNode {
	return LitNode{value: s}
}

// NOT RECOMMENDED. Raw string (unescaped)
func UnsafeRaw(s string) RawNode {
	return RawNode{value: s}
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
func Group(nodes ...Node) GroupNode {
	return GroupNode{children: nodes}
}

// OR -------------------------------------------------------------------------------------
type OrNode struct {
	children []Capturable
}

// Regex equivalent: (...|...)
func Or(nodes ...Capturable) OrNode {
	return OrNode{children: nodes}
}

// CHAR CLASS -----------------------------------------------------------------------------
type CharClassNode struct {
	children []CharClassable
	negate   bool
}

// Regex equivalent: [...]
func CharClass(nodes ...CharClassable) CharClassNode {
	return CharClassNode{children: nodes}
}

// Regex equivalent: [^...]
func NotCharClass(nodes ...CharClassable) CharClassNode {
	return CharClassNode{children: nodes, negate: true}
}

// REPEAT ---------------------------------------------------------------------------------
type AtLeastRepeatNode struct {
	child Repeatable
	num   uint
}
type ExactlyRepeatNode struct {
	child Repeatable
	num   uint
}
type BetweenRepeatNode struct {
	child Repeatable
	min   uint
	max   uint
}
type OptionalRepeatNode struct {
	child Repeatable
}

func atLeast[T Repeatable](node T, n uint) AtLeastRepeatNode {
	return AtLeastRepeatNode{child: node, num: n}
}
func zeroOrMore[T Repeatable](node T) AtLeastRepeatNode {
	return AtLeastRepeatNode{child: node, num: 0}
}
func exactly[T Repeatable](node T, n uint) ExactlyRepeatNode {
	return ExactlyRepeatNode{child: node, num: n}
}
func between[T Repeatable](node T, from, to uint) BetweenRepeatNode {
	return BetweenRepeatNode{child: node, min: from, max: to}
}
func optional[T Repeatable](node T) OptionalRepeatNode {
	return OptionalRepeatNode{child: node}
}

// Regex equivalent: ...* || ...+ || ...{n,}
func (node LitNode) AtLeast(n uint) AtLeastRepeatNode {
	return atLeast(node, n)
}
func (node MetaCharNode) AtLeast(n uint) AtLeastRepeatNode {
	return atLeast(node, n)
}
func (node RuneCharNode) AtLeast(n uint) AtLeastRepeatNode {
	return atLeast(node, n)
}
func (node GroupNode) AtLeast(n uint) AtLeastRepeatNode {
	return atLeast(node, n)
}
func (node OrNode) AtLeast(n uint) AtLeastRepeatNode {
	return atLeast(node, n)
}
func (node CharClassNode) AtLeast(n uint) AtLeastRepeatNode {
	return atLeast(node, n)
}

// Regex equivalent: ...*
func (node LitNode) ZeroOrMore() AtLeastRepeatNode {
	return zeroOrMore(node)
}
func (node MetaCharNode) ZeroOrMore() AtLeastRepeatNode {
	return zeroOrMore(node)
}
func (node RuneCharNode) ZeroOrMore() AtLeastRepeatNode {
	return zeroOrMore(node)
}
func (node GroupNode) ZeroOrMore() AtLeastRepeatNode {
	return zeroOrMore(node)
}
func (node OrNode) ZeroOrMore() AtLeastRepeatNode {
	return zeroOrMore(node)
}
func (node CharClassNode) ZeroOrMore() AtLeastRepeatNode {
	return zeroOrMore(node)
}

// Regex equivalent: ...{n}
func (node LitNode) Exactly(n uint) ExactlyRepeatNode {
	return exactly(node, n)
}
func (node MetaCharNode) Exactly(n uint) ExactlyRepeatNode {
	return exactly(node, n)
}
func (node RuneCharNode) Exactly(n uint) ExactlyRepeatNode {
	return exactly(node, n)
}
func (node GroupNode) Exactly(n uint) ExactlyRepeatNode {
	return exactly(node, n)
}
func (node OrNode) Exactly(n uint) ExactlyRepeatNode {
	return exactly(node, n)
}
func (node CharClassNode) Exactly(n uint) ExactlyRepeatNode {
	return exactly(node, n)
}

// Regex equivalent: ...{n,m} || ...?
func (node LitNode) Between(from, to uint) BetweenRepeatNode {
	return between(node, from, to)
}
func (node MetaCharNode) Between(from, to uint) BetweenRepeatNode {
	return between(node, from, to)
}
func (node RuneCharNode) Between(from, to uint) BetweenRepeatNode {
	return between(node, from, to)
}
func (node GroupNode) Between(from, to uint) BetweenRepeatNode {
	return between(node, from, to)
}
func (node OrNode) Between(from, to uint) BetweenRepeatNode {
	return between(node, from, to)
}
func (node CharClassNode) Between(from, to uint) BetweenRepeatNode {
	return between(node, from, to)
}

// Regex equivalent: ...?
func (node LitNode) Optional() OptionalRepeatNode {
	return optional(node)
}
func (node MetaCharNode) Optional() OptionalRepeatNode {
	return optional(node)
}
func (node RuneCharNode) Optional() OptionalRepeatNode {
	return optional(node)
}
func (node GroupNode) Optional() OptionalRepeatNode {
	return optional(node)
}
func (node OrNode) Optional() OptionalRepeatNode {
	return optional(node)
}
func (node CharClassNode) Optional() OptionalRepeatNode {
	return optional(node)
}
