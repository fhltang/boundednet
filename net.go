package boundednet

// For convenience, we allow 2^32 as an "address" since this allows us
// to express networks as closed-open intervals.
type Address uint64

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
	if this.Left == this.Right {
		return 0, 0, 0
	}
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

