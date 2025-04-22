package weixingzh

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
)

// UserInfoResult
// https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html#UinonId
type UserInfoResult struct {
	Subscribe interface{} `json:"subscribe"`
	Openid    string      `json:"openid"`
	Language  string      `json:"language"`

	SubscribeTime any    `json:"subscribe_time"`
	Unionid       string `json:"unionid"`
	Remark        string `json:"remark"`
	Groupid       any    `json:"groupid"`
	TagidList     []any  `json:"tagid_list"`

	SubscribeScene string `json:"scene"`
	QrScene        any    `json:"qr_scene"`     // 二维码扫码场景（开发者自定义）
	QrSceneStr     string `json:"qr_scene_str"` // 二维码扫码场景描述（开发者自定义）
}

type UserInfoList struct {
	UserInfoList []UserInfoResult `json:"user_info_list"`
}

const (
	BatchGetUserInfoUrl = "https://api.weixin.qq.com/cgi-bin/user/info/batchget"
)

func (r UserInfoList) CheckError() error {
	if len(r.UserInfoList) == 0 {
		return errors.New("parse batch user info failed")
	}
	for _, v := range r.UserInfoList {
		if v.Openid == "" {
			return errors.New("parse batch user info failed, openid empty")
		}
	}
	return nil
}

// BatchUserinfo 批量获取微信用户资料
func (s *Service) BatchUserinfo(ctx context.Context, openids []string) ([]UserInfoResult, error) {
	users := make([]map[string]string, len(openids))
	for i, openid := range openids {
		users[i] = map[string]string{
			"openid": openid,
			"lang":   "zh_CN",
		}
	}
	d := map[string][]map[string]string{
		"user_list": users,
	}

	b, _ := json.Marshal(d)
	var list UserInfoList
	e := s.postWithToken(ctx, &list, BatchGetUserInfoUrl, bytes.NewReader(b), nil)
	if e != nil {
		return nil, e
	}
	return list.UserInfoList, nil
}
