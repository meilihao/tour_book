# 微信开发

## 前提

###　设置微信公众平台的基本配置

参考 : [基本配置文档](https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421135319&token=&lang=zh_CN)

1. 开发者提交信息后，微信服务器将发送GET请求到填写的服务器地址URL上，开发者通过检验signature对请求进行校验，
若确认此次GET请求来自微信服务器，请原样返回echostr参数内容，则接入生效，成为开发者成功，否则接入失败，会提示比如`Token验证失败`等错误.

2. 启用基本配置中的`服务器配置`,启用后，用户发给公众号的消息以及开发者需要的事件推送，将被微信转发到该URL中.

### FAQ
1. redirect_uri域名与后台配置不一致，错误码10003
订阅号没有权限, 权限位置: 接口权限 - 网页授权 - 网页授权获取用户基本信息 - 网页授权域名, 填入域名(仅域名, 没有http schema); 同时检查程序里RedirectURL里的域名是否与填入域名一致.
<!-- 解决方法: 
		1. 通过"微信认证"
		1. 开通测试账号(仅测试): 开发者工具 - 公众平台测试账号 -->

1. 接口配置信息
[接入概述](https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421135319)
```go
func _VerifyWeixin(c echo.Context) error {
	signature := c.QueryParam("signature")
	timestamp := c.QueryParam("timestamp")
	nonce := c.QueryParam("nonce")
	token := "xxx"

	data := []string{token, timestamp, nonce}
	sort.Strings(data)

	tmp := fmt.Sprintf("%x", sha1.Sum([]byte(strings.Join(data, ""))))
	if tmp != signature {
		LogZ.Error("err sign")

		return nil
	}

	return c.String(200, c.QueryParam("echostr"))
}
```
1. weixin如何关联code和openid
推测: 微信中的用户访问`https://open.weixin.qq.com/connect/qrconnect`时会将其session相关信息放入该请求的header或cookie中,从而关联openid和回调的code.

### 微信公众号基本设置 token验证失败
验证token失败, nginx access.log显示http status = 499, 原因:
阿里云拦截了未在其上备案的请求`curl "http://xxx.com/callback/wechat/mp?signature=8019d9af0be83b3febe982c96b44184dcf52f1c0&echostr=4646699353577078623&timestamp=1539256159&nonce=1709717818"`,返回了备案提示的html.

### 扫二维码即关注并登录
想推送必须关注公众号, 普通流程是:
1. 关注二维码
1. 通过`网页授权获取用户基本信息`绑定uid

其实上面两个步骤可以合二为一, 即通过[生成带参数的二维码](https://mp.weixin.qq.com/wiki?action=doc&id=mp1443433542&t=0.376179226179156)实现, 可参考[微信带参二维码](https://www.jianshu.com/p/084d49ea16bb).