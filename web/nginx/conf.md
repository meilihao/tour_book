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

