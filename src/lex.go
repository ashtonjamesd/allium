package main

import (
	"unicode"
)

type TokenType int

const (
	Hashtag TokenType = iota
	Comma
	Identifier
	NewLine
	WhiteSpace
	Exclamation
	Eof
	None
)

type LexState struct {
	source  string
	current int
	tokens  []Token
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

func (l *LexState) isEnd() bool {
	return l.current >= len(l.source)
}

func (l *LexState) currentChar() rune {
	if l.isEnd() {
		return 0
	}
	return rune(l.source[l.current])
}

func (l *LexState) advance() {
	l.current++
}

func (l *LexState) recede() {
	l.current--
}

func (l *LexState) Tokenize() []Token {
	for !l.isEnd() {
		expr := l.parseChars()
		l.tokens = append(l.tokens, expr)

		l.advance()
	}

	l.tokens = append(l.tokens, newToken(Eof, ""))

	return l.tokens
}

func (l *LexState) parseChars() Token {
	switch {
	case unicode.IsLetter(l.currentChar()):
		return l.parseIdentifier()
	default:
		return l.parseSymbol()
	}
}

func (l *LexState) parseSymbol() Token {
	symbolMap := make(map[string]TokenType)
	symbolMap["#"] = Hashtag
	symbolMap[","] = Comma
	symbolMap["\n"] = NewLine
	symbolMap[" "] = WhiteSpace
	symbolMap["!"] = Exclamation

	c := string(l.currentChar())

	value, ok := symbolMap[c]
	if !ok {
		return newToken(None, "")
	}

	return newToken(value, c)
}

func (l *LexState) parseIdentifier() Token {
	var start = l.current
	for unicode.IsLetter(l.currentChar()) {
		l.advance()
	}

	var lexeme = l.source[start:l.current]
	l.recede()

	return newToken(Identifier, lexeme)
}

func newToken(lexType TokenType, value string) Token {
	var token Token
	token.value = value
	token.tokenKind = lexType

	return token
}
