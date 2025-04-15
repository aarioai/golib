package stdfmt_test

import (
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"testing"
)

func TestBase64(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := coding.RandASCIICode(i + 1)
		urlSafe := false
		withoutPadding := false
		for j := 0; j < 4; j++ {
			if j&1 == 1 {
				urlSafe = true
			}
			if j>>1 == 1 {
				withoutPadding = true
			}
			got := stdfmt.EncodeBase64(s, urlSafe, withoutPadding)
			decodeGot, err := stdfmt.DecodeBase64(got)
			if err != nil || s != string(decodeGot) {
				t.Errorf("decode failed, stdRaw: %s, got: %s, err: %v", got, decodeGot, err)
			}
		}
	}
}
