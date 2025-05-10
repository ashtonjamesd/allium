package parse

import (
	"allium/src/lex"
	"math"
)

type Parser struct {
	Nodes   []NodeInterface
	Tokens  []lex.Token
	Current int
}

func (p *Parser) Parse() []NodeInterface {
	for !p.isEnd() && !p.match(lex.Eof) {
		node := p.parseBlock()
		p.Nodes = append(p.Nodes, node)
	}

	return p.Nodes
}

// func (p *Parser) parseOrderedList() NodeInterface {
// 	numberToken := p.currentToken()
// 	p.advance()
// 	p.advance()

// 	for

// 	var node OrderedListNode
// 	return node
// }

func (p *Parser) parseBlock() NodeInterface {

	switch p.currentToken().TokenKind {
	case lex.Number:
		if p.peek().TokenKind == lex.Dot {
			return p.parseUnorderedListGroup()
		}
		return p.parseParagraph()
	case lex.Hashtag:
		return p.parseHeader()
	case lex.LeftSquareBracket, lex.Exclamation:
		return p.parseLink()
	case lex.Star:
		if p.peek().TokenKind == lex.WhiteSpace {
			return p.parseUnorderedListGroup()
		}
		if p.peek().TokenKind == lex.Star {
			p.advance()

			if p.peek().TokenKind == lex.Star {
				p.recede()
				return p.parseHorizontalRule()
			} else {
				p.recede()
			}
		}
		return p.parseParagraph()
	case lex.Minus:
		if p.peek().TokenKind == lex.WhiteSpace {
			return p.parseUnorderedListGroup()
		}
		return p.parseHorizontalRule()
	case lex.GreaterThan:
		return p.parseBlockQuote()
	case lex.BackTick:
		if p.peek().TokenKind != lex.BackTick {
			return p.parseInlineCode()
		}
		return p.parseInlineCodeBlock()
	default:
		return p.parseParagraph()
	}
}

func (p *Parser) parseInlineCodeBlock() NodeInterface {
	p.advance()
	if p.match(lex.BackTick) {
		p.advance()

		if p.match(lex.BackTick) {
			p.advance()

			content := ""
			for !p.isEnd() && !p.match(lex.BackTick) {
				content += p.currentToken().Value
				p.advance()
			}

			p.advance()
			p.advance()
			p.advance()

			var node InlineCodeBlockNode
			node.Content = content

			return node
		} else {
			p.recede()
			p.recede()
		}
	}

	return p.parseParagraph()
}

func (p *Parser) parseInlineCode() NodeInterface {
	p.advance()

	content := ""
	for !p.isEnd() && !p.match(lex.BackTick) {
		content += p.currentToken().Value
		p.advance()
	}
	p.advance()

	var node InlineCodeNode
	node.Content = content
	return node
}

func (p *Parser) parseBlockQuote() NodeInterface {
	p.advance()

	var nodes []NodeInterface
	for !p.isEnd() && !p.isNodeEnd() {
		node := p.parseInline()
		nodes = append(nodes, node)
	}

	var node BlockQuoteNode
	node.Nodes = nodes

	return node
}

func (p *Parser) isHorizontalRuleToken() bool {
	return p.match(lex.Star) || p.match(lex.Minus)
}

func (p *Parser) parseHorizontalRule() NodeInterface {
	p.advance()

	if p.isHorizontalRuleToken() {
		p.advance()

		if p.isHorizontalRuleToken() {
			p.advance()
			return HorizontalRuleNode{}
		} else {
			p.recede()
			p.recede()
		}
	} else {
		p.recede()
	}

	return p.parseParagraph()
}

func (p *Parser) isListToken() bool {
	return p.match(lex.Star) || p.match(lex.Minus) || p.match(lex.Number)
}

func (p *Parser) parseUnorderedListGroup() NodeInterface {
	var list ListNode

	isOrdered := false
	if p.match(lex.Number) {
		isOrdered = true
	}

	for !p.isEnd() && p.isListToken() && (p.peek().TokenKind == lex.WhiteSpace || p.peek().TokenKind == lex.Dot) {
		listItem := p.parseUnorderedListNode()
		list.Nodes = append(list.Nodes, listItem)

		for !p.isEnd() && p.isNodeEnd() {
			p.advance()
		}
	}

	list.IsOrdered = isOrdered
	return list
}

func (p *Parser) parseUnorderedListNode() NodeInterface {
	var node ListItemNode

	if p.isListToken() {
		p.advance()
	}

	if p.match(lex.Dot) {
		p.advance()
	}

	if p.match(lex.WhiteSpace) {
		p.advance()
	}

	for !p.isEnd() && !p.isNodeEnd() {
		token := p.parseInline()
		if _, ok := token.(NoNode); !ok {
			node.Nodes = append(node.Nodes, token)
		}
	}

	return node
}

func (p *Parser) isNodeEnd() bool {
	return p.match(lex.NewLine) || p.match(lex.CarriageReturn)
}

func (p *Parser) parseLink() NodeInterface {
	isImage := p.match(lex.Exclamation)
	if isImage {
		p.advance()
	}

	p.advance()

	linkText := ""
	for !p.isEnd() && !p.match(lex.RightSquareBracket) {
		linkText += p.currentToken().Value
		p.advance()
	}

	link := ""
	p.advance()
	p.advance()

	for !p.match(lex.RightParen) {
		link += p.currentToken().Value
		p.advance()
	}
	p.advance()

	if isImage {
		var node ImageNode
		node.Link = link
		node.LinkText = linkText

		return node
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
	for !p.isEnd() && !p.match(lex.NewLine) {
		if p.match(lex.Hashtag) {
			break
		}

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
	default:
		content := p.currentToken().Value
		p.advance()
		return TextNode{Content: content}
	}
}

func (p *Parser) isBoldItalicToken(token lex.Token) bool {
	return token.TokenKind == lex.Star || token.TokenKind == lex.Underscore
}

func (p *Parser) parseBoldItalic() NodeInterface {
	marker := p.currentToken().TokenKind
	markerCount := 0

	for !p.isEnd() && p.currentToken().TokenKind == marker {
		markerCount++
		p.advance()
	}

	var nodes []NodeInterface
	for !p.isEnd() {
		if p.currentToken().TokenKind == marker {
			count := 0
			start := p.Current
			for !p.isEnd() && p.currentToken().TokenKind == marker {
				count++
				p.advance()
			}

			if count == markerCount {
				break
			} else {
				p.Current = start
			}
		}

		child := p.parseInline()
		nodes = append(nodes, child)
	}

	var result NodeInterface
	switch markerCount {
	case 3:
		result = BoldNode{Nodes: []NodeInterface{ItalicNode{Nodes: nodes}}}
	case 2:
		result = BoldNode{Nodes: nodes}
	case 1:
		result = ItalicNode{Nodes: nodes}
	}

	return result
}

func (p *Parser) parseHeader() NodeInterface {
	level := 1

	p.advance()
	for p.match(lex.Hashtag) {
		level++
		p.advance()
	}
	if p.match(lex.WhiteSpace) {
		p.advance()
	}

	var children []NodeInterface
	for !p.isEnd() && !p.match(lex.NewLine) {
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

func (p *Parser) peek() lex.Token {
	p.advance()
	next := p.currentToken()
	p.recede()

	return next
}

func NewParser(tokens []lex.Token) *Parser {
	var parser Parser
	parser.Tokens = tokens

	return &parser
}
