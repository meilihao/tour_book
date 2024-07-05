# rocksdb
参考:
- [RocksDB系列](https://www.jianshu.com/p/061927761027)
- [RocksDB 第6课](https://www.modb.pro/db/385068)
- [Rocksdb基本用法](https://www.cnblogs.com/wanshuafe/p/11564148.html)
- [Basic Operations](https://github.com/facebook/rocksdb/wiki/Basic-Operations)
- [漫谈RocksDB(四)存储结构](https://www.modb.pro/db/112483)
- [Tuning RocksDB - Statistics](https://www.jianshu.com/p/ddf652aa4882)

    开启statistics会增加5%~10%的额外开销
- [rocksdb参数](https://tikv.org/docs/6.1/deploy/configure/tikv-configuration-file/#rocksdb)
- [Apache Flink中的RocksDB状态后端](https://zhuanlan.zhihu.com/p/332484994)
- [rocksdb-doc-cn](https://wanghenshui.github.io/rocksdb-doc-cn/)
- [WiscKey 发布的五年后，工业界用上了 KV 分离吗？](https://zhuanlan.zhihu.com/p/397466422)
- [TIDB TIKV数据存储到ROCKSDB探秘 与 ROCKSDB 本尊](https://cloud.tencent.com/developer/article/1857152)
- [rocksdb/USERS.md](https://github.com/facebook/rocksdb/blob/main/USERS.md)
- [RocksDB 笔记](https://blog.csdn.net/qq_32907195/article/details/117933955)
- [rocksdb-doc-cn](https://github.com/johnzeng/rocksdb-doc-cn)

其他软件使用的rocksdb版本:
- [apache/kvrocks](https://github.com/apache/kvrocks/blob/unstable/cmake/rocksdb.cmake)
- [facebook/mysql-5.6](https://github.com/facebook/mysql-5.6)
- [flink-state-backends](https://github.com/apache/flink/blob/master/flink-state-backends/flink-statebackend-rocksdb/pom.xml)

    1.18.0->6.20.3

> go+[badger](https://github.com/dgraph-io/badger)也是不错的选择, 特别是cgo问题无法解决的时候

> 推荐使用xxx.fb/xxx.fb.myrocks分支的代码, 下载方式: `https://codeload.github.com/facebook/rocksdb/tar.gz/refs/heads/xxx.fb`

RocksDB的目的是成为一套能在服务器压力下，真正发挥高速存储硬件（特别是Flash 和 RAM）性能的高效单点数据库系统. 它是一个C++库，允许存储任意长度二进制kv数据, 支持原子读写操作, 因此本质上来说它是一个可插拔式的存储引擎选择.

> rocksdb解决的是写多读少的场景需求; B+解决的是读多写少.

RocksDB大量复用了levedb的代码，并且还借鉴了许多HBase的设计理念, 同时Rocksdb也借用了一些Facebook之前就有的理念和代码.

RocksDB是一个嵌入式的K-V（任意字节流）存储. 所有的数据在引擎中是有序存储，可以支持Get(key)、Put（Key）、Delete（Key）和NewIterator()。RocksDB的基本组成是memtable、sstfile和logfile。memtable是一种**内存数据结构**，写请求会先将数据写到memtable中，然后可选地写入事务日志logfile(WAL)。logfile是一个**顺序写**的文件。当memtable被填满的时候，数据会flush到sstfile中，然后这个memtable对应的logfile也会安全地被删除。sstfile中的数据也是有序存储以方便查找.

> rocksdb支持Direct IO, 以绕过系统Page Cache，通过应用内存从存储设置中直接进行IO读写操作.

> rocksdb版本定义在`include/rocksdb/version.h`, history在`HISTORY.md`.

![](/misc/img/develop/4304640-891400b1777c999d.png)

## 编译
1. 参照[rocksdb INSTALL](https://github.com/facebook/rocksdb/blob/master/INSTALL.md), 选择平台安装依赖lib

    ```bash
    # apt install libgflags-dev libsnappy-dev zlib1g-dev libbz2-dev liblz4-dev  libzstd-dev libjemalloc-dev
    # make -j4 shared_lib # 调试模式用make dbg(LIB_MOD=shared) from `INSTALL.md`, make all也可以(LIB_MOD=shared)
    # make install-shared
    # ldconfig
    ```

    > 编译7.10.2发现gcc-c++需要支持c++17, 推荐gcc8, 仅[gcc](/shell/cmd/compile/gcc.md)

    > rocksdb的cmake是针对windows 64-bit的, 见CMakeLists.txt

    > 构建rocksdb时去除`PORTABLE=1`, 可能使得gdb coredump时获取更多信息
1. `cd rocksdb_source_root`, 查看Makefile, 选择`make static_lib/make shared_lib`进行编译

    如果构建环境存在jemalloc/tcmalloc, make会通过`build_tools/build_detect_platform <platform>`将相应的环境变量存入生成的make_config.mk中, 供自身使用

    > [`ROCKSDB_DISABLE_JEMALLOC=1 make shared_lib`的ROCKSDB_DISABLE_JEMALLOC可禁用jemalloc](https://github.com/facebook/rocksdb/issues/1442), 此时tcmalloc已安装(libtcmalloc-minimal4, 是libtcmalloc_minimal.so)也并不会启用tcmalloc, 是build_detect_platform没有探测到需要的`libtcmalloc.so`, 其实启用tcmalloc需要`atp install libgoogle-perftools-dev(libtcmalloc.so在libgoogle-perftools4里)/yum install gperftools gperftools-devel`, 禁用tcmalloc可用`ROCKSDB_DISABLE_TCMALLOC=1`.
1. 参考rocksdb的Makefile, 再执行`make install-static/make install-shared`即可. 如果安装位置需要还可使用`INSTALL_PATH=/usr/local make install-static/install-shared`, `INSTALL_PATH`默认已是`/usr/local`, 最终`librocksdb.a/librocksdb.so`会出现在`$INSTALL_PATH/lib`下
1. 设置环境变量

	```bash
	# vim ~/.bashrc
	export CPLUS_INCLUDE_PATH=${CPLUS_INCLUDE_PATH}:$INSTALL_PATH/include
	export LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:$INSTALL_PATH/lib
	export LIBRARY_PATH=${LIBRARY_PATH}:$INSTALL_PATH/lib
	```
1. 重启terminal即可

> 编译rocksdb源码下的examples: `g++ simple_example.cc -o test -std=c++11 -lpthread -lrocksdb -ldl -lrt -lsnappy -lgflags -lz -lbz2 -lzstd`

## 文件介绍
- *.log: 事务日志用于保存数据操作日志，可用于数据恢复
- LOG/LOG.old.* : 当前log/历史log

     `options.keep_log_file_num` : 可指定log文件个数(>=1, 因为当前LOG文件正在被使用)
     `options.max_log_file_size` : 日志轮转大小
     `options.log_file_time_to_roll` : 日志轮转时间
- *.sst: 数据持久换文件
- MANIFEST：数据库中的 MANIFEST 文件记录数据库状态即某个时刻SST文件的视图. 压缩过程会添加新文件并从数据库中删除旧文件，并通过将它们记录在 MANIFEST 文件中使这些操作持久化.

	RocksD有一个内建的机制来克服POSIX文件系统的各种限制(比如操作不是原子的, 不支持批量操作等)，这种机制就是通过一个MANIFEST文件记录**RocksDB状态改变的所有事物日志**.所以，MANIFEST文件可以在DB重启时恢复到最近一次的一致性状态.

	在RocksDB中任意时间存储引擎的状态都会保存为一个Version(也就是SST的集合)，而每次对Version的修改都是一个VersionEdit,而最终这些VersionEdit就是 组成manifest-log文件的内容.
- MANIFEST-* : Manifest的滚动日志文件
- CURRENT：指定当前正在使用的MANIFEST文件
- LOCK：rocksdb自带的文件锁，防止两个进程来打开数据库
- IDENTITY : 存放当前rocksdb的唯一标识

## Memtable
> rocksdb内存部分就地写; 磁盘部分追加写.

可插拔 memtable，RocksDB 的 memtable 的默认实现是一个 skiplist。skiplist 是一个有序集，当工作负载使用 range-scans 并且交织写入时，这是一个必要的结构。然而，一些应用程序不交织写入和扫描，而一些应用程序根本不执行范围扫描。对于这些应用程序，排序集可能无法提供最佳性能。因此，RocksDB 支持可插拔的 API，允许应用程序提供自己的 memtable 实现。

开发库提供了三个 memtable：skiplist memtable，vector memtable 和前缀散列（prefix-hash） memtable
- Vector memtable 适用于将数据批量加载到数据库中。每个写入在向量的末尾插入一个新元素; 当它是刷新 memtable 到存储的时候，向量中的元素被排序并写出到 L0 中的文件
- 前缀散列 memtable 允许对 gets，puts 和 scans-within-a-key-prefix 进行有效的处理。

只读的 MemTable和 MemTable 的数据结构完全一样，唯一的区别就是不允许再写入了.

内存结构选择:
1. 红黑树

    并发需要锁整棵树
2. 调表

    并发只需要锁最大高度和节点

参数:
- write-buffer-size = "128MB" # RocksDB memtable 的大小
- max-write-buffer-number = 5 # 最多允许几个 memtable 存在

## SSTFile(SSTTable)
> sst单个文件没有重复
> level0同层可能有重复; level1-level N 同层没有重复, 因为其由上层合并而来, 但层间可能有重复


RocksDB在磁盘上的file结构sstfile由block作为基本单位组成，一个sstfile结构由多个data block和meta block组成， 其中data block就是数据实体block，meta block为元数据block， 其中data block就是数据实体block，meta block为元数据block。 sstfile组成的block有可能被压缩(compression)，不同level也可能使用不同的compression方式。 sstfile如果要遍历block，会逆序遍历，从footer开始。

sst里面的数据已**按照key进行排序**能方便对其进行二分查找. 在SST文件内，还额外包含以下特殊信息：
- Bloom Fileter : 用于快速判断目标查询key是否存在于当前SST文件内
- Index / Partition Index，SST内部数据块索引文件快速找到数据块的位置

compaction输入的SST file并不是立即就从SST file集合中删除，因为有可能在这些SST file上正进行着get or iterator操作. 只有当冗余的SST file上没有任何操作的时候，才会执行真正的删除文件操作. [这些逻辑是通过引用计数来实现的](https://www.jianshu.com/p/b95db752178f).

写流程：
rocksdb写入时，直接以append方式写到log文件, 成功后应用到memtable，随即返回，因此非常快速. memtable/immute memtable触发阈值后， flush 到Level0 SST，Level0 SST触发阈值后，经合并操作(compaction)生成level 1 SST， level1 SST 合并操作生成level 2 SST，以此类推，生成level n SST.

![rocksdb 写入原理](/misc/img/rocksdb/1f950c8b30fe6992437242c368f0f8b1.png)

Get()流程：
1. 在MemTable中查找，无法命中转到下一流程
1. 在immutable_memtable中查找，查找不中转到下一流程
1. 在第0层SSTable中查找，无法命中转到下一流程

	对于L0 的文件，RocksDB 采用遍历的方法查找，所以为了查找效率 RocksDB 会控制 L0 的文件个数

1. 在剩余SSTable中查找

    对于 L1 层以及 L1 层以上层级的文件，每个 SSTable 没有交叠，可以使用二分查找快速找到 key 所在的 Level 以及 SSTfile.

> 如果启用 Level Style Compaction, L0 存储着 RocksDB 最新的数据，Lmax 存储着比较老的数据，L0 里可能存着重复 keys，但是其他层文件则不可能存在重复 key.

参数:
- target-file-size-base = "32MB" # sst 文件的大小

## RocksDB的典型场景（低延时访问）:
1. 需要存储用户的查阅历史记录和网站用户的应用
1. 支持大量写和删除操作的消息队列

## 衍生版
- Pika: 解决大数据量下, redis启动慢

## 功能
### Column Families
RocksDB支持将一个数据库实例分片为多个列族(column families, 类似表Table). 类似HBase，每个DB新建时默认带一个名为"default"的列族，如果一个操作没有携带列族信息，则默认使用这个列族. 如果WAL开启，当实例crash再恢复时，RocksDB可以保证用户一个一致性的视图. 通过WriteBatch API，可以实现跨列族操作的原子性.

所有 Column Family 共享 WAL、Current、Manifest 文件，但是每个 Column Family 有自己单独的 memtable & ssttable(sstfile)，即 log 共享而数据分离.

> memtable默认是基于 Skip-List 跳表实现的, 也支持HashSkipList.

CF 提供了对 DB 进行逻辑划分开来的方法，用户可以通过 CF 同时对多个 CF 的 KV 进行并行读写的方法，提高了并行度.
### Updates
如果k已经存在的话，则已有的v会被新的v覆盖

RocksDB 的写是异步的：仅仅把数据写进了操作系统的缓存区就返回了，而这些数据被写进磁盘是一个异步的过程. 如果为了数据安全，可以用`write_options.sync = true`把写过程改为同步写(异步写的吞吐率是同步写的一千多倍).

WriteBatch 默认使用了事务，确保批量写成功.

Merge接口是修改现有值(get->put/delete)的原子操作.

MergeOperator还可以用于非关联型数据类型的更新, 比如json.
```c++
 // Put/store the json string into to the database
    db_->Put(put_option_, "json_obj_key",
             "{ employees: [ {first_name: john, last_name: doe}, {first_name: adam, last_name: smith}] }");
    // Use a pre-defined "merge operator" to incrementally update the value of the json string
    db_->Merge(merge_option_, "json_obj_key", "employees[1].first_name = lucy");
    db_->Merge(merge_option_, "json_obj_key", "employees[0].last_name = dow");
```

### Gets、Iterators、Snapshots
RocksDB中的key和value完全是byte stream，key和value的大小没有任何限制. Get接口提供用户一种从DB中查询key对应value的方法，MultiGet提供批量查询功能. DB中的所有数据都是按照key有序存储，其中key的compare方法可以用户自定义.

Iterator方法提供用户RangeScan功能，首先seek到一个特定的key，然后从这个点开始遍历. Iterator也可以实现RangeScan的逆序遍历，当执行Iterator时，用户看到的是一个时间点的一致性视图.

Snapshot接口可以创建数据库在某一个时间点的快照. Get和Iterator接口也可以执行在某一个Snapshot上. 某种意义上，Iterator和Snapshot提供了DB在某个时间点的一个一致性视图，但是其实现原理却不一样.

快速短期/前台的scan操作比较适合用Iterator，长期/后台操作适合用Snapshot. 当使用Iterator时，会对数据库相应时间点的所有底层文件增加引用计数，直到Iterator结束或者释放了引用计数后，这些文件才允许被删除. Snapshot不关注数据文件是否被删除的问题，Compation进程会感知Snapshot的存在，会保证对应视图的数据不会被删除. 当实例重启时，Snapshot会丢失，这是因为RocksDB不会持久化Snapshot相关数据.

> RocksDB 自身会给 key 和 value 添加一个 C-style 的 `\0`，所以 slice 的指针指向的内存区域自身作为字符串输出没有问题.

### Transations
RocksDB提供了多个操作的事务性，支持悲观和乐观模式

### Prefix Iterator
大部分的LSM引擎都不支持高效的RangeScan操作，这是由于执行RangeScan操作时都要访问所有的数据文件导致。但是大部分用户并不仅仅是完全scan所有的数据，相反，很多情况下仅仅需要按照key的前缀字符串区遍历。RocksDB根据这些应用场景，优化了对应的底层实现。用户可以prefix_extractor来声明一个key_prefix，然后RocksDB为每一个key_prefix存储相应的blooms。配置了key_prefix的Iterator操作可以通过对应的bloom bits来避免检索不含有特定key prefix的数据文件，依次可以提高Iterator性能.

[Prefix seek](https://www.jianshu.com/p/9848a376d41d)是RocksDB的一种模式，主要影响Iterator的行为. 这种模式下，RocksDB的Iterator并不保证所有key是有序的，而只保证具有相同前缀的key是有序的. 这样可以保证具有相同特征的key（例如具有相同前缀的key）尽量地被聚合在一起.

SliceTransform中的transform 就是提取 key 的 prefix. in_domain 用来判断这个 key 是否符合提取的要求，如果返回了 true，则表明可以使用 transform 提取 prefix 并插入到 bloom filter 里面.

### Persistence
RocksDB有事物日志，所有的写操作首先写入内存表内，然后可选地写入到事物日志中。当DB重启时会重新执行事物日志中的所有操作，然后恢复到特定的数据状态。事物日志数据可以与DB数据文件配置成不同的目录下，这种情况适用于将数据文件写到一致性、性能高的快存中，同时可以将事物日志保存在读写性能相对比较慢的持久化存储上来保证数据的安全性。当写数据时可以配置WriteOption,来支持是否将写操作记录在事物日志中或者当用户执行commit时是否需要执行事物日志记录的sync操作。

> 一个WAL文件只有当所有的列族数据都已经flush到SST file之后才会被删除.

RocksDB中每一个提交的记录都是持久化的, 没有提交的记录保存在WAL  file中. 当DB正常退出时，在退出之前会提交所有没有提交的数据，所以总是能够保证一致性. 当RocksDB进程被kill或者服务器重启时，RocksDB需要恢复到一个一致性状态, 其中最重要的恢复操作之一就是replay所有WAL中没有提交的记录.

### Fault Torlerance
RocksDB通过checksum来检测磁盘数据损坏。每个sst file的数据块（4k-128k）都有相应的checksum值。写入存储的数据块内容不允许被修改。

### Multi-Threaded Compactions
当用户重复写入一个key时，在DB中会存在这个key的多个value，compaction操作就是来删除这个key的冗余数据。当一个key被删除时，compation也可以用来真正执行这个底层数据的删除工作，如果用户配置合适的话，compation操作可以多线程执行。DB的数据都存储在sstfile中，当内存表的数据满的时候，会将内存数据（去重、删除无效数据后）写入到L0 文件中。每隔一段时间小文件中的数据会重新merge到更大的文件中，这就是compation。LSM引擎的写吞吐直接依赖于compation的性能，特别是数据存储在SSD或者RAM的情况。RocksDB也支持多线程并行compaction.

### Avoiding Stalls
后台的compaction线程用来将内存数据flush到存储，当所有的后台线程都正在执行compaction时，瞬时大量写操作会很快将内存表写满，这就会引起写停顿。可以配置少一些的线程用于执行数据flush操作

### [Full Backups, Incremental Backups](https://www.jianshu.com/p/85b7610a73bf) and Replication
RocksDB支持增量备份，增量复制需要能够查找到所有的DB修改记录。GetUpdatesSince接口可以提供tail DB transction log的功能。RocksDB的tranction log记录在数据库目录中，当日志文件不再需要时就会move到归档目录。归档目录之所以存在是因为数据复制流比较落后时有可能需要检索过去某一个时间点的日志。GetSortedWalFiles可以返回所有的transction log文件列表.

正常情况下，backup数据是递增的. 开发者可以使用BackupEngine::CreateNewBackup() 创建一个新的backup，且只有新增的数据才会copy到backup 目录中.

### Block Cache -- Compressed and Uncompressed Data
RocksDB使用LRU cache提供block的读服务, 存储SST文件被经常访问的热点数据. block cache partition为两个独立的cache，其中一块可以cache未压缩RAM数据，另一块cache 压缩RAM数据。如果压缩cache配置打开的话，用户一般会开启direct io，以避免OS的也缓存重新cache相同的压缩数据。

### Table Cache
Table cache缓存了所有已打开的文件句柄，这些文件都是sstfile。用户可以设置table cache的最大值。

### Merge Operator
RocksDB原生地就支持三种记录类型，分别为Put、Delete和Merge。Merge可以合并多个Put和Merge记录为一个单独的记录

### Time to Live
开启ttl时, 每个 kv 被插入数据文件中的时候会带上创建时的机器 (int32_t)Timestamp 时间值, 但仅在compaction时, 如果 kv 满足条件`Timestamp+ttl<time_now`则会被淘汰掉.

## 数据结构
### Option
RocksDB通过Options类将配置信息传入引擎. 除此之外，还可以以下其他方法设置，分别为：
1. 通过option file生成一个option class
1. 从option string中获取option 信息

	每个option信息在option string中以`<option_name>:<option_value>`传入，多个option之间以`;`分割.
	开发者可以调用`GetXXXOptionsFromString()`解析option string.
1. 从string map中获取option信息

	开发者可通过`GetXXXOptionsFromMap()`解析option信息.

### Bloom Filter
在任意的keys集合中，应用一个算法并生成一个字节数组，这个字节数组就是Bloom filter. 对于任意一个key，通过Bloom filter可以得出两个结论：
1. 这个key有可能在集合中
1. 这个key肯定不在集合中

在RocksDB引擎中，如果设置了filter policy的话，每个新创建的SST file都会包含一个Bloom filter，这个Bloom filter可以确定要查找的key是否有可能在这个SST file中.

### 其他结构
#### [Block Cache](https://www.jianshu.com/p/64ff46550ee5)
Block Cache是RocksDB把数据缓存在内存中以提高读性能的一种方法. RocksDB中有两种cache的实现方式，分别为LRUCache和CLockCache. 这两种cache都会被分片，来降低锁压力.

默认情况下，会对key的所有字节进行hash计算来设置bloom filter。这可以通过设置BlockBasedTableOptions::whole_key_filtering为false来避免对全部字节进行计算。当Options.prefix_extractor设置后，针对每个key的前缀计算的hash值也添加到了bloom filter中. 由于key的前缀集合要小于key集合，因此计算key前缀生成的bloom filter会更小，当然也会提高误报率.

## 工具
参考:
- [Administration and Data Access Tool](https://www.jianshu.com/p/35a5d5792d65)

RocksDB提供以下3大类型的工具:
1. 性能测试工具

    Benchmark Tool
    Stress Tool，压力测试工具

    ```bash
    $ cd rocksdb
    $ make BENCHMARKS
    ```

    > db_bench可用`--db=xxx`指定db path

1. workload模拟工具

    用户数据访问行为模拟工具
    Workload生成工具

    [ldb](https://github.com/facebook/rocksdb/wiki/Administration-and-Data-Access-Tool)命令行工具提供了不同的数据访问和数据库管理命令.
    [sst_dump tool](https://github.com/facebook/rocksdb/wiki/Administration-and-Data-Access-Tool)可以dump数据然后分析SST file.

    ```bash
    # Linux
	$ cd rocksdb
	$ make tools
	$ cp ldb /usr/local/bin/
	$ cp sst_dump /usr/local/bin/
    ```

1. 性能分析工具，DB Analyzer

## API
- DeleteRange : 范围是`[)`

## 源码
- [RocksDB · 数据的读取(一)](http://mysql.taobao.org/monthly/2018/11/05/)
- [RocksDB解析](https://www.cnblogs.com/pdev/p/11277784.html)

## 组件
### Iterator
rocksdb迭代器为用户提供了大量的便捷操作和接口访问方式:
- NewIterator 创建一个迭代器，需要传入读配置项
- Seek 查找一个key
- SeekToFirst 迭代器移动到db的第一个key位置，一般用于顺序遍历整个db的所有key
- SeekToLast 迭代器移动到db的最后一个key位置， 一般用于反向遍历整个db的所有key
- SeekForPrev移动到当前key的上一个位置，一般用于遍历(limit, start]之间的key
- Next 迭代器移动到下一个key
- Prev迭代器移动到上一个key

## example
```go
options.SetPrefixExtractor(gorocksdb.NewFixedPrefixTransform(1))
db,err:=gorocksdb.OpenDb(options, "rdb")

opts:=gorocksdb.NewDefaultReadOptions()
opts.SetPrefixSameAsStart(true) // 确保下面it返回值的前缀都是'b'
opts.SetTotalOrderSeek(true) // cdp大数据量场景下, 该参数+`Seek()/Iterator.Next()`很慢, 特别是no space时
defer opts.Destroy()

it:=db.NewIterator()
defer it.Close()

it.Seek([]byte{'b'}) // Seek的参数超过NewFixedPrefixTransform设定时会按FixedPrefixTransform截取. 按前缀查找需设置SetTotalOrderSeek=true.
```

```go
// 查找最后一个可用[SetIterateUpperBoundin](https://stackoverflow.com/questions/60576749/rocksdb-iterator-seek-until-last-matching-prefix)
ro := NewDefaultReadOptions()
ro.SetIterateUpperBound([]byte("keya")) // 全尺寸key, targetkey < "keya"

iter, count := db.NewIterator(ro), 0
for iter.SeekToFirst(); iter.Valid(); iter.Next() {
    count++
}

ro.Destroy()
iter.Close()
```

## 调优
ref:
- [TiKV 配置文件描述](https://docs.pingcap.com/zh/tidb/v7.4/tikv-configuration-file#rocksdb)

### Block Cache 系列参数
ref:
- [Flink on RocksDB 参数调优指南](https://cloud.tencent.com/developer/article/1592441)
- [Apache Flink中的RocksDB状态后端](https://zhuanlan.zhihu.com/p/332484994)
- [Kvrocks 在 RocksDB 上的优化实践](https://blog.csdn.net/weixin_45583158/article/details/122227628)
- [RocksDB Tuning Guide](https://github.com/facebook/rocksdb/wiki/RocksDB-Tuning-Guide)
- [Advanced RocksDB State Backends Options](https://nightlies.apache.org/flink/flink-docs-release-1.16/docs/deployment/config/#advanced-rocksdb-state-backends-options)

Block 块是 RocksDB 保存在磁盘中的 SST 文件的基本单位，它包含了一系列有序的 Key 和 Value 集合，可以设置固定的大小. 经过[我们的长期验证](https://cloud.tencent.com/developer/article/1592441)，发现 Block Size 及其 Cache 的大小设置，对读写性能的影响最大.

- block_size

    block 的大小，默认值为4KB。在生产环境中总是会适当调大一些，一般16~32KB比较合适，对于机械硬盘可以再增大到128~256KB，充分利用其顺序读取能力。但是需要注意，如果 block 大小增大而 block cache 大小不变，那么缓存的 block 数量会减少，无形中会增加读放大.
- block_cache_size

    block cache 的大小，默认为8MB。由上文所述的读写流程可知，较大的 block cache 可以有效避免热数据的读请求落到 sstable 上，所以若内存余量充足，建议设置到128MB甚至256MB，读性能会有非常明显的提升。

### Index 和 Bloom Filter 系列参数
每个 SST 都可以有一个索引（Index）和 Bloom Filter（布隆过滤器），可以提升读性能，因为有了索引，不必顺序遍历整个 SST 文件，就可以定位具体的 Key 在哪里，因为已经保存了所有的 Key、Offset、Size 等元数据；而通过布隆过滤器，可以在假阳（False Positive）率很低的情况下，迅速判断某个 Key 是否在这个 SST 文件中，如果返回 False 就不再继续找索引了.

- filter_policy

    也就是 bloom filter，通常在点查 Get 的时候我们需要快速判断这个 key 在 SST 文件里面是否存在，如果 bloom filter 已经确定不存在了，就可以过滤掉这个 SST，减少没必要的磁盘读取操作了。我们使用 rocksdb::NewBloomFilterPolicy(bits_per_key) 来创建 bloom filter，bits_per_key 默认是 10，表示可能会有 1% 的误判率，bits_per_key 越大，误判率越小，但也会占用更多的 memory 和 space amplification
- cache_index_and_filter_blocks

    默认是 false，表示不在内存里缓存索引和过滤器 Block，而是用到了载入，不用就踢出去。如果设置为 true，则表示允许把这些索引和过滤器放到 Block Cache 中备用，这样可以提升局部数据存取的效率（无需磁盘访问就知道 Key 在不在，以及在哪里）。但是，如果启用了这个选项，必须同时把 pin_l0_filter_and_index_blocks_in_cache也设置为 true，否则可能会因为操作系统的换页操作，导致性能抖动.

    对于此参数，一定要注意 Block Cache 的总大小有限，如果允许 Index 和 Filter 也放进去，那么用来存放数据的空间就少了。因此，我们建议在 Key 具有局部热点（某些 Key 频繁访问，而其他的 Key 访问的很少）的情况下，才打开这两个参数；对于分布比较随机的 Key，这个参数甚至会起到反作用（随机 Key 时，读取性能大幅下降）
- optimize_filters_for_hits

    columnFamilyOptions.setOptimizeFiltersForHits设置为 true，则 RocksDB 不会给 L0 生成 Bloom Filter，据文档中描述，可以减少 90% 的 Filter 存储开销，有利于减少内存占用。但是，这个参数也仅仅适合于具有局部热点或者确信基本不会出现 Cache Miss 的场景，否则频繁的找不到，会拖累读取性能

打开 Partitioned Block 相关的配置， Partitioned Block 的原理是在 Index/Filter Block 加一层二级索引，当读 Index 或者 Filter 时，先将二级索引读入内存，然后根据二级索引找到所需的分区 Index Block，将其加载进 Block Cache 里面。

Partitioned Block 带来的优点如下：
1. 增加缓存命中率：大 Index /Filter Block 会污染缓存空间，将大的 Block 进行分区，允许以更细的粒度加载 Index/Filter Block，从而有效地利用缓存空间
1. 提高缓存效率：分区 Index/Filter Block 将变得更小，在 Cache 中锁的竞争将进一步降低，提升了高并发下的效率
1. 降低 IO 利用率：当索引 / 过滤器分区的缓存 Miss 时，只需要从磁盘加载一个小的分区，与读取整个 SST 文件的 Index/Filter Block 相比，这会使磁盘上的负载更小

Partitioned Block相关配置:
```conf
format_version = 5  
index_type = IndexType::kTwoLevelIndexSearch  # 使用 partition index 
NewBloomFilterPolicy(BITS, false) : 使用 Full Filter
BlockBasedTableOptions::partition_filters = true  # 使用partition filter, index_type 必须为 kTwoLevelIndexSearch
cache_index_and_filter_blocks = true
pin_top_level_index_and_filter = true
cache_index_and_filter_blocks_with_high_priority = true
pin_l0_filter_and_index_blocks_in_cache = true
optimize_filters_for_memory = true
```

### MemTable 系列参数
MemTable 是 RocksDB 重要的内存结构，也是 LSM 树的核心实现，通常使用 Skip List 跳表来实现。它可以看作是一个内存中的 Write Buffer（写缓冲），影响着写性能。

SkipList Memtable，相比 HashSkipList Memtable 跨多个前缀查找的性能更好，也更节省内存.

- whole_key_filtering

    针对 SkipList Memtable 打开whole_key_filtering选项，该选项会为 Memtable 中的 key 创建 Bloom Filter，这可以减少在跳表中的比较次数，从而降低点查询时的 CPU 使用率

    ```conf
    whole_key_filtering = true
    prefix_bloom_size_ratio=0.1
    ```

- write_buffer_size
    默认大小是 64 MB.  memtable 大小达到此阈值时，就会被标记为不可变. 

    通常来说，Write Buffer 越大，写放大效应越小，因而写性能也会改善, 但同时会增大 flush 后 L0、L1 层的压力。因此这个参数的调整，必须随着下面的几个参数一起来做，否则可能会达不到预期的效果。
- max_write_buffer_number

    columnFamilyOptions.setMaxWriteBufferNumber可控制内存中允许保留的 MemTable 最大个数, 包含活跃的和不可变的，超过这个个数后，就会被 Flush 刷写到磁盘上成为 SST 文件。

    这个参数的默认值是 2. 对于机械磁盘来说或者内存足够大，可以调大到 4 左右，以令 MemTable 的大小减小一些，降低 Flush 操作时造成 Write Stall 的概率.
- min_write_buffer_number_to_merge

    columnFamilyOptions.setMinWriteBufferNumberToMerge决定了 Write Buffer 合并的最小阈值，默认值为 1，对于机械硬盘来说可以适当调大，避免频繁的 Merge 操作造成的写停顿。

    根据我们的调优经验来看，这个参数调小、调大都会造成性能下滑，它的最佳值会在某个中间值附近，例如 2/3 等.
- max_bytes_for_level_base

    如果增加 Write Buffer Size，请一定要适当增加 L1 层的大小阈值（max_bytes_for_level_base），这个因子影响非常非常大。如果这个参数太小，那么每层能存放的 SST 文件就很少，层级就会变得很多，造成查找困难；如果这个参数太大，则会造成每层文件太多，那么执行 Compaction 等操作的耗时就会变得很长，此时容易出现 Write Stall（写停止）现象，造成写入中断。
- max_bytes_for_level_multiplier 

    决定了每层级的大小阈值的倍数关系. 如果不考虑其他因子的影响，如果 max_bytes_for_level_base = 1GB，max_bytes_for_level_multiplier = 5，那么 L1 的大小阈值是 1GB，L2 的大小阈值是 5GB，L3 的大小阈值是 25GB... 以此类推，所以 RocksDB 的默认分层存储叫做 Leveled 结构，有点类似于阶梯（如果启用 state.backend.rocksdb.compaction.level.use-dynamic-size 参数，则更加复杂）.

    这个 max_bytes_for_level_multiplier 参数对写入性能影响也是非常大的，请根据实际情况进行调整，没有一个统一的规则。
### Flush 和 Compaction 相关参数
ref:
- [RocksDB 7 终于解决了 Compaction 时性能下降问题](https://zhuanlan.zhihu.com/p/579468143)
- [带你全面了解compaction 的13个问题](https://tidb.net/blog/eedf77ff)
- [Rocksdb Compaction源码详解（二）：Compaction 完整实现过程 概览](https://blog.csdn.net/Z_Stand/article/details/107592966)

RocksDB 的后台进程中，有持续不断的 Flush 和 Compaction 操作。前者将 MemTable 的内容刷写到磁盘的 SST 文件中；后者则会对多个 SST 文件做归并和重整，删除重复值，并向更高的层级（Level）移动。例如 L0 -> L1 等。

频繁的 Flush 和 Compaction 操作，在写数据量大时，会严重影响性能，甚至造成写入的完全停顿，即 Write Stall.

RocksDB 的 compaction 策略，并且提到了读放大、写放大和空间放大的概念，对 RocksDB 的调优本质上就是在这三个因子之间取得平衡。而在 Flink 作业这种注重实时性的场合，则要重点考虑读放大和写放大。

- target_file_size

    ColumnFamilyOptions 的 setTargetFileSizeBase 方法可设置上一级的 SST 文件达到多大时触发 Compaction 操作，默认值是 2MB（每增加一级，阈值会自动乘以 target_file_size_multiplier）. 为了减少 Compaction 的频率，可以适当调大此参数，例如调整为 32MB 等，此参数对性能的影响也比较大.

    但是，调大这个参数以后，只是推迟了 Compaction 的时机，并没有真正减少数据量，因此可能会造成重复数据不能及时清理（影响读性能），或者一次需要清理超多的数据（影响写性能），因此这个参数比较缺乏灵活性.
- dynamic_level_bytes

    允许 RocksDB 对每层的存储的数据量阈值进行动态调整，不再是简单的 Level Base 的倍数关系，这样生成的 LSM 树更稳定。

    这个参数的默认值为 false，对于机械硬盘用户，建议设置为 true.
- compaction_style
    
    ColumnFamilyOptions 的 setCompactionStyle 方法允许用户调整 Compaction 的组织方式，默认值是 LEVEL（较为均衡），但也可以改为 UNIVERSAL 或 FIFO.

    相对于默认的 LEVEL 方式，UNIVERSAL 属于 Tiered 的一种，可以减少写放大效应，但是副作用是会增加空间放大和读放大效应，因此只适合写入量较大而读取不频繁，同时磁盘空间足够的场景。

    FIFO 则适合于将 RocksDB 作为时序数据库的场景，因为它是先入先出算法，可以批量淘汰掉过期的旧数据.
- compression_type

    ColumnFamilyOptions提供了 setCompressionType 方法, 可以指定对 Block 的压缩算法.

    如果追求性能，可以关闭压缩（NO_COMPRESSION)；否则建议使用 LZ4 算法，其次是 Zstd 算法. 启用压缩后，ReadOptions 的 verify_checksums 选项可以关闭，以提升读取速度（但是可能会受到磁盘坏块的影响）.

    对于大数据量、低 QPS 的场景，还会将最后一层设置为 ZSTD，进一步降低存储空间和减少成本，ZSTD 的优点是压缩比更高，压缩速度也较快.
- thread_num

    允许用户增加最大的后台 Compaction 和 Flush 操作的线程数，Flush 操作默认在高优先级队列，Compaction 操作则在低优先级队列。

    默认的后台线程数为 1，机械硬盘用户可以改为 4 等更大的值。

    如果后台所有的线程都在做 Compaction 操作时，如果这时候突然有很多写请求，就会引发写停顿（Write Stall）。写停顿可以通过日志或者监控指标来发现。
- write_batch_size 

    允许指定 RocksDB 批量写入时占用的最大内存量，默认为 2m，如果设置为 0 的话就会自动根据任务量进行调整。这个参数如果没有特别的需求，可以不调整
- max_background_compactions

    DBOptions 的setMaxBackgroundCompactions 可设置最大的后台并行 Compaction 作业数。

    如果 CPU 负载不高的话，建议增加这项的值，以允许更多的 Compaction 同时进行，减少读放大和空间放大，提升读取效率；但是如果设置的过大，可能会造成性能损耗，因为 Compaction 操作会带来停顿。

    L1-LN, 在非0Level上多个compactions可以被并行执行， max_background_compactions控制了最大并行数量.
- set_max_subCompactions

    DBOptions 允许用户通过 setMaxSubcompactions 选项，把 Compaction 操作拆分为并行的子任务。这个选项资料较少，我们测试来看，性能影响也不大，因此可以暂时忽略。

    大于1时，L0->L1, 会尝试把L0中数据文件分割开，用多线程合并到L1中
- set_level_0_file_num_compaction_trigger

    ColumnFamilyOption 的可以指定 L0 触发 Compaction 操作的文件个数阈值。默认值为 4，可以调大一些，以减少 Compaction 操作的频率（但是会带来 Compaction 时间的延长）。
- max_background_flushes

    DBOptions 的 setMaxBackgroundFlushes可设置后台最多同时进行的 Flush 任务数，默认值为 0. 如果设置为非 0 值，则表名把 Flush 任务放到高优先级队列。

    我们建议把这项的值设置为 1 即可，不需要太大。
- set_use_fsync

    DBOptions 提供了 setUseFsync 选项，让用户可以启用 fsync 同步写入磁盘，避免数据丢失，默认值是 false。

    我们建议 Flink 用户也设置为 false，因为 Flink 本身已经提供了 Checkpoint 机制来恢复状态，所以 WAL 和 fsync 等安全机制只会带来性能开销却不能带来好处。

    另外，机械硬盘用户，可以把另一个选项 setDisableDataSync 设置为 true，这样可以使用操作系统提供的 Buffer 来提升性能。
- set_num_levels

    ColumnFamilyOptions 的 setNumLevels 方法，可以设置总层级数。这个参数设置大一些没关系，最多就是高层次没有数据而已.

    默认值是 7，表示最多有 7 层。如果数据实在太多，可以设置为更大的值.

### 其他
- Key-Value 分离

    业界也基于WiscKey论文实现了 LSM 型存储引擎的 KV 分离，比如：RocksDB 的 BlobDB、PingCAP 的 Titan 引擎、百度的 UNDB 所使用的 Quantum 引擎.

    RocksDB 6.18 版本重新实现了 BlobDB（RocksDB 的 Key-Value 分离方案），将其集成到 RocksDB 的主逻辑中，并且也一直在完善和优化 BlobDB 相关特性.

    ```conf
    ColumnFamilyOptions.enable_blob_files = config_->RocksDB.enable_blob_files;
    min_blob_size = 4096
    blob_file_size = 128M
    blob_compression_type = lz4
    enable_blob_garbage_collection = true
    blob_garbage_collection_age_cutoff = 0.25
    blob_garbage_collection_force_threshold = 0.8
    ```

## FAQ
### 即使Put() use writeOptions.SetSync(true), Iterator遍历时部分数据(最近put的数据)无法访问到, 但Get()正常
```go
// github.com/tecbot/gorocksdb
writeOptions := gorocksdb.NewDefaultWriteOptions()
        writeOptions.SetSync(true)
```

解决方法:
```go
		// 手动flush
        fopts:=gorocksdb.NewDefaultFlushOptions()
        fopts.SetWait(true)
        db.Flush(fopts)
        defer fopts.Destroy()
```

推测: 部分数据在内存中, Iterator无法访问到???

### [Repairer 数据库文件](https://www.jianshu.com/p/8510f6c2562a)

### rocksdb内存组成
RocksDB的内存大致有如下四个区：
- Block Cache
- Indexes and bloom filters
- Memtables
- Blocked pinned by iterators

### 获取所使用的rocksdb version/config
ref:
- [MySQL · MyRocks · MyRocks参数介绍](https://developer.aliyun.com/article/397246)

查看db_path下的OPTIONS-<SN>中的section "Version"即可, 或源码的`include/rocksdb/version.h`

参考rocksdb 配置:
- [tikv](https://github.com/tikv/tikv/blob/master/docs/reference/configuration/rocksdb-option-config.md)

### core dumped
> driver : github.com/tecbot/gorocksdb

在x64使用6.10.1创建的db拷贝到arm64上用6.11.4打时, `gorocksdb.OpenDb()`崩溃了, 删除该db后重新运行程序, 此时由程序重新生成的db是正常的.

同时`db.Put(writOptions, []byte("123"), make([]byte, 1024))`即存储稍大点的value也会崩溃.

通过gdb, 问题定位在`util/crc32c_arm64.cc`的`t1 = (uint64_t)vmull_p64(crc1, k1)`上, 大value时该语句会报"Illegal instruction", 应该是cpu crc32硬件实现上有问题或有限制, [官方issue在这里](https://github.com/facebook/rocksdb/issues/7363).

解决方法: 将crc32c_arm64函数中的`#ifdef HAVE_ARM64_CRYPTO`块注释掉, 改用软件实现的crc32c算法.

`uint32_t crc32c_arm64(uint32_t crc, unsigned char const *data, unsigned len)`的上层封装入口就是`util/crc32c.cc`中的`Extend(uint32_t crc, const char* buf, size_t size)`:
```c++
// util/crc32c.h
// Return the crc32c of data[0,n-1]
inline uint32_t Value(const char* data, size_t n) {
  return Extend(0, data, n);
}

// util/crc32c_test.cc
TEST(CRC, Extend) {
  ASSERT_EQ(Value("hello world", 11),
            Extend(Value("hello ", 6), "world", 5));
}

// util/crc32c.cc
static Function ChosenExtend = Choose_Extend();
// crc就是将一段输入分多次计算时的前一段结果
uint32_t Extend(uint32_t crc, const char* buf, size_t size) {
  return ChosenExtend(crc, buf, size);
}
```

### github.com/tecbot/gorocksdb静态连接rocksdb
不推荐, 曾尝试过一个下午, 但还是不行.

### go driver
- [linxGnu/grocksdb](https://github.com/linxGnu/grocksdb), follow rocksdb latest, **推荐**

    `CGO_CFLAGS="-I/usr/local/include/rocksdb" CGO_LDFLAGS="-L/usr/local/lib -lrocksdb -lstdc++ -lm -lz -lsnappy -llz4 -lzstd" go build`


    崩溃示例:
    ```go
    func WrongCall(db *grocksdb.DB) {
        rOpt := grocksdb.NewDefaultReadOptions()
        wOpt := grocksdb.NewDefaultWriteOptions()

        k := []byte("k")
        tmp, _ := db.Get(rOpt, k)
        if tmp.Exists() {
            data := tmp.Data()
            tmp.Free()

            data[0] = 'a'
            if err := db.Put(wOpt, k, data); err != nil { // tmp已free, 再次使用data会导致cgo segmentation violation
                CheckErr(err)
            }
        }
    }
    ```
- [tecbot/gorocksdb](https://github.com/tecbot/gorocksdb), most using, **许久没更新, 不推荐**

    编译出的程序没法链接librocksdb.so时, 可参考[/go/cgo.md].

    **grocksdb的每个tag是与指定rocksdb version对应的(tag message上有提示)**.

### 同步写 与 异步写
参考:
- [RocksDB 笔记](http://alexstocks.github.io/html/rocksdb.html)

默认情况下，RocksDB 的写是异步的：仅仅把数据写进了操作系统的缓存区就返回了，而这些数据被写进磁盘是一个异步的过程. 如果为了数据安全，可以用如下代码把写过程改为同步写：
```c++
    rocksdb::WriteOptions write_options;
    write_options.sync = true;
    db->Put(write_options, …);
```

这个选项会启用 Posix 系统的 fsync(...) or fdatasync(...) or msync(..., MS_SYNC) 等函数.

**异步写的吞吐率是同步写的一千多倍**. 异步写的缺点是机器或者操作系统崩溃时可能丢掉最近一批写请求发出的由操作系统缓存的数据，但是 RocksDB 自身崩溃并不会导致数据丢失. 而机器或者操作系统崩溃的概率比较低，所以大部分情况下可以认为异步写是安全的.

RocksDB 由于有 WAL 机制保证，所以即使崩溃，其重启后会进行写重放保证数据一致性. 如果不在乎数据安全性，可以把 write_option.disableWAL 设置为 true，加快写吞吐率.

RocksDB 调用 Posix API fdatasync() 对数据进行异步写. 如果想用 fsync() 进行同步写，可以设置 Options::use_fsync 为 true.

### Rocksdb性能优化--写优化
参考:
- [官方RocksDB调优](https://github.com/facebook/rocksdb/wiki/RocksDB-Tuning-Guide), 翻译在[RocksDB参数调优](https://xiking.win/2018/12/05/rocksdb-tuning/)
- [Tuning RocksDB - Write Stalls](https://www.jianshu.com/p/a2892a161a7b)
- [RocksDB Tuning Guide](https://github.com/facebook/rocksdb/wiki/RocksDB-Tuning-Guide)

有人在HDD/NVME SSD/Ceph Rados上分别测试过Rocksdb的性能，这三者的同步写性能分别是(170MB/s, 1200MB/s, 120MB/s), 不过最终的测试结果却让人大跌眼镜， Rocksdb的性能分别为(<1MB/s, 50MB/s, 4MB/s)，最多只能发挥存储侧不到0.05%的性能.

自己测试性能(应用是自己写的cdp记录(by rocksdb)保存工具+同事的cdp截获驱动):
`/dev/sdi`大小:1G, 测试命令`dd if=/dev/zero of=/dev/sdi bs=4k count=2048`即写4M.

对照组: 不保存cdp信息耗时3s, 截获bio 2048个.

> options.SetWriteBufferSize = wbs, writeOptions.SetSync=sync, writeOptions.DisableWAL=wal
> sync=true也意味着wal=false, 设置wal=true会报错

实际组:
1. wbs=8M, sync=false, wal=false : 5m37s
1. wbs=64M, sync=false, wal=false : 6s
1. wbs=64M, sync=false, wal=true : 5s
1. wbs=128M, sync=false, wal=false : 5s
1. wbs=64M, sync=true, wal=false : 78s
1. wbs=128M, sync=true, wal=false : 92s ??? 缓存比上面大, 性能反而下降

> 这个测试数据量可能有问题, 没法将wbs写满, 时间有限, 先这样.

### Sync writes has to enable wal
writeOptions.SetSync=true时, writeOptions.DisableWAL必须为false.

### rocksdb delete a range of keys
参考:
- [Compaction](https://www.jianshu.com/p/a2c092b8d1ea)
- [带你全面了解 compaction 的13个问题](https://tidb.net/book/tidb-monthly/2022-06/usercase/compaction-question)
- [2020-10-03-Rocksdb删除问题总结.pdf](https://emperorlu.github.io/files/2020-10-03-Rocksdb%E5%88%A0%E9%99%A4%E9%97%AE%E9%A2%98%E6%80%BB%E7%BB%93.pdf)
- [rocksdb系列delete a range of keys](https://www.jianshu.com/p/cea1267628b6) from 官方文档的翻译[Delete A Range Of Keys](https://github.com/facebook/rocksdb/wiki/Delete-A-Range-Of-Keys)
- [iterator.seek() so slow](https://github.com/facebook/rocksdb/issues/10510)

    通过PerfContext输出的block_read_bytes和internal_delete_skipped_count, 发现有很多tombstone(即rocksdb的标志删除)导致`Seek()`慢

    `db.CompactRange(Range{nil, nil})`处理过程慢(原42g,后11g, 耗时11分钟), 处理后`Seek()`很快.
- [Rocksdb和Lethe 对Delete问题的优化](https://blog.csdn.net/Z_Stand/article/details/110270317)

最消极的做法: 遍历DB， 遇到特定范围里的key，直接调用Delete就可以了，这种方法适合于要删除的keys数量小的情况, 另一方面， 这种方法有两个显著问题：
1. 不能立刻收回资源，得等compaction完成后在能真正收回资源
2. 大量的tombstones（也就是标记为del的key） 会减缓迭代器效率

> [全量Compaction不可以停止，必须等待操作完成, 这是 RocksDB 的限制](https://docs.nebula-graph.com.cn/3.5.0/8.service-tuning/compaction/)

解决方法:
- 针对1:

    - 对要删除range key使用DeleteFilesInRange， 这个方法会直接删除只包含range keys的文件， 对于大块的range， 这个方法可以直接回收资源，值得注意的是：

        - 即使做完这个操作，一些在range keys范围内的数据依然存在于db中， 这个时候还需要一点的其他的opera
        - 这个操作直接忽视了snapshots， 导致可能通过snapshot 读不到本该可以读到的数据

    - compaction filter + compactRange, 一旦写个compaction filter, 做完compaction的时候，该删除的数据也没了. 在使用的时候， 可以设置CompactionFilter::IgnoreSnapshots = true， 这样的话， 就可以保证删除range keys，否则的话， 不能删除所有的keys， 这种方法也能很好的回收资源， 唯一的缺点是可能会增加I/O压力

- 针对2:

    - 如果从来不改写已经存在的keys，使用DB::SingleDelete而不是Delete就可以很快消灭tombstones， 这种Tombstone可以很快被清除而不是被一直推到last level
    - 使用NewCompactOnDeletionCollectorFactory()， 可以在有大量tombstones的时候加快compaction

        NewCompactOnDeletionCollectorFactory函数声明一个用户配置的表格属性收集器，这里的表格指的是sstable. 该函数需要传入的两个参数sliding_window_size和 deletion_trigger 表示删除收集器 的滑动窗口 和 触发删除的key的个数.

        在每个sst文件内维护一个滑动窗口，滑动窗口维护了一个窗口大小，在当前窗口大小内出现的delete-key的个数超过了窗口的大小，那么这个sst文件会被标记为need_compaction_，从而在下一次的compaction过程中被优先调度。达到将delete-key较多的sst文件快速下沉到底层，delete的数据被真正删除的目的，可以理解为这是一种更加激进的compaction调度策略

### RateLimiter
rocksdb通过RateLimiter来控制最大写入速度的上限, 因为在某些场景下，突发的大量写入会导致很大的read latency，从而影响系统性能.

RateLimiter参数：
- int64_t rate_bytes_per_sec：控制 compaction 和 flush 每秒总写入量的上限。一般情况下只需要调节这一个参数
- int64_t refill_period_us：控制 tokens 多久再次填满，譬如 rate_limit_bytes_per_sec 是 10MB/s，而 refill_period_us 是 100ms，那么每 100ms 的流量就是 1MB/s. 相当于将1s分成1s/refill_period_us个小窗口, 每个小窗口占用一定的token, 使得写入更均匀.
- int32_t fairness：用来控制 high 和 low priority 的请求，防止 low priority 的请求饿死.


### `setTotalOrderSeek(true)`
不设置时查找第一个key可能返回不一样, 原因未知.

### DeleteFilesInRange没生效
写几十个kv, 删除中间的一个range, 再执行DeleteFilesInRange, 打印数据发现没有生效, 推测是sst包含了未删除的key, 导致sst被忽略.

### [tikv如何调用rocksdb](https://asktug.com/t/topic/663579)
tikv通过[components/engine_rocks](https://github.com/tikv/tikv/blob/v6.1.0/components/engine_rocks/Cargo.toml#L56)调用rocksdb.

[rust-rocksdb](https://github.com/tikv/rust-rocksdb/blob/tikv-6.1/src/rocksdb.rs)

### [rocksdb perf_context性能perf](https://github.com/facebook/rocksdb/issues/261)
```go
grocksdb.SetPerfLevel(grocksdb.KEnableCount)
p:=gorocksdb.NewPerfContext()
p.Reset()

// do your thing -- SeekToFirst()

fmt.Println(p.Report(false))
p.Destroy()
```

```c++
rocksdb::SetPerfLevel(rocksdb::kEnableCount);
rocksdb::perf_context.Reset();

// do your thing -- SeekToFirst()

cout << rocksdb::perf_context.ToString();
```

### ttl
参考:
- [基于RocksDB实现精准的TTL过期淘汰机制](https://segmentfault.com/a/1190000021185954)
- [TTL](http://pegasus.incubator.apache.org/api/ttl)
- [MySQL · MyRocks · TTL特性介绍](http://mysql.taobao.org/monthly/2018/04/04/)

#### OpenDbWithTTL
这个是 RocksDB 本身支持的一种数据过期淘汰方案, 该方案是通过特定的 API 打开 DB(作用在column_family)，对写入该 DB 的全部 key 都遵循一个 TTL 过期策略，例如 TTL 为 3 天，那么写入该 DB 的 key 都会在写入的三天后自动过期. 该方案底层也是通过 compaction filter 实现的，也就是说过期数据虽然对用户不可见，但是磁盘空间并不会及时回收，另外该方案不灵活，无法针对每一条 key 设置 TTL.

每次put数据时，会调用DBWithTTLImpl::AppendTS将过期时间append到value最后.

在Compact时通过自定义的TtlCompactionFilter, 去判断数据是否可以清理, 具体参考DBWithTTLImpl::IsStale.

> RocksDB TTL在compact时才清理过期数据，所以，过期时间并不是严格的，会有一定的滞后，取决于compact的速度.

## FAQ
### grocksdb程序报`pure virtual method called\nterminate called without an active exception`
这个错误的主要原因是当前对象已经被销毁或者正在被销毁, 导致最终调用到其基类的虚方法上, 最终报出了这个错误.

> 其他场景(比如嵌入式)报该错可能是程序链接的so不存在或环境不提供某些方法.

解决方法: 记录db引用count, 使用前atomic+1,defer atomic-1, 等到再count为0后再关闭

### rocksdb 7.10.2报`Disk quota exceeded`后zfs fs扩容成功, 但rocksdb还是报该错
ref:
- [Can't auto recovery when encounter ENOSPC error from the filesystem](https://github.com/facebook/rocksdb/issues/10134)

官方bug.

### rocksdb刚启动就coredump
env:
- Ubuntu 20.04
- jemalloc: 5.2.1
- rocksdb: 7.10.2/8.1.1

通过`gdb <my_program_path> <coredump_path>`, 崩溃在`_start_thread` jemalloc分配内存时.

解决方法:
1. ~~使用其他版本jemalloc, 需要自编译: 官方apt repo的jemalloc需要libc6>=2.34, 而Ubuntu 20.04的是2.31~~

    后来发现, 构建rocksdb时优先使用了`/usr/local`下的自编译jemalloc, 而自编的jemalloc 5.3.0/5.2.1都会引发coredump, 移除自编译jemalloc而采用ubuntu官方jemalloc后不再coredump
1. 使用tcmalloc
1. 采用默认malloc, 经测试可行

### tikv调用rocksdb
ref:
- [求助！关于tikv调用rust-rocksdb](https://asktug.com/t/topic/663579)

`tikv components/engine_rocks` -> `https://github.com/tikv/rust-rocksdb` -> rocksdb


> AgateDB - A new storage engine created by PingCAP in an attempt to replace RocksDB from the Tikiv DB stack

### Flush报`Shutdown in progress`
Flush前报`CancelAllBackgroundWork(true)`导致, 将CancelAllBackgroundWork放到db close前即可

### options.Clone().Destroy()崩溃
env:
- oracle linux 7.9 x64
- rocksdb 8.1.1
- github.com/linxGnu/grocksdb v1.8.0
- jemalloc: from 官方repo

ps: ubuntu 20.04/22.04正常

```go
func main() {
    dir, err := os.MkdirTemp(".", "t-")
    CheckErr(err)

    givenNames := []string{"default", "write"}
    options := grocksdb.NewDefaultOptions()
    options.SetCreateIfMissing(true)
    options.SetCreateIfMissingColumnFamilies(true)

    oOptions := options.Clone()
    {
        oOptions.SetMemTablePrefixBloomSizeRatio(0.1)
        oOptions.SetPrefixExtractor(grocksdb.NewFixedPrefixTransform(1))

        bto := grocksdb.NewDefaultBlockBasedTableOptions()
        bto.SetBlockSize(32 << 20)
        bto.SetChecksum(0x1)
        bto.SetFilterPolicy(grocksdb.NewBloomFilterFull(1))
        bto.SetCacheIndexAndFilterBlocks(true)
        bto.SetCacheIndexAndFilterBlocksWithHighPriority(true)
        oOptions.SetBlockBasedTableFactory(bto)
    }

    dOptions := options.Clone()
    {
        dOptions.SetOptimizeFiltersForHits(false)
        dOptions.SetPrefixExtractor(grocksdb.NewNoopPrefixTransform())

        bto := grocksdb.NewDefaultBlockBasedTableOptions()
        bto.SetBlockSize(32 << 20)
        bto.SetChecksum(0x1)
        bto.SetWholeKeyFiltering(false)
        bto.SetCacheIndexAndFilterBlocks(true)
        bto.SetCacheIndexAndFilterBlocksWithHighPriority(true)
        dOptions.SetBlockBasedTableFactory(bto)
    }

    db, cfh, err := grocksdb.OpenDbColumnFamilies(options, dir, givenNames, []*grocksdb.Options{dOptions, oOptions})
    CheckErr(err)

    if len(cfh) != 2 {
        panic("cfh")
    }

    cfh[0].Destroy()
    cfh[1].Destroy()

    db.Close()
    dOptions.Destroy() // will segmentation violation, 崩溃原因是[`dOptions.SetPrefixExtractor(grocksdb.NewNoopPrefixTransform())`](https://github.com/facebook/rocksdb/issues/1095), 当设置SetPrefixExtractor后, 不调用`C.rocksdb_slicetransform_destroy(opts.cst)`即不调用`dOptions.Destroy(),oOptions.Destroy()`, 建议是仅注释C.rocksdb_slicetransform_destroy(), 避免相关资源内存泄露
    oOptions.Destroy()
    options.Destroy()
}
```

解决方法: 注释dOptions.Destroy()和oOptions.Destroy()

### `make dbg`报`error: format '%u' expects argument of type 'unsigned int', but argument 2 has type 'std::reference_wrapper<unsigned int>'`
ref:
- [build v7.10.2 failed with gcc 11/12](https://github.com/facebook/rocksdb/issues/11655)

env: gcc 7/9/11/12

解决方法: 在Makefile的`WARNING_FLAGS = ...`后追加`-Wformat=0`

### `make dbg`报`error: no match for 'operator=' (operand type are 'std::reference_wrapper<unsigned int>' and 'int')`
ref:
- [build v7.10.2 failed with gcc 11/12](https://github.com/facebook/rocksdb/issues/11655)

env: gcc 7/9/11/12

解决方法是修改代码:

```c++
class StressCacheKey {
 public:
  void Run() {
    if (FLAGS_sck_footer_unique_id) {
      // Proposed footer unique IDs are DB-independent and session-independent
      // (but process-dependent) which is most easily simulated here by
      // assuming 1 DB and (later below) no session resets without process
      // reset.
      FLAGS_sck_db_count = 1;
    }
...
```

改为:
```c++
class StressCacheKey {
 public:
  void Run() {
    if (FLAGS_sck_footer_unique_id) {
      // Proposed footer unique IDs are DB-independent and session-independent
      // (but process-dependent) which is most easily simulated here by
      // assuming 1 DB and (later below) no session resets without process
      // reset.
      unsigned int x = 1;
      FLAGS_sck_db_count = std::ref(x);
    }
```

### `make all`报`error: '__Y' may by use uninitialized [-Werror=maybe-uninitialized]`和`error: '__Y' may by use uninitialized [-Werror=uninitialized]`
ref:
- [build v7.10.2 failed with gcc 11/12](https://github.com/facebook/rocksdb/issues/11655)
- [Building from source with MKL-DNN fails on Fedora37](https://github.com/apache/mxnet/issues/21188)

env: gcc 12

是在编译`util/xxhash.o`时遇到的.

解决方法: 在Makefile的`WARNING_FLAGS = ...`后追加`-Wno-maybe-uninitialized -Wno-uninitialized`

> gcc 9编译`util/xxhash.o`没报该错

### rocksdb写入卡住
- [Write Stalls](https://github.com/facebook/rocksdb/wiki/Write-Stalls)

在rocksdb log里检索`ing writes`, 导致该问题有两种情况:
- `Stalling writes ...`: 触发写入限速
- `Stopping writes ...`: 触发写暂停

根据具体原因调整参数

### 获取db状态
[GetProperty](https://bravoboy.github.io/2019/07/06/rocksdb-log/)
