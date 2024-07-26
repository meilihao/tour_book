# k8s

## base
- [10 张图，说透 Kubernetes 架构原理，这会算是真懂了](https://my.oschina.net/morflameblog/blog/11572634)

## new
- [Kubernetes release](https://kubernetes.io/releases/)
- [Kubernetes dev roadmap](https://www.kubernetes.dev/resources/release/)

## 进阶
- [Kubernetes网络插件（CNI）基准测试的最新结果](https://tonybai.com/2019/04/18/benchmark-result-of-k8s-network-plugin-cni/)
- [Kubernetes网络插件（CNI）基准测试的最新结果202008](https://itnext.io/benchmark-results-of-kubernetes-network-plugins-cni-over-10gbit-s-network-updated-august-2020-6e1b757b9e49)
- [CNI 基准测试：Cilium 网络性能分析](https://kubesphere.io/zh/blogs/cilium-cni-benchmark/)
- [盘点Kubernetes网络问题的4种解决方案](https://blog.51cto.com/xjsunjie/2131176)
- [白话 Kubernetes Runtime](https://juejin.im/entry/5c8e5c28e51d4554ad53a1fc)
- [eBPF技术应用云原生网络实践系列之kubernetes网络](https://mlog.club/article/5493341)
- [为什么 Kubernetes 选择了 ETCD？](https://www.mgasch.com/2021/01/listwatch-part-1/)
- [virtual-kubelet 是一个开源的社区主导型项目，是Kubernetes kubelet的一种实现](https://luanlengli.github.io/2020/11/10/kubernetes%E5%9F%BA%E4%BA%8Evirtual-kubelet%E5%AE%9E%E7%8E%B0%E5%BC%B9%E6%80%A7Pod.html)

    它伪装成kubelet，与Kubernetes集群API通信, 实现Kubernetes API向阿里云的ECI、AWS的Fargate等serverless平台扩展.
- [OpenELB及其与 MetalLB 的对比](https://kubesphere.io/zh/blogs/openelb-joins-cncf-sandbox-project/)
- [kubernetes源码解析](https://www.cnblogs.com/lianngkyle/tag/kubernetes%E6%BA%90%E7%A0%81%E8%A7%A3%E6%9E%90/)

## ha
- [[k8s源码分析][kube-scheduler]scheduler之高可用及原理](https://www.jianshu.com/p/e30addc18560)
- [关于kube-controller-manager以及kube-scheduler的HA实现方式](https://w564791.gitbooks.io/kubernetes_gitbook/content/concept/leader-election.html)
- [谈谈k8s的leader选举--分布式资源锁](https://blog.csdn.net/weixin_39961559/article/details/81877056)
- [Kubernetes 源码剖析之 Leader 选举](https://wemp.app/posts/8ca1c89e-856e-4e37-bd20-5f34b43ddc40)

## dashboard
- [kuboard](https://kuboard.cn/install/v3/install-in-k8s.html)

## 部署工具
- kubekey
- sealos

## 存储
- [一文读懂 K8s 持久化存储流程](https://mp.weixin.qq.com/s?__biz=MzUzNzYxNjAzMg==&mid=2247490043&idx=1&sn=c09ad4a9bc790f4b742abd8ca1301ffb)
- [一文读懂容器存储接口 CSI](https://developer.aliyun.com/article/783464)
- [Kubernetes CSI(Container Storage Interface) 设计文档](http://anywhy.xyz/posts/2603569835/)
- [一图看懂CSI插件如何注册至Kubernetes](https://juejin.cn/post/7008041558997991461)
- [Kubernetes-CSI Documentation](https://github.com/kubernetes-csi/docs)
- [kubernetes ceph-csi分析](https://www.cnblogs.com/lianngkyle/p/14772131.html)
- [Kubernetes 1.23：树内存储向 CSI 卷迁移工作的进展更新](https://kubernetes.io/zh-cn/blog/2021/12/10/storage-in-tree-to-csi-migration-status-update/)

    在v1.27后, k8s所有csi in-tree代码都将被移除, 由out-of-tree代码取代.
- [container-storage-interface/spec](https://github.com/container-storage-interface/spec)