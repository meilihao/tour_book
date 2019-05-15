# k8s

## 概念
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

### node
node是具体负责运行容器的应用, 会监控并汇报容器状态, 同时会根据master的要求管理容器的生命周期. node由master管理.

### pod
pod是k8s最小的工作(调度,扩展,共享资源,管理生命周期)单位. 每个pod包含一个或多个容器, 它们会作为一个整体被master进行管理, 其原因是:
1. 可管理性: 某些容器的应用间存在紧密联系
1. 共享通信和存储: pod中的容器使用同一个网络namespace.

### service
定义了外界访问一组特定pod的方法.

### Namespace
Namespace将一个物理cluster逻辑上划分为多个虚拟的cluster, 不同Namespace间的资源是完全隔离的.

k8s默认会创建两个Namespace:
- default : 创建资源时如果不指定Namespace则会放入这里
- kube-system : k8s自己创建的系统资源会放入这里

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
- `kubectl describe pods/POD-NAME` # 获取pod的描述信息(简单), 生命周期事件
- `kubectl get pod myweb-fnncj --output json/yaml` # 获取pod的详细信息, 有状态信息
- `kubectl get event` # 查询所有事件
- `kubectl logs POD-NAME Container-NAME [-p]` # 查询pod中容器的日志,`-p`允许查询`Container-NAME`重建前的日志