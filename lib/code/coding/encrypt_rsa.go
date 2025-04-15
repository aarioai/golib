package coding

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

// RsaEncrypt 使用公钥加密，支持PEM和DER格式
// @warn 这个生成的是二进制字节，如果传递给客户端，应使用 RsaEncryptToBase64
// @test 保证不修改 pubkey
func RsaEncrypt(msg []byte, pubkey []byte, isDER bool) ([]byte, error) {
	var pub *rsa.PublicKey
	var err error

	if isDER {
		pub, err = rsaParsePubkeyFromDER(pubkey)
	} else {
		pub, err = rsaParsePubkeyFromPEM(pubkey)
	}
	if err != nil {
		return nil, fmt.Errorf("parse public key error: %w", err)
	}
	return rsa.EncryptPKCS1v15(rand.Reader, pub, msg)
}

// RsaDecrypt 使用私钥解密，支持PEM和DER格式
// @test 保证不修改 privkey
func RsaDecrypt(cipher []byte, privkey []byte, isDER bool) ([]byte, error) {
	var priv *rsa.PrivateKey
	var err error

	if isDER {
		priv, err = rsaParsePrivkeyFromDER(privkey)
	} else {
		priv, err = rsaParsePrivkeyFromPEM(privkey)
	}
	if err != nil {
		return nil, fmt.Errorf("parse private key error: %w", err)
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipher)
}

// RsaEncryptToBase64 加密并返回Base64编码结果，方便传递给客户端
// @test 保证不修改 pubkey
func RsaEncryptToBase64(msg []byte, pubkey []byte, isDER bool) (string, error) {
	cipher, err := RsaEncrypt(msg, pubkey, isDER)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipher), nil
}

// RsaDecryptFromBase64 从Base64解码后解密
// @test 保证不修改 privkey
func RsaDecryptFromBase64(cipherBase64 string, privkey []byte, isDER bool) ([]byte, error) {
	cipher, err := base64.StdEncoding.DecodeString(cipherBase64)
	if err != nil {
		return nil, fmt.Errorf("base64 decode error: %w", err)
	}
	return RsaDecrypt(cipher, privkey, isDER)
}

// RasToPKCS8 将DER格式(二进制或base64编码)转换为PKCS8 PEM格式
// @test 保证不修改 der
func RasToPKCS8(der []byte, isPrivate, isBase64DER bool) []byte {
	const lineLength = 64 // 每行64个字符

	prefix := "-----BEGIN PUBLIC KEY-----\n"
	suffix := "-----END PUBLIC KEY-----\n"
	if isPrivate {
		prefix = "-----BEGIN PRIVATE KEY-----\n"
		suffix = "-----END PRIVATE KEY-----\n"
	}

	derBase64 := der
	if !isBase64DER {
		// 使用string，避免修改 der 产生副作用
		derBase64 = make([]byte, base64.StdEncoding.EncodedLen(len(der)))
		base64.StdEncoding.Encode(derBase64, der)
	}

	var s strings.Builder
	s.Grow(len(prefix) + len(suffix) + len(derBase64) + len(derBase64)/lineLength + 1)
	s.WriteString(prefix)
	for i := 0; i < len(derBase64); i += lineLength {
		end := i + lineLength
		if end > len(derBase64) {
			end = len(derBase64)
		}
		s.Write(derBase64[i:end])
		s.WriteByte('\n')
	}
	s.WriteString(suffix)
	return []byte(s.String())
}

// rsaParsePrivkeyFromDER 从DER格式的私钥数据中解析出RSA私钥
// @test 保证不修改 der
func rsaParsePrivkeyFromDER(der []byte) (*rsa.PrivateKey, error) {
	var priv *rsa.PrivateKey
	// 尝试PKCS8格式
	pk8, err := x509.ParsePKCS8PrivateKey(der)
	if err == nil {
		var ok bool
		if priv, ok = pk8.(*rsa.PrivateKey); ok {
			return priv, nil
		}
	}

	// 尝试PKCS1格式解析
	if priv, err = x509.ParsePKCS1PrivateKey(der); err != nil {
		return nil, errors.New("parsed key is not an RSA private key")
	}
	return priv, nil
}

// rsaParsePrivkeyFromPEM 从PEM格式的私钥数据中解析出RSA私钥
// @test 保证不修改 privatePEM
func rsaParsePrivkeyFromPEM(privatePEM []byte) (*rsa.PrivateKey, error) {
	block, rest := pem.Decode(privatePEM)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	if len(rest) > 0 {
		return nil, errors.New("extra data after PEM block")
	}

	// 验证密钥类型
	if block.Type != "PRIVATE KEY" && block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected PEM block type: %s", block.Type)
	}
	return rsaParsePrivkeyFromDER(block.Bytes)
}

// rsaParsePubkeyFromPEM 解码pem字节
// PKCS1(PEM格式）被PKCS8取代 -----BEGIN RSA PRIVATE KEY---
// PKCS8(PEM格式）默认使用  -----BEGIN PRIVATE KEY---
// DER 二进制格式，计算最原始状态。可以跟PEM格式互相转换，更节省空间和计算量
// @test 保证不修改 pemBytes
func rsaParsePubkeyFromPEM(pemBytes []byte) (pub *rsa.PublicKey, err error) {
	block, rest := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	if len(rest) > 0 {
		return nil, errors.New("extra data after PEM block")
	}

	if block.Type != "PUBLIC KEY" && block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("unexpected PEM block type: %s", block.Type)
	}

	return rsaParsePubkeyFromDER(block.Bytes)
}

// rsaParsePubkeyFromDER 从DER格式的公钥数据中解析出RSA公钥
// @test 保证不修改 der
func rsaParsePubkeyFromDER(der []byte) (*rsa.PublicKey, error) {
	ifc, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	pub, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("parsed key is not an RSA public key")
	}
	return pub, nil
}
