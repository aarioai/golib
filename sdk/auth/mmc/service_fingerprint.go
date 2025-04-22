package mmc

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"strconv"
	"strings"
	"time"
)

type fingerprintHeader struct {
	X        int // 组件左上角坐标 x,y，及宽高
	Y        int
	Width    int
	Height   int
	XA       int // 组件内“双击同意协议”及图标左上角相对X,Y位置,及宽高
	YA       int
	WidthA   int
	HeightA  int
	ClientMs int64 // 客户端创建组件客户端时间ms
	ServerMs int64 // 服务端下发配置服务器时间ms
}
type fingerprintBodySegment struct {
	Tag                []byte
	X1                 int // 相对 fingerprintHeader.X 的位置
	Y1                 int
	XA                 int
	YA                 int
	ClientTimeOffsetMs int64 // 相对于 fingerprintHeader.ClientMs
}

// 这里base64是没有填充=的
func (s *Service) decodeFingerprintClientDeskey(rsaKeyBase64Raw []byte) ([]byte, error) {
	clientRSA, err := stdfmt.DecodeBase64(rsaKeyBase64Raw)
	if err != nil {
		return nil, err
	}
	rsaPrivkey, err := s.rsaPrivDER()
	if err != nil {
		return nil, err
	}
	return coding.RsaDecrypt(clientRSA, []byte(rsaPrivkey), true)
}

// 这里base64是没有填充=的
func decodeFingerprintRecord(desDataBase64Raw, deskey []byte) ([]byte, error) {
	desData, err := stdfmt.DecodeBase64(desDataBase64Raw)
	if err != nil {
		return nil, err
	}
	return coding.EcbDecrypt(desData, deskey)
}

// EncryptClientRecordToFingerprint 将客户端上传的 recode 转为 fingerprint
// record=`base64_withoutPadding(rsa(deskey)) base64_withoutPadding(des(record))`
// 第一部分为 RSA 加密后的deskey，deskey 必须很短，目前只能是8字节
// 第二部分为通过deskey加密后的数据
func (s *Service) EncryptClientRecordToFingerprint(ctx context.Context, encryptedRecord []byte, apollo, userAgent, ip string) ([]byte, *ae.Error) {
	rs := bytes.Split(encryptedRecord, []byte{' '})
	if len(rs) != 2 || len(rs[0]) == 0 || len(rs[1]) == 0 {
		return nil, ae.NewNotAcceptable().WithDetail("record format error")
	}

	clientDeskey, err := s.decodeFingerprintClientDeskey(rs[0])
	if err != nil {
		return nil, ae.NewNotAcceptable().WithDetail(err.Error())
	}
	record, err := decodeFingerprintRecord(rs[1], clientDeskey)
	if err != nil {
		return nil, ae.NewNotAcceptable().WithDetail("des decrypt `%s` with key `%s` error: %s", string(encryptedRecord), string(clientDeskey), err.Error())
	}
	xid := NewFingerprintId()
	uuid := strconv.FormatUint(xid, 36)
	// 这里不做深度校验，避免过度被识别算法。这里只做简单校验。
	dr, e := s.EncryptFingerprint(record, apollo, userAgent, ip, uuid)
	if e != nil {
		return nil, e
	}
	if ok := s.h.CacheFingerprintUUID(ctx, uuid); !ok {
		return nil, ae.ErrorCacheFailed
	}
	return dr, nil
}

func (s *Service) VerifyFingerprint(ctx context.Context, fp []byte, apollo, userAgent, ip string) (int64, error) {
	record, encryptTimeMs, apollo2, userAgent2, ip2, uuid2, e := s.DecryptFingerprint(fp)
	if !s.app.Check(ctx, e) {
		return 0, fmt.Errorf("decrypt fingerprint failed: %s", e.Text())
	}

	if err := s.verifyFingerprintMatches(ip, ip2, apollo, apollo2, userAgent, userAgent2); err != nil {
		return 0, fmt.Errorf("fingerprint miss match: %s", err.Error())
	}

	// 双击后，30分钟内有效 —— 重新发送验证码也需要
	if time.Now().Add(-fingerprintValidDuration).UnixMilli() > encryptTimeMs {
		return 0, fmt.Errorf("fingerprint expired")
	}
	if !s.h.CheckMmcFingerprintUUID(ctx, uuid2) {
		return 0, fmt.Errorf("fingerprint uuid not match")
	}

	fpHeader, _, err := parseFpRecord(record)
	if err != nil {
		return 0, fmt.Errorf("parse fingerprint record failed: %s", err.Error())
	}

	// @TODO 使用布隆过滤器

	return fpHeader.ServerMs, nil
}

func (s *Service) verifyFingerprintMatches(ip, ip2, apollo, apollo2, userAgent, userAgent2 string) error {
	if ip != ip2 || apollo != apollo2 {
		return fmt.Errorf("ip(%s <> %s) apollo(%s <> %s)", ip, ip2, apollo, apollo2)
	}

	ua1 := normalizeUserAgent(userAgent)
	ua2 := normalizeUserAgent(userAgent2)

	if ua1 != ua2 {
		return fmt.Errorf("user-agent %s (len:%d) != %s (len:%d)",
			ua1, len(ua1), ua2, len(ua2))
	}

	return nil
}

func normalizeUserAgent(ua string) string {
	return strings.ToLower(strings.ReplaceAll(ua, " ", ""))
}

// 这个由于只会绑定一次，而且用户可能会很久才会绑定手机号，所以就不用那么复杂，直接这样加密的
// 这里通过mmc key 加密
func (s *Service) EncryptFingerprint(record []byte, apollo, userAgent, ip, uuid string) ([]byte, *ae.Error) {
	if _, _, err := parseFpRecord(record); err != nil {
		return nil, ae.NewNotAcceptable("invalid fingerprint record").WithDetail(err.Error())
	}

	uab := []byte(userAgent)
	base64UserAgent := make([]byte, base64.StdEncoding.EncodedLen(len(uab)))
	base64.StdEncoding.Encode(base64UserAgent, uab)
	tm := strconv.FormatInt(time.Now().UnixNano(), 10)
	size := len(tm) + len(uuid) + len(ip) + len(base64UserAgent) + len(apollo) + len(record) + 4 + len(fingerprintSeparator)
	var dr bytes.Buffer
	dr.Grow(size)
	dr.WriteString(tm)
	dr.WriteByte(' ')
	dr.WriteString(uuid)
	dr.WriteByte(' ')
	dr.WriteString(ip)
	dr.WriteByte(' ')
	dr.Write(base64UserAgent)
	dr.WriteByte(' ')
	dr.WriteString(apollo)
	dr.WriteString(fingerprintSeparator)
	dr.Write(record)

	secret, err := s.gcmKey()
	if err != nil {
		return nil, NewE("get mmc secret failed " + err.Error())
	}
	cipher, err := coding.GcmEncryptToBase64(dr.Bytes(), []byte(secret))
	if err != nil {
		return nil, NewE("encrypt fingerprint failed " + err.Error())
	}
	return cipher, nil
}

func (s *Service) DecryptFingerprint(fingerprint []byte) (record []byte, encryptTimeMs int64, apollo, userAgent, ip, uuid string, e *ae.Error) {
	secret, err := s.gcmKey()
	if err != nil {
		e = NewE("get mmc secret failed " + err.Error())
		return
	}

	dr, err := coding.GcmDecryptFromBase64(fingerprint, []byte(secret))
	if err != nil {
		e = NewError(err)
		return
	}
	arr := bytes.Split(dr, []byte(fingerprintSeparator))
	if len(arr) != 2 {
		e = NewE("invalid fingerprint")
		return
	}
	record = arr[1]

	xarr := bytes.Split(arr[0], []byte{' '})
	if len(xarr) != 5 {
		e = NewE("invalid fingerprint")
		return
	}
	encryptTimeMs, err = strconv.ParseInt(string(xarr[0]), 10, 64)
	if err != nil {
		e = NewE("invalid fingerprint")
		return
	}
	uuid = string(xarr[1])
	ip = string(xarr[2])
	base64UserAgent := make([]byte, base64.StdEncoding.DecodedLen(len(xarr[3])))
	ul, err := base64.StdEncoding.Decode(base64UserAgent, xarr[3])
	if err != nil {
		e = NewE("invalid fingerprint")
		return
	}
	userAgent = string(base64UserAgent[:ul])
	apollo = string(xarr[4])
	return
}

func parseFpRecord(record []byte) (*fingerprintHeader, []fingerprintBodySegment, error) {
	// 使用\n分开
	segs := bytes.Split(record, []byte{'\n'})
	if len(segs) < 3 {
		return nil, nil, fmt.Errorf("invalid record length: %d", len(segs))
	}

	header, err := parseHeader(segs[0])
	if err != nil {
		return nil, nil, err
	}
	bodies, err := parseBodies(segs[1:])
	if err != nil {
		return nil, nil, err
	}

	return header, bodies, nil
}

func parseHeader(data []byte) (*fingerprintHeader, error) {
	headerSegs := bytes.Split(data, []byte{','})
	if len(headerSegs) != fpHeaderLength {
		return nil, fmt.Errorf("invalid header length: %d", len(headerSegs))
	}

	h := make([]int64, fpHeaderLength)
	for i := 0; i < fpHeaderLength; i++ {
		val, err := strconv.ParseInt(string(headerSegs[i]), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid header value: %s", string(headerSegs[i]))
		}
		h[i] = val
	}
	header := fingerprintHeader{
		X:        int(h[0]),
		Y:        int(h[1]),
		Width:    int(h[2]),
		Height:   int(h[3]),
		XA:       int(h[4]),
		YA:       int(h[5]),
		WidthA:   int(h[6]),
		HeightA:  int(h[7]),
		ClientMs: h[8],
		ServerMs: h[9],
	}
	return &header, nil
}
func parseBodies(segs [][]byte) ([]fingerprintBodySegment, error) {
	bodies := make([]fingerprintBodySegment, 0, len(segs))
	for _, seg := range segs {
		segParts := bytes.Split(seg, []byte{','})
		if len(segParts) != fpBodySegmentLength {
			return nil, fmt.Errorf("invalid body length: %d", len(segParts))
		}
		values := make([]int64, fpBodySegmentLength-1)
		for i := 1; i < fpBodySegmentLength-1; i++ {
			val, err := strconv.ParseInt(string(segParts[i]), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid body value: %s", string(segParts[i]))
			}
			values[i-1] = val
		}
		body := fingerprintBodySegment{
			Tag:                segParts[0],
			X1:                 int(values[0]),
			Y1:                 int(values[1]),
			XA:                 int(values[2]),
			YA:                 int(values[3]),
			ClientTimeOffsetMs: values[4],
		}
		bodies = append(bodies, body)
	}
	return bodies, nil
}

func (h fingerprintHeader) serialize() []byte {
	segs := make([][]byte, fpHeaderLength)
	segs[0] = []byte(types.FormatInt(h.X))
	segs[1] = []byte(types.FormatInt(h.Y))
	segs[2] = []byte(types.FormatInt(h.Width))
	segs[3] = []byte(types.FormatInt(h.Height))
	segs[4] = []byte(types.FormatInt(h.XA))
	segs[5] = []byte(types.FormatInt(h.YA))
	segs[6] = []byte(types.FormatInt(h.WidthA))
	segs[7] = []byte(types.FormatInt(h.HeightA))
	segs[8] = []byte(strconv.FormatInt(h.ClientMs, 10))
	segs[9] = []byte(strconv.FormatInt(h.ServerMs, 10))
	return bytes.Join(segs, []byte{','})
}

func (b fingerprintBodySegment) serialize() []byte {
	segs := make([][]byte, fpBodySegmentLength)
	segs[0] = b.Tag
	segs[1] = []byte(types.FormatInt(b.X1))
	segs[2] = []byte(types.FormatInt(b.Y1))
	segs[3] = []byte(types.FormatInt(b.XA))
	segs[4] = []byte(types.FormatInt(b.YA))
	segs[5] = []byte(types.FormatInt(b.ClientTimeOffsetMs))
	return bytes.Join(segs, []byte{','})
}
