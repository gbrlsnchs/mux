package radix

import (
	"net/http"
	"sort"
)

// Node is a radix tree node.
type Node struct {
	Handler  http.Handler
	edges    []*edge
	priority int
}

func newNode(handler http.Handler) *Node {
	return &Node{
		Handler: handler,
		edges:   make([]*edge, 0),
	}
}

func (n *Node) Len() int {
	return len(n.edges)
}

func (n *Node) Less(i, j int) bool {
	return n.edges[i].node.priority > n.edges[j].node.priority
}

func (n *Node) Swap(i, j int) {
	n.edges[i], n.edges[j] = n.edges[j], n.edges[i]
}

func (n *Node) sibling() *Node {
	return &Node{
		Handler:  n.Handler,
		edges:    n.edges,
		priority: n.priority,
	}
}

func (n *Node) sort() {
	sort.Sort(n)
	for _, edge := range n.edges {
		edge.node.sort()
	}
}
