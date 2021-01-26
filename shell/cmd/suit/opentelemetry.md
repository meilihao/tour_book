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
```bash
# docker pull docker.elastic.co/elasticsearch/elasticsearch:7.10.1
# docker volume create data-es
# docker run -d --restart=unless-stopped --net host -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" -v data-es:/usr/share/elasticsearch/data docker.elastic.co/elasticsearch/elasticsearch:7.10.1
```

### 1. jaeger

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

**当前(2020/11/07), opentelemetry-all-in-one(jaeger-collector)不支持opentelemetry go sdk发送的数据, 报["unknown service opentelemetry.proto.collector.metrics.v1.MetricsService"](https://github.com/jaegertracing/jaeger/issues/2620)**

> jaeger-collector 等价于 OpenTelemetry collector, 估计是jaeger定制的OpenTelemetry collector.

### 2. prometheus+grafana
参考[prometheus_grafana.md](/shell/cmd/suit/prometheus_grafana.md)

### 3. OpenTelemetry Collector
1. 下载[binary](https://github.com/open-telemetry/opentelemetry-collector/releases).
1. 运行otelcol_linux_amd64

    ```bash
    # cat << EOF > otel-config.yaml
    # [OpenTelemetry Collector Architecture](https://github.com/open-telemetry/opentelemetry-collector/blob/master/docs/design.md)
    # [Configuring the OpenTelemetry Collector:](https://www.sumologic.com/blog/configure-opentelemetry-collector/)
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
      # logging:
      #   logLevel: debug
      jaeger:
        endpoint: localhost:14250
        insecure: true # 否则需要设置tls
      prometheus: # 作为prometheus的client(opentelemetry-collector收集的其clients的metrics), 被prometheus server采集metrics
        endpoint: "0.0.0.0:8889"
        namespace: promexample
        const_labels:
          label1: value1

    service:

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

      extensions: [health_check, pprof, zpages]
    EOF
    # ./otelcol_linux_amd64 --config otel-config.yaml
    ```

otelcol端口:
- 1777 : pprof extension
- 8888 : prometheus server采集OpenTelemetry Collector's metrics
- 8889 : 作为prometheus的client(opentelemetry-collector收集的其clients的metrics), 被prometheus server采集metrics
- 13133 : health_check extension
- 55679 : zpages extension
- 55680 : grpc, 接收OpenTelemetry client的上传
- 55681 : http, 接收OpenTelemetry client的上传

> opentelemetry collector 支持级联by [otlpexporter](https://github.com/open-telemetry/opentelemetry-collector/tree/master/exporter/otlpexporter)/[otlphttpexporter](https://github.com/open-telemetry/opentelemetry-collector/tree/master/exporter/otlphttpexporter), 此时前一级的opentelemetry-collector也被成为opentelemetry agent.
