# ha
参考:
- [如何实现靠谱的分布式锁？（附SharkLock的设计选择)](https://cloud.tencent.com/developer/article/1348743)

常用的方法是基于分布式存储(比如etcd)实现的分布式锁. 比如[vault](/develop/apps/vault.md)和k8s的ha.