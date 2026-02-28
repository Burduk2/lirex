package lirex

type MetaCharNode struct {
	value string
}
type RuneCharNode struct {
	value string
}

var (
	// Regex equivalent: \s
	Whitespace = RuneCharNode{value: `\s`}
	// Regex equivalent: \S
	NonWhitespace = RuneCharNode{value: `\S`}
	// Regex equivalent: \t
	Tab = RuneCharNode{value: `\t`}
	// Regex equivalent: \n
	Newline = RuneCharNode{value: `\n`}
	// Regex equivalent: (\r\n|\n|\r)
	LineBreak = RuneCharNode{value: `(\r\n|\n|\r)`}
	// Regex equivalent: ' '
	Space = RuneCharNode{value: ` `}
	// Regex equivalent: \p{Z} => includes: all Unicode space separators (space, non-breaking space, thin space, ideographic space)
	SpaceSeparator = RuneCharNode{value: `\p{Z}`}

	// Regex equivalent: [a-z] => includes: all lowercase Latin letters
	LowerLatin = RuneCharNode{value: `[a-z]`}
	// Regex equivalent: [A-Z] => includes: all uppercase Latin letters
	UpperLatin = RuneCharNode{value: `[A-Z]`}
	// Regex equivalent: [a-zA-Z] => includes: all Latin letters
	Latin = RuneCharNode{value: `[a-zA-Z]`}
	// Regex equivalent: [a-zA-Z0-9] => includes: all Latin letters or digits
	LatinDigit = RuneCharNode{value: `[a-zA-Z0-9]`}
	// Regex equivalent: \p{Latin} => includes: all Latin script letters, including extended letters like ä, ö, ü, ß, ñ, ç, and accented characters
	ExtendedLatin = RuneCharNode{value: `\p{Latin}`}
	// Regex equivalent: \p{L} => includes: all Unicode letters (Latin, Cyrillic, Greek, Arabic, Chinese, etc.)
	Letter = RuneCharNode{value: `\p{L}`}
	// Regex equivalent: \p{Lu} => includes: all uppercase Unicode letters (A, Ä, Б, Ω, Ї, etc.)
	UpperLetter = RuneCharNode{value: `\p{Lu}`}
	// Regex equivalent: \p{Ll} => includes: all lowercase Unicode letters (a, ä, б, λ, ї, etc.)
	LowerLetter = RuneCharNode{value: `\p{Ll}`}
	// Regex equivalent: \p{Cyrillic} => includes: all Cyrillic letters (А, Б, В, а, б, в, ї, є, ё, etc.)
	Cyrillic = RuneCharNode{value: `\p{Cyrillic}`}
	// Regex equivalent: \p{Greek} => includes: all Greek letters (Α, Β, Γ, α, β, γ, etc.)
	Greek = RuneCharNode{value: `\p{Greek}`}
	// Regex equivalent: \p{Arabic} => includes: all Arabic letters, numerals, and common script characters
	Arabic = RuneCharNode{value: `\p{Arabic}`}
	// Regex equivalent: \p{Hebrew} => includes: all Hebrew letters (א, ב, ג, etc.)
	Hebrew = RuneCharNode{value: `\p{Hebrew}`}
	// Regex equivalent: \p{Han} => includes: all CJK ideographs (Chinese, Japanese, Korean)
	Han = RuneCharNode{value: `\p{Han}`}
	// Regex equivalent: [\p{L}\p{N}\p{P}\p{S}] => includes: all printable characters
	Printable = RuneCharNode{value: `[\p{L}\p{N}\p{P}\p{S}]`}

	// Regex equivalent: \p{Nd} => includes: all Unicode decimal digits (0–9 and other script digits)
	DigitUnicode = RuneCharNode{value: `\p{Nd}`}
	// Regex equivalent: \p{N} => includes: all Unicode numeric characters (digits, Roman numerals, fractions, superscripts)
	Numeric = RuneCharNode{value: `\p{N}`}
	// Regex equivalent: [0-9A-Fa-f] => includes: all hexadecimal digits
	HexDigit = RuneCharNode{value: `[0-9A-Fa-f]`}
	// Regex equivalent: \d => includes: all ASCII digits
	Digit = RuneCharNode{value: `\d`}
	// Regex equivalent: \D => excludes: all characters that are ASCII digits
	NonDigit = RuneCharNode{value: `\D`}

	// Regex equivalent: \p{P} => includes: all Unicode punctuation (.,!?'"()[]{}-—«»… etc.)
	Punctuation = RuneCharNode{value: `\p{P}`}
	// Regex equivalent: \p{S} => includes: all Unicode symbols ($€£¥+=<>©®™✓∞, emojis, and similar)
	Symbol = RuneCharNode{value: `\p{S}`}

	// Regex equivalent: \w => includes: [a-zA-Z0-9_]
	WordChar = RuneCharNode{value: `\w`}
	// Regex equivalent: \W => excludes: [a-zA-Z0-9_]
	NonWordChar = RuneCharNode{value: `\W`}

	// Regex equivalent: .
	AnyChar = MetaCharNode{value: `.`}
	// Regex equivalent: $
	LineEnd = MetaCharNode{value: `$`}
	// Regex equivalent: ^
	LineStart = MetaCharNode{value: `^`}

	// Regex equivalent: \r
	Return = RuneCharNode{value: `\r`}
	// Regex equivalent: \b
	WordBoundary = RuneCharNode{value: `\b`}
	// Regex equivalent: \B
	NonWordBoundary = RuneCharNode{value: `\B`}
)
