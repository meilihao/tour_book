# [minio](https://min.io)
参考:
- [Minio 中文docs](https://docs.min.io/cn)
- [Minio 的 benchmark](https://min.io/resources#benchmarks)

MinIO 是一个基于 Apache License v2.0 开源协议的对象存储服务. 它兼容亚马逊 S3 云存储服务接口，非常适合于存储大容量非结构化的数据，例如图片、视频、日志文件、备份数据和容器/虚拟机镜像等，而一个对象文件可以是任意大小，从几 kb 到最大 5T 不等.

通常minio用来存储的都是一些冷数据，更新不频繁，这种场景下用RS码就显得非常经济实惠了.

默认情况下, MinIO 将对象拆分成N/2数据和N/2 奇偶校验盘, 此时读读仲裁`N/2`, 写读仲裁`N/2+1`.

minio没有快照功能.

> minio出品自一个有着多年网络文件系统开发经验的团队，其初始创始团队都来自于原Glusterfs团队.

> [minio支持在bucket级别开启versioning, 作用于其内的每个object](https://docs.min.io/docs/minio-bucket-versioning-guide.html)

类似:
- [seaweedfs](https://github.com/chrislusf/seaweedfs)及[weed-fs使用简介](https://tonybai.com/2015/08/22/intro-of-using-weedfs/)
- [bfs](https://github.com/Terry-Mao/bfs)及[bfs:支撑Bilibili的小文件存储系统](https://mp.weixin.qq.com/s?__biz=MzAwMDU1MTE1OQ==&mid=406016886&idx=1&sn=f5aa286373fb981c9de904568fe7ddb2)

## [纠删码](https://docs.min.io/cn/minio-erasure-code-quickstart-guide.html)
Minio使用纠删码erasure code(Reed-Solomon code)和checksum([HighwayHash](https://github.com/minio/highwayhash))来保护数据免受硬件故障和无声数据损坏.

### [在纠删码模式下支持的存储类型](https://github.com/minio/minio/tree/master/docs/zh_CN/erasure/storage-class)
MinIO 支持两种存储类型:
1. 低冗余存储(REDUCED_REDUNDANCY)

    REDUCED_REDUNDANCY存储类型意味着奇偶校验盘比REDUCED_REDUNDANCY少
1. 标准存储(STANDARD)

    STANDARD存储类型意味着奇偶校验盘比REDUCED_REDUNDANCY多.

MinIO可创建每组4到16个磁盘组成的纠删码集合(推荐是8 from [基于MINIO的对象存储在探探的实践](https://github.com/gopherchina/conference/blob/master/2019/2.1%20%E5%9F%BA%E4%BA%8EMINIO%E7%9A%84%E5%AF%B9%E8%B1%A1%E5%AD%98%E5%82%A8%E6%96%B9%E6%A1%88%E5%9C%A8%E6%8E%A2%E6%8E%A2%E7%9A%84%E5%AE%9E%E8%B7%B5%20-%20%E4%BA%8E%E4%B9%90.pdf)), 所以提供的磁盘总数必须是其中一个数字的**倍数**.

> Minio纠删码是作用在对象级别，可以一次恢复一个对象，而RAID机制是作用在卷级别，数据恢复时间很长.

> Minio纠删码的设计目标是为了性能和尽可能的使用硬件加速.

> **一个EC集合至少4个盘**, 另外minio还有一个最大16的限制，它似乎是minio自己加的，RS码本身要求shard不要超256即可.

## minio部署
### 单节点
```bash
MINIO_ACCESS_KEY=minioadmin MINIO_SECRET_KEY=minioadmin ./minio server /mnt/data

mc config host add minio http://172.18.100.177:9000 minioadmin minioadmin --api s3v4 # minio即为要添加的alias
mc alias set myminio http://127.0.0.1:9000 minioadmin minioadmin
mc mb myminio/test # 创建buckets, 不能省略alias, 比如这里的myminio, 否则命令也会报创建成功, 但实际没有创建
mc cp a.sh myminio/test # put file
```

### 分布式
**分布式Minio至少需要4个硬盘，使用分布式Minio会自动引入纠删码功能**.

> 分布式Minio里所有的节点需要有同样的access秘钥和secret秘钥，这样这些节点才能建立联接, 且**需NTP保障时间一致**.

> 建议运行分布式MinIO设置的所有节点都是同构的，即相同的操作系统，相同数量的磁盘和相同的网络互连

> 每个对象被写入一个EC集合中，因此该对象分布在不超过16个磁盘上

部署1(单节点, 四挂载):
```bash
export MINIO_ACCESS_KEY="minio"
export MINIO_SECRET_KEY="minio123"
cd /usr/minio
mkdir data{1..4}
nohup ./minio server --address ":9001" http://127.0.0.1:9001/usr/minio/data1 http://127.0.0.1:9002/usr/minio/data2  http://127.0.0.1:9003/usr/minio/data3 http://127.0.0.1:9004/usr/minio/data4 > "/usr/minio/9001.log" 2>&1 & # /usr/minio/data1是export路径
nohup ./minio server --address ":9002" http://127.0.0.1:9001/usr/minio/data1 http://127.0.0.1:9002/usr/minio/data2  http://127.0.0.1:9003/usr/minio/data3 http://127.0.0.1:9004/usr/minio/data4 > "/usr/minio/9002.log" 2>&1 &
nohup ./minio server --address ":9003" http://127.0.0.1:9001/usr/minio/data1 http://127.0.0.1:9002/usr/minio/data2  http://127.0.0.1:9003/usr/minio/data3 http://127.0.0.1:9004/usr/minio/data4 > "/usr/minio/9003.log" 2>&1 &
nohup ./minio server --address ":9004" http://127.0.0.1:9001/usr/minio/data1 http://127.0.0.1:9002/usr/minio/data2  http://127.0.0.1:9003/usr/minio/data3 http://127.0.0.1:9004/usr/minio/data4 > "/usr/minio/9004.log" 2>&1 &
```

部署2(双节点、双挂载):
```bash
export MINIO_ACCESS_KEY="minio"
export MINIO_SECRET_KEY="minio123"
cd /usr/minio/
# --- peer_id都应能被对方访问到, local_ip是否能用127.0.0.1待测试???
# --- on local
./minio server --address :9001 http://<local_ip>:9001/usr/minio/data1 http://<local_ip>:9001/usr/minio/data2  http://<peer_id>:9002/usr/minio/data3 http://<peer_id>:9002/usr/minio/data4 > /usr/minio/minio1.log 2>&1 &
# --- on peer
./minio server --address :9002 http://<peer_id>:9001/usr/minio/data1 http://<peer_id>:9001/usr/minio/data2  http://<local_ip>:9002/usr/minio/data3 http://<local_ip>:9002/usr/minio/data4 > /usr/minio/minio2.log 2>&1 &
```

### [升级 MinIO](https://docs.min.io/cn/)
```bash
mc admin update <minio alias, e.g., myminio>
```

### mc
```bash
# 1. 性能采集
mc admin profile start -type cpu # possible values are 'cpu', 'mem', 'block', 'mutex', 'trace', 'threads' and 'goroutines' (default: "cpu,mem,block")
mc admin profile stop
pprof -http=0.0.0.0:2222 ./profiling-10.189.33.60\:9000.pprof #  使⽤ google pprof ⼯具进⾏可视化展示
mc watch local/rclone/source #  local是host配置, `rclone/source`是path, event type仅[这里](https://docs.min.io/docs/minio-bucket-notification-guide.html)
```

## source
### 基本概念
- drive : 简单理解为一块硬盘
- set : 一组drive的集合

    一个对象存储在一个set上
    一个集群可划分为多个set
    一个set包含的drive数量是固定的
        - 默认由系统根据集群规模自动计算得出
        - MINIO_ERASURE_SET_DRIVE_COUNT
    一个set中的drive尽可能分布在不同的节点上

### code
MinIO对外提供http接口的入口相关文件都在[`/cmd`](https://github.com/minio/minio/tree/master/cmd).

Object handle的接口是`cmd/object-api-interface.go#ObjectLayer`, 当前有两个实现(`cmd/server-main.go#newObjectLayer()`):
- NewFSObjectLayer(), 单disk
- newErasureServerPools(), 基于ec

`minio server`的调用链路是:
1. server-main.go#serverMain
1. [routers.go#configureServerHandler(endpointServerPools EndpointServerPools)](https://github.com/minio/minio/blob/master/cmd/routers.go#L86)
1. [api-router.go#registerAPIRouter(router *mux.Router)](https://github.com/minio/minio/blob/master/cmd/api-router.go#L82), 它注册了所有http相关的url处理函数, url分发由`github.com/gorilla/mux`处理.

比如PutObject:
```go
        // https://github.com/minio/minio/blob/master/cmd/api-router.go#L181
        // PutObject
        bucket.Methods(http.MethodPut).Path("/{object:.+}").HandlerFunc(
            collectAPIStats("putobject", maxClients(httpTraceHdrs(api.PutObjectHandler))))

        // https://github.com/minio/minio/blob/master/cmd/api-router.go#L169
        // GetObject
        bucket.Methods(http.MethodGet).Path("/{object:.+}").HandlerFunc(
            collectAPIStats("getobject", maxClients(httpTraceHdrs(api.GetObjectHandler))))
```

[api.PutObjectHandler](https://github.com/minio/minio/blob/master/cmd/object-handlers.go#L1311)函数是实现如何把一个对象上传到一个桶里.

[api.GetObjectHandler](https://github.com/minio/minio/blob/master/cmd/object-handlers.go#L302)函数是实现如何从一个桶中获取一个对象.

### ErasureCoding
ec代码在`cmd/erasure-*.go`, 入口在`cmd/erasure-coding.go#NewErasure`

## FAQ
### 1. `Disk /usr/minio/data1 is a root disk. Please ensure the disk is mounted properly, refusing to use root disk.`
使用了挂载在`/`下的export, 即分布式minio必须使用数据盘.

### Waiting for the first server to format the disks
清空所有组成minio cluster的盘(包括其下的隐藏文件)再重试即可.