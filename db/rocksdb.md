# rocksdb

RocksDB的目的是成为一套能在服务器压力下，真正发挥高速存储硬件（特别是Flash 和 RAM）性能的高效数据库系统. 它是一个C++库，允许存储任意长度二进制kv数据, 支持原子读写操作.

RocksDB大量复用了levedb的代码，并且还借鉴了许多HBase的设计理念, 同时Rocksdb也借用了一些Facebook之前就有的理念和代码.

RocksDB是一个嵌入式的K-V（任意字节流）存储。所有的数据在引擎中是有序存储，可以支持Get(key)、Put（Key）、Delete（Key）和NewIterator()。RocksDB的基本组成是memtable、sstfile和logfile。memtable是一种**内存数据结构**，写请求会先将数据写到memtable中，然后可选地写入事务日志logfile。logfile是一个**顺序写**的文件。当memtable被填满的时候，数据会flush到sstfile中，然后这个memtable对应的logfile也会安全地被删除。sstfile中的数据也是有序存储以方便查找。

RocksDB支持将一个数据库实例分片为多个列族(column families). 类似HBase，每个DB新建时默认带一个名为"default"的列族，如果一个操作没有携带列族信息，则默认使用这个列族. 如果WAL开启，当实例crash再恢复时，RocksDB可以保证用户一个一致性的视图. 通过WriteBatch API，可以实现跨列族操作的原子性.

所有 Column Family 共享一个 WAL 文件，但是每个 Column Family 有自己单独的 memtable & ssttable(sstfile)，即 log 共享而数据分离.

## 编译
1. 参照[rocksdb INSTALL](https://github.com/facebook/rocksdb/blob/master/INSTALL.md), 选择平台安装依赖lib
1. `cd rocksdb_source_root`, 选择`make static_lib/make shared_lib`进行编译
1. 参考rocksdb的Makefile, 再执行`make install-static/make install-shared`即可. 如果安装位置需要还可使用`INSTALL_PATH=/usr/local make install-static/install-shared`, `INSTALL_PATH`默认已是`/usr/local`, 最终`librocksdb.a/librocksdb.so`会出现在`$INSTALL_PATH/lib`下
1. 设置环境变量

	```bash
	# vim ~/.bashrc
	export CPLUS_INCLUDE_PATH=${CPLUS_INCLUDE_PATH}:$INSTALL_PATH/include
	export LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:$INSTALL_PATH/lib
	export LIBRARY_PATH=${LIBRARY_PATH}:$INSTALL_PATH/lib
	```
1. 重启terminal即可

## 文件介绍
- *.log: 事务日志用于保存数据操作日志，可用于数据恢复
- *.sst: 数据持久换文件
- MANIFEST：数据库中的 MANIFEST 文件记录数据库状态即某个时刻SST文件的视图. 压缩过程会添加新文件并从数据库中删除旧文件，并通过将它们记录在 MANIFEST 文件中使这些操作持久化.
- CURRENT：记录当前正在使用的MANIFEST文件
- LOCK：rocksdb自带的文件锁，防止两个进程来打开数据库

## Memtable
可插拔 memtable，RocksDB 的 memtable 的默认实现是一个 skiplist。skiplist 是一个有序集，当工作负载使用 range-scans 并且交织写入时，这是一个必要的结构。然而，一些应用程序不交织写入和扫描，而一些应用程序根本不执行范围扫描。对于这些应用程序，排序集可能无法提供最佳性能。因此，RocksDB 支持可插拔的 API，允许应用程序提供自己的 memtable 实现。开发库提供了三个 memtable：skiplist memtable，vector memtable 和前缀散列（prefix-hash） memtable。Vector memtable 适用于将数据批量加载到数据库中。每个写入在向量的末尾插入一个新元素; 当它是刷新 memtable 到存储的时候，向量中的元素被排序并写出到 L0 中的文件。前缀散列 memtable 允许对 gets，puts 和 scans-within-a-key-prefix 进行有效的处理。

## SSTFile(SSTTable)
RocksDB在磁盘上的file结构sstfile由block作为基本单位组成，一个sstfile结构由多个data block和meta block组成， 其中data block就是数据实体block，meta block为元数据block， 其中data block就是数据实体block，meta block为元数据block。 sstfile组成的block有可能被压缩(compression)，不同level也可能使用不同的compression方式。 sstfile如果要遍历block，会逆序遍历，从footer开始。

写流程：
rocksdb写入时，直接以append方式写到log文件以及memtable，随即返回，因此非常快速. memtable/immute memtable触发阈值后， flush 到Level0 SST，Level0 SST触发阈值后，经合并操作(compaction)生成level 1 SST， level1 SST 合并操作生成level 2 SST，以此类推，生成level n SST.

读流程：
按照 memtable --> Level 0 SST–> Level 1 SST --> … -> Level n SST的顺序读取数据, 这和记录的新旧顺序是一的. 因此只要在当前级别找到记录，就可以返回.

## RocksDB的典型场景（低延时访问）:
1. 需要存储用户的查阅历史记录和网站用户的应用
1. 支持大量写和删除操作的消息队列