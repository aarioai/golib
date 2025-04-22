package weixingzh

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"github.com/aarioai/golib/lib/code/coding"
	"io"
	"sort"
)

// JSSDK 使用步骤
// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#4
// 1. 绑定域名：公众号后台“JS接口安全域名”
// 2. 引入JS SDK文件
// 3. 通过 JS Config 注入权限验证设置

// SortedParams
func sortedParams(params map[string]string, ignoredSignKey string, joint bool) string {
	signStr := ""
	keys := make([]string, 0, len(params)-1)
	for k, _ := range params {
		if k != ignoredSignKey {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for i, k := range keys {
		if params[k] != "" {
			if joint {
				// 用 = & 连接
				if i > 0 {
					signStr += "&"
				}
				signStr += k + "=" + params[k]
			} else {
				signStr += k + params[k]
			}
		}
	}

	return signStr
}

// Sha1Signature Sha1 签名，大写结果
func sha1Signature(params map[string]string, ignoredSignKey string, bufsize int, joint bool) string {
	h1 := sha1.New()
	if bufsize > 0 {
		// specify memory size
		bufw := bufio.NewWriterSize(h1, bufsize)
		coding.WriteSortedParams(bufw, params, ignoredSignKey, joint)
		bufw.Flush()
	} else {
		// not specify memory size
		p := sortedParams(params, ignoredSignKey, joint)
		io.WriteString(h1, p)
	}
	bs := make([]byte, hex.EncodedLen(h1.Size()))
	hex.Encode(bs, h1.Sum(nil))
	return string(bytes.ToUpper(bs))
}
