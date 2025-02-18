# libvirt
ref:
- [libvirt 源码分析 - virsh](https://winddoing.github.io/post/dec26e6d.html)
- [virsh list所有vm state及其转换](https://docs.openeuler.org/zh/docs/20.03_LTS_SP3/docs/Virtualization/%E7%AE%A1%E7%90%86%E8%99%9A%E6%8B%9F%E6%9C%BA.html)
- [Domain XML format](https://avdv.github.io/libvirt/formatdomain.html)
- [虚拟化调试和优化指南](https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/7/html/virtualization_tuning_and_optimization_guide/sect-virtualization_tuning_optimization_guide-blockio-tuning)
- [如何安装 oVirt 管理程序](https://www.storagereview.com/zh-CN/news/how-to-install-ovirt-hypervisor)
- [Proxmox VE(Proxmox Virtual Environment)](https://pve.proxmox.com/wiki/Main_Page)
   支持集群
- [KubeSphere 虚拟化 KSV 安装体验](https://blog.csdn.net/networken/article/details/125009058)
- [kubevirt在360的探索之路（k8s接管虚拟化）](https://zyun.360.cn/blog/?p=691)
- [虚拟机迁移](https://kiosk007.top/post/%E8%99%9A%E6%8B%9F%E6%9C%BA%E6%80%8E%E4%B9%88%E5%81%9A%E5%88%B0%E7%83%AD%E8%BF%81%E7%A7%BB/)

   - [网络数据传输层](https://mp.weixin.qq.com/s?__biz=MzI3MDM0NjU3MA==&mid=2247486297&idx=1&sn=046c2b4c363668a4d22af38d81bef17e&chksm=ead33b7cdda4b26ac7cd2b62f76de0bcf160f59a37f0f171eacbb21f1c3ea72da11c03ec21d2&cur_album_id=2372160472910839809&scene=189#wechat_redirect)
- [fedora 40 设置虚拟化环境](https://docs.fedoraproject.org/en-US/quick-docs/virtualization-getting-started/)

目前使用最广泛的对kvm进行管理的工具和应用程序接口, 它也支持xen, vmware, virtualbox, hyper-v等平台虚拟化, 以及openvz, lxc等容器虚拟化.

libvirt对多种不同的Hypervisor的支持是通过一种基于驱动程序的架构来实现的. libvirt对不同的Hypervisor提供了不同的驱动： 对Xen有Xen的驱动， 对QEMU/KVM有QEMU驱动， 对VMware有VMware驱动. 在libvirt源代码中， 可以很容易找到qemu_driver.c、 xen_driver.c、 xenapi_driver.c、 VMware_driver.c、 vbox_driver.c这样的驱动源码.


sdk:
- libvirt.org/go/libvirt和libvirt.org/go/libvirtxml

   [libvirt.org/libvirt-go已被支持go module的libvirt.org/go/libvirt取代](https://libvirt.org/libvirt-go.html). sdk使用参考[libvirt-go中能够提供的虚机信息](https://blog.csdn.net/zhagzheguo/article/details/100050474)

安装: `sudo apt install qemu-system-x86 virt-manager libvirt-daemon libvirt-daemon-system virtinst libvirt-clients bridge-utils`

## 概念
- 节点（Node） 

   是一个物理机器, 上面可能运行着多个虚拟客户机. Hypervisor和Domain都运行在节点上.

   通常node上除了需要运行相应的Hypervisor以外， 还需要让libvirtd这个守护进程处于运行中的状态， 以便让客户端连接到libvirtd， 从而进行管理操作。 不过， 也并非所有的Hypervisor都需要运行libvirtd守护进程， 比如VMware ESX/ESXi就不需要在服务器端运行libvirtd， 依然可以通过libvirt客户端连接到VMware, 具体可参考[这里](http://libvirt.org/drvesx.html).

- Hypervisor也称虚拟机监控器（VMM）

   如KVM、 Xen、 VMware、 Hyper-V等， 是虚拟化中的一个底层软件层， 它可以虚拟化一个节点让其运行多个虚拟客户机（不同客户机可能有不同的配置和操作系统）。
- 域（Domain） 是在Hypervisor上运行的一个客户机操作系统实例

   域也被称为实例（instance， 如在亚马逊的AWS云计算服务中客户机就被称为实例） 、 客户机操作系统
（guest OS） 、 虚拟机（virtual machine） ， 它们都是指同一个概念.

## 功能
1. 域的管理

   包括对节点上的域的各个生命周期的管理， 如启动、 停止、 暂停、 保存、 恢复和动态迁移。 还包括对多种设备类型的热插拔操作， 包括磁盘、 网卡、 内存和
CPU。 当然不同的Hypervisor上对这些热插拔的支持程度有所不同.
1. 远程节点的管理

   只要物理节点上运行了libvirtd这个守护进程, 远程的管理程序就可以连接到该节点进程管理操作， 经过认证和授权之后， 所有的libvirt功能都可以被访
问和使用.

   libvirt支持多种网络远程传输类型， 如SSH、 TCP套接字、 Unix domainsocket、 TLS的加密传输等. 假设使用了最简单的SSH， 不需要额外的配置工作， 比如，在example.com节点上运行了libvirtd， 而且允许SSH访问， 在远程的某台管理机器上就可以用`virsh -c qemu+ssh://root@example.com/system`命令行来连接到example.com上， 从而管理其上的域.
1. 存储的管理

   任何运行了libvirtd守护进程的主机， 都可以通过libvirt来管理不同类型的存储， 如创建不同格式的客户机镜像（qcow2、 raw、 qde、 vmdk等） 、 挂载NFS共享存储系统、 查看现有的LVM卷组、 创建新的LVM卷组和逻辑卷、 对磁盘设备分区、 挂载iSCSI共享存储、 使用Ceph系统支持的RBD远程存储， 等等。 当然在libvirt中， 对存储的管理也是支持远程的.
1. 网络的管理

   任何运行了libvirtd守护进程的主机， 都可以通过libvirt来管理物理的和逻辑的网络接口。 包括列出现有的网络接口卡， 配置网络接口， 创建虚拟网络接口， 网络接口的桥接， VLAN管理， NAT网络设置， 为客户机分配虚拟网络接口， 等等.
1. 提供一个稳定、 可靠、 高效的应用程序接口， 以便可以完成前面的4个管理功能

   libvirt主要由3个部分组成:
   1. 应用程序编程接口库

      应用程序接口是为其他虚拟机管理工具（ 如virsh、virt-manager等） 提供虚拟机管理的程序库支持

   2. 一个守护进程（ libvirtd）

      libvirtd守护进程负责执行对节点上的域的管理工作， 在用各种工具对虚拟机进行管理时， 这个守护进程一定要处于运行状态中. 而且这个守护进程可以分为两种： 一种是root权限的libvirtd， 其权限较大， 可以完成所有支持的管理工作； 一种是普通用户权限的libvirtd， 只能完成比较受限的管理工作.

      ```bash
      apt install libvirt-daemon libvirt-daemon-system # libvirtd.service is in libvirt-daemon-system
      ```
   3. 一个默认命令行管理工具（virsh）

      virsh是libvirt项目中默认的对虚拟机管理的一个命令行工具

      ```bash
      apt install libvirt-clients # for virsh
      ```

## arch
libvirt API大致可划分为如下8个部分:
1. 连接Hypervisor相关的API： 以virConnect开头的一系列函数

   只有在与Hypervisor建立连接之后， 才能进行虚拟机管理操作， 所以连接Hypervisor的API是其他所有API使用的前提条件. 与Hypervisor建立的连接为其他API的执行提供了路径， 是其他虚拟化管理功能的基础.

   通过调用virConnectOpen函数可以建立一个连接， 其返回值是一个virConnectPtr对象， 该对象就代表到Hypervisor的一个连接； 如果连接出错， 则返回空值（NULL）. 而virConnectOpenReadOnly函数会建立一个只读的连接， 在该连接上可以使用一些查询的功能， 而不使用创建、 修改等功能. virConnectOpenAuth函数提供了根据认证建立的连接. virConnectGetCapabilities函数返回对Hypervisor和驱动的功能描述的XML格式的字符串. virConnectListDomains函数返回一列域标识符， 它们代表该Hypervisor上的活动域

   libvirt使用了在互联网应用中广泛使用的URI（Uniform Resource Identifier， 统一资源标识符） 来标识到某个Hypervisor的连接:
   - 本地uri: `driver[+transport]:///[path][?extral-param]`

      > driver是连接Hypervisor的驱动名称（如qemu、 xen、 xbox、 lxc等）, transport是选择该连接所使用的传输方式（可以为空， 也可以是"unix"） ， path是连接到服务器端上的某个路径， ?extral-param是可以额外添加的一些参数（如Unix domainsockect的路径）

      在libvirt中KVM使用QEMU驱动. QEMU驱动是一个多实例的驱动， 它提供了一个系统范围内的特权驱动（即`system`实例） 和一个用户相关的非特权驱动（即`session`实例）. 

      session实例是根据客户端的当前用户和用户组去服务器端寻找对应用户下的实例. 在建立session连接后， 可以查询和控制的域或其他资源都仅仅是在当前用户权限范围内的， 而不是整个节点上的全部域或其他资源.

      system实例需要系统特权账号`root`权限. 在建立system连接后， 由于它是具有最大权限的， 因此可以查询和控制整个节点范围内的域， 还可以管理该节点上特权用户才能管理的块设备、 PCI设备、 USB设备、 网络设备等系统资源.

      一般来说， 为了方便管理， 在公司内网范围内建立到system实例的连接进行管理的情况比较常见， 当然为了安全考虑, 赋予不同用户不同的权限就可以使用建立到.

      在libvirt中， 本地连接QEMU/KVM的几个URI示例如下：
      - `qemu:///session`： 连接到本地的session实例， 该连接仅能管理当前用户的虚拟化资源
      - `qemu+unix:///session`： 以Unix domain sockect的方式连接到本地的session实例， 该连接仅能管理当前用户的虚拟化资源
      - `qemu:///system`： 连接到本地的system实例， 该连接可以管理当前节点的所有特权用户可以管理的虚拟化资源
      - `qemu+unix:///system`： 以Unix domain sockect的方式连接到本地的system实例, 该连接可以管理当前节点的所有特权用户可以管理的虚拟化资源
   - 远程URI: `driver[+transport]://[user@][host][:port]/[path][?extral-param]`

      > transport表示传输方式， 其取值可以是ssh、 tcp、 libssh2等； user表示连接远程主机使用的用户名， host表示远程主机的主机名或IP地址， port表示连接远程主机的端口, 其余参数的意义与本地URI中介绍的完全一样.

      在libvirt中， 远程连接QEMU/KVM的URI示例如下：
      - `qemu+ssh://root@example.com/system?keyfile=/root/.ssh/example_key`： 通过ssh通道连接到远程节点的system实例，具有最大的权限来管理远程节点上的虚拟化资源. 建立该远程连接时， 需要经过ssh的用户名和密码验证或者基于密钥的验证.
      - `qemu+ssh://user@example.com/session`： 通过ssh通道连接到远程节点的使用user用户的session实例， 该连接仅能对user用户的虚拟化资源进行管理， 建立连接时同样需要经过ssh的验证。
      - `qemu://example.com/system`： 通过建立加密的TLS连接与远程节点的system实例相连接， 具有对该节点的特权管理权限。 在建立该远程连接时， 一般需要经过TLS x509安全协议的证书验证
      - `qemu+tcp://example.com/system`： 通过建立非加密的普通TCP连接与远程节点的system实例相连接， 具有对该节点的特权管理权限。 在建立该远程连接时， 一般需要经过SASL/Kerberos认证授权

      除了针对QEMU、 Xen、 LXC等真实Hypervisor的驱动之外， libvirt自身还提供了一个名叫“test”的傀儡Hypervisor及其驱动程序. test Hypervisor是在libvirt中仅仅用于测试和命令学习的目的， 因为在本地的和远程的Hypervisor都连接不上（或无权限连接）时， test这个Hypervisor却一直都会处于可用状态. 其用法: `virsh -c test:///default <cmd, 比如list>`
1. 域管理的API： 以virDomain开头的一系列函数

   虚拟机最基本的管理职能就是对各个节点上的域的管理， 故在libvirt API中实现了很多针对域管理的函数. 要管理域， 首先要获取virDomainPtr这个域对象， 然后才能对域进行操作.

   有很多种方式来获取域对象， 如virDomainPtrvirDomainLookupByID(virConnectPtr conn， int id)函数是根据域的id值到conn这个连接上去查找相应的域. 类似的， virDomainLookupByName、 virDomainLookupByUUID等函数分别是根据域的名称和UUID去查找相应的域. 在得到某个域的对象后， 就可以进行很多操作， 可以查询域的信息（如virDomainGetHostname、 virDomainGetInfo、virDomainGetVcpus、 virDomainGetVcpusFlags、 virDomainGetCPUStats等） ， 也可以控制域的生命周期（如virDomainCreate、 virDomainSuspend、 virDomainResume、virDomainDestroy、 virDomainMigrate等）
1. 节点管理的API： 以virNode开头的一系列函数

   域运行在物理节点之上， libvirt也提供了对节点进行信息查询和控制的功能. 节点管理的多数函数都需要使用一个连接Hypervisor的对象作为其中的一个传入参数， 以便可以查询或修改该连接上的节点信息. virNodeGetInfo函数是获取节点的物理硬件信息，virNodeGetCPUStats函数可以获取节点上各个CPU的使用统计信息，virNodeGetMemoryStats函数可以获取节点上的内存的使用统计信息，virNodeGetFreeMemory函数可以获取节点上可用的空闲内存大小。 还有一些设置或者控制节点的函数, 如virNodeSetMemoryParameters函数可以设置节点上的内存调度的参数，virNodeSuspendForDuration函数可以让节点（宿主机） 暂停运行一段时间
1. 网络管理的API： 以virNetwork开头的一系列函数和部分以virInterface开头的函数

   libvirt也对虚拟化环境中的网络管理提供了丰富的API。 libvirt首先需要创建virNetworkPtr对象， 然后才能查询或控制虚拟网络。 查询网络相关信息的函数有，virNetworkGetName函数可以获取网络的名称， virNetworkGetBridgeName函数可以获取该网络中网桥的名称, virNetworkGetUUID函数可以获取网络的UUID标识，virNetworkGetXMLDesc函数可以获取网络的以XML格式的描述信息， virNetworkIsActive函数可以查询网络是否正在使用中。控制或更改网络设置的函数有，virNetworkCreateXML函数可以根据提供的XML格式的字符串创建一个网络（返回virNetworkPtr对象） ， virNetworkDestroy函数可以销毁一个网络（同时也会关闭使用该网络的域） ， virNetworkFree函数可以回收一个网络（但不会关闭正在运行的域） ，virNetworkUpdate函数可根据提供XML格式的网络配置来更新一个已存在的网络。 另外，virInterfaceCreate、 virInterfaceFree、 virInterfaceDestroy、 virInterfaceGetName、virInterfaceIsActive等函数可以用于创建、 释放和销毁网络接口， 以及查询网络接口的名称和激活状态
1. 存储卷管理的API： 以virStorageVol开头的一系列函数

   libvirt对存储卷（volume） 的管理主要是对域的镜像文件的管理， 这些镜像文件的格式可能是raw、 qcow2、 vmdk、 qed等。 libvirt对存储卷的管理， 首先需要创建virStorageVolPtr这个存储卷对象， 然后才能对其进行查询或控制操作。 libvirt提供了3个函数来分别通过不同的方式来获取存储卷对象， 如virStorageVolLookupByKey函数可以根据全局唯一的键值来获得一个存储卷对象， virStorageVolLookupByName函数可以根据名称在一个存储资源池（storage pool） 中获取一个存储卷对象， virStorageVolLookupByPath函数可以根据它在节点上的路径来获取一个存储卷对象。 有一些函数用于查询存储卷的信息， 如virStorageVolGetInfo函数可以查询某个存储卷的使用情况， virStorageVolGetName函数可以获取存储卷的名称， virStorageVolGetPath函数可以获取存储卷的路径，virStorageVolGetConnect函数可以查询存储卷的连接。 一些函数用于创建和修改存储卷，如virStorageVolCreateXML函数可以根据提供的XML描述来创建一个存储卷，virStorageVolFree函数可以释放存储卷的句柄（但是存储卷依然存在）,virStorageVolDelete函数可以删除一个存储卷， virStorageVolResize函数可以调整存储卷的大小
1. 存储池管理的API： 以virStoragePool开头的一系列函数

   libvirt对存储池（pool） 的管理包括对本地的基本文件系统、 普通网络共享文件系统、 iSCSI共享文件系统、 LVM分区等的管理。 libvirt需要基于virStoragePoolPtr这个存储池对象才能进行查询和控制操作。 一些函数可以通过查询获取一个存储池对象， 如virStoragePoolLookupByName函数可以根据存储池的名称来获取一个存储池对象，virStoragePoolLookupByVolume可以根据一个存储卷返回其对应的存储池对象。virStoragePoolCreateXML函数可以根据XML描述来创建一个存储池（默认已激活） ，virStoragePoolDefineXML函数可以根据XML描述信息静态地定义一个存储池（尚未激活） ， virStorage PoolCreate函数可以激活一个存储池。 virStoragePoolGetInfo、virStoragePoolGetName、 virStoragePoolGetUUID函数可以分别获取存储池的信息、 名称和UUID标识。 virStoragePool IsActive函数可以查询存储池状态是否处于使用中，virStoragePoolFree函数可以释放存储池相关的内存（但是不改变其在宿主机中的状态） ，virStoragePoolDestroy函数可以用于销毁一个存储池（但并没有释放virStoragePoolPtr对象， 之后还可以用virStoragePoolCreate函数重新激活它） ， virStoragePoolDelete函数可以物理删除一个存储池资源（该操作不可恢复）
1. 事件管理的API： 以virEvent开头的一系列函数。

   libvirt支持事件机制， 在使用该机制注册之后， 可以在发生特定的事件（如域的启动、 暂停、 恢复、 停止等） 时得到自己定义的一些通知。
1. 数据流管理的API： 以virStream开头的一系列函数
   libvirt还提供了一系列函数用于数据流的传输

## news
- [从v6.0.0开始libvirt-python.spec仅支持python3, 不再支持python2](https://github.com/libvirt/libvirt-python/commit/b22e4f2441078aec048b9503fde2b45e78710ce1)

## 安装与配置
安装方法: `dnf install libvirt`
查看version: `libvirtd --version`
libvirt的C API的使用: `dnf install libvirt-devel`
libvirt的Python API的使用: `dnf install libvirt-python`

libvirt相关配置在`/etc/libvirt`:
1. libvirt.conf

   libvirt.conf文件用于配置一些常用libvirt连接（通常是远程连接） 的别名, 比如:
   ```conf
   uri_aliases = [
   "remote1=qemu+ssh://root@192.168.93.201/system",
   ]
   ```

   此时可使用`virsh -c remote1`进行远程管理
1. libvirtd.conf

   libvirtd.conf是libvirt的守护进程libvirtd的配置文件， 被修改后需要让libvirtd重新加载配置文件（或重启libvirtd） 才会生效. 在libvirtd.conf中配置了libvirtd启动时的许多设置， 包括是否建立TCP、 UNIX domain socket等连接方式及其最大连接数， 以及这些连接的认证机制， 设置libvirtd的日志级别等.

   **在默认情况下， libvirtd在监听一个本地的Unix domain socket**, 而没有监听基于网络的TCP/IP socket. 要让TCP、 TLS等连接生效， 需要在启动libvirtd时加上--listen参数（简写为-l） , 而默认的systemctl start libvirtd命令在启动libvirtd服务时并没带--listen参数.

   **libvirtd守护进程的启动或停止， 并不会直接影响正在运行中的客户机**. libvirtd在启动或重启完成时， 只要客户机的XML配置文件是存在的， libvirtd会自动加载这些客户的配置， 获取它们的信息。 当然，如果客户机没有基于libvirt格式的XML文件来运行（例如直接使用qemu命令行来启动的客户机） ， libvirtd则不能自动发现它.

   通过libvirt启动客户机， 经过文件解析和命令参数的转换， 最终也会调用qemu命令行工具来实际完成客户机的创建.
1. qemu.conf

   qemu.conf是libvirt对QEMU的驱动的配置文件， 包括VNC、 SPICE等， 以及连接它们时采用的权限认证方式的配置， 也包括内存大页、 SELinux、 Cgroups等相关配置.
1. qemu

   在qemu目录下存放的是使用QEMU驱动的域的配置文件, 比如`centos7u2-1.xml`. 同时该目录下的networks目录保存了创建一个域时默认使用的网络配置

## 构建
参考:
- [Centos7.6 下编译安装 Libvirt 7.5](https://blog.frytea.com/archives/546/)
- [KVM安装及使用指南](https://bbs.huaweicloud.com/forum/forum.php?mod=viewthread&tid=113876)

前提:
1. qemu

参考源码中的README
```
# -- https://github.com/archlinux/svntogit-community/blob/packages/libvirt/trunk/PKGBUILD
rm /usr/lib/aarch64-linux-gun/libvirt*
wget https://libvirt.org/sources/libvirt-6.0.0.tar.xz
tar -xf libvirt-6.0.0.tar.xz
cd libvirt-6.0.0
mkdir build && cd build
apt install gnutls-bin ebtables
apt install libgnutls-dev libpciaccess-dev libnl-3-dev libnl-route-3-dev libdevmapper-dev libyajl-dev
../configure --prefix=/usr --sysconfdir=/etc --localstatedir=/var --with-qemu # `../autogen.sh --system --with-qemu` # 根据当前环境自动选择编译选项 # configure没有识别到自编译的qemu, 因此需要追加`--with-qemu`; `-–system`参数会尽可能保证安装目录与原生发行版系统的一致性
make && make install
ldconfig
systemctl start libvirtd
virsh version # 验证virsh
reboot # 否则virt-manger连接可能报错
```

libvirt-7.6.0编译:
```bash
$ apt install meson ninja # meson至少0.54
$ pip3 install rst2html5
$ mkdir build
$ meson build --prefix=/usr # meson build -Dsystem=true
$ ninja -C build
$ sudo ninja -C build install
```

> libvirt 7.2.0开始要求`glib 2.56`

> libvirt 6.3.0编译html docs时会报错

## 备份
- [virtnbdbackup 虚拟机备份](https://blog.csdn.net/ruanchao2012/article/details/129787012)
- [Efficient live full disk backup](https://libvirt.org/kbase/live_full_disk_backup.html)
- [Internals of incremental backup handling in qemu](https://libvirt.org/kbase/internals/incremental-backup.html)
- [Checkpoint XML format](https://libvirt.org/formatcheckpoint.html)
- [Domain state capture using Libvirt](https://libvirt.org/kbase/domainstatecapture.html)

## FAQ
### libvirt5.6.0源码并编译安装
```bash
# 1. 安装edk2
wget https://www.kraxel.org/repos/firmware.repo -O /etc/yum.repos.d/firmware.repo
yum -y install edk2.git-aarch64

或
dnf install dnf-plugins-core
dnf config-manager --add-repo https://www.kraxel.org/repos/firmware.repo
dnf install edk2.git-ovmf-x64

# 1. 安装依赖包
yum -y install libxml2-devel readline-devel ncurses-devel libtasn1-devel gnutls-devel libattr-devel libblkid-devel augeas systemd-devel libpciaccess-devel yajl-devel sanlock-devel libpcap-devel libnl3-devel libselinux-devel dnsmasq radvd cyrus-sasl-devel libacl-devel parted-devel device-mapper-devel xfsprogs-devel librados2-devel librbd1-devel glusterfs-api-devel glusterfs-devel numactl-devel libcap-ng-devel fuse-devel netcf-devel libcurl-devel audit-libs-devel systemtap-sdt-devel nfs-utils dbus-devel scrub numad rpm-build git

# 1. 下载源码并安装
wget https://libvirt.org/sources/libvirt-5.6.0-1.fc30.src.rpm
rpm -i libvirt-5.6.0-1.fc30.src.rpm

# 1. 生成rpm包，如果编译失败，可以重试
cd /root/rpmbuild/SPECS/
rpmbuild -ba libvirt.spec

# 1. 安装rpm包
cd /root/rpmbuild/RPMS/aarch64/
yum -y install *.rpm

# 1. 修改配置
vim /etc/libvirt/qemu.conf

#784行添加以下代码
nvram = ["/usr/share/edk2.git/aarch64/QEMU_EFI-pflash.raw:/usr/share/edk2.git/aarch64/vars-template-pflash.raw"]

# 1. 关闭SELinux并重启服务
# 重启libvirtd服务
systemctl restart libvirtd
# 关闭SELinux
setenforce 0 # 避免不能启动虚拟机
```

### `failed to connect to the hypervisor` & `failed to connect socket to '/var/run/libvirt/libvirt-sock': No such file or directory`
原因: libvirt服务未启动，找不到libvirt-sock.

解决方法: `systemctl start libvirtd`

### `Cannot read CA certificate '/etc/pki/CA/cacert.pem': No such file or directory`
当连接指定主机名(`qemu://system`或`qemu://session`, 使用**两个正斜杠**)时, QEMU 传输默认为 TLS, 这会要求证书.

> 使用三个正斜杠连接到的是本地主机, 例如`qemu:///system`, 不用tls.

解决方法: 见[TLSCreateCACert](https://wiki.libvirt.org/page/TLSCreateCACert)

ca.info:
```conf
cn = libvirt.org
ca
cert_signing_key
expiration_days = 7000
```

### `Cannot read certificate '/etc/pki/libvirt/servercert.pem': No such file or directory`
解决方法: 见[TLSCreateServerCerts](https://wiki.libvirt.org/page/TLSCreateServerCerts)

server.info:
```conf
organization = libvirt.org
cn = xxx
tls_www_server
encryption_key
signing_key
expiration_days = 7000
```

### `systemctl start libvirtd`日志报`direct firewall backed requested, but /sbin/ebtables is not avaliable: No such file or directory`
`apt install ebtables`

### `journalctl`报`failed to connect socket to '/var/run/libvirt/virtqemud-sock': No such file or directory`
参考:
- [libvirt.spec.in](https://github.com/libvirt/libvirt/blob/master/libvirt.spec.in)

virt-manager访问本地qemu虚拟机时用到, 应该是 virtqemud 服务没起来导致的, 还可能是编译libvirt时没有追加`--with-qemu`导致没有编译出该服务.

### libvirtd调试
参考:
- [libvirt ：debug and logging](https://blog.csdn.net/ggg_xxx/article/details/80060672)

方法:
1. `LIBVIRT_DEBUG=1 libvirtd`
1. `vim /etc/libvirt/libvirtd.conf`将`log_level`设为2

> 1时log不友好. 

### `YAJL 2 is required to build QEMU driver`
`apt install libyajl-dev`

### `XDR is required for remote driver `
`apt install libtirpc-dev`

### `Cannot check QEMU binary /usr/libexec/qemu-kvm: No such file or directory`
virt-manger创建vm时提示"警告:KVM不可用", 此时libvirtd就报该错. [根据qemu文档: The KVM project used to maintain a fork of QEMU called qemu-kvm. All feature differences have been merged into QEMU upstream and the development of the fork suspended.](https://wiki.qemu.org/Features/KVM), 即qemu-kvm现在已被弃用, 其中的所有代码以合并入 qemu-system-x86_64.

当前未找到libvirtd将qemu-kvm切换到qemu-system-x86_64的方法, 只是将`/usr/libexec/qemu-kvm`指向了qemu-system-x86_64, 但启用kvm需要添加显示参数`--enable-kvm`

此时libvirtd还会报`this function is not supported by the connect driver: cannot detect host CPU mode for aarch64 architecture`(虚拟机和物理机都会, 且物理机的kvm-ok是提示`KVM acceleration can be used`), 这个问题确定是libvirt接口不兼容引起的, 可参考[鲲鹏服务器安装kvm虚拟机，cannot detect host CPU model for aarch64 architect](https://bbs.huaweicloud.com/forum/thread-75455-1-1.html)和[libvirt启动时出现cannot detect host CPU model for architecture](https://bbs.huaweicloud.com/forum/thread-59280-1-1.html)

### virt-manger创建vm时提示`Falied to setup UEFI: 不能为架构'aarch64'找到任何UEFI二进制路径\nInstall options are limited`
`virt-manager --debug`报`Did not find any UEFI binary path for ...`

```bash
apt install qemu-efi # curl -o /etc/yum.repos.d/firmware.repo https://www.kraxel.org/repos/firmware.repo && yum install -y edk2.git-aarch64 
cat >> /etc/libvirt/qemu.conf <<'EOF'
nvram = [
    "/usr/share/AAVMF/AAVMF_CODE.fd:/usr/share/AAVMF/AAVMF_VARS.fd" # for ubuntu
    "/usr/share/edk2.git/aarch64/QEMU_EFI-pflash.raw:/usr/share/edk2.git/aarch64/vars-template-pflash.raw" # for centos
]
EOF
systemctl restart libvirtd
```

ps:
centos KVM对UEFI的支持包如下:
- 32位的arm处理器需要edk2.git-arm
- 64位的arm处理器需要edk2.git-aarch64
- x86-64处理器需要edk2.git-ovmf-x64

OVMF_CODE是bootloader的镜像文件，而OVMF_VARS则是保存OVMF_CODE中变量的文件. 在UEFI启动页面可以设置一些参数. 而这些参数的保存则需要OVMF_VARS文件, 它们配合使用: vm可以共享OVMF_CODE, 但使用自己的OVMF_VARS(通常放在`/var/lib/libvirt/qemu/nvram`或`/etc/libvirt/qemu/nvram`), vm xml相关配置可参考[这里](https://libvirt.org/formatdomain.html#bios-bootloader).

支持安全启动的UEFI，需要开启secure=’yes’ ，不过并不是所有的machine都支持，目前只支持q35系列，且需在feature中需要添加SMM.

### libvirtd.service自动退出
Ubuntu16.04.6+飞腾主板+libvirt 6.0.0, systemd里没有报错日志, 也没有crash.

### `'HWCAP_CPUID' undeclared`
内核版本太低, 比如4.4, HWCAP_CPUID没有定义. libvirt从6.4.0开始引入它.

### `pip install libvirt-python`报`Perhaps you should add the directory containing `libvirt.pc' to the PKG_CONFIG_PATH environment variable`和`Package 'libvirt', required by 'virtual:world', not found`
`dnf install libvirt libvirt-devel`

### virt-install报`cannot access storage file`
```bash
$ sudo vim /etc/libvirt/qemu.conf
user = "root"
group = "root"
$ sudo systemctl restart libvirtd
```

### `virsh undefine xxx`报`cannot undefine domain with nvram`
`virsh dumpxml 25 | grep nvram`报`<nvram>/var/lib/libvirt/qemu/nvram/centos8.0_VARS.fd</nvram>`

解决方法: `virsh undefine xxx --nvram`

报错源码: qemuDomainUndefineFlags()

### `virsh undefine xxx`报`cannot undefine transient domain`
之前创建过同名的domain, 此时要先`virsh destroy xxx`再`virsh undefine xxx`

解决方法: `virsh undefine xxx --nvram`

### `type=direct,source=eth0,source_mode=bridge,model=e1000`无法ping通网关
ref
- [确保 vSphere 标准交换机的安全](https://docs.vmware.com/cn/VMware-vSphere/7.0/com.vmware.vsphere.security.doc/GUID-3507432E-AFEA-4B6B-B404-17A020575358.html)和[混杂模式运行](https://docs.vmware.com/cn/VMware-vSphere/7.0/com.vmware.vsphere.security.doc/GUID-92F3AB1F-B4C5-4F25-A010-8820D7250350.html)
- [VMware 嵌套虚拟机网络ping 不通](https://blog.csdn.net/huannanchunli/article/details/78741574)

env:
- gateway,172.16.25.1, huawei(from wireshark展示)
- vmware esxi, 6.5.0, 4564106, 172.16.25.20
- host: ip,172.16.25.157;gateway,172.16.25.1;os,oracle linux 7.9; on vmware esxi

   - net.ipv4.ip_forward=1
- vm: ip,172.16.25.159;gateway,172.16.25.1;os,windows server 2012

xml:
```xml
    <interface type='direct'>
      <mac address='52:54:00:ce:a4:e9'/>
      <source dev='eth0' mode='bridge'/>
      <target dev='macvtap1'/>
      <model type='virtio'/>
      <boot order='3'/>
      <alias name='net0'/>
      <address type='pci' domain='0x0000' bus='0x00' slot='0x03' function='0x0'/>
    </interface>
```

vm无法ping通172.16.25.1, 将驱动换成virtio问题仍旧, 留意到vm 网卡状态上`已发送`挺多, 但`已接收`=0.

在host上, 刚开始发现`tcpdump -i macvtap1 icmp -X -vvv`有来源是172.16.25.159的数据包, 但`tcpdump -i eth0 icmp -X -vvv`上却没有, 后来多次重复该测试后, 发现`tcpdump -i eth0 icmp -X -vvv`是小概率没有来源是172.16.25.159的数据包(该小概率可能是tcpdump处理慢, 需要等待一会), 但gateway不响应vm的ICMP echo request. 后来又配置了一台linux vm 172.16.25.160, 发现172.16.25.159与172.16.25.160互ping正常, 因此应是gateway与host间出了问题, 但172.16.25.157 ping gateway又是正常的.

参考[vmware的抓包位置](/shell/cmd/virt/vmware.md), 在`--switchport <id> --capture PortInput/PortOutput`和`--uplink <vmnicX> --capture PortInput/PortOutput`抓包, 发现gateway有响应25.159的icmp, 但数据包进入`vSwitch0`的`
VM Network`后被丢弃了. 通过`vSphere client`->vm所在节点->配置->网络->虚拟交换机->编辑->安全, 将"混杂模式"改为"接受", 25.159 ping 25.1变正常.

> vm清空arp cache, 再ping 172.16.25.1, 虽然ping不通但arp cache能生成且正确

### 无法ping通网关2
ref:
- [关于KVM虚拟化网桥bridge的一个mac/port映射表故障分析](https://blog.csdn.net/watermelonbig/article/details/125118931)

env:
- vm(v, 8.103)-> host(h, kvm, macvtap0->eth0(8.120))->物理机(p, macvtap1 -> br0(8.9, on enp2s0f0))->gw(8.1)

   8.120 ping 8.1正常

ps:
1. p, h均已设置`net.ipv4.ip_forward=1`, 并关闭防火墙
1. p的macvtap1, br0, enp2s0f0/host的macvtap0, eth0均设置`promisc on`
1. p的`ip -s link show` macvtap1, br0, enp2s0f0 均没有丢包

在p的br0上抓包发现有v发送到gw且gw响应给v的icmp, 但p的mactap1上抓包仅有v发送给gw的icmp, 应该在br0和mactap1间丢包了, 具体位置无法定位. ???

> 在p的enp2s0f0上抓包和在br0上结果一致, 因为此时enp2s0f0就相当于一根网线而已.

### vm虚拟机网络问题
1. 宿主机的ip不通，就要确认下虚拟机网卡的类型

   - 对于macvlan网卡，ping不通宿主机网卡是正常的，但可以ping通同网段的别的ip
   - 对于bridge网卡，使用brctl show命令检查虚拟的网卡（比如vnet0）是否挂到了网桥上，若没有，使用brctl addif br0 vnet0添加

1. . 同网段的ip不通，就要考虑网卡转发

   - 内核ip_forward参数(sysctl net.ipv4.ip_forward): 0表示没有开启，需要设置为1
   - 网桥报文是否开启iptables过滤(`cat /proc/sys/net/bridge/bridge-nf-call-iptables`): 0表示没有开启，若为1则需要检查iptables配置

1. 不同网段的ip不通，检查路由
  
### `virsh insall`报`unsupported configuration: ACPI requires UEFI on this architecture`
[aarch64 KVM只支持UEFI BIOS，编译源码时未安装edk2, 无法识别Firmware文件](https://support.huaweicloud.com/trouble-kunpengcpfs/kunpengkvm_09_0006.html)

解决方法:
1. 使用uefi
1. 使用`virsh insall --features acpi=off`, 禁用acpi

   有用seabios aarch64, 但经测试还是报该错误

### `virsh insall`报`Couldn't find kernel for install tree`
不是使用`--location /home/me/Downloads/ubuntu-18.10-desktop-amd64.iso`, 而要采用`--cdrom /home/me/Downloads/ubuntu-18.10-desktop-amd64.iso`

原因是找不到文件: install/vmlinuz, install/initrd.gz

### virt-install uefi + cdrom
`--boot uefi --boot cdrom --cdrom xxx.iso`, 仅virt-install时有效, vm restart后进入uefi shell.

原因是vm xml里的cdrom缺source标签, 通过为virt-install添加`--debug`打印vm xml可见仅在virt-install时使用了iso, 关闭vm后xml里的source被删除:
```xml
<disk type='file' device='cdrom'>
  <driver name='qemu' type='raw'/>
  <source file='/tmp/SLES-11-DVD-i586-GM-DVD1.iso'/>
  <target dev='sda' bus='scsi'/>
  <readonly/>
  <address type='drive' controller='0' bus='1' unit='0'/>
</disk>
```

推测它的操作可能是:
```bash
# virsh attach-disk <GuestName> /tmp/SLES-11-DVD-i586-GM-DVD1.iso sda --type cdrom --mode readonly # 模拟最初的cdrom配置, 有iso信息
# cat update-device.xml
<disk type='file' device='cdrom'>
  <driver name='qemu' type='raw'/>
  <source file='/tmp/SLES-11-DVD-i586-GM-DVD1.iso'/>
  <target dev='sda' bus='scsi'/> # dev需要与bus对应
  <readonly/>
  <address type='drive' controller='0' bus='0' target='0' unit='0'/>
</disk>
# virsh update-device <GuestName> update-device.xml # 删除了source标签
```

### windows在kvm上鼠标不同步(飘)
现象: 出现两个鼠标, 且无法重合.

添加有且仅添加一个usb/virtio(**推荐使用usb, 在xp下virtio还是有漂移**)设备: `<input type="tablet" bus="usb">`. x64环境下, ps2的mouse和keyboard都是默认设备且无法替换为usb/virtio总线或删除, 已在centos 7上验证.

修正后效果: 鼠标移动过程中, 可能出现残影的现象, 但停止移动后会立即重合.

> 添加多个usb鼠标windows可能会蓝屏

> 试过两个ps设备+一个usb鼠标, 但还是飘.

### Guest has not initialized the display (yet) 
- [虽然qemu machine i440fx/q35都支持 BIOS 和 UEFI, 但**uefi推荐使用q35**](https://blog.csdn.net/m0_47541842/article/details/113521732)
- iso里os的arch与qemu使用的arch不一致
- kylinv10 host(aarm64) + `vm(osVariant:ubuntu 19.10 + uefi + vga)` + Ubuntu 20.04-arm64.iso : 启动过程**很慢(超过90s, 同时cpu负载高)**且装机界面是字符型, 上下移动光标会出现花屏. 显卡model.type换成virtio后正常

   `host(aarm64) + vm(uefi + vga)`发现很慢或者甚至不出现装机界面, 因此uefi配合virtio或qxl为佳.

或用`virt-manager --debug`调试.

### ide+pc-i440fx-4.2+uefi+rhel 8.8 启动卡在uefi logo界面
> 有个环境ide+pc-i440fx-4.2+uefi+rhel 8.10正常, 且这个两个环境的img都已包含ata_piix驱动???

推测是ide与uefi的兼容或ide驱动问题. q35(`virsh domcapabilities --machine pc-q35-5.1 | xmllint --xpath '/domainCapabilities/devices/disk' -`的bus)直接不支持ide.

解决方法:
1. 将ide换成virtio, 再进入resume执行`dracut -f`(dracut添加virtio驱动), 重启后恢复正常.
2. 给内核启动参数追加`console=ttyS0[,115200]`, 很神奇的方法, 能成功

### `unsupported configuration: spice graphics are not supported with this QEMU`
qemu构建时没有选中spice.

### [vm 磁盘扩容](https://opengers.github.io/openstack/openstack-instance-disk-resize-and-convert/)


```bash
# qemu-img info /data/test/centos7.qcow2
image: /data/test/centos6.qcow2
file format: qcow2
virtual size: 40G (85899345920 bytes)
disk size: 1.2G
# qemu-img resize centos7.qcow2 +40G # 如果virtual size以足够, 就忽略这步
```

扩容vdisk上的fs见[/shell/cmd/suit/part/resize.md]

### vm disk测试性能容易卡死
ref:
- [Virtualization Tuning and Optimization Guide7.2. Caching](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/virtualization_tuning_and_optimization_guide/sect-virtualization_tuning_optimization_guide-blockio-caching)
vm disk xml没有设置cache(即使用了默认选项), 使用none后问题消失. 同时none通常是支持迁移的最佳和唯一选项.

> io: ssd推荐使用native, 其他使用threads.

### 修改vm xml磁盘的bus而修改target.dev时`virsh define`会报`Found duplicate dirve address for disk with target name 'hdb'`
将原先vdb 的bus从virtio改为ide, 报了该错, 将vdb改为hdc后正常.

### ping虚拟机网络断断续续
ref:
- [浅析虚拟化环境网卡绑定模式](https://www.shsnc.com/show-109-956-1.html)

env:
- host: 10.0.4.171. bond0(mode 6)
- vm: 10.0.4.188
- gateway: 10.0.4.1
- test: 10.0.4.170

现象:
1. vm ping gateway: 断断续续
1. test ping vm: 断断续续
两种现象同时出现, 同时恢复.

解决: 使用bond mod=1, 业界大多虚拟化环境采用该模式稳定及高可用性已得到了充分的证明.

> zstack: 使用bond0(mod 4), 在其上创建bridge设备br_bond0, 再将br_bond0给vm

绑定设备的数据包传输算法是由绑定的模式所决定的，绑定模式共有7种(mode-0 ~ mode-6)，其中mode-1 ~ mode-4支持虚拟机网络(使用网桥)和非虚拟机网络(无网桥)；mode-0、mode-5、mode-6只支持非虚拟机网络(无网桥), ovirt虚拟化平台默认使用的是mode-4.

### 使用e1000接管的windows server 2019 datacenter没有识别到网卡
`设备管理器->其他设备->以太网控制器`提示叹号

解决: 选择该控制器, 右键选`更新驱动程序`->`自动搜索更新的驱动程序软件`, 完成后就可以看到网卡`Intel(R) PRO/1000 MT Network Connection`了

推荐: 换virtio网卡

NeoKylin Linux Aadvanced Sever V7Update6 arm64遇到同样问题, 换virtio网卡后正常.

### 接管的windows 2012R2启动卡在转圈
ref:
- [How To Enable Boot Log in Windows 10](https://winaero.com/enable-boot-log-windows-10/)
- [windows server 2019 断电无法启动修复](https://cloud.tencent.com/developer/article/1739700)

推测卡在ntfs scan. 强关再启动即可, 且启动卡住时间段的日志不在`事件查看器->Windows 日志->系统`里. 开启bootlog后, 发现C:\Windows\ntbtlog.txt还是未记录到卡住时的日志.

### aarch64上vm开机报`cpu mode 'host-model' for aarch64 kvm domain on aarch64 host is not supported by hypervisor`
解决方法:
1. 通过`qemu-kvm -machine virt -cpu help`
1. cpu_mode 改成 host-passthrough

### aarch64上vm开机报`failed to find romfile "vgabios-stdvga.bin"`
已安装seavgabios-bin.

解决方法: `ln -s /usr/share/seavgabios/vgabios-stdvga.bin /usr/share/qemu/vgabios-stdvga.bin`

### aarch64上vm开机报`this qemu does not support 'qxl' video device`
qemu 的 configure 将 spice 改为 yes 并再次编译安装

### virt-xml update boot order报`unsupported configuration: per-device boot elements cannot be used together with os/boot elements`
env: virt-install v2.2/v1.5

当磁盘,光驱和网卡都没有设置boot order或没有这些可引导设备时, xml会自动在`<os>`加上`<boot dev='xxx'/>`, 该项与设置设备的boot order(包括添加带boot order的device)冲突.

手动删除`<boot dev='hd'/>`, 再使用virt-xml添加带boot_order的disk device还是会报该错.

解决方法:
1. 必须保留一个带boot order的光驱
2. 自己构建xml

### 运行非root使用virsh
```bash
sudo su root
usermod -a -G libvirt <userName>
exit

vim ~/.bashrc # 添加如下内容
export LIBVIRT_DEFAULT_URI="qemu:///system"
source ~/.bashrc
```

selinux环境需在虚拟机XML配置文件中的domain根元素中添加如下内容, 使qemu-system-x86_64进程可以访问磁盘镜像文件: `<seclabel type='dynamic' model='dac' relabel='yes'>`

> seclabel元素允许控制安全驱动程序的操作.

### 创建kvm时libvirtd报`unsupported configuration: more than 255 vCPUs require extended interrupt mode enable on the iommu device`
ref:
- [Fail to start a q35 guest with vcpu > 255](https://bugzilla.redhat.com/show_bug.cgi?id=1451282)
- [让KVM突破限制，支持512个vCPU](https://github.com/GiantVM/doc/blob/master/extend_kvm_cpu.md)

`virsh domcapabilities --machine q35`返回上限是288, 但xml使用256时就报该错.

### 启动vm报`Unable to add bridge eth0 port vnet0: Operation not supported`
eth0不是brigde device.

### 启动vm报`error creating macvtap interface macvtap@eth0 (52:54:00:56:84:7a): Device or resource busy`
eth0被用于创建bond0, 此时应使用bond0

### `type=direct,source=eth0,source_mode=bridge,model=e1000`
ref:
- [「KVM」- 实现 Host 与 Guest 访问 @20210314](https://blog.51cto.com/u_11101184/3137103)

   放弃 Guest 访问 Host 的做法: Host 本就处于虚拟化的底层, 不应该让 Guest 访问 Host 主机
- [Linux上虚拟网络与真实网络的映射](https://cloud.tencent.com/developer/article/1082914)

同网段ip的vm使用macvtap0, 与host不通, 但与其他同网段的ip相通

原因: 这实际上不是一个错误, 而是 macvtap 的定义行为. 由于主机的物理以太网连接到 macvtap 网桥的方式, 从guest进入该网桥并转发到物理接口的流量无法返回到host的 IP 堆栈

此外, 从host IP 堆栈发送到物理接口的流量无法返回到 macvtap 网桥以转发到guest

解决方法:
1. Creating a separate macvtap interface for the host
1. Using libvirt for creating an isolated network
1. 多网卡: vms和host使用独立网卡

### 启动vm报`unsupported configuration: disk type 'virtio' of 'vdb' does not support ejectable media'`
原因:cdrom使用了virtio.

cdrom bus获取:
1. virt domcapabilities ... 获取disk bus, 过滤fdc,usb,virtio,scsi
1. q35过滤ide

> 实践发现bus=scsi的iso起不来, 包括oracle 7.9 x64, windows2012r2

### vm启动后发现时间比host快8h
ref:
- [kvm虚拟化环境中的时区设置](https://opengers.github.io/virtualization/kvm-guest-clock-timezone/)

vm和host都是东8区.

处理方法:
- 如果guest OS是Linux系统，应该选用utc，guest OS在启动时便会向host同步一次utc时间，然后根据/etc/localtime中设置的时区，来计算系统时间
- 如果guest OS是windows系统，则应该选用localtime，guest OS在启动时向host同步一次系统时间

原因: 当前guest使用centos, 因此将其xml配置的clock offset的localtime改为utc即可.

### [windows接管报`No physical memory available at the location required for the Windows Boot Manager.`](https://answers.microsoft.com/en-us/windows/forum/all/no-physical-memory-is-available-at-the-location/33e1c1de-dd92-45bb-a9b2-b7be3611e6ed)
应是磁盘或内存问题. 经排查是创建vm xml时内存参数是2048Kib, 导致内存不足引发的.


### Unable to open /dev/kvm: No such file or director
步骤:
- 执行`systemd-detect-virt`探明环境

   - vmware: 开启嵌套虚拟化
   - none: 查cpu指令集是否支持虚拟化

### `virsh start xxx`报`internal error: qemu unexpectedly closed the monitor: Could not access KVM kernel module: Permission denied\n...qemu-system-x86_64: failed to initialize KVM: Permission denied`
`ls -al /dev/kvm`返回`crw-rw----+`m, 存在acl属性

解决(2种):
1. `setfacl -b /dev/kvm && chmod 0660 /dev/kvm`(重启后失效)
2. 修改/etc/libvirt/qemu.conf, 将user和group都设为root

### 通过snap接管的含lvm的vm启动失败
env: oracle linux 7.9
需要将`/root`分区所在的disk, 放在xml disks的第一个

### 添加ide cdrom失败`unsupported configuration: Only a single IDE controller is supported for this machine type`
ref:
- [qemuValidateDomainDeviceDefControllerIDE](https://github.com/libvirt/libvirt/blob/master/src/qemu/qemu_validate.c#L3403)
- [/usr/share/virt-manager/virtinst/devices/disk.py", line 917]

   virt-manager对ide设备的数量检查 from 添加ide cdrom报错

env:
- libvirtd (libvirt) 8.0.0
- 当前vm已有4个cdrom

[当前机器类型其一个ide控制器只能连接4个设备](https://bbs.sangfor.com.cn/forum.php?mod=viewthread&tid=136273). virt-manager添加超过4个ide cdrom时会提示"仅支持总线'ide'的4个磁盘"

当前观察到x86仅支持一个ide controller.

### `unsupported configuration: IDE controllers are unsupported for this QEMU binary or machine type`
arm64不支持IDE总线

### cdrom启动"No bootable device"
bus使用了scsi, 换成sata,ide就正常了.

### windows启动"No bootable device"
带引导信息的磁盘必须是xml的第一块磁盘

### window启动报`Windows 未能启动。原因可能是最近更改了硬件或软件`, 错误码是`0xc000000f/0xc000000e`, 信息是`未连接或无法访问所需设备`
初步怀疑是磁盘签名的问题

### 创建vm时将默认q35改为i440fx, 导致创建vm报错`must be model='pci-root' for this machine type, but model='pcie-root' was found instead`
ref:
- [virt-manager#398](https://github.com/virt-manager/virt-manager/issues/398)

env:
- virt-manager : 4.0.0
- vm os: centos 7

virt-manager根据vm os默认选中q35时xml已填充了部分controller数据, 该数据和之后的i440fx配置冲突.

解决方法: 将当前xml保存为tmp.xml, 在删除其中所有`<controller type='pci*'/>`, 再`virsh define tmp.xml`即可.

### `internal error: Failed to reserve port 5900`
vnc端口5900已被占用, 修改端口即可

### [virt-install报`Failed to get "write" lock`](https://serverfault.com/questions/1057939/not-able-to-start-virtual-machine-domain-in-kvm-failed-to-get-write-lock)
xml disk中多次使用了相同的盘, 此时`fuser, lsof, lslocks`均没法反映该情况

### virsh list报`error: Failed to initialize libvirt`
env:
- FusionCompute

应是认证问题, 切换到root后正常.

### 使用usb 2/3
ref:
- [虚拟机配置](https://docs.openeuler.org/zh/docs/22.03_LTS/docs/Virtualization/%E8%99%9A%E6%8B%9F%E6%9C%BA%E9%85%8D%E7%BD%AE.html)
- [usb](https://libvirt.org/formatdomain.html#controllers)

`virt-xml --build-xml --controller type=usb,model=qemu-xhci`, 由model控制:
- usb3 : qemu-xhci
- usb2 : ich9-ehci1

### debug vm boot
ref:
- [开虚拟机串口控制台](https://docs.redhat.com/zh_hans/documentation/red_hat_enterprise_linux/9/html/configuring_and_managing_virtualization/proc_opening-a-virtual-machine-serial-console_assembly_connecting-to-virtual-machines)

1. 编辑grub启动项, 追加`console=ttyS0[,115200]`, 按ctrl+x启动即可
1. `virsh console xxx`, 启动信息会输出在terminal里

> 遇到`vm centos 7.7 使用ide/sata`启动后图形界面卡住(底层大概是进入了dracut), 但其grub追加`console=ttyS0,115200`后能正常进入系统或使用virtio后进入dracut. 这种情况通常是vmware vm接管到kvm, 因initramfs里的驱动差异导致的. 解决方法: 启动时选恢复模式, 再执行`dracut -f`(dracut会自动识别硬件并更新initramfs里的驱动), 最后重启即可.

> 类似的案例: 接管openeuler 22.03 x64 on vmware 使用sata或virtio卡住并进入终端界面(非dracut)/scsi无启动设备

### virbr0和virbr0-nic
ref:
- [libvirt之virbr0和virbr0-nic](https://xiaoz.info/2020/01/08/libvirt-virbr0/)

libvirtd会自动创建一个virbr0, 它是一个virtual network switch(bridge device), 所有虚拟机都将连接到virbr0. 网桥virbr0，相当于VMware的 VMNET8，提供NAT的网卡.

默认virbr0使用NAT模式, 可以提供NAT模式上网. 默认情况下, virbr0分配地址192.168.122.1, 它可以为连接到它的其他虚拟接口提供 DHCP 服务.

virbr0包括两个端口：virbr0-nic 为网桥内部端口，vnet0 为虚拟机网关端口(192.168.122.1).

> 增加virbr0-nic接口是为了解决一个内核的bug(或者说是feature)。创建bridge后，当我们添加第一块虚拟NIC到bridge时，这块NIC的MAC地址会复制到bridge，作为bridge的MAC地址。当我们所有NIC从bridge移除之后，这时bridge会丢失原来的MAC地址。而再次加入另外的NIC时，bridge又会获取新的MAC地址，这个MAC地址获取的是新NIC的MAC地址. virbr0-nic其实是一个[dummy device](https://xiaoz.info/2020/01/08/libvirt-virbr0/).

> libvirt net的NAT 模式（默认）将 virbr0 当作一个 Firewall NAT 设备. 底层支撑使用了 iptables nat 的 MASQUERADE（地址伪装）规则类型，而非 SNAT 或 DNAT. IP 地址伪装规则，使得 VMs 可以使用 Host 的 IP 地址访问外部网络，但外部网络无法主动访问到 VMs，因为 VMs 对于外部网络而言是不可见的.

### vm id
一旦vm运行中, `virsh list --all`就会输出其id, 包括paused(暂停中).

virDomainGetID可能返回4294967295, 它即[`^uint32(0)=(unsigned int)-1)`](https://github.com/libvirt/libvirt/blob/master/tools/virsh-domain-monitor.c#L1231)等价于`virsh list`中的`-`.

### GPU 透传
打开 iommu, VT-d

### domain type
- kvm: linux kvm
- qemu: 完全虚拟化
- hvf: 苹果公司虚拟化

## uefi shell
- exit : 进入qemu machine(virt-4.0)的类似bios界面的字符uefi firmware settings界面.

   在grub shell执行fwsetup也可进入字符uefi firmware settings界面

## virtsh
virsh 属于 libvirt 的命令行工具, 与virt-manager类似, libvirt 是目前使用最为广泛的对 KVM 虚拟机进行管理的工具和 API, 它还可管理 VMware, VirtualBox, Hyper-V等.

Libvirt 分服务端和客户端, Libvirtd 是一个 daemon 进程, 是服务端, 可以被本地的 virsh 调用, 也可以被远程的 virsh 调用, virsh 相当于客户端.

virsh同时支持交互模式和非交互模式.

> virsh是用C语言编写的一个使用libvirt API的虚拟化管理工具. virsh程序的源代码在libvirt项目源代码的tools目录下， 实现virsh工具最核心的一个源代码文件是virsh.c

## viewer
- [remote-viewer](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/virtualization_deployment_and_administration_guide/sect-graphic_user_interface_tools_for_guest_virtual_machine_management-remote_viewer)

   ```bash
   remote-viewer spice://testguest:5900
   remote-viewer vnc://testguest2:5900
   ```

- [virt-viewer](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/virtualization_deployment_and_administration_guide/chap-graphic_user_interface_tools_for_guest_virtual_machine_management#sect-Graphic_User_Interface_tools_for_guest_virtual_machine_management-Using_virt_viewer_command_line)

   ```bash
   virt-viewer --connect qemu+ssh://root@192.168.88.151/system test
   ```

### 常用命令
ref:
- `man virsh`
- [<<KVM实战>>的4.2 virsh]
- [QEMU中的命令行参数及其monitor中的命令， 在virsh中的对应关系](http://wiki.libvirt.org/page/QEMUSwitchToLibvirt) 
- [热迁移虚拟机](https://docs.openeuler.org/zh/docs/20.03_LTS_SP3/docs/Virtualization/%E7%83%AD%E8%BF%81%E7%A7%BB%E8%99%9A%E6%8B%9F%E6%9C%BA.html)
- [centos7上使用virt-install命令创建kvm虚拟机](https://blog.51cto.com/u_11555417/2341874)
- [Networking](http://wiki.libvirt.org/page/Networking)
- [MacVTap](https://virt.kernelnewbies.org/MacVTap)
- [虚拟机 XML 配置示例](https://docs.redhat.com/zh_hans/documentation/red_hat_enterprise_linux/9/html/configuring_and_managing_virtualization/sample-virtual-machine-xml-configuration_viewing-information-about-virtual-machines)

如下命令启动虚拟机： `virsh create <name of virtual machine>` : 通过`virsh create <vmname>.xml`创建的虚拟机不会持久化，关机后会消失
启动虚拟机： `virsh start <name>`
列出所有虚拟机 (不管是否运行)： `virsh list --all`, `--all`包括没运行的vm, 则只输出运行中的vm
正常关闭 guest ： `virsh shutdown <virtual machine (name | id | uuid)>`
强制关闭 guest ： `virsh destroy <virtual machine (name | id | uuid)>`, 通常只需要几秒, 有次遇到是81s
挂起vm: `virsh suspend <name>`
恢复被挂起的vm: `virsh resumed <name>`
开机自启动vm: `virsh autostart <name>`
连接vm: `virsh console <name>`, 退出使用`ctrl+]`
保存虚拟机快照到文件： `virsh save <virtual machine (name | id | uuid)> <filename>`
从快照恢复虚拟机： `virsh restore <filename>`
查看虚拟机配置文件： `virsh dumpxml <virtual machine (name | id | uuid)`
删除vm的配置文件: `virsh undifine <name>`
根据配置文件定义vm: `virsh define <file.xml>`
列出全部 virsh 可用命令： `virsh help`
help: `virt-install <参数> ?`

    ```conf
    # virsh help
    command：

     Domain Management (help keyword 'domain'):
        attach-device                  从一个XML文件附加装置
        attach-disk                    附加磁盘设备. 即时生效，但系统重启后新硬盘会消失. 永久方法: 修改vm xml.

         virsh attach-disk 361way /data1/kvm.img vdc
        attach-interface               获得网络设备. 添加网卡:virsh attach-interface vm-yaohai --type bridge --source br0 --model virtio --config; 删除网卡(by mac): virsh detach-interface vm-yaohai --type bridge --mac 52:54:00:61:4c:f3 --config
        autostart                      自动开始一个域
        blkdeviotune                   设定或者查询块设备 I/O 调节参数。
        blkiotune                      获取或者数值 blkio 参数
        blockcommit                    启动块提交操作。
        blockcopy                      启动块复制操作。
        blockjob                       管理活跃块操作
        blockpull                      使用其后端映像填充磁盘。
        blockresize                    创新定义域块设备大小
        change-media                   更改 CD 介质或者软盘驱动器
        console                        连接到客户会话

            前提:
            1.  vim /etc/default/grub `GRUB_CMDLINE_LINUX_DEFAULT="console=tty0 console=ttyS0,115200n8"` && `grub2-mkconfig -o /boot/grub2/grub.cfg`

               GRUB_CMDLINE_LINUX_DEFAULT会片接在GRUB_CMDLINE_LINUX后
            1. `systemctl enable --now serial-getty@ttyS0.service`

            退出用`Ctrl + Shift + 5`
        cpu-stats                      显示域 cpu 统计数据
        create                         从一个 XML 文件创建一个域
        define                         从一个 XML 文件定义（但不开始）一个域
             virsh define xxx/foo.xml : reload xml. 手动删除foo.xml里的ps2设备, reload时会被重新添加. `virsh edit也无法删除ps2设备`
        desc                           显示或者设定域描述或者标题
        destroy                        强制关闭域
        detach-device                  从一个 XML 文件分离设备
        detach-device-alias            detach device from an alias
        detach-disk                    分离磁盘设备
        detach-interface               分离网络界面
        domdisplay                     域显示连接 URI
        domfsfreeze                    Freeze domain's mounted filesystems.
        domfsthaw                      Thaw domain's mounted filesystems.
        domfsinfo                      Get information of domain's mounted filesystems.
        domfstrim                      在域挂载的文件系统中调用 fstrim。
        domhostname                    输出域主机名
        domid                          把一个域名或 UUID 转换为域 id
        domif-setlink                  设定虚拟接口的链接状态
        domiftune                      获取/设定虚拟接口参数
        domjobabort                    忽略活跃域任务
        domjobinfo                     域任务信息
        domname                        将域 id 或 UUID 转换为域名
        domrename                      rename a domain
        dompmsuspend                   使用电源管理功能挂起域
        dompmwakeup                    从 pmsuspended 状态唤醒域
        domuuid                        把一个域名或 id 转换为域 UUID
        domxml-from-native             将原始配置转换为域 XML
        domxml-to-native               将域 XML 转换为原始配置
        dump                           把一个域的内核 dump 到一个文件中以方便分析
        dumpxml                        XML 中的域信息
        edit                           编辑某个域的 XML 配置
        event                          Domain Events
        inject-nmi                     在虚拟机中输入 NMI
        iothreadinfo                   view domain IOThreads
        iothreadpin                    control domain IOThread affinity
        iothreadadd                    add an IOThread to the guest domain
        iothreaddel                    delete an IOThread from the guest domain
        send-key                       向虚拟机发送序列号
        send-process-signal            向进程发送信号
        lxc-enter-namespace            LXC 虚拟机进入名称空间
        managedsave                    管理域状态的保存
        managedsave-remove             删除域的管理保存
        managedsave-edit               edit XML for a domain's managed save state file
        managedsave-dumpxml            Domain information of managed save state file in XML
        managedsave-define             redefine the XML for a domain's managed save state file
        memtune                        获取或者数值内存参数(KB)
                                             - hard_limit : 最大可用mem
                                             - soft_limit : 竞争时的mem
                                             - swap_hard_limit : 最大内存加swap
                                             - min_guarantee : 最低保证给vm使用的内存

                                             生效方法:
                                             - `--config` : 写入xml, 重启vm生效
                                             - `--live` : 影响正在运行的vm, vm停止后, 效果消失, **默认**
                                             - `--current` : 影响停止和正在运行的vm, 如果vm运行, vm停止后效果消失

                                             examples:
                                             - `virsh memtune c7 --hard-limit 9437184 --config` : 最大限制在9G, 重启vm生效


                                             限制host将swap分配给vm:
                                             ```xml
                                             <memoryBacking><locked/></memoryBacking> # 设置locked时必须有hard_limit
                                             ```
        perf                           Get or set perf event
        metadata                       show or set domain's custom XML metadata
        migrate                        将域迁移到另一个主机中
        migrate-setmaxdowntime         设定最大可耐受故障时间
        migrate-getmaxdowntime         get maximum tolerable downtime
        migrate-compcache              获取/设定压缩缓存大小
        migrate-setspeed               设定迁移带宽的最大值
        migrate-getspeed               获取最长迁移带宽
        migrate-postcopy               Switch running migration from pre-copy to post-copy
        numatune                       获取或者数值 numa 参数
        qemu-attach                    QEMU 附加
        qemu-monitor-command           QEMU 监控程序命令
        qemu-monitor-event             QEMU Monitor Events
        qemu-agent-command             QEMU 虚拟机代理命令
        reboot                         重新启动一个域
        reset                          重新设定域
        restore                        从一个存在一个文件中的状态恢复一个域
        resume                         重新恢复一个域
        save                           把一个域的状态保存到一个文件
        save-image-define              为域的保存状态文件重新定义 XML
        save-image-dumpxml             在 XML 中保存状态域信息
        save-image-edit                为域保存状态文件编辑 XML
        schedinfo                      显示/设置日程安排变量
        screenshot                     提取当前域控制台快照并保存到文件中
        set-lifecycle-action           change lifecycle actions
        set-user-password              set the user password inside the domain
        setmaxmem                      改变最大内存限制值
        setmem                         改变内存的分配
        setvcpus                       改变虚拟 CPU 的号
        shutdown                       关闭一个域
        start                          开始一个（以前定义的）非活跃的域
        suspend                        挂起一个域
        ttyconsole                     tty 控制台
        undefine                       取消定义一个域, 若虚拟机启动时使用了nvram文件，销毁该虚拟机需要指定nvram的处理策略: keep-nvram/nvram, 其他nvram是销毁. 该命令不会删除xml, 需要自行删除
        update-device                  从 XML 文件中关系设备
        vcpucount                      域 vcpu 计数
        vcpuinfo                       详细的域 vcpu 信息. 一个vm默认只能使用同一颗物理cpu的逻辑核.
        vcpupin                        控制或者查询域 vcpu 亲和性. `vcpupin 21 0 28`: 强制vcpu0绑定到物理逻辑核28.
        emulatorpin                    控制或查询域模拟器亲和性, 即vm可使用host的哪些逻辑cpu. `emulatorpin 21 26-32 --live`: 强制让vm只能在部分物理逻辑核之间调度
        vncdisplay                     查询vnc连接信息
        guestvcpus                     query or modify state of vcpu in the guest (via agent)
        setvcpu                        attach/detach vcpu or groups of threads
        domblkthreshold                set the threshold for block-threshold event for a given block device or it's backing chain element

     Domain Monitoring (help keyword 'monitor'):
        domblkerror                    在块设备中显示错误
        domblkinfo                     域块设备大小信息
        domblklist <name>              列出所有域块
        domblkstat                     获得域设备块状态
        domcontrol                     域控制接口状态
        domif-getlink                  获取虚拟接口链接状态
        domifaddr <name>               Get network interfaces' addresses for a running domain. 获取虚拟机的IP地址(NAT DHCP分配的)
        domiflist <name>               列出所有域虚拟接口, 比如mac
        domifstat                      获得域网络接口状态
        dominfo                        域信息
        dommemstat                     获取域的内存统计
        domstate                       域状态
        domstats                       get statistics about one or multiple domains
        domtime                        domain time
        list                           列出域

     Host and Hypervisor (help keyword 'host'):
        allocpages                     Manipulate pages pool size
        capabilities                   性能
        cpu-baseline                   计算基线 CPU
        cpu-compare                    使用 XML 文件中描述的 CPU 与主机 CPU 进行对比
        cpu-models                     CPU models
        domcapabilities                domain capabilities

         virsh domcapabilities --machine q35
         virsh domcapabilities --machine q35 | xmllint --xpath '/domainCapabilities/os' -
        freecell                       NUMA可用内存
        freepages                      NUMA free pages
        hostname                       打印管理程序主机名
        hypervisor-cpu-baseline        compute baseline CPU usable by a specific hypervisor
        hypervisor-cpu-compare         compare a CPU with the CPU created by a hypervisor on the host
        maxvcpus                       连接 vcpu 最大值
        node-memory-tune               获取或者设定节点内存参数
        nodecpumap                     节点 cpu 映射
        nodecpustats                   输出节点的 cpu 状统计数据。
        nodeinfo                       节点信息
        nodememstats                   输出节点的内存状统计数据。
        nodesuspend                    在给定时间段挂起主机节点
        sysinfo                        输出 hypervisor sysinfo
        uri                            打印管理程序典型的URI
        version                        显示版本

     Interface (help keyword 'interface'):
        iface-begin                    生成当前接口设置快照，可在今后用于提交 (iface-commit) 或者恢复 (iface-rollback)
        iface-bridge                   生成桥接设备并为其附加一个现有网络设备
        iface-commit                   提交 iface-begin 后的更改并释放恢复点
        iface-define                   define an inactive persistent physical host interface or modify an existing persistent one from an XML file
        iface-destroy                  删除物理主机接口（启用它请执行 "if-down"）
        iface-dumpxml                  XML 中的接口信息
        iface-edit                     为物理主机界面编辑 XML 配置
        iface-list                     物理主机接口列表
        iface-mac                      将接口名称转换为接口 MAC 地址
        iface-name                     将接口 MAC 地址转换为接口名称
        iface-rollback                 恢复到之前保存的使用 iface-begin 生成的更改
        iface-start                    启动物理主机接口（启用它请执行 "if-up"）
        iface-unbridge                 分离其辅助设备后取消定义桥接设备
        iface-undefine                 取消定义物理主机接口（从配置中删除）

     Network Filter (help keyword 'filter'):
        nwfilter-define                使用 XML 文件定义或者更新网络过滤器
        nwfilter-dumpxml               XML 中的网络过滤器信息
        nwfilter-edit                  为网络过滤器编辑 XML 配置
        nwfilter-list                  列出网络过滤器
        nwfilter-undefine              取消定义网络过滤器
        nwfilter-binding-create        create a network filter binding from an XML file
        nwfilter-binding-delete        delete a network filter binding
        nwfilter-binding-dumpxml       XML 中的网络过滤器信息
        nwfilter-binding-list          list network filter bindings

     Networking (help keyword 'network'):
        net-autostart                  自启动网络
        net-create                     从一个 XML 文件创建一个网络
        net-define                     define an inactive persistent virtual network or modify an existing persistent one from an XML file
        net-destroy                    销毁（停止）网络
        net-dhcp-leases <default>      print lease info for a given network: 给domain分配的ip(包括之前分配的还在租约内的ip)
        net-dumpxml                    XML 中的网络信息
        net-edit                       为网络编辑 XML 配置.  `virsh net-edit default`
        net-event                      Network Events
        net-info                       网络信息
        net-list                       列出网络. 获取default配置: `cat /etc/libvirt/qemu/networks/default.xml`
        net-name                       把一个网络UUID 转换为网络名
        net-start                      开始一个(以前定义的)不活跃的网络
        net-undefine                   undefine a persistent network
        net-update                     更新现有网络配置的部分
        net-uuid                       把一个网络名转换为网络UUID

     Node Device (help keyword 'nodedev'):
        nodedev-create                 根据节点中的 XML 文件定义生成设备
        nodedev-destroy                销毁（停止）节点中的设备
        nodedev-detach                 将节点设备与其设备驱动程序分离
        nodedev-dumpxml                XML 中的节点设备详情
        nodedev-list                   这台主机中中的枚举设备

         输出可作`--host-device`参数
        nodedev-reattach               重新将节点设备附加到他的设备驱动程序中
        nodedev-reset                  重置节点设备
        nodedev-event                  Node Device Events

     Secret (help keyword 'secret'):
        secret-define                  定义或者修改 XML 中的 secret
        secret-dumpxml                 XML 中的 secret 属性
        secret-event                   Secret Events
        secret-get-value               secret 值输出
        secret-list                    列出 secret
        secret-set-value               设定 secret 值
        secret-undefine                取消定义 secret

     Snapshot (help keyword 'snapshot'):
        snapshot-create                使用 XML 生成快照
        snapshot-create-as             使用一组参数生成快照
        snapshot-current               获取或者设定当前快照
        snapshot-delete                删除域快照
        snapshot-dumpxml               为域快照转储 XML
        snapshot-edit                  编辑快照 XML
        snapshot-info                  快照信息
        snapshot-list                  为域列出快照
        snapshot-parent                获取快照的上级快照名称
        snapshot-revert                将域转换为快照

     Storage Pool (help keyword 'pool'):
        find-storage-pool-sources-as   找到潜在存储池源
        find-storage-pool-sources      发现潜在存储池源
        pool-autostart                 自动启动某个池
        pool-build                     建立池
        pool-create-as                 从一组变量中创建一个池
        pool-create                    从一个 XML 文件中创建一个池
        pool-define-as                 在一组变量中定义池
        pool-define                    define an inactive persistent storage pool or modify an existing persistent one from an XML file
        pool-delete                    删除池
        pool-destroy                   销毁（删除）池
        pool-dumpxml                   XML 中的池信息
        pool-edit                      为存储池编辑 XML 配置
        pool-info                      存储池信息
        pool-list                      列出池
        pool-name                      将池 UUID 转换为池名称
        pool-refresh                   刷新池
        pool-start                     启动一个（以前定义的）非活跃的池
        pool-undefine                  取消定义一个不活跃的池
        pool-uuid                      把一个池名称转换为池 UUID
        pool-event                     Storage Pool Events

     Storage Volume (help keyword 'volume'):
        vol-clone                      克隆卷。
        vol-create-as                  从一组变量中创建卷
        vol-create                     从一个 XML 文件创建一个卷
        vol-create-from                生成卷，使用另一个卷作为输入。
        vol-delete                     删除卷
        vol-download                   将卷内容下载到文件中
        vol-dumpxml                    XML 中的卷信息
        vol-info                       存储卷信息
        vol-key                        为给定密钥或者路径返回卷密钥
        vol-list                       列出卷
        vol-name                       为给定密钥或者路径返回卷名
        vol-path                       为给定密钥或者路径返回卷路径
        vol-pool                       为给定密钥或者路径返回存储池
        vol-resize                     创新定义卷大小
        vol-upload                     将文件内容上传到卷中
        vol-wipe                       擦除卷

     Virsh itself (help keyword 'virsh'):
        cd                             更改当前目录
        echo                           echo 参数
        exit                           退出这个非交互式终端
        help                           打印帮助
        pwd                            输出当前目录
        quit                           退出这个非交互式终端
        connect                        连接（重新连接）到 hypervisor
    ```

virt-xml(需关机操作):
```
virt-xml testguest --remove-device --disk target=vdb
virt-xml --build-xml --disk virt-xml --build-xml --disk /mnt/1.qcow2,target=vdb # test option
virt-xml --build-xml --network type=bridge,source=br0
virt-xml --remove-device --disk target=sda
virt-xml --add-device --disk xxx
virt-xml vs002 --edit target=sda --disk path=''
virt-xml vs002 --edit target="sda" --disk="bonmap -sP 192.168.0.0/24**](https://documentation.suse.com/sles/15-SP1/html/SLES-all/cha-libvirt-config-virsh.html)
- [koan/virtinstall.py](https://github.com/cobbler/koan/blob/master/koan/virtinstall.py)

创建vm:
```bash
# --- virsh 5.5
qemu-img create -f qcow2 centos_kvm1.qcow2 16G
virt-install \
--virt-type=kvm \
--name=centos-kvm \
--hvm \
--vcpus=1 \
--memory=1024 \
--cdrom=/srv/kvm/CentOS-7-x86_64-Minimal-1810.iso \
--disk path=/srv/kvm/centos_kvm1.qcow2,size=16,format=qcow2 \
--graphics vnc,password=kvm,listen=0.0.0.0,port=5911 \
--network bridge=virbr0 \
--boot uefi \
--autostart
```

> 生成的xml在`/etc/libvirt/qemu/<name>.xml`

安装成功后使用任意一个可以访问 KVM 宿主机的带有桌面的设备上的 VNC viewer 进入 `<vm宿主机ip>:5911`, 输入密码 `kvm` 就可以进入虚拟机, 然后继续安装了.

install 常用参数说明展开目录:
- 一般选项
nmap -sP 192.168.0.0/24
   - name : 指定虚拟机名称
   - memory: 分配内存大小, 单位是MB
   - vcpus : 分配CPU核心数，最大与实体机CPU核心数相同
   - cpu=CPU：CPU模式及特性，如coreduo等；可以使用`qemu-system-x86_64 -cpu ?`来获取支持的CPU模式
   - virt-type : hypervisor类型, 可使用`virsh capabilities`获取
   - os-variant=rhel6, 可用`osinfo-query os`获取, 信息来源于`/usr/share/osinfo`, 较新的os xml(比如`/usr/share/osinfo/os/centos.org/centos-stream-9.xml`)包含了支持的设备列表`<devices>`标签
      osinfo-query支持family, eol-date等
   - machine : machine类型, 可用`qemu-system-x86_64 -machine help`获取
   - soundhw: 声卡类型, 可用`qemu-system-x86_64 -soundhw help`获取
- 安装方式

   - cdrom=xxx.iso : 指定安装镜像iso

      第一次挂载的光驱, 重启后自动消失是针对linux的功能(其实就是将`xml <source>属性的file置空`); windows安装需要多次重启, 因此不会自动消失
   - location : 安装源URL, 多用于网络安装, 支持FTP、HTTP及NFS等, 但也支持本地路径, 如`ftp://172.16.0.1/pub`, `/xxx/x.iso/(mounted的iso目录)`
   - --boot  cdrom,hd,network：指定引导次序, 可用`virt-insall --boot ?`查看
   - --boot kernel=KERNEL,initrd=INITRD,kernel_args=”console=/dev/ttyS0”：指定启动系统的内核及initrd文件
   - pxe : 基于PXE完成安装
   - --import : 跳过os安装过程, 用现有磁盘镜像来构建vm, 常用`--disk`联用
   - --boot uefi : uefi启动. 估计是设置了`--boot loader=xxx.fd`, 与virt-manager 创建vm-`customize configuration before install`-overview-hypervisor details中的fireware相同.

      前提: `dnf install edk2-ovmf/apt install ovmf`

      验证: vm启动后显示uefi log(tianocore)
   - --boot /usr/share/seabios/biso.bin

      前提: `apt install seabios`
- 网络配置

   - network OPTIONS: 网络配置, 参考[Understanding virtual networking](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/9/html/configuring_and_managing_virtualization/configuring-virtual-machine-network-connections_configuring-and-managing-virtualization#understanding-virtual-networking-overview_configuring-virtual-machine-network-connections)


      [网络模式, 主要分为4种模式](https://www.cnblogs.com/good-study/p/15742351.html):
      1. host-only：又称隔离模式，可以理解为vmware的仅主机模式，意思就是将所有的虚拟机组成一个局域网，不能和外界通信，不能访问Internet，其他主机也不能访问虚拟主机，安全性高
      1. nat模式，虚拟出一个有nat转换的网络设备，虚拟机内部自动获取ip地址, 然后通过nat转换访问互联网. 这种模式内（虚拟机）访问外（虚拟机外）可以, 外不可以访问内.
      1. 桥接模式,与真实的物理网卡绑定，虚拟出交换机用于通信. 这网络模式下客户机与宿主机处于同一个网络环境，类似于一台真实的宿主机，直接访问网络资源，设置好后客户机与互联网，客户机与主机之间的通信都很容易
      1. Router模式：路由模式，当使用路由模式时，虚拟交换机连接到主机物理机器的物理LAN，在不适用NAT的情况下来回传输流量. 虚拟交换机可以检查所有流量，并使用网络数据包中包含的信息来作出路由决策

      - type=direct,source=eth0,source_mode=bridge,model=virtio : macvtap

         ```xml
         <interface type='direct'>
            <mac address='52:54:00:ce:a4:e9'/>
            <source dev='eth0' mode='bridge'/>
            <model type='virtio'/>
         </interface>
         ```

         > 网桥virbr0，相当于VMware的 VMNET8，提供NAT的网卡
      - bridge=br0 = `type=bridge,source=br0,model=virtio` : 使用桥接br0
      
         ```xml
         <interface type='bridge'>
            <mac address='52:54:00:eb:d7:7d'/>
            <source bridge='br0'/>
            <model type='virtio'/>
         </interface>
         ```
      - NETWORK=NAME : 即nat

         ```xml
         <interface type='network'>
            <mac address='52:54:00:eb:d7:7d'/>
            <source network='default'/>
            <address type='pci' domain='0x0000' bus='0x00' slot='0x03' function='0x0'/>
         </interface>
         <interface type='network'>
            <source network='default' portgroup='engineering'/>
            <target dev='vnet7'/>
            <mac address="00:11:22:33:44:55"/>
            <virtualport>
              <parameters instanceid='09b11c53-8b5c-4eeb-8f00-d84eaa0aaa4f'/>
            </virtualport>

          </interface>
         ```

      选项:
      - model: netdev model, 可用`qemu-system-x86_64 -net nic,model=?`获取
      - mac=52:54:00:01:02:03 : 指定mac, 对于 QEMU 或 KVM 虚拟机, 它必须是`52:54:00`, **注意检查mac重复, 实际发现随机生成有较高的概率**

      桥接和NAT区别: VM发包时, 源addr会被NAT修改掉; 桥接不修改.
- 存储配置

   - disk : 指定虚拟机的磁盘存储位置
     
      - size : 磁盘大小，以 GB 为单位.
      - format：磁盘映像格式，如raw、qcow2、vmdk等
      - none: 没有磁盘, 常用于livecd
      - bus：磁盘总结类型，其值可以为ide、scsi、usb、virtio或xen
      - perms：访问权限，如rw、ro或sh（共享的可读写），默认为rw
      - cache：缓存模型，其值有none、writethrouth（缓存读）及writeback（缓存读写）

         ![](http://blog.chinaunix.net/attachment/201210/12/20940095_1350022176aC72.png)
      - sparse：磁盘映像使用稀疏格式，即不立即分配指定大小的空间
      - boot_order: 多个磁盘用于安装时guest时的尝试引导的顺序, 越小优先
- 图形配置

   - graphics TYPE,opt1=val1,opt2=val2 : 图形化显示配置
     # 全新安装虚拟机过程中可能会有很多交互操作，比如设置语言，初始化 root 密码等等.
     # graphics 选项的作用就是配置图形化的交互方式，可以使用 vnc（一种远程桌面软件）进行链接.
     # 我们这列使用命令行的方式安装，所以这里要设置为 none，但要通过 --extra-args 选项指定终端信息，
     # 这样才能将安装过程中的交互信息输出到当前控制台.
     
     TYPE：指定显示类型，可以为vnc、sdl、spice或none等，默认为vnc`

     - vnc启用VNC远程管理，一般安装系统都要启用.

         - port : 指定VNC监控端口，默认端口为5900，端口不能重复.
         - listen : 指定VNC绑定IP，默认绑定127.0.0.1，这里改为0.0.0.0
         - password: TYPE为vnc或spice时，为远程访问监听的服务进指定认证密码

         例如: `--graphics vnc,password=123456,port=5910`
   - [video](https://wiki.archlinux.org/title/QEMU_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87)#%E5%9B%BE%E5%BD%A2)

      - std: 标准vga, windows, 使用uiefi, 推荐
      - vmware: vmware兼容显卡, vmware虚拟机, ubuntu等推荐
      - cirrus : cirrus兼容显卡, 默认显卡, redhat等推荐使用
      - none : 等于没有vga卡, 无法通过`-vnc`访问它. 与使用`-nographic`不同, `-nographic`会让QEMU模拟VGA卡且只是关闭了SDL输出
      - virtio : 一个基于virgl的3D并行虚拟化图形驱动
      - qxl : 一个支持2D的并行虚拟化图形驱动, 非redhat, 使用uefi, 推荐

- 其他

   - extra-args : 根据不同的安装方式设置不同的额外选项

      比如kickstart安装参数: `--location https://mirrors.aliyun.com/centos/8-stream/BaseOS/x86_64/os/ --initrd-inject /path/to/ks.cfg  --extra-args="ks=file:/ks.cfg console=tty0 console=ttyS0,115200n8"`. `net.ifnames=0 biosdevname=0`是内核参数,将网卡设备名固定为eth0..eth1等
   - autostart : 指定虚拟机是否在物理启动后自动启动

      其实就是在`/etc/libvirt/qemu/autostart`下创建执行该vm xml的软连接
   - print-xml : 如果虚拟机不需要安装过程(--import、--boot)，则显示生成的XML而不是创建此虚拟机. 默认情况下，此选项仍会创建磁盘映像
   - --dry-run：执行创建虚拟机的整个过程，但不真正创建虚拟机、改变主机上的设备配置信息及将其创建的需求通知给libvirt
   - **`--debug`**：显示debug信息
   - --connect=CONNCT选项来指定连接至一个非默认的hypervisor
      
      - `qemu:///system` : If running on a bare metal kernel as root (needed for KVM installs)
      - `qemu:///session` : If running on a bare metal kernel as non-root
   - --metadata

      - uuid : 虚拟机的唯一编号. 如果没有指定，将生成一个随机UUID
   - --sound

      - none : 没有声卡
   - noautoconsole: 不自动连接到guest console. 即不阻塞virt-install
   - --check : check开关

      - `disk_size=off` : 不检查disk size
      - `all=off` : 全部不检查

## [vm生命周期及状态转换](https://docs.openeuler.org/zh/docs/22.03_LTS/docs/Virtualization/%E7%AE%A1%E7%90%86%E8%99%9A%E6%8B%9F%E6%9C%BA.html)
虚拟机主要有如下几种状态：

- 未定义（undefined）：虚拟机未定义或未创建，即libvirt认为该虚拟机不存在
- 关闭状态（shut off）：虚拟机已经被定义但未运行，或者虚拟机被终止
- 运行中（running）：虚拟机处于运行状态
- 暂停（paused）：虚拟机运行被挂起，其运行状态被临时保存在内存中，可以恢复到运行状态

   可用被destroy
- 保存（saved）：与暂停（paused）状态类似，其运行状态被保存在持久性存储介质中，可以恢复到运行状态
- 崩溃（crashed）：通常是由于内部错误导致虚拟机崩溃，不可恢复到运行状态

![状态转换](https://docs.openeuler.org/zh/docs/22.03_LTS/docs/Virtualization/figures/status-transition-diagram.png)

在同一个主机上，每个domain具有唯一标识，通过虚拟机名称Name、UUID、Id表示:
- Name : 虚拟机名称
- UUID : 通用唯一识别码
- Id : 虚拟机运行标识, 关闭状态的虚拟机无此标识

## 网络
KVM虚拟化支持Linux网桥、Open vSwitch网桥等多种类型的网桥. 数据传输路径为`虚拟机 -> 虚拟网卡设备 -> Linux网桥或Open vSwitch网桥 -> 物理网卡`.

搭建linux bridge:
```bash
yum install bridge-utils
brctl addbr br0
brctl addif br0 eth0 # 将物理网卡eth0绑定到Linux网桥
ifconfig eth0 0.0.0.0 # eth0与网桥连接后，不再需要IP地址，将eth0的IP设置为0.0.0.0
dhclient br0 / ifconfig br0 192.168.1.2 netmask 255.255.255.0 # 设置br0的IP地址. 如果有DHCP服务器，可以通过dhclient设置动态IP地址; 
如果没有DHCP服务器，给br0配置静态IP
```

搭建open vswitch bridge见[这里](https://docs.openeuler.org/zh/docs/22.03_LTS/docs/Virtualization/%E5%87%86%E5%A4%87%E4%BD%BF%E7%94%A8%E7%8E%AF%E5%A2%83.html).

## xml
- backingStore

    <backingStore> 元素通常用于表示虚拟磁盘的后端存储（Backing Store）. 后端存储是一种虚拟磁盘的特性，它指的是虚拟磁盘的底层数据来源或依赖关系.

   ```xml
   <disk type='file' device='disk'>
     <driver name='qemu' type='qcow2'/>
     <source file='/path/to/vm_disk.qcow2'/>
     <backingStore type='file' index='1'>
       <format type='qcow2'/>
       <source file='/path/to/base_image.qcow2'/>
     </backingStore>
   </disk>
   ```

   在这个例子中, vm_disk.qcow2 是虚拟机的磁盘镜像，它有一个后端存储（<backingStore>）指向 base_image.qcow2. 这表示 vm_disk.qcow2 是基于 base_image.qcow2 创建的，两者之间存在一种依赖关系.

## virtio
ref:
- [Chapter 10. KVM Paravirtualized (virtio) Drivers](https://docs.redhat.com/zh_hans/documentation/red_hat_enterprise_linux/6/html/virtualization_host_configuration_and_guest_installation_guide/chap-virtualization_host_configuration_and_guest_installation_guide-para_virtualized_drivers#chap-Virtualization_Host_Configuration_and_Guest_Installation_Guide-Para_virtualized_drivers)
- [KVM Paravirtualized (virtio) Drivers](https://access.redhat.com/articles/2488201)
- [Certified Guest Operating Systems in Red Hat OpenStack Platform, Red Hat Virtualization, Red Hat OpenShift Virtualization and Red Hat Enterprise Linux with KVM](https://access.redhat.com/zh_CN/articles/7065625)
- [优化 Windows 虚拟机](https://docs.redhat.com/zh_hans/documentation/red_hat_enterprise_linux/9/html/configuring_and_managing_virtualization/optimizing-windows-virtual-machines-on-rhel_installing-and-managing-windows-virtual-machines-on-rhel#installing-kvm-drivers-on-a-host-machine_installing-kvm-paravirtualized-drivers-for-rhel-virtual-machines)
- [virtio-win](https://fedorapeople.org/groups/virt/virtio-win/direct-downloads/archive-virtio/)/[Windows VirtIO Drivers](https://pve.proxmox.com/wiki/Windows_VirtIO_Drivers)/[Windows guests - build ISOs including VirtIO drivers](https://pve.proxmox.com/wiki/Windows_guests_-_build_ISOs_including_VirtIO_drivers)

virtio-win iso 与系统版本的对应关系(from chatgpt, 未找到其他信息来源):
- virtio-win-0.1.x 到 virtio-win-0.1.171：

   早期版本，主要支持 Windows XP、Windows Server 2003、Windows Vista、Windows 7、Windows Server 2008 和 Windows Server 2008 R2
- virtio-win-0.1.173 到 virtio-win-0.1.180：

   这些版本开始支持 Windows 8 和 Windows Server 2012
- virtio-win-0.1.181 到 virtio-win-0.1.209：

   支持 Windows 8.1 和 Windows Server 2012 R2
- virtio-win-0.1.211 到 virtio-win-0.1.221：

   支持 Windows 10 和 Windows Server 2016
- virtio-win-0.1.223 及以上版本：

   提供对 Windows 10 的更新支持，并支持 Windows Server 2019 和 Windows Server 2022

### 获取支持的网卡
参考virt-manager的[default_model](https://github.com/virt-manager/virt-manager/blob/main/virtinst/devices/interface.py#L328)或[interface_recommended_models](https://github.com/virt-manager/virt-manager/blob/main/virtManager/addhardware.py#L592)

### disk type 'virtio' of 'vdb' does not support ejectable media'
[cdrom不支持virtio](https://lists.libvirt.org/archives/list/devel@lists.libvirt.org/message/4U6V62GKYPOCBVY5B3KM5JAP4RVLUCTZ/)

### arm64新建kvm uefi boot manager无法识别到光驱(bus=scsi)和网卡
正常uefi的Boot Manager如果存在光驱和网卡, 那么会识别到`UEFI QEMU QEMU CD-ROM`, `UEFI PXEv4 (MAC:<mac>)`和`UEFI PXEv6 (MAC:<mac>)`

再新建一个vm, 却能正常识别光驱. 用好的vm的nvram覆盖问题vm的, 问题依旧.

后来对比xml, 发现问题xml `controller type=scsi model=lsilogic`, 而正常xml是`controller type=scsi model=virtio-scsi`, 推测应该是当初选错了os(即xml libosinfo, **后来修正过libosinfo但xml controller model应该是创建时就固定了**), 导致使用了错误的scsi控制器而uefi无法识别该控制器.

### amd64新建kvm uefi boot manager无法识别到磁盘(bus=virtio)
disk file(raw)是原其他非uefi虚拟机的系统盘, 换bios引导即可

### arm64新建kvm+uefi 无法先从光驱(bus=sata)启动, 变成了从PXE启动
换scsi后正常

> 不知是没sata controller还是根本不支持, 未验证

### disk添加backingStore
```xml
<disk type='file' device='disk'>
  <driver name='qemu' type='qcow2'/>
  <source file='/tmp/pull4.qcow2'/>
   <backingStore type='file'>
      <format type='qcow2'/>
      <source file='/tmp/pull0.qcow2'/>
      <backingStore/> # 不能丢, 表示backing结束
   </backingStore>
  <target dev='vda' bus='virtio'/>
  <address type='pci' domain='0x0000' bus='0x00' slot='0x0a' function='0x0'/>
</disk>
```

### virt-manager vm界面黑屏
排查:
1. 防火墙

   可用[remote-viewer vnc://testguest2:5900](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/7/html/virtualization_deployment_and_administration_guide/sect-graphic_user_interface_tools_for_guest_virtual_machine_management-remote_viewer#sect-Graphic_User_Interface_tools_for_guest_virtual_machine_management-remote_viewer)测试

### 接管vm进图形界面黑屏
原机是gnome桌面.

解决:
1. 在单用户模式中换成字符界面

   `ln -s /usr/lib/systemd/system/graphical.target /etc/systemd/system/default.target`

### 加密磁盘
ref:
- [聊聊Kvm Qcow2和Ceph Rbd虚拟机磁盘加密](https://www.51cto.com/article/710126.html)
- [libvirt-qemu-磁盘加密之一：qcow2](https://blog.csdn.net/isclouder/article/details/79191665)

在libvirt 4.5版本之前，除了luks加密之外，还支持qcow加密的. 在 QEMU 中使用 qcow 加密卷在 QEMU 2.3 中开始逐步淘汰.

### `channel type="unix"`
- host socket在channel.source.path

   每次关闭会导致channel.source.path变化
- vm 内部则是串口在/dev/virtio-ports/<channel.target.name>

   如果vm app打开了这个serial, 则channel.target.state显示为connected

如果vm内部安装并启用了qemu-guest-agent(vm xml需要指定channel.target.name=org.qemu.guest_agent.0, 因为qemu-guest-agent.service即qemu-ga进程默认使用了它), 在host侧通过`virsh qemu-agent-command centos '<cmd>'`即可操作vm, 比如`{"execute":"guest-info"}`

> 其他工具: `echo '{"execute": "guest-ping"}' | socat - UNIX-CONNECT:/tmp/qemu-serial.sock`

### 添加virtio 磁盘报"No more available PCI slots"
ref:
- [PCI topology and hotplug](https://libvirt.org/pci-hotplug.html)
- PCI 地址计算方法见 [Device Addresses](https://libvirt.org/formatdomain.html#device-addresses) 和 [PCI addresses in domain XML and guest OS](https://libvirt.org/pci-addresses.html)

PCI 插槽不够了，没法挂载新设备, 所以需要新增 PCI 插槽

解决方法(未验证):
1. 添加`<controller type='pci' index='10' model='pcie-root-port' />`, 注意index不能重复
2. 添加`<controller type='pci' model='pci-bridge'/>`

### 兆芯KH-40000/KX-U6780A + oracle linux 7.9没有/dev/kvm
ref:
- [开胜® KH-40000系列处理器](https://www.zhaoxin.com/prod_view.aspx?nid=3&typeid=582&id=2411)

   cpu本身支持虚拟化

现象:
1. `egrep '(vmx|svm)' /proc/cpuinfo`找到vmx
2. `dmesg |grep -i vmx`提示`kvm: VMX is disabled on CPU 0`, 但百敖uefi已确认开启了虚拟化(即现象1)
3. `lsmod |grep kvm`有发现kvm模块, 但没有kvm_xxx

`modprobe -v kvm_intel`报`could not inert 'kvm_intel': Input/output error`, 系统日志报`kvm: VMX is disabled on CPU 0`

### vm启动报"value format, found hv_relaxed"
ref:
- [Hyper-V Enlightenments](https://www.qemu.org/docs/master/system/i386/hyperv.html)

vm(host是arm64)从光驱启动时使用了windows x64 iso

