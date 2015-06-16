### `error: the HTTP rewrite module requires the PCRE library`

```
sudo apt-get install libpcre3 libpcre3-dev
sudo dnf -y install pcre-devel
```

### `error: SSL modules require the OpenSSL library`

```
sudo apt-get install libssl-dev
sudo dnf -y install  openssl-devel
```

### `nginx: [emerg] unknown log format "main" in /usr/local/nginx/conf/nginx.conf`

打开access_log选项后未启用日志格式main，启用log_format	main(nginx默认日志格式)即可.

### `An error occurred`

可能情况:

`proxy_pass`指定的地址不可用.