package main

import "unicode"

type TokenType int

const (
	Hashtag TokenType = iota
	Identifier
	None
)

type LexState struct {
	source  string
	current int
	exprs   []Token
}

type Token struct {
	tokenKind TokenType
	value     string
}

func NewLexer(source string) *LexState {
	var lexer LexState
	lexer.source = source

	return &lexer
}

func isEnd(lexer *LexState) bool {
	return lexer.current >= len(lexer.source)
}

func current(lexer *LexState) byte {
	if isEnd(lexer) {
		return 0
	}
	return lexer.source[lexer.current]
}

func match(c byte, lexer *LexState) bool {
	return current(lexer) == c
}

func advance(lexer *LexState) {
	lexer.current++
}

func Tokenize(lexer *LexState) []Token {
	var exprs []Token

	for !isEnd(lexer) {
		expr := parse_chars(lexer)
		exprs = append(exprs, expr)

		advance(lexer)
	}

	return exprs
}

func parse_chars(lexer *LexState) Token {
	if match('#', lexer) {
		return newToken(Hashtag, "#")
	} else if unicode.IsLetter(rune(current(lexer))) {
		return newToken(Identifier, string(current(lexer)))
	}

	var none Token
	none.tokenKind = None

	return none
}

func newToken(lexType TokenType, value string) Token {
	var token Token
	token.value = value
	token.tokenKind = lexType

	return token
}
