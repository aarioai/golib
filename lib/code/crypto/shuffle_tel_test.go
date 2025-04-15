package crypto_test

import (
	"bytes"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/lib/code/crypto"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"testing"
)

var telPairs = []string{
	"10000000000", // 特殊情况，全0
	"99999999999", // 特殊情况，全是9
	"15000777962",

	"15000777963",
	"18900777999",

	"19999999999",
	"0755-12345678",
	"0755-9203213",
	"0999-9999999",
	"0999-99999999",
	"0564-6582404",
	"0564—6582404",
	"0564——6582404",
	"（0564）—6582404",
	"（0564)—6582404",
	"（0564)———6582404",
	"021-1258",
	"021-12580",
}

func TestDesenseTelCiphertext(t *testing.T) {
	stdKey := stdfmt.Base62Digits + "~!@#$%^&*()_+-=[{}]:;<>,.?/"
	key := coding.Shuffle([]byte(stdKey))
	for _, tel := range telPairs {
		for i := 0; i < 2; i++ {
			scatter := i == 1
			ciphertext, _, _ := crypto.ShuffleEncryptTel(tel, key, scatter)

			_, err := ciphertext.Desensitize()
			if err != nil {
				t.Errorf("failed DesenseTelCiphertext(%s <= %s): %v", ciphertext, tel, err)
			}
		}
	}
}
func TestEncodeTel(t *testing.T) {
	stdKey := stdfmt.Base62Digits + "~!@#$%^&*()_+-=[{}]:;<>,.?/"
	key := coding.Shuffle([]byte(stdKey))
	keyClone := bytes.Clone(key)
	for _, s := range telPairs {
		for i := 0; i < 2; i++ {
			scatter := i == 1
			gotCipher, tel, err := crypto.ShuffleEncryptTel(s, key, scatter)
			if err != nil {
				t.Errorf("failed ShuffleEncryptTel(%s, scatter:%v): %v", s, scatter, err)
				continue
			}
			if !bytes.Equal(key, keyClone) {
				t.Error("ShuffleEncryptTel " + afmt.ErrmsgSideEffect(key))
				continue
			}

			if len(gotCipher) < 6 || len(gotCipher) > crypto.TelCipherSafeLen {
				t.Logf("%s(scatter: %v) => %s (len:%d)\n", s, scatter, gotCipher, len(gotCipher))
			}
			gotTel, err := gotCipher.Decrypt(key)
			if err != nil {
				t.Errorf("failed Decrypt(%s, scatter:%v  <== %s): %v", gotCipher, scatter, s, err)
				continue
			}
			if gotTel.String() != tel.String() {
				t.Errorf("Decrypt(%s, scatter:%v <==%s) got `%s`, want `%s`", gotCipher, scatter, s, gotTel.String(), tel.String())
				continue
			}
			if !bytes.Equal(key, keyClone) {
				t.Error("Decrypt " + afmt.ErrmsgSideEffect(key))
				continue
			}

			// 测试不scatter key，每次是否加密结果是否都一致
			gotCipher2, _, err := crypto.ShuffleEncryptTel(tel.Local(), key, false)
			if err != nil {
				t.Errorf("ShuffleEncryptTel(%s) %v", tel.Local(), err)
				continue
			}
			if (!scatter && gotCipher != gotCipher2) || (scatter && gotCipher == gotCipher2) {
				t.Errorf("ShuffleEncryptTel(%s, scatter:%v), ecncrypt twice got `%s` and `%s`", tel.Local(), scatter, gotCipher2, gotCipher)
				continue
			}

		}

	}

}
