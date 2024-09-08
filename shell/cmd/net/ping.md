# ping

> [hping:一个具有可嵌入tcl脚本功能的 TCP/IP包伪造工具](https://dev.to/cylon/linuxwang-luo-ming-ling-ji-jin-3c10)

## 描述

ping过ICMP（Internet控制消息协议）工作, 可用来测试本机与目标主机是否联通、联通速度如何、稳定性如何.

## 选项
- -c count :   ping指定次数后停止ping
- -I : 指定源ip, windows则使用`-S`
- -l : 指定网卡
- -i interval : 设定间隔几秒发送一个ping包，默认一秒ping一次
- -w : 以秒为单位, 整个程序会话的deadline时间. 使用`-w`时`-c`将被忽略
- -W : 以秒为单位, 设置单个ping resp的超时时间

## 例子
```bash
ping -c 5 -i 0.6 qq.com # 每隔0.6秒ping一次，一共ping 5次
```