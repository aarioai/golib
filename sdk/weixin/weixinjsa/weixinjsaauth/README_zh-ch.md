# 小程序授权获取user access_token
> https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-login/code2Session.html

 ```
 小程序 --wx.login() 获取 code--> [用户授权]
       --wx.request(code)-----> [服务端] ---appid+secret+code---> [微信服务]
                                        <--session_key+openid---
 
 小程序登录之后，可以保存登录状态在小程序 local storage。下次进入的时候，可以直接调用接口判断是否登录，也可以通过 wx.request() 携带自定义登录状态
 服务端直接根据登录状态，返回对应openid 和 session_key，减少与微信服务器通信。
 ```