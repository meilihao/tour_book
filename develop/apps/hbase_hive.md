# hbase_hive

Hbase和Hive在大数据架构中处在不同位置，Hbase主要**解决实时数据查询**问题，Hive主要**解决数据处理和计算**问题，一般是配合使用.

## 区别与关系
区别:
Hbase： Hadoop database的简称，也就是基于Hadoop数据库，是一种NoSQL数据库，主要适用于海量明细数据（十亿、百亿）的随机实时查询，如日志明细、交易清单、轨迹行为等。

Hive：Hive是Hadoop数据仓库，严格来说，不是数据库，主要是让开发人员能够通过SQL来计算和处理HDFS上的结构化数据，适用于离线的批量数据计算。

	1. 通过元数据来描述Hdfs上的结构化文本数据，通俗点来说，就是定义一张表来描述HDFS上的结构化文本，包括各列数据名称，数据类型是什么等，方便我们处理数据，当前很多SQL ON Hadoop的计算引擎均用的是hive的元数据，如Spark SQL、Impala等；
	1. 基于第一点，通过SQL来处理和计算HDFS的数据，Hive会将SQL翻译为Mapreduce来处理数据；

关系:

在大数据架构中，Hive和HBase是协作关系，数据流一般如下：
1. 通过ETL工具将数据源抽取到HDFS存储
1. 通过Hive清洗、处理和计算原始数据
1. HIve清洗处理后的结果，如果是面向海量数据随机查询场景的可存入Hbase
1. 数据应用从HBase查询数据

## 备份
- [hbase，hive备份](https://blog.csdn.net/J16730zzl/article/details/113620490)