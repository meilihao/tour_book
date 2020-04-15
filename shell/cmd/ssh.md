# ssh
SSH（Secure Shell）是一种能够以安全的方式提供远程登录的协议，也是目前远程管理
Linux 系统的首选方式.

sshd 是基于SSH协议(目前推荐版本是v2)开发的一款远程管理服务程序，不仅使用起来方便快捷，而且能够提供两种安全验证的方法：
- 基于口令的验证—用账户和密码来验证登录
- 基于密钥的验证—需要在本地生成密钥对，然后把密钥对中的公钥上传至服务器，并与服务器中的公钥进行比较；该方式相较来说更安全

sshd 服务的配置信息保存在/etc/ssh/sshd_config中.

sshd 服务配置文件中包含的参数以及作用:
参数 作用
Port 22 默认的 sshd 服务端口
ListenAddress 0.0.0.0 设定 sshd 服务器监听的 IP 地址
Protocol 2 SSH 协议的版本号
HostKey /etc/ssh/ssh_host_key SSH 协议版本为 1 时，DES 私钥存放的位置
HostKey /etc/ssh/ssh_host_rsa_key SSH 协议版本为 2 时，RSA 私钥存放的位置
HostKey /etc/ssh/ssh_host_dsa_key SSH 协议版本为 2 时，DSA 私钥存放的位置
PermitRootLogin yes 设定是否允许 root 管理员直接登录, 推荐使用`no`
StrictModes yes 当远程用户的私钥改变时直接拒绝连接
MaxAuthTries 6 最大密码尝试次数
MaxSessions 10 最大终端数
PasswordAuthentication yes 是否允许密码验证
PermitEmptyPasswords no 是否允许空密码登录（很不安全）