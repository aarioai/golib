# 说明文档

## 用户网页授权
> https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html

1. 第一步：用户同意授权，获取code
2. 第二步：通过code换取网页授权access_token
3. 第三步：刷新access_token（如果需要）
4. 第四步：拉取用户信息(需scope为 snsapi_userinfo)  --> 可以合并为 Code2UserInfo()