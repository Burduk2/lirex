package lirex

type HelperNode struct {
	name       string
	groupNames []string
	value      string
}
type helpersMap struct {
	Domain             HelperNode
	Email              HelperNode
	InternationalPhone HelperNode
	CreditCard         HelperNode
	FullUrl            HelperNode
}

func compile(node Node, groups ...string) HelperNode {
	return HelperNode{
		value:      Exp(node).MustCompile(Options{}).String(),
		name:       groups[0],
		groupNames: groups,
	}
}

var Helpers = helpersMap{
	Domain:             compile(domain(true), "Domain"),
	Email:              compile(email(), "Email", "Email_localPart", "Email_domain"),
	InternationalPhone: compile(phone(), "Phone", "Phone_countryCode", "Phone_areaCode"),
	CreditCard:         compile(creditCard(), "CreditCard"),
	FullUrl:            compile(url(), "FullUrl"),
}

func domain(capture bool) Node {
	exp := Seq(
		LatinDigit,
		Group(
			CharClass(LatinDigit, Lit("-.")).ZeroOrMore(),
			LatinDigit,
		).Optional(),
		Lit("."),
		LatinDigit.Between(1, 63),
	)
	if capture {
		return Capture("Domain", exp)
	}
	return Group(exp)
}
func email() CaptureNode {
	return Capture("Email",
		Capture("Email_localPart",
			Group(
				CharClass(WordChar, Lit("%+")),
				Group(
					CharClass(WordChar, Lit(".%+-")).ZeroOrMore(),
					CharClass(WordChar, Lit("%+-")),
				).Optional(),
			).Between(1, 64),
		),
		Lit("@"),
		Capture("Email_domain", domain(false)),
	)
}
func phone() CaptureNode {
	return Capture("Phone",
		Lit("+"),
		Capture("Phone_countryCode", Digit.Between(1, 3)),
		Whitespace.Optional(),
		Capture("Phone_areaCode",
			Or(
				Digit.Exactly(3),
				Group(Lit("("), Digit.Exactly(3), Lit(")")),
			),
		),
		CharClass(Digit, Lit(" .-")).Between(0, 10+3),
		Digit,
	)
}
func creditCard() CaptureNode {
	return Capture("CreditCard",
		Group(
			Digit.Exactly(4),
			CharClass(Lit(" -")).Optional(),
		).Exactly(4),
	)
}
func url() CaptureNode {
	return Capture("FullUrl",
		Latin,
		CharClass(LatinDigit, Lit("+-.")).ZeroOrMore(),
		Lit("://"),
		domain(false),
		Group(Lit(":"), Digit.AtLeast(1)).Optional(),
		Group(Lit("/"), CharClass(WordChar, Lit("%/-._~")).ZeroOrMore()).Optional(),
		Group(Lit("?"), CharClass(WordChar, Lit("=&%+-._~")).ZeroOrMore()).Optional(),
		Group(Lit("#"), CharClass(WordChar, Lit("%-._~")).ZeroOrMore()).Optional(),
	)
}

var ReservedGroupNames = map[string]struct{}{
	"Domain":            {},
	"Email":             {},
	"Email_localPart":   {},
	"Email_domain":      {},
	"Phone":             {},
	"Phone_countryCode": {},
	"Phone_areaCode":    {},
	"CreditCard":        {},
	"FullUrl":           {},
}
