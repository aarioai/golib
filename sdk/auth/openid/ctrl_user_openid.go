package openid

//func (s *Service) GrantUserOpenid(ictx iris.Context) {
//	defer ictx.Next()
//	_, resp, _ := httpsvr.New(ictx)
//
//	uid, _, _ := broker.CtxGetUid(ictx)
//
//	openid, ttl, err := arpc.New(c.app).EncodeSvcOpenid(conf.Svc, uid)
//	if err != nil {
//		resp.WriteErr(err)
//		return
//	}
//	type uo struct {
//		Openid    string `json:"openid"`
//		ExpiresIn int    `json:"expires_in"`
//	}
//	resp.Write(uo{openid, int(ttl.Seconds())})
//}
