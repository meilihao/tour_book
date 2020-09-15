# targetcli
参考:
- [Managing storage devices#Getting started with iSCSI](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html/managing_storage_devices/getting-started-with-iscsi_managing-storage-devices)


```bash
# yum -y install targetd targetcli-fb

# apt install targetcli-fb # ubuntu 18.04
# apt install targetcli # [ubuntu 16.04](https://packages.ubuntu.com/search?suite=xenial&section=all&arch=any&keywords=targetcli&searchon=contents), 不推荐
```

> targetcli的官方git repo已不再维护, 推荐使用**[targetcli-fb](https://github.com/open-iscsi/targetcli-fb)**

targetcli 是用于管理 iSCSI 服务端存储资源的专用配置命令将 iSCSI 共享资源的配置内容抽象成“目录”的形式.

targetcli概念:
- backstore(后端存储)：后端真正的存储设备（实物）
- target(目标)：共享磁盘名（虚拟）
        target共享名的要求：iqn规范
            iqn规范 = iqn.yyyy-mm.主机域名反写:任意字串,  例: iqn.2018-02.com.example:data
- lun(逻辑单元)：Logic Unit Number ，绑定、关联存储设备

在执行 targetcli 命令后就能看到交互式的配置界面了, 利用 ls 查看目录参数的结构,使用 cd 切换到不同的目录中等.

**iSCSI target 名称是由系统自动生成的,这是一串用于描述共享资源的唯一字符串**.

iSCSI 协议是通过客户端名称进行验证的, 而该名称也是 iSCSI 客户端的唯一标识,而且必须与服务端配置文件中访问控制列表中的信息一致,否则客户端在尝试访问存储共享设备时,系统会弹出验证失败的保存信息. 因此用户在访问存储共享资源时不需要输入密码,只要 iSCSI 客户端的名称与服务端中设置的访问控制列表中某一名称条目一致即可.
iSCSI 协议是通过客户端的名称来进行验证,

acls 参数目录用于存放允许访问 iSCSI 服务端共享存储资源的客户端名称. 推荐在刚刚系统生成的 iSCSI target 后面追加上类似于:client 的参数,这样既能保证客户端的名称具有唯一性,又非常便于管理和阅读. "all"表示所有客户端都可以访问.

在 portals 参数目录中写上服务器的 IP 地址, 以便对外提供服务.

targetcli(服务端)使用步骤:
1. `/backstores/block> create disk0 /dev/md0` 创建磁盘映射,  `disk0`是后端存储名, `/dev/md0`是后端存储磁盘路径

    ```
    /backstores/block> create disk0 /dev/md0
    Created block storage object disk0 using /dev/md0.
    /backstores/block> ls
    o- block ...................................................................................................... [Storage Objects: 1]
      o- disk0 ............................................................................. [/dev/md0 (8.0GiB) write-thru deactivated]
        o- alua ....................................................................................................... [ALUA Groups: 1]
          o- default_tg_pt_gp ........................................................................... [ALUA state: Active/optimized]
    ```
1. `/iscsi> create` 创建iscsi target, 名称自动生成, 比如这里的`iqn.2003-01.org.linux-iscsi.linuxprobe.x8664:sn.d497c356ad80`

    ```
    /iscsi> create
    Created target iqn.2003-01.org.linux-iscsi.linuxprobe.x8664:sn.d497c356ad80.
    Created TPG 1.
    Global pref auto_add_default_portal=true
    Created default portal listening on all IPs (0.0.0.0), port 3260.
    /iscsi> ls
    o- iscsi .............................................................................................................. [Targets: 1]
      o- iqn.2003-01.org.linux-iscsi.linuxprobe.x8664:sn.d497c356ad80 ........................................................ [TPGs: 1]
        o- tpg1 ................................................................................................. [no-gen-acls, no-auth]
          o- acls ............................................................................................................ [ACLs: 0]
          o- luns ............................................................................................................ [LUNs: 0]
          o- portals ...................................................................................................... [Portals: 1]
            o- 0.0.0.0:3260 ....................................................................................................... [OK]
    ```
1. 创建lun
    ```
    /iscsi> cd iqn.2003-01.org.linux-iscsi.linuxprobe.x8664:sn.d497c356ad80/
    /iscsi/iqn.20....d497c356ad80> ls
    o- iqn.2003-01.org.linux-iscsi.linuxprobe.x8664:sn.d497c356ad80 .... [TPGs: 1]
        o- tpg1 ............................................. [no-gen-acls, no-auth] # Target Portal Group
            o- acls ........................................................ [ACLs: 0]
            o- luns ........................................................ [LUNs: 0]
            o- portals .................................................. [Portals: 0]
    /iscsi/iqn.20....d497c356ad80> cd tpg1/luns
    /iscsi/iqn.20...d80/tpg1/luns> create /backstores/block/disk0 # 创建lun，关联共享名和后端存储设备
    Created LUN 0.
    /iscsi/iqn.20...d80/tpg1/luns> ls
    o- luns .................................................................................................................. [LUNs: 1]
      o- lun0 ............................................................................. [block/disk0 (/dev/nbd0) (default_tg_pt_gp)]
    ```
1. 创建ACL
    ```
    /iscsi/iqn.20...d80/tpg1/acls> create iqn.2003-01.org.linux-iscsi.linuxprobe.x8664:sn.d497c356ad80:client # 指定允许访问的clients的iqn
    Created Node ACL for iqn.2003-01.org.linux-iscsi.linuxprobe.x8664:sn.d497c356ad80:client
    Created mapped LUN 0.
    /iscsi/iqn.20...d80/tpg1/acls> ls
    o- acls .................................................................................................................. [ACLs: 1]
      o- iqn.2003-01.org.linux-iscsi.linuxprobe.x8664:sn.d497c356ad80:client .......................................... [Mapped LUNs: 1]
        o- mapped_lun0 ......................................................................................... [lun0 block/disk0 (rw)]
    ```
1. 配置服务端口
    ```
    /iscsi/iqn.20...c356ad80/tpg1> cd portals
    /iscsi/iqn.20.../tpg1/portals> create 192.168.10.10  ip_port=3261
    Created network portal 192.168.10.10:3261.
    /iscsi/iqn.20.../tpg1/portals> ls
    o- portals ............................................................................................................ [Portals: 2]
      o- 0.0.0.0:3260 ............................................................................................................. [OK]
      o- 192.168.10.10:3261 ....................................................................................................... [OK]
    ```
1. 输入 exit 命令来退出配置, 重启 iSCSI 服务端程序`systemctl restart targetd`

## iscsi客户端(initiator)
```bash
# yum install iscsi-initiator-utils -y
# apt install open-iscsi -y
```

> [open-iscsi's git repo](https://github.com/open-iscsi/open-iscsi)

```bash
# 查询iSCSI 客户端的唯一标识iqn
# vim /etc/iscsi/initiatorname.iscsi
InitiatorName=iqn.2003-01.org.linux-iscsi.linuxprobe.x8664:sn.d497c356ad80:client # 编辑 iSCSI 客户端中的 initiator 名称文件,写入服务端的访问控制列表名称
# systemctl restart iscsid
# systemctl enable iscsid
```

iSCSI 客户端访问并使用共享存储资源的思路: 先发现,再登录,最后挂载并使用.

iscsiadm 是用于管理、查询、插入、更新或删除 iSCSI数据库配置文件的命令行工具,用户需要先使用这个工具扫描发现远程 iSCSI 服务端,然后查看找到的服务端上有哪些可用的共享存储资源.

由于 udev 服务是按照系统识别硬盘设备的顺序来命名硬盘设备的,当客户端主机同时使用多个远程存储资源时,如果下一次识别远程设备的顺序发生了变化,则客户端挂载目录中的文件也将随之混乱. 为了防止发生这样的问题,我们应该在/etc/fstab 配置文件中使用设备的 UUID 唯一标识符进行挂载,这样,不论远程设备资源的识别顺序再怎么变化,系统也能正确找到设备所对应的目录.

注意: 由于/dev/sdb是一块网络存储设备，而iSCSI协议是基于TCP/IP网络传输数据的，因此必须在/etc/fstab配置文件中添加上_netdev参数，表示当系统联网后再进行挂载操作，以免系统开机时间过长或开机失败.

iscsiadm:
- -m xxx : 动作: discovery,扫描并发现可用的存储资源;node, 为将客户端所在主机作为一台节点服务器
- -t st : 为执行扫描操作的类型,
- -p 192.168.10.10 : 指定 iSCSI 服务端的IP 地址
- -T iqn.2003-01.org.linux-iscsi.linux.x8664:sn.d497c356ad80 : 指定要使用的存储资源
- --login 或-l : 进行登录验证
- -u : 卸载iscsi设备

### example
```bash
# iscsiadm -m discovery -t st -p 192.168.10.10
# iscsiadm -m node -T iqn.2003-01.org.linux-iscsi.linux.x8664:sn.d497c356ad80 -p 192.168.10.10 --login # 在 iSCSI 客户端成功登录之后,会在客户端主机上多出一块名为`/dev/sd${xxx}` 的设备文件. `-T`表示要挂载的盘. 如果target使用了多张网卡时会存在多路径问题, 挂载磁盘数=target提供的磁盘数*路径数
# iscsiadm -m session -P 3 # 获取挂载信息, `-P`, 信息的详细level, 越大越详细.
# mkfs.xfs /dev/sdb
# mkdir /iscsi
# mount /dev/sdb /iscsi
# blkid | grep /dev/sdb
# vim /etc/fstab
UUID=eb9cbf2f-fce8-413a-b770-8b0f243e8ad6 /iscsi xfs defaults,_netdev 0 0 # 由于iscsi 磁盘是一块网络存储设备,而 iSCSI 协议是基于TCP/IP 网络传输数据的, 因此必须在/etc/fstab 配置文件中添加上_netdev 参数,表示当系统联网后再进行挂载操作,以免系统开机时间过长或开机失败.
# umount /iscsi   # 如果磁盘正在挂载使用，建议先卸载再登出
# iscsiadm -m node -T iqn.2003-01.org.linux-iscsi.linux.x8664:sn.d497c356ad80 -u # 登出
```

## FAQ
### 查找iSCSI client挂载生成的盘符
1. 找出所有iscsi盘: `lsblk -SJo TRAN,NAME`, 将tran是iscsi的所有盘找出, 假设这里仅有一块sdo
1. 找到对应的sgN: `ll /sys/block/sdo/device/scsi_generic`或`sg_map -i`
1. 找到关联的iqn号: `sg_inq -p 0x83 /dev/sgN|grep iqn`与iscsi挂载时所用的iqn做匹配即可

### 查看target iblock的lun序号
在target端查找磁盘的T10 VPD Unit Serial Number(即scsi serial number, LUN序列号)

    ```bash
    # cat /sys/kernel/config/target/core/iblock_xxx/${iblock_name}/wwwn/vpd_unit_serial # iblock_name是targetcli's backstores/iblock中对于的名称
    T10 VPD Unit Serial Number: xxx # xxx为lun序列号, 创建iblock时自行生成
    ```