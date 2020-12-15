# kubesphere

## 部署
```bash
# all-in-one部署, kk v1.0.1
$ sudo KKZONE=cn kk create cluster --with-kubernetes v1.18.6 --with-kubesphere v3.0.0
$ cat my.yaml # 使用`kk create cluster -f`部署时的配置文件
apiVersion: kubekey.kubesphere.io/v1alpha1
kind: Cluster
metadata:
  name: mytest
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
    version: v1.18.6
    imageRepo: kubesphere
    clusterName: cluster.local
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
    version: v3.0.0
spec:
  local_registry: ""
  persistence:
    storageClass: ""
  authentication:
    jwtSecret: ""
  etcd:
    monitoring: true
    endpointIps: localhost
    port: 2379
    tlsEnable: true
  common:
    es:
      elasticsearchDataVolumeSize: 20Gi
      elasticsearchMasterVolumeSize: 4Gi
      elkPrefix: logstash
      logMaxAge: 7
    mysqlVolumeSize: 20Gi
    minioVolumeSize: 20Gi
    etcdVolumeSize: 20Gi
    openldapVolumeSize: 2Gi
    redisVolumSize: 2Gi
  console:
    enableMultiLogin: false  # enable/disable multi login
    port: 30880
  alerting:
    enabled: false
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
    logsidecarReplicas: 2
  metrics_server:
    enabled: true
  monitoring:
    prometheusMemoryRequest: 400Mi
    prometheusVolumeSize: 20Gi
  multicluster:
    clusterRole: none  # host | member | none
  networkpolicy:
    enabled: false
  notification:
    enabled: false
  openpitrix:
    enabled: false
  servicemesh:
    enabled: false
```

根据[network plugin性能测试的结果](https://itnext.io/benchmark-results-of-kubernetes-network-plugins-cni-over-10gbit-s-network-updated-august-2020-6e1b757b9e49), 目前calico是最佳选择, 修改`plugin: calico`即可.

> 实践发现, 即使使用calico, cilium还是会被安装, 应该有其他kubesphere组件依赖了cilium.

## kk
```bash
$ kk create config --with-kubernetes v1.18.6 --with-kubesphere v3.0.0 -f mycluster.yaml # 创建cluster配置kubesphere
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