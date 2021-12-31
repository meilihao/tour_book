# prometheus
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

prometheus+grafana是一站式通用监控架构的最佳方案之一.

> CNCF的Observability and Analysis(监控和可观察性)囊括了主流的监控解决方案.

opentelemetry监控三要素:
- metrics

  它是某段时间内发生的可聚合(aggregatable)数据点的集合.
- logging
- tracing

通过metrics->tracing->logging顺序分析问题, 通常比直接查log更高效.

prometheus遵从MDD思想.

> MDD(metrics-driven development)主张整个应用开发过程由指标驱动, 通过实时指标来驱动快速, 精确和细粒度的软件迭代.

分布式环境一般存在服务发现, 此时监控使用pull更方便, 因为push的话需要配置监控参数(其实通过服务发现也可自动配置参数).

> 当前阿里使用了prometheus.

## 架构
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

## 参数
- --config.file=prometheus.yml : 指定配置文件
- --web.listen-address=:9090 : 指定web访问端口, 此时必须指定`--config.file=prometheus.yml`
- --web.enable-lifecycle : 可通过http管理prometheus
    
    - `curl -X POST http://localhost:9090/-/reload` : 触发重载prometheus.yml
    - `curl -X POST http://localhost:9090/-/quit` : prometheus退出

- --web.external-url : 使用外部url和代理
- --web.route-prefix 

## prometheus.yml
```yaml
# my global config
global: # 全局配置信息. 它定义的内容可被scrape_configs中的每个job单独覆盖
  scrape_interval: 15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

# Alertmanager configuration
alerting:
  alertmanagers:
    - static_configs: # 静态配置alertmanager的地址, 也可依赖服务发现动态识别
        - targets:
          # - alertmanager:9093

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
  # - "first_rules.yml"
  # - "second_rules.yml"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: "prometheus"

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs: # 静态配置alertmanager的地址, 也可依赖服务发现动态识别
      - targets: ["localhost:9090"]

  - job_name: 'otel-collector'
    scrape_interval: 10s
    static_configs:
    - targets: ['0.0.0.0:8888']
  
  - job_name: 'clients'
    scrape_interval: 10s
    static_configs:
    - targets: ['0.0.0.0:8889']

  - job_name: 'node'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9100']
```

## promql
prometheus默认查询地址是`localhost:9090/graph`.

promql命令:
- up : 查看每个监控job的监控状态: 1, ok; 0, bad.