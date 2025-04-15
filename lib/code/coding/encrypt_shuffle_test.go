package coding_test

import (
	"bytes"
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"math/rand/v2"
	"testing"
)

func TestScatter(t *testing.T) {
	for i := 0; i < 10; i++ {
		key := []byte(coding.RandASCIICode(8 * i))
		seed := rand.IntN(128)
		got := coding.Scatter(key, seed)
		if len(key) != len(got) {
			t.Fatalf("failed Scatter(%s, %d): length mismatch, want %d, got %s(len:%d)", string(key), seed, len(key), string(got), len(got))
		}

		gotNew := coding.Unscatter(got, seed)
		if !bytes.Equal(key, gotNew) {
			t.Fatalf("failed Unscatter(%s, %d), got %s, want %s", string(got), seed, string(gotNew), string(key))
		}
	}

}
func TestObfuscateBytes(t *testing.T) {
	key := []byte(stdfmt.ReadableAsciiCodes)
	for i := 0; i < 10; i++ {
		shift := rand.IntN(16)
		// 需要混淆的字符串，必须要在key里面
		obStr := []byte(coding.Rand(rand.IntN(128), key))
		want := string(obStr)
		err := coding.ShuffleEncrypt(obStr, shift, key)
		if err != nil {
			t.Errorf("ShuffleEncrypt(%s, %d,%s) error  %v", want, shift, string(key), err)
		}
		deStr := make([]byte, len(obStr))
		copy(deStr, obStr)
		err = coding.ShuffleDecrypt(deStr, shift, key)
		if err != nil {
			t.Errorf("ShuffleDecrypt(%s, %d,%s) error  %v", string(deStr), shift, string(key), err)
		}
		if string(deStr) != want {
			t.Errorf("ShuffleDecrypt(%s, %d,%s) got %s, want %s", string(deStr), shift, string(key), string(deStr), want)

		}
	}
}
