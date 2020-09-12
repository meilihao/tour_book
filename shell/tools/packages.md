# packages.md

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

## FAQ
### Ubuntu 如何解压 rpm 文件
```
# rpm2cpio xxx.rpm | cpio -div # `apt install rpm2cpio`
```