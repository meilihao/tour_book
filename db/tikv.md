# tikv
参考：
- [TiKV 是如何存取数据的](https://pingcap.com/blog-cn/how-tikv-store-get-data/)

TiKV 提供两套 API，一套叫做 RawKV，另一套叫做 TxnKV. TxnKV 对应的就是[Percolator](https://pingcap.com/blog-cn/how-tikv-store-get-data/)，而 RawKV 则不会对事务做任何保证，而且比 TxnKV 简单很多.