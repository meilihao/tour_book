# ncat
ref:
- [ncat 取代了rhel 7的 netcat](https://docs.redhat.com/zh-cn/documentation/red_hat_enterprise_linux/7/html/networking_guide/sec-managing_data_using_the_ncat_utility)

ncat是Nmap项目中netcat的版本, 是nc的一个增强版本，提供了更多的功能，特别是在安全和加密方面, **推荐使用**.

## 选项
- -z : 指定要扫描的端口范围

## example
```bash
# ncat [-u] -l 8080 # 监听8080, `-u`是使用udp
# ncat 10.0.11.60 8080 # 连接10.0.11.60:8080
# ncat -l 8080 > outputfile
# ncat -l 10.0.11.60 8080 < inputfile
# ncat -l 8080 < inputfile # 同上, 传输方向相反
# ncat -l 10.0.11.60 8080 > outputfile
# ncat -l --proxy-type http localhost 8080 # 在 localhost 端口 8080 中创建 HTTP 代理服务器
# ncat -z 10.0.11.60 80-90 # 端口扫描
# ncat -e /bin/bash -k -l 8080 --ssl # 基于ssl的通信
# ncat --ssl 10.0.11.60 8080 
```

# nc

## 描述

nc是netcat的简写，有着网络界的瑞士军刀美誉。因为它短小精悍、功能实用，被设计为一个简单、可靠的网络工具. **推荐使用ncat**.

参考:
- [nc命令用法举例](https://www.cnblogs.com/nmap/p/6148306.html)

## 选项

- -l : 指定nc处于侦听模式, 意味着nc被当作server，侦听并接受连接
- -s : 指定发送数据的源IP地址，适用于多网卡的情况
- -u : 指定nc使用UDP协议，默认为TCP
- -v : 输出交互或出错信息，新手调试时尤为有用
- -w : 超时秒数，后面跟数字
- -z : 表示zero，表示扫描时不发送任何数据

## example
```bash
$ telnet 10.0.1.161 9999 # telnet作为client, telnet是运行于tcp协议的
$ nc -l 8080 # tcp监听8080
$ nc -ul 8080 # udp监听8080
$ nmap 10.0.1.161 -p9999 # nmap作为client
$ nc -vz 127.0.0.1 8080 # 测试tcp是否可连接
$ nc -uz 127.0.0.1 8080 # 测试udp是否可连接. 检测UDP端口的时候不会立即返回测试结果，可能需要等待几秒钟
$ nc -vz -w 2 10.0.1.161 9999 # nc作为client
$ nc -vz -w 2 10.0.1.161 9998-9999 # 检查连续的两个端口
$ nc -l -p 8888 -c "nc 192.168.19.153 22" # 8888转发到192.168.19.153:22
```

> ubuntu 24.04的nc没有`-c`, 但其ncat有

client:
```bash
$ dd if=/dev/zero bs=9000 count=1000 > /dev/tcp/$target_host/$port
$ cat < /dev/urandom > /dev/tcp/$target_host/$port
```

测试unix-socket:
```bash
nc -lU /tmp/socket_test # server
nc -U /tmp/socket_test # client
```