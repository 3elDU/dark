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
	convertTree(dom, doc)

	return dom, nil
}

func convertTree(tree *dom.Dom, node *html.Node) dom.NodeHandle {
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

	tree.AddNode(parent)

	for child := range node.ChildNodes() {
		cChild := convertTree(tree, child)
		tree.AppendChild(parent.Id(), cChild)
	}

	return parent.Id()
}

func convertAttrs(attrs []html.Attribute) dom.AttributeList {
	out := dom.AttributeList{}

	for _, attr := range attrs {
		out[attr.Key] = attr.Val
	}

	return out
}
