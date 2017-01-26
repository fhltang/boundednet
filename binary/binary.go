// A binary-tree-recursive solution of the bounded netwrok problem.

package binary

import (
	"fmt"
	bn "github.com/fhltang/boundednet"
	"sort"
)

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

type Node struct {
	// Left and right child nodes.
	Left, Right *Node

	// Size of minimal solution with up to M networks.  Index 0 is
	// size of minimal solution for M=1.
	MinSize []int

	// Minimal solution for M=1
	Network bn.Network

	// Used to recover minimal solution.  Minimal solution with at
	// most M networks is attained by taking a minimal solution
	// with at most `LeftSolution[M-1]` networks from the left
	// subnet together with a minimal solution with at most `M -
	// LeftSolution[M-1]` networks from the right subnet.
	LeftSolution []int
}

func (this *Node) String() string {
	if this == nil {
		return "NIL"
	}
	return fmt.Sprintf("%v", *this)
}

type Solver struct {
	// Inputs.
	Input []bn.Network
	M     int

	// Binary tree
	Tree *Node
}

func (this *Solver) Solve(input []bn.Network, m int) []bn.Network {
	this.Init(input, m)
	this.Tree = this.BuildTree(0, len(this.Input), 0)
	this.ComputeMinSize(this.Tree)

	return this.Backtrack(this.Tree, this.M)
}

func (this *Solver) Init(input []bn.Network, m int) {
	this.Input = bn.NormaliseInput(input)
	this.M = m
}

func (this *Solver) BuildTree(i, j, depth int) *Node {
	if i == j {
		return nil
	}

	node := &Node{
		MinSize:      make([]int, max(0, this.M-depth)),
		LeftSolution: make([]int, max(0, this.M-depth)),
	}

	node.Network = bn.LeastNetwork(this.Input, i, j)

	if j-i > 1 {
		midAddr := (node.Network.Left + node.Network.Right) / 2
		midIdx := i + sort.Search(j-i, func(k int) bool {
			return midAddr <= this.Input[i+k].Left
		})
		node.Left = this.BuildTree(i, midIdx, depth+1)
		node.Right = this.BuildTree(midIdx, j, depth+1)
	}

	return node
}

func (this *Solver) ComputeMinSize(node *Node) {
	if node.Left != nil && len(node.MinSize) > 1 {
		this.ComputeMinSize(node.Left)
	}
	if node.Right != nil && len(node.MinSize) > 1 {
		this.ComputeMinSize(node.Right)
	}

	node.MinSize[0] = node.Network.Size()

	for m := 1; m < len(node.MinSize); m++ {
		node.MinSize[m], node.LeftSolution[m] = node.MinSize[m-1], 0
		if node.Left == nil || node.Right == nil {
			continue
		}
		for i := 1; i <= m; i++ {
			if (i-1) < len(node.Left.MinSize) && (m-i) < len(node.Right.MinSize) {
				presolutionSize := node.Left.MinSize[i-1] + node.Right.MinSize[m-i]
				if presolutionSize < node.MinSize[m] {
					node.MinSize[m], node.LeftSolution[m] = presolutionSize, i
				}
			}
		}
	}
}

func (this *Solver) Backtrack(node *Node, m int) []bn.Network {
	result := make([]bn.Network, 0, m)
	return this.backtrack(result, node, m)
}

func (this *Solver) backtrack(dest []bn.Network, node *Node, m int) []bn.Network {
	if m == 1 {
		return append(dest, node.Network)
	}

	dest = this.backtrack(dest, node.Left, node.LeftSolution[m-1])
	return this.backtrack(dest, node.Right, m-node.LeftSolution[m-1])
}

func Solve(input []bn.Network, m int) []bn.Network {
	solver := Solver{}
	return solver.Solve(input, m)
}
