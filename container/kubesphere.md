# kubesphere
参考:
- [KubeSphere 在华为云 ECS 高可用实例](https://kubesphere.io/zh/docs/installing-on-linux/public-cloud/install-kubesphere-on-huaweicloud-ecs/)
- [使用 Cilium 作为网络插件部署 Kubernetes + KubeSphere](https://kubesphere.io/zh/blogs/cilium-kubesphere-kubernetes/)
- [查询kubekey支持的k8s版本](https://hub.docker.com/r/kubesphere/kube-apiserver/tags?page=1&ordering=last_updated)
- [使用 KubeKey 升级](https://kubesphere.io/zh/docs/upgrade/upgrade-with-kubekey/)

## 部署
```bash
# all-in-one部署, kk v2.1.0, [k8s支持版本](https://github.com/kubesphere/kubekey/blob/v2.1.0/docs/kubernetes-versions.md)
$ sudo apt install socat ipset conntrack ipvsadm
$ sudo KKZONE=cn kk create config --container-manager containerd --with-kubernetes v1.24.0 --with-kubesphere v3.2.1 -f my.yaml
$ cat my.yaml # 使用`kk create cluster -f`部署时的配置文件
apiVersion: kubekey.kubesphere.io/v1alpha2
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
    ## Internal loadbalancer for apiservers 
    # internalLoadbalancer: haproxy

    domain: lb.kubesphere.local
    address: ""
    port: 6443
  kubernetes:
    version: v1.24.0
    clusterName: cluster.local
    autoRenewCerts: true
  etcd:
    type: kubekey
  network:
    plugin: cilium
    kubePodsCIDR: 10.233.64.0/18
    kubeServiceCIDR: 10.233.0.0/18
    ## multus support. https://github.com/k8snetworkplumbingwg/multus-cni
    multusCNI:
      enabled: false
  registry:
    privateRegistry: ""
    namespaceOverride: ""
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
    version: v3.2.1
spec:
  persistence:
    storageClass: ""
  authentication:
    jwtSecret: ""
  zone: ""
  local_registry: ""
  namespace_override: ""
  # dev_tag: ""
  etcd:
    monitoring: true
    endpointIps: localhost
    port: 2379
    tlsEnable: true
  common:
    core:
      console:
        enableMultiLogin: true
        port: 30880
        type: NodePort
    # apiserver:
    #  resources: {}
    # controllerManager:
    #  resources: {}
    redis:
      enabled: false
      volumeSize: 2Gi
    openldap:
      enabled: false
      volumeSize: 2Gi
    minio:
      volumeSize: 20Gi
    monitoring:
      # type: external
      endpoint: http://prometheus-operated.kubesphere-monitoring-system.svc:9090
      GPUMonitoring:
        enabled: false
    gpu:
      kinds:         
      - resourceName: "nvidia.com/gpu"
        resourceType: "GPU"
        default: true
    es:
      # master:
      #   volumeSize: 4Gi
      #   replicas: 1
      #   resources: {}
      # data:
      #   volumeSize: 20Gi
      #   replicas: 1
      #   resources: {}
      logMaxAge: 7
      elkPrefix: logstash
      basicAuth:
        enabled: false
        username: ""
        password: ""
      externalElasticsearchHost: ""
      externalElasticsearchPort: ""
  alerting:
    enabled: false
    # thanosruler:
    #   replicas: 1
    #   resources: {}
  auditing:
    enabled: false
    # operator:
    #   resources: {}
    # webhook:
    #   resources: {}
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
    # operator:
    #   resources: {}
    # exporter:
    #   resources: {}
    # ruler:
    #   enabled: true
    #   replicas: 2
    #   resources: {}
  logging:
    enabled: false
    containerruntime: docker
    logsidecar:
      enabled: true
      replicas: 2
      # resources: {}
  metrics_server:
    enabled: false
  monitoring:
    storageClass: ""
    # kube_rbac_proxy:
    #   resources: {}
    # kube_state_metrics:
    #   resources: {}
    # prometheus:
    #   replicas: 1
    #   volumeSize: 20Gi
    #   resources: {}
    #   operator:
    #     resources: {}
    #   adapter:
    #     resources: {}
    # node_exporter:
    #   resources: {}
    # alertmanager:
    #   replicas: 1
    #   resources: {}
    # notification_manager:
    #   resources: {}
    #   operator:
    #     resources: {}
    #   proxy:
    #     resources: {}
    gpu:
      nvidia_dcgm_exporter:
        enabled: false
        # resources: {}
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
参考:
- [fix containerd still uses "k8s.gcr.io/pause"](https://github.com/kubesphere/website/pull/1924)

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

经排查是pod ks-controller-manager-f457f6957-2ps85 Pending导致, [通过修改kubelet的配置减少它的部分保留计算资源来启动该pod](k8s.md).