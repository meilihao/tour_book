# http压测

## 如何压测
- 测试http1.1是使用了apache的ab命令
```
$ ab -k -t 180 -c 6 -n 1000000 https://golang.d.openhello.net
```
参数:
- `-k` 使用http的 KeepAlive属性，保持连接处于活动状态
- `-t` 超时时间
- `-c` 并发数
- `-n` 请求数量

- 测试http2.0是使用了nghttp2的h2load命令
```
$ h2load -c 6 -T 180 -n 1000000 https://golang.d.openhello.net
```
参数:
- `-c` 并发数
- `-T` 超时时间限制
- `-n` 请求数量

- 使用[hey](https://github.com/rakyll/hey)

## 其他
### 安装工具
1. ab
ab是apache自带的命令，安装后可以直接使用.

```sh
$ yum install apache # Centos/RedHat
$ apt-get install apache2 # Ubuntu/Debian
```
1. nghttp2
```
$ yum install nghttp2
```