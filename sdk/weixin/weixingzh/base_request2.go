package weixingzh

import (
	"context"
	"github.com/aarioai/airis/pkg/httpc"
	"github.com/aarioai/golib/sdk/weixin/weixingzh/base"
	"io"
)

// 第二种返回数据结构：https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/Nontax_Bill/API_list.html#2.1
// 这里若判断是access token错误，就会重新拉取
func (s *Service) forceGetWithToken2(ctx context.Context, target base.ResultInterface, link string, forceRetry bool, headers ...map[string]string) (err error) {
	if link, err = s.urlWithToken(ctx, link, forceRetry); err != nil {
		return err
	}
	var body []byte

	if _, body, err = httpc.Get(link, nil, headers...); err != nil {
		return err
	}
	return base.ParseResult2(body, target)
}

// 这里若判断是access token错误，就会重新拉取
func (s *Service) getWithToken2(ctx context.Context, target base.ResultInterface, link string, headers ...map[string]string) error {
	return s.forceGetWithToken2(ctx, target, link, false, headers...)
}

func (s *Service) forcePostWithToken2(ctx context.Context, target base.ResultInterface, link string, data io.Reader, forceRetry bool, headers ...map[string]string) (err error) {
	if link, err = s.urlWithToken(ctx, link, forceRetry); err != nil {
		return err
	}
	var body []byte
	if _, body, err = httpc.Post(link, data, headers...); err != nil {
		return err
	}
	return base.ParseResult2(body, target)
}

// 这里若判断是access token错误，就会重新拉取
func (s *Service) postWithToken2(ctx context.Context, target base.ResultInterface, link string, data io.Reader, headers ...map[string]string) error {
	return s.forcePostWithToken2(ctx, target, link, data, false, headers...)
}
