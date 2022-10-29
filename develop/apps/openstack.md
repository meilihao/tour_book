# openstack

## cinder
参考:
- [<<OpenStack架构分析与实践>>]
- [<<OpenStack设计与实现(第2版)>>]

cinder负责openstack的块存储, 它本身不是一种存储技术, 而是提供一个中间抽象层, 然后通过调用不同存储后端类型的驱动接口来管理相对应的后端存储.

cinder组件:
- cinder-api : 一个WSGI应用, 将外部的请求通过rabbitmq路由到cinder中的相关服务, 然后响应相应的请求.

    入口: `cinder/cmd/api:main()`
- cinder-scheduler : 负责筛选出合适的cinder-volume进行块设备的创建. 它默认的调度算法有按容量, 按AZ调度, 按volume类型, 按用户自定义的属性等.
- cinder-volume : 管理cinder中的后端块存储设备, 由其中的volume manager对接不同的存储后端驱动.
- cinder-backup : 将块存储中的volume备份到其他地方, 比如openstack的对象存储(swift), google cloud storage等.

> consistency group即对volume进行逻辑分组, 从而实现对组内成员进行统一操作, 比如打快照.

> cinder配置在`/etc/cinder`中.

cinder驱动状态:
- [Available Drivers](https://docs.openstack.org/cinder/rocky/drivers.html)
- [Cinder Driver Support Matrix](https://docs.openstack.org/cinder/rocky/reference/support-matrix.html)

## Neutron
ref:
- [LBaaS 实现机制 - 每天5分钟玩转 OpenStack（125）](https://developer.aliyun.com/article/463372)

Neutron 是用 Haproxy 来实现负责LBaaS(负载均衡)的.