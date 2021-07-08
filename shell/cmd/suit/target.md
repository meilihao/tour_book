# targetcli
参考:
- [Managing storage devices#Getting started with iSCSI](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html/managing_storage_devices/getting-started-with-iscsi_managing-storage-devices)
- [FC 和iSCSI的使用差异](https://www.huaweicloud.com/articles/57a2047190a7768914f2a0040d3da26f.html)


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
- tpg : Target Portal Group, 实际是为了方便管理target下众多LUN,而添加的一个虚拟层

在执行 targetcli 命令后就能看到交互式的配置界面了, 利用 ls 查看目录参数的结构,使用 cd 切换到不同的目录中等.

**iSCSI target 名称是由系统自动生成的,这是一串用于描述共享资源的唯一字符串**.

iSCSI 协议是通过客户端名称进行验证的, 而该名称也是 iSCSI 客户端的唯一标识,而且必须与服务端配置文件中访问控制列表中的信息一致,否则客户端在尝试访问存储共享设备时,系统会弹出验证失败的保存信息. 因此用户在访问存储共享资源时不需要输入密码,只要 iSCSI 客户端的名称与服务端中设置的访问控制列表中某一名称条目一致即可.
iSCSI 协议是通过客户端的名称来进行验证,

acls 参数目录用于存放允许访问 iSCSI 服务端共享存储资源的客户端名称. 推荐在刚刚系统生成的 iSCSI target 后面追加上类似于:client 的参数,这样既能保证客户端的名称具有唯一性,又非常便于管理和阅读. "all"表示所有客户端都可以访问.

在 portals 参数目录中写上服务器的 IP 地址, 以便对外提供服务.

## iscsi
参考:
- [targetcli配置iSCSI](https://www.cnblogs.com/luxiaodai/p/9851214.html)

targetcli(服务端)使用步骤:
1. `/backstores/block> create disk0 /dev/md0` 创建磁盘映射,  `disk0`是后端存储名, `/dev/md0`是后端存储磁盘路径, 删除backstores对象时target中相应的lun, mapped_lun也会被删除

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
1. 创建lun(当luns存在多个lun时, iscsiadm login可一次性挂载这些luns, 同理iscsiadm logout也会一次性卸载这些盘)
    ```
    /iscsi> cd iqn.2003-01.org.linux-iscsi.linuxprobe.x8664:sn.d497c356ad80/
    /iscsi/iqn.20....d497c356ad80> ls
    o- iqn.2003-01.org.linux-iscsi.linuxprobe.x8664:sn.d497c356ad80 .... [TPGs: 1]
        o- tpg1 ............................................. [no-gen-acls, no-auth] # Target Portal Group
            o- acls ........................................................ [ACLs: 0]
            o- luns ........................................................ [LUNs: 0]
            o- portals .................................................. [Portals: 0]
    /iscsi/iqn.20....d497c356ad80> cd tpg1/luns
    /iscsi/iqn.20...d80/tpg1/luns> create /backstores/block/disk0 # 创建lun，关联共享名和后端存储设备, lun id must in [0,255]
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

> 在交互模式下默认创建完配置exit退出时会主动将配置保存到配置文件`/etc/rtslib-fb-target/saveconfig.json`或`/target/saveconfig.json`中，重启后生效, 该配置路径可通过`targetcli saveconfig`修改.

ps: 执行saveconfig也会输出配置文件路径相关的信息.

## fc
targetcli(服务端)使用步骤:
1. 在/backstores/block下创建磁盘映射disk0, 同iscsi
1. 创建光纤target

    ```
    /> qla2xxx> create 21:00:00:1b:32:81:6e:f1    //本机的wwwn，获取方式见FAQ
    ```
1. 创建lun

    ```
    /qla2xxx/21:00:00:1b:32:81:6e:f1/luns> create  /backstores/block/disk0
    ```
1. 创建acl

    ```
    /qla2xxx/21:00:00:1b:32:81:6e:f1/acls> create 21:00:00:1b:32:98:7d:1b   //将Lun映射给192.168.1.88对应的wwwn
    ```
1. 保存配置

    ```
    /> saveconfig # 必须在顶层执行
    ```

## targetcli CHAP（质询握手身份验证协议）
参考:
- [Applying CHAP Authentication to iSCSI ACLs](https://docs.softnas.com/display/SD/Applying+CHAP+Authentication+to+iSCSI+ACLs)

配置targetcli CHAP认证, 分为全局配置和局部配置:
- /iscsi 下为全局配置

    ```bash
    cd /iscsi
    set discovery_auth enable=1 # 启用发现 CHAP 验证

    # 用于单向验证, 客户端登入验证用户和密码
    set discovery_auth userid=InUser
    set discovery_auth password=InPassword

    # 用于双向验证, 设置反向验证用户名和密码，当只设置单向验证时，请取消下面的设置
    set discovery_auth mutual_userid=OutUse
    set discovery_auth mutual_password=OutPassword
    ```
- 在 iscsi/iqn.2019-10.cc.pipci.iscsi:debian.tgt1/tpg1/ 下为单个Target的配置，配置只对单个IQN生效为局部配置

    > 其实Normal Authentication可通过TPG节点和ACL节点中的set auth命令配置, 这取决于 TPG 的 generate_node_acls 属性设置.
如果generate_node_acls是1，则将使用tpg范围的设置, 如果generate_node_acls为0，则将使用用户创建的ACLs的设置. **在targetcli 2.1.fb48(deepin 20.2-kernel 5.12)/2.1.51(ubuntu 20.04-kernel 5.4)上实践发现normal auth仅设置在acl下才有效.**

    ```bash
    cd iscsi/
    set discovery_auth enable=0 # 关闭服务端的发现验证
     
    cd /iscsi/iqn.2018-07.com.holoem.iscsi:target/tpg1/
     
    set attribute authentication=1 # 开启验证
    
    # 这两项也可以不设置
    set attribute generate_node_acls=1
    set attribute cache_dynamic_acls=1 # 设置强制使用 TPG 的身份验证
    
    # 下面两项配置位置取决于generate_node_acls: 1在tpg, 0在acl
    set auth userid=InAuthUser password=InAuthPassword # 启用客户端登入验证,并设置登入用户名和密码, kv的value上不能用`"`包裹否则不能生效. 取消时必须逐个取消(比如``set auth userid=``), 一起取消会发生错乱

    # 设置反向验证用户名和密码，当只设置单向验证时，请取消该设置
    set auth mutual_userid=OutAuthUser mutual_password=OutAuthPassword
     
    cd /
    saveconfig # 保存配置
    exit
     
    systemctl restart target # 重启服务
    ```

全局配置下只能设置发现(discovery)认证，局部配置只能设置登录(normal)认证，其中每种认证又分为单向认证和双向认证, 无论那种认证都是在target端配置的:
- 单向认证是target server认证initiator, 即initiator端在发现target端的时候，要提供正确的认证才能发现在target端的iSCSI服务
- 双向认证是target server和inititor互相认证, 即在单向认证的基础上，target端需要正确设置initiator端设置的认证才能被initiator端发现

> 设置双向认证必须建立在单向认证的基础上，因为在initiator登录的时候要先进行单项认证.

具体配置参考[这里](https://www.cnblogs.com/pipci/p/11622014.html)和[认证](https://wiki.archlinux.org/index.php/ISCSI_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87)/LIO_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87))

配置后需要saveconfig保存targetcli配置以及`systemctl restart target.service`.


全局时的iscsi initiator配置:
```bash
vim /etc/iscsi/iscsid.conf # 编辑客户端 iscsiadm 守护程序配置文件
# 设置启用发现验证
discovery.sendtargets.auth.authmethod = CHAP
 
# 设置发现服务端时需要的用户和密码
discovery.sendtargets.auth.username = InUser
discovery.sendtargets.auth.password = InPassword
 
# 设置当服务端访向客户端验证时的用户和密码
discovery.sendtargets.auth.username_in = OutUser
discovery.sendtargets.auth.password_in = OutPassword
 
systemctl restart iscsi
systemctl restart iscsid # 重启 iscsi和iscsid 服务
```

在 tpg 下置客户端 iscsiadm CHAP 双向认证:
```bash
vim /etc/iscsi/iscsid.conf # 编辑客户端 iscsiadm 守护程序配置文件
 
# 设置启用登陆验证
node.session.auth.authmethod = CHAP
 
# 设置登入服务端 targetcli 的用户名和密码
node.session.auth.username = InAuthUser
node.session.auth.password = InAuthPassword
 
# 设置当服务端访向客户端验证时的用户和密码
node.session.auth.username_in = OutAuthUser
node.session.auth.password_in = OutAuthPassword
 
systemctl restart iscsi
systemctl restart iscsid # 重启 iscsi和iscsid 服务
```

> 在targetcli 2.1.fb48(deepin 20.2-kernel 5.12)/2.1.51(ubuntu 20.04-kernel 5.4)上实践发现配置/etc/iscsi/iscsid.conf的normal chap后, `sudo iscsiadm -m node -T iqn.2003-01.org.linux-iscsi.chen-aliyun.x8664:sn.186dcf05e59b |grep auth`仍为空. 但针对某个target设置normal认证(见iscsiadm example)是生效的.

## targetcli cmd模式
```bash
# targetcli /backstores/block create name=disk1 dev=/dev/nbd1 [wwn=bb3f4d39-881a-4932-9e3e-9537ba9be9f4] # wwn会保存到/sys/kernel/config/target/core/iblock_1/disk1/wwn/vpd_unit_serial中
# targetcli /backstores/block help create
```

## targetcli 部分参数/属性
- demo_mode_write_protect=0 : 设置为 demo 模式, 即无需要配置 ACL 权限控制列表即可连接
- authentication=0 : 关闭密码认证(chap信息未删除不会影响密码认证的关闭)
- generate_node_acls=1 # 此项为将 Initiator 自动添加进 ACL 认证列表，设置generate_node_acls时必须与authentication同时配置, 它可控制Normal Authentication设置在TPG节点还是ACL节点
- cache_dynamic_acls=1 # 此项为记录 ACL 信息, 方便 SCSI 保持连接，设置成 demo 模式后，默认即为1

## targetcli backstores
backstores分类:
- block/iblock(旧版使用) : 通常能提供最好的性能，可以使用其他任何类型的磁盘设备

    /backstores> iblock/ create name=block_backend dev=/dev/sdb
- fileio : 不要使用 buffered FILEIO，默认是non-buffered 模式

    /backstores> fileio/ create name=file_backend file_or_dev=/usr/src/fileio size=2G

    如果新建的FILEIO 中，参数 buffered =True，就可以使用buffer cache ，将明显提高其有效性能
    同时伴随的风险是一系列数据的整体风险：如果系统崩溃，一个 unflushed buffer cache将导致整个后
    备存储不能挽回的损坏.
- [pscsi(parallel SCSI)](https://en.wikipedia.org/wiki/Parallel_SCSI): 物理scsi设备，不推荐使用, 建议使用 block 代替

    /backstores> pscsi/ create name=pscsi_backend dev=/dev/sr0
- ramdisk : RAM 硬盘后备存储

    /backstores> ramdisk/ create name=rd_backend size=1GB

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
# iscsiadm -m discoverydb -t st -p 192.168.10.10 -o show # 输出discovery信息(含认证)
# iscsiadm -m node -T iqn.2003-01.org.linux-iscsi.linux.x8664:sn.d497c356ad80 -p 192.168.10.10 --login # 此时是禁用CHAP的情况 ,在 iSCSI 客户端成功登录之后,会在客户端主机上多出一块名为`/dev/sd${xxx}` 的设备文件. `-T`表示要挂载的盘. 如果target使用了多张网卡时会存在多路径问题, 挂载磁盘数=target提供的磁盘数*路径数
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

针对某个target设置chap:
```bash
sudo iscsiadm -m node -T iqn.2003-01.org.linux-iscsi.fyhdesktop29.x8664:sn.0d690d398ec5 |grep auth # 查看是否设置了auth
sudo iscsiadm -m node -T iqn.2003-01.org.linux-iscsi.fyhdesktop29.x8664:sn.0d690d398ec5 -o update --name=node.session.auth.authmethod --value=CHAP
sudo iscsiadm -m node -T iqn.2003-01.org.linux-iscsi.fyhdesktop29.x8664:sn.0d690d398ec5 -o update --name=node.session.auth.username --value=user
sudo iscsiadm -m node -T iqn.2003-01.org.linux-iscsi.fyhdesktop29.x8664:sn.0d690d398ec5 -o update --name=node.session.auth.password --value=password
sudo cat /etc/iscsi/nodes/iqn.2003-01.org.linux-iscsi.fyhdesktop29.x8664\:sn.0d690d398ec5/127.0.0.1\,3260\,1/default # 配置位置, `127.0.0.1`是target server ip, `1`未知, 但`iscsiadm -m discovery -t st -p 127.0.0.1`时都能找到
```

## FAQ
### 查找iSCSI initiator挂载生成的盘符
方法1, **推荐**:
1. 在target端查找磁盘的T10 VPD Unit Serial Number(即scsi serial number, LUN序列号)

    ```bash
    # cat /sys/kernel/config/target/core/iblock_xxx/${iblock_name}/wwwn/vpd_unit_serial # iblock_name是targetcli's backstores/iblock中对于的名称
    T10 VPD Unit Serial Number: xxx # xxx为lun序列号, 创建block时自行生成
    ```
1. 在initiator端执行`ll /dev/disk/by-id |grep xxx`即可

方法2:
1. 找出所有iscsi盘: `lsblk -SJo TRAN,NAME`, 将tran是iscsi的所有盘找出, 假设这里仅有一块sdo
1. 找到对应的sgN: `ll /sys/block/sdo/device/scsi_generic`或`sg_map -i`
1. 找到关联的iqn号: `sg_inq -p 0x83 /dev/sgN|grep iqn`与iscsi挂载时所用的iqn做匹配即可或通过`sg_inq -p 0x83 /dev/sgN|grep "vendor specific"`与target端的`T10 VPD Unit Serial Number`做匹配

> 方法2仅测试过target端是一个target提供一个lun的情况, 而一个target提供若干lun的情况未测试.

### 查找FC initiator挂载生成的盘符
方法1,**推荐**:
1. 在target端查找磁盘的T10 VPD Unit Serial Number(即scsi serial number, LUN序列号)

    ```bash
    # cat /sys/kernel/config/target/core/iblock_xxx/${iblock_name}/wwwn/vpd_unit_serial # iblock_name是targetcli's backstores/iblock中对于的名称
    T10 VPD Unit Serial Number: xxx # xxx为lun序列号, 创建block时自行生成
    ```
1. 在initiator端执行`ll /dev/disk/by-id |grep xxx`即可

方法2:
1. 找出所有fc盘: `lsblk -SJo TRAN,NAME`, 将tran是fc的所有盘找出, 假设这里仅有一块sdo
1. 找到对应的sgN: `ll /sys/block/sdo/device/scsi_generic`或`sg_map -i`
1. 找到关联的naa: `sg_inq -p 0x83 /dev/sgN|grep naa`与target的wwpn做匹配, 此时只能确定该lun由指定target提供而不能一一对应. 但通过`sg_inq -p 0x83 /dev/sgN|grep "vendor specific"(=`sg_inq -p 0x80 /dev/sgN`)`与target端的`T10 VPD Unit Serial Number`做匹配即可一一对应.

### 不设置acl
在ACL配置目录执行 set attribute generate_node_acls=0使用自定义的acl实现访问控制，则需要设置访问权限控制列表acl（默认就是这种），acl参数目录用于存放能够访问target端共享存储资源的initiator的iqn. 在客户端访问时，只要iscsi客户端的iqn名称与服务端设置的访问控制列表中的iqn名称一致即可访问. 如果不想使用ACL可以在ACL配置目录执行 set attribute generate_node_acls=1使用自动生成acl节点，这样不添加initiator的iqn也允许initiator访问.
```
/iscsi/iqn.20...ian.tgt1/tpg1> set attribute generate_node_acls=1 # 配置成自动生成acl节点
```

一旦配置成自动生成acl节点，当initiator认证成功后，再配置成自定义的acl实现访问控制是无效的 只有重启系统后恢复正常，我感觉这个是因为有认证记忆的功能.

### iscsiadm -m node xxx 无法login, 报"initiator reported error ( 24 - ..."
开启了CHAP认证, 禁用即可: `.../tpg1> set attribute authentication=0` 或 iscsiadm添加chap信息(见iscsiadm example)

此时target端是报: `kernel: Initiator is requesting CSG: 1, has not been successfully authenticated, and the Target is enforcing iSCSI Authentication, login failed.`

ps: `/iscsi`设置`set discovery_auth enable=0`, 但tpgX设置`set attribute authentication=1`时, tpgX还是开启了chap认证.

### 获取光纤信息
和以太网卡的MAC地址一样，HBA上也有独一无二的标识：WWN（World Wide Name）, FC HBA上的WWN有两种：
1. Node WWN（WWNN）：每块HBA有其独有的Node WWN
2. Port WWN（WWPN）：每块HBA卡上每个port有其独一无二的Port WWN

由于通信是通过port进行的，因此多数情况下需要使用WWPN而不是WWNN. WWN的长度为8bytes，用16进制表示并用冒号分隔，例如：`50:06:04:81:D6:F3:45:42:23`. 通常说的光纤WWN均指WWPN.

```bash
# lspci  | grep -i fibre # 查看fc HBA卡, 通常一块光纤卡有两个光纤口
# cat /sys/class/fc_host/host<N>/node_name # 查看fc HBA卡WWNN信息
# cat /sys/class/fc_host/host<N>/port_name # 查看fc HBA卡WWPN信息
# cat /sys/class/fc_host/host<N>/port_state # 查看fc 插口的状态: Online表示插有光纤且与对端(光纤卡或光纤交换机)联通
Online
# cat /sys/class/fc_host/host<N>/port_type # 查看fc 插口的连接类型: LPort是与其他HBA卡相连; NPort是与光纤交换机相连
# cat /sys/class/fc_host/host<N>/supported_speeds # 查看port支持的速率
# systool -v -c fc_host # 获取详细的光纤卡信息, from `apt install sysfsutils`
```

### 光纤initiator发现的方法
1. `echo 1 > /sys/class/fc_host/host<N>/issue_lip`, **推荐** # 此时会通过issue_lip重置HBA链路，重新扫描整个链路并配置SCSI target. 该操作是一种异步操作类型，具体完成时间需要参考system log.
1. `echo "- - -" |tee -a /sys/class/scsi_host/*/scan` # `- - -`分别代表通道，SCSI目标ID和LUN, 此时破折号充当通配符，表示"重新扫描所有内容"

> 有时明明fc target配置正确但fc client还是不能扫出新盘: 有坏的fc链路占用了相同的盘符(比如sdc), 导致不能扫出. 解决方法:1. `rescan-scsi-bus.sh -r`; 2. `reboot`

### Could not create Qla2xxxFabricModule in configFS | Could not create Target in configFS | 看不到FC fabric
`modprobe tcm_qla2xxx`

### 光纤initiator无法发现target
qla2xxx.ko模式可能不对.

qla2xxx.ko支持target模式和initiator模式, 在存储服务器上必须根据target模式加载, 而initiator端需要initiator模式，可以参考下面的命令重新加载：
```bash
# cat /sys/module/qla2xxx/parameters/qlini_mode # 查看当前qla2xxx.ko的模式
# systool -m qla2xxx -v |grep qlini_mode # 也可通过systool查看当前qla2xxx.ko的模式
# modprobe -r qla2xxx
# modprobe qla2xxx qlini_mode="disabled" # 只支持target模式, 除非重新加载qla2xxx驱动
# modprobe qla2xxx qlini_mode="enabled"  # 只支持initiator模式, 除非重新加载qla2xxx驱动
```

> 也可通过/etc/modprobe.d/qla2xxx.conf指定qla2xxx驱动参数, 比如`options qla2xxx qlini_mode="enabled"`.

> 其实qlini_mode默认是"exclusive"模式: 默认支持initiator模式, 通过操作target驱动提供的configfs接口, 可切换到target模式, 还可以再切回initiator模式.

### rm -rf "/sys/kernel/config/target/core/iblock_0", 删除失败
target configfs与普通的文件系统有一定的差异导致删除失败.

解决方法:
```python
os.rmdir("/sys/kernel/config/target/core/iblock_0") # 这样即可
```

### 查看/设置属性
`get/set attribute`只在特定targetcli path下有效(本质是`/sys/kernel/config/target/*`下的相应目录下存在attrib文件夹), 比如(可能不仅这些):

```
/backstores/block/disk0> get attribute

/iscsi/iqn.20...d80/tpg1> get attribute
/iscsi/iqn.20...d80/tpg1> set attribute demo_mode_write_protect=0

/qla2xxx/21:00:00:1b:32:81:6e:f1> get attribute
/qla2xxx/21:00:00:1b:32:81:6e:f1> set attribute demo_mode_write_protect=0
```

### arm64 github.com/Datera/targetcli 3.0-pre4 iscsi无法添加portals
```bash
# vim /usr/local/lib/python2.7/dist-packages/rtslib/utils.py
def list_eth_names(max_eth=1024):
...
    if os.uname()[4].endswitch("_64") # arm64的输出是aarch64, 去掉"_64"中的"_"即可
```

### `Abort command issued nexus` error messages in /var/log/syslog
参考: [“Abort command issued nexus” error messages in /var/log/messages file](https://www.thegeekdiary.com/abort-command-issued-nexus-error-messages-in-var-log-messages-file/)

现象：fio测试multipath时出现报错, 用`multipath -F`解除多路径进行单盘测试. 在fc client端测试(`dd if=/dev/sdc of=/dev/null bs=1M count=100`)过程中发现`cat /sys/class/fc_remote_ports/*/port_state`fc盘的链路有时一直时Online有时Blocked(过会会恢复成Online), 且syslog报错: `Abort command issued nexus=8:2:0 -- 2002`, 根据参考的资料的结论是: 这些日志条目的存在表明 fc链路 I/O 性能严重下降, 从而导致一系列问题.

在target端进行同样的dd测试发现原盘同样很慢, 因此断定是原盘的问题.

### 获取fc链路的client port_name
1. 通过`/sys/class/fc_host/xxx`确认名称
1. 读取`/sys/class/scsi_host/xxx/sns_table`, 其中包含就是与该host通信的client port_name

### CHAP user or password not set for Initiator ACL
- case1:

    chap应设置在target的acl下, 比如在`/iscsi/iqn.2003-01.org.linux-iscsi.fyhdesktop29.x8664:sn.0d690d398ec5/tpg1/acls/iqn.1993-08.org.debian:01:7ed7bee79b74`下设置`set auth userid=user password=password`

- case2:

    target server已取消CHAP, 但initiator端仍未取消, 因此在Initiator端注释chap配置即可

# tgtadm
参考:
- [github.com/longhorn/go-iscsi-helper](https://github.com/longhorn/go-iscsi-helper/blob/master/iscsi/target.go)

安装tgt:
```bash
apt-get install tgt
ll /etc/tgt # 配置文件位置
```

常用命令:
```bash
tgtadm --mode target --op show
```