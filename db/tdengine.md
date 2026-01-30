# TDengine
ref:
- [TDengine 发布历史及下载链接](https://docs.taosdata.com/releases/tdengine/)
- [深入解析：工业时序数据库TDengine 架构、性能与实战全解析](https://www.cnblogs.com/jzssuanfa/p/19306344)

书:
1. 时序大数据平台TDengine核心原理与实战: [即TDengine TSDB 文档的精简版](https://docs.taosdata.com/)

## 概念
### datatype(数据类型)
ref:
- [数据类型](https://docs.taosdata.com/reference/taos-sql/datatype/)

### 组件
ref:
- [组件介绍](https://docs.taosdata.com/operation/intro/)

1. taosd: tdengine数据库引擎

  负责处理所有与数据相关的操作，包括数据写入、查询和管理等
1. taosAdapter: 应用和tdengine之间的桥梁

  支持用户通过 RESTful 接口和 WebSocket 连接访问 TDengine TSDB 服务，实现数据的便捷接入和处理.

  taosAdapter 能够与各种数据收集代理工具（如 Telegraf、StatsD、collectd 等）无缝对接，从而将数据导入 TDengine TSDB。此外，它还提供了与 InfluxDB/OpenTSDB 兼容的数据写入接口，使得原本使用 InfluxDB/OpenTSDB 的应用程序能够轻松移植到 TDengine TSDB 上，无须进行大量修改

1. taosKeeper: tdengine监控指标的导出工具

  旨在方便用户对 TDengine TSDB 的运行状态和性能指标进行实时监控
1. taosX, taosX agent: 数据管道工具, 仅tdengine enterprise提供

  1. taosX: 旨在为用户提供一种无须编写代码即可轻松对接第三方数据源的方法，实现数据的便捷**导入**
  1. taosX agent: 与 taosX 协同工作，负责接收 taosX 下发的外部数据源导入任务。taosX Agent 能够启动连接器或直接从外部数据源获取数据，随后将采集的数据转发给 taosX 进行处理

  在边云协同场景中，taosX Agent 通常部署在边缘侧，尤其适用于那些外部数据源无法直接通过公网访问的情况。通过在边缘侧部署 taosX Agent，可以有效地解决网络限制和数据传输延迟等问题，确保数据的实时性和安全性
1. taosExploer: 可视化图形管理工具
1. taosc: tdengine客户端驱动

![TDengine TSDB 产品生态的拓扑架构](https://docs.taosdata.com/assets/images/tdengine-topology-f603c27e26f6e8162c0ad79349ea8771.png)

[通过与各类应用程序、可视化和 BI（Business Intelligence，商业智能）工具以及数据源集成，TDengine TSDB 为用户提供了灵活、高效的数据处理和分析能力，以满足不同场景下的业务需求](https://docs.taosdata.com/operation/intro/).

## 开发
ref:
- [开发指南](https://docs.taosdata.com/develop/)
- [TDengine TSDB Go Connector](https://docs.taosdata.com/reference/connector/go/)
- [第三方工具](https://docs.taosdata.com/third-party/)
- [技术内幕](https://docs.taosdata.com/tdinternal/)
- [实践案例](https://docs.taosdata.com/application/)

对于 WebSocket 连接和原生连接，连接器都提供了相同或相似的 API 操作数据库，只是初始化连接的方式稍有不同，用户在使用上不会感到什么差别.

推荐使用websocket.

### 数据模型
一张超级表至少包含一个时间戳列、一个或多个采集量列以及一个或多个标签列.

在 TDengine TSDB 中，表代表具体的数据采集点，而超级表则代表一组具有相同属性的数据采集点集合.

子表是数据采集点在逻辑上的一种抽象表示，它是隶属于某张超级表的具体表。用户可以将超级表的定义作为模板，并通过指定子表的标签值来创建子表。这样，通过超级表生成的表便被称为子表。超级表与子表之间的关系主要体现在以下几个方面:
1. 一张超级表包含多张子表，这些子表具有相同的表结构，但标签值各异
1. 子表的表结构不能直接修改，但可以修改超级表的列和标签，且修改对所有子表立即生效
1. 超级表定义了一个模板，自身并不存储任何数据或标签信息

在 TDengine TSDB 中，查询操作既可以在子表上进行，也可以在超级表上进行。针对超级表的查询，TDengine TSDB 将所有子表中的数据视为一个整体，首先通过标签筛选出满足查询条件的表，然后在这些子表上分别查询时序数据，最终将各张子表的查询结果合并。本质上，TDengine TSDB 通过对超级表查询的支持，实现了多个同类数据采集点的高效聚合

## 部署
ref:
- [TDengine Download Center](https://tdengine.com/downloads)
- [快速体验](https://docs.taosdata.com/get-started/)
- [安装部署](https://docs.taosdata.com/operation/install/) + [集群部署](https://docs.taosdata.com/operation/deployment/)
- [多级存储与对象存储](https://docs.taosdata.com/operation/multi/)

docker:
```bash
# docker run -d \
  -v $(pwd)/data/taos/dnode/data:/var/lib/taos \
  -v $(pwd)/data/taos/dnode/log:/var/log/taos \
  -p 6030:6030 -p 6041:6041 -p 6043:6043 -p 6060:6060 \
  -p 6044-6049:6044-6049 \
  -p 6044-6045:6044-6045/udp \
  -p 6050:6050 -p 6055:6055 \
  --name tsdb-demo \
  tdengine/tsdb-ee
# docker exec -it tsdb-demo taos
taos> show databases;
> show variables; # 查看和修改变量时需要注意是全局还是局部变量
> show dnodes;
> show dnode 1 variables;
> show dnode 1 variables like '%git%';
> SHOW VARIABLES like '%time%';
```

访问: http://localhost:6060 root/taosdata

## sql
```bash
> SELECT _wstart, _wend, sum(annual_bilateral) FROM dee.ml_1384 where ts > '2025-12-15T00:00:01.000+08:00' INTERVAL(1d, 57601s); -- 按天具体数据, 当天零点属于前一天. 57601s = 16h1s, 但tdengine不支持"16h1s"
> SELECT sum(annual_bilateral) FROM dee.ml_1384 where ts > '2025-12-31T00:00:00.000+08:00' and ts <= '2026-01-01T00:00:00.000+08:00'
```

可用函数:
- to_iso8601(ts)

## FAQ
### taos显示不完全
解决方法:
1. `set max_binary_display_width 1000;`, 设置一个长度,立即生效
2. vim taos.cfg 文件中的设置项`maxBinaryDisplayWidth:1000`

### interval
GROUP BY time(1d), tsdb 会从 UTC 时间的凌晨 00:00:00 开始对齐, 而GROUP BY time(1d, 16h), 是utc的[... 16:00:00Z, ... 15:59:59Z) -> 北京时间[... 00:00:00+8:00, ... 23:59:59+8:00]