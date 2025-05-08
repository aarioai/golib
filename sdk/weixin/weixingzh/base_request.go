package weixingzh

import (
	"context"
	"github.com/aarioai/airis/pkg/httpc"
	"github.com/aarioai/golib/sdk/weixin/weixingzh/base"
	"io"
	"strings"
)

// 第一种返回数据结构：

func (s *Service) urlWithToken(ctx context.Context, link string, force bool) (string, error) {
	token, err := s.GrantAccessToken(ctx, force)
	if err != nil {
		return "", err
	}
	b := '?'
	if strings.IndexByte(link, '?') > 0 {
		b = '&'
	}
	return link + string(b) + "access_token=" + token.Token, nil
}

// 这里若判断是access token错误，就会重新拉取
func (s *Service) forceGetWithToken(ctx context.Context, target base.ResultInterface, link string, forceRetry bool, headers ...map[string]string) (err error) {
	if link, err = s.urlWithToken(ctx, link, forceRetry); err != nil {
		return err
	}
	var body []byte

	if _, body, err = httpc.Get(ctx, link, nil, headers...); err != nil {
		return err
	}
	var werr *base.Error
	if werr, err = base.ParseResult(body, target); err == nil {
		return nil
	}

	// 没重试过， 且是 access token 错误，就重试一次
	if werr != nil && !forceRetry && (werr.Code == -1 || werr.Code == 40001 || werr.Code == 40014) {
		return s.forceGetWithToken(ctx, target, link, true, headers...)
	}
	return err
}

// 这里若判断是access token错误，就会重新拉取
func (s *Service) getWithToken(ctx context.Context, target base.ResultInterface, link string, headers ...map[string]string) error {
	return s.forceGetWithToken(ctx, target, link, false, headers...)
}

func (s *Service) forcePostWithToken(ctx context.Context, target base.ResultInterface, link string, data io.Reader, forceRetry bool, headers ...map[string]string) (err error) {
	if link, err = s.urlWithToken(ctx, link, forceRetry); err != nil {
		return err
	}
	var body []byte
	if _, body, err = httpc.Post(ctx, link, data, headers...); err != nil {
		return err
	}
	var werr *base.Error
	if werr, err = base.ParseResult(body, target); err == nil {
		return nil
	}
	// 没重试过， 且是 access token 错误，就重试一次
	if werr != nil && !forceRetry && (werr.Code == -1 || werr.Code == 40001 || werr.Code == 40014) {
		return s.forcePostWithToken(ctx, target, link, data, true, headers...)
	}
	return err
}

// 这里若判断是access token错误，就会重新拉取
func (s *Service) postWithToken(ctx context.Context, target base.ResultInterface, link string, data io.Reader, headers ...map[string]string) error {
	return s.forcePostWithToken(ctx, target, link, data, false, headers...)
}
