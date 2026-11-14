package nodes

import "github.com/3elDU/dark/dom"

type baseNode struct{}

func (b baseNode) TagName() string {
	return ""
}
func (b baseNode) String() string {
	return ""
}
func (b baseNode) Text() string {
	return ""
}
func (b baseNode) Attributes() dom.AttributeList {
	return nil
}
