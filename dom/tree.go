// Package dom implements parsing and structures to represent the final DOM tree
package dom

import (
	"iter"
	"strings"
)

// A handle to the node. Useful only with the [Dom] structure.
// By itself is an index into it.
type NodeHandle int

const InvalidHandle NodeHandle = -1

// The document tree.
type Dom struct {
	nodes []*Node
	root  NodeHandle
}

func NewDom() *Dom {
	return &Dom{
		nodes: []*Node{},
		root:  InvalidHandle,
	}
}

// Add a node to this tree
func (dom *Dom) AddNode(node *Node) NodeHandle {
	dom.nodes = append(dom.nodes, node)
	node.id = NodeHandle(len(dom.nodes) - 1)

	// set the first node as root
	if node.id == 0 {
		dom.root = 0
	}

	return node.id
}

// Get a node by its handle
func (dom *Dom) GetNode(id NodeHandle) *Node {
	return dom.nodes[id]
}

// Add a node to this tree, as a child of another node. Returns a handle to the child
func (dom *Dom) AppendNewChild(phandle NodeHandle, child *Node) NodeHandle {
	dom.AddNode(child)

	child.parent = phandle
	parent := dom.GetNode(phandle)
	parent.children = append(parent.children, child.id)

	return child.id
}

// Add an existing node as a child to another node
func (dom *Dom) AppendChild(phandle, chandle NodeHandle) {
	parent := dom.GetNode(phandle)
	child := dom.GetNode(chandle)

	if child.parent == phandle {
		// child is already a parent of given node, do nothing
		return
	}

	// remove node from old parent's children
	if child.parent != InvalidHandle {
		oldParent := dom.GetNode(child.parent)

		id := oldParent.childIndex(chandle)
		if id != -1 {
			oldParent.children = append(oldParent.children[:id], oldParent.children[id+1:]...)
		}
	}

	child.parent = phandle
	parent.children = append(parent.children, chandle)
}

// Add many children to a node
func (dom *Dom) AppendChildren(phandle NodeHandle, children ...NodeHandle) {
	for _, chandle := range children {
		dom.AppendChild(phandle, chandle)
	}
}

func (dom *Dom) ChildNodes(node *Node) iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for i := 0; i < len(node.children) && yield(dom.GetNode(node.children[i])); i++ {
		}
	}
}

func (dom *Dom) GetRootNode() *Node {
	if dom.root == InvalidHandle {
		return nil
	}

	return dom.nodes[dom.root]
}

func (dom *Dom) QueryAllByTagName(tagName string) []*Node {
	matching := []*Node{}

	for _, node := range dom.nodes {
		if node.Data.TagName() == tagName {
			matching = append(matching, node)
		}
	}

	return matching
}

func (dom *Dom) GetTextContent(node *Node) string {
	switch node.Data.Type() {
	case TextNode:
		return node.Data.Text()
	case CommentNode:
		return node.Data.Text()
	case ElementNode:
		var sb strings.Builder

		for child := range dom.ChildNodes(node) {
			sb.WriteString(dom.GetTextContent(child))
		}

		return sb.String()
	default:
		return ""
	}
}

type Node struct {
	// ID of this node in the parent tree
	id       NodeHandle
	parent   NodeHandle
	children []NodeHandle
	Data     NodeData
}

func NewNode(data NodeData) *Node {
	return &Node{
		id:       InvalidHandle,
		parent:   InvalidHandle,
		children: []NodeHandle{},
		Data:     data,
	}
}

func (n Node) Id() NodeHandle {
	return n.id
}

// -1 if child does not belong to this node
func (n Node) childIndex(child NodeHandle) int {
	for i, c := range n.children {
		if child == c {
			return i
		}
	}

	return -1
}

type NodeData interface {
	TagName() string
	// String representation of the node for debugging
	String() string
	Type() NodeType
	Text() string
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
