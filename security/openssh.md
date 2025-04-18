# SSH 使用密钥登录并禁止口令登录实践

[原文](http://wsgzao.github.io/post/ssh/)

### 扩展阅读

- [云服务器 ECS Linux SSH 基于密钥交换的自动登录原理简介及配置说明](https://help.aliyun.com/knowledge_detail/41493.html)
- [SSH原理与运用](http://www.ruanyifeng.com/blog/2011/12/ssh_remote_login.html)
- [Linode](https://www.linode.com/docs/networking/ssh/use-public-key-authentication-with-ssh)
- [Google Authenticator实现SSH登陆二次验证](https://www.infvie.com/ops-notes/google-authenticator-sshd.html)
- [Linux下使用ssh动态验证码登陆机器](https://blog.51cto.com/noodle/1831697)

## 最新
- [OpenSSH 8.2](http://www.openssh.com/txt/release-8.2)增加了[对 FIDO/U2F 硬件身份验证器的支持](https://www.debian.cn/archives/3683)
- [Cloudflare 开源 OpenPubkey SSH（OPKSSH），将单点登录集成到 SSH](https://www.oschina.net/news/341113/cloudflare-open-sourcing-openpubkey)

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

ssh-copy-id添加的pub key是追加到authorized_keys末尾的, 且不检查是否重复.

### 修改SSH配置文件

```
#编辑sshd_config文件
vi /etc/ssh/sshd_config

Protocol 2

#RSAAuthentication yes # Protocol 1支持的配置
PubkeyAuthentication yes # Protocol 2支持的配置
#指定公钥数据库文件
#AuthorsizedKeysFile .ssh/authorized_keys

PermitRootLogin no # prohibit-password: 允许root以非密码方式登入, 比如key

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

## ssh协议
SSH协议是建立在不安全的网络之上的进行远程安全登陆的协议. 它是一个协议族，其中有三个子协议，分别是：
1. 传输层协议[SSH-TRANS]:提供服务器验证、完整性和保密性功能,建立在传统的TCP/IP协议之上
2. 验证协议[SSH-USERAUTH]:向服务器验证客户端用户，有基于用户名密码和公钥两种验证方式，建立在传输层协议[SSH-TRANS]之上
3. 连接协议[SSH-CONNECT]:将加密隧道复用为若干逻辑信道。它建立在验证协议之上

## ssh通信过程 by `ssh aliyun -vvvv`
SSH建立连接的过程主要分为下面几个阶段：
1. SSH协议版本协商阶段. openSSH推荐仅使用SSH2

	1. client 根据username@server或`.ssh/config`中的server addr与server三次握手建立连接
	1. 服务器通过建立好的连接向客户端发送一个包含SSH版本信息的报文，格式为“SSH-<SSH协议大版本号>.<SSH协议小版本号>-<软件版本号>”，软件版本号主要用于调试
	1. 客户端将自己决定使用的版本号发给服务器，服务器判断客户端使用的版本号自己是否支持，从而决定是否能够继续完成SSH连接
	1. 如果协商成功，则进入密钥和算法协商阶段

1. 密钥和算法协商阶段，SSH支持多种加密算法，双方根据自己和对端支持的算法进行协商，最终决定要使用的算法

	1. 服务器端和客户端分别发送算法协商报文给对端，报文中包含自己支持的公钥算法列表，加密算法列表，MAC（Message Authentication Code，消息验证码）算法列表，压缩算法列表等
	1. 和版本协商阶段类似，服务器端和客户端根据自己和对端支持的算法来决定最终要使用的各个算法, 以客户端支持的协议为主
	1. 服务器端和客户端利用Diffie-Hellman密钥交换算法，主机密钥对等参数，生成共享密钥和会话ID. 会话密钥用于在后续的通信过程中两端对传输的数据进行加密和解密，而会话ID用于认证过程
1. 认证阶段，服务器对客户端进行身份验证

	1. 客户端向服务器端发送认证请求，请求中包含用户名，认证方法，密码或密钥
	1. 服务器端对客户端进行认证，如果认证失败，则向客户端发送失败消息，其中包含可以再次认证的方法列表
	1. 客户端再次使用支持的认证方法中的一种进行认证，直到达到认证次数上限被服务器终止连接，或者认证成功为止

	SSH支持的两种认证方式：
    - 密码认证：客户端通过用户名/密码进行认证，将使用会话密钥加密后的用户名和密码发送给服务器，服务器解密后与系统保存的用户名和密码进行对比，并向客户端返回认证成功或失败的消息
    - 密钥认证：采用数字签名来进行认证，目前可以通过RSA,ECDSA,ed25519实现数字签名，客户端通过用户名，公钥以及公钥算法等信息来与服务器完成验证
1. 会话请求阶段，完成认证后，客户端会向服务器端发送会话请求

	1. 服务器等待客户端请求
    1. 认证完成后，客户端向服务器发送会话请求
    1. 服务器处理客户端请求，完成后，会向客户端回复SSH_SMSG_SUCCESS报文，双方进入交互会话阶段. 如果请求未被成功处理，则服务器返回SSH_SMSG_FAILURE报文，表示请求处理失败或者不能识别客户端请求
1. 交互会话阶段，会话请求通过后，服务器端和客户端进行信息的交互

	1. 客户端将要执行的命令加密发送给服务器
	1. 服务器收到后，解密命令，执行后将结果加密返回客户端
	1. 客户端将返回结果解密后显示到终端上

## FAQ
### 普通用户无法用ssh通信, 但root可以

```bash
# --- ssh -V
server: OpenSSH_5.3p1, OpenSSL 1.0.0-fips 29 Mar 2010
client: OpenSSH_7.3p1 Ubuntu-1, OpenSSL 1.0.2g  1 Mar 2016
```
root能用pubkey登录,普通用户不能.

解决: 创建普通用户指定shell时,其名称错误.

### 添加ssh key后无法登陆
检查:
- /home/$USER/.ssh权限: 700
- /home/$USER/.authorized_keys权限: 600
- /home/$USER/.authorized_keys内容: 比如key的开头少字母, 这种情况通常在使用vim粘贴时出现

### ssh支持的算法
```bash
ssh -Q cipher       # List supported ciphers
ssh -Q mac          # List supported MACs
ssh -Q key          # List supported public key types
ssh -Q kex          # List supported key exchange algorithms
ssh -Q sig          # List supported signature algorithms
```

### 为什么需要known_hosts
SSH client就是通过known_hosts中的host key来验证Server的身份的.

### 查看Server host key即远程主机的key fingerprint
```
chen@xxx:/etc/ssh$ ssh-keygen  -lf ssh_host_ecdsa_key.pub
256 SHA256:kM9uQJBdQt9JGlDkuIh4bIJSWjF5EPnTpcq5X1pMmVw root@iZuf6hftd4ce4kf92zb5ycZ (ECDSA)
$ ssh-keygen -E md5 -lf meilihao_github.pub
2048 MD5:4f:32:da:5c:d2:4c:25:a4:ea:dd:08:c9:aa:31:dc:22 563278383@qq.com (RSA) # 即github.com SSH keys上显示的Fingerprint
```

> `ssh-keygen -lf`也适用于known_hosts和authorized_keys文件

重新生成server host key:
逐个替换`/etc/ssh/ssh_host_xxx`或使用`dpkg-reconfigure`命令

```
# rm -v /etc/ssh/ssh_host_*
# dpkg-reconfigure openssh-server
```

### 获取sever上openssh的公钥
```
ssh-keyscan -t ed25519 -p 22 xxx.com
```
获取的是`/etc/ssh`下对应类型的公钥`ssh_host_${type}_key.pub`, type有rsa, ecdsa, ed25519.

### xshell导入linux上生成的密钥报`用户密钥导入失败!`
原因: 密码格式是OPENSSH PRIVATE KEY, 而xshell需要RSA PRIVATE KEY.

转换命令: `ssh-keygen -p -N "" -m pem -f /path/to/key`, 转换后的内容会覆盖源文件

> key(带密码)的转换: `ssh-keygen -p -P "old passphrase" -N "new passphrase" -m pem -f path/to/key`

其他方法:
```bash
puttygen id_rsa -o id_rsa.ppk
puttygen id_rsa.ppk -O private-openssh -o id_rsa.pem
puttygen id_rsa.pem -L # 查看结果
```
或直接:
```bash
puttygen id_rsa -O private-openssh -o id_rsa.pem # ubuntu 20.04 不可行, 但上面的方法可以
```

### ssh root 密码无法登入
env: openssh 7.4, 且authorized_keys已添加pub key

当前配置:
```bash
PermitRootLogin yes
StrictModes yes
PubkeyAuthentication yes
PasswordAuthentication yes
AuthenticationMethods publickey,password
```

直接使用`ssh root@192.168.1.2`直接被拒绝; 使用`ssh -i <private key> root@192.168.1.2`时出现密码提示, 输入密码后成功登入.

修改方法: 注释`AuthenticationMethods`并重启ssh, `ssh -i <private key> root@192.168.1.2`直接登入成功.

openssh 9.2没有AuthenticationMethods, 保证其他如`当前配置`, 即可使用root+password登入.

其实root登入关键是配置PermitRootLogin, PasswordAuthentication.

### 无法使用密码登入
开启`PasswordAuthentication yes`并重启sshd

### sshd日志
添加`-d`

注意: 仅用于调试. 因为实际使用中发现, ssh成功登入后会导致sshd退出.

### sshd启动报`sshd: no hostkeys availabled -- exiting`
是ssh.service的`ExecStartPre=sshd -t`抛出的

解决方案: `ssh-keygen -A`