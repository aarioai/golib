package coding

import "encoding/hex"

// EncodeHex 将字节切片编码为一串十六进制字符
// @example EncodeHex("hello world") ==>"68656c6c6f20776f726c64"
// @test 保证不修改 hexKey
func EncodeHex(text []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(text)))
	hex.Encode(dst, text)
	return dst
}

// DecodeHex 将十六进制字符串解码为字节切片，
// @example DecodeHex([]byte("68656c6c6f20776f726c64")) ==> hello world
// @test 保证不修改 hexKey
func DecodeHex(hexText []byte) ([]byte, error) {
	dst := make([]byte, hex.DecodedLen(len(hexText)))
	n, err := hex.Decode(dst, hexText)
	if err != nil {
		return nil, err
	}
	return dst[:n], err
}
