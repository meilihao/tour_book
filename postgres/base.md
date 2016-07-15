## 数据类型

### 整型
| 类型名称 | 范围 |存储需求|
|--------|--------|----------|
|  smallint      |        |2B|
|  int      |        |4B|
|  bigint      |        |8B|

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
|  char      | 固定长度非二进制字符串,不足不空白      |
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
运算符优先级: [官方文档](http://www.postgres.cn/docs/9.4/sql-syntax-lexical.html#SQL-PRECEDENCE)
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

## pg函数

### 数学函数
[数学函数](http://www.postgres.cn/docs/9.4/functions-math.html)
[字符串函数](http://www.postgres.cn/docs/9.4/functions-string.html)
>char_length(str),一个多字节字符算一个字符．
>length(str),使用utf8编码时，一个汉字3字节．

[时间/日期函数](http://www.postgres.cn/docs/9.4/functions-datetime.html)
[条件判断函数](http://www.postgres.cn/docs/9.4/functions-conditional.html)
[系统信息函数](http://www.postgres.cn/docs/9.4/functions-info.html)
[加密函数](http://www.postgres.cn/docs/9.4/functions-binarystring.html)

## 数据查询/select
>1. 除非必要应尽量避免使用`select *`,否则会因获取不需要的列数据而降低性能．
>2. 有时IN操作符可以实现OR操作符的效果，此时推荐使用IN,因为其语法更简明且执行速度更快，而且还支持更复杂的嵌套查询．
>3. group by通常和集合函数一起使用，例如，max(),min(),count(),sum(),avg().
>4. having和where均用于过滤数据．having用在数据分组后来过滤分组;where用于选择特定的记录．
>5．`limit x [offset y]`:［从表的第ｙ＋1条记录开始］选取ｘ条记录．
>6. `count(column)`会忽略column字段值为空值(NULL)的行;`count(*)`则所有记录均不忽略．
>7. `sum()`忽略列值为NULL的行．
>8. `max()`不仅适用于数值类型，也可用于字符类型，日期值．

### [连接查询](http://www.postgres.cn/docs/9.4/queries-table-expressions.html)
- 内连接：
格式：`select ... from t1 inner join t2 on t1.xxx=t2.xxx`
内连接也可用where子句来实现，性能可能更高．
自连接是内连接的特例．
- 外连接:
`left join(返回包括左表中的所有记录和右表中连接字段相等的记录)`和`right join(返回包括右表中的所有记录和左表中连接字段相等的记录)`
- 全连接
`FULL OUTER JOIN`,显示符合条件的数据行，同时显示左右不符合条件的数据行，相应的左右两边显示NULL，即显示左连接、右连接和内连接的并集.
- 交叉连接
`T1 CROSS JOIN T2`,对每个来自T1和T2 的行进行组合（也就是，一个笛卡尔积），连接成的表将包含这样的行： 所有T1里面的字段后面跟着所有T2 里面的字段。如果两表分别有 N 和 M 行，连接成的表将有 N*M 行.`FROM T1 CROSS JOIN T2`等效于`FROM T1,T2`.

### [子查询](http://www.postgres.cn/docs/9.4/functions-subquery.html)
子查询常用操作符：any/some,all,in,exists.
- any/some
`expression operator ANY (subquery)`,左边表达式使用operator对子查询结果的每一行进行一次计算和比较， 其结果必须是布尔值.如果至少获得一个真值，则ANY结果为"真"; 如果全部获得假值，则结果是"假"(包括子查询没有返回任何行的情况).

- all
`expression operator ALL (subquery)`,左边表达式使用operator对子查询结果的每一行进行一次计算和比较， 其结果必须是布尔值。如果全部获得真值，ALL结果为"真" (包括子查询没有返回任何行的情况)。如果至少获得一个假值，则结果是"假"。`NOT IN等效于<> ALL`.
- EXISTS
`EXISTS (subquery)`,判断subquery是否返回行。如果它至少返回一行，那么EXISTS的结果就为 "真"；如果子查询没有返回任何行，那么EXISTS的结果是"假"。`NOT EXISTS`和`EXISTS`作用相反．
- IN
`expression IN (subquery)`,左边表达式对子查询结果的每一行进行一次计算和比较。如果找到任何相等的子查询行， 则IN结果为"真"。如果没有找到任何相等行，则结果为"假" (包括子查询没有返回任何行的情况)。`NOT IN`和`IN`作用相反．

### 正则查询
[模式匹配#POSIX 正则表达式](http://www.postgres.cn/docs/9.4/functions-matching.html)

## 索引
索引是对数据库表中一列或多列的值进行排序的一种结构，其包含着对数据表里所有记录的引用指针，是提高数据库性能的常用方法．

PostgreSQL提供了好几种索引类型：B-tree, Hash, GiST, SP-GiST和GIN ．每种索引类型都比较适合某些特定的查询类型，因为它们用了不同的算法。缺省时， CREATE INDEX命令将创建 B-tree 索引，它适合大多数情况．
- B-tree
适合处理那些能够按顺序存储的数据,特别是在一个建立了索引的字段涉及到使用`<,<=,=,>=,>`.
- Hash(不推荐)
只能处理简单的等于比较。当一个索引了的列涉及到使用= 操作符进行比较的时候.
>不推荐原因：性能比B-tree弱，且不走WAL日志，因此数据库崩溃时需使用REINDEX重建Hash索引
- GiST
GiST 索引不是单独一种索引类型，而是一种架构，可以在这种架构上实现很多不同的索引策略。 因此，可以使用 GiST 索引的操作符高度依赖于索引策略(操作符类)。
- SP-GiST
SP-GiST索引类似于GiST索引，提供一个支持不同类型检索的架构。 SP-GiST允许实现许多各种不同的非平衡的基于磁盘的数据结构，例如四叉树，k-d树和基数树(字典树)。
- GIN
GIN 索引是反转索引，它可以处理包含多个键的值(比如数组)。与 GiST和SP-GiST 类似， GIN 支持用户定义的索引策略，可以使用 GIN 索引的操作符根据索引策略的不同而不同。

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
SetOp	INTERCECT，EXCEPT
有启动时
