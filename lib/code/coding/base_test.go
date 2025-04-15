package coding_test

import (
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"testing"
)

func TestShuffle(t *testing.T) {
	nums := coding.CloneNumberRunes()
	lows := coding.CloneLowercaseRunes()
	uppers := coding.CloneUppercaseRunes()
	numlowers := coding.CloneNumLowerRunes()
	alphas := coding.CloneAlphabetRunes()
	if len(nums) != len([]rune(stdfmt.Numbers)) {
		t.Errorf("numberRunes changed unexpectedly")
	}
	if len(lows) != len([]rune(stdfmt.Lowercases)) {
		t.Errorf("lowercaseRunes changed unexpectedly")
	}
	if len(uppers) != len([]rune(stdfmt.Uppercases)) {
		t.Errorf("uppercaseRunes changed unexpectedly")
	}
	if len(numlowers) != len([]rune(stdfmt.Base36Digits)) {
		t.Errorf("numLowerRunes changed unexpectedly")
	}
	if len(alphas) != len([]rune(stdfmt.Base62Digits)) {
		t.Errorf("alphabetRunes changed unexpectedly")
	}
}
