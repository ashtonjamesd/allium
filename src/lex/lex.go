package lex

import (
	"unicode"
)

type LexState struct {
	source  string
	current int
	tokens  []Token
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
	case unicode.IsDigit(l.currentChar()):
		return l.parseNumeric()
	default:
		return l.parseSymbol()
	}
}

func (l *LexState) parseSymbol() Token {
	symbolMap := make(map[string]TokenType)
	symbolMap["#"] = Hashtag
	symbolMap[","] = Comma
	symbolMap["\n"] = NewLine
	symbolMap["\r"] = CarriageReturn
	symbolMap[" "] = WhiteSpace
	symbolMap["!"] = Exclamation
	symbolMap["*"] = Star
	symbolMap["_"] = Underscore
	symbolMap["["] = LeftSquareBracket
	symbolMap["]"] = RightSquareBracket
	symbolMap["("] = LeftParen
	symbolMap[")"] = RightParen
	symbolMap["\t"] = Tab
	symbolMap["-"] = Minus
	symbolMap[">"] = GreaterThan
	symbolMap["`"] = BackTick
	symbolMap["."] = Dot

	c := string(l.currentChar())

	value, ok := symbolMap[c]
	if !ok {
		return newToken(None, c)
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

func (l *LexState) parseNumeric() Token {
	var start = l.current
	for unicode.IsDigit(l.currentChar()) {
		l.advance()
	}

	var lexeme = l.source[start:l.current]
	l.recede()

	return newToken(Number, lexeme)
}

func newToken(lexType TokenType, value string) Token {
	var token Token
	token.Value = value
	token.TokenKind = lexType

	return token
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

func NewLexer(source string) *LexState {
	var lexer LexState
	lexer.source = source

	return &lexer
}
