package boundednet_test

import (
	"fmt"
	bn "github.com/fhltang/boundednet"
	"github.com/fhltang/boundednet/binary"
	"github.com/fhltang/boundednet/snoc"
	"reflect"
	"testing"
)

func TestEmptyNetworkValid(t *testing.T) {
	if !bn.EmptyNetwork().Valid() {
		t.Fail()
	}
}

func TestParseNetwork(t *testing.T) {
	left := uint64(192) << 24 + uint64(168) << 16 + uint64(1) << 8
	right := uint64(192) << 24 + uint64(168) << 16 + uint64(2) << 8
	expected := bn.Network{bn.Address(left), bn.Address(right)}

	net := bn.ParseNetwork("192.168.1.0/24")
	if expected != net {
		t.Error("expecting", expected, "got", net)
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

func TestNormaliseInput(t *testing.T) {
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
			input := bn.NormaliseInput(tc.Input)
			if !reflect.DeepEqual(tc.Expected, input) {
				t.Error("Expected", tc.Expected,
					"got", input)
			}
		})
	}

}

func TestSolvers(t *testing.T) {
	type Solver struct {
		Name string
		Solve bn.Solver
	}
	solvers := []Solver{
		{"snoc", snoc.Solve},
		{"binary", binary.Solve},
	}

	type Case struct {
		Input []bn.Network
		Expected [][]bn.Network
	}
	cases := []Case{
		{
			[]bn.Network{
				bn.Network{0, 1},
				bn.Network{1, 2},
				bn.Network{32, 36},
				bn.Network{60, 64},
			},
			[][]bn.Network{
				[]bn.Network{
					bn.Network{0, 64},
				},
				[]bn.Network{
					bn.Network{0, 2},
					bn.Network{32, 64},
				},
				[]bn.Network{
					bn.Network{0, 2},
					bn.Network{32, 36},
					bn.Network{60, 64},
				},
			},
		},
		{
			[]bn.Network{
				bn.Network{100 << 24, 101 << 24},
				bn.Network{200 << 24, 201 << 24},
			},
			[][]bn.Network{
				[]bn.Network{
					bn.Network{0, 1 << 32},
				},
				[]bn.Network{
					bn.Network{100 << 24, 101 << 24},
					bn.Network{200 << 24, 201 << 24},
				},
			},
		},
		{
			[]bn.Network{
				bn.ParseNetwork("192.168.0.0/24"),
				bn.ParseNetwork("192.168.1.0/24"),
				bn.ParseNetwork("192.168.3.0/24"),
				bn.ParseNetwork("192.168.4.0/23"),
				bn.ParseNetwork("192.168.16.0/21"),
				bn.ParseNetwork("194.0.0.0/8"),
				bn.ParseNetwork("200.0.0.11/32"),
				bn.ParseNetwork("200.0.0.1/32"),
				bn.ParseNetwork("200.0.0.13/32"),
				bn.ParseNetwork("200.0.0.3/32"),
				bn.ParseNetwork("200.0.0.5/32"),
				bn.ParseNetwork("200.0.0.7/32"),
				bn.ParseNetwork("200.0.0.9/32"),
			},
			[][]bn.Network{
				[]bn.Network{
					bn.ParseNetwork("192.0.0.0/4"),
				},
				[]bn.Network{
					bn.ParseNetwork("192.0.0.0/6"),
					bn.ParseNetwork("200.0.0.0/28"),
				},
				[]bn.Network{
					bn.ParseNetwork("192.168.0.0/19"),
					bn.ParseNetwork("194.0.0.0/8"),
					bn.ParseNetwork("200.0.0.0/28"),
				},
				[]bn.Network{
					bn.ParseNetwork("192.168.0.0/21"),
					bn.ParseNetwork("192.168.16.0/21"),
					bn.ParseNetwork("194.0.0.0/8"),
					bn.ParseNetwork("200.0.0.0/28"),
				},
				[]bn.Network{
					bn.ParseNetwork("192.168.0.0/22"),
					bn.ParseNetwork("192.168.4.0/23"),
					bn.ParseNetwork("192.168.16.0/21"),
					bn.ParseNetwork("194.0.0.0/8"),
					bn.ParseNetwork("200.0.0.0/28"),
				},
				[]bn.Network{
					bn.ParseNetwork("192.168.0.0/23"),
					bn.ParseNetwork("192.168.3.0/24"),
					bn.ParseNetwork("192.168.4.0/23"),
					bn.ParseNetwork("192.168.16.0/21"),
					bn.ParseNetwork("194.0.0.0/8"),
					bn.ParseNetwork("200.0.0.0/28"),
				},
				[]bn.Network{
					bn.ParseNetwork("192.168.0.0/24"),
					bn.ParseNetwork("192.168.1.0/24"),
					bn.ParseNetwork("192.168.3.0/24"),
					bn.ParseNetwork("192.168.4.0/23"),
					bn.ParseNetwork("192.168.16.0/21"),
					bn.ParseNetwork("194.0.0.0/8"),
					bn.ParseNetwork("200.0.0.0/28"),
				},
			},
		},
	}
	for _, solver := range solvers {
		for _, tc := range cases {
			for m := 1; m <= len(tc.Expected); m++ {
				t.Run(fmt.Sprintf("%s M=%d", solver.Name, m), func(t *testing.T) {
					solution := solver.Solve(tc.Input, m)
					if !reflect.DeepEqual(tc.Expected[m-1], solution) {
						t.Error("Expected", tc.Expected[m-1],
							"got", solution)
					}
				})
			}

		}
	}
}
