package coding

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"errors"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"io"
)

// @see _说明.md
// AES 是 加强版 DES，都是对称加密
// AES ECB(text, key)   	--> 相同明文，密文相同。明文需要填充，有 padding oracle 攻击风险
// AES CBC(text, key, iv)  	-->	相同明文，密文可能不同。明文需要填充，有 padding oracle 攻击风险
// AES GCM(text, key{16/32字节})    	--> 明文不需要填充。

const (
	DesKeyLength   = 8 //  支持 任意ASCII符号
	X3DESKeyLength = 24

	GcmKeyLength   = 16 // AES-128  支持 任意ASCII符号
	X2GCMKeyLength = 32 // AES-256
)

// GcmEncryptToBase64 必须要用 hex.EncodeToString([]byte) 转为可读字符串
// GCM 加密，比CBC更安全。Cipher 支持大小写、密码长度更长。key长度必须是16或32字节
// @test 保证不修改 key
func GcmEncryptToBase64(text []byte, key []byte) ([]byte, error) {
	// 进去怎么进去，出来就怎么出来就行了。如果是用 hex.DecodeString() 转换的，decrypt取出来的时候，就需要用 hex.EncodeToString()
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	seal := gcm.Seal(nonce, nonce, text, nil)
	return stdfmt.EncodeBase64(seal, true, true), nil
}

// GcmDecryptFromBase64
// @test 保证不修改 key
func GcmDecryptFromBase64(ciphertext []byte, key []byte) ([]byte, error) {
	var err error
	if ciphertext, err = stdfmt.DecodeBase64(ciphertext); err != nil {
		return nil, err
	}
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	var nonce []byte
	nonce, ciphertext = ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
func newGCM(key []byte) (cipher.AEAD, error) {
	if len(key) != GcmKeyLength && len(key) != X2GCMKeyLength {
		return nil, ErrInvalidCipherKeyLen
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	return gcm, nil
}

// CbcEncrypt DES CBC 模式加密
// iv 初始化向量，用于设定对称算法初始向量。长度须与密钥相同。
// 返回的值，需要进行 base64处理才能显示为正常字符串
// @test 保证不修改 key 和 iv
func CbcEncrypt(text, key, iv []byte) ([]byte, error) {
	block, err := newDESBlock(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	src := PadPKCS7(text, blockSize)
	ciphertext := make([]byte, len(src))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(ciphertext, src)
	return ciphertext, nil
}

// CbcDecrypt
// @test 保证不修改 key 和 iv
func CbcDecrypt(ciphertext, key, iv []byte) ([]byte, error) {
	block, err := newDESBlock(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	dst := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(dst, ciphertext)
	var ok bool
	if dst, ok = UnpadPKCS7(dst); !ok {
		return nil, ErrInvalidPKCS7Padding
	}
	return dst, nil
}

// CbcEncryptToBase64
// @test 保证不修改 key
func CbcEncryptToBase64(src, key, iv []byte) ([]byte, error) {
	raw, err := CbcEncrypt(src, key, iv)
	if err != nil {
		return nil, err
	}
	return stdfmt.EncodeBase64(raw, true, true), nil
}

// CbcDecryptFromBase64
// @test 保证不修改 key
func CbcDecryptFromBase64(src, key, iv []byte) ([]byte, error) {
	erytext, err := stdfmt.DecodeBase64(src)
	if err != nil {
		return nil, err
	}
	return CbcDecrypt(erytext, key, iv)
}

// EcbEncrypt
// 一般用于对客户端临时校验，客户端生成的key都是临时的，因此使用ECB更快捷
func EcbEncrypt(text, key []byte) ([]byte, error) {
	block, err := newDESBlock(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	src := PadPKCS7(text, blockSize)
	ciphertext := make([]byte, len(src))
	dst := ciphertext // 这一步是必须的
	for len(src) > 0 {
		block.Encrypt(dst, src[:blockSize])
		src = src[blockSize:]
		dst = dst[blockSize:]
	}
	return ciphertext, nil
}

// EcbDecrypt
// 一般用于对客户端临时校验
func EcbDecrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := newDESBlock(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(ciphertext)%blockSize != 0 {
		return nil, ErrInvalidBlockSize
	}
	text := make([]byte, len(ciphertext))
	dst := text
	for len(ciphertext) > 0 {
		block.Decrypt(dst, ciphertext[:blockSize])
		ciphertext = ciphertext[blockSize:]
		dst = dst[blockSize:]
	}
	var ok bool
	if text, ok = UnpadPKCS7(text); !ok {
		return nil, ErrInvalidPKCS7Padding
	}
	return text, nil
}

// EcbEncryptToBase64
// @test 保证不修改 key
func EcbEncryptToBase64(src, key []byte) ([]byte, error) {
	raw, err := EcbEncrypt(src, key)
	if err != nil {
		return nil, err
	}
	return stdfmt.EncodeBase64(raw, true, true), nil
}

// EcbDecryptFromBase64
// @test 保证不修改 key
func EcbDecryptFromBase64(src, key []byte) ([]byte, error) {
	raw, err := stdfmt.DecodeBase64(src)
	if err != nil {
		return nil, err
	}
	return EcbDecrypt(raw, key)
}
func newDESBlock(key []byte) (cipher.Block, error) {
	switch len(key) {
	case DesKeyLength:
		return des.NewCipher(key)
	case X3DESKeyLength:
		return des.NewTripleDESCipher(key)
	}
	return nil, ErrInvalidCipherKeyLen
}
