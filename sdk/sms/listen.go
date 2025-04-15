package sms

import (
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/golib/sdk/sms/mo"
)

var (
	vericodeLog = make(chan mo.SmsVericodeLog)
)

func (s *Service) Init(ctx acontext.Context) {
	go func() {
		for {
			select {
			case log := <-vericodeLog:
				s.handleVericodeLog(ctx, log)
			case <-ctx.Done():
				return
			}
		}
	}()
}

// @TODO 放入kafka队列
func (s *Service) handleVericodeLog(ctx acontext.Context, t mo.SmsVericodeLog) {
	_, e := s.mongo.ORM(t).Insert(ctx)
	s.app.Check(ctx, e)
}
