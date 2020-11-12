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

## 选项
- -t : 强制分配伪终端. 可以在远程机器上执行任何基于屏幕(screen-based)的程序, 所以非常有用, 例如菜单服务. 并联的 -t 选项(即`-tt`)强制分配终端, 即使没有本地终端.
- -p : 指定端口
- -q : 隐藏ssh自身输出, 比如"connection to xx.xxx.xx.xxx closed"
- -o ConnectTimeout=3 : 3s超时

## example
```bash
ssh root@192.168.16.40 -t "cd /proc/cpuinfo"
ssh root@192.168.16.40 bash -c "cd /proc/cpuinfo" # 推荐, 因为ssh执行`-t "ls -lL /home/ubuntu | awk '{print $9}'"`和远端直接执行的输出不一致
```

## FAQ
### 不检查host key, 即不检查fingerprint
仅在安全网络下这样配置, 比如内网.

```bash
# ssh -o "UserKnownHostsFile=/dev/null" -o "StrictHostKeyChecking=no" user@host
# vim ~/.ssh/config
Host *
   StrictHostKeyChecking no # 初次连接时不检查host key, 但该主机的公钥还是会追加到文件 ~/.ssh/known_hosts 中
   UserKnownHostsFile=/dev/null # host key因服务器系统重装，服务器间IP地址交换，DHCP，虚拟机重建，中间人劫持等出现变更也不提示. 即不加入KnownHostsFile
```

### 查看ssh-agent已缓存的key
`ssh-add -l`

### ssh登录后立即退出
login script有问题, 登录时禁止执行即可: `ssh -t user@host bash --noprofile`.