# [minio](https://min.io)
参考:
- [Minio 中文docs](https://docs.min.io/cn)
- [Minio 的 benchmark](https://min.io/resources#benchmarks)

MinIO 是一个基于 Apache License v2.0 开源协议的对象存储服务. 它兼容亚马逊 S3 云存储服务接口，非常适合于存储大容量非结构化的数据，例如图片、视频、日志文件、备份数据和容器/虚拟机镜像等，而一个对象文件可以是任意大小，从几 kb 到最大 5T 不等.

通常minio用来存储的都是一些冷数据，更新不频繁，这种场景下用RS码就显得非常经济实惠了.

默认情况下, MinIO 将对象拆分成N/2数据和N/2 奇偶校验盘, 此时读读仲裁`N/2`, 写读仲裁`N/2+1`.

minio没有快照功能.

## [纠删码](https://docs.min.io/cn/minio-erasure-code-quickstart-guide.html)
Minio使用纠删码erasure code(Reed-Solomon code)和checksum([HighwayHash](https://github.com/minio/highwayhash))来保护数据免受硬件故障和无声数据损坏.

### [在纠删码模式下支持的存储类型](https://github.com/minio/minio/tree/master/docs/zh_CN/erasure/storage-class)
MinIO 支持两种存储类型:
1. 低冗余存储(REDUCED_REDUNDANCY)

    REDUCED_REDUNDANCY存储类型意味着奇偶校验盘比REDUCED_REDUNDANCY少
1. 标准存储(STANDARD)

    STANDARD存储类型意味着奇偶校验盘比REDUCED_REDUNDANCY多.

MinIO可创建每组4到16个磁盘组成的纠删码集合, 所以提供的磁盘总数必须是其中一个数字的**倍数**.

> Minio纠删码是作用在对象级别，可以一次恢复一个对象，而RAID机制是作用在卷级别，数据恢复时间很长.

> Minio纠删码的设计目标是为了性能和尽可能的使用硬件加速.

> **一个EC集合至少4个盘**, 另外minio还有一个最大16的限制，它似乎是minio自己加的，RS码本身要求shard不要超256即可.

## 部署
### 单节点
```bash
MINIO_ACCESS_KEY=minioadmin MINIO_SECRET_KEY=minioadmin ./minio server /mnt/data

mc alias set myminio http://127.0.0.1:9000 minioadmin minioadmin
mc mb myminio/test # 创建buckets, 不能省略alias, 比如这里的myminio, 否则命令也会报创建成功, 但实际没有创建
mc cp a.sh myminio/test # put file
```

### 分布式
**分布式Minio至少需要4个硬盘，使用分布式Minio会自动引入纠删码功能**.

> 分布式Minio里所有的节点需要有同样的access秘钥和secret秘钥，这样这些节点才能建立联接, 且需NTP保障时间一致.

> 建议运行分布式MinIO设置的所有节点都是同构的，即相同的操作系统，相同数量的磁盘和相同的网络互连

> 每个对象被写入一个EC集合中，因此该对象分布在不超过16个磁盘上

部署1(当节点, 四挂载):
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
./minio server --address :9001 http://127.0.0.1:9001/usr/minio/data1 http://127.0.0.1:9001/usr/minio/data2  http://127.0.0.1:9002/usr/minio/data3 http://127.0.0.1:9002/usr/minio/data4 > /usr/minio/minio1.log 2>&1 &
./minio server --address :9002 http://127.0.0.1:9001/usr/minio/data1 http://127.0.0.1:9001/usr/minio/data2  http://127.0.0.1:9002/usr/minio/data3 http://127.0.0.1:9002/usr/minio/data4 > /usr/minio/minio2.log 2>&1 &
```

### [升级 MinIO](https://docs.min.io/cn/)
```bash
mc admin update <minio alias, e.g., myminio>
```

## FAQ
1. `Disk /usr/minio/data1 is a root disk. Please ensure the disk is mounted properly, refusing to use root disk.`
使用了挂载在`/`下的export, 即分布式minio必须使用数据盘.