# Kubernetes 集群 DNS 服务发现原理
参考:
- [Kubernetes 集群 DNS 服务发现原理](https://developer.aliyun.com/article/779121?spm=a2c6h.12873581.0.dArticle779121.32cc88143C3J6a)

前提:
1. 拥有一个 Kubernetes 集群
1. 能通过 kubectl 连接 Kubernetes 集群

Kubernetes 集群中部署了一套 DNS 服务，通过 coredns 的服务名暴露 DNS 服务, 可执行以下命令查看 coredns 的服务详情:
```bash
kubectl get svc coredns -n kube-system
NAME      TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)                  AGE
coredns   ClusterIP   10.233.0.3   <none>        53/UDP,53/TCP,9153/TCP   12h
```

该服务的后端是两个名为 coredns 的 Pod, 可执行以下命令查看 coredns 的 Pod 详情:
```bash
$ kubectl get deployment coredns -n kube-system
NAME      READY   UP-TO-DATE   AVAILABLE   AGE
coredns   2/2     2            2           12h
```

> kube-dns已被coredns取代


## 集群内域名解析原理
Kubernetes 集群节点上 kubelet 有--cluster-dns=${dns-service-ip} 和 --cluster-domain=${default-local-domain} 两个 dns 相关参数，分别被用来设置集群DNS服务器的IP地址和主域名后缀.

查看集群 default 命名空间下 dnsPolicy:ClusterFirst 模式的 Pod 内的 DNS 域名解析配置文件 /etc/resolv.conf 内容：
```conf
nameserver 172.24.0.10
search default.svc.cluster.local svc.cluster.local cluster.local
options ndots:5
```

各参数描述如下：
- nameserver: 定义 DNS 服务器的 IP 地址
- search: 设置域名的查找后缀规则，查找配置越多，说明域名解析查找匹配次数越多。集群匹配有 default.svc.cluster.local、svc.cluster.local、cluster.local 3个后缀，最多进行8次查询 (IPV4和IPV6查询各四次) 才能得到正确解析结果
- option: 定义域名解析配置文件选项，支持多个KV值. 例如该参数设置成ndots:5，说明如果访问的域名字符串内的点字符数量超过ndots值，则认为是完整域名，并被直接解析；如果不足ndots值，则追加search段后缀再进行查询

根据上述文件配置，在 Pod 内尝试解析：
- 同命名空间下服务，如 kubernetes：添加一次 search 域，发送kubernetes.default.svc.cluster.local. 一次 ipv4 域名解析请求到 172.24.0.10 进行解析即可
- 跨命名空间下的服务，如 coredns.kue-system：添加两次 search 域，发送 coredns.kue-system.default.svc.cluster.local. 和 coredns.kue-system.svc.cluster.local. 两次 ipv4 域名解析请求到 172.24.0.10 才能解析出正确结果
- 集群外服务，如 aliyun.com：添加三次 search 域，发送 aliyun.com.default.svc.cluster.local.、aliyun.com.svc.cluster.local.、aliyun.com.cluster.local. 和 aliyun.com 四次 ipv4 域名解析请求到 172.24.0.10 才能解析出正确结果

## Pod dnsPolicy
Kubernetes 集群中支持通过 dnsPolicy 字段为每个 Pod 配置不同的 DNS 策略. 目前支持四种策略：
- ClusterFirst：通过集群 DNS 服务来做域名解析，Pod 内 /etc/resolv.conf 配置的 DNS 服务地址是集群 DNS 服务的 kube-dns 地址。该策略是集群工作负载的默认策略
- None：忽略集群 DNS 策略，需要您提供 dnsConfig 字段来指定 DNS 配置信息
- Default：Pod 直接继承集群节点的域名解析配置, 即在集群直接使用节点的 /etc/resolv.conf 文件
- ClusterFirstWithHostNetwork：强制在 hostNetWork 网络模式下使用 ClusterFirst 策略（默认使用 Default 策略）

## CoreDNS
CoreDNS 目前是 Kubernetes 标准的服务发现组件，dnsPolicy: ClusterFirst 模式的 Pod 会使用 CoreDNS 来解析集群内外部域名.
在命名空间 kube-system 下，集群有一个名为 coredns 的 configmap(` kubectl get configmap coredns -n kube-system -o yaml`或`kubectl describe configmap coredns -n kube-system`). 其 Corefile 字段的文件配置内容如下（CoreDNS 功能都是通过 Corefile 内的插件提供）.
```yaml
  Corefile: |
    .:53 {
        errors
        health {
           lameduck 5s
        }
        ready
        kubernetes cluster.local in-addr.arpa ip6.arpa {
           pods insecure
           fallthrough in-addr.arpa ip6.arpa
           ttl 30
        }
        prometheus :9153
        forward . /etc/resolv.conf
        cache 30
        loop
        reload
        loadbalance
    }
```

其中，各插件说明：
- errors：错误信息到标准输出
- health：CoreDNS自身健康状态报告，默认监听端口8080，一般用来做健康检查. 可以通过http://localhost:8080/health获取健康状态
- ready：CoreDNS插件状态报告，默认监听端口8181，一般用来做可读性检查. 可以通过http://localhost:8181/ready获取可读状态. 当所有插件都运行后，ready状态为200
- kubernetes：CoreDNS kubernetes插件，提供集群内服务解析能力
- prometheus：CoreDNS自身metrics数据接口. 可以通过http://localhost:9153/metrics获取prometheus格式的监控数据
- forward（或proxy）：将域名查询请求转到预定义的DNS服务器. 默认配置中，当域名不在kubernetes域时，将请求转发到预定义的解析器（/etc/resolv.conf）中. 默认使用宿主机的/etc/resolv.conf配置
- cache：DNS缓存
- loop：环路检测，如果检测到环路，则停止CoreDNS
- reload：允许自动重新加载已更改的Corefile. 编辑ConfigMap配置后，请等待两分钟以使更改生效
- loadbalance：循环DNS负载均衡器，可以在答案中随机A、AAAA、MX记录的顺序