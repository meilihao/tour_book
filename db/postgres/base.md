## 体系结构
ref:
- [史上最详细的PostgreSQL体系架构介绍](https://cloud.tencent.com/developer/article/1873795)

### 存储结构
PG数据存储结构分为：逻辑存储结构和物理存储存储. 其中：逻辑存储结构是内部的组织和管理数据的方式；物理存储结构是操作系统中组织和管理数据的方式

#### 逻辑存储结构
所有数据库对象都有各自的oid(object identifiers),oid是一个无符号的四字节整数，相关对象的oid都存放在相关的system catalog表中，比如数据库的oid和表的oid分别存放在pg_database,pg_class表中.  

在逻辑存储结构的相关术语：
- 数据库集群-Database cluster
  也叫数据库集簇。它是指有单个PostgreSQL服务器实例管理的数据库集合，组成数据库集群的这些数据库使用相同的全局配置文件和监听端口、共用进程和内存结构。

  一个DataBase Cluster可以包括：多个DataBase、多个User、以及Database中的所有对象。

- 数据库-Database

  在PostgreSQL中，数据库本身也是数据库对象，并且在逻辑上彼此分离.

- 数据库对象-Database object

  如：表、视图、索引、序列、函数等等。在PostgreSQL中的所有数据库对象都由各自的对象标识符（OID）进行内部的管理。例如，数据库的OID存储在pg_database系统表中，可以通过`select oid,datname from pg_database;`进行查询.

  而数据库中的表、索引、序列等数据库对象的OID则存在了pg_class系统表中，例如可以通过`select oid,relname,relkind,relfilenode from pg_class where relname ='testtable1';`查询前面创建的testtable1表的OID.

- 表空间-tablespace

  数据库在逻辑上分成多个存储单元，称作表空间. 表空间用作把逻辑上相关的结构放在一起. 数据库逻辑上是由一个或多个表空间组成. 初始化的时候，会自动创建pg_default和pg_global两个表空间.

  在PostgreSQL中, 表空间是一个目录, 里面存储的是它所包含的数据库的各种物理文件.
- 模式-Schema

  当创建一个数据库时，会为其创建一个名为public的默认Schema. Schema是数据库中的命名空间，在数据库中创建的所有对象都是在Schema中创建, 一个用户可以从同一个客户端连接中访问不同的Schema. 而不同的Schema中可以有多个同名的Table、Index、View、Sequence、Function等等数据库对象, 可以通过`\dn`来查看当前数据库的Schema.

  > Schema将数据库对象组织成逻辑组，从而提高可管理性、安全性和性能

  **在PostgreSQL中, 数据库的创建是通过克隆数据库模板来实现的, 这与SQL SERVER是同样的机制**. 由于CREATE DATABASE dbname并没有指明数据库模板, 所以系统将默认克隆template1数据库, 得到新的数据库dbname.

  ```sql
  CREATE DATABASE dbname TEMPLATE template1 TABLESPACE tablespacename;
  ALTER DATABASE dbname OWNER TO custom;
  ```

- 段-segment

  一个段是分配给一个逻辑结构（一个表、一个索引或其他对象）的一组区，是数据库对象使用的空间的集合；段可以有表段、索引段、回滚段、临时段和高速缓存段等。

- 区-extent

  区是数据库存储空间分配的一个逻辑单位，它由连续数据块所组成。第一个段是由一个或多个盘区组成。当一段中间所有空间已完全使用，PostgreSQL为该段分配一个新的范围。

- 块-block（Page）
  数据块是PostgreSQL 管理数据文件中存储空间的单位，为数据库使用的I/O的最小单位，是最小的逻辑部件。默认值8K。

### 物理存储结构
在执行initdb的时候会初始化一个目录，通常我们都会在系统配置相关的环境变量$PGDATA来表示，初始化完成后，会再这个目录生成相关的子目录以及一些文件。在postgresql中，表空间的概念并不同于其他关系型数据库，这里一个Tablespace对应的都是一个目录.

而PostgreSQL的物理存储结构主要是指硬盘上存储的文件，包括：数据文件、日志文件、参数文件、控制文件、redo日志（WAL）。下面分别进行介绍:

- 数据文件（表文件）

  顾名思义，数据文件用于存储数据。文件名以OID命名，对于超出1G的表数据文件，PostgreSQL会自动将其拆分为多个文件来存储，而拆分的文件名将由pg_class中的relfilenode字段来决定, 例如`select oid,relname,relkind,relfilenode from pg_class where relname ='testtable1';`

  在PostgreSQL中，将保存在磁盘中的块（Block）称为Page。数据的读写是以Page为最小单位，每个Page默认的大小是8K。在编译PostgreSQL时指定BLCKSZ大小将决定Page的大小。每个表文件由逗哥BLCKSZ字节大小的Page组成。在分析型数据库中，适当增加BLCKSZ大小可以小幅度提升数据库的性能。
- 日志文件

  PostgreSQL日志文件的类型，分为以下几种:
  - 运行日志（pg_log）

    默认没有开启，开启后会自动生成
  - 重做日志（pg_xlog）

    pg_xlog 这个日志是记录的Postgresql的WAL信息，默认存储在目录$PGDATA/pg\_wal/，是一些事务日志信息(transaction log)。默认单个大小是16M，源码安装的时候可以更改其大小（./configure --with-wal-segsize=target_value 参数，即可设置）这些日志会在定时回滚恢复(PITR)， 流复制(Replication Stream)以及归档时能被用到，这些日志是非常重要的，记录着数据库发生的各种事务信息，不得随意删除或者移动这类日志文件，不然你的数据库会有无法恢复的风险.
  -  事务日志（pg_xact） 

    pg_xact是事务提交日志，记录了事务的元数据。默认开启。内容一般不能直接读。默认存储在目录$PGDATA/pg_xact.
  - 服务器日志

    如果用pg_ctl启动的时候没有指定-l参数来指定服务器日志，错误可能会输出到cmd前台。服务器日志记录了数据库的重要信息.
- 参数文件

  主要包括postgresql.conf、pg_hba.conf和pg_ident.conf这三个参数文件:
  - postgresql.conf

    PostgreSQL的主要参数文件，有很详细的说明和注释，和Oracle的pfile，MySQL的my.cnf类似。默认在$PGDATA下。很多参数修改后都需要重启。9.6之后支持了alter system来修改，修改后的会存在$PGDATA/postgresql.auto.conf下，可以reload或者 restart来使之生效.
  - pg_hba.conf

    这个是黑白名单的设置
  - pg_ident.conf

    pg_ident.conf是用户映射配置文件，用来配置哪些操作系统用户可以映射为数据库用户。结合pg_hba.conf中，method为ident可以用特定的操作系统用户和指定的数据库用户登录数据库.
- 控制文件

  控制文件记录了数据库运行的一些信息，比如数据库id，是否open，wal的位置，checkpoint的位置等等。controlfile是很重要的文件。  

  控制文件的位置：$PGDATA/global/pg_control，可以使用命令bin/pg_controldata查看控制文件的内容
- redo日志（WAL）

  默认保存在$PGDATA/pg_wal目录下. 文件名称为16进制的24个字符组成，每8个字符一组，每组的意义是：时间线    逻辑ID 物理ID. 可通过`select pg_switch_wal();`进行WAL的手动切换.

### 进程结构
通过`ps -ef | grep postgres`可列出所有的PostgreSQL的进程:
-  Postmaster进程

  主进程Postmaster是整个数据库实例的总控制进程，负责启动和关闭数据库实例。用户可以运行postmaster，postgres命令加上合适的参数启动数据库。实际上，postmaster命令是一个指向postgres的链接.

  更多时候我们使用pg_ctl启动数据库，pg_ctl也是通过运行postgres来启动数据库，它只是做了一些包装，让我们更容易启动数据库，所以，主进程Postmaster实际是第一个postgres进程，此进程会fork一些与数据库实例相关的辅助子进程，并管理它们.

  当用户与PostgreSQL数据库建立连接时，实际上是先与Postmaster进程建立连接。此时，客户端程序会发出身份证验证的消息给Postmaster进程，Postmaster主进程根据消息中的信息进行客户端身份验证。如果验证通过，它会fork一个子进程postgres为这个连接服务，fork出来的进程被称为服务进程，查询pg_stat_activity表可以看到的pid，就是这些服务进程的pid.
- SysLogger进程

  在postgresql.conf里启用    运行日志（pg_log）后，会有SysLogger进程。SysLogger会在日志文件达到指定的大小时关闭当前日志文件，产生新的日志文件.
- BgWriter后台写进程

  BgWriter是PostgreSQL中在后台将脏页写出到磁盘的辅助进程，引入该进程主要为达到如下两个目的：
  1. 首先，数据库在进行查询处理时若发现要读取的数据不在缓冲区中时要先从磁盘中读入要读取的数据所在的页面，此时如果缓冲区已满，则需要先选择部分缓冲区中的页面替换出去。如果被替换的页面没有被修改过，那么可以直接丢弃；但如果要被替换的页已被修改，则必需先将这页写出到磁盘中后才能替换，这样数据库的查询处理就会被阻塞。通过使用BgWriter定期写出缓冲区中的部分脏页到磁盘中，为缓冲区腾出空间，就可以降低查询处理被阻塞的可能性。
  1. 其次，PostgreSQL在定期作检查点时需要把所有脏页写出到磁盘，通过BgWriter预先写出一些脏页，可以减少设置检查点时要进行的IO操作，使系统的IO负载趋向平稳。通过BgWriter对共享缓冲区写操作的统一管理，避免了其他服务进程在需要读入新的页面到共享缓冲区时，不得不将之前修改过的页面写出到磁盘的操作。
- WalWriter预写日志写进程  

  该进程用于保存WAL预写日志。预写式日志WAL（Write Ahead Log，也称为Xlog）的中心思想是对数据文件的修改必须是只能发生在这些修改已经记录到日志之后，也就是先写日志后写数据。如果遵循这个过程，那么就不需要在每次事务提交的时候都把数据块刷回到磁盘，这一点与Oracle数据库是完全一致的.   
- PgArch归档进程

  从PostgreSQL 8.x开始，有了PITR（Point-In-Time-Recovery）技术，该技术支持将数据库恢复到其运行历史中任意一个有记录的时间点；PITR的另一个重要的基础就是对WAL文件的归档功能。PgArch辅助进程的目标就是对WAL日志在磁盘上的存储形式进行归档备份。但在默认情况下，PostgreSQL是非归档模式，因此看不到PgArch进程。PgArch进程通过postgresql.conf文件中的参数进行配置.
- AutoVacuum自动清理进程

  在PG数据库中，对数据进行UPDATE或者DELETE操作后，数据库不会立即删除旧版本的数据，而是标记为删除状态。这是因为PG数据库具有多版本的机制，如果这些旧版本的数据正在被另外的事务打开，那么暂时保留他们是很有必要的。当事务提交后，旧版本的数据已经没有价值了，数据库需要清理垃圾数据腾出空间，而清理工作就是AutoVacuum进程进行的。postgresql.conf文件中与AutoVacuum进程相关的参数.
- PgStat统计信息收集进程
  
  PgStat进程是PostgreSQL数据库的统计信息收集器，用来收集数据库运行期间的统计信息，如表的增删改次数，数据块的个数，索引的变化等等。收集统计信息主要是为了让优化器做出正确的判断，选择最佳的执行计划。postgresql.conf文件中与PgStat进程相关的参数.
- CheckPoint检查点进程

### 内存结构
PostgreSQL的内存结构，分为：本地内存和共享内存

## 数据类型

### 整型
| 类型名称 | 范围 |存储需求|
|--------|--------|----------|
|  smallint / int2      |        |2B|
|  int / int4     |        |4B|
|  bigint / int8   |        |8B|

### 浮点类型
| 类型名称 | 说明 |存储需求|
|--------|--------|----------|
|  real      | 6 位十进制数字精度       |4B|
|  double precision	      |     15 位十进制数字精度  |8B|

>infinity,正无穷打;-infinity,负无穷打;NaN,不是一个数字.

### 任意精度类型
| 类型名称 | 说明 |
|--------|--------|
|numeric|numeric(M,N).M称为精度,即总的位数;N是标度,即小数的位数.如果用户数据的精度超出指定精度,则会四舍五入|

### 日期和时间类型
| 类型名称 | 含义 |存储需求|
|--------|--------|----------|
|  time      |  只用于一日内时间,格式:HH:MM:SS      |8B|
|  date      |  只用于日期,格式:YYYY-MM-DD      |4B|
|  timestamp      |  日期和时间,格式:YYYY-MM-DD HH:MM:SS|8B|

> 1. 不合法的日期和时间类型存入时会被替换为相应零值.
> 2. current_time和now()均表示当前系统时间,在存入数据库时,只保留符合相应类型的部分

### 字符串类型
| 类型名称 | 说明 |
|--------|--------|
|  char      | 固定长度非二进制字符串,不足补空白, **推荐使用varchar**      |
|  varchar	      |    变长非二进制字符串,有长度限制  |
|  text	      |     变长非二进制字符串,无长度限制  |

### 二进制类型
| 类型名称 | 说明 |存储需求|
|--------|--------|----------|
|  bytea      | 变长的二进制字符串       |1或4字节加上实际的二进制字符串|

### 布尔类型
| 类型名称 | 说明 |存储需求|
|--------|--------|----------|
|  boolean      | true/false       |1B|

### 数组类型
允许将字段定义为定长或变长的一维或多维数组.不过目前pg并不强制数组的长度,所以声明长度和不声明长度是一样的.

## 常见运算符
运算符优先级: [官方文档](http://www.postgres.cn/docs/10/sql-syntax-lexical.html#SQL-PRECEDENCE)
### 比较运算符
规则:
1. 有一个或两个参数为NULL,结果为空
2. 两个参数均为字符串,按照字符串进行比较.
3. 两个参数均为数值,按照数值比较.
4. 一个字符串和一个数值时,将字符串转换为数值(无法转换则直接报错)再比较

其他规则:
- `x BETWEEN min AND max`==`min<=x<=max`
- LEAST运算符,返回在两个及以上参数中的最小值,**会忽略NULL**.
- GREATEST运算符,返回在两个及以上参数中的最大值,**会忽略NULL**.
- IN运算符,判断操作数是否在IN列表中,返回值为t或f.当左侧表达式为NULL或找不到匹配项且列表中有NULL时,返回空值.
- LIKE运算符,匹配字符串,返回值为t或f.左侧表达式或右侧匹配条件中任何一方为NULL,返回空值.LIKE的通配符:`%`匹配任意数目(包括零)的字符;`_`只能匹配一个字符.

### 逻辑运算符
- NOT运算符,操作数必须是布尔变量(即1,0,t,f,y,n,true,false,yes,no的字符串且大小写不敏感).当操作数为NULL时,返回空值.
- AND运算符,两个操作数都为true且都不为NULL时,返回t;至少有一个操作数为false时,返回f;其他情况都返回空值.
- OR运算符,两个操作数都不为NULL且任一操作数为true时,返回t,否则返回f;一个操作数为NLL且另一个为true时,返回t,否则返回空值;两个均为NULL时返回空值.

### 位运算符
<table><thead><tr><th>操作符</th><th>描述</th><th>例子</th><th>结果</th></tr></thead><tbody><tr><td>||</td><td>连接</td><td>B'10001' || B'011'</td><td>10001011</td></tr><tr><td>&amp;</td><td>位与</td><td>B'10001' &amp; B'01101'</td><td>00001</td></tr><tr><td>|</td><td>位或</td><td>B'10001' | B'01101'</td><td>11101</td></tr><tr><td>#</td><td>位异或</td><td>B'10001' # B'01101'</td><td>11100</td></tr><tr><td>`~`</td><td>位非</td><td>`~ B'10001'`</td><td>01110</td></tr><tr><td>&lt;&lt;</td><td>位左移</td><td>B'10001' &lt;&lt; 3</td><td>01000</td></tr><tr><td>&gt;&gt;</td><td>位右移</td><td>B'10001' &gt;&gt; 2</td><td>00100</td></tr></tbody></table>

ps: `&, |,#`的位串操作数必须等长

### [运算符优先级](http://www.postgres.cn/docs/11/sql-syntax-lexical.html#SQL-PRECEDENCE)


## sql
psql:
```sql
-- 查询sql命令的解释
\h alter table
-- 打开编辑器
\e
```

```sql
-- 修改数据库名称
alter database mytest rename to mytest1;
-- 修改数据库所有人
alter database mytest owner to user1;
-- 修改字段类型
alter table tb alter column name type varchar(10); -- 注意有`type`
-- 修改字段名称
alter table tb rename xxx to xxx2;
-- 添加字段
alter table tb add column a int8;
-- 删除外键
alter table tb drop constraint  ${外键名};
-- 删除表
drop table if exists td;
-- 插入时间, t2 is time
insert into users (t2) values (CURRENT_TIME);
```

## pg函数

### 数学函数
- [数学函数](http://www.postgres.cn/docs/10/functions-math.html)
  - random() : 范围 0.0 <= x < 1.0 中的随机值
  - trunc（value,precision）: 按精度(precision)截取某个数字,**不进行舍入操作**
  - round（value,precision）: 根据给定的精度(precision)输入数值, **会进行舍入操作**, 如果precision为负数, 则将保留value值到小数点左边precision位
  - ceil (value) : 产生大于或等于指定值（value）的最小整数
  - floor（value）: 与 ceil（）相反，产生小于或等于指定值（value）的最小整数
  - sign(value) : 取参数的符号（-1, 0, +1）
  - abs(x) : 绝对值
  - mod(x,y) : 返回x被y除后的余数

- [字符串函数](http://www.postgres.cn/docs/10/functions-string.html)
  - char_length(str) : 一个多字节字符算一个字符．
  - length(str) : 返回字符串的字节长度, 比如使用utf8编码时，一个汉字3字节．
  - concat(s1,s2,...) : 返回拼接后的字符串(无拼接字符)
  - concat(x,s1,s2,...) : 返回使用x拼接后的字符串
  - left(s,n) : 返回字符串s的最左边n个字符
  - right(s,n) : 返回字符串s的最右边n个字符
  - substring(s, n, len) : 从字符串s返回一个长度为len的子字符串, 起始位置是n. 如果n<0,表示倒数第n个字符.

- [时间/日期函数](http://www.postgres.cn/docs/10/functions-datetime.html)
  - current_date : 获取系统当前日期
  - current_time : 获取系统当前具体时间(时分秒)
  - current_timestamp : 获取系统当前时间(日期+具体时间)
  - extract(type from date)从当前时间中提前部分
- [条件判断函数](http://www.postgres.cn/docs/10/functions-conditional.html)
  - case [value] when v1 then r1 [when v2 then r2] [else rn] end : if...else....没有value时, vn应为bool表达式.
- [系统信息函数](http://www.postgres.cn/docs/10/functions-info.html)
  - version()
  - user/current_user
- [加密函数](http://www.postgres.cn/docs/10/functions-binarystring.html)
  - md5('xxx')
  - decode/encode(str, encodeType) : 使用encodeType(base64,hex)编/解码str

## 数据查询/select
- 除非必要应尽量避免使用`select *`,否则会因获取不需要的列数据而降低性能．
- 有时IN操作符可以实现OR操作符的效果，此时推荐使用IN,因为其语法更简明且执行速度更快，而且还支持更复杂的嵌套查询．
- group by通常和集合函数一起使用，例如，max(),min(),count(),sum(),avg().
- having和where均用于过滤数据．having用在数据分组即group by后来过滤分组;where用于选择特定的记录．
-`limit x [offset y]`:［从表的第ｙ＋1条记录开始］选取ｘ条记录．
- `count(column)`会忽略column字段值为空值(NULL)的行;`count(*)`则所有记录均不忽略．
- `sum()`忽略列值为NULL的行．
- `max()`不仅适用于数值类型，也可用于字符类型，日期值．
- `order by`与`limit`联用时, `limit`必须在后面
- `distinct`是应用于给出的所有列

### [连接查询](http://www.postgres.cn/docs/10/queries-table-expressions.html)
- 内连接:

  格式：`select ... from t1 inner join t2 on t1.xxx=t2.xxx`
  内连接也可用where子句来实现，性能可能更高．
  自连接是内连接的特例．
- 外连接:

  `left join(返回包括左表中的所有记录和右表中连接字段相等的记录)`和`right join(返回包括右表中的所有记录和左表中连接字段相等的记录)`
- 全连接

  `FULL OUTER JOIN`,显示符合条件的数据行，同时显示左右不符合条件的数据行，相应的左右两边显示NULL，即显示左连接、右连接和内连接的并集.
- 交叉连接

  `T1 CROSS JOIN T2`,对每个来自T1和T2 的行进行组合（也就是，一个笛卡尔积），连接成的表将包含这样的行： 所有T1里面的字段后面跟着所有T2 里面的字段。如果两表分别有 N 和 M 行，连接成的表将有 N*M 行.`FROM T1 CROSS JOIN T2`等效于`FROM T1,T2`.

### [子查询](http://www.postgres.cn/docs/10/functions-subquery.html)
子查询常用操作符：any/some,all,in,exists.
- any=some

  `expression operator ANY (subquery)`,左边表达式使用operator对子查询结果的每一行进行一次计算和比较， 其结果必须是布尔值.如果至少获得一个真值，则ANY结果为"真"; 如果全部获得假值，则结果是"假"(包括子查询没有返回任何行的情况).

- all

  `expression operator ALL (subquery)`,左边表达式使用operator对子查询结果的每一行进行一次计算和比较， 其结果必须是布尔值。如果全部获得真值，ALL结果为"真" (包括子查询没有返回任何行的情况)。如果至少获得一个假值，则结果是"假"。`NOT IN等效于<> ALL`.

- EXISTS

  `EXISTS (subquery)`,判断subquery是否返回行。如果它至少返回一行，那么EXISTS的结果就为 "真"；如果子查询没有返回任何行，那么EXISTS的结果是"假"。`NOT EXISTS`和`EXISTS`作用相反．

- IN

  `expression IN (subquery)`,左边表达式对子查询结果的每一行进行一次计算和比较。如果找到任何相等的子查询行， 则IN结果为"真"。如果没有找到任何相等行，则结果为"假" (包括子查询没有返回任何行的情况)。`NOT IN`和`IN`作用相反．

### 合并查询
- union [all]: 列数,数据类型都必须相同. all的作用是不删除重复行,也不对结果进行排序

### 正则查询
[模式匹配#POSIX 正则表达式](http://www.postgres.cn/docs/10/functions-matching.html)
- ~ : 匹配正则, 且区分大小写
- ~* : 匹配正则, 不区分大小写
- !~ : 不匹配正则, 且区分大小写
- !~* : 不匹配正则, 不区分大小写

## 索引
索引是对数据库表中一列或多列的值进行排序的一种结构，其包含着对数据表里所有记录的引用指针，是提高数据库性能的常用方法．其由存储引擎实现.

> 在使用分组和排序子句进行查询时使用索引, 也可显著减少查询中分组和排序的时间.

PostgreSQL提供了好几种索引类型：B-tree, Hash, GiST, SP-GiST和GIN．每种索引类型都比较适合某些特定的查询类型，因为它们用了不同的算法, 缺省时， CREATE INDEX命令将创建 B-tree 索引，它适合大多数情况.

索引类型:
- B-tree, 默认

  适合处理那些能够按顺序存储的数据,特别是在一个建立了索引的字段涉及到使用`<,<=,=,>=,>`.

  它适用于范围查询和通用索引.
- Hash(不推荐)

  只能处理简单的等于比较。当一个索引了的列涉及到使用= 操作符进行比较的时候.
  
  不推荐原因：性能比B-tree弱，且不走WAL日志，因此数据库崩溃时需使用REINDEX重建Hash索引

- GiST(广义搜索树)

  支持复杂数据类型和自定义搜索算法. 它对于几何形状、文本搜索和其他复杂数据类型非常有用

  GiST 索引不是单独一种索引类型，而是一种架构，可以在这种架构上实现很多不同的索引策略。 因此，可以使用 GiST 索引的操作符高度依赖于索引策略(操作符类)。

- SP-GiST(空间分区广义搜索树)

  主要针对空间数据和其他分区数据集进行了优化, 对空间索引和层次数据非常有效

  SP-GiST索引类似于GiST索引，提供一个支持不同类型检索的架构。 SP-GiST允许实现许多各种不同的非平衡的基于磁盘的数据结构，例如四叉树，k-d树和基数树(字典树)。

- GIN(广义倒排索引)

  索引设计用于索引数组数据类型和全文搜索，并被认为适合全文搜索和多值字段

  GIN 索引是反转索引，它可以处理包含多个键的值(比如数组)。与 GiST和SP-GiST 类似， GIN 支持用户定义的索引策略，可以使用 GIN 索引的操作符根据索引策略的不同而不同。
- BRIN(块范围索引)

  对于具有顺序排列数据的大型表非常高效, 并被认为是理想的用于具有块级相关性的数据，例如时间序列数据

设计索引的准则:
- 索引并非越多越好,一个表中如果有大量的索引,不仅占用的存储空间将增大,而且会影响Insert,Delete,Update等语句的性能,因为当表中的数据更改时，所有索引都须进行适当的调整和更新。
- 避免对经常更新的表进行过多的索引，并且索引应保持较窄，就是说，列要尽可能少。而对经常用于查询的字段应该创建索引,但要避免添加不必要的字段.
- 数据量小的表最好不要使用索引,由于数据较少,查询花费的时间可能比遍历索引的时间还要短,索引可能不会产生优化效果.
- 在条件表达式中经常用到的、不同值较多的列上建立检索,在不同值少的列上不要建立索引.比如学生表的'性别'字段只有男,女两种,因此就没必要建立索引,此时如果建立索引不但不会提高查询效率,反而会严重降低更新速度.
- 当唯一性是某种数据本身的特性时,指定唯一索引,使用该索引能够确保定义的列的数据完整性,提高查询速度.
- 在频繁进行排序或分组的列上建立索引时,如果待排序的列有多个,可以在这些列上建立组合索引.

### 索引分类
1. 普通索引,最基本的索引,没有唯一性之类的限制,其作用只是加快对数据的访问速度.
2. 唯一索引,和普通索引类似,但索引列的值必须唯一,允许有NULL,如果是组合索引,则列值的组合必须唯一.其作用是减少查询索引列操作的执行时间,特别是大表的时候.
3. 单列索引,在数据表中的某个字段上创建的索引,一个表中可以有多个单列索引.
4. 组合索引,在数据表中的多个字段上创建的索引.

## 视图
视图是一个虚拟表,是从一个或多个表中导出的,它的行为与普通表非常类似,可以帮助用户屏蔽真实表结构变化带来的影响,并提高安全性(添加限定条件,屏蔽特定的行和列).

## 触发器
由业务层处理, 更自由友好.

## 事务
pg的有事务管理器负责事务, 可分为两部分:
1. 锁管理器 : 主要提供在事务的写阶段并发控制所需要的各种锁, 从而保证事务的各种隔离级别
1. 日志管理器 : 主要记录事务执行的状态和数据的变化过程

## 权限
在PostgreSQL 里没有区分用户和角色的概念，"CREATE USER" 为 "CREATE ROLE" 的别名，这两个命令几乎是完全相同的，唯一的区别是"CREATE USER" 命令创建的用户默认带有LOGIN属性，而"CREATE ROLE" 命令创建的用户默认不带LOGIN属性.

## TOAST
当数据大小(比如使用 TEXT 数据类型的列)超过阈值时，由 TOAST 管理.

TOAST通过将大数据拆分为可管理的块来最小化存储大数据的性能影响. 这改善了对大对象的访问时间.

通过将大属性存储在单独的表中，TOAST减少了主表所需的空间，优化了磁盘使用.

使用TOAST是一种有效管理大数据类型而不影响性能的方法.

> 查询pg_toast命名空间以检查TOAST表: `SELECT relname FROM pg_class WHERE relname LIKE 'pg_toast_%'`

## 临时表
临时表用于存储仅在数据库会话或事务期间需要的瞬态数据. 它们对于管理中间结果、执行复杂计算和简化查询逻辑非常有价值.

临时表在会话或事务结束时也会自动删除, 但也支持手动删除它们.

创建临时表: `CREATE TEMP TABLE sales.temp_orders AS SELECT * FROM sales.orders WHERE order_date = '2023-03-22'`

## CTE(公共表表达式)
with 子句也被称为CTE, 用于定义可以在 SELECT 或 DELETE 语句中引用的临时结果集. CTE 特别适合将复杂查询分解为更简单、更易管理的部分.

> CTE 可以链在一起以创建更复杂的查询, 即同时定义多个CTE.

```sql
WITH sales_today AS (
SELECT order_id, customer_id, total_due
FROM sales.orders
WHERE order_date = '2023-03-22'
)

SELECT COUNT(order_id) AS total_orders,
SUM(total_due) AS total_sales
FROM sales_today;
```

## 递归查询
支持使用 WITH RECURSIVE 子句的递归查询，这允许查询引用自己的输出. 递归查询对于查询层次数据结构（例如组织结构图或产品类别）特别有用.

## pg_dump
```
$ pg_dump -U postgres -f /home/chen/test_backup test # 备份数据库test中的所有表
$ pg_dump -U postgres -t t1 [-t tn,...] -f /home/chen/test_backup test # 备份数据库test中的指定表
$ pg_dumpall -U postgres -f /home/chen/db_backup test # 备份所有数据库
$ psql -d test -U postgres -f /home/chen/test_backup # 将备份的数据库还原(还原文件应是create,insert语句的文本文件)
$ pg_restore -d test -U postgres -C /home/chen/test_backup # 将备份的数据库还原, `-C`表示在恢复数据库之前先创建它
```

sql:
```
$ pg_dump -U postgres test > xxx.sql
$ psql -U postgres -d test < db_backup.sql
```

## SQL说明
```sql
CREATE DATABASE mytest
  WITH OWNER = postgres --新数据库的所有者
       ENCODING = 'UTF8' --创建新数据库使用的字符编码
       TABLESPACE = pg_default --和新数据库关联的表空间名字
       LC_COLLATE = 'en_US.UTF-8' --(区域支持)字符串排序规则
       LC_CTYPE = 'en_US.UTF-8' --(区域支持)字符分类,比如哪些字符是字母，哪些是数字，大小写等
       CONNECTION LIMIT = -1; --连接数限制,默认-1表示不限制

COMMENT ON DATABASE mytest --为数据库添加注释
  IS '创建的第一个数据库';

ALTER DATABASE mytest
  RENAME TO mytest2; --数据库重命名
ALTER DATABASE mytest2
  OWNER TO qor; --修改数据库的所有者
```

## explain
explain(仅用于select)用于分析一个语句的执行计划, 即显示语句如何查询表.

格式:
```sql
-- ANALYZE : 显示执行时间
-- VERBOSE : 计划树完整的内部表现形式,而不仅仅是摘要
-- cost : 代表语句的执行成本(即计算开销), 包括语句的花费时间, 扫描的行数等.
EXPLAIN [ ANALYZE ] [ VERBOSE ] 语句
```

> pg预估成本: `cost=0.00..5.04`意味着PostgreSQL希望花费`5.04`的任意计算单位来找到这些值,而`0.00`是该节点起始工作成本(即启动成本).
> 时间成本: `actual time=0.049..0.049`表示此步骤的开始时间是0.049,结束时间0.049,单位为毫秒,因此此实际执行时间是0,实际时间是每次迭代的平均值,可以将值乘以循环次数以获得真实的执行时间.

连接查询更高效: pg不需要在内存中创建临时表来完成查询工作
子查询效率慢: pg需要为内层查询的结果建立临时表以供外层查询语句查询.

> 显示表的统计信息: `ANALYZE [ VERBOSE ] [ table_and_columns [, ...] ]`

[PostgreSQL执行计划的解释](http://blog.csdn.net/ls3648098/article/details/7602136)
执行计划运算类型	操作说明	是否有启动时间
Seq Scan	扫描表	无启动时间
Index Scan	索引扫描	无启动时间
Bitmap Index Scan	索引扫描	有启动时间
Bitmap Heap Scan	索引扫描	有启动时间
Subquery Scan	子查询	无启动时间
Tid Scan	ctid = …条件	无启动时间
Function Scan	函数扫描	无启动时间
Nested Loop	循环结合	无启动时间
Merge Join	合并结合	有启动时间
Hash Join	哈希结合	有启动时间
Sort	排序，ORDER BY操作	有启动时间
Hash	哈希运算	有启动时间
Result	函数扫描，和具体的表无关	无启动时间
Unique	DISTINCT，UNION操作	有启动时间
Limit	LIMIT，OFFSET操作	有启动时间
Aggregate	count, sum,avg, stddev集约函数	有启动时间
Group	GROUP BY分组操作	有启动时间
Append	UNION操作	无启动时间
Materialize	子查询	有启动时间
SetOp	INTERCECT，EXCEPT 有启动时

## 扩展
### 插件
- pglogical: 双向复制

### 特殊表
1. temporary table : 临时表. 会话隔离, 本会话创建的临时表，不能被其他会话看到; 临时表的生命周期最长就是会话生命周期
1. unlogged table : 为临时数据设计的(不写WAL)，写入性能较高，但是当postgresql进程崩溃时会丢失数据

## 表空间
表空间即PostgreSQL存储数据文件的位置, 其中包括数据库对象, 如: 索引、表等.

PostgreSQL使用表空间映射逻辑名称和磁盘物理位置, 默认提供了两个表空间：
- pg_default 表空间存储用户数据, 比如存储系统目录对象、用户表、用户表index、和临时表、临时表index、内部临时表的默认空间, 对应存储目录`$PADATA/base/`
- pg_global 表空间存储全局数据, 比如存放系统字典表pg_database、pg_authid、pg_tablespace等表以及它们的索引, 对应存储目录`$PADATA/global/`

自定义表空间, 是用户创建的表空间, 对应文件系统目录`$PADATA/pg_tblspc/`, 当手动创建表空间时, 该目录下会自动生成一个软链接, 指向表空间设定的路径.

利用表空间可以控制PostgreSQL的磁盘布局，它有两方面的优势：
1. 首先，如果集群中的某个分区超出初始空间, 可以在另一个分区上创建新的表空间并使用, 后期可以重新配置系统
1. 其次，可以使用统计优化数据库性能. 比如可以把频繁访问的索引或表放在高性能的磁盘上, 如固态硬盘; 把归档数据放在较慢的设备上

创建语法:
```sql
CREATE TABLESPACE tablespace_name
OWNER user_name
LOCATION directory_path;

create database testdb tablespace tablespace_name; -- 此时删除tablespace需要先删除testdb
CREATE TABLE foo(i int) TABLESPACE tablespace_name;
CREATE INDEX newtab_val_idx ON newtab (val) TABLESPACE mytbsp;
```
表空间名称不能以pg开头，因为这些名称为系统表空间保留. directory_path是表空间使用空目录的绝对路径，PostgreSQL的用户必须拥有该目录的权限可以进行读写操作.

在Oracle数据库中, 一个表空间只属于一个数据库使用; 而一个数据库可以拥有多个表空间, 属于"一对多"的关系.
在PostgreSQL集群中, 一个表空间可以让多个数据库使用; 而一个数据库可以使用多个表空间, 属于"多对多"的关系

> 表索引不会继承表的表空间.

## config
### core
- shared_buffers: 它决定了分配给共享内存缓冲区中缓存数据的内存量. 通过增加此值，可以通过减少磁盘读取的频率来提高读取性能, 建议设置为系统 RAM 的 25-40%
- work_mem: 指定了用于排序和哈希操作的内存量. 更高的值可以提高查询性能, 特别是对于涉及排序和聚合的复杂操作. 建议根据工作负载要求和可用内存进行调整
- maintenance_work_mem: 为维护任务分配内存，例如清理和索引. 通过增加这个值可以加快维护过程，减少停机时间. 对于维护需求频繁的大型数据库，通常设置得更高
- effective_cache_size: 估计操作系统磁盘缓存的大小，指导查询规划器的决策. 它帮助规划器做出关于查询执行计划的明智决策，从而提高性能. 它设置为系统 RAM 的大约 75%
- max_connections: 限制了与服务器的并发连接数. 这个值的平衡对于资源管理很重要，防止服务器过载. 它是根据预期的用户负载和应用程序需求设置的
- checkpoint_timeout: 指定了执行检查点的频率. 更长的间隔减少了I/O开销，但增加了崩溃后的恢复时间
- autovacuum : 控制表和索引的自动清理以防止膨胀. 通过启用此参数，它有助于通过清理死行(dead rows)来维护数据库性能
- wal_buffers: 为 WriteAhead Logging (WAL) 缓冲区分配的内存量. 更大的缓冲区可以提高事务吞吐量，特别是在写入密集型环境中. 对于高事务量的系统，通常设置得更高.

### other
- wal_compression: 通过压缩 WAL 记录的数据部分来帮助减少 WAL 文件的大小, 以有效管理存储空间; 还可以降低 I/O 操作，从而提高整体系统性能
- wal_level: 默认设置为replica，适合大多数用例
- max_wal_size: 设置检查点之间的最大 WAL 数据量. 较高的值可以提高性能，但可能会延长恢复时间.
- min_wal_size: 定义保留的最小 WAL 空间量，有助于防止高负载期间 WAL 空间耗尽
- synchronous_commit: 控制事务提交是否等待 WAL 记录刷新到磁盘. 禁用此功能可以提高写入性能，但存在崩溃时数据丢失的风险
- archive_mode = on: 开启持续归档

  - archive_command : 单个归档目标地址

    - archive_command = 'rsync -z %p user@192.0.2.1:/path/to/archive/%f': 远程归档
  - archive_destinations: 多个归档目标地址(from pg 15)

  > 清理无用归档: `find /path/to/archive -type f -mtime +30 -delete`

## 迁移
- pg_dump/pg_restore
- pgloader
- FDW

## 监控wal和replcation
```sql
SELECT client_addr, pg_xlog_location_diff(sent_location,
replay_location) AS replication_lag
FROM pg_stat_replication;

SELECT archived_count, last_archived_wal,
last_archived_time
FROM pg_stat_archiver;

SELECT checkpoints_timed, checkpoints_req,
buffers_checkpoint, buffers_clean, buffers_backend
FROM pg_stat_bgwriter;
```