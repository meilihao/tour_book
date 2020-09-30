# prometheus & grafana
![架构](https://prometheus.io/assets/architecture.png)

## 部署grafana
```
# mkdir -p /var/lib/grafana
# chmod 777 /var/lib/grafana
# docker run -d  --net=host -p 3000:3000 -v /var/lib/grafana:/var/lib/grafana --name grafana  grafana/grafana
```

访问http://localhost:3000, 初始密码: admin/admin

> grafana配置位置: /etc/grafana/grafana.ini

## 部署prometheus
```
# mkdir -p /var/lib/prometheus
# chmod 777 /var/lib/prometheus
# docker run -d --net=host -p 9090:9090 -v /var/lib/prometheus:/prometheus --name prometheus prom/prometheus
```

访问http://localhost:9090

其他配置:
- --config.file=prometheus.yml : 指定配置文件
- --web.listen-address=:9090 : 指定web访问端口, 此时必须指定`--config.file=prometheus.yml`

## 部署node_exporter
```
# ./node_exporter
```

访问: http://${ip}:9100/metrics即可.