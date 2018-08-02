# dig

## 描述

Domain Information Groper,域名查询工具(DNS lookup utility),可以用来测试域名系统工作是否正常

## 语法格式

```
grep [OPTIONS] [查询内容] [querytype]
```

## 选项

- @<DNS服务器地址> : 指定进行域名解析的域名服务器
- -4 : 使用IPv4
- -6 : 使用IPv6
- -b <IP>: 当本机具有多个IP地址，指定使用本机的哪个IP地址向域名服务器发送域名查询请求
- -t <类型> : 指定要查询的DNS数据类型
- -x <IP> : 执行逆向域名查询,可以查询IP地址到域名的映射关系
- +nssearch : 查找域名的权威域名服务器
- +short : 提供简要应答
- +trace : 显示DNS的整个查询过程
- +tcp : 使用tcp查询dns,默认是udp

querytype : A/AAAA/PTR/MX/ANY等值，默认是查询A记录

常见的DNS记录类型(querytype):
- A : 地址记录（Address），返回域名指向的IPv4地址
- AAAA : 地址记录（Address），返回域名指向的IPv6地址
- NS : 域名服务器记录（Name Server），返回管理下一级域名信息的服务器地址.该记录只能设置为域名，不能设置为IP地址
- MX : 邮件记录（Mail eXchange），返回接收电子邮件的服务器地址
- CNAME : 规范名称记录（Canonical Name），返回另一个域名，即当前查询的域名是另一个域名的跳转
- PTR : 逆向查询记录（Pointer Record），只用于从IP地址查询域名

一般来说，为了服务的安全可靠，至少应该有两条NS记录，而A记录和MX记录也可以有多条，这样就提供了服务的冗余性，防止出现单点失败.

DNS解析完整过程:
1. 用户向本地DNS发出解析请求
2. 本地DNS向根服务器（根服务器地址在本地有个静态列表）请求谁是该域名的顶机DNS（并把结果缓存在本地）
3. 根DNS告诉本地DNS谁是该域的顶级DNS
4. 本地DNS向顶级DNS请求谁是权威DNS（并缓存在本地）
5. 应答谁是该域权威DNS
6. 本地DNS向权威DNS请求域名解析（并把结果缓存在本地）
7. 权威服务器应答
8. 本地DNS把结果告诉用户

## 其他

- [可视化的DNS查询工具 : dnsgraph，绘制根域名到指定域名的所有可能路径](http://ip.seveas.net/dnsgraph/)

## 结果

- 第一段 : 查询参数和统计(Got answer)
- 第二段 : 查询内容(QUESTION SECTION),`IN`后面的`A`指`A记录`
- 第三段 : DNS服务器的答复(ANSWER SECTION),`IN`前面的数字是`TTL值`,表示缓存时间，即`n`秒之内不用重新查询.
- 第四段 : AUTHORITY SECTION,显示被查询域名的NS记录（Name Server的缩写），即哪些服务器负责管理该域名的DNS记录;没有AUTHORITY SECTION说明结果来自被查询DNS服务器的缓存，而不是权威域名服务器
- 第五段 : 是AUTHORITY SECTION中域名服务器的IP地址，这是随着前一段一起返回的
- 第六段 : 被查询的DNS服务器的一些传输信息,比如`SERVER: 127.0.1.1#53`,本机的DNS服务器是`127.0.1.1`，查询端口是53（DNS服务器的默认端口）;`MSG SIZE  rcvd`,DNS服务器回复的字节数.

## 例
```
$ dig www.baidu.com
$ dig # 查询root dns
$ dig . # 查询本地DNS服务器
$ dig @223.5.5.5 www.baidu.com A
```

## 参考
- [DNS 查询之 10 个 dig 详细例子](https://www.tuicool.com/articles/MvMrI36)
