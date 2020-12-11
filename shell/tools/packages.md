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

###　构建源
- [CentOS 的 SRPM](http://vault.centos.org/)
- [Fedora Package Sources](https://src.fedoraproject.org/)
- [Arch Package](https://www.archlinux.org/packages/)

## FAQ

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