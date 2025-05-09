package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	toHTML = "tohtml"
	toMD   = "tomd"
)

func main() {
	convertFlag := flag.String("convert", "", "Conversion type: tohtml or tomd")
	flag.Parse()

	if *convertFlag == "" {
		fmt.Println("Usage: go run ./src --convert=[tohtml | tomd]")
		return
	}

	switch *convertFlag {
	case toHTML:
		err := convertToHTML("example/example.md")
		if err != nil {
			fmt.Printf("Error converting to HTML: %v\n", err)
		}
	case toMD:
		fmt.Println("tomd not currently supported.")
	default:
		fmt.Printf("Invalid convert type: %s\n", *convertFlag)
	}
}

func convertToHTML(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	lexer := NewLexer(string(data))
	tokens := lexer.Tokenize()
	PrintTokens(tokens)

	parser := NewParser(tokens)
	exprs := parser.Parse()
	PrintNodes(exprs)

	gen := NewGenerator(exprs)
	gen.GenerateHtml()

	return nil
}
