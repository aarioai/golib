package coding_test

import (
	"bytes"
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"slices"
	"testing"
)

// 512位公钥DER base64格式
const rsa512PubDERB64 = "MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAN33FbO2q6vSt8y+O+5NU9m6oYSvfr8I7URzN2Oy29KPlrvwAqUdjygjAeFN5/nZnyHAQjSC1gyQEVDUm4oK7QUCAwEAAQ=="

const rsa512PubkeyPKCS8 = `-----BEGIN PUBLIC KEY-----
MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAN33FbO2q6vSt8y+O+5NU9m6oYSvfr8I
7URzN2Oy29KPlrvwAqUdjygjAeFN5/nZnyHAQjSC1gyQEVDUm4oK7QUCAwEAAQ==
-----END PUBLIC KEY-----
`

const rsa512PrivkeyDERB64 = "MIIBVQIBADANBgkqhkiG9w0BAQEFAASCAT8wggE7AgEAAkEA3fcVs7arq9K3zL477k1T2bqhhK9+vwjtRHM3Y7Lb0o+Wu/ACpR2PKCMB4U3n+dmfIcBCNILWDJARUNSbigrtBQIDAQABAkEApM4Umv8Cr+0g8zA8J0/a9kqQKohzPzxNjwlNEwV2GfuJDwmpYZLEDl3HOYwjN0YQtoWIUdBR2aXWFa03O+PE8QIhAPL7EfylxlQWkZUm7vSJLuhizwswDqCEUXskH9byRym7AiEA6du/RqYHgfBUeRWZk1B2e5aRbcb7KnlIqwgip6FoeD8CIGlyLd8frg8l8C3zRHYY5qNw5fsr8t0ULywqhCrK37krAiEAxjtowzk/yexvnogpu08McCysr/Jou5M9fwURYykWBj8CIEm7obEFroWvL3TUfOz3+tDDV5D8cTprT2jaQyuyPrhl"

const rsa512PrivkeyPKCS8 = `-----BEGIN PRIVATE KEY-----
MIIBVQIBADANBgkqhkiG9w0BAQEFAASCAT8wggE7AgEAAkEA3fcVs7arq9K3zL47
7k1T2bqhhK9+vwjtRHM3Y7Lb0o+Wu/ACpR2PKCMB4U3n+dmfIcBCNILWDJARUNSb
igrtBQIDAQABAkEApM4Umv8Cr+0g8zA8J0/a9kqQKohzPzxNjwlNEwV2GfuJDwmp
YZLEDl3HOYwjN0YQtoWIUdBR2aXWFa03O+PE8QIhAPL7EfylxlQWkZUm7vSJLuhi
zwswDqCEUXskH9byRym7AiEA6du/RqYHgfBUeRWZk1B2e5aRbcb7KnlIqwgip6Fo
eD8CIGlyLd8frg8l8C3zRHYY5qNw5fsr8t0ULywqhCrK37krAiEAxjtowzk/yexv
nogpu08McCysr/Jou5M9fwURYykWBj8CIEm7obEFroWvL3TUfOz3+tDDV5D8cTpr
T2jaQyuyPrhl
-----END PRIVATE KEY-----
`

func TestRsaEncryptPEM(t *testing.T) {
	// 加密
	pubkey := []byte(rsa512PubkeyPKCS8)
	text := []byte("Hello, Aario!")
	cipher, err := coding.RsaEncrypt(text, pubkey, false)
	if err != nil {
		t.Error("RsaEncrypt", err)
		return
	}
	if !slices.Equal(pubkey, []byte(rsa512PubkeyPKCS8)) {
		t.Error("RsaEncrypt changed pubkey slice")
		return
	}

	// 解密
	privkey := []byte(rsa512PrivkeyPKCS8)
	newText, err := coding.RsaDecrypt(cipher, privkey, false)
	if err != nil {
		t.Fatal("RsaPEMDecrypt", err)
	}
	if !bytes.Equal(text, newText) {
		t.Error("RsaDecrypt", string(newText), "!=", string(text))
		return
	}
	if !slices.Equal(privkey, []byte(rsa512PrivkeyPKCS8)) {
		t.Error("RsaDecrypt changed privkey slice")
		return
	}
}
func TestRsaEncryptDER(t *testing.T) {
	// 加密
	pubDER, err := stdfmt.DecodeBase64([]byte(rsa512PubDERB64))
	if err != nil {
		t.Fatal("DecodeBase64", err)
	}
	pubDERClone := bytes.Clone(pubDER)
	text := []byte("Hello, Aario!")
	cipher, err := coding.RsaEncrypt(text, pubDER, true)
	if err != nil {
		t.Error("RsaEncrypt", err)
		return
	}
	if !slices.Equal(pubDER, pubDERClone) {
		t.Error("RsaEncrypt changed public DER slice")
		return
	}

	// 解密
	privDER, err := stdfmt.DecodeBase64([]byte(rsa512PrivkeyDERB64))
	if err != nil {
		t.Error("DecodeBase64", err)
		return
	}
	privDERClone := bytes.Clone(privDER)
	got, err := coding.RsaDecrypt(cipher, privDERClone, true)
	if err != nil {
		t.Error("RsaDecrypt", err)
		return
	}
	if !bytes.Equal(text, got) {
		t.Error("RsaDecrypt", string(got), "!=", string(text))
		return
	}
	if !slices.Equal(privDER, privDERClone) {
		t.Error("RsaDecrypt changed private DER slice")
		return
	}
}
func TestRsaEncryptToBase64(t *testing.T) {
	// 加密
	pubkey := []byte(rsa512PubkeyPKCS8)
	text := []byte("Hello, Aario!")
	cipher, err := coding.RsaEncryptToBase64(text, pubkey, false)
	if err != nil {
		t.Fatal("RsaEncryptToBase64", err)
	}
	if !slices.Equal(pubkey, []byte(rsa512PubkeyPKCS8)) {
		t.Error("RsaEncryptToBase64 changed pubkey slice")
		return
	}

	// 解密
	privkey := []byte(rsa512PrivkeyPKCS8)
	newText, err := coding.RsaDecryptFromBase64(cipher, privkey, false)
	if err != nil {
		t.Fatal("RsaPEMDecryptFromBase64", err)
	}
	if !bytes.Equal(text, newText) {
		t.Error("RsaDecryptFromBase64", string(newText), "!=", string(text))
		return
	}
	if !slices.Equal(privkey, []byte(rsa512PrivkeyPKCS8)) {
		t.Error("RsaDecryptFromBase64 changed privkey slice")
		return
	}
}
func TestRasToPKCS8(t *testing.T) {
	der := []byte(rsa512PubDERB64)
	pkcs8 := coding.RasToPKCS8(der, false, true)
	pem := string(pkcs8)
	if pem != rsa512PubkeyPKCS8 {
		t.Errorf("convert rsa der base64 to pkcs8 failed(length:%d <> %d): \n%s", len(rsa512PubkeyPKCS8), len(pem), pem)
		return
	}
	if !slices.Equal(der, []byte(rsa512PubDERB64)) {
		t.Error("RasToPKCS8 changed der slice")
		return
	}
}
