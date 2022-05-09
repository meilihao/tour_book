# user

# useradd
创建新的用户.

创建过程:
1. 读取/etc/login.defs和/etc/default/useradd规则
1. 使用自定义规则覆盖默认规则
1. 向/etc/passwd和/etc/group添加用户和用户组.
1. 向/etc/passwd/etc/group添加记录
1. 判断是否在/etc/default/useradd设定的目录下创建用户主目录
1. 复制/etc/skel的全部内容到新用户的主目录

可用`id ${user}`查看用户信息

## 格式
- -c : 注释信息
- -d : 指定用户的家目录（默认为/home/username）
- -D : 查看/修改`/etc/default/useradd`(创建新用户时的一些默认属性)
- -e : 账户的到期时间，格式为 YYYY-MM-DD. 
- -f inactive : 指定账户过期时长后永久停用. 0, 立即停用; -1, 关闭此功能. 默认是-1.
- -u : 指定该用户的默认 UID 
- -g : 指定一个初始的用户基本组(也叫主组)（必须已存在）
- -G : 指定一个或多个扩展用户组(必须已存在, 多个时用`,`分隔)
- -N : 不创建与用户同名的基本用户组
- -s : 指定该用户的默认 Shell 解释器

## example
```
# useradd -g a -G b,c -d /opt/z myname
```

# usermod
修改用户的属性

## 格式
- -c 填写用户账户的备注信息
- -d -m 参数-m 与参数-d 连用，可重新指定用户的家目录并自动把旧的数据转移过去
- -e 账户的到期时间，格式为 YYYY-MM-DD 
- -g 变更所属用户组
- -G 变更扩展用户组(必须已存在, 多个时用`,`分隔), 会覆盖已有的支持组.
- -l : 修改用户的登录名为新名称
- -L 锁定用户禁止其登录系统
- -U 解锁用户，允许其登录系统
- -s 变更默认终端
- -u 修改用户的 UID 

## example
```sh
$ usermod -G root chen # 将用户 chen 加入到 root 用户组中
```

# userdel
删除用户, 会自动清理其已加入的group.

## 格式
- -f 强制删除用户
- -r 同时删除用户及用户主目录

# id
获取帐号自身信息及组信息

## example
```bash
$ id -nG chen # 获取用户的支持组
```

## FAQ
### 多机器账户同步
```
/etc/subuid
/etc/subgid
/etc/passwd
/etc/group
/etc/shadow
/etc/gshadow
```

### passwd: unrecognized option '--stdin'
`echo "username:cleartext_password" | sudo chpasswd`

### 用户名限制
- [by systemd]https://systemd.io/USER_NAMES/)
- see useradd/groupadd

### 删除用户报`userdel: user xxx is currently used by process 4655`
因为有进程在占用着用户, 导致无法删除，解决办法就是用`ps -u +用户名`找出占用的进程, kill掉之后在执行删除命令.

或**推荐**使用`systemctl stop user@<uid>`一步到位.