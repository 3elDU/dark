package nodes

import (
	"fmt"
	"strings"

	"github.com/3elDU/dark/dom"
)

type TextData struct {
	baseNode
	text string
}

func NewTextData(text string) TextData {
	return TextData{
		baseNode{},
		text,
	}
}

func (t TextData) String() string {
	if len(t.text) > 1024 {
		return fmt.Sprintf("omitted, text was too long (%v characters)", len(t.text))
	}
	return strings.TrimSpace(t.text)
}
func (TextData) Type() dom.NodeType {
	return dom.TextNode
}
func (t TextData) Text() string {
	return t.text
}
