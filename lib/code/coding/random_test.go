package coding_test

import (
	"bytes"
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"math/rand/v2"
	"testing"
)

func TestRandASCIICode(t *testing.T) {
	size := rand.IntN(128)
	got := coding.RandASCIICode(size)
	if size != len(got) {
		t.Errorf("RandASCIICode wrong length: got %d, want %d", len(got), size)
		return
	}
}

func TestRandAlphabets(t *testing.T) {
	var s string
	var num string
	var n int

	charSet := []byte(stdfmt.Base62Digits)

	for i := 0; i < 100; i++ {
		n = (i+4)%10 + 2

		charSetBefore := bytes.Clone(charSet)
		s = coding.Rand(n, charSet)
		if len(s) != n {
			t.Errorf("generate %d alphabests,bad %s", n, s)
			return
		}
		charSetAfter := bytes.Clone(charSet)
		if len(charSetBefore) != len(charSetAfter) {
			t.Errorf("Rand changed rune slice")
			return
		}

		alphabetBefore := coding.CloneShuffledAlphabets()
		s = coding.RandAlphabets(n, 0)
		if len(s) != n {
			t.Errorf("generate %d alphabests,bad %s", n, s)
			return
		}
		alphabetAfter := coding.CloneShuffledAlphabets()
		if len(alphabetBefore) != len(alphabetAfter) {
			t.Errorf("RandAlphabets changed CloneShuffledAlphabets rune slice")
			return
		}

		lowsBefore := coding.CloneShuffledNumLowers()
		s = coding.RandNumLowers(n)
		if len(s) != n {
			t.Errorf("generate %d number-lower,bad %s", n, s)
			return
		}
		lowsAfter := coding.CloneShuffledNumLowers()
		if len(lowsBefore) != len(lowsAfter) {
			t.Errorf("RandNumLowers changed CloneShuffledNumLowers rune slice")
			return
		}

		numsBefore := coding.CloneShuffledNums()
		num = coding.RandNum(n, 0)
		if len(num) != n {
			t.Errorf("generate %d num,bad %s", n, num)
			return
		}
		numsAfter := coding.CloneShuffledNums()
		if len(numsBefore) != len(numsAfter) {
			t.Errorf("CloneShuffledNums changed CloneShuffledNums rune slice")
			return
		}
	}
}
