# scp(secure copy)

## 描述

scp命令是一个基于 SSH 协议在网络之间进行安全传输的命令.

scp如果指定了用户名，命令执行后仅需要再输入密码;没有指定用户名，命令执行后需要输入用户名和密码,也可能直接提示"Permission denied",所以**推荐指定用户名**. 更推荐使用基于`~/.ssh/config`的安全密钥验证来传输文件.

> 使用ssh-copy-id -i .ssh/id_rsa.pub  xxx@192.168.x.xxx来上传公钥

## 选项

- -v : 显示详细的处理信息,比如进度等.
- -P : 指定远程主机的 sshd 端口号
- -r : 用于传送文件夹
- -6 : 使用 IPv6 协议

## 例

### 从 本地 复制到 远程:

#### 复制文件

    scp local_file remote_username@remote_ip:remote_folder #仅指定了远程的目录，文件名字不变

或者

    scp local_file remote_username@remote_ip:remote_file #指定了文件名

或者

    scp local_file remote_ip:remote_folder

或者

    scp local_file remote_ip:remote_file

例如:

    scp /home/space/music/1.mp3 root@www.xxx.com:/home/root/others/music
    scp /home/space/music/1.mp3 root@www.xxx.com:/home/root/others/music/002.mp3

#### 复制目录

    scp -r local_folder remote_username@remote_ip:remote_folder

或者

    scp -r local_folder remote_ip:remote_folder

### 从 远程 复制到 本地

从 远程 复制到 本地，只要将 从 本地 复制到 远程 的命令 的 后2个参数 调换顺序 即可；