# simple-admin
## 部署
1. 先安装pnpm, pg, redis, 并建立core.yaml中DatabaseConf.DBName指定的数据库(否则无法启动)
1. 部署simple-admin-core
```bash
$ mkdir simple-admin
$ git clone -b v1.4.1 --depth 1 https://github.com/suyuan32/simple-admin-core.git
$ cd simple-admin-core
$ go mod tidy
```

需先启动rpc, 再启动api.

rpc:
```bash
$ cd rpc
$ go build
$ vim etc/core.yaml
Name: core.rpc
ListenOn: 0.0.0.0:9101

DatabaseConf:
  Type: postgres
  Host: 127.0.0.1
  Port: 5432
  DBName: simple_admin
  Username: postgres # set your username
  Password: postgres # set your password
  MaxOpenConn: 100
  SSLMode: disable
  CacheTime: 5

Log:
  ServiceName: coreRpcLogger
  Mode: file
  Path: logs/core/rpc
  Encoding: json
  Level: info
  Compress: false
  KeepDays: 7
  StackCoolDownMillis: 100

RedisConf:
  Host: 127.0.0.1:6379

# Prometheus:
#   Host: 0.0.0.0
#   Port: 4001
#   Path: /metrics

CasbinConf:
  ModelText: |
    [request_definition]
    r = sub, obj, act
    [policy_definition]
    p = sub, obj, act
    [role_definition]
    g = _, _
    [policy_effect]
    e = some(where (p.eft == allow))
    [matchers]
    m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act

# Tracing Analysis

#Telemetry:
#  Name: core-rpc
#  Endpoint: localhost:4317
#  Sampler: 1.0
#  Batcher: otlpgrpc # grpc
$ ./rpc -f etc/core.yaml
```

api:
```bash
$ cd api
$ go build
$ vim etc/core.yaml
Name: core.api
Host: 0.0.0.0
Port: 9100
Timeout: 6000

Auth:
  AccessSecret: jS6VKDtsJf3z1n2VKDtsJf3z1n2
  AccessExpire: 259200

Log:
  ServiceName: coreApiLogger
  Mode: file
  Path: logs/core/api
  Level: info
  Compress: false
  KeepDays: 7
  StackCoolDownMillis: 100

Captcha:
  KeyLong: 5
  ImgWidth: 240
  ImgHeight: 80
  Driver: digit

DatabaseConf:
  Type: postgres
  Host: 127.0.0.1
  Port: 5432
  DBName: simple_admin
  Username: postgres # set your username
  Password: postgres # set your password
  MaxOpenConn: 100
  SSLMode: disable
  CacheTime: 5

ProjectConf:
  DefaultRoleId: 1  # default role id when register
  DefaultDepartmentId: 1  # default department id when register
  DefaultPositionId: 1 # default position id when register
  EmailCaptchaExpiredTime: 600 # the expired time for email captcha
  SmsTemplateId: xxx  # The template id for sms
  SmsAppId: xxx  # the app id for sms
  SmsSignName: xxx  # the signature name of sms
  SmsParamsType: json # the type of sms param, support json and array
  RegisterVerify: captcha # register captcha type, support captcha, email, sms, sms_or_email
  LoginVerify: captcha # login captcha type, support captcha, email, sms, sms_or_email, all
  ResetVerify: email # reset captcha type, support email, sms, sms_or_email
  AllowInit: true # if false, ban initialization

# Prometheus:
#   Host: 0.0.0.0
#   Port: 4000
#   Path: /metrics

CasbinConf:
  ModelText: |
    [request_definition]
    r = sub, obj, act
    [policy_definition]
    p = sub, obj, act
    [role_definition]
    g = _, _
    [policy_effect]
    e = some(where (p.eft == allow))
    [matchers]
    m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act

I18nConf:
  Dir: # set the path of locale if you need to load files

CROSConf:
  Address: '*'    # if it is *, allow all ip and address. e.g. http://example.com

RedisConf:
  Host: 127.0.0.1:6379
  Db: 0

CoreRpc:
  Endpoints:
    - 127.0.0.1:9101 # the same as rpc address
  Enabled: true

JobRpc:
  Endpoints:
    - 127.0.0.1:9105
  Enabled: false

McmsRpc:
  Endpoints:
    - 127.0.0.1:9106
  Enabled: false
  Timeout: 5000

# Tracing Analysis

#Telemetry:
#  Name: core-api
#  Endpoint: localhost:4317
#  Sampler: 1.0
#  Batcher: otlpgrpc # grpc
$ ./api -f etc/core.yaml
```

1. 部署simple-admin-backend-ui
```
$ git clone -b v1.4.1 --depth 1 https://github.com/suyuan32/simple-admin-backend-ui.git
$ cd simple-admin-backend-ui
$ pnpm install
$ pnpm serve
```

先访问`http://localhost:5173/init`初始化db, 再访问`http://localhost:5173/`即可.

> backend UI 的代理配置放在 vite.config.ts

1. 部署simple-admin-member-api

```bash
$ git clone -b v1.2.5 --depth 1 https://github.com/suyuan32/simple-admin-member-api.git
$ git clone -b v1.2.5 --depth 1 https://github.com/suyuan32/simple-admin-member-rpc.git
```

rpc:
```bash
$ cd simple-admin-member-rpc
$ go build
$ vim etc/mms.yaml
Name: mms.rpc
ListenOn: 0.0.0.0:9103

DatabaseConf:
  Type: postgres
  Host: 127.0.0.1
  Port: 5432
  DBName: simple_admin
  Username: postgres # set your username
  Password: postgres # set your password
  MaxOpenConn: 100
  SSLMode: disable
  CacheTime: 5

RedisConf:
  Host: 127.0.0.1:6379
  Db: 0

Log:
  ServiceName: mmsRpcLogger
  Mode: file
  Path: logs/mms/rpc
  Encoding: json
  Level: info
  Compress: false
  KeepDays: 7
  StackCoolDownMillis: 100

# Prometheus:
#   Host: 0.0.0.0
#   Port: 4003
#   Path: /metrics
$ ./simple-admin-member-rpc -f etc/mms.yaml
```

api:
```bash
$ cd simple-admin-member-api
$ go build
$ vim etc/mms.yaml
Name: mms.api
Host: 0.0.0.0
Port: 9104
Timeout: 30000

Auth:
  AccessSecret: jS6VKDtsJf3z1n2VKDtsJf3z1n2 # the same as core
  AccessExpire: 259200

CROSConf:
  Address: '*'

Log:
  ServiceName: mmsApiLogger
  Mode: file
  Path: logs/mms/api
  Level: info
  Compress: false
  KeepDays: 7
  StackCoolDownMillis: 100

# Prometheus:
#   Host: 0.0.0.0
#   Port: 4004
#   Path: /metrics

# The database config of casbin
DatabaseConf:
  Type: postgres
  Host: 127.0.0.1
  Port: 5432
  DBName: simple_admin
  Username: postgres # set your username
  Password: postgres # set your password
  MaxOpenConn: 100
  SSLMode: disable
  CacheTime: 5

CasbinConf:
  ModelText: |
    [request_definition]
    r = sub, obj, act
    [policy_definition]
    p = sub, obj, act
    [role_definition]
    g = _, _
    [policy_effect]
    e = some(where (p.eft == allow))
    [matchers]
    m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act

ProjectConf:
  UseCaptcha: false  # If it is false, you can log in without captchaId and captcha fields
  DefaultRankId: 1  # Default rankId when register
  EmailCaptchaExpiredTime: 600 # The expired time for email captcha
  SmsTemplateId: xxx  # The template id for sms
  SmsAppId: xxx  # The app id for sms
  SmsSignName: xxx  # The signature name of sms
  SmsParamsType: json # the type of sms param, support json and array
  RegisterVerify: captcha # verify captcha type, support captcha, email, sms, sms_or_email
  LoginVerify: captcha
  ResetVerify: email # support  email, sms, sms_or_email
  WechatMiniOauthProvider: wechat_mini # the oauth provider for wechat mini program
  AllowInit: true

RedisConf:
  Host: 127.0.0.1:6379
  Type: node

MmsRpc:
  Endpoints:
    - 127.0.0.1:9103
  Enabled: true

CoreRpc:
  Endpoints:
    - 127.0.0.1:9101
  Enabled: true

McmsRpc:
  Endpoints:
    - 127.0.0.1:9106
  Enabled: false
  Timeout: 5000
$ ./simple-admin-member-api -f etc/mms.yaml
```

访问`http://localhost:5173/init`初始化db, 再在`角色管理->编辑角色->权限管理`页面添加 菜单和API接口权限, 最后重新登入即可. 


## FAQ
### `empty etcd hosts`
注释了api/etc/core.yaml里的JobRpc和McmsRpc, 但它们是必须的, 不能注释.

### pnpm serve报`If you want to bypass this version check, you can set the "package-manager-strict" configuration to "false" or set the "COREPACK_ENABLE_STRICT" environment variable to "0"`
本机pnpm版本是9.1.0, 而package.json使用了`"packageManager": "pnpm@9.0.6"`

解决方法: `COREPACK_ENABLE_STRICT=0 pnpm serve`