# tc
ref:
- [linux 网卡限速（利用tc，iptables limit模块等）](https://blog.51cto.com/linzb/1766218)
- [Tc 流控 HOWTO 文档](http://www.tldp.org/HOWTO/html_single/Traffic-Control-HOWTO/)
- [《Linux 高级路由与流量控制手册（2012）》第九章：用 tc qdisc 管理 Linux 网络带宽](https://arthurchiao.art/blog/lartc-qdisc-zh/)
- [利用Cgroup限制网络带宽](https://guanjunjian.github.io/2017/11/29/study-14-cgroup-network-control-group/)
- [wondershaper - tc的shell wrapper]

TC 是linux自带的模块，可以用来控制网速, from `yum install iproute`.

tc/qdisc 是 Cilium/eBPF 依赖的最重要的网络基础设施之一.

Linux 流量控制方法: 控发不控收(因linux 对接收队列的控制不够好), 所以只能对产生瓶颈网卡处的发包速率进行控制, 具体是通过队列 (queueing discipline)来实现控制网络的收发速度.

在 linux 中,TC 有二种控制方法 CBQ(Class Based Queueing，基于类的排队) 和 HTB(Hierarchical Token Bucket，层级令牌桶). HTB 是设计用来替换 CBQ 的. 它是一个层次式的过滤框架.

Martin Devera (devik) 意识到CBQ 太复杂了，并且没有针对很多典型的场景进 行优化. 因此他设计了 HTB，这种层级化的方式对下面这些场景很适用：
- 有一个固定总带宽，想将其分割成几个部分，分别用作不同目的
- 每个部分的带宽是有保证的（guaranteed bandwidth）
- 还可以指定每个部分向其他部分借带宽

## 选项
格式: `tc [qdisc/class/filter] [add/del/replace] dev 网卡名字 其他参数`

**tc限速主要是将数据包发送到不同类型的队列中**，然后由队列控制发送. 限速队列主要由两种:
1. 一种是无类队列, 可以对某个网络 接口（interface）上的所有流量进行无差别整形, 包括对对数据进行`重新调度（reschedule）`,`增加延迟（delay）`,`丢弃（drop）`.

	其中包括以下:
	- pfifo_fast(先进先出)

		对所有包都一视同仁.

		pfifo_fast 有三个所谓的 “band”（可理解为三个队列），编号分别为 0、1、2：

		- 每个 band 上分别执行 FIFO 规则
		- 但是，如果 band 0 有数据，就不会处理 band 1；同理，band 1 有数据时， 不会去处理 band 2
		- 内核会检查数据包的 TOS 字段，将“最小延迟”的包放到 band 0

	- TBF ( 令牌桶过滤器)

		TBF 是一个简单 qdisc，对于没有超过预设速率的流量直接透传，但也能容忍超过预 设速率的短时抖动（short bursts in excess of this rate）。 TBF 非常简洁，对网络和处理器都很友好（network- and processor friendly）。 如果只是想实现接口限速，那 TBF 是第一选择.

	- SFQ(随机公平队列) 

		随机公平排队（SFQ）是公平排队算法族的一个简单实现。相比其他算法，SFQ 精准性要差 一些，但它所需的计算量也更少，而结果几乎是完全公平的（almost perfectly fair）.
	- ID (前 向随机丢包)
	- etc.


	选择使用哪种 qdisc 时，下面几点可供参考:

	- 单纯对出向流量限速（slow down outgoing traffic），推荐使用 TBF。如果是 针对大带宽进行限速，需要将 bucket 调大。
	- 如果带宽已经打满，想确保带宽没有被任何单个 session 占据，推荐使用 SFQ
	- If you have a big backbone and know what you are doing, consider Random Early Drop (see Advanced chapter).
	- 对（不再转发的）入向流量整形，使用 Ingress Policer。顺便说一句，入向整形称为 ‘policing’，而不是 ‘shaping’。
	- 对需要本机转发的流量整形

		- 如果目的端是单个设备，那在目的端设备上使用 TBF。
		- 如果目的端是多个设备（同一个入向设备分流到多个目的设备），使用 Ingress Policer。
	- 如果你不需要整形，只是想看看网络接口（interface）是否过载（so loaded that it has to queue）， 使用 pfifo queue（注意不是 pfifo_fast）。pfifo 内部没有 bands，但会记录 backlog 的大小。

1. 另外一种是分类队列, 一个排队规则中又包含其他 排队规则（qdisc-containing-qdiscs）.

	其中由引出了class(类, 表示控制策略)，filter(Classifiers, 过滤器, 用来将用户划入到具体的控制策略中即不同的 class 中)的概念.

	**如果想对不同类型的流量做不同处理，那 classful qdisc 非常有用**

	TC 可以使用的过滤器有:
	- fwmark 分类器
	- u32 分类器
	- 基于路由的分类器
	- RSVP 分类器 (分别用于 IPV6、IPV4) 
	- etc.

	其中，fwmark 分类器允许我们使用 Linux netfilter 代码选择流量，而 u32 分类器允许我们选择基于 ANY 头的流量 . 需要注意的是，filter(过滤器) 是在 QDisc 内部，它们不能作为主体

qdisc(queueing discipline, 排队规则),class,filter三者直接关系如下：
每创建一个class，都会有一个默认的qdisc,该qdisc挂在class作为子节点上。filter挂在队列上，主要决定让数据包向子节点类传递.


无论是队列，还是 class 和 filter 都有 ID 之类的标志符，一般都有 parent(父，上层的)，注意 ID 具有接口本地性，不同的网络接口可以有相同的 ID.

## 术语
- Queueing Discipline (qdisc，排队规则)

	管理设备队列（queues of devices）的算法，可以是管理入向（incoing/ingress ）队列，也可以是管理出向队列（outgoing/egress）。

- root qdisc（根排队规则）

- attach 到网络设备的那个 qdisc。

- Classless qdisc（无类别排队规则）

	对所有包一视同仁，同等对待。

- Classful qdisc（有类别排队规则）

	一个 classful qdisc 会包含多个类别（classes）。每个类别（class）可以进一步包 含其他 qdisc，可以是 classful qdisc，也可以是 classless qdisc。

	严格按定义来说，pfifo_fast 属于有类别排队规则（classful），因为它内部包 含了三个 band，而这些 band 实际上是 class。但从用户配置的视角来说，它是 classless 的，因为这三个内部 class 用户是无法通过 tc 命令配置的。

- Classes（类别）

	每个 classful qdisc 可能会包含几个 class，这些都是 qdisc 内部可见的。对于每 个 class，也是可以再向其添加其他 class 的。因此，一个 class 的 parent 可以 是一个 qdisc，也可以是另一个 class。

	Leaf class 是没有 child class 的 class。这种 class 中 attach 了一个 qdisc ，负责该 class 的数据发送。

	创建一个 class 时会自动 attach 一个 fifo qdisc。而当向这个 class 添加 child class 时，这个 fifo qdisc 会被自动删除。对于 leaf class，可以用一个更合适的 qdisc 来替换掉这个fifo qdisc。你甚至能用一个 classful qdisc 来替换这个 fifo qdisc，这样就可以添加其他 class了。

- Classifier（分类器）

	每个 classful qdisc 需要判断每个包应该放到哪个 class。这是通过分类器完成的。

- Filter（过滤器）

	分类过程（Classification）可以通过过滤器（filters）完成。过滤器包含许多的判 断条件，匹配到条件之后就算 filter 匹配成功了。

- Scheduling（调度）

	在分类器的协助下，一个 qdisc 可以判断某些包是不是要先于其他包发送出去，这 个过程称为调度，可以通过例如前面提到的 pfifo_fast qdisc 完成。调度也被 称为重排序（reordering），但后者容易引起混淆。

- Shaping（整形）

	在包发送出去之前进行延迟处理，以达到预设的最大发送速率的过程。整形是在 egress 做的（前面提到了，ingress 方向的不叫 shaping，叫 policing，译者注）。 不严格地说，丢弃包来降低流量的过程有时也称为整形。

- Policing（执行策略，决定是否丢弃包）

	延迟或丢弃（delaying or dropping）包来达到预设带宽的过程。 在 Linux 上， policing 只能对包进行丢弃，不能延迟 —— 没有“入向队列”（”ingress queue”）。

- Work-Conserving qdisc（随到随发 qdisc）

	work-conserving qdisc 只要有包可发送就立即发送。换句话说，只要网卡处于可 发送状态（对于 egress qdisc 来说），它永远不会延迟包的发送。

- non-Work-Conserving qdisc（非随到随发 qdisc）

	某些 qdisc，例如 TBF，可能会延迟一段时间再将一个包发送出去，以达到期望的带宽 。这意味着它们有时即使有能力发送，也不会发送。


```
                Userspace programs
                     ^
                     |
     +---------------+-----------------------------------------+
     |               Y                                         |
     |    -------> IP Stack                                    |
     |   |              |                                      |
     |   |              Y                                      |
     |   |              Y                                      |
     |   ^              |                                      |
     |   |  / ----------> Forwarding ->                        |
     |   ^ /                           |                       |
     |   |/                            Y                       |
     |   |                             |                       |
     |   ^                             Y          /-qdisc1-\   |
     |   |                            Egress     /--qdisc2--\  |
  --->->Ingress                       Classifier ---qdisc3---- | ->
     |   Qdisc                                   \__qdisc4__/  |
     |                                            \-qdiscN_/   |
     |                                                         |
     +----------------------------------------------------------+

Thanks to Jamal Hadi Salim for this ASCII representation.
```

上图中的框代表 Linux 内核。最左侧的箭头表示流量从外部网络进入主机, 然后进入 Ingress Qdisc，这里会对包进行过滤（apply Filters），根据结果决定是否要丢弃这个包。这个过程称为 “Policing”. 这个过程发生在内核处理的很早阶段，在穿过大部分内核基础设施之前。因此在这里丢弃包是很高效的，不会消耗大量 CPU.

如果判断允许这个包通过，那它的目的端可能是本机上的应用（local application），这种情况下它会进入内核 IP 协议栈进行进一步处理，最后交给相应的用户态程序。另外，这 个包的目的地也可能是其他主机上的应用，这种情况下就需要通过这台机器 Egress Classifier 再发送出去。主机程序也可能会发送数据，这种情况下也会通过 Egress Classifier 发送.

Egress Classifier 中会用到很多 qdisc. 默认情况下只有一个：pfifo_fast qdisc ，它永远会接收包，这称为“入队”（”enqueueing”）。

此时包位于 qdisc 中了，等待内核召唤，然后通过网络接口（network interface）发送出去。 这称为“出队”（”dequeueing”）。

以上画的是单网卡的情况。在多网卡的情况下，每个网卡都有自己的 ingress 和 egress hooks.

## qdisc
- 每个接口都有一个 egress "root qdisc"。默认情况下，这个 root qdisc 就是前面提到的 classless pfifo_fast qdisc.

	理论上 egress 流量是本机可控的，所以需要配备一个 qdisc 来提供这种控制能力

- 每个 qdisc 和 class 都会分配一个相应的 handle（句柄），可以指定 handle 对 qdisc 进行配置。
- 每个接口可能还会有一个 ingress qdisc，用来对入向流量执行策略（which polices traffic coming in）。

	理论上 ingress 基本是不受本机控制的，主动权在外部，所以不一定会有 qdisc.

关于 handle：
- 每个 handle 由两部分组成: `<major>:<minor>`
- 按照惯例，root qdisc 的 handle 为`1:`, 这是`1:0`的简写.
- 每个 qdisc 的 minor number 永远是 0

关于 class：

- 每个 class 的 major number 必须与其 parent 一致。
- major number 在一个 egress 或 ingress 内必须唯一。
- minor number 在一个 qdisc 或 class 内必须唯一。

上面的解释有点模糊，可对照`man 8 tc`的解释：

所有 qdiscs、classes 和 filters 都有 ID，这些 ID 可以是指定的，也可以是自动分的。

ID 格式 major:minor，major 和 minor 都是 16 进制数字，不超过 2 字节。 两个特殊值：

	- root 的 major 和 minor 初始化全 1。
	- 省略未指定的部分将为全 0。

下面分别介绍以上三者的 ID 规范:

- qdisc：qdisc 可能会有 children。

	- major 部分：称为 handle，表示的 qdisc 的唯一性
	- minor 部分：留给 class 的 namespace
- class：class 依托在 qdisc 内

	- major 部分：继承 class 所在的 qdisc 的 major
	- minor 部分：称为 classid，在所在的 qdisc 内唯一就行
- filter：由三部分构成，只有在使用 hashed filter hierarchy 时才会用到


一个典型的 handle 层级如下:
```
                     1:   root qdisc
                      |
                     1:1    child class
                   /  |  \
                  /   |   \
                 /    |    \
                 /    |    \
              1:10  1:11  1:12   child classes
               |      |     |
               |     11:    |    leaf class
               |            |
               10:         12:   qdisc
              /   \       /   \
           10:1  10:2   12:1  12:2   leaf classes
```

一个包可能会被链式地分类如下（get classified in a chain）：`1: -> 1:1 -> 1:12 -> 12: -> 12:2`, 最后到达 attach 到 class 12:2 的 qdisc 的队列。在这个例子中，树的每个“节点”（ node）上都 attach 了一个 filter，每个 filter 都会给出一个判断结果，根据判断结果 选择一个合适的分支将包发送过去。这是常规的流程。

但下面这种流程也是有可能的：`1: -> 12:2`, 在这种情况下，attach 到 root qdisc 的 filter 决定直接将包发给 `12:2`


### 包是如何从 qdisc 出队（dequeue）然后交给硬件的
当内核决定从 qdisc dequeue packet 交给接口（interface）发送时，它会

- 向 root qdisc 1: 发送一个 dequeue request
- 1: 会将这个请求转发给 1:1，后者会进一步向下传递，转发给 10:、11:、12:
- 每个 qdisc 会查询它们的 siblings，并尝试在上面执行 dequeue() 方法

总结: **嵌套类（nested classes）只会和它们的 parent qdiscs 通信，而永远不会直 接和接口交互。内核只会调用 root qdisc 的 dequeue() 方法**

## HTB
![](/misc/img/net/tc/htb-borrow.png)

HTB 的工作方式与 CBQ 类似，但不是借助于计算空闲时间（idle time）来实现整形。 在内部，它其实是一个 classful TBF（令牌桶过滤器）—— 这也是它叫层级令牌桶（HTB） 的原因.

## tc命令
tc可以使用以下命令对QDisc、类和过滤器进行操作：

- add：
 
 	在一个节点里加入一个QDisc、类或者过滤器。添加时，需要传递一个祖先作为参数，传递参数时既可以使用ID也可以直接传递设备的根。如果要建立一个QDisc或者过滤器，可以使用句柄(handle)来命名；如果要建立一个类，可以使用类识别符(classid)来命名。
- remove：
 
 	删除有某个句柄(handle)指定的QDisc，根QDisc(root)也可以删除。被删除QDisc上的所有子类以及附属于各个类的过滤器都会被自动删除。
- change：
 
 	以替代的方式修改某些条目。除了句柄(handle)和祖先不能修改以外，change命令的语法和add命令相同。换句话说，change命令不能一定节点的位置。
- replace：
 
 	对一个现有节点进行近于原子操作的删除／添加。如果节点不存在，这个命令就会建立节点。
- link：
 
 	只适用于DQisc，替代一个现有的节点

### example
ref:
- [TC(Traffic Control)命令—linux自带高级流控](https://cloud.tencent.com/developer/article/1409664)

```bash
tc qdisc ls # 查看所有interface上的tc配置
tc qdisc show dev eth0 # 查看eth0上的tc配置
tc qdisc add dev enp0s3 root netem delay 1000ms # 设置数据包延迟1000ms(非精确, 但不会少于1000ms)发送
tc qdisc add dev enp0s3 root netem delay 1000ms 500ms # # 设置数据包延迟1000ms 正负 500ms.
tc qdisc change dev enp0s3 root netem delay 1000ms # 修改
tc qdisc add dev enp0s3 root netem loss 50% # 丢包50%
tc qdisc add dev enp0s3 root netem corrupt 50% # Packet Corruption 包损50%
tc qdisc add dev enp0s3 root netem duplicate 50% # Packet Duplicates 重复包50%
tc qdisc add dev eth1 root tbf rate 5120kbit latency 50ms burst 5120 # 限速. tbf, 算法; rate, 5120kbit表示限制宽带为5M; burst, 峰值为5M
tc qdisc del dev eth1 root tbf
tc qdisc del dev eth0 root # 删除eth0上的 egress "root qdisc". grep ' htb '
tc qdisc del dev eth0 ingress # 删除eth0上的ingress. grep ' ingress '
tc qdisc add dev eth0 root handle 1: htb # 建立root qdisc

# --- 针对端口进行限速, 场景: 在使用git拉去代码时很容易跑满带宽，为了控制带宽的使用
tc -s qdisc ls dev eth0 # 查看现有的队列
tc -s class ls dev eth0 # 查看现有的分类
tc qdisc add dev eth0 root handle 1:0 htb default 1 # 添加一个tbf队列，绑定到eth0上，命名为1:0 ，默认归类为1. handle：为队列命名或指定某队列
tc class add dev eth0 parent 1:0 classid 1:1 htb rate 10Mbit burst 15k # 为eth0下的root队列1:0添加一个分类并命名为1:1，类型为htb，带宽为10M. rate: 是一个类保证得到的带宽值.如果有不只一个类,请保证所有子类总和是小于或等于父类
tc class add dev eth0 parent 1:1 classid 1:10 htb rate 10Mbit ceil 10Mbit burst 15k # 为1:1类规则添加一个名为1:10的类，类型为htb，带宽为10M
tc qdisc add dev eth0 parent 1:10 handle 10: sfq perturb 10 # 为了避免一个会话永占带宽,添加随即公平队列sfq. perturb：是多少秒后重新配置一次散列算法，默认为10秒; sfq,他可以防止一个段内的一个ip占用整个带宽
tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip sport 22 flowid 1:10 # 使用u32创建过滤器

# --- 针对ip进行限速 情景：因为带宽资源有限（20Mbit≈2Mbyte），使用git拉取代码的时候导致带宽资源告警，所以对git进行限速，要求：内网不限速；外网下载速度为1M左右
#清空原有规则
tc qdisc del dev eth0 root

#创建根序列. 默认使用1的class
tc qdisc add dev eth0 root handle 1: htb default 1

#创建一个主分类绑定所有带宽资源（20M）
tc class add dev eth0 parent 1:0 classid 1:1 htb rate 20Mbit burst 15k

#创建子分类. ceil是一个类最大能得到的带宽值. `classid1:10`表示创建一个标识为1:11的类别
tc class add dev eth0 parent 1:1 classid 1:10 htb rate 20Mbit ceil 10Mbit burst 15k
tc class add dev eth0 parent 1:1 classid 1:20 htb rate 20Mbit ceil 20Mbit burst 15k

#避免一个ip霸占带宽资源
tc qdisc add dev eth0 parent 1:10 handle 10: sfq perturb 10
tc qdisc add dev eth0 parent 1:20 handle 20: sfq perturb 10

#创建过滤器. 在过滤器中加入prio，指定规则的优先级(prio 越小，优先级越高), 以避免内网ip被限速
#对所有ip限速
tc filter add dev eth0 protocol ip parent 1:0 prio 2 u32 match ip dst 0.0.0.0/0 flowid 1:10
#对内网ip放行
tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip dst 12.0.0.0/8 flowid 1:20
```