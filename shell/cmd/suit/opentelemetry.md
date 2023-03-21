# opentelemetry
参考:
- [OpenTelemetry – 云原生下可观测性的新标准](https://www.daqianduan.com/17340.html)

> elemetry : [təˈlemətri], 遥测

OpenTelemetry合并了OpenTracing和OpenCensus项目，提供了一组API和库来标准化遥测数据的采集和传输. OpenTelemetry提供了一个安全，厂商中立的工具，这样就可以按照需要将数据发往不同的后端.

OpenTelemetry的核心工作目前主要集中在3个部分：
- 规范的制定和协议的统一，规范包含数据传输、API的规范，协议的统一包含：HTTP W3C的标准支持及GRPC等框架的协议标准
- 多语言SDK的实现和集成，用户可以使用SDK进行代码自动注入和手动埋点，同时对其他三方库（Log4j、LogBack等）进行集成支持
- 数据收集系统的实现，当前是基于OpenCensus Service的收集系统，包括Agent和Collector

  - Exporters：可以将数据发往一个可选择的后端

  [OpenTelemetry exporters list](https://opentelemetry.io/registry/)
  - Collectors：厂商中立的实现，用于处理和导出遥测数据

OpenTelemetry的自身定位很明确：数据采集和标准规范的统一，对于数据如何去使用、存储、展示、告警，官方是不涉及的. 目前社区推荐使用Prometheus + Grafana做Metrics存储、展示，使用Jaeger做分布式跟踪的存储和展示.

OpenTelemetry的终极目标是：实现Metrics、Tracing、Logging的融合及大一统，作为APM的数据采集终极解决方案:
- Tracing：提供了一个请求从接收到处理完成整个生命周期的跟踪路径，一次请求通常过经过N个系统，因此也被称为分布式链路追踪
- Metrics：例如cpu、请求延迟、用户访问数等Counter、Gauge、Histogram指标
- Logging：传统的日志，提供精确的系统记录

这三者的组合可以形成大一统的APM解决方案：
- 基于Metrics告警发现异常
- 通过Tracing定位到具体的系统和方法
- 根据模块的日志最终定位到错误详情和根源
- 调整Metrics等设置，更精确的告警/发现问题

在OpenTelemetry中使用Context为Metrics、Logging、Tracing提供统一的上下文，三者均可以访问到这些信息，同时Context可以随着请求链路的深入，不断往下传播:
- Context数据在Task/Request的执行周期中都可以被访问到
- 提供统一的存储层，用于保存Context信息，并保证在各种语言和处理模型下都可以工作（例如单线程模型、线程池模型、CallBack模型、Go Routine模型等）
- 多种维度的关联基于元信息(标签)实现，元信息由业务确定，例如：通过Env来区别是测试还是生产环境等
- 提供分布式的Context传播方式，例如通过W3C的traceparent/tracestate头、GRPC协议等

## opentelemetry-collector
参考:
- [Opentelemetry Collector的配置和使用](https://cloud.tencent.com/developer/article/1735992)

collector通过pipeline处理service中启用的数据. 

### pipeline
pipeline由接收遥测数据的组件构成，包括：
- [receivers](https://github.com/open-telemetry/opentelemetry-collector/blob/master/receiver/README.md)

	receivers定义了数据如何进入OpenTelemetry Collector, 必须配置一个或多个receiver.
- [processors](https://github.com/open-telemetry/opentelemetry-collector/blob/master/processor/README.md)

	processors运行在数据的接收和导出之间. 虽然processors是可选的，但有时候会建议使用processors.
- [exporters](https://github.com/open-telemetry/opentelemetry-collector/blob/master/exporter/README.md)

	exporters指定了如何将数据发往一个或多个后端/目标, 必须配置一个或多个exporter.

### service
service部分用于配置opentelemetry-collector根据receivers, processors, exporters和extensions的配置会启用哪些特性.

service分为两部分：
- [extensions](https://github.com/open-telemetry/opentelemetry-collector/blob/master/extension/README.md)

	用于监控opentelemetry-collector的健康状态. extensions是可选的
- [pipelines](https://github.com/open-telemetry/opentelemetry-collector/blob/master/docs/pipelines.md)

	一个pipeline是一组 receivers, processors, 和exporters的集合. 必须在service之外定义每个receiver/processor/exporter的配置，然后将其包含到pipeline中.

注：每个receiver/processor/exporter都可以用到多个pipeline中。当多个pipeline引用processor(s)时，每个pipeline都会获得该processor(s)的一个实例，这与多个pipeline中引用receiver(s)/exporter(s)的情况不同(所有pipelines仅能获得receiver/exporter的一个实例)。

pipelines有两类：
- metrics: 采集和处理metrics数据
- traces: 采集和处理trace数据

## OpenTelemetry API
参考:
- [^From 0 to Insight with OpenTelemetry in Go](https://www.honeycomb.io/blog/from-0-to-insight-with-opentelemetry-in-go/)
- [OpenTelemetry: 经得起考验的工具](https://www.cnblogs.com/charlieroro/p/13862471.html)
- [opentelemetry-java/QUICKSTART.md](https://github.com/open-telemetry/opentelemetry-java/blob/master/QUICKSTART.md)
- [Setup and Install OpenTelemetry Go](https://opentelemetry.lightstep.com/go/setup-and-installation/)
- [通过Open Telemetry上报Java应用数据](https://www.alibabacloud.com/help/zh/doc-detail/196899.htm)
- [opentracing-io](https://wu-sheng.gitbooks.io/opentracing-io/content/pages/spec.html)
- [^opentelemetry-go api status](https://opentelemetry.io/docs/go/) : go sdk未实现log api

应用开发者会使用 Open Telemetry API对其代码进行插桩，库作者会用它(在库中)直接编写桩功能. API不处理操作问题，也不关心如何将数据发送到厂商后端:

API分为四个部分：
- A Tracer API

  Tracer API 支持生成spans，可以给span分配一个traceId，也可以选择性地加上时间戳. 一个Tracer会给spans打上名称和版本. 当查看数据时，名称和版本会与一个Tracer关联，通过这种方式可以追踪生成的sapan
- A Metrics API

  [Metric API](https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/metrics/api.md)提供了多种类型的Metric instruments(桩功能)，如Counters 和Observers. Counters 允许对度量进行计算，Observers允许获取离散时间点上的测量值. 例如，可以使用Observers 观察不在Span上下文中出现的数值，如当前CPU负载或磁盘上空闲的字节数.
- A Context API

  [Context API](https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/context/context.md)会在使用相同"context"的spans和traces中添加上下文信息，如[W3C Trace Context](https://www.w3.org/TR/trace-context/), Zipkin B3首部, 或 [New Relic distributed tracing](https://docs.newrelic.com/docs/understand-dependencies/distributed-tracing/get-started/introduction-distributed-tracing) 首部. 此外该API允许跟踪spans是如何在一个系统中传递的. 当一个trace从一个处理传递到下一个处理时会更新上下文信息. Metric instruments可以访问当前上下文.
- 语义规范

  OpenTelemetry API包含一组[语义规范](https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/overview.md#semantic-conventions.md)，该规范包含了命名spans，属性以及与spans相关的错误. 通过将该规范编码到API接口规范中，OpenTelemetry 项目保证所有的instrumentation(不论任何语言)都包含相同的语义信息. 对于希望为所有用户提供一致的APM体验的厂商来说，该功能非常有价值.

> 跨span传递trace的原理: [传递traceparent](https://github.com/open-telemetry/opentelemetry-go/tree/master/propagation/trace_context.go)

## demo
1. [Instrumentation](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/master/instrumentation)
1. [lightstep / opentelemetry-examples go demo](https://github.com/lightstep/opentelemetry-examples/tree/main/go)
1. [go demo](https://github.com/open-telemetry/opentelemetry-go/tree/master/example/otel-collector)

	演示了如何从OpenTelemetry-Go SDK 中导出trace和metric数据，并将其导入opentelemetry-collector，最后通过collector将trace数据传递给Jaeger，将metric数据传递给Prometheus.
1. [OpenTelemetry Collector Demo](https://github.com/open-telemetry/opentelemetry-collector/tree/master/examples)和[如何在客户端和服务器中正确使用OpenTelemetry导出器和OpenTelemetry收集器？](https://mlog.club/article/5893955)
1. [OpenTelemetry 支持split driver即分开导出trace/metric](https://github.com/open-telemetry/opentelemetry-go/blob/master/exporters/otlp/example_test.go)

## 完整部署
测试client: [exporter_otelcol](https://github.com/meilihao/demo/tree/master/opentelemetry/exporter_otelcol)

### 0. es
ref:
- [es下载地址](https://www.elastic.co/cn/downloads/past-releases#elasticsearch)

```bash
# docker pull docker.elastic.co/elasticsearch/elasticsearch:7.10.1
# docker volume create data-es
# docker run -d --restart=unless-stopped --net host -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" -v data-es:/usr/share/elasticsearch/data docker.elastic.co/elasticsearch/elasticsearch:7.10.1
```

> 开启es xpack security时, [jaeger配置方法见这里](https://github.com/jaegertracing/jaeger/issues/3382).

### 1. jaeger
ref:
- [jaeger 存储要求](https://www.jaegertracing.io/docs/1.36/#technical-specs)

  es 7.x可以, es 8.3.1时jaeger-collector启动报错

[演示部署(数据在内存)](https://www.jaegertracing.io/docs/1.21/getting-started/).
```bash
# docker run -d --name jaeger \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 14250:14250 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.21
```

> [jaeger使用的port](https://www.jaegertracing.io/docs/1.29/getting-started/)

实际部署参考[start_jaeger.sh](/container/dockerfile/jaeger_elasticsearch)

访问http://localhost:16686即可看到jaeger ui.

或直接部署jaeger-opentelemetry-collector:
```bash
# mkdir config && cd config
# cat << EOF > config.yaml
exporters:
  jaeger_elasticsearch:
    es:
      server-urls: http://elasticsearch:9200
      num-replicas: 1
processors:
  attributes:
    actions:
      - key: user
        action: delete
service:
  pipelines:
    traces:
      processors: [attributes]
EOF
# docker run --rm -it --net host -v ${PWD}:/config \
    -e SPAN_STORAGE_TYPE=elasticsearch \
    jaegertracing/opentelemetry-all-in-one \
    --config-file=/config/config.yaml \
    --es.server-urls=http://localhost:9200 \
    --es.num-shards=1
# ss -anlpt |grep 5568 # 看到55680/55681
# ss -anlpt |grep 16686 # 16686 for jaeger ui
```

**当前(2020/11/07), opentelemetry-all-in-one(jaeger-collector)不支持opentelemetry go sdk发送的metrics数据, 报["unknown service opentelemetry.proto.collector.metrics.v1.MetricsService"](https://github.com/jaegertracing/jaeger/issues/2620)**

> jaeger-collector 等价于 OpenTelemetry collector, 估计是jaeger定制的OpenTelemetry collector.

### 2. prometheus+grafana
参考[prometheus_grafana.md](/shell/cmd/suit/prometheus_grafana.md)

### 3. OpenTelemetry Collector(v0.41.0)
1. 下载[otelcol-contrib binary](https://github.com/open-telemetry/opentelemetry-collector-releases/releases).
1. 运行otelcol_linux_amd64

    ref:
    - [opentelemetry-collector examples](https://github.com/open-telemetry/opentelemetry-collector/blob/main/examples)
    - [OpenTelemetry Collector Architecture](https://github.com/open-telemetry/opentelemetry-collector/blob/main/docs/design.md)
    - [Configuring the OpenTelemetry Collector:](https://www.sumologic.com/blog/configure-opentelemetry-collector/)

    **[opentelemetry-collector git repo默认不再包含`exporter/jaegerexporter`(已被移到[`opentelemetry-collector-contrib `](https://github.com/open-telemetry/opentelemetry-collector-contrib)), PR是[这个](](https://github.com/open-telemetry/opentelemetry-collector/issues/3474)), 因此可用ocb工具精简binary. 同时根据[otel发行版Distributions](https://opentelemetry.io/docs/collector/distributions/), otel官方提供的otelcol-contrib已包含全部opentelemetry-collector-contrib组件**.

    > otel精简构建可参考[OpenTelemetry Collector builder](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder)+[distributions/otelcol-contrib/manifest.yaml](https://github.com/open-telemetry/opentelemetry-collector-releases/blob/main/distributions/otelcol-contrib/manifest.yaml)

    ```bash
    # --- 下载相应版本的ocb
    # cat << EOF > otelcol-builder.yaml
    exporters:
      - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/exporter/jaegerexporter v0.41.0
      - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter v0.41.0
    # ocb --config otelcol-builder.yaml # 参考 opentelemetry-collector-releases的distributions/otelcol-contrib/manifest.yaml
    # mv /tmp/otelcol-distribution1631641147/otelcol-custom /usr/local/bin/otelcol
    # cat << EOF > otel-config.yaml
    extensions:
      health_check:
      pprof:
        endpoint: 0.0.0.0:1777
      zpages:
        endpoint: 0.0.0.0:55679

    receivers:
      otlp:
        protocols:
          grpc:
          # http:

      # # Collect own metrics
      # prometheus:
      #   config:
      #     scrape_configs:
      #       - job_name: 'otel-collector'
      #         scrape_interval: 10s
      #         static_configs:
      #           - targets: ['0.0.0.0:8888']

    processors:
      batch:

    # logging is print in stdout
    exporters:
      logging:
        logLevel: debug
      jaeger:
        endpoint: localhost:14250
        tls:
          insecure: true # 否则需要设置tls
      prometheus: # 作为prometheus的client(opentelemetry-collector收集的其clients的metrics), 被prometheus server采集metrics
        endpoint: "0.0.0.0:8889"
        namespace: promexample
        const_labels:
          label1: value1

    service:
      extensions: [health_check, pprof, zpages]
      pipelines:

        # traces:
        #   receivers: [otlp]
        #   processors: [batch]
        #   exporters: [logging, jaeger]

        # metrics:
        #   receivers: [otlp]
        #   processors: [batch]
        #   exporters: [logging, prometheus]

        traces:
          receivers: [otlp]
          processors: [batch]
          exporters: [jaeger]

        metrics:
          receivers: [otlp]
          processors: [batch]
          exporters: [prometheus]
        logs:
          receivers: [otlp]
          processors: [batch]
          exporters: [logging]
    EOF
    # otelcol --config otel-config.yaml
    ```

    > 通过[otel receiver/otlpreceiver](https://github.com/open-telemetry/opentelemetry-collector/blob/main/receiver/otlpreceiver)和[otel exporter/otlpexporter](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/otlpexporter)推测exporters的otlp是用于级联opentelemetry-collector.

otelcol端口:
- 4317 : grpc, 接收OpenTelemetry client的上传(以前port是55680)
- 4318 : http, 接收OpenTelemetry client的上传(以前port是55681)
- 1777 : pprof extension
- 8888 : OpenTelemetry Collector's exposed metrics
- 8889 : 作为prometheus的client(opentelemetry-collector收集的其clients的metrics), 被prometheus server采集metrics
- 13133 : health_check extension
- 55679 : zpages extension
- 14250 : 安装otelcol-contrib 0.54.0后, 起来后会占用jaeger的14250端口, 使用上述otel-config.yaml重启后没再使用14250

> opentelemetry collector 支持级联by [otlpexporter](https://github.com/open-telemetry/opentelemetry-collector/tree/master/exporter/otlpexporter)/[otlphttpexporter](https://github.com/open-telemetry/opentelemetry-collector/tree/master/exporter/otlphttpexporter), 此时前一级的opentelemetry-collector也被成为opentelemetry agent.

## FAQ
### jaeger-ui `http://xxx:16686/search`刷新报`all shards failed [type=search_phase_execution_exception]`
具体报错接口是`curl http://openhello.net:16686/api/services`, 新旧版jaeger数据导致, 删除旧数据重新刷新即可.

删除elasticsearch所有数据:
```
curl -X DELETE 'http://localhost:9200/_all'
```

### 清空elasticsearch数据后, opentelemetry发送的数据未进入es
步骤:
```bash
# -- jaeger-collector和jaeger-query运行中
systemctl stop elasticsearch
rm -rf /var/lib/elasticsearch
systemctl start elasticsearch
# -- 发送opentelemetry数据给otel-collector
# -- 访问jaeger ui: http://xxx:16686/search, 没有数据
```

解决方法: 重启jaeger-collector和jaeger-query后解决.

### ocb构建报`IncludeCore is deprecated. Starting from v0.41.0, you need to include all components explicitly`
不是错误, 只是WARN并打印了stack信息.

### jaeger-collector写es报`[TOO_MANY_REQUESTS/12/disk usage exceeded flood-stage watermark, index has read-only-allow-delete block]`
是因为一次请求中批量插入的数据条数巨多，以及短时间内的请求次数巨多引起ES节点服务器存储超过限制，ES主动给索引上锁.

释放空间后执行`curl -XPUT -H "Content-Type: application/json" http://localhost:9200/_all/_settings -d '{"index.blocks.read_only_allow_delete": null}'`来解除索引只读状态

> es存储阈值是95%, 当前是96%.

### 异步trace
ref:
- [分布式链路追踪教程(一)---Opentracing 基本概念](https://www.lixueduan.com/posts/tracing/01-opentracing/)
- [Opentracing概念和术语](https://wu-sheng.gitbooks.io/opentracing-io/content/pages/spec.html)
- [Replace span relationship with a potentially remote parent context](https://github.com/open-telemetry/opentelemetry-go/pull/451/files)
- [Go HTTP Server 基于 OpenTelemetry 使用 Jaeger - 代码实操](https://xie.infoq.cn/article/1ee204583dae4fa8c14a9946f)
- [OpenTelemetry 规范阅读](https://jckling.github.io/2021/04/02/Jaeger/OpenTelemetry%20%E8%A7%84%E8%8C%83%E9%98%85%E8%AF%BB/)
- [Sync and Async children (FOLLOWS_FROM)](https://github.com/open-telemetry/opentelemetry-specification/issues/65)
- [Links between spans](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/overview.md#links-between-spans)
- [How to trace two asynchronous go routines with open telemetry](https://stackoverflow.com/questions/70438619/how-to-trace-two-asynchronous-go-routines-with-open-telemetry)

Opentracing 定义了两种引用关系:ChildOf和FollowFrom:
- ChildOf: 父span的执行依赖子span的执行结果时, 此时子span对父span的引用关系是ChildOf. 比如对于一次RPC调用, 服务端的span（子span）与客户端调用的span（父span）是ChildOf关系.
- FollowFrom：父span的执不依赖子span执行结果时, 此时子span对父span的引用关系是FollowFrom. FollowFrom常用于异步调用的表示, 例如消息队列中consumerspan与producerspan之间的关系.

opentelemetry使用Link描述该引用关系.

opentelemetry-go使用ContextWithRemoteSpanContext实现该效果.