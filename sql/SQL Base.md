## 定义
DB是存储数据的有序集合.
table是一个二维数组的集合,是存储数据和操作数据的逻辑结构.
SQL(Structured Query Language) 是用于访问和处理数据库的标准的计算机语言.
SQL= 数据定义语言 (DDL - create,drop,alter)， 数据操作语言 (DML - select,update,delete,insert)， 数据控制语言 (DCL - grant,revoke,commit,rollback)

### 运算符优先级(Operator Precedence)

- [Postgres](http://www.postgresql.org/docs/9.4/static/sql-syntax-lexical.html)
- [MySQL](https://dev.mysql.com/doc/refman/5.6/en/operator-precedence.html)

### SQL解析顺序

    FROM>WHERE>SELECT
