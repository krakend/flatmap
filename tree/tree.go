package tree

import (
	"errors"
	"strings"
)

const wildcard = "*"

type Tree struct {
	root *node
	sb   *strings.Builder
}

var errNoNilValuesAllowed = errors.New("no nil values allowed")

func New(v interface{}) (*Tree, error) {
	if v == nil {
		return nil, errNoNilValuesAllowed
	}

	tr := &Tree{
		root: &node{},
		sb:   new(strings.Builder),
	}

	tr.Add([]string{}, v)

	return tr, nil
}

func (t *Tree) Add(ks []string, v interface{}) {
	if v == nil {
		return
	}
	t.root.Add(ks, v)
}

func (t *Tree) Get(ks []string) interface{} {
	return t.root.Get(ks...)
}

func (t *Tree) Del(ks []string) {
	t.root.Del(ks...)
}

type nodeAndPath struct {
	n *node
	p []string
}

type edgeToMove struct {
	nodeAndPath
	e *edge
}

func (t *Tree) Move(src, dst []string) {
	next := []nodeAndPath{{n: t.root, p: []string{}}}
	acc := []nodeAndPath{}

	prefixLen := len(src)

	if prefixLen > 1 {
		for _, step := range src[:prefixLen-1] {
			if step == wildcard {
				for _, n := range next {
					for _, e := range n.n.edges {
						acc = append(acc, nodeAndPath{n: e.n, p: append(n.p, e.label)})
					}
				}
			} else {
				for _, n := range next {
					for _, e := range n.n.edges {
						if step == e.label {
							acc = append(acc, nodeAndPath{n: e.n, p: append(n.p, e.label)})
							break
						}
					}
				}
			}
			next = acc
			acc = []nodeAndPath{}
		}
	}

	edgesToMove := []edgeToMove{}
	lenDst := len(dst)
	isEdgeRelabel := prefixLen == lenDst
	// extract from the actual tree
	for _, n := range next {
		for i, e := range n.n.edges {
			if e.label != src[prefixLen-1] {
				continue
			}

			if isEdgeRelabel {
				e.label = dst[prefixLen-1]
				break
			}

			edgesToMove = append(edgesToMove, edgeToMove{nodeAndPath: n, e: e})
			switch i {
			case 0:
				n.n.edges = n.n.edges[1:]
			case len(n.n.edges) - 1:
				n.n.edges = n.n.edges[:len(n.n.edges)-1]
			default:
				n.n.edges = append(n.n.edges[:i], n.n.edges[i+1:]...)
			}
			break
		}
	}

	if isEdgeRelabel {
		return
	}

	// insert the extracted values
	if prefixLen > lenDst {
		for _, n := range edgesToMove {
			parent := t.root
			for i, path := range dst[:lenDst-1] {
				l := path
				if path == wildcard {
					l = n.p[i]
				}
				found := false
				for _, e := range parent.edges {
					if e.label != l {
						continue
					}
					found = true
					parent = e.n
					break
				}
				if !found {
					break
				}
			}

			n.e.n.SetDepth(parent.depth + 1)
			n.e.label = dst[lenDst-1]
			parent.edges = append(parent.edges, n.e)
		}
		return
	}

	// TODO: create the required intermediate nodes if prefixLen < lenDst
}

func (t *Tree) Sort() {
	t.root.sort()
}

// String returns a string representation of the tree structure.
func (t *Tree) String() string {
	t.sb.Reset()
	t.sb.WriteByte('\n')
	t.root.writeTo(t.sb)
	return t.sb.String()
}
