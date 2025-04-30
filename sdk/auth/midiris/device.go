package midiris

import (
	"errors"
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/aa/httpsvr/request"
	"github.com/aarioai/golib/enumz"
	"github.com/aarioai/golib/typez"
	"github.com/kataras/iris/v12"
	"github.com/mssola/useragent"
	"sync"
)

var deviceInfoMtx sync.Mutex

func ClientIP(ictx iris.Context) atype.IP {
	ipr := acontext.ClientIP(ictx)
	return atype.ToIP(ipr)
}

func deviceInfo(ictx iris.Context) (*typez.DeviceInfo, bool) {
	ua, err := ictx.Values().GetUint16("DeviceInfo.UA")
	if err != nil {
		return nil, false
	}
	var d typez.DeviceInfo
	d.IP = acontext.ClientIP(ictx)
	d.UA, _ = enumz.NewUA(ua)
	d.PSID = ictx.Values().GetString("DeviceInfo.PSID")
	d.UDID = ictx.Values().GetString("DeviceInfo.UDID")
	d.OAID = ictx.Values().GetString("DeviceInfo.OAID")
	d.UUID = ictx.Values().GetString("DeviceInfo.UUID")
	d.Model = ictx.Values().GetString("DeviceInfo.Model")
	d.DpWidth, _ = ictx.Values().GetUint16("DeviceInfo.DpWidth")
	d.DpHeight, _ = ictx.Values().GetUint16("DeviceInfo.DpHeight")
	d.DipWidth, _ = ictx.Values().GetUint16("DeviceInfo.DipWidth")
	d.OS = ictx.Values().GetString("DeviceInfo.Os")
	d.Agent = ictx.Values().GetString("DeviceInfo.Agent")
	d.Lang = ictx.Values().GetString("DeviceInfo.Lang")
	d.Info = ictx.Values().GetString("DeviceInfo.Info")
	return &d, true
}

func setDeviceInfo(ictx iris.Context, d typez.DeviceInfo) {
	ictx.Values().Set("DeviceInfo.UA", d.UA.Uint16())

	if d.PSID != "" {
		ictx.Values().Set("DeviceInfo.PSID", d.PSID)
	}
	if d.UDID != "" {
		ictx.Values().Set("DeviceInfo.UDID", d.UDID)
	}
	if d.OAID != "" {
		ictx.Values().Set("DeviceInfo.OAID", d.OAID)
	}
	if d.UUID != "" {
		ictx.Values().Set("DeviceInfo.UUID", d.UUID)
	}
	if d.Model != "" {
		ictx.Values().Set("DeviceInfo.Model", d.Model)
	}
	if d.DpWidth > 0 {
		ictx.Values().Set("DeviceInfo.DpWidth", d.DpWidth)
	}
	if d.DpHeight > 0 {
		ictx.Values().Set("DeviceInfo.DpHeight", d.DpHeight)
	}
	if d.DipWidth > 0 {
		ictx.Values().Set("DeviceInfo.DipWidth", d.DipWidth)
	}
	if d.OS != "" {
		ictx.Values().Set("DeviceInfo.OS", d.OS)
	}
	if d.Agent != "" {
		ictx.Values().Set("DeviceInfo.Agent", d.Agent)
	}
	if d.Lang != "" {
		ictx.Values().Set("DeviceInfo.Lang", d.Lang)
	}
	if d.Info != "" {
		ictx.Values().Set("DeviceInfo.Info", d.Info)
	}
}
func parseUserAgent(ictx iris.Context) (ua enumz.UA, model, os, agent, raw string) {
	if raw = ictx.GetHeader("User-Agent"); raw == "" {
		return
	}
	uag := useragent.New(raw)
	ua = enumz.UserAgentToUA(uag)
	model = uag.Model()
	os = uag.OS()
	browserName, browserVer := uag.Browser()
	agent = browserName + " " + browserVer
	return
}
func parseDeviceInfo(ictx iris.Context) (*typez.DeviceInfo, error) {
	apollo := request.QueryWild(ictx, enumz.ParamApollo)
	if apollo == "" {
		return nil, errors.New("miss device info")
	}
	di, err := typez.DecodeDeviceInfo(apollo)
	if err != nil {
		return nil, err
	}
	di.IP = acontext.ClientIP(ictx)
	return di, nil
}

// DeviceInfo 通过HTML访问，middleware.SetCommonViewData 会执行一次这个
// 如果是访问HTML，同时请求API，API里面也有获取 GetDeviceInfo() ，那么则相当于两个独立请求执行这个，debug时候输出多个是此原因
func DeviceInfo(ictx iris.Context) typez.DeviceInfo {
	deviceInfoMtx.Lock()
	defer deviceInfoMtx.Unlock()
	dp, ok := deviceInfo(ictx)
	if ok {
		return *dp
	}
	ua, model, os, agent, info := parseUserAgent(ictx)
	devi, err := parseDeviceInfo(ictx)
	if err != nil {
		devi = &typez.DeviceInfo{}
	}
	if ua.Valid() {
		devi.UA = ua
		devi.Model = model
		devi.OS = os
		devi.Agent = agent
		devi.Info = info
	}

	setDeviceInfo(ictx, *devi)
	return *devi
}
