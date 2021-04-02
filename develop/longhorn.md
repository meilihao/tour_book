# longhorn
参考:
- [Longhorn全解析及快速入门指南](http://user.1949idc.com/customer/shownews595.html)
- [longhorn-engine 源码分析](https://peteryj.github.io/2020/10/11/longhorn-engine-code-analysis/)

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

## 功能
大多数现有的分布式存储系统通常采用复杂的控制器软件来服务于从数百到数百万不等的volume. 但Longhorn不同，每个控制器上只有一个volume，Longhorn将每个volume都转变成了微服务. 每个控制器会有若干replica容器.

控制器的功能类似于典型的镜像RAID控制器，对其副本进行读写操作并监控副本的健康状况.

## FAQ
### 使用tgt的原因
tgt 是运行在用户态的iscsi网关，它支持多种后端，并且可扩展.

由于运行在用户态，在容器化流行的今天，它反而比Linux LIO显得更加灵活.