# oracle
ref:
- [Oracle Database版本](https://en.wikipedia.org/wiki/Oracle_Database)
- [Oracle Database Documentation](https://docs.oracle.com/en/database/oracle/oracle-database/index.html)
- [Oracle Database 12c安装](https://blog.51cto.com/u_13728740/2293922)
- [Oracle之体系结构详解，基本操作管理及客户端远程连接](https://blog.51cto.com/u_13728740/2298336)

> Oracle Database Express Edition 是Oracle Database的简装版, 只有一个数据库实例XE, 因此也叫Oracle Database XE. 它拥有正式版的所有功能, 只是在内存和数据大小上做了限制, 适合初学者用来学习Oracle.

## 组件
- Oracle Database : 部署oracle db单机
- Oracle Grid Infrastructure : 部署oracle db cluster

## 部署
- [CentOS 7.4下安装Oracle 11.2.0.4数据库的方法](https://cloud.tencent.com/developer/article/1721536)

## cmd
```bash
# --- 登录db
su - oracle
sqlplus / as sysdba
> ? -- help
> conn 切换用户授权
> help index; -- 命令列表
> shutdown immediate; -- 停止oracle
> startup; -- 启动oracle
> SELECT *FORM TAB; -- 查看所有表
> DESC xxx; -- 查看xxx的表结构
> rename a to b; -- 修改表名
> show user; -- 查看当前连接用户
> @a.sql; -- 执行外部sql脚本
> select * from all_users; -- 查看所有用户
> select tablespace_name from user_tablespaces; -- 查询当前用户拥有的所的有表空间
> select * from database_properties where property_name=’DEFAULT_TEMP_TABLESPACE’; -- 查询默认临时表空间
> create tablespace animal
datafile 'animal.dbf' size 10M [autoextend on next 10m]; -- 创建表空间(Oracle自动将表空间名字全部转为大写). datafile 指定表空间对应的数据文件; size 后定义的是表空间的初始大小; autoextend on  自动增长,当表空间存储都占满时,自动增长; next 后指定的是一次自动增长的大小
> ho ls -lh /oracle/app/oracle/product/11.2.0/db_1/dbs/animal.dbf; -- 不退出sqlplus查看表空间信息
> alter database datafile '/oracle/app/oracle/oradata/school01.dbf' resize 80m; -- 调整数据文件大小
> alter tablespace animal add datafile '/oracle/app/oracle/oradata/school02.dbf' size 80m; -- 向表空间添加文件
> alter tablespace animal read only; -- 修改为只读权限
> alter tablespace animal read write; -- 修改为读写(默认)
> create user csy identified by csy
default tablespace ANIMAL; -- 创建用户名和密码均是csy的账号
> grant connect,resource,dba to csy; -- 赋予用户dba权限. Oracle数据库中常用角色:connect,连接角色.基本角色; resource,开发者角色;dba,超级管理员角色
> grant create session to csy ;
> create table dog
(
    name varchar(12),
    age varchar(12)
)
tablespace animal; -- 创建名为dog的表
> select tablespace_name, table_name from user_tables
where tablespace_name = 'ANIMAL'; -- 查看ANIMAL表空间下的所有表
> --修改表结构
> alter table person add gender number(1); ----添加一列
> alter table person add (name1 number(1),name2 number(1)); ----添加多列
> alter table persion modify gender char(1);  ----修改列类型
> alter table person rename column gender to sex; ----修改列名称
> alter table person drop column sex; ----删除一列
> DROP TABLESPACE animal INCLUDING CONTENTS AND DATAFILES; -- 删除数据库
> drop user csy cascade; -- 删除用户
> show con_name; -- 查看当前所在的容器
> show pdbs; -- 查看所有容器. CDB：默认的数据库; PDB：容器型数据库
> -- 索引
> create index 索引名称 on 表名（列名）; -- B树索引
> create unique index 索引名称 on 表名（列名）; -- 唯一索引/非唯一索引
> create index 索引名称 on 表名（列名）reverse; -- 反向索引
> create bitmap index 索引名称 on 表名（列名）; -- 位图索引
> create index 索引名称 on 表名（upper（列名））; -- 其他索引 //大写函数索引
> select index_name,index_type,table_name,tablespace_name from user_indexes; -- 查看索引
> select index_name,table_name,column_name from user_ind_columns where index_name like ‘EMP%’; -- 查看索引相关信息
> alter index 索引名称 rebuild; -- 重建索引
> alter index 索引名称 rebuild tablespace 表空间; -- 重建索引
> alter index 索引名称 coalesce; -- 合并索引碎片
> drop index 索引名称; -- 删除索引
> -- 创建物化视图
> create materialized view mtview
build immediate //创建完成立马生成新数据
refresh fast //刷新数据
on commit //提交
enable query rewrite //开启查询重写功能
as
select * from info; // info表必须存在主键
> drop materialized view mtview; //删除物化视图
> -- 序列
create sequence id_seq
start with 10 //初始值
increment by 1 //增量
maxvalue 1000 //最大值
nocycle //肺循环
cache 50; //缓存
> insert into info values (id_seq.nextval,‘tom’,80,to_date(‘2018-04-10’,‘yyyy-mm-dd’)); -- 插入数据时调用序列
> select id_seq.currval from dual; -- 查询序列当前值
> alter sequence id_seq cache 30; -- 更改序列
> select sequence_name,increment_by,cache_size from user_sequences; -- 查看序列信息
> drop sequence id_seq; -- 删除序列
```

## 数据监控em
- https://localhost:1158/em

## Oracle数据类型
- Varchar,varchar2	表示一个字符串
- number	

	- number(n) 表示一个整数,长度为n
	- number(m,n) 表示一个小数,总长度是m.小数是n,整数是m-n

- data	表示日期类型
- clob	大对象,表示大文本数据类型,可存4G
- blob	大对象,表示二进制数据,可存4G

## FAQ
### sqlplus报`ORA-01034: ORACLE not available`
出现ORA-01034的原因是多方面的：主要是Oracle当前的服务不可用, 用`startup;`启动即可

### sqlplus下使用backspace回删出现`^H`
方法:
1. 在sqlplus里面用ctrl+backspace代替backspace
2. 在使用sqlplus前执行`stty erase ^H`

	通过`stty -a`查看终端设置, 其中会有这样的一个字段`erase = ^?;`表示终端的清除字符的方式是Ctrl+Backspace, 可将它放入`.bashrc`

### sqlplus按方向键不支持显示历史命令
可安装软件rlwrap回调sqlplus中执行过的命令来解决

```bash
# dnf install rlwrap
# vim ~/.bashrc
alias sqlplus='rlwrap sqlplus'
...
```

### mysql和oracle 概念区别
ref:
- [谈谈mysql和oracle的使用感受 -- 差异](https://www.cnblogs.com/yougewe/p/13662695.html)

MySQL是一个以用户为中心的概念, 一个用户下, 拥有多个数据库, 一个数据库下拥有多个数据库表.

Oracle中，一个RDMS拥有多个实例(一般只有一个)，一个实例可以拥有多个用户，而一个用户可以拥有多个表空间（注意表空间多个用户可以交叉使用，但是前提是要有权限！），一个表空间可以有多个数据库文件（注意，数据库文件只能存在于一个表空间中，并不能交叉使用！）

> Oracle其实没有数据库，实际是表空间，是通过实例给每个用户分配不同的表空间，这样的表空间就是每个用户的数据库。所以直接创建表空间和对应的用户就可以了.

> 一个Oracle实例(Oracle instance) 有一系列的后台进程(backguound Processes)和内存结构(memory structures)组成.

表空间是Oracle对物理数据库上相关数据文件(ORA或者DBF文件)的逻辑映射. 一个数据库在逻辑上被划分1到若干个表空间,每个表空间包含了在逻辑上相关联的一组结构,每个数据库至少有一个表空间(称之为system表空间)

数据文件是数据库的物理存储单位. 数据库的数据是存储在表空间中的. 真正数据是在某一个或者多个数据文件中. 而一个表空间可以由一个或多个数据文件组成, 一个数据文件只能属于表空间, 一旦数据文件被加入到某个表空间后, 就不能删除这个文件, 如果要删除某个数据文件,只能删除其所属于的表空间才行.

### Temporary tablespace的作用
临时表空间主要用途是在数据库进行排序运算(如创建索引、order by及group by、distinct、union/intersect/minus/、sort-merge及join、analyze命令)、管理索引(如创建索引、IMP进行数据导入)、访问视图等操作时提供临时的运算空间，当运算完成之后系统会自动清理.

当临时表空间不足时，表现为运算速度异常的慢，并且临时表空间迅速增长到最大空间（扩展的极限），并且一般不会自动清理了.

如果临时表空间没有设置为自动扩展，则临时表空间不够时事务执行将会报ora-01652 无法扩展临时段的错误，当然解决方法也很简单:1、设置临时数据文件自动扩展，或者2、增大临时表空间.

### dns
`conn 用户名/密码@主机字符串`

### 数据库开启的三个状态
```bash
//开启三阶段：启动实例---------装载数据库(alter database mount)--------打开数据库(alter database open)
1：startup nomount （alter database mount; alter database open;）
2：startup mount （alter database open;）
3：startup
```