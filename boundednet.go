package boundednet

type Address uint32

type Network struct {
	A int
	K uint

	// close-open interval
	Left Address
	Right Address
}

func NewNetwork(a int, k uint) Network {
	return Network{
		A: a,
		K: k,
		Left: Address(a << k),
		Right: Address((a+1) << k),
	}
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

	return this.Backtrack()
}

func (this *BacktrackingSolver) Init(input []Network, m int) {
	this.Input = input
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
				this.leastNetwork[j][i] = Network{}
			} else {
				left, right := this.Input[i], this.Input[j-1]
				for left.A != right.A || left.K != right.K {
					if right.K < left.K {
						right = NewNetwork(int(right.A / 2), right.K + 1)
					} else if left.K < right.K {
						left = NewNetwork(int(left.A / 2), left.K + 1)
					} else if left.A < right.A {
						right = NewNetwork(int(right.A / 2), right.K + 1)
					} else if right.A < left.A {
						left = NewNetwork(int(left.A / 2), left.K + 1)
					}
				}
				this.leastNetwork[j][i] = left
			}
		}
	}
}

func (this *BacktrackingSolver) ComputeTable() {
	// TODO: compute table
}

func (this *BacktrackingSolver) Backtrack() []Network {
	output := make([]Network, this.M)
	// TODO: do backtracking
	return output
}
