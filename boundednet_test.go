package boundednet_test

import (
	"fmt"
	bn "github.com/fhltang/boundednet"
	"reflect"
	"testing"
)

func TestEmptyNetworkValid(t *testing.T) {
	if !bn.EmptyNetwork().Valid() {
		t.Fail()
	}
}

func TestValid(t *testing.T) {
	validCases := []bn.Network{
		bn.Network{Left: 6, Right: 8},
		bn.Network{Left: 8, Right: 12},
		bn.Network{Left: 5, Right: 6},
	}
	for _, tc := range validCases {
		t.Run(fmt.Sprintf("Valid [%d, %d)", tc.Left, tc.Right),
			func(t *testing.T) {
				if !tc.Valid() {
					t.Fail()
				}
			})
	}

	notValidCases := []bn.Network{
		bn.Network{Left: 2, Right: 1},
		bn.Network{Left: 8, Right: 11},
	}
	for _, tc := range notValidCases {
		t.Run(fmt.Sprintf("Not valid [%d, %d)", tc.Left, tc.Right),
			func(t *testing.T) {
				if tc.Valid() {
					t.Fail()
				}
			})
	}
}

func TestToNonEmptyNetwork(t *testing.T) {
	type Case struct {
		Input    bn.Network
		Expected bn.NonEmptyNetwork
	}
	cases := []Case{
		{
			bn.Network{Left: 6, Right: 8},
			bn.NonEmptyNetwork{A: 3, K: 1},
		}, {
			bn.Network{Left: 8, Right: 12},
			bn.NonEmptyNetwork{A: 2, K: 2},
		}, {
			bn.Network{Left: 5, Right: 6},
			bn.NonEmptyNetwork{A: 5, K: 0},
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("[%d, %d)", tc.Input.Left, tc.Input.Right),
			func(t *testing.T) {
				if tc.Input.ToNonEmptyNetwork() != tc.Expected {
					t.Fail()
				}
			})
	}
}

func TestInit(t *testing.T) {
	type Case struct {
		Name     string
		Input    []bn.Network
		Expected []bn.Network
	}
	cases := []Case{
		{
			"SortedNoOverlap",
			[]bn.Network{
				bn.Network{0, 1},
				bn.Network{1, 2},
				bn.Network{32, 36},
				bn.Network{60, 64},
			},
			[]bn.Network{
				bn.Network{0, 1},
				bn.Network{1, 2},
				bn.Network{32, 36},
				bn.Network{60, 64},
			},
		},
		{
			"SortedWithFullOverlap",
			[]bn.Network{
				bn.Network{0, 1},
				bn.Network{0, 2},
				bn.Network{32, 36},
				bn.Network{60, 64},
			},
			[]bn.Network{
				bn.Network{0, 2},
				bn.Network{32, 36},
				bn.Network{60, 64},
			},
		},
		{
			"UnsortedNoOverlap",
			[]bn.Network{
				bn.Network{32, 36},
				bn.Network{60, 64},
				bn.Network{0, 1},
				bn.Network{1, 2},
			},
			[]bn.Network{
				bn.Network{0, 1},
				bn.Network{1, 2},
				bn.Network{32, 36},
				bn.Network{60, 64},
			},
		},
		{
			"UnsortedWithOverlap",
			[]bn.Network{
				bn.Network{32, 36},
				bn.Network{60, 64},
				bn.Network{35, 36},
				bn.Network{0, 1},
				bn.Network{1, 2},
			},
			[]bn.Network{
				bn.Network{0, 1},
				bn.Network{1, 2},
				bn.Network{32, 36},
				bn.Network{60, 64},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			solver := bn.BacktrackingSolver{}
			solver.Init(tc.Input, 1)
			if !reflect.DeepEqual(tc.Expected, solver.Input) {
				t.Error("Expected", tc.Expected,
					"got", solver.Input)
			}
		})
	}

}

func TestPrecomputeLeastNetwork(t *testing.T) {
	input := []bn.Network{
		bn.Network{0, 1},
		bn.Network{1, 2},
		bn.Network{32, 36},
		bn.Network{60, 64},
	}

	solver := bn.BacktrackingSolver{}
	solver.Init(input, 1)
	solver.PrecomputeLeastNetwork()

	type Case struct {
		Left     int
		Right    int
		Expected bn.Network
	}

	cases := []Case{
		{0, 1, bn.Network{0, 1}},
		{0, 2, bn.Network{0, 2}},
		{0, 3, bn.Network{0, 64}},
		{0, 4, bn.Network{0, 64}},
		{2, 4, bn.Network{32, 64}},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("L=%d, R=%d", tc.Left, tc.Right),
			func(t *testing.T) {
				n := solver.LeastNetwork(tc.Left, tc.Right)
				if tc.Expected != n {
					t.Error("Expecting", tc.Expected,
						"got", n)
				}
			})
	}
}

func TestComputeTable(t *testing.T) {
	input := []bn.Network{
		bn.Network{0, 1},
		bn.Network{1, 2},
		bn.Network{32, 36},
		bn.Network{60, 64},
	}

	solver := bn.BacktrackingSolver{}
	solver.Init(input, 4)
	solver.PrecomputeLeastNetwork()
	solver.ComputeTable()

	type Case struct {
		N               int
		M               int
		ExpectedMinSize int
	}

	cases := []Case{
		{4, 4, 10},
		{2, 1, 2},
		{3, 2, 6},
		{4, 2, 34},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("N=%d M=%d", tc.M, tc.N), func(t *testing.T) {
			r, c := tc.M-1, tc.N-1
			minSize := solver.Table[r][c].MinSize
			if tc.ExpectedMinSize != minSize {
				t.Error("Expecting MinSize ",
					tc.ExpectedMinSize, " got ", minSize)
			}
		})

	}
}

func TestBacktrack(t *testing.T) {
	input := []bn.Network{
		bn.Network{0, 1},
		bn.Network{1, 2},
		bn.Network{32, 36},
		bn.Network{60, 64},
	}

	solver := bn.BacktrackingSolver{}
	solver.Init(input, 4)
	solver.PrecomputeLeastNetwork()
	solver.ComputeTable()

	type Case struct {
		M        int
		N        int
		Expected []bn.Network
	}
	cases := []Case{
		{1, 4, []bn.Network{
			bn.Network{0, 64},
		}},
		{2, 4, []bn.Network{
			bn.Network{0, 2},
			bn.Network{32, 64},
		}},
		{3, 4, []bn.Network{
			bn.Network{0, 2},
			bn.Network{32, 36},
			bn.Network{60, 64},
		}},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("M=%d", tc.M), func(t *testing.T) {
			solution := solver.Backtrack(tc.M, tc.N)
			if !reflect.DeepEqual(tc.Expected, solution) {
				t.Error("Expected", tc.Expected,
					"got", solution)
			}
		})

	}
}
