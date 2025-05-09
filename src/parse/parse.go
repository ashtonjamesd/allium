package parse

import (
	"allium/src/lex"
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
	Level   int
	Content []NodeInterface
}

type ParagraphNode struct {
	Content []NodeInterface
}

type NoNode struct{}

type ItalicNode struct {
	Content string
}

type TextNode struct {
	Content string
}

type BoldNode struct {
	Content string
}

type LinkNode struct {
	LinkText string
	Link     string
}

type UnorderedList struct {
	Nodes []NodeInterface
}

type UnorderedListNode struct {
	Nodes []NodeInterface
}

type WhiteSpaceNode struct{}
type NewLineNode struct{}

type Parser struct {
	Nodes   []NodeInterface
	Tokens  []lex.Token
	Current int
}

func (p *Parser) currentToken() lex.Token {
	if p.isEnd() {
		return lex.Token{}
	}
	return p.Tokens[p.Current]
}

func (p *Parser) match(tokenType lex.TokenType) bool {
	return p.currentToken().TokenKind == tokenType
}

func (p *Parser) advance() {
	p.Current++
}

func (p *Parser) recede() {
	p.Current--
}

func (p *Parser) isEnd() bool {
	return p.Current >= len(p.Tokens)
}

// func (p *Parser) peek() Token {
// 	p.advance()
// 	next := p.currentToken()
// 	p.recede()

// 	return next
// }

func (p *Parser) Parse() []NodeInterface {
	for !p.isEnd() && p.currentToken().TokenKind != lex.Eof {
		node := p.parseBlock()
		p.Nodes = append(p.Nodes, node)
	}

	return p.Nodes
}

func (p *Parser) parseBlock() NodeInterface {
	switch p.currentToken().TokenKind {
	case lex.Hashtag:
		return p.parseHeader()
	case lex.LeftSquareBracket:
		return p.parseLink()
	case lex.Star:
		return p.parseUnorderedList()
	default:
		return p.parseParagraph()
	}
}

func (p *Parser) parseUnorderedList() NodeInterface {
	var list UnorderedList

	for !p.isEnd() && p.currentToken().TokenKind == lex.Star {
		listItem := p.parseUnorderedListNode()
		list.Nodes = append(list.Nodes, listItem)

		for !p.isEnd() && (p.currentToken().TokenKind == lex.NewLine || p.currentToken().TokenKind == lex.CarriageReturn) {
			p.advance()
		}
	}

	return list
}

func (p *Parser) parseUnorderedListNode() NodeInterface {
	var node UnorderedListNode

	if p.match(lex.Star) {
		p.advance()
	}

	if p.match(lex.WhiteSpace) {
		p.advance()
	}

	for !p.isEnd() &&
		p.currentToken().TokenKind != lex.NewLine &&
		p.currentToken().TokenKind != lex.CarriageReturn {

		token := p.parseInline()
		if _, ok := token.(NoNode); !ok {
			node.Nodes = append(node.Nodes, token)
		}
	}

	return node
}

func (p *Parser) parseLink() NodeInterface {
	p.advance()

	linkText := ""
	for !p.isEnd() && p.currentToken().TokenKind != lex.RightSquareBracket {
		linkText += p.currentToken().Value
		p.advance()
	}

	link := ""
	p.advance()
	p.advance()

	for p.currentToken().TokenKind != lex.RightParen {
		link += p.currentToken().Value
		p.advance()
	}

	var node LinkNode
	node.LinkText = linkText
	node.Link = link

	return node
}

func (p *Parser) parseParagraph() NodeInterface {
	if p.match(lex.CarriageReturn) {
		p.advance()
		p.advance()
		return NoNode{}
	}

	var node ParagraphNode
	var content []NodeInterface

	var lastNode NodeInterface
	for !p.isEnd() && p.currentToken().TokenKind != lex.NewLine {
		currentNode := p.parseInline()

		if _, ok := lastNode.(NewLineNode); ok {
			if _, ok := currentNode.(NewLineNode); ok {
				// p.recede()
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

	node.Content = content
	return node
}

func (p *Parser) parseInline() NodeInterface {
	switch p.currentToken().TokenKind {
	case lex.Star, lex.Underscore:
		return p.parseBoldItalic()
	case lex.WhiteSpace:
		p.advance()
		return WhiteSpaceNode{}
	case lex.NewLine, lex.CarriageReturn:
		p.advance()
		p.advance()
		return NewLineNode{}
	case lex.Identifier, lex.Comma, lex.Exclamation:
		content := p.currentToken().Value
		p.advance()
		return TextNode{Content: content}
	default:
		p.advance()
		return NoNode{}
	}
}

func (p *Parser) isBoldItalicToken(token lex.Token) bool {
	return token.TokenKind == lex.Star || token.TokenKind == lex.Underscore
}

func (p *Parser) parseBoldItalic() NodeInterface {
	p.advance()
	isBold := p.isBoldItalicToken(p.currentToken())

	if isBold {
		p.advance()
	}

	content := ""
	for !p.isEnd() && p.currentToken().TokenKind != lex.Eof {
		content += p.currentToken().Value
		p.advance()

		if p.isBoldItalicToken(p.currentToken()) {
			break
		}
	}

	p.advance()

	if isBold {
		p.advance()

		var node BoldNode
		node.Content = content
		return node
	}

	var node ItalicNode
	node.Content = content
	return node
}

func (p *Parser) parseHeader() NodeInterface {
	level := 1

	p.advance()
	for p.currentToken().TokenKind == lex.Hashtag {
		level++
		p.advance()
	}
	if p.currentToken().TokenKind == lex.WhiteSpace {
		p.advance()
	}

	var children []NodeInterface
	for !p.isEnd() && p.currentToken().TokenKind != lex.NewLine {
		expr := p.parseInline()

		if _, ok := expr.(NewLineNode); ok {
			break
		}

		if p.match(lex.NewLine) {
			break
		}

		children = append(children, expr)
	}

	var node HeaderNode
	node.Level = int(math.Min(float64(level), 6))
	node.Content = children

	return node
}

func NewParser(tokens []lex.Token) *Parser {
	var parser Parser
	parser.Tokens = tokens

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
	fmt.Printf("%sParagraphNode: %s\n", spaces(indent), n.Content)
}

func (n ItalicNode) Print(indent int) {
	fmt.Printf("%sItalicNode: _%s_\n", spaces(indent), n.Content)
}

func (n BoldNode) Print(indent int) {
	fmt.Printf("%sBoldNode: **%s**\n", spaces(indent), n.Content)
}

func (n TextNode) Print(indent int) {
	fmt.Printf("%sTextNode: '%s'\n", spaces(indent), n.Content)
}

func (n HeaderNode) Print(indent int) {
	fmt.Printf("%sHeaderNode (level %d):\n", spaces(indent), n.Level)
	for _, child := range n.Content {
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

func (n UnorderedListNode) Print(indent int) {
	fmt.Printf("%sUnorderedListNode: \n", spaces(indent))
	for _, child := range n.Nodes {
		printNode(child, indent+2)
	}
}

func (n UnorderedList) Print(indent int) {
	fmt.Printf("%sUnorderedList: \n", spaces(indent))
	for _, child := range n.Nodes {
		printNode(child, indent+2)
	}
}

func (n LinkNode) Print(indent int) {
	fmt.Printf("%sLinkNode: '%s' : '%s'\n", spaces(indent), n.LinkText, n.Link)
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
	case LinkNode:
		node.Print(indent)
	case UnorderedListNode:
		node.Print(indent)
	case UnorderedList:
		node.Print(indent)
	default:
		fmt.Printf("%sUnknown node type\n", spaces(indent))
	}
}
