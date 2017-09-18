# openssl

参考:
- [OpenSSL](https://netkiller.github.io/cryptography/openssl/index.html)

### openssl ciphers
CipherSuite 包含多种技术，例如认证算法（Authentication）、加密算法（Encryption）、消息认证码算法（Message Authentication Code，简称为 MAC）、密钥交换算法（Key Exchange）和密钥衍生算法（Key Derivation Function）.

SSL 的 CipherSuite 协商机制具有良好的扩展性，每个 CipherSuite 都需要在 IANA 注册，并被分配两个字节的标志。全部 CipherSuite 可以在 IANA 的 [TLS Cipher Suite Registry](https://www.iana.org/assignments/tls-parameters/tls-parameters.xhtml#tls-parameters-4) 页面查看

在配置 CipherSuite 时，请务必参考权威文档，如：[Mozilla 的推荐配置](https://wiki.mozilla.org/Security/Server_Side_TLS#Recommended_configurations)、[CloudFlare 使用的配置](https://github.com/cloudflare/sslconfig/blob/master/conf)

参考:
- [The Basics of How to Work with Cipher Settings](https://drjohnstechtalk.com/blog/2011/09/the-basics-of-how-to-work-with-ciphers/)
- [TLS协议分析 与 现代加密通信协议设计](http://ju.outofmemory.cn/entry/210548)
- [cipher list的格式](http://timd.cn/2016/06/29/nginx-https/)

### openssl verison
- 获取openssl配置文件位置: `openssl version -a`的`OPENSSLDIR`项目.

### openssl list
- -digest-algorithms      List of message digest algorithms 列出所有可用的信息摘要算法
- -cipher-commands        List of cipher commands
- -cipher-algorithms      List of cipher algorithms # 列出所有可用的加密算法
- -public-key-algorithms  List of public key algorithms # 列出所有可用的公钥加密算法

### openssl s_client

ps:
- 如果使用`openssl s_client`的`-cipher`参数测试密码套件时要在其**前面**指定相应的tls版本比如`-tls1_2`,否则会输出`Cipher is (NONE)`.

### openssl ecparam
- -list_curves           Prints a list of all curve 'short names' 列出所有可用的EC曲线

### openssl x509
- -x509toreq             基于当前证书创建一个CSR文件

## Error
### version `OPENSSL_XXX' not found (required by openssl)
为使用tls1.3,手动编译openssl,运行openssl命令报该错误:
```
$ ~/openssl-master ./config enable-tls1_3 
$ ~/openssl-master make
$ ~/openssl-master sudo make install
$ openssl version
openssl: /usr/lib/x86_64-linux-gnu/libssl.so.1.1: version `OPENSSL_1_1_1' not found (required by openssl)
openssl: /usr/lib/x86_64-linux-gnu/libcrypto.so.1.1: version `OPENSSL_1_1_1' not found (required by openssl)
```

原因:
系统原本有`OpenSSL 1.1.0e  16 Feb 2017`,现在切换到`OPENSSL_1_1_1`,依赖的so还是旧版本.

解决(3种):
1. 在`/etc/ld.so.conf.d/x86_64-linux-gnu.conf`最上面添加`/usr/local/lib(openssl的默认lib位置)`,再用`sudo ldconfig`使配置生效即可,**推荐**.
1. 添加配置编译参数`./config enable-tls1_3 -Wl, -rpath=/usr/local/lib`,指定该程序的动态库搜索路径.
1. 通过环境变量`LD_LIBRARY_PATH`指定动态库搜索路径

### Can't locate Test/Harness.pm in @INC
`make`后运行`make test`时报该错误.

阅读官方repo的[INSTALL](https://github.com/openssl/openssl/blob/master/INSTALL),上面有关于perl的要求,具体在[NOTES.PERL](https://github.com/openssl/openssl/blob/master/NOTES.PERL)里.