package lex

type TokenType int

const (
	Hashtag TokenType = iota
	Comma
	Identifier
	Number
	NewLine
	CarriageReturn
	WhiteSpace
	Exclamation
	Star
	Underscore
	LeftSquareBracket
	RightSquareBracket
	LeftParen
	RightParen
	GreaterThan
	BackTick
	Dot
	Tab
	Minus
	Eof
	None
)

type Token struct {
	TokenKind TokenType
	Value     string
}
