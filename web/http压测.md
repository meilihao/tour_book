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

### neokylin安装php报缺httpd-mmn, 但Google httpd-mmn没找到具体package
参考: [Run-Time Dependencies](https://fedoraproject.org/wiki/PackagingDrafts/ApacheHTTPModules#Run-Time_Dependencies)

使用`yumdownloader httpd-mmn`发现就是httpd, 再用`dnf install httpd-2.4.34-18...aarch64.rpm`, 最后再用`dnf install php`即可(也可先用yumdownloader下载再安装).

> [Module Magic Number (MMN)](http://httpd.apache.org/docs/2.0/glossary.html), 其实httpd-mmn就是httpd, 用于确保mod_xxx仅与提供相同二进制模块接口的 httpd 包一起使用