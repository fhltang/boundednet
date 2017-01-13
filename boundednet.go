package boundednet

import (
	"sort"
)

type Address uint32

// A network represented as a closed-open interval [Left, Right).
type Network struct {
	Left  Address
	Right Address
}

func EmptyNetwork() Network {
	return Network{}
}

func (this Network) Size() int {
	return int(this.Right - this.Left)
}

// Normalise network into form [x*2^k, y*2^k) for largest k.
func (this Network) Normalise() (int, int, uint) {
	k := uint(0)
	x := int(this.Left)
	y := int(this.Right)
	for x != y && x%2 == 0 && y%2 == 0 {
		k++
		x = x >> 1
		y = y >> 1
	}
	return x, y, k
}

// Determine if a network is valid.
func (this Network) Valid() bool {
	if this.Right < this.Left {
		return false
	}

	if this.Size() == 0 {
		return true
	}

	// Try to find a, k such that this == [a*2^k, (a+1)*2^k)
	a, a1, _ := this.Normalise()
	return a+1 == a1
}

// Alternative representation of network [A*2^K, (A+1)*2^K)
type NonEmptyNetwork struct {
	A int
	K uint
}

func (this Network) ToNonEmptyNetwork() NonEmptyNetwork {
	if !this.Valid() || this.Size() == 0 {
		panic("Not a valid non-empty network")
	}

	// Find a, k such that this == [a*2^k, (a+1)*2^k)
	a, _, k := this.Normalise()
	return NonEmptyNetwork{A: a, K: k}
}

func (this NonEmptyNetwork) ToNetwork() Network {
	return Network{
		Left:  Address(this.A * (1 << this.K)),
		Right: Address((this.A + 1) * (1 << this.K)),
	}
}

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
	output := make([]Network, 0, m)
	row, col := m-1, n-1
	for {
		cell := this.Table[row][col]
		if cell.Network.Size() > 0 {
			output = append(output, cell.Network)
		}
		if row == 0 {
			break
		}
		row, col = cell.NextRow, cell.NextCol
	}
	// reverse output
	left, right := 0, len(output)-1
	for left < right {
		output[left], output[right] = output[right], output[left]
		left, right = left+1, right-1
	}
	return output
}
