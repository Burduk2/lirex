package lirex

type helpersMap struct {
	Domain             CaptureNode
	Email              CaptureNode
	InternationalPhone CaptureNode
	CreditCard         CaptureNode
	FullUrl            CaptureNode
}

var Helpers = helpersMap{
	Domain:             domain(),
	Email:              email(),
	InternationalPhone: phone(),
	CreditCard:         creditCard(),
	FullUrl:            url(),
}

func domain() CaptureNode {
	return Capture("Domain",
		Or(
			LatinDigit,
			Group(
				LatinDigit,
				CharClass(LatinDigit, Lit("-.")).ZeroOrMore(),
				LatinDigit,
			),
		),
		Lit("."),
		LatinDigit.Between(1, 63),
	)
}
func email() CaptureNode {
	return Capture("Email",
		Capture("Email_localPart",
			Or(
				CharClass(WordChar, Lit("%+")),
				Group(
					CharClass(WordChar, Lit("%+")),
					CharClass(WordChar, Lit(".%+-")).ZeroOrMore(),
					CharClass(WordChar, Lit("%+-")),
				),
			).Between(1, 64),
		),
		Lit("@"),
		Capture("Email_domain", domain()),
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
		Capture("FullUrl_protocol", Or(Lit("http"), Lit("https"), Lit("ftp"))),
		Lit("://"),
		Capture("FullUrl_domain", domain()),
		Group(Lit(":"), Capture("Url_port", Digit.AtLeast(1))).Optional(),
		Capture("FullUrl_path", Group(Lit("/"), CharClass(WordChar, Lit("%/-._~")).ZeroOrMore()).Optional()),
		Capture("FullUrl_query", Group(Lit("?"), CharClass(WordChar, Lit("=&%+-._~")).ZeroOrMore()).Optional()),
		Capture("FullUrl_fragment", Group(Lit("#"), CharClass(WordChar, Lit("%-._~")).ZeroOrMore()).Optional()),
	)
}
