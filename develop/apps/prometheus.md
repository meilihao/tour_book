# prometheus
监控的客户:
- 技术

  定位问题
- 业务

  保障业务可持续

监控目的:
- 趋势分析
- 对照分析
- 告警
- 故障分析和定位
- 数据可视化

监控层次划分:
- 端(即末端)监控

  针对用户体验进行监控, 比如页面加载数据, api调用成功率. 阿里有SPM(super position model, 超级位置模型), SCM(super content model, 超级内容模型), 黄金令箭等理论支撑的端监控实践.

- 业务层监控

  比如qps, 日活等

- 应用层监控

  链路追踪
- 中间件监控
- 系统层监控

监控应用程序主要有两种方法： 探针（probing, 从应用外部检查其状态） 和内省（introspection, 应用自身提供其状态, 内容相对比探针丰富）.

prometheus+grafana是一站式通用监控架构的最佳方案之一.

> CNCF的Observability and Analysis(监控和可观察性)囊括了主流的监控解决方案.

opentelemetry监控三要素:
- metrics

  它是某段时间内发生的可聚合(aggregatable)观察点(observation)的集合.

  观察点通常包括值、 时间戳， 有时也涵盖描述观察点的一系列属性（如源或标签）. 观察的集合称为时间序列.

  指标类型:
  - 测量型: 本质上是特定度量的快照. 比如cpu利用率
  - 计数型: 随着时间增加而不会减少的数字, 但可能被重置. 比如设备收发包的字节数或登录次数
  - 直方图: 对观察点进行采样的指标类型， 可以展现数据集的频率分布
- logging
- tracing

通过metrics->tracing->logging顺序分析问题, 通常比直接查log更高效.

prometheus遵从MDD思想.

> MDD(metrics-driven development)主张整个应用开发过程由指标驱动, 通过实时指标来驱动快速, 精确和细粒度的软件迭代.

分布式环境一般存在服务发现, 此时监控使用pull更方便, 因为push的话需要配置监控参数(其实通过服务发现也可自动配置参数).

> 当前阿里使用了prometheus.

## 架构
Prometheus通过抓取或拉取应用程序中暴露的时间序列数据来工作. 时间序列数据通常由应用程序本身通过客户端库或称为exporter的代理来作为HTTP端点暴露.

Prometheus称其可以抓取的指标来源为端点（endpoint）, 端点通常对应单个进程、 主机、 服务或应用程序.

为了抓取端点数据, Prometheus定义了名为目标（target） 的配置, 这是执行抓取所需的信
息. 一组目标被称为作业（job）, 作业通常是具有相同角色的目标组.

prometheus组成:
1. prometheus server

    prometheus核心, 主要功能包括:
    1. 抓取

        周期性从job, exporter, pushgateway获取metrics
    1. 存储

        抓取到的监控数据通过一定的规则清理和数据整理(抓取前使用服务发现提供的relabel_configs方法筛选metric, 抓取后使用作业内的metrics_relabel_configs方法筛选metric)
    1. 查询.
1. pushgateway

    prometheus是以pull为主的监控系统, 可通过pushgateway实现push. 它本质是prometheus server无法抓取的资源的解决方案.
1. job/exporter

    是prometheus target, 是prometheus监控的对象
1. service discovery

    通过它发现agent. 通过它可实现不重启prometheus server的情况下动态发现target
1. alertmanager

    独立于prometheus的一个告警组件.
1. dashboard

> prometheus cluster可使用thanos.

## 数据模型
Prometheus时序数据的数据结构：
1. 指标名称（Metric Name） ：描述数据的类型，例如 node_cpu_seconds_total
1. 标签（Labels） ：以键值对的形式为数据提供上下文，例如 {cpu="0", mode="user", instance="192.168.1.10:9100"}
1. 样本值（Sample Value） ：实际的数值数据，通常为浮点数，例如 1345.67
1. 时间戳（Timestamp） ：每个数据点都有一个关联的时间戳，例如 1701409200000

时间序列数据模型结合了时间序列名称和称为标签（label） 的键/值对, 这些标签提供了维度. 每个时间序列由时间序列名称和标签的组合唯一标识.

[时间序列名称](https://prometheus.io/docs/practices/naming/)通常描述收集的时间序列数据的一般性质, 比如total_website_visits为网站访
问的总数. 名称可以包含ASCII字符、 数字、 下划线和冒号.

标签为Prometheus数据模型提供了维度。 它们为特定时间序列添加上下文. 比如total_website_visits时间序列可以使用能够识别网站名称、 请求IP或其他特殊标识的标签.

标签共有两大类： 插桩标签（instrumentation label） 和目标标签（target label）. 插桩标签来自被
监控的资源. 目标标签由Prometheus在抓取期间和之后添加.

可通过`http://localhost:9090/graph`内置web服务界面查看指标.

relable时机:
1. 对来自服务发现的目标进行重新标记

  这对于将来自服务发现的元数据标签中的信息应用于指标上的标签来说非常有用
1. 在抓取之后且指标被保存于存储系统之前

> 在抓取之前使用relabel_configs; 在抓取之后使用metric_relabel_configs

## 服务发现
访问`https://localhost:9090/service-discovery`, 可通过Web界面查看服务发现目标的完整列
表及其元数据标签

- 用户提供的静态资源列表
- 基于文件的发现. 比如使用配置管理工具生成在Prometheus中可以自动更新的资源列表

  用file_sd_configs块替换prometheus.yml文件中的static_configs块
- 自动发现. 比如查询Consul等数据存储， 在Amazon或Google中运行实例， 或使用DNS SRV记录
来生成资源列表

  比如发现amazon ec2的ec2_sd_configs, 基于dns发现的dns_sd_configs.

## 参数
- --config.file=prometheus.yml : 指定配置文件
- --web.listen-address=:9090 : 指定web访问端口, 此时必须指定`--config.file=prometheus.yml`
- --web.enable-lifecycle : 可通过http管理prometheus
    
    - `curl -X POST http://localhost:9090/-/reload` : 触发重载prometheus.yml
    - `curl -X POST http://localhost:9090/-/quit` : prometheus退出

- --web.external-url : 使用外部url和代理
- --web.route-prefix
- --storage.tsdb.path: 时间序列保存位置, 建议使用ssd
- --storage.tsdb.retention.time: 时间序列保留时间
- --storage.tsdb.retention.size: 时间序列最大大小

## PromQL
- process_resident_memory_bytes : Prometheus进程的内存使用情况
- `{quantile="0.25"}` # 带`quantile="0.25"`标签的指标
- `go_gc_duration_seconds{quantile="0.25"}` # 时间序列名称是go_gc_duration_seconds且带`quantile="0.25"`标签的指标
- `go_gc_duration_seconds{quantile!="0.25"}`
- promhttp_metric_handler_requests_total: Prometheus服务器中抓取数据所产生的HTTP请
求总数, 结果是按resp code分类
- `sum(promhttp_metric_handler_requests_total)`: 汇总promhttp_metric_handler_requests_total
- `sum(promhttp_metric_handler_requests_total) by (job)`: 按job汇总promhttp_metric_handler_requests_total
- `rate(promhttp_metric_handler_requests_total[5m])`: 5m内promhttp_metric_handler_requests_total的增长率

  rate()函数用来计算一定范围内时间序列的每秒平均增长率， 只能与计数器一起使用, 最适合用于增长较慢的计数器或用于警报的场景.

  时间单位:
  - s : 秒
  - m : 分钟
  - h : 小时
  - d : 天
  - w : 周

## prometheus.yml
```yaml
# my global config
global: # 全局配置信息. 它定义的内容可被scrape_configs中的每个job单独覆盖
  scrape_interval: 15s # Set the scrape interval to every 15 seconds. Default is every 1 minute. 抓取数据间隔
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s). 指定Prometheus评估规则的频率

# Alertmanager configuration
alerting: # 设置Prometheus的警报
  alertmanagers:
    - static_configs: # 静态配置alertmanager的地址, 也可依赖服务发现动态识别比如dns_sd_configs
        - targets:
          # - alertmanager:9093

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files: # 指定包含记录规则或警报规则的文件列表
  # - "first_rules.yml"
  # - "second_rules.yml"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs: # 指定Prometheus抓取的所有目标
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: "prometheus"

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs: # 设置抓取目标
      - targets: ["localhost:9090"] # Prometheus服务器自身

  - job_name: 'otel-collector'
    scrape_interval: 10s
    static_configs:
    - targets: ['0.0.0.0:8888']
  
  - job_name: 'clients'
    scrape_interval: 10s
    static_configs:
    - targets: ['0.0.0.0:8889']

  - job_name: 'node'
    scheme: http
    metrics_path: /metrics
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9100']
    params:
      collect[]: # 仅抓取cpu, 即curl http://localhost:9100/metrics?collect[]=cpu
        - cpu
  - job_name: 'xxx'
    static_configs:
      - targets: ['localhost:9101']
    metric_relabel_configs:
      - source_labels: [__name__] # 删除指标, __name__是表示指标名称的预留标签
        regex: 'node_netstat_Icmp_OutMsgs' # re2语法
        action: drop
      - source_labels: [pod] # 将所有的 metrics 中名为`pod`的label 值拷贝到名为`pod_name`的label中, 即pod_name=$1
        separator: ;
        regex: (.+)
        target_label: pod_name
        replacement: $1
      - regex: 'kernelVersion' # 删除标签
        action: labeldrop
```

> honor_labels: 默认false, 通过在其前面添加exported_前缀来重命名现有标签, 否则忽略覆盖已存在的标签

> 重启Prometheus服务器或进行SIGHUP可reload `rule_files`

> scrape_configs job 也支持file_sd_configs

目前主要有两种规则:
- 记录规则（recording rule）

  记录规则是一种根据已有时间序列计算新时间序列（特别是聚合时间序列） 的方法, 比如预先计算使用频繁且开销大的表达式, 并将结果保存为一个新的时间序列数据

- 警报规则（alerting rule）

  允许定义警报条件

## promql
prometheus默认查询地址是`localhost:9090/graph`.

promql命令:
- up : 查看每个监控job的监控状态: 1, ok; 0, bad.

## 高可用
ref:
- [高可用prometheus集群方案选型分享](https://developer.aliyun.com/article/993139)
- [高可用Prometheus：Thanos 实践](https://www.kubernetes.org.cn/7217.html)

Prometheus最简单的容错解决方案是并行运行两个配置相同的Prometheus服务器， 并且这两个服务器同时处于活动状态. 该配置生成的重复警报可以交由上游Alertmanager使用其分组（及抑制） 功能进行处理.

## exporter
### node_exporter
```bash
# ./node_exporter --version
# ./node_exporter --web.listen-address=:9100
```

`[no-]xxx`选项表示关闭相应的收集器

textfile收集器: node_exporter通过扫描指定目录中的文件，提取所有格式为Prometheus指标的字符串， 然后暴露它们以便抓取. 它使用`--
collector.textfile.directory`指定

### 从日志提取信息
mtail

### Blackbox exporter
探针, 允许通过HTTP、 HTTPS、 DNS、 TCP和ICMP来探测端点