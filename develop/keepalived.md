# keepalived
linux下轻量级ha解决方案, 主要通过虚拟路由冗余来实现ha. 与heartbeat相比, 没有其功能强大, 但在某些场景下更简单合适.

keepalived起初是为LVS设计, 专门监控(根据3~5层的交换机制)集群中各节点的状态, 通过剔除故障节点和加入正常节点来实现ha.

keepalived支持VRRP(Virtual Router Redundancy Protocol, 虚拟路由器冗余协议, 为了解决静态路由器出现的单点故障).