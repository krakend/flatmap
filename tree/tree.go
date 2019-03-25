package tree

import (
	"errors"
)

const wildcard = "*"

var errNoNilValuesAllowed = errors.New("no nil values allowed")

type Tree struct {
	root *node
}

func New(v interface{}) (*Tree, error) {
	if v == nil {
		return nil, errNoNilValuesAllowed
	}

	tr := &Tree{
		root: &node{},
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

func (t *Tree) Move(src, dst []string) {
	next := []nodeAndPath{{n: t.root, p: []string{}}}
	acc := []nodeAndPath{}

	prefixLen := len(src)

	// lookup for the nodes matching the src pattern without the last segment

	if prefixLen > 1 {
		for _, step := range src[:prefixLen-1] {
			if step == wildcard {
				for _, nap := range next {
					for _, e := range nap.n.edges {
						acc = append(acc, nodeAndPath{n: e.n, p: append(nap.p, e.label)})
					}
				}
			} else {
				for _, nap := range next {
					for _, e := range nap.n.edges {
						if step == e.label {
							acc = append(acc, nodeAndPath{n: e.n, p: append(nap.p, e.label)})
							break
						}
					}
				}
			}
			next = acc
			acc = []nodeAndPath{}
		}
	}

	// extract edges from the collected nodes that match the last segment of the src

	edgesToMove := []edgeToMove{}
	lenDst := len(dst)

	// if len(src) == len(dst) this is just a relabeling
	isEdgeRelabel := prefixLen == lenDst

	for _, nap := range next {
		for i, e := range nap.n.edges {
			if e.label != src[prefixLen-1] {
				continue
			}

			if isEdgeRelabel {
				e.label = dst[prefixLen-1]
				break
			}

			edgesToMove = append(edgesToMove, edgeToMove{nodeAndPath: nap, e: e})

			numOfEdges := len(nap.n.edges)
			if numOfEdges == 1 {
				nap.n.edges = []*edge{}
				break
			}

			switch i {
			case 0:
				nap.n.edges = nap.n.edges[1:]
			case numOfEdges - 1:
				nap.n.edges = nap.n.edges[:len(nap.n.edges)-1]
			default:
				nap.n.edges = append(nap.n.edges[:i], nap.n.edges[i+1:]...)
			}
			break
		}
	}

	if isEdgeRelabel {
		return
	}

	// insert the extracted values

	// if len(src) > len(dst) this is a promotion
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

	// this is an embedding, so intermediate nodes should be found or created
	found := false
	for _, em := range edgesToMove {
		root := em.n
		for _, k := range dst[prefixLen-1 : lenDst-1] {
			found = false
			for _, e := range root.edges {
				if e.label != k {
					continue
				}
				found = true
				root = e.n
				break
			}
			if found {
				continue
			}
			child := newNode(root.depth + 1)
			root.edges = append(root.edges, &edge{label: k, n: child})
			root = child
		}
		em.e.label = dst[lenDst-1]
		em.e.n.SetDepth(root.depth + 1)
		root.edges = append(root.edges, em.e)
	}
}

func (t *Tree) Sort() {
	t.root.sort()
}

type nodeAndPath struct {
	n *node
	p []string
}

type edgeToMove struct {
	nodeAndPath
	e *edge
}
