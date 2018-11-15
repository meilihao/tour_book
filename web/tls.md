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