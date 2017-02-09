package boundednet_test

import (
	"fmt"
	bn "github.com/fhltang/boundednet"
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
