# secure protocol

## tls

## https server key的作用
- rsa握手: 解密client 生成的encrypted premaster secret
- dh握手: 签名server dh parameter(x25519是其pub key)

## 部署https
[acme.sh](https://github.com/Neilpang/acme.sh)

## ssh v2

## FAQ
### 查看tls 1.3流量
wireshark:
打开【编辑】-【首选项】-【Protocols】->【SSL】，然后设置 【(Pre)-Master-Secret log filename】, 重启wireshark.

浏览器(设置SSLKEYLOGFILE):
```sh
$ env SSLKEYLOGFILE=/home/chen/.config/tls.log ./firefox
```