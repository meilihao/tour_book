# packages.md

## 源
- [neokylin v10](http://update.cs2c.com.cn:8080/NS/V10/V10SP2/os/adv/lic/)

## tools
- [command-not-found](https://command-not-found.com/)
- [cmd package search by man](http://manpages.ubuntu.com/manpages/focal/man8/), 比如[cgdisk](http://manpages.ubuntu.com/manpages/focal/man8/cgdisk.8.html)
- [pkgs.org](https://pkgs.org)
- rpm

    - [rpmfind](https://rpmfind.net/)

### 构建源
- [CentOS 的 SRPM](http://mirror.centos.org/)
- [Fedora Package Sources](https://src.fedoraproject.org/)
- [Arch Package, arch官方repo](https://www.archlinux.org/packages/)
- [aur.archlinux.org - AUR=Arch User Repository, 创建 AUR 的初衷是方便用户维护和分享新软件包，并由官方定期从中挑选软件包进入 community 仓库](https://aur.archlinux.org/packages/)

### 其他打包工具
- [Arch Linux 软件包制作入门](https://linux.cn/article-13843-1.html)
- [`.run/.bin`制作 - makeself](https://github.com/megastep/makeself)

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