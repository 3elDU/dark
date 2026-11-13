package nodes

import "github.com/3elDU/dark/dom"

type DocumentData struct{}

func NewDocumentData() DocumentData {
	return DocumentData{}
}

func (d DocumentData) String() string {
	return "<html>"
}
func (DocumentData) Type() dom.NodeType {
	return dom.DocumentNode
}
func (d DocumentData) Text() string {
	return ""
}
func (d DocumentData) Attributes() []dom.Attr {
	return nil
}
