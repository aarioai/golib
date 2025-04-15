package crypto_test

import (
	"github.com/aarioai/golib/lib/code/crypto"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"strings"
	"testing"
)

func TestEncodeEmail(t *testing.T) {
	pairs := []string{
		"a@luexu.com",
		"Aario@luexu.com",
		"AarioAi@gmail.com",
		"x-yzdf_tom2320df@gmail.com",
		"x-yzdf_tom2320df@sina.com.cn",
		"x-yzdf_tom2320df@xzf.com.cn",
	}
	key := []byte(strings.ReplaceAll(stdfmt.ReadableAsciiCodes, string(crypto.EmailCipherSeparator), ""))

	for _, email := range pairs {
		for i := 0; i < 2; i++ {
			scatter := i == 0
			etel, err := crypto.ShuffleEncryptEmail(email, key, scatter)
			if err != nil {
				t.Errorf("ShuffleEncryptEmail(%s, scatter:%v) error %v\n", email, scatter, err)
				continue
			}
			te, err := etel.Decrypt(key)
			if err != nil || te != email {
				t.Errorf("Decrypt %s got %s, want %s  %v\n", etel, te, email, err)
				continue
			}
		}
	}
}
