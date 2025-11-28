# longhorn
参考:
- [**云原生分布式存储系统：Longhorn 初体验**](https://mp.weixin.qq.com/s/nrct5JIP8AepcP8wsVsFEQ)
- [Longhorn 云原生分布式块存储解决方案设计架构和概念](https://www.cnblogs.com/hacker-linner/p/15151778.html)
- [Longhorn全解析及快速入门指南](https://blog.51cto.com/u_12462495/1921097)
- [longhorn 源码分析](https://peteryj.github.io/categories/#longhorn-ref)
- [Longhorn 1.6.0 发布：引领性能革新与平台扩展](https://zhuanlan.zhihu.com/p/683919487)
- [https://github.com/longhorn/longhorn/wiki/Roadmap](https://github.com/longhorn/longhorn/wiki/Roadmap)
- [[FEATURE] Support v2 Data Engine - Experimental](https://github.com/longhorn/longhorn/issues/6229)
- [Longhorn系列-00-从iSCSI解析Longhorn](https://runzhliu.cn/longhorn%E7%B3%BB%E5%88%97-00-%E4%BB%8Eiscsi%E8%A7%A3%E6%9E%90longhorn/)
- [Longhorn 微服务化存储初试](https://mritd.com/2021/03/06/longhorn-storage-test/)
- [Longhorn 企业级云原生容器存储解决方案-部署篇](https://www.cnblogs.com/hacker-linner/p/15156477.html)
- [Longhorn 云原生容器分布式存储](https://www.cnblogs.com/hacker-linner/p/15221394.html)

## 组件
ref:
- [longhorn-spdk-engine](https://github.com/longhorn/longhorn-spdk-engine)

    参考README还是无法构建longhorn-spdk(没有`main()`, 根据[scripts/build](https://github.com/longhorn/longhorn-spdk-engine/blob/main/scripts/build)无法编译)

    - [Reimplement Longhorn Engine with SPDK](https://fossies.org/linux/longhorn/enhancements/20221213-reimplement-longhorn-engine-with-SPDK.md)
    - [How to setup a RAID1 block device with SPDK - longhorn/longhorn-spdk-engine GitHub Wiki](https://github-wiki-see.page/m/longhorn/longhorn-spdk-engine/wiki/How-to-setup-a-RAID1-block-device-with-SPDK)
- [Longhorn全解析及快速入门指南](http://user.1949idc.com/customer/shownews595.html)
- [longhorn-engine 源码分析](https://peteryj.github.io/2020/10/11/longhorn-engine-code-analysis/)
- [Longhorn 的正确使用姿势：如何处理增量 replica 与其中的 snapshot/backup](https://www.it120.vip/yq/10152.html)

由两部分组成:
1. longhorn manager
1. longhorn engine

    ![longhorn engine](/misc/img/develop/longhorn-engine.png)

安装要求:
1. 节点

    open-iscsi
1. fs

    ext4/xfs

操作iscsi的工具: [go-iscsi-helper](https://github.com/longhorn/go-iscsi-helper)

限制:
1. volume max未知

    1TB卷消耗256MB的内存读取索引(为了提高读取性能).
1. snap max 254

## 构建
- [Dockerfile.dapper](https://github.com/longhorn/longhorn-engine/blob/master/Dockerfile.dapper)

## 功能
大多数现有的分布式存储系统通常采用复杂的控制器软件来服务于从数百到数百万不等的volume. 但Longhorn不同，每个控制器上只有一个volume，Longhorn将每个volume都转变成了微服务. 每个控制器会有若干replica容器.

控制器的功能类似于典型的镜像RAID控制器，对其副本进行读写操作并监控副本的健康状况.

## FAQ
### 使用tgt的原因
ref:
- [NewFrontend](https://github.com/longhorn/longhorn-engine/blob/master/pkg/controller/init_frontend.go#L24)

    frontend类型(支持tgt-iscsi, tgt-blockdev, rest, socket这四种), 实际用户仅可使用, 见[The frontend (only Open-iSCSI/tgt are supported at this moment)](https://github.com/longhorn/longhorn-engine/blob/master/README.md)和[what is frontend?](https://github.com/longhorn/longhorn/issues/852)

tgt 是运行在用户态的iscsi网关，它支持多种后端，并且可扩展.

由于运行在用户态，在容器化流行的今天，它反而比Linux LIO显得更加灵活.
