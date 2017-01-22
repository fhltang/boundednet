package boundednet_test

import (
	"fmt"
	bn "github.com/fhltang/boundednet"
	"github.com/fhltang/boundednet/snoc"
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

func TestSnocSolver(t *testing.T) {
	input := []bn.Network{
		bn.Network{0, 1},
		bn.Network{1, 2},
		bn.Network{32, 36},
		bn.Network{60, 64},
	}

	solvers := []bn.Solver{
		snoc.Solve,
	}

	type Case struct {
		M        int
		Expected []bn.Network
	}
	cases := []Case{
		{1, []bn.Network{
			bn.Network{0, 64},
		}},
		{2, []bn.Network{
			bn.Network{0, 2},
			bn.Network{32, 64},
		}},
		{3, []bn.Network{
			bn.Network{0, 2},
			bn.Network{32, 36},
			bn.Network{60, 64},
		}},
	}
	for _, solve := range solvers {
		for _, tc := range cases {
			t.Run(fmt.Sprintf("M=%d", tc.M), func(t *testing.T) {
				solution := solve(input, tc.M)
				if !reflect.DeepEqual(tc.Expected, solution) {
					t.Error("Expected", tc.Expected,
						"got", solution)
				}
			})

		}
	}
}
