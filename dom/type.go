package dom

type NodeType uint

const (
	TextNode NodeType = iota
	DocumentNode
	ElementNode
	CommentNode
	DoctypeNode
)
