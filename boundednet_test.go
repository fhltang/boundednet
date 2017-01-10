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

	type Cases struct {
		Left int
		Right int
		Expected bn.Network
	}

	cases := []Cases{
		{0, 1, bn.NewNetwork(0, 0)},
		{0, 2, bn.NewNetwork(0, 1)},
		{0, 3, bn.NewNetwork(0, 6)},
		{0, 4, bn.NewNetwork(0, 6)},
		{2, 4, bn.NewNetwork(1, 5)},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("L=%d, R=%s", tc.Left, tc.Right), func(t *testing.T) {
			n:= solver.LeastNetwork(tc.Left, tc.Right)
			if tc.Expected != n {
				t.Error("Expecting ", tc.Expected, " got ", n)
			}
		})
	}
}
