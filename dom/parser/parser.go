package parser

import (
	"strings"

	"github.com/3elDU/dark/dom"
	"github.com/3elDU/dark/dom/nodes"
	"golang.org/x/net/html"
)

func ParseDocument(page string) (*dom.Dom, error) {
	doc, err := html.Parse(strings.NewReader(page))
	if err != nil {
		return nil, err
	}

	dom := dom.NewDom()
	tree := convertTree(doc)
	dom.ReplaceRoot(tree)

	return dom, nil
}

func convertTree(node *html.Node) *dom.Node {
	var data dom.NodeData

	switch node.Type {
	case html.TextNode:
		data = nodes.NewTextData(node.Data)
	case html.DocumentNode:
		data = nodes.NewDocumentData()
	case html.ElementNode:
		data = nodes.NewElementData(node.Data, convertAttrs(node.Attr))
	case html.CommentNode:
		data = nodes.NewCommentData(node.Data)
	case html.DoctypeNode:
		data = nodes.NewDoctypeData()
	}

	parent := dom.NewNode(data)

	for child := range node.ChildNodes() {
		cChild := convertTree(child)

		if err := parent.AppendChildren(cChild); err != nil {
			// an error here means a bug in x/net/html
			panic(err)
		}
	}

	return parent
}

func convertAttrs(attrs []html.Attribute) dom.AttributeList {
	out := dom.AttributeList{}

	for _, attr := range attrs {
		out[attr.Key] = attr.Val
	}

	return out
}
