# bareos
参考:
- [备份/恢复系统BAREOS的安装、设置和使用（二）](https://blog.csdn.net/laotou1963/article/details/82711776)
- [OSBConf 2015 | Backup of VMware snapshots with Bareos by Philipp Storz & Stephan Dühr](https://www.youtube.com/watch?v=pDNhfK9MO0g)

Bareos 由 bacula fork而來.

[Bareos组成](https://docs.bareos.org/IntroductionAndTutorial/WhatIsBareos.html#bareos-components-or-services):
- bconsole : 全功能cli, 与Director进行通信
- webui : 只用于备份和恢复功能, 同时支持基于Web的bconsole界面
- Director : Bareos中控, 它计划并监督所有备份, 还原, 验证和存档操作.
- Storage Daemon : 在Bareos上作为备份数据的存储空间, 允许多个

    应Bareos Director的请求负责从Bareos File后台驻留程序接受数据, 并将文件属性和数据存储到物理备份介质或卷中
- File Daemon : **在客户机**, 管理本地文件的备份和恢复.

    它会应Director的请求它会找到要备份的文件, 并将指定数据发送到Bareos Storage Daemon
- Catalog : 目录服务由负责维护所有备份文件的文件索引和卷数据库.

    目录服务允许系统管理员或用户快速查找和还原任何所需的文件.

> Bareos推荐使用postgres, mysql/mariadb已废弃.

> 要成功执行保存或还原, 必须配置并运行以下四个守护程序：Director daemon, File daemon, Storage daemon 以及 Catalog service(即DB).

> [Bareos所有相关package](https://docs.bareos.org/IntroductionAndTutorial/WhatIsBareos.html#bareos-packages)

> [bareos网络连接概览](https://docs.bareos.org/TasksAndConcepts/NetworkSetup.html#network-connections-overview)

## 备份
ref:
- [备份族谱来了！细数各种备份类型（一）](https://www.talkwithtrend.com/Article/247499)
- [备份族谱来了！细数各种备份类型（二）](https://www.talkwithtrend.com/Article/247515)

![](/misc/img/develop/apps/1236661_158376518962878.png)

### 根据备份层级划分
根据层级可以分为:
- 文件级备份（File-Level Backup）

    将文件从一个地方复制到另外一个地方（备份介质）存储起来.

    对于全量备份后的增量备份，一般备份软件会通过比较文件的归档位或修改时间来检测自上一次备份以来文件是否发生过变化，以便于只备份更改过的文件.

    文件级备份比较简单，应用也比较广泛，但它存在一些缺点:
    - 对于比较大的文件，即使只是文件的一小部分变更了，那也必须备份整个文件，消耗备份介质
    - 某些操作系统（比如Windows）会阻止你打开正在使用的文件，这种场景下就需要辅助其他技术来支持备份
    - 需要对所有的文件进行打开/关闭操作，当小文件很多时备份的时间就会很长
    - 增量备份的时候必须遍历所有文件，以检查出那些变更的文件
- 映像级备份（Image-Level Backup）, 也叫块级备份（Block-Level Backup）

    块级备份通常伴随着快照技术，比如LVM的卷快照、windows下的vss技术、虚拟化平台针对虚拟机的快照等等. 有了快照技术，一般就可以针对快照后变化的数据块进行跟踪，增量备份时只备份变化的数据块。相比文件级备份，块级备份的效率往往更高，增量备份也更加有效.

    块级备份存在一些问题：

    - 通常需要备份整个计算机或整个磁盘，但可能你的磁盘并没有完全得到利用，空余了一大部分，空余的备份也会纳入备份范围
    - 文件系统层常常存在删除的文件以及一些碎片空间，这些空间对于块级备份是不可知的，也会纳入备份范围
    - 块级备份的风险更大. 每次增量备份都必须依赖于前面的备份数据不出现问题，否则会导致整个备份链都错误
    - 无法备份网络文件系统上的数据
- 字节级备份（Byte-Level Backup）

    字节级的备份依赖于监听文件系统层和应用层的IO，这种技术无法做全量备份，一般用于文件系统层的增量数据捕获，或者容灾技术。由于增量同步的数据传输量小，因此在云迁移云灾备等方案中表现还比较不错.
### 根据对业务的影响划分
备份过程需要读取生产数据，总会对生产系统带来一些影响. 基于对业务的影响，可以分为冷备份和热备份。一般来讲，冷备份和热备份描述的是针对于数据库的备份.

数据库处于关闭状态下的备份属于冷备份，有的时候也叫做离线备份或者脱机备份。 在备份过程中数据库不会产生新的数据。 使用了冷备份一方面备份操作比较简单、速度快，另一方面比较安全，维护简单。 冷备份的主要缺点在于： 备份过程数据库无法工作，数据库只能用于备份。

数据库处于运行状态下的备份属于热备份，也叫在线备份或者异步备份。 热备份要求数据库在Archivelog方式下工作，并且需要有比较大的档案空间。 热备份解决了运行状态下可以备份的问题，但是它的主要缺陷在于必须小心谨慎，确保备份成功，一旦失败，后果就比较严重了.

### 根据备份频率划分
根据备份频率一般可以分为定时备份和实时备份。定时备份就是指在某个确定的时间点进行备份. 一般灾备厂商会提供多种定时策略以满足用户多样的需求。实时备份就是指不间断的备份，能够实时的监控到用户数据的变化，并将变化的数据备份到灾备介质中.

#### 实时备份
实时备份顾名思义，表示实时的，不间断的备份。一般来讲，实时备份需要利用CDP（Continual Data Protection）技术，CDP一方面用于实时复制，另一方面用于实时备份。注意：实时复制一般是指与生产端与灾备端保持数据同步，当生产端出现故障时，灾备端可以直接拉起业务，保证业务的连续性；实时备份是指将数据保存到备份介质中，以便于数据可以恢复到任意历史时间点。

#### 定时备份
定时备份可以区分多种备份类型，比如全量备份，增量备份等等。然而随着数据量越来越大，备份场景越来越复杂，简单的全量备份和增量备份已经不能满足用户的需求，由此衍生出很多其他的备份类型，如：反向增量备份、合成备份等等.


[数据备份主要模式](https://www.baculasystems.com/glossary/enterprise-backup-types-differences-and-policies-full-incremental-differential-and-synthetic-full-backups/)主要包括全量备份、增量备份，以及差异备份等:
- 全量备份（Full backup）: 把数据完全复制一份

    全量备份可靠是可靠，但是消耗的介质空间是相当大的, 解决这个问题的办法最简单的就是限制副本数.

    ![](/misc/img/develop/apps/1236661_158384678850987.png)

- 增量备份（Incremental Backup）: 一种传统的数据备份技术，它以上一次完全备份或增量备份为基础，此后每次只备份相对于上一次备份操作以来变更(新创建,更新过或删除)的数据
    
    缺点:
    1. 每次增量备份都依赖于前一个备份副本，中间一个副本出了问题可能影响整个备份链，因此不如全量备份可靠
    1. 备份对应的就是恢复，恢复的过程是先恢复全量副本，再依次恢复增量副本。如果增量副本过多，那么会影响恢复效率

    要解决上面的问题，一般都会采用周期性全量备份的策略

    ![](/misc/img/develop/apps/1236661_158384681086375.png)

- 差异备份（Differential Backup）: 备份自从上次完全备份后被修改过的文件, 即指除第一次全量备份外，后面的备份都只备份上次全量到当前时间点的变化数据

    从差量备份中恢复也是很快的，因为只需要两份磁带——最后一次完全备份和最后一次差异备份.

    它与增量备份很类似，一般也是采用一样的技术，差异点就在于：**增量备份是基于上一次备份时间点（不管全量还是增量）获取变化数据，而差量备份是基于上一次全量备份时间点获取变化数据**. 差量备份的出现就是为了解决对增量备份链不信任的问题

    ![](/misc/img/develop/apps/1236661_158384683926482.png)
- 多级增量备份: 增量和差量备份的结合体, 它主要用于更加精细的备份策略

    在多级增量备份中，一般把全量备份定义为Level 0，把增量备份定义为Level 1, Level 2，Level 3等。Level n 仅备份自Level n-1以来的增量数据.

    多级增量备份实际应用的话不太多见，很多灾备厂商目前甚至都不提供这样的功能.

    ![](/misc/img/develop/apps/1236661_158384686159509.png)
- 反向增量备份: 将本次要备份的增量数据替换到全量副本中，然后将全量副本中对应的数据替换出来，作为历史增量数据

    ![](/misc/img/develop/apps/1236661_158384687957938.png)

    采用反向增量备份方式的客户，最终备份介质中数据副本的形态如下图。在这种方式下，客户可以以最快的速度恢复最新的时间点.

    ![](/misc/img/develop/apps/1236661_158384689527809.png)
- 永久增量备份，就是指除第一次备份外，后续所有备份都采用增量备份的方式

    因为随着数据量的爆发式增长, 已经无法完成全备了.

    永久增量备份其实牺牲了客户对全量备份的信任和依赖，是客户数据量增长过快场景下的妥协.
- 合成完全备份(Synthetic Full Backups): 将备份介质中的全量和增量数据进行整合，合并成一个新的全量副本

    执行一次完全备份，后续每次均进行增量备份，备份软件将会基于第一份全量备份数据，和随后的增量备份集进行合并，生成新的全量备份数据，然后，再定时把新生成的全量备份集和增量数据集进行合并，再生成新的全量备份数据，以此循环处理.

    合成完全备份是在恢复过程中充当完全备份的备份，但在备份过程中却不充当完全备份.

    一般灾备厂商也会提供合成备份的策略供客户自行配置，比如：多长时间进行一次合成、产生了多少个副本后进行一次合成等等。常规的合成备份由于会生成一份全新的全量副本，经过合成备份后，合成之前的数据副本可以删除。这样的方式也称作前向永久增量。这样的方式无法解决反向增量备份原本要解决的问题，也就是恢复效率的问题.

    ![](/misc/img/develop/apps/1236661_158384693591159.png)

## 编译
参考:
- [Configure (cmake) build settings](https://docs.bareos.org/DeveloperGuide/BuildAndTestBareos/systemtests.html)
### v21.1.2
env: oracle linux 7.9

前提: 根据FAQ安装gcc9

```bash
# --- 准备环境
yum install cmake3 # 在repo `ol7_developer_EPEL`
ls -s /usr/bin/cmake3 /usr/bin/cmake
yum install redhat-lsb redhat-release jannson-devel # jannson-devel在repo `ol7_optional_latest`需启用
yum-builddep bareos.spec --define "centos_version 790"
yum-builddep python-bareos.spec --define "centos_version 790"
# --- 开始构建
cd /root/rpmbuild/SPECS
# spec文件来源于[官方](http://download.bareos.org/bareos/release/21/RHEL_7/src/)
# 修正:
# - `droplet 0`, droplet已不在维护
# - `# BuildRequires: lsb-release`, 没有lsb-release包. 查看Ubuntu 20.04的该包, 其提供了命令lsb_release, 但它已在oracle linux的redhat-lsb-core包里. 
rpmbuild -bb bareos.spec --define "centos_version 790"
rpmbuild -bb python-bareos.spec --define "centos_version 790"
cp -r ../RPMS/noarch/* ../RPMS/x86_64
ll ../RPMS/x86_64 # 所需rpms
```

### v20.0.3
env: Ubuntu 20.04

```bash
# git clone -b Release/20.0.3 --depth=1 git@github.com:bareos/bareos.git # 应使用`git clone xxx`方式获取bareos源码, 因为cmake需要从git tag/log获取信息.
# cmake -P write_version_files.cmake
# apt install libreadline-dev libpq-dev chrpath
# mkdir build && cd build
# cmake -Dpostgresql=yes -Dtraymonitor=no -Dmysql=no -Dsqlite3=no .. # make install时用, 而非deb打包时, cmake参数参考`debian/rules`
dpkg-checkbuilddeps
# generate changelog from [here](https://github.com/bareos/bareos/blob/15f82cd288f295f4ae13c3f27775eb2df46f2c98/.travis.yml)
NOW=$(LANG=C date -R -u)
BAREOS_VERSION=$(cmake -P get_version.cmake | sed -e 's/-- //')
printf "bareos (%s) unstable; urgency=low\n\n  * See https://docs.bareos.org/release-notes/\n\n -- nobody <nobody@example.com>  %s\n\n" "${BAREOS_VERSION}" "${NOW}" | tee debian/changelog # 或从官方deb中拷贝一份
vim debian/rules # ~~根据上面的cmake参数定制deb打包编译bareos时需要的参数~~. 不能修改参数, 只能装全依赖, 因为deb打包时dh_install并没有根据参数(比如`-Dmysql=no -Dsqlite3=no`, `-Dtraymonitor=no`等)忽略相关依赖文件. 不禁用traymonitor会报很多错. 删除override_dh_install并未使用的libbareoscats-mysql, libbareoscats-sqlite3
fakeroot debian/rules binary # 来自`/.travis/travis_before_script.sh`. 编译完成后在打包时会报错， 因为生成的debian/control还是包含了mysql， sqlit3和traymonitor, 同时还要去掉debian/bareos-[director|filedaemon|storage].install中的traymonitor相关项, 修改后再次执行该命令即可.
# --- 制作apt local repo
# mkdir bareos-apt
# mv bareos*.deb bareos-apt
# cd bareos-apt
# apt-ftparchive packages . > Packages
# vim /etc/apt/sources.list
deb [trusted=yes] file:///root/bareos-apt/ ./ # 放在第一行， 优先使用. "[trusted=yes]"在apt 1.1开始支持
...
# apt-get update [--allow-insecure-repositories] # allow-insecure-repositories会忽略报错: `The repository 'file:/xxx xxx/ Release' does not have a Release file`
```

> 使用`dpkg-scanpackages bareos-apt | gzip> bareos-apt/Packages.gz`, `apt update`时会报"File not found - /root/bareos-apt/Packages"

总结: bareos debian/rules 是编译全部组件的, 禁用部分选项则需要修改相关的deb构建脚本. 

仅cmake编译(非fakeroot打包编译)的缺陷:
1. arm没有vmware插件, 因为依赖的vmware不提供arm so
1. xxx.service 没有User/Group, 可使用root
1. 数据库配置在`/usr/local/etc/bareos-dir.d/catalog/Mycatalog.conf`, 且默认使用sqlite3, 需改用postgres
1. `bareos-dir -t -f -d 500 -v`发现database bareos不存在. 需手动配置db见[这里](https://docs.bareos.org/IntroductionAndTutorial/InstallingBareos.html#other-platforms)

```bash
# --- kylin v10
# dnf install kylin-lsb kylin-release jansson-devel # 其他依赖rpmbuild时根据提示安装. jansson-devel kylin默认没有,导致没法启用".api json"
# rpmdev-setuptree
# cd ~/rpmbuild
# cp bareos-Release-20.0.3.tar.gz SOURCES/bareos-20.0.3.tar.gz
# cp bareos-Release-20.0.3/core/platforms/packaging/bareos.spec SPECS/bareos.spec # 或使用官方[src.rpm里的spec](https://download.bareos.org/bareos/release/20/CentOS_7/src/).
# --- set bareos compile platform, 见[core/platforms](https://github.com/bareos/bareos/tree/master/core/platforms), 这里应该参照centos把platform指定为redhat
# vim SOURCES/bareos-20.0.3.tar.gz # 先解压bareos-20.0.3.tar.gz再编辑再重新打包, 会在执行`rpmbuild -bb bareos.spec`时因为解压处理软连接问题而报错
# 修改:
#     - core/cmake/distname.sh: CentOS) -> CentOS|Kylin)
#     - core/cmake/BareosGetDistInfo.cmake: COMMAND ${CMAKE_CURRENT_LIST_DIR}/distname.sh -> COMMAND bash ${CMAKE_CURRENT_LIST_DIR}/distname.sh # 因为vim编辑distname.sh后丢失可执行权限; **通过压缩软件打开bareos-20.0.3.tar.gz,再用编辑器编辑,保存时可自动借助压缩软件重新打包功能使得不丢失权限**
# vim SPECS/bareos.spec
# 修正: Version: 20.0.3; Release: 3%{?dist}; user/group: bareos-> root
#      build_qt_monitor 0, build_sqlite3 0; build_mysql 0; systemd_support 1; droplet 0 (kylin没有droplet, 但bareos repo包含了libdroplet的源码)
#      redhat-lsb->kylin-lsb; lsb-release->kylin-release
#      Requires: php-zip -> Requires: php-common; Requires: libzip # kylin repo没有php-zip， 参考centos7/8, 它由php-common+libzip替代
#      按照官方bareos.spec 修正changelog
# cd SPECS && rpmbuild -bb bareos.spec # [rpmbuild报error: `Installed (but unpackaged) file(s)`, 因为glusterfs没启用但还是编译出了相关文件](/shell/tools/packages.md)
# --- 编译python-bareos
# cd ../BUILD/bareos
# mv python-bareos python-bareos-20.0.3
# cp python-bareos-20.0.3/packaging/python-bareos.spec ../../SPECS
# tar -czf python-bareos-20.0.3.tar.gz python-bareos-20.0.3
# cp python-bareos-20.0.3.tar.gz ../../SOURCES
# cd ../..
# vim SPECS/python-bareos.spec
# 修正：Version:        20.0.3; Release:        3%{?dist}
#      按照官方bareos.spec 修正changelog
# cd SPECS && rpmbuild -bb python-bareos.spec
# --- 制作本地dnf repo
# cat /etc/yum.repos.d/local.repo 
[local-media]
name=Kylinv10 - bareos
baseurl=file:///root/rpmbuild/RPMS/aarch64/
gpgcheck=0
enabled=1
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-kylin
# createrepo /root/rpmbuild/RPMS/aarch64
# yum clean all
# yum makecache # 会报`cannot download repodata/repomd.xml`， 因此不能省略createrepo的步骤.
```

oracle linux 7.9 x86启用vmware插件:
1. 根据[官方文档](https://docs.bareos.org/bareos-21/TasksAndConcepts/Plugins.html#vmware-plugin)提供的[vmware vddk](https://code.vmware.com/web/sdk/7.0/vddk)地址, 注册并获取vddk
1. 配置vddk, 配置方法源自`/usr/include/jansson.h`和`core/cmake/BareosFindAllLibraries.cmake`

    ```bash
    tar -xf VMware-vix-disklib-7.0.3-19513565.x86_64.tar.gz
    cd vmware-vix-disklib-distrib
    cp include/* /usr/include
    mkdir -p /usr/lib/VMware-vix-disklib
    cp -r lib64/* /usr/lib/VMware-vix-disklib
    ```

    再将VMware-vix-disklib-7.0.3-19513565.x86_64.tar.gz删减到仅保留其中的lib64目录, 并重命名为VMware-vix-disklib-7.0.3-19513565.only_redistributable_libs.x86_64.tar.gz, 最后放入`rpmbuild/SOURCES`
1. 编译bareos并开启vmware配置

    `rpmbuild -bb bareos.spec --define "rhel_version 790" --define "vmware 1"`
1. 从[官方bareos-vmware-vix-disklib-7.0.1_16860560-1.el7_9.src.rpm](https://download.bareos.org/bareos/release/21/RHEL_7/src/bareos-vmware-vix-disklib-7.0.1_16860560-1.el7_9.src.rpm)提取bareos-vmware-vix-disklib.spec, 并修改其VMware-vix-disklib到相应版本, 再编译即可

    其实就是将vmware官方提供的so打包成rpm

## 概念
- volume : Bareos将在其上写入备份数据的单个物理磁带（或可能是单个文件）
- pool : 定义接收备份数据的多个volume（磁带或文件）组成的逻辑组


## [部署](https://docs.bareos.org/IntroductionAndTutorial/InstallingBareos.html#install-the-bareos-software-packages)
```bash
# -- pg
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
# -- bareos
wget -q http://download.bareos.org/bareos/release/20/xUbuntu_20.04/Release.key -O- | apt-key add -
wget -O /etc/apt/sources.list.d/bareos.list http://download.bareos.org/bareos/release/20/xUbuntu_20.04/bareos.list
# -- install
sudo apt install postgresql-12 postgresql-client-12 pgadmin4

vim ${pg}/pg_hba.conf # kylinv10在/var/lib/pgsql/data/pg_hba.conf， 且kylin环境需要将原有两条`host all ... ident`的ident替换为md5, 否则下面的psql测试密码登入将无法登入
local bareos bareos md5 # 插在最前面. bareos默认使用本地pg, 因此添加该匹配规则

sudo -u postgres psql # 进入psql
> alter user postgres with password 'postgres'; # 为postgres创建密码
psql -h localhost -p 5432 -U postgres -W # 测试密码登入

systemctl restart postgresql

# -- [官方安装文档](https://docs.bareos.org/IntroductionAndTutorial/InstallingBareos.html#section-installbareospackages)
apt install bareos bareos-database-postgresql # 输入db密码. bareos-database-postgresql会利用dbconfig-common mechanism, 在apt install过程中配置db, 相关配置会保存在`/etc/dbconfig-common/bareos-database-common.conf`. 可用`dpkg-reconfigure bareos-database-common`重新配置, 手动配置db见[这里](https://docs.bareos.org/IntroductionAndTutorial/InstallingBareos.html#other-platforms), 它是根据bareos-database-common.conf来进行db初始化.
# --- rpm: 因为redhat环境在bareos-database-postgresql时没有交互配置db的过程, 因此[需要参照官方文档执行scripts来初始化db](https://docs.bareos.org/IntroductionAndTutorial/InstallingBareos.html#other-platforms), db初始化完成后需要利用pg postgres账号修改pg bareos账号的密码, 再更新bareos db conf: /etc/bareos/bareos-dir.d/catalog/MyCatalog.conf.

systemctl restart bareos-dir
systemctl restart bareos-sd
systemctl restart bareos-fd

bareos-dir -t -f -d 500 -v # 测试bareos-dir是否正常, 包括与pg的连接
bareos-sd -t -f -d 500 -v
bareos-fd -t -f -d 500 -v
bareos-dbcheck -B # 作用同上, 显示db的连接信息

apt install bareos-webui # 默认基于php+apache2, 推荐使用nginx
# -- 配置webui, 也可使用[bconsole configure子命令](https://docs.bareos.org/IntroductionAndTutorial/InstallingBareosWebui.html#create-a-restricted-consoles)
cp /etc/bareos/bareos-dir.d/console/admin.conf.example /etc/bareos/bareos-dir.d/console/admin.conf
vim /etc/bareos/bareos-dir.d/console/admin.conf # 设置bareos-dir admin用于bareos-webui. 如果bconsole reload失败则需要: chown bareos:bareos /etc/bareos/bareos-dir.d/console/admin.conf
vim /etc/bareos/bareos-dir.d/profile/webui-admin.conf # 按需求修改CommandACL, 比如删除`!purge, !prune, !configure`.
systemctl restart bareos-dir # 不能省略, 否则可能webui无法登入(账号正确)
systemctl restart apache2 # 访问http://HOSTNAME/bareos-webui即可使用webui

systemctl enable bareos-dir bareos-sd bareos-fd postgresql
```

bareos-webui也可使用[nginx](https://docs.bareos.org/IntroductionAndTutorial/InstallingBareosWebui.html#nginx), 但访问地址要变为`http://bareos:9100/`
```bash
# apt/yum install nginx php-fpm # nginx 1.18.0
# systemctl enable php-fpm
# systemctl start php-fpm
# mkdir /etc/nginix/snippets
# cat /etc/nginix/snippets/fastcgi-php.conf
fastcgi_split_path_info ^(.+?\.php)(/.*)$;
try_files $fastcgi_script_name = 404;
set $path_info $fastcgi_path_info;
fastcgi_param PATH_INFO $path_info;
fastcgi_index index.php;
include fastcgi_params;
# vim /etc/nginx/bareos-webui.conf # 具体配置参考官网doc
# -- 修正for kylinv10 with php7.2：
#    1. fastcgi_pass unix:/var/run/php5-fpm.sock; -> fastcgi_pass unix:/var/run/php5-fpm.sock;
# -- 修正for ubuntu 20.04 with php7.4：
#    1. fastcgi_pass unix:/var/run/php5-fpm.sock; -> fastcgi_pass unix:/var/run/php/php-fpm.sock;
# systemctl restart nginx
```

### bareos-fd部署
1. 需备份的机器(client端, 使用9102端口, 等待来自bareos-dir的连接)安装客户端软件bareos-filedaemon

    - `apt install bareos-filedaemon`

    > [Windows client下载地址](http://download.bareos.org/bareos/release/20/windows/), `netstat -ano|findstr "9102"`
1. bareos director配置bareos-dir
```bash
$ bconsole
* configure add client name=client2-fd address=192.168.0.2 password=secret [TlsRequire=yes] # 注册client, 会创建`/etc/bareos/bareos-dir.d/client/client2-fd.conf`和`/etc/bareos/bareos-dir-export/client/client2-fd/bareos-fd.d/director/bareos-dir.conf`(bareos-fd访问bareos-dir的授权, **如果其中不包含Address-<dir_ip>时请添加**)
reload # 不能丢
exit
```
1. 配置clients

    需备份的机器(client端, 使用9102端口, 等待来自bareos-dir的连接)安装客户端软件

    - linux

        ```bash
        # apt install bareos-filedaemon
        # scp dareos-server:/etc/bareos/bareos-dir-export/client/client2-fd/bareos-fd.d/director/bareos-dir.conf /etc/bareos/bareos-fd.d/director # director对client的授权
        # vim /etc/bareos/bareos-fd.d/director/director-mon.conf # 用/etc/bareos/bareos-dir.d/console/bareos-mon.conf文件中的Password替换该文件中的Password
        # systemctl restart bareos-fd
        ```

    - windows

        需要设置的参数:
        - Client Name: bconsole注册clients时的名称, 最好是clients os的hostname
        - Director Name: 不修改
        - Password: dareos-server:/etc/bareos/bareos-dir-export/client/client2-fd/bareos-fd.d/director/bareos-dir.conf中的Password即`[md5]xxx`
        - Network Address: 注册client时本机的ip
        - Client Monitor Password: 用/etc/bareos/bareos-dir.d/console/bareos-mon.conf文件中的Password

1. 测试client by bconsole

  ```bash
  status client=client2-fd
  ```

  如果无法链接到windows的bareos client上(windows日志里均是提示tls相关错误), 先卸载该client, 卸载时必须选择**不保留配置**, 再重新安装并配置入正确的参数即可.

## bconsole cmd
```bash
* reload # 重载配置
* status client # 测试client connection
* status storage
* show client=l130 [verbose]
* show fileset[=xxx]
* list clients
* list pools
* list volumes [pool=xxx]
* list jobs
* llist jobs # 更详细的`list jobs`
* .jobs # 更精简的`list jobs`, 只有job name
* configure add console name=admin password=pwd111111 profile=webui-admin # 注册bconsole
* configure add client name=client2-fd address=192.168.0.2 password=secret # 注册client, 需要重启bareos-dir
* setdebug client=bareos-fd level=200 # [测试client](https://docs.bareos.org/TasksAndConcepts/TheWindowsVersionOfBareos.html#enable-debuggging)
* configure add job name=client2-job client=client2-fd jobdefs=DefaultJob # 添加job
* restore # 常用选项3来还原指定id. 选择文件的命令在[restore-command](https://docs.bareos.org/TasksAndConcepts/TheRestoreCommand.html#restore-command)即`mark (xxx|*)`, 被选中的文件名前会带`*`

# --- cancel
* cancle all # 取消所有job
* cancel 20 #  取消任务ID=20 的任务

# --- list命令列出各种备份状态信息
list Jobs     #列出所有备份记录状态
list jobid=2  #列出jobid等于2有状态信息
 
list Job=t3_full       #列出Job名称等于t3_full的任务信息
list jobname=t3_full   #列出Job名称等于t3_full的任务信息
list joblog jobid=78   #列出jobid=78的详细备份日志信息
list jobmedia/volumes jobid=78 #列出jobid=78的状态信息与所在Volume信息
list files jobid=78    #列出jobid=78的状态信息与所备份的数据信息
list jobs jobname=xxx client=xxx jobstatus=x joblevel=x last
 
list clients           #列出备份的客户端
list jobtotals         #列出所有作业任务使用的空间大小
 
list media pool=dbpool   #查看dbpool属性的media
list Volume Pool=dbpool  #查看dbpool属性的Volume
 
list pool    #查看定义的dbpool属性
llist pool   #查看定义的dbpool属性(更详细)

llist backups client="xxx" filset="any" order=desc limit=200 # 显示该客户端的所有(不限制fileset)备份任务的前200条. v20.2 order参数不生效
llist jobs job="xxx" jobstatus=x limit=200 offset=100 [days=14]# **llist jobs不支持order**
llist jobid=2160 # 输出jobid=2160的信息

> llist = long list, 即使用与list相同的参数, 但会列出所选记录的完整内容(from db)

# --- show查看配置信息
show Job=t3_full   #查看Job名称等于t3_full的配置信息
show pools         #查看池的信息
show pools=dbpool  #查看dbpool池的信息
show filesets
show clients
show storages [verbose]
show schedule
show jobs
show message

# --- status当着状态信息
status #查看状态信息
status client=t3-fd  #客户端名称t3-fd的状态信息
status client   # 查看 client  的状态
status dir      # 查看director 的状态
status storage  # 查看 storage 的状态 

# --- run执行job任务. bareos storage空间满后会阻塞分配到其上的job
run  # 未指定job时需要选择job, 即进入交互模式操作
run job=t3_full yes   #手动执行job为t3_full任务作业
rerun jobid=xxx yes
cancel jobid=xxx yes
enable jobid=xxx yes
disable jobid=xxx yes

# ---- restore
restore # 开始进入交互式restore流程
1 # 选择最近20个job
3 # 允许使用指定的jobid
cwd is: / # 进入bvfs选择文件,
mark <file> # 选中文件, 被选中文件的前面会出现`*`. `unmark <file>`取消选中
done # 完成选中

# --- estimate : 对某次任务进行评估. 它会连接到客户端，并输出这次任务的fileset 中 文件数,和这次备份任务所占的空间
estimate job=t3_full listing client=t3-fd  #估算下这个备份有多少文件,需要多大容量. 作业任务t3_full,客户端t3-fd
# --- delete删除备份
delete JobId=79  #删除jobid等于79的备份
list JobId=79    #查看就没有这个备份包了,但在status中还是会出这个,实际存储中空间并没有减小.

# --- 特殊的几个命令
.jobs [type=R]     #查看定义的job作业任务名称. `type=R`是restore job.
.clients  #查看定义的客户端名称
.filesets #查看定义的备份资源FS的名称
.msgs     #查看定义的日志消息记录的名称
.pools    #查看定义的pool池属性名称
.storage  #查看定义的storage数据的存储方式的名称

# --- 清理
purge # 是一个危险命令, 能清除一个客户端的所有备份任务，文件，和卷
purge volume storage=File pool=Full + "*<mediaid>" # 清理指定volume
prune # 这个命令和 purge 相似，但安全很多，它只会清除过期的文件，任务，和卷 
```

```bash
printf "list clients\nquit" | bconsole
printf ".api json\nlist clients" | bconsole
```

## webui
### 还原
- 客户端：从下拉菜单中选择备份所属的客户端
- 备份作业：从下拉菜单中选择需要的备份作业
- 合并所有客户端文件集：自动把该客户端该作业和该作业以前的所有备份（**含不同作业**）集合在一起供恢复文件使用; 如选"否", 只从选择的备份中恢复文件
- 合并所有相关作业：如选"是", 自动把该客户端该作业和该作业以前的所有**同一作业**的备份集合在一起供恢复文件使用; 如选"否", 只从选择的备份中恢复文件
- 还原到客户端：从下拉菜单中选择恢复文件的目标客户端
- 还原作业：从下拉菜单中选择预定义的还原作业
- 替换客户端上的文件：选择同名文件的覆盖规则. 可选规则为：总是、从不、比现有文件旧和比现有文件新
- 要恢复到客户端的位置：指定恢复文件的目标路径
- 文件选择：点击文件/路径前的`□`来选择是否要恢复此文件/路径; 如选择路径, 在该路径下的所有文件都会被恢复

## 流程图/时序

[连接流程](https://docs.bareos.org/TasksAndConcepts/NetworkSetup.html#network-connections-overview)
如果dir-fd之间的连接有正在执行的job, 那么fd不会复用该连接而是在需要时发起新连接, 用于发送其他命令, 比如`cancel`.

[job的网络时序图](https://docs.bareos.org/DeveloperGuide/netprotocol.html#network-sequence-diagrams)

[使用sd device的流程图](https://docs.bareos.org/DeveloperGuide/reservation.html#usedevicecmd)

[job执行的流程图](https://docs.bareos.org/DeveloperGuide/jobexec.html)

## api
bareos console支持非交互式的[点命令](https://docs.bareos.org/DeveloperGuide/api.html#dot-commands), 同时支持json输出(执行`.api json compact=yes`即可, compact=yes表示压缩空格).

### python-bareos
[python-bareos](https://github.com/bareos/bareos/tree/master/python-bareos/)是bareos官方的python sdk, 用于与bareos-dir通信.

### Bareos REST API (based on python-bareos)
参考:
- [README](https://github.com/bareos/bareos/tree/master/rest-api#readme)

```bash
wget https://github.com/bareos/bareos/archive/refs/tags/Release/20.0.1.tar.gz
tar -xf 20.0.1.tar.gz && cd bareos/rest-api
pip3 install -r requirements.txt
vim api.ini # 配置Director并设置secret_key
uvicorn [--host 0.0.0.0 --port 8000] bareos-restapi:app --reload
```

Serve the Swagger UI to explore the REST API: http://127.0.0.1:8000/docs
Alternatively you can use the redoc format: http://127.0.0.1:8000/redoc

原理: 用账户名和密码创建[bareos.bsock.DirectorConsoleJson](https://pypi.org/project/python-bareos/), 再将DirectorConsoleJson和用户名关联, 返回包含该用户名的JWT, 调用restful api with JWT即使用`DirectorConsoleJson.call(cmd)`执行拼接好的cmd.

**需要更新账号对应的profile**, 否则部分rest api接口会报错, 比如创建client, 具体见FAQ的`configure: is an invalid command.`.

> 将`rest-api/bareos-restapi.py`中的print打印的注释去掉即可看到rest-api执行过程中向bareos-dir.d发送的cmd了.

> 页面有cdn资源依赖. 该功能由fastapi提供, [离线资源加载看这里](https://fastapi.tiangolo.com/advanced/extending-openapi/#self-hosting-javascript-and-css-for-docs), 在自身项目上引入fastapi资源来解决. 注意不能忘记这两属性`FastAPI(docs_url=None, redoc_url=None)`, 否则应用还是使用fastapi默认的渲染函数.

> 只需设置`http://127.0.0.1:8000/docs`页面的"Authorize"按钮里的username和password即可使用openapi的`try it out`

#### url map
根据`/<module>/<action>` -> ``bareos-webui/module/<module>/src/<module>/Controller/<module>Controller.php#<action>Action`映射的, 比如
- `/restore/` : bareos-webui/module/Restore/src/Restore/Controller/RestoreController.php#indexAction
- `/restore/filebrowser` : bareos-webui/module/Restore/src/Restore/Controller/RestoreController.php#filebrowserAction

### 要点
建议在`.bvfs_lsfiles`和`bvfs_lsdirs`查询中使用 pathid 而不是 path, 查询`/`除外.

文件还原时checked_files和checked_directories传参:
- 半选中状态目录不传(即目录内容没有被全选)
- 选中状态全传(可以剪枝:目录内容全选时可只传该目录id)

## plugin
> [官方 plugins](https://github.com/bareos/bareos/tree/master/core/src/plugins/filed), [官方 contrib plugins](https://github.com/bareos/bareos/tree/master/contrib)和[开源plugins:"bareos-tasks-plugins"(其他它已包含在contrib plugins中)](https://github.com/marcolertora/bareos-tasks-plugins)

bareos原生支持dir, storage, filedaemon的插件扩展. 使用插件前必须在配置中启用它们, **修改后需要重启服务**, 当前支持python 2/3. **bareos 20开始推荐使用python3, 虽然官方20.0.1目前plugins都是python2的**.

> **一个client无法同时加载python2和python3插件, 可创建两个client分别加载python2/3来解决**.

> 前提: `apt install bareos-{director,storage,filedaemon}-python3-plugin`或`apt install bareos-{director,storage,filedaemon}-python2-plugin`, 都装时先安装python3的.

> [Porting existing Python plugins和Switching to Python 3](https://docs.bareos.org/TasksAndConcepts/Plugins.html)

> [bpluginfo](https://docs.bareos.org/Appendix/BareosPrograms.html#bpluginfo)可用于查看plugin相关信息, 比如`bpluginfo -v /usr/lib/bareos/plugins/python3-fd.so`

> 插件依赖的python package在`core/src/plugins/{dir,file,store}d/python/pyfiles`下, 会由`bareos-{directoor,storage,filedaemon}-python-plugins-common`安装在`/usr/lib/bareos/plugins`下

因为最常用的是fd-plugins, 这里重点介绍. 其他两种请参考[bareos docs](https://docs.bareos.org/TasksAndConcepts/Plugins.html)

区分python2/3 plugin的方法: 查看bareos-fd-xxx.py的load_bareos_plugin的返回值, bRC_OK是python3的, `bRCs['bRC_OK']`是python2的.

python3 plugin启用方法:
1. 在client conf使用

    ```conf
    # vim /etc/bareos/bareos-fd.d/client/xxx.conf
    FileDaemon {                          
        Name = client-fd
        ...
        Plugin Names = "python3" # 其实就是用于指定bareos plugins目录下的`xxx-fd.so`
        Plugin Directory = /usr/lib/bareos/plugins
    }
    ```
1. 在fileset上使用

    ```conf
    # vim /etc/bareos/bareos-dir.d/fileset/xxx.conf
    Plugin = "python:module_path …"
    ```

    **不能省略module_path**, 否则会报"No module named 'xxx'"(即使bareos-fd已配置Plugin Directory)

### fd-plugins
以[官方MySQL Plugin](https://github.com/bareos/bareos/tree/master/contrib/fd-plugins/mysql-python)举例:
1. 配置

    - client安装bareos-filedaemon-python-plugin
    - client的`bareos-fd.d/client/myself.conf`的`Plugin Directory`指向fd plugins目录`/usr/lib/bareos/plugins`
    - director中的`bareos-dir.d/fileset/mysql.conf`: `Include.Plugin = "python:module_path=/usr/lib/bareos/plugins:module_name=bareos-fd-mysql"`

        插件参数拼接在Plugin中以`:`分隔即可, 比如`Plugin = "python:module_path=/usr/lib/bareos/plugins:module_name=bareos-fd-mysql:mysqlhost=dbhost:mysqluser=bareos:mysqlpassword=bareos"`

        **module_path即client的fd plugins目录**

    > bareos-fd-mysql插件中的[_mysqlbackups_](https://docs.bareos.org/Appendix/Howtos.html#backup-of-mysql-databases-using-the-python-mysql-plugin)是虚拟目录, 说明fd plugins可将io流(mysqldump的输出)发送到storage中.

    还原时, 显示的是`_mysqlbackups_`的文件

其他官方插件:
- [`bareos-fd-local-fileset`](https://github.com/aussendorf/bareos-fd-python-plugins/wiki): 备份时动态将filename=/etc/bareos/extra-files中的文件列表加入fileset

fd-plugins其实就是操作fileset, fliter或添加需要备份的文件列表.

## FAQ
### bconsole配置
`/etc/bareos/bconsole.conf`

### bareos-sd配置
> 修改bareos-sd的配置后, 必须重启bareos-sd. 在重启bareos-sd前, 请首先使用`bareos-sd -t -v`检查bareos-sd配置文件, 如它没有任何输出, 说明配置文件没有任何语法问题.

`/etc/bareos/bareos-sd.d`:
- device : [数据存储位置](https://docs.bareos.org/Configuration/StorageDaemon.html#device-resource)

    ```conf
    # HDD 存储设备
    Device {
      Name = FileStorage                  # 设备名称
      Media Type = File                   # 媒体类型, [必须唯一, 否则还原时可能找不到备份所使用的pool](https://bugs.bareos.org/view.php?id=1455)
      Archive Device = /bareos/hdd        # Ubuntu下的备份文件目录（或mount point）
      LabelMedia = yes;                   # lets Bareos label unlabeled media
      Random Access = yes;                # 可随机读写
      AutomaticMount = yes;               # 自动加载
      RemovableMedia = no;                # 媒体介质不可移除
      AlwaysOpen = yes;                   # 建议总是打开, FIFO存储设备除外
      Description = "File device. A connecting Director must have the same Name and MediaType"
    }

    # 磁带存储设备
    Device {
      Name = TapeStorage                  # 设备名称
      Media Type = File2
      Archive Device = /bareos/tape       # Ubuntu下的mount point
      LabelMedia = yes;                   # lets Bareos label unlabeled media
      Random Access = no;                 # 不能随机读写
      AutomaticMount = no;                # 不自动加载
      RemovableMedia = yes;               # 媒体介质可移除
      AlwaysOpen = yes;                   # 按需打开
    }
    ```
- director

    - bareos-dir.conf : 管理storage对director的授权

        - Password : 授权director访问sd的密码. 在director创建storage时会用到. 同理bareos-fd也有该文件和该字段, 作用相同.
    - bareos-mon.conf : 管理storage对bareos traymonitor的授权
- message : storage message管理
    
    - Standard.conf : bareos-sd日志处理
- storage

    - bareos-sd.conf : bareos-sd配置

### bareos-dir配置
> 修改bareos-dir的配置后(比如添加fileset), 必须重启Director. 在重启Director前, 请首先使用`bareos-dir -t -v`检查bareos-dir配置文件. 如命令没有任何输出, 说明配置文件没有任何语法问题.

> 创建文件时注意owner需要是bareos, 否则`systemctl restart bareos-dir`会因为权限导致执行失败.

`/etc/bareos/bareos-dir.d`:
- catalog : 备份/还原索引信息来源

    - MyCatalog.conf : db配置
- client : clients信息

    - xxx.conf : client注册信息
- console

    - admin.conf : web ui访问的授权
    - bareos-mon.conf : monitor访问bareos-dir的授权
- director

    - bareos-dir.conf : bareos-dir配置
- fileset : 备份文件组(定义如何备份一组文件)配置

    example:
    ```conf
    FileSet {                                     # fileset 开始标志
      Name = "LinuxAll"                           # 该 fileset 的名字，这个名字会在备份任务中使用
      Description = "备份所有系统，除了不需要备份的。"
      Include {                                   # 备份中需要包含的文件
        Options {                                 # 选项
          Compression = LZ4                       # [压缩](https://docs.bareos.org/Configuration/Director.html)
          Signature = MD5                         # 每个文件产生MD5校验文件
          One FS = No                             # 所有指定的文件（含子目录是mountpoint）都会被备份
          # One FS = Yes                          # 指定的文件（含子目录）如不在同一文件系统下不会被备份
          #
          # 需要备份的文件系统类型列表
          FS Type = btrfs                         # btrfs 文件系统需要备份
          FS Type = ext2                          # ext2 文件系统需要备份
          FS Type = ext3                          # ext3 文件系统需要备份
          FS Type = ext4                          # ext4 文件系统需要备份
          FS Type = reiserfs                      # reiserfs 文件系统需要备份
          FS Type = jfs                           # jfs 文件系统需要备份
          FS Type = xfs                           # xfs 文件系统需要备份
          FS Type = zfs                           # zfs 文件系统需要备份
        }
        File = /                                  # 所有目录和文件
      }
      # 定义不需要备份的文件和目录
      Exclude {                                   # 备份中不应该包含的文件
        # 无需备份文件/目录列表
        File = /var/lib/bareos                    # /var/lib/bareos 下放的是bareos的临时文件
        File = /var/lib/bareos/storage            # /var/lib/bareos/storage 下放的是备份文件
        File = /proc                              # /proc 无需备份
        File = /tmp                               # /tmp无需备份
        File = /var/tmp                           # /var/tmp无需备份
        File = /.journal                          # /.journal 无需备份
        File = /.fsck                             # /.fsck无需备份
      }
    }

    FileSet {
      Name = "Windows电脑备份[A-Z]:/QMDownload"
      Enable VSS = yes                                  # 当YES时，当文件正在被写时也能被备份；如NO，被写文件不会被备份
      Include {
        Options {
          Signature = MD5
          Drive Type = fixed                            # 只备份固定磁盘, only for windows
          IgnoreCase = yes                              # 忽略字母的大小写, only for windows
          WildFile = "[A-Z]:/pagefile.sys"              # 指定文件：从磁盘A到Z下的/pagefile.sys. Wild是通配符的简写
          WildDir = "[A-Z]:/RECYCLER"                   # 指定文件：从磁盘A到Z下的
          WildDir = "[A-Z]:/$RECYCLE.BIN"               # 指定文件：从磁盘A到Z下的
          WildDir = "[A-Z]:/System Volume Information"  # 指定文件：从磁盘A到Z下的
          Exclude = yes                                 # 另一种方式指定不备份上述指定文件
        }
        File ="C:/QMDownload"                    # 备份目录C:/QMDownload
      }
    ```

    - win.conf

        ```conf
        # all office files in users (c:/ and d:/)
        # for win 7     = D
        # for win 10    = C 


        FileSet {
          Name = "Win7_office"
          
          # volume shadow copy service
          Enable VSS = yes
          Include {
          
          # location
            File = "D:/Users"
            File = "D:/My Documents"
          
          Options {
            # config
            Signature = MD5
            compression = LZ4
            IgnoreCase = yes
            noatime = yes
            
            # Word
            WildFile = "*.doc"
            WildFile = "*.dot"
            WildFile = "*.docx"
            WildFile = "*.docm"

            # Excel
            WildFile = "*.xls"
            WildFile = "*.xlt"
            WildFile = "*.xlsx"
            WildFile = "*.xlsm"
            WildFile = "*.xltx"
            WildFile = "*.xltm"

            # Powerpoint
            WildFile = "*.ppt"
            WildFile = "*.pot"
            WildFile = "*.pps"
            WildFile = "*.pptx"
            WildFile = "*.pptm"
            WildFile = "*.ppsx"
            WildFile = "*.ppsm"
            WildFile = "*.sldx"

            # access
            WildFile = "*.accdb"
            WildFile = "*.mdb"
            WildFile = "*.accde"
            WildFile = "*.accdt"
            WildFile = "*.accdr"

            # publisher
            WildFile = "*.pub"

            # open office
            WildFile = "*.odt"
            WildFile = "*.ods"
            WildFile = "*.odp"

            # pdf
            WildFile = "*.pdf"
            
            # flat text / code
            WildFile = "*.xml"
            WildFile = "*.log"
            WildFile = "*.rtf"
            WildFile = "*.tex"
            WildFile = "*.sql"
            WildFile = "*.txt"
            WildFile = "*.tsv"
            WildFile = "*.csv"
            WildFile = "*.php"
            WildFile = "*.sh"
            WildFile = "*.py"
            WildFile = "*.r"
            WildFile = "*.rProj"
            WildFile = "*.js"
            WildFile = "*.html"
            WildFile = "*.css"
            WildFile = "*.htm"
          } 

          # exclude everything else
            Options {
            
            # all files not in include
            RegExFile = ".*"
            
            # default user profiles
            WildDir = "[C-D]:/Users/All Users/*"
            WildDir = "[C-D]:/Users/Default/*"
            
            # explicit don't backup
            WildDir = "[C-D]:/Users/*/AppData"
            WildDir = "[C-D]:/Users/*/Music"
            WildDir = "[C-D]:/Users/*/Videos"
            WildDir = "[C-D]:/Users/*/Searches"
            WildDir = "[C-D]:/Users/*/Saved Games"
            WildDir = "[C-D]:/Users/*/Favorites"
            WildDir = "[C-D]:/Users/*/Links"
          
            # application specific
            WildDir = "[C-D]:/Users/*/MicrosoftEdgeBackups"
            WildDir = "[C-D]:/Users/*/Documents/R"
            WildDir = "*iCloudDrive*"
            WildDir = "*.svn/*"
            WildDir = "*.git/*"
            WildDir = "*.metadata/*"
            WildDir = "*cache*"
            WildDir = "*temp*"
            WildDir = "*OneDrive*"
            WildDir = "*RECYCLE.BIN*"
            WildDir = "[C-D]:/System Volume Information"
            Exclude = yes
          }
           
          }
        }

        FileSet {
          Name = "Win10_office"
          
          # volume shadow copy service
          Enable VSS = yes
          Include {
          
          # location
            File = "C:/Users"
          
          Options {
            # config
            Signature = MD5
            compression = LZ4
            IgnoreCase = yes
            noatime = yes
            
            # Word
            WildFile = "*.doc"
            WildFile = "*.dot"
            WildFile = "*.docx"
            WildFile = "*.docm"

            # Excel
            WildFile = "*.xls"
            WildFile = "*.xlt"
            WildFile = "*.xlsx"
            WildFile = "*.xlsm"
            WildFile = "*.xltx"
            WildFile = "*.xltm"

            # Powerpoint
            WildFile = "*.ppt"
            WildFile = "*.pot"
            WildFile = "*.pps"
            WildFile = "*.pptx"
            WildFile = "*.pptm"
            WildFile = "*.ppsx"
            WildFile = "*.ppsm"
            WildFile = "*.sldx"

            # access
            WildFile = "*.accdb"
            WildFile = "*.mdb"
            WildFile = "*.accde"
            WildFile = "*.accdt"
            WildFile = "*.accdr"

            # publisher
            WildFile = "*.pub"

            # open office
            WildFile = "*.odt"
            WildFile = "*.ods"
            WildFile = "*.odp"

            # pdf
            WildFile = "*.pdf"
            
            # flat text / code
            WildFile = "*.xml"
            WildFile = "*.log"
            WildFile = "*.rtf"
            WildFile = "*.tex"
            WildFile = "*.sql"
            WildFile = "*.txt"
            WildFile = "*.tsv"
            WildFile = "*.csv"
            WildFile = "*.php"
            WildFile = "*.sh"
            WildFile = "*.py"
            WildFile = "*.r"
            WildFile = "*.rProj"
            WildFile = "*.js"
            WildFile = "*.html"
            WildFile = "*.css"
            WildFile = "*.htm"
          } 

          # exclude everything else
            Options {
            
            # all files not in include
            RegExFile = ".*"
            
            # default user profiles
            WildDir = "[C-D]:/Users/All Users/*"
            WildDir = "[C-D]:/Users/Default/*"
            
            # explicit don't backup
            WildDir = "[C-D]:/Users/*/AppData"
            WildDir = "[C-D]:/Users/*/Music"
            WildDir = "[C-D]:/Users/*/Videos"
            WildDir = "[C-D]:/Users/*/Searches"
            WildDir = "[C-D]:/Users/*/Saved Games"
            WildDir = "[C-D]:/Users/*/Favorites"
            WildDir = "[C-D]:/Users/*/Links"
          
            # application specific
            WildDir = "[C-D]:/Users/*/MicrosoftEdgeBackups"
            WildDir = "[C-D]:/Users/*/Documents/R"
            WildDir = "*iCloudDrive*"
            WildDir = "*.svn/*"
            WildDir = "*.git/*"
            WildDir = "*.metadata/*"
            WildDir = "*cache*"
            WildDir = "*temp*"
            WildDir = "*OneDrive*"
            WildDir = "*RECYCLE.BIN*"
            WildDir = "[C-D]:/System Volume Information"
            Exclude = yes
          }
           
          }
        }
        ```
- jobdefs : 备份任务定义, 可被多个作业重复调用, 类似于job template

    ```conf
    JobDefs {
      Name = "TestJob"                                          # 测试任务
      Type = Backup                                             # 类型：备份（Backup）
      Level = Incremental                                       # 方式：递进（Incremental）
      Client = bareos-fd                                        # 被备份客户端：bareos-fd （在Client中定义）
      FileSet = "TestSet"                                       # 备份文件组：TesetSet （在FileSet中定义）
      Schedule = "WeeklyCycle"                                  # 备份周期：WeeklyCy（在schedule中定义）
      Storage = File                                            # 备份媒体： File（在Storage中定义）
      Messages = Standard                                       # 消息方式：Standard（在Message中定义）
      Pool = Incremental                                        # 存储池：Incremental（在pool中定义） 
      Priority = 10                                             # 优先级：10
      Write Bootstrap = "/var/lib/bareos/%c.bsr"                # 
      Full Backup Pool = Full                  # Full备份, 使用 "Full" 池（在storage中定义）
      Differential Backup Pool = Differential  # Differential备份, 使用 "Differential" 池（在storage中定义）
      Incremental Backup Pool = Incremental    # Incremental备份, 使用 "Incremental" 池（在storage中定义）
    }
    ```


    备份类型:
    - Full : 备份整个文件
    - Incremental : 备份状态变化的文件
    - Differential : 备份修改了（modified标志变化）的文件
- job : 任务配置

    任务类型分: Backup(备份)/Restore(还原), 默认存在的backup-bareos-fd.conf和BackupCatalog.conf是备份job, RestoreFiles.conf是还原job.

    ```conf
    Job {
      Name = "backup-test-on-bareos-fd"              # 任务名
      JobDefs = "TestJob"                            # 使用已定义的备份任务TestJob （在jobdefs中定义）
      Client = "bareos-fd"                           # 客户端名称： bareos-fd（在client中定义）
    }
    ```
- storage : 备份保存位置的配置

    ```conf
    Storage {
      Name = File
      Address = bareos                # director-sd名字, 使用FQDN (不要使用 "localhost" ).
      Password = "JgwtSYloo93DlXnt/cjUfPJIAD9zocr920FEXEV0Pn+S" # 来自sd daemon的director/bareos-dir.conf#Password
      Device = FileStorage            # 在bareos-sd中定义
      Media Type = File
    }
    ```

    > Device, Media Type项必须与bareos-sd定义的一致
- pool : pool配置

    - full : 完整备份

        ```conf
        Pool {
          Name = Full
          Pool Type = Backup
          Recycle = yes                       # Bareos 自动回收重复使用 Volumes（Volume备份文件标记）
          AutoPrune = yes                     # 自动清除过期的Volumes
          Volume Retention = 365 days         # Volume有效时间
          Maximum Volume Bytes = 50G          # Volume最大尺寸
          Maximum Volumes = 100               # 单个存储池允许的Volume数量
          Label Format = "Full-"              # Volumes 将被标记为 "Differential-<volume-id>"
        }
        ```
    - incremental : 增量备份, 备份所有状态变化的文件. 前提是有full备份, 否则转为full备份.

        ```conf
        Pool {
          Name = Incremental
          Pool Type = Backup
          Recycle = yes                       # Bareos 自动回收重复使用 Volumes（Volume备份文件标记）
          AutoPrune = yes                     # 自动清除过期的Volumes
          Volume Retention = 30 days          # Volume有效时间
          Maximum Volume Bytes = 1G           # Volume最大尺寸
          Maximum Volumes = 100               # 单个存储池允许的Volume数量
          Label Format = "Incremental-"       # Volumes 将被标记为 "Differential-<volume-id>"
        }
        ```
    - differential : 差异备份, 备份所有modified标志变化的文件. 前提是有full备份, 否则转为full备份.

        ```conf
        Pool {
          Name = Differential
          Pool Type = Backup
          Recycle = yes                       # Bareos 自动回收重复使用 Volumes（Volume备份文件标记）
          AutoPrune = yes                     # 自动清除过期的Volumes
          Volume Retention = 90 days          # Volume有效时间
          Maximum Volume Bytes = 10G          # Volume最大尺寸
          Maximum Volumes = 100               # 单个存储池允许的Volume数量
          Label Format = "Differential-"      # Volumes 将被标记为 "Differential-<volume-id>"
        }
        ```
    - scratch: 当系统找不到需要的volume时, 自动使用该pool. 该pool名称不可修改, 其他pool名称没有重命名限制.
- schedule: 计划配置

    ```conf
    Schedule {
      Name = "WeeklyCycle"
      Run = Full 1st sat at 21:00                   # 每月第一个周六/晚九点, 完整备份
      Run = Differential 2nd-5th sat at 21:00       # 其余周六/晚九点, 差异备份
      Run = Incremental mon-fri at 21:00            # 周一至周五, 递增备份
      Run = daily at 21:10 # 每天21:10备份
      Run = Incremental hourly at 0:22 # 每小时0:22备份
    }
    ```
- message : 提示信息(job完成后如何发送提示信息)的配置

    ```conf
    Messages {
      Name = Standard
      Description = "Reasonable message delivery -- send most everything to email address and to the console."
      # operatorcommand = "/usr/bin/bsmtp -h localhost -f \"\(Bareos\) \<%r\>\" -s \"Bareos: Intervention needed for %j\" %r"
      # mailcommand = "/usr/bin/bsmtp -h localhost -f \"\(Bareos\) \<%r\>\" -s \"Bareos: %t %e of %c %l\" %r"
      operator = root@localhost = mount                                 # 执行operatorcommand命令, 用户：root@localhost, 操作：mount
      mail = root@localhost = all, !skipped, !saved, !audit             # 执行mailcommand, 用户：root@localhost, 操作：所有（除skipped, saved和audit）, **注释该行即可取消发送email**
      console = all, !skipped, !saved, !audit                           # 所有操作, 除skipped, saved和audit
      append = "/var/log/bareos/bareos.log" = all, !skipped, !saved, !audit  # 所有操作, 除skipped, saved和audit
      catalog = all, !skipped, !saved, !audit                           # 所有操作, 除skipped, saved和audit
       # 可用参数
      # %% = %
      # %c = Client’s name
      # %d = Director’s name
      # %e = Job Exit code (OK, Error, ...)
      # %h = Client address
      # %i = Job Id
      # %j = Unique Job name
      # %l = Job level
      # %n = Job name
      # %r = Recipients
      # %s = Since time
      # %t = Job type (e.g. Backup, ...)
      # %v = Read Volume name (Only on director side)
      # %V = Write Volume name (Only on director side)
      # console：定义发送到console的信息
      # append：定义发送到日志文件的信息
      # catalog：定义发送到数据库的信息
    }
    ```
- profile : 定义一组访问控制用于针对不同控制台或角色

### fileset
- `One FS=no` : no, 不检查是否在同一个fs上; yes, 检查是否在同一个fs上
- `FS Type=ext4` : 支持备份的fs类型
- `File=/` : 备份开始位置
- `Exclude {}` : 排除位置
- `WildDir` : 指定文件
- `Exclude = yes`: 排除`WildDir`指定的文件

### backup参数
```conf
Run Backup job
JobName:  backup-test-on-bareos-fd
Level:    Full
Client:   lswin7-1-fd
Format:   Native
FileSet:  TestSet
Pool:     Full (From Job FullPool override)
Storage:  File (From Job resource)
When:     2018-10-05 10:39:59
Priority: 10
OK to run? (yes/mod/no):
```
### restore参数
```conf
Run Restore job
JobName:         RestoreFiles
Bootstrap:       /var/lib/bareos/client1.restore.3.bsr
Where:           /tmp/bareos-restores
Replace:         Always
FileSet:         Full Set
Backup Client:   client1
Restore Client:  client1
Format:          Native
Storage:         File
When:            2013-06-28 13:30:08
Catalog:         MyCatalog
Priority:        10
Plugin Options:  *None*
OK to run? (yes/mod/no):
```

### bconsole命令行调用形式
bconsole是交互式命令, 无法直接后接子命令的形式试用, 因此使用:
```bash
bconsole -c ./bconsole.conf <<END_OF_DATA
show pool
quit
END_OF_DATA
```

[组合使用(备份+还原)](https://docs.bareos.org/TasksAndConcepts/BareosConsole.html#running-the-console-from-a-shell-script):
```bash
bconsole <<END_OF_DATA
@output /dev/null
messages
@output /tmp/log1.out
label volume=TestVolume001
run job=Client1 yes
wait
messages
@#
@# now do a restore
@#
@output /tmp/log2.out
restore current all
yes
wait
messages
@output
quit
END_OF_DATA
```

### job执行过程中报`BnetHost2IpAddrs() for host "ubuntu-18" failed: ERR=`
ubuntu-18是storage daemon的参数在`/etc/bareos/bareos-dir.d/storage/File.conf`的`Address`.

file daemon备份时, 从dareos-dir获取storage参数, 因为网络中没有dns, 因此无法获取到storage的ip.

解决方法: 将Address的参数换成ip即可.

> 错误来源: `/var/log/bareos/bareos.log`或 webui中job的log

### job备份windows文件时报`no drive letters found for generating vss snapshots`
fileset中备份文件路径错误.

错误路径: `File=/c/dsDefault.log`, 正确路径: `File="C:/dsDefault.log"`

### job备份Windows 10文件报`error:14094417:SSL routines:ssl3_read_bytes:sslv3 alert illegal parameter`, `TLS negotiation failed(while probing client protocol)`和`Network error during CRAM MD5 With 192.168.0.197`

此时Windows 10 log报"SSL routines:tls_psk_do_binder:binder does not verify", `TLS negotiation failed`.

解决方法: 卸载并重新安装bareos windows client, 安装时填入正确的参数即可.

> 出问题时安装是使用默认参数(即错误参数), 安装完成后修正`C:\Program Files\Bareos\defaultconfigs\bareos-fd.d\director`下的`*.conf`并重启`Bareos File Backup Service`进行配置的.

### `status client=xxx` bareos.log报`Network error during CRAM MD5 with 192.168.0.130\nUnable to authenticate with File daemon at "192.168.0.130:9102"`, client端报:`TLS negotiation failed`, `error:1408F119:SSL routines:ssl3_get_record:decryption failed or bad record mac`
client在win10上.

已尝试重装client, 或重新配置client的director-dir.conf, 均无效.

这个错误的原因是openssl 在对收到的包做完整性校验时发现收到的报数据不对. 调查时需要从两端同时抓取wireshrk包分析，看到底哪里将数据破坏掉了.

### 修改Director邮件发送命令
参考:
- [备份/恢复系统BAREOS的安装、设置和使用（四）](https://blog.csdn.net/laotou1963/article/details/82939355)

在Director默认使用bsmtp发送邮件, 由于bsmtp的局限性，无法使用一般外部商业SMTP服务，我们必须对此进行修改。在示例中，我们对/etc/bareos/bareos-dir.d/message/Standard.conf做修改，您可以参照示例，对其他的邮件发送配置做对应的修改。

配置文件中的默认邮件命令为：
`mailcommand = "/usr/bin/bsmtp -h localhost -f \"\(Bareos\) \<%r\>\" -s \"Bareos: %t %e of %c %l\" %r"`

改为: `mailcommand = "/usr/local/bin/sendmail -c %c -d %d -e %e -h %h -i %i -j %j -n %n -r %r -t %t -s \"%s\"  -l %l -v \"%v\" -V \"%V\%"`

`/user/local/bin/sendmail`是自定义的发送邮件脚本程序. 以`%`开头的是在Bareos中可用的参数，可把所有可用参数全部传递到脚本程序.

> ps: `%s、%v和%V`用`" "`包起来的原因是，这些参数有可能为空，如不把它们包起来，当它们为空时，会造成参数处理问题.

```bash
#!/usr/bin/env bash
# available mailcommand parameters
# %% = %
# %c = Client’s name
# %d = Director’s name
# %e = Job Exit code (OK, Error, ...)
# %h = Client address
# %i = Job Id
# %j = Unique Job name
# %l = Job level
# %n = Job name
# %r = Recipients
# %s = Since time
# %t = Job type (e.g. Backup, ...)
# %v = Read Volume name (Only on director side)
# %V = Write Volume name (Only on director side)

bareos_admin="admin@lswin.cn"
mail_receiver="s.zhang@lswin.cn"

# get input opts
while getopts ":c:d:e:h:i:j:l:n:r:s:t:v:V:" o; do
  case "${o}" in
    c)
       client_name=${OPTARG}
       ;;
    d)
       director_name=${OPTARG}
       ;;
    e)
       job_exit_code=${OPTARG}
       ;;
    h)
       client_address=${OPTARG}
       ;;
    i)
       job_id=${OPTARG}
       ;;
    j)
       unique_job_name=${OPTARG}
       ;;
    l)
       job_level=${OPTARG}
       ;;
    n)
       job_name=${OPTARG}
       ;;
    r)
       recipients=${OPTARG}
       ;;
    s)
       since_time=${OPTARG}
       ;;
    t)
       job_type=${OPTARG}
       ;;
    v)
       read_volume_name=${OPTARG}
       ;;
    V)
       write_volume_name=${OPTARG}
       ;;
    *)
       ;;
    esac
done

# 建立邮件 Subject
ubject="BAREOS任务执行"
if [[ "$job_exit_code" == "OK" ]]
then
  Subject=$Subject"完成通知"
else
  Subject=$Subject"失败通知！"
fi

# 建立邮件内容
Content="\"任务 "$job_name" 执行简况:\n 任务ID："$job_id"\n 任务名字："$unique_job_name"\n 任务类型："$job_type
if [[ ! -z "$job_level" && "$job_type" == "Backup" ]]; then Content=$Content"\n 备份级别："$job_level; fi
Content=$Content"\n 完成情况："$job_exit_code"\n 主控端名字："$director_name"\n 客户端名字："$client_name"\n 客户端地址："$client_address
if [[ ! -z "$read_volume_name" && "$job_type" == "RestoreFiles" ]]; then Content=$Content"\n 读取卷名字："$read_volume_name; fi
if [[ ! -z "$write_volume_name" && "$job_type" == "Backup" ]]; then Content=$Content"\n 写入卷名字："$write_volume_name; fi
Content=$Content"\""

# 建立邮件发送命令
cmd="echo -e $Content | /usr/bin/mail -s \"Subject: $Subject\" -r $bareos_admin $mail_receiver"

# 执行邮件发送命令
eval $cmd

exit 0
```

email example:
```conf
Subject: BAREOS任务执行完成通知

发件人：admin <admin@lswin.cn>      
时   间：2018年10月18日(星期四) 上午10:26  纯文本 |  
收件人：
S Zhang<s.zhang@lswin.cn>
任务 backup-bareos-fd 执行简况:
 任务ID：52
 任务名字：backup-bareos-fd.2018-10-18_10.26.39_12
 任务类型：Backup
 备份级别：Full
 完成情况：OK
 主控端名字：bareos-dir
 客户端名字：bareos-fd
 客户端地址：localhost
 写入卷名字：Full-0001

# ----
Subject: BAREOS任务执行失败通知！

发件人：admin <admin@lswin.cn>      
时   间：2018年10月18日(星期四) 上午10:45  纯文本 |  
收件人：
S Zhang<s.zhang@lswin.cn>
任务 backup-test-on-bareos-fd 执行简况:
 任务ID：53
 任务名字：backup-test-on-bareos-fd.2018-10-18_10.42.13_17
 任务类型：Backup
 备份级别：Full
 完成情况：Error
 主控端名字：bareos-dir
 客户端名字：lscms-fd
 客户端地址：lscms.lswin.cn

 # ---
Subject: BAREOS任务执行完成通知

发件人：admin <admin@lswin.cn>      
时   间：2018年10月18日(星期四) 上午10:45  纯文本 |  
收件人：
S Zhang <s.zhang@lswin.cn>
任务 RestoreFiles 执行简况:
 任务ID：54
 任务名字：RestoreFiles.2018-10-18_10.43.18_37
 任务类型：Restore
 完成情况：OK
 主控端名字：bareos-dir
 客户端名字：bareos-fd
 客户端地址：localhos

# ---
Subject: BAREOS任务执行失败通知！

发件人：admin <admin@lswin.cn>      
时   间：2018年10月18日(星期四) 上午10:45  纯文本 |  
收件人：
S Zhang<s.zhang@lswin.cn>
任务 RestoreFiles 执行简况:
 任务ID：55
 任务名字：RestoreFiles.2018-10-18_10.44.20_01
 任务类型：Restore
 完成情况：Error
 主控端名字：bareos-dir
 客户端名字：lswin7-1-fd
 客户端地址：lswin7-1.lswin.cn
```
### BVFS
BVFS（Bareos虚拟文件系统）提供了一个API来浏览目录中的备份文件并选择文件进行恢复.

### bareos webui如何获取data
以job列表页`localhost:9100/job/`举例, 找到其ajax req(`localhost:9100/job/getData/?data=jobs&period=7&sort=jobid&order=desc`)

在bareos webui root(`/usr/share/bareos-webui/module/Job`)下执行`grep -r getData`, 在`src/Job/Controller/JobController.php`中找到`getDataAction()`, 再在其中找到关键函数`getJobs`.

执行`grep -r getJobs`, 在`src/Job/Model/JobModel.php`中找到它, 看其实现基本可推断是基于bsock, 通过`$bsock->send_command()->send()`逆推, 在`src/Job/Controller/JobController.php`中找到`$this->bsock=$this->getServiceLocator()->get('director')`.

在`/usr/share/bareos-webui`执行`grep -r "send_command" |grep -v "bsock"`, 在`vendor/Bareos/library/Bareos/BSock/BareosBSock.php`找到其实现(需考虑send_command有参数列表). 在找到它的上层函数send(), 发现它是操作`fwrite($this->socket,...)`, 找到socket定义: [`stream_socket_client()`](https://php.golaravel.com/function.stream-socket-client.html).

截获bareos cmd: 在BareosBSock.php的send()开头添加打印语句:`error_log("[". date("Y-m-d H:i:s", time()) ."] : $msg \n", 3, "/tmp/bareos_cmd.log");`.

### bareos python sdk截获cmd
1. 根据bareos-restapi.py的`current_user.jsonDirector.call()`找到`self.jsonDirector = bareos.bsock.DirectorConsoleJson`
1. 为`DirectorConsoleJson.call()`添加打印即可, 比如`pprint(command)`

### log
使用`-d 500 -v`参数, 可打印详细日志

bareos-dird log在`/var/log/bareos/bareos.log`
bareos-fd log在systemd.

/var/log/bareos/bareos-audit.log是bareos dir的审计日志, 比如bconsole执行的命令.

### 使用官方plugin [bareos-fd-mysql](https://docs.bareos.org/Appendix/Howtos.html#backup-mysql-python)执行job时报`... PluginSave: Command plugin "<python plugin>" required, but is not loaded`
fd `/etc/bareos/bareos-fd.d/client/myself.conf`配置:
```
Client {
  ...

  # remove comment from "Plugin Directory" to load plugins from specified directory.
  # if "Plugin Names" is defined, only the specified plugins will be loaded,
  # otherwise all filedaemon plugins (*-fd.so) from the "Plugin Directory".
  #
  Plugin Directory = "/usr/lib/bareos/plugins"
  Plugin Names = "python"

  ...
}
```

使用`-d 500`参数, 打印详细日志可见, fd log提示`field/fd_plugins.cc:1750-0 No plugin loaded`.

结合myself.conf和日志调试发现, 只要启用了`Plugin Names`即使其value为空, 均会按`Plugin Names`指定的名称去load plugin. 将`Plugin Names`注释默认加载全部插件即可.

### 使用自编译bareos 20.0.1 arm版本, linux备份还原正常, 官方对应版本的windows client无法备份
dir, sd, fd均无报错.

### `configure: is an invalid command.`
通过bareos resp-api创建client报错, 通过修改`bareos-restapi.py`打开print来获取到具体调用的命令, 发现相同的命令在bconsole执行成功, 且监控`/var/log/bareos/bareos-audit.log`发现报`Audit acl failure for Command configure`.

查看`bareos-dir.d/console/admin.conf`发现它使用了`bareos-dir.d/profile/web-admin.conf`, 而web-admin.conf的acl中禁用了configure.

解决方法:
更新web-admin.conf的acl, 取消禁用configure. 需要重启bareos-dir.service(bconsole的reload命令无效).

### job‘s jobstatus
定义在`/usr/share/bareos-webui/public/js/bootstrap-table-formatter.js`, 对应的翻译在`/usr/share/bareos-webui/module/Application/language/cn_CN.po`.

或在[BareosDirPluginPrometheusExporter.py](https://github.com/bareos/bareos/blob/master/contrib/dir-plugins/prometheus/BareosDirPluginPrometheusExporter.py)

### bareos client
- golang

    - [barethoven](https://github.com/myENA/barethoven)

### systemd显示bareos-sd运行中但实际bareos-sd未执行(未监听端口)
bareos-sd所在host宕机重启后出现该现象. 原因: bareos-sd的pidfile是持久化的, 宕机后该pidfile未清理.

修改bareos-sd.service的PIDFile=/run/xxx.pid, 发现`systemctl start bareos-sd`无法启动.

解决方法: 监控bareos-sd是否监听了端口, 否则执行`systemctl restart bareos-sd`

### run job时joblog卡住, `status storage=xxx`显示`Device is BLOCKED waiting for mount of volume "Full-0010"`
解决方法:
1. list volume
2. purge volume=Full-0010 yes
3. 在Full-0010所在storage执行`systemctl restart bareos-sd`

### 使用官方bareos-webui nginx配置可能访问`localhost:9100`空白
env: php-fpm 7.2

安装php-fpm后会生成`/etc/nginx/default.d/php.conf`, bareos-webui.conf中的`location ~ \.php$`需要使用`php.conf`配置中的`fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;`:
```
 include snippets/fastcgi-php.conf;

                # php5-cgi alone:
                # pass the PHP
                # scripts to FastCGI server
                # listening on 127.0.0.1:9000
                #fastcgi_pass 127.0.0.1:9000;

                # php5-fpm:
                fastcgi_pass unix:/var/run/php5-fpm.sock;

                # APPLICATION_ENV:  set to 'development' or 'production'
                #fastcgi_param APPLICATION_ENV development;
                fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name; # 脚本文件请求的路径,也就是说当访问127.0.0.1/index.php的时候，需要读取网站根目录下面的index.php文件，如果没有配置这一配置项时，nginx不回去网站根目录下访问.php文件，所以返回空白 
                fastcgi_param APPLICATION_ENV production;
```

### bconsole执行`status client=bareos-fd`报`Probing client protocol... (result will be saved until config reload)`
看`/var/log/bareos/bareos.log`提示`Unable to authenticate with File daemon at "localhost:9102"`, 是dir中client的配置的Password字段错误, 用`/etc/bareos/bareos-fd.d/director/bareos-dir.conf`中Password替换即可.

### bareos-tray-monitor的开机自启
`cat /etc/xdg/autostart/bareos-tray-monitor.desktop`

### oracle linux 7.9构建bareos 21.1.2报`Target xxx requires the language dialect "CXX17"`
gcc版本不够.

```bash
yum install centos-release-scl -y
yum install devtoolset-9 -y

# 临时覆盖系统原有的gcc引用
scl enable devtoolset-9 bash
gcc --version

# 永久
echo "source /opt/rh/devtoolset-9/enable" >>/etc/profile
```

### bareos 21.1.2执行备份vmware vm报`Fatal error: filed/fd_plugins.cc:670 PluginSave: Command plugin "python:module_path=..." requested, but is not loaded`
bareos-fd的client/myself.conf是`Plugin Names = "python"`, 而vmware plugin`bareos-fd-vmware.py`是python3, 因此将其改为`Plugin Names = "python3"`即可.

### [bareos备份vmware](https://docs.bareos.org/TasksAndConcepts/Plugins.html#vmware-plugin)
1. 先用`vmware_cbt_tool.py`将要备份的 VM 启用 CBT（更改块跟踪）

    没开启并备份时joblog会报`No snapshot was taken, skipping snapshot removal`
2. 其他的参考文档

### bareos vmware备份的vm不还原到vmware
ref:
- [Restore VmWare VM by bareos](http://www.voleg.info/bareos-restore-vmware.html)
- [Backup VM ESXi using Bareos](https://sudonull.com/post/76101-VM-ESXi-backup-using-Bareos-SIM-Networks-Blog)

还原时默认会被还原到原有vm位置并覆盖它的存储, 前提时该vm已关键.

还原成文件的方法: `run ... pluginoptions=python:localvmdk=yes`.

> 通过bconsole手动还原时选择修改restore job的"Plugin Options"为`python:localvmdk=yes`.

### bareos vmware不能同一时刻多个client备份同一台vm

### bareos vmware如何避免还原时需到vmware环境下使用vmkfstools转换格式(未完成)
ref:
- [KVM虚拟机迁移到VMWare ESXi](https://blog.csdn.net/avatar_2009/article/details/117769202)
- [通过qemu-img工具转换镜像格式](https://support.huaweicloud.com/bestpractice-ims/ims_bp_0030.html)
- [Virtual Disk Types](https://vdc-repo.vmware.com/vmwb-repository/dcr-public/6335f27c-c6e9-4804-95b0-ea9449958403/c7798a8b-4c73-41d9-84e8-db5453de7b17/doc/vddkDataStruct.5.3.html)

```bash
# qemu-img info centos6.9-64bit.vmdk # centos6.9-64bit.vmdk是bareos还原vm到本地时的文件
...
    create type: monolithicSparse
...
# vim -R <bareos src>/core/src/vmware/vadp_dumper/bareos_vadp_dumper.cc # 有vmfs_thin, 修改BareosFdPluginVMware.py启用`bareos_vadp_dumper_opts["dump"] = "-S -D -M -t vmfs_thin "`, 经测试后无效果.
```

在 ESX/ESXi 主机上, VMDK 文件的子格式类型为 VMFS_FLAT 或 VMFS_THIN(适合放在nfs上), `qemu-img convert`不支持这两种格式.

> ESXi 格式的虚拟磁盘由两个单独的文件组成: 一个数据文件和一个磁盘描述符文件.

> VMware Workstation 和 VMware ESXi 的 VMware 虚拟磁盘格式是另一回事. VMware Workstation 格式的虚拟磁盘具有内置于单个 VMDK 文件中的磁盘描述符.