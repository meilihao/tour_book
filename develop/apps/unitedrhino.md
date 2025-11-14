# 联犀
ref:
- [mqtt认证主要流程](https://doc.unitedrhino.com/use/fe695a/)

## 部署
### docker
ref:
- [安装教程](https://doc.unitedrhino.com/use/046431/)

运行things/deploy/docker/run-all/docker-compose.yml时, docker compose会引入同级目录下的`.env`用于注入docker-compose.yml所需的环境变量.

> 获取things启动日志: 进入things容器, 执行`./thingsEEsvr -h 2>&1 |tee -a s.log`

apisvr实际加载的配置是`deploy/docker/conf/things/etc/commonDocker.yaml` + `deploy/docker/conf/things/etc/api.yaml`, 就是将api.yaml的内容拼接在commonDocker.yaml后面, 加载核心逻辑见`/data/tmpfs/unitedrhino/share/utils/conf.go#Load()`. 因此修改log level就是修改commonDocker.yaml的`Log`配置即可

apisvr的NewServiceContext()会根据`c.DmRpc.Enable`选择用服务模式还是直连模式.

## 设备接入
步骤:
1. 创建设备

### mqtt
ref:
- [设备接入指引](https://doc.unitedrhino.com/use/adfab2/)
- [物模型协议](https://doc.unitedrhino.com/use/b4c708/)

	包含订阅参数说明

步骤:
1. 在设备详情页点击`生成MQTT三元组`([原理](https://doc.unitedrhino.com/use/93d0ca/))生成mqtt参数
2. 用mqttx连接到emqx
3. 发布数据

	```log
	订阅响应topic: $thing/down/property/69/temp1
	
	发布topic: $thing/up/property/69/temp1

	{                     
	    "method": "report",            
	    "msgToken": "afwegafeegfa",   
	    "params": { 
	      "Temp": 1
	    }
	}
	```

	> 订阅topic校验逻辑在`service/dgsvr/internal/logic/deviceauth/accessAuthLogic.go#deviceAuth.AccessAuth()`

## 数据处理
### 设备上报
入口在`service/dgsvr/internal/repo/event/subscribe/subDev/mqtt.go#handle(ctx, message.Topic(), message.Payload())`, 具体处理函数是`service/dmsvr/internal/event/deviceMsgEvent/deviceMsg.go#(l *DeviceMsgHandle) Thing()`->`service/dmsvr/internal/domain/protocol/script.go#(s *ScriptTrans) UpBeforeTrans()`->`service/dmsvr/internal/event/deviceMsgEvent/thing.go#(l *ThingLogic) Handle()`

replay逻辑:
1. `service/dmsvr/internal/event/deviceMsgEvent/deviceMsg.go#(l *DeviceMsgHandle) Thing()`里的`l.deviceResp(resp)`->`l.svcCtx.PubDev.PublishToDev(l.ctx, respMsg)`
1. `service/dmsvr/internal/repo/event/publish/pubDev/nats.go#(n *pubDevClient) PublishToDev(ctx context.Context, reqMsg *deviceMsg.PublishMsg)`

添加mqtt handle流程:
1. `service/dgsvr/dgdirect/direct.go#GetSvcCtx()`的`startup.PostInit(svcCtx)`
1. `service/dgsvr/internal/startup/init.go#PostInit()`->`sd.SubDevMsg()`
1. `service/dgsvr/internal/repo/event/subscribe/subDev/mqtt.go#(d *MqttClient) SubDevMsg(handle Handle)`-> `d.subscribeWithFunc(cli, TopicThing,...)`