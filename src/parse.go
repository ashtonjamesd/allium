package main

import "fmt"

type NodeType int

const (
	GenericNodeType NodeType = iota
	HeaderNodeType
)

type NodeInterface interface {
	Type() NodeType
	String() string
}

type HeaderNode struct {
	level   int
	content string
}

func (h HeaderNode) Type() NodeType {
	return HeaderNodeType
}

func (h HeaderNode) String() string {
	return h.content
}

type GenericNode struct {
	content string
}

func (g GenericNode) Type() NodeType {
	return GenericNodeType
}

func (g GenericNode) String() string {
	return g.content
}

type Parser struct {
	nodes   []NodeInterface
	tokens  []Token
	current int
}

func (p *Parser) currentToken() Token {
	return p.tokens[p.current]
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

func (p *Parser) Parse() []NodeInterface {
	for i, token := range p.tokens {
		fmt.Printf("%d: %s %d\n", i, token.value, int(token.tokenKind))
	}

	for !p.isEnd() && p.currentToken().tokenKind != Eof {
		node := p.parseToken()
		p.nodes = append(p.nodes, node)

		p.advance()
	}

	return p.nodes
}

func (p *Parser) parseToken() NodeInterface {
	switch p.currentToken().tokenKind {
	case Hashtag:
		return p.parseHeader()
	default:
		return p.parseText()
	}
}

func (p *Parser) parseText() NodeInterface {
	var node GenericNode
	node.content = p.currentToken().value
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

	content := ""
	for !p.isEnd() && p.currentToken().tokenKind != NewLine {
		content += p.currentToken().value
		p.advance()
	}

	var node HeaderNode
	node.level = level
	node.content = content

	return node
}

func NewParser(tokens []Token) *Parser {
	var parser Parser
	parser.tokens = tokens

	return &parser
}
