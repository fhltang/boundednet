package boundednet

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// For convenience, we allow 2^32 as an "address" since this allows us
// to express networks as closed-open intervals.
type Address uint64

// A closed-open interval [Left, Right).
type Interval struct {
	Left  Address
	Right Address
}

// A network represented as a closed-open interval [Left, Right).
type Network Interval

func EmptyNetwork() Network {
	return Network{}
}

func ParseNetwork(netmask string) Network {
	addrAndMask := strings.Split(netmask, "/")
	if len(addrAndMask) != 2 {
		panic(fmt.Sprintf("Cannot parse input %s", netmask))
	}
	addr, mask := addrAndMask[0], addrAndMask[1]

	ones, err := strconv.ParseUint(mask, 10, 0)
	if err != nil {
		panic(fmt.Sprintf("Cannot parse number %s", mask))
	}

	bytes := strings.Split(addr, ".")
	if len(bytes) != 4 {
		panic(fmt.Sprintf("Cannot parse address %s", bytes))
	}
	var left, right uint64
	for _, b := range bytes {
		left = left << 8
		byte, err := strconv.ParseUint(b, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Cannot parse byte value %s", b))
		}
		left = left + byte
	}
	bitMask := uint64(((1 << ones) - 1) << (32 - ones))
	left = left & bitMask
	right = left + (uint64(1) << (32 - ones))
	return Network{Address(left), Address(right)}
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

// Canonicalise a list of Interval.
//
// Given a list of Interval, return the shortest list of intervals in
// increasing order whose union is the same as that of the input list.
type CanonicalOrder []Interval

func (this CanonicalOrder) Len() int      { return len(this) }
func (this CanonicalOrder) Swap(i, j int) { this[i], this[j] = this[j], this[i] }
func (this CanonicalOrder) Less(i, j int) bool {
	if this[i].Left == this[j].Left {
		return this[i].Right > this[j].Right
	}
	return this[i].Left < this[j].Left
}

func Canonical(input []Interval) []Interval {
	if len(input) == 0 {
		return []Interval{}
	}

	result := make([]Interval, len(input))
	copy(result, input)
	sort.Sort(CanonicalOrder(result))
	i := 0
	for j, next := range result[1:] {
		prev := result[i]
		if next.Left > prev.Right {
			i++
			result[i], prev = result[j+1], next
		} else if next.Right > prev.Right {
			result[i].Right = next.Right
		}
	}
	result = result[:i+1]
	return result
}

// Determine if a list of `Interval`s is a subset of another.
func Subset(x, y []Interval) bool {
	cy := Canonical(y)
	for _, intvl := range x {
		i := sort.Search(len(cy), func(j int) bool {
			return cy[j].Left > intvl.Left
		})
		if i == 0 {
			return false
		}
		if intvl.Right > cy[i-1].Right {
			return false
		}
	}
	return true
}

// Determine the size of the union of a list of `Interval`s.
func FootprintSize(input []Interval) int {
	size := 0
	for _, intvl := range Canonical(input) {
		size = size + int(intvl.Right-intvl.Left)
	}
	return size
}
