package auth

import (
	"fmt"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/golib/enumz"
	"github.com/kataras/iris/v12"
)

// HeadUserToken detect user token validate
func (s *Service) HeadUserToken(ictx iris.Context) {
	defer ictx.Next()
	_, _, _, ttl, _, e := s.ParseUserAuthorization(ictx)
	if e != nil && e.Code == ae.Unauthorized {
		e = ae.ErrorLoginTimeout
	}
	if e == nil {
		ictx.StatusCode(200)
		return
	}
	ictx.StatusCode(e.Code)
	ictx.Header(enumz.HeaderError, e.Text())
	ictx.Header(enumz.HeaderData, fmt.Sprintf(`{"ttl":%d}`, ttl))
}
