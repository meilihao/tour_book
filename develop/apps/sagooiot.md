# sagooiot
ref:
- [sagooiot物联网核心平台](https://opendeep.wiki/sagoo-cloud/sagooiot/mindmap)
- [SagooIoT V1.x到V2.x升级说明](https://iotdoc.sagoo.cn/blog/sagooiot-updatav1tov2)

    所有项目所需仓库和拉取分支说明

## source
### 配置
tdengine是必配的, tsd.database可用TdEngine或Influxdb

### mqtt
topic定义: `pkg/iotModel/sagooProtocol/topic.go`
消费注册入口: `network/core/router.go#RegisterSubTopicHandler()`

    `PropertyRegisterSubRequestTopic`未使用, 只有处理设备上报批量属性的逻辑`network/core/logic/model/up/property/batch/batch.go#core.RegisterSubTopicHandler()`和事件中处理属性的逻辑`network/core/logic/model/up/event/event.go#core.RegisterSubTopicHandler()`->`reporter.ReportProperty(ctx, data)`

    **ps**: 根据官方[演示web](https://zhgy.sagoo.cn/iotmanager/device/instance/WOMEN)设备的`Topic列表`, 属性上报是放在event上报中, 例如`/sys/mn_byq_product/WOMEN/thing/event/property/post`