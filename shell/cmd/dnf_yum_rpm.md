# yum
rehat开发的包管理软件, 已被dnf取代.

# rpm
## 常用的 RPM 软件包命令
安装软件的命令格式 rpm -ivh filename.rpm 
升级软件的命令格式 rpm -Uvh filename.rpm 
卸载软件的命令格式 rpm -e filename.rpm 
查询软件描述信息的命令格式 rpm -qpi filename.rpm 
列出软件文件信息的命令格式 rpm -qpl filename.rpm 
查询文件属于哪个 RPM 的命令格式 rpm -qf filename

## 常见的 Yum 命令
yum repolist all 列出所有仓库
yum list all 列出仓库中所有软件包
yum info 软件包名称 查看软件包信息
yum install 软件包名称 安装软件包
yum reinstall 软件包名称 重新安装软件包
yum update 软件包名称 升级软件包
yum remove 软件包 移除软件包
yum clean all 清除所有仓库缓存
yum check-update 检查可更新的软件包
yum grouplist 查看系统中已经安装的软件包组
yum groupinstall 软件包组 安装指定的软件包组
yum groupremove 软件包组 移除指定的软件包组
yum groupinfo 软件包组 查询指定的软件包组信息

## examples
```bash
# --- 查看系统已安装软件相关命令
# rpm -qa # 查询系统已安装的rpm包
# rpm -qf /绝对路径/file_name # 查询系统中一个已知的文件属于哪个rpm包
# rpm -ql 软件名 # 查询已安装的软件包的相关文件的安装路径
# rpm -qi 软件名 # 查询一个已安装软件包的信息
# rpm -qc 软件名 # 查看已安装软件的配置文件
# rpm -qd 软件名 # 查看已安装软件的文档的安装位置
# rpm -qR 软件名 # 查看已安装软件所依赖的软件包及文件

# --- 查看系统未安装软件相关命令
# rpm -qpi rpm包 # 查看软件包的详细信息
# rpm -qpl rpm包 # 查看软件包所包含的目录和文件
# rpm -qpd rpm包 # 查看软件包的文档所在的位置
# rpm -qpc rpm包 # 查看软件包的配置文件（若没有，则标准输出就为空）
# rpm -qpR rpm包 # 查看软件包的依赖关系

# --- 其他
# rpm -i --nodeps xxx.rpm # `--nodeps`安装时不检查依赖
# rpm -q --provides openssl-libs | grep libcrypto.so.10 # 查看openssl-libs中的libcrypto.so.10版本
```

## FAQ
### Ubuntu 如何解压 rpm 文件
```
# rpm2cpio xxx.rpm | cpio -div # `apt install rpm2cpio`
```

### 从rpm提前spec
`rpmrebuild -e -p --notest-install rsyslog-8.39.0-4.el7.x86_64.rpm`

### rpm Build-ID
参考:
- [Releases/FeatureBuildId](https://fedoraproject.org/wiki/Releases/FeatureBuildId)

在新版本的 Fedora 27 以及 Redhat 8 中，增加了对于 build-id 的支持，在使用 rpmbuild 时默认会自动添加，会在 /usr/lib/.build-id 目录下生成相关的文件.

可以通过 --define "_build_id_links none" 参数取消文件的生成.

增加 build-id 的目的是为了可快速找到正确的二进制文件以及 Debuginfo.