package nodes

import "github.com/3elDU/dark/dom"

type CommentData struct {
	comment string
}

func NewCommentData(comment string) CommentData {
	return CommentData{comment}
}

func (d CommentData) String() string {
	return "<!-- " + d.comment + "-->"
}
func (CommentData) Type() dom.NodeType {
	return dom.CommentNode
}
func (CommentData) Text() string {
	return ""
}
func (CommentData) Attributes() []dom.Attr {
	return nil
}
