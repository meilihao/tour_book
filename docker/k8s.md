# k8s

## doc
- [Kubernetes指南](https://feisky.gitbooks.io/kubernetes/)
- [Kubernetes中文手册](https://www.kubernetes.org.cn/docs)

## install
```sh
apt-get update && apt-get install -y apt-transport-https
curl https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | apt-key add - 
cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb http://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main
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

## cmd
- `kubectl logs -f POD-NAME` # 获取pod日志
- `kubectl exec -it POD-NAME sh` # 进入pod的容器
- `kubectl describe pods/POD-NAME` # 获取pod的描述信息(简单), 生命周期事件
- `kubectl get pod myweb-fnncj --output json/yaml` # 获取pod的详细信息, 有状态信息
- `kubectl get event` # 查询所有事件
- `kubectl logs POD-NAME Container-NAME [-p]` # 查询pod中容器的日志,`-p`允许查询`Container-NAME`重建前的日志