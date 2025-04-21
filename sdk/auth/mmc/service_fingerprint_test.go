package mmc_test

import (
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"testing"
)

// 512位公钥DER base64格式
const rsa512PubDERB64 = "MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAN33FbO2q6vSt8y+O+5NU9m6oYSvfr8I7URzN2Oy29KPlrvwAqUdjygjAeFN5/nZnyHAQjSC1gyQEVDUm4oK7QUCAwEAAQ=="

const rsa512PrivkeyPEM = `-----BEGIN PRIVATE KEY-----
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

const (
	deskey          = "vLiQZlIU"
	clientRSABase64 = "B/8F3WPKi8X/Mb92H4ezpkffUbJ9l1JoNfWbHOWGoQ0DbypPRBOE0tU/YnTkFBQC2MyM2s2RpHKWxNnUgJVa3w"
	recordDESBase64 = "D0TpKPx6k6WoCMzBcOEJ6PRPL12iPnU0L6oSAsXYTOdUhpvNxW16xmYSpCyO9sVife24PhKWoi4mwM40EfxqfPzq8ehcJMfE1uCeeiU0Wi+A5xnItiyItOWm3RN9ghHE9Qe2SKVeQnNUfESA5K2ryYK8LZLIRLLtihp6QIE5biQWEYiZNp8uHlVQ3Mtdct+8wLoqhsoGbUdi0CZusJ/etFVuyvzioDroOP2Cw5kqp/qxuHgqFlIpDpdcTYTQfBMS/DcqzaFp6lBBLffEioL51zNj4HSCwsKmI4Npi7nUrH+J3e+XCSlSOGhL1dZLghxzTi19d5rFIgbMwsOLG5SmUOWm3RN9ghHErOky6HT7VcgTinEIgCVbedgGvvdBTY+FBviDe0tjJ6U66uGwcJCWDeoW6whLADjrZY7OFA7HE/WzAVmFWuF9zzkIAu0EiXoM78ytp2U8IhIjn3QV+Ux430fv820dzQtg/DcqzaFp6lDXIAjmgkpUi0wA1K8UVkjPe3xbEHzU3sClMxW1cvjGiYdBI0hvR0FrHa07y9thzBXexnnkCh9s4+Wm3RN9ghHEyrcVXrvkXsgmBJfgZ0cVk6FACiKwzwa0zGvDFnsyl6g+ops2y2cfXpWu8RQet8m12pl/0LeV3yVmoPjT7O7iRZ4v1p0WX3EU6PKGp/lFqO/Yb1mPQrnkZIEm+mIOWSYK/DcqzaFp6lB0jj2Z5WaMh+dsA1sI/T+eD5fugARXqUVxJL003O5a0IQGwYrUu9byAELkITEqhKeOqBitKJxLu+Wm3RN9ghHEBlhVLCPOO+PXlFf5hdb7vlEZNlLDqnwAjilOe6kr0CIFaOAOKp+HrZTz+H8BDaC1yQS0zABqMKklUgek/SzR3dnsfiiL58pX8o+1ZqiKKlwErnUdepuBfuth0Fj4R6Rq/DcqzaFp6lCDIbqAFPL+xFN/K9QFjcoDqRHM8CWEBh9GUUnhO9NKExt0r6yeiYRZpVdgoRCqq+DoQdqlb+o1j+Wm3RN9ghHEECg+4TPhFwbF6Rhk3J5QR7xvCxTJFhRz8wsCoPpkhng+ops2y2cfXlTzXPbcSQe0KpZ76eVyKF8Qa2e6aixHhwBjRyMnDHxd"
)

// RSA解密获得deskey
func testDecryptDeskey(t *testing.T) ([]byte, error) {
	privblock, _ := pem.Decode([]byte(rsa512PrivkeyPEM))
	if privblock == nil {
		return nil, errors.New("decode private key fail")
	}
	privkeyDER := privblock.Bytes

	// 测试base64解码
	clientRSA, err := stdfmt.DecodeBase64([]byte(clientRSABase64))
	if err != nil {
		return nil, fmt.Errorf("decode base64 fail %s", err.Error())
	}
	// 测试RSA解密
	got, err := coding.RsaDecrypt(clientRSA, privkeyDER, true)
	if err != nil || string(got) != deskey {
		pubkeyPEM := coding.RasToPKCS8([]byte(rsa512PubDERB64), false, true)
		t.Log("\n=========\n" + string(pubkeyPEM) + "\n===========\n")
		return nil, fmt.Errorf("invalid rsa encrypt got (%s) %v", string(got), err)
	}
	return got, nil
}
func testDecryptRecord(desBase64 []byte, deskey []byte) ([]byte, error) {
	desData, err := stdfmt.DecodeBase64(desBase64)
	if err != nil {
		return nil, err
	}
	return coding.EcbDecrypt(desData, deskey)
}

// 测试客户端JSEncrypt传递的参数
func TestFingerprintJSEncrypt(t *testing.T) {
	gotDeskey, err := testDecryptDeskey(t)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = testDecryptRecord([]byte(recordDESBase64), gotDeskey)
	if err != nil {
		t.Error(err)
		return
	}
	//t.Log(string(record))
}
