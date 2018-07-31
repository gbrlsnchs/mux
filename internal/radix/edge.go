package radix

type edge struct {
	label []byte
	node  *Node
}

func newEdge(label []byte, node *Node) *edge {
	return &edge{label: label, node: node}
}
