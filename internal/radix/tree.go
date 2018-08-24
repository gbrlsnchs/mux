package radix

import (
	"bytes"
	"net/http"
)

const (
	placeholder = ':'
	delim       = '/'
)

// Tree is a radix tree.
type Tree struct {
	root *Node
}

// New creates a radix tree and initializes its root.
func New() *Tree {
	return &Tree{
		root: &Node{edges: make([]*edge, 0)},
	}
}

// Add adds a new node to the tree according to the label.
func (t *Tree) Add(label []byte, handler http.Handler) {
	if len(label) == 0 || handler == nil {
		return
	}

	tnode := t.root
	for {
		var next *edge
		var slice []byte
		sum := 0
	walk:
		for _, edge := range tnode.edges {
			slice = edge.label
			for i := range slice {
				if i < len(label) && slice[i] == label[i] {
					sum++
					continue
				}
				break
			}
			if sum > 0 {
				label = label[sum:]
				slice = slice[sum:]
				next = edge
				break walk
			}
		}

		if next != nil {
			tnode = next.node
			tnode.priority++

			// Match the whole word.
			if len(label) == 0 {
				// The label is exactly the same as the edge's label,
				// so just replace its node's value.
				//
				// Example:
				// 	(root) -> ("tomato", v1)
				// 	becomes
				// 	(root) -> ("tomato", v2)
				if len(slice) == 0 {
					tnode.Handler = handler
					goto sort
				}
				// The label is a prefix of the edge's label.
				//
				// Example:
				// 	(root) -> ("tomato", v1)
				// 	then add "tom"
				// 	(root) -> ("tom", v2) -> ("ato", v1)
				tnode.edges = []*edge{newEdge(slice, tnode.sibling())}
				next.label = next.label[:len(next.label)-len(slice)]
				tnode.Handler = handler
				goto sort
			}

			// Add a new node but break its parent into prefix and
			// the remaining slice as a new edge.
			//
			// Example:
			// 	(root) -> ("tomato", v1)
			// 	then add "tornado"
			// 	(root) -> ("to", nil) -> ("mato", v1)
			// 	                      +> ("rnado", v2)
			if len(slice) > 0 {
				tnode.edges = []*edge{
					newEdge(slice, tnode.sibling()),  // The suffix that is split into a new node.
					newEdge(label, newNode(handler)), // The new node.
				}
				next.label = next.label[:len(next.label)-len(slice)]
				tnode.Handler = nil
				goto sort
			}
			continue
		}

		tnode.edges = append(tnode.edges, newEdge(label, newNode(handler)))
	}
sort:
	t.root.sort()
}

// Get retrieves an http.Handler and a map of parameters according to the matching label.
func (t *Tree) Get(label []byte) (*Node, map[string]string) {
	if len(label) == 0 {
		return nil, nil
	}
	tnode := t.root
	var params map[string]string
	for tnode != nil && len(label) > 0 {
		var next *edge
	walk:
		for _, edge := range tnode.edges {
			slice := edge.label
			for {
				// Check if there are any placeholders.
				// If there are none, then use the whole word for comparison.
				phIndex := bytes.IndexByte(slice, placeholder)
				if phIndex < 0 {
					phIndex = len(slice)
				}

				prefix := slice[:phIndex]

				// If "slice" (until placeholder) is not prefix of
				// "label", then keep walking.
				if !bytes.HasPrefix(label, prefix) {
					continue walk
				}

				label = label[len(prefix):]

				// If "slice" is the whole label,
				// then the match is complete and the algorithm
				// is ready to go to the next edge.
				if len(prefix) == len(slice) {
					next = edge

					break walk
				}

				// Check whether there is a delimiter.
				// If there isn't, then use the whole world as parameter.
				var delimIndex int
				slice = slice[phIndex:]
				if delimIndex = bytes.IndexByte(slice[1:], delim) + 1; delimIndex <= 0 {
					delimIndex = len(slice)
				}

				key := slice[1:delimIndex] // Remove the placeholder from the map key.
				slice = slice[delimIndex:]
				if delimIndex = bytes.IndexByte(label[1:], delim) + 1; delimIndex <= 0 {
					delimIndex = len(label)
				}

				if params == nil {
					params = make(map[string]string)
				}
				params[string(key)] = string(label[:delimIndex])

				label = label[delimIndex:]
				if len(slice) == 0 && len(label) == 0 {
					next = edge
					break walk
				}
			}
		}
		if next != nil {
			tnode = next.node
			continue
		}
		tnode = nil
	}
	return tnode, params
}
