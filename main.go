package main

import (
	"fmt"
	lx "lirex/lirex"
)

func main() {
	expr := lx.Exp(
		lx.Lit("<%"),
		lx.Whitespace.Optional(),
		lx.Capture("insides", lx.AnyChar.AtLeast(1)),
		lx.Whitespace.AtLeast(0),
		lx.Lit("%>"),
		lx.Group().AtLeast(1),
		lx.Or(lx.LowerLatin, lx.UpperLatin).AtLeast(1),
		lx.Not(lx.Whitespace.AtLeast(0)),
	).Build(lx.Options{})

	fmt.Println(expr)
}

// no bc bad for complex expressions
// lx.New().
//   Lit("<%").
//   Whitespace().AtLeast(0).
//   Capture("insides", lx.AnyChar.AtLeast(1)).
//   Whitespace().AtLeast(0).
//   Lit("%>").
//   Build()
