#  passwd
ref:
- [Linux 密码安全策略设置](https://docs.azure.cn/zh-cn/articles/azure-operations-guide/virtual-machines/linux/aog-virtual-machines-linux-howto-site-linux-password-security-policy-settings)

修改用户密码、过期时间、认证信息等

实际上linux要求的密码验证机制是在/etc/login.defs中规定最小密码字符数; 同时还要受到/etc/pam.d/passwd的限定.                                                                                                                                                                                                                                                                     

## 格式
- -l 锁定用户，禁止其登录
- -u 解除锁定，允许用户登录
- --stdin 允许通过标准输入修改用户密码，如`echo "NewPassWord" | passwd --stdin Username `
- -d 使该用户可用空密码登录系统
- -e 强制用户在下次登录时修改密码
- -S 显示用户的密码是否被锁定，以及密码所采用的加密算法名称

## FAQ
### `passwd xxx`报`The password contains less than 3 character classes`
密码需要由 3 个类别（数字，小写字母，大写字母，其他）的字符组成, 由`/etc/security/pwquality.conf`的`minclass  = 3`决定, 按照`xxx`内容修改minclass即可.

### `passwd xxx`报`The password fails th dictionary check - it is based on on dictionary word`
输入的密码是常用密码, 被认为是不安全的.

解决方法:
1. 换复杂密码
1. 覆盖默认cracklib字典(未测试)

	默认cracklib在/usr/share/cracklib里, 覆盖默认cracklib字典:
	```bash
	# touch /usr/share/words
	# create-cracklib-dict /usr/share/words # 利用空文件覆盖默认cracklib字典, 建议备份/usr/share/cracklib
	```

### 直接替换密码
```bash
chmod +w /etc/shadow
sed -i 's@^root:.*@root:$6$2ApN2xUpr/SSRPcp$EAaZFSEymmd8NDTjkxR8...9Fuv...XRdjgs7p/sYrU.yERj4/:19018:0:99999:7:::@' /etc/shadow
chmod -w /etc/shadow
```

其实就是先用passwd生成指定密码, 再用其替换其他环境即可.