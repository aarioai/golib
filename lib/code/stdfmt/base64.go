package stdfmt

import (
	"bytes"
	"encoding/base64"
	"github.com/aarioai/airis/aa/ae"
)

// 移入airis 的要求是：必须airis 本身需要用到。如果用不到，就不要移过去！！！

// EncodeBase64 将src编码为base64，如果urlSafe为true，则返回URL-Safe的base64
// StdEncoding: [a-Z\d+/=]  尾部填充 =
// URLEncoding: [a-Z\d-_=]  替换 + 为 -, / 为 _ 尾部填充 =
// types.Base64Digits [a-Z\d_~]
func EncodeBase64[T []byte | string](src T, urlSafe, withoutPadding bool) []byte {
	byteSrc := []byte(src)
	if len(byteSrc) == 0 {
		return nil
	}
	encoder := base64.StdEncoding
	if urlSafe {
		encoder = base64.URLEncoding
	}
	if withoutPadding {
		encoder = encoder.WithPadding(base64.NoPadding)
	}
	dst := make([]byte, encoder.EncodedLen(len(byteSrc)))
	encoder.Encode(dst, byteSrc)
	return dst
}

// DecodeBase64 通用base64解码（含URL-Safe类型和标准类型，以及填充与不填充）
func DecodeBase64[T []byte | string](src T) ([]byte, error) {
	byteSrc := []byte(src)
	l := len(byteSrc)
	if l == 0 {
		return nil, ae.ErrEmptyInput
	}
	encoder := base64.StdEncoding
	// 检测是否为URL安全的base64
	isURLSafe := bytes.IndexByte(byteSrc, '-') >= 0 || bytes.IndexByte(byteSrc, '_') >= 0
	if isURLSafe {
		encoder = base64.URLEncoding
	}
	// 尾部填充 =，保证base64字符串长度是4的整倍数
	paddingRest := l % 4
	if paddingRest != 0 {
		encoder = encoder.WithPadding(base64.NoPadding)
	}
	dbuf := make([]byte, encoder.DecodedLen(len(byteSrc)))
	n, err := encoder.Decode(dbuf, byteSrc)
	return dbuf[:n], err
}
