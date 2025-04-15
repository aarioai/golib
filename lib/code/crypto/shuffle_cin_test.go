package crypto_test

import (
	"github.com/aarioai/golib/lib/code/crypto"
	"testing"
)

func TestEncodeCin(t *testing.T) {
	key := []byte("cC13oL)ZUNEmX6!PlWRHzFDkyqi8.a%MB94e|#ns>G0(bIj&{S7gY$=<tv/Vh-w~}*uQ_;2A@Od,]Jp5[Kfr^:T?x")
	cins := []string{
		"610628199209080042",
		"410205198408208434",
		"54032519490819523X",
		"53030219690918165X",
		"450406201008100327",
		"511524200805066820",
	}
	for _, cin := range cins {
		distri, birthDate, sex, ciphertext, err := crypto.ShuffleEncryptCin(cin, key)
		if err != nil {
			t.Errorf("ShuffleEncryptCin(%s) => %d, %s, %d, %s  %v", cin, distri, birthDate, sex, ciphertext, err)
			return
		}
		got, err := ciphertext.Decrypt(distri, birthDate, key)
		if err != nil || got != cin {
			t.Errorf("Decrypt got %s, want %s %v", got, cin, err)
		}
	}
}
