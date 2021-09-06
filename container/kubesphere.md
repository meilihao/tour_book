# kubesphere
参考:
- [KubeSphere 在华为云 ECS 高可用实例](https://kubesphere.io/zh/docs/installing-on-linux/public-cloud/install-kubesphere-on-huaweicloud-ecs/)
- [使用 Cilium 作为网络插件部署 Kubernetes + KubeSphere](https://kubesphere.io/zh/blogs/cilium-kubesphere-kubernetes/)

## 部署
```bash
# all-in-one部署, kk v1.1.1, [k8s支持版本](https://github.com/kubesphere/kubekey/blob/release-1.1/docs/kubernetes-versions.md)
$ sudo apt install socat ipset conntrack ipvsadm
$ sudo KKZONE=cn kk create config --with-kubernetes v1.20.6 --with-kubesphere v3.1.1 -f my.yaml
$ cat my.yaml # 使用`kk create cluster -f`部署时的配置文件
apiVersion: kubekey.kubesphere.io/v1alpha1
kind: Cluster
metadata:
  name: sample
spec:
  hosts:
  - {name: chen-aliyun, address: 192.168.0.43, internalAddress: 192.168.0.43, user: root, password: xxx}
  roleGroups:
    etcd:
    - chen-aliyun
    master: 
    - chen-aliyun
    worker:
    - chen-aliyun
  controlPlaneEndpoint:
    domain: lb.kubesphere.local
    address: ""
    port: 6443
  kubernetes:
    version: v1.20.6
    imageRepo: kubesphere
    clusterName: cluster.local
    containerManager: containerd
  network:
    plugin: cilium
    kubePodsCIDR: 10.233.64.0/18
    kubeServiceCIDR: 10.233.0.0/18
  registry:
    registryMirrors: []
    insecureRegistries: []
  addons: []


---
apiVersion: installer.kubesphere.io/v1alpha1
kind: ClusterConfiguration
metadata:
  name: ks-installer
  namespace: kubesphere-system
  labels:
    version: v3.1.1
spec:
  persistence:
    storageClass: ""       
  authentication:
    jwtSecret: ""
  zone: ""
  local_registry: ""        
  etcd:
    monitoring: false      
    endpointIps: localhost  
    port: 2379             
    tlsEnable: true
  common:
    redis:
      enabled: false
    redisVolumSize: 2Gi 
    openldap:
      enabled: false
    openldapVolumeSize: 2Gi  
    minioVolumeSize: 20Gi
    monitoring:
      endpoint: http://prometheus-operated.kubesphere-monitoring-system.svc:9090
    es:  
      elasticsearchMasterVolumeSize: 4Gi   
      elasticsearchDataVolumeSize: 20Gi   
      logMaxAge: 7          
      elkPrefix: logstash
      basicAuth:
        enabled: false
        username: ""
        password: ""
      externalElasticsearchUrl: ""
      externalElasticsearchPort: ""  
  console:
    enableMultiLogin: true 
    port: 30880
  alerting:       
    enabled: false
    # thanosruler:
    #   replicas: 1
    #   resources: {}
  auditing:    
    enabled: false
  devops:           
    enabled: false
    jenkinsMemoryLim: 2Gi     
    jenkinsMemoryReq: 1500Mi 
    jenkinsVolumeSize: 8Gi   
    jenkinsJavaOpts_Xms: 512m  
    jenkinsJavaOpts_Xmx: 512m
    jenkinsJavaOpts_MaxRAM: 2g
  events:          
    enabled: false
    ruler:
      enabled: true
      replicas: 2
  logging:         
    enabled: false
    logsidecar:
      enabled: true
      replicas: 2
  metrics_server:             
    enabled: false
  monitoring:
    storageClass: ""
    prometheusMemoryRequest: 400Mi  
    prometheusVolumeSize: 20Gi  
  multicluster:
    clusterRole: none 
  network:
    networkpolicy:
      enabled: false
    ippool:
      type: none
    topology:
      type: none
  openpitrix:
    store:
      enabled: false
  servicemesh:    
    enabled: false  
  kubeedge:
    enabled: false
    cloudCore:
      nodeSelector: {"node-role.kubernetes.io/worker": ""}
      tolerations: []
      cloudhubPort: "10000"
      cloudhubQuicPort: "10001"
      cloudhubHttpsPort: "10002"
      cloudstreamPort: "10003"
      tunnelPort: "10004"
      cloudHub:
        advertiseAddress: 
          - ""           
        nodeLimit: "100"
      service:
        cloudhubNodePort: "30000"
        cloudhubQuicNodePort: "30001"
        cloudhubHttpsNodePort: "30002"
        cloudstreamNodePort: "30003"
        tunnelNodePort: "30004"
    edgeWatcher:
      nodeSelector: {"node-role.kubernetes.io/worker": ""}
      tolerations: []
      edgeWatcherAgent:
        nodeSelector: {"node-role.kubernetes.io/worker": ""}
        tolerations: []
$ sudo kk create cluster -f config.yaml
```

根据[network plugin性能测试的结果](https://itnext.io/benchmark-results-of-kubernetes-network-plugins-cni-over-10gbit-s-network-updated-august-2020-6e1b757b9e49), 目前calico是最佳选择, 修改`plugin: calico`即可.

> 实践发现, 即使使用calico, cilium还是会被安装, 应该有其他kubesphere组件依赖了cilium.

## kk
```bash
$ export KKZONE=cn # 因为无法访问 https://storage.googleapis.com
$ kk create config --with-kubernetes v1.18.6 --with-kubesphere v3.0.0 -f mycluster.yaml # 创建cluster配置kubesphere
$ vim mycluster.yaml # 自定义配置
$ kk create cluster -f mycluster.yaml # 根据yaml配置部署kubesphere
$ kk create config --from-cluster -f my.yaml --kubeconfig=/etc/kubernetes/admin.conf # 根据已有的k8s/kubesphere集群创建配置(该配置不含安装kubesphere的信息), 因此需要结合`kk create config`输出的配置merge `installer.kubesphere.io/v1alpha`(安装kubesphere的配置)才可用
```

## FAQ
### kubectl报`The connection to the server localhost:8080 was refused - did you specify the right host or port?`
在all-in-one部署成功后, `Step 4: 验证安装结果`时kubectl报错.

原因: kubectl命令需要使用kubernetes-admin来运行, all-in-one部署时使用了sudo, 普通用户执行kubectl时没有所需的KUBECONFIG, 将主节点的`/etc/kubernetes/admin.conf`(这里仅有一个节点即当前节点)让用户可读即可.

修复:
```bash
$ sudo ls -l /etc/kubernetes/admin.conf # 检查权限
-rw------- 1 root root 5467 Dec  1 22:03 /etc/kubernetes/admin.conf
$ sudo chmod +r /etc/kubernetes/admin.conf # 允许其他用户可读
$ echo "export KUBECONFIG=/etc/kubernetes/admin.conf" >> ~/.bash_profile
$ source ~/.bash_profile
```

### 组件异常, 事件里提示"0/1 nodes are available: 1 Insufficient cpu."
参考:
- [浅析Kubernetes资源管理](https://www.infoq.cn/article/rrsrvv093hh6f1ymkcez)

可通过`kubectl get events --namespace=kubesphere-monitoring-system  |grep "prometheus-k8s-0"`或`kubectl describe pod prometheus-k8s-0 -n kubesphere-monitoring-system`查看对应的event.

原因: cpu不足无法调度, 可通过`kubectl describe nodes`查看该node的资源限制. 目前除了添加资源外无解(除非手动修改deployment/pod配置的resources.requests).

<<<<<<< HEAD
### `crictl pull kubesphere/pause:3.2`报`Unimplemented desc = unknown service runtime.v1alpha2.ImageService`
参考:
- [容器运行时笔记](https://gobomb.github.io/post/container-runtime-note/)
- [unknown service runtime.v1alpha2.ImageService](https://github.com/kubernetes-sigs/cri-tools/issues/710)

因为 contianerd 没有 启用 CRI 插件，所以无法使用 crictl 连接.

解决方法: 注释`/etc/containerd/config.toml`里的`disabled_plugins = ["cri"]`并`systemctl restart containerd`.

### `kk create cluster`报`Unable to fetch the kubeadm-config ConfigMap from cluster: failed to get config map: Get "https://lb.kubesphere.local:6443/api/v1/namespaces/kube-system/configmaps/kubeadm-config?timeout=10s": dial tcp 192.168.0.43:6443: connect: connection refused`
其实是调用的`kubeadm reset`报的, 不是错误. `sudo rm -rf /etc/kubernetes*`后重新执行`kk create cluster`就不报该提示了.

### `kk create cluster`报`No kubeadm config, using etcd pod spec to get data directory`
其实是调用的`kubeadm reset`报的, 不是错误.

### `"Container runtime network not ready" networkReady="NetworkReady=false reason:NetworkPluginNotReady message:Network plugin returns error: cni plugin not initialized"`
其实这段输出并没有什么问题，等会安装flannel网络插件后，Network就会初始化了.

### `Failed to initialize CSINode: error updating CSINode annotation: timed out waiting for the condition; caused by: nodes "chen-aliyun" not found`
因为node "chen-aliyun"还没由注册, 后续会存在注册操作的日志比如`Successfully registered node chen-aliyun`.

### KubeSphere安装报`KubeSphere startup timeout`, systemd日志里由很多`Container runtime network not ready: NetworkReady=false reason:NetworkPluginNotReady message:Network plugin returns error: cni plugin not initialized`
参考:
- [KubeSphere startup timeout](https://github.com/kubesphere/kubekey/issues/654)

报`cni plugin not initialized`的原因是`etc/cni/net.d` 中没有定义 CNI 网络, 必须将配置文件写入该目录以告诉 CNI 驱动程序如何配置网络. 通过`kubectl describe node`查看是k8s
网络插件没安装成功导致node not ready.

```bash
# kubectl get pod -A
NAMESPACE     NAME                                  READY   STATUS    RESTARTS   AGE
kube-system   cilium-htv9g                          1/1     Running   1          5h30m
kube-system   cilium-operator-75f898cccc-64fpl      1/1     Running   1          5h30m
kube-system   cilium-operator-75f898cccc-cpkpp      1/1     Running   1          5h30m
...
```
查看cilium的pod, 显然已运行.

问题复现(删除并再创建cluster):
```bash
# kk create cluster -f config.k8s.yaml # onyly install k8s, 且部署成功
# kk delete cluster -f config.k8s.yaml
# kk create cluster -f config.k8s.yaml
```

解决方法: `kk delete cluster -f config.k8s.yaml`并**reboot**, 再执行`kk create cluster -f config.k8s.yaml`.

### kubesphere部署完成第一次登录修改初始密码或稍后修改都报错:`Internal error occurred: failed calling webhook "users.iam.kubesphere.io": Post "https://ks-controller-manager.kubesphere-system.svc:443/validate-email-iam-kubesphere-io-v1alpha2?timeout=30s": dial tcp 10.233.29.13:443: connect: connection refused`
参考:
- [帐户无法登录](https://kubesphere.io/zh/docs/faq/access-control/cannot-login/)

经排查是pod ks-controller-manager-f457f6957-2ps85 Pending导致, 通过修改kubelet的保留计算资源启动它.