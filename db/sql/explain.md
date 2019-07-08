# explain
参考:
- [如何分析一条 SQL 的性能](https://www.tuicool.com/articles/qaA7Jju)
- [数据库允许空值(null)，往往是悲剧的开始](https://mp.weixin.qq.com/s?__biz=MjM5ODYxMDA5OQ==&mid=2651962495&idx=1&sn=74e9e0dc9d03a872fd5bce5769f6c22a&chksm=bd2d09a38a5a80b50da3b67c03da8417426cbb201427557959fa91e9094a848a14e0db214370&scene=21#wechat_redirect)

大多数情况下都是从 RDMS 的 慢查询日志 中揪出来一些查询效率比较慢的 sql 来使用 explain 分析.

## mysql
explain 返回的字段:
- type : mysql 访问数据的方式
    - all : 全表扫描
    - index : 遍历索引
    - range : 索引区间查询
    - ref : 普通索引等值扫描, 对于前表的每一行 (row), 后表可能有多于一行的数据被扫描
    - eq_ref : pk或者unique(not null)上的join查询，等值匹配. 对于前表的每一行 (row) ，后表只有一行被扫描
    - const : pk或者unique上的等值查询
    - system : 系统表, 少量数据，往往不需要进行磁盘IO
    
    效率从最好到最差的一个排序:`system > const > eq_ref > ref > range > index > all`
- key : 查询过程实际会用到的索引名称
- rows : 查询过程中可能需要扫描的行数，这个数据不一定准确，是mysql 抽样统计的一个数据
- Extra : 一些额外的信息，通常会显示是否使用了索引，是否需要排序，是否会用到临时表等