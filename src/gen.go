package main

import (
	"fmt"
	"os"
)

type Generator struct {
	nodes        []NodeInterface
	headerCount  int
	previousNode NodeInterface
}

func NewGenerator(nodes []NodeInterface) Generator {
	var gen = Generator{}
	gen.nodes = nodes

	return gen
}

func (g *Generator) GenerateHtml() {
	_ = os.Mkdir("out", 0755)
	file, err := os.Create("out/out.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, node := range g.nodes {
		g.convert_node(file, node)
	}
}

func (g *Generator) convert_node(file *os.File, node NodeInterface) {
	switch node := node.(type) {
	case HeaderNode:
		g.headerCount++

		fmt.Fprintf(file, "<h%d id =\"%s-%d\">", node.level, "header", g.headerCount)
		for _, headerNode := range node.content {
			g.convert_node(file, headerNode)
		}
		fmt.Fprintf(file, "</h%d>\n", node.level)
	case ParagraphNode:
		fmt.Fprintf(file, "<p>")
		for _, paragraphNode := range node.content {
			g.convert_node(file, paragraphNode)
		}
		fmt.Fprintf(file, "</p>\n")
	case ItalicNode:
		fmt.Fprintf(file, "<em>%s</em>", node.content)
	case BoldNode:
		fmt.Fprintf(file, "<strong>%s</strong>", node.content)
	case WhiteSpaceNode:
		fmt.Fprintf(file, " ")
	case TextNode:
		fmt.Fprintf(file, node.content)
	case NoNode:
		fmt.Fprintf(file, "")
	case NewLineNode:
		// if _, wasNewLine := g.previousNode.(NewLineNode); !wasNewLine {
		// fmt.Fprintf(file, "<br>")
		fmt.Fprintf(file, "\n")
		// }
	default:
		fmt.Printf("Unknown node type: %s", node)
	}

	g.previousNode = node
}
