package auth

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/golib/enumz"
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/sdk/auth/configz"
	"github.com/aarioai/golib/typez"
	"math/rand"
	"strconv"
	"time"
)

const (
	userTokenMd5Len    = 8 // 截取前8位
	staffLength        = 38
	userTokenBaseShift = 9
)

/*
    DES：数据加密标准，密钥偏短（56位）、生命周期短（避免被破解）。
3DES：密钥长度112位或168位，通过增加迭代次数提高安全性 。处理速度慢、密钥计算时间长、加密效率不高 。
AES：高级数据加密标准，能够有效抵御已知的针对DES算法的所有攻击 。密钥建立时间短、灵敏性好、内存需求低、安全性高 。

    在RSA加解密算法中提及到RSA加密明文会受密钥的长度限制，这就说明用RSA加密的话明文长度是有限制的，而在实际情况我们要进行加密的明文长度或许会大于密钥长度，这样一来我们就不得不舍去RSA加密了。对此，DES加密则没有此限制。

　　鉴于以上两点(个人观点)，单独的使用DES或RSA加密可能没有办法满足实际需求，所以就采用了RSA和DES加密方法相结合的方式来实现数据的加密。

　　其实现方式即：
    客户端随机生成 des密钥
　　1、信息(明文)采用DES密钥加密。
　　2、使用RSA加密前面的DES密钥信息。
　　最终将混合信息进行传递。
　　而接收方接收到信息后：
　　1、用RSA解密DES密钥信息。
　　2、再用RSA解密获取到的密钥信息解密密文信息。
　　最终就可以得到我们要的信息(明文)。
*/

// 1. 转化为 36进制字符串（strconv.FormatUint(99999999, 36)） AuthAt2、Factor2、Biz2、UA2、Uid2，Vuid2
// 2、进行填充(_) 得到  AuthAt3<7字符> + Factor3<7字符> + Biz3<4字符> + UA3<2字符> + Uid3<9位>+ Vuid3<9位>
func uPad(x uint64, n int) []byte {
	s := strconv.FormatUint(x, 36)
	return rightPad([]byte(s), n, '_')
}

func unPad(_svc, _xuid, _vuid, _ua, _authAt, _factor []byte) (svc typez.Svc, uid, vuid uint64, ua enumz.UA, authAt, factor int64, err error) {
	var x uint64
	var ok bool

	vuid, err = strconv.ParseUint(string(rightUnpad(_vuid, '_')), 36, 64)
	if err != nil {
		return
	}

	if x, err = strconv.ParseUint(string(rightUnpad(_authAt, '_')), 36, 64); err != nil {
		return
	}
	authAt = int64(x)
	x, err = strconv.ParseUint(string(rightUnpad(_xuid, '_')), 36, 64)
	if err != nil {
		return
	}
	if x == 0 {
		uid = vuid // mock 账号
	} else {
		uid = x + uint64(authAt)
	}
	if x, err = strconv.ParseUint(string(rightUnpad(_factor, '_')), 36, 64); err != nil {
		return
	}
	factor = int64(x)

	if x, err = strconv.ParseUint(string(rightUnpad(_svc, '_')), 36, 24); err != nil {
		return
	}
	svc = typez.Svc(x)

	if x, err = strconv.ParseUint(string(rightUnpad(_ua, '_')), 36, 16); err != nil {
		return
	}

	if ua, ok = enumz.NewUA(uint16(x)); !ok {
		err = errors.New("bad token for user-agent " + string(_ua))
		return
	}
	return
}

// 模拟 rand.Perm 洗牌算法
// staffLength + 1
// 生成：${1位随机数}混淆数据
func randPerm(svc, xuid, vuid, ua, authAt, factor []byte, psid string, odd bool) []byte {
	const base = 36               //base36
	n := rand.Int31n(int32(base)) // 一个随机数[0,36) 即 [0,35]  --> [0, 0b100011]
	swapRange := 1                // 交换间隔；奇数隔1位交换、偶数逐次交换
	if odd {
		swapRange++
		n |= 1 // [warn] 最大数为奇数的取奇数方法 -- > 这里比较特殊，  [0,35] | 1  的范围是 [1,35] 不会越界。 若是 [0,36] | 1 ===> [1, 37] 会越界
	} else {
		n = n >> 1 << 1 // 取偶数
	}
	r := strconv.FormatUint(uint64(n), base)[0]

	arr := [6][]byte{factor, authAt, xuid, svc, ua, vuid} // 手动随机一下
	l := staffLength + 1 + len(psid)
	var buf bytes.Buffer
	buf.Grow(l)
	buf.WriteByte(r)
	for _, aa := range arr {
		buf.Write(aa)
	}
	buf.WriteString(psid)
	ps := buf.Bytes()
	mid := l / 2 // 从第三位开始，进行轴对称交换
	for i := 2; i < mid; i += swapRange {
		ps[i], ps[staffLength-i] = ps[staffLength-i], ps[i]
	}
	return ps
}

func unPerm(ps []byte) (svc, xuid, vuid, ua, authAt, factor []byte, psid string, odd bool, err error) {
	r := ps[0] // 随机数
	var n uint64
	n, err = strconv.ParseUint(string(r), 36, 8)
	if err != nil {
		return
	}
	swapRange := int(n&1) + 1 // n&1 判断 n 是否为奇数
	odd = swapRange == 2
	mid := len(ps) / 2 // 从第三方开始，进行轴对称交换
	for i := 2; i < mid; i += swapRange {
		ps[i], ps[staffLength-i] = ps[staffLength-i], ps[i]
	}
	factor = ps[1:8]
	authAt = ps[8 : 8+7]
	xuid = ps[8+7 : 8+7+9]
	svc = ps[8+7+9 : 8+7+9+4]
	ua = ps[8+7+9+4 : 8+7+9+4+2]
	vuid = ps[8+7+9+4+2 : 8+7+9+4+2+9]
	psid = string(ps[8+7+9+4+2+9:])
	return
}
func replaceUserTokenUnderlines(atoken []byte) []byte {
	for i := userTokenMd5Len; i < len(atoken); i++ {
		if atoken[i] == '_' {
			atoken[i] = 'A' + atoken[i-userTokenMd5Len]%26 // 'A'=65, md5 实际保留 0:mn 位， A-Z 是 26个字母
		}
	}
	return atoken
}

// expire_in 由服务端限定即可，这样避免客户端access token被拦截，服务端无能为力。
// 随机码1位，决定排序
// access token 包含：AuthAt<int64>、Factor<int64>、SvcId<0-36^4-1 = 1679615>、UA<0-1295>、Vuid<uint64>、Uid<uint64> Flag<uint8>
// AuthAt 用于判断用户设置密码的时候，是否需要重新登陆；ExpiresIn 用户判断token是否过期
// 		1. uPad:  svc<3字符>, uid-authAt<9位 -> 避免被侦测到>, Vuid3<9位>,  UA3<2字符>,AuthAt<7字符>, Factor<7字符>,
// 		2. md5(svc + _xuid + _vuid + _ua + _authAt + _factor)，截取前8位
//      3、生成随机码，决定顺序
// 		4、整体移位混淆
//      5、将_字符，替换成  ascii('A'+ md5s[当前_所在字节位置 % 32]%26)

// 通过加上uid，大大增加token暴力攻击难度
// token 增加：biz 和 平台，保证不同业务，不同平台（微信、iOS、Android）都单独唯一登陆
//
//	用 []byte 是指针

// encryptUserToken
//
//	@Description:
//	@receiver s
//	@param svc
//	@param uid
//	@param vuid
//	@param ua
//	@param psid 设备唯一码
//
// @param authAt
// @param factor 自增数字
// @param secureLogin: bool  是否使用密码、验证码、第三方授权登录；区分有些通过识别设备唯一码psid不安全登录
// @return string
// @return *ae.Error
func (s *Service) encryptUserToken(svc typez.Svc, uid, vuid uint64, ua enumz.UA, psid string, authAt, factor int64, secureLogin bool) (string, *ae.Error) {
	//SvcId<0-36^4-1 = 1679615>、Ua<0-35>、 ExpiresIn<最多699天，60460000s>
	//if uid < uint64(authAt) {
	//	return "", ae.NewE("bad uid: %d to encrypt user token", uid)
	//}
	if svc > 1679615 || !ua.Valid() {
		return "", NewE("encrypt parameter is out of range (%d, %d)", svc, ua)
	}

	// 避免UID被侦察到总是不变；mock UID 会小于认证时间，这种情况暴露也无所谓
	var xuid uint64 // 0 表示 mock， mock uid = vuid
	if uid > uint64(authAt) {
		xuid = uid - uint64(authAt) // uid 一定大于 秒数
	}

	_svc := uPad(uint64(svc), 4)
	_xuid := uPad(xuid, 9)
	_vuid := uPad(vuid, 9)
	_ua := uPad(uint64(ua), 2)
	_authAt := uPad(uint64(authAt), 7)
	_factor := uPad(uint64(factor), 7)

	var b bytes.Buffer
	b.Grow(staffLength + len(configz.UserTokenCryptMd5Key))
	b.Write(_svc)
	b.Write(_xuid)
	b.Write(_vuid)
	b.Write(_ua)
	b.Write(_authAt)
	b.Write(_factor)
	b.WriteString(configz.UserTokenCryptMd5Key)
	h := md5.Sum(b.Bytes())
	md5s := hex.EncodeToString(h[:])

	//s.app.Log.Debug(context.Background(), "encode: svc: (%d, %s) uid:(%d, %s) vuid(%d, %s) ua(%d, %s) auth_at(%s, %s) factor(%d, %s) md5key: %s", _svc, string(_svc), uid, string(_xuid), vuid, string(_vuid), ua, string(_ua), authAt, string(_authAt), factor, string(_factor), md5key)

	p := randPerm(_svc, _xuid, _vuid, _ua, _authAt, _factor, psid, secureLogin) // 39
	var atb bytes.Buffer
	atb.Grow(len(p) + 8)
	atb.WriteString(md5s[0:userTokenMd5Len])
	atb.Write(p)
	atoken := atb.Bytes()
	base := []byte(configz.UserTokenShuffleBase)
	err := coding.ShuffleEncrypt(atoken, userTokenBaseShift, base)
	if err != nil {
		return "", ae.NewError(err)
	}
	atoken = replaceUserTokenUnderlines(atoken)
	return string(atoken), nil
}

func (s *Service) DbgDecryptUserToken(ctx context.Context, token string) (svc typez.Svc, uid, vuid uint64, ua enumz.UA, psid string, authAt int64, expiresIn int64, factor int64, secureLogin bool, e *ae.Error) {
	expiresIn = configz.UserTokenTTLs
	svc, uid, vuid, ua, psid, authAt, factor, secureLogin, e = s.decryptUserToken(ctx, token)
	return
}

// 这个pass只能用于服务内，因为没有对客户端传来的做验证。客户端请求需要用下面 ParseSsoAccessToken
func (s *Service) decryptUserToken(ctx context.Context, atoken string) (svc typez.Svc, uid, vuid uint64, ua enumz.UA, psid string, authAt int64, factor int64, secureLogin bool, e *ae.Error) {

	// 如果用指针需要深度复制，而string是const，强转为 []byte 会重新开辟新的内存空间
	token := []byte(atoken)
	for i, b := range token {
		// 下划线是 95，小写字符是从97开始；A-Z 是 65-90；  数字0-9是：48-57
		if b < 91 && b > 64 {
			token[i] = '_'
		}
	}
	base := []byte(configz.UserTokenCryptMd5Key)
	err := coding.ShuffleDecrypt(token, userTokenBaseShift, base)
	if err != nil {
		e = ae.NewError(err)
		return
	}

	var (
		_svc, _xuid, _vuid, _ua, _authAt, _factor []byte
	)
	_svc, _xuid, _vuid, _ua, _authAt, _factor, psid, secureLogin, err = unPerm(token[userTokenMd5Len:])
	if err != nil {
		s.app.Log.Info(ctx, "bad token "+atoken+" :"+err.Error())
		e = ae.ErrorUnauthorized
		return
	}

	var b bytes.Buffer
	b.Grow(staffLength + len(configz.UserTokenCryptMd5Key))
	b.Write(_svc)
	b.Write(_xuid)
	b.Write(_vuid)
	b.Write(_ua)
	b.Write(_authAt)
	b.Write(_factor)
	b.WriteString(configz.UserTokenCryptMd5Key)
	h := md5.Sum(b.Bytes())
	md5s := hex.EncodeToString(h[:])
	//s.app.Log.Debug(context.Background(), "decode: %s %s %s %s %s %s %s", string(_svc), string(_xuid), string(_vuid), string(_ua), string(_authAt), string(_factor), md5key)

	if bytes.Compare([]byte(md5s[0:userTokenMd5Len]), token[0:userTokenMd5Len]) != 0 {
		e = ae.NewBadParam("access_token")
		return
	}
	// // 		4、整体移位混淆 AuthAt4<7字符> +Factor3<7字符>+Biz4<4字符> + UA4<2字符> +  Uid4<9位> +  md5<3字符> ，得到32字符串

	svc, uid, vuid, ua, authAt, factor, err = unPad(_svc, _xuid, _vuid, _ua, _authAt, _factor)
	if uid == 0 {
		s.app.Log.Info(ctx, "bad token: %s, no uid", atoken)
		e = ae.ErrorUnauthorized
		return
	}
	now := time.Now().Unix()
	expiresIn := authAt + configz.UserTokenTTLs + configz.UserTokenTimeWindow - now
	if expiresIn < 0 {
		e = ae.ErrorUnauthorized
		return
	}
	return
}

func rightPad(s []byte, length int, pad byte) []byte {
	ls := len(s)
	g := length - ls
	if g > 0 {
		b := make([]byte, length)
		copy(b, s)
		for i := ls; i < length; i++ {
			b[i] = pad
		}
		return b
	}
	return s
}
func rightUnpad(s []byte, pad byte) []byte {
	for i, b := range s {
		if b == pad {
			return s[:i]
		}
	}
	return s
}
