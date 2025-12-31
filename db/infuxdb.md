# infuxdb
ref:
- [InfluxData 文档](https://docs.influxdb.org.cn/)

InfluxDB 企业版由 META 节点和 DATA 节点 2 个逻辑单元组成的，而且这两个节点是 2 个单独的程序, 因为它们的场景不同:
1. META 节点存放的是系统运行的关键元信息，比如数据库（Database）、表（Measurement）、保留策略（Retention policy）等。它的特点是一致性敏感，但读写访问量不高，需要一定的容错能力

	对于 META 节点来说，节点数的多少代表的是容错能力，一般 3 个节点就可以了，因为从实际系统运行观察看，能容忍一个节点故障就可以了.

	META 节点需要强一致性，实现 CAP 中的 CP 模型, 它使用 Raft 算法实现 META 节点的一致性.
1. DATA 节点存放的是具体的时序数据。它有这样几个特点：最终一致性、面向业务、性能越高越好，除了容错，还需要实现水平扩展，扩展集群的读写性能

	对 DATA 节点而言，节点数的多少则代表了读写性能，一般而言，在一定数量以内（比如 10 个节点）越多越好，因为节点数越多，读写性能也越高，但节点数量太多也不行，因为查询时就会出现访问节点数过多而延迟大的问题.

	DATA 节点存放的是具体的时序数据，对一致性要求不高，实现最终一致性就可以了。但是，DATA 节点也在同时作为接入层直接面向业务，考虑到时序数据的量很大，要实现水平扩展，所以必须要选用 CAP 中的 AP 模型，因为 AP 模型不像 CP 模型那样采用一个算法（比如 Raft 算法）就可以实现了，也就是说，AP 模型更复杂.

## 时序数据特点
时序数据在读写、存储和分析处理方面有下面这些特点: 
1. 时序数据是持续地写入,一般是采用固定的频率,没有写入量忽大忽小的明显变化。数量非常大,而且并发写入的需求也很高。但是数据很少做更新,旧数据除了特殊情况下的修改,基本是不需要更新的写入操作 
1. 时序数据的读取很少,相比写入的高并发和高频率,读取的需求主要是进行数据分析的用,而分析应用的并发访问量是比较少的
1. 时序数据时效性很强,一般是越新的数据价值就越大,旧数据会迅速失去价值
1. 时序数据的数据分析主要关心的是新数据,旧数据被查询和分析的概率不高。旧数据一般是粗颗粒度的读取分析。而且在被分析时,一般是基于时间范围读取分析,读取某一条录的几率很小

## 概念
influxdb行协议是 InfluxDB 数据库独创的一种数据格式，它由纯文本构成，只要数据符合这种格式，就能使用 InfluxDB 的 HTTP API 将数据写入数据库.

一行数据的构成:
1. measurement: 测量名称, 类似关系型数据库中的表
1. tag set: 标签集, 可选, 一种索引类型, 标识和过滤数据, 通常用于存储维度信息, 比如地理位置, 设备id等
1. field set: 字段集, 一种数据类型, 用于存储实际的数值数据, 比如温度
1. timestamp: 时间戳, 可选, 默认精度是ns

其他:
1. bucket: 基本存储单元, 类似关系型数据库的实例, 用于组织和存储时间序列数据. 它是数据的逻辑容器, 用于区分和隔离不同类型或来源的数据
1. org: 组织概念, 用于隔离和管理用户, 应用程序或项目. 每个用户可以属于一个或多个组织, bucket属于一个org
1. InfluxQL: InfluxDB 3 现在除了 InfluxQL（一种为时间序列查询定制的类 SQL 语言）之外，还支持原生 SQL 用于查询

	InfluxDB 2.0 中引入的语言 Flux 在 InfluxDB 3 中不受支持
1. point: 类似关系型数据库的行, 包含influxdb行协议中除measurement外的所有字段
1. series: influxdb使用series来管理数据. (measurement, tag_set, 一个filed)组成一个series.

### datatype(数据类型)
ref:
- [data type](https://docs.influxdata.com/influxdb3/core/reference/glossary/#data-type)

tag始终是string, field支持:
- string
- boolean
- float (64-bit)
- integer (64-bit)
- unsigned integer (64-bit)
- time : unix时间戳

## 场景
1. 监控和运维
1. 物联网
1. 金融数据/日志分析

## 使用
InfluxDB 3.0 不再包含像 InfluxDB 2.x 那样内置仪表盘（Dashboard）功能, 可使用Grafana.

```bash
# influxdb3 create token --admin --host http://localhost:8181 # 创建管理员令牌

New token created successfully!

Token: apiv3_J1sa1xGrrGEHfJ2zdl0nKigFhol4lSi5TvOpLrKk0AUPPkz1CyngRPDZEOkjXJCEm1AYfzKk8uJaEsf0MkF0Ww
HTTP Requests Header: Authorization: Bearer apiv3_J1sa1xGrrGEHfJ2zdl0nKigFhol4lSi5TvOpLrKk0AUPPkz1CyngRPDZEOkjXJCEm1AYfzKk8uJaEsf0MkF0Ww

IMPORTANT: Store this token securely, as it will not be shown again.

# export INFLUXDB3_AUTH_TOKEN=apiv3_J1sa1xGrrGEHfJ2zdl0nKigFhol4lSi5TvOpLrKk0AUPPkz1CyngRPDZEOkjXJCEm1AYfzKk8uJaEsf0MkF0Ww
# influxdb3 show tokens # 查看tokens
# influx -version
# influx
# > show databases
# > use <database_name>
# > show measurements: 查看数据库中的表
# > show tag keys from <measurement_name>: 查看指定表的标签键
# > show field keys from <measurement_name>: 查看指定表的字段键
# > select ua from device_prop where time >='2025-08-11 22:50:00' and time <='2025-08-11 23:10:00' and uuid='1757904264136532677'; # influxdb time都是基于UTC
```

查询思路:
1. 指定bucket
1. 指定数据的时间范围
1. 指定(measurement, tag_set, 一个filed)

ps: 可以将measurement, tag_set, filed, 时间视为索引

## data一致性
1. 自定义副本数
1. Hinted-handoff
	
	解决场景: 一个节点接收到写请求时，需要将写请求中的数据转发一份到其他副本所在的节点, 在这个过程中，远程 RPC 通讯是可能会失败的.

	Hinted-handoff 实现:
	1. 写失败的请求，会缓存到本地硬盘上
	1. 周期性地尝试重传
	1. 相关参数信息，比如缓存空间大小 (max-szie)、缓存周期（max-age）、尝试间隔（retry-interval）等，是可配置的
1. 反熵
1. Quorum NWR

## FAQ
### inflexdb 1.6 select时将time转成可读时间
```bash
# apt install python3-influxdb
# vim s.py
import pandas as pd
from influxdb import InfluxDBClient

# 设置 pandas 显示选项
pd.set_option('display.max_rows', 1000)  # 最大行数
pd.set_option('display.max_columns', 100)  # 最大列数
pd.set_option('display.width', None)  # 自动调整宽度
pd.set_option('display.max_colwidth', None)  # 显示完整列内容

client = InfluxDBClient(host='localhost', port=8086, database='test')

results = client.query("SELECT * FROM device_prop where uuid='xxx' order by time desc LIMIT 100")

df = pd.DataFrame(list(results.get_points()))

# 将utc时间戳转换为可读时间格式
df['time'] = pd.to_datetime(df['time']) 
df['time_local'] = df['time'].dt.tz_convert('Asia/Shanghai') # 如果原始数据是 UTC 时间，转换为本地时间
# df['time_local'] = df['time'].dt.tz_localize('Asia/Shanghai') # 如果原始数据没有时区（naive），直接设成本地时间

print(df)
# python3 s.py
```

### `SELECT * FROM device_prop limit 1`高cpu占用
加其他where条件