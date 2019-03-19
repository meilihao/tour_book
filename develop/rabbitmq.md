# rabbitmq

## FAQ
1. 持久化消息
默认情况下的交换机和队列以及消息是非持久化的. 如果消息想要从Rabbitmq崩溃中恢复，那么消息必须满足以下条件:
  1. 把它的投递默认选项设置为持久化
  1. 发送到持久化的交换机
  1. 到达持久化的队列

1. topic如何bind queue
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