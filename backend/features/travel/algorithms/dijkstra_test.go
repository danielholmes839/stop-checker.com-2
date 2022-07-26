package algorithms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type node struct {
	id string
}

func (n node) ID() string {
	return n.id
}

type edge struct {
	to     string
	weight int
}

type path struct {
	*node
	prev   *path
	weight int
}

func TestAlgorithm(t *testing.T) {
	/*
		A--1-->B---3
		|          |
		|          v
		|		   D
		|          ^
		|          |
		---2-->C---1
	*/

	nodes := map[string]*node{
		"A": {id: "A"},
		"B": {id: "B"},
		"C": {id: "C"},
		"D": {id: "D"},
	}

	edges := map[string][]edge{
		"A": {{to: "B", weight: 1}, {to: "C", weight: 2}},
		"B": {{to: "D", weight: 3}},
		"C": {{to: "D", weight: 1}},
	}

	config := &DijkstraConfig[*path]{
		Destination: "D",
		Initial: []*path{{
			node:   nodes["A"],
			prev:   nil,
			weight: 0,
		}},
		Expand: func(n *path) []*path {
			expanded := []*path{}
			for _, edge := range edges[n.id] {
				expanded = append(expanded, &path{
					node:   nodes[edge.to],
					prev:   n,
					weight: n.weight + edge.weight,
				})
			}
			return expanded
		},
		Compare: func(a, b *path) bool {
			return a.weight < b.weight
		},
	}

	solution, _ := Dijkstra(config)

	assert.Equal(t, 3, solution.Node.weight)
}
