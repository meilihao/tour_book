## 检索

### 在where子句中使用别名列

将查询作为内联视图即可,不可漏掉下面的`x`(每个派生出来的表都必须有一个自己的别名).

    select * from (select SAL AS salary from EMP ) x where salary >1000;

### 连接列值

    select ename||' works as a'||job as msg from emp where deptno=10;//psql 使用双竖线作为连接运算符
    select concat(ename,' works as a',job) as msg from EMP where DEPTNO=10; //mysql使用concat函数

其实psql也支持concat,其`||`即是concat函数的简写形式,推荐使用concat,方便移植.

### 在select语句中使用条件逻辑

`case表达式`可以针对查询的返回值执行条件逻辑，该表达式也支持别名.

```sql
select ename,sal,
         case when sal <=2000 then 'under' 
         when sal>=4000 then 'over' 
         else 'mid' 
       end as status 
from emp;
```
**注意:在这里是用单引号表示字符序列.因为psql中会将双引号内容当做一个表的名字或者字段名字来对待.**

### 限制返回的行数

    select ename from emp limit 2;

### 从表中随机返回n行

mysql(使用内置函数:rand()):

    select ename,job from emp order by rand() limit 5;

postgres(使用内置函数:random()):

    select ename,job from emp order by random() limit 5;

ps:
`order by`子句,默认以升序排序,以逗号分隔排序列,优先次序是从左到右.其指定数字常量时,是按照select列表中的相应列来排序;使用随机函数时是先按照函数给每一行计算一个结果,再按结果排序.

### 将空值转换为实际值

使用coalesce函数用实际值替换空值.

    select coalesce(comm,0) from emp;

COALESCE表达式是 CASE 表达式的语法快捷方式:

    select case when comm is null then 0 else comm end from emp;

### 按模式搜索

like模式:`%`,匹配任意字符序列;`_`,匹配单个字符.

## 排序

### 按子串排序

按照job的最后三个字符排序:

     select ename,job from emp order by substr(job,length(job)-2);

>substr(string string,num start,num length):
string为字符串；
start为起始位置；
length为长度
注：**sql(mysql,postgres等)中的start是从1开始的**

### 对字母数字混合的数据排序

按照混合数据中的数字或字符来排序.

    create view v as select concat(ename,' ',deptno) as data from emp;

按data中的deptno排序:

    select * from v order by replace(data,replace(translate(data,'0123456789','##########'),'#',''),'')

按data中的ename排序:

    select * from v order by replace(translate(data,'0123456789','##########'),'#','')

>replace(string text, from text, to text) : 把字符串string里出现地所有子字符串from 替换成子字符串to
>
>     replace('abcdefabcdef', 'cd', 'XX')	abXXefabXXef
>
>translate(string text, from text, to text) : 把在string中包含的任何匹配from中的字符转化为对应的在to中的字符,即目标字符to和源字符from都可以同时指定多个.如果from比to长， 删掉在from中出现的额外的字符.
>
>     translate('12345', '143', 'ax')==a2x5
>[Postgres # 9.4. 字符串函数和操作符](http://www.postgres.cn/docs/9.3/functions-string.html)

因mysql不支持translate(),这个问题无解决方案.

### 处理排序空值

使用case表达式来"标记"一个值是否为NULL(排除空值的干扰),再排序.

所有空值在后且comm升序:

    select ename,sal,comm from (select ename,sal,comm, case when comm is null then 0 else 1 end as is_null from emp) x order by is_null desc,comm

### 根据数据项的键排序

如果job是"SALESMAN",要根据comm列排序,否则按sal列排序:

    select ename,sal,job,comm from emp order by case when job='SALESMAN' then comm else sal end;

## 多表操作

### 记录集的叠加

显示emp中deptno=10的员工名字和部门编号与dept中所有部门的名称和编号

    select ename as ename_and_dname,deptno from emp where deptno=10 union all select '---',null union all select dname,deptno from dept;

>UNION 操作符用于合并两个或多个 SELECT 语句的结果集.
**注意**:UNION 内部的 SELECT 语句必须拥有相同数量的列,列也必须拥有相似的数据类型,同时，每条 SELECT 语句中的列的顺序必须相同.UNION ALL 命令和 UNION 命令几乎是等效的，不过 UNION ALL 命令会列出所有的值(即包括重复项),而UNION会去重.**使用UNION时postgres会对新结果重新排序(排序方式未知)**,mysql还是按照各自select语句的结果直接合并;使用UNION ALL时,postgres和mysql均是按照各自select语句的结果直接合并.除非有必要,一般使用UNION ALL即可．

### 内联接

部门10中所有员工名称及其工作地点．

    select e.ename,d.loc from emp e,dept d where e.deptno=d.deptno and e.deptno=10;

相等于下面使用显式JOIN子句(INNER关键字可选)的语句:

    select e.ename,d.loc from emp e inner join dept d on (e.deptno=d.deptno) where e.deptno=10;

inner join(等值连接) 只返回两个表中联结字段相等的行，具体过程：

执行`select e.ename,d.loc,e.deptno as emp_deptno,d.deptno as despt_deptno from
 emp e,dept d  where e.deptno=10;`,获取表的笛卡尔积(行的所有可能组合),再取`emp.deptno=dept.deptno`相等的行.

### 从一个表中查找另一个表没有的值

从表dept中查找不存在与表emp关联的所有部门.

pg(支持差集操作EXCEPT):

    select deptno from dept except select deptno from emp;

except:获取第一个结果集,再从中去掉第二个结果集中也有的行.其限制条件:两个select列表中的值的数目和数据类型必须匹配.**except不会返回重复行,而且跟使用NOT IN的子查询不同,NULL值不会有问题.**

mysql:

    select distinct deptno from dept where deptno no in (select deptno from emp);

在mysql中使用not in时一定要注意NULL值问题.

执行`select deptno from dept where deptno not in (10,50,NULL);`(等价于`select deptno from dept wehre not (deptno=10 or deptno=50 or deptno=NUll);`),结果为空.

原因:sql中`任意值=NULL`,`任意值!=NULL`,`false or NULL`,`not NULL`的结果均为NULL.

- deptno=10=>not (true or false or NUll)=>not (true)=>false
- deptno=20=>not (false or false or NULL)=>not (NULL)=>NULL

>补充:
>IN和NOT IN的本质是OR运算.
>`NULL in (NULL)`和`12 in (13,NULL)`的结果均为NULL.

解决`NOT IN 和NULL`有关问题的方法:

1. 使用is:

       select deptno from dept wehre not (deptno=10 or deptno=50 or deptno is NUll);

2. 使用NOT EXISTS和相关子查询.

       select d.deptno from dept d where not exists (select null from emp e where d.deptno=e.deptno);

>补充:
>exists : 强调的是是否返回结果集，不要求知道返回什么.exists(sql 返回结果集为真);not exists(sql 不返回结果集为真)
>exists 与 in 最大的区别在于 in引导的子句只能返回一个字段.

### 在一个表中查询与其他表不匹配的记录

查询没有职员的部门信息.

    select d.* from dept d left outer join emp e on d.deptno=e.deptno where e.deptno is null;

### 查找重复的记录
select * from VAT2User a
where (a.UserId,a.TNo) in (select UserId,TNo from VAT2User group by UserId,TNo having count(*) > 1)
