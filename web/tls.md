# tls

## FAQ
### WebSocket connection failed: Error in connection establishment: net::ERR_CERT_AUTHORITY_INVALID
使用wss连接server时, chrome报该错: 由于证书是自签名的，浏览器认为不安全，拦截了请求.

解决(假设被拦截的域名是`wss://www.wss.com`):
1. 新开一个Tab页面, 访问该域名`https://www.wss.com`
2. 你会发现浏览器告警："您的连接不是私密连接.......", 此时点`高级`, 继续点击`继续前往 www.wss.com（不安全）`
3. 页面会提示"...HTTP ERROR 400"，不用管，这是因为用HTTPS协议访问WSS服务所致，不用管，到这里就可以解决错误了.
4. 刷新错误网页,重新尝试即可.

ps:
Firefox 和 Chrome 浏览器对SSL证书拒绝的错误提示是不一样的,Firefox报`Firefox 无法建立到 wss://www.wss.com/ 服务器的连接`,但是解决步骤完全一样样

### https证书的extendedKeyUsage
在生成 HTTPS 服务器端证书时，注意要加上秘钥扩展 extendedKeyUsage = serverAuth,1.3.6.1.5.5.7.3.1, 这样生成的秘钥才可以用来在服务器和客户端之间进行认证，不然会提示鉴权失败. 在生成 HTTPS 服务器端证书时，需要填写 Common Name (e.g. server FQDN or YOUR name), 即访问服务的域名信息，如果有很多子域名，可以用`*`代替，如`*.test.com`. 如果需要客户端和服务器端双向认证， 在生成客户端证书时， 注意要加上秘钥扩展 extendedKeyUsage = clientAuth,1.3.6.1.5.5.7.3.2.