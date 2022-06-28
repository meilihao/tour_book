# backup
ref:
- [PostgreSQL 最佳实践 - 在线增量备份与任意时间点恢复](https://developer.aliyun.com/article/59359)
- [Postgres数据库归档热备份](https://www.escapelife.site/posts/b47f1fcb.html)
- [PostgreSQL Plugin](https://docs.bareos.org/master/TasksAndConcepts/Plugins.html#postgresql-plugin)
- [BareosFdPluginPostgres.py](https://github.com/bareos/bareos/blob/master/core/src/plugins/filed/python/postgres/BareosFdPluginPostgres.py)
- [PostgreSQL 指南：内幕探索](https://cloud.tencent.com/developer/article/1459719)

冷备份, 以及逻辑备份都是某一个时间点的备份, 没有增量的概念.

如果数据库在运行过程中发生故障, 使用逻辑备份只能将数据库还原到备份时刻, 无法恢复到故障发生前的那个时刻.

在使用 PostgreSQL 数据库在做写入操作时，对数据文件做的任何修改信息，首先会写入 WAL 日志(预写日志), 然后才会对数据文件做物理修改。当数据库服务器掉重启时，PostgreSQL 数据库在启动时会优先读取 WAL日志，对数据文件进行恢复。因此，从理论上讲，如果我们有一个数据库的基础备份(也称为全备)，再配合 WAL 日志，是可以将数据库恢复到过去的任意时间点上.

在PostgreSQL中, 自8.0版本开始提供了在线的全量物理备份，整个数据库集簇（即物理备份数据）的运行时快照被称为基础备份.

PostgreSQL还在8.0版中引入了时间点恢复（Point-In-Time Recovery，PITR）. 这一功能可以将数据库恢复至任意时间点，这通过使用一个基础备份和由持续归档生成的归档日志来实现.

> 从PostgreSQL10 开始将"pg_xlog"目录重命名为"pg_wal"

在PG12以前, recovery.conf该文件的存在即触发恢复模式

从PG12开始，recovery.conf被废弃, recovery.conf文件中的参数放到了postgresql.conf配置文件中, 如果文件存在，服务器将不会启动. 同时recovery.conf由下面两个新文件进行替换:
- recovery.signal, 告诉PostgreSQL进入正常的归档恢复
- standby.signal, 告诉PostgreSQL进入standby模式.

如果两个文件都存在，则standby.signal优先. **恢复完成后将删除recovery.signal或standby.signal文件，但其中的参数postgresql.conf仍然保留。只要PostgreSQL不在恢复中，它们就会被忽略，但是最好用“#”注释禁用它们**.

## 备份
WAL归档，其实就是把pg_wal里面的日志备份出来，当系统故障后可以通过归档的日志文件对数据进行恢复.

1. 配置postgresql.conf:
```conf
wal_level = replica # 该参数的可选的值有minimal，replica和logical，wal的级别依次增高，在wal的信息也越多. 由于minimal这一级别的wal不包含从基础的备份和wal日志重建数据的足够信息，在该模式下，无法开启wal日志归档. replica在9.5及之前版本为hot_standby.
archive_mode = on # 打开归档备份
archive_command = 'test ! -f /var/lib/postgresql/12/wal_archive/%f && cp %p /var/lib/postgresql/12/wal_archive/%f' # 用于将日志拷贝到其他的地方, **注意备份的目录权限(这里即wal_archive)需要设置为数据库启动的用户权限**. 在shell脚本或命令中可以用 “%p” 表示将要归档的wal文件包含完整路径的信息的文件名，用“%f” 代表不包含路径信息的wal文件的文件名. 修改archive_command参数不需要重启，只需要reload数据库配置(`SELECT pg_reload_conf();`)即可
```
1. restart pg
1. 全量备份

	```bash
	# 主数据库目录
	/var/lib/postgresql/12/main

	# 创建备份目录
	$ mkdir -pv /var/lib/postgresql/12/pg_wal_backup/{basebackup, wal}

	# 进入备份目录做一次基础备份
	# 用于获取正在运行的PostgreSQL数据库集群的基本备份
	# 不会影响数据库的其他客户端，可用于时间点恢复或流复制备用服务器的起点
	$ pg_basebackup -U postgres -h 127.0.0.1 -p 5432 -W -D /var/lib/postgresql/12/pg_wal_backup/basebackup
	```

## 恢复
恢复大致流程如下：
1. 如果服务器仍在运行，停止它
1. 如果具有足够的空间，将整个集簇数据目录和表空间复制到一个临时位置，稍后将用到它们
1. 移除所有位于集簇数据目录和正在使用的表空间根目录下的文件和子目录
1. 在集簇数据目录中创建一个 recovery.conf 恢复命令文件
1. 启动服务器，服务器将会进入到恢复模式并且进而根据需要读取归档 WAL 文件。
1. 检查数据库的内容来确保你已经恢复到了期望的状态


```bash
# --- 配置postgresql.conf开启wal归档, 参考上面的`备份`
# --- 建立模拟数据
$ sudo su postgres -c psql
# create database app;
# \c app;
# create table test01(id int primary key,name varchar(20));
# insert into test01 values(1,'a'), (2,'b'), (3,'c');
# \q
$ psql -U postgres -h 127.0.0.1 -p 5432 -c "select pg_start_backup(now()::text)"
$ rsync -acvz -L --exclude "pg_wal" --exclude "pg_xlog" --exclude "pg_log" $PGDATA /xxx
$ psql -U postgres -h 127.0.0.1 -p 5432 -c "select pg_stop_backup()"
# --- 模拟故障
$ psql -U postgres -h 127.0.0.1 -p 5432 -c "select pg_create_restore_point('$(date +\"\%Y\%m\%d\%H\%M\")');" # 创建一个还原点
$ psql -U postgres -h 127.0.0.1 -p 5432 -d app -c "delete from test01;" # 模拟误操作: 删表
$ psql -U postgres -h 127.0.0.1 -p 5432 -d app -c "select pg_switch_wal();" # 确保之前操作均已归档
$ sudo systemctl stop postgresql
# --- 删除当前数据库并将备份数据导入
$ rm -rf /var/lib/postgresql/12/main/*
$ cp -r  /var/lib/postgresql/12/pg_wal_backup/basebackup/* /var/lib/postgresql/12/main
$ cp -r  /var/lib/postgresql/12/pg_wal_backup/wal/* /var/lib/postgresql/12/main/pg_wal
# --- 准备recovery.conf
$ cp /usr/share/postgresql/12/recovery.conf.sample /var/lib/postgresql/12/main/recovery.conf
$ chmod 0600 /var/lib/postgresql/12/main/recovery.conf
$ vim /var/lib/postgresql/12/main/recovery.conf
restore_command='cp /var/lib/postgresql/12/main/pg_wal/%f %p' # 从归档目录恢复日志
recovery_target_name='2019-12-12 12:12:00' # 指定归档时间点，如没指定恢复到故障前的最后一完成的事务
recovery_target_timeline='latest' # 指定归档时间线，latest代表最新的时间线分支，如没指定恢复到故障前的pg_control里面的时间线
# ---- 修改postgresql.conf配置文件将wal备份停止
$ vim /etc/postgresql/12/main/postgresql.conf
archive_mode='off'
# ---- 启动db
$ sudo systemctl start postgresql # 监控日志. 将recovery.conf变成了recovery.done文件. 恢复后pg只接收只读连接, 需要执行"select pg_wal_replay_resume"解除该状态
$ psql -U postgres -h 127.0.0.1 -p 5432 -d app -c "select pg_wal_replay_resume();" # 解除暂停状态. 当然recovery.conf可以移除，避免下次数据重启数据再次恢复到该还原点
```

## 数据复制方式
PostgreSQL 中复制的三种方法:

|POSTGRES 的类型|谁这样做|主要好处|
|简单的流式复制|本地/EC2|更易于设置; 高 I/O 性能; 大容量存储|
|复制块设备|RDS/Azure|适用于数据在云环境中的持久性|
|从 WAL 重建|Heroku/Citus|后台节点重建; 启用 fork 和 PITR|

1. 流复制

	最常见的方法. 所谓流复制，就是备服务器通过tcp流从主服务器中同步相应的数据，主服务器在WAL记录产生时即将它们以流式传送给备服务器，而不必等到WAL文件被填充

	有一个主节点，主节点具有表的数据和预写日志(WAL). 当修改 Postgres 中的行时，更改首先会被提交到仅附加重做日志，此重做日志称为预写日志或WAL. 然后，此 Postgres 的 WAL 日志将流式传输到辅助节点.

	在此方法中，当构建新的辅助节点时，新的辅助节点需要从主节点重播整个状态 - 从时间开始。然后，重放操作可能在主节点上引入显着负载。如果数据库的主节点提供服务，则此负载变得更加重.
1. 复制块设备

	在此方法中，更改将写入持久 volume. 然后，此 volume 将同步镜像到另一个 volume 中. 这种方法的好处是它适用于所有关系数据库, 可以将它用于 MySQL，PostgreSQL 或 SQL Server.

	但是，Postgres 中的磁盘镜像复制方法还要求复制表和 WAL 日志数据. 此外，现在每次写入数据库都需要同步通过网络. 不能错过任何一个字节，因为这可能会使数据库处于损坏状态.
1. 从 WAL 重建(并切换到流复制)

	将复制和灾难恢复过程彻底改变。写入在主节点，主节点每天执行完整数据库备份，每 60 秒执行一次增量备份。

	当需要构建新的辅助节点时，辅助节点会从备份重建其整个状态。这样，不会在主数据库上引入任何负载。可以启动新的辅助节点并从 S3 / Blob 存储重建它们. 当辅助节点足够接近主节点时，可以从主节点开始流式传输 WAL 日志并赶上它. 在正常状态下，辅助节点跟随主节点.

	在这种方法中，预写日志优先。这种设计适用于更加云原生的架构。可以随意调出或击落副本，而不会影响关系数据库的性能, 还可以根据需要使用同步或异步复制.


## 时间线
### 引入时间线
PostgreSQL 引入了时间线的概念. 每当归档文件恢复完成后，创建一个新的时间线用来区别新生成的 WAL 记录, 以避免当某次recovery后重新生成的wal归档与recovery前生成的wal归档的文件重名.

时间线 ID 号是 WAL 文件名组成之一，因此一个新的时间线不会覆盖由以前的时间线生成的 WAL. 每个时间线类似一个分支，在当前时间线的操作不会对其他时间线 WAL 造成影响，有了时间线，就可以恢复到之前的任何时间点.

![/misc/img/sql/pg/postgres-timeline.png]

WAL 文件名由时间线和日志序号组成，源码实现如下:
```c
#define XLogFileName(fname, tli, log, seg) \
    snprintf(fname, XLOG_DATA_FNAME_LEN + 1, "%08X%08X%08X", tli, log, seg)
```

### 新时间线的出现场景
1. 即时恢复(PITR)配置 recovery.conf 文件

	设置好 recovery.conf 文件后，启动数据库，将会产生新的 timeline，而且会生成一个新的 history 文件. 恢复的默认行为是沿着与当前基本备份相同的时间线恢复. 如果想恢复到某些时间线，则需要指定的 recovery.conf 目标时间线 recovery_target_timeline，不能恢复到早于基本备份分支的时间点.
1. standby promote

	搭建一个 PG 主备，然后停止主库，在备库机器执行promote. 这时候备库将会升为主备，同时产生一个新的 timeline，同样生成一个新的 history 文件.

### history 文件
每次创建一个新的时间线，PostgreSQL 都会创建一个“时间线历史”文件，文件名类似 .history，它里面的内容是由原时间线 history 文件的内容再追加一条当前时间线切换记录. 假设数据库恢复启动后，切换到新的时间线 ID＝5，那么文件名就是 00000005.history，该文件记录了自己从什么时间哪个时间线什么原因分出来的，该文件可能含有多行记录，每个记录的内容格式如下所示:
```conf
* <parentTLI> <switchpoint> <reason>
*      parentTLI       ID of the parent timeline
*      switchpoint     XLogRecPtr of the WAL position where the switch happened
*      reason          human-readable explanation of why the timeline was changed
```

当数据库在从包含多个时间线的归档中恢复时，这些 history 文件允许系统选取正确的 WAL 文件. 当然，它也能像 WAL 文件一样被归档到 WAL 归档目录里。历史文件只是很小的文本文件，所以保存它们的代价很小.

当在 recovery.conf 指定目标时间线 tli 进行恢复时，程序首先寻找 .history 文件，根据 .history 文件里面记录的时间线分支关系，找到从 pg_control 里面的 startTLI 到 tli 之间的所有时间线对应的日志文件，再进行恢复.


## checkpoint
Checkpoint的概念太重要了，无论在Oracle，MySQL这个概念，**它主要功能是在检查点时刻， 脏数据全部刷新到磁盘，以实现数据的一致性和完整性**. PostgreSQL为什么要设计Checkpoint呢？跟Oracle一样，其主要目的是缩短崩溃恢复时间. PostgreSQL在崩溃恢复时会以最近的Checkpoint为基础，不断应用这之后的WAL日志.

## FAQ
### 触发WAL归档
1. `pg_switch_xlog()`

	执行pg_switch_xlog()后，WAL会切换到新的日志，这时会将老的WAL日志归档

	pg_switch_xlog执行完后会在$PGDATA/pg_wal/archive_status下生成相应wal已归档的标志.
1. WAL日志被写满后会触发归档

	WAL 日志文件的默认大小为 16MB，这个值可以在编译 PostgreSQL 时通过参数`--with-wal-segsize` 参数更改，编译后不能修改.
1. 设置 archive_timeout 参数

	假如设置 archive_timeout=60，那么每 60s 会触发一次 WAL 日志切换，同时触发日志归档. 这里有个隐含的假设，就是当前 WAL 日志中仍有未归档的 WAL 日志内容.

	尽量不要把 archive_timeout 设置的很小，如果设置的很小，它会膨胀归档存储。因为，强制归档的日志，即使没有写满，也会是默认的 wal-segsize M 大小(比如wal日志写满的默认大小为16M).

### 全量备份当前数据库
`pg_dump -U postgres -p 5432 -d draft -Z9 \
      -f /data/pg_wal_backup/pg_dump/app_`date +%Y%m%d%H%M%S`.sql.gz`

### pg_dump和pg_basebackup区别
pg_dump创建一个逻辑备份，即一系列 SQL 语句，当执行这些语句时，会创建一个逻辑上与原始数据库类似的新数据库.

pg_basebackup创建物理备份，即**整个数据库集群的文件的副本**而不是单个database.

[pg_dump、pg_dumpall](https://www.cnblogs.com/haha029/p/15650047.html)都是是逻辑备份，前者支持多种备份格式，后者只支持sql文本. 因为**pg_dump和pg_dumpall不会产生物理备份， 因此不能用于连续归档方案**.

> pg_basebackup整合了 pg_start_backup和pg_stop_backup命令.

pg_dump:
```bash
--- 备份postgres库
# pg_dump -h 127.0.0.1 -p 5432 -U postgres -f postgres.sql --column-inserts

--- 备份postgres库并tar打包
# pg_dump -h 127.0.0.1 -p 5432 -U postgres -f postgres.sql.tar -Ft
 
--- 只备份postgres库对象数据
# pg_dump -U postgres -d postgres -f postgres.sql -Ft --data-only --column-inserts
 
--- 只备份postgres库对象结构
# pg_dump -U postgres -d postgres -f postgres.sql -Ft --schema-only

--- 导入SQL文件
# psql -f postgre.sql postgres postgres
```

pg_dumpall:
```bash
--- 备份postgres库，转储数据为带列名的INSERT命令
# pg_dumpall -d postgres -U postgres -f postgres.sql --column-inserts

--- 备份postges库，转储数据为INSERT命令
# pg_dumpall -d postgres -U postgres -f postgres.sql --inserts
 
--- 导入数据
# psql -f postgres.sql
```

### pg_start_backup
pg_start_backup备份方法跟pg_basebackup不同，被称为非排他低级基础备份.

执行基础备份(即全备)最简单方式是使用pg_basebackup, 它使用普通文件或者tar归档创建基础备份. 如果需更多的灵活性, 也可以使用 低水平API pg_start_backup 创建基础备份.

> 非排他：就是**备份的时候不影响其读写**

> 低级基础：就是使用低级API来实现. 像pg_basebackup是全自动的；而pg_start_backup需要你自己指定什么时候开始，什么时候结束，不是全自动的，被称为低级API

> [pg_start_backup/pg_stop_backup src](https://github.com/postgres/postgres/blob/master/src/backend/access/transam/xlogfuncs.c)

pg_start_backup备份流程:
1. pg_start_backup('label',false,false)

	$PGDATA中记录一个标签文件叫backup_label, 其包含标签名label(需唯一), 以及执行这条指令的启动时间, START WAL LOCATION, CHECKPOINT LOCATION等.

	它会执行一次checkpoint操作.

	参数:
	1. label：唯一标识这次备份操作的任意字符串
	1. 第二个参数：是否尽快完成备份, 默认情况下（就是设置为false）

		pg_start_backup需要较长的时间来完成。因为会执行一个检查点，该检查点所需要的IO将会分散到一段时间内，默认是检查点间隔（checkpoint_completion_target参数）的一半，可以最小化对其他会话查询的影响。如果想尽快的开始备份，则把false改为true。这将会立即发出一个检查点并且使用尽可能多的I/O.
	1. 第三个参数：默认false的情况下，意思是告诉pg_start_backup开始一次非排他基础备份，默认就行
2. 进行文件系统级别的备份（用tar、cpio、rsync等命令，非pg_dump、pg_dumpall）

	备份时pg_wal,pg_log,pg_xlog可不用备份.

	比如:
	1. `rsync -acvz -L --exclude "pg_wal" --exclude "pg_xlog" --exclude "pg_log" $PGDATA /xxx`
	1. cp
	1. `pg_basebackup -D ./basebackup_20130527 -F p -X stream -h 172.16.3.33 -p 1999 -U replica`

	确保备份包含了$PGDATA下的所有文件. **如果使用了不在此目录下的表空间，也得把这些表空间备份. 并且确保备份将表空间的软连接归档为链接，否则恢复过程将破坏表空间.**

3. pg_stop_backup

	终止备份模式，并且执行一个自动切换到下一个wal段。进行切换的原因是将在备份期间生成的最新wal段文件安排为可归档.

	pg_stop_backup结束时会输出WAL位置, 通过`select pg_xlogfile_name('F/3001A20');`可获取对应的WAL文件, **不知pg_stop_backup是否会触发wal归档, 但很多blog会在其后立即执行一次`select pg_switch_wal();`**

这样就完成了一次基础备份. 该基础备份跟pg_basebackup差不多，不过pg_basebackup都封装好了，也不用用操作系统命令进行备份，它自己就全部搞完了.

pg_start_backup特性:
1. 允许其他并发备份进行，即包括其他pg_start_backup的备份，也包括pg_basebackup的备份
2. wal归档启用并正常工作
3. 使用pg_start_backup的用户必须有该函数的权限，和连接到服务器上的权限，一般用postgres用户
4. 不在乎连接到哪个数据库，只要能发出命令就行
5. 调用pg_start_backup的会话连接必须保持到备份结束，否则备份将自动被终止

### recovery.conf
PostgreSQL 支持指定还原点的位置, 即数据库恢复到什么位置停下来.

4个recovery.conf参数控制恢复停止在哪个位置:
1. recovery_target_name

	指pg_create_restore_point(text)创建的还原点, 如果有重名的还原点, 那么在recovery过程中第一个遇到的还原点即停止.

2. recovery_target_time

	指XLOG中记录的recordXtime(xl_xact_commit_compact->xact_time), 配合recovery_target_inclusive使用,

	如果在同一个时间点有多个事务回滚或提交, 那么recovery_target_inclusive=false则恢复到这个时间点第一个回滚或提交的事务(含), recovery_target_inclusive=true则恢复到这个时间点最后一个回滚或提交的事务(含).   

	如果时间点上刚好只有1个事务回滚或提交, 那么recovery_target_inclusive=true和false一样, 恢复将处理到这个事务包含的xlog信息(含).   

	如果时间点没有匹配的事务提交或回滚信息, 那么recovery_target_inclusive=true和false一样, 恢复将处理到这个时间后的下一个事务回滚或提交的xlog信息(含).

3. recovery_target_xid

	指XLogRecord->xl_xid, 可以配合recovery_target_inclusive使用, 但是recovery_target_inclusive只影响日志的输出, 并不影响恢复进程截至点的选择, 截至都截止于这个xid的xlog位置. 也就是说无论如何都包含了这个事务的xlog信息的recovery.

	这里需要特别注意xid的信息体现在结束时, 而不是分配xid时. 所以恢复到xid=100提交|回滚点, 可能xid=102已经先提交了. 那么包含xid=102的xlog信息会被recovery. 
4. recovery_target_inclusive

	- 如果在同一个时间点有多个事务回滚或提交，那么recovery_target_inclusive=false则恢复到这个时间点第一个回滚或提交的事务(含)，recovery_target_inclusive=true则恢复到这个时间点最后一个回滚或提交的事务(含)。
	- 如果时间点上刚好只有1个事务回滚或提交，那么recovery_target_inclusive=true和false一样，恢复将处理到这个事务包含的xlog信息(含)。
	- 如果时间点没有匹配的事务提交或回滚信息，那么recovery_target_inclusive=true和false一样，恢复将处理到这个时间后的下一个事务回滚或提交的xlog信息(含)。

### pg_waldump
pg_waldump是PG 用来对 wal日志进行查看.
