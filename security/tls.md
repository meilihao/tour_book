# tls
参考:
- [HTTPS 温故知新（四） —— 直观感受 TLS 握手流程(下)](https://halfrost.com/https_tls1-3_handshake/)
- [TLS 1.3科普——新特性与协议实现](https://zhuanlan.zhihu.com/p/28850798)
- [HTTP - TLS1.3 初次解读](https://cloud.tencent.com/developer/article/2123171)

## tls协议
TLS 1.3包括3个子协议——alert、handshake、record:

- handshake : 负责协商使用的TLS版本、加密算法、哈希算法、密钥材料和其他与通信过程有关的信息，对服务器进行身份认证，对客户端进行可选的身份认证，最后对整个握手阶段信息进行完整性校验以防范中间人攻击，是整个TLS协议的核心
- record : 负责对接收到的报文进行加密解密，将其分片为合适的长度后转发给其他协议层。
- alert : 负责处理消息传输与握手阶段中的异常情况。 

## format
TLS_AES_128_GCM_SHA256, 分解:
1. protocol : TLS
1. AEAD cipher mode : AES_128_GCM
1. HKDF hash algorithm : sha256

## 改进
1. 密钥交换：ClientHello含有key exchange信息, server收到ClientHello处理key exchange后所有的数据都会被加密

## FAQ
### tlsv1.3 Kx=any Au=any
为了降低非前向加密连接和Bleichenbacher漏洞所带来的风险，RSA加密已从TLS 1.3中删除，仅支持三种key exchange:
- (EC)DHE (Diffie-Hellman over either finite fields or elliptic curves)
- PSK-only

    client写死一个key，server写死一个相同的key
- PSK with (EC)DHE

因此TLS 1.3定义了一组经过测试的DH参数，从而无需与服务器协商参数.

Digital Signature (Authentication) algorithms:
- ECDSA / EdDSA
- RSA
- a pre-shared key (PSK)

### PSK（pre_shared_key）——新的密钥交换暨身份认证机制
0-RTT：客户端和服务端的一次交互（客户端发一个报文，服务端回应一个报文）叫做一个RTT，TLS 1.2普遍采用2-RTT的握手过程，服务器延迟明显. 因此TLS 1.3引入了一种0-RTT的机制，即在刚开始TLS密钥协商的时候，就能附送一部分经过加密的数据传递给对方.

为了实现0-RTT，需要双方在刚开始建立连接的时候就已经持有一个对称密钥，这个密钥在TLS 1.3中称为PSK（pre_shared_key）.

PSK是TLS 1.2中的rusumption机制的一个升级，TLS握手结束后，服务器可以发送一个NST（new_session_ticket）的报文给客户端，该报文中记录PSK的值、名字和有效期等信息，双方下一次建立连接可以使用该PSK值作为初始密钥材料。因为PSK是从以前建立的安全信道中获得的，只要证明了双方都持有相同的PSK，不再需要证书认证，就可以证明双方的身份，因此，PSK也是一种身份认证机制。 

ps: 0-RTT的实现有一定的安全缺陷，自身没有抗重放攻击的机制. 在TLS 1.3中提出了几个对性能消耗比较大的可能的解决方法.

### HKDF（HMAC_based_key_derivation_function）——新的密钥导出函数

经过密钥协商得出来的密钥材料因为随机性可能不够，协商的过程能被攻击者获知，需要使用一种密钥导出函数来从初始密钥材料（PSK或者DH密钥协商计算出来的key）中获得安全性更强的密钥. HKDF正是TLS 1.3中所使用的这样一个算法，使用协商出来的密钥材料和握手阶段报文的哈希值作为输入，可以输出安全性更强的新密钥.

HKDF包括extract_then_expand的两阶段过程. extract过程增加密钥材料的随机性，在TLS 1.2中使用的密钥导出函数PRF实际上只实现了HKDF的expand部分，并没有经过extract，而直接假设密钥材料的随机性已经符合要求.

### AEAD（Authenticated_Encrypted_with_associated_data）——唯一保留的加密方式

TLS协议的最终目的是协商出会话过程使用的对称密钥和加密算法，双方最终使用该密钥和对称加密算法对报文进行加密. AEAD将完整性校验和数据加密两种功能集成在同一算法中完成，是TLS 1.3中唯一支持的加密方式.

TLS 1.2还支持流加密和CBC分组模式的块加密方法，使用MAC来进行完整性校验数据，这两种方式均被证明有一定的安全缺陷. 但是有研究表明AEAD也有一定局限性：使用同一密钥加密的明文达到一定长度后，就不能再保证密文的安全性。因此，TLS 1.3中引入了密钥更新机制，一方可以（通常是服务器）向另一方发送Key Update（KU）报文，对方收到报文后对当前会话密钥再使用一次HKDF，计算出新的会话密钥，使用该密钥完成后续的通信.

### ssl_ciphers选择
筛选命令(只包含tls1.2):
```sh
$ openssl ciphers -V 'ALL'|grep "1.2"|egrep -v "Kx=DH|Kx=PSK|Kx=ECDHEPSK|RSAPSK|Camellia"|egrep -v "Enc=AESGCM\(256\)|Enc=AESCCM\(256\)|Enc=AESCCM8"|grep -v "Mac=SHA384"|egrep -v "Enc=AES\(256\)"|column  -t
TLS_AES_128_GCM_SHA256 TLS_CHACHA20_POLY1305_SHA256 TLS_AES_256_GCM_SHA384 \ # tls1.3
ECDHE-ECDSA-AES128-GCM-SHA256 ECDHE-ECDSA-CHACHA20-POLY1305 \ # ECDHE+ECDSA+AEAD
ECDHE-RSA-AES128-GCM-SHA256 ECDHE-RSA-CHACHA20-POLY1305 \ # ECDHE+RSA+AEAD
ECDHE-ECDSA-AES128-SHA256 ECDHE-RSA-AES128-SHA256 # ECDHE+!AEAD
ECDHE-ECDSA-AES128-SHA ECDHE-RSA-AES128-SHA # TLSv1 for win7,旧Android
```

通过[ssllabs](https://www.ssllabs.com/ssltest/analyze.html)对比发现`ECDHE-ECDSA-*`和`ECDHE-RSA-*`支持的设备跨度是一样的,因此仅保留`ECDSA`即可:
```
TLS_AES_128_GCM_SHA256 TLS_CHACHA20_POLY1305_SHA256 TLS_AES_256_GCM_SHA384 \ # tls1.3
ECDHE-ECDSA-AES128-GCM-SHA256 ECDHE-ECDSA-CHACHA20-POLY1305 \ # ECDHE+ECDSA+AEAD
ECDHE-ECDSA-AES128-SHA256 # ECDHE+!AEAD
ECDHE-ECDSA-AES128-SHA # TLSv1 for win7,旧Android
```

> 在配置 CipherSuite 时，请务必参考权威文档，如：[CloudFlare 使用的配置](https://github.com/cloudflare/sslconfig/blob/master/conf);[Mozilla 的推荐配置](https://wiki.mozilla.org/Security/Server_Side_TLS#Recommended_configurations)
> ssl_ecdh_curve选择: `ssl_ecdh_curve   X25519:P-256:P-384:P-224:P-521;`

### Create ED25519 certificates for TLS with OpenSSL
ref:
- [openssl ca vs openssl x509 comparison [With Examples]](https://www.golinuxcloud.com/openssl-ca-vs-openssl-x509-comparison/)
- [Create ED25519 certificates for TLS with OpenSSL](https://blog.pinterjann.is/ed25519-certificates.html)
- [PKI Certificate Authority questions. (ED25519) Design, best practices, how to.](https://www.reddit.com/r/sysadmin/comments/1417lhz/pki_certificate_authority_questions_ed25519/)

ps:
- 添加使用cfssl
- 使用`openssl x509`替换`openssl ca`

    openssl ca维护了证书db, 不允许签发同一证书多次, 除非先撤销之前证书.

```bash
# openssl version
OpenSSL 3.0.2 15 Mar 2022 (Library: OpenSSL 3.0.2 15 Mar 2022)
# mkdir private certs csr newcerts
# touch index.txt
# echo 1001 > serial
# --- Generate ED25519 private key for Root cert
# openssl genpkey -algorithm ED25519 -out private/ca.key.pem
# cat openssl-25519.cnf # 仅用于演示创建root ca
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no
[req_distinguished_name]
C = DE
CN = www.example.com
[v3_req]
keyUsage = keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = www.example.com
DNS.2 = example.com
# cat openssl-25519v2.cnf
# OpenSSL intermediate CA configuration file.
# Copy to `/root/ca/cert-archive/intermediate/openssl.cnf`.

[ ca ]
# `man ca`
default_ca = CA_default

[ CA_default ]
# Directory and file locations.
dir               = /home/chen/test/ed25519
certs             = $dir/certs
crl_dir           = $dir/crl
new_certs_dir     = $dir/newcerts
database          = $dir/index.txt
serial            = $dir/serial
RANDFILE          = $dir/private/.rand

# The root key and root certificate.
private_key       = $dir/private/ca.key.pem
certificate       = $dir/certs/ca.cert.pem

# For certificate revocation lists.
crlnumber         = $dir/crlnumber
crl               = $dir/crl/intermediate.crl.pem
crl_extensions    = crl_ext
default_crl_days  = 3652

# SHA-1 is deprecated, so use SHA-2 instead.
default_md        = sha256

name_opt          = ca_default
cert_opt          = ca_default
default_days      = 7300
preserve          = no
policy            = policy_loose

[ policy_strict ]
# The root CA should only sign intermediate certificates that match.
# See the POLICY FORMAT section of `man ca`.
countryName             = match
stateOrProvinceName     = match
organizationName        = match
organizationalUnitName  = optional
commonName              = supplied
emailAddress            = optional

[ policy_loose ]
# Allow the intermediate CA to sign a more diverse range of certificates.
# See the POLICY FORMAT section of the `ca` man page.
countryName             = optional
stateOrProvinceName     = optional
localityName            = optional
organizationName        = optional
organizationalUnitName  = optional
commonName              = supplied
emailAddress            = optional
subjectAltName          = optional

[ req ]
distinguished_name = req_distinguished_name
x509_extensions     = server_cert
req_extensions      = v3_req

[ v3_req ]
# Extensions for a certificate request (`man x509v3_config`).
subjectAltName = @alt_names

[ req_distinguished_name ]
# See <https://en.wikipedia.org/wiki/Certificate_signing_request>.
countryName                     = Country Name (2 letter code)
stateOrProvinceName              = State or Province Name
localityName                    = Locality Name
0.organizationName              = Organization Name
organizationalUnitName          = Organizational Unit Name
commonName                      = Common Name
emailAddress                    = Email Address

# Optionally, specify some defaults.
countryName_default             = US
stateOrProvinceName_default     = NOTEXIST
localityName_default            = SOMEWHERE
0.organizationName_default      = Example Co.
organizationalUnitName_default  = IT
emailAddress_default            = not@exist.us

[ v3_ca ]
# Extensions for a typical CA (`man x509v3_config`).
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always,issuer
basicConstraints = critical, CA:true, pathlen:4
keyUsage = critical, digitalSignature, cRLSign, keyCertSign
subjectAltName = @alt_names

[ v3_intermediate_ca ]
# Extensions for a typical intermediate CA (`man x509v3_config`).
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always,issuer
basicConstraints = critical, CA:true, pathlen:1
keyUsage = critical, digitalSignature, cRLSign, keyCertSign
subjectAltName = @alt_names

[ usr_cert ]
# Extensions for client certificates (`man x509v3_config`).
basicConstraints = CA:FALSE
nsCertType = client
nsComment = "OpenSSL Generated Client Certificate"
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid,issuer:always
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth

[ server_cert ]
# Extensions for server certificates (`man x509v3_config`).
basicConstraints = CA:FALSE
nsCertType = server
nsComment = "OpenSSL Generated Server Certificate"
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid,issuer:always
keyUsage = critical, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[ alt_names ]
IP.1 = 10.10.60.1
DNS.1 = appajava.server1.test.int.local
DNS.2 = server1.test.int.local

[ crl_ext ]
# Extension for CRLs (`man x509v3_config`).
authorityKeyIdentifier=keyid:always
# --- generate self signed root ca from config file, 有两种方法:
## ---req和x509分开处理
# openssl req -new -out csr/ca.cert.csr -key private/ca.key.pem -config openssl-25519.cnf
# openssl req -in csr/ca.cert.csr -text -noout
# openssl x509 -req -days 700 -in csr/ca.cert.csr -signkey private/ca.key.pem -out certs/ca.cert.pem
## ---req和x509一起处理by`-x509`
# openssl req -config openssl-25519v2.cnf -key private/ca.key.pem -new -x509 -days 7300 -sha256 -extensions v3_ca -out certs/ca.cert.pem
# openssl x509 -in certs/ca.cert.pem -text -noout
# --- INTERMEDIATE, 也可将openssl-25519v2.cnf拷贝一份, 作为intermediate自己的签名配置, 记得修改`[ CA_default ]`里的private_key和certificate
# --- Generate ED25519 private key for intermediate cert
# openssl genpkey -algorithm ED25519 -out private/intermediate_ca.key.pem
# Genereate CSR for intermediate cert
openssl req -config openssl-25519v2.cnf -new -sha256 -key private/intermediate_ca.key.pem  -extensions v3_intermediate_ca -out csr/intermediate.csr.pem # `Common Name`必填, 可用subj指定. sha256是CSR签名时用的摘要算法
# --- Sign the intermediate cert with the root CA
# openssl ca -config openssl-25519v2.cnf -extensions v3_intermediate_ca -days 3650 -notext -md sha256 -in csr/intermediate.csr.pem -out certs/intermediate.cert.pem
# openssl x509 -in certs/intermediate.cert.pem -text -noout
# --- SERVER
# --- Generate ED25519 private key for server cert
# openssl genpkey -algorithm ED25519 -out private/server.key.pem
# --- Genereate CSR for server cert
# openssl req -config openssl-25519v2.cnf -extensions v3_req -key private/server.key.pem -new -sha256 -out csr/server.csr.pem # `Common Name`必填
# --- Sign the intermediate cert with the intermediate CA
# openssl ca -config openssl-25519v2.cnf -extensions server_cert -days 3750 -notext -md sha256 -in csr/server.csr.pem -out certs/server.cert.pem
```