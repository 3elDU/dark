package nodes

import (
	"strings"

	"github.com/3elDU/dark/dom"
)

type ElementData struct {
	tagName string
	attrs   []dom.Attr
}

func NewElementData(tagName string, attrs []dom.Attr) ElementData {
	return ElementData{tagName, attrs}
}

func (e ElementData) String() string {
	var sb strings.Builder

	sb.WriteRune('<')
	sb.WriteString(e.tagName)

	if len(e.attrs) > 0 {
		sb.WriteRune(' ')
	}

	for i, attr := range e.attrs {
		sb.WriteString(attr.Key)
		sb.WriteRune('=')
		sb.WriteRune('"')
		sb.WriteString(attr.Value)
		sb.WriteRune('"')

		if i < len(e.attrs)-1 {
			sb.WriteRune(' ')
		}
	}

	sb.WriteRune('>')

	return sb.String()
}
func (ElementData) Type() dom.NodeType {
	return dom.ElementNode
}
func (ElementData) Text() string {
	return ""
}
func (e ElementData) Attributes() []dom.Attr {
	return e.attrs
}
