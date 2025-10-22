# README
status: doing

## prepare
```bash
# --- 编译: README.md + .github/workflows/Dockerfile
# git clone --depth 1 -b v1.74.0 git@github.com:jaegertracing/jaeger-ui.git
# npm ci
# npm run build
```

没使用jaegertracing/jaeger-query, 因为只需要jaeger-ui的静态资源

jaeger-query架构:
```
┌─────────────┐   查询请求    ┌──────────────┐   读取数据    ┌─────────────┐
│   用户       │ ──────────→ │ jaeger-query  │ ──────────→ │  存储后端    │
│   (浏览器)   │              │              │              │ (Cassandra/ │
└─────────────┘   返回结果    └──────────────┘              │  ES/内存等)  │
                                   │                      └─────────────┘
                                   │ 提供UI
                             ┌─────────────┐
                             │   Web UI    │
                             │  静态资源   │
                             └─────────────┘
```

[VictoriaTraces/Querying/Visualization in Jaeger UI](https://docs.victoriametrics.com/victoriatraces/querying/jaeger-frontend/)

## docker-compose.yaml
结合Victoria{Traces, Metrics, Logs}/blob/master/deployment/docker/compose-vt-single.yml, 以Traces的compose-vt-single.yml为主

部署:
```bash
# docker network create jaeger
# cd VictoriaStack
# cp -r /data/tmpfs/jaeger-ui-1.74.0/packages/jaeger-ui/build nginx/html/jaeger-ui/
# docker compose -f docker-compose.yaml up -d
```

jaeger访问: http://127.0.0.1:8080
grafana: http://localhost:3000/login admin/admin

VictoriaLogs导航: http://localhost:9428
VictoriaLogs web ui: http://localhost:9428/select/vmui

VictoriaMetrics导航: http://localhost:8428
VictoriaMetrics web ui: http://localhost:8428/vmui

alertmanager web ui: http://localhost:9093
