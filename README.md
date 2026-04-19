# lirex

`lirex` is a small Go library for building regular expressions from composable typed nodes instead of hand-writing regex strings.

It targets Go's standard `regexp` package and helps with:

- safer literal escaping via `Lit(...)`
- readable composition with `Exp`, `Group`, `Or`, `CharClass`, and `Capture`
- reusable helper expressions such as email, URL, phone, and domain
- extracting named captures with `FindCaptures`

## Status

This repository currently exposes the library code directly and uses the module path from [go.mod](/Volumes/Au/code/Libs/lirex/go.mod). In this repo, examples import the package as:

```go
import lx "lirex/lirex"
```

If you publish this module under a full remote path later, update `go.mod` and imports accordingly.

## Installation

Add the package to your project in the usual Go way once the module path matches where you host it.

For local development in this repository:

```go
import lx "lirex/lirex"
```

## Quick Start

```go
package main

import (
	"fmt"
	lx "lirex/lirex"
)

func main() {
	re := lx.Exp(
		lx.LineStart,
		lx.Capture("name",
			lx.Latin.Between(1, 20),
		),
		lx.Lit("@"),
		lx.Capture("domain",
			lx.Helpers.Domain,
		),
		lx.LineEnd,
	).MustCompile(lx.Options{})

	fmt.Println(re.String())
	fmt.Println(re.MatchString("alice@example.com"))
}
```

Produces a compiled `*regexp.Regexp` using Go's standard library.

## Core Concepts

Build expressions with `Exp(...)`:

```go
re, err := lx.Exp(
	lx.LineStart,
	lx.Lit("item-"),
	lx.Digit.Exactly(4),
	lx.LineEnd,
).Compile(lx.Options{})
```

Use `Lit(...)` for escaped text and `UnsafeRaw(...)` only when you intentionally want to inject raw regex:

```go
lx.Lit("a+b")        // escaped
lx.UnsafeRaw(`a+b`)  // raw regex
```

Compose alternatives and groups:

```go
lx.Or(
	lx.Lit("cat"),
	lx.Lit("dog"),
)

lx.Group(
	lx.Lit("http"),
	lx.Lit("s").Optional(),
)
```

Build character classes:

```go
lx.CharClass(lx.Latin, lx.Digit, lx.Lit("_-"))
lx.NotCharClass(lx.Newline, lx.Return)
```

Add repetition to repeatable nodes:

```go
lx.Digit.AtLeast(1)       // +
lx.Digit.ZeroOrMore()     // *
lx.Digit.Exactly(4)       // {4}
lx.Digit.Between(2, 5)    // {2,5}
lx.Lit("-").Optional()    // ?
```

## Named Captures

Use `Capture(name, ...)` to create named groups:

```go
re := lx.Exp(
	lx.Capture("year", lx.Digit.Exactly(4)),
	lx.Lit("-"),
	lx.Capture("month", lx.Digit.Exactly(2)),
	lx.Lit("-"),
	lx.Capture("day", lx.Digit.Exactly(2)),
).MustCompile(lx.Options{})
```

Then extract all named matches with `FindCaptures`:

```go
captures, ok := lx.FindCaptures(re, "2026-04-19 2027-01-05")
if ok {
	fmt.Println(captures["year"])  // [2026 2027]
	fmt.Println(captures["month"]) // [04 01]
}
```

## Built-in Helpers

The package exposes reusable helper patterns under `lx.Helpers`:

- `Helpers.Domain`
- `Helpers.Email`
- `Helpers.InternationalPhone`
- `Helpers.CreditCard`
- `Helpers.FullUrl`

Example:

```go
re := lx.Exp(
	lx.Helpers.Email,
).MustCompile(lx.Options{
	CaseInsensitive: true,
})
```

Helpers reserve capture group names internally. If you define your own captures, avoid the reserved helper names such as `Email`, `Domain`, `Phone`, and `FullUrl`.

## Available Character and Meta Nodes

Common predefined nodes include:

- `Whitespace`, `NonWhitespace`, `Tab`, `Newline`, `LineBreak`
- `Latin`, `LowerLatin`, `UpperLatin`, `LatinDigit`
- `Letter`, `UpperLetter`, `LowerLetter`
- `Digit`, `NonDigit`, `WordChar`, `NonWordChar`
- `AnyChar`, `LineStart`, `LineEnd`, `WordBoundary`

Several Unicode-oriented script and class nodes are also available, including `ExtendedLatin`, `Cyrillic`, `Greek`, `Arabic`, `Hebrew`, `Han`, `AnyDecimal`, `NumberLike`, `Punctuation`, and `Symbol`.

## Compile Options

`Options` controls regexp flags and compiler behavior:

```go
type Options struct {
	CaseInsensitive bool
	Multiline bool
	DotMatchesNewline bool
	ShowWarnings bool
	AllowRedundant bool
}
```

- `CaseInsensitive` adds `(?i)`
- `Multiline` adds `(?m)`
- `DotMatchesNewline` adds `(?s)`
- `ShowWarnings` prints warnings for redundant constructs when allowed
- `AllowRedundant` permits empty or unnecessary group-like constructs that would otherwise return errors

## Notes and Current Limitations

- The package compiles to Go's `regexp` engine semantics.
- `UnsafeRaw(...)` validates the raw fragment by compiling it, but it still bypasses `lirex` escaping guarantees.
- `Explain(...)` exists, but its structured explanation output is not implemented yet.
- Helper nodes are intended to be used once per expression; reusing the same helper multiple times in one expression currently raises an error.

## License

See [LICENSE](/Volumes/Au/code/Libs/lirex/LICENSE).
