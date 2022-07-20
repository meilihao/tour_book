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

## FAQ
### db init
ref:
- [GitLab在CockroachDB和YugabyteDB上的兼容性对比之系统初始化](https://www.tuicool.com/articles/Ibi6nuI)

`sudo -u git -H bundle exec rake gitlab:setup RAILS_ENV=production`

YugabyteDB支持PostgreSQL Extension，CockroachDB不支持, 且它们均实现了pg协议:
1. CockroachDB v22.1不支持Extension功能，导致GitLab无法初始化数据库
1. YugabyteDB v2.9不支持Gin Index（Generalized inverted indexes）, 但最新版本v2.13后支持, 可以正常访问GitLab页面以及注册用户