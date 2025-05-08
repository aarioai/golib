package weixingzh

import (
	"context"
	"errors"
	"fmt"
	"github.com/aarioai/airis/pkg/httpc"
	"time"
)

type AccessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
}

const GrantAccessTokenURL = "https://api.weixin.qq.com/sns/oauth2/access_token"

func (s *Service) readTokenCache(ctx context.Context) (AccessToken, error) {
	rdb, err := s.rdb()
	if err != nil {
		return AccessToken{}, err
	}
	ttl, err := rdb.TTL(ctx, s.cacheClientCredential).Result()
	if err != nil || ttl.Seconds() < 1.0 {
		return AccessToken{}, errors.New("access token expired")
	}
	accessToken, err := rdb.Get(ctx, s.cacheClientCredential).Result()
	if err != nil {
		return AccessToken{}, err
	}
	token := AccessToken{
		Token:     accessToken,
		ExpiresIn: int64(ttl.Seconds()),
	}
	return token, nil
}

func (s *Service) saveTokenCache(ctx context.Context, t AccessToken) error {
	rdb, err := s.rdb()
	if err != nil {
		return err
	}
	return rdb.SetEx(ctx, s.cacheClientCredential, t.Token, time.Duration(t.ExpiresIn)*time.Second).Err()
}

// GrantAccessToken 获取基础API access token (grant type: client_credential)
func (s *Service) GrantAccessToken(ctx context.Context, force bool) (AccessToken, error) {
	if !force {
		if tk, err := s.readTokenCache(ctx); err == nil {
			return tk, nil
		}
	}

	var token AccessToken
	statusCode, err := httpc.GetJson(ctx, &token, GrantAccessTokenURL, map[string]string{
		"grant_type": "client_credential",
		"appid":      s.appid,
		"secret":     s.secret,
	})
	if err != nil {
		return AccessToken{}, err
	}
	if token.Token == "" {
		return AccessToken{}, fmt.Errorf("grant access token is empty, status code: %d", statusCode)
	}

	token.ExpiresIn -= 60
	s.app.CheckErrors(ctx, s.saveTokenCache(ctx, token))
	return token, nil
}
