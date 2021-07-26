# acl
- [Access Control Lists in Linux](https://documentation.suse.com/sles/12-SP4/html/SLES-all/cha-security-acls.html)
- [POSIX Access Control Lists on Linux](https://www.usenix.org/legacy/publications/library/proceedings/usenix03/tech/freenix03/full_papers/gruenbacher/gruenbacher_html/main.html)

配置fs 的acl.

ACL是由一系列的Access Entry所组成的，每一条Access Entry定义了特定的类别可以对文件拥有的操作权限.
Access Entry有三个组成部分：
- Entry tag type
- qualifier (optional) : 定义了特定用户和拥护组对于文件的权限. 只有user和group才有qualifier，其他的都为空.
- permission

Entry tag type它有以下几个类型：
- ACL_USER_OBJ： 	相当于Linux里file_owner的permission
- ACL_USER： 	    定义了额外的用户可以对此文件拥有的permission
- ACL_GROUP_OBJ： 	相当于Linux里group的permission
- ACL_GROUP： 	    定义了额外的组可以对此文件拥有的permission
- ACL_MASK： 	    定义了ACL_USER,ACL_GROUP_OBJ和ACL_GROUP的最大权限
- ACL_OTHER： 	    相当于Linux里other的permission

> mask即ACL_MASK，　｀#effective:...｀即当前mask限制后的权限
> 当权限位只包含"-"时，可用`-`代替`---`
> 使用`-m`时， 指定的user，group必须存在，否则报错
> [POSIX/NFSv4 ACL注意事项](https://www.alibabacloud.com/help/zh/doc-detail/143242.htm)
> [使用POSIX ACL进行权限管理](https://help.aliyun.com/document_detail/143010.html)

ACL 定义的权限是 ugo 权限的超集.
ACL 权限与 ugo 权限的对应关系:
- 文件的 owner 权限对应于 ACL 权限中的 ACL_USER_OBJ 条目。
- **当 ACL 权限中具有 ACL_MASK 条目时，文件的 group 权限对应于 ACL 权限中的 ACL_MASK 条目。否则，当 ACL 权限中具没有 ACL_MASK 条目时，文件的 group 权限对应于 ACL 权限中的 ACL_GROUP_OBJ 条目**.

	acl mask是ACL的默认的权限掩码, 是ACL_USER(named user), ACL_GROUP(named group), ACL_GROUP_OBJ(owning group)能够被授予的最大权限.
	effective(实际权限) = 当前账号对应权限 & acl mask.


- 文件的 other 权限对应于 ACL 权限中的 ACL_OTHER_OBJ 条目。

posix acl按以下顺序检查：所有者、命名用户、拥有组或命名组、其他.

文件的 ugo 权限总是与对应的 ACL 条目保持一致。修改文件的 ugo 权限会导致修改相关的 ACL 条目，同样的，修改这些 ACL 条目会导致修改对应的 guo 权限:
- 当添加了 ACL_USER 或 ACL_GROUP 后，必须有一个对应的 ACL_MASK 条目. 在当前的情况下，ACL_MASK 是被自动创建的，它的权限被设置成了 group 的权限
- 在有ACL或ACL未被完全清空的情况下，改变mask就是改变属组的权限，改变属组的权限就是改变mask.

### 新建文件的 default ACL
一个文件的 access ACL 会在通过 creat()、mkdir()、mknod()、mkfifo() 和 open() 函数创建该文件时被初始化.

如果一个目录被设置了 default ACL，那么将会由文件创建函数的 mode 参数和目录的 default ACL 共通决定新文件的 ACL 权限(此时 umask 被忽略)：
- 新的文件继承父目录的 default ACL 作为自己的 access ACL
- 修改与 ugo 权限对应的 access ACL 条目，使其不包含文件创建函数的 mode 参数不包含的权限

如果一个目录没有被设置 default ACL，那么将由文件创建函数的 mode 参数和 umask 共同决定新文件的 ACL 权限：
- 新建文件的 access ACL 包含 ACL_USER_OBJ, ACL_GROUP_OBJ, 和 ACL_OTHER 条目。这些条目的权限被设置为由 umask 决定的权限
- 修改与 ugo 权限对应的 access ACL 条目，使其不包含文件创建函数的 mode 参数不包含的权限

### ACL MASK
[为什么需要 ACL_MASK 条目](https://www.cnblogs.com/sparkdev/p/9694015.html).


## setfacl
参考:
- [ACL 是什么](https://www.cnblogs.com/sparkdev/p/5536868.html)
- [Linux ACL 权限之进阶篇, 有acl mask的详细说明](https://www.cnblogs.com/sparkdev/p/9694015.html)
- [setfacl命令说明](https://wangchujiang.com/linux-command/c/setfacl.html)
- [FilePermissionsACLs](https://help.ubuntu.com/community/FilePermissionsACLs)
- [Handling ACLs](https://www.pks.mpg.de/~mueller/docs/suse10.2/html/opensuse-manual_en/manual/sec.acls.handle.html)
- [权限检查算法 ACCESS CHECK ALGORITHM](http://man7.org/linux/man-pages/man5/acl.5.html)

### 格式

	setfacl [-bkRd] [{-m|-x} acl参数] 文件/目录名

### 选项
- -m ：配置后面的 acl 参数给文件/目录使用，不可与 -x 合用
- -n : 不重新计算mask. setfacl默认在未明确给出mask时会重新计算acl mask. mask就是ACL的默认的权限掩码
- -x ：删除后续的 acl 参数，不可与 -m 合用
- -b ：移除所有的 ACL 配置参数, 包括default, 会递归清除acl.
- -k ：移除默认的 ACL 参数
- -R ：递归配置 acl
- -d ：配置“默认 acl 参数”，**只对目录有效**，在该目录新建的数据会引用此默认值. 它能让我们创建的子文件或者子文件夹继承父文件夹的权限设置.
- --set : 会先清除掉原有的 ACL 权限，然后添加新的权限. 需要注意的是**一定要包含UGO的设置**，不能象-m一样只是添加ACL就可以了
- --test: 测试模式，不会改变任何文件的acl规则，操作后的acl规则将被列出
- --restore=file: 从文件恢复备份的acl规则，此参数只能独立执行，除了--test,备份的acl文件可由 getfacl -R产生

## example
```
# su admini -c 'touch dir0/file' # for 测试权限
# apt install acl
# getfacl test # 查看acl
# getfacl --omit-header ./test.sh
# setfacl -b . # 移除所有acl
# setfacl -k . # 移除所有default acl
# setfacl -m m::--- . # 修改mask
# setfacl [-R] -m u:zhangy:rw- test    #  添加/修改一个用户权限, `-R`:递归修改
# setfacl -m u::r-- a # 没有指定用户时即修改文件所有者的权限
# setfacl -m g:zhangying:r-w test      # 添加/修改一个组权限
# setfacl -x u:tank test    # 清除tank用户在test文件的acl规则
# setfacl -x d:u:tank test # 清除tank用户在test文件的default acl规则
# setfacl -m d:u:user1:rwx /test <=> setfacl -d -m u:user1:rwx /test # Default ACL是指对于一个目录进行Default ACL设置，并且在此目录下建立的文件都将继承此目录的ACL
# setfacl --set u::rw,u:testu1:rw,g::r,o::- file1 # --set选项会把原有的ACL项都删除，用新的替代(此时会设置mask)
# ### 禁用对用户组的自动授予权限
# setfacl -m group::--- /srv/samba/example/
# setfacl -m default:group::--- /srv/samba/example/
# --- 让目录中创建的新文件系统对象继承权限, 需设置default
# setfacl -m default:group:"DOMAIN\Domain Admins":rwx /srv/samba/example/
# setfacl -m default:group:"DOMAIN\Domain Users":r-x /srv/samba/example/
# setfacl -m default:other::--- /srv/samba/example/
# ### 备份和恢复 ACL 权限
# getfacl -R acldir > acldir.acl
# setfacl -R -b acldir
# setfacl --restore acldir.acl
```

## FAQ
### `setfacl: Option -m: Invalid argument near character 5`
需设置的user/group不存在

### [File system and ACL support](https://www.ibm.com/support/knowledgecenter/en/SSEQVQ_8.1.7/client/c_bac_aclsupt.html)
参考:
- [NAS NFS ACL by aliyun](https://www.alibabacloud.com/help/zh/doc-detail/143242.htm?spm=a2c63.p38356.b99.25.1e294085E97WtS)

**POSIX ACL是NFSv3协议能够扩展支持的权限控制协议. POSIX ACL对mode权限控制进行了扩展**，能够对owner、group、other以外的特定用户和群组设置权限，也支持权限继承. 详细介绍请参见[acl - Linux man page](https://linux.die.net/man/5/acl?spm=a2c63.p38356.879954.5.1a20545b5bjSCq).
NFSv4 ACL是NFSv4协议能够扩展支持的权限控制协议，提供比POSIX ACL更细粒度的权限控制. 详细介绍请参见[nfs4_acl - Linux man page](https://linux.die.net/man/5/nfs4_acl?spm=a2c63.p38356.879954.6.1a20545b5bjSCq).

使用NFSv3协议挂载含有NFSv4 ACL的文件系统，挂载后NFSv4 ACL会被转化为POSIX ACL. 同时也可以用NFSv4协议挂载含有POSIX ACL的文件系统，挂载后POSIX ACL会被转化为NFSv4 ACL. 但由于NFS4 ACL和POSIX ACL并不完全兼容，加上mode和ACL之间的互操作也无法尽善尽美，另外NAS NFSv3挂载不支持锁，所以建议在使用NFS ACL功能时尽量只使用NFSv4协议挂载并设置NFS4 ACL，不使用mode和POSIX ACL.

> 强烈建议使用NFSv4 ACL之后请勿使用mode. 禁用mode: `chmod -R 000 xxx`
> NAS NFSv4 ACL只支持Allow不支持Deny，所以建议将everyone的权限设置到最低，因为被everyone允许的权限对任何用户都适用.

### 查看ext4支持POSIX ACL
```bash
# dumpe2fs -h /dev/sda1|grep -i acl
Default mount options:                 user_xattr    acl
```

### writer1在writer中, 明明有权限却无法进入
```bash
$ cd d2
-bash: cd: d2: 权限不够
$ getfacl d2
# file: d2
# owner: writer1
# group: users
# flags: -s-
user::---
group::---
group:reader:r-x
group:writer:rwx
mask::rwx
other::---
default:user::rwx
default:group::---
default:group:reader:r-x
default:group:writer:rwx
default:mask::rwx
default:other::---

$ id writer1
uid=1001(writer1) gid=100(users) 组=100(users),1001(writer)
```

参考文件权限检查的算法(Access Check Algorithm, from `man acl`), effective user ID 与文件 owner 匹配, 但ACL_USER_OBJ 不包含请求所需的权限, 因此被拒绝.


## NFSv4 ACL
参考:
- [NFSv4 ACL和POSIX ACL相关的特性](https://help.aliyun.com/document_detail/143008.html)
- [NFS 4 ACL Tool](https://www.server-world.info/en/note?os=CentOS_7&p=nfs&f=5)
- [使用NFSv4 ACL进行权限管理](https://www.alibabacloud.com/help/zh/doc-detail/143009.htm?spm=a2c63.p38356.b99.28.40225118iM9BWN)

NFSv4 ACL是目前新的ACL, 比POSIX_ACL功能强大. 目前xfs已支持(, ext4好像需要挂载参数未测试), 仅在nfsv4 client上有用, 因为权限就是用户产生的, 不应该在nfs server端修改.

```sh
apt/yum install nfs4-acl-tools  # Commandline and GUI ACL utilities for the NFSv4 client
```

**nfs4-acl-tools仅在nfsv4 client(即用户用nfsv4挂载的目录上)端有用**, 原因: [Design of the linux NFSv4 ACL implementation, linux服务器导出的nfs文件系统都不支持NFSv4 ACL](http://wiki.linux-nfs.org/wiki/index.php/ACLs) .

> RichACL是linux kernel 对NFSv4 ACL规范的实现. [ext4: Add richacl support patch @  13 Feb 2017 10:34:26 -0500](https://patchwork.kernel.org/patch/9570019/). [xfs: Add richacl support @ 11 Oct 2015 23:24:52 +0000](https://patchwork.kernel.org/patch/7371021/) by `mkfs.xfs -m richacl=1`. 但[**RichACL已中止开发**](https://github.com/andreas-gruenbacher/richacl/issues/9)

> nfs4-acl-tools 解析已有posix acl时会自动将其转成NFSv4 ACL.

> [zfs 还未支持NFSv4](https://github.com/openzfs/zfs/pull/9709)

> Richacls使用ext4扩展文件属性（xattrs）来存储ACL.

> NFS4 ACL和POSIX ACL并不完全兼容.

[cp、tar、rsync工具迁移NFSv4 ACL的方法](https://access.redhat.com/solutions/3628891):
```bash
cp --preserve=xattr
tar --xattrs
rsync -X # requires rsync-3.1.2-10.el7 or later on RHEL 7. RHEL 8 has rsync 3.1.3-4-el8 by default.
```

### acl 自动继承(automatic inheritance)
对目录进行ACL更改后，该更改将传播到启用了自动继承的任何文件或目录下，除非还设置了“ protected”标志. 只要显式设置了文件的ACL或模式，就会设置“保护”标志. 这样可以避免继承覆盖已设置为其他权限的权限.

### 查看文件属性
setfattr可以设置EA(Extended Attributes)，getfattr可以获取EA，attr是一个综合命令，可以设置、获取、删除、列举所有等操作.

```
# apt install attr # 文件系统的挂载选项需要包括user_xattr选项
# getfattr -n system.nfs4_acl . # 通过Extended Attributes查看NFSv4 ACL. `system`是namespace, 文件属性上可能有多个namespace.
# getfattr -m.* -d . # 列举所有属性，包含命名空间
# attr -lq . # 读取所有属性, 但在一些fs(nfs mount with xfs)上不生效, 但getfattr有输出.
```