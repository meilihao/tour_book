# mysql
参考:
- [MySQL5.7的官方文档(中文版)](https://www.docs4dev.com/docs/zh/mysql/5.7/reference/manual-info.html)
- [MariaDB general release maintenance periods](https://mariadb.org/about/#maintenance-policy)

MySQL 的逻辑架构图:
![](/misc/img/mysql/dfff6efbab0d51a715a36f867daeacf8.png)
mysql的查询流程:
![](/misc/img/mysql/20c7d975ab3bd4617c6b4f46b720fea3.png)
[MySQL架构总览](https://wwxiong.com/2017/06/sql-optimization-up/):
![](/misc/img/mysql/mysql.png)

MySQL可以分为Server层和存储引擎层两部分:
1. Server层包括连接器、~~查询缓存(MySQL 8.0 版本后移除:随着并发性的增加, 查询缓存成为让人诉病的瓶颈)~~、分析器、优化器、执行器等，涵盖MySQL的大多数核心服务功能，以及所有的内置函数（如日期、时间、数学和加密函数等），所有跨存储引擎的功能都在这一层实现，比如存储过程、触发器、视图等.

  - 分析器: 解析sql以创建内部数据结构(解析树)
  - 优化器: 对解析树进行各种优化,包括重写查询、决定表的读取顺序,以及选择合适的索引等. 用户可以通过特殊关键 字向 优化器传递提示, 从而影响优化器的决策过程
1. 存储引擎层负责数据的存储和提取. 其架构模式是插件式的，支持InnoDB(**推荐**)、~~MyISAM(不推荐)~~、Memory等多个存储引擎. 现在最常用的存储引擎是InnoDB，它从MySQL 5.5.5版本开始成为了默认存储引擎

查询语句的执行流程如下：`权限校验（如果命中缓存）--->查询缓存--->分析器--->优化器--->权限校验--->执行器--->引擎`
更新语句执行流程如下：`分析器---->权限校验---->执行器--->引擎---redo log(prepare 状态)--->binlog--->redo log(commit 状态)`

## 安装
- [mariadb 离线包下载](https://mariadb.com/downloads/)

## 数据类型
[](https://oss.javaguide.cn/github/javaguide/mysql/summary-of-mysql-field-types.png)

DECIMAL 和 FLOAT 的区别是：DECIMAL 是定点数，FLOAT/DOUBLE 是浮点数。DECIMAL 可以存储精确的小数值，FLOAT/DOUBLE 只能存储近似的小数值.

如果预期长度范围可以通过 VARCHAR 来满足，建议避免使用 TEXT. 数据库规范通常不推荐使用 BLOB 和 TEXT 类型，这两种类型具有一些缺点和限制, 比如没有默认值, 检索效率较低. 必须使用时建议把 BLOB 或是 TEXT 列分离到单独的扩展表中.

DATETIME 类型没有时区信息，TIMESTAMP 和时区有关.

MySQL 中没有专门的布尔类型，而是用 TINYINT(1) 类型来表示布尔值.

除非有特别的原因使用 NULL 值，否则应该总是让字段保持 NOT NULL:
- 索引 NULL 列需要额外的空间来保存，所以要占用更多的空间
- 进行比较和计算时要对 NULL 值做特别的处理

Blob和text有什么区别:
1. Blob用于存储二进制数据，而Text用于存储大字符串。
1. Blob值被视为二进制字符串（字节字符串）,它们没有字符集，并且排序和比较基于列值中的字节的数值。
1. text值被视为非二进制字符串（字符字符串）。它们有一个字符集，并根据字符集的排序规则对值进行排序和比较

## 查询

> 由执行器检查是否有查询权限.
>
> 慢查询日志中的rows_examined字段表示这个语句执行过程中执行器扫描了多少行. 在有些场景下，执行器调用一次，在引擎内部则扫描了多行，因此引擎扫描次数跟rows_examined并不是完全相同的.

### 查询缓存
大多数情况下建议**不要使用查询缓存**: **查询缓存的失效非常频繁，只要有对一个表的更新，这个表上所有的查询缓存都会被清空; 同时MySQL 8.0版本已砍掉查询缓存功能**.

## 日志
MySQL 日志 主要包括错误日志、查询日志、慢查询日志、事务日志、二进制日志几大类, 其中比较重要的还要属二进制日志 binlog（归档日志）和事务日志 redo log（重做日志）和 undo log（回滚日志）

MySQL InnoDB 引擎使用 redo log(重做日志) 保证事务的持久性，使用 undo log(回滚日志) 来保证事务的原子性.

## 更新
update语句执行流程(浅色框表示是在InnoDB内部执行的，深色框表示是在执行器中执行):
![](/misc/img/mysql/2e5bff4910ec189fe1ee6e2ecc7b4bbe.png)

更新涉及两个重要的日志模块：redo log（重做日志）和 binlog（归档日志）.将redo log的写入拆成了两个步骤：prepare和commit，这就是"两阶段提交", 这是让这两个log状态保持逻辑上的一致.

数据写入流程: 开始事务-> undo log->redo log-> 修改buffer pool -> [bin log] -> 提交事务

### redo log
WAL的全称是Write-Ahead Logging，它的关键点就是先写日志，再写磁盘. 它工作在**存储引擎层**, 是物理日志, 记录内容是`在某个数据页上做了什么修改`. 它让 MySQL 拥有了**崩溃恢复**能力.

具体来说，当有一条记录需要更新的时候，InnoDB引擎就会先把记录写到redo log 里面，并更新内存，这样更新就算完成了. 同时，InnoDB引擎会在适当(空闲)的时候，将这个操作记录更新到磁盘里面, 并清理redo log.

InnoDB 将 redo log 刷到磁盘上有几种情况：
1. 事务提交：当事务提交时，log buffer 里的 redo log 会被刷新到磁盘（可以通过innodb_flush_log_at_trx_commit参数控制）
1. log buffer 空间不足时：log buffer 中缓存的 redo log 已经占满了 log buffer 总容量的大约一半左右，就需要把这些日志刷新到磁盘上
1. 事务日志缓冲区满：InnoDB 使用一个事务日志缓冲区（transaction log buffer）来暂时存储事务的重做日志条目。当缓冲区满时，会触发日志的刷新，将日志写入磁盘
1. Checkpoint（检查点）：InnoDB 定期会执行检查点操作，将内存中的脏数据（已修改但尚未写入磁盘的数据）刷新到磁盘，并且会将相应的重做日志一同刷新，以确保数据的一致性
1. 后台刷新线程：InnoDB 启动了一个后台线程，负责周期性（每隔 1 秒）地将脏页（已修改但尚未写入磁盘的数据页）刷新到磁盘，并将相关的重做日志一同刷新
1. 正常关闭服务器：MySQL 关闭的时候，redo log 都会刷入到磁盘里去

总之，InnoDB 在多种情况下会刷新重做日志，以保证数据的持久性和一致性.

硬盘上存储的 redo log 日志文件不只一个，而是以一个日志文件组的形式出现的，每个的redo日志文件大小都是一样的. 它采用的是环形数组形式，从头开始写，写到末尾又回到头循环写.

InnoDB的redo log是固定大小的, 从头开始写，写到末尾就又回到开头的循环写, 对此它使用了两个游标:
- write pos是当前记录的位置，一边写一边后移并再循环
- checkpoint是当前要清理的位置，也是往后推移并再循环的，但清理记录前会确保记录已更新到数据文件.

`write pos -> checkpoint`之间即为空闲可用部分. 如果write pos追上checkpoint，表示WAL已满，这时候不能再执行新的更新，得停下来先清理一些记录，把checkpoint推进一下.

> 在 MySQL 8.0.30 之前可以通过 innodb_log_files_in_group 和 innodb_log_file_size 配置日志文件组的文件数和文件大小，但在 MySQL 8.0.30 及之后的版本中，这两个变量已被废弃，即使被指定也是用来计算 innodb_redo_log_capacity 的值。而日志文件组的文件数则固定为 32，文件大小则为 innodb_redo_log_capacity / 32

有了redo log，InnoDB就可以保证即使数据库发生异常重启，之前提交的记录都不会丢失，这个能力称为crash-safe.

参数`innodb_flush_log_at_trx_commit`可控制N次事务后redo log持久化到磁盘:
- 0 : 每次事务提交时不进行刷盘操作. 这种方式性能最高，但是也最不安全. 如果 MySQL 挂了或宕机了，可能会丢失最近 **1 秒**内的事务
- 1 : 每次事务提交时都将进行刷盘操作, **默认**. 这种方式性能最低，但是也最安全，因为只要事务提交成功，redo log 记录就一定在磁盘里，不会有任何数据丢失 
- 2 : 每次事务**提交**时都只把 log buffer 里的 redo log 内容写入 page cache（文件系统缓存）. 这种方式的性能和安全性都介于前两者中间

  如果仅仅只是 MySQL 挂了不会有任何数据丢失，但是宕机可能会有1秒数据的丢失

InnoDB 存储引擎有一个后台线程，每隔1 秒，就会把 redo log buffer 中的内容写到文件系统缓存（page cache），然后调用 fsync 刷盘. 因此一个没有提交事务的 redo log 记录，也可能会刷盘. 这也是mysql故障可能会仅丢失最近 **1 秒**内的数据的原因.

除了后台线程每秒1次的轮询操作，还有一种情况，当 redo log buffer 占用的空间即将达到 innodb_log_buffer_size 一半的时候，后台线程会主动刷盘.

### binlog
binlog（归档日志）保证了 MySQL 集群架构的数据一致性.

它工作在Server层, 是逻辑日志, 记录内容是语句的原始逻辑, 是记录所有涉及更新数据的逻辑操作，并且是顺序写.

不管用什么存储引擎，只要发生了表数据更新，都会产生 binlog 日志.

MySQL 数据库的数据备份、主备、主主、主从都离不开 binlog，需要依靠 binlog 来同步数据，保证数据一致性.

参数`sync_binlog`表示N次事务后binlog持久化到磁盘, 通常是`1`以保证MySQL异常重启之后binlog不丢失.

binlog 日志有三种格式，可以通过binlog_format参数指定:
- statement : 记录的内容是SQL语句原文

  `update T set update_time=now() where id=1`会导致与原库的数据不一致
- row : 记录的内容不再是简单的SQL语句了，还包含操作的具体数据

  通常情况下都是指定为row，这样可以为数据库的恢复与同步带来更好的可靠性. 但是这种格式，需要更大的容量来记录，比较占用空间，恢复与同步时会更消耗 IO 资源，影响执行速度

  row格式记录的内容看不到详细信息，要通过mysqlbinlog工具解析出来
- mixed

  MySQL 会判断这条SQL语句是否可能引起数据不一致，如果是，就用row格式，否则就用statement格式

binlog 的写入时机: 事务执行过程中，先把日志写到binlog cache，事务提交的时候，再把binlog cache写到 binlog 文件中.

一个事务的 binlog 不能被拆开，无论这个事务多大，也要确保一次性写入，所以系统会给每个线程分配一个块内存作为binlog cache. 通过binlog_cache_size参数控制单个线程 binlog cache 大小，如果存储内容超过了这个参数就要暂存到Swap.

主从复制的原理:
1. 生成binlog
1. 分发binlog
1. 回放binlog

异步、同步和半同步复制:
1. 异步复制（Asynchronous replication, 默认），主库在执行完客户端提交的事务后会立即将结果返给给客户端，并不关心从库是否已经接收并处理。原理最简单，性能最好，但是主从之间数据不一致的概率很大
1. 全同步复制（Fully synchronous replication），指当主库执行完一个事务，所有的从库都执行了该事务才返回给客户端。因为需要等待所有从库执行完该事务才能返回，所以全同步复制的性能必然会收到严重的影响
1. 半同步复制（Semisynchronous replication），介于异步复制和全同步复制之间，主库在执行完客户端提交的事务后不是立刻返回给客户端，而是等待至少一个从库接收到并写到relay log中才返回给客户端。相对于异步复制，半同步复制牺牲了一定的性能，提高了数据的安全性

### redo log 与 binlog
不同点:
- redo log是InnoDB引擎特有的；binlog是MySQL的Server层实现的，所有引擎都可以使用
- redo log是物理日志，记录的是`在某个数据页上做了什么修改`；binlog是逻辑日志，记录的是这个语句的原始逻辑，比如`给ID=2这一行的c字段加1`.
- redo log是循环写的，空间固定；binlog是可以追加写入的. `追加写`是指binlog文件写到一定大小后会切换到下一个，并不会覆盖以前的日志.
- redo log 与 binlog 的写入时机: 在执行更新语句过程，会记录 redo log 与 binlog 两块日志，以基本的事务为单位，redo log 在事务执行过程中可以不断写入，而 binlog 只有在提交事务时才写入

如果执行过程中写完 redo log 日志后，binlog 日志写期间发生了异常, 为了解决两份日志之间的逻辑一致问题，InnoDB 存储引擎使用两阶段提交方案. 即将 redo log 的写入拆成了两个步骤prepare和commit.

使用两阶段提交后，写入 binlog 时发生异常也不会有影响，因为 MySQL 根据 redo log 日志恢复数据时，发现 redo log 还处于prepare阶段，并且没有对应 binlog 日志，就会回滚该事务. 如果redo log 设置commit阶段发生异常, 并不会回滚事务, 因为redo log 是处于prepare阶段，但是能通过事务id找到对应的 binlog 日志，所以 MySQL 认为是完整的，就会提交事务恢复数据.

### undo log
每一个事务对数据的修改都会被记录到 undo log ，当执行事务过程中出现错误或者需要执行回滚操作的话，MySQL 可以利用 undo log 将数据恢复到事务开始之前的状态。undo log 属于逻辑日志，记录的是 SQL 语句

每一个事务对数据的修改都会被记录到 undo log ，当执行事务过程中出现错误或者需要执行回滚操作的话，MySQL 可以利用 undo log 将数据恢复到事务开始之前的状态.

undo log 属于逻辑日志，记录的是 SQL 语句，比如说事务执行一条 DELETE 语句，那 undo log 就会记录一条相对应的 INSERT 语句。同时，**undo log 的信息也会被记录到 redo log 中，因为 undo log 也要实现持久性保护**。并且，undo-log 本身是会被删除清理的，例如 INSERT 操作，在事务提交之后就可以清除掉了；UPDATE/DELETE 操作在事务提交不会立即删除，会加入 history list，由后台线程 purge 进行清理.

undo log 是采用 segment（段）的方式来记录的，每个 undo 操作在记录的时候占用一个 undo log segment（undo 日志段），undo log segment 包含在 rollback segment（回滚段）中。事务开始时，需要为其分配一个 rollback segment。每个 rollback segment 有 1024 个 undo log segment，这有助于管理多个并发事务的回滚需求.

通常情况下， rollback segment header（通常在回滚段的第一个页）负责管理 rollback segment。rollback segment header 是 rollback segment 的一部分，通常在回滚段的第一个页。history list 是 rollback segment header 的一部分，它的主要作用是记录所有已经提交但还没有被清理（purge）的事务的 undo log。这个列表使得 purge 线程能够找到并清理那些不再需要的 undo log 记录.

## 事务
MySQL的事务启动方式：
1. 显式启动事务语句，`begin/start transaction`+`commit`/`rollback`
  具体操作: 保存 autocommit 的当前状态，然后 start transaction，直到 commit or rollback 结束本次事务，再恢复之前保存的 autocommit 的状态
1. `set autocommit=0`会关闭这个线程的自动提交, 这相当于开启一个全局的事务且这个事务持续存在直到主动执行commit 或 rollback，或者断开连接. 在 mysql 的事务中，默认 `autocommit = 1`，每一次 sql 操作都被认为是一个单次的事务，被隐式提交.

> 查询autocommit状态: `show variables like 'autocommit';`, 状态值: 0=OFF; 1=ON

事务的启动时机:
`begin/start transaction`并不是一个事务的起点,在执行到它们之后的第一个操作 InnoDB 表的语句,事务才真正启动. 如果想要马上启动一个事务,可以使用`start transaction with consistent snapshot`命令.

InnoDB 里面每个事务有一个唯一的事务 ID,叫作 transaction id. 它是在事务开始的时候向InnoDB 的事务系统申请的,是按申请顺序严格递增的. 每行数据也都是有多个版本的. 每次事务更新数据的时候,都会生成一个新的数据版本,并且把 transaction id 赋值给这个数据版本的事务 ID即row trx_id.

每行数据的版本并不是物理上真实存在的,而是每次需要的时候根据当前版本和 undo log 计算出来的.

### 事务隔离的实现
ref:
- [](Innodb中的事务隔离级别和锁的关系)(https://tech.meituan.com/2014/08/20/innodb-lock.html)

  《高性能MySQL》第三版

MySQL 的隔离级别基于锁和 MVCC 机制共同实现的.

MVCC 是一种并发控制机制，用于在多个并发事务同时读写数据库时保持数据的一致性和隔离性。它是通过在每个数据行上维护多个版本的数据来实现的.

> 查看隔离级别: `SELECT @@tx_isolation/SELECT @@transaction_isolation (from v8.0)`

事务隔离的机制是通过视图（read-view即consistent read view, 一致性读视图）来实现的并发版本控制（MVCC），不同的事务隔离级别创建读视图的时间点不同:
- 可重复读是每个事务重建读视图，整个事务存在期间都用这个视图
- 读已提交是每条 SQL (开始执行时)创建read-view，隔离作用域仅限该条 SQL 语句.
- 读未提交是不创建，直接返回记录上的最新值
- 串行化隔离级别下直接用**加锁**的方式来避免并行访问
这里的视图可以理解为数据副本，每次创建视图时，是将当前已持久化的数据创建副本，后续直接从副本读取，从而达到数据隔离效果

> read-view用于支持RC ( Read Committed ,读提交)和 RR ( Repeatable Read ,可重复读)隔离级别的实现
> 每一次的修改操作，并不是直接对行数据进行操作,而是新加一行
> 在MySQL中，实际上每条记录在更新的时候都会同时记录一条回滚操作(undo 日志). 通过回滚操作可以得到前一个状态的值.

大多数对数据的变更操作包括INSERT/DELETE/UPDATE，其中INSERT操作在事务提交前只对当前事务可见，因此产生的Undo日志可以在事务提交后直接删除，而对于UPDATE/DELETE则需要维护多版本信息，在InnoDB里，UPDATE和DELETE操作产生的Undo日志被归成一类，即update_undo(from [MySQL · 引擎特性 · InnoDB undo log 漫游](http://mysql.taobao.org/monthly/2015/04/01/)).

**尽量不要使用长事务**: 长事务意味着系统里面会存在很老的事务视图. 由于这些事务随时可能需要访问数据库里面的任何数据，所以这个事务提交之前，数据库里面它可能用到的回滚记录都必须保留，这就会导致大量占用存储空间. 除了对回滚段的影响，长事务还占用锁资源，也可能拖垮整个库.

> 长事务查询: `select * from information_schema.innodb_trx where TIME_TO_SEC(timediff(now(),trx_started))>60`
> 监控 information_schema.Innodb_trx表，设置长事务阈值，超过就报警/或者kill(通过pt-kill)
> 通过SET MAX_EXECUTION_TIME命令，来控制每个语句执行的最长时间，避免单个语句意外执行太长时间.
> 使用的是MySQL 5.6或者更新版本，把innodb_undo_tablespaces设置成2（或更大的值）. 这样如果真的出现大事务导致回滚段过大，设置后清理起来更方便

MVCC 在 MySQL 中实现所依赖的手段主要是: 隐藏字段、read view、undo log. 在内部实现中，InnoDB 通过数据行的 DB_TRX_ID 和 Read View 来判断数据的可见性，如不可见，则通过数据行的 DB_ROLL_PTR 找到 undo log 中的历史版本. 每个事务读到的数据版本可能是不一样的，在同一个事务中，用户只能看到该事务创建 Read View 之前已经提交的修改和该事务本身做的修改.

InnoDB 存储引擎为每行数据添加了三个 隐藏字段：
- DB_TRX_ID（6字节）：表示最后一次插入或更新该行的事务 id. 此外，delete 操作在内部被视为更新，只不过会在记录头 Record header 中的 deleted_flag 字段将其标记为已删除
- DB_ROLL_PTR（7字节） 回滚指针，指向该行的 undo log 。如果该行未被更新，则为空
- DB_ROW_ID（6字节）：如果没有设置主键且该表没有唯一非空索引时，InnoDB 会使用该 id 来生成聚簇索引

ReadView 主要是用来做可见性判断，里面保存了`当前对本事务不可见的其他活跃事务`主要有以下字段:
- m_low_limit_id：目前出现过的最大的事务 ID+1，即下一个将被分配的事务 ID。大于等于这个 ID 的数据版本均不可见
- m_up_limit_id：活跃事务列表 m_ids 中最小的事务 ID，如果 m_ids 为空，则 m_up_limit_id 为 m_low_limit_id。小于这个 ID 的数据版本均可见
- m_ids：Read View 创建时其他未提交的活跃事务 ID 列表。创建 Read View时，将当前未提交事务 ID 记录下来，后续即使它们修改了记录行的值，对于当前事务也是不可见的。m_ids 不包括当前事务自己和已提交的事务（正在内存中）
- m_creator_trx_id：创建该 Read View 的事务 ID

undo log 主要有两个作用：
1. 当事务回滚时用于将数据恢复到修改前的样子
1. 另一个作用是 MVCC ，当读取记录时，若该记录被其他事务占用或当前版本对该事务不可见，则可以通过 undo log 读取之前的版本数据，以此实现非锁定读

在RR级别中，通过MVCC机制，虽然让数据变得可重复读，但读到的数据可能是历史数据，不是数据库当前的数据！这在一些对于数据的时效特别敏感的业务中，就很可能出问题. 对于这种读取历史数据的方式，也叫它快照读 (snapshot read)，而读取数据库当前版本数据的方式，叫当前读 (current read).

事务的隔离级别实际上都是**定义了当前读的级别**，MySQL为了减少锁处理（包括等待其它锁）的时间，提升并发能力，引入了快照读的概念，使得select不用加锁。而update、insert这些“当前读”，就需要另外的机制来解决了.

为了解决当前读中的幻读问题，MySQL事务使用了Next-Key锁(`=行锁和GAP（间隙锁）`). 行锁防止别的事务修改或删除，GAP锁防止别的事务新增，行锁和GAP锁结合形成的的Next-Key锁共同解决了RR级别在写数据时的幻读问题. 注意: 如果使用的是没有索引的字段, 即使没有匹配到任何数据, 那么会给全表加入gap锁, 除非该事务提交, 否则其它事务无法插入任何数据.

> 间隙锁是两个值间(不包含边界值)的空隙加锁, 用于解决RR幻读的机制, 仅在RR下生效.

## index(索引)
索引是为了提高数据查询的效率，起到书的目录一样的作用.

索引类型:
1. 哈希表, 一种以键-值（key-value）存储数据的结构. 适用于只有等值查询的场景, 不适合区间查询.
1. 有序数组, 在等值查询和范围查询场景中的性能就都非常优秀, 但不适合更新(比如插入和删除), 因此它只适用于静态存储引擎(数据不变化)
1. 树/跳表等.

> 在MySQL中，索引是在存储引擎层实现的，并没有统一的索引标准，即不同存储引擎的索引的工作方式并不一样.

常见索引列建议:
1. 出现在 SELECT、UPDATE、DELETE 语句的 WHERE 从句中的列
1. 包含在 ORDER BY、GROUP BY、DISTINCT 中的字段
1. 不要将符合 1 和 2 中的字段的列都建立一个索引，通常将 1、2 中的字段建立联合索引效果更好
1. 多表 join 的关联列

### innodb索引
InnoDB使用了B+树索引模型.

> 查看索引: `show index from ${table_name} \G`;

根据叶子节点的内容，索引类型分为主键索引和非主键索引:
- 主键索引的叶子节点存的是**整行数据**. 在InnoDB里，主键索引也被称为聚簇索引（clustered index）.
- 非主键索引的叶子节点内容是**主键的值**. 在InnoDB里，非主键索引也被称为二级索引（secondary index）.

基于主键索引和普通索引的查询区别:
- 如果语句是select * from T where ID=500，即主键查询方式，则只需要搜索ID这棵B+树
- 如果语句是select * from T where k=5，即普通索引查询方式，则需要先搜索k索引树，得到ID的值为500，再到ID索引树搜索一次. 这个过程称为回表.

**即基于非主键索引的查询需要多扫描一棵索引树. 因此，我们在应用中应该尽量使用主键查询**.

> 主键长度越小，普通索引的叶子节点(存储主键的值)就越小，普通索引占用的空间也就越小.
> 从性能和存储空间方面考量，自增主键往往是更合理的索引

优化回表:
- 使用覆盖索引
  只查询索引上的叶子节点(主键值,即自增主键或联合主键)内容, 比如`select ID from T where k between 3 and 5 -- k上有普通索引`, 此时不需要回表. 也就是说，在这个查询里面，索引k已经“覆盖了”我们的查询需求，其称为覆盖索引

  > 覆盖索引：就是包含了所有查询字段 (where、select、order by、group by 包含的字段) 的索引

  对于频繁的查询，优先考虑使用覆盖索引
- 联合索引时利用 B+索引的“最左前缀”，来定位记录
  最左前缀可以是联合索引的最左N个字段，也可以是字符串索引的最左M个字符.
  本质是**B+索引项是按照索引定义里面出现的字段顺序排序的, 即索引对顺序是敏感的**, 而且**MySQL的查询优化器会自动调整where子句的条件顺序以使用适合的索引**.

  > 创建一个联合索引(key1,key2,key3)时，相当于创建了（key1）、(key1,key2)和(key1,key2,key3)三个索引

- [MySQL 5.6 引入的索引下推优化（index condition pushdown)](https://www.cnblogs.com/hollischuang/p/11155447.html)， 可以在索引遍历过程中，**对索引中包含的字段先做判断，直接过滤掉不满足条件的记录**，减少回表次数.

  索引下推的 下推 其实就是指将部分上层（Server 层）负责的事情，交给了下层（存储引擎层）去处理

  在建立联合索引的时候，安排索引内的字段顺序的方法:
  1. 如果通过调整顺序，可以少维护一个索引，那么这个顺序往往就是需要优先考虑采用的
  1. 保证准确情况下, 节省空间
    ```sql
    CREATE TABLE `geek` (
      `a` int(11) NOT NULL,
      `b` int(11) NOT NULL,
      `c` int(11) NOT NULL,
      `d` int(11) NOT NULL,
      PRIMARY KEY (`a`,`b`), // 有了(`a`,`b`)后就不需要单独在a上建立索引了
      KEY `c` (`c`),
      KEY `ca` (`c`,`a`), // key `ca`与key `c`的存储内容一致都是(`c`,`a`,`b`,),因此可删除`ca`来节省空间
    )
    ```

MySQL 5.7 可以通过查询 sys 库的 schema_unused_indexes 视图来查询哪些索引从未被使用=

## 锁
根据加锁的范围，MySQL里面的锁大致可以分成全局锁、表级锁和行锁三类.

在5.5中, information_schema 库中增加了三个关于锁的表（MEMORY引擎）：
- innodb_trx         # 当前运行的所有事务
- innodb_locks       # 当前出现的锁
- innodb_lock_waits  # 锁等待的对应关系

特点:
- 行锁: 开销大，加锁慢；会出现死锁；锁定粒度最小，发生锁冲突的概率最低，并发度也最高
- 表级锁: 开销小，加锁快；不会出现死锁；锁定粒度大，发出锁冲突的概率最高，并发度最低

锁由存储引擎实现. InnoDB支持行级锁(row-level locking)和表级锁,默认为行级锁.

**InnoDB行锁是通过给索引上的索引项加锁来实现的. InnoDB这种行锁实现特点意味着：只有通过索引条件检索数据，InnoDB才使用行级锁，否则，InnoDB将使用表锁！**

> Oracle是通过在数据块中对相应数据行加锁来实现的

在InnoDB中，锁是逐步获得的，就造成了死锁的可能.

在MySQL中，行级锁并不是直接锁记录，而是锁索引。索引分为主键索引和非主键索引两种，如果一条sql语句操作了主键索引，MySQL就会锁定这条主键索引；如果一条语句操作了非主键索引，MySQL会先锁定该非主键索引，再锁定相关的主键索引。 在UPDATE、DELETE操作时，MySQL不仅锁定WHERE条件扫描过的所有索引记录，而且会锁定相邻的键值，即所谓的next-key locking.

当两个事务同时执行，一个锁住了主键索引，在等待其他相关索引。另一个锁定了非主键索引，在等待主键索引。这样就会发生死锁。

发生死锁后，InnoDB一般都可以检测到，并使一个事务释放锁回退，另一个获取锁完成事务。

有多种方法可以避免死锁，这里只介绍常见的三种
1. 如果不同程序会并发存取多个表，尽量约定以相同的顺序访问表，可以大大降低死锁机会
1. 在同一个事务中，尽可能做到一次锁定所需要的所有资源，减少死锁产生概率
1. 对于非常容易产生死锁的业务部分，可以尝试使用升级锁定颗粒度，通过表级锁定来减少死锁产生的概率

### 全局锁
`flush tables with read lock`(FTWRL), 让整个库处于只读状态, 之后其他线程的以下语句会被阻塞：数据更新语句（数据的增删改）、数据定义语句（包括建表、修改表结构等）和更新类事务的提交语句.

全局锁的典型使用场景是，做全库逻辑备份, 但支持事务的存储引擎有更好的方法:`mysqldump --single-transaction`. 它会启动一个事务，来确保拿到一致性视图, 而由于MVCC的支持，这个过程中数据是可以正常更新的, 但此时不要做DDL.

既然要全库只读，不使用`set global readonly=true`的原因:
1. 在有些系统中，readonly的值会被用来做其他逻辑，比如用来判断一个库是主库还是备库
1. 在异常处理机制上有差异. 如果执行FTWRL命令之后由于客户端发生异常断开，那么MySQL会自动释放这个全局锁，整个库回到可以正常更新的状态. 而将整个库设置为readonly之后，如果客户端发生异常，则数据库就会一直保持readonly状态,风险较高.

### 表级锁
MySQL里面表级别的锁有两种：一种是表锁，一种是元数据锁（meta data lock，MDL).

表锁的语法是 lock tables … read/write, 可以用`unlock tables`主动释放锁，也可以在客户端断开的时候自动释放.

> `lock/unlock tables`是在服务器级别而不在存储引擎中实现.  InnoDB 支持行级锁, 所以没必要使用 LOCK TABLES了.

> Lock table和事务之间的交互非常复杂, 并且在一些server版本中存在意想不到的行为. 建议除了在禁用autocommit的事务中使用外, 无论什么存储引擎都禁止显式使用lock tables

`LOCK TABLES products READ；` : 读锁定, 之后无论是当前线程还是其他线程均只能读操作，写操作全部被堵塞.
`LOCK TABLES products WRITE；`写锁定，之后只有当前线程可以进行读操作和写操作，**其他线程读操作和写操作均被堵塞**.

MDL（metadata lock), 不需要显式使用，在访问一个表的时候会被自动加上, 作用是保证读写的(DDL的)正确性, 但事务不提交，也会一直占着MDL锁.

如何安全地给小表加字段:
1. 在MySQL的information_schema 库的 innodb_trx 表中，查询当前执行中的事务. 如果要做DDL变更的表刚好有长事务在执行，要考虑先暂停DDL，或者kill掉这个长事务
1. 对于热点表, kill可能未必管用，因为新的请求马上就来. 可在alter table语句里面设定等待时间，如果在这个指定的等待时间里面能够拿到MDL写锁最好，拿不到也不要阻塞后面的业务语句，先放弃之后再重试即可.
    ```sql
    -- [maridb starting with 10.3.0](https://mariadb.com/kb/en/library/wait-and-nowait/) + alisql
    ALTER TABLE tbl_name NOWAIT add column ...
    ALTER TABLE tbl_name WAIT N add column ...
    ```

### 行级锁
MySQL 的行锁是在引擎层由各个引擎自己实现的. 但并不是所有的引擎都支持行锁,比如MyISAM 引擎就不支持行锁. 不支持行锁意味着并发控制只能使用表锁, 因此同一张表上任何时刻只能有一个更新在执行,这就会影响到业务并发度. InnoDB 是支持行锁的, 这也是 MyISAM 被 InnoDB 替代的重要原因之一.

在 InnoDB 事务中, 在事务执行期间, 随时都可以获取锁, 但锁只有在提交或回滚后才会释放, 并且所有的锁会同时释放, 这个就是两阶段锁协议(two-phase locking protocol).

>  InnoDB 支持通过特定的语句(`select ... for share/update`)进行显式锁定,这些语句不属于 SQL 规范.

因此如果事务中需要锁多个行,要**把最可能造成锁冲突、最可能影响并发度的锁尽量往后放**,比如顾客 A 要在影院 B 购买电影票, 涉及到以下操作：
- A. 从顾客 A 账户余额中扣除电影票价
- B. 给影院 B 的账户余额增加这张电影票价 // 并发高
- C. 记录一条交易日志

最合理的执行顺序是: `C,A,B`

#### 死锁和死锁检测
当并发系统中不同线程出现循环资源依赖,涉及的线程都在等待别的线程释放资源时,就会导致这几个线程都进入无限等待的状态,称为死锁.

处理死锁策略:
1. 直接进入等待,直到超时, 由innodb_lock_wait_timeout设置超时时间, 默认是50s.
1. 开启死锁检测(`innodb_deadlock_detect = on`, **推荐**),发现死锁后,主动回滚死锁链条中的某一个事务,让其他事务得以继续执行, 但该功能要耗费大量的 CPU 资源.

> 推荐在业务层用mq控制并发.

## cmd
- `mysql -h$ip -P$port -u$user -p` : 登录mysql
- [`show processlist`](https://www.docs4dev.com/docs/zh/mysql/5.7/reference/show-processlist.html) : 显示正在运行的线程

  其中的Command列显示为“Sleep”的行表示该连接是一个空闲连接

  client空闲太长时间时，连接器就会自动将它断开. 这个时间是由参数`wait_timeout`控制的，默认值是8小时.

  MySQL在执行过程中临时使用的内存是由该连接对象管理的, 这些资源会在连接断开的时候才释放. 随着持续使用, 内存占用太大，被系统强行杀掉（OOM），从现象看就是MySQL异常重启了, 解决方法:
  1. 定期断开长连接. 或者程序里面判断执行过一个占用内存的大查询后，断开连接重连
  1. 如果用的是MySQL 5.7或更新版本，可以在每次执行一个比较大的操作后，通过执行 mysql_reset_connection来重新初始化连接资源. 这个过程不需要重连和重新验证权限，但是会将连接恢复到刚刚创建时的状态.

- [`flush privileges`](https://dev.mysql.com/doc/refman/8.0/en/privilege-changes.html): 刷新内存中的权限内容

  当mysqld启动时，所有的权限都会被加载到内存中.

  如果使用GRANT/REVOKE/SET PASSWORD/RENAME USER命令来更改数据库中的权限表，mysqld服务器将会注意到这些变化并立即加载更新后的权限表至内存中，即权限生效.

  如果使用INSERT/UPDATE/DELETE语句更新权限表，则内存中的权限表不会感知到数据库中权限的更新，必须重启服务器或者使用FLUSH PRIVILEGES命令使更新的权限表加载到内存中，即权限需在重启服务器或者FLUSH PRIVILEGES之后方可生效.

  授权表重新加载会影响每个现有客户端会话的权限:
  1. 表和列权限更改将随客户端的下一个请求生效
  1. 数据库权限更改将在客户端下次执行`USE db_name`时生效
  1. 对于已连接的客户端，全局权限和密码不受影响, 这些更改仅在后续连接的会话中生效

  > 因此最保险的方法是`FLUSH PRIVILEGES`后重新连接.

## 窗口函数
窗口函数对一组查询行执行类似于聚合的操作。但是，聚合操作将查询行分组为一个单独的结果行，而窗口函数为每个查询行生成一个结果.

相关函数:
- DENSE_RANK(): 返回当前行在其分区中的排名，没有间隙。对等项被视为并列并获得相同的排名。此函数为对等组分配连续的排名；结果是大于一的组不产生不连续的排名号码.

## 存储过程
ref:
- [MySQL中的存储过程（详细篇）](https://zhuanlan.zhihu.com/p/679169773)

存储过程：（PROCEDURE）是事先经过编译并存储在数据库中的一段SQL语句的集合. **强烈不推荐使用**

基本语法:
- 存储过程中的参数分别是 in，out，inout三种类型:

  1. in代表输入参数（默认情况下为in参数），表示该参数的值必须由调用程序指定
  1. ou代表输出参数，表示该参数的值经存储过程计算后，将out参数的计算结果返回给调用程序
  1. inout代表即时输入参数，又是输出参数，表示该参数的值即可有调用程序制定，又可以将inout参数的计算结果返回给调用程序
- 存储过程中的语句必须包含在BEGIN和END之间。
- DECLARE中用来声明变量，变量默认赋值使用的DEFAULT，语句块中改变变量值，使用SET 变量=值

```sql
SHOW   {  PROCEDURE  |  FUNCTION  }  status -- 查看列表
SHOW   CREATE   {  PROCEDURE  |  FUNCTION  }  sp_name -- 看不到检查权限
```

## 表空间
> [mysql通过配置项：innodb_file_per_table指定MySQL使用独立表空间, MySQL5.6.6及以后的版本默认值是ON, MySQL5.6.5以前的版本默认值是OFF](http://uusama.com/922.html).
> [The CREATE TABLESPACE statement is not supported by MariaDB](https://mariadb.com/kb/en/create-tablespace/)

  创建tablespace语句(必须包含`ENGINE=xxx`)会成功, 但不会有实际效果.

```sql
mysql> create tablespace big_data_in_mysql add datafile 'first.ibd' [ENGINE [=] engine_name]; -- 未指定存储目录, 所以使用默认存储路径. 这个表空间所对应的元数据存放在first.ibd这个文件中
mysql> create tablespace test_tablespace add datafile '/var/lib/test_mysql_tablespace/first.ibd';
mysql> select * from information_schema.INNODB_SYS_TABLESPACES; -- 查看表空间
mysql> create table t2(c1 int primary key) tablespace test_tablespace;
mysql> DROP TABLESPACE tablespace_name [ENGINE [=] engine_name];
```

## FAQ
### 从 innodb 的索引结构分析,为什么索引的 key 长度不能太长
key 太长会导致一个页当中能够存放的 key 的数目变少,间接导致索引树的页数目变多,索引层次增加,从而影响整体查询变更的效率.

### MySQL 的数据如何恢复到任意时间点
**恢复到任意时间点以定时的做全量备份,以及备份增量的 binlog 日志为前提**. 恢复到任意时间点首先将全量备份恢复之后,再此基础上回放增加的 binlog 直至指定的时间点.

### 不推荐外键
外键与级联更新适用于单机低并发，不适合分布式、高并发集群；级联更新是强阻塞，存在数据库更新风暴的风险；外键影响数据库的插入速度.

### 主从延迟
只要不是主从双写, 该问题就存在.

原因:
1. 主从复制本身机制
2. 网络原因: 带宽不足或网络不稳定
3. 负载过高, master处理不过来: 修改复制配置
4. 级联复制

### GROUP_CONCAT
```sql
SELECT
    user_no,
    SUBSTRING_INDEX(GROUP_CONCAT( DISTINCT role_name ORDER BY role_id desc ),',',2) AS role_name
FROM
    report_user_role_info
GROUP BY
    user_no;
```

### with rollup
在最后一行添加一条汇总记录

```sql
SELECT coalesce(name, '总金额'),name, SUM(money) as money FROM  test GROUP BY name WITH ROLLUP;
```

### 子查询提取with as
```sql
with
e as (select * from scott.emp),
d as (select * from scott.dept)
select * from e, d where e.deptno = d.deptno;
```

### 优雅处理数据`插入/更新`时`主键或唯一键`冲突
- ignore: 有则忽略, 无则插入

  `insert ignore into ...`
- replace: 有则删除再插入, 无则插入

  `replace into ...`
- on duplicate key update: 有则更新, 无则插入

  `insert into t (...) values (...) on duplicate key update size = size + 10`