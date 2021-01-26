# kingbase v8(人大金仓)

> kingbase基于postgres

## 集群备份+启停
```
# /opt/Kingbase/ES/V8/Server/bin/kingbase_monitor.sh stop # 停止
# cp -r <数据目录>/data /opt/Kingbase # 备份数据
# cp -r /backup/data_backup /opt/Kingbase
# /opt/Kingbase/ES/V8/Server/bin/kingbase_monitor.sh start # 启动
```