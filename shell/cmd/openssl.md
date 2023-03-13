# openssl

## example
```bash
$ openssl s_client -connect baidu.com:443 -showcerts : 查看SSL握手信息（比如cipher协商结果）,证书信息等
...
# --- 上面命令的执行时会进入交互式会话中, 输入以下信息即可访问baidu页面
[...]
GET / HTTP/2.0
HOST: baidu.com

# --- 访问stmp
$ perl -MMIME::Base64 -e 'print encode_base64("username");'
$ perl -MMIME::Base64 -e 'print encode_base64("password");'
$ openssl s_client -starttls smtp -connect email.example.com:587
> ehlo example.com
> auth login
##paste your user base64 string here##
##paste your password base64 string here##
> mail from: noreply@example.com
> rcpt to: admin@example.com
> data
> Subject: Test 001
This is a test email.
.
> quit
# openssl req -new -key example-fd.priv.key -x509 -out example-fd.pub.key -days 365 -subj "/C=CN/ST=LiaoNing/L=DaLian/O=devops/OU=unicorn/CN=devops.com" # 非交互
```

## rsa
```bash
# openssl rsautl -inkey my_private_key -decrypt -oaep -in my_encrypted_file # default is PKCS#1 v1.5
```

## pkcs12
ref:
- [使用 OpenSSL 生成自签名证书](https://www.ibm.com/docs/zh/api-connect/2018.x?topic=overview-generating-self-signed-certificate-using-openssl)

```bash
openssl req -newkey rsa:2048 -nodes -keyout key.pem -x509 -days 365 -out certificate.pem
openssl x509 -text -noout -in certificate.pem # 检查已创建的证书
openssl pkcs12 -inkey key.pem -in certificate.pem -export -out certificate.p12 # 将密钥和证书组合在 PKCS#12 (P12) 中, 此过程需要输入密码
openssl pkcs12 -in certificate.p12 -noout -info # 验证P12 文件
```