package parser

import "unicode"

const (
	OpPlus     = '+'
	OpMinus    = '-'
	OpDivide   = '/'
	OpMultiply = '*'
	OpEqual    = '='

	OpenParen  = '('
	CloseParen = ')'

	Space      = ' '
	Dot        = '.'
	Underscore = '_'
)

func containsDot(number []rune) bool {
	for _, char := range number {
		if char == Dot {
			return true
		}
	}
	return false
}

func isLetter(char rune) bool {
	return unicode.IsLetter(char) || char == Underscore
}

const (
	KindOpPlus     = "KindOpPlus"
	KindOpMinus    = "KindOpMinus"
	KindOpDivide   = "KindOpDivide"
	KindOpMultiply = "KindOpMultiply"
	KindOpEqual    = "KindOpEqual"

	KindParentheses = "KindParentheses"

	KindInteger = "KindInteger"
	KindFloat   = "KindFloat"
	KindString  = "KindString"

	KindVar = "KindVar"
)
