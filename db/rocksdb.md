# rocksdb
参考:
- [RocksDB系列](https://www.jianshu.com/p/061927761027)

RocksDB的目的是成为一套能在服务器压力下，真正发挥高速存储硬件（特别是Flash 和 RAM）性能的高效单点数据库系统. 它是一个C++库，允许存储任意长度二进制kv数据, 支持原子读写操作, 因此本质上来说它是一个可插拔式的存储引擎选择.

RocksDB大量复用了levedb的代码，并且还借鉴了许多HBase的设计理念, 同时Rocksdb也借用了一些Facebook之前就有的理念和代码.

RocksDB是一个嵌入式的K-V（任意字节流）存储. 所有的数据在引擎中是有序存储，可以支持Get(key)、Put（Key）、Delete（Key）和NewIterator()。RocksDB的基本组成是memtable、sstfile和logfile。memtable是一种**内存数据结构**，写请求会先将数据写到memtable中，然后可选地写入事务日志logfile(WAL)。logfile是一个**顺序写**的文件。当memtable被填满的时候，数据会flush到sstfile中，然后这个memtable对应的logfile也会安全地被删除。sstfile中的数据也是有序存储以方便查找.

> rocksdb支持Direct IO, 以绕过系统Page Cache，通过应用内存从存储设置中直接进行IO读写操作.

![](/misc/img/develop/4304640-891400b1777c999d.png)

## 编译
1. 参照[rocksdb INSTALL](https://github.com/facebook/rocksdb/blob/master/INSTALL.md), 选择平台安装依赖lib
1. `cd rocksdb_source_root`, 选择`make static_lib/make shared_lib`进行编译

    如果构建环境存在jemalloc/tcmalloc, make会通过`build_tools/build_detect_platform <platform>`将相应的环境变量存入生成的make_config.mk中, 供自身使用
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
可插拔 memtable，RocksDB 的 memtable 的默认实现是一个 skiplist。skiplist 是一个有序集，当工作负载使用 range-scans 并且交织写入时，这是一个必要的结构。然而，一些应用程序不交织写入和扫描，而一些应用程序根本不执行范围扫描。对于这些应用程序，排序集可能无法提供最佳性能。因此，RocksDB 支持可插拔的 API，允许应用程序提供自己的 memtable 实现。

开发库提供了三个 memtable：skiplist memtable，vector memtable 和前缀散列（prefix-hash） memtable
- Vector memtable 适用于将数据批量加载到数据库中。每个写入在向量的末尾插入一个新元素; 当它是刷新 memtable 到存储的时候，向量中的元素被排序并写出到 L0 中的文件
- 前缀散列 memtable 允许对 gets，puts 和 scans-within-a-key-prefix 进行有效的处理。

## SSTFile(SSTTable)
RocksDB在磁盘上的file结构sstfile由block作为基本单位组成，一个sstfile结构由多个data block和meta block组成， 其中data block就是数据实体block，meta block为元数据block， 其中data block就是数据实体block，meta block为元数据block。 sstfile组成的block有可能被压缩(compression)，不同level也可能使用不同的compression方式。 sstfile如果要遍历block，会逆序遍历，从footer开始。

sst里面的数据按照key进行排序能方便对其进行二分查找. 在SST文件内，还额外包含以下特殊信息：
- Bloom Fileter : 用于快速判断目标查询key是否存在于当前SST文件内
- Index / Partition Index，SST内部数据块索引文件快速找到数据块的位置

compaction输入的SST file并不是立即就从SST file集合中删除，因为有可能在这些SST file上正进行着get or iterator操作. 只有当冗余的SST file上没有任何操作的时候，才会执行真正的删除文件操作. [这些逻辑是通过引用计数来实现的](https://www.jianshu.com/p/b95db752178f).

写流程：
rocksdb写入时，直接以append方式写到log文件以及memtable，随即返回，因此非常快速. memtable/immute memtable触发阈值后， flush 到Level0 SST，Level0 SST触发阈值后，经合并操作(compaction)生成level 1 SST， level1 SST 合并操作生成level 2 SST，以此类推，生成level n SST.

![rocksdb 写入原理](/misc/img/rocksdb/1f950c8b30fe6992437242c368f0f8b1.png)

Get()流程：
1. 在MemTable中查找，无法命中转到下一流程
1. 在immutable_memtable中查找，查找不中转到下一流程
1. 在第0层SSTable中查找，无法命中转到下一流程

	对于L0 的文件，RocksDB 采用遍历的方法查找，所以为了查找效率 RocksDB 会控制 L0 的文件个数

1. 在剩余SSTable中查找

    对于 L1 层以及 L1 层以上层级的文件，每个 SSTable 没有交叠，可以使用二分查找快速找到 key 所在的 Level 以及 SSTfile.

> 如果启用 Level Style Compaction, L0 存储着 RocksDB 最新的数据，Lmax 存储着比较老的数据，L0 里可能存着重复 keys，但是其他层文件则不可能存在重复 key.

## RocksDB的典型场景（低延时访问）:
1. 需要存储用户的查阅历史记录和网站用户的应用
1. 支持大量写和删除操作的消息队列

## 功能
### Column Families
RocksDB支持将一个数据库实例分片为多个列族(column families, 类似表Table). 类似HBase，每个DB新建时默认带一个名为"default"的列族，如果一个操作没有携带列族信息，则默认使用这个列族. 如果WAL开启，当实例crash再恢复时，RocksDB可以保证用户一个一致性的视图. 通过WriteBatch API，可以实现跨列族操作的原子性.

所有 Column Family 共享 WAL、Current、Manifest 文件，但是每个 Column Family 有自己单独的 memtable & ssttable(sstfile)，即 log 共享而数据分离.

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

1. workload模拟工具

    用户数据访问行为模拟工具
    Workload生成工具

    [ldb](https://github.com/facebook/rocksdb/wiki/Administration-and-Data-Access-Tool)命令行工具提供了不同的数据访问和数据库管理命令.
    [sst_dump tool](https://github.com/facebook/rocksdb/wiki/Administration-and-Data-Access-Tool)可以dump数据然后分析SST file.

    ```bash
    # Linux
	$ cd rocksdb
	$ make ldb sst_dump
	$ cp ldb /usr/local/bin/
	$ cp sst_dump /usr/local/bin/
    ```

1. 性能分析工具，DB Analyzer

## 源码
- [RocksDB · 数据的读取(一)](http://mysql.taobao.org/monthly/2018/11/05/)
- [RocksDB解析](https://www.cnblogs.com/pdev/p/11277784.html)

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
查看db_path下的OPTIONS-<SN>中的section "Version"即可

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
- [tecbot/gorocksdb](https://github.com/tecbot/gorocksdb), most using
- [linxGnu/grocksdb](https://github.com/linxGnu/grocksdb), follow rocksdb latest

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
- [rocksdb系列delete a range of keys](https://www.jianshu.com/p/cea1267628b6)

最消极的做法: 遍历DB， 遇到特定范围里的key，直接调用Delete就可以了，这种方法适合于要删除的keys数量小的情况, 另一方面， 这种方法有两个显著问题：
1. 不能立刻收回资源，得等compaction完成后在能真正收回资源
2. 大量的tombstones（也就是标记为del的key） 会减缓迭代器效率

解决方法:
- 针对1:

    - 对要删除range key使用DeleteFilesInRange， 这个方法会直接删除只包含range keys的文件， 对于大块的range， 这个方法可以直接回收资源，值得注意的是：

        - 即使做完这个操作，一些在range keys范围内的数据依然存在于db中， 这个时候还需要一点的其他的opera
        - 这个操作直接忽视了snapshots， 导致可能通过snapshot 读不到本该可以读到的数据

    - compaction filter + compactRange, 一旦写个compaction filter, 做完compaction的时候，该删除的数据也没了. 在使用的时候， 可以设置CompactionFilter::IgnoreSnapshots = true， 这样的话， 就可以保证删除range keys，否则的话， 不能删除所有的keys， 这种方法也能很好的回收资源， 唯一的缺点是可能会增加I/O压力

- 针对2:

    - 如果从来不改写已经存在的keys，使用DB::SingleDelete而不是Delete就可以很快消灭tombstones， 这种Tombstone可以很快被清除而不是被一直推到last level
    - 使用NewCompactOnDeletionCollectorFactory()， 可以在有大量tombstones的时候加快compaction
### RateLimiter
rocksdb通过RateLimiter来控制最大写入速度的上限, 因为在某些场景下，突发的大量写入会导致很大的read latency，从而影响系统性能.

RateLimiter参数：
- int64_t rate_bytes_per_sec：控制 compaction 和 flush 每秒总写入量的上限。一般情况下只需要调节这一个参数
- int64_t refill_period_us：控制 tokens 多久再次填满，譬如 rate_limit_bytes_per_sec 是 10MB/s，而 refill_period_us 是 100ms，那么每 100ms 的流量就是 1MB/s. 相当于将1s分成1s/refill_period_us个小窗口, 每个小窗口占用一定的token, 使得写入更均匀.
- int32_t fairness：用来控制 high 和 low priority 的请求，防止 low priority 的请求饿死.
