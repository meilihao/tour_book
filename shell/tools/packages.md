# packages.md

## 源
- [neokylin v10](http://update.cs2c.com.cn:8080/NS/V10/V10SP1-2500/os/adv/lic/)

## tools
- [command-not-found](https://command-not-found.com/)
- [cmd package search by man](http://manpages.ubuntu.com/manpages/focal/man8/), 比如[cgdisk](http://manpages.ubuntu.com/manpages/focal/man8/cgdisk.8.html)
- [pkgs.org](https://pkgs.org)

## 打包(source -> deb/rpm)
- [BuildingTutorial](https://wiki.debian.org/BuildingTutorial#Building_the_modified_package)
- [How RPM packages are made: the source RPM](https://fedoramagazine.org/how-rpm-packages-are-made-the-source-rpm/)
- [重建一个源代码 RPM](https://wiki.centos.org/zh/HowTos/RebuildSRPM)
- [Easy way to create a Debian package and local package repository](https://linuxconfig.org/easy-way-to-create-a-debian-package-and-local-package-repository)
- [apt source包打包](https://www.debian.org/doc/manuals/apt-howto/ch-sourcehandling.zh-cn.html)
- [Debian 新维护者手册](https://www.debian.org/doc/manuals/maint-guide/)
- [RPM打包流程、示例及常见问题](https://bbs.huaweicloud.com/forum/thread-38327-1-1.html)
- [RPM 包的构建 - 实例](https://www.cnblogs.com/michael-xiang/p/10500704.html)

###　构建源
- [CentOS 的 SRPM](http://mirror.centos.org/)
- [Fedora Package Sources](https://src.fedoraproject.org/)
- [Arch Package, arch官方repo](https://www.archlinux.org/packages/)
- [aur.archlinux.org - AUR=Arch User Repository, 创建 AUR 的初衷是方便用户维护和分享新软件包，并由官方定期从中挑选软件包进入 community 仓库](https://aur.archlinux.org/packages/)


## dpkg-buildpackage
选项:
- -nc : doesn't call the clean target, 因此无需重新编译(但可能还是需要编译少量内容)
- -uc : don't sign the changes file
- -us : unsigned source package

## deb debug package
打deb包时, 通过将可执行文件的符号表通过剥离成独立的 dbg 包, 称为 debug package. 正常情况下 -dbg.deb 不会安装.

新版本（debhelper/9.20151219 or newer in Debian）的 debhelper 已经把 -dbg.deb 改为 -dbgsym.deb，详情请见[DebugPackage](https://wiki.debian.org/DebugPackage).

### rpmbuild
参考:
- [RPM打包原理、示例、详解及备查](https://cloud.tencent.com/developer/article/1444873)

```bash
# yum install -y rpm-build rpmdevtools
# rpmdev-setuptree # 构建rpm build环境, 默认在 $HOME 目录下多了一个叫做 rpmbuild的目录
# cp ~/rpmbuild/
# rpmdev-newspec -o SPECS/xxx.spec 生成SPEC 文件的模板
# cd SPECS
# rpmbuild -bb xxx.spec # 开始构建, 也可使用`rpmbuild -bb xxx.spec --define "centos_version 700"`指定define参数
```

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

rpmbuild 的选项:
- -bp 只解压源码及应用补丁, 执行到pre
- -bc 只进行编译, 执行到 build段
- -bi 只进行安装到%{buildroot}, 执行install段
- -bb 只生成二进制 rpm 包
- -bs 只生成源码 rpm 包
- -ba 生成二进制 rpm 包和源码 rpm 包
- -bl 检测有文件没包含
- --target 指定生成 rpm 包的平台，默认会生成 i686 和 x86_64 的 rpm 包，但一般我只需要 x86_64 的 rpm 包

可以先rpmbuild -bp ,再-bc 再-bi 如果没问题，rpmbuild -ba 生成src包与二进制包.

## FAQ
### 通过deb-src构建deb
```bash
# vim /etc/apt/source.list # 添加deb-src源
deb-src http://pl.archive.ubuntu.com/ubuntu/ natty main restricted
# apt update
# apt build-dep ccache # 安装构建ccache所需的依赖
# apt-get -b source ccache # 获取ccache源码并构建
# dpkg -i ccache*.deb
```

### dpkg-buildpackage报Warning "Compatibility levels before 9 are deprecated"
将项目`debian/compat`中的数字改为9即可

### rpmbuild error: `installed (but unpackaged) file(s) found`
解决方法有2:
1. 在/usr/lib/rpm/macros文件中有一个定义:`%_unpackaged_files_terminate_build 1`，把1改为0只警告, **推荐**
1. 找到 /usr/lib/rpm/macros 中`%__check_files  /usr/lib/rpm/check-files %{buildroot}`注释掉

## desktop
### UbuntuDDE
参考:
- [Install Deepin Desktop Environment on Ubuntu 20.04](https://computingforgeeks.com/install-deepin-desktop-environment-on-ubuntu/)

```bash
$ sudo add-apt-repository ppa:ubuntudde-dev/stable
$ sudo apt update
$ sudo apt install ubuntudde-dde
$ sudo reboot
```