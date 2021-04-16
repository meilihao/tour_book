# getent

getent命令帮助用户在administrative databases中查找相关信息. administrative databases包括：
- passwd – 可用于确认用户名、用户名、主目录和用户的全名
- group – 系统已知的有关Unix组的所有信息
- services – 系统上配置的所有Unix服务
- networks – 网络信息–您的系统所属的网络
- protocols –你的系统知道的关于网络协议的一切

## example
```bash
$ getent hosts ubuntu # 查找hostname对应的IP
127.0.1.1       ubuntu
192.168.0.2     ubuntu
$ getent hosts myhost.mydomain.com # 执行反响DNS查询（即根据域名查找对应IP）
15.77.3.40       myhost.mydomain.com myhost
$ getent passwd greys # 根据用户名查找UID
greys:x:1000:1000:Gleb Reys,,,:/home/greys:/bin/bas
$ getent passwd 1000 # 根据UID查找用户名
greys:x:1000:1000:Gleb Reys,,,:/home/greys:/bin/bash
$ getent passwd `whoami` # 获取当前登陆用户的信息
root:x:0:0:root:/root:/bin/bash
$ getent services 22 # 查找那个服务在使用特定端口
ssh                   22/tcp
```

## FAQ
### getent passwd支持ldap
```bash
$ vim /etc/nsswitch.conf
...
passwd: compat systemd files ldap
group: compat systemd files ldap
shadow: compat files ldap
gshadow: files
...
$ sudo apt install nslcd # 安装过程需要输入ldap信息, nslcd配置在`/etc/nslcd.conf`
$ cat /etc/nslcd.conf
...
uri ldap://192.168.0.245:389
base dc=xx,dc=cn
...
log /var/log/nslcd.log debug # 或 log syslog debug
$ sudo systemctl restart nslcd # 存在相近的服务nscd, 不要输错
$ getent passwd chen # 测试是否会输出ldap账户
chen:x:1000:1000::/home/chen:/bin/bash
```