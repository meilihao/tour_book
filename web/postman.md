# postman

## FAQ
1. SSL Error: UNABLE_TO_VERIFY_LEAF_SIGNATURE
nginx配置项ssl_certificate没有包含完整的证书链(证书链顺序: 域名证书->根证书).
1. SSL Error: ELF_SIGNED_CERT_IN_CHAINSELF_SIGNED_CERT_IN_CHAIN
File - Settings - General - 关闭`SSL certificate verification`

## js操作req参数

url: `http://localhost/data={{data}}`

在`Pre-request Script`输入:
```js
var encoded = encodeURIComponent(`
	"path":"xxx"
`);
pm.environment.set("data", encoded);
```