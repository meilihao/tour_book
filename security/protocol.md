# secure protocol

## tls

## https server key的作用
- rsa握手: 解密client 生成的encrypted premaster secret
- dh握手: 签名server dh parameter(x25519是其pub key)

> tls 1.3 握手过程被加密待解密(wireshark2.6.2还不行)

## ssh v2
