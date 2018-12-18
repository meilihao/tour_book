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
1. 将`group_vars/all.yml`里的`dev_mode`设为`True`, 这样bootstrap.yml 阶段会跳过磁盘检测、CPU、内存容量检测, **开发阶段推荐此配置**
1. 修改`roles/check_system_optional/defaults/main.yml`里的`tidb_min_cpu`

参考:
- [Deploy 2.0GA failed #6423](https://github.com/pingcap/tidb/issues/6423)
- [新增tidb节点报错 #491](https://github.com/pingcap/tidb-ansible/issues/491)

1. 卡在`wait xxx up`
通常是权限问题, 即`/home/tidb`下的文件不是tidb所有引起, 通过`sudo chown -R tidb:tidb *`进行修正即可.