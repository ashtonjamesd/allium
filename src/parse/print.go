package parse

import "fmt"

func (n ParagraphNode) Print(indent int) {
	fmt.Printf("%sParagraphNode: \n", spaces(indent))
	for _, child := range n.Content {
		printNode(child, indent+2)
	}
}

func (n ItalicNode) Print(indent int) {
	fmt.Printf("%sItalicNode: \n", spaces(indent))
	for _, child := range n.Nodes {
		printNode(child, indent+2)
	}
}

func (n BoldNode) Print(indent int) {
	fmt.Printf("%sBoldNode: \n", spaces(indent))
	for _, child := range n.Nodes {
		printNode(child, indent+2)
	}
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

func (n HorizontalRuleNode) Print(indent int) {
	fmt.Printf("%sHorizontalRuleNode: '%s'\n", spaces(indent), "---")
}

func (n ListItemNode) Print(indent int) {
	fmt.Printf("%sUnorderedListNode: \n", spaces(indent))
	for _, child := range n.Nodes {
		printNode(child, indent+2)
	}
}

func (n ListNode) Print(indent int) {
	fmt.Printf("%sUnorderedList: \n", spaces(indent))
	for _, child := range n.Nodes {
		printNode(child, indent+2)
	}
}

func (n LinkNode) Print(indent int) {
	fmt.Printf("%sLinkNode: '%s' : '%s'\n", spaces(indent), n.LinkText, n.Link)
}

func (n ImageNode) Print(indent int) {
	fmt.Printf("%sImageNode: '%s' : '%s'\n", spaces(indent), n.LinkText, n.Link)
}

func (n InlineCodeBlockNode) Print(indent int) {
	fmt.Printf("%sInlineCodeBlockNode: '%s'\n", spaces(indent), n.Content)
}

func (n InlineCodeNode) Print(indent int) {
	fmt.Printf("%sInlineCodeNode: '%s'\n", spaces(indent), n.Content)
}

func (n BlockQuoteNode) Print(indent int) {
	fmt.Printf("%sBlockQuoteNode: \n", spaces(indent))
	for _, child := range n.Nodes {
		printNode(child, indent+2)
	}
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
	case ListItemNode:
		node.Print(indent)
	case ListNode:
		node.Print(indent)
	case HorizontalRuleNode:
		node.Print(indent)
	case BlockQuoteNode:
		node.Print(indent)
	case ImageNode:
		node.Print(indent)
	case InlineCodeBlockNode:
		node.Print(indent)
	case InlineCodeNode:
		node.Print(indent)
	default:
		fmt.Printf("%sUnknown node type\n", spaces(indent))
	}
}

func PrintNodes(nodes []NodeInterface) {
	for _, node := range nodes {
		printNode(node, 0)
	}
}

func spaces(n int) string {
	return fmt.Sprintf("%*s", n, "")
}
