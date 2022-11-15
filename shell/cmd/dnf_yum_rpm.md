# dnf
ref:
- [Fedora 39 将使用 DNF5 作为默认打包工具](https://www.oschina.net/news/209467/fedora-39-dnf5-plan)

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

# dnf --showduplicates list [available] $pkgname # 显示软件的多版本列表
# dnf list # 列出用户系统上的所有来自软件库的可用软件包和所有已经安装在系统上的软件包
# dnf list | grep mariadb # 也支持显示多版本
# dnf list installed # 已安装包的列表
# dnf list available # 列出来自所有可用软件库的可供安装的软件包
# dnf search httpd # 查找包
# dnf install [--assumeno] httpd # 安装包, `--assumeno`类似dryrun
# dnf reinstall httpd # 重装包
# dnf download httpd # 下载包, 但不安装. 用`dnf download --resolve samba`/`yumdownloader --resolve samba`可同时下载依赖
# dnf download --source httpd # 下载src.rpm
# dnf info httpd # 查看包的详细信息
# dnf remove httpd # 卸载http包
# dnf remove --duplicates # 删除重复软件包的旧版本
# dnf downgrade acpid # 回滚某个特定软件的版本
# dnf autoremove # 去除不要的孤立(不再被其他包依赖)依赖包
# dnf provides /bin/bash # 查找某一文件的提供者, 也可是qmake-qt5之类的文件名即不是路径
# dnf whatprovides "/usr/bin/qemu-kvm" # 查找某一文件的提供者
# dnf whatprovides libmysqlclient* # 查找某一包的提供者
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
# dnf repoquery --list qt5-qtbase-devel # 列出该软件包提供的所有文件
# dnf -qy module disable postgresql # Disable the built-in PostgreSQL module. Module：是代表着一组通常一起安装的RPM包, 一个典型的module包含应用，依赖库，文档库，帮助组件等.
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
ref:
- [YUM COMMAND CHEAT SHEET ](https://access.redhat.com/sites/default/files/attachments/rh_yum_cheatsheet_1214_jcs_print-1.pdf)

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
# rpm2cpio xxx.rpm | cpio -div # `apt install rpm2cpio`, rpm解压
# rpm -q --whatprovides /etc/pki/CA # 查找文件的提供者
# rpm -qf /etc/pki/CA # 同上
# rpm -ivh filename.rpm # 安装软件
# rpm -ivh 源码包名*.src.rpm # 安装至 ~/rpmbuild 目录
# rpm -Uvh filename.rpm # 升级软件, `-U`表示升级, 当已安装的version和要更新的version相同时, 原rpm安装到系统里的文件不变动, 但新rpm里的升级script会被执行.
# rpm -e [--nodeps] [-vvh --test] appname # 卸载软件. `--nodeps`表示不卸载依赖; `--test`=dry run; `-vvh`=detail log
# rpm -i --nodeps xxx.rpm # `--nodeps`安装时不检查依赖
# rpm --reinstall xxx.rpm # 重复安装 from rpm v4.12.0
# rpm -q --provides openssl-libs | grep libcrypto.so.10 # 查看openssl-libs中的libcrypto.so.10版本
# rpm -qp --scripts ./packagecloud-test-1.1-1.x86_64.rpm # 查看preinstall 和 postinstall scripts
# rpm -qpi xxx.rpm # 查看rpm描述信息
# rpm -qpl xxx.rpm # 查看rpm文件信息
# rpm -q --scripts <pkg> # 查看已安装的pkg的rpm scripts
# rpm -e --test centos-release  # 检查谁依赖了这个包
# repoquery -i php-intl # repoquery from yum-utils. 获取包信息, 包括来源repo. `-i`,展示详情
# yum/dnf list installed | grep @epel # 已安装包的来源repo
# dnf repo-pkgs <repoid> list installed # 同上
# createrepo /root/rpmbuild/RPMS/aarch64 # 创建repo
```

# yum
rehat开发的包管理软件, 已被dnf取代.

## examples
```bash
# yum deplist httpd # 列出依赖
# yum repolist all # 列出所有仓库
# yum list all # 列出仓库中所有软件包
# yum info [--showduplicates]  # 软件包名称 查看软件包信息, `--showduplicates`显示软件多版本
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
# yum whatprovides libmysqlclient* # 查找某一包的提供者
# yum whatprovides '*bin/grep' # 查找某一文件的提供者
# yum localinstall *.rpm
```

# rpmlint
检查rpm spec 的信息，给予提示以及改进，同时也支持对于rpm 文件处理. 比如rpmlint 会在将 ELF 以外的文件保存至 /usr/lib 目录时返回警告.

```bash
yum install -y rpmlint
rpmlint SPECS/rong.spec
rpmlint RPMS/x86_64/dalong-demo-1-1.x86_64.rpm
```

## rpmbuild
参考:
- [RPM打包原理、示例、详解及备查](https://cloud.tencent.com/developer/article/1444873)
- [RPM包制作之rpmbuild命令说明](http://www.someapp.cn/article/3.mhtml)
- [How RPM packages are made: the source RPM](https://fedoramagazine.org/how-rpm-packages-are-made-the-source-rpm/)
- [重建一个源代码 RPM](https://wiki.centos.org/zh/HowTos/RebuildSRPM)
- [RPM打包流程、示例及常见问题](https://bbs.huaweicloud.com/forum/thread-38327-1-1.html)
- [RPM 包的构建 - SPEC 基础知识](https://www.cnblogs.com/michael-xiang/p/10480809.html)
- [RPM 包的构建 - 实例](https://www.cnblogs.com/michael-xiang/p/10500704.html)
- [How to create an RPM package/zh-cn](https://fedoraproject.org/wiki/How_to_create_an_RPM_package/zh-cn)

```bash
# yum install -y rpm-build rpmdevtools
# rpmdev-setuptree # 构建rpm build环境, 默认在 $HOME 目录下多了一个叫做 rpmbuild的目录. rpmbuild默认路径是由在/usr/lib/rpm/macros里的宏变量`%_topdir`来定义. `echo '%_topdir %(echo $HOME)/rpmbuild' > ~/.rpmmacros`可更改该路径.
# cp ~/rpmbuild/
# rpmdev-newspec -o SPECS/xxx.spec 生成SPEC 文件的模板
# cd SPECS
# rpmbuild -bb xxx.spec # 开始构建, 也可使用`rpmbuild -bb xxx.spec --buildroot="xxx" --define "_topdir xxx" --define "centos_version 700"`, 指定define参数: _topdir指定编译的根目录. 指定buildroot原因: 默认的buildroot是BUILDROOT下的子目录, 该目录的文件名包含该os的特定信息.
```

有些项目的xxx.spec不包含编译步骤, 它们是通过组合自编译 + xxx.spec(仅打包, 通过将需要的文件追加到`%files`来指定打包需要的文件)来实现的.

rpmdev-setuptree生成的文件说明:
|默认位置|宏代码|名称|用途|
|~/rpmbuild/SPECS|%_specdir|Spec 文件目录|保存 RPM 包配置（.spec）文件|
|~/rpmbuild/SOURCES|%_sourcedir|源代码目录|保存源码包（如 .tar 包）和所有 patch 补丁|
|~/rpmbuild/BUILD|%_builddir|构建目录|源码包被解压至此，并在该目录的子目录完成编译|
|~/rpmbuild/BUILDROOT|%_buildrootdir|最终安装目录|保存 %install 阶段安装的文件|
|~/rpmbuild/RPMS|%_rpmdir|标准 RPM 包目录|生成/保存二进制 RPM 包|
|~/rpmbuild/SRPMS|%_srcrpmdir|源代码 RPM 包目录|生成/保存源码 RPM 包(SRPM)|

配置在SPEC文件中的，具体来说各个阶段：
|阶段|读取的目录|写入的目录|具体动作|
|%prep|%_sourcedir|%_builddir|读取位于 %_sourcedir 目录的源代码和 patch. 之后，解压源代码至 %_builddir 的子目录并应用所有 patch.|
|%build|%_builddir|%_builddir|编译位于 %_builddir 构建目录下的文件。通过执行类似 ./configure && make 的命令实现。|
|%install|%_builddir|%_buildrootdir|读取位于 %_builddir 构建目录下的文件并将其安装至 %_buildrootdir 目录。这些文件就是用户安装 RPM 后，最终得到的文件。注意一个奇怪的地方: 最终安装目录 不是 构建目录。通过执行类似 make install 的命令实现。|
|%check|%_builddir|%_builddir|检查软件是否正常运行。通过执行类似 make test 的命令实现。很多软件包都不需要此步。|
|bin|%_buildrootdir|%_rpmdir|读取位于 %_buildrootdir 最终安装目录下的文件，以便最终在 %_rpmdir 目录下创建 RPM 包。在该目录下，不同架构的 RPM 包会分别保存至不同子目录， noarch 目录保存适用于所有架构的 RPM 包。这些 RPM 文件就是用户最终安装的 RPM 包。|
|src|%_sourcedir|%_srcrpmdir|创建源码 RPM 包（简称 SRPM，以.src.rpm 作为后缀名），并保存至 %_srcrpmdir 目录。SRPM 包通常用于审核和升级软件包。|

> 特别需要注意的是：%install 部分使用的是绝对路径，而 %file 部分使用则是相对路径.

在 rpmbuild 中，对上表中的每个宏代码都有对应的目录：

|宏代码|名称|默认位置|用途|
|%_specdir|	Spec 文件目录|	~/rpmbuild/SPECS|	保存 RPM 包配置（.spec）文件|
|%_sourcedir|	源代码目录	|~/rpmbuild/SOURCES|	保存源码包（如 .tar 包）和所有 patch 补丁|
|%_builddir|构建目录|~/rpmbuild/BUILD|源码包被解压至此，并在该目录的子目录完成编译|
|%_buildrootdir|最终安装目录|~/rpmbuild/BUILDROOT|保存 %install 阶段安装的文件|
|%_rpmdir|标准 RPM 包目录|~/rpmbuild/RPMS|生成/保存二进制 RPM 包|
|%_srcrpmdir|源代码 RPM 包目录|~/rpmbuild/SRPMS|生成/保存源码 RPM 包(SRPM)|

rpmbuild 的选项:
- -bp 只解压源码及应用补丁, 执行到pre
- -bc 只进行编译, 执行到 build段
- -bi 只进行安装到%{buildroot}, 执行install段
- -bb 只生成二进制 rpm 包
- -bs 只生成源码 rpm 包
- -ba 生成二进制 rpm 包和源码 rpm 包
- -bl 检测有文件没包含
- --target 指定生成 rpm 包的平台，默认会生成 i686 和 x86_64 的 rpm 包，但一般我只需要 x86_64 的 rpm 包
- --buildroot ： 替换%buildroot值

可以先rpmbuild -bp ,再-bc 再-bi 如果没问题，rpmbuild -ba 生成src包与二进制包.

#### spec
参考:
- [How to create an RPM package/zh-cn](https://fedoraproject.org/wiki/How_to_create_an_RPM_package/zh-cn)

```bash
rpmdev-newspec -o name.spec
```

spec:
```conf
Release: 3%{?dist} # 发行编号, 初始值为 1%{?dist}, 每次制作新包时, 应该递增该数字
Version: xxx # 不能包含`-`
%prep
%setup -q                                    # 解压源码并切换到目录
%setup -c -n bareos                          # 解压源码到bareos目录
```

scripts section:
- %pre 安装前执行的脚本
- %post 安装后执行的脚本
- %preun 卸载前执行的脚本
- %postun 卸载后执行的脚本 : 无需定义清理安装文件的脚本, 软件包被卸载时会由包管理器来自动清理(按安装时由包管理器创建的文件记录, 非包管理器创建的文件会被忽略).
- %pretrans 在事务开始时执行脚本
- %posttrans 在事务结束时执行脚本

```
%files # **推荐使用完整的绝对路径避免安装时与其他包的路径冲突**
%defattr (-,root,root,0755)                         ← 设定默认权限
/opt/xxx                                            # 即/opt/xxx及其所有子文件
%config(noreplace) /etc/my.cnf                      ← 表明是配置文件，noplace表示替换文件
%doc %{src_dir}/Docs/ChangeLog                      ← 表明这个是文档
%attr(644, root, root) %{_mandir}/man8/mysqld.8*    ← 分别是权限，属主，属组
%attr(755, root, root) %{_sbindir}/mysqld
```

## FAQ

### 从rpm提前spec
`rpmrebuild -e -p --notest-install rsyslog-8.39.0-4.el7.x86_64.rpm`

### rpm Build-ID
参考:
- [Releases/FeatureBuildId](https://fedoraproject.org/wiki/Releases/FeatureBuildId)

在新版本的 Fedora 27 以及 Redhat 8 中，增加了对于 build-id 的支持，在使用 rpmbuild 时默认会自动添加，会在 /usr/lib/.build-id 目录下生成相关的文件.

可以通过`--define "_build_id_links none"`参数取消文件的生成.

增加 build-id 的目的是为了可快速找到正确的二进制文件以及 Debuginfo.

### dnf install httpd报"nothing provides httpd-mmn"
```bash
# dnf download httpd
# dnf download httpd-devel
# dnf install httpd-*.rpm
```

### rpmbuild error: `installed (but unpackaged) file(s) found`
解决方法有2:
1. 在/usr/lib/rpm/macros文件中有一个定义:`%_unpackaged_files_terminate_build 1`，把1改为0只警告, **推荐**
1. 找到 /usr/lib/rpm/macros 中`%__check_files  /usr/lib/rpm/check-files %{buildroot}`注释掉

### [yum锁定软件版本](https://www.onitroad.com/jc/linux/centos/centos-redhat-fedora-yum-lock-package-version-command.html)
有两种方法：
- 将`--exclude`指令传递给yum命令，以定义要从更新或安装中排除的软件包列表

    ```bash
    # yum --exclude httpd,php xxx
    # cat /etc/yum.repo.d/xxx.repo
    ...
    exclude=python-3*       # Exclude Single Package
    exclude=httpd php       # Exclude Multiple Packages
    ```
- `yum/dnf versionlock`命令版本锁定rpm软件包命令

> [`dnf install python3-dnf-plugin-versionlock`](https://www.getpagespeed.com/server-setup/centos-rhel-8-how-to-prevent-a-package-from-upgrading)

> `yum -y install yum-versionlock`

### rpmbuild报"Arch dependent binaries in noarch package"
在amd64上打包arm64的包(golang程序), spec的`BuildArch`是`noarch`.

解决方法: 在spec文件里`%define _binaries_in_noarch_packages_terminate_build   0`

### rpm安装报`file /opt from install of <xxx.rpm> conflicts with file from package <pkg>`
原spec:
```conf
%files
%defattr(-,root,root,-)
/*
```

改为:
```conf
%files
%defattr(-,root,root,-)
/opt/*
```

### `rpm -i`安装同名软件包的多个版本时会提示"file xxx from install of yyy conflicts with file from package zzz"
使用`rpm -i --force`时, `rpm -qa`会查到多个软件包版本.

使用`yum/dnf install`时会先移除旧软件包再安装.

### `rpm -e`遇到script error
`rpm -e --noscripts xxx`, 删除软件包时不执行script.

### yum/dnf upgrade script执行顺序
参考:
- [升级和安装的rpm过程中 spec 文件中脚本调用顺序和参数](https://blog.csdn.net/kyle__shaw/article/details/115461583)
- [How to execute a script at %pre, %post, %preun or %postun stage (spec file) while installing/upgrading an rpm](https://www.golinuxhub.com/2018/05/how-to-execute-script-at-pre-post-preun-postun-spec-file-rpm/)

在执行这些脚本时，都会有相同的传入值 `$1`, 来判断具体执行的是以下的哪步操作:
![](/misc/img/shell/13-05-20182B15-53-21.png.webp)

安装/升级/删除过程中，具体执行的操作的顺序:
|install| upgrade| un-install|
|pre $1=1 |   %pre $1=2 |  %preun $1=0|
|copy files | copy files | remove files|
|%post $1=1 | %post $1=2 | %postun $1=0|
||%preun $1=1 from old RPM.   ||
||delete files only found in old package  ||
||%postun $1=1 from old RPM.  ||

v1->v2的脚本执行顺序:
1. 执行 v2 的 %pre
1. 释放 v2 中的文件
1. 执行 v2 的 %post
1. 执行 v1 的 %preun
1. 删除 v1 中特有的文件
1. 执行 v1 中的 %postun

### 构建rpm报:`error: cannot open Packages database in /var/lib/rpm`
rpmdb本地数据存储文件损坏.

```bash
# --- 清理YUM仓库本地数据存储文件
mv /var/lib/rpm/__db* /tmp # 或 for i in `ls /var/lib/rpm | grep 'db.'`;do mv $i $i.bak;done 或 rm -f /var/lib/rpm/__db*
# --- 执行如下命令，清理yum缓存
rpm --rebuilddb
yum clean all
```

### dnf remove xxx时依赖xxx的包也被移除了
1. 使用rpm移除(未测试)
1. 修改/etc/dnf/dnf.conf 

    将`clean_requirements_on_remove=True`改为`clean_requirements_on_remove=False`

### oracle linux 7存在两套repo, 包不能混用
from [Oracle Linux 7 Repositories](https://yum.oracle.com/oracle-linux-7.html)

1. latest: 发行版

```
[ol7_optional_latest]
name=Oracle Linux $releasever Optional Latest ($basearch)
baseurl=http://yum.oracle.com/repo/OracleLinux/OL7/optional/latest/$basearch/
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-oracle # 校验文件所使用的公钥
gpgcheck=1 # 是否校验文件
enabled=1
```

2. developer: 技术预览版

```
[ol7_optional_developer]
name=Oracle Linux $releasever Optional developer ($basearch)
baseurl=http://yum.oracle.com/repo/OracleLinux/OL7/optional/developer/$basearch/
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-oracle
gpgcheck=1
enabled=1
```

### 安装spec文件的deps
```bash
yum-builddep my-package.spec [--define "centos_version 790"]
yum-builddep my-package.src.rpm

dnf builddep my-package.spec
dnf builddep my-package.src.rpm

rpmbuild --rebuild xx.src.rpm # build src.rpm
```

### `yum-builddep`时报`cannot install both readline-7 or readline-8`
> 当前系统已安装readline-8

错误提示中没找到具体包冲突, 推测yum-builddep找到了两个版本的readline-devel, 但不知如何选择导致报错.

直接安装readline-devel即可.

### 用yum cache制作离线源
```
cat << EOF >> /etc/yum.conf
cachedir=/var/cache/yum/$basearch/$releasever
keepcache=1   #修改此参数为 1 ，代表打开缓存
debuglevel=2
logfile=/var/log/yum.log
EOF
```

其他方法, 参考[制作离线yum源](https://www.wanpeng.life/1903.html):
```bash
mkdir -p yum_cache
yumdownloader --resolve --destdir=yum_cache docker
createrepo yum_cache
cat <<EOF>> /etc/yum.repos.d/docker.repo
[docker-local]
name=dokcer-ce
baseurl=file:///root/yum_cache
gpgcheck=0
enabled=1
EOF
yum repolist
```

安装单个工具: `yum install -y iotop --downloaddir=/root/centos-repo`

### `%{?dist}`含义
不加问号，如果 dist 有定义，那么就会用定义的值替换，否则就会保 %{dist};
加问好，如果 dist 有定义，那么也是会用定义的值替换，否则就直接移除这个tag %{?dist}

### `rpmbuild -bb rocksdb.spec`报错
env: rocksdb.spec from [fedora 36](https://src.fedoraproject.org/rpms/rocksdb), build on oracle linux 7.9

- `Unknown tag: %forgemeta` : 删除它
- `forgesetup\n... no job control` : 用`%setup -q`代替
- `%{set_build_flags}\n... no job control` : 删除它.

    [set_build_flags用途](https://src.fedoraproject.org/rpms/redhat-rpm-config/blob/rawhide/f/buildflags.md?text=True)

### 清理旧kernel
```bash
# rpm -aq |grep kernel # 查看kernel所有包
# uname -a # 查看正在使用的kernel
# rpm -q kernel # 查看全部kernel
# dnf remove kernel-5.14.0-96.el9.x86_64
# dnf remove kernel-core-5.14.0-96.el9.x86_64
```

一键清理所有非使用的kernel: `dnf remove $(rpm -qa | grep kernel | grep -v $(uname -r))`