package coding_test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aarioai/golib/lib/code/coding"
	"math/rand/v2"
	"testing"

	"github.com/aarioai/airis/pkg/afmt"
)

func TestEncryptGCM(t *testing.T) {
	key := []byte(coding.RandASCIICode(coding.GcmKeyLength - 5))
	text := []byte("Hello, Aario!")
	if _, err := coding.GcmEncryptToBase64(text, key); err == nil {
		t.Errorf("GcmEncryptToBase64 invalid key length %d, want error", len(key))
	}
	// 16位key加密
	key = []byte(coding.RandASCIICode(coding.GcmKeyLength))
	ciphertext, err := coding.GcmEncryptToBase64(text, key)
	if err != nil {
		t.Errorf("GcmEncryptToBase64 error: %s", err)
		return
	}

	// 16位key解密
	if _, err = coding.GcmDecryptFromBase64(ciphertext, key); err != nil {
		t.Errorf("GcmDecryptFromBase64 error: %s", err)
		return
	}
}
func TestEncryptCBC(t *testing.T) {
	text := []byte(coding.RandAlphabets(rand.IntN(100)+5, 0))
	//密钥，长度必须8位或24位
	key := []byte(coding.RandAlphabets(8, 0))
	iv := []byte(coding.RandAlphabets(8, 0))

	if err := testEncryptCBC(text, key, iv); err != nil {
		t.Error(err)
	}
	if err := testEncryptCBCToBase64(text, key, iv); err != nil {
		t.Error(err)
	}
}
func TestEncryptECB(t *testing.T) {
	text := []byte(coding.RandAlphabets(rand.IntN(100)+5, 0))
	//密钥，长度必须8位或24位
	key := []byte(coding.RandAlphabets(8, 0))

	if err := testEncryptECB(text, key); err != nil {
		t.Error(err)
	}
	if err := testEncryptECBToBase64(text, key); err != nil {
		t.Error(err)
	}
}

func testEncryptCBC(text, key, iv []byte) error {
	keyClone := bytes.Clone(key)
	ivClone := bytes.Clone(iv)

	// 测试加密
	ciphertext, err := coding.CbcEncrypt(text, key, iv)
	if err != nil {
		return err
	}
	if !bytes.Equal(iv, ivClone) {
		return errors.New("CbcEncrypt iv " + afmt.ErrmsgSideEffect(iv))
	}
	if !bytes.Equal(key, keyClone) {
		return errors.New("CbcEncrypt key " + afmt.ErrmsgSideEffect(key))
	}
	// 测试解密
	got, err := coding.CbcDecrypt(ciphertext, key, iv)
	if err != nil {
		return err
	}
	if !bytes.Equal(text, got) {
		return fmt.Errorf("CbcDecrypt encypt %s to %s, decrypt got %s", string(text), string(ciphertext), string(got))
	}
	if !bytes.Equal(iv, ivClone) {
		return errors.New("CbcDecrypt iv " + afmt.ErrmsgSideEffect(iv))
	}
	if !bytes.Equal(key, keyClone) {
		return errors.New("CbcDecrypt key " + afmt.ErrmsgSideEffect(key))
	}
	return nil
}
func testEncryptCBCToBase64(text, key, iv []byte) error {
	keyClone := bytes.Clone(key)
	ivClone := bytes.Clone(iv)

	// 测试加密
	ciphertext, err := coding.CbcEncryptToBase64(text, key, iv)
	if err != nil {
		return err
	}
	if !bytes.Equal(iv, ivClone) {
		return errors.New("CbcEncryptToBase64 iv " + afmt.ErrmsgSideEffect(iv))
	}
	if !bytes.Equal(key, keyClone) {
		return errors.New("CbcEncryptToBase64 key " + afmt.ErrmsgSideEffect(key))
	}
	// 测试解密
	got, err := coding.CbcDecryptFromBase64(ciphertext, key, iv)
	if err != nil {
		return err
	}
	if !bytes.Equal(text, got) {
		return fmt.Errorf("CbcDecryptFromBase64 encypt %s to %s, decrypt got %s", string(text), string(ciphertext), string(got))
	}
	if !bytes.Equal(iv, ivClone) {
		return errors.New("CbcDecryptFromBase64 iv " + afmt.ErrmsgSideEffect(iv))
	}
	if !bytes.Equal(key, keyClone) {
		return errors.New("CbcDecryptFromBase64 key " + afmt.ErrmsgSideEffect(key))
	}
	return nil
}

func testEncryptECB(text, key []byte) error {
	keyClone := bytes.Clone(key)

	// 测试加密
	ciphertext, err := coding.EcbEncrypt(text, key)
	if err != nil {
		return err
	}
	if !bytes.Equal(key, keyClone) {
		return errors.New("EcbEncrypt key " + afmt.ErrmsgSideEffect(key))
	}
	// 测试解密
	got, err := coding.EcbDecrypt(ciphertext, key)
	if err != nil {
		return err
	}
	if !bytes.Equal(text, got) {
		return fmt.Errorf("EcbDecrypt encypt %s to %s, decrypt got %s", string(text), string(ciphertext), string(got))
	}
	if !bytes.Equal(key, keyClone) {
		return errors.New("EcbDecrypt key " + afmt.ErrmsgSideEffect(key))
	}
	return nil
}

func testEncryptECBToBase64(text, key []byte) error {
	keyClone := bytes.Clone(key)

	// 测试加密
	ciphertext, err := coding.EcbEncryptToBase64(text, key)
	if err != nil {
		return err
	}
	if !bytes.Equal(key, keyClone) {
		return errors.New("EcbEncryptToBase64 key " + afmt.ErrmsgSideEffect(key))
	}
	// 测试解密
	got, err := coding.EcbDecryptFromBase64(ciphertext, key)
	if err != nil {
		return err
	}
	if !bytes.Equal(text, got) {
		return fmt.Errorf("EcbDecryptFromBase64 encypt %s to %s, decrypt got %s", string(text), string(ciphertext), string(got))
	}
	if !bytes.Equal(key, keyClone) {
		return errors.New("EcbDecryptFromBase64 key " + afmt.ErrmsgSideEffect(key))
	}
	return nil
}
