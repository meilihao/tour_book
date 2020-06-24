# tidb

## 部署 
### [tiup(**推荐**)](https://github.com/pingcap-incubator/tiup-cluster)
```
# tiup cluster list # 查看已部署的集群
# generate rsa key for root@xxx
# tiup cluster deploy tidb-test v4.0.0-rc ./topology.yaml --user root -i /home/chen/.ssh/tidb_rsa
# tiup cluster destroy tidb-test # 它会使用deploy创建的ssh key(可用tiup cluster list获取该key)去destroy
# tiup cluster display tidb-test # 显示cluster的状态
# mysql -u root -h 47.100.15.8 -P 4000 # 默认没密码
> SET PASSWORD FOR 'root'@'%' = 'xxx';
# tiup cluster upgrade tidb-test v4.0.1 [--force] # cluster升级
```

> **部署完记得修改grafana, tidb的密码**

> tidb dashboard可用`tiup cluster display tidb-test --dashboard`查看, from v4.0.1 pd-server的2379已支持绑定到`*`.

```yaml
# Global variables are applied to all deployments and as the default value of
# them if the specific deployment value missing.
global:
  user: "tidb"
  ssh_port: 22
  deploy_dir: "/home/tidb/tidb-deploy"
  data_dir: "/home/tidb/tidb-data"

monitored:
  deploy_dir: "/home/tidb/tidb-deploy/monitored-9100"
  data_dir: "/home/tidb/tidb-data/monitored-9100"
  log_dir: "/home/tidb/tidb-deploy/monitored-9100/log"

server_configs:
  tidb:
    log.slow-threshold: 300
    log.level: warn
    binlog.enable: false
    binlog.ignore-error: false
  tikv:
    readpool.storage.use-unified-pool: true
    readpool.coprocessor.use-unified-pool: true
  pd:
    schedule.leader-schedule-limit: 4
    schedule.region-schedule-limit: 2048
    schedule.replica-schedule-limit: 64
    replication.enable-placement-rules: true


pd_servers:
  - host: 172.19.136.22
    # ssh_port: 22
    # name: "pd-1"
    client_port: 12379
    peer_port: 12380
    # deploy_dir: "deploy/pd-2379"
    # data_dir: "data/pd-2379"
    # log_dir: "deploy/pd-2379/log"
    # numa_node: "0,1"
    # # Config is used to overwrite the `server_configs.pd` values
    # config:
    #   schedule.max-merge-region-size: 20
    #   schedule.max-merge-region-keys: 200000
tidb_servers:
  - host: 172.19.136.22
    # ssh_port: 22
    # port: 4000
    # status_port: 10080
    # deploy_dir: "deploy/tidb-4000"
    # log_dir: "deploy/tidb-4000/log"
    # numa_node: "0,1"
    # # Config is used to overwrite the `server_configs.tidb` values
    # config:
    #   log.level: warn
    #   log.slow-query-file: tidb-slow-overwritten.log
tikv_servers:
  - host: 172.19.136.22
    # ssh_port: 22
    # port: 20160
    # status_port: 20180
    # deploy_dir: "deploy/tikv-20160"
    # data_dir: "data/tikv-20160"
    # log_dir: "deploy/tikv-20160/log"
    # numa_node: "0,1"
    # # Config is used to overwrite the `server_configs.tikv` values
    #  config:
    #    server.labels:
    #      zone: sh
    #      dc: sha
    #      rack: rack1
    #      host: host1
tiflash_servers: # 需要很大内存, 当前见到是1.2g
  - host: 172.19.136.22
    # ssh_port: 22
    # tcp_port: 9000
    # http_port: 8123
    # flash_service_port: 3930
    # flash_proxy_port: 20170
    # flash_proxy_status_port: 20292
    # metrics_port: 8234
    # deploy_dir: deploy/tiflash-9000
    # data_dir: deploy/tiflash-9000/data
    # log_dir: deploy/tiflash-9000/log
    # numa_node: "0,1"
    # # Config is used to overwrite the `server_configs.tiflash` values
    #  config:
    #    logger:
    #      level: "info"
    #  learner_config:
    #    log-level: "info"
monitoring_servers:
  - host: 172.19.136.22
grafana_servers:
  - host: 172.19.136.22
alertmanager_servers:
  - host: 172.19.136.22
```

### [~~tidb-ansible~~, 作废, 已切换到tiup](https://github.com/pingcap/tidb-ansible)
使用**用户名tidb**进行部署.

#### 滚动升级
```sh
ansible-playbook local_prepare.yml -c local
ansible-playbook rolling_update.yml -c local
```

#### Error
1. `ansible-playbook bootstrap.yml`报错
fatal: [127.0.0.1]: UNREACHABLE! => {"changed": false, "msg": "Failed to connect to the host via ssh: Permission denied (publickey).\r\n", "unreachable": true}

```
$ sudo ansible-playbook bootstrap.yml -c local
```

1. fail when NTP service is not running or ntpstat is not synchronised to NTP server

```
$ sudo apt-get install ntpstat
```

1. The default max number of file descriptors is too low 65535 should be 100000

```
$ sudo vim /etc/security/limits.conf # 账号需重新登录
```

1. This machine does not have sufficient CPU to run TiDB, at least 8 cores.
解决方法:
    - 将`group_vars/all.yml`里的`dev_mode`设为`True`, 这样bootstrap.yml 阶段会跳过磁盘检测、CPU、内存容量检测, **开发阶段推荐此配置**
    - 修改`roles/check_system_optional/defaults/main.yml`里的`tidb_min_cpu` (**在v4.0时不用修改, dev_mode=True后会跳过**)

参考:
- [Deploy 2.0GA failed #6423](https://github.com/pingcap/tidb/issues/6423)
- [新增tidb节点报错 #491](https://github.com/pingcap/tidb-ansible/issues/491)

1. 卡在`wait xxx up`
通常是权限问题, 即`/home/tidb`下的文件不是tidb所有引起, 通过`sudo chown -R tidb:tidb *`进行修正即可.

1. TiDB-Ansible部署的pd_server 端口和 rancher 部署etcd的端口存在冲突, 比如2380
解决方法:
    - 修改`group_vars/pd_servers.yml`里的pd_client_port和pd_peer_port
    - [如何自定义端口](https://pingcap.com/docs-cn/dev/how-to/deploy/orchestrated/ansible/)

或直接修改inventory.ini:
```yaml
[pd_servers]
172.19.136.22 pd_client_port=12379 pd_peer_port=12380
```

1. Check if the operating system supports EPOLLEXCLUSIVE : epollexclusive is not available
EPOLLEXCLUSIVE是4.5+内核新添加的一个 epoll 的标识, 需内核支持

> [惊群效应](https://mcgrady-forever.github.io/2018/03/19/network-thundering-herd/)
> [Ngnix 是如何解决 epoll 惊群的](https://simpleyyt.com/2017/06/25/how-ngnix-solve-thundering-herd/)

1. Make sure NTP service is running and ntpstat is synchronised to NTP server
```sh
$ sudo systemctl start ntpd
```

1. [wait for region replication complete](https://github.com/pingcap/tidb-ansible/issues/846)