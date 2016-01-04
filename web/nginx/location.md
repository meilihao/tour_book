## NGINX location 在配置中的优先级

[原文](http://www.bo56.com/nginx-location%E5%9C%A8%E9%85%8D%E7%BD%AE%E4%B8%AD%E7%9A%84%E4%BC%98%E5%85%88%E7%BA%A7/)
**nginx有多个`server{...}`时需注意server_name指令和error.log中请求信息的`server: xxx`是否匹配,否则可能将请求匹配到其他server的location上而导致错误**
### location表达式类型

- ~ 表示执行一个正则匹配，区分大小写
- ~* 表示执行一个正则匹配，不区分大小写
- ^~ 表示普通字符匹配。使用前缀匹配。如果匹配成功，则不再匹配其他location。
- = 进行普通字符精确匹配。也就是完全匹配（即不支持正则）。
- @ 它定义一个命名的 location，使用在内部定向时，例如 error_page, try_files


#### location优先级说明

在nginx的location和配置中location的顺序没有太大关系。和location表达式的类型有关。相同类型的表达式，字符串长的会优先匹配。

以下是按优先级排列说明（[官方](http://wiki.nginx.org/HttpCoreModuleChs#location)）：

1. 等号类型（=）的优先级最高。一旦匹配成功，则不再查找其他匹配项。
1. 常规字符串匹配类型（按前缀匹配）。
1. ^~类型表达式。一旦匹配成功，则不再查找其他匹配项。
1. 正则表达式类型（`~ ~*`）的优先级次之。如果有多个location的正则能匹配的话，则使用正则表达式最长的那个。当`~ ~*`表达式相同时，以先后顺序来匹配。

**总结:**

- 不包含正则的 location 在配置文件中的顺序不会影响匹配顺序。而包含正则表达式的 location 会按照配置文件中定义的顺序进行匹配

### location优先级示例

配置项如下:

```shell
location = / {
    # 仅仅匹配请求 /
    [ configuration A ]
}
location / {
    # 匹配所有以 / 开头的请求。
    # 但是如果有更长的同类型的表达式，则选择更长的表达式。
    # 如果有正则表达式可以匹配，则优先匹配正则表达式。
    [ configuration B ]
}
location /documents/ {
    # 匹配所有以 /documents/ 开头的请求。
    # 但是如果有更长的同类型的表达式，则选择更长的表达式。
    # 如果有正则表达式可以匹配，则优先匹配正则表达式。
    [ configuration C ]
}
location ^~ /images/ {
    # 匹配所有以 /images/ 开头的表达式，如果匹配成功，则停止匹配查找。
    # 所以，即便有符合的正则表达式location，也不会被使用
    [ configuration D ]
}
location ~* \.(gif|jpg|jpeg)$ {
    # 匹配所有以 gif jpg jpeg结尾的请求。
    # 但是 以 /images/开头的请求，将使用 Configuration D
    [ configuration E ]
}
```
请求匹配示例

```shell
/ -> configuration A
/index.html -> configuration B
/documents/document.html -> configuration C
/images/1.gif -> configuration D
/documents/1.jpg -> configuration E
```
注意，以上的匹配和在配置文件中定义的顺序无关

## 其他

### 
```
location / {
        return 200 '123';
	# because nginx default content-type is application/octet-stream,
    	# browser will offer to "save the file"...
    	# if you want to see reply in browser, uncomment next line 
        add_header Content-Type text/plain;
}
```
