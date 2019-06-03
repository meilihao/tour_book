# k8s
参考:
- [Kubernetes核心概念总结](http://dockone.io/article/8866)
- [Kubernetes架构为什么是这样的？](https://www.tuicool.com/articles/J7Rbimu)
- [k8s yaml apiVersion](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.14/)
- [kubernetes/examples](https://github.com/kubernetes/examples)

Kubernetes 最主要的设计思想是从更宏观的角度，以统一的方式来定义任务之间的各种关系，并且为将来支持更多种类的关系留有余地.

## 概念
### 基于k8s的开发模式
而当应用本身发生变化时，开发人员和运维人员可以依靠容器镜像来进行同步；当应用部署参数发生变化时，YAML 配置文件就是他们相互沟通和信任的媒介.

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

#### deployment
部署pod并分布到各个node上, 每个node允许存在多个副本.

#### daemonset
每个node至多运行一个副本, k8s本身就在用daemonset运行部分组件`kubectl get daemonsets --all-namespaces`.

daemonset的典型应用场景:
1. 存储: ceph
1. 日志收集: flunentd, logstash
1. 监控: prometheus node exporter, collected

#### job
从程序的运行形态上来区分，我们可以将Pod分为两类：
- 长时运行服务（JBoss、MySQL等）
- 一次性任务（数据计算、测试）

RC创建的Pod都是长时运行的服务，而Job创建的Pod都是一次性任务.

在Job的定义中，restartPolicy（重启策略）只能是Never和OnFailure. Job可以控制一次性任务的Pod的完成次数（Job.spec.completions）和并发执行数（Job.spec.parallelism），当Pod成功执行指定次数后，即认为Job执行完毕.

可用`kubectl get pods --show-all`查看已`Completed`的pod
job的并行性可通过设置parallelism.
job设置completion可指定完成job所需要的pod总数,即job执行次数.
job还可以是定时Job, yaml的kind = CronJob.

### node
node是具体负责运行容器的应用, 会监控并汇报容器状态, 同时会根据master的要求管理容器的生命周期. node由master管理.

### pod
pod是k8s最小的工作(调度,扩展,共享资源,管理生命周期)单位. 它是一个或多个相关容器的集合, 其原因是:
1. 可管理性: 某些容器的应用间存在紧密联系
1. 共享通信和存储: pod中的容器使用同一个网络namespace.

围绕着容器和 Pod 不断向真实的技术场景扩展，我们就能够摸索出一幅如下所示的 Kubernetes 项目核心功能的“全景图”:
![Kubernetes 项目核心功能的“全景图](https://static001.geekbang.org/resource/image/16/06/16c095d6efb8d8c226ad9b098689f306.png)

> pod类似于进程组的或虚拟机的角色, 而容器就是里面的进程
> 凡是调度、网络、存储，以及安全相关的属性，基本上是 Pod 级别的; 凡是跟容器的 Linux Namespace 相关的属性，也一定是 Pod 级别的, 因为Pod 的设计就是要让它里面的容器尽可能多地共享 Linux Namespace，仅保留必要的隔离和限制能力
> Pod是一组共享了某些资源的容器, Pod 里的所有容器，共享的是同一个 Network Namespace，并且可以声明共享同一个 Volume

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
$ kubectl rollout pause deployment nginx-deployment2  #暂停升级
$ kubectl rollout resume deployment nginx-deployment2  #继续升级
$ kubectl rollout undo deployment nginx-deployment2  #升级回滚
$ kubectl scale deployment nginx-deployment --replicas 10  #弹性伸缩Pod数量
```

#### Horizontal Pod Autoscaler
Horizontal Pod Autoscaler的操作对象是Replication Controller、ReplicaSet或Deployment对应的Pod，根据观察到的实际资源使用量与用户的期望值进行比对，做出是否需要增减实例数量的决策.

### volume
在Kubernetes中，当Pod重建的时候，数据是会丢失的，Kubernetes也是通过数据卷挂载来提供Pod数据的持久化的. Kubernetes数据卷是对Docker数据卷的扩展，Kubernetes数据卷是Pod级别的，可以用来实现Pod中容器的文件共享. 目前，Kubernetes支持的数据卷类型如下：
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
Persistent Volume是由系统管理员配置创建的一个数据卷（目前支持HostPath、GCE Persistent Disk、AWS Elastic Block Store、NFS、iSCSI、GlusterFS、RBD），它代表了某一类存储插件实现；而对于普通用户来说，通过Persistent Volume Claim可请求并获得合适的Persistent Volume，而无须感知后端的存储实现. Persistent Volume和Persistent Volume Claim相互关联，有着完整的生命周期管理：
- 准备：系统管理员规划或创建一批Persistent Volume；
- 绑定：用户通过创建Persistent Volume Claim来声明存储请求，Kubernetes发现有存储请求的时候，就去查找符合条件的Persistent Volume（最小满足策略）。找到合适的就绑定上，找不到就一直处于等待状态；
- 使用：创建Pod的时候使用Persistent Volume Claim；
- 释放：当用户删除绑定在Persistent Volume上的Persistent Volume Claim时，Persistent Volume进入释放状态，此时Persistent Volume中还残留着上一个Persistent Volume Claim的数据，状态还不可用；
- 回收：是否的Persistent Volume需要回收才能再次使用。回收策略可以是人工的也可以是Kubernetes自动进行清理（仅支持NFS和HostPath）

#### 本地数据卷
EmptyDir、HostPath只能用于本地文件系统, 所以当Pod发生迁移的时候，数据便会丢失. 该类型Volume的用途是：Pod中容器间的文件共享、共享宿主机的文件系统.

1. EmptyDir
如果Pod配置了EmpyDir数据卷，在Pod的生命周期内都会存在，当Pod被分配到 Node上的时候，会在Node上创建EmptyDir数据卷，并挂载到Pod的容器中。只要Pod 存在，EmpyDir数据卷都会存在（容器删除不会导致EmpyDir数据卷丟失数据），但是如果Pod的生命周期终结（Pod被删除），EmpyDir数据卷也会被删除，并且永久丢失。

EmpyDir数据卷非常适合实现Pod中容器的文件共享, 比如可以通过一个专职日志收集容器，在每个Pod中和业务容器中进行组合，来完成日志的收集和汇总.

1. HostPath
HostPath数据卷允许将容器宿主机上的文件系统挂载到Pod中. 如果Pod需要使用宿主机上的某些文件，可以使用HostPath

#### 网络数据卷
Kubernetes提供了很多类型的数据卷以集成第三方的存储系统，包括一些非常流行的分布式文件系统，也有在IaaS平台上提供的存储支持，这些存储系统都是分布式的，通过网络共享文件系统.

网络数据卷能够满足数据的持久化需求，Pod通过配置使用网络数据卷，每次Pod创建的时候都会将存储系统的远端文件目录挂载到容器中，数据卷中的数据将被水久保存，即使Pod被删除，只是除去挂载数据卷，数据卷中的数据仍然保存在存储系统中，且当新的Pod被创建的时候，仍是挂载同样的数据卷. 网络数据卷包含以下几种：NFS、iSCISI、GlusterFS、RBD（Ceph Block Device）、Flocker、AWS Elastic Block Store、GCE Persistent Disk.

#### 信息数据卷
Kubernetes中有一些数据卷，主要用来给容器传递配置信息. 比如Secret（处理敏感配置信息，密码、Token等）、Downward API（通过环境变量的方式告诉容器Pod的信息）、Git Repo（将Git仓库下载到Pod中），都是将Pod的信息以文件形式保存，然后以数据卷方式挂载到容器中，容器通过读取文件获取相应的信息.

### StatefulSet
适合于有状态服务, 比如数据库服务MySQL和PostgreSQL，集群化管理服务ZooKeeper、etcd等.

StatefulSet做的只是将确定的Pod与确定的存储关联起来保证状态的连续性.

### ConfigMap
很多生产环境中的应用程序配置较为复杂，可能需要多个Config文件、命令行参数和环境变量的组合. 并且这些配置信息应该从应用程序镜像中解耦出来，以保证镜像的可移植性以及配置信息不被泄露.

ConfigMap包含了一系列的键值对，用于存储被Pod或者系统组件（如controller）访问的信息. 这与secret的设计理念有异曲同工之妙，它们的主要区别在于ConfigMap通常不用于存储敏感信息，而只存储简单的文本信息.

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
管理 cluster的各种资源, 保证资源处于预期状态

controller分类:
- replication controller : 管理 deployment, statefulset, daemonset
- endpoints controller
- namespace controller : 管理namespace
- serviceaccounts controller

> 如果认为kube-apiserver是前台, 那么kube-controller-manager就是后台

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

**安装太复杂, 使用[rancher部署,推荐](https://www.cnrancher.com/kubernetes-installation/)来代替**

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
为一组具有相同功能的容器应用提供一个统一的入口, 并将请求进行负载均衡地分发到pod上, 其屏蔽了pod ip的变化, 并通过Label来关联pod.

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
- `kubectl apply -f nginx-deployment.yaml` # 统一进行 Kubernetes 对象的创建和更新操作, 是 Kubernetes“声明式 API”所推荐的使用方法. 作为用户，你不必关心当前的操作是创建还是更新，你执行的命令始终是 kubectl apply，而 Kubernetes 则会根据 YAML 文件的内容变化，自动进行具体的处理
- `kubectl api-versions` # 查看api version支持的资源版本

## 生态
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