## 部署
使用官方的[tidb-ansible](https://github.com/pingcap/tidb-ansible)

### 滚动升级
```sh
sudo ansible-playbook local_prepare.yml -c local
sudo ansible-playbook rolling_update.yml -c local
```

### Error
1. `ansible-playbook bootstrap.yml`报错
fatal: [127.0.0.1]: UNREACHABLE! => {"changed": false, "msg": "Failed to connect to the host via ssh: Permission denied (publickey).\r\n", "unreachable": true}

```
$ sudo ansible-playbook bootstrap.yml -c local 
```

2. fail when NTP service is not running or ntpstat is not synchronised to NTP server

```
$ sudo apt-get install ntpstat
```

3. The default max number of file descriptors is too low 65535 should be 100000

```
$ sudo vim /etc/security/limits.conf # 账号需重新登录
```