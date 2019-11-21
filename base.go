package burovka

import (
	"math"
	"sync"
)

var wg sync.WaitGroup

// Vertex
type Vertex struct {
	Name     string
	Order    int32
	Edges    []*Edge
	Adjacent map[*Vertex]*Edge
	minEdge  *Edge
}

// Edge struct for Naming, weight set etc
type Edge struct {
	Name   string
	Weight int64
	Order  int32
}

// FindMST gets simple weighted undirected graph as an input
// returns minimum spanning tree as an output
// Algorithm is based on Burovka's theorem for faster parallel computing
func FindMST(swug []*Vertex) []*Vertex {
	l := len(swug)
	half := l / 2
	treeOne := swug[:half]
	treeTwo := swug[half:]

	// process 1st part
	wg.Add(1)
	go func() {
		for _, v := range treeOne {
			v.findMinEdge()

			if len(v.Adjacent) > 0 { // check adjacency and try to find foreign
				minForeign := int64(math.MaxInt64)
				// search the min connecting edge between 2 trees
				for _, vTwo := range treeTwo {
					for _, vOne := range treeOne {
						if val, ok := vOne.Adjacent[vTwo]; ok && minForeign > val.Weight {
							minForeign = val.Weight
						}
					}
				}
			}
		}
		wg.Done()
	}()

	// process 2nd part - here we don't need to search the min connecting edge between 2 trees
	wg.Add(1)
	go func() {
		for _, v := range treeTwo {
			v.findMinEdge()
		}
		wg.Done()
	}()
	wg.Wait()

	return append(treeOne, treeTwo...)
}

// finds min edge for every Vertex
func (v *Vertex) findMinEdge() {
	min := int64(math.MaxInt64)
	for _, e := range v.Edges {
		if e.Weight < min {
			min = e.Weight
			v.minEdge = e
		}
	}
}
