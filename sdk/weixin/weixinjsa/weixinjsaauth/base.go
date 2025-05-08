package weixinjsaauth

type Response struct {
	SessionKey string `json:"session_key"`
	Unionid    int64  `json:"unionid"`
	Openid     string `json:"openid"`
}

type ResponseWithError struct {
	Response Response
	Errcode  int32  `json:"errcode"`
	Errmsg   string `json:"errmsg"`
}
