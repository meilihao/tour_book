# prometheus & grafana
![架构](https://prometheus.io/assets/architecture.png)

## 部署prometheus
```
# mkdir -p /var/lib/prometheus
# chmod 777 /var/lib/prometheus
# docker run -d --net=host -p 9090:9090 -v /etc/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml -v /var/lib/prometheus:/prometheus --name prometheus prom/prometheus
```

> prometheus.yml可通过不指定`-v /etc/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml`时进入容器来获取 

访问http://localhost:9090

其他配置:
- --config.file=prometheus.yml : 指定配置文件
- --web.listen-address=:9090 : 指定web访问端口, 此时必须指定`--config.file=prometheus.yml`

## 部署node_exporter
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

## 部署grafana
```
# mkdir -p /var/lib/grafana
# chmod 777 /var/lib/grafana
# docker run -d  --net=host -p 3000:3000 -v /var/lib/grafana:/var/lib/grafana --name grafana  grafana/grafana
```

访问http://localhost:3000, 初始密码: admin/admin

> grafana配置位置: /etc/grafana/grafana.ini

### Grafana 配置数据源
点击左侧菜单栏-设置-data sources-"Add data source"-选择"Prometheus", 在"Settings" tag页输入Prometheus配置信息, 再选中"Dashboards" tag页Import "Prometheus 2.0 Stats", 再保存即可.

### Grafana dashboards添加
点击左侧菜单栏中的"+"图标下的import，在页面中输入grafana labs的[dashboards](https://grafana.com/grafana/dashboards)下的插件id（node_exporter：1860， mysqld_exporter: 6239， jvm: 4701/9568），会自动跳转至配置页，选择数据源为prometheus，然后点击import即可.

> 不联网环境使用`https://grafana.com/grafana/dashboards/1860`右侧页面`Download JSON`链接下载的json配置即可.

> Grafana首页的Dashboards tag页仅显示已使用过的dashboard, 初次可用左侧菜单栏的搜索按钮进行查找.

## FAQ
### 让Prometheus reload配置激活新的job的方法
1. `send SIGHUP signal`

    `kill -HUP <pid>`
1. `send a HTTP POST to the Prometheus web server`

    需要追加`--web.enable-lifecycle`选项, 即`/prometheus --config.file=prometheus.yml --web.enable-lifecycle` + `curl -X POST http://localhost:9090/-/reload`
### [exporters](https://prometheus.io/docs/instrumenting/exporters/)