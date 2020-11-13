package eif

import (
	"fmt"
	"math/rand"
)

// Node is a tree node
type Node struct {
	terminus  bool
	plane     []float64
	intercept []float64
	l, r      *Node
	endData   [][]float64
}

func (n *Node) isInLeftDivision(point []float64) bool {
	if len(n.plane) != len(point) || len(n.plane) != len(n.intercept) {
		panic(fmt.Sprintf("wrong dimensions: %d %d %d", len(n.plane), len(n.intercept), len(point)))
	}

	var d float64
	for i, v := range point {
		d += (v - n.intercept[i]) * n.plane[i]
	}
	return d <= 0
}

func newNode(dims int, data [][]float64, d int) *Node {
	if len(data) <= 1 || d == 0 {
		return &Node{
			terminus: true,
			endData:  data,
		}
	}

	plane := make([]float64, dims)
	intercept := make([]float64, dims)
	max := append([]float64{}, data[0]...)
	min := append([]float64{}, data[0]...)
	for _, v := range data {
		for i, vv := range v {
			if vv > max[i] {
				max[i] = vv
			}
			if vv < min[i] {
				min[i] = vv
			}
		}
	}

	for i := range max {
		plane[i] = rand.NormFloat64()
		intercept[i] = min[i] + rand.Float64()*(max[i]-min[i])
	}

	var l, r [][]float64
	node := &Node{false, plane, intercept, nil, nil, nil}

	for _, v := range data {
		if node.isInLeftDivision(v) {
			l = append(l, v)
		} else {
			r = append(r, v)
		}
	}

	d--
	node.l = newNode(dims, l, d)
	node.r = newNode(dims, r, d)

	return node
}

func (n *Node) depth(p []float64, e int) int {
	if n.terminus {
		return e
	}

	if n.isInLeftDivision(p) {
		return n.l.depth(p, e+1)
	}

	return n.r.depth(p, e+1)
}

// Depth gets the depth of p below this node
func (n *Node) Depth(p []float64) int {
	return n.depth(p, 1)
}

// NewTree constructs a tree (node)
func NewTree(data [][]float64, maxDepth int) *Node {
	return newNode(len(data[0]), data, maxDepth)
}
