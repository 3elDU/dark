// Package dom implements parsing and structures to represent the final DOM tree
// Nil pointers are used extensively in place of *Node values, meaning the abscence of node.
package dom

import (
	"errors"
	"iter"
	"slices"
	"strings"
)

// The document tree.
type Dom struct {
	// The root node. Empty for a fresh DOM tree.
	// [Dom.ReplaceRoot] needs to be called to set the root.
	Root *Node
}

func NewDom() *Dom {
	return &Dom{
		Root: nil,
	}
}

func (dom *Dom) ReplaceRoot(root *Node) {
	dom.Root = root
}

func (dom *Dom) GetElementsByTagName(tagName string) []Element {
	matching := []Element{}

	for node := range dom.Root.Walk() {
		el, ok := node.AsElement()
		if !ok {
			continue
		}

		if el.TagName() == tagName {
			matching = append(matching, el)
		}
	}

	return matching
}

type Node struct {
	parent   *Node
	children []*Node
	Data     NodeData
}

// Creates a new node with no parent and empty list of children
func NewNode(data NodeData) *Node {
	return &Node{
		parent:   nil,
		children: []*Node{},
		Data:     data,
	}
}

func (n *Node) GetParent() *Node {
	return n.parent
}

func (n *Node) Children() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for _, c := range n.children {
			if !yield(c) {
				return
			}
		}
	}
}

var ErrHierarchyRequest = errors.New("HierarchyRequestError")

// Add one or more children to this node.
// If a child node already had a parent, it is reparented to this node.
//
// [ErrHierarchyRequest] is returned when a cycle is detected.
func (n *Node) AppendChildren(children ...*Node) error {
	for _, child := range children {
		// attempt to add a node to itself
		if child == n {
			return ErrHierarchyRequest
		}
		// attempt to add an ancestor to it's descendant children
		for ancestor := range n.WalkAncestors() {
			if ancestor == child {
				return ErrHierarchyRequest
			}
		}

		// remove node from previous parent
		if child.parent != nil {
			child.parent.RemoveChild(child)
		}

		child.parent = n
		n.children = append(n.children, child)
	}

	return nil
}

func (n *Node) RemoveChild(child *Node) {
	idx := slices.Index(n.children, child)
	if idx == -1 {
		return
	}

	n.children = slices.Delete(n.children, idx, idx+1)
	child.parent = nil
}

// Walk returns the iterator for every node beneath this node, including itself.
func (n *Node) Walk() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		n.walk(yield)
	}
}

// reuses the same yield function
func (n *Node) walk(yield func(*Node) bool) bool {
	if !yield(n) {
		return false
	}

	for _, child := range n.children {
		if !child.walk(yield) {
			return false
		}
	}

	return true
}

// WalkAncestors iterates over the node's ancestor nodes until a nil parent is found
func (n *Node) WalkAncestors() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		n.walkAncestors(yield)
	}
}

// reuses the same yield function
func (n *Node) walkAncestors(yield func(*Node) bool) bool {
	if n.parent == nil {
		return false
	}

	if !yield(n.parent) {
		return false
	}

	return n.parent.walkAncestors(yield)
}

func (n *Node) GetTextContent() string {
	switch n.Data.Type() {
	case TextNode:
		return *n.Data.Value()
	case CommentNode:
		return *n.Data.Value()
	case ElementNode:
		var sb strings.Builder

		for child := range n.Walk() {
			if child.Data.Type() == TextNode {
				sb.WriteString(*child.Data.Value())
			}
		}

		return sb.String()
	default:
		return ""
	}
}

type NodeData interface {
	// String representation of the node for debugging
	String() string
	Type() NodeType
	// https://developer.mozilla.org/en-US/docs/Web/API/Node/nodeValue
	//
	// tl;dr: string for text and comment, nil for every other node type
	Value() *string
}

type Element struct {
	*Node
}

func (e *Element) TagName() string {
	return e.Data.(ElementData).TagName()
}

func (e *Element) Attributes() AttributeList {
	return e.Data.(ElementData).Attributes()
}

// Data specific to an [ElementNode]
type ElementData interface {
	NodeData

	TagName() string
	Attributes() AttributeList
}

type AttributeList map[string]string

// Get the attribute value. Returns an empty string if an attribute doesn't exist.
func (al AttributeList) Get(name string) string {
	if val, ok := al[name]; ok {
		return val
	}
	return ""
}

// Whether the element has an attribute with the given name (a defined empty attribute returns true)
func (al AttributeList) Has(name string) bool {
	_, has := al[name]
	return has
}

// Attempts to cast this node to an element. ok is false if the node is not an element
func (n *Node) AsElement() (e Element, ok bool) {
	data, ok := n.Data.(ElementData)
	if !ok || data.Type() != ElementNode {
		return Element{}, false
	}

	return Element{n}, true
}
