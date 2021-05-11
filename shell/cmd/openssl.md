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
```