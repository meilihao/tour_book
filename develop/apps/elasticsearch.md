# Elasticsearch

- [安装文档](https://www.elastic.co/guide/en/elasticsearch/reference/current/install-elasticsearch.html)
- [2万字详解，吃透 ES](https://www.tuicool.com/articles/NFjEfeZ)

> elasticsearch.service: ES_HOME = `/usr/share/elasticsearch`; ES_PATH_CONF = `/etc/elasticsearch`

> 数据目录: ES_PATH_CONF/elasticsearch.yaml中的path.data; log: elasticsearch.yaml中的path.logs

当前发现问题: es 8.3.1安装后使用默认生成的密码报授权失败且无法重置elastic用户的密码.

## 授权
ref:
- [Secure settings](https://www.elastic.co/guide/en/elasticsearch/reference/8.3/secure-settings.html)
- [Built-in users](https://www.elastic.co/guide/en/elasticsearch/reference/current/built-in-users.html)
- [密码使用策略](https://github.com/elastic/elasticsearch/pull/77036)

安装elasticsearch 8.3.1时, 提示了es默认启用了authentication.

在安装 Elasticsearch 时，如果 elastic 用户还没有密码，它将使用默认的 bootstrap. bootstrap 是一种临时密码，使你可以运行设置所有内置用户密码的工具.

默认情况下，bootstrap 密码来自于随机化的 keystore.seed 设置，该设置在安装过程中添加到了密钥库中。 你不需要知道或更改此 bootstrap 密码。 但是，如果您在密钥库中定义了 bootstrap.password 设置，则将使用该值. 为内置用户（尤其是 elastic 用户）设置密码后，bootstrap 密码将不再使用.

对秘钥库的任意修改都只有重启 Elasticsearch 才会生效.

es 官方有security开启的操作流程:
1. 检查使用license,因为不同的license能够使用的权限验证方式也是不一样的。
1. 检查集群，保证每一个node都进行了设置 xpack.security.enabled: true
1. 这个会开启security设置，也就开启了node之间使用ssl和基于用户的访问模式
1. 想要跑在专用的jvm上（这个一般用不到）
1. 开启node之间的transport层面的ssl/tls
1. 启动es
1. 为build-in user设置password
1. 选择一种realm的管理方式，basic只能使用 file,native两种方式，其他的都是收费的
1. 创建一些role和user来进行使用
1. 打开audit功能呢（可选的，实际上这个是收费的）

配置项说明:
- xpack.security.enabled : 会开启集群的node之间的encrypt通信的要求，同时也会开启基于user-password的访问特性

## 命令
```bash
# cd $ES_HOME/bin
# --- elasticsearch-keystore : 管理keystore
# ./elasticsearch-keystore has-passwd # 检查elasticsearch keystore是否有密码保护
# ./elasticsearch-keystore passwd # 设置密码保护
# ./elasticsearch-keystore show keystore.seed
# ./elasticsearch-keystore add user.name # 添加项
# ./elasticsearch-keystore list # 罗列项
...
user.name
# ./elasticsearch-keystore remove user.name # 移除项
# ./elasticsearch-keystore upgrade # 升级elasticsearch keystore
# ./elasticsearch-keystore create # 创建elasticsearch keystore
# -- 设置elasticsearch keystore权限
# chown root:elasticsearch /etc/elasticsearch/elasticsearch.keystore
# chmod 0660 /etc/elasticsearch/elasticsearch.keystore
# echo "demopassword" | ./elasticsearch-keystore add -x "bootstrap.password" # 设置bootstrap.password. `-x`, 从stdin读取密码
# --- elasticsearch-users : 管理user
# ./elasticsearch-users useradd  chencc -p chencc -r superuser
# ./elasticsearch-users list
```

## FAQ
### Cannot allocate memory
调小启动内存
```sh
# vim /etc/elasticsearch/jvm.options

#-Xms2g
#-Xmx2g
-Xms512m
-Xmx512m
```

### 验证elasticsearch是否已运行
`curl -X GET "localhost:9200/?pretty"`

### 访问开启了`xpack.security.http.ssl.enabled:true`的es的方法
1. `curl -k` : 跳过cert检查
1. `curl --cacert /etc/elasticsearch/certs/http_ca.crt` : 用ca cert检查

### 访问es报`unable to authenticate user [elastic] for REST request [/]`
需要带上 authenticate 信息, 比如`curl --user elastic:xxx -XGET 'localhost:9200/_cat/health?v&pretty'`

### `curl "localhost:9200"`返回"Empty reply from server"
权限问题. 比如rpm安装时, 因内存不足导致es无法启动, 手动修改`/etc/elasticsearch/jvm.options` jvm参数为`-Xms512m\n-Xmx1g`并重启后发现使用es 8.3.1安装时console输出的初始密码无法使用.

解决方法:
1. 关闭xpack.security

	修改`/etc/elasticsearch/elasticsearch.yml`将`xpack.security.enabled`设为false, 再重启es即可
2. 使用`curl -u elastic "localhost:9200"`+密码访问
3. 修改(或重置)密码, 再按照方法2访问

### 修改密码
使用elasticsearch-reset-passowrd或elasticsearch-setup-password是可能报`Failed to determine the health of the cluster. Cluster health is currently RED`, 可以加`-f`参数解决.

- 已知密码:

	- 手动修改: `/usr/share/elasticsearch/bin/elasticsearch-setup-passwords interactive`
- 忘记密码

	在es 8.3.1安装后elasticsearch-reset-password无法修改elastic密码.

验证: `curl -u elastic 'http://localhost:9200'`