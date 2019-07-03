# 分布式
![](/misc/img/distributed/bg2018071607.jpg)

参考:
- [分布式系统一致性（ACID、CAP、BASE、二段提交、三段提交、TCC、幂等性）原理详解](https://juejin.im/post/5c9443406fb9a070fe0dd9a9)
- [分布式事务数据库事务CAP定理BASE理论分布式事务案例](https://cloud.tencent.com/developer/article/1346890)
- [分布式系统的事务处理](http://coolshell.cn/articles/10910.html)

## cap : 分布式系统的三个指标
- Consistency : 一致性, 这里特指强一致

    所有节点上的数据**时刻**保持同步. 一致性严谨的表述是**原子读写**, 即所有读写都应该看起来是原子的或串行的. 所有的读写请求都好像是经全局排序过的一样,写后面的读一定能读到前面所写的内容.

    > 一致性与结果的正确性没有关系,而是系统对外主现的状态是否一致. 例如,所有节,最都达成一个错误的共识也是一致性的一种表现.
- Availability : 可用性

    任何**非故障节点**都应该在**有限的时间**内给出请求的响应,**不论请求是否成功**.
- Partition tolerance : 分区容忍性

    当发生网络分区时(即**节点之间无法通信**),在丢失**任意多**消息的情况下,系统仍然能够正常工作.

### Consistency 和 Availability 的矛盾
一般来说，**分区容错无法避免**，因此可以认为 CAP 的 P 总是成立. **CAP 定理告诉我们，剩下的 C 和 A 无法同时做到**.

#### AP
一旦发生网络分区(P),节点之间将无法通信,为了满足高可用(A),每个节点只能用本地数据提供服务, 这样就会导致数据的不一致(!C).

BASE(Basic Availability, Soft state, Eventually Consistency)理论是对CAP中的一致性和可用性进行一个权衡的结果，该理论的核心思想就是：我们无法做到强一致，但每个应用都可以根据自身的业务特点，采用适当的方式来使系统达到最终一致性(Eventualconsistency). 例如, NoSQL 数据库(Cassandra 、 CouchDB 等)往往会放宽对一致性的要求(**满足最终一致性**即可),以此来换取基本的可用性.

#### CP
如果要求数据在各个服务器上是强一致的(C),然而网络分区(P)会导致**同步**时间无限延长,那么如此一来可用性就得不到保障了(!A). 坚持事务 ACID 的传统数据库以及对结果一致性非常敏感的应用(例如,金融业务)通常会做出这样的选择.

#### CA
如果不存在网络分区,那么强一致性和可用性是可以同时满足的, 但该系统不一定是分布式系统了.

### 拜占庭错误(Byzantine Fa ilure)
它在计算机科学领域特指分布式系统中的某些恶意节点扰乱系统的正常运行,包括选择性不传递消息,选择性伪造消息等. 很显然, 拜占庭错误是一个 overly pessimistic 模型(最悲观 、 最强的错误模型).

进程失败错误(fail - stop Failure ,如同若机)则是一个 overly optimistic 模型(最乐观、最弱的错误模型) . 这个模型假设当某个节点出错时, 这个节点会停止运行,并且其他所有节点都知道这个节点发生了错误.

结论:

一个 RSM(Replicated State Machine, 复制状态机)系统要容忍 N 个拜占庭错误,至少需要 2N+l 个复制节点. 如果只是把错误的类型缩小到进程失败,则至少需要 N+l 个复制节点才能容错.

综上所述, 对于一个通用的、具有复制状态机语义的分布式系统,如果要做到 N 个节点的容错,理论上最少需要 2N+l 个复制节点. 这也是典型的一致性协议都要求半数以上(N/2+ 1)的服务器可用才能做出一致性决定的原因.

### FLP 不可能性
参考:
- [分布式理论梳理——FLP定理](https://my.oschina.net/duofuge/blog/1512344)
- [FLP 不可能原理](https://www.iminho.me/wiki/docs/blockchain_guide/distribute_system-flp.md)

FLP **定理**实际上说明了在允许节点失效的场景下,基于**异步通信**方式的分布式协议,无法确保在有限的时间内达成一致性.

> FLP 不可能原理告诉我们: 不要浪费时间，去试图为异步分布式系统设计面向任意场景的共识算法.
>
> 在分布式系统中,“异步通信”与“同步通信”的最大区别是没有时钟、不能时间同步、不能使用超时、不能探测失败、消息可任意延迟、消息可乱序等.

根据 FLP 定理,实际的一致性协议( Paxos 、 Ra武等)在理论上都是有缺陷的, 最大的问题是理论上存在不可终止性! 至于 Paxos 和 Raft 协议在工程的实现上都做了调整(例如, Paxos 和 Raft 都通过随机的方式显著降低了发生算法无法终止的概率), 并使用同步假设(或保证 safety ,或保证 liveness ), 从而规避了理论上存在的问题.

### [Paxos 算法与 Raft 算法](https://www.iminho.me/wiki/docs/blockchain_guide/distribute_system-paxos.md)
参考:
- [raft动画](http://thesecretlivesofdata.com/raft/)
- [微信自研生产级paxos类库PhxPaxos实现原理介绍](http://mp.weixin.qq.com/s?__biz=MzI4NDMyNTU2Mw==&mid=2247483695&idx=1&sn=91ea422913fc62579e020e941d1d059e)

因为Paxos的难理解, 这里只考虑Raft.

Raft 算法主要使用两种方法来提高可理解性: 问题分解和减少状态空间.
- 问题分解

    Raft 算法把问题分解成了领袖选举(leader election )、日志复制( log replication ) 、 安全性( safety )和成员关系变化
(membership changes )这几个子问题:

    - 领袖选举:在一个领袖节点发生故障之后必须重新给出一个新的领袖节点
    - 日志复制:领袖节点从客户端接收操作请求,然后将操作日志复制到集群中的其他服务器上,并且强制要求其他服务器的日志必须和自己的保持一致.
    - 安全性: Raft 关键的安全特性是状态机安全原则( State Machine Safety ). 即如果一个服务器已经将给定索引位置的日志条目应用到状态机中,则所有其他服务器不会在该索引位置应用不同的条目.
    - 成员关系变化:配置发生变化的时候,集群能够继续工作
- 减少状态空间

    Raft 算法通过减少需要考虑的状态数量来简化状态空间. 这将使得整个系统更加一致并且能够尽可能地消除不确定性. 需要特别说明的是,日志条目之间不允许出现空洞,并且还要限制日志出现不一致的可能性. 尽管在大多数情况下, Raft 都在试图消除不确定性以减少状态空间. 但在一种场景下(选举), Raft 会用随机方法来简化选举过程中的状态空间.


Raft 算法采用的是非对称节点关系模型, 在一个由 Raft 协议组织的集群中,一共包含如下 3 类角色:
- Leader (领袖)
- Candidate (候选人)
- Follower (群众)

## etcd
etcd 应用场景包括但不限于分布式数据库、服务注册与发现 、 分布式锁 、 分布式消息队列 、 分布式系统选主等.

etcd server 默认使用 2380 端口监听集群中其他 server 的请求.

> [搭建本地etcd集群](https://doczhcn.gitbook.io/etcd/index/index/local_cluster)

### [命令选项](https://doczhcn.gitbook.io/etcd/index/index-1/configuration)
> "proxy" 仅支持 v2 API, **不推荐**.

etcd选项:
- 成员标记
    - --name : 指定member名称
    - --data-dir : 数据目录的路径, 默认是`${name}.etcd`
    - --listen-client-urls : etcd server 向客户端暴露的通信地址列表
    - --listen-peer-urls : etcd server 向其他 member 暴露的通信地址列表
- 集群标记

    - --advertise-client-urls : 列出这个成员的客户端URL, 通告给集群中的其他成员. 

        用途: etcd server 之间是可以重定向请求的,比如, Follower 节点可将client的写请求重定向给 Leader 节点.
    - --initial-advertise-peer-urls : 列出member的peer urls 以便通告给集群的其他成员
    - --initial-cluster : 指定集群各节点所谓的 advertised peer URLs
    
        它们需要与 etcd 集群对应节点的`--initial-advertise-peer-urls`配置值相匹配
    - --initial-cluster-token : 在启动期间用于 etcd 集群的初始化集群的标识
    - --initial-cluster-state : 初始化集群状态("new" or "existing"). 
    
        在初始化静态(initial static)或者 DNS 启动 (DNS bootstrapping) 期间为所有成员设置为 new. 如果这个选项被设置为 existing , etcd 将试图加入已有的集群.

etcdctl选项:
- --endpoints : 指定 etcd server 的地址, 默认访问 `localhost:2379`

### etcd 集群的启动方式
1. 静态配置

    比较适用于线下环境, 前提条件: `集群节点个数已知` 和 `集群各节点的地址已知`.

    需要设置:
    ```
    --advertise-client-urls
    --initial-advertise-peer-urls
    --initial-cluster-token
    --initial-cluster
    --initial-cluster-state new
    ```

    以 `--initial-cluster` 开头的命令行参数将在 etcd 随后的运行中被忽略. 可以在初始化启动进程之后随意的删除环境变量或者命令行标记.
1. 服务发现
    1. etcd 自发现

        通过etcd自发现的方式启动 etcd 集群需要事先准备一个 etcd 集群.

        > 如果没有现成的集群可用，可以使用托管在 discovery.etcd.io 的公共发现服务

        [步骤](https://github.com/coreos/etcd/blob/master/Documentation/dev-internal/discovery_protocol.md):
        1. 设置discovery URL
        1. 用`--discovery`=discovery URL来启动etcd

            需要配置:
            ```
            --name // 每个成员必须有指定不同的名字标记, 否则发现会因为重复名字而失败
            --initial-advertise-peer-urls 
            --advertise-client-urls
            --discovery
            ```
    1. dns 自发现

        DNS SRV records 可以作为发现机制使用

        `-discovery-srv`可以用于设置 DNS domain name，在这里可以找到发现 SRV 记录.

        需要配置:
        ```
        --discovery-srv
        --initial-advertise-peer-urls
        --initial-cluster-token
        --initial-cluster-state
        --advertise-client-urls
        ```

        > 集群client/peer urls可以使用 IP 地址或域名来启动, **推荐域名**

