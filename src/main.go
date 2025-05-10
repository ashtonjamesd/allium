package main

import (
	gen "allium/src/convert"
	"allium/src/lex"
	"allium/src/parse"
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
	pathFlag := flag.String("path", "", "Source path")
	outputFlag := flag.String("output", "", "output path")
	flag.Parse()

	if *convertFlag == "" || *pathFlag == "" || *outputFlag == "" {
		fmt.Println("Usage: go run ./src --convert=[tohtml | tomd] --path=<source> --output=*.[html | md]")
		return
	}

	switch *convertFlag {
	case toHTML:
		err := convertToHTML(*pathFlag, *outputFlag)
		if err != nil {
			fmt.Printf("Error converting to HTML: %v\n", err)
		}
	case toMD:
		fmt.Println("tomd not currently supported.")
	default:
		fmt.Printf("Invalid convert type: %s\n", *convertFlag)
	}
}

func convertToHTML(filepath string, outputPath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	lexer := lex.NewLexer(string(data))
	tokens := lexer.Tokenize()
	// lex.PrintTokens(tokens)

	parser := parse.NewParser(tokens)
	exprs := parser.Parse()
	// parse.PrintNodes(exprs)

	gen := gen.NewGenerator(exprs)
	gen.GenerateHtml(outputPath)

	fmt.Printf("Finished converting Markdown to HTML")
	fmt.Printf("Output at %s", outputPath)

	return nil
}
