package coding

import (
	"bytes"
)

// PadPKCS7
// PKCS7 是PKCS5的超集，PKCS5 只支持填充8字节，而PKCS7支持1-255字节
func PadPKCS7(src []byte, blockSize int) []byte {
	paddingLength := blockSize - len(src)%blockSize
	// 把填充宽度作为字符，填充到尾部
	padding := bytes.Repeat([]byte{byte(paddingLength)}, paddingLength)
	return append(src, padding...)
}
func UnpadPKCS7(src []byte) ([]byte, bool) {
	length := len(src)
	end := length - int(src[length-1])
	if end == 0 {
		return src, true
	}
	if end > 0 {
		return src[:end], true
	}
	return nil, false
}
