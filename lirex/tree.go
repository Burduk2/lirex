package lirex

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type ExpTreeNode []Node

func Exp(nodes ...Node) ExpTreeNode {
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
	AllowRedundant    bool
}
type CompileContext struct {
	groupNames     map[string]struct{}
	helpersUsed    map[string]struct{}
	showWarnings   bool
	allowRedundant bool
}
type ExplainContext struct {
	indent uint
}

func (tree ExpTreeNode) Compile(opts Options) (*regexp.Regexp, error) {
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

	ctx := &CompileContext{
		groupNames:     make(map[string]struct{}),
		helpersUsed:    make(map[string]struct{}),
		showWarnings:   opts.ShowWarnings,
		allowRedundant: opts.AllowRedundant,
	}
	result, err := compileNodes(tree, ctx)
	if err != nil {
		return nil, err
	}
	if result == "" {
		return nil, fmt.Errorf("Lirex Compile: Expression resolved to empty string.")
	}
	return regexp.Compile(mode + result)
}
func (tree ExpTreeNode) MustCompile(opts Options) *regexp.Regexp {
	result, err := tree.Compile(opts)
	if err != nil {
		panic(err)
	}
	return result
}

func yellow(s string) string {
	return "\033[33m" + s + "\033[0m"
}
func red(s string) string {
	return "\033[31m" + s + "\033[0m"
}
func green(s string) string {
	return "\033[32m" + s + "\033[0m"
}
func bold(s string) string {
	return "\033[1m" + s + "\033[0m"
}
func (tree ExpTreeNode) Explain(opts Options) {
	ctx := &CompileContext{groupNames: make(map[string]struct{}), showWarnings: true}

	var b strings.Builder
	for _, child := range tree {
		t := reflect.TypeOf(child)
		prefix := t.String()
		compiled, err := child.compile(ctx)
		if err != nil {
			compiled = red("Failed to compile: " + err.Error())
		}
		fmt.Println(compiled)
		b.WriteString(prefix + ": " + bold(compiled) + "\n")
	}

	// ctx := &CompileContext{
	// 	groupNames:     make(map[string]struct{}),
	// 	helpersUsed:    make(map[string]struct{}),
	// 	showWarnings:   opts.ShowWarnings,
	// 	allowRedundant: opts.AllowRedundant,
	// }
	result, err := compileNodes(tree, ctx)
	if err != nil {
	}
	if result == "" {
	}
	// return regexp.Compile(mode + result)
	fmt.Println(b.String())
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
