package boundednet

import (
	"sort"
)

type Address uint32

type Network struct {
	// close-open interval
	Left Address
	Right Address
}

func EmptyNetwork() Network {
	return Network{}
}

func (this Network) Size() int {
	return int(this.Right - this.Left)
}

func (this Network) Valid() bool {
	if this.Right < this.Left {
		return false
	}

	if this.Size() == 0 {
		return true
	}

	// Try to find a, k such that this == [a*2^k, (a+1)*2^k)
	a := this.Left
	a1 := this.Right
	for a != a1 && a % 2 == 0 && a1 % 2 == 0 {
		a = a >> 1
		a1 = a1 >> 1
	}
	return a + 1 == a1
}

type NonEmptyNetwork struct {
	// Rerpresents network [A*2^K, (A+1)*2^K)
	A int
	K uint
}

func (this Network) ToNonEmptyNetwork() NonEmptyNetwork {
	if !this.Valid() || this.Size() == 0 {
		panic("Not a valid non-empty network")
	}

	// Try to find a, k such that this == [a*2^k, (a+1)*2^k)
	k := uint(0)
	a := int(this.Left)
	a1 := int(this.Right)
	for a != a1 && a % 2 == 0 && a1 % 2 == 0 {
		k++
		a = a >> 1
		a1 = a1 >> 1
	}
	return NonEmptyNetwork{A: a, K: k}
}

func (this NonEmptyNetwork) ToNetwork() Network {
	return Network{Left: Address(this.A * (1 << this.K)), Right: Address((this.A + 1) * (1 << this.K))}
}

type Solver interface {
	Solve([]Network, int) []Network
}

type TableCell struct {
	MinSize int
	NextRow int
	NextCol int
	Network Network
}

type BacktrackingSolver struct {
	Input []Network
	M int
	
	leastNetwork [][]Network
	Table [][]TableCell
}

func (this *BacktrackingSolver) Solve(input []Network, m int) []Network {
	this.Init(input, m)

	this.PrecomputeLeastNetwork()

	this.ComputeTable()

	return this.Backtrack(this.M, len(this.Input))
}

type ByLeftWidth []Network
func (this ByLeftWidth) Len() int { return len(this) }
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
	for j := 0; j < len(this.Input) - 1; j++ {
		if this.Input[j].Right <= this.Input[j+1].Left {
			if  i < j {
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
	this.leastNetwork = make([][]Network, 0, len(this.Input) + 1)
	for j := 0; j <= len(this.Input); j++ {
		this.leastNetwork = append(this.leastNetwork, make([]Network, j + 1))
		for i := 0; i <= j; i++ {
			if i == j {
				this.leastNetwork[j][i] = EmptyNetwork()
			} else {
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
				this.leastNetwork[j][i] = left.ToNetwork()
			}
		}
	}
}

func (this *BacktrackingSolver) ComputeTable() {
	this.Table = make([][]TableCell, this.M)
	for m := 0; m < this.M; m++ {
		this.Table[m] = make([]TableCell, len(this.Input))
		for k := 0; k < len(this.Input); k++ {
			if m == 0 {
				network := this.LeastNetwork(0, k+1)
				this.Table[m][k] = TableCell{
					MinSize: network.Size(),
					Network: network,
				}
			} else {
				this.Table[m][k].MinSize = 1<<32
				for n := 0; n <= k; n++ {
					network:= this.LeastNetwork(n+1, k+1)
					presolutionSize := network.Size() + this.Table[m-1][n].MinSize
					if presolutionSize < this.Table[m][k].MinSize {
						this.Table[m][k] = TableCell{
							MinSize: presolutionSize,
							Network: network,
							NextRow: m-1,
							NextCol: n,
						}
					}
				}
			}
		}
	}
}

func (this *BacktrackingSolver) Backtrack(m, n int) []Network {
	output := make([]Network, 0, m)
	row, col := m - 1, n - 1
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
	left, right := 0, len(output) - 1
	for left < right {
		output[left], output[right] = output[right], output[left]
		left, right = left + 1, right - 1
	}
	return output
}
