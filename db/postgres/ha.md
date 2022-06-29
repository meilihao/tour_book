# ha
ref:
- <<PostgresSQL实战>>

## base
### 主备切换的触发方式
PostgreSQL热备(HOT-Standby)如果主库出现异常，备库有两种方法切换到主库:
1. 备库配置文件 recovery.conf 中有配置项 trigger_file ，它是激活从库的触发文件，当它存在就会激活从库

	`trigger_file = '/var/lib/pgsql/10/data/trigger_standby'时, `touch /var/lib/pgsql/10/data/trigger_standby`即可激活备库
1. 在备库上执行 pg_ctl promote 命令激活(from 9.1)

### 主备状态标志
数据字典表pg_stat_replication 、命令 pg_controldata、进程、自带的函数 pg_is_in_recovery() 都可以区别或判断实例的主备状态:
1. pg_controldata

	pg运行时执行`pg_controldata -D /var/lib/pgsql/data | grep 'Database cluster state'`, 输出:
	- `Database cluster state:               in production`是主库
	- `Database cluster state:               in archive recovery`是备库
1. `select pg_is_in_recovery()`: 备库返回`t`; 主库返回`f`
1. 进程

	` ps -ef | grep "wal"`输出"wal sender process"是主库;  "wal receiver process"是备库

备库激活后可以插入数据, 即变为可读写. 这时配置文件 recovery.conf 变为 recovery.done
### 模拟宕机
- `systemctl stop postgresql.service/pg_ctl stop`
- `关机`

## 配置
Hot-Standby 切换步骤比较多，有些配置可以提前做好的，例如 .pgpass, pg_hba.conf 等, 其他注意点:
1. 主，备切换时，务必先将主库关闭，否则一旦从库被激活时，而主库尚未关闭，会有问题

	pg关闭会做一次checkpoint, 完成后通知wal发送进程要关闭, wal进程会将截止此次checkpoint的wal日志流发送给备库的wal接收进程, 备节点接收到主库最后发送的wal并应用后, 达到与主库一致的状态.

	如果激活主库时忘记关闭老主库, 老主库此时可能不能直接切换成备库, 此时需要pg_rewind. pg_rewind与pg_basebakup的主要区别是它只复制变化的数据, 而不是全量同步
1. 主要文件权限, 所有者应是postgres
1. 使用replication slot时, 主备分别指定不同primary_slot_name

**生产环境中一主一备通常不采用同步复制, 因为备库宕机会对生产系统造成严重影响. 在一主多备的架构中才会采用同步复制.**

### pg 10 主备切换/Hot-Standby
env:
```
# --- 均关闭了防火墙
test1 192.168.1.11 # 最开始是主
test2 192.168.1.12
```

> 没有特别说明, 仅在主库上操作, 因为备库会同步主库数据目录.

1. `postgresql-setup --initdb`
1. 启动主库, 配置管理员密码和创建复制账号

	1. 修改postgres用户密码: `ALTER USER postgres WITH PASSWORD 'postgres';`
	1. 建立用户repuser: `CREATE ROLE repuser WITH REPLICATION LOGIN ENCRYPTED PASSWORD '123456';`
	1. 在从机查看从库是否可以访问主节点: `psql -h 192.168.1.11 -U postgres`
1. 编辑pg_hba.conf, 在文件最前面添加
	```
	host replication repuser all md5 # all 表示不限制ip
	host all all all md5 #
	```

1. 创建pg_archive

	```bash
	mkdir -p /var/lib/pgsql/data/pg_archive
	chown postgres:postgres /var/lib/pgsql/data/pg_archive
	```

1. 编辑postgresql.conf

	备份postgresql.conf为postgresql.conf.origin, 再修改postgresql.conf, 再将postgresql.conf备份为postgresql.conf.bak:
	```bash
	listen_addresses = '*'
	wal_level = replica
	synchronous_commit = on
	synchronous_standby_names = '12' # `-.`均是非法字符. 使用异步复制时该字段为空
	archive_mode = on
	archive_command = 'test ! -f /var/lib/pgsql/data/pg_archive/%f && cp %p /var/lib/pgsql/data/pg_archive/%f'
	max_wal_senders = 10
	wal_keep_segments = 64
	wal_sender_timeout = 60s
	max_connections = 1000 # 主备最好相等
	hot_standby = on
	wal_log_hints = on
	full_page_writes = on
	# --- for standby
	max_standby_streaming_delay = 30s
	wal_receiver_status_interval = 10s
	hot_standby_feedback = on
	```
1. 创建recovery.conf

	```bash
	# vim /var/lib/pgsql/data/recovery.conf.sample
	standby_mode = on # 表明是从库
	primary_conninfo = 'host=xxx port=5432 user=repuser password=123456 application_name=11' # 未设置synchronous_standby_names时, pg会忽略application_name
	recovery_target_timeline = 'latest'
	#trigger_file = '/tmp/pg.trigger.5432'
	# chown postgres:postgres /var/lib/pgsql/data/recovery.conf.sample
	```

	> recovery.conf.sample来自/usr/share/pgsql/recovery.conf.sample

	> host是对端ip

	> 指定触发文件，文件存在时，将触发从库提升为主库，前提是必须设置”standby_mode = on”；如果不设置此参数，也可采用”pg_ctl promote“触发从库切换成主库.

	> recovery.conf没有配置明文密码时, 是从`~/.pgpass`获取的

	从库follow新主库时, 修改recovery.conf里面的host, 然后重启pg即可; 主库降为备库时, 修改recovery.conf后, 还需要用到pg_rewind, 比如`pg_rewind --target-pgdata=/var/lib/pgsql/data --source-server='host=192.168.1.12 port=5432 user=postgres password=postgres' -P`, 如果pg_rewind失败就需要重新做备库了.
1. 重启主库
1. 备库同步主库数据

	1. 从库安装完成后, 无需初始化; 若已初始化, 则清空其/var/lib/pgsql/data目录
	1. `pg_basebackup -h 192.168.1.11 -p 5432 -U repuser -F p -Xs -P -D /var/lib/pgsql/data/ && chown postgres:postgres -R /var/lib/pgsql/data`, 这步需要输入repuser的密码

		pg_basebackup参数说明:
		1. `-F` : 指定了输出的格式，支持p（原样输出）或者t（tar格式输出）
		1. -`X` : 用于在备份中包括所需的预写日志文件（WAL文件）. 取值为f/fetch 时与参数-x 含义一样. 取s/stream时表示备份开始后，启动另一个流复制进程接收WAL日志, 此方式需要主库max_wal_senders参数大于或等于2

			- f : 指wal日志在基础备份完成后被传送到备节点, 这时主库上的wal_keep_segments需要设得较大, 以免备份过程中产生的wal还没发送到备库就已被覆盖, 如果出现这种情况创建基准备份会失败. 该模式下主库会启动一个基准备份wal发送进程
			- s : 主库上除了启动原红一个基准备份wal发送进程外还会额外启动一个wal发送进程用于发送主库产生的wal增量日志流, 这种方式避免了f方式过程中主库wal被覆盖的情况, 生产模式推荐这种方式, 特别是比较繁忙的库或大库.
		1. `-P` :表示允许在备份过程中实时的打印备份的进度(估算值)
		1. `-R` : 表示会在备份结束后自动生成recovery.conf文件，这样也就避免了手动创建, 通常是使用手动创建的
		1. `-D` : 指定把备份写到哪个目录，前提是: **做基础备份之前从库的数据目录需要手动清空**
		1. `-l` : 表示指定一个备份的标识
		1. `-v` : 输出详细log
		1. `-C` : 在开始备份之前，允许创建由`-S`选项命名的复制插槽, 比如`pg_basebackup -h 10.20.20.1 -D /var/lib/pgsql/12/data -U replicator -P -v  -R -X stream -C -S pgstandby1`
		1. `-S` : 指定复制插槽名称

		> pg_basebackup能自动处理表空间问题.

1. 修改配置, 并启动备库即可完成一主一从的复制配置
	1. 将postgresql.conf.bak复制为postgresql.conf:

		如果使用了synchronous_standby_names, 那么将其设置为'11'
	1. 将recovery.conf.sample中的host换成`192.168.1.11`, 并将其复制为recovery.conf

		如果使用了synchronous_standby_names, 那么将application_name设置为'12'

测试:
1. 在主库查看同步节点并写入测试数据: 

	```psql
	select * from pg_stat_replication; # 查看同步节点
	create table test(id int primary key,name varchar(20),salary real);
	insert into test values(10,'i love you',10000.00);
	insert into test values(2,'li si',12000.00);
	```

	启用synchronous_standby_names做同步复制时, test2关闭后, test1上执行插入会卡住, 直到test2 pg启动.

	pg_stat_replication:
	- pid: wal发送进程的pid
	- state:
		- startup : wal进程在启动中
		- catchup : 备库正在追赶主库
		- streaming: 备库已追上主库, 并且主库向备库发送wal日志流
		- backup: 通过pg_basebackup正在进行备份
		- stopping: wal发送进程正在关闭
		- sent_lsn: wal发送进程最近发送的wal日志位置
		- write_lsn: 备库最近写入wal日志位置, 这时wal日志流还在os缓存中, 还没有写入备库wal日志文件
		- flush_lsn: 备库最近写入wal日志位置, 这时wal日志流已写入备库wal日志文件
		- replay_lsn: 备库最近应用的wal日志位置
		- write_lag: 主库上wal日志落盘后等待备库接收wal日志(这时wal日志流还在os缓存中, 还没有写入备库wal日志文件)并返回确认信息的时间
		- flush_lag: 主库上wal日志落盘后等待备库接收wal日志(这时wal日志流已写入备库wal日志文件, 但还没有应用wal日志)并返回确认信息的时间
		- replay_lag: 主库上wal日志落盘后等待备库接收wal日志(这时wal日志流已写入备库wal日志文件, 并已应用wal日志)并返回确认信息的时间
		- sync_priority: 基于优先级的模式中备库被选择成为同步的优先级, 对基于quorum的宣讲模式此字段则无影响
	- sync_state:
		- async : 备库使用异步复制
		- sync : 备库使用同步复制
		- potential : 当前为异步复制, 如果当前的同步备库宕机后, 异步备库可升级为同步备库
		- quorum : 备库为quorum standbys的候选

	write_lag, flush_lag, replay_lag是pg 10新增, 用于衡量主备延迟的重要指标.

	查看从服务(WAL接收器进程)状态: `psql -c "\x" -c "SELECT * FROM pg_stat_wal_receiver;"`

	pg_stat_wal_receiver:
	- pid : wal接收进程的pid
	- status : wal接收进程的状态
	- receive_start_lsn: wal接收进程启动后使用的第一个wal日志位置
	- receive_start_tli: wal的timeline
	- received_lsn: 最近接收并写入wal日志文件的wal位置
	- received_tli： 最近接收并写入wal日志文件的wal timeline
	- last_msg_send_time: 备库接收到发送进程最后一个消息后,向主库发回确认消息的发送时间
	- last_msg_receipt_time: 备库接收到发送进程最后一个消息的接收时间
	- latest_end_lsn:
	- latest_end_time: 
	- conninfo: wal接收进程使用的连接串

	监控流复制的相关函数:
	- pg_last_wal_receive_lsn() : 备库最近接收的wal日志位置
	- pg_last_wal_replay_lsn() : 备库最近应用的wal日志位置
	- pg_last_xact_replay_timestamp(): 备库最近事务的应用时间
	- pg_current_wal_lsn(): 主库wal当前写入位置
	- pg_wal_lsn_diff('3/940001B0','3/94000A0') : 计算主备的wal日志位置的偏移量
1. 主备切换

	1. 在test1执行`systemctl stop postgresql`
	1. 用postgres用户在test2执行`pg_ctl promote`, 并插入数据
	1. 在test1, 将recovery.conf.sample复制为recovery.conf, 启动pg, 能看到上一步插入的新数据

1. 再次切换

	1. 在test2执行`systemctl stop postgresql`
	1. 用postgres用户在test1执行`pg_ctl promote`, 并插入数据
	1. 在test2, 执行`pg_rewind --target-pgdata=/var/lib/pgsql/data --source-server='host=192.168.1.11 port=5432 user=postgres password=postgres' -P` , 修正postgresql.conf和recovery.conf.sample, 再将recovery.conf.sample还原为recovery.conf并启动pg, 就能看到上一步插入的新数据了

		成功执行pg_rewind后, 它会原样将对端文件原样拷贝过来(排除了postmaster.pid), 因此可能本地的recovery.conf, recovery.done会被删除. 因此需要修正postgresql.conf和recovery.conf.sample

### pg 12
ref:
- [PostgreSQL 12主从复制](https://www.cxymm.net/article/c410425783/120665685)

#### 修改监听端口
```psql
su - postgres
psql -c "ALTER SYSTEM SET listen_addresses TO '*';"/

ls -l /var/lib/pgsql/12/data/  # 会生成一个postgresql.auto.conf文件
```

ALTER SYSTEM SET会将配置保存在一个postgresql.conf.auto中，与postgresql.conf并存，系统会优先使用.auto配置


## FAQ
### ha切换脚本
ref:
- [postgresql 主备及切换-恢复方案](https://www.yisu.com/zixun/34836.html)

```bash
#!/bin/bash
# from https://www.yisu.com/zixun/34836.html
PRIMARY_IP="192.168.10.2"
STANDBY_IP="192.168.10.3"
PGDATA="/DATA/postgresql/data"
SYS_USER="root"
PG_USER="postgresql"
PGPREFIX="/opt/pgsql"

pg_status()
{
        ssh ${SYS_USER}@$1 /
        "su - ${PG_USER} -c '${PGPREFIX}/bin/pg_controldata -D ${PGDATA} /
        | grep cluster' | awk -F : '{print \$2}' | sed 's/^[ \t]*\|[ \t]*$//'"
}

# recover to primary
recovery_primary()
{
        ssh ${SYS_USER}@$1 /
        "su - ${PG_USER} -c '${PGPREFIX}/bin/pg_ctl promote -D ${PGDATA}'"
}

# primary to recovery
primary_recovery()
{
        ssh ${SYS_USER}@$1 /
        "su - ${PG_USER} -c 'cd ${PGDATA} && mv recovery.done recovery.conf'"
}

send_mail()
{
        echo "send SNS"
}

case "`pg_status ${PRIMARY_IP}`" in
        "shut down")
                case "`pg_status ${STANDBY_IP}`" in
                        "in archive recovery")
                                primary_recovery ${PRIMARY_IP}
                                recovery_primary ${STANDBY_IP}
                                ;;
                        "shut down in recovery"|"in production")
                                send_mail
                                ;;
                esac
                ;;
        "in production")
                case "`pg_status ${STANDBY_IP}`" in
                        "shut down in recovery"|"shut down"|"in production")
                                send_mail
                                ;;
                esac
                echo "primary"
                ;;
        "in archive recovery")
                case "`pg_status ${STANDBY_IP}`" in
                        "shut down")
                                primary_recovery ${STANDBY_IP}
                                recovery_primary ${PRIMARY_IP}
                                ;;
                        "shut down in recovery"|"in archive recovery")
                                send_mail
                                ;;
                esac
                echo "recovery"
                ;;
        "shut down in recovery")
                case "`pg_status ${STANDBY_IP}`" in
                        "shut down in recovery"|"shut down"|"in archive recovery")
                                send_mail
                                ;;
                esac
                echo "recovery down"
                ;;
esac
```

## FAQ
### pg_rewind报`target server needs to use either data checksums or "wal_log_hints = on"`
在initdb时启用checksum或配置`wal_log_hints = on`

已验证在pg_basebackup前配置了`wal_log_hints = on`是可行的.

### pg_rewind报`fatal: target server must be shut down cleanly`
有两种情况:
1. 之前pg未成功关闭, 先尝试启动pg, 再关闭, 并重新执行pg_rewind命令
1. 之前已成功执行pg_rewind, 再次执行就报该错了

### 将原先主库变成备库报`...could not start WAL streaming: ERROR:  requested starting point 3/33000000 on timeline 1 is not in this server's history\n...new timeline 2 forked off current database system timeline 1 before current recovery point...`
原因是：原备用服务器在能够赶上主服务器之前已经升级timeline, 导致现在主服务器不能充当备用服务器的角色

先执行pg_rewind再启动pg

### 执行pg_rewind报"source and target cluster are on the same timeline\n no rewind required"
设置recovery.conf后直接重启即可

### pg启动后报`highest timeline 3 of the primary is behind recovery timeline 4`
在主库test1不停机的情况下, 对备库test2执行`pg_ctl promote`, 此后再在备库上执行插入, 再停止备库pg, 重建recovery.conf并启动pg, 发现备库日志报该错误. 此时主库的pg_stat_replication返回信息中没有该备库.

执行`pg_rewind --target-pgdata=/var/lib/pgsql/data --source-server='host=192.168.1.11 port=5432 user=postgres password=postgres' -P`, 并修改postgresql.conf和recovery.conf后启动pg, 发现主库pg_stat_replication返回信息中有该备库了, 且test2启动中有报`consistent recovery state reached at 0/3037EE8\ninvalid record length at 0/3037EE8: wanted 24, got 0`, 该错的意思是: [standby尝试读取并重放存放在standby的WAL文件, 然后它发现了无效的WAL记录, 即它不再能够在本地读取WAL记录了](https://www.postgresql.org/message-id/CAHGQGwFvv0pxaf_iZ1FU1H%3Dd%3DexhPUoM0ss-9GkDWRP%3DFureMg%40mail.gmail.com), 但它会从test1重新同步wal, 因此在test1插入, test2也能显示出来, 忽略该提示日志即可.

### 如何处理主库wal被覆盖导致备库不可用
1. 调大wal_keep_segments, 但也要确保pg_wal不会撑爆存储
1. 在主库上开启归档, 如果没有足够空间时, 至少在备库听见维护时临时开启主库归档. 这样, 当备库启动时, 如果所需的wal已被覆盖, 至少可从归档中获取所需wal, 比重做备库省时省力
1. 主库设置复制槽(replication slots), 通过该特性, 主库将知道备库的复制状态, 即使备库宕机也一样. 因此主库不会删除备库还未接收的wal. 如果备库停机时间长, 主库的pg_wal可能撑爆存储. 所有**如果设置了复制槽, 建议将pg_wal单独放在大容易磁盘上**.

	主库查看复制槽状态`select * from pg_replication_slots;`:
	- slot_name: 复制槽名称
	- plugin: 如果是物理复制, 显示为空
	- slot_type: physical/logical
	- database: 如果是物理复制, 显示为空; 逻辑复制槽显示dbname
	- active: 是否使用中, t是使用中
	- active_pid: 使用复制槽session的pid, 即主库wal发送进程的pid
	- xmin: db保留的最老事务

	配置方法:
	1. 主库创建复制槽: `select * from pg_create_physical_replication_slot('phy_slot1');`
	1. 备库的recovery.conf需要添加`primary_slot_name='phy_slot1'`