# mqtt
ref:
- [**MQTT 教程**](https://www.emqx.com/zh/mqtt-guide)
- [MQTT 协议快速入门 2025：基础知识和实用教程](https://www.emqx.com/zh/blog/the-easiest-guide-to-getting-started-with-mqtt)
- [MQTT 服务器（MQTT Broker）：工作原理与快速入门指南](https://www.emqx.com/zh/blog/the-ultimate-guide-to-mqtt-broker-comparison)
- [MQTT 客户端库 & SDKs](https://www.emqx.com/zh/mqtt-client-sdk)
- [物联网首选协议，关于 MQTT 你需要了解这些](https://www.emqx.com/zh/blog/what-is-the-mqtt-protocol)
- [**MQTT 服务器**](https://www.mqtt.cn/mqtt-servers)

    mqtt学习网站

MQTT是一种轻量级、基于发布-订阅模式的消息传输协议，适用于资源受限的设备和低带宽、高延迟或不稳定的网络环境。它在物联网应用中广受欢迎，能够实现传感器、执行器和其它设备之间的高效通信.

mqtt特点:
1. 轻量级：物联网设备通常在处理能力、内存和能耗方面受到限制。MQTT 开销低、报文小的特点使其非常适合这些设备，因为它消耗更少的资源，即使在有限的能力下也能实现高效的通信。
1. 可靠：物联网网络常常面临高延迟或连接不稳定的情况。MQTT 支持多种 QoS 等级、会话感知和持久连接，即使在困难的条件下也能保证消息的可靠传递，使其非常适合物联网应用。
1. 安全通信：MQTT 提供传输层安全（TLS）和安全套接层（SSL）加密功能。此外，MQTT 还通过用户名/密码凭证或客户端证书提供身份验证和授权机制，以保护网络及其资源的访问。
1. 双向通信：MQTT 的发布-订阅模式为设备之间提供了无缝的双向通信方式。客户端既可以向主题发布消息，也可以订阅接收特定主题上的消息，从而实现了物联网生态系统中的高效数据交换，而无需直接将设备耦合在一起。这种模式也简化了新设备的集成，同时保证了系统易于扩展。
1. 连续、有状态的会话：MQTT 提供了客户端与 Broker 之间保持有状态会话的能力，这使得系统即使在断开连接后也能记住订阅和未传递的消息。此外，客户端还可以在建立连接时指定一个保活间隔，这会促使 Broker 定期检查连接状态。如果连接中断，Broker 会储存未传递的消息（根据 QoS 级别确定），并在客户端重新连接时尝试传递它们。这个特性保证了通信的可靠性，降低了因间断性连接而导致数据丢失的风险。
1. 大规模物联网设备支持：MQTT 的轻量级特性、低带宽消耗和对资源的高效利用使其成为大规模物联网应用的理想选择。通过采用发布-订阅模式，MQTT 实现了发送者和接收者的解耦，从而有效地减少了网络流量和资源使用。此外，协议对不同 QoS 等级的支持使得消息传递可以根据需求进行定制，确保在各种场景下获得最佳的性能表现。
1. 语言支持：物联网系统包含使用各种编程语言开发的设备和应用。MQTT 具有广泛的语言支持，使其能够轻松与多个平台和技术进行集成，从而实现了物联网生态系统中的无缝通信和互操作性

## lib
ref:
- [libraries](https://github.com/mqtt/mqtt.org/wiki/libraries)

## 工作原理
MQTT 是基于发布-订阅模式的通信协议，由 MQTT 客户端通过主题（Topic）发布或订阅消息，通过 MQTT Broker 集中管理消息路由，并依据预设的服务质量等级(QoS)确保端到端消息传递可靠性.

### 组件
1. MQTT 客户端

    任何运行 MQTT 客户端库的应用或设备都是 MQTT 客户端
1. MQTT Broker

    MQTT Broker 是负责处理客户端请求的关键组件，包括建立连接、断开连接、订阅和取消订阅等操作，同时还负责消息的转发

### QoS
ref:
- [MQTT QoS 0、1、2 解析：快速入门指南](https://www.emqx.com/zh/blog/introduction-to-mqtt-qos)

MQTT 提供了三种服务质量（QoS），在不同网络环境下保证消息的可靠性:

1. QoS 0：最多交付一次

    如果当前客户端不可用，它将丢失这条消息。
2. QoS 1：至少交付一次

    包含了简单的重发机制，发布者发送消息之后等待接收者的 ACK，如果没收到 ACK 则重新发送消息。这种模式能保证消息至少能到达一次，但无法保证消息重复
3. QoS 2：只交付一次

    设计了重发和重复消息发现机制，保证消息到达对方并且严格只到达一次
    
    与 QoS 1 相比，QoS 2 新增了 PUBREL 报文和 PUBCOMP 报文的流程，也正是这个新增的流程带来了消息不会重复的保证. 因此每一次的 QoS 2 消息投递，都要求发送方与接收方进行至少两次请求/响应流程. 

### MQTT 的工作流程
1. 客户端使用 TCP/IP 协议与 Broker 建立连接
1. 客户端既可以向特定主题发布消息，也可以订阅主题以接收消息
1. MQTT Broker 接收发布的消息，并将这些消息转发给订阅了对应主题的客户端

### topic
ref:
- [MQTT 主题与通配符（Topics & Wildcards）入门手册](https://www.emqx.com/zh/blog/advanced-features-of-mqtt-topics)

MQTT 协议根据主题来转发消息。主题通过 / 来区分层级，类似于 URL 路径.

MQTT 主题支持以下两种通配符:
1. `+`：表示单层通配符，例如 `a/+` 匹配 `a/x` 或 `a/y`

    `+`用于单个主题层级匹配的通配符. 在使用单层通配符时，单层通配符必须占据整个层级
2. `#`：表示多层通配符，例如 `a/#` 匹配 `a/x`、`a/b/c/d`

    `#`是用于匹配主题中任意层级的通配符. 多层通配符表示它的父级和任意数量的子层级，在使用多层通配符时，它必须占据整个层级并且必须是主题的最后一个字符

注意: 通配符主题**只能用于订阅，不能用于发布**

example:
```conf
+ 有效
sensor/+ 有效
sensor/+/temperature 有效
sensor+ 无效（没有占据整个层级）

# 有效，匹配所有主题
sensor/# 有效
sensor/bedroom# 无效（没有占据整个层级）
sensor/#/temperature 无效（不是主题最后一个字符）
```

#### 系统主题
以 `$SYS/` 开头的主题为系统主题，系统主题主要用于获取 MQTT 服务器自身运行状态、消息统计、客户端上下线事件等数据。目前，MQTT 协议暂未明确规定 $SYS/ 主题标准，但大多数 MQTT 服务器都遵循该标准建议。

例如，EMQX 服务器支持通过以下主题获取集群状态。

主题	说明
$SYS/brokers	EMQX 集群节点列表
$SYS/brokers/emqx@127.0.0.1/version	EMQX 版本
$SYS/brokers/emqx@127.0.0.1/uptime	EMQX 运行时间
$SYS/brokers/emqx@127.0.0.1/datetime	EMQX 系统时间
$SYS/brokers/emqx@127.0.0.1/sysdescr	EMQX 系统信息
EMQX 还支持客户端上下线事件、收发流量、消息收发、系统监控等丰富的系统主题，用户可通过订阅 $SYS/# 主题获取[所有系统主题消息](https://docs.emqx.com/zh/emqx/v5.0/observability/mqtt-system-topics.html#%E5%AE%A2%E6%88%B7%E7%AB%AF%E4%B8%8A%E4%B8%8B%E7%BA%BF%E4%BA%8B%E4%BB%B6)

#### 共享订阅
共享订阅是 MQTT 5.0 引入的新特性，用于在多个订阅者之间实现订阅的**负载均衡**，MQTT 5.0 规定的共享订阅主题以 $share 开头。

> 虽然 MQTT 协议在 5.0 版本才引入共享订阅，但是 EMQX 从 MQTT 3.1.1 版本开始就支持共享订阅

## tools
- hbmqtt :  hbmqtt_sub 和 hbmqtt_pub

  ```bash
  # hbmqtt_sub --url mqtt://mqtt.eclipse.org:1883 -t /geektime/iot
  # hbmqtt_pub --url mqtt://mqtt.eclipse.org:1883 -t /geektime/iot -m Hello,World!
  ```

## EMQX
访问管理控制台： 通过 http://localhost:18083/ 登录，默认账号: admin 密码: public

## mosquitto
```bash
# mosquitto_sub -t /db/change # 订阅topic
```

## FAQ
### 程序行为异常, 定位是mqtt问题
1. 检查是否使用了相同的clientid

### mqtt client: `not currently connected and ResumeSubs not set`
多个client使用了相同的clientid