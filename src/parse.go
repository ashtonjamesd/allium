package main

import (
	"fmt"
	"math"
)

type NodeType int

const (
	GenericNodeType NodeType = iota
	HeaderNodeType
)

type NodeInterface any

type HeaderNode struct {
	level   int
	content []NodeInterface
}

type ParagraphNode struct {
	content []NodeInterface
}

type NoNode struct{}

type ItalicNode struct {
	content string
}

type TextNode struct {
	content string
}

type BoldNode struct {
	content string
}

type WhiteSpaceNode struct{}
type NewLineNode struct{}

type Parser struct {
	nodes   []NodeInterface
	tokens  []Token
	current int
}

func (p *Parser) currentToken() Token {
	if p.isEnd() {
		return Token{}
	}
	return p.tokens[p.current]
}

func (p *Parser) match(tokenType TokenType) bool {
	return p.currentToken().tokenKind == tokenType
}

func (p *Parser) advance() {
	p.current++
}

func (p *Parser) recede() {
	p.current--
}

func (p *Parser) isEnd() bool {
	return p.current >= len(p.tokens)
}

// func (p *Parser) peek() Token {
// 	p.advance()
// 	next := p.currentToken()
// 	p.recede()

// 	return next
// }

func (p *Parser) Parse() []NodeInterface {
	for !p.isEnd() && p.currentToken().tokenKind != Eof {
		node := p.parseBlock()
		p.nodes = append(p.nodes, node)

		p.advance()
	}

	return p.nodes
}

func (p *Parser) parseBlock() NodeInterface {
	switch p.currentToken().tokenKind {
	case Hashtag:
		return p.parseHeader()
	default:
		return p.parseParagraph()
	}
}

func (p *Parser) parseParagraph() NodeInterface {
	if p.match(NewLine) {
		return NoNode{}
	}

	var node ParagraphNode
	var content []NodeInterface

	var lastNode NodeInterface
	for !p.isEnd() && p.currentToken().tokenKind != NewLine {
		currentNode := p.parseInline()

		if _, ok := lastNode.(NewLineNode); ok {
			if _, ok := currentNode.(NewLineNode); ok {
				p.recede()
				break
			}
		}

		content = append(content, currentNode)
		lastNode = currentNode
	}

	for len(content) > 0 {
		if _, ok := content[len(content)-1].(NewLineNode); ok {
			content = content[:len(content)-1]
		} else {
			break
		}
	}

	filtered := make([]NodeInterface, 0, len(content))
	for _, n := range content {
		switch n.(type) {
		case NewLineNode, NoNode:

		default:
			filtered = append(filtered, n)
		}
	}

	if len(filtered) == 0 {
		return NoNode{}
	}

	node.content = content
	return node
}

func (p *Parser) parseInline() NodeInterface {
	switch p.currentToken().tokenKind {
	case Star, Underscore:
		return p.parseBoldItalic()
	case WhiteSpace:
		p.advance()
		return WhiteSpaceNode{}
	case NewLine, CarriageReturn:
		p.advance()
		p.advance()
		return NewLineNode{}
	case Identifier, Comma, Exclamation:
		content := p.currentToken().value
		p.advance()
		return TextNode{content: content}
	default:
		p.advance()
		return NoNode{}
	}
}

func (p *Parser) isBoldItalicToken(token Token) bool {
	return token.tokenKind == Star || token.tokenKind == Underscore
}

func (p *Parser) parseBoldItalic() NodeInterface {
	p.advance()
	isBold := p.isBoldItalicToken(p.currentToken())
	if isBold {
		p.advance()
	}

	content := ""
	for !p.isEnd() && p.currentToken().tokenKind != Eof {
		content += p.currentToken().value
		p.advance()

		if p.isBoldItalicToken(p.currentToken()) {
			break
		}
	}

	p.advance()

	if isBold {
		p.advance()

		var node BoldNode
		node.content = content
		return node
	}

	var node ItalicNode
	node.content = content
	return node
}

func (p *Parser) parseHeader() NodeInterface {
	level := 1

	p.advance()
	for p.currentToken().tokenKind == Hashtag {
		level++
		p.advance()
	}
	if p.currentToken().tokenKind == WhiteSpace {
		p.advance()
	}

	var children []NodeInterface
	for !p.isEnd() && p.currentToken().tokenKind != NewLine {
		expr := p.parseInline()
		children = append(children, expr)

		if p.match(NewLine) {
			break
		}
	}

	var node HeaderNode
	node.level = int(math.Min(float64(level), 6))
	node.content = children

	return node
}

func NewParser(tokens []Token) *Parser {
	var parser Parser
	parser.tokens = tokens

	return &parser
}

func PrintNodes(nodes []NodeInterface) {
	for _, node := range nodes {
		printNode(node, 0)
	}
}

func spaces(n int) string {
	return fmt.Sprintf("%*s", n, "")
}

func (n ParagraphNode) Print(indent int) {
	fmt.Printf("%sParagraphNode: %s\n", spaces(indent), n.content)
}

func (n ItalicNode) Print(indent int) {
	fmt.Printf("%sItalicNode: _%s_\n", spaces(indent), n.content)
}

func (n BoldNode) Print(indent int) {
	fmt.Printf("%sBoldNode: **%s**\n", spaces(indent), n.content)
}

func (n TextNode) Print(indent int) {
	fmt.Printf("%sTextNode: '%s'\n", spaces(indent), n.content)
}

func (n HeaderNode) Print(indent int) {
	fmt.Printf("%sHeaderNode (level %d):\n", spaces(indent), n.level)
	for _, child := range n.content {
		printNode(child, indent+2)
	}
}

func (n WhiteSpaceNode) Print(indent int) {
	fmt.Printf("%sWhiteSpaceNode: '%s'\n", spaces(indent), " ")
}

func (n NewLineNode) Print(indent int) {
	fmt.Printf("%sNewLineNode: '%s'\n", spaces(indent), "\\n")
}

func (n NoNode) Print(indent int) {
	fmt.Printf("%sNoNode: '%s'\n", spaces(indent), "none")
}

func printNode(n NodeInterface, indent int) {
	switch node := n.(type) {
	case ParagraphNode:
		node.Print(indent)
	case TextNode:
		node.Print(indent)
	case ItalicNode:
		node.Print(indent)
	case BoldNode:
		node.Print(indent)
	case HeaderNode:
		node.Print(indent)
	case WhiteSpaceNode:
		node.Print(indent)
	case NewLineNode:
		node.Print(indent)
	case NoNode:
		node.Print(indent)
	default:
		fmt.Printf("%sUnknown node type\n", spaces(indent))
	}
}
