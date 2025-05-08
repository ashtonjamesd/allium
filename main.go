package main

import (
	"fmt"
	"os"
)

func main() {
	dat, err := os.ReadFile("example.md")
	if err != nil {
		fmt.Println(err)
		return
	}

	lex := NewLexer(string(dat))
	tokens := Tokenize(lex)

	for i, token := range tokens {
		fmt.Printf("%d: %s %d\n", i, token.value, int(token.tokenKind))
	}

}
