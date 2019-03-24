package tree

import (
	"fmt"
	"sort"
	"strings"
)

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

	newNode := &node{
		edges: []*edge{},
		depth: n.depth + 1,
	}
	n.edges = append(n.edges, &edge{label: ks[0], n: newNode})
	newNode.Add(ks[1:], v)
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

		n.edges = []*edge{}
		return
	}

	for i, e := range n.edges {
		if e.label == ks[0] {
			if lenKs == 1 {
				if i == 0 {
					n.edges = n.edges[1:]
					return
				}
				if i == len(n.edges)-1 {
					n.edges = n.edges[:i]
					return
				}
				n.edges = append(n.edges[:i], n.edges[i+1:]...)
				return
			}
			e.n.Del(ks[1:]...)
			return
		}
	}
}

func (n *node) Get(ks ...string) interface{} {
	lenKs := len(ks)
	if n.IsLeaf() && lenKs > 0 {
		return nil
	}

	if lenKs == 0 {
		return n.expand()
	}

	if lenKs == 1 {
		for _, e := range n.edges {
			if e.label == ks[0] {
				return e.n.Get()
			}
		}
		return nil
	}

	for _, e := range n.edges {
		if e.label == ks[0] {
			return e.n.Get(ks[1:]...)
		}
	}
	return nil
}

// Depth returns the node's depth.
func (n *node) Depth() int {
	return n.depth
}

func (n *node) SetDepth(d int) {
	n.depth = d
	for _, e := range n.edges {
		e.n.SetDepth(d + 1)
	}
}

// IsLeaf returns whether the node is a leaf.
func (n *node) IsLeaf() bool {
	return len(n.edges) == 0
}

func (n *node) expand() interface{} {
	if n.IsLeaf() {
		return n.Value
	}

	if n.isCollection {
		res := make([]interface{}, len(n.edges))
		for i, e := range n.edges {
			res[i] = e.n.Get()
		}
		return res
	}

	res := map[string]interface{}{}
	for _, e := range n.edges {
		res[e.label] = e.n.Get()
	}
	return res
}

func (n *node) flatten(i interface{}) {
	switch v := i.(type) {
	case map[string]interface{}:
		n.isCollection = false
		for k, e := range v {
			n.Add([]string{k}, e)
		}
	case []interface{}:
		n.isCollection = true
		// update(append(ks, "#"), len(vs))
		for i, e := range v {
			n.Add([]string{fmt.Sprintf("%d", i)}, e)
		}
	default:
		n.isCollection = false
		n.Value = v
	}
}

func (n *node) sort() {
	s := &sorter{
		n: n,
	}
	sort.Sort(s)
	for _, e := range n.edges {
		e.n.sort()
	}
}

func (n *node) writeTo(bd *strings.Builder) {
	for i, e := range n.edges {
		e.writeTo(bd, []bool{i == len(n.edges)-1})
	}
}

type sorter struct {
	n *node
}

func (s *sorter) Len() int {
	return len(s.n.edges)
}

func (s *sorter) Less(i, j int) bool {
	if s.n.isCollection {
		return i < j
	}
	return s.n.edges[i].label < s.n.edges[j].label
}

func (s *sorter) Swap(i, j int) {
	s.n.edges[i], s.n.edges[j] = s.n.edges[j], s.n.edges[i]
}
