package weixingzh

import (
	"context"
	"errors"
	"time"
)

type JsTicket struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"` // 有效时长，一般是7200s
}

// 获取JSSDK ticket
// jsapi_ticket是公众号用于调用微信JS接口的临时票据。ticket用于加强安全性。ticket的有效期目前为2个小时，需定时刷新。建议公众号开发者使用中控服务器统一获取和刷新ticket。
// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#58
// https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/Nontax_Bill/API_list.html#2.1%20%E8%8E%B7%E5%8F%96ticket
type jsTicketResult struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"` // 有效时长，一般是7200s
}

const JsTicketExpiresInterval = int64(10) // 10s

func (r jsTicketResult) CheckError() error {
	if r.Errcode != 0 || r.Ticket == "" {
		return errors.New("check js ticket error")
	}
	return nil
}

// ticketType: jsapi|wx_card ....
func (s *Service) jsTicketCacheName(ticketType string) string {
	return s.cacheJSSDKTicketPrefix + ":" + ticketType
}

// JsTicket  获取JS ticket
func (s *Service) JsTicket(ctx context.Context, ticketType string, force bool) (JsTicket, error) {
	if !force {
		if tkt, err := s.readJsTicketCache(ctx, ticketType); err == nil {
			return tkt, nil
		}
	}
	link := "https://api.weixin.qq.com/cgi-bin/ticket?getticket=" + ticketType
	var tkt jsTicketResult
	if err := s.getWithToken2(ctx, &tkt, link); err != nil {
		return JsTicket{}, err
	}

	ticket := JsTicket{
		Ticket:    tkt.Ticket,
		ExpiresIn: tkt.ExpiresIn,
	}
	s.app.CheckErrors(ctx, s.saveJsTicketCache(ctx, ticketType, ticket))
	return ticket, nil
}

func (s *Service) readJsTicketCache(ctx context.Context, ticketType string) (JsTicket, error) {
	rdb, err := s.rdb()
	if err != nil {
		return JsTicket{}, err
	}
	k := s.jsTicketCacheName(ticketType)
	ticket, _ := rdb.Get(ctx, k).Result()
	if ticket == "" {
		return JsTicket{}, errors.New("ticket missing or expired")
	}
	ttl, _ := rdb.TTL(ctx, k).Result()
	if ttl.Seconds() < 10 {
		return JsTicket{}, errors.New("ticket expired")
	}
	tkt := JsTicket{
		Ticket:    ticket,
		ExpiresIn: int64(ttl.Seconds()),
	}
	return tkt, nil
}

func (s *Service) saveJsTicketCache(ctx context.Context, ticketType string, t JsTicket) error {
	rdb, err := s.rdb()
	if err != nil {
		return err
	}
	k := s.jsTicketCacheName(ticketType)
	ttl := time.Duration(t.ExpiresIn-JsTicketExpiresInterval) * time.Second
	return rdb.SetEx(ctx, k, t.Ticket, ttl).Err()
}
