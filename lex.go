package main

import "fmt"

type LexState struct {
	source  string
	current int
}

func NewLexer(source string) *LexState {
	var lex LexState
	lex.source = source

	return &lex
}

func isEnd(state *LexState) bool {
	return state.current >= len(state.source)
}

func current(state *LexState) byte {
	if isEnd(state) {
		return 0
	}
	return state.source[state.current]
}

func advance(state *LexState) {
	state.current++
}

func Tokenize(state *LexState) {

	for !isEnd(state) {
		ch := current(state)
		fmt.Printf("%c\n", ch)

		advance(state)
	}
}

func parse_char() {

}
