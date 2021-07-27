# apt

## 描述

debian,ubuntu等发行版的包管理.

> apt删除一个包并不会删除**已修改的用户配置文件**, 以防用户意外删除了包. 如果想通过 apt 命令删除一个应用及其配置文件，请在之前删除过的应用程序上使用 purge 子命令.

## 格式

## 例
```bash
# apt search drbd-utils # 类似`apt-cache search`, 但更好用
# apt-cache madison pouch # 列出软件包的所有版本
# apt install pouch=1.0.0-0~ubuntu # 安装指定版本的软件包
# apt-get install --reinstall aptitude # 重新安装package
# apt-get install --only-upgrade samba # 仅更新单个package
# apt list -a cifs-utils # package all version
# apt-cache policy cifs-utils # package all version, 推荐
# rmadison cifs-utils # package all version, 推荐
# apt-cache depends -i samba # 查看依赖
# apt install --simulate samba # 仅模拟不安装
# apt install --download-only samba # 仅下载, 不安装
# apt list --installed # 查看已安装的package
# apt-cache show cpio # 查看软件依赖
# apt purge package_name # remove命令卸载指定软件包，但是留下一些包文件. 如果想彻底卸载软件包，包括它的文件，使用purge替换remove
# apt build-dep xxx # 获取构建包xxx的依赖(包括源码), 前提是取消`/etc/apt/sources.list*`中相应deb-src源的注释
```

> apt-file也可用于查找文件

# dpkg

## example
```
$ sudo dpkg -i --force-bad-verify  acl_2.2.52-3_amd64.deb # 跳过签名验证, `--ignore-depends=<x1>,<x2>`忽略依赖
$ dpkg -S file # 这个文件属于哪个已安装软件包
$ dpkg -L package # 列出软件包中的所有文件
$ dpkg -s package # 列出软件包中的描述信息
$ echo "PACKAGE hold" | sudo dpkg --set-selections  ##锁定软件包
$ dpkg --get-selections | grep hold  ##显示锁定的软件包列表
$ echo "PACKAGE install" | sudo dpkg --set-selections  ##解除对软件包的锁定
$ dpkg --info xxx.deb | grep Depends # 查看deb包的依赖
$ dpkg-deb -c xxx.deb # 查看deb内容
```

按文件搜索package也可直接使用[debian package服务](https://www.debian.org/distrib/packages)


## deb
参考:
- [Ubuntu下制作deb包的方法详解](https://blog.csdn.net/gatieme/article/details/52829907)
- [对一个deb包的解压、修改、重新打包全过程方法](https://blog.csdn.net/yygydjkthh/article/details/36695243), 重新打包后apt install遇到`Failed to fetch ....deb  Size mismatch`, 改用dpkg安装即可, 因为"Packages.gz"里的元信息是旧的.


deb包本身有三部分组成：
- 数据包，包含实际安装的程序数据，文件名为 data.tar.XXX

    data.tar.gz包含的是实际安装的程序数据，而在安装过程中，该包里的数据会被直接解压到根目录(即 / )，因此在打包之前需要根据文件所在位置设置好相应的文件/目录树
- 安装所需的信息及控制脚本包, 包含deb的安装说明，标识，脚本等，文件名为 control.tar.gz

    一般有 5 个文件：
    控制文件    描述
    control     用了记录软件标识，版本号，平台，依赖信息等数据
    preinst     在解包data.tar.gz前运行的脚本
    postinst    在解包数据(即安装)后运行的脚本
    prerm   卸载时, 在删除文件之前运行的脚本
    postrm  卸载时, 在删除文件之后运行的脚本

- 最后一个是deb文件的一些二进制数据，包括文件头等信息，一般看不到，在某些软件中打开可以看到

deb本身可以使用不同的压缩方式. tar格式并不是一种压缩格式，而是直接把分散的文件和目录集合在一起，并记录其权限等数据信息. 之前提到过的 data.tar.XXX，这里 XXX 就是经过压缩后的后缀名. deb默认使用的压缩格式为gzip格式，所以最常见的就是 data.tar.gz. 常有的压缩格式还有 bzip2 和 lzma，其中 lzma 压缩率最高，但压缩需要的 CPU 资源和时间都比较长.

### dpkg重新打包
```bash
#解压出包中的文件到extract目录下
$ dpkg -X ../openssh-client_6.1p1_i386.deb extract/
#解压出包的控制信息extract/DEBIAN/下：
$ dpkg -e ../openssh-client_6.1p1_i386.deb extract/DEBIAN/
$ dpkg-deb -b extract/ build/ # build存放打包好的deb
$ ll build/
-rw-r--r-- 1 ufo ufo 1020014  7月  3 20:20 openssh-client_6.1p1_i386.deb # 验证方法为：再次解开重新打包的deb文件，查看在etc/ssh/sshd_config文件是否已经被修改
$ dpkg-scanpackages . | gzip -9c > Packages.gz #  制作本地软件源
```

## snap

## FAQ
### dpkg-deb: error: archive '<file>.deb' has premature member 'data.tar.gz' before
dpkg的bug: [dpkg无法解析tar.xz格式-xz compressed control.tar files not supported](https://bugs.launchpad.net/ubuntu/+source/dpkg/+bug/1730627)

升级dpkg版本(>=1.17.5), 即`apt install dpkg`.

原因:
对于软件安装包的提供者而言, 一定是希望安装包具有更好的兼容性. 最好可以使用xz压缩data部分, 仍然用gzip打control部分. 旧版的dpkg-deb, 默认会把control和data分开用不同的格式打包, control默认始终使用gzip的格式打包. 而新版的dpkg-deb(1.19.0)之后都会使用相同的格式压缩control和data. 如果你指定了-Z xz , 那就都是xz. 

还好, dpkg-deb提供了一个参数：--no-uniform-compression加上这一句就可以了. 

默认是：--uniform-compression, 代表使用统一的格式进行压缩. 加上--no-uniform-compression后不再统一, control使用gz压缩. 

### apt install报`Size mismatch`
下载到的deb软件包信息和源信息列表Packages记录(Packages.gz)的数据不相符, 可用`dpkg -i`安装

### apt install 安装的deb的缓存位置
ubuntu中由apt-get获得的文件包保存在/var/cache/apt/archives

### 删除snap
```bash
snap list; sudo snap remove xxx
sudo apt install ubuntu-software
sudo snap remove snap-store
sudo apt purge snapd
sudo rm -rf /var/cache/snapd
sudo rm -rf ~/snap
```

### 路径`debian/rules`
dpkg-buildpackage的构建目录结构

### 删除`dpkg -l`显示状态为的"rc"包
```bash
aptitude search ~c # list the residual packages
sudo aptitude purge ~c # purge them
```

### apt log
`/var/log/apt/term.log`