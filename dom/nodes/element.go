package nodes

import (
	"strings"

	"github.com/3elDU/dark/dom"
)

type ElementData struct {
	tagName string
	attrs   dom.AttributeList
}

func NewElementData(tagName string, attrs dom.AttributeList) ElementData {
	return ElementData{tagName, attrs}
}

func (e ElementData) TagName() string {
	return e.tagName
}

func (e ElementData) String() string {
	var sb strings.Builder

	sb.WriteRune('<')
	sb.WriteString(e.tagName)

	if len(e.attrs) > 0 {
		sb.WriteRune(' ')
	}

	i := 0
	for k, attr := range e.attrs {
		sb.WriteString(k)
		sb.WriteRune('=')
		sb.WriteRune('"')
		sb.WriteString(attr)
		sb.WriteRune('"')

		if i < len(e.attrs)-1 {
			sb.WriteRune(' ')
		}
		i++
	}

	sb.WriteRune('>')

	return sb.String()
}
func (ElementData) Type() dom.NodeType {
	return dom.ElementNode
}
func (ElementData) Value() *string {
	return nil
}
func (e ElementData) Attributes() dom.AttributeList {
	return e.attrs
}
