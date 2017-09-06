# gitlab

## 安装
[Docker 安装](https://docs.gitlab.com.cn/omnibus/docker/)

## 环境
- 账号
Username: root
Password: 5iveL!fe

- postgres
```sh
$ su - gitlab-psql
$ psql -h /var/opt/gitlab/postgresql -d gitlabhq_production
// 即
$ sudo -u gitlab-psql psql -h /var/opt/gitlab/postgresql -d gitlabhq_production
```

- redis
```sh
$ redis-cli -s /var/opt/gitlab/redis/redis.socket
```