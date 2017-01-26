package boundednet

import (
	"sort"
)

type Solver func([]Network, int) []Network

type ByLeftWidth []Network
func (this ByLeftWidth) Len() int      { return len(this) }
func (this ByLeftWidth) Swap(i, j int) { this[i], this[j] = this[j], this[i] }
func (this ByLeftWidth) Less(i, j int) bool {
	if this[i].Left == this[j].Left {
		return this[i].Right > this[j].Right
	}
	return this[i].Left < this[j].Left
}

func NormaliseInput(input []Network) []Network {
	result := make([]Network, len(input))
	copy(result, input)

	// Sort input.
	sort.Sort(ByLeftWidth(result))

	// Remove overlapping networks.
	i := 0
	for j := 0; j < len(result)-1; j++ {
		if result[j].Right <= result[j+1].Left {
			if i < j {
				result[i+1] = result[j+1]
			}
			i++
		}
	}
	result = result[:i+1]
	return result
}

func LeastNetwork(networks []Network, i, j int) Network {
	if i == j {
		return EmptyNetwork()
	}

	left := networks[i].ToNonEmptyNetwork()
	right := networks[j-1].ToNonEmptyNetwork()
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

