# dnf
rehat开发的新一代包管理软件, 用于取代dnf.

DNF包管理器克服了YUM包管理器的一些瓶颈，提升了包括用户体验，内存占用，依赖分析，运行速度等多方面的内容. DNF使用 RPM, libsolv 和 hawkey 库进行包管理操作.

DNF配置文件的位置:
1. Main Configuration: /etc/dnf/dnf.conf
1. Repository: /etc/yum.repos.d/
1. Cache Files: /var/cache/dnf

> dnf 命令每次使用时都会更新元信息，所以 update 和 upgrade 子命令是可以互换的. apt保存了一个需要定期更新的缓存信息, 其由update命令更新, 而upgrade命令来更新packages.

## examples
```bash
# dnf --version
# dnf help clean # 获取有关某条命令的使用帮助
# dnf help # 查看所有的 DNF 命令及其用途
# dnf history # 查看 DNF 命令的执行历史
# dnf history info # 显示有关历史的详细信息. 如果未指定，则显示最近一次历史信息
# dnf history info 3 # 查看有关给定ID的历史详细信息
# dnf history redo 3 # 对指定的ID历史操作重复执行
# dnf history undo 3 # 执行与指定历史ID执行的所有操作相反的操作
# dnf history rollback 7 # 撤消在历史ID之后执行的所有操作

# dnf list # 列出用户系统上的所有来自软件库的可用软件包和所有已经安装在系统上的软件包
# dnf list installed # 已安装包的列表
# dnf list available # 列出来自所有可用软件库的可供安装的软件包
# dnf search httpd # 查找包
# dnf install httpd # 安装包
# dnf reinstall httpd # 重装包
# dnf download httpd # 下载包, 但不安装
# dnf info httpd # 查看包的详细信息
# dnf remove httpd # 卸载http包
# dnf remove --duplicates # 删除重复软件包的旧版本
# dnf downgrade acpid # 回滚某个特定软件的版本
# dnf autoremove # 去除不要的孤立(不再被其他包依赖)依赖包
# dnf provides /bin/bash # 查找某一文件的提供者
# dnf install /path/to/file.rpm
# dnf install https://xyz.com/file.rpm

# # dnf mark命令允许始终将指定的程序包保留在系统上，并且在运行自动删除命令时不从系统中删除此程序包
# dnf mark install nano # 指定的软件包标记为由用户安装.
# dnf mark remove nano # 取消将指定的软件包标记为由用户安装

# dnf check-update # 检测系统上的所有系统包的更新
# dnf check-update nano # 检查对指定软件包的更新. args: `--changelog`, 显示changelog
# dnf update # 更新系统中的所有安装包
# dnf list updates # 检查可用更新
# dnf list obsoletes # 列出系统上已安装的已废弃的软件包
# dnf list recent # 列出最近添加到仓库中的软件包
# dnf list autoremove # 列出将被dnf autoremove命令删除的软件包
# dnf update httpd # 更新特定的软件包
# dnf upgrade nano-2.9.8-1 # 将给定的一个或多个软件包升级到指定的版本
# dnf updateinfo list # 检查系统上更新公告的信息
# dnf clean all # 清除所有缓存的软件包
# dnf clean dbcache/metadata/packages # 默认情况下，当执行各种dnf操作时，dnf会将包和存储库元数据之类的数据缓存到`/var/cache/dnf`目录中. 该缓存在一段时间内会占用大量空间。这将允许删除所有缓存的数据(dbcache: 缓存文件; metadata, repo数据; packages, 包)
# dnf distro-sync # 更新软件包到最新的稳定发行版
# dnf upgrade-minimal # 将每个软件包更新为提供错误修正，增强功能或安全修复程序的最新版本
# dnf upgrade-minimal [Package_Name] # 将给定的一个或多个软件包更新为提供错误修正，增强或安全修复的最新版本
# dnf check # 检查本地包装，并生成有关已检测到的任何问题的信息. 可以通过选项限制`packagedb`检查–dependencies，–duplicates，–obsoleted或–provides
# dnf makecache # 用于下载和启用系统上当前启用的仓库的所有数据

# dnf group summary # 显示了系统上已安装并可用的组数量
# dnf group list [-v] # 列出安装组包（Group packages）
# dnf group install 'System Tools' # 安装特定的组包
# dnf group install 'System Tools' # 同上
# dnf group update 'System Tools' # 更新组包
# dnf group update 'System Tools' # 同上
# dnf group remove ‘Educational Software’ # 删除一个软件包组
# dnf group remove 'Development Tools' # 同上
# dnf group info 'Development Tools' # 查看指定的软件包组信息

# dnf repolist # 仅列出系统上可用的仓库. args: `-v`, 显示有关每个存储库的详细信息
# dnf repolist all # 列出所有的仓库(包括不可用)
# dnf repolist enabled # 列出系统上已启用的仓库
# dnf repolist disabled # 列出系统上禁用的仓库
# dnf –enablerepo=epel install phpmyadmin # 从特定的软件包库安装特定的软件
# dnf repoquery htop # 在启用的存储库中搜索给定的程序包并显示信息, 等效于`rpm -q`
```

设置 DNF自动更新:
```bash
# dnf install dnf-automatic # 安装DNF自动更新工具
# # 编辑/etc/dnf/automatic.conf文件并替换apply_updates = yes而不是apply_updates = no. 在配置文件中进行更改后，启用`dnf-automatic-timer`服务
# systemctl enable dnf-automatic.timer
# systemctl start dnf-automatic.timer
```

缺点:
1. 在 DNF 中没有 –skip-broken 命令，并且没有替代命令供选择
1. 在 DNF 中没有判断哪个包提供了指定依赖的 resolvedep 命令
1. 在 DNF 中没有用来列出某个软件依赖包的 deplist 命令
1. 当在 DNF 中排除了某个软件库，那么该操作将会影响到之后所有的操作，不像在 YUM 那样，排除操作只在升级和安装软件时才起作用

# rpm

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
# rpm -ivh filename.rpm # 安装软件
# rpm -Uvh filename.rpm # 升级软件
# rpm -e filename.rpm # 卸载软件
# rpm -i --nodeps xxx.rpm # `--nodeps`安装时不检查依赖
# rpm -q --provides openssl-libs | grep libcrypto.so.10 # 查看openssl-libs中的libcrypto.so.10版本
```

# yum
rehat开发的包管理软件, 已被dnf取代.

## examples
```bash
# yum repolist all # 列出所有仓库
# yum list all # 列出仓库中所有软件包
# yum info # 软件包名称 查看软件包信息
# yum install # 软件包名称 安装软件包
# yum reinstall # 软件包名称 重新安装软件包
# yum update # 软件包名称 升级软件包
# yum remove # 软件包 移除软件包
# yum clean all # 清除所有仓库缓存
# yum check-update # 检查可更新的软件包
# yum grouplist # 查看系统中已经安装的软件包组
# yum groupinstall # 软件包组 安装指定的软件包组
# yum groupremove # 软件包组 移除指定的软件包组
# yum groupinfo # 软件包组 查询指定的软件包组信息
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