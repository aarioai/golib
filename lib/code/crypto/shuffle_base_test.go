package crypto_test

import (
	"bytes"
	"errors"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/golib/lib/code/crypto"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"math"
	"math/rand"
	"testing"
)

func TestShuffleEncryptDigits(t *testing.T) {

	for i := 0; i < 1000; i++ {
		for j := 0; j < 2; j++ {
			scatter := j == 0
			var key []byte
			kn := rand.Intn(len(types.Base64Digits)) + 10
			if kn < 36 {
				key = []byte(types.Base64Digits[:kn+1])
			} else {
				key = []byte(stdfmt.ReadableAsciiCodes)
			}
			var n uint64
			if i == 1 {
				n = math.MaxUint64
			} else if i > 1 {
				n = rand.Uint64()
			}

			var wantErr error
			if n < crypto.ShuffleEncryptMin {
				wantErr = ae.ErrInputTooSmall
			}
			keyClone := bytes.Clone(key)
			ciphertext, err := crypto.ShuffleEncryptDigits(n, key, scatter)
			if !errors.Is(wantErr, err) {
				t.Fatalf("ShuffleEncryptDigits(%d, key base:%d,%v) = %v; want %v", n, len(key), scatter, err, wantErr)
			}
			if err != nil {
				continue // 上面判断过了
			}
			if !bytes.Equal(keyClone, key) {
				t.Fatalf("ShuffleEncryptDigits changed key")
			}
			got, err := ciphertext.Decrypt(key)
			if err != nil {
				t.Fatalf("failed decrypt(%s, %s): %v", ciphertext, string(key), err)
			}
			if got != n {
				t.Fatalf("failed decrypt(%s) got %d, want %d", ciphertext, got, n)
			}
			if !bytes.Equal(keyClone, key) {
				t.Fatalf("Decrypt changed key")
			}

		}
	}
}
