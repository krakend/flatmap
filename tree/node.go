package tree

import (
	"sort"
	"strconv"
)

type edge struct {
	label string
	n     *node
}

func newNode(depth int) *node {
	return &node{
		edges: []*edge{},
		depth: depth,
	}
}

type node struct {
	Value        interface{}
	isCollection bool
	edges        []*edge
	depth        int
}

func (n *node) Add(ks []string, v interface{}) {
	if len(ks) == 0 {
		n.flatten(v)
		return
	}

	for _, e := range n.edges {
		if e.label == ks[0] {
			e.n.Add(ks[1:], v)
			return
		}
	}

	child := newNode(n.depth + 1)
	n.edges = append(n.edges, &edge{label: ks[0], n: child})
	child.Add(ks[1:], v)
}

func (n *node) Del(ks ...string) {
	lenKs := len(ks)

	if lenKs == 0 || n.IsLeaf() {
		return
	}

	if ks[0] == wildcard {
		if lenKs > 1 {
			for _, e := range n.edges {
				e.n.Del(ks[1:]...)
			}
			return
		}

		for i := range n.edges {
			n.edges[i] = nil
		}
		n.edges = n.edges[:0]
		return
	}

	for i, e := range n.edges {
		if e.label == ks[0] {
			if lenKs == 1 {
				copy(n.edges[i:], n.edges[i+1:])
				n.edges[len(n.edges)-1] = nil
				n.edges = n.edges[:len(n.edges)-1]
				return
			}
			e.n.Del(ks[1:]...)
			return
		}
	}
}

func (n *node) Get(ks ...string) interface{} {
	lenKs := len(ks)
	lenEdges := len(n.edges)

	if lenEdges == 0 && lenKs > 0 {
		return nil
	}

	if lenKs == 0 {
		return n.expand()
	}

	if ks[0] == wildcard {
		res := make([]interface{}, lenEdges)
		for i, e := range n.edges {
			res[i] = e.n.Get(ks[1:]...)
		}
		return res
	}

	for _, e := range n.edges {
		if e.label == ks[0] {
			return e.n.Get(ks[1:]...)
		}
	}
	return nil
}

func (n *node) Depth() int {
	return n.depth
}

func (n *node) SetDepth(d int) {
	n.depth = d
	for _, e := range n.edges {
		e.n.SetDepth(d + 1)
	}
}

func (n *node) IsLeaf() bool {
	return len(n.edges) == 0
}

func (n *node) expand() interface{} {
	children := len(n.edges)
	if children == 0 {
		return n.Value
	}

	if n.isCollection {
		res := make([]interface{}, children)
		for i, e := range n.edges {
			res[i] = e.n.Get()
		}

		return res
	}

	res := make(map[string]interface{}, children)
	for _, e := range n.edges {
		res[e.label] = e.n.Get()
	}
	return res
}

func (n *node) flatten(i interface{}) {
	switch v := i.(type) {
	case map[string]interface{}:
		n.isCollection = false
		if len(v) == 0 {
			n.Value = v
			break
		}

		for k, e := range v {
			n.Add([]string{k}, e)
		}
	case []interface{}:
		n.isCollection = true
		if len(v) == 0 {
			n.Value = v
			break
		}

		for i, e := range v {
			n.Add([]string{strconv.Itoa(i)}, e)
		}
	default:
		n.isCollection = false
		n.Value = v
	}
}

func (n *node) sort() {
	if n.IsLeaf() {
		return
	}

	if !n.isCollection {
		sort.Sort(n)
	}

	for _, e := range n.edges {
		e.n.sort()
	}
}

func (n *node) Len() int {
	return len(n.edges)
}

func (n *node) Less(i, j int) bool {
	if n.isCollection {
		return i < j
	}
	return n.edges[i].label < n.edges[j].label
}

func (n *node) Swap(i, j int) {
	n.edges[i], n.edges[j] = n.edges[j], n.edges[i]
}
