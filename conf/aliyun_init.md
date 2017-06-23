# aliyun init

1. 登录
```
$ ssh -p 22 root@127.0.0.1
```

1. [禁用内核更新](https://help.aliyun.com/knowledge_detail/41199.html)
```
# /etc/yum.conf
exclude=kernel* centos-release*
```

再更新系统并重启.

1. 新建用户
```
// 新建用户
# useradd username
// 添加sudo权限
# vim /etc/sudoers
// username  ALL=(ALL)   ALL
// 修改密码
# passwd username
```

1. 配置ssh
[参考](security/2015_07_08_001.md),但推荐使用`ed25519`.