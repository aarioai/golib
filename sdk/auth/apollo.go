package auth

import (
	"context"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/golib/sdk/auth/dtoz"
	"github.com/aarioai/golib/typez"
	"time"
)

var apolloUpdatedAt int64
var apolloDailyUpdated time.Time

// GrantApollo 通过设备码登录（不安全）
// @TODO 还未完成
func (s *Service) GrantApollo(ctx context.Context, di typez.DeviceInfo) (*dtoz.Token, *ae.Error) {
	if !di.Valid() {
		return nil, ae.ErrorPreconditionRequired
	}
	apollo := di.Encode('B')
	var updatedAt int64
	now := time.Now()
	var expiresIn int64
	if apolloUpdatedAt > 0 && now.Before(apolloDailyUpdated) {
		updatedAt = apolloUpdatedAt
		expiresIn = apolloDailyUpdated.Unix() - now.Unix()
	} else {
		var ok bool
		updatedAt, ok = s.h.LoadApolloUpdatedAt(ctx)
		if !ok {
			updatedAt = now.Unix()
		}
		apolloUpdatedAt = updatedAt
		apolloDailyUpdated = now.AddDate(0, 0, 1)
		expiresIn = 86400 // one day
	}

	token := dtoz.Token{
		AccessToken:  "",
		ExpiresIn:    expiresIn,
		Scope:        nil,
		State:        "",
		TokenType:    "",
		Conflict:     false,
		RefreshAPI:   "",
		RefreshToken: "",
		RefreshTTL:   0,
		Secure:       false,
		ValidateAPI:  "",
		Attach: map[string]any{
			"apollo":     apollo,
			"updated_at": updatedAt, // 早于这个日期的缓存，除了access_token以外，其余全部清空
		},
	}
	return &token, nil
}
