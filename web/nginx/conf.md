### nginx.conf  `location`配置

配置二级目录时,后端server(假设server应用路由都是从`/`根开始)收到的url请求.测试url:

- `localhost/blog`
- `localhost/blog/`
- `localhost/blog/t`

```
location /blog/ {
           proxy_pass http://127.0.0.1:9001/;
        }
/*---
output(推荐,前提是应用中url需支持二级目录):
/
/
/t
---*/
```
```
location /blog {
           proxy_pass http://127.0.0.1:9001/;
        }
/*---
output(error):
//
//
//t
---*/
```
```
location /blog {
           proxy_pass http://127.0.0.1:9001;
        }
/*---
output(前提是应用中url需支持二级目录,且需对url做StripPrefix处理,即截去收到的url的多余的前缀):
/blog/
/blog/
/blog/t
---*/
```
```
location /blog/ {
           proxy_pass http://127.0.0.1:9001;
        }
/*---
output(前提是应用中url需支持二级目录,且需对url做StripPrefix处理):
/blog/
/blog/
/blog/t
---*/
```

### 静态文件匹配

```
location ^~ /static/ {
           root /home/chen/xxx.com/;
        }
/*---
请求将匹配到xxx.com下的static目录.
---*/
```

扩展:

nginx配置下有两个指定目录的执行，root和alias

```shell
location /img/ {
    alias /var/www/image/;
}
#若按照上述配置的话，则访问/img/目录里面的文件时，ningx会自动去/var/www/image/目录找文件
location /img/ {
    root /var/www/image;
}
#若按照这种配置的话，则访问/img/目录下的文件时，nginx会去/var/www/image/img/目录下找文件。]
```
alias是一个目录别名的定义，root则是最上层目录的定义。

还有一个重要的区别是alias后面必须要用“/”结束，否则会找不到文件的;而root则可有可无(**推荐不加**).

