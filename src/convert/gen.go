package gen

import (
	"allium/src/parse"
	"fmt"
	"os"
)

type Generator struct {
	Nodes        []parse.NodeInterface
	HeaderCount  int
	PreviousNode parse.NodeInterface
}

func NewGenerator(nodes []parse.NodeInterface) Generator {
	var gen = Generator{}
	gen.Nodes = nodes

	return gen
}

func (g *Generator) GenerateHtml() {
	_ = os.Mkdir("out", 0755)
	file, err := os.Create("out/out.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, node := range g.Nodes {
		g.convert_node(file, node)
	}
}

func (g *Generator) convert_node(file *os.File, node parse.NodeInterface) {
	switch node := node.(type) {
	case parse.HeaderNode:
		g.HeaderCount++

		fmt.Fprintf(file, "<h%d id =\"%s-%d\">", node.Level, "header", g.HeaderCount)
		for _, headerNode := range node.Content {
			g.convert_node(file, headerNode)
		}
		fmt.Fprintf(file, "</h%d>\n", node.Level)
	case parse.ParagraphNode:
		fmt.Fprintf(file, "<p>")
		for _, paragraphNode := range node.Content {
			g.convert_node(file, paragraphNode)
		}
		fmt.Fprintf(file, "</p>\n")
	case parse.ItalicNode:
		fmt.Fprintf(file, "<em>%s</em>", node.Content)
	case parse.BoldNode:
		fmt.Fprintf(file, "<strong>%s</strong>", node.Content)
	case parse.WhiteSpaceNode:
		fmt.Fprintf(file, " ")
	case parse.TextNode:
		fmt.Fprintf(file, node.Content)
	case parse.NoNode:
		fmt.Fprintf(file, "")
	case parse.NewLineNode:
		// if _, wasNewLine := g.previousNode.(NewLineNode); !wasNewLine {
		// fmt.Fprintf(file, "<br>")
		fmt.Fprintf(file, "\n")
		// }
	default:
		fmt.Printf("Unknown node type: %s", node)
	}

	g.PreviousNode = node
}
