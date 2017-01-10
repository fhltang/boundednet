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
	this.leastNetwork = make([][]Network, 0, len(this.Input))
	for i := 0; i < len(this.Input); i++ {
		this.leastNetwork = append(this.leastNetwork, make([]Network, i + 1))
		// TODO: do the precomputation
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
