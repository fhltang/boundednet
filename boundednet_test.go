package boundednet_test

import (
	"github.com/fhltang/boundednet"
	"testing"
)

func TestNewNetwork(t *testing.T) {
	n := boundednet.NewNetwork(1, 4)
	t.Logf("[%d, %d]", n.Left, n.Right)
	t.Fail()
}

