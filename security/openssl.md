# openssl

参考:
- [OpenSSL](https://netkiller.github.io/cryptography/openssl/index.html)
- [SSL/TLS CipherSuite 介绍](https://blog.helong.info/blog/2015/01/23/ssl_tls_ciphersuite_intro/)
- [cfssl创建证书](https://www.centos.bz/2017/09/k8s%E9%83%A8%E7%BD%B2%E4%B9%8B%E4%BD%BF%E7%94%A8cfssl%E5%88%9B%E5%BB%BA%E8%AF%81%E4%B9%A6/)
- [Kubernetes安装之证书验证](https://jimmysong.io/posts/kubernetes-tls-certificate/)
- [TLS 网站/API 安全评估](https://myssl.com/)

### openssl ciphers
CipherSuite 包含多种技术，例如认证算法（Authentication）、加密算法（Encryption）、消息认证码算法（Message Authentication Code，简称为 MAC）、密钥交换算法（Key Exchange）和密钥衍生算法（Key Derivation Function）.

SSL 的 CipherSuite 协商机制具有良好的扩展性，每个 CipherSuite 都需要在 IANA 注册，并被分配两个字节的标志。全部 CipherSuite 可以在 IANA 的 [TLS Cipher Suite Registry](https://www.iana.org/assignments/tls-parameters/tls-parameters.xhtml#tls-parameters-4) 页面查看


查看openssl支持的ciphes: `openssl ciphers -V | column -t`

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

### crt
```bash
openssl x509 -noout -text -in ca.crt
# --- 验证
openssl verify selfsign.crt # 验证自签名crt
openssl verify -CAfile ca.crt myserver.crt # 验证非自签名crt
```

### rsa
参考:
- [ 利用openssl进行RSA加密解密](https://www.cnblogs.com/alittlebitcool/archive/2011/09/22/2185418.html)

```sh
$ openssl rsa -noout -text -in my.key # 查看私钥内容
$ openssl rsa -noout -text -pubin -in my.key.pub # 查看公钥内容
$ openssl rsautl -encrypt -in hello -inkey test_pub.key -pubin -out hello.en # 公钥加密
$ openssl rsautl -decrypt -in hello.en -inkey test.key -out hello.de # 私钥解密
$ cat license.dat |base64 -d - |openssl rsautl -verify -in - -inkey licens.key.pub -pubin -out clear_text # 公钥解密. -pubin：表明输入的是一个公钥文件，默认输入为私钥文件
```

### rsa解密
参考:
- [RSA填充方式](https://www.jianshu.com/p/205abb4b9dc6)

```python
# from PKCS1_v1_5.new(RSA.importKey(pubKey)).verify(h, base64.b64decode(signature))
def rsaPublicDecrypt(S):
    key = RSA.importKey(pubKey)

    modBits = Crypto.Util.number.size(key.n)
    k = ceil_div(modBits,8) # Convert from bits to bytes
    # Step 1
    if len(S) != k:
        return ""
        # Step 2a (O2SIP) and 2b (RSAVP1)
        # Note that signature must be smaller than the module
        # but RSA.py won't complain about it.
        # TODO: Fix RSA object; don't do it here.
    m = key.encrypt(S, 0)[0]
    if len(m) > 0 and m[0] == '\x01': # rsa encrypt的输出编码(RSA_padding_add_PKCS1_type_1): 01 FF ... FF 00 + data
            pos = m.find('\x00')
            if pos > 0:
                m=m[pos+1:]
    return str(m)
```

## FAQ
### 网站不支持ALPN提示"No ALPN negotiated"
[ALPN介绍](https://imququ.com/post/enable-alpn-asap.html),检查是否支持ALPN的命令:
```sh
$ openssl s_client -alpn h2 -servername imququ.com -connect golang.d.openhello.net:443 < /dev/null | grep 'ALPN'
```
测试结果:
- 如果提示`unknown option -alpn`，说明本地的 openSSL 版本太低，请升级到`1.0.2+`,再进行测试.
- 如果结果包含`ALPN protocol: h2`，说明服务端支持 ALPN.
- 如果结果包含`No ALPN negotiated`，说明服务端不支持 ALPN，浏览器无法协商到 HTTP/2，需要尽快升级openssl.

本质就是**是否支持 ALPN 完全取决于服务端使用的 OpenSSL 版本**.

由于 openssl 是公共基础库，大量其他软件都对它有依赖，如果直接升级系统自带的 openssl，很容易引发各种问题. 更为稳妥的做法是在编译 Web Server 时自己指定 openssl 的位置,比如[编译nginx](https://imququ.com/post/enable-tls-1-3.html).

> 其实也可使用 Qualys SSL Labs's SSL Server Test 这个在线工具来测试,测试项目为`ALPN`.

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

### [openssl rsa/pkey(查看私钥、从私钥中提取公钥、查看公钥)](https://www.cnblogs.com/wyzhou/p/9738964.html)
### rsa验签(public key加密/private key解密)
```bash
openssl genrsa -out mykey
openssl rsa -in mykey -pubout -out mykey.pub
openssl rand -hex 64 -out myfile # Generate the random password file 
md5sum myfile | openssl rsautl -inkey mykey -sign > checksum.signed
openssl rsautl -inkey mykey.pub -pubin -in checksum.signed
```

或者(实际与上面一样, 就是参数有点小变化):
```bash
openssl rsautl -sign -in file -inkey mykey -out sig # Sign some data using a private key. 实际就是加密file, sig是加密后的文件
hexdump -C sig
openssl rsautl -verify -in sig -inkey mykey.pub -pubin # Recover the signed data. 实际就是解密sig. `-inkey`是公钥时需用`-pubin`指明
openssl rsautl -verify -in sig -inkey mykey.pub -pubin -raw -hexdump # Examine the raw signed data. 输出是PKCS#1格式
```

### rsa+aes加密文件
largefile.pdf是原始文件, largefile.pdf.enc是加密后的文件

```bash
openssl rand -hex 64 -out key.bin
openssl enc -aes-256-cbc -salt -in largefile.pdf -out largefile.pdf.enc -pass file:./bin.key
openssl rsautl -encrypt -inkey publickey.pem -pubin -in key.bin -out key.bin.enc # 加密aes key
openssl rsautl -decrypt -inkey privatekey.pem -in key.bin.enc -out key.bin # 解出aes key
openssl enc -d -aes-256-cbc -in largefile.pdf.enc -out largefile.pdf -pass file:./bin.key # 解密文件
```

### ERR_CERT_COMMON_NAME_INVALID(chrome)/SSL_ERROR_BAD_CERT_DOMAI(firefox)
向chrome/firefox导入自签名证书并重启浏览器后, 网站还是显示不安全, 这是因为虽然已信任ca证书, 但网站证书里的主题备用名(SAN,Subject Alternative Name)与实际使用的ip/域名不匹配.

SAN无法对所有ip进行授权, 只能使用多SAN的形式指定多个ip地址.

其他方案:
- 让用户先手动访问目标地址授权不拦截该种不安全的https, 然后再进行操作. 场景: 不通过管理节点中转, 而是直接向存储节点上传iso. 