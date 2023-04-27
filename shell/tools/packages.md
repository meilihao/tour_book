# packages.md

## 源
- [neokylin v10](http://update.cs2c.com.cn:8080/NS/V10/V10SP2/os/adv/lic/)

    kylin具体版本信息在`/etc/.kyinfo`里

## tools
- [command-not-found](https://command-not-found.com/)
- [cmd package search by man](http://manpages.ubuntu.com/manpages/focal/man8/), 比如[cgdisk](http://manpages.ubuntu.com/manpages/focal/man8/cgdisk.8.html)
- [pkgs.org](https://pkgs.org)
- rpm

    - [rpmfind](https://rpmfind.net/)

        能获取到`src.rpm`
- [repology.org/projects](https://repology.org/projects/)
- [Arch Linux 的最佳 GUI 包管理器](https://linux.cn/article-15727-1.html)

### 构建源
- [CentOS 的 SRPM](http://mirror.centos.org/)
- [Fedora Package Sources](https://src.fedoraproject.org/)

    点击每个项目主页的`Stable version`的列表链接或`Builds Status`的链接即可在里面找到`src.rpm`
- [Arch Package, arch官方repo](https://www.archlinux.org/packages/)
- [aur.archlinux.org - AUR=Arch User Repository, 创建 AUR 的初衷是方便用户维护和分享新软件包，并由官方定期从中挑选软件包进入 community 仓库](https://aur.archlinux.org/packages/)

### 其他打包工具
- [Arch Linux 软件包制作入门](https://linux.cn/article-13843-1.html)
- [`.run/.bin`制作 - makeself](https://github.com/megastep/makeself)
- [flatpak](https://linux.cn/article-15007-1.html)

windows:
- [NSIS](https://nsis.sourceforge.io/Download)

    - [安装程序打包工具NSIS](https://wuziqingwzq.github.io/other/2018/01/08/NSIS1.html)
- [HM NIS EDIT](http://hmne.sourceforge.net/)
    NSIS脚本编辑器

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

## releases
- [ubuntu](https://wiki.ubuntu.com/Releases)

## 大版本升级
- ubuntukylin

    ```bash
    $ sudo apt update
    $ sudo apt upgrade
    $ sudo do-release-upgrade --allow-third-party -d
    ```
## arch安装
[`archinstall`](https://www.debugpoint.com/archinstall-guide/)

## FAQ
### do-release-upgrade 离线升级
env: 有apt repo proxy(且支持ubuntu jammy, 比如`nexus repository manager`是可以限制代理的Ubuntu版本的), 但无法访问changelogs.ubuntu.com

do-release-upgrade 会读取文件 /etc/update-manager/meta-release 以查找发布信息的meta

离线升级方法, 这里以20.04->22.04举例:
1. `do-release-upgrade -d`

    报无法获取`https://changelogs.ubuntu.com/meta-release-lts-development`, 将其手动下载到内网并保存到`/etc/update-manager/meta-release-lts-development`


    推测生效位置:
    ```bash
    # vim /usr/lib/ubuntu-release-upgrader/check-new-release
    m = MetaReleaseCore(useDevelopmentRelease=options.devel_release,
                      useProposed=options.proposed_release)
    # this will timeout eventually
    m.downloaded.wait()
    ```
1. 修改`/etc/update-manager/meta-release`将`URI_LTS = https://changelogs.ubuntu.com/meta-release-lts`改为`URI_LTS = file:///etc/update-manager/meta-release-lts`
1. 再次执行`do-release-upgrade -d`即可

### do-release-upgrade调试
```bash
export DEBUG_UPDATE_MANAGER=1
do-release-upgrade
```

### ubuntu hwe
hwe: Ubuntu LTS enablement（也叫 HWE 或 Hardware Enablement）stacks 用于支持不断更新的硬件技术, 能够为现存的 Ubuntu LTS 提供更新的内核与图形支持, 适用于桌面版、服务器版甚至 Cloud 版和虚拟镜像.

`apt install linux-generic-hwe-22.04`