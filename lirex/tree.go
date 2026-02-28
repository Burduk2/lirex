package lirex

import (
	"regexp"
)

type ExpTree []Node

func Exp(nodes ...Node) ExpTree {
	return nodes
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

func (tree ExpTree) Compile(opts Options) (*regexp.Regexp, error) {
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
		return nil, err
	}
	return regexp.Compile(mode + result)
}
func (tree ExpTree) MustCompile(opts Options) *regexp.Regexp {
	result, err := tree.Compile(opts)
	if err != nil {
		panic(err)
	}
	return result
}

func (tree ExpTree) Explain() {

}

func FindCaptures(re *regexp.Regexp, str string) (map[string][]string, bool) {
	all := re.FindAllStringSubmatch(str, -1)
	if len(all) == 0 {
		return nil, false
	}
	names := re.SubexpNames()
	if len(names) <= 1 {
		return nil, false
	}

	captures := make(map[string][]string, len(names)-1)
	for _, match := range all {
		for i := 1; i < len(match); i++ { // skip index 0 (whole match)
			name := names[i]
			captures[name] = append(captures[name], match[i])
		}
	}

	return captures, true
}
