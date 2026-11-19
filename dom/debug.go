package dom

import (
	"fmt"
	"strings"
)

// Walk the DOM and print it
func (d *Dom) Print() {
	walkAndPrint(0, d.Root)
}

func walkAndPrint(depth int, node *Node) {
	str := node.Data.String()

	// skip empty tags
	if str != "" {
		fmt.Println(strings.Repeat(" ", depth*4) + str)
	}

	for _, child := range node.children {
		walkAndPrint(depth+1, child)
	}
}
