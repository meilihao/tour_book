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

