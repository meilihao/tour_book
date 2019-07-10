# mq
参考:
- [90%的Java程序员，都扛不住这波消息中间件的面试四连炮！](http://www.liuhaihua.cn/archives/587877.html)

使用MQ的场景有挺多的，但是比较核心的有3个：异步、解耦、削峰填谷.

## rabbitmq
rabbitmq解决消息丢失:
1. 持久化机制
1. 关闭RabbitMQ消费者的自动提交ack,在消费者处理完这条消息之后再手动提交ack

## rocketmq

### FAQ
#### broker卡死
内存和swap被耗光导致系统卡死.

解决方法:
修改broker的启动参数
```sh
$ vim vim rocketmq-all-4.3.0-bin-release/bin/runbroker.sh
...
# JAVA_OPT="${JAVA_OPT} -server -Xms8g -Xmx8g -Xmn4g" # 原始参数
JAVA_OPT="${JAVA_OPT} -server -Xms4g -Xmx4g -Xmn1g"
...
```

#### rocketmq-console: org.apache.rocketmq.remoting.exception.RemotingConnectException: connect to <null> failed
解决:
1. `export NAMESRV_ADDR=localhost:9876` 之后再运行`mvn spring-boot:run`
1. 编辑`rocketmq-console/src/main/resources/application.properties`的`rocketmq.config.namesrvAddr=`