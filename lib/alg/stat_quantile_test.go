package alg_test

import (
	"github.com/aarioai/golib/lib/alg"
	"testing"
)

func TestStandardDeviation(t *testing.T) {
	pairs := map[float64]float64{
		0.95: 1.96,
	}
	for confident, deviation := range pairs {
		devi := alg.StandardQuantile(confident)
		if devi != deviation {
			t.Errorf("confident %f 's deviation supposed to be %f, but get %f", confident, deviation, devi)
		}
	}
}
