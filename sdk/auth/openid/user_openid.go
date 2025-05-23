package openid

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/golib/enumz"
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/sdk/auth/configz"
	"github.com/aarioai/golib/typez"
	"github.com/kataras/iris/v12"
	"strconv"
	"time"
)

func parseUserOpenidSvc(s string) (typez.Svc, bool) {
	// 只有base36
	var x string
	for _, b := range s {
		if b < 'A' || b > 'Z' {
			x += string(b)
		}
	}
	v, err := strconv.ParseUint(x, 36, 32)
	if err != nil {
		return 0, false
	}
	return typez.Svc(v), true
}

// 将 appid secret 转为 des cbc
func (s *Service) UserOpenidDesCBCKey(appid, secret string) ([]byte, []byte, *ae.Error) {
	str := []byte(appid + secret)
	length := coding.DesKeyLength
	key := str[len(str)-length:]
	iv := str[:length]
	if len(key) != length || len(iv) != length {
		return nil, nil, NewCode(ae.VariantAlsoNegotiates, "invalid key or iv")
	}
	return key, iv, nil
}

// EncodeUserOpenid svc  需要用作解密时候参数，所以独立出来 svc+ DES( + uid|时间戳|appid)
// 结构： [len:1][svc:$len][ds]
//
//	         ds => [hash:16][factor]
//					factor => [len-ts:1][ts:N]|[xu:N]
func (s *Service) EncodeUserOpenid(ctx context.Context, svc typez.Svc, uid uint64) (string, time.Duration, *ae.Error) {
	appid, secret, e := s.secretHandler(ctx, s.app, svc)
	if e != nil {
		return "", 0, e
	}
	desKey, desIV, e := s.UserOpenidDesCBCKey(appid, secret)
	if e != nil {
		return "", 0, e
	}
	expAt := time.Now().Add(configz.OpenidTTL).Unix() // 过期时间
	ts := strconv.FormatInt(expAt, 36)
	// DES 加密尽量不要用固定的字符串，固定的要与时间变化下
	u := "0" // uid=0时
	if uid > 0 {
		xu := int64(uid) - expAt // uid 可能为0；不为0时，uid 一定比 秒数大
		u = strconv.FormatInt(xu, 36)
	}

	lts := strconv.FormatUint(uint64(len(ts)), 36)
	// 避免过多 | ，让人猜出来分隔符
	factor := lts + ts + u

	// 将appid hash 化
	h := md5.Sum([]byte(appid + factor))
	ha := hex.EncodeToString(h[:])
	d := ha[8:24] + factor // 16字符hash
	ds, err := coding.CbcEncryptToBase64([]byte(d), desKey, desIV)
	if err != nil {
		return "", 0, ae.NewError(err)
	}
	// svc 需要用作解密参数，所以需要独立
	sv := strconv.FormatUint(uint64(svc), 36)
	sv = coding.RandPad(sv, configz.OpenidEncodeSvcLen, 'A', 'B', false)
	return sv + string(ds), configz.OpenidTTL, nil
}

func (s *Service) DecodeUserOpenid(ctx context.Context, openid string) (string, typez.Svc, uint64, *ae.Error) {
	if len(openid) < configz.OpenidEncodeSvcLen {
		return "", 0, 0, ae.ErrorPreconditionFailed
	}
	svc, ok := parseUserOpenidSvc(openid[0:configz.OpenidEncodeSvcLen])
	if !ok {
		return "", 0, 0, ae.ErrorPreconditionFailed
	}

	appid, secret, e := s.secretHandler(ctx, s.app, svc)
	if e != nil {
		return "", 0, 0, e
	}

	desKey, desIV, e := s.UserOpenidDesCBCKey(appid, secret)
	if e != nil {
		return "", 0, 0, e
	}
	d, err := coding.CbcDecryptFromBase64([]byte(openid[configz.OpenidEncodeSvcLen:]), desKey, desIV)
	if err != nil || len(d) < 16 {
		return "", 0, 0, ae.ErrorVariantAlsoNegotiates
	}
	aha := string(d[0:16]) // appid 的 hash
	factor := string(d[16:])
	h := md5.Sum(append([]byte(appid), factor...))
	ha := hex.EncodeToString(h[:])
	if aha != ha[8:24] {
		return "", 0, 0, ae.ErrorPreconditionFailed
	}
	var lts uint64
	if lts, _ = strconv.ParseUint(factor[0:1], 36, 8); lts == 0 {
		return "", 0, 0, ae.ErrorPreconditionFailed
	}
	now := time.Now().Unix()
	ts, _ := strconv.ParseInt(factor[1:lts+1], 36, 64)
	if ts < now {
		return "", 0, 0, ae.ErrorPreconditionFailed
	}

	xu, _ := strconv.ParseInt(factor[lts+1:], 36, 64)
	if xu == 0 {
		return "", 0, 0, ae.ErrorPreconditionFailed
	}
	uid := uint64(xu + ts)
	return appid, svc, uid, nil
}

func (s *Service) ParseUserOpenid(ictx iris.Context, uid uint64) (appid string, svc typez.Svc, openUid uint64, e *ae.Error) {
	openid := ictx.GetHeader(enumz.HeaderOpenid)
	if openid == "" {
		e = ae.New(ae.PreconditionRequired, "require openid")
		return
	}
	ctx := acontext.FromIris(ictx)

	appid, svc, openUid, e = s.DecodeUserOpenid(ctx, openid)
	if e != nil {
		s.app.Log.Warn(acontext.FromIris(ictx), "openid:%s, %s", openid, e.Text())
		return
	}
	if openUid > 0 && uid != openUid {
		e = ae.New(ae.PreconditionFailed, "bad openid")
		return
	}
	return
}
