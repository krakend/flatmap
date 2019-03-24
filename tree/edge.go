package tree

import (
	"bytes"
	"fmt"
	"strings"
)

const tabSize = 4

type edge struct {
	label string
	n     *node
}

func (e *edge) writeTo(bd *strings.Builder, tabList []bool) {
	length := len(tabList)
	isLast, tlist := tabList[length-1], tabList[:length-1]
	for _, hasTab := range tlist {
		if hasTab {
			bd.Write(bytes.Repeat([]byte(" "), tabSize))
			continue
		}
		bd.WriteRune('│')
		bd.Write(bytes.Repeat([]byte(" "), tabSize-1))
	}
	if !isLast {
		bd.WriteRune('├')
	} else {
		bd.WriteRune('└')
	}
	bd.WriteString("── ")
	bd.WriteString(e.label)
	if e.n.IsLeaf() {
		fmt.Fprintf(bd, "\t%+v", e.n.Value)
	} else if e.n.isCollection {
		bd.WriteString(" []")
	}
	bd.WriteByte('\n')
	for i, next := range e.n.edges {
		if len(tabList) < next.n.depth { // runs only for the first edge
			tabList = append(tabList, i == len(e.n.edges)-1)
		} else {
			tabList[next.n.depth-1] = i == len(e.n.edges)-1
		}
		next.writeTo(bd, tabList)
	}
}
