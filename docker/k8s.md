# k8s
env: k8s 1.14.1 / Rancher v2.2.3

参考:
- [Kubernetes核心概念总结](http://dockone.io/article/8866)
- [Kubernetes架构为什么是这样的？](https://www.tuicool.com/articles/J7Rbimu)
- [k8s yaml apiVersion](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.14/)
- [kubernetes/examples](https://github.com/kubernetes/examples)

Kubernetes 最主要的设计思想是从更宏观的角度，以统一的方式来定义任务之间的各种关系，并且为将来支持更多种类的关系留有余地.

## 概念
### 基于k8s的开发模式
而当应用本身发生变化时，开发人员和运维人员可以依靠容器镜像来进行同步；当应用部署参数发生变化时，YAML 配置文件就是他们相互沟通和信任的媒介.

Kubernetes“一切皆对象”的设计思想.

### 容器
是一种沙盒技术. 容器技术的核心功能，就是通过约束和修改进程的动态表现，从而为其创造出一个“边界”. “敏捷”和“高性能”是容器相较于虚拟机最大的优势.

对于 Docker 等大多数 Linux 容器来说，Cgroups 技术是用来制造**约束**的主要手段，而Namespace 技术则是用来修改**进程视图(容器进程看待整个操作系统的视图)**的主要方法

一个正在运行的 Linux 容器，其实可以被“一分为二”地看待：
- 一组是“容器镜像”（Container Image），是容器的静态视图
- 一个由 Namespace+Cgroups 构成的隔离环境, 是“容器运行时”（Container Runtime），是容器的动态视图

> Docker 容器的本质是Namespace 做隔离，Cgroups 做限制，镜像(rootfs) 做文件系统
> Cgroups是进程设置资源限制的方法, Namespace是进程隔离的方法.
> 容器是一个“单进程”模型
> 容器本身的设计，就是希望容器和应用能够同生命周期，这个概念对后续的容器编排非常重要. 否则，一旦出现类似于“容器是正常运行的，但是里面的应用早已经挂了”的情况，编排系统处理起来就非常麻烦了
> 容器所谓的一致性：无论在本地、云端，还是在一台任何地方的机器上，用户只需要解压打包好的容器镜像，那么这个应用运行所需要的完整的执行环境就被重现出来了
> 对一个应用来说，操作系统本身才是它运行所需要的最完整的“依赖库”

#### 容器的持久化存储
容器最典型的特征之一是无状态, 而容器的持久化存储，就是用来保存容器存储状态的重要手段. 存储插件会在容器里挂载一个基于网络或者其他机制的**远程数据卷**，使得在容器里创建的文件，实际上是保存在远程存储服务器上，或者以分布式的方式保存在多个节点上，而与当前宿主机没有任何绑定关系.

Kubernetes 项目上创建的所有 Pod 就能够通过 Persistent Volume（PV）和 Persistent Volume Claim（PVC）的方式在容器里挂载持久化的数据卷了.

### 控制器模式
在 Kubernetes 中使用一种 API 对象（Deployment）管理另一种 API 对象（Pod）的方法

### 编排（Orchestration）
指用户如何通过某些工具或者配置来完成一组虚拟机/容器以及关联资源的定义、配置、创建、删除等工作，然后由云计算平台按照这些指定的逻辑来完成的过程.

### cluster
cluster是计算,存储和网络资源的集合. k8s使用这些资源运行各种基于容器的服务.

### master
master是cluster的大脑, 主要职责是调度. 为了高可用, 通常会运行多个master

### controller
k8s通常不直接创建pod, 而是通过controller来运行管理pod.

controller分类:
- Deployment : 管理pod的多副本(通过ReplicaSet), 确保pod按照期望的状态运行
- ReplicaSet : 实现了Pod的多副本管理,包括滚动更新. 使用Deployment时会自动创建ReplicaSet, 因此我们通常不直接使用它
- DaemonSet : 用于node最多只运行一个pod副本的场景
- StatefuleSet : 保证pod的每个副本在整个生命周期中的名称是不变的(因故障需删除并重启除外), 同时会保证副本按照固定的顺序启动,更新或删除
- Job

> 一个 ReplicaSet 对象，其实就是由副本数目的定义和一个 Pod 模板组成的; Deployment 控制器实际操纵的，正是 ReplicaSet 对象，而不是 Pod 对象

#### deployment
部署pod并分布到各个node上, 每个node允许存在多个副本.

`kubectl get deployment ${deployment_name}`输出:
- DESIRED : 期望状态是READY的副本数(spec.replicas 的值)
- CURRENT : 当前处于Running状态的副本总数
- UP-TO-DATE : 当前已完成更新的副本数(即当前处于最新版本的 Pod 的个数)
- AVAILABLE : 当前处于READY的副本数(Running状态+最新版本+处于Ready（健康检查正确）状态), 是用户所期望的最终状态

滚动更新的控制参数:
- maxSurge : 控制滚动更新中副本总数超过DESIRED的上限, 默认是`25%`, 取值方法是roundUp. 值越大, 滚动更新中创建的新副本就越多
- maxUnavailable : 控制滚动更新中, 不可用副本占DESIRED的上限, 默认是`25%`, 取值方法是roundDown. 值越大, 滚动更新中销毁的旧副本就越多

![Deployment、ReplicaSet 和 Pod 的关系](https://static001.geekbang.org/resource/image/79/f6/79dcd2743645e39c96fafa6deae9d6f6.png)

#### daemonset
通过 nodeAffinity 和 Toleration保证每个node至多运行一个副本. k8s本身就在用daemonset运行部分组件`kubectl get daemonsets --all-namespaces`.

daemonset的典型应用场景:
1. 存储: ceph
1. 日志收集: flunentd, logstash
1. 监控: prometheus node exporter, collected

在创建每个 Pod 的时候，DaemonSet 会自动给这个 Pod 加上一个 nodeAffinity，从而保证这个 Pod 只会在指定节点上启动. 同时它还会自动给这个 Pod 加上一个 Toleration，从而忽略节点的 unschedulable"污点".

> DaemonSet支持滚动更新
> 在 DaemonSet 上，我们一般都应该加上 resources 字段，来限制它的 CPU 和内存使用，防止它占用过多的宿主机资源

#### job
从程序的运行形态上来区分，我们可以将Pod分为两类：
- 长时运行服务（JBoss、MySQL等）
- 一次性任务（数据计算、测试）

RC创建的Pod都是长时运行的服务，而Job创建的Pod都是一次性任务.

在Job的定义中，restartPolicy（重启策略）只能是Never和OnFailure. restartPolicy=Never，那么离线作业失败后 Job Controller 就会不断地尝试创建一个新 Pod;restartPolicy=OnFailure，那么离线作业失败后，Job Controller 就不会去尝试创建新的 Pod, 而是会不断地尝试重启 Pod 里的容器

可用`kubectl get pods --show-all`查看已`Completed`的pod
job的并行性可通过设置parallelism.
job设置completion可指定完成job所需要的pod总数,即job执行次数.
`spec.backoffLimit`设置job的重试次数; Job Controller 重新创建 Pod 的间隔是呈指数增加的，即下一次重新创建 Pod 的动作会分别发生在 10 s、20 s、40 s …后
`spec.activeDeadlineSeconds` 字段可以设置最长运行时间
在 Job 对象中，负责并行控制的参数有两个：
- spec.parallelism，它定义的是一个 Job 在任意时间最多可以启动多少个 Pod 同时运行；
- :spec.completions，它定义的是 Job 完成时 至少要完成的 Pod 数目，即 Job 的最小完成数

Job Controller 控制的对象，直接就是 Pod. Job Controller 在控制循环中进行的调谐（Reconcile）操作，是根据实际在 Running 状态 Pod 的数目、已经成功退出的 Pod 的数目，以及 parallelism、completions 参数的值共同计算出在这个周期里，应该创建或者删除的 Pod 数目，然后调用 Kubernetes API 来执行这个操作.

使用 Job 对象的方法: 外部管理器 +Job 模板, 最简单普遍.

###cronjob
job还可以是定时Job, yaml的kind = CronJob. CronJob 是一个 Job 对象的控制器（Controller）
schedule 字段定义的、一个标准的Unix Cron格式的表达式。

比如，"*/1 * * * *"。

这个 Cron 表达式里 */1 中的 * 表示从 0 开始，/ 表示“每”，1 表示偏移量。所以，它的意思就是：从 0 开始，每 1 个时间单位执行一次.
Cron 表达式中的五个部分分别代表：分钟、小时、日、月、星期。

所以，上面这句 Cron 表达式的意思是：从当前开始，每分钟执行一次.

通过 spec.concurrencyPolicy 字段来定义具体的处理策略:
- concurrencyPolicy=Allow，这也是默认情况，这意味着这些 Job 可以同时存在
- concurrencyPolicy=Forbid，这意味着不会创建新的 Pod，该创建周期被跳过
- concurrencyPolicy=Replace，这意味着新产生的 Job 会替换旧的、没有执行完的 Job

如果某一次 Job 创建失败，这次创建就会被标记为“miss”。当在指定的时间窗口(spec.startingDeadlineSeconds)内，miss 的数目达到 100 时，那么 CronJob 会停止再创建这个 Job

### node
node是具体负责运行容器的应用, 会监控并汇报容器状态, 同时会根据master的要求管理容器的生命周期. node由master管理.

### pod
pod是k8s最小的工作(调度,扩展,共享资源,管理生命周期)单位. 它是一个或多个相关容器的集合, 其原因是:
1. 可管理性: 某些容器的应用间存在紧密联系
1. 共享通信和存储: pod中的容器使用同一个网络namespace.

围绕着容器和 Pod 不断向真实的技术场景扩展，我们就能够摸索出一幅如下所示的 Kubernetes 项目核心功能的“全景图”:
![Kubernetes 项目核心功能的“全景图](https://static001.geekbang.org/resource/image/16/06/16c095d6efb8d8c226ad9b098689f306.png)

> pod类似于进程组的或虚拟机的角色, 而容器就是里面的进程
> 凡是调度、网络、存储，以及安全相关的属性，基本上是 Pod 级别的; 凡是跟容器的 Linux Namespace 相关的属性或容器要共享宿主机的 Namespace，都一定是 Pod 级别的, 因为Pod 的设计就是要让它里面的容器尽可能多地共享 Linux Namespace，仅保留必要的隔离和限制能力
> Pod是一组共享了某些资源的容器, Pod 里的所有容器，共享的是同一个 Network Namespace，并且可以声明共享同一个 Volume
> 容器进程返回值非零, k8s会认为容器发生故障就会按照Pod的restartPolicy进行处理

restartPolicy:
- Always ：在任何情况下，只要容器不在运行状态，就自动重启容器
- OnFailure : 只在容器 异常时才自动重启容器
- Never : 从来不重启容器

> Pod 的恢复过程，永远都是发生在当前节点上，而不会跑到别的节点上去. 事实上，一旦一个 Pod 与一个节点（Node）绑定，除非这个绑定发生了变化（pod.spec.node 字段被修改），否则它永远都不会离开这个节点包括(宿主机宕机). 而如果想让 Pod 出现在其他的可用节点上，就必须使用 Deployment Controller 来管理 Pod.

restartPolicy 和 Pod 里容器的状态，以及 Pod 状态的对应关系:
- 只要 Pod 的 restartPolicy 指定的策略允许重启异常的容器（比如：Always），那么这个 Pod 就会保持 Running 状态，并进行容器重启。否则，Pod 就会进入 Failed 状态 .
- 对于包含多个容器的 Pod，只有它里面所有的容器都进入异常状态后，Pod 才会进入 Failed 状态. 在此之前，Pod 都是 Running 状态. Pod 的 READY 字段会显示正常容器的个数.

在 Kubernetes 项目里，Pod 的实现需要使用一个中间容器，这个容器叫作 Infra 容器. 在这个 Pod 中，Infra 容器永远都是第一个被创建的容器，而其他用户定义的容器，则通过 Join Network Namespace 的方式，与 Infra 容器关联在一起.

> Infra 容器一定要占用极少的资源，所以它使用的是一个非常特殊的镜像，叫作：k8s.gcr.io/pause. 这个镜像是一个用汇编语言编写的、永远处于“暂停”状态的容器，解压后的大小也只有 100~200 KB 左右.
> 在 Pod 中，所有 Init Container 定义的容器，都会比 spec.containers 定义的用户容器先启动. 之后Init Container 容器会按顺序逐一启动，而直到它们都启动并且退出了，用户容器才会启动
> sidecar 模式指在一个 Pod 中启动一个辅助容器，来完成一些独立于主进程（主容器）之外的工作. 最典型的例子是 Istio 这个微服务治理项目

基本操作:
- 创建: kubectl create -f xxx.yaml
- 查询: kubectl get pod ${pod_name}, kubectl describe pod ${pod_name}
- 删除: kubectl delete pod ${pod_name}
- 更新: kubectl replace xxx_new.yaml

yaml设置:
- 启动命令: spec.containers.command
- 环境变量: spec.containers.env.name/value
- 端口桥接: spec.containers.ports.containerPort/protocol/hostIP/hostPort
  使用hostPort时需要注意端口冲突的问题，不过Kubernetes在调度Pod的时候会检查宿主机端口是否冲突，比如当两个Pod均要求绑定宿主机的80端口，Kubernetes将会将这两个Pod分别调度到不同的机器上
- Host网络: spec.hostNetwork=true
  一些特殊场景下，容器必须要以host方式进行网络设置(如接收物理机网络才能够接收到的组播流)
- 数据持久化: spec.containers.volumeMounts.mountPath
- 重启策略: 当Pod中的容器终止退出后，重启容器的策略. 这里的所谓Pod的重启，实际上的做法是容器的重建，之前容器中的数据将会丢失，如果需要持久化数据，那么需要使用数据卷进行持久化设置.
  Pod支持三种重启策略:
    - Always(默认策略，当容器终止退出后，总是重启容器)
    - OnFailure(当容器终止且异常退出时，重启)
    - Never(从不重启)

> Pod共享PID、Network、IPC、UTS. 除此之外，Pod中的容器可以访问共同的数据卷来实现文件系统的共享

默认情况下 Master 节点是不允许运行用户 Pod 的. 而 Kubernetes 做到这一点依靠的是 Kubernetes 的 Taint/Toleration 机制. 它的原理是一旦某个节点被加上了一个 Taint即被“打上了污点”，那么所有 Pod 就都不能在这个节点上运行.

为节点打上“污点”（Taint）的命令是：
```sh
$ kubectl taint nodes ${node_name} foo=bar:NoSchedule // NoSchedule意味着这个 Taint 只会在调度新 Pod 时产生作用，而不会影响已经在 node 上运行的 Pod
```

Toleration : pod忽略Taint, 只要在 Pod 的.yaml 文件中的 spec 部分，加入 tolerations 字段即可：
```yaml
apiVersion: v1
kind: Pod
...
spec:
  tolerations:
  - key: "foo" # 这个 Pod 能“容忍”所有键值对为 foo=bar 的 Taint
    operator: "Equal" # 还可以使用`operator: "Exists"`
    value: "bar"
    effect: "NoSchedule"
```

移除taint:
```
$ kubectl taint nodes --all foo- // 只需要在taint的键后面加上了一个短横线`-`即可
```

#### yaml Container 属性
ImagePullPolicy定义了镜像拉取的策略

ImagePullPolicy
- Always(默认) : 每次创建 Pod 都重新拉取一次镜像
- Never : Pod 永远不会主动拉取这个镜像
- IfNotPresent : 只在宿主机上不存在这个镜像时才拉取

Lifecycle定义了Container Lifecycle Hooks:
- postStart : 在容器启动后，立刻执行一个指定的操作. 虽然是在 Docker 容器 ENTRYPOINT 执行之后，但它并不严格保证顺序. 即在 postStart 启动时，ENTRYPOINT 有可能还没有结束.
  如果 postStart 执行超时或者错误，Kubernetes 会在该 Pod 的 Events 中报出该容器启动失败的错误信息，导致 Pod 也处于失败的状态
- preStop 发生的时机则是容器被杀死之前（比如，收到了 SIGKILL 信号）. 而需要明确的是，preStop 操作的执行是同步的. 所以，它会阻塞当前的容器杀死流程，直到这个 Hook 定义操作完成之后，才允许容器被杀死，这跟 postStart 不一样

#### Pod 的生命周期

Pod 生命周期的变化主要体现在 Pod API 对象的Status 部分，这是它除了 Metadata 和 Spec 之外的第三个重要字段. 其中pod.status.phase就是 Pod 的当前状态，它有如下几种可能的情况：
- Pending : Pod 的 YAML 文件已经提交给了 Kubernetes，API 对象已经被创建并保存在 Etcd 当中. 但是这个 Pod 里有些容器因为某种原因而不能被顺利创建, 比如，调度不成功.
- Running : Pod 已经调度成功，跟一个具体的节点绑定. 它包含的容器都已经创建成功，并且至少有一个正在运行中.
- Succeeded : Pod 里的所有容器都正常运行完毕，并且已经退出了. 这种情况在运行Job任务时最为常见.
- Failed : Pod 里至少有一个容器以不正常的状态（非 0 的返回码）退出. 这个状态的出现意味着得想办法 Debug 这个容器的应用，比如查看 Pod 的 Events 和日志
- Unknown : 异常状态，意味着 Pod 的状态不能持续地被 kubelet 汇报给 kube-apiserver，这很有可能是主从节点（Master 和 Kubelet）间的通信出现了问题

更进一步地，Pod 对象的 Status 字段还可以再细分出一组 Conditions. 这些细分状态的值包括：PodScheduled、Ready、Initialized，以及 Unschedulable. 它们主要用于描述造成当前 Status 的具体原因是什么.
比如，Pod 当前的 Status 是 Pending，对应的 Condition 是 Unschedulable，这就意味着它的调度出现了问题.
而其中，Ready 这个细分状态非常值得我们关注：它意味着 Pod 不仅已经正常启动（Running 状态），而且已经可以对外提供服务了. 这两者之间（Running 和 Ready）是有区别的.

> [type Pod struct](https://github.com/kubernetes/api/blob/master/core/v1/types.go)

#### health check
kubelet 就会根据指定 Probe 的返回值决定这个容器的状态，而不是直接以容器进行是否运行（来自 Docker 返回的信息）作为依据.

livenessProbe支持exec, HTTP, TCP.

- livenessProbe : 告诉k8s何时通过重启容器实现自愈
- readinessProbe :　告诉k8s何时可以将容器加入Service的负载均衡池对外提供服务, 常用于Scale Up/Rolling Update中.

两者比较:
1. 默认均通过判断容器进程的返回值是否为零来判断探测是否成功; 默认连续3次非零则启用应对策略
1. liveness失败是重启容器, readiness失败是将容器设为不可用, 不再接收Service转发的请求
1. 两者独立无依赖, 可组合使用

#### Pod生命周期
![pod生命周期](http://dockone.io/uploads/article/20190520/c8e551e53f7e7e2a3af022c4ea672fe9.png)

Pod被分配到一个Node上之后，就不会离开这个Node，直到被删除. 当某个Pod失败，首先会被Kubernetes清理掉，之后ReplicationController将会在其它机器上（或本机）重建Pod，**重建之后Pod的ID发生了变化**，那将会是一个新的Pod. 所以**Kubernetes中Pod的迁移，实际指的是在新Node上重建Pod**

生命周期回调函数：PostStart（容器创建成功后调用该回调函数）、PreStop（在容器被终止前调用该回调函数）
```yaml
containers:
- image: sample:v2  
 name: war
 lifecycle：
  posrStart:
   exec:
     command:
      - “cp”
      - “/sample.war”
      - “/app”
  prestop:
   httpGet:
    host: monitor.com
    psth: /waring
    port: 8080
    scheme: HTTP
```

PodPreset（Pod 预设置 form v1.11): PodPreset 里定义的内容，只会在 Pod 被创建之前追加在这个对象本身上, 不推荐, 使用yaml模板更直观.

### service
定义了外界访问一组特定pod的方法. 它提供了一套简化的服务代理和发现机制，天然适应微服务架构.

在Kubernetes中，在受到RC调控的时候，Pod副本是变化的，对应的虚拟IP也是变化的，比如发生迁移或者伸缩的时候。这对于Pod的访问者来说是不可接受的. Kubernetes中的Service是一种抽象概念，它定义了一个Pod逻辑集合以及访问它们的策略，**Service同Pod的关联同样是居于Label来完成的**. Service的目标是提供一种桥梁， 它会为访问者提供一个固定访问地址，用于在访问时重定向到相应的后端，这使得非 Kubernetes原生应用程序，在无须为Kubemces编写特定代码的前提下，轻松访问后端.

Service同RC一样，都是通过Label来关联Pod的. 当Pod发生变化时（增加、减少、重建等），Service会及时更新. 这样一来，Service就可以作为Pod的访问入口，起到代理服务器的作用，而对于访问者来说，通过Service进行访问，无需直接感知Pod.

需要注意的是，Kubernetes分配给Service的固定IP是一个虚拟IP，并不是一个真实的IP，**在外部是无法寻址的**. 真实的系统实现上，Kubernetes是通过Kube-proxy组件来实现的虚拟IP路由及转发. 所以在之前集群部署的环节上，我们在每个Node上均部署了Proxy这个组件，从而实现了Kubernetes层级的虚拟转发网络.

#### Service代理外部服务
Service不仅可以代理Pod，还可以代理任意其他后端，比如运行在Kubernetes外部的MySQL、Oracle等. 这是通过定义两个同名的Service和Endpoints来实现的：

```yaml
redis-service.yaml
apiVersion: v1
kind: Service
metadata:
name: redis-service
spec:
ports:
- port: 6379
targetPort: 6379
protocol: TCP
```

```yaml
redis-endpoints.yaml
apiVersion: v1
kind: Endpoints
metadata:
name: redis-service
subsets:
- addresses:
- ip: 10.0.251.145
ports:
- port: 6379
  protocol: TCP
```

基于文件创建完Service和Endpoints之后，在Kubernetes的Service中即可查询到自定义的Endpoints:
```bash
[root@k8s-master demon]# kubectl describe service redis-service
Name:            redis-service
Namespace:        default
Labels:            <none>
Selector:        <none>
Type:            ClusterIP
IP:            10.254.52.88
Port:            <unset>    6379/TCP
Endpoints:        10.0.251.145:6379
Session Affinity:    None
No events. // 进行 Debug 的重要依据. 如果有异常发生，events往往可以看到非常详细的错误信息
[root@k8s-master demon]# etcdctl get /skydns/sky/default/redis-service
{"host":"10.254.52.88","priority":10,"weight":10,"ttl":30,"targetstrip":0} 
```

### Service内部负载均衡
当Service的Endpoints包含多个IP的时候，及服务代理存在多个后端，将进行请求的负载均衡, 默认的负载均衡策略是轮训或者随机（由kube-proxy的模式决定）. 同时，Service上通过设置Service.spec.sessionAffinity=ClientIP，来实现基于源IP地址.
会话保持.

### 发布Service
Service的虚拟IP是由Kubernetes虚拟出来的内部网络，外部是无法寻址到的。但是有些服务又需要被外部访问到，例如Web前端. 这时候就需要加一层网络转发，即外网到内网的转发. Kubernetes提供了NodePort、LoadBalancer、Ingress三种方式.

NodePort的原理是Kubernetes会在每一个Node上暴露出一个端口：nodePort，外部网络可以通过（任一Node）[NodeIP]:[NodePort]访问到后端的Service.

LoadBalancer，在NodePort基础上，Kubernetes可以请求底层云平台创建一个负载均衡器，将每个Node作为后端，进行服务分发. 该模式需要底层云平台（例如GCE）支持.

Ingress，是一种HTTP方式的路由转发机制，由Ingress Controller和HTTP代理服务器组合而成. Ingress Controller实时监控Kubernetes API，实时更新HTTP代理服务器的转发规则. HTTP代理服务器有GCE Load-Balancer、HaProxy、Nginx等开源方案.

### Servicede自发性机制
Kubernetes中有一个很重要的服务自发现特性. 一旦一个Service被创建，该Service的Service IP和Service Port等信息都可以被注入到Pod中供它们使用. Kubernetes主要支持两种Service发现机制：环境变量和DNS.

1. 环境变量

Kubernetes创建Pod时会自动添加所有可用的Service环境变量到该Pod中，如有需要．这些环境变量就被注入Pod内的容器里. 需要注意的是，**环境变量的注入只发送在Pod创建时，且不会被自动更新**. 这个特点暗含了Service和访问该Service的Pod的创建时间的先后顺序，即任何想要访问Service的Pod都需要在Service已经存在后创建，否则与Service相关的环境变量就无法注入该Pod的容器中，这样先创建的容器就无法发现后创建的Service. 

1. DNS

Kubernetes集群现在支持增加一个可选的组件——DNS服务器. 这个DNS服务器使用Kubernetes的watchAPI，不间断的监测新的Service的创建并为每个Service新建一个DNS记录. 如果DNS在整个集群范围内都可用，那么所有的Pod都能够自动解析Service的域名. Kube-DNS搭建及更详细的介绍请见[基于Kubernetes集群部署skyDNS服务](http://www.cnblogs.com/zhenyuyaodidiao/p/6500992.html). 

### 多个Service如何避免地址和端口冲突
此处设计思想是，Kubernetes通过为每个Service分配一个唯一的ClusterIP，所以当使用ClusterIP：port的组合访问一个Service的时候，不管Port是什么，这个组合是一定不会发生重复的. 另一方面，kube-proxy为每个Service真正打开的是一个绝对不会重复的随机端口，用户在Service描述文件中指定的访问端口会被映射到这个随机端口上. 这就是为什么用户可以在创建service时随意指定访问端口. 

### Service目前存在的不足
Kubernetes使用iptables和kube-proxy解析Service的入口地址，在中小规模的集群中运行良好，但是当Service的数量超过一定规模时，仍然有一些小问题. 首当其冲的便是Service环境变量泛滥，以及Service与使用Service的Pod两者创建时间先后的制约关系. 目前来看，很多使用者在使用Kubernetes时往往会开发一套自己的Router组件来替代Service，以便更好地掌控和定制这部分功能. 

### Deployment
Kubernetes提供的一种更加简单的维护RC和Pod的机制. 通过在Deployment中描述所期望的集群状态，Deployment Controller会将现在的集群状态在一个可控的速度下逐步更新到所期望的集群状态.

90%的功能与Replication Controller完全一样，可以看做新一代的Replication Controller. 但是，它又具备了Replication Controller之外的新特性：
- Replication Controller全部功能
- 事件和状态查看：可以查看Deployment的升级详细进度和状态
- 回滚：当升级Pod镜像或者相关参数的时候发现问题，可以使用回滚操作回滚到上一个稳定的版本或者指定的版本
- 版本记录：每一次对Deployment的操作，都能保存下来，给予后续可能的回滚使用
- 暂停和启动：对于每一次升级，都能够随时暂停和启动
- 多种升级方案：
    - Recreate : 删除所有已存在的Pod，重新创建新的
    - RollingUpdate : 滚动升级，逐步替换的策略，同时滚动升级时，支持更多的附加参数，例如设置最大不可用Pod数量，最小升级间隔时间等等

> 前提假设: 一个应用的所有 Pod，是完全一样的, 它们互相之间没有顺序，也无所谓运行在哪台宿主机上. 需要的时候，Deployment 就可以通过 Pod 模板创建新的 Pod；不需要的时候，Deployment 就可以“杀掉”任意一个 Pod
> Deployment 实际上是一个两层控制器. 首先，它通过ReplicaSet 的个数来描述应用的版本；然后，它再通过ReplicaSet 的属性（比如 replicas 的值），来保证 Pod 的副本数量

#### 滚动升级
相比于RC，Deployment直接使用`kubectl edit deployment ${deployment_name}(推荐)`或者`kubectl set`方法就可以直接升级（原理是Pod的template发生变化，例如更新Label、更新镜像版本等操作会触发Deployment的滚动升级）.

创建Deployment:
```sh
$ kubectl create -f nginx-deploy-v1.yaml --record // `--record`便于之后使用查看history
$ kubectl rollout history deployment nginx-deployment // 查看Deployment的变更信息, 这些信息得以保存是创建时候加的`--record`选项的作用
```
Deployment的一些基础命令:
```sh
$ kubectl describe deployments  #查询详细信息，获取升级进度
$ kubectl rollout pause deployment nginx-deployment2  #暂停升级, 此时对 Deployment 的所有修改都不会触发新的"滚动更新", 也不会创建新的 ReplicaSet
$ kubectl rollout resume deployment nginx-deployment2  #继续升级, kubectl rollout resume 之前，在 kubectl rollout pause 之后的这段时间里对 Deployment 进行的所有修改，最后只会触发一次"滚动更新"
$ kubectl rollout undo deployment nginx-deployment2  #升级回滚
$ kubectl scale deployment nginx-deployment --replicas 10  #弹性伸缩Pod数量
kubectl rollout status deployment nginx-deployment # 实时查看 Deployment 对象的状态变化
```

> 对 Deployment 进行的每一次更新操作，都会生成一个新的 ReplicaSet 对象, 可通过`spec.revisionHistoryLimit`解决
> 滚动更新要求你一定要使用 Pod 的 Health Check 机制检查应用的运行状态，而不是简单地依赖于容器的 Running 状态. 否则容器已经变成 Running 了，但服务很有可能尚未启动，“滚动更新”的效果也就达不到了
> 将一个集群中正在运行的多个 Pod 版本，交替地逐一升级的过程，就是“滚动更新”
> `kubectl get rs`的ReplicaSet 的 DESIRED、CURRENT 和 READY 字段的含义，和 Deployment 中是一致的. Deployment 只是在 ReplicaSet 的基础上，添加了 UP-TO-DATE 这个跟版本有关的状态字段.

> k8s使用ControllerRevision(from v1.7)来记录Controller 的版本, 其在回滚时非常有用.

#### Horizontal Pod Autoscaler
Horizontal Pod Autoscaler的操作对象是Replication Controller、ReplicaSet或Deployment对应的Pod，根据观察到的实际资源使用量与用户的期望值进行比对，做出是否需要增减实例数量的决策.

### volume
参考:
- [k8s学习笔记之持久化存储](https://zhuanlan.zhihu.com/p/29706309)

在Kubernetes中，当Pod重建的时候，数据是会丢失的，Kubernetes也是通过数据卷挂载来提供Pod数据的持久化的. Kubernetes数据卷是对Docker数据卷的扩展，Kubernetes数据卷是Pod级别的，可以用来实现Pod中容器的文件共享. 目前，[Kubernetes支持的数据卷类型](https://kubernetes.io/docs/concepts/storage/volumes/#types-of-volumes)如下：
- EmptyDir
- HostPath
- GCE Persistent Disk
- AWS Elastic Block Store
- NFS
- iSCSI
- Flocker
- GlusterFS
- RBD
- Git Repo
- Secret
- Persistent Volume Claim
- Downward API

Kubernetes中提供了存储消费模式: Persistent Volume和Persistent Volume Claim机制.
[Persistent Volume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#types-of-persistent-volumes)是由系统管理员配置创建的一个数据卷（目前支持HostPath、GCE Persistent Disk、AWS Elastic Block Store、NFS、iSCSI、GlusterFS、RBD），它代表了某一类存储插件实现；而对于普通用户来说，通过Persistent Volume Claim可请求并获得合适的Persistent Volume，而无须感知后端的存储实现. Persistent Volume和Persistent Volume Claim相互关联，有着完整的生命周期管理：
- 准备：系统管理员规划或创建一批Persistent Volume；
- 绑定：用户通过创建Persistent Volume Claim来声明存储请求，Kubernetes发现有存储请求的时候，就去查找符合条件的Persistent Volume（最小满足策略）。找到合适的就绑定上，找不到就一直处于等待状态；
- 使用：创建Pod的时候使用Persistent Volume Claim；
- 释放：当用户删除绑定在Persistent Volume上的Persistent Volume Claim时，Persistent Volume进入释放状态，此时Persistent Volume中还残留着上一个Persistent Volume Claim的数据，状态还不可用；
- 回收：是否的Persistent Volume需要回收才能再次使用。回收策略可以是人工的也可以是Kubernetes自动进行清理（仅支持NFS和HostPath）

PersistentVolume的Access Modes:
- ReadWriteOnce : the volume can be mounted as read-write by a single node (单node的读写) 
- ReadOnlyMany : the volume can be mounted read-only by many nodes (多node的只读) 
- ReadWriteMany : the volume can be mounted as read-write by many nodes (多node的读写) 

pv的reclaim policy:
- Retain : 管理员手动回收
- Recycle : 基本的数据擦除 (`rm -rf /thevolume/*`)
- Delete : 删除相关联的后端存储卷， 后端存储比如AWS EBS, GCE PD, Azure Disk, or OpenStack Cinder

pv的volumn phase:
- Available : 资源可用， 还没有被声明绑定
- Bound : 被声明绑定
- Released : 绑定的声明被删除了，但是还没有被集群重声明
- Failed : 自动回收失败

注意: 只有本地盘和nfs支持数据盘Recycle 擦除回收， AWS EBS, GCE PD, Azure Disk, and Cinder 存储卷支持Delete策略

PVC 要真正被容器使用起来，就必须先和某个符合条件的 PV 进行绑定, 要检查的条件包括两部分：
1. PV 和 PVC 的 spec 字段. 比如PV 的存储（storage）大小，就必须满足 PVC 的要求
1. PV 和 PVC 的 storageClassName 字段必须一样

PVC 可以理解为持久化存储的“接口”，它提供了对某种持久化存储的描述，但不提供具体的实现；而这个持久化存储的实现部分则由 PV 负责完成.
PersistentVolumeController将一个 PV 与 PVC 进行“绑定”，其实就是将这个 PV 对象的名字，填在了 PVC 对象的 spec.volumeName 字段上

PV 对象如何变成容器里的一个持久化存储(两阶段处理):
1. Attach : 为虚拟机挂载远程磁盘的操作, 有的不需要挂载设备可以跳过, 比如NFS.
  该阶段是由 Volume Controller 负责维护的，这个控制循环的名字叫作：AttachDetachController, 其运行在 Master 节点上的: Attach 操作只需要调用公有云或者具体存储项目的 API，并不需要在具体的宿主机上执行操作
1. Mount: 将磁盘设备格式化并挂载到 Volume 宿主机目录的操作
  必须发生在 Pod 对应的宿主机上，所以它必须是 kubelet 组件的一部分。这个控制循环的名字，叫作：VolumeManagerReconciler，它运行起来之后，是一个独立于 kubelet 主循环的 Goroutine

上述关于 PV 的“两阶段处理”流程，是靠独立于 kubelet 主控制循环（Kubelet Sync Loop）之外的两个控制循环来实现的

pv 存在:
- 静态供给(Static Provisioning) : 集群管理员创建多个PV, 存在于Kubernetes API中，可用于消费. 它们携带可供集群用户使用的真实存储的详细信息
- 动态供给(Dynamical Provisioning) :  StorageClass 对象的作用，其实就是创建 PV 的模板, 只有同属于一个 StorageClass 的 PV 和 PVC，才可以绑定在一起

具体地说，StorageClass 对象会定义如下两个部分内容：
1. PV 的属性 : 比如，存储类型、Volume 的大小等等
1. 创建这种 PV 需要用到的存储插件 : 比如，Ceph 等等
有了这样两个信息之后，Kubernetes 就能够根据用户提交的 PVC，找到一个对应的 StorageClass 了。然后，Kubernetes 就会调用该 StorageClass 声明的存储插件，创建出需要的 PV


#### 本地数据卷
EmptyDir、HostPath只能用于本地文件系统, 所以当Pod发生迁移的时候，数据便会丢失. 该类型Volume的用途是：Pod中容器间的文件共享、共享宿主机的文件系统.

1. EmptyDir
如果Pod配置了EmpyDir数据卷，在Pod的生命周期内都会存在，当Pod被分配到 Node上的时候，会在Node上创建EmptyDir数据卷，并挂载到Pod的容器中. 只要Pod 存在，EmpyDir数据卷都会存在（容器删除不会导致EmpyDir数据卷丟失数据），但是如果Pod的生命周期终结（Pod被删除），EmpyDir数据卷也会被删除，并且永久丢失.

EmpyDir数据卷非常适合实现Pod中容器的文件共享, 比如可以通过一个专职日志收集容器，在每个Pod中和业务容器中进行组合，来完成日志的收集和汇总.

1. HostPath
HostPath数据卷允许将容器宿主机上的文件系统挂载到Pod中, 增加了Pod与节点的耦合. 如果Pod需要使用宿主机上的某些文件，可以使用HostPath, 比如kube-apiserver, kube-controller-manager.

#### 网络数据卷
Kubernetes提供了很多类型的数据卷以集成第三方的存储系统，包括一些非常流行的分布式文件系统，也有在IaaS平台上提供的存储支持，这些存储系统都是分布式的，通过网络共享文件系统.

网络数据卷能够满足数据的持久化需求，Pod通过配置使用网络数据卷，每次Pod创建的时候都会将存储系统的远端文件目录挂载到容器中，数据卷中的数据将被水久保存，即使Pod被删除，只是除去挂载数据卷，数据卷中的数据仍然保存在存储系统中，且当新的Pod被创建的时候，仍是挂载同样的数据卷. 网络数据卷包含以下几种：NFS、iSCISI、GlusterFS、RBD（Ceph Block Device）、Flocker、AWS Elastic Block Store、GCE Persistent Disk.

#### 信息数据卷
Kubernetes中有一些数据卷，主要用来给容器传递配置信息. 比如Secret（处理敏感配置信息，密码、Token等）、Downward API（通过环境变量的方式告诉容器Pod的信息）、Git Repo（将Git仓库下载到Pod中），都是将Pod的信息以文件形式保存，然后以数据卷方式挂载到容器中，容器通过读取文件获取相应的信息.

#### Projected Volume(投射数据卷 from Kubernetes v1.11)
kubernetes 支持的 Projected Volume 一共有四种：
- Secret
- ConfigMap
- Downward API
- ServiceAccountToken

其实，Secret、ConfigMap，以及 Downward API 这三种 Projected Volume 定义的信息，还可以通过环境变量的方式出现在容器里. 但是，通过环境变量获取这些信息的方式，**不具备自动更新的能力**. 所以，建议使用 Volume 文件的方式访问这些信息.

> Kubernetes 项目的 Projected Volume 其实只有三种，因为第四种 ServiceAccountToken，只是一种特殊的 Secret 而已

#### Secret
把 Pod 想要访问的加密数据存放到 Etcd 中, 避免直接在配置中保存敏感信息, 再通过在 Pod 里挂载 Volume 的方式，访问这些 Secret 里保存的信息(以文件的形式出现在了容器的 Volume 目录里)
支持以Volume(支持动态更新, 由kubelet 定时维护这些 Volume)或Env方式加载.

格式: key = base64(value)

```
$ echo -n admin | base64
YWRtaW4=
$ echo -n YWRtaW4= | base64 --decode
admin
```

### ConfigMap
很多生产环境中的应用程序配置较为复杂，可能需要多个Config文件、命令行参数和环境变量的组合. 并且这些配置信息应该从应用程序镜像中解耦出来，以保证镜像的可移植性以及配置信息不被泄露.

ConfigMap包含了一系列的键值对，用于存储被Pod或者系统组件（如controller）访问的信息.

与 Secret 类似，它与 Secret 的区别在于，ConfigMap 保存的是不需要加密的、应用所需的配置信息. 而  其用法几乎与 Secret 完全相同.

#### Downward API
让 Pod 里的容器能够直接获取到这个 Pod 本身的信息,比如metadata

注意: Downward API 能够获取到的信息，一定是**Pod 里的容器进程启动之前就能够确定下来的信息**. 而如果你想要获取 Pod 容器运行后才会出现的信息，比如，容器进程的 PID，那就肯定不能使用 Downward API 了，而应该考虑在 Pod 里定义一个 sidecar 容器

#### ServiceAccountToken
Service Account 就是 Kubernetes 系统内置的一种"服务账户"，它是 Kubernetes 进行权限分配的对象.

Service Account 的授权信息和文件，实际上保存在它所绑定的一个特殊的 Secret 对象里的. 这个特殊的 Secret 对象，就叫作ServiceAccountToken, 它是Pod 里使用 Kubernetes 的 Client时进行授权的方法.

每一个 运行在k8s里的Pod都已经自动声明一个类型是 Secret、名为 default-token-xxxx 的 Volume, 它就是ServiceAccountToken.

> 把 Kubernetes 客户端以容器的方式运行在集群里，然后使用 default Service Account 自动授权的方式，被称作"InClusterConfig".

### StatefulSet
适合于有状态服务, 比如数据库服务MySQL和PostgreSQL，集群化管理服务ZooKeeper、etcd等.

StatefulSet 将应用状态抽象为了两种情况：
- 拓扑状态 : 应用的多个实例之间不是完全对等的关系, 它们必须按照某些顺序启动，比如应用的主节点 A 要先于从节点 B 启动. 而如果把 A 和 B 两个 Pod 删除掉，它们再次被创建出来时也必须严格按照这个顺序才行. 并且新创建出来的 Pod，必须和原来 Pod 的网络标识一样，这样原先的访问者才能使用同样的方法，访问到这个新 Pod. 
  解决方案: DNS
- 存储状态 : 应用的多个实例分别绑定了不同的存储数据. 对于这些应用实例来说，Pod A 第一次读取到的数据，和隔了十分钟之后再次读取到的数据，应该是同一份，哪怕在此期间 Pod A 被重新创建过. 这种情况最典型的例子就是一个数据库应用
  解决方案: PV+PVC, 此时 PVC 的名字会被分配一个与这个 Pod 完全一致的编号, 命名规则`${PVC 名字}-${Pod_Name}`

StatefulSet做的只是将确定的Pod与确定的存储关联起来保证状态的连续性. 通过使用 Pod 模板创建 Pod 的时候，对它们进行编号，并且按照编号顺序逐一完成创建/调谐工作

Kubernetes 将 Pod 的拓扑状态（比如：哪个节点先启动，哪个节点后启动），按照 Pod 的`${StatefulSet_Name}-${序号, 从0开始}`的方式固定了下来. Kubernetes 还为每一个 Pod 提供了一个固定并且唯一的访问入口即这个 Pod 对应的 DNS 记录. 这些状态，在 StatefulSet 的整个生命周期里都会保持不变，绝不会因为对应 Pod 的删除或者重新创建而失效. 

注意: 解析到的 Pod 的 IP 地址，并不是固定的. 对于`有状态应用`实例的访问必须使用 DNS 记录或者 hostname 的方式，而绝不应该直接访问这些 Pod 的 IP 地址.

StatefulSet 的控制器直接管理的是 Pod, 通过 Headless Service为这些有编号的 Pod在 DNS 服务器中生成带有同样编号的 DNS 记录. StatefulSet 还为每一个 Pod 分配并创建一个同样编号的 PVC以维持存储状态.

设计思想: StatefulSet 其实就是一种特殊的 Deployment，而其独特之处在于，它的每个 Pod 都被编号了。而且，这个编号会体现在 Pod 的名字和 hostname 等标识信息上，这不仅代表了 Pod 的创建顺序，也是 Pod 的重要网络标识. 有了这些编号后，StatefulSet 就使用 Kubernetes 里的两个标准功能：Headless Service 和 PV/PVC，实现了对 Pod 的拓扑状态和存储状态的维护.

### Namespace
Namespace将一个物理cluster逻辑上划分为多个虚拟的cluster, 不同Namespace间的资源是完全隔离的.

k8s默认会创建两个Namespace:
- default : 创建资源时如果不指定Namespace则会放入这里
- kube-system : k8s自己创建的系统资源会放入这里

### 容器
一种轻量级的虚拟化技术．容器间是隔离(基于Linux Namespace实现)的.

### 镜像
k8s镜像下载策略：
- Always：每次都下载最新的镜像
- Never：只使用本地镜像，从不下载
- IfNotPresent：只有当本地没有的时候才下载镜像

### Helm
k8s应用打包工具, 支持:
- 从零创建chart
- 与存储chart的repo交互, 拉取, 保存和更新 chart
- 安装卸载chart
- 更新, 回滚和测试 chart

概念:
- chart : 一个应用的信息集合, 包含各种k8s 对象的配置模板, 参数定义, 依赖关系, 文档等.
- release : chart 的运行实例. chart能够多次安装到同一个集群

## 架构
k8s cluster 由 master 和 node 组成.

master运行着kube-apiserver、
kube-scheduler、kube-controller-manager、etcd和pod网络.

node运行着pod、pod网络、kubelet、kube-proxy、Runtime.

### kube-dns
为cluster提供dns服务

### kube-apiserver
提供k8s api(RESTFul API), 是k8s cluster的前端接口(即系统管理命令的统一入口).

> kubectl(k8s的客户端)就是通过kube-apiserver提供的api与k8s cluster交互

### kube-scheduler
负责将pod调度到哪个node上.

scheduler 考虑维度: cluster的拓扑结构, 节点负载, 应用对高可用,性能,数据亲和性等

### kube-controller-manager
是一系列控制器的集合. 它管理 cluster的各种资源, 保证资源处于预期状态

controller分类:
- replication controller : 管理 deployment, statefulset, daemonset
- endpoints controller
- namespace controller : 管理namespace
- serviceaccounts controller

类似 Deployment controller，实际上都是由上半部分的控制器定义（包括期望状态），加上下半部分的被控制对象的模板组成的.

> 如果认为kube-apiserver是前台, 那么kube-controller-manager就是后台
> k8s源码的`pkg/controller`里是对应的各种controller.

controllers都遵循 Kubernetes 项目中的一个通用编排模式：控制循环（control loop）.
```
for {
  实际状态 := 获取集群中对象 X 的实际状态（Actual State）
  期望状态 := 获取集群中对象 X 的期望状态（Desired State）
  if 实际状态 == 期望状态{
    什么都不做
  } else {
    执行编排动作，将实际状态调整为期望状态 // 这步通常被叫作调谐（Reconcile）. 这个调谐的过程，则被称作"Reconcile Loop"（调谐循环）或者"Sync Loop"（同步循环）
  }
}
```go

实际状态往往来自于 Kubernetes 集群本身: kubelet 通过心跳汇报的容器状态和节点状态，或者监控系统中保存的应用监控数据，或者控制器主动收集的它自己感兴趣的信息
期望状态一般来自于用户提交的 YAML 文件

#### Deployment Controller
deployment controller创建pod的过程:
1. 通过kubectl创建Deployment
1. Deployment创建ReplicaSet
1. ReplicaSet创建pod

#### Replication Controller(RC)
应用托管在Kubernetes之后，Kubernetes需要保证应用能够持续运行，这是RC的工作内容，它会确保任何时间Kubernetes中都有指定数量的Pod在运行. 在此基础上，RC还提供了一些更高级的特性，比如滚动升级、升级回滚等

RC与Pod的关联是通过Label来实现的, 可在Pod.metadata.labeks中进行设置. 另外Lable是不具有唯一性的，为了更准确的标识一个Pod，应该为Pod设置多个维度的label.

修改了对应Pod的Label，就会使Pod脱离RC的控制. 同样，在RC运行正常的时候，若试图继续创建同样Label的Pod，是创建不出来的. 因为RC认为副本数已经正常了，再多起的话会被RC删掉的.

##### 弹性伸缩
弹性伸缩是指适应负载变化，以弹性可伸缩的方式提供资源。反映到Kubernetes中，指的是可根据负载的高低动态调整Pod的副本数量(通过修改RC中Pod的副本是来实现的).

```bash
# 扩容Pod的副本数目到10
$ kubectl scale relicationcontroller yourRcName --replicas=10

# 缩容Pod的副本数目到1
$ kubectl scale relicationcontroller yourRcName --replicas=1
```

#### 滚动升级
滚动升级是一种平滑过渡的升级方式，通过逐步替换的策略，保证整体系统的稳定，在初始升级的时候就可以及时发现、调整问题，以保证问题影响度不会扩大.

```bash
$ kubectl rolling-update my-rcName-v1 -f my-rcName-v2-rc.yaml --update-period=10s
```

升级开始后，首先依据提供的定义文件创建V2版本的RC，然后每隔10s（--update-period=10s）逐步的增加V2版本的Pod副本数，逐步减少V1版本Pod的副本数. 升级完成之后，删除V1版本的RC，保留V2版本的RC，即实现了滚动升级.

升级过程中，发生了错误中途退出时，可以选择继续升级。Kubernetes能够智能的判断升级中断之前的状态，然后紧接着继续执行升级。当然，也可以进行回退：
```bash
$ kubectl rolling-update my-rcName-v1 -f my-rcName-v2-rc.yaml --update-period=10s --rollback
```

回退的方式实际就是升级的逆操作，逐步增加V1.0版本Pod的副本数，逐步减少V2版本Pod的副本数.

#### 新一代副本控制器Replica set
这里所说的Replica set，可以被认为 是“升级版”的Replication Controller. 也就是说, Replica set也是用于保证与Label Selector匹配的Pod数量维持在期望状态. 区别在于，Replica set引入了对基于子集的selector查询条件，而Replication Controller仅支持基于值相等的selecto条件查询, 这是目前从用户角度肴，两者唯一的显著差异.

社区引入这一API的初衷是用于取代vl中的Replication Controller，也就是说．**当v1版本被废弃时，Replication Controller就完成了它的历史使命**，而由Replica set来接管其工作. 虽然Replica set可以被单独使用，但是目前它多被Deployment用于进行pod的创建、更新与删除. Deployment在滚动更新等方面提供了很多非常有用的功能.

### RBAC
在 Kubernetes 项目中，负责完成授权（Authorization）工作的机制，就是 RBAC.
- Role/RoleBinding
- ClusterRole/ClusterRoleBinding # 对于非 Namespaced（Non-namespaced）对象（比如：Node），或者，某一个 Role 想要作用于所有的 Namespace 的时候;在 Kubernetes 中已经内置了很多个为系统保留的 ClusterRole，它们的名字都以 system: 开头
- ServiceAccount #  Kubernetes 负责管理的“内置用户”

Kubernetes 还拥有“用户组”（Group）的概念.

Kubernetes 还提供了四个预先定义好的 ClusterRole 来供用户直接使用:
- cluster-admin : 整个 Kubernetes 项目中的最高权限（verbs=*）
- admin
- edit
- view

### Operator
一个相对更加灵活和编程友好的管理“有状态应用”的解决方案. 它实际上是利用了 Kubernetes 的自定义 API 资源（CRD）来描述想要部署的“有状态应用”, 然后在自定义控制器里根据自定义 API 对象的变化，来完成具体的部署和运维工作

### etcd
保持 k8s cluster的配置信息和各种资源的状态

### pod网络
负责pod间的通信, 比如flannel, canal, calico, weave

### node
pod运行的地方, 其上运行的相关组件有kubelet,kube-proxy和pod网络.

#### kubelet
是node的agent, 负责维护和管理该Node上面的所有容器. scheduler选中该node后会将pod的具体配置信息(image,volume等)发送到kubelet, 由kubelet依据这些信息创建和运行容器, 并向master报告运行状态. 本质上，它负责使Pod得运行状态与期望的状态一致.

此外，kubelet 还通过 gRPC 协议同一个叫作 Device Plugin 的插件进行交互. 这个插件是 Kubernetes 项目用来管理 GPU 等宿主机物理设备的主要组件，也是基于 Kubernetes 项目进行机器学习训练、高性能作业支持等工作必须关注的功能.

而kubelet 的另一个重要功能，则是调用网络插件和存储插件为容器配置网络和持久化存储. 这两个插件与 kubelet 进行交互的接口，分别是 CNI（Container Networking Interface）和 CSI（Container Storage Interface）.

> 如果容器不是通过Kubernetes创建的，它并不会管理.

#### kube-proxy
实现了Kubernetes中的服务发现和反向代理功能. 

反向代理方面: 负责将访问service的tcp/udp数据流转发到具体的容器上. 有多个容器副本时, 它还能实现负载均衡(默认基于Round Robin算法)
服务发现方面: 使用etcd的watch机制，监控集群中Service和Endpoint对象数据的动态变化，并且维护一个Service到Endpoint的映射关系，从而保证了后端Pod的IP变化不会对访问者造成影响. 另外kube-proxy还支持session affinity.

#### pod
##### 用label控制pod位置
默认情况下, scheduler会将pod调度到所有可用的node. 但k8s可以通过label来实现调度pod到指定node:
1. 将指定node打上label: `kubectl label node xxx disktype=ssd`
1. 在yaml的pod模板的spec里指定nodeSelector
    ```yaml
    apiVersion: v1
    kind: Deployment
    ...
    spec:
      containers:
      - name: nginx
        image: nginx
        imagePullPolicy: IfNotPresent
      nodeSelector:
        disktype: ssd
    ```
nodeSelector即将废弃, 用nodeAffinity代替, 它支持支持更加丰富的语法, 比如`operator: In`

在部署 Kubernetes 集群的时候，能够先部署 Kubernetes 本身、再部署网络插件的根本原因: 网络插件是DaemonSet类型. 在 Kubernetes 项目中，当一个节点的网络插件尚未安装时，这个节点就会被自动加上名为node.kubernetes.io/network-unavailable的“污点”. 而通过这样一个 Toleration，调度器在调度DaemonSet创建的 Pod 时，就会忽略当前节点上的“污点”，从而成功地将网络插件的 Agent 组件调度到这台机器上启动起来.

#### HostAliases
定义了 Pod 的 hosts 文件（比如 /etc/hosts）里的内容
```yaml
apiVersion: v1
kind: Pod
...
spec:
  hostAliases:
  - ip: "10.1.1.1"
    hostnames:
    - "a.remote" // Pod's /etc/hosts 会存在`10.1.1.1 a.remote`的记录
...
```

#### 共享linux namespace
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  shareProcessNamespace: true # 容器间共享 PID Namespace
  containers:
  - name: nginx
    image: nginx
  - name: shell
    image: busybox
    stdin: true # 必须设置tty和stdin, 否则kubectl attach时会报错: Unable to use a TTY - container es-node did not allocate one
    tty: true
```
或
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  hostNetwork: true # 容器与主机间共享 Network Namespace
  hostIPC: true # 容器与主机间共享 IPC Namespace
  hostPID: true # 容器与主机间共享 PID Namespace
  containers:
  - name: nginx
    image: nginx
  - name: shell
    image: busybox
    stdin: true
    tty: true
```

#### kubectl exec 与 kubectl attach
[kubectl attach](http://kubernetes.kansea.com/docs/user-guide/kubectl/kubectl_attach/)连接到现有容器中一个正在运行的进程
```
kubectl attach POD -c CONTAINER // -c 容器名. 如果省略，则默认选择第一个 pod
```

kubectl exec在容器中执行命令:
```
kubectl exec POD [-c CONTAINER] -- COMMAND [args...]
```


#### Runtime
容器运行环境，目前Kubernetes支持Docker和rkt两种容器

> k8s对象命名方式: 父对象名 + "-" + 随机字符串(字母+数字)
> 出于安全考虑, 默认下k8s不会将pod调度到master节点,可使用[`kubectl taint`修改](https://kubernetes.io/zh/docs/concepts/configuration/taint-and-toleration/)

## doc
- [Kubernetes v1.10.x HA 全手动安装教程](https://www.kubernetes.org.cn/3814.html)
- [Kubernetes指南](https://feisky.gitbooks.io/kubernetes/)
- [Kubernetes中文手册](https://www.kubernetes.org.cn/docs)
- [k8s yaml定义](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.9/)
- [Kubernetes网络原理及方案](http://www.youruncloud.com/blog/131.html)
- [详解 Kubernetes Deployment 的实现原理](https://www.tuicool.com/articles/bQZnYjA)

## install
```sh
apt-get update && apt-get install -y apt-transport-https
curl https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | apt-key add -
cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb http://mirrors.aliyun.com/kubernetes/apt/ kubernetes-stretch main
EOF  
apt-get update
apt-get install -y kubelet kubeadm kubectl
```

> [k8s官网doc](https://kubernetes.io/docs/setup/independent/install-kubeadm/)给的安装源是https://packages.cloud.google.com，但国内访问不了，此时我们可以使用[阿里云的仓库镜像](https://opsx.alibaba.com/mirror)

**安装太复杂, 使用[rancher部署,推荐](https://www.cnrancher.com/kubernetes-installation/)来代替**.

> [rancher卸载](https://www.cnrancher.com/docs/rancher/v2.x/cn/configuration/admin-settings/remove-node/)

## rancher部署
[docker 版本支持](http://rancher.com/docs/rancher/v1.6/en/hosts/#supported-docker-versions)

其他参考:
- [yonyoucloud/install_k8s](https://github.com/yonyoucloud/install_k8s)
- [centos7 使用二进制包搭建kubernetes 1.9.0集群](https://www.58jb.com/html/180.html)

## binary部署(未完成)
如何获取k8s v1.9.5 binary:
- 根据github k8s repo的CHANGELOG-x.y.md,下载`kubernetes.tar.gz`
- 根据解压后的`server/README`指示来获取binary

```
$ cd server
$ tar zxvf kubernetes-server-linux-amd64.tar.gz
$ cd kubernetes/server/bin
$ mkdir log
# on master
$ sudo ./kube-apiserver --v=0 --etcd_servers=http://localhost:12379 --insecure-bind-address=0.0.0.0 --insecure-port=8080 --service-cluster-ip-range=10.254.0.0/16 >> log/kube-apiserver.log 2>&1 &
$ ./kube-controller-manager --v=0 --master=http://0.0.0.0:8080 >> log/kube-controller-manager.log 2>&1 &
$ ./kube-scheduler --v=0 --master=http://0.0.0.0:8080 >> log/kube-scheduler.log 2>&1 &
$ sudo ./kube-proxy --v=0 --master=http://0.0.0.0:8080 >> log/kube-proxy.log 2>&1 &
# on node
$ ./kubelet --v=0 --kubeconfig= --address=0.0.0.0 --api-servers=http://0.0.0.0:8080 >> log/kubelet.log 2>&1 & ## 不会配置kubeconfig, 停止部署
```

### 问题
1. [kubernetes-dashboard:9090/# Service unavailable](https://github.com/rancher/rancher/issues/10650), k8s部分镜像在google服务上,无法访问, 使用[rancher_cn提供的镜像替换(推荐)](https://www.cnrancher.com/kubernetes-installation/),也可参考[rancher安装Kubernetes](https://anjia0532.github.io/2017/11/10/rancher-k8s/)

逐个修改出问题的pod`kubectl --namespace=kube-system edit deployment [heapster|kube-dns|...]`, 修改配置后的pod会自行重试.
kube-dns则要麻烦些, 要先缩容`kubectl --namespace=kube-system scale deployment kube-dns --replicas=0`,再扩容`kubectl --namespace=kube-system scale deployment kube-dns --replicas=1`

1. Failed to deploy addon execute job: Job.batch "rke-network-plugin-deploy-job" is invalid: spec.template.spec.nodeName: Invalid value: "chen-PC"
完整错误：
```
Failed to deploy addon execute job: Job.batch 
"rke-network-plugin-deploy-job" is invalid: spec.template.spec.nodeName:
 Invalid value: "chen-PC": a DNS-1123 subdomain must consist of lower 
case alphanumeric characters, '-' or '.', and must start and end with an
 alphanumeric character (e.g. 'example.com', regex used for validation 
is '[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*')
```

不知rancher/agent怎么获取hostname的, 我用`hostnamectl`和`hostname`工具都修改了本机的hostname(`chen-PC`->`chen-pc`), 删除原有node的rancher/agent,重新注册就OK了.

1. Runtime network not ready: NetworkReady=false reason:NetworkPluginNotReady message:docker: network plugin is not ready: cni config uninitialized
未知,rancher自行解决了.

1. error: the server doesn't have a resource type "cronjobs"
`.kube/config`应在`$HOME`下.

### Kubernetes通过yaml配置文件创建实例时总是重新pull镜像
参考: [Kubernetes通过yaml配置文件创建实例时不使用本地镜像的原因](https://www.58jb.com/html/154.html)

```
spec: 
  containers: 
    - name: nginx 
      image: image: reg.docker.lc/share/nginx:latest 
      imagePullPolicy: IfNotPresent   #或者使用Never 
```
imagePullPolicy:
- Always, 默认
- IfNotPresent ：如果本地存在镜像就优先使用本地镜像。
- Never：直接不再去拉取镜像了，使用本地的；如果本地不存在就报异常了

### 命令 vs 配置文件(yaml)
- 命令 : 简单, 直观, 快捷. 适用于临时环境(测试,实验)
- 配置文件:
    1. 正式, 丰富. 适合正式的, 规模化部署
    1. 可以像管理代码一样进行管理

## 术语
### pod
pod是k8s的基本处理单元, 其包含一个特殊的Pause容器(即根容器)和一个或多个紧密相关的用户业务容器.

设计pod的原因:
1. 引入业务无关且不易死亡的Pause容器来代表整个pod的状态
1. pod里的多个业务容器共享Pause容器的IP和其挂载的volume,简化了密切相关的业务容器间的通信及文件共享问题

### service
参考:
- [kubernetes的网络策略(kube-proxy)流程探究](https://www.jianshu.com/p/13b86daf56dc)
- [华为云在 K8S 大规模场景下的 Service 性能优化实践](https://zhuanlan.zhihu.com/p/37230013)
- [Kubernetes 从1.10到1.11升级记录(续)：Kubernetes kube-proxy开启IPVS模式](https://blog.frognew.com/2018/10/kubernetes-kube-proxy-enable-ipvs.html)

为一组具有相同功能的容器应用提供一个统一的入口, 并将请求进行负载均衡地分发到pod上, 其屏蔽了pod ip的变化, 并通过Label来关联pod.

Endpoint: service通过selector和pod建立关联, k8s会根据service关联到pod的podIP信息组合成一个endpoint.  可通过`kubectl describe service ${service_name}`查看Endpoints.

Service类型:
- ClusterIP(默认) : 只有Cluster内的节点和Pod可以访问
- NodePort(不推荐) : Service通过Cluster节点上的端口(`NodeIP:NodePort`)对外提供服务.
  k8s还是会分配一个ClusterIP, `EXTERNAL-IP`变为`nodes`, `PORT(S)`变为`${Cluster_Port}:${Node_Port,自动分配时的范围:30000~32767}/protocol`.
- LoadBalancer: 借助cloud provider创建一个外部的负载均衡器进行请求转发, 比如转发到ClusterIP/NodePort. 目前[cloud provider(具体实现见代码)](https://github.com/kubernetes/kubernetes/blob/master/pkg/cloudprovider/providers/providers.go)有GCE, AWS, Aliyun等.

kube-proxy 有三种实现 service 的方案, userspace, iptables 和 ipvs:
- userspace(废弃) : 在用户空间监听一个端口，所有的 service 都转发到这个端口，然后 kube-proxy 在内部应用层对其进行转发。因为是在用户空间进行转发，所以效率也不高
- iptables(默认, 推荐) : 完全使用 iptables 来实现 service，是目前默认的方式，也是推荐的方式，效率很高（只有内核中 netfilter 一些损耗）
- ipvs模式 (**最推荐**): IPVS(IP Virtual Server)是lvs项目的一部分，作为Linux内核的一部分，提供4层负载均衡器的功能，即传输层负载均衡

> Kubernetes 原生的 Service 负载均衡基于 Iptables 实现，其规则链会随 Service 的数量呈线性增长，在大规模场景下对 Service 性能影响严重.
> IPVS 在 CPU/内存两个维度的指标都要远远低于 Iptables.

#### ClusterIP
它是一个虚拟ip, 由k8s 节点上的iptables管理.

```sh
$ sudo iptables-save |grep "httpd-svc"
-A KUBE-SERVICES ! -s 10.42.0.0/16 -d 10.43.5.99/32 -p tcp -m comment --comment "default/httpd-svc: cluster IP" -m tcp --dport 8080 -j KUBE-MARK-MASQ // 允许cluster内的pod(源地址来自10.42.0.0/16) 访问httpd-svc
-A KUBE-SERVICES -d 10.43.5.99/32 -p tcp -m comment --comment "default/httpd-svc: cluster IP" -m tcp --dport 8080 -j KUBE-SVC-RL3JAE4GN7VOGDGP //  其他源地址访问httpd-svc, 跳转到规则 KUBE-SVC-RL3JAE4GN7VOGDGP
$ sudo iptables-save |grep KUBE-SVC-RL3JAE4GN7VOGDGP // KUBE-SVC-RL3JAE4GN7VOGDGP的规则
-A KUBE-SVC-RL3JAE4GN7VOGDGP -m statistic --mode random --probability 0.50000000000 -j KUBE-SEP-46JOUSXCF6FAU35R // 一半概率跳到规则 46JOUSXCF6FAU35R
-A KUBE-SVC-RL3JAE4GN7VOGDGP -j KUBE-SEP-UOPS3DAR6EOBA5G6 // 剩下一半概率跳到规则 KUBE-SEP-UOPS3DAR6EOBA5G6
$ sudo iptables-save |egrep "KUBE-SEP-46JOUSXCF6FAU35R|KUBE-SEP-UOPS3DAR6EOBA5G6" // 将请求分别转发到后端的两个pod
-A KUBE-SEP-46JOUSXCF6FAU35R -s 10.42.0.7/32 -j KUBE-MARK-MASQ
-A KUBE-SEP-46JOUSXCF6FAU35R -p tcp -m tcp -j DNAT --to-destination 10.42.0.7:80
-A KUBE-SEP-UOPS3DAR6EOBA5G6 -s 10.42.0.8/32 -j KUBE-MARK-MASQ
-A KUBE-SEP-UOPS3DAR6EOBA5G6 -p tcp -m tcp -j DNAT --to-destination 10.42.0.8:80
```

> cluster每个节点都配置了相同的iptables规则, 确保整个cluster都能够通过Service的cluster ip访问Service.

#### NodePort
与ClusterIP类似(唯一不同之处是NodePort在节点上开了一个占位端口), 将访问当前节点Node_Port端口的请求路由到 ClusterIP 上
```
$ sudo ss -anlpt|grep 31920 // 31920 is Node_Port
LISTEN     0      32768       :::31920                   :::*                   users:(("kube-proxy",pid=2961,fd=6)) // 占用nodePort端口: 为了防止主机上的其它进程使用了该nodePort端口而导致访问冲突
$ sudo iptables -S -t nat|grep 31920
-A KUBE-NODEPORTS -p tcp -m comment --comment "default/httpd-svc:" -m tcp --dport 31920 -j KUBE-MARK-MASQ
-A KUBE-NODEPORTS -p tcp -m comment --comment "default/httpd-svc:" -m tcp --dport 31920 -j KUBE-SVC-RL3JAE4GN7VOGDGP // 访问当前节点31920端口的请求会应用到KUBE-SVC-RL3JAE4GN7VOGDGP, 即NodePort 会路由到 ClusterIP 上.
$ sudo iptables-save |egrep "KUBE-SVC-RL3JAE4GN7VOGDGP"
-A KUBE-SVC-RL3JAE4GN7VOGDGP -m statistic --mode random --probability 0.50000000000 -j KUBE-SEP-46JOUSXCF6FAU35R
-A KUBE-SVC-RL3JAE4GN7VOGDGP -j KUBE-SEP-UOPS3DAR6EOBA5G6
$ sudo iptables-save |egrep "KUBE-SEP-46JOUSXCF6FAU35R|KUBE-SEP-UOPS3DAR6EOBA5G6" // 到具体的pod
-A KUBE-SEP-46JOUSXCF6FAU35R -s 10.42.0.7/32 -j KUBE-MARK-MASQ
-A KUBE-SEP-46JOUSXCF6FAU35R -p tcp -m tcp -j DNAT --to-destination 10.42.0.7:80
-A KUBE-SEP-UOPS3DAR6EOBA5G6 -s 10.42.0.8/32 -j KUBE-MARK-MASQ
-A KUBE-SEP-UOPS3DAR6EOBA5G6 -p tcp -m tcp -j DNAT --to-destination 10.42.0.8:80
``` 

#### 通过dns访问Service
dns功能由kube-dns提供, 每当有新的Service被创建, kube-dns会添加该Service的DNS记录. 这样Cluster中的Pod就可以通过`<Serviced_Name>.<Namespace_Name>`访问Service了, 同namespace访问可直接用`<Serviced_Name>`即可. 它分为两种处理方法:
- Normal Service : `my-svc.my-namespace.svc.cluster.local`解析到的是 my-svc 这个 Service 的 VIP(Virtual IP)
- Headless Service: `my-svc.my-namespace.svc.cluster.local`解析到的是 my-svc 代理的某一个 Pod 的 IP 地址, 即Headless Service 不需要分配一个 VIP, 而是会以 DNS 记录的方式暴露出它所代理的 Pod

Headless Service的DNS 记录格式:`<pod-name>.<svc-name>.<namespace>.svc.cluster.local`

> Headless Service定义: spec.clusterIP=None

如何获取dns所在服务的ip:
- 进入pod, 查看`/etc/resolv.conf`的`nameserver`
- `kubectl get service -n kube-system | grep -i dns`

如何获取Service的dns记录:
- 进入pod, 执行`nslookup httpd-svc`即可. // 在kube-dns和namespace system下的pod不能用`nslookup httpd-svc.defualt`访问namespce default 下的Service???



#### 外网访问Service


### 如何访问
参考:
1. [Kubernetes的三种外部访问方式：NodePort、LoadBalancer 和 Ingress](http://dockone.io/article/4884)
1. [Publishing services - service types](https://kubernetes.io/docs/concepts/services-networking/service/)

Kubernetes 暴露服务的方式目前只有三种：LoadBlancer Service、NodePort Service、Ingress.

- NodePort
NodePort 服务通过外部访问服务的最基本的方式. 在集群的每个node上暴露一个端口，然后将这个端口映射到某个具体的service来实现的，虽然每个node的端口有很多(默认的取值范围是 30000-32767)，但是由于安全性和易用性(服务多了就乱了，还有端口冲突问题)实际使用可能并不多.
- HostPort
它直接将容器的端口与所调度的node上的端口路由，这样用户就可以通过宿主机的IP加上来访问Pod了. 这样做有个缺点，因为Pod重新调度的时候该Pod被调度到的宿主机可能会变动，这样就变化了，用户必须自己维护一个Pod与所在宿主机的对应关系
- ClusterIP
ClusterIP 服务是 Kubernetes **默认**的服务类型. 如果你在集群内部创建一个服务，则**在集群内部的其他应用程序可以进行访问，但是不具备集群外部访问的能力**.
clusterIP 是一个虚拟的IP， cluster ip 仅作用于kubernetes service 这个对象.
在k8s上是由kubernetes proxy负责实现cluster ip路由和转发.
- Load Balancer
LoadBalancer 服务是暴露服务至互联网最标准的方式, 其只能在service上定义, 而且所有通往指定端口的流量都会被转发到对应的服务. 它没有过滤条件，没有路由, 这意味着你几乎可以发送任何种类的流量到该服务，像 HTTP，TCP，UDP，Websocket，gRPC 或其它任意种类.
- Ingress
使用nginx等开源的反向代理负载均衡器实现对外暴露服务，可以理解Ingress就是用于配置域名转发，在nginx中就类似upstream，它与ingress-controller结合使用，通过ingress-controller监控到pod及service的变化，动态地将ingress中的转发信息写到诸如nginx、apache、haproxy等组件中实现方向代理和负载均衡.

Kubernetes Ingress提供了负载平衡器的典型特性：HTTP路由，粘性会话，SSL终止，SSL直通，TCP和UDP负载平衡等。目前并不是所有的Ingress controller都实现了这些功能，需要查看具体的Ingress controller文档, 而且Ingress controller直接将流量转发给后端Pod，不需再经过kube-proxy的转发，比LoadBalancer方式更高效.
- targetPort
targetPort 是Pod上的端口.

## cmd
- `kubectl run alpine --rm -ti --image=alpine /bin/sh` # 创建调试pod
- `kubectl logs -f POD-NAME` # 获取pod日志
- `kubectl edit deployment ${deployment-name}` # 查看deployment配置和运行状态
- `kubectl exec -it POD-NAME sh` # 进入pod的容器
- `kubectl describe node Node-NAME` # 获取node的描述信息
- `kubectl describe pod POD-NAME` # 获取pod的描述信息(简单), 生命周期事件
- `kubectl describe deployment DeploymentName` # 获取deployment的描述信息(简单), 生命周期事件
- `kubectl describe replicaset ReplicasetName` # 获取replicaset的描述信息(简单), 生命周期事件
- `kubectl describe pod nginx-ingress-controller-hv2n6 --namespace=ingress-nginx` # 查看指定namespace的指定pod的状态
- `kubectl get pod myweb-fnncj --output json/yaml` # 获取pod的详细信息, 有状态信息
- `kubectl get pods --all-namespaces [-o wide]` # 获取所有pod的状态,加`-o wide`时还会输出更多信息, 比如ip和node host
- `kkubectl get pod -l app=nginx` # 获取所有lable是`app=nginx`的pods
- `kubectl get events` # 查询所有事件
- `kubectl get pods [--show-all]` # 查询所有pod
- `kubectl get deployments` # 获取所有deployment
- `kubectl get nodes [--show-labels]` # 获取所有node
- `kubectl get replicaset` # 获取所有replicaset
- `kubectl logs POD-NAME Container-NAME [-p]` # 查询pod中容器的日志,`-p`允许查询`Container-NAME`重建前的日志
- `kubectl apply -f nginx-deployment.yaml [--record]` # 统一进行 Kubernetes 对象的创建和更新操作, 是 Kubernetes“声明式 API”(kubectl apply 命令)所推荐的使用方法. 作为用户，你不必关心当前的操作是创建还是更新，你执行的命令始终是 kubectl apply，而 Kubernetes 则会根据 YAML 文件的内容变化，自动进行具体的处理. `--record`用于记录revision历史
- `kubectl api-versions` # 查看api version支持的资源版本
- `kubectl rollout history deployment ${deployment_name} [--revision=2]` # 查看revision历史记录/指定历史版本
- `kubectl rollout undo deployment ${deployment_name} --to-revision=${num}` # 回滚到指定revison版本

## Dynamic Admission Control/Initializer
当一个 Pod 或者任何一个 API 对象被提交给 APIServer 之后，总有一些“初始化”性质的工作需要在它们被 Kubernetes 项目正式处理之前进行, 而这个“初始化”操作的实现，借助的是一个叫作 Admission 的功能. Kubernetes 项目为我们额外提供了一种“热插拔”式的 Admission 机制，它就是 Dynamic Admission Control，也叫作：Initializer. 典型案例就是Istio 项目.



## 声明式API
kubectl replace 的执行过程，是使用新的 YAML 文件中的 API 对象，替换原有的 API 对象；而 kubectl apply，则是执行了一个对原有 API 对象的 PATCH 操作, 类似地，kubectl set image 和 kubectl edit 也是对已有 API 对象的修改.

声明式 API是 Kubernetes 项目编排能力“赖以生存”的核心所在

在 Kubernetes 项目中，一个 API 对象在 Etcd 里的完整资源路径，是由：Group（API 组）、Version（API 版本）和 Resource（API 资源类型）三个部分组成的.

Kubernetes 是如何对 Resource、Group 和 Version 进行解析，从而在 Kubernetes 项目里找到 具体对象的定义:
1. Kubernetes 会匹配 API 对象的组
  对于 Kubernetes 里的核心 API 对象，比如：Pod、Node 等是不需要 Group 的（即它们 的Group 是`""`）, 对于这些 API 对象来说，Kubernetes 会直接在 /api 这个层级进行下一步的匹配过程
  而对于 非核心 API 对象来说，Kubernetes 就必须在 /apis 这个层级里查找它对应的 Group，进而根据 Group 的名字，找到 `/apis/${group_name}`. 这些 API Group 的分类是以对象功能为依据的
1. Kubernetes 会进一步匹配到 API 对象的版本号
  在 Kubernetes 中，同一种 API 对象可以有多个版本，这正是 Kubernetes 进行 API 版本化管理的重要手段
1. Kubernetes 会匹配 API 对象的资源类型

![整个 Kubernetes 里的所有 API 对象是以树形结构来组织的](https://static001.geekbang.org/resource/image/70/da/709700eea03075bed35c25b5b6cdefda.png)
![](APIServer 创建CronJob 对象的过程)(https://static001.geekbang.org/resource/image/df/6f/df6f1dda45e9a353a051d06c48f0286f.png)

[Kubernetes Deep Dive: Code Generation for CRD（Custom Resource Definition）](https://blog.openshift.com/kubernetes-deep-dive-code-generation-customresources/)
在自定义控制器里面，可以通过对自定义 API 对象和默认 API 对象进行协同，从而实现更加复杂的编排功能.

## 网络模型
k8s基于扁平地址空间的网络模型, 采用Container Networking Interface(CNI)规范, 每个Pod有自己的ip, Pod间不需要配置NAT就能直接通信. 同一个Pod贡献容器IP, 能通过localhost通信.

> [网络方案安装](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/)
> [Kube-OVN](https://github.com/alauda/kube-ovn) : 基于 OVN 的 Kubernetes 网络系统
> [ovn-kubernetes提供了一个ovs OVN网络插件，支持underlay和overlay两种模式](https://www.bookstack.cn/read/sdn-handbook/ovs-ovn-kubernetes.md)


## 生态
### 监控
Prometheus Operator, 可使用Helm安装.

### 日志
k8s Elasticsearch addon

### kubeadm
参考:
- [kubeadm init](https://k8smeetup.github.io/docs/reference/setup-tools/kubeadm/kubeadm-init/)
- [kubeadm 实现细节](http://docs.kubernetes.org.cn/829.html)

通过两条指令完成一个 Kubernetes 集群的部署：
```
$ kubeadm init // 创建一个 Master 节点
$ kubeadm join <Master 节点的 IP 和端口 > // 将一个 Node 节点加入到当前集群中
```

有人会问为什么不用容器部署 Kubernetes 呢？那问题就变成了`如何容器化 kubelet`.
除了跟容器运行时打交道外，kubelet 还需在配置容器网络、管理容器数据卷时直接操作宿主机. 到目前为止，在容器里运行 kubelet，依然没有很好的解决办法，因此 kubeadm 把 kubelet 直接运行在宿主机上，然后使用容器部署其他的 Kubernetes 组件.

kubeadm 目前还不能用于生产环境: 因为 kubeadm 目前最欠缺的是如何一键部署一个高可用的 Kubernetes 集群，即Etcd、Master 组件都应该是多节点集群，而不是现在这样的单点.

如果有部署规模化生产环境的需求，推荐使用kops或rancher. 但使用 kubeadm 部署一个 Kubernetes 集群，对于理解 Kubernetes 组件的工作方式和架构确是恰到好处.

#### kubeadm init 的工作流程
1. 执行一系列的检查以确定这台机器可以用来部署 Kubernetes, 这个过程叫`Preflight Checks`:
  Linux 内核的版本必须是否是 3.10 以上？
  Linux Cgroups 模块是否可用？
  机器的 hostname 是否标准？在 Kubernetes 项目里，机器的名字以及一切存储在 Etcd 中的 API 对象，都必须使用标准的 DNS 命名（RFC 1123）
  用户安装的 kubeadm 和 kubelet 的版本是否匹配？
  机器上是不是已经安装了 Kubernetes 的二进制文件？
  Kubernetes 的工作端口 10250/10251/10252 端口是不是已经被占用？
  ip、mount 等 Linux 指令是否存在？
  Docker 是否已经安装？

1. 生成证书
在通过了 Preflight Checks 之后，kubeadm 会生成 Kubernetes 对外提供服务所需的各种证书和对应的目录. Kubernetes 对外提供服务时，除非专门开启“不安全模式”(不推荐)，否则都要通过 HTTPS 才能访问 kube-apiserver. 这就需要为 Kubernetes 集群配置好证书文件.

kubeadm 为 Kubernetes 项目生成的证书文件都放在 Master 节点的 /etc/kubernetes/pki 目录下. 在这个目录下，最主要的证书文件是 ca.crt 和对应的私钥 ca.key.

此外，用户使用 kubectl 获取容器日志等 streaming 操作时，需要通过 kube-apiserver 向 kubelet 发起请求，这个连接也必须是安全的. kubeadm 为这一步生成的是 apiserver-kubelet-client.crt 文件，对应的私钥是 apiserver-kubelet-client.key.

除此之外，Kubernetes 集群中还有 Aggregate APIServer 等特性，也需要用到专门的证书等等. 当然也可选择不让 kubeadm 生成这些证书，而是拷贝现有的证书到如下证书的目录里：
/etc/kubernetes/pki/ca.{crt,key}, 这时kubeadm 就会跳过证书生成的步骤，把它完全交给用户处理.

1. 生成配置文件
证书生成后，kubeadm 接下来会为其他组件生成访问 kube-apiserver 所需的配置文件. 这些文件的路径是：/etc/kubernetes/xxx.conf：
```sh
ls /etc/kubernetes/
admin.conf  controller-manager.conf  kubelet.conf  scheduler.conf
```
这些文件里面记录了当前这个 Master 节点的服务器地址、监听端口、证书目录等信息. 这样对应的客户端（比如 scheduler，kubelet 等），可以直接加载相应的文件，使用里面的信息与 kube-apiserver 建立安全连接.

1. 容器化部署master的组件
接下来，kubeadm 会为 Master 组件生成 Pod 配置文件, 如果没有提供外部 etcd，也会为 etcd 生成一个额外的静态 Pod manifest 文件. 再之后 kube-apiserver、kube-controller-manager、kube-scheduler, etcd(可选)都会被使用 Pod 的方式部署起来.

这时 Kubernetes 集群尚不存在，难道 kubeadm 会直接执行 docker run 来启动这些容器吗？当然不是.

在 Kubernetes 中，有一种特殊的容器启动方法叫做“Static Pod”. 它允许你把要部署的 Pod 的 YAML 文件放在一个指定的目录里. 之后当这台机器上的 kubelet 启动时，它会自动检查这个目录，加载所有的 Pod YAML 文件，然后在这台机器上启动它们. 从这一点也可以看出，kubelet 在 Kubernetes 项目中的地位非常高，在设计上它就是一个完全独立的组件，而其他 Master 组件则更像是辅助性的系统容器.

在 kubeadm 中，Master 组件的 YAML 文件会被生成在 /etc/kubernetes/manifests 路径下:
```sh
$ ls /etc/kubernetes/manifests/
etcd.yaml  kube-apiserver.yaml  kube-controller-manager.yaml  kube-scheduler.yaml
```

如果要修改一个已有集群的配置就需要修改这些 YAML 文件.

Master 容器启动后，kubeadm 会通过 `localhost:6443/healthz` 这个 Master 组件的健康检查 URL，直到检测到 Master 组件完全运行起来.

1. 创建bootstrap token
kubeadm 就会为集群生成一个 bootstrap token, 只要持有这个 token，任何一个安装了 kubelet 和 kubadm 的节点都可以通过 kubeadm join 加入到这个集群当中. 这个 token 的值和使用方法会，会在 kubeadm init 结束后被打印出来.

在 token 生成之后，kubeadm 会将 ca.crt 等 Master 节点的重要信息，通过 ConfigMap 的方式保存在 Etcd 当中，供后续部署 Node 节点使用, 这个 ConfigMap 的名字是 cluster-info.

1. 安装默认插件
kubeadm init 的最后一步是安装默认插件。Kubernetes 默认 kube-proxy 和 DNS 这两个插件是必须安装的. 它们分别用来提供整个集群的服务发现和 DNS 功能. 这两个插件也只是两个容器镜像而已，所以 kubeadm 只要用 Kubernetes 客户端创建两个 Pod 就可以了.

#### kubeadm join 的工作流程
使用kubeadm init 生成的 bootstrap token 就可以在任意一台安装了 kubelet 和 kubeadm 的机器上执行 kubeadm join.

为什么执行 kubeadm join 需要这样一个 token: 任何一台机器想要成为 Kubernetes 集群中的一个节点，就必须在集群的 kube-apiserver 上注册, 但要想跟 apiserver 打交道，这台机器就必须要获取到相应的证书文件（CA 文件）. 可是为了能够一键安装，我们就不能让用户去 Master 节点上手动拷贝这些文件. 因此 kubeadm 至少需要发起一次“不安全模式”的访问到 kube-apiserver，从而拿到保存在 ConfigMap 中的 cluster-info（它保存了 APIServer 的授权信息）, 而 bootstrap token扮演的就是这个过程中的安全验证的角色.

只要有了 cluster-info 里的 kube-apiserver 的地址、端口、证书，kubelet 就可以以“安全模式”连接到 apiserver 上，这样一个新的节点就部署完成了.

接下来，你只要在其他节点上重复这个指令就可以了.

#### kubeadm 的部署参数
推荐在使用 kubeadm init 部署 Master 节点使用配置文件：
```sh
$ kubeadm init --config kubeadm.yaml  // 为kubeadm 提供一个 YAML 文件（比如kubeadm.yaml），通过它配置参数
```