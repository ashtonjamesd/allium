package main

import (
	"fmt"
	"os"
)

func GenerateHtml(nodes []NodeInterface) {
	_ = os.Mkdir("out", 0755)
	file, err := os.Create("out/out.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, node := range nodes {
		convert_node(file, node)
	}
}

func convert_node(file *os.File, node NodeInterface) {
	switch node := node.(type) {
	case HeaderNode:
		fmt.Fprintf(file, "<h%d>%s</h%d>\n", node.level, node.content, node.level)
	case GenericNode:
		fmt.Fprintf(file, "<p>%s</p>\n", node.content)
	default:
		fmt.Println("Unknown node type")
	}
}
