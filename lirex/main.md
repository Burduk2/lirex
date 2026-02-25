package lirex

import (
	"strings"
)

type NodeKind int

const (
	LitNode NodeKind = iota
	CaptureNode
	RepeatNode
	GroupNode
	SeqNode
	NotNode
	NotAheadNode
	NotBehindNode
)

type Node interface{}

type litNode struct {
	value string
	raw   bool
}

type LirexpNode struct {
	kind      NodeKind
	value     string       // for literal
	unescaped bool         // for literal
	native    bool         // for built in meta chars
	name      string       // for capture
	children  []LirexpNode // for sequence/repeat
	min, max  uint         // for repeats
	noLimit   bool         // for repeats
}
type SequenceNode LirexpNode
type TextNode LirexpNode

type Options struct {
	CaseInsensitive   bool
	Multiline         bool
	DotMatchesNewline bool
}

var (
	Whitespace    = metaChar(`\s`)
	NonWhitespace = metaChar(`\S`)
	Tab           = metaChar(`\t`)
	Newline       = metaChar(`\n`)

	LowerLatin  = metaChar(`[a-z]`)
	UpperLatin  = metaChar(`[A-Z]`)
	Latin       = metaChar(`[a-zA-Z]`)
	WordChar    = metaChar(`\w`)
	NonWordChar = metaChar(`\W`)
	Digit       = metaChar(`\d`)
	NonDigit    = metaChar(`\D`)
	AnyChar     = metaChar(`.`)
	LineEnd     = metaChar(`$`)
	LineStart   = metaChar(`^`)

	Return          = metaChar(`\r`)
	WordBoundary    = metaChar(`\b`)
	NonWordBoundary = metaChar(`\B`)
)

func metaChar(c string) TextNode {
	return TextNode{kind: LitNode, value: c, unescaped: true, native: true}
}

func Seq(exps ...LirexpNode) SequenceNode {
	return SequenceNode{kind: SeqNode, children: exps}
}

func (exps SequenceNode) Build(opts Options) string {
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

	// parts := make([]string, len(exps))
	// for i, e := range exps {
	// 	parts[i] = string(e)
	// }

	// pattern := fmt.Sprintf(`%s%s`, mode, strings.Join(parts, ""))
	return "pattern"
}

func Lit(str string) TextNode {
	return TextNode{kind: LitNode, value: str}
}
func compileLit(str string) LirexpNode {
	sensitiveChars := `.*+-!?:#<>()[]{}^$|\`
	var builder strings.Builder
	for _, char := range str {
		if strings.ContainsRune(sensitiveChars, char) {
			builder.WriteRune('\\')
		}
		builder.WriteRune(char)
	}

	literal := builder.String()
	if len([]rune(literal)) > 1 {
		literal = "(?:" + literal + ")"
	}
	return LirexpNode{kind: LitNode, value: literal}
}
func UnsafeRaw(str string) LirexpNode {
	return LirexpNode{kind: LitNode, value: str, unescaped: true}
}

func (exp TextNode) AtLeast(n uint) LirexpNode {
	return LirexpNode{kind: RepeatNode, children: []LirexpNode{exp}, min: n, noLimit: true}
}
func (exp TextNode) Exactly(n uint) LirexpNode {
	return LirexpNode{kind: RepeatNode, children: []LirexpNode{exp}, min: n, max: n}
}
func (exp TextNode) Between(from, to uint) LirexpNode {
	return LirexpNode{kind: RepeatNode, children: []LirexpNode{exp}, min: from, max: to}
}

func Capture(name string, exps ...LirexpNode) LirexpNode {
	return LirexpNode{kind: CaptureNode, name: name, children: exps}
}
func Group(exps ...LirexpNode) LirexpNode {
	return LirexpNode{kind: GroupNode, children: exps}
}

func Not(exps ...LirexpNode) LirexpNode {
	return LirexpNode{kind: NotNode, children: exps}
}
func NotAhead(exps ...LirexpNode) LirexpNode {
	return LirexpNode{kind: NotAheadNode, children: exps}
}
func NotBehind(exps ...LirexpNode) LirexpNode {
	return LirexpNode{kind: NotBehindNode, children: exps}
}
