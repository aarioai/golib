package sms

import (
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/golib/sdk/sms/mo"
)

var (
	vericodeLog   = make(chan mo.SmsVericodeLog)
	mongoEntities = []index.Entity{
		mo.SmsVericodeLog{},
		mo.SmsVerifyLog{},
	}
)

func (s *Service) Init(ctx acontext.Context) {
	s.initMongodb(ctx)
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

func (s *Service) initMongodb(ctx acontext.Context) {
	// create tables and indexes
	for _, t := range mongoEntities {
		ae.PanicOn(s.mongo.ORM(t).CreateIndexes(ctx))
	}
}

// @TODO 放入kafka队列
func (s *Service) handleVericodeLog(ctx acontext.Context, t mo.SmsVericodeLog) {
	t.CreatedAt = atype.Now(s.loc)
	_, e := s.mongo.ORM(t).Insert(ctx)
	s.app.Check(ctx, e)
}
