# k8s

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

## ha
- [[k8s源码分析][kube-scheduler]scheduler之高可用及原理](https://www.jianshu.com/p/e30addc18560)
- [关于kube-controller-manager以及kube-scheduler的HA实现方式](https://w564791.gitbooks.io/kubernetes_gitbook/content/concept/leader-election.html)
- [谈谈k8s的leader选举--分布式资源锁](https://blog.csdn.net/weixin_39961559/article/details/81877056)
- [Kubernetes 源码剖析之 Leader 选举](https://wemp.app/posts/8ca1c89e-856e-4e37-bd20-5f34b43ddc40)