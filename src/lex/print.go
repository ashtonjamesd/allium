package lex

import (
	"fmt"
	"strings"
)

func (t TokenType) String() string {
	switch t {
	case Hashtag:
		return "Hashtag"
	case Comma:
		return "Comma"
	case Identifier:
		return "Identifier"
	case NewLine:
		return "NewLine"
	case WhiteSpace:
		return "WhiteSpace"
	case Exclamation:
		return "Exclamation"
	case Star:
		return "Star"
	case Underscore:
		return "Underscore"
	case CarriageReturn:
		return "CarriageReturn"
	case BackTick:
		return "BackTick"
	case Eof:
		return "Eof"
	case Minus:
		return "Minus"
	case GreaterThan:
		return "GreaterThan"
	case LeftSquareBracket:
		return "LeftSquareBracket"
	case RightSquareBracket:
		return "RightSquareBracket"
	case LeftParen:
		return "LeftParen"
	case RightParen:
		return "RightParen"
	case None:
		return "None"
	case Dot:
		return "Dot"
	case Number:
		return "Number"
	default:
		return "Unknown"
	}
}

func PrintTokens(tokens []Token) {
	for i, token := range tokens {
		escapedValue := escapeSpecialChars(string(token.Value))
		fmt.Printf("%d: %s %s\n", i, escapedValue, token.TokenKind)
	}
}

func escapeSpecialChars(s string) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	s = strings.ReplaceAll(s, "\b", "\\b")
	s = strings.ReplaceAll(s, "\f", "\\f")
	return s
}
