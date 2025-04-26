# 2FA
ref:
- [一文搞懂 OTP 双因素认证](https://monchickey.com/post/2023/09/18/otp-two-factor-authentication/)

## OTP
ref:
- [Key Uri Format](https://github.com/google/google-authenticator/wiki/Key-Uri-Format)

OTP 的实现主要有两种
- HOTP（HMAC-Based One-Time Password Algorithm）：基于 HMAC 的一次性密码

    安全令牌硬件（如某些银行的动态密保，游戏账号的密码令牌）通常是基于HOTP方式实现. 其表现为使用专有硬件，为某一特定产品或系统服务.
- TOTP（Time-Based One-Time Password Algorithm）：基于时间戳的一次性密码, **主流**

    app: Google Authenticator

    totp格式: `otpauth://totp/GitHub:xxx?secret=yyy&issuer=GitHub`

## pam
- [Configure SSH to use two-factor authentication](https://ubuntu.com/tutorials/configure-ssh-2fa#1-overview)
- [在 Ubuntu 22.04 系统上为 SSH 开启基于时间的 TOTP 认证](https://www.cnblogs.com/wx2020/p/17643066.html)
- [Ubuntu 20.04 开启并使用二步验证教程 (Two-Factor Authentication)](https://zhuanlan.zhihu.com/p/662638117)