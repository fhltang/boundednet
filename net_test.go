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

func TestParseNetwork(t *testing.T) {
	left := uint64(192)<<24 + uint64(168)<<16 + uint64(1)<<8
	right := uint64(192)<<24 + uint64(168)<<16 + uint64(2)<<8
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

func TestCanonical(t *testing.T) {
	type Case struct {
		Name string
		Input []bn.Interval
		Expected []bn.Interval
	}

	cases := []Case{
		{"empty", []bn.Interval{}, []bn.Interval{}},
		{"singleton", []bn.Interval{{1, 2}}, []bn.Interval{{1, 2}}},
		{"reverse", []bn.Interval{{4, 5}, {1, 2}}, []bn.Interval{{1, 2}, {4, 5}}},
		{"adjacent", []bn.Interval{{1, 2}, {2, 3}}, []bn.Interval{{1, 3}}},
		{"adjacent_reverse", []bn.Interval{{2, 3}, {1, 2}}, []bn.Interval{{1, 3}}},
		{"contains", []bn.Interval{{1, 3}, {1, 2}}, []bn.Interval{{1, 3}}},
		{"contains_reverse", []bn.Interval{{1, 2}, {1, 3}}, []bn.Interval{{1, 3}}},
		{"overlaps", []bn.Interval{{1, 3}, {2, 4}}, []bn.Interval{{1, 4}}},
		{"overlaps_reverse", []bn.Interval{{2, 4}, {1, 3}}, []bn.Interval{{1, 4}}},
		{"three", []bn.Interval{{2, 4}, {1, 3}, {5, 8}}, []bn.Interval{{1, 4}, {5, 8}}},
		{"overlaps_overlaps",
			[]bn.Interval{{2, 4}, {1, 3}, {6, 8}, {5, 7}},
			[]bn.Interval{{1, 4}, {5, 8}}},
		{"overlaps_merge",
			[]bn.Interval{{2, 4}, {1, 3}, {6, 8}, {4, 7}},
			[]bn.Interval{{1, 8}}},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if !reflect.DeepEqual(tc.Expected, bn.Canonical(tc.Input)) {
				t.Fail()
			}
		})
	}
}
