package nodes

import "github.com/3elDU/dark/dom"

type CommentData struct {
	baseNode
	comment string
}

func NewCommentData(comment string) CommentData {
	return CommentData{
		baseNode{},
		comment,
	}
}

func (c CommentData) String() string {
	return "<!-- " + c.comment + "-->"
}
func (CommentData) Type() dom.NodeType {
	return dom.CommentNode
}
func (c CommentData) Text() string {
	return c.comment
}
