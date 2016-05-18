## 用户

用户列表:`/etc/passwd`
格式:`name(用户名):password(口令):id(用户标识号):group id(组标识号):comment(注释性描述):主目录:login shell`

`password`部分被`x`替代表示密码实际存储在`/etc/shadow`.

`login shell`中的`/usr/sbin/nologin`和`/bin/false`:

- `/usr/sbin/nologin`,礼貌地拒绝登录(会显示一条提示信息),但可以使用其他服务,比如ftp.
- `/bin/false`,什么也不做只是返回一个错误状态,然后立即退出.即用户会无法登录,并且不会有任何提示,是最严格的禁止login选项，一切服务都不能用.

> 如果存在`/etc/nologin`文件,则系统**只允许root用户**登录,其他用户全部被拒绝登录,并向他们显示`/etc/nologin`文件的内容

## 用户组

用户列表:`/etc/group`
格式:`group name:password:group id:user list`

`password`部分被`x`替代表示没有密码.

### 命令

- groups : 一个用户属于哪些组.

## 权限

### SUID

SUID权限**仅对二进制可执行文件有效**,让用户在执行时具有文件所有者的权限,且该权限仅在执行该文件的过程中有效.

```
chmod u+s xxx
```

举例: passwd命令

### SGID

#### SGID对目录

- 用户在此目录下的有效用户组将会变成该目录的用户组
- 如果用户在该目录下具有 w 的权限,则其所创建的新文件的用户组与此目录的用户组相同

#### SGID对文件

- SGID 对二进制可执行文件有效
- 用户对该文件具备 x 的权限
- 在执行的过程中将会获得该文件群组的支持

举例: locate命令

### sticky bit

是在other用户的权限上设置,可以理解为防删除位,**仅对目录有效**.
如果具有可执行权限，设置sticky bit后是t；如果没有可执行权限的话，设置sticky bit后是T.

```
chmod o+t xxx
```

- 对一个目录设置了sticky-bit之后,比如`rwxrwxrwt`，存放在该目录的文件仅准许其属主(或root)执行删除、 移动等操作.

举例: /tmp目录

###
