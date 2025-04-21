package dtoz

import "github.com/aarioai/airis/aa/atype"

type UserOpenidResponse struct {
	Openid    string                `json:"openid"`
	ExpiresIn atype.DurationSeconds `json:"expires_in"`
}
