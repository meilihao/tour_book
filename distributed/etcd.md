# etcd
etcd 应用场景包括但不限于分布式数据库、服务注册与发现 、 分布式锁 、 分布式消息队列 、 分布式系统选主等.

etcd server 默认使用 2380 端口监听集群中其他 server 的请求.

> [搭建本地etcd集群](https://doczhcn.gitbook.io/etcd/index/index/local_cluster)

## [命令选项](https://doczhcn.gitbook.io/etcd/index/index-1/configuration)
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

## etcd 集群的启动方式
1. [静态配置](http://play.etcd.io/install)

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

## 常用操作
> 先设置`export ETCDCTL_API=3`

```sh
$ etcdctl put ${key} ${value} // 写入key
$ etcdctl get [-w] [--rev=N] ${key} [--print-value-only] // 读取key, 默认输出key和value. `-w`设置输出格式, 比如`fileds`会输出更详细的信息(包括Revision, Lease等). `--rev`指定版本(etcd 集群全局的版本号)
$ etcdctl get --from-key ${key} // 输出大于或等于key的 kvs
$ etcdctl get ${key} ${key2} // 输出[key, key2)的kv
$ etcdctl get --prefix [--limit=N] ${key-prefix} // 遍历所有以 key-prefix 为前缀的 key, `--limit=N`限制输出数量
$ etcdctl del [--prev-kv] ${key} // 删除key. `--prev-kv`删除时返回value
$ etcdctl del ${key} ${key2} //删除[key, key2)的key
$ etcdctl del --prefix ${prefix-key}//删除的前缀是prefix-key的key
$ etcdctl del --from-key ${key} //删除`字典序>=key`的keys
$ etcdctl watch [--rev=N] ${key} // watch key的变化. `--rev`指定开始watch的Revision
$ etcdctl watch --prev-kv ${key} // `--prev-kv`是watch时返回变化前kv+最新的kev
$ etcdctl watch ${key} ${key2}
$ etcdctl watch --prefix ${key}
$ etcdctl compact N // 压缩所有 key版本号 N 之前的所有数据. 为了让client能够访问 key 过去任意版本的 value, etcd 会一直保存 key 所有历史版本的 value. 然而, etcd 所占的磁盘空间有限制.
$ etcdctl lease grant N // 授予租约，TTL为N秒, 创建时就开始生效了. 一旦租约的 TTL 到期，租约就过期并且所有关联的 key 都将被删除.
$ etcdctl put --lease=32695410dcc0ca06 ${key} ${value} // 为key添加租约
$ $ etcdctl lease revoke 32695410dcc0ca06 // 撤销租约. 租约被撤销后将会删除绑定在上面的所有 key.
$ etcdctl lease keep-alive 32695410dcc0ca06 // 续租,会阻塞命令, 每次续租都发生在该租约快过期时.
$ etcdctl lease timetolive [--keys] 694d5765fc71500b // 查询租约的 TTL 以及剩余时间. `--keys`同时返回该租约关联的keys
```

## mvcc
- main ID: 在etcd中每个事务的唯一id,全局递增不重复.
- sub ID: 在事务中的连续多个修改操作会从0开始编号,这个编号就是sub ID
- revision: 由(mainID,subID)组成的唯一标识

## Env
```conf
ETCD_NAME=etcd3 # etcd节点名称, 每个节点都应不同
ETCD_DATA_DIR=/etc/etcd/data # etcd数据存储目录

ETCD_CERT_FILE=/etc/etcd/pki/etcd_server.crt # etcd为client提供服务的server crt
ETCD_KEY_FILE=/etc/etcd/pki/etcd_server.key
ETCD_TRUSTED_CA_FILE=/etc/kubernetes/pki/ca.crt
ETCD_CLIENT_CERT_AUTH=true # 是否启用客户端证书认证
ETCD_LISTEN_CLIENT_URLS=https://192.168.18.5:2379 # 为client提供服务的url地址
ETCD_ADVERTISE_CLIENT_URLS=https://192.168.18.5:2379 # 广播给集群中其他成员自己为客户端提供服务的地址

ETCD_PEER_CERT_FILE=/etc/etcd/pki/etcd_server.crt # etcd为peer提供服务的server crt
ETCD_PEER_KEY_FILE=/etc/etcd/pki/etcd_server.key
ETCD_PEER_TRUSTED_CA_FILE=/etc/kubernetes/pki/ca.crt
ETCD_LISTEN_PEER_URLS=https://192.168.18.5:2380 # 为本集群其他节点提供服务的url地址
ETCD_INITIAL_ADVERTISE_PEER_URLS=https://192.168.18.5:2380 # 广播给集群中其他成员自己为peer提供服务的地址

ETCD_INITIAL_CLUSTER_TOKEN=etcd-cluster # 集群名称
ETCD_INITIAL_CLUSTER="etcd1=https://192.168.18.3:2380,etcd2=https://192.168.18.4:2380,etcd3=https://192.168.18.5:2380" # 集群各节点的endpoint列表
ETCD_INITIAL_CLUSTER_STATE=new # new, 初始集群状态; existing,集群已存在时使用
```