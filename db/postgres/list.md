## base
- [30个实用SQL语句，玩转PostgreSQL](https://mp.weixin.qq.com/s?__biz=Mzg3MjA5OTkzMw==&mid=2247484562&idx=1&sn=6774d5e3131fbc74a4f2ed9df03ca5fd)
- [An Overview of PostgreSQL Indexes](https://www.enterprisedb.com/postgres-tutorials/overview-postgresql-indexes)

## 进阶
- [PostgreSQL黑科技大集会](https://yq.aliyun.com/articles/2727)
- [PostgreSQL on Linux 最佳部署手册](http://mp.weixin.qq.com/s/FR65pyRmpEFFVvoJ28uBUg)
- [Postgres 索引类型探索之旅](https://linux.cn/article-9035-1.html)
- [PostgreSQL数据目录深度揭秘](https://www.tuicool.com/articles/aiYZny6)

### 相关主页

- [开源知识库](http://code.csdn.net/openkb/p-PostgreSQL)
- [postgres常用命令](http://developer.51cto.com/art/201401/426180.htm)
- [PostgreSQL 9.3.1 中文手册](http://www.postgres.cn/docs/9.3/index.html)
- [PostgreSQL手册](http://pgsqlcn.com)

## ha
- [PostgreSQL部署|基于Stream复制的手工主从切换](https://www.modb.pro/db/404682)
- [PostgreSQL流复制之主备切换](https://www.modb.pro/db/235078)
- [pacemaker+drbd+postgres](https://www.insight-ltd.co.jp/tech_blog/postgresql/440/)
- [PolarDB for PostgreSQL高可用原理](https://developer.aliyun.com/article/789048)
- [Pigsty : pg集群方案](https://www.oschina.net/news/197066/pigsty-1-5-released)
- [Patroni + Etcd + PostgreSQL 部署集群](https://www.modb.pro/db/107608)
- [基于Patroni的PostgreSQL高可用环境部署](https://developer.aliyun.com/article/775029)

	含"自动切换和脑裂防护"说明
- [PostgresSQL HA高可用架构实战](https://blog.51cto.com/u_14977574/2548233)

	基于Corosync +Pacemaker, [resource-agents/heartbeat/pgsql](https://github.com/ClusterLabs/resource-agents/blob/main/heartbeat/pgsql)
- [PG高可用之主从流复制+keepalived 的高可用](https://bbs.huaweicloud.com/blogs/330678)
- [PostgreSQL如何保障数据的一致性](https://chenhuajun.github.io/2017/09/02/PostgreSQL%E5%A6%82%E4%BD%95%E4%BF%9D%E9%9A%9C%E6%95%B0%E6%8D%AE%E7%9A%84%E4%B8%80%E8%87%B4%E6%80%A7.html)

	PG通过synchronous_commit参数设置复制的持久性级别, 下面这些级别越往下越严格，从remote_write开始就可以保证单机故障不丢数据了:
    - off
    - local
    - remote_write
    - on
    - remote_apply

	每个级别的含义参考手册: [19.5. 预写式日志](ttp://www.postgres.cn/docs/9.6/runtime-config-wal.html#RUNTIME-CONFIG-WAL-SETTINGS)或[Evolution of Fault Tolerance in PostgreSQL: Synchronous Commit](https://www.2ndquadrant.com/en/blog/evolution-fault-tolerance-postgresql-synchronous-commit/)

	> [整体来说MySQL的日志同步上还是没有PostgreSQL做的严谨，在金融系统中，PostgreSQL的日志同步级别都是设置为on，即日志接收到，apply，然后等待数据刷盘才返回commit的ack](https://www.cnblogs.com/kuang17/p/11331969.html)

	> MySQL通过半同步复制在很大程度上降低了failover丢失数据的概率。MySQL的主库在等待备库的应答超时时半同步复制会自动降级成异步，此时发生failover会丢失数据
- [PostgreSQL复制槽实操](https://www.modb.pro/db/29737)

## 发行版
- [Pigsty 开箱即用的 PostgreSQL 发行版](https://www.oschina.net/p/pigsty)/[Pigsty 近况与 v1.4 前瞻](https://www.oschina.net/news/184665)/[Vonng/pigsty](https://github.com/Vonng/pigsty)
- [pg支持周期](https://www.postgresql.org/support/versioning/)
