package midiris

import (
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/golib/enumz"
	"github.com/aarioai/golib/sdk/auth/configz"
	"github.com/aarioai/golib/typez"
	"github.com/kataras/iris/v12"
	"strconv"
)

func SetUid(ictx iris.Context, svc typez.Svc, uid, vuid uint64) bool {
	if uid == 0 {
		return false
	}

	userId := uid
	_, ok := ictx.Values().Set(enumz.IctxParamUid, uid)
	if !ok {
		return false
	}

	if vuid > 0 {
		userId = vuid
		_, ok = ictx.Values().Set(enumz.IctxParamVuid, vuid)
		if !ok {
			return false
		}
	}

	acontext.IrisWithRemoteUser(ictx, strconv.FormatUint(userId, 10))

	// 这里设置一定要用  uint32 类型，不然用 GetUint32 取不出来
	if svc.Valid() {
		_, ok = ictx.Values().Set(enumz.IctxParamSvc, uint32(svc))
	}

	return ok
}

// Uid
// @note token 依赖于 UA，使用Chrome调试切换PC/Phone的时候，token 不可复用！！！
func Uid(ictx iris.Context) (svc typez.Svc, uid, vuid uint64, ok bool) {
	var err error
	// UID 获取已经授权的UID。权限验证统一在middleware层处理，controller层不用再多此一举了
	uid, err = ictx.Values().GetUint64(enumz.IctxParamUid)
	if err != nil || uid == 0 {
		return 0, 0, 0, false
	}
	var sv uint32
	sv, err = ictx.Values().GetUint32(enumz.IctxParamSvc)
	if err == nil && sv > 0 {
		svc = typez.Svc(sv)
	}

	if !configz.RequireVuid {
		return svc, uid, 0, true
	}
	vuid, _ = ictx.Values().GetUint64(enumz.IctxParamVuid)
	return svc, uid, vuid, vuid > 0
}

//func ParseUserOpenid(app *aa.App, ictx iris.Context, uid uint64) (appid string, svc typez.Svc, e *ae.Error) {
//	//ctx.URLParam("lastname") == ctx.Request().URL.Query().Get("lastname")
//	openid := ictx.GetHeader(enumz.HeaderOpenid)
//	if openid == "" {
//		e = ae.New(403, "bad openid")
//	}
//	var (
//		err          error
//		applicantUid uint64
//	)
//	appid, svc, applicantUid, err = arpc.New(app).DecodeSvcOpenid(openid)
//	if err != nil {
//		app.Log.Warn(acontext.FromIris(ictx), "openid:%s, %s", openid, err.Error())
//		e = ae.New(403, err.Error())
//		return
//	}
//	if applicantUid > 0 && uid != applicantUid {
//		e = ae.New(403, "bad openid applicant")
//		return
//	}
//	return
//}
