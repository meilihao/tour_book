# 定义
- DB是存储数据的有序集合.
- table是一个二维数组的集合,是存储数据和操作数据的逻辑结构.
- SQL(Structured Query Language) 是用于访问和处理数据库的标准的计算机语言.
- SQL= 数据定义语言 (DDL - create,drop,alter)， 数据操作语言 (DML - - select,update,delete,insert)， 数据控制语言 (DCL - grant,revoke,commit,rollback)

## 运算符优先级(Operator Precedence)

- [Postgres](http://www.postgresql.org/docs/9.4/static/sql-syntax-lexical.html)
- [MySQL](https://dev.mysql.com/doc/refman/5.6/en/operator-precedence.html)

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
- 事务是需要在同一个处理单元中执行的一系列操作的集合.
- 事务特性ACID:原子性(Atomicity),一致性(Consistency),隔离性(Isolation)和持久性(Durability).

### 视图
- 视图中保存的是从表中取出数据所使用的`SELECT`语句.
- 多重视图会降低SQL的性能,不推荐.
- 通过视图更新数据有诸多限制,不推荐.

### 子查询
- 子查询即是一次性的视图.
- 多重子查询会降低SQL的性能,不推荐.
- 子查询可以通过`AS`关键字来命名,或用空格来分隔.
- 标量子查询(scalar subquery)指有且仅有一行一列的结果,其可以用在`SELECT`,`GROUP BY`,`HAVING`,`ORDER BY`子句等地方.
- 关联子查询就是指子查询与主查询之间有条件关联,不能独自执行.子查询的执行的次数依赖于外部查询，外部查询每执行一行，子查询执行一次,性能不佳.
- 在细分的组内进行比较时,需要使用关联子查询.

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
- 交叉连接`a CROSS JOIN b`,没有`ON`子句,即求笛卡儿积.
- `FULL JOIN`是`LEFT JOIN`与`RIGHT JOIN`的并集.
- 不推荐将连接条件写在`WHERE`中.
- 使用`JOIN`时,`SELECT`子句中的列需按照`<表名>.<列名>`的格式来书写.

参考 : ![图解SQL的JOIN操作](http://jbcdn2.b0.upaiyun.com/2013/05/SQL-Joins.jpg)

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
