package coding

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"github.com/aarioai/airis/pkg/arrmap"
	"io"
	"strings"
)

// SortedParams
func SortedParams(params map[string]string, ignoredSignKey string, joint bool) string {
	keys := arrmap.SortedKeysFunc(params, func(key string, value string) (string, bool) {
		return key, key != ignoredSignKey
	})
	var signStr strings.Builder
	for i, k := range keys {
		if params[k] != "" {
			if joint {
				// 用 = & 连接
				if i > 0 {
					signStr.WriteByte('&')
				}
				signStr.WriteString(k)
				signStr.WriteByte('=')
				signStr.WriteString(params[k])
			} else {
				signStr.WriteString(k)
				signStr.WriteString(params[k])
			}
		}
	}

	return signStr.String()
}

// WriteSortedParams
func WriteSortedParams(w *bufio.Writer, params map[string]string, ignoredSignKey string, joint bool) {
	keys := arrmap.SortedKeysFunc(params, func(key string, value string) (string, bool) {
		return key, key != ignoredSignKey
	})
	for i, k := range keys {
		v := params[k]
		if v == "" {
			continue
		}
		if joint {
			// 用 = & 连接
			if i > 0 {
				w.WriteByte('&')
			}
			w.WriteString(k)
			w.WriteByte('=')
			w.WriteString(v)
		} else {
			w.WriteString(k)
			w.WriteString(v)
		}
	}
	return
}

// Sha1Signature Sha1 签名，大写结果
func Sha1Signature(params map[string]string, ignoredSignKey string, bufsize int, joint bool) string {
	h1 := sha1.New()
	if bufsize > 0 {
		// specify memory size
		bufw := bufio.NewWriterSize(h1, bufsize)
		WriteSortedParams(bufw, params, ignoredSignKey, joint)
		bufw.Flush()
	} else {
		// not specify memory size
		p := SortedParams(params, ignoredSignKey, joint)
		io.WriteString(h1, p)
	}
	bs := make([]byte, hex.EncodedLen(h1.Size()))
	hex.Encode(bs, h1.Sum(nil))
	return string(bytes.ToUpper(bs))
}

// HmacSignature HMAC 签名
// e.g. HmacSignature(sha1.New(), )
//func HmacSignature(h func() hash.Hash, params map[string]interface{}) string {
//	p := []byte("s")
//	mac := hmac.New(h, p)
//	mac.Write()
//}
