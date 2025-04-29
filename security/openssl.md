# openssl

参考:
- [OpenSSL](https://netkiller.github.io/cryptography/openssl/index.html)
- [SSL/TLS CipherSuite 介绍](https://blog.helong.info/blog/2015/01/23/ssl_tls_ciphersuite_intro/)
- [cfssl创建证书](https://www.centos.bz/2017/09/k8s%E9%83%A8%E7%BD%B2%E4%B9%8B%E4%BD%BF%E7%94%A8cfssl%E5%88%9B%E5%BB%BA%E8%AF%81%E4%B9%A6/)
- [Kubernetes安装之证书验证](https://jimmysong.io/posts/kubernetes-tls-certificate/)
- [TLS 网站/API 安全评估](https://myssl.com/)
- [Openssl密码管理(文件加密)](http://foolishflyfox.xyz/blog/2020/04/12/linux/password-save/)

    `openssl aes-128-cbc -in passwd.txt -out passwd.txt.aes`

    > 支付宝/微信均使用aes cbc加密内容
- [使用 OpenSSL 加密和解密文件](https://linux.cn/article-13368-1.html)
- [浅析自签名证书应用](https://m.freebuf.com/articles/es/285265.html)

## 优化
- [3.4.3 SSL 层优化实践](https://www.thebyte.com.cn/http/ssl-performance.html)

    考虑建议使用 TLS1.3 + ECC 证书方式

## openssl

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

### dgst
生成Message Digest(信息摘要)

openssl dgst -sha256 freenas-v1.db

#### hmac
openssl dgst -sha256 -hmac "secretkey2" -out hmac_result.txt freenas-v1.db 

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

## pem
ref:
- [PKCS1与PKCS8的小知识](https://www.jianshu.com/p/a428e183e72e)

    PKCS1仅针对rsa
- [X.509、PKCS文件格式介绍](https://chanjarster.github.io/post/x509-pkcs-file-formats/)

    ASN.1是一种用来定义数据结构的接口描述语言, 其有一套关联的编码规则，这些编码规则用来规定如何用二进制来表示数据结构, DER是其中一种.

    PEM是一个用来存储和发送密码学key、证书和其他数据的文件格式的事实标准。许多使用ASN.1的密码学标准（比如X.509和PKCS）都使用DER编码，而DER编码的内容是二进制的，不适合与邮件传输（早期Email不能发送附件），因此使用PEM把二进制内容转换成ASCII码.

    PKCS8用于加密或非加密地存储Private Certificate Keypairs（不限于RSA）

    PKCS #12定义了通常用来存储Private Keys和Public Key Certificates（即X.509）的文件格式，使用基于密码的对称密钥进行保护. 注意Private Keys和Public Key Certificates是复数形式, 这意味着PKCS #12文件实际上是一个Keystore, PKCS #12文件可以被用做Java Key Store（JKS）.
- [pkcs8](https://github.com/youmark/pkcs8)

    Go package implementing functions to parse and convert private keys in PKCS#8 format, as defined in RFC5208 and RFC5958

RSA私钥存在PKCS1和PKCS8两种格式，通过openssl生成的私钥格式为PKCS1,公钥格式为PKCS8:
PKCS1:
```
-----BEGIN RSA PRIVATE KEY-----
...
-----END RSA PRIVATE KEY-----
```

PKCS8:
```
-----BEGIN PUBLIC KEY-----
...
-----END PUBLIC KEY-----
```

pem相关命令:
```bash
$ --- pkcs1 转 pkcs8
$ openssl pkcs8 -topk8 -inform PEM -in pkcs1.pem -outform PEM -out pkcs8.pem -nocrypt # 不加密
$ openssl pkcs8 -topk8 -in private-key.p1.pem -out private-key.p8.pem # 加密
$ --- pkcs8 加密/非加密互转
$ openssl pkcs8 -topk8 -in private-key.p8.nocrypt.pem -out private-key.p8.crypt.pem # 非加密->加密
$ openssl pkcs8 -topk8 -in private-key.p8.crypt.pem -out private-key.p8.nocrypt.pem -nocrypt # 加密 -> 非加密
$ --- pkcs8 转 pkcs1
$ openssl rsa -in private-key.p8.pem -out private-key.p1.pem
$ --- pem 转 der
$ openssl rsa -in pkcs1.pem -out pkcs1.der -outform DER
$ openssl pkcs8 -topk8 -inform PEM -in pkcs1.pem -outform DER -nocrypt -out pkcs8.der
$ --- 查看der
$ openssl asn1parse -i -in pkcs8.der -inform DER
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
openssl verify -CAfile ca.crt server-signed-by-ca.crt # 验证server证书
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

### AES密钥KEY和初始化向量IV
初始化向量IV可以有效提升安全性, 目的是防止同样的明文块，始终加密成同样的密文块. 在实际的使用场景中，它不能像密钥KEY那样直接保存在配置文件或固定写死在代码中，一般正确的处理方式为：在加密端将IV设置为一个16位的随机值，然后和加密文本一起返给解密端即可.

密钥KEY：AES标准规定区块长度只有一个值，固定为128Bit，对应的字节为16位。AES规定密钥长度只有三个值，128Bit、192Bit、256Bit，对应的字节为16位、24位和32位，密钥KEY不能公开传输，用于加密和解密数据；
初始化向量IV：该字段可以公开，用于将加密随机化。同样的明文被多次加密也会产生不同的密文，避免了较慢的重新产生密钥的过程，初始化向量与密钥相比有不同的安全性需求，因此IV通常无须保密。然而在大多数情况中，不应当在使用同一密钥的情况下两次使用同一个IV，**推荐初始化向量IV为16位的随机值**.

### [添加/删除ca cert](https://manuals.gfi.com/en/kerio/connect/content/server-configuration/ssl-certificates/adding-trusted-root-certificates-to-the-server-1605.html)
```bash
# --- debian/ubuntu
# cp ca.pem /usr/local/share/ca-certificates/
# update-ca-certificates
# ---- Remove your CA.
# rm /usr/local/share/ca-certificates/ca.pem
# update-ca-certificates --fresh
```

### 查看系统支持的ca cert
```bash
$ awk -v cmd='openssl x509 -noout -subject' '
    /BEGIN/{close(cmd)};{print | cmd}' < /etc/ssl/certs/ca-certificates.crt
$ awk -v decoder='openssl x509 -noout -subject -enddate 2>/dev/null' '
  /BEGIN/{close(decoder)};{print | decoder}
' < /etc/ssl/certs/ca-certificates.crt
```

### 内网https证书
ref:
- [Requirements and restrictions on IP addresses in SSL certificates](https://www.geocerts.com/support/ip-address-in-ssl-certificate)
- [Can an SSL Certificate Be Issued For an IP Address? ](https://sectigostore.com/page/ssl-certificate-for-ip-address/)
- [All publicly trusted SSL Certificates issued to internal names and reserved IP addresses
will expire before November 1, 2015.](https://www.digicert.com/kb/advisories/internal-names.htm)
- [Guidance on Internal Names](https://cabforum.org/working-groups/server/internal-names/)

证书服务商的证书只支持公网ip/域名(即可验证所有权), 内网https证书都是使用自签名即不可信.