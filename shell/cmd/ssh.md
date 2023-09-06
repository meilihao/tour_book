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
ssh root@192.168.16.40 -t "ls -lL /home/ubuntu | awk '{print \$9}'" # 简单命令时推荐, 此时需要转义否则结果会与预期不符.
ssh root@192.168.16.40 'bash -s' < a.sh # 复杂命令时, 推荐
ssh root@192.168.16.40 bash -c "ls -lL /home/ubuntu | awk '{print \$9}'" # 此时的执行结果不正确, 因此不能使用`bash -c "xxx"`的形式
ssh aliyun "nohup sleep 10 &" # 使用nohup也会卡10s, 推测是输出会被写入ssh conn导致
ssh aliyun "nohup sleep 10 >/dev/null 2>&1 &" # 不卡
ssh -p 22 -C -f -N -g -L 9200:192.168.1.19:9200 ihavecar@192.168.1.19 # 将发往本机（192.168.1.15）的 9200 端口访问转发到 192.168.1.19 的 9200 端口
w # 查看ssh进程及其terminal, src ip, cmd
```

## FAQ
### 不检查host key, 即不检查fingerprint
仅在安全网络下这样配置, 比如内网.

```bash
# ssh -o "UserKnownHostsFile=/dev/null" -o "StrictHostKeyChecking=no" -o ConnectTimeout=10 user@host
# vim ~/.ssh/config
Host *
   StrictHostKeyChecking no # 初次连接时不检查host key, 但该主机的公钥还是会追加到文件 ~/.ssh/known_hosts 中
   UserKnownHostsFile=/dev/null # host key因服务器系统重装，服务器间IP地址交换，DHCP，虚拟机重建，中间人劫持等出现变更也不提示. 即不加入KnownHostsFile
```

### 查看ssh-agent已缓存的key
`ssh-add -l`

### ssh登录后立即退出
login script有问题, 登录时禁止执行即可: `ssh -t user@host bash --noprofile`.

### ssh无法输入密码
`sshpass -p 'xxx' ssh root@xxx "ls"`

### ssh获取cmd的exit code
ssh成功连接到remote并执行cmd后, ssh返回的exit code就是cmd执行后的exit code. 验证方法: `ssh aliyun "exit 13"`

### ssh client报`packet_write_wait: Connection to 47.111.xxx.xxx port 22: Broken pipe`/ssh 超时自动断开
只需在客户端设置(**推荐**)：
- 全局设置：`/etc/ssh/ssh_config`
- 当前user设置：~/.ssh/config`

```
# 在ssh_config开头处添加
Host *
    ServerAliveInterval 300
    ServerAliveCountMax 2
...
```

这些设置让ssh client 每5分钟发送一个空包到另一端, 如果它在尝试了2次后，仍没有收到任何响应，则放弃, 即断开连接.

> 也有可能是防火墙掐掉空闲连接导致: [Linux使用ssh超时断开连接的真正原因](http://bluebiu.com/blog/linux-ssh-session-alive.html)

其他方式:
1. 服务端设置
    找到/etc/ssh/sshd_config

    # 30表示30s给客户端发送一次心跳
    ClientAliveInterval 30
    # 此客户端没有返回心跳3次，则会断开连接
    ClientAliveCountMax 3
    # TCP保持连接不断开
    TCPKeepAlive yes



### ssh-add无法添加ed25519 key
```
$ ssh-add ~/.ssh/my_ed25519
Enter passphrase for /home/chen/.ssh/my_ed25519:
Bad passphrase, try again for /home/chen/.ssh/my_ed25519:
Could not add identity "/home/chen/.ssh/my_ed25519": communication with agent failed

$ ssh -V
OpenSSH_7.5p1 Debian-5, OpenSSL 1.0.2l  25 May 2017
$ echo $SSH_AUTH_SOCK
/run/user/1000/keyring/ssh
```

虽然ssh-add无法添加, 但`ssh xxx`还是可正常使用

### 执行`ssh-add`报`Could not open a connection to your authentication agent`
ssh commands need to know how to talk to the ssh-agent, they know that from the SSH_AUTH_SOCK environment variable.

> **一般系统内置的terminal启动时都会注入SSH_AUTH_SOCK和SSH_AGENT_PID, 而tilix之类的就没有**.

> 在Ubuntu下, ssh-agent 通常从 /etc/X11/Xsession.d/90x11-common_ssh-agent 启动, 其`has_option`控制来自`/etc/X11/Xsession.options`

```bash
$ eval "$(ssh-agent -s)" # 直接开启一个ssh-agent进程, 因为它是独立进程，所以即使用户退出当前shell连接，它依然存在. `ssh-agent -s`是启动ssh-agent并输出启用其环境变量的bash script.
```
或者`ssh-agent $SHELL`, 它会在当前的shell中启动一个子shell，ssh-agent程序运行在这个子shell中，退出当前的子shell，ssh-agent会随之消失， 可用pstree验证.

上述两种方法仅对当前执行该命令的terminal有效.

获取SSH_AUTH_SOCK:
```bash
# ssh-agent 
SSH_AUTH_SOCK=/tmp/ssh-iHWAYGzwmxGH/agent.6571; export SSH_AUTH_SOCK;
SSH_AGENT_PID=6572; export SSH_AGENT_PID;
echo Agent pid 6572;
```

这个每次terminal都要运行一遍启动ssh-agent的命令好像有点麻烦，值得高兴的是大多数的linux发行版都在登录图形界面时都会启动一个ssh-agent进程， 用户不需要任何操作，可以使用ps -ef | grep ssh-agent查看. 如果系统没有这个功能，请在~/.xsession文件中加入:
```conf
ssh-agent gnome-session # 请使用当前系统的窗口管理器取代gnome-session, 可用`ps aux|grep -i session`查找对应的session
```

或在`~/.bashrc`添加`export $(gnome-keyring-daemon --start --components=secrets,ssh)`来启用ssh-agent.

**最佳方法是在`/etc/profile`中追加`eval "$(ssh-agent)"`**.

### OpenSSH8无法用root登录
OpenSSH8默认仅支持非root用户用password登录.

```bash
# vim /etc/ssh/sshd_config
...
PermitRootLogin yes
# PermitRootLogin prohibit-password
...
```

PermitRootLogin选项:
- yes : 允许root登录
- prohibit-password : 允许root登录, 但是禁止root用密码登录

### `PermitEmptyPasswords yes`导致ssh连接出现`Authentication failed`
env:
- kylin server v10
- openssh 8.2p1

创建/删除nas(业务代码)的同时用python2 paramiko ssh+password访问节点会报`Authentication failed`, 与同事确认将`PermitEmptyPasswords`设为`no`后问题消失.


推荐使用ssh key取代ssh password来解决.

### go ssh+password连接esxi报`ssh: unable to authenticate, attempted methods [none], no supported methods remain`
ref:
- [SSH in Go: unable to authenticate, attempted methods [none], no supported methods remain](https://stackoverflow.com/questions/47102080/ssh-in-go-unable-to-authenticate-attempted-methods-none-no-supported-method)
- [x/crypto/ssh: failed to connect using keyboard interactive auth method](https://github.com/golang/go/issues/32108)

通过`ssh -T -vvv root@xxx`, 查看esxi ssh server支持的ssh method是keyboard-interactive, 而不是password.