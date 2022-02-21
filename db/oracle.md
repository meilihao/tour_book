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
# lsnrctl status : 查看服务器端listener进程的状态
# --- 登录db
su - oracle
> sqlplus --不显露密码的登陆方式
Enter user-name：sys
Enter password：password as sysdba --以sys用户登陆的话 必须要加上 as sysdba 子句
> -- 直接登入db
sqlplus / as sysdba
> exit
> -- 先进入sqlplus再登入
> sqlplus /nolog -- /nolog是不登陆(no login)到数据库服务器的意思. 如果没有/nolog参数，sqlplus会提示输入用户名和密码
> connect / as sysdba -- 连接db by 用户授权
> select user from dual; -- 查看当前用户
> ? -- help
> help index; -- 命令列表
> shutdown immediate; -- 停止oracle
> startup; -- 启动oracle
> SELECT *FORM TAB; -- 查看所有表
> DESC xxx; -- 查看xxx的表结构
> rename a to b; -- 修改表名
> show user; -- 查看当前连接用户
> @a.sql; -- 执行外部sql脚本
> select name from v$database; -- 获取当前数据库实例
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
> grant create session to csy;
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
> show sga -- 查看instance是否已启动
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

## 备份/还原
ref:
- [Oracle备份的几种方式](https://www.cnblogs.com/lcword/p/11775657.html)
- [rman备份与恢复](https://zhuanlan.zhihu.com/p/143866731)
- [Oracle Databases Enterprise Bacula Plugin Quick Guide](https://www.bacula.lat/oracle-databases-enterprise-bacula-plugin-quick-guide)
- [Database Backup and Recovery User's Guide](https://docs.oracle.com/cd/E11882_01/backup.112/e10642/toc.htm)
- [Oracle exp/imp数据导入导出工具基本用法](https://www.cnblogs.com/pandachen/p/5935078.html)
- [expdp impdp 数据库导入导出命令详解](https://blog.51cto.com/shitou118/310033)

> 备份需要sysdba权限

EXP和IMP是客户端工具程序，它们既可以在客户端使用，也可以在服务端使用.
EXPDP和IMPDP是服务端的工具程序，他们只能在ORACLE服务端使用，不能在客户端使用.
IMP只适用于EXP导出的文件，不适用于EXPDP导出文件；IMPDP只适用于EXPDP导出的文件，而不适用于EXP导出文件.
rman: RMAN可以进行增量备份, **推荐使用**.

备份:
```bash
# exp help=y # exp help
# exp \'sys/xxx as sysdba\' file=db.dmp full=y -- 1.将数据库完全导出. 用操作系统权限认证的oracle sys管理员身份, xxx是任意密码
# exp csy/csy file=db.dmp full=y -- 1.将数据库完全导出. 用csy账户
# exp system/manager@orcl file=db.dmp full=y -- 1.将数据库完全导出，设置full选项
# exp system/manager@orcl file=db.dmp rows=n full=y -- 2、导出数据库结构，不导出数据，设置rows选项
# exp system/manager@orcl file=db1.dmp,db2.dmp filesize=50M full=y -- 3、当导出数据量较大时，可以分成多个文件导出，设置filesize选项
# exp system/manager@orcl file=Test_bak.dmp owner=(system,sys) -- 4.将数据库中system用户与sys用户的表导出，设置owner选项
# exp system/manager@orcl file=Test_bak.dmp tables=(t_result,t_khtime) -- 5.将数据库中的表t_result,t_khtime导出，设置tables选项
# exp kpuser/kpuser@orcl file=Test_bak.dmp tables=(T_SCORE_RESULT) query=\" where updatedate>to_date('2016-9-1 18:32:00','yyyy-mm-dd hh24:mi:ss')\" -- 6、将数据库中的表T_SCORE_RESULT中updatedate日期字段大于某个值的数据导出，设置query选项
```

还原:
```bash
# imp system/manager@orcl file=Test_bak.dmp ignore=y -- 1、导入dmp文件，如果表已经存在，会报错且不导入已经存在的表，设置ignore选项
# imp kpuser/kpuser@orclfile=kpuser.dmp tables=(T_SCORE_RESULT) -- 2、导入dmp文件中部分指定的表，设置tables选项
# -- 3、导入一个或一组指定用户所属的全部表、索引和其他对象，设置fromuser选项
# imp system/manager@orcl file=kpuser.dmp fromuser=kpuser //kpuser必须存在
# imp system/manager@orcl file=users.dmp fromuser=(kpuser,kpuser1,test) //kpuser,kpuser1,test用户必须存在
# -- 4、将数据导入指定的一个或多个用户，设置fromuser和touser选项
# imp system/manager file=kpuser.dmp fromuser=kpuser touser=kpuser1 //kpuser1必须存在
# imp system/manager file=users.dmp fromuser=(kpuser,kpuser1) touser=(kpuser2, kpuser3) //kpuser2、kpuser3必须存在
```

### rman
> RMAN-SBT是指rman备份到tape.

前提:
1. `SELECT LOG_MODE FROM SYS.V$DATABASE;`/`archive log list`, db在ARCHIVELOG模式, 默认是NOARCHIVELOG

	启用ARCHIVELOG模式:
	```
	> archive log list -- 查看Database Archiving Mode, **推荐**
	> shutdown -- 关闭db
	> exit
	# mkdir -p /mnt/archive
	# chown oracle:oinstall /mnt/archive 
	# sqlplus / as sysdba
	> startup mount -- 以加载方式启动
	> alter database archivelog; -- 修改归档模式
	> alter system Set LOG_ARCHIVE_DEST_1='LOCATION=/mnt/archive' -- /mnt/archive 要存在
	> archive log list -- 检查参数
	> shutdown immediate -- 关闭db
	> connect / as sysdba
	> startup
	# --- 另一个terminal
	# rman target / log a.log -- 指定log后, rman日志会输出到a.log而不是terminal
	> backup database; / backup database format "/home/oracle/%u";
	```

	> log_archive_dest_1会在`{instance}/dbs/xxx.ora`里


备份数据库指定文件:
1. 获取指定文件的file_id

	- 通过数据字典dba_data_files查询出表空间对应的数据文件及其序号: `Select file_name, file_id, tablespace_name from dba_data_files;`
	- 查看某个表对应的序号及表空间: `Select file_name, file_id, tablespace_name from dba_data_files where file_id in (select distinct file_id from dba_extents where segment_name='表名');`
1. 备份

	```
	# rman target /
	> backup datafile 2,7 format "/home/oracle/%u"; -- 2,7为要备份文件的file_id
	```

## FAQ
### sqlplus报`ORA-01034: ORACLE not available`
出现ORA-01034的原因是多方面的：主要是Oracle当前的服务不可用, 用`startup;`启动即可

### sqlplus下使用backspace回删出现`^H`
方法:
1. 在sqlplus里面用ctrl+backspace代替backspace
2. 在使用sqlplus前执行`stty erase ^H`

	通过`stty -a`查看终端设置, 其中会有这样的一个字段`erase = ^?;`表示终端的清除字符的方式是Ctrl+Backspace, 可将它放入`.bashrc`

### sqlplus不支持方向键
可安装软件rlwrap回调sqlplus中执行过的命令来解决.

```bash
# dnf install rlwrap
# vim ~/.bashrc
alias sqlplus='rlwrap sqlplus'
...
```

自编译:
```bash
# dnf install readline-devel
# wget https://github.com/hanslub42/rlwrap/releases/download/v0.43/rlwrap-0.43.tar.gz
# tar -xf rlwrap-0.43.tar.gz
# cd rlwrap-0.43
# ./configure && make && make install
# vim ~/.bashrc
alias sqlplus='rlwrap sqlplus'
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

STARTUP 选项说明：
- NOMOUNT—开启实例，不加载数据库.允许访问数据库，仅用于创建数据库或重建控制文件
- MOUNT—开启实例，并加载数据库，但不打开数据库。允许DBA进行操作，但是不允许普通的数据库访问。
- OPEN—开启实例，加载数据库，打开数据库,等同STARTUP
- FORCE-在启动或关闭遇到问题时，强制启动实例
- OPEN RECOVER—在完成完整的备份后启动实例

### db登入方式
1. `sqlplus / as sysdba` : =`sqlplus sys/xxx as sysdba`(xxx为任意密码).这是以操作系统权限认证的oracle sys管理员登陆，不需要listener进程
2. `sqlplus sys/oracle` : 非管理员用户登录. 这种连接方式只能连接本机数据库，同样不需要listener进程
3. `sqlplus scott/oracle@orcl` : 非管理员用户使用tns别名登录. 这种方式需要listener进程处于可用状态, 最普遍的通过网络连接
3. `sqlplus sys/oracle@orcl as sysdba` : 管理员用户使用tns别名登录. 这种方式需要listener进程处于可用状态

以上连接方式使用sys用户或者其他通过密码文件验证的用户都不需要数据库处于可用状态，操作系统认证也不需要数据库可用，普通用户因为是数据库认证，所以数据库必需处于open状态

> 当给某个用户赋予权限的时候,可以直接对其赋予权限. 也可以先将若干权限形成一个集合体, 再将这个集合体整体赋予该用户. 这里这个权限的集合体就是角色(role), 比如sysdba.

### [print_table 实现 sqlplus 类似 mysql \G 及 psql \x 的功能](https://icode.best/i/31745333641226)
```sql
> create or replace procedure print_table( p_query in varchar2 )
AUTHID CURRENT_USER
is
	l_theCursor integer default dbms_sql.open_cursor;
	l_columnValue varchar2(4000);
	l_status integer;
	l_descTbl dbms_sql.desc_tab;
	l_colCnt number;
begin
	execute immediate
	'alter session set nls_date_format=''yyyy-mm-dd hh24:mi:ss'' ';

	dbms_sql.parse( l_theCursor, p_query, dbms_sql.native );
	dbms_sql.describe_columns( l_theCursor, l_colCnt, l_descTbl );

	for i in 1 .. l_colCnt loop
		dbms_sql.define_column(l_theCursor, i, l_columnValue, 4000);
	end loop;

	l_status := dbms_sql.execute(l_theCursor);

	while ( dbms_sql.fetch_rows(l_theCursor) > 0 ) 
	loop
		for i in 1 .. l_colCnt loop
			dbms_sql.column_value
			( l_theCursor, i, l_columnValue );
			dbms_output.put_line
			( rpad( l_descTbl(i).col_name, 30 )
			|| ': ' || 
			l_columnValue );
		end loop;
		dbms_output.put_line( '-----------------' );
	end loop;
	execute immediate 'alter session set nls_date_format=''dd-MON-rr'' ';
	
	exception
		when others then
		execute immediate 'alter session set nls_date_format=''dd-MON-rr'' ';
		raise;
end;
/
> set serveroutput on;
> exec print_table('select * from v$database'); -- 测试效果
```