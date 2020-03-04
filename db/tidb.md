## 部署
使用官方的[tidb-ansible](https://github.com/pingcap/tidb-ansible), 使用**用户名tidb**进行部署.

### 滚动升级
```sh
ansible-playbook local_prepare.yml -c local
ansible-playbook rolling_update.yml -c local
```

### Error
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