# 定义
- DB是存储数据的有序集合.
- table是一个二维数组的集合,是存储数据和操作数据的逻辑结构.
- SQL(Structured Query Language) 是用于访问和处理数据库的标准的计算机语言.
- SQL的四部分:
  - 数据定义语言 (Data Definition Language, DDL : create,drop,alter)
  - 数据操作语言 (Data Manipulation Language, DML : select,update,delete,insert)
  - 事务控制语言 (Transaction Control Language, TCL: commit,rollback) 
  - 数据控制语言 (Data Control Language, DCL : grant,revoke)

SQL 数据库非常适合需要强数据一致性、定义良好的模式和复杂关系的应用程序. 其典型用例包括电子商务平台、金融系统和内容管理系统。例如，MySQL 提供ACID（原子性、一致性、隔离性与持久性）合规性，使其适合需要事务完整性的应用程序.

NoSQL 数据库非常适合优先考虑可扩展性、灵活性和高性能数据检索的应用程序. 它们在实时分析、社交媒体平台和 IoT（物联网）应用程序等用例中表现出色.

## 运算符优先级(Operator Precedence)

- [Postgres](http://www.postgresql.org/docs/9.4/static/sql-syntax-lexical.html)
- [MySQL](https://dev.mysql.com/doc/refman/5.6/en/operator-precedence.html)

## 范式
1. 第一范式（1NF）

在任何一个关系数据库中，第一范式（1NF）是对关系模式的基本要求，不满足第一范式（1NF）的数据库就不是关系数据库.
所谓第一范式（1NF）是指数据库表的每一列都是不可分割的基本数据项，同一列中不能有多个值，即实体中的某个属性不能有多个值或者**不能有重复的属性**.

2. 第二范式（2NF）

第二范式（2NF）是在第一范式（1NF）的基础上建立起来的，即满足第二范式（2NF）必须先满足第一范式（1NF）. 第二范式（2NF）要求数据库表中的每行必须可以被惟一地区分. 为实现区分通常需要为表加上一个列，以存储各个实例的惟一标识, 这个惟一属性列被称为主键.
第二范式（2NF）要求实体的**属性完全依赖于主键**. 简而言之，第二范式就是非主属性非部分依赖于主关键字。

3. 第三范式（3NF）

满足第三范式（3NF）必须先满足第二范式（2NF）. 简而言之，第三范式（3NF）要求一个数据库表中不包含已在其它表中已包含的非主键信息. 简而言之，第三范式就是属性不依赖于其它非主属性,即**消除传递依赖**.

## 表与表之间的关系（relation）
- 一对一（one-to-one）：一种对象与另一种对象是一一对应关系，比如一个学生只能在一个班级
- 一对多（one-to-many）： 一种对象可以属于另一种对象的多个实例，比如一张唱片包含多首歌
- 多对多（many-to-many）：两种对象彼此都是"一对多"关系，比如一张唱片包含多首歌，同时一首歌可以属于多张唱片

one_to_one               reverse one_to_one
one_to_many              reverse many_to_one(foreign_key)
many_to_one(foreign_key) reverse one_to_many
many_to_many             reverse many_to_many

# SQL

## SQL规则

### 命名/大小写

- 只使用英文状态下的**英文字母,数字,下划线**作为数据库,表和字段的名称,且以英文字母开头.
- sql关键字不区分大小写,**推荐使用大写**.

### 书写

- SQL中的字符串,日期时间或数字被称为常数,字符串和日期时间常数需使用单引号包裹,数字则直接书写即可.
- 别名(AS)通常**使用双引号包裹**(推荐)或直接书写,如果名称包含空格等字符时必须使用双引号包裹.
- 常数作为列的内容时,字符串和日期时间常数需使用单引号包裹,数字则直接书写.
- 注释: 单行注释用`-- `开头;多行注释用`/* */`包裹.注释可穿插在SQL语句中.
- `SELECT>FROM>WHERE>GROUP BY>HAVING`

## SQL语句

### ALTER TABLE

```sql
ALTER TABLE <表名> ADD COLUMN <列的定义>;
ALTER TABLE <表名> DROP COLUMN <列名>;
```

### 表重命名

```sql
// Postgres
ALTER TABLE [ IF EXISTS ] tbl_name RENAME TO new_tbl_name
// mysql
RENAME TABLE tbl_name TO new_tbl_name
```

### GROUP BY/HAVING

- `GROUP BY`字句中指定的列称为聚合键或分组列.聚合键中含NULL时,在结果中会以`空行`的形式来表示.
- `GROUP BY`允许使用别名.
- `GROUP BY`的结果是无序的.
- 使用聚合函数和`GROUP BY`时,`SELECT`和`HAVING`子句只能包含`常数,聚合函数和聚合键`.
- 聚合键所对应的条件不应该写在`HAVING`子句中,而应该写在`WHERE`子句中,且写在`WHERE`中更有性能优势.

### ORDER BY
- 排序列包含NULL时,会在开头或末尾显示.
- `ORDER BY`允许使用别名.
- `ORDER BY`允许使用未包含在`SELECT`中的列和聚合函数.

### INSERT
- 列名或值用逗号分隔,并用`()`包裹,这种形式称为清单,`INSERT`包含列清单和值清单.
- `INSERT INTO...SELECT`用于表复制.

### DELETE
- `DELETE`只能使用`WHERE`子句.
- Postgres和MySQL支持`TRUNCATE <表名>`来清空表,速度比`DELETE`快.

### UPDATE

- Postgres的`UPDATE`不支持`LIMIT`,但可将`LIMIT`放在`From`或`WHERE`的子句里.

### 事务
ref:
- [事务处理](https://icyfenix.cn/architect-perspective/general-architecture/transaction/)

事务的核心是故障恢复+隔离级别, 前者保证了数据库存储数据不会丢失，后者保证并发读写数据时的完整性.

- 事务是需要在同一个处理单元中执行的一系列操作的集合.
- 事务特性ACID:原子性(Atomicity),一致性(Consistency),隔离性(Isolation)和持久性(Durability).

原子性保证了事务内的所有操作是不可分割的，也就是它们要么全部成功，要么全部失败，不存在部分成功的情况. 成功的标志是在事务的最后会有提交（Commit）操作，它成功后会被认为整个事务成功. 而失败会分成两种情况，一种是执行回滚（Rollback）操作，另一种就是数据库进程崩溃退出.

事务的隔离性是不同的事务在运行的时候可以互相不干扰，就像没有别的事务发生一样.

持久性就是事务一旦被提交，那么它对数据库的修改就可以保留下来. 这里要注意这个“保存下来”不仅仅意味着别的事务能查询到，更重要的是在数据库面临系统故障、进程崩溃等问题时，提交的数据在数据库恢复后，依然可以完整地读取出来.

WAL日志在理论上可以无限增长，但实际上没有意义. 因为一旦数据从缓存中被刷入磁盘，该操作之前的日志就没有意义了，此时日志就可以被截断（Trim），从而释放空间, 而这个被截断的点，一般称为检查点. 检查点之前的页缓存中的脏页需要被完全刷入磁盘中.

#### 物理日志 Redo Log 与逻辑日志 Undo Log
事务对数据的修改其实是一种状态的改变，比如将 3 改为 5。这里我们将 3 称为前镜像（before-image），而 5 称为后镜像（after-image）, 可以得到如下公式：
- 前镜像+redo log=后镜像
- 后镜像+undo log=前镜像

redo log 存储了页面和数据变化的所有历史记录，称它为物理日志。而 undo log 需要一个原始状态，同时包含相对这个状态的操作，所以又被称为逻辑日志. 使用 redo 和 undo 就可以将数据向前或向后进行转换，这其实就是事务操作算法的基础.

redo 和 undo 有两种写入策略：steal 和 force.

steal 策略是说允许将事务中未提交的缓存数据写入数据库，而 no-steal 则是不能。可以看到如果是 steal 模式，说明数据从后镜像转变为前镜像了，这就需要 undo log 配合，将被覆盖的数据写入 undo log，以备事务回滚的时候恢复数据，从而可以恢复到前镜像状态。

force 策略是说事务提交的时候，需要将所有操作进行刷盘，而 no-force 则不需要。可以看到如果是 no-force，数据在磁盘上还是前镜像状态。这就需要 redo log 来配合，以备在系统出现故障后，从 redo log 里面恢复缓存中的数据，从而能转变为后镜像状态。

从上可知，当代数据库存储引擎大部分都有 undo log 和 redo log，那么它们就是 steal/no-force 策略的数据库.

### [隔离级别](https://juejin.im/post/5b90cbf4e51d450e84776d27)
- 脏读(dirty read/Read uncommitted)：一个事务读取了另一个事务尚未提交的修改
- 不可重复读(non-repeatable read/Read committed)：一个事务对**同一行**数据读取两次，得到不同结果, 即读到其他事务已提交的数据
- 幻读(phantom read/Repeatable read)：事务在操作过程中进行了两次查询，第二次的结果**包含了第一次未出现的新行或部分行消失**, 即被其他事务增删

    - 在 MySQL的REPEATABLE-READ中，InnoDB 通过 MVCC（多版本并发控制） +  Next-Key Lock(行锁（Record Lock）和间隙锁（Gap Lock）的结合) 防止幻读
    - pg REPEATABLE-READ也解决了幻读
- 串行化(Serializable)：一个事务在执行过程中完全看不到其他事务对数据库所做的更新．`写`会加`写锁`，`读`会加`读锁`,当出现读写锁冲突的时候，后访问的事务必须等前一个事务执行完成，才能继续执行.

    InnoDB 存储引擎在分布式事务的情况下一般会用到 SERIALIZABLE 隔离级别

标准的REPEATABLE-READ(可重复读) ：对同一字段的多次读取结果都是一致的，除非数据是被本身事务自己所修改，可以阻止脏读和不可重复读，但幻读仍有可能发生.

> 现在为止:所有的数据库都避免脏读
>
> 不可重复读是由于数据修改引起的，幻读是由数据插入或者删除引起的
>
> 串行化:可避免脏读、不可重复读、幻读的发生
>
> 随着事务隔离级别变得越来越严格，数据库对于并发执行事务的性能也逐渐下降
>
> MySQL的默认隔离级别就是Repeatable read, 查看方法:`show [global]  variables like "%isolation%";`

#### 悲观锁和乐观锁
悲观锁，数据库总是认为别人会去修改它所要操作的数据，因此在数据库处理过程中将数据加锁. 其实现依靠数据库底层.

乐观锁，总是认为别人不会去修改，只有在提交更新的时候去检查数据的状态. 通常基于数据版本（ Version ）记录机制(给数据增加一个字段来标识数据的版本)实现.

乐观控制使用的场景是并行事务不太多的情况，也就是只需要很少的时间来解决冲突, 对应方法: 提交前冲突检查.

基于锁的控制是典型的悲观控制.

### 视图
- 视图是虚拟的表, 是从真实表中取出数据所使用的`SELECT`语句. 使用视图可以简化复杂的sql操作，隐藏具体的细节，保护数据.
- 多重视图会降低SQL的性能,**不推荐**.
- 通过视图更新数据有诸多限制,**不推荐**.

### 子查询
- 子查询即是一次性的视图.
- 多重子查询会降低SQL的性能,不推荐.
- 子查询可以通过`AS`关键字来命名,或用空格来分隔.
- 标量子查询(scalar subquery)指有且仅有一行一列的结果,其可以用在`SELECT`,`GROUP BY`,`HAVING`,`ORDER BY`子句等地方.
- 关联子查询就是指子查询与主查询之间有条件关联,不能独自执行.子查询的执行的次数依赖于外部查询，外部查询每执行一行，子查询执行一次,性能不佳.
- 在细分的组内进行比较时,需要使用关联子查询.

子查询性能差的原因：子查询的结果集无法使用索引. 建议将子查询优化为 join 操作.

### 函数
参考 : SQL基础教程.MICK 的第6章.

### 谓词
- 谓词(predicate)即返回值是真值(即布尔值)的函数,作用是"判断是否存在满足某种条件的记录".
- 通常使用关联子查询作为`EXIST`的参数.

### UNION
- `UNION`会去除重复行,会对结果进行排序,但`UNION ALL`会保留重复行且不排序.
- `UNION`的列数和列的类型必须一致.
- `ORDER BY`只能用在`UNION`的最后,不能用于`SELECT`子句中.

### JOIN
- `JOIN`分内连接和外连接.
- 内连接`a INNER JOIN b ON a.xxx=b.xxx`
- 外连接`a {LEFT|RIGHT} [OUTER] JOIN b ON a.xxx=b.xxx`,不存在的列用NULL填充.
- 交叉连接`a CROSS JOIN b`,没有`ON`子句,即求笛卡儿积. **cross join后加条件只能用where,不能用on**.
- `FULL JOIN`是`LEFT JOIN`与`RIGHT JOIN`的并集. 当某行在另一个表中没有匹配行时，则另一个表的选择列用NULL值填充. 如果表之间有匹配行，则整个结果集行包含基表的数据.
- 不推荐将连接条件写在`WHERE`中.
- 使用`JOIN`时,`SELECT`子句中的列需按照`<表名>.<列名>`的格式来书写.

参考 : ![图解SQL的JOIN操作](/misc/img/sql/PxKTxj8.png)

example:
```sql
select * from table1 join table2 on table1.id=table2.id <=> select a.*,b.* from table1 a,table2 b where a.id=b.id <=> select * from table1 cross join table2 where table1.id=table2.id
select * from table1 cross join table2 <=> select * from table1,table2
```

### 窗口函数
- 通过`PARTITION BY`分组的记录集合被称为"窗口",代表"范围".
- 窗口函数兼具**分组和排序**两种功能,只能在`SELECT`中使用,格式是:
```sql
function (expression) OVER (
  [ PARTITION BY expression_list ] -- 指定分组的对象范围,将窗口缩小到特定的集合内
  [ ORDER BY order_list [ frame_clause ] ]
  [ frame_clause ])
```
- [frame_clause](http://www.postgres.cn/docs/9.4/sql-expressions.html#SYNTAX-WINDOW-FUNCTIONS)用于指定统计范围(其以当前记录作为基准),格式是：
```sql
{ RANGE | ROWS } frame_start
{ RANGE | ROWS } BETWEEN frame_start AND frame_end
```
- 窗口函数是对`WHERE`子句或者`GROUP BY`子句处理后的"结果"进行的操作.
- 相同的order_list,窗口函数的计算结果相同
- 没有`PARTITION BY`时,将整个表作为一个窗口,且聚合类窗口函数会按照order_list进行叠加处理.
- 窗口函数中的`ORDER BY`用于窗口内的排序,`FROM`后的`ORDER BY`用于查找结果的排序,作用不同,且可同时使用.
- 窗口函数包括:
1. 能够作为窗口函数的聚合函数(SUM,AVG,COUNT,MAX,MIN)
2. 专用窗口函数:RANK,DENSE_RANK,ROW_NUMBER等.
- RANK,DENSE_RANK,ROW_NUMBER区别:
1. ROW_NUMBER是在其分区中的当前行号
2. RANK有间隔的当前行排名,order_list的内容相同时行号相同,行号有间隔
3. DENSE_RANK没有间隔的当前行排名,order_list的内容相同时行号相同,行号连续

参考 : [Postgres的窗口函数](http://www.postgres.cn/docs/9.4/functions-window.html#FUNCTIONS-WINDOW-TABLE)

### GROUPING
- `ROLLUP`可以同时计算出合计和小计.
- 通过`GROUPING`函数可以简单地分辨出原始数据中的NULL和超级分组记录中的NULL.

参考 : [Postgres的GROUPING](https://www.postgresql.org/docs/devel/static/queries-table-expressions.html#QUERIES-GROUPING-SETS)

### 其他
- `EXISTS`会针对基础表的每条记录进行子查询操作,**不推荐使用**(因为基础表量大时很耗资源).
- `DISTINCT`必须紧跟在`SELECT`后(即第一个列名之前),如果其后存在多列时会将多列作为一个整体来去重.
- `SELECT`子句中可使用表达式,比如四则运算.
- 所有包含NULL的计算,其结果均为NULL.
- `<>`和`!=`作用相同,但`<>`属于SQL标准,推荐使用.
- 聚合函数会对NULL以外的数据进行合计,但`COUNT`除外.`COUNT({ [ [ ALL | DISTINCT ] expression ] | * })`有且仅有一个参数,
  `COUNT(*)`会统计包含NULL的行,而`COUNT(<列名>)`不会,`COUNT(整数)`==`COUNT(*)`,推荐`COUNT(pk)`;`DISTINCT`的结果会包含NULL,且支持所有的聚合函数.
- `MAX/MIN`支持所有可排序的数据类型,即不仅支持数字,也支持字符串,时间等.
- 聚合函数不能使用别名.
- `WHERE`字句不能使用聚合函数,只有`SELECT`,`ORDER BY`和`HAVING`子句可以.
- `WHERE`用来筛选数据行,`HAVING`用来指定分组的条件.
- [SQL解析顺序](http://www.jellythink.com/archives/924):
    ```
    (7)     SELECT
    (8)     DISTINCT <select_list>
    (1)     FROM <left_table>
    (3)     <join_type> JOIN <right_table>
    (2)     ON <join_condition>
    (4)     WHERE <where_condition>
    (5)     GROUP BY <group_by_list>
    (6)     HAVING <having_condition>
    (9)     ORDER BY <order_by_condition>
    (10)    LIMIT <limit_number>
    ```

    即`from -> on -> join -> where -> group -> having ->select -> distinct -> order by`

### 查询优化器
在一条单表查询语句真正执行之前，MySQL的查询优化器会找出执行该语句所有可能使用的方案，对比之后找出成本最低的方案. 这个成本最低的方案就是所谓的执行计划. 优化过程大致如下：
1. 根据搜索条件，找出所有可能使用的索引
2. 计算全表扫描的代价
3. 计算使用不同索引执行查询的代价
4. 对比各种执行方案的代价，找出成本最低的那一个

### NULL
参考:
- [MySQL中IS NULL、IS NOT NULL、!=不能用索引？胡扯！](https://www.tuicool.com/articles/yUF77zM)

值为NULL的二级索引记录都被放在了B+树的最左边，这是因为设计InnoDB规定: NULL值是列中最小的值.

IS NULL、IS NOT NULL、!=这些条件都可能使用到索引, InnoDB如何判断走索引还是全表: 成本, 优化器根据`index dive`或依据统计数据估算需要扫描的二级索引(即非主键索引)记录条数，如果这个条数占整个记录条数的比例特别大，那么就趋向于使用全表扫描执行查询，否则趋向于使用这个索引执行查询.

> index dive: 在查询真正执行前优化器就率先访问索引来计算需要扫描的索引记录数量的方式.

### index
**数据库索引**，是数据库管理系统中一个排序的数据结构(以某种方式引用/指向数据)，以实现高效查询. 索引的实现通常使用B树及其变种B+树.

为表设置索引要付出代价的：
1. 在插入, 修改和删除数据时要花费更多的时间, 因为索引也要随之变动
1. 增加了数据库的存储空间

优点:
1. 通过创建唯一性索引，可以保证数据库表中每一行数据的唯一性
1. 通过使用索引，可以在查询的过程中, 使用优化器, 大大加快数据的检索速度
1. 在使用分组和排序子句进行数据检索时，同样可以显著减少查询中分组和排序的时间
1. 加速表和表之间的连接

创建索引的考虑:
- 应该
  1. 在经常需要搜索的列上，可以加快搜索的速度
  1. 在作为主键的列上，强制该列的唯一性和组织表中数据的排列结构
  1. 在经常用在join的列上，这些列主要是一些外键，可以加快连接的速度
  1. 在经常需要根据范围进行搜索的列上创建索引，因为索引已经排序，其指定的范围是连续的
  1. 在经常需要排序的列上创建索引，因为索引已经排序，这样查询可以利用索引的排序，加快排序查询时间
  1. 在经常使用在WHERE子句中的列上面创建索引，加快条件的判断速度
- 不应该
  1. 对于那些在查询中很少使用或者参考的列不应该创建索引. 这是因为，既然这些列很少使用到，因此有索引或者无索引，并不能提高查询速度. 相反，由于增加了索引，反而降低了系统的维护速度和增大了空间需求.
  1. 对于那些只有很少数据值的列(列的基数低)也不应该增加索引. 这是因为，由于这些列的取值很少，例如user表的性别列，在查询的结果中，结果集的数据行占了表中数据行的很大比例，即需要在表中搜索的数据行的比例很大. 增加索引，并不能明显加快检索速度.
  1. 对于那些定义为text, image和bit数据类型的列不应该增加索引. 这是因为，这些列的数据量要么相当大，要么取值很少.
  1. 当修改性能远远大于检索性能时，不应该创建索引. 这是因为，修改性能和检索性能是互相矛盾的. 当增加索引时，会提高检索性能，但是会降低修改性能。当减少索引时，会提高修改性能，降低检索性能. 因此，当修改性能远远大于检索性能时，不应该创建索引
  1. 查询中使用了函数或表达式

索引分类:
- 按数据结构维度:
  - BTree索引（B-Tree或B+Tree索引）// O(LogN),相当于二分查找

    B+Tree索引是B-Tree的改进版本: 数据都在叶子节点上，并且增加了顺序访问指针，每个叶子节点都指向相邻的叶子节点的地址.

    **聚簇索引的叶子节点是整行数据, 非聚簇索引的叶子节点是主键的值**.

    > InnoDB 使用 B+Tree 作为索引结构
  - Hash索引

    Hash 索引不支持顺序和范围查询

    基于哈希表实现，哈希索引适合等值查询，但是不无法进行范围查询, 没办法利用索引完成排序, 不支持多列联合索引的最左匹配规则; 如果有大量重复键值得情况下，哈希索引的效率会很低(哈希碰撞)
  - full-index全文索引

    一般不会使用，效率较低，通常使用搜索引擎如 ElasticSearch 代替
  - 空间索引: R-Tree

    一般不会使用，仅支持 geometry 数据类型，优势在于范围查找，效率较低，通常使用搜索引擎如 ElasticSearch 代替
- 按应用维度:
    - 唯一索引

        索引列(允许多列)的值必须唯一，但允许有NULL值, 其不允许其中任何两行具有相同索引值的索引
    - 主键索引
        不允许NULL值的唯一索引
    - 普通索引

        一个索引只包含单个列，一个表可以有多个单列索引
    - 复合索引

        一个索引包含多个列
    - 全文索引：对文本的内容进行分词，进行搜索
    - 前缀索引：对文本的前几个字符创建索引
- 根据底层存储方式(数据的物理顺序与键值的逻辑（索引）顺序关系)/
    - 聚集索引
        **并不是一种单独的索引类型，而是一种数据存储方式**. 具体细节取决于不同的实现，**InnoDB的聚簇索引其实就是在同一个结构中保存了B+Tree索引和数据行**

        一个聚集索引定义了**表中数据的物理存储顺序**, 因为行记录只能按照一个维度进行排序，所以**一张表只能有一个聚集索引**. 与非聚集索引相比，聚集索引通常提供更快的数据访问速度.

        聚集索引一般是表中的主键索引，如果表中没有显示指定主键，则会选择表中的第一个不允许为NULL的唯一索引，如果还是没有的话，就采用Innodb存储引擎为每行数据内置的6字节ROWID作为聚集索引.

        - 优点：
            1. 查询速度非常快：聚簇索引的查询速度非常的快，因为整个 B+ 树本身就是一颗多叉平衡树，叶子节点也都是有序的，定位到索引的节点，就相当于定位到了数据。相比于非聚簇索引， 聚簇索引少了一次读取数据的 IO 操作
            1. 对排序查找和范围查找优化：聚簇索引对于主键的排序查找和范围查找速度非常快
            
        - 缺点:
            1. 依赖于有序的数据：因为 B+ 树是多路平衡树，如果索引的数据不是有序的，那么就需要在插入时排序，如果数据是整型还好，否则类似于字符串或 UUID 这种又长又难比较的数据，插入或查找的速度肯定比较慢
            1. 更新代价大：如果对索引列的数据被修改时，那么对应的索引也将会被修改，而且聚簇索引的叶子节点还存放着数据，修改代价肯定是较大的，所以对于主键索引来说，主键一般都是不可被修改的

        > 在 MySQL 中，InnoDB 引擎的表的 `.ibd`文件就包含了该表的索引和数据
    - 非聚集索引
        不是聚簇索引的索引, 即索引结构和数据分开存放的索引，二级索引（辅助索引）就属于非聚簇索引

        二级索引（Secondary Index）的叶子节点存储的数据是主键的值. 唯一索引、普通索引、前缀索引等索引都属于二级索引.

        - 优点：更新代价比聚簇索引要小, 因为其叶子节点是不存放数据.
        - 缺点：
            - 依赖于有序的数据：跟聚簇索引一样，非聚簇索引也依赖于有序的数据
            - 可能会二次查询（回表）：这应该是非聚簇索引最大的缺点了

> 覆盖索引: 一个查询语句的执行只用从索引中就能够取得，不必回表.

> MySQL主要提供2种方式的索引：B-Tree索引，Hash索引, 全文索引.

参考:
- [Postgresql、MySQL相关的四种索引类型：B-Tree，Hash，Gist，GIN](https://my.oschina.net/wangchuandui/blog/1504260)

> 一般使用磁盘I/O次数评价索引结构的优劣

相关问题:
1. 为什么索引结构默认使用B-Tree，而不是hash，二叉树，红黑树

    - hash ：虽然可以快速定位，但是没有顺序，IO复杂度高
    - 二叉树 ：树的高度不均匀，**不能自平衡，查找效率跟数据有关（树的高度）**，并且IO代价高
    - 红黑树 ：树的高度随着数据量增加而增加，IO代价高
1. 为什么官方建议使用自增长主键作为索引

    结合B+Tree的特点，自增主键是连续的，在插入过程中尽量减少页分裂，即使要进行页分裂，也只会分裂很少一部分. 它还能减少数据的移动，每次插入都是插入到最后.

    总之就是减少分裂和移动的频率.

索引在什么情况下会失效:
- 条件中有or
- 使用like模糊查询以%开头的
- 在索引列上进行计算，使用函数，隐式转化
- 对于组合索引，不遵循最左匹配原则
- 在索引字段上使用is null / is not null判断时会导致索引失败

#### 数据库优化的思路
思路来源通常是监控, 比如慢查询日志.

1. EXPLAIN查看SQL执行计划

    执行计划结果中共有 12 列，各列代表的含义总结如下表(列名+含义)：
    - id: SELECT 查询的序列标识符

        用于标识每个 SELECT 语句的执行顺序
    - select_type: SELECT 关键字对应的查询类型

        - SIMPLE：简单查询，不包含 UNION 或者子查询
        - PRIMARY：查询中如果包含子查询或其他部分，外层的 SELECT 将被标记为 PRIMARY
        - SUBQUERY：子查询中的第一个 SELECT
        - UNION：在 UNION 语句中，UNION 之后出现的 SELECT
        - DERIVED：在 FROM 中出现的子查询将被标记为 DERIVED
        - UNION RESULT：UNION 查询的结果

    - table : 用到的表名
    - partitions: 匹配的分区，对于未分区的表，值为 NULL
    - type: 表的访问方法
    - possible_keys: 可能用到的索引
    - key: 实际用到的索引
    - key_len: 所选索引的长度
    - ref: 当使用索引等值查询时，与索引作比较的列或常量
    - rows: 预计要读取的行数
    - filtered: 按表条件过滤后，留存的记录数的百分比
    - Extra: 附加信息

        常见的值如下：
        - Using filesort：在排序时使用了外部的索引排序，没有用到表内索引进行排序
        - Using temporary：MySQL 需要创建临时表来存储查询的结果，常见于 ORDER BY 和 GROUP BY
        - Using index：表明查询使用了覆盖索引，不用回表，查询效率非常高
        - Using index condition：表示查询优化器选择使用了索引条件下推这个特性
        - Using where：表明查询使用了 WHERE 子句进行条件过滤。一般在没有使用到索引的时候会出现
        - Using join buffer (Block Nested Loop)：连表查询的方式，表示当被驱动表的没有使用索引的时候，MySQL 会先将驱动表读出来放到 join buffer 中，再遍历被驱动表与驱动表进行查询
        
        当 Extra 列包含 Using filesort 或 Using temporary 时，MySQL 的性能可能会存在问题，需要尽可能避免

1. SQL语句优化
    1. 应尽量避免在 where 子句中使用!=或<>操作符，否则将引擎放弃使用索引而进行全表扫描
    1. 应尽量避免在 where 子句中对字段进行 null 值判断，否则将导致引擎放弃使用索引而进行全表扫描.
        ```sql
        select id from t where num is null
        -- 可以在num上设置默认值0，确保表中num列没有null值，然后再查询：
        select id from t where num=0
        ```
    1. SQL语句中IN包含的值不应过多

        对于IN里的**连续数值**，能用between就不要用in.
    1. 很多时候用 exists 代替 in 是一个好的选择
    1. SELECT语句务必指明字段名称
    1. 当只需要一条数据的时候，使用limit 1
    1. 如果排序字段没有用到索引，就尽量少排序
    1. 如果限制条件中其他字段没有索引，尽量少用or

        or两边的字段中，如果有一个不是索引字段，而其他条件也不是索引字段，会造成该查询不走索引的情况. 很多时候使用union all或者是union（必要的时候）的方式来代替“or”会得到更好的效果
    1. 尽量用union all代替union

        union和union all的差异主要是前者需要将结果集合并后再进行唯一性过滤操作
    1. 不使用ORDER BY RAND()

        ```sql
        select id from `dynamic` t1 join (select rand() * (select max(id) from `dynamic`) as nid) t2 on t1.id > t2.nidlimit 1000;
        ```
    1. 区分in和exists、not in和not exists

        ```sql
        select * from 表A where id in (select id from 表B)
        --等价于
        select * from 表A where exists(select * from 表B where 表B.id=表A.id)
        ```

        区分in和exists主要是造成了驱动顺序的改变（这是性能变化的关键），如果是exists，那么以外层表为驱动表，先被访问，如果是IN，那么先执行子查询. 所以IN适合于外表大而内表小的情况；EXISTS适合于外表小而内表大的情况.

        关于not in和not exists，推荐使用not exists，不仅仅是效率问题，not in可能存在逻辑问题.

        ```sql
        -- 结果集: A - (A ∩ B)
        select colname … from A表 where a.id not in (select b.id from B表)
        --等价于
        select colname … from A表 Left join B表 on where a.id = b.id where b.id is null
        ```
    1. 使用合理的分页方式以提高分页的效率

        ```sql
        -- 随着表数据量的增加，直接使用limit分页查询会越来越慢
        select id,name from product limit 866613, 20
        -- 优化：可以取前一页的最大行数的id，然后根据这个最大的id来限制下一页的起点
        select id,name from product where id> 866612 limit 20
        ```
    1. 分段查询

        在一些用户选择页面中，可能一些用户选择的时间范围过大，造成查询缓慢。主要的原因是扫描行数过多. 这个时候可以通过程序，分段进行查询，循环遍历，将结果合并处理进行展示.
    1. 不建议使用`%`前缀模糊查询

        LIKE`%name`或者LIKE`%name%`，这种查询会导致索引失效而进行全表扫描。但是可以使用LIKE `name%`(利用 B+索引的“最左前缀”).

        查询%name%可使用全文索引:
        ```sql
        ALTER TABLE `dynamic_201606` ADD FULLTEXT INDEX `idx_user_name` (`user_name`);
        select id,fnum,fdst from dynamic_201606 where match(user_name) against('zhangsan' in boolean mode);
        ```
    1. 避免在where子句中对字段进行表达式操作

        ```sql
        -- 进行了算术运算，这会造成引擎放弃使用索引
        select user_id,user_project from user_base where age*2=36;
        -- 建议业务层处理
        select user_id,user_project from user_base where age=18;
        ```
    1. 避免隐式类型转换, 该操作可能会导致优化器生成执行计划有问题

        where子句中出现column字段的类型和传入的参数类型不一致的时候发生的类型转换, 而不使用索引
    1. 对于联合索引来说，要遵守最左前缀法则

        举列来说索引含有字段id、name、school，可以直接用id字段，也可以id、name这样的顺序，但是name;school都无法使用这个索引. 所以在创建联合索引的时候一定要注意索引字段顺序，常用的查询字段放在最前面
    1. 必要时可以使用force index来强制查询走某个索引
    1. 注意范围查询语句

        对于联合索引来说，如果存在范围查询，比如between、>、<等条件时，会造成后面的索引字段失效
    1. 关于JOIN优化

        - `LEFT JOIN` : 左表为驱动表
        - `INNER JOIN` : MySQL会自动找出那个数据少的表作用驱动表
        - `RIGHT JOIN` : 右表为驱动表

        建议利用小表去驱动大表.

        合理利用索引：被驱动表的索引字段作为on的限制字段
    1. 用Where子句替换HAVING子句: 因为HAVING 只会在检索出所有记录之后才对结果集进行过滤

2. 索引优化

3. 数据库结构优化

    1. 范式优化： 比如消除冗余, 以节省空间
    1. 反范式优化：比如适当加冗余等（减少join）
    1. 拆分表

      分区将数据在物理上分隔开，不同分区的数据可以制定保存在处于不同磁盘上的数据文件里. 这样，当对这个表进行查询时，只需要在表分区中进行扫描，而不必进行全表扫描，明显缩短了查询时间，另外处于不同磁盘的分区也将对这个表的数据传输分散在不同的磁盘I/O，一个精心设置的分区可以将数据传输对磁盘I/O竞争均匀地分散开. 对数据量大的时时表可采取此方法, 比如可按月自动建表分区.

      拆分其实又分垂直拆分(列的拆分)和水平拆分(行的拆分).
      水平拆分: 订单表通过完成状态拆分为已完成订单和未完成订单

4. 服务器硬件优化

#### drop,delete与truncate的区别
drop直接删掉表; truncate删除表中数据，再插入时自增长id又从1开始; delete删除表中数据，可以加where字句:

1. DELETE语句执行删除的过程是每次从表中删除一行，并且同时将该行的**删除操作作为事务记录在日志中保存以便进行进行回滚操作**. TRUNCATE TABLE 则一次性地从表中删除所有的数据并**不把单独的删除操作记录记入日志保存，删除行是不能恢复的**, 并且在删除的过程中**不会激活与表有关的删除触发器**, 执行速度快.
1. 表和索引所占空间.当表被TRUNCATE 后，这个表和索引所占用的空间会恢复到初始大小，而DELETE操作不会减少表或索引所占用的空间.drop语句将表所占用的空间全释放掉.
1. 一般而言，速度: drop > truncate > delete
1. 应用范围.TRUNCATE 只能对TABLE；DELETE可以是table和view
1. TRUNCATE 和DELETE只删除数据，而DROP则删除整个表（结构和数据）
1. truncate与不带where的delete ：只删除数据，而不删除表的结构（定义）. drop语句将删除表的结构被依赖的约束（constrain),触发器（trigger)索引（index);依赖于该表的存储过程/函数将被保留，但其状态会变为：invalid.
1. delete语句为DML（data maintain Language),这个操作会被放到 rollback segment中,事务提交后才生效.如果有相应的 tigger,执行的时候将被触发.
1. truncate、drop是DLL（data define language),操作立即生效，原数据不放到 rollback segment中，不能回滚
1. 在没有备份情况下，谨慎使用 drop 与 truncate.要删除部分数据行采用delete且注意结合where来约束影响范围.回滚段要足够大.要删除表用drop;若想保留表而将表中数据删除，如果于事务无关，用truncate即可实现.如果和事务有关，或想触发trigger,还是用delete.
1. Truncate table 表名 速度快,而且效率高,因为: truncate table 在功能上与不带 WHERE 子句的 DELETE 语句相同：二者均删除表中的全部行.但 TRUNCATE TABLE 比 DELETE 速度快，且使用的系统和事务日志资源少.DELETE 语句每次删除一行，并在事务日志中为所删除的每行记录一项.TRUNCATE TABLE 通过释放存储表数据所用的数据页来删除数据，并且只在事务日志中记录页的释放.
1. TRUNCATE TABLE 删除表中的所有行，但表结构及其列、约束、索引等定义保持不变.新行标识所用的计数值重置为该列的种子.如果想保留标识计数值，请改用 DELETE.如果要删除表定义及其数据，请使用 DROP TABLE 语句.
1. 对于由 FOREIGN KEY 约束引用的表，不能使用 TRUNCATE TABLE，而应使用不带 WHERE 子句的 DELETE 语句.由于 TRUNCATE TABLE 不记录在日志中，所以它不能激活触发器.

#### 存储过程与触发器的区别
注意: 推荐业务层实现: 避免绑定db; 业务不需要知道db的内容.

触发器与存储过程非常相似，触发器也是SQL语句集，两者唯一的区别是**触发器不能用EXECUTE语句调用，而是在用户执行SQL语句时自动触发（激活）执行**. 触发器是在一个修改了指定表中的数据时执行的存储过程. 通常通过创建触发器来强制实现不同表中的逻辑相关数据的引用完整性和一致性. 由于用户不能绕过触发器，所以可以用它来强制实施复杂的业务规则，以确保数据的完整性. **触发器不同于存储过程，触发器主要是通过事件执行触发而被执行的，而存储过程可以通过存储过程名称名字而直接调用**. 当对某一表进行诸如UPDATE、INSERT、DELETE这些操作时，RDMS就会自动执行触发器所定义的SQL语句，从而确保对数据的处理必须符合这些SQL语句所定义的规则.

## FAQ
### 索引失效的情况
1. 索引字段使用 or 时, 要想使用or，又想让索引生效，只能将or条件中的每个列都加上索引
1. 使用函数, 比如`SUM`,`AVG`,`DAY`等
1. 在MYSQL使用不等于（`<>`,`!=`）的时候
1. like以通配符开头（`'%abc...'`）的时候
1. 如果mysql估计使用全表扫描要比使用索引快,则不使用索引
1. 对于多列索引不遵循最左匹配规则的时候
1. where的数据类型与定义的不一致
1. 在JOIN操作中主键和外键的数据类型不相同
1. 使用 is null / is not null 判断时
