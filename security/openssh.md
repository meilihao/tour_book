# SSH 使用密钥登录并禁止口令登录实践

[原文](http://wsgzao.github.io/post/ssh/)

### 扩展阅读

- [云服务器 ECS Linux SSH 基于密钥交换的自动登录原理简介及配置说明](https://help.aliyun.com/knowledge_detail/41493.html)
- [SSH原理与运用](http://www.ruanyifeng.com/blog/2011/12/ssh_remote_login.html)
- [Linode](https://www.linode.com/docs/networking/ssh/use-public-key-authentication-with-ssh)

## 最新
[OpenSSH 8.2](http://www.openssh.com/txt/release-8.2)增加了[对 FIDO/U2F 硬件身份验证器的支持](https://www.debian.cn/archives/3683).

## 前言

无论是个人的VPS还是企业允许公网访问的服务器，如果开放22端口的SSH密码登录验证方式，被众多黑客暴力猜解捅破菊花也可能是经常发生的惨剧。企业可以通过防火墙来做限制，普通用户也可能借助修改22端口和强化弱口令等方式防护，但目前相对安全和简单的方案则是让SSH使用密钥登录并禁止口令登录。

**这是最相对安全的登录管理方式**

## 生成PublicKey

建议**使用ed25519并设置passphrase密码短语**，以Linux生成为例

```ssh
// 私钥 (id_rsa) 与公钥 (id_rsa.pub)
#生成SSH密钥对
// `-C` : 加入注释标识方便管理
[echo -e 'y\n' |] ssh-keygen -t rsa -b 2048 -C "xxx" -P "" [-f  $HOME/.ssh/id_rsa] # `-P ""`表示不用输入passphrase; `-f`指定key的保存路径; `echo -e 'y\n' |`(等价于`yes |`)用于强制覆盖已有的同名key
ssh-keygen -t ed25519 -C "xxx"

Generating public/private rsa key pair.
#建议直接回车使用默认路径
Enter file in which to save the key (/root/.ssh/id_rsa):
#输入密码短语（留空则直接回车）
Enter passphrase (empty for no passphrase):
#重复密码短语
Enter same passphrase again:
Your identification has been saved in /root/.ssh/id_rsa.
Your public key has been saved in /root/.ssh/id_rsa.pub.
The key fingerprint is:
aa:8b:61:13:38:ad:b5:49:ca:51:45:b9:77:e1:97:e1 root@localhost.localdomain
The key's randomart image is:
+--[ RSA 2048]----+
|    .o.          |
|    ..   . .     |
|   .  . . o o    |
| o.  . . o E     |
|o.=   . S .      |
|.*.+   .         |
|o.*   .          |
| . + .           |
|  . o.           |
+-----------------+
```

```
# openssl默认生成的密钥格式为PEM
# 导出RSA公钥格式
openssl rsa -in <私钥> -pubout
# 导出ssh2公钥格式
ssh-keygen -y -f <私钥>
# 将pem公钥转成ssh2格式
ssh-keygen -i -m PKCS8 -f pubkey.pem
```

> [OpenSSH 6.5才开始支持ed25519](https://www.openssh.com/txt/release-6.5)

### 复制密钥对

> authorized_keys需要ssh2公钥格式

也可以手动在客户端建立目录和authorized_keys，注意修改权限

```
#复制公钥到无密码登录的服务器上,22端口改变可以使用下面的命令,**推荐**
ssh-copy-id -i ~/.ssh/id_rsa.pub root@192.168.15.241
// 或
// 该方法要注意.ssh及其子文件的权限问题
cat ~/.ssh/id_rsa.pub | ssh root@120.26.38.248 'cat >> ~/.ssh/authorized_keys'
```

### 修改SSH配置文件

```
#编辑sshd_config文件
vi /etc/ssh/sshd_config

Protocol 2

RSAAuthentication yes
PubkeyAuthentication yes
#指定公钥数据库文件
#AuthorsizedKeysFile .ssh/authorized_keys

PermitRootLogin no

#禁用密码验证(**RSA登录成功后再禁用**)
PasswordAuthentication no

#禁止空密码登录
PermitEmptyPasswords no
```

重启SSH服务前建议多保留一个会话以防不测
```
#RHEL/CentOS系统
service sshd restart
#ubuntu系统
service ssh restart
#debian系统
/etc/init.d/ssh restart
```

```
// 调试
ssh -Tv -i $rsa_primary_key root@192.168.15.241
// 调试,更多信息
ssh -Tvvv -i $rsa_primary_key root@192.168.15.241
```

## FAQ
### 普通用户无法用ssh通信, 但root可以

```bash
# --- ssh -V
server: OpenSSH_5.3p1, OpenSSL 1.0.0-fips 29 Mar 2010
client: OpenSSH_7.3p1 Ubuntu-1, OpenSSL 1.0.2g  1 Mar 2016
```
root能用pubkey登录,普通用户不能.

解决: 创建普通用户指定shell时,其名称错误.

### ssh支持的算法
```bash
ssh -Q cipher       # List supported ciphers
ssh -Q mac          # List supported MACs
ssh -Q key          # List supported public key types
ssh -Q kex          # List supported key exchange algorithms
ssh -Q sig          # List supported signature algorithms
```
