# ceph
存储根据其类型，可分为块存储，对象存储和文件存储. 在主流的分布式存储技术中，HDFS/GPFS/GFS属于文件存储，Swift属于对象存储，而Ceph可支持块存储、对象存储和文件存储，故称为统一存储.

![Ceph的主要架构](/misc/img/ceph/1f950c8b30fe6992437242c368f0f8b1.png)

架构:
1. Ceph的最底层是RADOS（分布式对象存储系统），它具有可靠、智能、分布式等特性，实现高可靠、高可拓展、高性能、高自动化等功能，并最终存储用户数据. RADOS系统主要由两部分组成，分别是OSD和Monitor
1. RADOS之上是LIBRADOS，LIBRADOS是一个库，它允许应用程序通过访问该库来与RADOS系统进行交互，支持多种编程语言，比如C、C++、Python等
1. 基于LIBRADOS层开发的有三种接口，分别是RADOSGW、librbd和MDS
	
	1. RADOSGW是一套基于当前流行的RESTFUL协议的网关，支持对象存储，兼容S3和Swift
	1. librbd提供分布式的块存储设备接口，支持块存储
	1. MDS提供兼容POSIX的文件系统，支持文件存储

Ceph的核心组件包括Client客户端、MON监控服务、MDS元数据服务、OSD存储服务，各组件功能如下：

1. MON监控服务：负责监控整个集群，维护集群的健康状态，维护展示集群状态的各种图表，如OSD Map、Monitor Map、PG Map和CRUSH Map
1. MDS元数据服务：负责保存文件系统的元数据，管理目录结构
1. OSD存储服务：主要功能是存储数据、复制数据、平衡数据、恢复数据，以及与其它OSD间进行心跳检查等. 一般情况下一块硬盘对应一个OSD

## 组件
- Client : 负责存储协议的接入，节点负载均衡
- Monitors: Ceph 监视器(ceph-mon) 

    维护着集群状态的各种运行图，包括如OSD Map、Monitor Map、PG Map和CRUSH Map，这些运行图都是很要紧的集群状态，对于各种 Ceph 守护进程的相互协作必不可少. 监视器还负责管理守护进程和客户端之间的认证. 考虑到冗余性和高可用性，一般都要求至少有三个监视器.
- Managers: Ceph 管理器守护进程（ ceph-mgr ）

    负责持续跟踪运行时指标和 Ceph 当前的状态，包括存储利用率、当前的性能指标、和系统负载. Ceph 管理器守护进程还托管着基于 python 的插件，用于管理和展示 Ceph 集群信息，包括一个基于网页的 Ceph 管理器仪表盘 和 REST API. 为保障高可用性，一般要求至少有两个管理器.
- Ceph OSDs: Ceph OSD （对象存储守护进程， ceph-osd ）

    负责存储数据、处理数据复制、恢复、重均衡、以及向 Ceph 监视器和管理器提供些监控信息，如检查其它 Ceph OSD 守护进程的心跳. 一般情况下一块硬盘对应一个OSD. 因此为保障冗余性和高可用性，一般需要至少 3 个 Ceph OSD.
- MDSs: Ceph 元数据服务器（ MDS ， ceph-mds ）, Ceph 文件系统客户端依赖它.

    为 Ceph 文件系统存储元数据（也就是说， Ceph 块设备和 Ceph 对象存储不使用 MDS ）. 元数据服务器有益于 POSIX 文件系统用户执行基本命令（像 ls 、 find 等等），避免了给 Ceph 存储集群增加过重的负担.

Ceph 把数据保存为逻辑存储池内的对象. 根据 CRUSH 算法， Ceph 可计算出哪个归置组应该持有指定对象，然后进一步计算出哪个 OSD 守护进程持有归置组，正因为有了 CRUSH 算法， Ceph 存储集群才具备动态伸缩、重均衡和动态恢复功能.

> 新版ceph的cephfs基于bluestore, 不再依赖其他文件系统.

## FAQ
### DRBD vs Ceph
CEPH 是一种开源软件，旨在在统一系统中提供高度可扩展的对象，块和基于文件的存储.

DRBD 是一种基于软件和网络(tcp/ip和RDMA)的块复制存储解决方案, 基于块设备, 相对而言, 其在ceph的块设备之上.