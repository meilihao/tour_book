# vault

## ha
参考:
- [High Availability Mode (HA)](https://www.vaultproject.io/docs/concepts/ha)

ha基于分布式锁.

如果vault的数据存储后端支持ha, 那么vault自动启用HA mode.

> 如果vault的数据存储后端实现了[physical.HABackend 接口](https://github.com/hashicorp/vault/blob/master/sdk/physical/physical.go#L51), 那么即可认为它支持ha, 比如[etcd](https://github.com/hashicorp/vault/blob/master/physical/etcd/etcd3.go#L42).

见[cluster_test.go#TestCluster_ForwardRequests](https://github.com/hashicorp/vault/blob/master/vault/cluster_test.go#L171).