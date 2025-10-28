# VictoriaMetricsStack
# VictoriaMetrics
ref:
- [VictoriaMetrics 中文教程（10）集群版介绍](https://www.cnblogs.com/ulricqin/p/18512032)

	如果数据量低于每秒一百万个数据点，建议使用单节点版本，而不是集群版本

	VictoriaMetrics 支持 replication，挂掉部分节点不影响数据安全，不过，建议不要开启，而是交由云盘等底层存储来保证数据的持久性。如果没有条件使用云存储，再考虑开启 replication

VictoriaMetrics 是一种快速、经济高效且可扩展的监控解决方案和时间序列数据库, 其目标是替换 Prometheus.

除了单节点 VictoriaMetrics之外，VictoriaMetrics 生态系统还包含以下组件：
- vmagent - 轻量级代理，用于通过 PULL 和 PUSH 的协议接收指标，将其转换并发送到已配置的与 Prometheus 兼容的远程存储系统，例如 VictoriaMetrics
- vmalert - 用于处理与 Prometheus 兼容的警报和记录规则的服务
- vmalert-tool - 用于验证警报和记录规则的工具
- vmauth - 针对 VictoriaMetrics 产品优化的授权代理和负载均衡器
- vmgateway - 具有每个租户速率限制功能的授权代理
- vmctl - 用于在不同存储系统之间迁移和复制数据以获取指标的工具
- vmbackup 、 vmrestore 和 vmbackupmanager - 用于为 VictoriaMetrics 数据创建备份和从备份恢复的工具
- vminsert 、 vmselect 和 vmstorage —— VictoriaMetrics 集群的组件

	组件名称,角色（功能）,主要职责
	vminsert,写入层 (Insert),"接收来自 Prometheus Remote Write 或其它协议（如 InfluxDB line protocol, OpenTSDB 等）的指标数据。它负责将数据分发到集群中的多个 vmstorage 节点
	vmstorage,存储层 (Storage),实际存储指标数据块（TSDB），并在本地执行数据压缩、去重和清理（Retention）。它是集群的持久化层
	vmselect,查询层 (Select),负责接收用户的查询请求（如 PromQL）。它从集群中的所有 vmstorage 节点并行获取所需的数据，执行查询逻辑（如聚合、连接），并将最终结果返回给用户
- VictoriaLogs - 用户友好、经济高效的日志数据库
- VictoriaTrace - 代替jaeger

## victoria-metrics-prod
参数:
- -storageDataPath : VictoriaMetrics 将所有数据存储在此目录中, 默认路径是当前工作目录中的 victoria-metrics-data
- -retentionPeriod - 存储数据的保留期限, 较旧的数据将被自动删除. 默认保留期为 1 个月（31 天）, 最短保留期限为24小时或1天.


## FAQ
### victoria-metrics-prod和vmagent怎么都有"--promscrape.config"参数?
参数所在组件,作用目的,配置内容,数据流向
vmagent,主要数据采集（外部指标）,业务目标的抓取配置 (scrape_configs),外部应用 → vmagent → vminsert
victoria-metrics-prod (VM Single/Cluster),组件自监控（内部指标）,自身 /metrics 端点的抓取配置,vmstorage 自身 → vmstorage 数据库

### victoria-metrics-prod和 vminsert 、 vmselect 和 vmstorage的关系
victoria-metrics-prod 这个名称通常指的是 VictoriaMetrics 单体 (Single-Node) 版本，而 vminsert、vmselect 和 vmstorage 是 VictoriaMetrics 集群 (Cluster) 版本 的核心组件

特性,VM 单体版 (victoria-metrics-prod=vminsert + vmstorage + vmselect),VM 集群版 (vminsert + vmstorage + vmselect)
组件数量,1 个进程,至少 3 个进程（通常更多）
高可用性,无（单点故障）,有（可通过部署多副本实现 HA）
可扩展性,垂直扩展（增加 CPU/内存/磁盘）,水平扩展（增加 vminsert/vmselect/vmstorage 节点）
数据流,数据进入 → 进程内处理 → 存储到本地,vminsert → 网络分发 → vmstorage (存储) ← 网络查询 ← vmselect