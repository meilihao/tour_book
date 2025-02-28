# netdata
ref:
- [Netdata 系统监控工具](https://xiashuo.xyz/posts/devops/netdata/)

netdata 是一个用于系统和应用的分布式实时性能和健康监控工具
ps: [linux离线安装NETDATA](https://blog.csdn.net/E08640104/article/details/138911328)和[最新的资源包](https://api.github.com/repos/netdata/netdata/releases/latest)

## 安装
```bash
# --- repo 支持
# yum install netdata
# --- repo 不支持, 比如kylinos
# ./netdata-latest.gz.run -h
# ./netdata-latest.gz.run --accept
```

访问: http://<ip>:19999/v1

> `http://<ip>:19999/v2`/`http://<ip>:19999/v3`界面类似, 但难用, 比如查看`/`剩余空间时

## FAQ
### 数据保存位置
ref:
- [更改Netdata存储指标的时间](https://github.com/netdata/localization/blob/master/zh/docs/tutorials/longer-metrics-storage.md)

/var/cache/netdata

### 允许外部ip请求
``` bash
# vim /etc/netdata/netdata.conf
...
[web]
    bind to = * # 127.0.0.1 ::1
..
# systemctl restart netdata
```

### 获取当前配置
`curl http://localhost:19999/netdata.conf`

### 修改保存时长
``` bash
# vim /etc/netdata/netdata.conf
...
[global]
    dbengine multihost disk space = 256 # 256MB
..
# systemctl restart netdata
```