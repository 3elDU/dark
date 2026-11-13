package nodes

import (
	"strings"

	"github.com/3elDU/dark/dom"
)

type TextData struct {
	text string
}

func NewTextData(text string) TextData {
	return TextData{text}
}

func (t TextData) String() string {
	return strings.TrimSpace(t.text)
}
func (TextData) Type() dom.NodeType {
	return dom.TextNode
}
func (t TextData) Text() string {
	return t.text
}
func (t TextData) Attributes() []dom.Attr {
	return nil
}
