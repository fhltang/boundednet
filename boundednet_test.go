package boundednet_test

import (
	"fmt"
	bn "github.com/fhltang/boundednet"
	"github.com/fhltang/boundednet/binary"
	"github.com/fhltang/boundednet/snoc"
	"reflect"
	"testing"
)

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

type Problem struct {
	Addresses     []bn.Interval
	FootprintSize int
}

func MakeProblem(x []bn.Network) *Problem {
	intervals := bn.IntervalSlice(x)
	return &Problem{Addresses: bn.Canonical(intervals), FootprintSize: bn.FootprintSize(intervals)}
}

func (this Problem) IsPresolution(m int, x []bn.Network) bool {
	if len(x) > m {
		return false
	}

	return bn.Subset(this.Addresses, bn.IntervalSlice(x))
}

func TestSolvers(t *testing.T) {
	type Solver struct {
		Name  string
		Solve bn.Solver
	}
	solvers := []Solver{
		{"snoc", snoc.Solve},
		{"binary", binary.Solve},
	}

	type Case struct {
		Input    []bn.Network
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
				[]bn.Network{
					bn.ParseNetwork("192.168.0.0/24"),
					bn.ParseNetwork("192.168.1.0/24"),
					bn.ParseNetwork("192.168.3.0/24"),
					bn.ParseNetwork("192.168.4.0/23"),
					bn.ParseNetwork("192.168.16.0/21"),
					bn.ParseNetwork("194.0.0.0/8"),
					bn.ParseNetwork("200.0.0.0/29"),
					bn.ParseNetwork("200.0.0.8/29"),
				},
				[]bn.Network{
					bn.ParseNetwork("192.168.0.0/24"),
					bn.ParseNetwork("192.168.1.0/24"),
					bn.ParseNetwork("192.168.3.0/24"),
					bn.ParseNetwork("192.168.4.0/23"),
					bn.ParseNetwork("192.168.16.0/21"),
					bn.ParseNetwork("194.0.0.0/8"),
					bn.ParseNetwork("200.0.0.0/29"),
					bn.ParseNetwork("200.0.0.8/29"),
					bn.ParseNetwork("200.0.0.13/32"),
				},
				[]bn.Network{
					bn.ParseNetwork("192.168.0.0/24"),
					bn.ParseNetwork("192.168.1.0/24"),
					bn.ParseNetwork("192.168.3.0/24"),
					bn.ParseNetwork("192.168.4.0/23"),
					bn.ParseNetwork("192.168.16.0/21"),
					bn.ParseNetwork("194.0.0.0/8"),
					bn.ParseNetwork("200.0.0.0/30"),
					bn.ParseNetwork("200.0.0.4/30"),
					bn.ParseNetwork("200.0.0.8/29"),
					bn.ParseNetwork("200.0.0.13/32"),
				},
				[]bn.Network{
					bn.ParseNetwork("192.168.0.0/24"),
					bn.ParseNetwork("192.168.1.0/24"),
					bn.ParseNetwork("192.168.3.0/24"),
					bn.ParseNetwork("192.168.4.0/23"),
					bn.ParseNetwork("192.168.16.0/21"),
					bn.ParseNetwork("194.0.0.0/8"),
					bn.ParseNetwork("200.0.0.0/30"),
					bn.ParseNetwork("200.0.0.4/30"),
					bn.ParseNetwork("200.0.0.9/32"),
					bn.ParseNetwork("200.0.0.11/32"),
					bn.ParseNetwork("200.0.0.13/32"),
				},
			},
		},
	}
	for _, solver := range solvers {
		for _, tc := range cases {
			for m := 1; m <= len(tc.Expected); m++ {
				t.Run(fmt.Sprintf("%s M=%d", solver.Name, m), func(t *testing.T) {
					referenceSolution := tc.Expected[m-1]
					problem := MakeProblem(tc.Input)
					// Sanity check that reference solution is a presolution.
					if !problem.IsPresolution(m, referenceSolution) {
						t.Error("You dumbass.  The reference solution isn't even a presolution.")
					}

					solution := solver.Solve(tc.Input, m)
					if !problem.IsPresolution(m, solution) {
						t.Error("Result is not a presolution.")
					}
					if bn.FootprintSize(bn.IntervalSlice(solution)) > bn.FootprintSize(bn.IntervalSlice(referenceSolution)) {
						t.Error("Solution", solution, "is not minimal")
					}
				})
			}

		}
	}
}
