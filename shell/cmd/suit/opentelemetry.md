# opentelemetry

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

## demo
1. [go demo](https://github.com/open-telemetry/opentelemetry-go/tree/master/example/otel-collector)

	演示了如何从OpenTelemetry-Go SDK 中导出trace和metric数据，并将其导入opentelemetry-collector，最后通过collector将trace数据传递给Jaeger，将metric数据传递给Prometheus.
1. [OpenTelemetry Collector Demo](https://github.com/open-telemetry/opentelemetry-collector/tree/master/examples)和[如何在客户端和服务器中正确使用OpenTelemetry导出器和OpenTelemetry收集器？](https://mlog.club/article/5893955)

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
