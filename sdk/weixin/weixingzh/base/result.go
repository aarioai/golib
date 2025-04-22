package base

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type ResultInterface interface {
	CheckError() error // 不能直接通过unmarshal err==nil 就断定解析成功了，还需要独立判断
}

type Error struct {
	Code int    `json:"errcode"`
	Msg  string `json:"errmsg"`
}

func ParseResult2(body []byte, target ResultInterface) error {
	if err := json.Unmarshal(body, &target); err != nil {
		return err
	}
	if target == nil {
		return fmt.Errorf("unmarshal response failed: %s", string(body))
	}
	return target.CheckError()
}

func ParseResult(body []byte, target ResultInterface) (*Error, error) {
	if bytes.Index(body, []byte("errcode")) <= 0 {
		return nil, ParseResult2(body, target)
	}
	werr := &Error{}
	if err := json.Unmarshal(body, &werr); err != nil {
		werr.Code = 500
		werr.Msg = err.Error()
	}
	return werr, fmt.Errorf("parse response error code: %d, msg: %s", werr.Code, werr.Msg)
}
