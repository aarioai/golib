package crypto_test

import (
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/lib/code/crypto"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"math/rand/v2"
	"testing"
)

func TestShuffleEncryptSid(t *testing.T) {
	key := coding.Shuffle([]byte(stdfmt.Base62Digits + "_-"))
	for i := 0; i < 1000; i++ {
		for j := 0; j < 2; j++ {
			scatter := j == 0
			n := rand.IntN(20) + 4
			sid := coding.Rand(n, key)
			ciphertext, err := crypto.ShuffleEncryptSid(sid, key, scatter)
			if err != nil {
				t.Fatal(err)
			}
			got, err := ciphertext.Decrypt(key)
			if err != nil {
				t.Fatal(err)
			}
			if got != sid {
				t.Fatalf("got %s, want %s", got, sid)
			}
		}
	}
}
