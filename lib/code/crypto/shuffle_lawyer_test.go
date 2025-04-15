package crypto_test

import (
	"github.com/aarioai/golib/lib/code/crypto"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"testing"
)

func TestEncryptLawyerLic(t *testing.T) {
	pairs := []string{
		"13101201811999999",
		"13101201811045709",  // 测试标准长度
		"131012018110000001", // 测试18位长度
		"131012018119999999", // 测试18位长度
	}
	key := []byte(stdfmt.ReadableAsciiCodes)

	for _, lic := range pairs {
		lt, d, y, ciphertext, err := crypto.ShuffleEncryptLawyerLic(lic, key)
		if err != nil {
			t.Errorf("ShuffleEncryptLawyerLic %s  %v", lic, err)
			return
		}
		got, err := ciphertext.Decrypt(lt, d, y, key)
		if err != nil || got != lic {
			t.Errorf("Decrypt lawyer lic got %s, want %s %v", got, lic, err)
			return
		}
	}
}
func TestEncryptLawyerCert(t *testing.T) {
	pairs := []string{
		"A20146106283000",
	}
	key := []byte("K#Le)R<:WO(.{hVlB;Ns4fAX@~SGyZPw!vU=}gYE7ibkpH?M_9zdJ5^%t]n*0>[j1IF/&3Da6Cqc8ru,xo|QT2m$-")

	for _, cert := range pairs {
		ty, y, d, ciphertext, err := crypto.ShuffleEncryptLawyerCert(cert, key)
		if err != nil {
			t.Errorf("faild ShuffleEncryptLawyerCert(%s) %v", cert, err)
			return
		}
		got, err := ciphertext.Decrypt(ty, y, d, key)
		if err != nil || got != cert {
			t.Errorf("Decrypt got %s, want %s %v", got, cert, err)
			return
		}
	}
}
