# directive

## try_files
```
try_files file [...] (uri|@tag|=code) # 至少需要两个参数
```

按顺序检查文件是否存在，并返回第一个存在的文件.如果所有文件都不存在，则**内部重定向到最后一个参数(该参数的指向必须存在,否则可能出现循环)所指位置"**.

其受 root 和 **index** 语句影响

### 执行顺序
- 如果`$uri`是文件,一切ok.
- 如果`$uri`是目录(即`/`结尾):
伪代码:
```golang
// $uri是目录
status:=0

for _,v:=range $index{
    if IsFile($uri+v){ // true
        $uri=$uri+v // 这也是try_files的第一个值通常是$uri的原因,因为该文件肯定存在
        status=404

        n:=len($try_files)
        for j,vv:=range $try_files{
            if j==n-1{
                goto Location # nginx内部重定向
            }

            if IsFile(strings.TrimRight(vv,"/")){
                status = 200
                
                break
            }
        }
    }else{
        if IsDir(strings.TrimRight($uri,"/")){
            status=403

            continue
        }else{ // No such file or directory
            status=404

            break
        }
    }
}
```

举例: `curl https://golang.d.openhello.net/blog/`且文件`/var/tmp/golang/blog/index.html`存在.
```sh
// resp status 200
// nginx conf
location / {
      add_header  Content-Type 'text/html; charset=utf-8';
      
      set $t "${uri}.html";
      if ($uri ~ /$) {
         set $t "${uri}index.html";
      }

      root /var/tmp/golang;
      try_files $uri $t =404;
    }

// nginx.error.log
2017/08/20 19:25:52 [notice] 11424#11424: *1 "/$" matches "/blog/", client: 219.82.134.60, server: golang.d.openhello.net, request: "GET /blog/ HTTP/1.1", host: "golang.d.openhello.net"
2017/08/20 19:25:52 [notice] 11424#11424: *1 "/$" does not match "/blog/index.html", client: 219.82.134.60, server: golang.d.openhello.net, request: "GET /blog/ HTTP/1.1", host: "golang.d.openhello.net"

// strace log
write(4, "2017/08/20 19:25:52 [notice] 114"..., 186) = 186
stat("/var/tmp/golang/blog/index.html", {st_mode=S_IFREG|0755, st_size=26502, ...}) = 0
gettid()                                = 11424
write(4, "2017/08/20 19:25:52 [notice] 114"..., 203) = 203
stat("/var/tmp/golang/blog/index.html", {st_mode=S_IFREG|0755, st_size=26502, ...}) = 0
open("/var/tmp/golang/blog/index.html", O_RDONLY|O_NONBLOCK) = 17
fstat(17, {st_mode=S_IFREG|0755, st_size=26502, ...}) = 0
pread(17, "\n<!DOCTYPE html>\n<html>\n<head>\n\t"..., 26502, 0) = 26502
```

执行顺序:
1. 依次检查目录下index directive列举的文件,如果index存在,$uri=$uri+index,再执行try_files,$uri存在.

```sh
// resp status 404
// nginx conf
location / {
      add_header  Content-Type 'text/html; charset=utf-8';

      set $t "${uri}.html";
      if ($uri ~ /$) {
         set $t "${uri}index.html";
      }

      root /var/tmp/golang;
      try_files $t =404;
    }

// nginx.error.log
2017/08/20 19:46:13 [notice] 11463#11463: *1 "/$" matches "/blog/", client: 219.82.134.60, server: golang.d.openhello.net, request: "GET /blog/ HTTP/1.1", host: "golang.d.openhello.net"
2017/08/20 19:46:13 [notice] 11463#11463: *1 "/$" does not match "/blog/index.html", client: 219.82.134.60, server: golang.d.openhello.net, request: "GET /blog/ HTTP/1.1", host: "golang.d.openhello.net"
2017/08/20 19:46:13 [error] 11463#11463: *1 open() "/etc/nginx/html404" failed (2: No such file or directory), client: 219.82.134.60, server: golang.d.openhello.net, request: "GET /blog/ HTTP/1.1", host: "golang.d.openhello.net"

// strace log
write(4, "2017/08/20 19:46:13 [notice] 114"..., 186) = 186
stat("/var/tmp/golang/blog/index.html", {st_mode=S_IFREG|0755, st_size=26502, ...}) = 0
gettid()                                = 11463
write(4, "2017/08/20 19:46:13 [notice] 114"..., 203) = 203
stat("/var/tmp/golang/blog/index.html.html", 0x7ffdc7eb29d0) = -1 ENOENT (No such file or directory)
open("/etc/nginx/html404", O_RDONLY|O_NONBLOCK) = -1 ENOENT (No such file or directory)
```

执行顺序:
1. 依次检查目录下index directive列举的文件,如果index存在,$uri=$uri+index,再执行try_files,$t不存在.

## module
- return : ngx_http_rewrite_module