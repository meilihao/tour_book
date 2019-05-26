# k8s
参考:
- [Kubernetes核心概念总结](http://dockone.io/article/8866)

## 概念

### 容器
是一种沙盒技术. 容器技术的核心功能，就是通过约束和修改进程的动态表现，从而为其创造出一个“边界”. “敏捷”和“高性能”是容器相较于虚拟机最大的优势.

对于 Docker 等大多数 Linux 容器来说，Cgroups 技术是用来制造**约束**的主要手段，而Namespace 技术则是用来修改**进程视图(容器进程看待整个操作系统的视图)**的主要方法

> Cgroups是进程设置资源限制的方法, Namespace是进程隔离的方法.
> 容器是一个“单进程”模型
> 容器本身的设计，就是希望容器和应用能够同生命周期，这个概念对后续的容器编排非常重要. 否则，一旦出现类似于“容器是正常运行的，但是里面的应用早已经挂了”的情况，编排系统处理起来就非常麻烦了
> 容器所谓的一致性：无论在本地、云端，还是在一台任何地方的机器上，用户只需要解压打包好的容器镜像，那么这个应用运行所需要的完整的执行环境就被重现出来了
> 对一个应用来说，操作系统本身才是它运行所需要的最完整的“依赖库”

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
- ReplicaSet : 实现了Pod的多副本管理. 使用Deployment时会自动创建ReplicaSet, 因此我们通常不直接使用它
- DaemonSet : 用于node最多只运行一个pod副本的场景
- StatefuleSet : 保证pod的每个副本在整个生命周期中的名称是不变的(因故障需删除并重启除外), 同时会保证副本按照固定的顺序启动,更新或删除
- Job

#### job
从程序的运行形态上来区分，我们可以将Pod分为两类：
- 长时运行服务（JBoss、MySQL等）
- 一次性任务（数据计算、测试）

RC创建的Pod都是长时运行的服务，而Job创建的Pod都是一次性任务.

在Job的定义中，restartPolicy（重启策略）只能是Never和OnFailure. Job可以控制一次性任务的Pod的完成次数（Job.spec.completions）和并发执行数（Job.spec.parallelism），当Pod成功执行指定次数后，即认为Job执行完毕.

### node
node是具体负责运行容器的应用, 会监控并汇报容器状态, 同时会根据master的要求管理容器的生命周期. node由master管理.

### pod
pod是k8s最小的工作(调度,扩展,共享资源,管理生命周期)单位. 它是一个或多个相关容器的集合, 其原因是:
1. 可管理性: 某些容器的应用间存在紧密联系
1. 共享通信和存储: pod中的容器使用同一个网络namespace.

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
No events.
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

> 如果容器不是通过Kubernetes创建的，它并不会管理.

#### kube-proxy
实现了Kubernetes中的服务发现和反向代理功能. 

反向代理方面: 负责将访问service的tcp/udp数据流转发到具体的容器上. 有多个容器副本时, 它还能实现负载均衡(默认基于Round Robin算法)
服务发现方面: 使用etcd的watch机制，监控集群中Service和Endpoint对象数据的动态变化，并且维护一个Service到Endpoint的映射关系，从而保证了后端Pod的IP变化不会对访问者造成影响. 另外kube-proxy还支持session affinity.

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
- `kubectl exec -it POD-NAME sh` # 进入pod的容器
- `kubectl describe pod POD-NAME` # 获取pod的描述信息(简单), 生命周期事件
- `kubectl describe deployment DeploymentName` # 获取deployment的描述信息(简单), 生命周期事件
- `kubectl describe replicaset ReplicasetName` # 获取replicaset的描述信息(简单), 生命周期事件
- `kubectl describe pod nginx-ingress-controller-hv2n6 --namespace=ingress-nginx` # 查看指定namespace的指定pod的状态
- `kubectl get pod myweb-fnncj --output json/yaml` # 获取pod的详细信息, 有状态信息
- `kubectl get pods --all-namespaces [-o wide]` # 获取所有pod的状态,加`-o wide`时还会输出更多信息, 比如ip和node host
- `kubectl get event` # 查询所有事件
- `kubectl get deployment` # 获取所有deployment
- `kubectl get nodes` # 获取所有node
- `kubectl get replicaset` # 获取所有replicaset
- `kubectl logs POD-NAME Container-NAME [-p]` # 查询pod中容器的日志,`-p`允许查询`Container-NAME`重建前的日志