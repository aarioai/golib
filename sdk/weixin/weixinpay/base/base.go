package base

import (
	"bytes"
	"errors"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/pkg/types"
	"strings"
)

func ToMoney(v *int64) atype.Money {
	if v == nil {
		return 0
	}
	return atype.Money(*v) * atype.Cent
}

func FromMoney(v atype.Money) *int64 {
	cent := v.ToCent()
	return &cent
}

// NewOutTradeNo 生成微信支付 OrderBatch
// 有些支付要求每次提交过去的订单ID都要唯一，为此以 订单 batch ID（或订单ID）  和 retry_at 作为支付依据
// totalAmount 用于判断金额
// attach time.Now().Format("150405") 订单支付超时24小时，所以唯一值保持24小时内即可
// 长度 6~32字符
func NewOutTradeNo(orderBatch uint64, totalAmount atype.Money, attach string) string {
	batch := types.FormatUint(orderBatch)
	total := totalAmount.String()
	var id bytes.Buffer
	id.Grow(len(batch) + 1 + len(total) + 1 + len(attach))
	id.WriteString(batch)
	id.WriteByte('_')
	id.WriteString(total)
	id.WriteByte('_')
	id.WriteString(attach)
	return id.String()
}

func ExtractOutTradeNo(outTradeNo string) (uint64, atype.Money, string, error) {
	a := strings.Split(outTradeNo, "_")
	if len(a) != 3 {
		return 0, 0, "", errors.New("invalid weixin outTradeNo: " + outTradeNo)
	}

	batch := a[0]
	total := a[1]
	attach := a[2]
	orderBatch, err1 := types.ParseUint64(batch)
	totalAmount, err2 := types.ParseInt64(total)
	err := ae.FirstError(err1, err2)
	return orderBatch, atype.Money(totalAmount), attach, err
}
