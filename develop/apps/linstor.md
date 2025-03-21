# linstor

[LINBIT DRBD9 Stack官方PPA](https://launchpad.net/~linbit/+archive/ubuntu/linbit-drbd9-stack)

## 构建
LINBIT未直接提供构建每个linstor组件的方法, 因此构建方法是: 每个项目的`README + .gitlab-ci.yml + Dockerfile` + google

### server

```bash
sdk install gradle 6.8.2
gradle -v
sudo apt install openjdk-11-jdk # base https://github.com/LINBIT/linstor-server/blob/master/build.gradle的sourceCompatibility, 推荐用openjdk-8-jdk, 我这是因为当前已安装openjdk-11-jdk
sudo apt default-jdk-headless dh-systemd python3-all
git clone --depth=1 git@github.com:LINBIT/linstor-server.git -b v1.11.1 # 必须使用git repo否则gradlew会报错
cd linstor-server
make check-submods
```

使用方法:
```bash
./gradlew getProtoc
./gradlew assemble
ll build/distributions # 构建结果, 非deb/rpm
具体使用方法见https://github.com/LINBIT/linstor-server/blob/master/Dockerfile.test-controller
# ./gradlew clean # gradle重新构建时可用于清理

# 继续构建deb:
删除https://github.com/LINBIT/linstor-server/blob/master/debian/control中gradle的限制 # 上面已安装gradle
LD_LIBRARY_PATH='' dpkg-buildpackage -rfakeroot -b -uc # 自定义的LD_LIBRARY_PATH可能影响dpkg-buildpackage的构建
```

### client
```bash
pip3 install python-linstor natsort
git clone --depth=1 git@github.com:LINBIT/linstor-client.git -b v1.6.1
cd linstor-client
NO_DOC="-no-doc" make debrelease # from `.gitlab-ci.yml`
删除https://github.com/LINBIT/linstor-client/blob/master/debian/control中的python-linstor, python3-natsort依赖
make deb
```

### drbd
```bash
到https://launchpad.net/~npalix/+archive/ubuntu/coccinelle下载coccinelle
apt install spatch # spatch需>= 1.0.8, deepin 20.1的是1.0.4
git clone --depth=1 git@github.com:LINBIT/drbd.git -b drbd-9.0.27-1
cd drbd
make clean
make
sudo make install
modprobe drbd
# sudo apt install module-assistant # 构建deb需要
# LD_LIBRARY_PATH='' dpkg-buildpackage -rfakeroot -nc -uc -us # 打出的deb内容不对
```

### drbd-utils
```bash
apt install po4a clitest
git clone --depth=1 git@github.com:LINBIT/drbd-utils.git -b v9.15.1
cd drbd-utils
make -f Makefile.in check-submods
./autogen.sh
./configure --prefix=/usr --localstatedir=/var --sysconfdir=/etc --without-83support --without-84support # 根据autogen.sh的输出选择合适的configure参数, 部分参数来自https://github.com/LINBIT/drbd-utils/blob/master/.gitlab-ci.yml
make debrelease
make
sudo make install # 无需deb时, 直接install即可
vim debian/rules
configure_with = --prefix=/usr --localstatedir=/var --sysconfdir=/etc --without-83support --without-84support # 注释文中configure_with操作, 并使用上面的参数替代
# configure_with = --prefix=/usr --localstatedir=/var --sysconfdir=/etc \
#       --sbindir=/usr/sbin --with-udev --with-xen \
#       --with-pacemaker --with-rgmanager --without-bashcompletion

# ifneq ($(WITH_SYSTEMD),)
# configure_with += --with-systemdunitdir=/lib/systemd/system \
#       --with-initscripttype=both
# # alternatively: --with-initscripttype=systemd,
# # if you want to drop the sysv script from the package. 
# # Need to adjust below, in that case.
# endif

# ifeq (,$(filter noprebuiltman,$(DEB_BUILD_OPTIONS)))
# configure_with += --with-prebuiltman
# endif
LD_LIBRARY_PATH='' dpkg-buildpackage -rfakeroot -nc -uc # 使用`-nc`避免重新构建
```

## linstor-client
参考:
- [DRBD9 and LINSTOR the easy way](https://pub.nethence.com/storage/drbd-linstor)

```bash
linstor node restore <node> # 重新注册node
linstor physical-storage list # 罗列node上的disk
linstor storage-pool list # 已注册到linstor的pool
zpool create mypool /dev/sdb # 在每个node上执行
linstor storage-pool create zfs <node> mypool <pool_name on node> # 将底层pool注册为linstor pool. zpool需先手动创建, 不推荐使用`linstor physical-storage create-device-pool`, 因为linstor本身不维护zfs/lvm pool. zfsthin是zfs pool与zfs类型没区别, 但建vol时都是thin vol.

pvcreate /dev/vdb
vgcreate drbdpool /dev/vdb
linstor storage-pool create lvm nodeA drbdpool drbdpool # 注册lvm

lvcreate -L 800G --thinpool drbdpool pve # 创建lvm thin pool, 它是建立在vg上的
linstor storage-pool create lvmthin pve1 drbdpool pve/drbdpool

linstor resource-definition create demo # 创建名为demo的resource-definition
linstor volume-definition create demo 15G # 指定demo大小
linstor volume-definition list # volume-definition列表 
linstor resource create nodeA demo --storage-pool mypool # 手动创建resource
linstor resource create demo --auto-place 2 # 自动创建resource副本, `--layer-list storage`可只创建底层vol而没有drbd

linstor resource list # 资源列表
linstor volume list # volume列表, 包drbd devicename

# resource-group是volume-definition的父对象，其中对资源组所做的所有属性更改都将由其资源定义的子级继承
# 继承设置的层次结构: 卷定义 设置优先于 卷组 设置， 资源定义 设置优先于 资源组 设置
linstor resource-group create my_group --storage-pool mypool --place-count 3 # `--place-count`分布在n个node上
linstor resource-group drbd-options --verify-alg crc32c my_verify_group # 设置drbd选项
linstor volume-group create my_group
linstor resource-group spawn-resources my_group my_res 5G # 依据my_group创建resource, 此时drbd role都是secondary
```

## linstor-gateway
linstor-gateway编译出来后重命名为linstor-iscsi/linstor-nfs即可使用.

## FAQ
###
```bash
$ LD_LIBRARY_PATH='' dpkg-buildpackage -rfakeroot -b -uc # for drbd-utils
...
make -C documentation/v9 doc
make[2]: 进入目录“/home/chen/test/drbd-utils/documentation/v9”
test -f drbdsetup.8
test -f drbd.conf.5
test -f drbd.8
make[2]: *** [../../documentation/common/Makefile_v9_com_post:179：drbdsetup.8] 错误 1
make[2]: *** 正在等待未完成的任务....
make[2]: *** [../../documentation/common/Makefile_v9_com_post:49：drbd.conf.5] 错误 1
make[2]: *** [../../documentation/common/Makefile_v9_com_post:49：drbd.8] 错误 1
make[2]: 离开目录“/home/chen/test/drbd-utils/documentation/v9”
make[1]: *** [Makefile:112：doc] 错误 2
make[1]: 离开目录“/home/chen/test/drbd-utils”
dh_auto_build: make -j4 returned exit code 2
make: *** [debian/rules:10：build] 错误 2
dpkg-buildpackage: error: debian/rules build subprocess returned exit status 2
```

`make -C documentation/v9 doc`执行报错, 导致drbdsetup.8等文件未生成, 原因是未执行`make debrelease`. 恢复到git clone时的状态, 按照drbd-utils的流程重新构建即可.

### linstor-client输出api调用信息
client sub cmds在`linstor_client/commands/*.py`里, 它们都是调用了`/usr/lib/python3/dist-packages/linstor/linstorapi.py`里的接口, 因此:
```python
logging.basicConfig(level=logging.WARNING) => logging.basicConfig(level=logging.DEBUG)

...
    def _rest_request_base(self, apicall, method, path, body=None, reconnect=True):
        ....
            self._rest_conn.request(
                method=method,
                url=path,
                body=json.dumps(body) if body is not None else None,
                headers=headers
            )
            self._logger.info("method: {}, url: {}, body: {}, headers: {}".format(method, path, body, headers)) # append

    def __convert_rest_response(self, apicall, response, path):
        ...
            data = json.loads(resp_data)
            self._logger.info("resp.data: {}".format(data)) # append
```

> linstor -m node list # `-m`输出原始resp json

### linstor-server输出api调用信息
server api源码入口在`linstor-server/controller/src/main/java/com/linbit/linstor/api/rest/v1/*.java`, 发现代码未埋点, 改用logger

[修改`/etc/linstor/linstor.toml`](https://github.com/LINBIT/linbit-documentation/blob/master/UG9/en/administration-linstor.adoc#logging):
```toml
[logging]
   level="TRACE"
```

### ha
参考:
- [Highly available LINSTOR Controller with Pacemaker](https://www.linbit.com/blog/linstor-controller-pacemaker/)
- [LINSTOR high availability](https://www.linbit.com/drbd-user-guide/linstor-guide-1_0-cn/#s-linstor_ha)

默认情况下, linsor cluster只能有一个活动的controller(上面注册了nodes), 其他nodes即使安装了controller, 执行`linstor node list`也是返回空.

linstor ha是通过drbd复制实现(h2数据库)的, 见[LINSTOR high availability](https://www.linbit.com/drbd-user-guide/linstor-guide-1_0-en/#s-linstor_ha)

> 默认重启linstor-satellite会清除drbd resource, 需添加env LS_KEEP_RES=linstor避免.

ha部署(base drbd-reactor(原名: drbdd)):
1. 将所有satellites注册到某个controller作为primary node
1. 准备linstor_db
```bash
# linstor resource-definition create linstor_db
# linstor volume-definition create linstor_db 200M
# linstor resource create linstor_db -s mypool --auto-place 3
```

> 不要手动分配linstor_db(drbd resource), 否则它的drbd minor会和以后controller分配的drbd minor重复从而引发错误

1. 迁移linstor db
```bash
# systemctl disable --now linstor-controller # all nodes
# systemctl disable linstor-controller # all nodes
# cat << EOF | sudo tee -a /etc/systemd/system/var-lib-linstor.mount # copy /etc/systemd/system/var-lib-linstor.mount to all other nodes
[Unit]
Description=Filesystem for the LINSTOR controller

[Mount]
# you can use the minor like /dev/drbdX or the udev symlink
What=/dev/drbd/by-res/linstor_db/0
Where=/var/lib/linstor
EOF
# mv /var/lib/linstor{,.orig}
# mkfs.ext4 /dev/drbd/by-res/linstor_db/0
# systemctl start var-lib-linstor.mount
# cp -r /var/lib/linstor.orig/* /var/lib/linstor # 保留备份, 避免意外
# systemctl start linstor-controller # primary node
# apt install drbdd # all nodes
# cat << EOF | sudo tee -a /etc/drbdd.toml  # all nodes, linstor_db is drbd resource
[promoter.resources.linstor_db]
start = ["var-lib-linstor.mount", "linstor-controller.service"]
EOF
# systemctl enable drbdd # all nodes
# systemctl restart drbdd # all nodes
# systemctl edit linstor-satellite # all node
[Service]
Environment=LS_KEEP_RES=linstor
[Unit]
After=drbdd.service
# systemctl restart linstor-satellite # all node
```

## LINBIT® CloudStack® HCI
### 部署
ref:
- [Ready-made CloudStack Deployment on Hyperconverged Infrastructure with SDS, HA, and DR Capabilities](https://linbit.com/blog/ready-made-cloudstack-deployment-on-hyperconverged-infrastructure-with-sds-ha-and-dr-capabilities/)

要求:
1. 抹除除系统盘(系统盘是安装系统时处理)外所有磁盘上的分区和文件系统, 特别是lvm信息, 因为HCI要求是空盘, lvm分区被识别后就无法抹除了, 会报`Device or resource busy`

其他:
1. 部署过程是由linstor-applian实现的
1. 各安装步骤中的ip range不能重合
1. 安装完成iso后, 建议打快照, 方便部署失败后重来

配置步骤:
1. `#step3 : Enter a Name for Your Cluster`

    会将`<ClusterName><序号>`作为hostname

    - `failed to list nodes from LINSTOR: failed to connect: could not connect to any controller`: `tail -f /var/log/message`发现是linstor-controller服务没有运行, 但首次安装好iso后, 第一次启动vm时, linstor-controller本身未自动启动 (linstor-controller将由负责ha的drbd-reactor管理)

        未知, 建议重新安装iso
    - `failed to list nodes from LINSTOR: failed to list resources of resource definition "linstor_db": 404 Not Found`

        未知, 建议重新安装iso
1. `CloudStack Node Setup`

    1. system ip range: 10~49
    2. Instance ip range: 100~249
1. `Making the LINSTOR controller service highly available`

    ip: 2 # cloudstack portal ip
1. `CloudStack Setup`

    1. service ip: 3
    2. secondary storage ip: 4
1. `#step10: Setting up CloudStack # Wait for CloudStack API to come online`: 耗时很长, 因为cloudstack-management启动很慢

ps:
1. 当Replica count配置为3时, 至少两个节点online, 否则drbd不会mount drbd设备
1. ha环境/usr/share/linstor-server/bin/linstor-database export-db会报错, 因为停止linstor-controller后, /var/lib/linstor会被卸载, 而db刚好在这上面. 解决方法: 直接拷贝db即可

### 架构
每个节点都有linstor-satellite, 在设置高可用性LINSTOR集群时, 但只能有一个linstor-controller处于活动状态.
linstor-gui/linstor均使用linstor-controller提供的LINSTOR REST API

### cmds
```bash
# linstor [--machine-readable] node list # 节点状态 `--machine-readable`:json输出
# linstor volume list # 卷列表, 包含所有节点
# linstor volume-definition list # 当前节点volume-definition列表
# linstor node info # 列出集群中卫星节点上支持的存储提供程序和存储层, 将显示两个表. 第一个表将显示可用的LINSTOR存储提供程序后端, 第二个表将显示可用的LINSTOR存储层
# linstor storage-pool list # 支持的存储池
# linstor resource list # 显示资源, resource是[(volumes, node),...]的集合, volume间用VolNr区分. resource在各节点上使用相关的port和ResourceName
# linstor node list-properties test1 # 查看node属性
# linstor resource-definition  list # 查看创建的资源
# linstor resource-definition list-properties linstor_db # 查看resource-definition属性
# drbd-reactorctl status linstor_db # 查看resource-definition状态
# linstor node interface list test # 查看node上的interfaces
# linstor controller set-log-level --global TRACE # 修改所有节点log level为TRACE
```

Resource, Resource Definition, and Resource Group关系:
1. Resource Definition是创建的资源
1. Resource是Resource Definition涉及的子资源, 是(volumes, node)的集合, 这里的volumes包含多个卷时即可理解为传统意义上的卷组
1. Resource Group是创建资源时依赖的属性

DRBD Consistency Groups: 一个资源中的多个卷就是一个一致性组

### res
- `/etc/drbd.d/linstor-resources.res` -> `/var/lib/linstor.d`
- mysql账号: cloud/cloud
- linstor配置: /etc/linstor/linstor.toml
- metrics: http://192.168.122.2:3370/metrics