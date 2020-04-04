# FAQ
## error
1. authentication handshake failed: x509: certificate has expired or is not yet valid
	服务器时间错误, 可使用`ntpdate ${time_server}`纠正

### aes cbc解密后缺少开头的部分数据, 其他数据正常
aes cbc加密时仅将加密后的数据一起输出, 解密时将数据的前16B当做iv处理, 导致数据前半部有缺失.

因此aes cbc加解密时对iv的处理需要一致.