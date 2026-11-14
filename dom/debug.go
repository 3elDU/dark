package dom

import (
	"fmt"
	"strings"
)

// Walk the DOM and print it
func (d *Dom) Print() {
	walkAndPrint(d, 0, d.GetRootNode())
}

func walkAndPrint(dom *Dom, depth int, node *Node) {
	str := node.Data.String()

	// skip empty tags
	if str != "" {
		fmt.Println(strings.Repeat(" ", depth*4) + str)
	}

	for child := range dom.ChildNodes(node) {
		walkAndPrint(dom, depth+1, child)
	}
}
