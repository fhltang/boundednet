package boundednet_test

import (
	"fmt"
	bn "github.com/fhltang/boundednet"
	"testing"
)

func TestNewNetwork(t *testing.T) {
	n := bn.NewNetwork(1, 4)
	t.Logf("[%d, %d]", n.Left, n.Right)
	t.Fail()
}

func TestPrecomputeLeastNetwork(t *testing.T) {
	input := []bn.Network{
		bn.NewNetwork(0, 0),
		bn.NewNetwork(1, 0),
		bn.NewNetwork(8, 2),
		bn.NewNetwork(15, 2),
	}

	solver := bn.BacktrackingSolver{}
	solver.Init(input, 1)
	solver.PrecomputeLeastNetwork()

	type Case struct {
		Left int
		Right int
		Expected bn.Network
	}

	cases := []Case{
		{0, 1, bn.NewNetwork(0, 0)},
		{0, 2, bn.NewNetwork(0, 1)},
		{0, 3, bn.NewNetwork(0, 6)},
		{0, 4, bn.NewNetwork(0, 6)},
		{2, 4, bn.NewNetwork(1, 5)},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("L=%d, R=%d", tc.Left, tc.Right), func(t *testing.T) {
			n:= solver.LeastNetwork(tc.Left, tc.Right)
			if tc.Expected != n {
				t.Error("Expecting ", tc.Expected, " got ", n)
			}
		})
	}
}

func TestComputeTable(t *testing.T) {
	input := []bn.Network{
		bn.NewNetwork(0, 0),
		bn.NewNetwork(1, 0),
		bn.NewNetwork(8, 2),
		bn.NewNetwork(15, 2),
	}

	solver := bn.BacktrackingSolver{}
	solver.Init(input, 4)
	solver.PrecomputeLeastNetwork()
	solver.ComputeTable()

	type Case struct {
		N int
		M int
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
			r, c := tc.M - 1, tc.N -1
			minSize := solver.Table[r][c].MinSize
			if tc.ExpectedMinSize != minSize {
				t.Error("Expecting MinSize ", tc.ExpectedMinSize, " got ", minSize)
			}
		})
		
	}
}
