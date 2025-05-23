package enumz

const (
	HeaderApollo         = "X-Apollo"
	HeaderAuthorization  = "Authorization" // header 里面
	HeaderCSRFToken      = "X-CSRF-token"
	HeaderData           = "X-Data" // HEAD 下附加的数据信息，json格式
	HeaderDebugToken     = "X-Debug-Token"
	HeaderError          = "X-Error" // HEAD下，服务端返回客户端的header错误信息
	HeaderMMCFingerprint = "X-MMC-Fingerprint"
	HeaderMockToken      = "X-Mock-Token"
	HeaderOpenid         = "X-Openid"
	HeaderStringify      = "X-Stringify"

	// 注意：手机端IP可能会经常变换，因此不要过度依赖
	// @warn 尽量不要通过自定义header传参，因为可能某个web server会基于安全禁止某些无法识别的header
	ParamApollo = "apollo" //  阿波罗计划，设备信息

	ParamAccessToken          = "access_token"
	ParamWeixinToken          = "wx_access_token"
	ParamGrantType            = "grant_type" // token grant type
	ParamRefreshToken         = "code"       // 配合 grant type使用
	ParamDbgWeixinAccessToken = "dbg_wx_token"
	ParamDbgAccessToken       = "dbg_access_token"
	ParamTaWeixinCode         = "code"
	ParamTaWeixinState        = "state"
	ParamAppid                = "appid" // 不用header，这样对于一些git hook等兼容性更强
	ParamSign                 = "sign"  // 不用header，这样对于一些git hook等兼容性更强
	ParamCallback             = "callback"

	ParamLogout = "logout"
	ParamOpenid = "openid"

	// 用户分享 这个是 f=$mission_b36(vuid)
	ParamAg = "ag"

	IctxClientDebug                = "ClientDebug"
	IctxClientMock                 = "ClientMock"
	IctxParamFingerprintServerTime = "FingerprintServerTime"
	IctxParamUid                   = "Uid"
	IctxParamVuid                  = "Vuid"
	IctxParamSvc                   = "Svc"
	IctxUserTokenTTL               = "UserTokenTTL"
	IctxUserToken                  = "UserToken"

	IctxParamTokenByT3 = "TokenByT3"
	////IctxParamAImageXhost = "AImageXhost"
	////IctxParamEnv        = "Env"
	//IctxParamAdminLevel = "AdminLevel" // 是否是私有账号， 不可被统计追踪
	//IctxParamTrace      = "Trace"      // 是否跟踪
	//IctxParamTab        = "Tab"        // tab
	//ViewKeyDeviceInfo   = "DeviceInfo"
)
