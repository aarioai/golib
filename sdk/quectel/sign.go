package quectel

import (
	"crypto/sha1"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
)

func (c *Client) sign(param url.Values) string {
	keys := make([]string, 0)
	for key, _ := range param {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	str := c.appSecret
	for _, key := range keys {
		str += key + param.Get(key)
	}
	str += c.appSecret

	md5Ctx := sha1.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)
	return strings.ToLower(hex.EncodeToString(cipherStr))
}
