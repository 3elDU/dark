package nodes

import "github.com/3elDU/dark/dom"

type DocumentData struct {
	baseNode
}

func NewDocumentData() DocumentData {
	return DocumentData{}
}

func (d DocumentData) String() string {
	return "<html>"
}
func (DocumentData) Type() dom.NodeType {
	return dom.DocumentNode
}
