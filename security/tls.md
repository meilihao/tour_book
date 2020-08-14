# tls
## format
TLS_AES_128_GCM_SHA256, 分解:
1. protocol : TLS
1. AEAD cipher mode : AES_128_GCM
1. HKDF hash algorithm : sha256

## FAQ
### tlsv1.3 Kx=any Au=any
为了降低非前向加密连接和Bleichenbacher漏洞所带来的风险，RSA加密已从TLS 1.3中删除，仅支持三种key exchange:
- (EC)DHE (Diffie-Hellman over either finite fields or elliptic curves)
- PSK-only
- PSK with (EC)DHE

因此TLS 1.3定义了一组经过测试的DH参数，从而无需与服务器协商参数.

Digital Signature (Authentication) algorithms:
- ECDSA / EdDSA
- RSA
- a pre-shared key (PSK)
