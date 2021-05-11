# telnet
通过telnet协议与远程主机通信或探测其端口是否开启, 默认使用23端口.

建议使用`openssl s_client -connect example.com:443 [-showcerts]`代替.

## 格式

    telnet 目标主机

## example
```
# telnet 192.168.0.121 22
```