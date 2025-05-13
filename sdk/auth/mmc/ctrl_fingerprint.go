package mmc

import (
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/httpsvr"
	"github.com/aarioai/airis/aa/httpsvr/request"
	"github.com/aarioai/airis/aa/httpsvr/response"
	"github.com/aarioai/golib/enumz"
	"github.com/aarioai/golib/sdk/irisz"
	"github.com/aarioai/golib/typez"
	"github.com/kataras/iris/v12"
)

// PostFingerprint Controller 人机对抗指纹，提供给mmcid加密
func (s *Service) PostFingerprint(ictx iris.Context) {
	defer ictx.Next()
	r, resp, ctx := httpsvr.New(ictx)
	defer resp.CloseWith(r)
	// 不对apollo传错提示
	apollo := r.QueryWildFast(enumz.ParamApollo)
	userAgent := r.UserAgent()
	ip := acontext.ClientIP(ictx)
	record, e0 := r.BodyBytes("record")
	if e := resp.FirstError(e0); e != nil {
		return
	}

	fp, e := s.EncryptClientRecordToFingerprint(ctx, record, apollo, userAgent, ip)
	if e != nil {
		resp.WriteE(e)
		return
	}
	resp.Write(map[string]string{
		"fingerprint": string(fp),
	})
}

// AssertFingerprint Middleware 断言判断 fingerprint 是否验证成功
//  SpringMVC 中的Interceptor 拦截器也是相当重要和相当有用的，它的主要作用是拦截用户的请求并进行相应的处理。比如通过它来进行权限验证，或者是来判断用户是否登陆，或者是像12306 那样子判断当前时间是否是购票时间。

// Man-machine confrontation  人机对抗
/*
  1. 客户端双击完成验证的时候，通过fingerprint 接口 生成 mmcid 和 fingerprint
  2. mmcid 相当于
*/

func (s *Service) AssertFingerprint(ictx iris.Context) {
	ctx := acontext.FromIris(ictx)
	if s.disable {
		s.app.Log.Warn(ctx, "sdk_auth_mmc: fingerprint is disabled")
		ictx.Next()
		return
	}
	fingerprint := ictx.GetHeader(enumz.HeaderMMCFingerprint)
	if fingerprint == "" {
		response.JsonE(ictx, ae.NewBadParam(enumz.HeaderMMCFingerprint))
		return
	}

	// 不对apollo传错提示
	apollo := request.QueryWild(ictx, enumz.ParamApollo)
	_, err := typez.DecodeDeviceInfo(apollo)
	if err != nil {
		s.app.Log.Warn(ctx, "query apollo got %s, parse device info failed: %s\n", apollo, err.Error())
		response.JsonE(ictx, ae.New(ae.PreconditionFailed, "客户端安全验证未通过，请联系技术人员处理"))
		return
	}

	userAgent := ictx.GetHeader("User-Agent")
	ip := acontext.ClientIP(ictx)

	unixMs, err := s.VerifyFingerprint(ctx, []byte(fingerprint), apollo, userAgent, ip)
	if err != nil {
		response.JsonE(ictx, ae.New(ae.FailedDependency, "验证失败，请先双击同意协议，或联系技术人员处理"))
		return
	}
	irisz.SetFingerprintServerTime(ictx, unixMs)
	ictx.Next()
}
