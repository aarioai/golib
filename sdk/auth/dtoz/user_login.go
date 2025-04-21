package dtoz

type Token struct {
	// @doc https://www.rfc-editor.org/rfc/rfc6749#section-4.2.2
	// 标准参数： access_token, expires_in, scope, state, token_type
	AccessToken string         `json:"access_token"`
	ExpiresIn   int64          `json:"expires_in"` //   time to live in seconds, recommended，时间间隔，而不是timestamp。
	Scope       map[string]any `json:"scope"`      // private:true   私有，或者告诉客户端；不要上报该账号数据
	State       string         `json:"state"`      // 透传回客户端
	TokenType   string         `json:"token_type"` // Bearer  --> 客户端上传header: Authorization: Bearer $access_token

	// 下面是非标准参数
	Conflict     bool           `json:"conflict"` // 是否冲突；一般之前登录后，通过第三方授权登录发生vuid不一致，使用授权登录的vuid取代之前登录的vuid
	RefreshAPI   string         `json:"refresh_api"`
	RefreshToken string         `json:"refresh_token"` // optional
	RefreshTTL   int64          `json:"refresh_ttl"`   // time to live in seconds
	Secure       bool           `json:"secure"`        // 是否安全 --> 如果是通过psid登录的，则就不安全
	ValidateAPI  string         `json:"validate_api"`  // validate client's access token is still available
	Attach       map[string]any `json:"attach"`
}
