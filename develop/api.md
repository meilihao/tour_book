# API & SDK

- [GoWechat 微信平台API](https://github.com/yaotian/gowechat)


## 设计
请求后端处理流程:
1. 程序拿到req, 将程序id(或其他标识)写入 resp header 表明请求已进入后端而不是由中间链路拦截了, 这避免了中间链路其他程序的干扰

请求前端处理流程:
1. 检查http status
1. 处理具体resp // status = 200

resp格式:
```json
{
    "error": "", // 存在"error"属性表示存在错误
    "...": ... // 业务属性
}
```