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
- [Debian 新维护者手册](https://www.debian.org/doc/manuals/maint-guide/)

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