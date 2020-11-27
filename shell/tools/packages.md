# packages.md

## 源
- [neokylin v10](http://update.cs2c.com.cn:8080/NS/V10/V10SP1-2500/os/adv/lic/)

## tools
- [command-not-found](https://command-not-found.com/)
- [cmd package search by man](http://manpages.ubuntu.com/manpages/focal/man8/), 比如[cgdisk](http://manpages.ubuntu.com/manpages/focal/man8/cgdisk.8.html)

## 打包(source -> deb/rpm)
- [BuildingTutorial](https://wiki.debian.org/BuildingTutorial#Building_the_modified_package)
- [How RPM packages are made: the source RPM](https://fedoramagazine.org/how-rpm-packages-are-made-the-source-rpm/)
- [重建一个源代码 RPM](https://wiki.centos.org/zh/HowTos/RebuildSRPM)

###　构建源
- [CentOS 的 SRPM](http://vault.centos.org/)
- [Fedora Package Sources](https://src.fedoraproject.org/)
- [Arch Package](https://www.archlinux.org/packages/)

## rpm
```bash
# rpm -q --provides openssl-libs | grep libcrypto.so.10 # 查看openssl-libs中的libcrypto.so.10版本
# rpm -i --nodeps xxx.rpm # `--nodeps`安装时不检查依赖
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