# prometheus & grafana
![架构](https://prometheus.io/assets/architecture.png)

参考:
- [prometheus-book](https://yunlzheng.gitbook.io/prometheus-book/)

# prometheus
Prometheus是一款开源的业务监控和时序数据库.

## 基本概念
Prometheus定义了4中不同的指标类型(metric type)：
- Counter（计数器）: 只增不减的计数器

  Counter类型的指标其工作方式和计数器一样,只增不减（除非系统发生重置）. 常见的监控指标,如http_requests_total,node_cpu都是Counter类型的监控指标. 一般在定义Counter类型指标的名称时推荐使用_total作为后缀. Counter是一个简单但有强大的工具,例如我们可以在应用程序中记录某些事件发生的次数,通过以时序的形式存储这些数据,我们可以轻松的了解该事件产生速率的变化. PromQL内置的聚合操作和函数可以用户对这些数据进行进一步的分析：
  - 通过rate()函数获取HTTP请求量的增长率：`rate(http_requests_total[5m])`
  - 查询当前系统中,访问量前10的HTTP地址：`topk(10, http_requests_total)`

- Gauge（仪表盘）: 可增可减的仪表盘

  与Counter不同,Gauge类型的指标侧重于反应系统的当前状态. 因此这类指标的样本数据可增可减. 常见指标如：node_memory_MemFree（主机当前空闲的内容大小）,node_memory_MemAvailable（可用内存大小）都是Gauge类型的监控指标. 通过Gauge指标,用户可以直接查看系统的当前状态： node_memory_MemFree 对于Gauge类型的监控指标, 通过PromQL内置函数delta()可以获取样本在一段时间返回内的变化情况. 例如,计算CPU温度在两个小时内的差异： delta(cpu_temp_celsius{host="zeus"}[2h]) 还可以使用deriv()计算样本的线性回归模型,甚至是直接使用predict_linear()对数据的变化趋势进行预测.例如,预测系统磁盘空间在4个小时之后的剩余情况： predict_linear(node_filesystem_free{job="node"}[1h], 4 * 3600)

- Histogram（直方图）/ Summary（摘要）: 使用Histogram和Summary分析数据分布情况

  Histogram和Summary主用用于统计和分析样本的分布情况. 在大多数情况下人们都倾向于使用某些量化指标的平均值, 例如CPU的平均使用率, 页面的平均响应时间. 这种方式的问题很明显,以系统API调用的平均响应时间为例：如果大多数API请求都维持在100ms的响应时间范围内,而个别请求的响应时间需要5s,那么就会导致某些WEB页面的响应时间落到中位数的情况,而这种现象被称为长尾问题. 为了区分是平均的慢还是长尾的慢,最简单的方式就是按照请求延迟的范围进行分组. 例如,统计延迟在0~10ms之间的请求数有多少而10~20ms之间的请求数又有多少.通过这种方式可以快速分析系统慢的原因.Histogram和Summary都是为了能够解决这样问题的存在,通过Histogram和Summary类型的监控指标, 可以快速了解监控样本的分布情况. 例如,指标prometheus_tsdb_wal_fsync_duration_seconds的指标类型为Summary. 它记录了Prometheus Server中wal_fsync处理的处理时间.

  与Summary类型的指标相似之处在于Histogram类型的样本同样会反应当前指标的记录的总数(以_count作为后缀)以及其值的总量（以_sum作为后缀）. 不同在于Histogram指标直接反应了在不同区间内样本的个数,区间通过标签len进行定义. 同时对于Histogram的指标,我们还可以通过histogram_quantile()函数计算出其值的分位数. 不同在于Histogram通过histogram_quantile函数是在服务器端计算的分位数. 而**Sumamry的分位数则是直接在客户端计算完成**.因此对于分位数的计算而言,Summary在通过PromQL进行查询时有更好的性能表现,而Histogram则会消耗更多的资源. 反之对于客户端而言Histogram消耗的资源更少.在选择这两种方式时用户应该按照自己的实际场景进行选择.

## 部署prometheus
非容器:
```bash
# cd prometheus-2.47.2
# ./prometheus --version
# mkdir -p /etc/prometheus
# cp prometheus.yml /etc/prometheus
# ./promtool check config /etc/prometheus/prometheus.yml # 检查配置文件
# ./prometheus --config.file "/etc/prometheus/prometheus.yml"
```

容器:
```bash
# mkdir -p /var/lib/prometheus
# chmod 777 /var/lib/prometheus
# cat << EOF |tee /etc/prometheus/prometheus.yml
# my global config
global:
  scrape_interval:     15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

# Alertmanager configuration
alerting:
  alertmanagers:
  - static_configs:
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
  - job_name: 'prometheus'

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
    - targets: ['localhost:9090']
  - job_name: 'otel-collector'
    scrape_interval: 10s
    static_configs:
    - targets: ['0.0.0.0:8888']
  - job_name: 'clients'
    scrape_interval: 10s
    static_configs:
    - targets: ['0.0.0.0:8889']
EOF
# docker volume create data-prometheus
# docker run -d --restart=unless-stopped -p 9090:9090 -v /etc/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml -v data-prometheus:/prometheus --name prometheus prom/prometheus --config.file=/etc/prometheus/prometheus.yml --storage.tsdb.retention.time=15d storage.tsdb.retention.size=8GB
```

> prometheus.yml可通过不指定`-v /etc/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml`时进入容器来获取 

访问http://localhost:9090

systemd demo:
```bash
# cat > /etc/sysconfig/prometheus << EOF
OPTIONS="--config.file=/opt/prometheus/prometheus.yml"
EOF
# cat > /lib/systemd/system/prometheus.service << EOF
[Unit]
Description=prometheus server daemon
Documentation=https://prometheus.io/
After=network.target

[Service]
Type=simple
User=root
ExecStart=/opt/prometheus/prometheus $OPTIONS
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF
```

## 部署node_exporter
> 默认已开启pprof

> node_exporter-1.8.2.linux-arm64解析fibrechannel遇到出现极高cpu占用, 卡住`/metrics`, 仅用该collector后正常 by `go tool pprof http://localhost:9100/debug/pprof/profile?seconds=10`

```
# ./node_exporter
```

访问: http://${ip}:9100/metrics即可.

编辑prometheus.yml并在scrape_configs节点下添加以下内容, 再重启prometheus即可:
```yaml
scrape_configs:
  # 采集node exporter监控数据
  - job_name: 'node'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9100']
```

此时访问http://localhost:9090，进入到Prometheus Server, 选择顶部导航栏的 Status --> Targets 中可以看到多了一个新的名为"node"的job且State为"Up"即表示添加job成功.

systemd部署见[node_exporter.service](https://github.com/prometheus/node_exporter/blob/master/examples/systemd/node_exporter.service)

systemd demo:
```bash
# cat > /etc/sysconfig/node_exporter << EOF
OPTIONS="--collector.textfile.directory /var/lib/node_exporter/textfile_collector"
EOF
# cat > /lib/systemd/system/node_exporter.socket << EOF
[Unit]
Description=Node Exporter

[Socket]
ListenStream=9100

[Install]
WantedBy=sockets.target
EOF
# cat > /lib/systemd/system/node_exporter.service << EOF
[Unit]
Description=Node Exporter
Requires=node_exporter.socket

[Service]
User=root
EnvironmentFile=/etc/sysconfig/node_exporter
ExecStart=/usr/local/bin/node_exporter --web.systemd-socket $OPTIONS

[Install]
WantedBy=multi-user.target
EOF
```

# grafana
参考:
- [玩转 Grafana 可视化系统](https://www.infvie.com/ops-notes/play-grafana.html)

grafana 是一款采用 go 语言编写的开源应用，主要用于大规模指标数据的可视化展现，是网络架构和应用分析中最流行的时序数据展示工具，目前已经支持绝大部分常用的时序数据库.

Grafana支持许多不同的数据源. 每个数据源都有一个特定的查询编辑器,该编辑器定制的特性和功能是公开的特定数据来源. 官方支持以下数据源:Graphite，Elasticsearch，InfluxDB，Prometheus，Cloudwatch，MySQL和OpenTSDB等.

每个数据源的查询语言和能力都是不同的, 可以把来自多个数据源的数据组合到一个dashborad，但每一个panel被绑定到一个特定的数据源,它就属于一个特定的组织.

## 基本概念
DashBoard：仪表盘，就像汽车仪表盘一样可以展示很多信息，包括车速，水箱温度等. Grafana的DashBoard就是以各种图形的方式来展示从Datasource拿到的数据.

Row：行，DashBoard的基本组成单元，一个DashBoard可以包含很多个row. 一个row可以展示一种信息或者多种信息的组合(即由Panels组成), 比如系统内存使用率，CPU五分钟及十分钟平均负载等, 所以在一个DashBoard上可以集中展示很多内容.

Panel：面板，实际上就是row展示信息的方式，支持表格（table），列表（alert list），热图（Heatmap）等多种方式.

Query Editor：查询编辑器，用来指定获取哪一部分数据. 类似于sql查询语句，比如你要在某个row里面展示test这张表的数据，那么Query Editor里面就可以写成`select * from test`. 这只是一种比方，[实际上每个DataSource获取数据的方式都不一样，所以写法也不一样](http://docs.grafana.org/features/datasources/).

Organization：组织，org是一个很大的概念，每个用户可以拥有多个org，grafana有一个默认的main org. 用户登录后可以在不同的org之间切换，前提是该用户拥有多个org. 不同的org之间完全不一样，包括datasource，dashboard等都不一样. 创建一个org就相当于开了一个全新的视图，所有的datasource，dashboard等都要再重新开始创建.

User：用户，这个概念应该很简单. Grafana里面用户有三种角色admin,editor,viewer:
- admin权限最高，可以执行任何操作，包括创建用户，新增Datasource，创建DashBoard
- editor角色不可以创建用户，不可以新增Datasource，可以创建DashBoard

  在2.1版本及之后新增了一种角色read only editor（只读编辑模式），这种模式允许用户修改DashBoard，但是不允许保存.
- viewer角色仅可以查看DashBoard.

每个user可以拥有多个organization.

## 部署grafana
```
# docker volume create data-grafana # 使用docker volume: docker run ... -v grafana-storage:/var/lib/grafana ...
# docker run -d --restart=unless-stopped --net=host -p 3000:3000 -v data-grafana:/var/lib/grafana --name grafana  grafana/grafana
```

访问http://localhost:3000, 初始密码: admin/admin, 可通过更新Grafana配置文
件的`[security]`部分来控制

> grafana配置位置: /etc/grafana/grafana.ini

非docker版安装参考[Download Grafana](https://grafana.com/grafana/download).

### Grafana 配置数据源
点击左侧菜单栏-设置-data sources-"Add data source"-选择"Prometheus", 在"Settings" tag页输入Prometheus配置信息, 再选中"Dashboards" tag页Import "Prometheus 2.0 Stats", 再保存即可.

### Grafana dashboards添加
点击左侧菜单栏中的"+"图标下的import，在页面中输入grafana labs的[dashboards](https://grafana.com/grafana/dashboards)下的插件id（node_exporter：1860， mysqld_exporter: 6239， jvm: 4701/9568），会自动跳转至配置页，选择数据源为prometheus，然后点击import即可.

> 不联网环境使用`https://grafana.com/grafana/dashboards/1860`右侧页面`Download JSON`链接下载的json配置即可.

> Grafana首页的Dashboards tag页仅显示已使用过的dashboard, 初次可用左侧菜单栏的搜索按钮进行查找.

# alertmanager
```bash
# ./alertmanager --version
# ./alertmanager --config.file="alertmanager.yml"
```

## 警报触发
警报可能有以下三种状态：
- Inactive： 警报未激活
- Pending： 警报已满足测试表达式条件， 但仍在等待for子句中指定的持续时间
- Firing： 警报已满足测试表达式条件， 并且Pending的时间已超过for子句的持续时间

  警报现在处于Firing状态, 会将其推送到Alertmanager

Pending到Firing的转换可以确保警报更有效, 且不会来回浮动. 没有for子句的警报会自动从Inactive转换为Firing, 只需要一个评估周期即可触发. 带有for子句的警报将首先转换为Pending， 然后转换为Firing, 因此至少需要两个评估周期才能触发.

silence: 特定时期内忽略警报

设置silence：
- 通过Alertmanager Web控制台
- 通过amtool命令行工具

## alertmanager.yml
group_by控制的是Alertmanager分组警报的方式

# FAQ
### 让Prometheus reload配置激活新的job的方法
1. `send SIGHUP signal`

    `kill -HUP <pid>`
1. `send a HTTP POST to the Prometheus web server`

    需要追加`--web.enable-lifecycle`选项, 即`/prometheus --config.file=prometheus.yml --web.enable-lifecycle` + `curl -X POST http://localhost:9090/-/reload`
### [exporters](https://prometheus.io/docs/instrumenting/exporters/)
### [prometheus 监控服务硬盘使用量过大，如何处理](https://asktug.com/t/topic/2588)
参考:
- [prometheus storage](https://prometheus.io/docs/prometheus/latest/storage/)

prometheus(`/home/tidb/tidb-deploy/prometheus-9090/scripts/run_prometheus.sh`)使用`--storage.tsdb.retention="15d"和--storage.tsdb.retention.size="2GB"`参数

> `--storage.tsdb.retention`默认是`15d`

> [prometheus不支持将storage.tsdb.retention加入prometheus.yml的原因: storage属于不能动态刷新的配置](https://github.com/prometheus/prometheus/issues/6188).

> curl -s http://localhost:9090/api/v1/status/runtimeinfo | jq '.data.storageRetention': 查看storage.tsdb.retention 

### grafana添加"Data Sources / Prometheus"报`HTTP Error Bad Gateway`
尝试使用`curl http://<prometheus_sever>/metrics`测试, 通常是当前浏览器无法访问到`http://<prometheus_sever>/metrics`导致的, 比如grafana, prometheus部署在aliyun, 此时用`localhost:9090`作为prometheus url就会报该错.

### grafana添加prometheus源时报`Error reading Prometheus: Post "/api/v1/query":  unsupported protocol scheme`
配置中的HTTP节的URL必须填写内容.

### grafana share, 比如面板内嵌(embe)
ref:
- [Share a panel](https://grafana.com/docs/grafana/latest/sharing/share-panel/)

方法需两步:
1. 在grafana.ini的`[security]`中启用`allow_embedding=true`, 否则会报`Refused to display '<grafana url>' in a frame because it set 'X-Frame-Options' to 'deny'.`
2. 在grafana.ini的`[auth.anonymous]`中启用`enabled=true`, 否则嵌入的iframe就需要登入.