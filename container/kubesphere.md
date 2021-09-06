# kubesphere
参考:
- [KubeSphere 在华为云 ECS 高可用实例](https://kubesphere.io/zh/docs/installing-on-linux/public-cloud/install-kubesphere-on-huaweicloud-ecs/)
- [使用 Cilium 作为网络插件部署 Kubernetes + KubeSphere](https://kubesphere.io/zh/blogs/cilium-kubesphere-kubernetes/)
- [查询kubekey支持的k8s版本](https://hub.docker.com/r/kubesphere/kube-apiserver/tags?page=1&ordering=last_updated)
- [使用 KubeKey 升级](https://kubesphere.io/zh/docs/upgrade/upgrade-with-kubekey/)

## 部署
```bash
# all-in-one部署, kk v1.1.1, [k8s支持版本](https://github.com/kubesphere/kubekey/blob/release-1.1/docs/kubernetes-versions.md)
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
其实就是kube-apiserver未启动导致.

查看container debug log发现错误:
```log
Sep 04 10:28:04 chen-aliyun containerd[58865]: time="2021-09-04T10:28:04.520051030+08:00" level=error msg="RunPodSandbox for &PodSandboxMetadata{Name:kube-apiserver-chen-aliyun,Uid:a80ce4480dd54dda09001653c9dce4df,Namespace:kube-system,Attempt:0,} failed, error" error="failed to get sandbox image \"k8s.gcr.io/pause:3.2\": failed to pull image \"k8s.gcr.io/pause:3.2\": failed to pull and unpack image \"k8s.gcr.io/pause:3.2\": failed to resolve reference \"k8s.gcr.io/pause:3.2\": failed to do request: Head \"https://k8s.gcr.io/v2/pause/manifests/3.2\": dial tcp 108.177.97.82:443: i/o timeout"
```

检查containerd:
```bash
# containerd --version
containerd containerd.io 1.4.9 e25210fe30a0a703442421b0f60afac609f950a3
# cat /etc/containerd/config.toml |grep san
    sandbox_image = "kubesphere/pause:3.2"
# systemctl restart containerd.service
# containerd config dump |grep san
    sandbox_image = "k8s.gcr.io/pause:3.2"
# crictl images |grep pause
docker.io/kubesphere/pause                     3.2                 80d28bedfe5de       299kB
docker.io/kubesphere/pause                     3.5                 ed210e3e4a5ba       301kB
```

将containerd 升级到1.5.5, sandbox_image还是未生效, 但得到一个warning:
```bash
# containerd config dump |grep san
WARN[0000] deprecated version : `1`, please switch to version `2`
    sandbox_image = "k8s.gcr.io/pause:3.5"
```

将`version = 2`加到`/etc/containerd/config.toml`的开头, 看到配置生效, 再重启containerd即可.
```bash
# containerd config dump |grep san
    sandbox_image = "kubesphere/pause:3.5"
# systemctl restart containerd.service
```

### `kk create cluster`报`No kubeadm config, using etcd pod spec to get data directory`
其实是调用的`kubeadm reset`报的, 不是错误.

### `"Container runtime network not ready" networkReady="NetworkReady=false reason:NetworkPluginNotReady message:Network plugin returns error: cni plugin not initialized"`
其实这段输出并没有什么问题，等会安装flannel网络插件后，Network就会初始化了.

### [kubekey v1.2.0-alpha.3部署k8s v1.22.1失败](https://github.com/kubesphere/kubekey/issues/640)
```log
Warning: spec.template.metadata.annotations[scheduler.alpha.kubernetes.io/critical-pod]: non-functional in v1.16+; use the "priorityClassName" field instead
daemonset.apps/cilium created
deployment.apps/cilium-operator created
ERRO[11:44:41 CST] Failed to deploy local-volume.yaml: Failed to exec command: sudo -E /bin/sh -c "/usr/local/bin/kubectl apply -f /etc/kubernetes/addons/local-volume.yaml" storageclass.storage.k8s.io/local unchanged
serviceaccount/openebs-maya-operator unchanged
deployment.apps/openebs-localpv-provisioner unchanged
unable to recognize "/etc/kubernetes/addons/local-volume.yaml": no matches for kind "ClusterRole" in version "rbac.authorization.k8s.io/v1beta1"
unable to recognize "/etc/kubernetes/addons/local-volume.yaml": no matches for kind "ClusterRoleBinding" in version "rbac.authorization.k8s.io/v1beta1": Process exited with status 1  node=192.168.0.43
```

应该是`rbac.authorization.k8s.io/v1beta1 ClusterRole/ClusterRoleBinding`已在v1.22.1删除导致.

切换到kk 1.1.1部署k8s v1.20.6
```log
Warning: rbac.authorization.k8s.io/v1beta1 ClusterRole is deprecated in v1.17+, unavailable in v1.22+; use rbac.authorization.k8s.io/v1 ClusterRole
clusterrole.rbac.authorization.k8s.io/openebs-maya-operator created
Warning: rbac.authorization.k8s.io/v1beta1 ClusterRoleBinding is deprecated in v1.17+, unavailable in v1.22+; use rbac.authorization.k8s.io/v1 ClusterRoleBinding
clusterrolebinding.rbac.authorization.k8s.io/openebs-maya-operator created
```