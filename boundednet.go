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
