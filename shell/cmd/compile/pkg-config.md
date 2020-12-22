# pkg-config

## 选项

- --list-all 输出已知lib

## FAQ
### OpenSSL >= 1.0.1e and associated developement headers required
编译freeswitch-1.8.2时碰到报错:
```txt
...
checking for openssl... yes
checking openssl_CFLAGS... -I/usr/local/include
checking openssl_LIBS... -L/usr/local/lib -lssl -lcrypto
  adding "-DHAVE_OPENSSL" to SWITCH_AM_CFLAGS
checking for SSL_CTX_set_tlsext_use_srtp in -lssl... no
configure: error: OpenSSL >= 1.0.1e and associated developement headers required
```

可明明openssl已经安装:
```sh
$ sudo dpkg -l|grep openssl
ii  openssl                                                1.1.1-1                              amd64        Secure Sockets Layer toolkit - cryptographic utility
```

检查freeswith的configure, 发现它使用`pkg-config`检查openssl的设置.
```sh
$ pkg-config --libs  openssl
-L/usr/local/lib -lssl -lcrypto
$ pkg-config --cflags   openssl
-I/usr/local/include
$ pkg-config --debug  openssl # 查找openssl.pc文件位置
$ locate openssl.pc
/usr/lib/x86_64-linux-gnu/pkgconfig/openssl.pc
/usr/local/lib/pkgconfig/openssl.pc
```

发现`/usr/local/lib/pkgconfig/openssl.pc`,应该是以前自编译openssl生成的,当前使用了系统自带的openssl, 它比当前的`/usr/lib/x86_64-linux-gnu/pkgconfig/openssl.pc`更优先了, 删除相关文件即可:
```sh
$ cd /usr/local/lib/pkgconfig
$ sudo rm openssl.pc libssl.pc libcrypto.pc
```

> 其实也可通过环境变量`PKG_CONFIG_PATH`更正pkg-config配置文件的查找顺序.