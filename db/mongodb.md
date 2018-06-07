# mongodb
推荐gui client: studio3t.

## 操作
### 登录
1. 直接登录,类似mysql
```sh
$ mongo --host s-xxx-pub.mongodb.rds.aliyuncs.com:3717 -u root -p 123456 --authenticationDatabase mydb
```

2. 先连接后验证
```sh
$ mongo --host s-xxx-pub.mongodb.rds.aliyuncs.com:3717                        16:50:51
...
mongos> use mydb
switched to db mydb
mongos> db.auth("root","123456")
1 // 输出 1 表示验证成功
```

## Error