# mysql
参考:
- [MySQL5.7的官方文档(中文版)](https://www.docs4dev.com/docs/zh/mysql/5.7/reference/manual-info.html)

MySQL 的逻辑架构图:
![](/misc/img/mysql/dfff6efbab0d51a715a36f867daeacf8.png)
mysql的查询流程:
![](/misc/img/mysql/20c7d975ab3bd4617c6b4f46b720fea3.png)
[MySQL架构总览](https://wwxiong.com/2017/06/sql-optimization-up/):
![](/misc/img/mysql/mysql.png)

MySQL可以分为Server层和存储引擎层两部分:
1. Server层包括连接器、查询缓存、分析器、优化器、执行器等，涵盖MySQL的大多数核心服务功能，以及所有的内置函数（如日期、时间、数学和加密函数等），所有跨存储引擎的功能都在这一层实现，比如存储过程、触发器、视图等.
1. 存储引擎层负责数据的存储和提取. 其架构模式是插件式的，支持InnoDB(**推荐**)、~~MyISAM(不推荐)~~、Memory等多个存储引擎. 现在最常用的存储引擎是InnoDB，它从MySQL 5.5.5版本开始成为了默认存储引擎

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
