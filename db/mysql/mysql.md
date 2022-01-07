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
1. Server层包括连接器、查询缓存、分析器、优化器、执行器等，涵盖MySQL的大多数核心服务功能，以及所有的内置函数（如日期、时间、数学和加密函数等），所有跨存储引擎的功能都在这一层实现，比如存储过程、触发器、视图等.
1. 存储引擎层负责数据的存储和提取. 其架构模式是插件式的，支持InnoDB(**推荐**)、~~MyISAM(不推荐)~~、Memory等多个存储引擎. 现在最常用的存储引擎是InnoDB，它从MySQL 5.5.5版本开始成为了默认存储引擎

## 安装
- [mariadb 离线包下载](https://mariadb.com/downloads/)

## 查询

> 由执行器检查是否有查询权限.
>
> 慢查询日志中的rows_examined字段表示这个语句执行过程中执行器扫描了多少行. 在有些场景下，执行器调用一次，在引擎内部则扫描了多行，因此引擎扫描次数跟rows_examined并不是完全相同的.

### 查询缓存
大多数情况下建议**不要使用查询缓存**: **查询缓存的失效非常频繁，只要有对一个表的更新，这个表上所有的查询缓存都会被清空; 同时MySQL 8.0版本已砍掉查询缓存功能**.

## 更新
update语句执行流程(浅色框表示是在InnoDB内部执行的，深色框表示是在执行器中执行):
![](/misc/img/mysql/2e5bff4910ec189fe1ee6e2ecc7b4bbe.png)

更新涉及两个重要的日志模块：redo log（重做日志）和 binlog（归档日志）.将redo log的写入拆成了两个步骤：prepare和commit，这就是"两阶段提交", 这是让这两个log状态保持逻辑上的一致.

### redo log
WAL的全称是Write-Ahead Logging，它的关键点就是先写日志，再写磁盘. 它工作在存储引擎层.

具体来说，当有一条记录需要更新的时候，InnoDB引擎就会先把记录写到redo log 里面，并更新内存，这样更新就算完成了. 同时，InnoDB引擎会在适当(空闲)的时候，将这个操作记录更新到磁盘里面, 并清理redo log.

InnoDB的redo log是固定大小的, 从头开始写，写到末尾就又回到开头的循环写, 对此它使用了两个游标:
- write pos是当前记录的位置，一边写一边后移并再循环
- checkpoint是当前要清理的位置，也是往后推移并再循环的，但清理记录前会确保记录已更新到数据文件.

`write pos -> checkpoint`之间即为空闲可用部分. 如果write pos追上checkpoint，表示WAL已满，这时候不能再执行新的更新，得停下来先清理一些记录，把checkpoint推进一下.

有了redo log，InnoDB就可以保证即使数据库发生异常重启，之前提交的记录都不会丢失，这个能力称为crash-safe.

参数`innodb_flush_log_at_trx_commit`可控制N次事务后redo log持久化到磁盘, 通常是`1`以保证MySQL异常重启之后数据不丢失.

### binlog
它工作在Server层.

参数`sync_binlog`表示N次事务后binlog持久化到磁盘, 通常是`1`以保证MySQL异常重启之后binlog不丢失.

### redo log 与 binlog
不同点:
- redo log是InnoDB引擎特有的；binlog是MySQL的Server层实现的，所有引擎都可以使用
- redo log是物理日志，记录的是`在某个数据页上做了什么修改`；binlog是逻辑日志，记录的是这个语句的原始逻辑，比如`给ID=2这一行的c字段加1`.
- redo log是循环写的，空间固定；binlog是可以追加写入的. `追加写`是指binlog文件写到一定大小后会切换到下一个，并不会覆盖以前的日志.

## 事务
MySQL的事务启动方式：
1. 显式启动事务语句，`begin/start transaction`+`commit`+`rollback`
  具体操作: 保存 autocommit 的当前状态，然后 start transaction，直到 commit or rollback 结束本次事务，再恢复之前保存的 autocommit 的状态
1. `set autocommit=0`会关闭这个线程的自动提交, 这相当于开启一个全局的事务且这个事务持续存在直到主动执行commit 或 rollback，或者断开连接. 在 mysql 的事务中，默认 `autocommit = 1`，每一次 sql 操作都被认为是一个单次的事务，被隐式提交.

> 查询autocommit状态: `show variables like 'autocommit';`, 状态值: 0=OFF; 1=ON

事务的启动时机:
`begin/start transaction`并不是一个事务的起点,在执行到它们之后的第一个操作 InnoDB 表的语句,事务才真正启动. 如果想要马上启动一个事务,可以使用`start transaction with consistent snapshot`命令.

InnoDB 里面每个事务有一个唯一的事务 ID,叫作 transaction id. 它是在事务开始的时候向InnoDB 的事务系统申请的,是按申请顺序严格递增的. 每行数据也都是有多个版本的. 每次事务更新数据的时候,都会生成一个新的数据版本,并且把 transaction id 赋值给这个数据版本的事务 ID即row trx_id.

每行数据的版本并不是物理上真实存在的,而是每次需要的时候根据当前版本和 undo log 计算出来的.

### 事务隔离的实现
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

## index(索引)
索引是为了提高数据查询的效率，起到书的目录一样的作用.

索引类型:
1. 哈希表, 一种以键-值（key-value）存储数据的结构. 适用于只有等值查询的场景, 不适合区间查询.
1. 有序数组, 在等值查询和范围查询场景中的性能就都非常优秀, 但不适合更新(比如插入和删除), 因此它只适用于静态存储引擎(数据不变化)
1. 树/跳表等.

> 在MySQL中，索引是在存储引擎层实现的，并没有统一的索引标准，即不同存储引擎的索引的工作方式并不一样.

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
- 联合索引时利用 B+索引的“最左前缀”，来定位记录
  最左前缀可以是联合索引的最左N个字段，也可以是字符串索引的最左M个字符.
  本质是**B+索引项是按照索引定义里面出现的字段顺序排序的, 即索引对顺序是敏感的**, 而且**MySQL的查询优化器会自动调整where子句的条件顺序以使用适合的索引**.

  > 创建一个联合索引(key1,key2,key3)时，相当于创建了（key1）、(key1,key2)和(key1,key2,key3)三个索引

- [MySQL 5.6 引入的索引下推优化（index condition pushdown)](https://www.cnblogs.com/hollischuang/p/11155447.html)， 可以在索引遍历过程中，**对索引中包含的字段先做判断，直接过滤掉不满足条件的记录**，减少回表次数.

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

## 锁
根据加锁的范围，MySQL里面的锁大致可以分成全局锁、表级锁和行锁三类.

在5.5中, information_schema 库中增加了三个关于锁的表（MEMORY引擎）：
- innodb_trx         # 当前运行的所有事务
- innodb_locks       # 当前出现的锁
- innodb_lock_waits  # 锁等待的对应关系

### 全局锁
`flush tables with read lock`(FTWRL), 让整个库处于只读状态, 之后其他线程的以下语句会被阻塞：数据更新语句（数据的增删改）、数据定义语句（包括建表、修改表结构等）和更新类事务的提交语句.

全局锁的典型使用场景是，做全库逻辑备份, 但支持事务的存储引擎有更好的方法:`mysqldump --single-transaction`. 它会启动一个事务，来确保拿到一致性视图, 而由于MVCC的支持，这个过程中数据是可以正常更新的, 但此时不要做DDL.

既然要全库只读，不使用`set global readonly=true`的原因:
1. 在有些系统中，readonly的值会被用来做其他逻辑，比如用来判断一个库是主库还是备库
1. 在异常处理机制上有差异. 如果执行FTWRL命令之后由于客户端发生异常断开，那么MySQL会自动释放这个全局锁，整个库回到可以正常更新的状态. 而将整个库设置为readonly之后，如果客户端发生异常，则数据库就会一直保持readonly状态,风险较高.

### 表级锁
MySQL里面表级别的锁有两种：一种是表锁，一种是元数据锁（meta data lock，MDL).

表锁的语法是 lock tables … read/write, 可以用`unlock tables`主动释放锁，也可以在客户端断开的时候自动释放.

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

在 InnoDB 事务中,行锁是在需要的时候才加上的,但并不是不需要了就立刻释放,而是要等到事务结束时才释放, 这个就是两阶段锁协议.

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

## FAQ
### 从 innodb 的索引结构分析,为什么索引的 key 长度不能太长
key 太长会导致一个页当中能够存放的 key 的数目变少,间接导致索引树的页数目变多,索引层次增加,从而影响整体查询变更的效率.

### MySQL 的数据如何恢复到任意时间点
**恢复到任意时间点以定时的做全量备份,以及备份增量的 binlog 日志为前提**. 恢复到任意时间点首先将全量备份恢复之后,再此基础上回放增加的 binlog 直至指定的时间点.