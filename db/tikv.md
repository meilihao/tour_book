# tikv
参考：
- [TiKV 是如何存取数据的](https://pingcap.com/blog-cn/how-tikv-store-get-data/)

TiKV 提供两套 API，一套叫做 RawKV，另一套叫做 TxnKV. TxnKV 对应的就是[Percolator](https://pingcap.com/blog-cn/how-tikv-store-get-data/)，而 RawKV 则不会对事务做任何保证，而且比 TxnKV 简单很多.

## [api](https://tikv.org/docs/6.1/develop/clients/introduction/)
ref:
- [请问下 TiDB 会使用到 mvcc/txn 接口的场景是什么？](https://asktug.com/t/topic/152923)

TiKV 提供两个可以交互的 API：
|API| 描述| 原子性| 使用场景|
|原始| 低级键值 API，用于直接与各个键值对交互|单个键|需要低延迟，并且不涉及分布式事务|
|事务性| 高级键值 API，用于提供 ACID 语义|多个键|需要分布式事务|

**不支持在同一个键空间上同时使用原始 API 和事务 API**

### [api version](https://docs.pingcap.com/zh/tidb/stable/tikv-configuration-file#api-version-%E4%BB%8E-v610-%E7%89%88%E6%9C%AC%E5%BC%80%E5%A7%8B%E5%BC%95%E5%85%A5)
- keyspace: 需要v2

    使用步骤:
    1. 使用pd-ctl创建keyspace
    1. sdk NewClient()时指定WithKeyspace()

### rawkv api
- Scan: 范围是`[start, end)`

    不考虑limit, `Scan(context.TODO(), nil, nil, rawkv.MAxRawKVScanLimit)`是获取全部数据

## pd-ctl
```bash
region --jq=".regions[] | {id: .id, peer_stores: [.peers[].store_id] | select(length != 3)}" # 副本数不为 3 的所有 Region
```

## FAQ
### titan配置
见tikv data-dir下的last_tikv.toml