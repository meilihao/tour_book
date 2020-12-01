# rabbitmq
作用, 在分布式系统下具备异步,削峰,负载均衡等一系列高级功能:
1. 拥有持久化的机制, 队列中的信息也可以保存下来, 防止丢失
1. 实现消费者和生产者之间的解耦
1. 对于高并发场景下，利用消息队列可以使得同步访问变为串行访问达到一定量的限流，利于数据库的操作
1. 可以使用消息队列达到异步请求的效果

场景:
1. 服务间高度解耦, 异步通信
1. 顺序消费
1. 定时任务
1. 流量削峰

## 架构
![RabbitMQ 流程](/misc/img/develop/435188-20180605151314266-1010797270.png)

## 管理

### rabbitmq_management(dashboard)
默认安装了rabbitmq_management, 只需要执行一条命令即可：
```bash
$ sudo rabbitmq-plugins enable rabbitmq_management # 可能需要重启rabbitmq使其生效
```

rabbitmq_management功能:
1. 主页地址: http://server-name:15672/ (默认用户名guest,密码guest)
1. API文档: http://server-name:15672/api/
1. 下载 rabbitmqadmin: http://server-name:15672/cli/

### cmd
1. rabbitmqctl list_queues : 查看mq上的queue
1. rabbitmqctl list_users : 已有帐号, 缺省的guest账户只能在本地豋录, `[xxx]`表示role
1. rabbitmqctl change_password  Username  'Newpassword'
1. rabbitmqctl status : 状态
1. rabbitmqctl list_connections : 连接状态

## 配置
配置文件位置: `/etc/rabbitmq`
默认log文件位置: `/var/lib/rabbitmq`

> `http://server-name:15672/#/`上就有相关信息.

### 用户权限管理
权限管理中主要包含四步：
1. 新建用户

	rabbitmqctl add_user rabbitadmin 123456 : 新建帐号
1. 配置权限

	rabbitmqctl set_permissions -p / rabbitadmin ".*" ".*" ".*" # set_permissions [-p <vhostpath>] <user> <conf> <write> <read>

	<conf> <write> <read>的位置分别用正则表达式来匹配特定的资源，如'^(amq\.gen.*|amq\.default)$'可以匹配server生成的和默认的exchange，'^$'不匹配任何资源

	权限细则:
    - exchange和queue的declare与delete分别需要exchange和queue上的配置权限
    - exchange的bind与unbind需要exchange的读写权限
    - queue的bind与unbind需要queue写权限exchange的读权限 发消息(publish)需exchange的写权限
    - 获取或清除(get、consume、purge)消息需queue的读权限

    其他:
    - sudo rabbitmqctl list_user_permissions rabbitadmin : 查看用户权限

1. 配置角色

	rabbitmqctl set_user_tags rabbitadmin administrator : 分配用户标签, 设置为管理用户, 如果未设置权限, 此时还是不能访问virtual hosts, 该情况可在dashboard的admin tab上看到.

	权限区别:
	- none

		不能访问 management plugin

    - management

	    用户可以通过AMQP做的任何事外加：
	    列出自己可以通过AMQP登入的virtual hosts
	    查看自己的virtual hosts中的queues, exchanges 和 bindings
	    查看和关闭自己的channels 和 connections
	    查看有关自己的virtual hosts的“全局”的统计信息，包含其他用户在这些virtual hosts中的活动

    - policymaker

	    management可以做的任何事外加：
	    查看、创建和删除自己的virtual hosts所属的policies和parameters

    - monitoring

	    management可以做的任何事外加：
	    列出所有virtual hosts，包括他们不能登录的virtual hosts
	    查看其他用户的connections和channels
	    查看节点级别的数据如clustering和memory使用情况
	    查看真正的关于所有virtual hosts的全局的统计信息

    - administrator

	    policymaker和monitoring可以做的任何事外加:
	    创建和删除virtual hosts
	    查看、创建和删除users
	    查看创建和删除permissions
	    关闭其他用户的connections
1. 删除用户

	rabbitmqctl.bat delete_user rabbitadmin

## FAQ
### 持久化消息
默认情况下的交换机和队列以及消息是非持久化的. 如果消息想要从Rabbitmq崩溃中恢复，那么消息必须满足以下条件:
  1. 把它的投递默认选项设置为持久化
  1. 发送到持久化的交换机
  1. 到达持久化的队列

### topic如何bind queue
```go
// producer
err = ch.Publish(
		"example.topic", // exchange
		routingKey,     // routing key
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, //消息持久化
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
```

```go
// consumer
err = ch.QueueBind(
		q.Name, // queue name
		bindingKey,     // binding key
		"example.topic", // exchange
		false,
		nil,
	)
```

> bindingKey: `#`表示0个或若干个词，`*`表示一个词.
> ![Topic Exchange](https://img-blog.csdn.net/20150914161921517)
### 如何确保消息正确地发送至 RabbitMQ？ 如何确保消息接收方消费了消息？
- 发送方确认模式

	将信道设置成 confirm 模式（发送方确认模式），则所有在信道上发布的消息都会被指派一个唯一的 ID。
	一旦消息被投递到目的队列后，或者消息被写入磁盘后（可持久化的消息），信道会发送一个确认给生产者（包含消息唯一 ID）。
	如果 RabbitMQ 发生内部错误从而导致消息丢失，会发送一条 nack（notacknowledged，未确认）消息。

	发送方确认模式是异步的，生产者应用程序在等待确认的同时，可以继续发送消息。当确认消息到达生产者应用程序，生产者应用程序的回调方法就会被触发来处理确认消息。

- 接收方确认机制

	消费者接收每一条消息后都必须进行确认（消息接收和消息确认是两个不同操作）。只有消费者确认了消息，RabbitMQ 才能安全地把消息从队列中删除。
	这里并没有用到超时机制，RabbitMQ 仅通过 Consumer 的连接中断来确认是否需要重新发送消息。也就是说，只要连接不中断，RabbitMQ 给了 Consumer 足够长的时间来处理消息。保证数据的最终一致性；

下面罗列几种特殊情况:
1. 如果消费者接收到消息，在确认之前断开了连接或取消订阅，RabbitMQ 会认为消息没有被分发，然后重新分发给下一个订阅的消费者。（可能存在消息重复消费的隐患，需要去重）
1. 如果消费者接收到消息却没有确认消息，连接也未断开，则 **RabbitMQ 认为该消费者繁忙，将不会给该消费者分发更多的消息**。

### 如何避免消息重复投递或重复消费？

在消息生产时，MQ 内部针对每条生产者发送的消息生成一个 inner-msg-id，作为去重的依据（消息投递失败并重传），避免重复的消息进入队列；在消息消费时，要求消息体中必须要有一个 bizId（对于同一业务全局唯一，如支付 ID、订单 ID、帖子 ID 等）作为去重的依据，避免同一条消息被重复消费。

### 消息基于什么传输？

由于 TCP 连接的创建和销毁开销较大，且并发数受系统资源限制，会造成性能瓶颈。RabbitMQ 使用信道的方式来传输数据。信道是建立在真实的 TCP 连接内的虚拟连接，且每条 TCP 连接上的信道数量没有限制。

### 消息如何分发？

若该队列至少有一个消费者订阅，消息将以循环（round-robin）的方式发送给消费者。每条消息只会分发给一个订阅的消费者（前提是消费者能够正常处理消息并进行确认）。通过路由可实现多消费的功能

### 消息怎么路由？

消息提供方->路由->一至多个队列消息发布到交换器时，消息将拥有一个路由键（routing key），在消息创建时设定。通过队列路由键，可以把队列绑定到交换器上。消息到达交换器后，RabbitMQ 会将消息的路由键与队列的路由键进行匹配（针对不同的交换器有不同的路由规则）；

常用的交换器主要分为一下三种：
1. fanout：如果交换器收到消息，将会广播到所有绑定的队列上
1. direct：如果路由键完全匹配，消息就被投递到相应的队列
1. topic：可以使来自不同源头的消息能够到达同一个队列。 使用 topic 交换器时，可以使用通配符

### 如何确保消息不丢失？

消息持久化，当然前提是队列必须持久化

RabbitMQ 确保持久性消息能从服务器重启中恢复的方式是，将它们写入磁盘上的一个持久化日志文件，当发布一条持久性消息到持久交换器上时，Rabbit 会在消息提交到日志文件后才发送响应。一旦消费者从持久队列中消费了一条持久化消息，RabbitMQ 会在持久化日志中把这条消息标记为等待垃圾收集。如果持久化消息在被消费之前 RabbitMQ 重启，那么 Rabbit 会自动重建交换器和队列（以及绑定），并重新发布持久化日志文件中的消息到合适的队列。

### RabbitMQ 的集群

镜像集群模式

你创建的 queue，无论元数据还是 queue 里的消息都会存在于多个实例上，然后每次你写消息到 queue 的时候，都会自动把消息到多个实例的 queue 里进行消息同步。

好处在于，你任何一个机器宕机了，没事儿，别的机器都可以用。坏处在于，第一，这个性能开销也太大了吧，消息同步所有机器，导致网络带宽压力和消耗很重！第二，这么玩儿，就没有扩展性可言了，如果某个 queue 负载很重，你加机器，新增的机器也包含了这个 queue 的所有数据，并没有办法线性扩展你的 queue.

### mq 的缺点
1. 系统可用性降低

系统引入的外部依赖越多，越容易挂掉，本来你就是 A 系统调用 BCD 三个系统的接口就好了，人 ABCD 四个系统好好的，没啥问题，你偏加个 MQ 进来，万一MQ 挂了咋整？MQ 挂了，整套系统崩溃了.

1. 系统复杂性提高

加个 MQ 进来，你怎么保证消息没有重复消费？怎么处理消息丢失的情况？怎么保证消息传递的顺序性？头大头大，问题一大堆，痛苦不已

1. 一致性问题

A 系统处理完了直接返回成功了，人都以为你这个请求就成功了；但是问题是，要是 BCD 三个系统那里，BD 两个系统写库成功了，结果 C 系统写库失败了，这数据就不一致了。

## FAQ
### epmd error for host xxx：address (cannot connect to host/port)
`/etc/hosts`错误导致解析host报错, 是之前修改ip导致, 修正即可.

### rabbit 3.5 启动`no proc`
`/var/log/rabbitmq/start_log`日志:
```
BOOT FAILED
===========

Error description:
   noproc

Log files (may contain more information):
   /var/log/rabbitmq/rabbit@localhost.log
   /var/log/rabbitmq/rabbit@localhost-sasl.log
...
```

经搜索应该是rabbitmq 3.5与erlang 22不兼容导致, 直接升级rabbitmq到`3.8`.

### rabbit 3.8.9 gust登录报错: `user can only log in via localhost`
rabbitmq从3.3.0开始禁止使用guest/guest权限通过除localhost外的访问. 建议参考本文的`用户权限管理`添加自定义帐号.

### 执行`rabbitmq-plugins enable rabbitmq_management`报`cannot_read_enabled_plugins_file /etc/rabbitmq/enabled_plugins eacces`
因为`/etc/rabbitmq/enabled_plugins`的权限是`-rw-------`, 执行
`umask 0022; rabbitmq-plugins enable rabbitmq_management`即可, 执行后权限变为`-rw-r--r--`