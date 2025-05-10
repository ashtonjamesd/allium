package parse

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

type ImageNode struct {
	LinkText string
	Link     string
}

type NoNode struct{}

type ItalicNode struct {
	Nodes []NodeInterface
}

type TextNode struct {
	Content string
}

type BoldNode struct {
	Nodes []NodeInterface
}

type LinkNode struct {
	LinkText string
	Link     string
}

type ListNode struct {
	Nodes     []NodeInterface
	IsOrdered bool
}

type ListItemNode struct {
	Nodes []NodeInterface
}

type BlockQuoteNode struct {
	Nodes []NodeInterface
}

type InlineCodeBlockNode struct {
	Content string
}

type InlineCodeNode struct {
	Content string
}

type HorizontalRuleNode struct{}

type WhiteSpaceNode struct{}
type NewLineNode struct{}
