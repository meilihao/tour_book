# ceph
Ceph是一种为优秀的性能、可靠性和可扩展性而设计的统一的、分布式的存储系统.

存储根据其类型，可分为块存储，对象存储和文件存储. 在主流的分布式存储技术中，HDFS/GPFS/GFS属于文件存储，Swift属于对象存储，而Ceph可支持块存储、对象存储和文件存储，故称为统一存储.

![Ceph的主要架构](/misc/img/ceph/Image00152.jpg)

架构:
1. Ceph的最底层和核心是RADOS（分布式对象存储系统），它具有可靠、智能(自动化)、分布式等特性，实现高可靠、高可拓展、高性能、高自动化等功能，并最终提供了一个可无限扩展的存储集群. 

	RADOS由大量的存储设备节点组成，每个节点拥有自己的硬件资源（CPU、内存、硬盘、网络），并运行着操作系统和文件系统.

	ceph采用具有计算能力的设备（比如普通的服务器）作为存储系统的存储节点, 是为了充分发挥存储设备自身的计算能力.

	RADOS集群主要由两种节点组成：为数众多的OSD，负责完成数据存储和维护；若干个Monitor，负责完成系统状态检测和维护. OSD和Monitor之间互相传递节点的状态信息，共同得出系统的总体运行状态，并保存在一个全局的数据结构中，即所谓的集群运行图（Cluster Map）里. 集群运行图与RADOS提供的特定算法相配合，便实现了Ceph的许多优秀特性.

	据定义，OSD可以被抽象为系统和守护进程（OSD Daemon）两个部分:
	1. OSD的系统部分本质上就是一台安装了操作系统和文件系统的计算机
	1. 每个OSD占用一定的计算能力、一定数量的内存, 一块硬盘（在通常情况下一个OSD对应一块硬盘）, 及足够的网络带宽.

		每个OSD拥有一个自己的OSD Daemon. 这个Daemon负责完成OSD的所有逻辑功能，包括与Monitor和其他OSD（事实上是其他OSD的Daemon）通信，以维护及更新系统状态，与其他OSD共同完成数据的存储和维护操作，与客户端通信完成各种数据对象操作，等等.

	RADOS集群从Ceph客户端接收数据（无论是来自Ceph块设备、Ceph对象存储、Ceph文件系统，还是基于librados的自定义实现），然后存储为对象. 每个对象是文件系统中的一个文件，它们存储在OSD的存储设备上，由OSD Daemon处理存储设备上的读/写操作.

	OSD在扁平的命名空间内把所有数据存储为对象（也就是没有目录层次）. 对象包含一个标识符、二进制数和由名/值对组成的元数据，元数据语义完全取决于Ceph客户端.

1. RADOS之上是LIBRADOS，LIBRADOS是一个库，它允许应用程序通过访问该库来与RADOS系统进行交互，支持多种编程语言，比如C、C++、Python等

	librados库实际上是对RADOS进行抽象和封装，并向上层提供API的，以便可以基于RADOS（而不是整个Ceph）进行应用开发. 特别要注意的是，RADOS是一个对象存储系统，因此，librados库实现的API也**只是针对对象存储功能**的.
1. 基于LIBRADOS层开发的有三种接口，分别是RADOSGW、RBD与Ceph FS
	
	1. RADOSGW是一套基于当前流行的RESTFUL协议的网关，支持对象存储，兼容S3和Swift
	1. RBD提供分布式的块存储设备接口，支持块存储. 常用于在虚拟化的场景下为虚拟机创建存储卷。Red Hat已经将RBD驱动集成在KVM/QEMU中，以提高虚拟机的访问性能.
	1. Ceph FS提供兼容POSIX的文件系统，支持文件存储

	因此严格意义上讲，Ceph只提供对象存储接口，所谓的块存储接口和文件系统存储接口都算是对象存储接口应用程序.

## 组件
- Client : 负责存储协议的接入，节点负载均衡
- Monitors: Ceph 监视器(ceph-mon) 

    负责系统状态检测和维护. 它维护着集群状态的各种运行图，包括如OSD Map、MDS Map、Monitor Map、PG Map和CRUSH Map，这些运行图都是很要紧的集群状态，对于各种 Ceph 守护进程的相互协作必不可少. 监视器还负责管理守护进程和客户端之间的认证. 考虑到冗余性和高可用性，一般都要求至少有三个监视器.
- Managers: Ceph 管理器守护进程（ ceph-mgr ）

    负责持续跟踪运行时指标和 Ceph 当前的状态，包括存储利用率、当前的性能指标、和系统负载. Ceph 管理器守护进程还托管着基于 python 的插件，用于管理和展示 Ceph 集群信息，包括一个基于网页的 Ceph 管理器仪表盘 和 REST API. 为保障高可用性，一般要求至少有两个管理器.
- Ceph OSDs: Ceph OSD （对象存储守护进程， ceph-osd ）

    负责数据存储和维护(数据复制、恢复、重均衡, 以及向 Ceph 监视器和管理器提供些监控信息，如检查其它 Ceph OSD 守护进程的心跳等). 一般情况下一块硬盘对应一个OSD. 因此为保障冗余性和高可用性，一般需要至少 3 个 Ceph OSD.
- MDSs: Ceph 元数据服务器（ MDS ， ceph-mds ）, Ceph 文件系统客户端依赖它.

    为 Ceph 文件系统存储元数据（也就是说， Ceph 块设备和 Ceph 对象存储不使用 MDS ）比如目录结构. 元数据服务器有益于 POSIX 文件系统用户执行基本命令（像 ls 、 find 等等），避免了给 Ceph 存储集群增加过重的负担.

Ceph 把数据保存为逻辑存储池内的对象. 根据 CRUSH 算法， Ceph 可计算出哪个归置组应该持有指定对象，然后进一步计算出哪个 OSD 守护进程持有归置组，正因为有了 CRUSH 算法， Ceph 存储集群才具备动态伸缩、重均衡和动态恢复功能.

> 新版ceph的cephfs基于bluestore, 不再依赖其他文件系统.

## OSD
OSD的状态直接影响数据的重新分配，所以监测OSD的状态是Monitor的主要工作之一.

OSD状态用两个维度表示：up或down（OSD Daemon与Monitor连接是否正常）；in或out（OSD是否含有PG）. 因此，对于任意一个OSD，共有4种可能的状态:
- up & out：OSD Daemon与Monitor通信正常，但是没有PG分配到该OSD上。这种状态一般是OSD Daemon刚刚启动时
- up & in：OSD Daemon工作的正常状态，有PG分配到OSD上
- down & in：OSD Daemon不能与Monitor或其他OSD进行正常通信，这可能是因为网络中断或Daemon进程意外退出
- down & out：OSD无法恢复，Monitor决定将OSD上的PG进行重新分配。之所以会出现该状态，是考虑OSD可能会在短时间内恢复，尽量减少数据的再分配

OSD状态是通过心跳（Heartbeat）检测的:
- Peer OSD之间的心跳包

	Peer OSD是指该OSD上所有PG的副本所在的OSD. 同时由于Ceph提供公众网络（Public Network）（OSD与客户端通信）和集群网络（Cluster Network）（OSD之间的通信），所以Peer OSD之间的心跳包也分为前端（公众网络）和后端（集群网络），这样可最大限度地监测OSD及公众网络和集群网络的状态，及时上报Monitor. 同时考虑到网络的抖动问题，可以设置Monitor在决定OSD下线之前需要收到多少次的报告.
- OSD与Monitor之间的心跳包

	这个心跳包可以看作是Peer OSD之间心跳包的补充. 如果OSD不能与其他OSD交换心跳包，那么就必须与Monitor按照一定频率进行通信，比如OSD状态是up & out时就需要这种心跳包

### Ceph寻址流程
![](/misc/img/ceph/Image00156.jpg)

Ceph寻址流程涉及的概念:
- File：此处的File就是用户需要存储或访问的文件. 对于一个基于Ceph开发的对象存储应用而言，这个File也就对应于应用中的“对象”，也就是用户直接操作的“对象”.
- Object：此处的Object是RADOS所看到的“对象”. Object与File的区别是，Object的最大尺寸由RADOS限定（通常为2MB或4MB），以便实现底层存储的组织管理. 因此，当上层应用向RADOS存入尺寸很大的File时，需要将File切分成统一大小的一系列Object（最后一个的大小可以不同）进行存储.
- PG（Placement Group）：顾名思义，PG的用途是对Object的存储进行组织和位置映射的. 具体而言，一个PG负责组织若干个Object（可以为数千个甚至更多），但一个Object只能被映射到一个PG中，即PG和Object之间是“一对多”的映射关系. 同时，一个PG会被映射到n 个OSD上，而每个OSD上都会承载大量的PG，即PG和OSD之间是“多对多”的映射关系. 在实践当中，n 至少为2，如果用于生产环境，则至少为3. 一个OSD上的PG可达到数百个. 事实上，PG数量的设置关系到数据分布的均匀性问题.
- OSD：OSD的数量事实上也关系到系统的数据分布均匀性，因此不应该太少. 在实践当中，至少也应该是数百个的量级才有助于Ceph系统发挥其应有的优势.

具体映射:
1. File→Object映射
	
	目的: 将用户要操作的File映射为RADOS能够处理的Object，其十分简单，本质上就是按照Object的最大尺寸对File进行切分，相当于磁盘阵列中的条带化过程.

	这种切分的好处有两个：一是让大小不限的File变成具有一致的最大尺寸、可以被RADOS高效管理的Object；二是让对单一File实施的串行处理变为对多个Object实施的并行化处理.

	每一个切分后产生的Object将获得唯一的oid，即Object ID.

	ino是待操作File的元数据，可以简单理解为该File的唯一ID, 因此ino的唯一性必须得到保证. ono则是由该File切分产生的某个Object的序号. 而oid就是将这个序号简单连缀在该File ID之后得到的.
1. Object → PG映射

	公式: `hash(oid) & mask -> pgid`

	首先，使用Ceph系统指定的一个静态哈希算法计算oid的哈希值，将oid映射为一个近似均匀分布的伪随机值. 然后，将这个伪随机值和mask按位相与，得到最终的PG序号（pgid）.

	根据RADOS的设计，给定PG的总数为m （m 应该为2的整数幂），则mask的值为m- 1. 因此，哈希值计算和按位与操作的整体结果事实上是从所有m 个PG中近似均匀地随机选择1个. 基于这一机制，当有大量Object和大量PG时，RADOS能够保证Object和PG之间的近似均匀映射. 又因为Object是由File切分而来的，大部分Object的尺寸相同，因此，这一映射最终保证了各个PG中存储的Object的总数据量近似均匀.

	这里反复强调了“大量”，意思是只有当Object和PG的数量较多时，这种伪随机关系的近似均匀性才能成立，Ceph的数据存储均匀性才有保证. 为保证“大量”的成立，一方面，Object的最大尺寸应该被合理配置，以使得同样数量的File能够被切分成更多的Object；另一方面，Ceph也推荐PG总数应该为OSD总数的数百倍，以保证有足够数量的PG可供映射.
1. PG → OSD映射

	将作为Object的逻辑组织单元的PG映射到数据的实际存储单元OSD上. RADOS采用一个名为CRUSH的算法，将pgid代入其中，然后得到一组共n 个OSD. 这n 个OSD共同负责存储和维护一个PG中的所有Object. n 的数值可以根据实际应用中对于可靠性的需求而配置，在生产环境下通常为3. 具体到每个OSD，则由其上运行的OSD Daemon负责执行映射到本地的Object在本地文件系统中的存储、访问、元数据维护等操作.

	和“Object → PG”映射中采用的哈希算法不同，CRUSH算法的结果不是绝对不变的，而会受到其他因素的影响。其影响因素主要有两个:
	1. 当前系统状态(集群运行图). 当系统中的OSD状态、数量发生变化时，集群运行图也可能发生变化，而这种变化将会影响到PG与OSD之间的映射关系
	1. 存储策略配置. 这里的策略主要与安全相关. 利用策略配置，系统管理员可以指定承载同一个PG的3个OSD分别位于数据中心的不同服务器或机架上，从而进一步改善存储的可靠性.

	Ceph正是利用了CRUSH算法的动态特性，可以将一个PG根据需要动态迁移到不同的OSD组合上，从而自动化地实现高可靠性、数据分布再平衡等特性.

从整个过程可以看到，这里没有任何的全局性查表操作需求. 至于唯一的全局性数据结构：集群运行图, 它的维护和操作都是轻量级的，不会对系统的可扩展性、性能等因素造成影响.

引入PG的好处至少有两方面：
1. 一方面实现了Object和OSD之间的动态映射，从而为Ceph的可靠性、自动化等特性的实现留下了空间

	如果Object直接映射到一组OSD上, 这种算法是某种固定映射的哈希算法, 那么osd损坏或新增osd, object没法迁移或再平衡到新osd.
1. 另一方面也有效简化了数据的存储组织，大大降低了系统的维护与管理成本

	如果Object直接映射到一组OSD上, 这种算法是某种动态算法, 比如仍然采用CRUSH算法. 在Ceph的现有机制中，一个OSD平时需要和与其共同承载同一个PG的其他OSD交换信息，以确定各自是否工作正常，是否需要进行维护操作. 如果没有pg, 则一个OSD需要和与其共同承载同一个Object的其他OSD交换信息, 由于每个OSD上承载的Object可能高达数百万个, 因此，同样长度的一段时间内，一个OSD大约需要进行的OSD间信息交换将暴涨至数百万次乃至数千万次. 而这种状态维护成本显然过高.

这种分层或分级的设计思路在很多复杂系统的寻址问题上都有应用，比如操作系统里的内存管理多级页表的使用，英特尔MPX（Memory Protection Extensions）技术里引入的Bound Directory等.

### 存储池
ceph存储池是一个逻辑概念，是对存储对象的逻辑分区. Ceph安装后，会有一个默认的存储池，用户也可以自己创建新的存储池. 一个存储池包含若干个PG及其所存储的若干个对象.

Ceph客户端从监视器获取一张集群运行图，并把对象写入存储池. 存储池的大小或副本数、CRUSH存储规则和归置组数量决定Ceph如何放置数据. ceph中通过`ceph osd pool crate`创建存储池.

创建存储池命令支持的参数如下:
- 设置数据存储的方法属于多副本模式还是纠删码模式. 如果是多副本模式，则可以设置副本的数量；如果是纠删码模式，则可以设置数据块和非数据块的数量（纠删码存储池把各对象存储为K +M 个数据块，其中有K 个数据块和M个编码块）. 默认为多副本模式（即存储每个对象的若干个副本），如果副本数为3，则每个PG映射到3个OSD节点上
- 设置PG的数目. 合理设置PG的数目，可以使资源得到较优的均衡
- 设置PGP的数目. 在通常情况下，与PG数目一致. 当需要增加PG数目时，用户数据不会发生迁移，只有进一步增加PGP数目时，用户数据才会开始迁移
- 针对不同的存储池设置不同的CRUSH存储规则. 比如可以创建规则，指定在选择OSD时，选择拥有固态硬盘的OSD节点

另外，通过存储池，还可以进行如下操作:
1. 提供针对存储池的功能，如存储池快照等
1. 设置对象的所有者或访问权限

PGP是存储池PG的OSD分布组合个数. PG数目的增加会引起PG的分裂，新的PG仍然在原来的OSD上，而PGP数目的增加则会引起部分PG的分布发生变化, 但是不会引起PG内对象的变动. 可参考[ceph分布式存储-PG和PGP的区别](https://cloud.tencent.com/developer/article/1664635)

## monitor
由若干个Monitor组成的监视器集群共同负责整个Ceph集群中所有OSD状态的发现与记录，并且形成集群运行图的主副本，包括集群成员、状态、变更，以及Ceph存储集群的整体健康状况. 该集群运行图会被扩散至全体OSD及客户端，OSD使用集群运行图进行数据的维护，而客户端使用集群运行图进行数据的寻址.

Ceph客户端读或写数据前必须先连接到某个Ceph监视器上，获得最新的集群运行图副本.

在集群中，各个Monitor的功能总体上是一样的，其之间的关系可以被简单理解为主从备份关系. Monitor并不主动轮询各个OSD的当前状态, 正相反，OSD需要向Monitor上报状态信息. 常见的上报有两种情况：一是新的OSD被加入集群，二是某个OSD发现自身或其他OSD发生异常. 在收到这些上报信息后，Monitor将更新集群运行图的信息并加以扩散.

集群运行图实际上是多个Map的统称，包括Monitor Map、OSDMap、PG Map、CRUSH Map及MDS Map等，各运行图维护着各自运行状态的变更. 其中CRUSH Map用于定义如何选择OSD，内容包含了存储设备列表、故障域树状结构（设备的分组信息，如设备、主机、机架、行、房间等）和存储数据时如何利用此树状结构的规则.

crush map中所有非叶子节点称为桶（Bucket），所有Bucket的ID号都是负数，和OSD的ID进行区分. 选择OSD时，需要先指定一个Bucket，然后选择它的一个子Bucket，这样一级一级递归，直到到达设备（叶子）节点. 目前有5种算法来实现子节点的选择，包括Uniform、List、Tree、Straw、Straw2. 这些算法的选择影响了两个方面的复杂度：在一个Bucket中，找到对应的节点的复杂度及当一个Bucket中的OSD节点丢失或增加时，数据移动的复杂度.

![不同Bucket算法复杂度比较](misc/img/io/Image00166.jpg)

其中，Uniform与item具有相同的权重，而且Bucket很少出现添加或删除item的情况，它的查找速度是最快的. Straw/Straw2不像List和Tree一样都需要遍历，而是让Bucket包含的所有item公平竞争. 这种算法就像抽签一样，所有的item都有机会被抽中（只有最长的签才能被抽中，每个签的长度与权重有关）.

除了存储设备的列表及树状结构，CRUSH Map还包含了存储规则，用来指定在每个存储池中选择特定OSD的Bucket范围，还可以指定备份的分布规则. CRUSH Map有一个默认存储规则，如果用户创建存储池时没有指定CRUSH规则，则使用该默认规则, 但是用户可以自定义规则，指定给特定存储池.

### 与monitor通信
1. Monitor与客户端的通信

客户端包括RBD客户端、RADOS客户端、Ceph FS客户端/MDS. 根据通信内容分为获取OSDMap和命令行操作:
- 命令行操作

命令行操作主要包括集群操作命令（OSD、Monitor、MDS的添加和删除，存储池的创建和删除等）、集群信息查询命令（集群状态、空间利用率、IOps和带宽等）. 这些命令都是由Monitor直接执行或通过Monitor转发到OSD上执行的.

- 获取OSDMap

客户端与RADOS的读/写不需要Monitor的干预，客户端通过哈希算法得到Object所在的PG信息，然后查询OSDMap就可以得到PG的分布信息，就可以与Primary OSD进行通信了. 那么客户端与Monitor仅仅是当需要获取最新OSDMap时才会进行通信.

	- 客户端初始化时。
	- 某些特殊情况，会主动获取新的OSDMap：OSDMap设置了 CEPH_OSDMAP_PAUSEWR/PAUSERD（Cluster暂停所有读/写），每一次的读/写都需要获取OSDMap；OSDMap设置了Cluster空间已满或存储池空间已满，每一次写都需要获取OSDMap；找不到相应的存储池或通过哈希算法得到PG，但是在OSDMap中查不到相关PG分布式信息（说明PG删除或PG创建）.

1. Monitor与OSD的通信

相比Monitor与客户端的通信，Monitor与OSD的通信会复杂得多，内容如下:
- Monitor需要知道OSD的状态，并根据状态生成新的OSDMap. 所以OSD需要将OSD的Down状态向Monitor报告.
- OSD和Monitor之间存在心跳机制，通过这种方式来判断OSD的状态.
- OSD定时将PG信息发送给Monitor。PG信息包括PG的状态（Active、degraded等）、Object信息（数目、大小、复制信息、Scrub/Repair信息、IOps和带宽等）。Monitor通过汇总这些信息就可以知道整个系统的空间使用率、各个存储池的空间大小、集群的IOps和带宽等实时信息.
- OSD的操作命令是客户端通过Monitor传递给OSD的。比如osd scrub/deep scrub、pg scrub/deep scrub等.
- OSD初始化或Client/Primary OSD所包含的OSDMap的版本高于当前的OSDMap

## 数据操作流程
Ceph的读/写操作采用Primary-Replica模型，客户端只向Object所对应OSD set的Primary发起读/写请求，这保证了数据的强一致性. 当Primary收到Object的写请求时，它负责把数据发送给其他副本，只有这个数据被保存在所有的OSD上时，Primary才应答Object的写请求，这保证了副本的一致性.

![Object写入流程](/misc/img/ceph/Image00168.jpg)

当某个客户端需要向Ceph集群写入一个File时，首先需要在本地完成寻址流程，将File变为一个Object，然后找出存储该Object的一组共3个OSD. 这3个OSD具有各自不同的序号，序号最靠前的那个OSD就是这一组中的Primary OSD，而后两个则依次是Secondary OSD和Tertiary OSD.

找出3个OSD后，客户端将直接和Primary OSD进行通信，发起写入操作（步骤1）. Primary OSD收到请求后，分别向Secondary OSD和Tertiary OSD发起写入操作（步骤2和步骤3）. 当Secondary OSD和Tertiary OSD各自完成写入操作后，将分别向Primary OSD发送确认信息（步骤4和步骤5）. 当Primary OSD确认其他两个OSD的写入完成后，则自己也完成数据写入，并向客户端确认Object写入操作完成（步骤6）.

之所以采用这样的写入流程，本质上是为了保证写入过程中的可靠性，尽可能避免出现数据丢失的情况. 同时，由于客户端只需要向Primary OSD发送数据，因此，在互联网使用场景下的外网带宽和整体访问延迟又得到了一定程度的优化.

当然，这种可靠性机制必然导致较长的延迟，特别是，如果等到所有的OSD都将数据写入磁盘后再向客户端发送确认信号，则整体延迟可能难以忍受。因此，Ceph可以分两次向客户端进行确认。当各个OSD都将数据写入内存缓冲区后，就先向客户端发送一次确认，此时客户端即可以向下执行。待各个OSD都将数据写入磁盘后，会向客户端发送一个最终确认信号，此时客户端可以根据需要删除本地数据。

分析上述流程可以看出，在正常情况下，客户端可以独立完成OSD寻址操作，而不必依赖于其他系统模块。因此，大量的客户端可以同时和大量的OSD进行并行操作。同时，如果一个File被切分成多个Object，这多个Object也可被并行发送至多个OSD上。

从OSD的角度来看，由于同一个OSD在不同的PG中的角色不同，因此，其工作压力也可以被尽可能均匀地分担，从而避免单个OSD变成性能瓶颈.

如果需要读取数据，客户端只需完成同样的寻址过程，并直接和Primary OSD联系. 在目前的Ceph设计中，被读取的数据默认由Primary OSD提供，但也可以设置允许从其他OSD中获取，以分散读取压力从而提高性能.

## Cache Tiering(存储分层技术)
Cache Tiering的理论基础，就是存储的数据是有热点的，数据并不是均匀访问的.

Cache Tiering的做法就是，用固态硬盘等相对快速、昂贵的存储设备组成一个存储池作为缓存层存储热数据，然后用相对慢速、廉价的设备作为存储后端存储冷数据（Storage层或Base层）. 缓存层使用多副本模式，Storage层可以使用多副本或纠删码模式.

在Cache Tiering中有一个分层代理，当保存在缓存层的数据变冷或不再活跃时，该代理把这些数据刷到Storage层，然后把它们从缓存层中移除，这种操作称为刷新（Flush）或逐出（Evict）.

Ceph的对象管理器（Objecter，位于osdc即OSD客户端模块）决定往哪里存储对象，分层代理决定何时把缓存内的对象“刷回”Storage层，所以缓存层和Storage层对Ceph客户端来说是完全透明的. 需要注意的是，Cache Tiering是基于存储池的，在缓存层和Storage层之间的数据移动是两个存储池之间的数据移动.

![Cache Tiering](/misc/img/ceph/Image00169.jpg)

目前Cache Tiering主要支持如下几种模式:
- 写回模式：对于写操作，当请求到达缓存层完成写操作后，直接应答客户端，之后由缓存层的代理线程负责将数据写入Storage层。对于读操作则看是否命中缓存层，如果命中直接在缓存层读，没有命中可以重定向到Storage层访问，如果Object近期有访问过，说明比较热，可以提升到缓存层中。
- forward模式：所有的请求都重定向到Storage层访问。
- readonly模式：写请求直接重定向到Storage层访问，读请求命中缓存层则直接处理，没有命中缓存层需要从Storage层提升到缓存层中进而完成请求，下次再读取直接命中缓存。
- readforward模式：读请求都重定向到Storage层中，写请求采用写回模式。
- readproxy模式：读请求发送给缓存层，缓存层去Storage层中读取，获得Object后，缓存层自己并不保存，而是直接发送给客户端，写请求采用写回模式。
- proxy模式：对于读/写请求都是采用代理的方式，不是转发而是代表客户端去进行操作，缓存层自己并不保存。

这里提及的重定向、提升与代理等几种操作的具体含义如下:
- 重定向：客户端向缓存层发送请求，缓存层应答客户端发来的请求，并告诉客户端应该去请求Storage层，客户端收到应答后，再次发送请求给Storage层请求数据，并由Storage层告诉客户端请求的完成情况。
- 代理：客户端向缓存层发送读请求，如果未命中，则缓存层自己会发送请求给Storage层，然后由缓存层将获取的数据发送给客户端，完成读请求。在这个过程中，虽然缓存层读取到了该Object，但不会将其保存在缓存层中，下次仍然需要重新向Storage层请求。
- 提升：客户端向缓存层发送请求，如果缓存层未命中，则会选择将该Object从Storage层中提升到缓存层中，然后在缓存层进行读/写操作，操作完成后应答客户端请求完成。在这个过程中，和代理操作的区别是，在缓存层会缓存该Object，下次直接在缓存中进行处理。

## 块存储
基于RADOS与librados库，Ceph通过RBD提供了一个标准的块设备接口，提供基于块设备的访问模式.

Ceph中的块设备称为Image，是精简配置的，即按需分配，大小可调且将数据条带化存储到集群内的多个OSD中.

条带化是指把连续的信息分片存储于多个设备中. 当多个进程同时访问一个磁盘时，可能会出现磁盘冲突的问题. 大多数磁盘系统都对访问次数（每秒的I/O操作）和数据传输率（每秒传输的数据量，TPS）有限制，当达到这些限制时，后面需要访问磁盘的进程就需要等待，这时就是所谓的磁盘冲突. 避免磁盘冲突是优化I/O性能的一个重要目标，而优化I/O性能最有效的手段是将I/O请求最大限度地进行平衡.

条带化就是一种能够自动将I/O负载均衡到多个物理磁盘上的技术. 通过将一块连续的数据分成多个相同大小的部分，并把它们分别存储到不同的磁盘上，条带化技术能使多个进程同时访问数据的不同部分而不会造成磁盘冲突，而且能够获得最大限度上的I/O并行能力.

条带化能够将多个磁盘驱动器合并为一个卷，这个卷所能提供的速度比单个盘所能提供的速度要快很多. Ceph的块设备就对应于LVM的逻辑卷，块设备被创建时，可以指定如下参数实现条带化:
- stripe-unit：条带的大小
- stripe-count：在多少数量的对象之间进行条带化

![Ceph块设备条带化](/misc/img/ceph/Image00170.jpg)

如上图, 当stripe-count为3时，表示块设备上地址［0， object-size*stripe_count-1］到对象位置的映射. 每个对象被分成stripe_size大小的条带，按stripe_count分成一组，块设备在上面依次分布. 块设备上［0， stripe_size-1］对应Object1上的［0，stripe_size-1］，块设备上［stripe_size， 2*stripe_size-1］对应Object2上的［0， stripe_size-1］，以此类推.

当处理大尺寸图像、大Swift对象（如视频）的时候，就能看到条带化到一个对象集合（Object Set）中的多个对象能带来显著的读/写性能提升. 当客户端把条带单元并行地写入相应对象时，就会有明显的写性能提升，因为对象映射到了不同的PG，并进一步映射到不同OSD，可以并行地以最大速度写入. 由于到单一磁盘的写入受制于磁头移动（如6ms寻道时间）和存储设备带宽（如100MB/s），Ceph把写入分布到多个对象（它们映射到不同PG和OSD中），这样就减少了寻道次数，并利用了多个驱动器的吞吐量，以达到更高的读/写速度.

使用Ceph的块设备有两种路径:
- 通过Kernel Module：即创建了RBD设备后，把它映射到内核中，成为一个虚拟的块设备，这时这个块设备同其他通用块设备一样，设备文件一般为/dev/rbd0，后续直接使用这个块设备文件就可以了，可以把/dev/rbd0格式化后挂载到某个目录，也可以直接作为裸设备进行使用.
- 通过librbd：即创建了RBD设备后，使用librbd、librados库访问和管理块设备. 这种方式直接调用librbd提供的接口，实现对RBD设备的访问和管理，不会在客户端产生块设备文件.

其中第二种方式主要是为虚拟机提供块存储设备. 在虚拟机场景中，一般会用QEMU/KVM中的RBD驱动部署Ceph块设备，宿主机通过librbd向客户机提供块存储服务. QEMU可以直接通过librbd，像访问虚拟块设备一样访问Ceph块设备.

## Ceph FS
Ceph FS是一个可移植操作系统接口兼容的分布式存储系统，与通常的网络文件系统一样，要访问Ceph FS，需要有对应的客户端. Ceph FS支持两种客户端：Ceph FS FUSE和Ceph FS Kernel. 也就是说有两种使用Ceph FS的方式：一是通过Kernle Module，Linux内核里包含了Ceph FS的实现代码；二是通过FUSE（用户空间文件系统）的方式. 通过调用libcephfs库来实现Ceph FS的加载，而libcephfs库又调用librados库与RADOS进行通信.

之所以会通过FUSE的方式实现Ceph FS的加载，主要是考虑Kernel客户端的功能、稳定性、性能都与内核绑定，在不能升级内核的情况下，很多功能可能就不能使用, 而FUSE基本就不存在这个限制.

![Ceph FS架构](/misc/img/ceph/Image00172.jpg)

Ceph FS架构上图所示. 上层是支持客户端的Ceph FS Kernel Object、Ceph FS FUSE、Ceph FS Library等，底层还是基础的OSD和Monitor，此外添加了元数据服务器（MDS）. Ceph FS要求Ceph存储集群内至少有一个元数据服务器，负责提供文件系统元数据（目录、文件所有者、访问模式等）的存储与操作. MDS只为Ceph FS服务，如果不需要使用Ceph FS，则不需要配置MDS.

Ceph FS从数据中分离出元数据，并存储于MDS，文件数据存储于集群中的一个或多个对象. MDS（称为ceph-mds的守护进程）存在的原因是，简单的文件系统操作，比如ls、cd这些操作会不必要地扰动OSD，所以把元数据从数据里分离出来意味着Ceph文件系统既能提供高性能服务，又能减轻存储集群负载.

### 1. Multi Active MDS

从Luminous(v12.xxx)开始，Multi Active MDS已经稳定，这样可以提高Ceph FS元数据的处理能力.

目前在Multi Active MDS功能基础上，并没有一步实现所谓的动态分布式数据管理. 而是折中实现了静态划分绑定：允许用户将不同的目录绑定在不同的MDS上，以此达到比较好的负载均衡.

当目录变得越来越大或访问频率越来越高时，目录所在的MDS就会变成瓶颈. 最简单的方法就是将目录分割，放在不同的MDS上. 但是又希望这种分割是自动的，是对用户透明的. 同样我们也希望随着删除、目录内容减少或访问频率减少，目录可以合并在一起。可以根据目录长度和访问频率两个维度来进行分割或合并:
- 目录长度

	当目录长度超过mds_bal_split_size（默认为10 000）后，就会进行分割。但是在正常情况下并不会马上分割，因为分割动作会影响正常操作，所以会在mds_bal_fragment_interval秒后分割。如果目录长度超过mds_bal_fragment_fast_ factor就会马上分割。分割子目录数是2 ^mds_bal_split_bits。

	mds_bal_fragment_size_max是目录片段大小的硬限制。如果达到，客户端会在片段中创建文件时收到ENOSPC错误。在正确配置的系统上，永远不应该在普通目录上达到此限制。

	当目录片段的大小小于mds_bal_merge_size时，会开始进行合并。

- 访问频率

	MDS为每个目录维护单独的时间衰减负载计数器（mds_decay_halflife），用于对目录片段进行读/写操作。写操作（包括元数据I/O，如重命名、删除和创建）导致写操作计数器增加，并与mds_bal_split_wr进行比较，如果超过则触发拆分。同样读操作导致读操作计数器增加，并与mds_bal_split_rd进行比较，如果超过则触发拆分。

	需要注意的是，根据访问频率进行分割后，它们仅基于大小阈值（mds_bal_merge_size）进行合并，因此根据访问频率进行分割可能导致目录永远保持碎片.

### 2. 配额
Ceph FS允许给系统内的任意目录以大小或文件数目的形式设置配额. 配额以xattr的方式存放：ceph.quota.max_files & ceph.quota.max_bytes. Ceph FS虽然支持quota，但是有以下局限性:
- Ceph FS的配额功能依赖于挂载它的客户端的合作，在达到上限时要停止写入；无法阻止篡改过的或对抗性的客户端，它们可以想写多少就写多少。在客户端完全不可信时，用配额防止多占空间是靠不住的。
- 配额是不准确的。因为配额是客户端和服务器端合作，需要两种数据同步，所以在达到配额限制一小段时间后，正在写入文件系统的进程才会被停止，很难避免它们超过配置的限额、多写入一些数据等问题。超过配额多大幅度主要取决于时间长短，而非数据量。
- 内核客户端还没有实现配额功能。用户空间客户端（libcephfs、ceph-fuse）已经支持配额，但是Linux内核客户端还没有实现配额功能。
- quota是针对目录进行设置的，并没有根据UID/GID进行设置。
- 基于路径限制挂载时必须谨慎地配置配额。客户端必须能够访问配置了配额的那个目录的索引节点，这样才能执行配额管理。如果某一个客户端被MDS能力限制只能访问一个特定路径（如/home/user），并且它无权访问配置了配额的父目录（如/home），这个客户端就不会按配额执行。所以，基于路径进行访问控制时，最好在限制了客户端的那个目录（如/home/user）或者在它下面的子目录上配置配额。

## 后端存储ObjectStore
ObjectStore是Ceph OSD中最重要的概念之一，它完成实际的数据存储，封装了所有对底层存储的I/O操作. I/O请求从客户端发出后，最终会使用ObjectStore提供的API进行相应的处理.

ObjectStore也有不同的实现方式，目前主要有FileStore、BlueStore、MemStore、KStore. 可以在配置文件中通过osd_objectstore字段来指定ObjectStore的类型. MemStore和元数据主要用于测试，其中MemStore将所有数据全部放在内存中，KStore将元数据与Data全部存放到KVDB中. MemStore和KStore都不具备生产环境的要求.

### FileStore
FileStore是基于Linux现有的文件系统，将Object存放在文件系统上的，也就是利用传统的文件系统操作实现ObjectStore API. 每个Object会被FileStore看作是一个文件，Object的属性（xattr）会利用文件的属性存取，因为有些文件系统（如ext4）对属性的长度有限制，所以超出长度的属性会作为omap存储.

FileStore目前支持的文件系统有XFS/ext4/Btrfs，推荐使用XFS.

1. 写操作

对于FileStore来说，一个简单的put操作，比如“rados -p rbd put test test-object”，包含了下面的4个操作，这4个操作在Ceph里合在一起称为事务. 事务要求是原子性的，或全部成功，或全部失败，不允许有中间过程.

- Touch：对应系统调用create。
- Setattrs：对应系统调用setattrs。
- Write：对应系统调用write。
- Omap_setkeys：对应RocksDB的transaction。

因为ObjectStore这4个操作组成的事务是原子性的，但是文件系统本身并不支持，所以需要引入journal.

2. journal

因为底层文件系统并不支持事务的原子性要求，所以就需要引入journal, 来达到事务的原子性. 简单来说就是将事务作为一个entry事先存放在journal上，然后将事务的操作在文件系统上执行. 如果文件系统正常完成，中间没有出现故障，那么就通知journal删除这个entry. 如果出现故障，那么在OSD 重启的时候，从journal上读取未删除的entry进行重做.

journal可以是一个普通文件或软链接。使用软链接这种形式是为journal性能考虑的，journal是关键路线，其性能直接影响整个ObjectStore的性能。所以允许使用额外的块设备。在生产环境中，也是这样建议的，用固态硬盘、NVMe或Optane作为journal。

FileStore将journal作为一个环形缓冲区（Ring Buffer），避免删除操作带来的消耗，进行删除操作时仅仅将有效位置前移即可。每个事务编码作为一个entry存放到journal中，然后等待FileStore完成后，将entry删除。如果journal由于某种原因满了，那么就必须要等待FileStore完成。所以创建OSD时，需要考虑journal的大小，太小容易满，从而导致操作阻塞。

journal的优点是保证事务的原子性，但是带来的消耗是巨大的。任何操作都要先写到journal，再写到文件系统上，数据量至少是两倍的增长。同时XFS也带有journal，写放大很严重，很容易导致journal成为瓶颈。而且XFS/Btrfs是通用文件系统，满足所有可移植操作系统接口的要求，所以在某些操作上性能就没有那么高。这时就需要采用一个新的实现形式来替代FileStore，以获得更高的性能。

### BlueStore
FileStore最初只是针对机械盘进行设计的，并没有对固态硬盘的情况进行优化，而且写数据之前先写journal也带来了一倍的写放大. 为了解决FileStore存在的问题，Ceph社区推出了BlueStore. BlueStore去掉了journal，通过直接管理裸设备的方式来减少文件系统的部分开销，并且也对固态硬盘进行了单独的优化.

![BlueStore架构](/misc/img/ceph/Image00179.jpg)

BlueStore和传统的文件系统一样，由3个部分组成：数据管理、元数据管理、空间管理（即所谓的Allocator）. 与传统文件系统不同的是，数据和元数据可以分开存储在不同的介质中.

BlueStore不再基于ext4/XFS等本地文件系统，而是直接管理裸设备，为此在用户态实现了BlockDevice，使用Linux AIO直接对裸设备进行I/O操作，并实现了Allocator对裸设备进行空间管理. 至于元数据则以Key/Value对的形式保存在KV数据库里（默认为RocksDB）. 但是RocksDB并不是基于裸设备进行操作的，而是基于文件系统（RocksDB可以将系统相关的处理抽象成Env，BlueStore实现了一个BlueRocksEnv，来为RocksDB提供底层系统的封装）进行操作的，为此BlueStore实现了一个小的文件系统BlueFS，在将BlueFS挂载到系统中的时候将所有的元数据加载到内存中.

1. KVDB

使用KVDB的原因是它可以快速实现ObjectStore的事务的语义. BlueStore默认使用的KVDB是RocksDB，但不表示只能用RocksDB. 只要能提供相应接口的，任何KVDB都可以使用.

RocksDB是Facebook公司基于LevelDB开发的KVDB。用于解决LevelDB单线程合并导致性能不佳的问题，并同时丰富API。RocksDB后端存储必须基于一个文件系统，所以对于RocksDB有两种选择：使用标准Linux文件系统，比如ext4/XFS等；实现一套简单的用户空间文件系统，满足RocksDB的需求。同样基于性能的考虑，Ceph实现了BlueFS，而没有使用XFS/Btrfs.

MemDB是专门为测试性能而产生的一个KVDB，其实就是用std:map在内存中存放元数据，通过了解MemDB可以快速了解实现一个KVDB需要的功能。MemDB在正常关闭情况下会将std::map的内容写到本地一个文件中，在重启时被读取出来。

2. Allocator

由于BlueStore和BlueFS直接操作裸设备，所以和传统文件系统一样，磁盘空间也需要自己来管理。不过这部分并没有创新，都是按照最小单元进行格式化的，使用BitMap来进行分配。

3. BlueFS

BlueFS是专门为RocksDB开发的文件系统，根据RocksDB WAL和SST（Sorted Sequence Table）的不同特性，可以配置两个单独的块设备。BlueFS自己管理裸设备，所以也有自己的Allocator，不同于BlueStore，BlueFS的元数据是作为一个内部特殊的文件进行管理的，这是因为它的文件和目录都不是很大。

4. 元数据

和传统文件系统的元数据一样，包含文件的基本信息：name、mtime、size、attributes，以及逻辑和物理块的映射信息，同时包含裸设备的块的分配信息。

ext2/3/4在存放大量小文件时可能会出现元数据空间不够，但是数据空间大量空闲的情况。BlueFS会共享BlueStore的块设备，在其原有空间不够的情况下，会从这个块设备分配空间。用这种方式来解决元数据和数据空间出现不匹配的问题。

5. BlockDevice

在使用BlueFS+RocksDB的情况下最多有3个块设备：Data存放BlueStore数据；WAL存放RocksDB内部journal；DB存放RocksDB SST.

- 1个块设备。只有Data设备，BlueStore和BlueFS共享这个设备。
- 2个块设备。有2种配置：Data+WAL，RocksDB SST存放在Data上；Data+DB，WAL放在DB上。
- 3个块设备。Data+DB+WAL，原则上数据存放在Data上，RocksDB journal存放在WAL上，RocksDB SST存放在DB上。
- 空间互借的原则。mkfs时不知道元数据到底有多少，KVDB空间大小也很难确定，所以就会允许当空间出现不足时，可以按照WAL→DB→Data这个原则进行分配。

6. 驱动

驱动是指如何读/写块设备，根据不同的设备种类可以选择不同的方式，以获取更高的性能.

- 系统调用读/写。

	这是最常用的方式，对任何块设备都适用。使用Linux系统调用读/写或aio_submit来访问块设备.

- SPDK

SPDK是专门为NVMe开发的用户空间驱动，放弃原有系统调用模式，是基于以下两个原因：①NVMe设备的IOps很高，如果按照传统的系统调用，需要先从用户空间进入内核空间，完成后再从内核空间返回用户空间，这些开销对NVMe而言就比较大了；②内核块层是通用的，而作为存储产品这些设备是专用的，不需要传统的I/O调度器。

- PMDK

PMDK是专门为英特尔AEP设备开发的驱动。AEP可以看成是带电的DRAM，掉电不丢数据。不同于NAND Flash，写之前需要擦除。基本存储单元不是扇区，而是和DRMA一样为Byte。接口不是PCIe，而是和DRAM一样直接插在内存槽上的。目前Kernel已经有驱动，将AEP作为一个块设备使用，也针对AEP的特性，实现了DAX的访问模式（AEP设备不使用Page Cache）。

基于性能的考虑，PMDK专门作为AEP的驱动出现。考虑到AEP存储单元是Byte，所以其将AEP设备mmap，这样就可以直接按字节访问AEP了。由于使用的是DAX技术，所以mmap跳过Page Cache，直接映射到AEP。这样就实现了用户空间直接访问AEP，极大地提高了I/O性能.

## SeaStore
SeaStore是下一代的ObjectStore，目前仅仅有一个设计雏形. 因为SeaStore基于SeaStar，所以暂时称为SeaStore，但是不排除后面有更合适的名字. SeaStore有以下几个目标:
- 专门为NVMe设备设计，而不是PMEM 和硬盘驱动器。
- 使用SPDK访问NVMe而不再使用Linux AIO。
- 使用SeaStar Future编程模型进行优化，以及使用share-nothing机制避免锁竞争。
- 网络驱动使用DPDK来实现零拷贝。

由于Flash设备的特性，重写时必须先要进行擦除操作。垃圾回收擦除时，并不清楚哪些数据有效，哪些数据无效（除非显式调用discard线程），但是文件系统层是知道这一点的。所以Ceph希望将垃圾回收功能可以提到SeaStore来做。SeaStore的设计思路主要有以下几点:
- SeaStore的逻辑段（segment）应该与硬件segment（Flash擦除单位）对齐。
- SeaStar是每个线程一个CPU核，所以将底层按照CPU核进行分段。
- 当空间使用率达到设定上限时，就会进行回收。当segment完全回收后，就会调用discard线程通知硬件进行擦除。尽可能保证逻辑段与物理段对齐，避免逻辑段无有效数据，但是底层物理段存在有效数据，那么就会造成额外的读/写操作。同时由discard带来的消耗，需要尽量平滑处理回收工作，减少对正常读/写的影响。
- 用一个公用的表管理segment的分配信息，所有元数据用B-Tree进行管理。

## crush算法
为一个分布式存储系统实现数据分布算法不简单，至少需要考虑下述情况:
- 实现数据的随机分布，并在读取时能快速索引
- 能够高效地重新分布数据，在设备加入、删除和失效时最小化数据迁移
- 能够根据设备的物理位置合理地控制数据的失效域
- 支持常见的镜像、磁盘阵列、纠删码等数据安全机制
- 支持不同存储设备的权重分配，来表示其容量大小或性能

CRUSH算法在设计时就考虑了上述5种情况. CRUSH算法根据输入x 得到随机的n 个有序的位置，即OSD.k，并保证在相同的元数据下，对于输入x 的输出总是相同的. Ceph只需要在集群中维护并同步少量的元数据，各个节点就能独立计算出所有数据的位置，并且保证输出结果对于同样的输入x 是相同的, 且没有任何中心节点.

CRUSH元数据包含了CRUSH Map、OSDMap和CRUSH Rule. 其中，CRUSH Map保存了集群中所有设备或OSD存储节点的位置信息(物理组织结构)和权重设置，使CRUSH算法能够感知OSD的实际分布和特性，并通过用户定义的CRUSH Rule来保证算法选择出来的位置能够合理分布在不同的失效域中. 而OSDMap保存了各个OSD的运行时状态，能够让CRUSH算法感知存储节点的失效、删除和加入情况，产生最小化的数据迁移，提高Ceph在各种情况下的可用性和稳定性.

CRUSH算法即通过一系列精心设计的哈希算法去访问和遍历CRUSH Map，按照用户定义的规则选择正常运行的OSD来保存数据对象.

### 1. CRUSH Map与设备的物理组织
CRUSH Map本质上是一种有向无环图（DAG），用来描述OSD的物理组织和层次结构. 

![CRUSH Map](/misc/img/ceph/Image00181.jpg)
其结构如上图，所有的叶子节点表示OSD设备，而所有的非叶子节点表示桶(bucket). 桶根据层次来划分可以定义不同的类型（CRUSH Type或Bucket Class），如根节点、机架、电源等.

在默认情况下，Ceph会创建两种类型的桶，分别是根节点和主机，然后把所有的OSD设备都放在对应的主机类型桶中，再把所有主机类型桶放入一个根节点类型桶中. 在更复杂的情况下，比如要防止由于机架网络故障或电源失效而丢失数据，就需要用户自行创建桶的类型层次并建立对应的CRUSH Map结构.

查看当前集群的CRUSH Map可用命令`ceph osd crush tree`, 其中，负数ID的行表示CRUSH桶，非负数ID的行表示OSD设备，CLASS表示OSD设备的Device Class，TYPE表示桶类型，即CRUSH Type.

### 2. CRUSH Map 叶子
CRUSH Map中所有的叶子节点即OSD设备，每个OSD设备在CRUSH Map中具有名称、Device Class和全局唯一的非负数ID. 其中，默认的Device Class有硬盘驱动器、固态硬盘和NVMe 3种，用于区分不同的设备类型. Ceph可以自动识别OSD的Device Class类型，当然也可以由用户手动创建和指定. 当前Ceph内部会为每个Device Class维护一个Shadow CRUSH Map，在用户规则中指定选择某一个Device Class，比如固态硬盘时，CRUSH算法会自行基于对应的Shadow CRUSH Map执行. 可以使用以下命令查看Device Class和Shadow CRUSH Map：
```bash
# ceph osd crush class ls # 列出所有的device class
# ceph osd class ls-osd <class> # 列出属于某个device class的osd设备
# ceph osd crush tree --show-Shadow # 显示shadow crush map
```

### 3. CRUSH Map桶
CRUSH Map中所有的非叶子节点即桶，桶也具有名称、Bucket Class和全局唯一的负数ID。属于同一种Bucket Class的桶往往处于CRUSH Map中的同一层次，其在物理意义上往往对应着同一类别的失效域，如主机、机架等。

作为保存其他桶或设备的容器，桶中还可以定义具体的子元素列表、对应的权重（Weight）、CRUSH算法选择子元素的具体策略，以及哈希算法。其中，权重可以表示各子元素的容量或性能，当表示为容量时，其值默认以TB为单位，可以根据不同的磁盘性能适当微调具体的权重。CRUSH算法选择桶的子元素的策略又称为Bucket Type，默认为Straw方式，它与CRUSH算法的实现有关，我们只需要知道不同的策略与数据如何重新分布、计算效率和权重的处理方式密切相关. 桶中的哈希算法默认值为0，其意义是rjenkins1，即Robert Jenkin's Hash. 它的特点是可以保证即使只有少量的数据变化，或者有规律的数据变化也能导致哈希值发生巨大的变化，并让哈希值的分布接近均匀. 同时，其计算方式能够很好地利用32位或64位处理器的计算指令和寄存器，达到较高的性能.

![在CRUSH Map中，Bucket Class与桶的具体定义](/misc/img/ceph/Image00184.jpg)

### 4. OSDMap与设备的状态
在运行时期，Ceph的Monitor会在OSDMap中维护一种所有OSD设备的运行状态，并在集群内同步. 其中，OSD运行状态的更新是通过OSD-OSD和OSD-Monitor的心跳完成的. 任何集群状态的变化都会导致Monitor中维护的OSDMap版本号（Epoch）递增，这样Ceph客户端和OSD服务就可以通过比较版本号大小来判断自己的Map是否已经过时，并及时进行更新.

OSD设备的具体状态可以是在集群中（in）或不在集群中（out），以及正常运行（up）或处于非正常运行状态（down）. 其中OSD设备的in、out、up和down状态可以任意组合，只是当OSD同时处于in和down状态时，表示集群处于不正常状态. 在OSD快满时，也会被标记为full. 我们可以通过以下命令查询OSDMap的状态，或者手动标记OSD设备的状态：
```bash
# ceph osd stat # 查看osd状态和osdmap epoch
# --- 手动标记osd设备的状态
# ceph osd up <osd-ids>
# ceph osd down <osd-ids>
# ceph osd in <osd-ids>
# ceph osd out <osd-ids>
```

### CRUSH中的规则与算法细节
#### 1. CRUSH Rule基础

仅仅了解了OSD设备的位置和状态，CRUSH算法还是无法确定数据该如何分布。由于具体使用需求和场景的不同，用户可能会需要截然不同的数据分布方式，而CRUSH Rule就提供了一种方式，即通过用户定义的规则来指导CRUSH算法的具体执行. 其场景主要如下所示:
- 数据备份的数量：规则需要指定能够支持的备份数量。
- 数据备份的策略：通常来说，多个数据副本是不需要有顺序的；但是纠删码不一样，纠删码的各个分片之间是需要有顺序的。所以CRUSH算法需要了解各个关联的副本之间是否存在顺序性。
- 选择存储设备的类型：规则需要能够选择不同的存储设备类型来满足不同的需求，比如高速、昂贵的固态硬盘类型设备，或者低速、廉价的硬盘驱动器类型设备。
- 确定失效域：为了保证整个存储集群的可靠性，规则需要根据CRUSH Map中的设备组织结构选择不同的失效域，并依次执行CRUSH算法。

Ceph集群通常能够自动生成默认的规则，但是默认规则只能保证集群数据备份在不同的主机中。实际情况往往更加精细和复杂，这就需要用户根据失效域自行配置规则，保存在CRUSH Map中.
![](/misc/img/ceph/Image00186.jpg)

crush rule定义如上图. 其中，规则能够支持的备份数量是由min_size和max_size确定的，type确定了规则所适用的备份策略. Ceph在执行CRUSH算法时，会通过ruleset对应的唯一ID来确定具体执行哪条规则，并通过规则中定义的step来选择失效域和具体的设备.

所有规则的详细定义可用`ceph osd crush rule dump`查看.

#### 2. CRUSH Rule的step take与step emit
CRUSH Rule执行步骤中的第一步和最后一步分别是step take与step emit. step take通过桶名称来确定规则的选择范围，对应CRUSH Map中的某一个子树. 同时，也可以选择Device Class来确定所选择的设备类型，如固态硬盘或硬盘驱动器，CRUSH算法会基于对应的Shadow CRUSH Map来执行接下来的step choose. step take的具体定义是`step take <bucket_name> [class <Device Class>]`.

step emit非常简单，即表示步骤结束，输出选择的位置.

#### 3. step choose与CRUSH算法原理
CRUSH Rule的中间步骤为step choose，其执行过程即对应CRUSH算法的核心实现. 每一step choose需要确定一个对应的失效域，以及在当前失效域中选择子项的个数. 由于数据备份策略的不同（如镜像与纠删码），step choose还要确定选择出来的备份位置的排列策略. 其定义如下：
![](/misc/img/ceph/Image00189.jpg)

此外，CRUSH Map中的桶定义也能影响CRUSH算法的执行过程. 例如，CRUSH算法还需要考虑桶中子项的权重来决定它们被选中的概率，同时，在OSDMap中的运行状态发生变化时，尽量减少数据迁移. 具体的子元素选择算法是由桶定义里面的Bucket Type来确定的.

桶定义还能决定CRUSH算法在执行时所选择的哈希算法. 哈希算法往往会导致选择冲突的问题. 类似地，当哈希算法选择出OSD设备后，可能会发现其在OSDMap被标记为不正常的运行状态. 这时，CRUSH算法需要有自己的一套机制来解决选择冲突和选择失败问题:

1. 选择方式、选择个数与失效域
在step的配置中，可以定义在当前步骤下选择的Bucket Class，即失效域，以及选择的具体个数n. 例如，让数据备份分布在不同的机架中，代码如下：
![](/misc/img/ceph/Image00190.jpg)

或者是让数据分布在两个电源下面的两个主机中，代码如下：
![](/misc/img/ceph/Image00191.jpg)

其中，当选择个数n的值为0时，表示选择与备份数量一致的桶；当n的值为负数时，表示选择备份数量减去n个桶；当n的值为正数时，即选择n个桶.

chooseleaf可以被看作choose的一种简写方式，它会在选择完指定的Bucket Class后继续递归直到选择出OSD设备. 例如，让数据备份分布在不同的机架的规则也可以写成：
![](/misc/img/ceph/Image00192.jpg)

2. 选择备份的策略：firstn与indep
CRUSH算法的选择策略主要和备份的方式有关. firstn对应的是以镜像的方式备份数据。镜像备份的主要特征是各数据备份之间是无顺序的，即当主要备份失效时，任何从属备份都可以转为主要备份来进行数据的读/写操作。其内部实现可以理解为CRUSH算法维护了一个基于哈希算法选择出来的设备队列，当一个设备在OSDMap中标记为失效时，该设备上的备份也会被认为失效。这个设备会被移出这个虚拟队列，后续的设备会作为替补。firstn的字面意思是选出虚拟队列中前n 个设备来保存数据，这样的设计可以保证在设备失效和恢复时，能够最小化数据迁移量.

而indep对应的是以纠删码的方式来备份数据的。纠删码的数据块和校验块是存在顺序的，也就是说它无法像firstn一样去替换失效设备，这将导致后续备份设备的相对位置发生变化。而且，在多个设备发生临时失效后，无法保证设备恢复后仍处于原来的位置，这就会导致不必要的数据迁移。indep通过为每个备份位置维护一个独立的（Independent）虚拟队列来解决这个问题。这样，任何设备的失效就不会影响其他设备的正常运行了；而当失效设备恢复运行时，又能保证它处于原来的位置，降低了数据迁移的成本.

虚拟队列是通过计算索引值来实现的.

3. 选择桶的子元素的方式：Bucket Type

CRUSH算法在确定了最终选择的索引值后，并不是按照索引值从对应的桶中直接选出子桶或子设备的，而是提供了多个选择，让用户能够根据不同情况进行配置。子元素的选择算法通过Bucket Type来配置，分别有Uniform、List、Tree、Straw和Straw2 5种。它们各自有不同的特性，可以在算法复杂度、对集群设备增减的支持，以及对设备权重的支持3个维度进行权衡。

> straw在迁移数据时, 会导致部分设备之间的数据迁移量； 当换成Straw2算法后，在相同的条件下整体数据迁移量能被控制在12.40%，约等于1/8的数据量。这基本上能与理论上最低的迁移量保持一致

根据Bucket Type的算法实现可以将子元素的选择算法分为三大类。首先是Uniform，它假定整个集群的设备容量是均匀的，并且设备数量极少变化。它不关心子设备中配置的权重，直接通过哈希算法将数据均匀分布到集群中，所以也拥有最快的计算速度O(1)，但是其适用场景极为有限.

其次是分治算法。List会逐一检查各个子元素，并根据权重确定选中对应子元素的概率，拥有O(n)的复杂度，优点是在集群规模不断增加时能够最小化数据迁移，但是在移除旧节点时会导致数据重新分布。Tree使用了二叉搜索树，让搜索到各子元素的概率与权重一致。它拥有O （logn ）的算法复杂度，并且能够较好地处理集群规模的增减。但是Tree算法在Ceph实现中仍有缺陷，已经不再推荐使用。

分治算法的问题在于各子元素的选择概率是全局相关的，所以子元素的增加、删除和权重的改变都会在一定程度上影响全局的数据分布，由此带来的数据迁移量并不是最优的。第三类算法的出现即为了解决这些问题。Straw会让所有子元素独立地互相竞争，类似于抽签机制，让子元素的签长基于权重分布，并引入一定的伪随机性，这样就能解决分治算法带来的问题。Straw的算法复杂度为O(n)，相对List耗时会稍长，但仍是可以被接受的，甚至成为了默认的桶类型配置。然而，Straw并没有能够完全解决最小数据迁移量问题，这是因为子元素签长的计算仍然会依赖于其他子元素的权重。Straw2的提出解决了Straw存在的问题，在计算子元素签长时不会依赖于其他子元素的状况，保证数据分布遵循权重分布，并在集群规模变化时拥有最佳的表现.

## CRUSH算法实践
1. 创建和维护CRUSH Map桶结构

Ceph集群在创建时能够自动生成默认的CRUSH Map和规则，即为每个OSD设备生成对应的CRUSH叶子，为每个主机生成对应的CRUSH桶，包含在一个名为root的根桶中. 默认规则的名称是replicated_rule，它将数据镜像分布到不同的OSD设备中. 我们可以通过以下命令获得默认的CRUSH Map字面定义：
```bash
# ceph osd getcrushmap > crush.bin
# crushtool -d crush.bin -o crush.txt
```

从集群中直接获得的CRUSH Map是一个二进制文件（crush.bin），需要通过crushtool解码成为可读的文本文件（crush.txt）.

获得CRUSH Map文本文件之后，就可以根据实际情况修改CRUSH算法结构，例如，支持新的失效域，以及指定新的CRUSH算法规则等. 之后上传新的CRUSH Map文件（crush_new.txt）来应用这些改动：
```bash
# crushtool -c crush_new.txt -o crush_new.bin
# ceph osd setcrushmap -i crush_new.bin
# ceph osd tree # 查看改动情况
```

此外，还可以用以下命令直接修改集群中保存的CRUSH Map：
```bash
# ceph osd crush add-bucket <bucket_name> <Bucket Class>
# ceph osd crush rename-bucket <old_bucket_name> <new_bucket_name>
# ceph osd crush swap-bucket <bucket_name> <替换的桶名>
# ceph osd crush create-or-move/move <bucket_name> <目的Bucket Class>=<目的桶>
# ceph osd crush rm <bucket_name>
```

2. 维护CRUSH设备的位置

当启动OSD服务时，Ceph需要知道该OSD在CRUSH Map中的具体位置，也称CRUSH Location. 默认地，Ceph会自动识别OSD所在的主机名称，并将它放入对应的主机桶中. 此外，还能通过多种方式维护OSD设备的CRUSH位置.

其中，最普遍的方式是在ceph.conf文件中配置crush_location，用来描述当前主机上所有设备的位置，例如：
```conf
[osd]
crush location = row=a rack=2 chassis=a2a host=a2a1
```

第二种方式是配置osd_crush_location_hook脚本，这样Ceph就可以通过执行该脚本自动获得当前主机的CRUSH Location了.

此外，也可以手动执行命令更新OSD位置:`ceph osd crush set <osd 设备id> <osd 设备权重> <CRUSH Location>`

若需要避免OSD服务在启动时自动覆盖已定义的设备位置，需要在ceph.conf文件中将osd_crush_update_on_start的值设为false.

3. 调整 CRUSH算法
CRUSH算法发展到现在，其实已经经历了多次迭代，提供了一些额外特性进行调整. 这里主要介绍CRUSH Tunables及CRUSH算法如何更精细地控制主OSD的选择.

	1. Tunables

	CRUSH算法经历过几次内部调整，每次的调整都会带来新的微调参数，比如子元素选择失败的次数、允许的桶类型和一些算法细节的优化开关等，这些调整选项就称为CRUSH Tunables. 它们在不同的Ceph版本之间会有不兼容的问题，而且大部分调整都是为了解决CRUSH算法的原有缺陷，每个版本都有自己的最优配置. 所以用户一般不需要知道具体如何调整CRUSH Tunables，一般可以选择名为optimal的描述文件，在出现兼容性问题时，也可以选择旧版本的描述文件: `ceph osd crush tunables optimal`

	2. 主OSD设备的亲和性

	在使用镜像策略时，由于主OSD承载了与客户端的连接和主要的读/写任务，所以Ceph对主OSD的性能要求会比较高. 在一个拥有不同性能设备的Ceph集群中，我们只需要知道集群哪些设备适合作为主OSD使用，哪些是不适合的，可以使用以下命令来实现:`ceph osd primary-affinity <osd 设备id> <权重: 0~1>`, 其中权重0表示不适合作为主osd设备, 1表示适合做主osd设备.

	3. 控制主OSD设备的选择范围

	除了使用亲和性设置，我们还可以使用自定义的CRUSH Rule定义更清晰的规则，以确定主OSD的选择范围，命令如下：
	![](/misc/img/ceph/Image00205.jpg)

4. 测试CRUSH算法
这里将介绍用于协助调试和优化Ceph集群中的CRUSH Map和CRUSH Rule的工具

	1. crushtool工具

	crushtool能协助编码及解码二进制的CRUSH Map，实际上它还能协助对CRUSH Map进行测试. 例如，测试编写的规则在当前的CRUSH Map中是否会出现运行失败的情况，命令如下：
	```bash
	# ceph osd getcrushmap > crush.bin
	# crushtool -i crush.bin --test --show-statistics
	```

	上述命令会为crush.bin中定义的所有规则中所有支持的镜像数量（numrep从min_size至max_size）进行测试，每轮测试会生成从0至1023共1024个输入x，进入CRUSH算法执行. `--show-statistics`会输出当前测试的执行结果.

	crushtool能显示具体的错误结果`crushtool -i crush.bin --test --show-bad-mappings`

	`crushtool -i crush.bin --test --num-rep 3 --show-utilization`可测试数据分布是否均匀.

	`crushtool -i crush.bin --test --num-rep 3 --show-choose-tries`可测试规则执行过程中发生的冲突和重试情况.

	2. crush小程序

	除了官方的crushtool，还可以选择第三方的工具对CRUSH Map进行测试，如crush，其安装命令如下：`pip install crush`.

	`crush analyze --type device --rule data --crushmap crush.json`可用于分析Map是否能按照权重均匀地分布数据, 类似于类似crushtool的--show-utilization测试.

	crush工具的一个更为强大的地方是它能够用于检测和比较OSD设备的增加或移除导致的数据迁移量问题. 例如，当Bucket Type配置为Straw算法时: `crush compare --rule data --origin straw.json --destination straw-add.json`.

## CRUSH算法在Ceph中的应用
Ceph并不是直接通过CRUSH算法将数据对象一一映射到OSD中的，这样做将非常复杂与低效, 而是通过两个逻辑概念，即存储池和PG，很好地将CRUSH算法的复杂性隐藏起来，并向用户提供了一个直观的对象存储接口，即RADOS.

1. 存储池
Ceph中的存储池是由管理员创建的，用于隔离和保存数据的逻辑结构. 不同的存储池可以关联不同的CRUSH规则，指定不同的备份数量，以镜像或纠删码的方式保存数据，或是指定不同的存取权限和压缩方式. 用户还可以以存储池为单位做快照，甚至建立上层应用，比如块存储、文件存储和对象存储. 至于存储池本身是如何保证数据安全性和分布的，这是由存储池关联的CRUSH规则决定的，用户不需要深入了解这些细节，只需要记住数据名称，选择合适的存储池即可在RADOS中执行读/写操作. 创建一个存储池的命令是`ceph osd pool create <pool_name> <pg_num> <replicated|erasure> <crush rule name>`

2. PG
存储池在创建时需要指定PG的数量。PG是存储池中保存实际数据对象的逻辑概念。Ceph中的数据对象总是唯一地保存在某个存储池中的一个PG里，或者存储池通过维护固定数量的PG来管理和保存所有的用户数据对象。

和存储池概念不同的是，**用户无法决定对象具体保存在哪个PG中，这是由Ceph自行计算的**。PG的意义在于为存储空间提供了一层抽象，它在概念上对应操作系统的虚拟地址空间，将数据对象的逻辑位置和其具体保存在哪些OSD设备的细节相互隔离。这样能够实现将各存储池的虚拟空间相互隔离，还能更好地利用和维护OSD设备的物理空间. Ceph首先将数据对象通过哈希算法分布到不同的逻辑PG中，再由CRUSH算法将PG合理分布到所有的OSD设备中。也就是说，CRUSH算法实际上是以PG而不是数据对象为单位确定数据分布和控制容灾域的，而在逻辑上一组PG组成了一个连续的虚拟空间来保存存储池中的数据对象。

## FAQ
### DRBD vs Ceph
CEPH 是一种开源软件，旨在在统一系统中提供高度可扩展的对象，块和基于文件的存储.

DRBD 是一种基于软件和网络(tcp/ip和RDMA)的块复制存储解决方案, 基于块设备, 相对而言, 其在ceph的块设备之上.