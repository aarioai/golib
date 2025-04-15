package alg_test

import (
	"github.com/aarioai/golib/lib/alg"
	"testing"
)

func TestWilsonRank(t *testing.T) {
	t.Log(alg.LowerBound(0, 0, 0.95))
	t.Log(alg.LowerBound(0, 1, 0.95))
	t.Log(alg.LowerBound(1, 1, 0.95))
	t.Log(alg.LowerBound(20, 50, 0.95))
	t.Log(alg.LowerBound(10, 12, 0.95))
	t.Log(alg.LowerBound(30, 50, 0.95))
	t.Log(alg.LowerBound(30000000000000, 50000000000000, 0.95))
}
