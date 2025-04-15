package crypto_test

import (
	"github.com/aarioai/golib/lib/code/crypto"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"testing"
)

func TestShuffleEncryptUSCC(t *testing.T) {
	pairs := []string{
		"9144030071526726XG",
		"11440300MB2C6448XB",
		"31440000X192794198",
		"A1440000X192794198",
		"Y1440000X192794198",
		"N1440000X192794198",
	}
	key := []byte(stdfmt.ReadableAsciiCodes)

	for _, uscc := range pairs {
		usccType, distri, ciphertext, err := crypto.ShuffleEncryptUSCC(uscc, key)
		if err != nil {
			t.Errorf("ShuffleEncryptUSCC(%s) -> %d, %d, %s %v", uscc, distri, usccType, ciphertext, err)
			return
		}
		got, err := ciphertext.Decrypt(usccType, distri, key)
		if err != nil || got != uscc {
			t.Errorf("Decrypt got %s, want %s %v", got, uscc, err)
			return
		}
	}
}
