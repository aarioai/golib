package test

import (
	"bytes"
	"fmt"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/golib/lib/code/coding"
	"math/rand/v2"
)

func TestRSA[T ~string | ~[]byte](rsaKeyPairs [][2]func() (T, error)) *ae.Error {
	for i, rsaKeyPair := range rsaKeyPairs {
		prefix := fmt.Sprintf("self-test rsa %d ", i)
		pubkeyDERB64, err := rsaKeyPair[0]()
		if err != nil {
			return ae.NewE(prefix + "failed get pubkey: " + err.Error())
		}
		// 转化为 PKCS8
		pem := coding.RasToPKCS8([]byte(pubkeyDERB64), false, true)
		text := []byte(coding.RandASCIICode(rand.IntN(100)))
		// 加密为base64格式
		cipherBase64, err := coding.RsaEncryptToBase64(text, pem, false)
		if err != nil {
			return ae.NewE(prefix + "failed encrypt: " + err.Error())
		}

		// 解密
		privkeyDER, err := rsaKeyPair[1]()
		if err != nil {
			return ae.NewE(prefix + "failed get privkey: " + err.Error())
		}
		newText, err := coding.RsaDecryptFromBase64(cipherBase64, []byte(privkeyDER), true)
		if err != nil {
			return ae.NewE(prefix + "failed decrypt: " + err.Error())
		}
		if !bytes.Equal(text, newText) {
			return ae.NewE(prefix + "failed: not equal")
		}
	}
	return nil
}
