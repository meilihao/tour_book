# 联犀
ref:
- [mqtt认证主要流程](https://doc.unitedrhino.com/use/fe695a/)

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