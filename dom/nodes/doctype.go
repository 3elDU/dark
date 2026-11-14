package nodes

import "github.com/3elDU/dark/dom"

type DoctypeData struct {
	baseNode
}

func NewDoctypeData() DoctypeData {
	return DoctypeData{}
}

func (d DoctypeData) String() string {
	return "<!doctype html>"
}
func (DoctypeData) Type() dom.NodeType {
	return dom.DoctypeNode
}
