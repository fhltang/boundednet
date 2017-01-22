package boundednet

import (
	"sort"
)

type Solver interface {
	Solve([]Network, int) []Network
}

type TableCell struct {
	// Size of a minimal solution.
	MinSize int

	// A minimal solution is obtained by combining Network with a
	// subsolution (NextRow, NextCol).
	//
	// Used by backtracking algorithm.
	NextRow int
	NextCol int
	Network Network
}

type BacktrackingSolver struct {
	// Inputs.
	Input []Network
	M     int

	// Precomputed values of LeastNetwork().
	leastNetwork [][]Network

	// Table used in dynamic programming solution.
	Table [][]TableCell
}

func (this *BacktrackingSolver) Solve(input []Network, m int) []Network {
	this.Init(input, m)
	this.PrecomputeLeastNetwork()
	this.ComputeTable()
	return this.Backtrack(this.M, len(this.Input))
}

type ByLeftWidth []Network

func (this ByLeftWidth) Len() int      { return len(this) }
func (this ByLeftWidth) Swap(i, j int) { this[i], this[j] = this[j], this[i] }
func (this ByLeftWidth) Less(i, j int) bool {
	if this[i].Left == this[j].Left {
		return this[i].Right > this[j].Right
	}
	return this[i].Left < this[j].Left
}

func (this *BacktrackingSolver) Init(input []Network, m int) {
	this.Input = make([]Network, len(input))
	copy(this.Input, input)

	// Sort input.
	sort.Sort(ByLeftWidth(this.Input))

	// Remove overlapping networks.
	i := 0
	for j := 0; j < len(this.Input)-1; j++ {
		if this.Input[j].Right <= this.Input[j+1].Left {
			if i < j {
				this.Input[i+1] = this.Input[j+1]
			}
			i++
		}
	}
	this.Input = this.Input[:i+1]

	this.M = m
}

func (this *BacktrackingSolver) LeastNetwork(i, j int) Network {
	return this.leastNetwork[j][i]
}

func (this *BacktrackingSolver) PrecomputeLeastNetwork() {
	this.leastNetwork = make([][]Network, 0, len(this.Input)+1)
	for j := 0; j <= len(this.Input); j++ {
		this.leastNetwork = append(
			this.leastNetwork, make([]Network, j+1))
		for i := 0; i <= j; i++ {
			this.leastNetwork[j][i] = this.computeLeastNetwork(i, j)
		}
	}
}

func (this BacktrackingSolver) computeLeastNetwork(i, j int) Network {
	if i == j {
		return EmptyNetwork()
	}

	left := this.Input[i].ToNonEmptyNetwork()
	right := this.Input[j-1].ToNonEmptyNetwork()
	for left != right {
		if right.K < left.K {
			right = NonEmptyNetwork{A: right.A >> 1, K: right.K + 1}
		} else if left.K < right.K {
			left = NonEmptyNetwork{A: left.A >> 1, K: left.K + 1}
		} else if left.A < right.A {
			right = NonEmptyNetwork{A: right.A >> 1, K: right.K + 1}
		} else if right.A < left.A {
			left = NonEmptyNetwork{A: left.A >> 1, K: left.K + 1}
		}
	}
	return left.ToNetwork()
}

func (this *BacktrackingSolver) ComputeTable() {
	this.Table = make([][]TableCell, this.M)
	for m := 0; m < this.M; m++ {
		this.Table[m] = make([]TableCell, len(this.Input))
		for k := 0; k < len(this.Input); k++ {
			this.Table[m][k] = this.computeCell(m, k)
		}
	}
}

func (this BacktrackingSolver) computeCell(m, k int) TableCell {
	if m == 0 {
		network := this.LeastNetwork(0, k+1)
		return TableCell{
			MinSize: network.Size(),
			Network: network,
		}
	}

	minimalSolution := TableCell{MinSize: 1 << 32} // max
	for n := 0; n <= k; n++ {
		network := this.LeastNetwork(n+1, k+1)
		presolutionSize := network.Size() + this.Table[m-1][n].MinSize
		if presolutionSize < minimalSolution.MinSize {
			minimalSolution = TableCell{
				MinSize: presolutionSize,
				Network: network,
				NextRow: m - 1,
				NextCol: n,
			}
		}
	}
	return minimalSolution
}

func (this *BacktrackingSolver) Backtrack(m, n int) []Network {
	output := make([]Network, m)
	head := m
	row, col := m-1, n-1
	for {
		cell := this.Table[row][col]
		if cell.Network.Size() > 0 {
			head--
			output[head] = cell.Network
		}
		if row == 0 {
			break
		}
		row, col = cell.NextRow, cell.NextCol
	}
	return output[head:]
}
