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
`order by`子句指定数字常量时,是按照select列表中的相应列来排序;使用随机函数时是先按照函数给每一行计算一个结果,再按结果排序.

### 将空值转换为实际值

使用coalesce函数用实际值替换空值.

    select coalesce(comm,0) from emp;

COALESCE表达式是 CASE 表达式的语法快捷方式:

    select case when comm is null then 0 else comm end from emp;

### 按模式搜索

like模式:`%`,匹配任意字符序列;`_`,匹配单个字符.