package binary_test

import (
	bn "github.com/fhltang/boundednet"
	"github.com/fhltang/boundednet/binary"
	"reflect"
	"testing"
)

func TestBuildTree(t *testing.T) {
	input := []bn.Network{
		{0, 4},
		{14, 16},
		{128, 256},
	}

	solver := binary.Solver{}
	solver.Init(input, 3)
	tree := solver.BuildTree(0, len(input), 0)
	expected := &binary.Node{
		Left: &binary.Node{
			Left: &binary.Node{
				Left:         nil,
				Right:        nil,
				MinSize:      []int{0},
				Network:      bn.Network{0, 4},
				LeftSolution: []int{0},
			},
			Right: &binary.Node{
				Left:         nil,
				Right:        nil,
				MinSize:      []int{0},
				Network:      bn.Network{14, 16},
				LeftSolution: []int{0},
			},
			MinSize:      []int{0, 0},
			Network:      bn.Network{0, 16},
			LeftSolution: []int{0, 0},
		},
		Right: &binary.Node{
			Left:         nil,
			Right:        nil,
			MinSize:      []int{0, 0},
			Network:      bn.Network{128, 256},
			LeftSolution: []int{0, 0},
		},
		MinSize:      []int{0, 0, 0},
		Network:      bn.Network{0, 256},
		LeftSolution: []int{0, 0, 0},
	}

	if !reflect.DeepEqual(*expected, *tree) {
		t.Error("Expected", *expected, "found", *tree)
	}
}

func TestComputeMinSize(t *testing.T) {
	input := []bn.Network{
		{0, 4},
		{14, 16},
		{128, 256},
	}
	solver := binary.Solver{}
	solver.Init(input, 3)
	tree := solver.BuildTree(0, len(input), 0)
	solver.ComputeMinSize(tree)

	expected := &binary.Node{
		Left: &binary.Node{
			Left: &binary.Node{
				Left:         nil,
				Right:        nil,
				MinSize:      []int{4},
				Network:      bn.Network{0, 4},
				LeftSolution: []int{0},
			},
			Right: &binary.Node{
				Left:         nil,
				Right:        nil,
				MinSize:      []int{2},
				Network:      bn.Network{14, 16},
				LeftSolution: []int{0},
			},
			MinSize:      []int{16, 6},
			Network:      bn.Network{0, 16},
			LeftSolution: []int{0, 1},
		},
		Right: &binary.Node{
			Left:         nil,
			Right:        nil,
			MinSize:      []int{128, 128},
			Network:      bn.Network{128, 256},
			LeftSolution: []int{0, 0},
		},
		MinSize:      []int{256, 144, 134},
		Network:      bn.Network{0, 256},
		LeftSolution: []int{0, 1, 2},
	}

	if !reflect.DeepEqual(*expected, *tree) {
		t.Error("Expected", *expected, "found", *tree)
	}
}

func TestBacktrack(t *testing.T) {
	input := []bn.Network{
		{0, 4},
		{14, 16},
		{128, 256},
	}
	solver := binary.Solver{}
	solver.Init(input, 3)
	solver.Tree = solver.BuildTree(0, len(input), 0)
	solver.ComputeMinSize(solver.Tree)

	type Case struct {
		M        int
		Expected []bn.Network
	}

	cases := []Case{
		{1, []bn.Network{bn.Network{0, 256}}},
		{2, []bn.Network{bn.Network{0, 16}, bn.Network{128, 256}}},
		{3, []bn.Network{bn.Network{0, 4}, bn.Network{14, 16}, bn.Network{128, 256}}},
	}

	for _, tc := range cases {
		result := solver.Backtrack(solver.Tree, tc.M)
		if !reflect.DeepEqual(tc.Expected, result) {
			t.Error("Expected", tc.Expected, "found", result)
		}
	}
}
