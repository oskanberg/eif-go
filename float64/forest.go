package eif

import "math"

const (
	numberOfTreesDefault = 10
)

func cFactor(n float64) float64 {
	return 2.0*(math.Log(n-1)+0.5772156649) - (2.0 * (n - 1.) / (n * 1.0))
}

// Forest is a group of trees
type Forest struct {
	trees []*Node
	c     float64
}

// Score calculates the anomaloussnes score (between 0.0 and 1.0) of a given point.
// Higher scores imply more anomaloussnes. Note that scores are affected by maximum tree depth.
func (f *Forest) Score(p []float64) float64 {
	var totalDepth int
	for _, t := range f.trees {
		totalDepth += t.Depth(p)
	}

	avg := float64(totalDepth) / float64(len(f.trees))
	return math.Pow(2, -avg/f.c)
}

// ForestParams is the set of metadata used to construct a forest
type forestParams struct {
	numberOfTrees int
	maxTreeDepth  int
}

// ForestOpt is an option to adapt ForestParams
type ForestOpt func(*forestParams)

// WithTrees is a construction option, setting the amount of trees to generate
func WithTrees(n int) ForestOpt {
	return func(f *forestParams) {
		f.numberOfTrees = n
	}
}

// WithMaxTreeDepth is a construction option, which determines how deep trees are constructed
// In general, running with smaller values is faster, but loses fidelity in anomaly scores
func WithMaxTreeDepth(n int) ForestOpt {
	return func(f *forestParams) {
		f.maxTreeDepth = n
	}
}

// NewForest creates a new forest, with optional construction parameters.
// By default 10 trees are constructed, and a maximum depth determined by the expected depth
// of an unsuccesful binary tree search with the given data size.
func NewForest(data [][]float64, opts ...ForestOpt) *Forest {
	sampleSize := float64(len(data))
	params := forestParams{
		numberOfTrees: numberOfTreesDefault,
		maxTreeDepth:  int(math.Ceil(math.Log2(sampleSize))),
	}
	for _, fn := range opts {
		fn(&params)
	}

	trees := make([]*Node, params.numberOfTrees)
	for i := range trees {
		trees[i] = NewTree(data, params.maxTreeDepth)
	}

	return &Forest{
		trees: trees,
		c:     cFactor(sampleSize),
	}
}
