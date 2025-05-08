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
	Tokenize(lex)

}
