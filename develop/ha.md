# cluster
cluster是一组协同工作的服务集合, 用来提供比单一服务更稳定, 高效及扩展性的服务平台. 它通常由两个及以上的server组成.

集群特定:
- 高可用性

    出现服务故障时, 自动从故障节点切换到备用节点(常用做法是漂移vip), 从而提供不间断的服务, 保障业务的持续性.
- 可扩展性

    动态加入节点, 增强cluster的整体性能

常见ha(high availability cluster)的形式: 双机热备(active/standby), 双机互备, 多机互备等.
常见lb(load balance cluster): haproxy ,lvs集群, F5 Networks等.
常见分布式计算集群: hadoop, spark.

# ha(High Availability)
参考:
- [pacemaker+corosync/heartbeat对比及资源代理RA脚本](https://www.cnblogs.com/clsblog/p/6202869.html)
- [<<DRBD权威指南——基于Corosync+Heartbeat技术构建网络RAID>>]
- [SUSE Linux Enterprise High Availability Extension](https://www.novell.com/zh-cn/documentation/sle_ha/book_sleha/)
- [中标麒麟高可用集群软件V7.0产品白皮书](http://www.kylinos.cn/support/document/34.html)

[驱动、开发者和Linux厂商，以及整个开源高可用集群社区，都已经转移到了基于Corosync 2.x+Pacemaker的HA堆栈上](http://www.linux-ha.org/wiki/Site_news), [Heartbeat](http://www.linux-ha.org)已名存实亡.

## ha概念
- node

    运行心跳进程的一个独立主机(host), 即为节点(node).
- resource

    一个node可控制的实体, 当该节点发生故障时, 这些资源可被其他node接管, 常见的资源有:
    - 磁盘, 文件系统
    - ip
    - 应用服务
- event

    cluster中可能发生的事
- action

    event发生时ha的响应行为, 通常使用脚本来处理.

## example
- [ZFS High-Availability NAS](https://github.com/ewwhite/zfs-ha/wiki)

## corosync
coreosync在传递信息的时候可以通过一个简单的配置文件来定义信息传递的方式和协议等. 因此corosync是用于高可用环境中的提供通讯服务的，它位于高可用集群架构中的底层（Message Layer），扮演着为各节点（node）之间提供心跳信息传递这样的一个角色.

log位置: `/etc/corosync/corosync.conf`中的`logging.logfile`

配置文件`/etc/corosync/corosync.conf`:
```conf
totem {                      //节点间的通信协议，主要定义通信方式，通信协议版本，加密算法
        version: 2           //协议版本
        crypto_cipher: aes256     //加密密码类型
        crypto_hash: sha1         //加密为sha1
        interface {      //定义集群心跳信息及传递的接口中，可以有多组
                ringnumber: 0       //环数量，如果一个主机有多块网卡，避免心跳信息环发送
                bindnetaddr: 192.168.10.0   //绑定的网段（本机网段为192.168.10.0/24）
                mcastaddr: 239.255.10.9     //组播地址
                mcastport: 5405             //组播端口
                ttl: 1                      //数据包的ttl值，由于不跨三层设备，这里用默认的1
        }
}
logging {                   //跟日志相关配置
        fileline: off        //是否记录fileline
        to_stderr: no          //表示是否需要发送到错误输出
        to_logfile: yes         //是否记录在日志文件中
        logfile: /var/log/cluster/corosync.log    //日志文件目录
        to_syslog: no   //是否将日志发往系统日志
        debug: off     //是否输出debug日志
        timestamp: on   //是否打开时间戳
        logger_subsys {
                subsys: QUORUM
                debug: off
        }
}
quorum {          //投票系统
        provider: corosync_votequorum  //支持哪种投票方式           
                  expected_votes：2  //总投票数
　　　　　　　　　　　two_nodes: 1  //是否为2节点集群（两节点特殊} 

}

nodelist {       //节点列表
　　　　node {
　　　　　　ring0_addr: Node1.contoso.com   //节点名称：主机名或IP
　　　　　　nodeid: 1                     //节点编号
　　　　　}
　　　node {
　　　　　ring0_addr: Node2.contoso.com
　　　　　nodeid: 2
　　　　}
}
```

## pacemaker
参考:
- [Pacemaker Administration](https://clusterlabs.org/pacemaker/doc/en-US/Pacemaker/2.0/html-single/Pacemaker_Administration/index.html)
- [从头开始搭建集群 在Fedora上面创建主/主和主备集群](https://clusterlabs.org/pacemaker/doc/deprecated/zh-CN/Pacemaker/1.1-plugin/html-single/Clusters_from_Scratch/index.html)
- [Pacemaker详解](https://www.cnblogs.com/chimeiwangliang/p/7975911.html)

pacemaker是一个开源的高可用集群资源管理器(CRM)，位于HA集群架构中资源管理、资源代理(RA)这个层次，它不能提供底层心跳信息传递的功能，要想与对方节点通信需要借助底层的心跳传递服务，将信息通告给对方. 通常它与corosync的结合方式有两种:
- pacemaker作为corosync的插件运行
- pacemaker作为独立的守护进程运行

pacemaker负责仲裁指定谁是活动节点、IP地址的转移、本地资源管理系统.

当故障节点修复后，资源返回来称为failback，当故障节点修复后，资源仍在备用节点，称为failover.

[pacemaker组建重命名(1.x->2.x)](https://wiki.clusterlabs.org/wiki/Pacemaker_2.0_Daemon_Changes):
- attrd 	pacemaker-attrd

    node attribute manager
- cib	pacemaker-based

    群集信息库管理者,是一个xml文件, 描述了配置，节点，资源，状态, 限制条件等集群信息.

    集群中只有一个主 CIB 文件,是通过 DC 来维护的. 所有节点的 CIB 文件都是主 CIB 文件的副本. 如果需要修改集权配置,就必须通过 DC 来对主 CIB 文件进
行修改.

    > cibadmin -Q : 可查看cib信息
- crmd 	pacemaker-controld

    集群资源管理守护进程. 主要是消息代理的PEngine和LRM，还选举一个领导者DC, 它会统筹集群的操作, 比如启动/停止/隔离/转移资源等.
- lrmd 	pacemaker-execd

    本地资源管理守护进程(local resource agent executor). 它提供了一个通用的接口, 支持直接调用资源代理（脚本）来管理资源.
- stonithd 	pacemaker-fenced

    STONITH(Shoot the Other Node in the Head), 强制使节点下线, 以防数据被恶意节点或并行访问破坏.

    因为如果某个节点没有响应，但并不代表它没有在提供服务，100%确保数据安全的**唯一**做法就是在允许另外一个节点操作数据之前，使用STONITH来隔离该节点, 保证其真的下线了.

    所有 STONITH 资源默认存放在每个节点的/usr/lib/stonith/plugins 目录下, 通过这些 STONITH 资源可是实现节点关机、节点重启等功能.

    使用`stonith -L`可查看cluster支持的STONITH.
- pacemaker_remoted 	pacemaker-remoted

    remote resource agent executor
- pengine 	pacemaker-schedulerd
    
    策略引擎, action scheduler

    主要负责将CRM发过来的一些信息按照配置文件中的各种设置（基于目前的状态和配置）计算集群的下一个状态.

    当主控节点(DC)需要对集群做出更改时, 策略引擎需根据配置文件中的各种设置及当前状态来计算下一状态中群集和列表(资源)需要采取的行动. DC 发送消息给有关联的 CRM,然后 CRM 调用本地资源管理(LRM)完成对资源的修改.

![Pacemaker Architecture 2.x](https://clusterlabs.org/pacemaker/doc/en-US/Pacemaker/2.0/html-single/Pacemaker_Administration/images/pcmk-internals.png)

pacemaker工作原理:
CIB使用XML表示集群的集群中的所有资源的配置和当前状态. CIB的内容会被自动在整个集群中同步，使用PEngine计算集群的理想状态，生成指令列表，然后输送到DC（指定协调员）. Pacemaker 集群中所有节点选举的DC节点作为主决策节点, 如果当选DC节点宕机，它会在所有的节点上， 迅速选出一个新的DC. DC将PEngine生成的策略，传递给其他节点上的LRMd（本地资源管理守护程序）或CRMD通过集群消息传递基础结构. 当集群中有节点宕机，PEngine重新计算的理想策略. 在某些情况下，可能有必要关闭节点，以保护共享数据或完整的资源回收. 为此，Pacemaker配备了stonithd设备. STONITH可以将其它节点`爆头`，通常是实现与远程电源开关. Pacemaker会将STONITH设备，配置为资源保存在CIB中，使他们可以更容易地监测资源失败或宕机.

相关概念:
- 资源粘性：资源粘性表示资源是否倾向于留在当前节点，如果为正整数，表示倾向，负数表示移离，-inf表示正无穷，inf表示正无穷.

    资源黏性值范围及其作用：
    - 0：默认选项. 资源放置在系统中的最适合位置。这意味着当负载能力“较好”或较差的节点变得可用时才转移资源
    此选项的作用基本等同于自动故障恢复，只是资源可能会转移到非之前活动的节点上
    - 大于0：资源更愿意留在当前位置，但是如果有更合适的节点可用时会移动。值越高表示资源越愿意留在当前位置
    - 小于0：资源更愿意移离当前位置。绝对值越高表示资源越愿意离开当前位置
    - INFINITY：如果不是因节点不适合运行资源（节点关机、节点待机、达到migration-threshold 或配置更改）而强制资源转移，资源总是留在当前位置. 此选项的作用几乎等同于完全禁用自动故障回复
    - -INFINITY：资源总是移离当前位置

    当某个高可用集群即包含资源粘性又包含位置约束，一旦该节点发生故障后，资源就会转移到另一个节点上去, 但是当之前的节点恢复正常时，需要比较所有的资源粘性之和与所有位置约束之和谁大谁小，这样资源才会留在大的一方

- 资源类型：

    资源是集群管理的最小单位对象.
    根据用户的配置，资源有不同的种类，其中最为简单的资源是原始资源(primitive Resource),此外还有相对高级和复杂的资源组(Resource Group)和克隆资源(Clone Resource)等集群资源概念:
    - primitive（native）：基本资源，原始资源
    - group : 捆绑约束将不同的资源捆绑在一起作为一个逻辑整体来调度

        特点:
        - 资源按照其指定的先后顺序启动，还会按照其指定顺序的相反顺序停止.
        - 如果资源组中的某个资源无法在任何节点启动运行，那么在该资源后指定的任何资源都将无法运行
        - 资源组中后追加的资源不影响原有group资源的运行

        资源组具有组属性，并且资源组会继承组成员的部分属性，主要被继承的资源属性包括 Priority、Targct-role、Is-managed等，资源属性决定了资源在集群中的行为规范，以及资源管理器可以对其进行哪些操作:
        - Priority:资源优先级，其默认值是0，如果集群无法保证所有资源都处于运行状态，则低优先权资源会被停止，以便让高优先权资源保持运行状态
        - Target-role：资源目标角色，其默认值是started'表示集群应该让这个资源处于何种状态，允许值为:

　　　　　　- Stopped：表示强制资源停止；
　　　　　　- Started:表示允许资源启动，但是在多状态资源的情况下不能将其提升为 Master资源
　　　　　　- Master:允许资源启动，并在适当时将其提升为 Master

        - is-managed：其默认值是true，表示是否允许集群启动和停止该资源，false表示不允许
        - Resource-stickiness：默认值是0，表示该资源保留在原有位置节点的倾向程度值
        - Requires：默认值为 fencing，表示资源在什么条件下允许启动

        相关命令:
        - `pcs resource group add group_name resource_id  ... [resource_id] [--before resource_id] --after resource_id`
        - `pcs resource create resource_id Standard:Provider:Type丨 type [ resource_options] [op operation_action operation_options] --group group_name`
        - `pcs resource group remove group_name resource_id ...`
        - `pcs resource group list`
    - clone : 将primitive/group资源克隆n份且同时运行在n个节点上

        语法:
        - `pcs resource create resource_id standard:provider: type| type [resource options] --clone[meta clone_options]`
        - `pcs resource clone resource_id group_name [clone_optione] ...`

        clone_options:
        - Priority/Target-role/ls-manage：这三个克隆资源属性是从被克隆的资源中继承而来的，具体意义可以参考被cloned资源属性的解释
        - Clone-max：该选项值表示需要存在多少资源副本才能启动资源，默认为该集群中的节点数, 即克隆资源时默认会在集群中的每个在线节点上都存在一个副本
        - Clone-node-max：表示在单一节点上能够启动多少个资源副本，默认值为1
        - Notify：表示在停止或启动克隆资源副本时，是否在开始操作前和操作完成后告知其他所有资源副本，允许值为 False和 True，默认值为 False　　
        - Globally-unique：表示是否允许每个克隆副本资源执行不同的功能，允许值为 False和 True。如果其值为 False,则不管这些克隆副本资源运行在何处，它们的行为都是完全和同的，因此每个节点中有且仅有一个克隆副本资源处于 Active状态。其值为 True,则运行在某个节点上的多个资源副本实例或者不同节点上的多个副本实例完全不一样。如果Clone-node-max取值大于1，即一个节点上运行多个资源副本，那么 Globally-unique的默认值为 True，否则为 False.
        - Ordered：表示是否顺序启动位于不同节点上的资源副本，true为顺序启动，false为并行启动，默认值是 False.
        - Interleave：该属性值主要用于改变克隆资源或者 Masters资源之间的 ordering约束行为， Interleave可能的值为 True和 False,如果其值为 False,则位于相同节点上的后一个克隆资源的启动或者停止操作需要等待前一个克隆资源启动或者停止完成才能进行，而如果其值为 True,则后一个克隆资源不用等待前一个克隆资源启动或者停止完成便可进行启动或者停止操作。 Interleave的默认值为 False

    - master/slave : 将primitive克隆2份、其中master和slave节点各运行一份，且只能在这2个节点上运行

        主从资源，如drdb

    Pacemaker中的资源类型使用standard, provider（仅当standard为ocf使用）和agent来进行标识，格式如下： `<standard>:[provider]:<agent>`.

- RA类型： 资源代理(resource agent)是一种标准化的集群接口, 每一个原始资源(primitive Resource)都有一个资源代理, packmaker**通过该接口对集群资源进行操作**.

    - Lsb(linux standard base resource agents)： 一般位于/etc/rc.d/init.d/目录下的支持start|stop|status等参数的服务脚本都是lsb
    - systemd
    - ocf：Open cluster Framework，开放集群架构, 是对LSB资源代理的扩展, **已成为使用最多的资源类别**, 在`/usr/lib/ocf/resource.d/provider`

        OCF 资源脚本至少包含 start, stop, status,monitor 以及 meta-data 执行动作。其中,meta-data 动作给出如何配置该脚本
    - heartbeat：heartbaet V1版本
    - stonith：专为配置stonith设备而用

    > 在多数情况下，资源 RA以 shell脚本的形式提供，当然也可以使用其他语言来实现 RA
    > OCF标准还严格定义了操作执行后的状态码，集群资源管理器将会根据资源代理返回的状态码来对执行结果做出判断

- fence device的原理及作用

    fence device用于强制隔离设备. 如果某个节点没有反应，并不代表没有数据访问, 能够 100% 确定数据安全的唯一方法是使用 SNOITH 隔离该节点，这样才能确定在允许从另一个节点访问数据前，该节点已确实离线.

    FENCE设备是RHCS集群中必不可少的一个组成部分，通过FENCE设备可以避免因出现不可预知的情况而造成的“脑裂”现象. FENCE设备的出现，就是为了解决类似这些问题，Fence设备主要就是通过服务器或存储本身的硬件管理接口或者外部电源管理设备，来对服务器或存储直接发出硬件管理指令，将服务器重启或关机，或者与网络断开连接.

    当意外原因导致主机异常或者宕机时，备机会首先调用FENCE设备，然后通过FENCE设备将异常主机重启或者从网络隔离，当FENCE操作成功执行后，返回信息给备机，备机在接到FENCE成功的信息后，开始接管主机的服务和资源. 这样通过FENCE设备，将异常节点占据的资源进行了释放，保证了资源和服务始终运行在一个节点上.

    RHCS的FENCE设备可以分为两种：内部FENCE和外部FENCE，常用的内部FENCE有IBMRSAII卡，HP的iLO卡，Dell的DRAC, 还有IPMI的设备等，外部fence设备有UPS、SANSWITCH、NETWORKSWITCH等.

    对于外部fence 设备，可以做拔电源的测试，因为备机可以接受到fence device返回的信号，备机可以正常接管服务.
    对于内部fence 设备，不能做拔电源的测试，因为主机断电后，备机接受不到主板芯片做为fence device返备的信号，就不能接管服务，`crm status`会看到资源的属主是unknow,查看日志会看到持续报fence failed的信息

pacemaker支持的集群类型:
- Active/Passive(主从)

    ![](https://clusterlabs.org/pacemaker/doc/en-US/Pacemaker/2.0/html-single/Pacemaker_Administration/images/pcmk-active-passive.png)

    每个节点上都部署有相同的服务实例，但是正常情况下只有一个节点上的服务实例处于激活状态，只有当前活动节点发生故障后，另外的处于 standby状态的节点上的服务才会被激活，这种模式通常意味着需要部署额外的且正常情况下不承载负载的硬件.
- Active/Active

    ![](https://clusterlabs.org/pacemaker/doc/en-US/Pacemaker/2.0/html-single/Pacemaker_Administration/images/pcmk-active-active.png)

    故障节点上的访问请求或自动转到另外一个正常运行节点上，或通过负载均衡器在剩余的正常运行的节点上进行负载均衡。这种模式下集群中的节点通常部署了相同的软件并具有相同的参数配置，同时各服务在这些节点上并行运行.
- N+1

    即多准备一个额外的备机节点，当集群中某一节点故障后该备机节点会被激活从而接管故障节点的服务. 在不同节点安装和配置有不同软件的集群中，即集群中运行有多个服务的情况下，该备机节点应该具备接管任何故障服务的能力，而如果整个集群只运行同一个服务，则N+1模式便退变为 Active/Passive模式
- N+M

    在单个集群运行多种服务的情况下，N+1模式下仅有的一个故障接管节点可能无法提供充分的冗余，因此，集群需要提供 M（M>l）个备机节点以保证集群在多个服务同时发生故障的情况下仍然具备高可用性， M的具体数目需要根据集群高可用性的要求和成本预算来权衡.
- N-to-1

    允许接管服务的备机节点临时成为活动节点（此时集群已经没有备机节点），但是，当故障主节点恢复并重新加人到集群后，备机节点上的服务会转移到主节点上运行，同时该备机节点恢复 standby状态以保证集群的高可用
- N-to-N

    是 Active/Active模式和N+M模式的结合， N-to-N集群将故障节点的服务和访问请求分散到集群其余的正常节点中，在N-to-N集群中并不需要有Standby节点的存在、但是需要所有Active的节点均有额外的剩余可用资源.

集群的全生命周期管理工具:
- pcs

    专用于pacemaker+corosync的设置工具,有CLI和web-based GUI界面. 
- crmsh, **推荐**

    是管理pacemaker的命令行界面, 是基于ssh进行远程管理

推荐: 先使用pcs安装配置集群，然后使用crmsh来管理进群, 因为crmsh管理集群比较方便.

> [pcs与crmsh命令比较](https://www.cnblogs.com/gzxbkk/p/7305227.html)或[pcs-crmsh-quick-ref.md](https://github.com/ClusterLabs/pacemaker/blob/master/doc/pcs-crmsh-quick-ref.md)

资源规则:
资源规则（Rule）使得 pacemaker集群资源具备了更强的动态调节能力，资源规则最常见的使用方式就是在集群资源运行时设置一个合理的粘性值（Resource-stickness)'以防止资源意外回切.

Rule通常在 Constraint命令中配置，如下语句是配置资源 Rule的语法格式：`pcs constraint rule add constraint_id [rule_type] [score=score] (id=rule-id] expression丨date_expression丨date-spec options`.

如果忽略 score值，则使用默认值 INFINITY,如果忽略 ID,则自动从 Constraint_id生成一个规则 ID,而 Rule-type可以是字符表达式或者日期表达式.

rule expression支持:
- defined|not_defined attribute
- attribute lt|gt|Ite|gte|eq|ne value
- date [start=start] [end=end] operation=gt|lt|in-range
- date-spec date_spec_options

```bash
# pcs property set --node computel1 osprole=compute # 设置node属性
# pcs resource constraint location resource_id rule [rule_id] [role=master|slave] [score=score expression] # 通过Rule的节点属性表达式来确定资源的Location
# pcs constraint location nova-compute-clone rule resource-discovery=exclusive score=0 osprole eq compute # rule_id = `resource-discovery=exclusive`, expression =`osprole eq compute`
# pcs constraint rule remove rule_id # 删除资源 Rule语法
```

## corosync & pacemaker高可用解决方案
参考:
- [Corosync+Pacemaker+crmsh构建Web高可用集群](https://www.cnblogs.com/cloudos/p/8336529.html)
- [Red Hat Enterprise Linux 8 High Availability Add-On Administration]
- [从头开始搭建集群 在Fedora上面创建主/主和主备集群](https://clusterlabs.org/pacemaker/doc/deprecated/zh-CN/Pacemaker/1.1-plugin/html-single/Clusters_from_Scratch/index.html)
- [Configuring and managing high availability clusters from rhel8](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html/configuring_and_managing_high_availability_clusters/index)
- [High Availability Add-On 参考 from rhel7](https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/7/html/high_availability_add-on_reference/index)

环境准备for 每个node:
```bash
# ---每个node都执行
# chronyc makestep # 同步时间, corosync的要求, 最好也设置crontab.
# vim /etc/hosts # 设置hosts
...
192.168.0.41 node11
192.168.0.42 node12
# apt install nginx
# vim /var/www/html/index.nginx-debian.html # 修改nginx default page, 在<title>中加入每个node的hostname
# curl localhost # 验证上一步的修改是否已生效
# apt install pacemaker corosync pcs crmsh # 使用了crmsh，就不再需要安装heartbeat
# systemctl start pcsd && systemctl enable pcsd
# echo hacluster:123456 | chpasswd # for ubuntu. `echo "123456" | passwd --stdin hacluster` for centos, 为hacluster用户(created by pcs)设置密码
# ---在任一节点执行
# pcs cluster auth node11  node12 -u hacluster # 认证节点
# pcs cluster setup --force --name mycluster node11  node12 # 配置集群, 会自动创建/etc/corosync/corosync.onf, 并同步到所有node
# pcs cluster start --all # 启动cluster
# pcs cluster enable --all # 设置自动启动
# pcs status # 检测cluster status, 类似`crm_mon -1`
...
Current DC: node11 (version 1.1.18-2b07d5c5a9) - partition with quorum # 当前的仲裁节点(Current DC, DC即为Designated Co-ordinator）为node11, 这个节点负责向集群中的节点发出一系列指令，使各个资源按照定义（存储在cib数据库中）启动或停止
...
# corosync-cfgtool -s # 查看当前心跳的位置, 仅显示当前node的状态
# corosync-cmapctl # 检查集群成员关系
# pcs status corosync # 检查corosync状态
# pcs property set stonith-enabled=false # pcs/crmsh可能会因为stonith resource未定义而报错`error: unpack_resources: ... no STONITH resources have been defined`, 由于没有stonith设备此时可以通过修改集群全局属性将此error先屏蔽
# --- 基础的集群已经配置完毕
```

> pcs支持https://node:2224的管理web, 可通过`Add Existing`添加现有集群或使用`Create New`创建新集群.

查看crm ra支持的资源:
```bash
crm(live)# cd ra
crm(live)ra# classes # 查看cluster支持的类型
lsb
ocf / .isolation heartbeat pacemaker redhat
service
stonith
systemd
crm(live)ra# list lsb # 查看该类别下可用的ra
crm(live)ra# info ocf:heartbeat:IPaddr # 查看该ra的help
```

常见资源：
－　primitive webip ocf:heartbeat:IPaddr２ parms ip="192.168.10.100" cidr_netmask="24"
－　primitive webserver systemd:nginx op start timeout=100s op stop timeout=100s
－　primitive webstore ocf:heartbeat:Filesystem params device="192.168.10.9:/data/web/htdocs" directory="/var/www/html" fstype="nfs" op start timeout=60s op stop timeout=60s

创建/修改资源:
```bash
# crm
crm(live)# configure
crm(live)configure# primitive webip ocf:heartbeat:IPaddr2 params ip=192.168.10.160 cidr_netmask=24 [--group webservice] # 定义webip资源
crm(live)configure# edit webip # 修改webip资源
crm(live)configure# primitive webserver systemd:nginx # 定义webserver资源
crm(live)configure# group webservice webip webserver # 资源默认会平均分配在各个不同的几点. 定义组资源，组名为webservice，使webip和webserver服务在同一节点上. **组中在前面的先启动**
crm(live)configure# verify # 验证配置
crm(live)configure# commit # 提交配置
...
crm(live)node# standby node11 # 模拟node下线, 查看webip所在
crm(live)node# online node11 # node重新上线
crm(live)node# cd ..
crm(live)# exit
# crm_verify -L -V # 检查crm配置是否正确
# systemctl stop pacemaker corosync # on node12, 其上的pacemaker resource会全部撤销
# ip addr # on node11, 看到webip出现
```

或
```bash
# pcs cluster stop --all # 将集群节点都停掉
# pcs resource create vip ocf:heartbeat:IPaddr2 ip=192.168.56.24 cidr_netmask=24 op monitor interval=30s
# pcs resource create nfs-root exportfs clientspec=192.168.0.0/255.255.0.0 options=rw,sync,no_root_squash directory=/mnt/nfs/exports --group nfsgroup
# pcs resource list # 查看pacemaker支持的资源
# pcs resource describe ocf:heartbeat:IPaddr2 # 查看该类型资源的描述
# pcs stonith create my_fence fence_kdump test1 # 创建fence_kdump
# pcs stonith show my_fence # 查看fence状态
# pcs cluster standy node11 # 将node状态置为standby
# pcs cluster stop node11 # 将node状态置为offline
# crm node standby # 将node设为standby
# crm node list # 查看所有node
# crm configure show # 查看当前配置
# corosync-quorumtool -l # 显示所有节点的信息与票数
# pcs cluster sync # 同步所有节点信息
# corosync-quorumtool -siH # 票数细节, 类似`pcs quorum status`
# pcs stonith list # 查看可用fence插件
# crm resource move webservice node12 # 手动转移资源
# crm_resource --list-raw # 资源列表
# crm configure show ${resource} # 查看resoure的配置
# crm_resource --locate --resourece ${resource} # 查看resoure所在node
# crm_failcount --resource ${resource} --node ${node} # 查看资源的故障计数
# crm resource clean ${resource} [${node}] # 清理资源的status, 比如failcount
# crm_resource --resource ${resource} --move --host ${node} # 转移资源, 不阻塞, 因此需要轮询资源是否已在指定节点且在该节点的failcount为0
# cibadmin --modify --xml-text '<op id="xxx-monitor-30" enabled="true">' # 启用资源上的monitor操作
# crm resource manage xxx # 允许crm管理资源
# crm_resource --resource xxx --get-parameter is-managed --meta # 检查资源是否已被crm管理
```

corosync+pacemaker集群默认对节点高可用，但是对于节点上资源的运行状态无法监控，因此，需要配置集群对于资源的监控，在资源因意外情况下，无法提供服务时，对资源提供高可用.

```bash
# crm configure
crm(live)configure# monitor webip 30s:20s # interval=30s, timeout=20s
crm(live)configure# monitor webserver 30s:100s
crm(live)configure# verify
crm(live)configure# commit
```

### pacemaker资源约束
用以指定在哪些群集节点上运行资源，以何种顺序装载资源，以及特定资源依赖于哪些其它资源.
pacemaker提供了三种资源约束方法：
- Resource Location（资源位置）：定义资源可以、不可以或尽可能在哪些节点上运行即限定了资源仅在哪些节点上启动运行, 类似k8s node上的tags

    资源对节点的倾向程度，通常可以使用一个分数（score）来定义，当score为正值时，表示资源倾向与此节点；负值表示资源倾向逃离于此节点. 也可以将score定义为-inf(负无穷大)和inf（正无穷大）.
- Resource Collocation（资源排列）：捆绑约束将不同的资源捆绑在一起作为一个逻辑整体来调度

    通常也是使用一个score来定义的. 当score是正值表示资源可以在一起；否则表示不可以在一起
- Resource Order（资源顺序）：顺序约束定义集群资源在节点上的运行顺序(启动的顺序, 关闭顺序与启动顺序相反)

定义约束时，还需要指定分数. 各种分数是集群工作方式的重要组成部分, 其实，从迁移资源到决定在已降级集群中停止哪些资源的整个过程是通过以某种方式修改分数来实现的. 分数按每个资源来计算，资源分数为负的任何节点都无法运行该资源. 在计算出资源分数后，集群选择分数最高的节点. INFINITY（无穷大）目前定义为 1,000,000. 加减无穷大遵循以下3个基本规则：
- 任何值 + 无穷大 = 无穷大
- 任何值 - 无穷大 = -无穷大
- 无穷大 - 无穷大 = -无穷大

定义资源约束时，也可以指定每个约束的分数. 分数表示指派给此资源约束的值. 分数较高的约束先应用，分数较低的约束后应用. 通过使用不同的分数为既定资源创建更多位置约束，可以指定资源要故障转移至的目标节点的顺序.

去除group属性，也可通过约束来完成资源的组合:
```bash
# crm configure
crm(live)configure# delete webservice
crm(live)configure# location webip_on_node1 webip inf: node11
crm(live)configure# colocation webip_with_webserver_with_webstore inf: webip webserver webstore
crm(live)configure# order webstore_befroe_webserver_befroe_webip  Mandatory: webstore webserver webip # Mandatory, Optional, and Serialize分别表示强制, 可选, 串行. inf即INFINITY
crm(live)configure# verify   //检查语法是否有错误
crm(live)configure# commit   //提交配置
```

### 法定票数问题
quorum policy是集群存货的要素. 假设有三个节点且每个节点有一张票，那么需要两个节点存活才算是一个集群: 算法为 N/2+1, N 代表票数.

修改node默认票数:
```conf
# from `man 5 votequorum`
nodelist {
    node {
        ring0_addr: 192.168.1.1
        quorum_votes: 3
    }
    ....
}
```

在双节点集群中，由于票数是偶数，当心跳出现问题`脑裂`时，两个节点都将**达不到法定票数，默认quorum策略会关闭集群服务**, 为了避免这种情况，可以增加票数为奇数（增加ping节点），或者调整默认quorum策略为`ignore` by `property no-quorum-policy="ignore"`

Pacemaker 中的法定人数功能可以防止`脑裂`.

#### pcs
pcs 常用命令如下:
- cluster : 配置集群选项和节点
- resource : 创建和管理集群资源
- stonith : 将 fence 设备配置为与 Pacemaker 一同使用
- constraint : 管理资源限制
- property : 设定 Pacemaker 属性
- status [commands]: 查看当前集群和资源状态

    如果未指定 commands 参数，这个命令会显示有关该集群和资源的所有信息. 指定 resources、groups、cluster、nodes 或 pcsd 选项则可以只显示具体集群组件的状态
- config : 以用户可读格式显示完整集群配置

pcs集群创建步骤:
- 认证组成集群的节点
- 配置和同步集群节点
- 在集群节点中启动集群服务

```bash
# pcs cluster cib # 查看原始集群配置
# pcs cluster cib bak.xml # 将原始 xml 从 CIB 保存到名为 bak.xml
# pcs -f bak.xml resource create VirtualIP ocf:heartbeat:IPaddr2  # 在 bak.xml 中创建一个资源，但不会将该资源添加到当前运行的集群配置中
# pcs cluster cib-push bak.xml # 将 bak.xml 的当前内容 push 回 CIB 中
# pcs config # 显示当前集群的完整配置
# --- 备份和恢复集群配置
# pcs config backup ${filename} # 备份集群配置
# pcs config restore [--local] [filename] # 恢复集群配置. 如果没有指定文件名，则会使用标准输入. `--local` 则只会恢复当前节点中的文件
# --- 创建cluster
# pcs cluster auth [node] [...] [-u username] [-p password] # 认证集群节点. 每个节点中的 pcs 管理员的用户名必须为 hacluster并使用相同的 hacluster 密码. 授权令牌会保存在 ~/.pcs/tokens（或 /var/lib/pcsd/tokens）文件中
# pcs cluster setup --name cluster_ name node1 [node2] [...] # 配置集群配置文件，并将配置同步到指定节点中
# pcs cluster start [--all|node[,...]] # `--all`表示会在所有节点中启动集群服务, 如果未指定任何节点，则只在本地节点中启动集群服务
# pcs cluster enable [--all|node[,...]] # 将集群服务配置为在指定节点启动时运行
# pcs cluster disable [--all|node[,...]] # 将集群服务配置为在指定节点启动时不运行
# pcs cluster stop [--all|node[,...]] # `--all`即在所有节点中停止集群服务. 如果未指定任何节点，则只在本地节点中停止集群服务
# pcs cluster kill # 强制停止本地节点中的集群服务
# pcs cluster destroy # 删除集群配置
# pcs cluster status # 仅显示显示集群状态
# --- 管理node
# pcs cluster node add|remove ${node} # 添加或删除集群节点. 添加时 corosync.conf 会同步到集群的所有节点中，包括新添加的节点
# pcs cluster standby node | --all # 让指定节点进入待机模式, 该指定节点不再托管资源. 当前节点中活跃的所有资源都将被移动到另一个节点中
# pcs cluster unstandby node | --all # 将指定节点从待机模式中移除. 运行这个命令后，该指定节点就可以托管资源
# --- 设定用户权限, 前提user是组 haclient 的成员
# pcs property set enable-acl=true --force # 启用 Pacemaker ACL
# pcs acl role create read-only description="Read access to cluster" read xpath /cib # 使用只读权限为 cib 创建名为 read-only 的角色. 写权限是`write xpath /cib`
# pcs acl user create rouser read-only # 在 pcs ACL 系统中创建用户 rouser，并为那个用户分配 read-only 角色
# --- status
# pcs status # 显示集群状态(pcs cluster status)及其资源的状态(pcs status resources)
# --- 配置 STONITH, 部分设备支持 fencing 拓扑功能(支持包含多个设备的 fencing 节点, 它们用优先级来指明尝试stonith的顺序, 不常用), 相关命令是`pcs stonith level`
# pcs stonith list [filter] # 查看所有可用 STONITH 代理列表
# pcs stonith describe stonith_agent # 查看指定 STONITH 代理
# pcs stonith create MyStonith fence_virt pcmk_host_list=f1 op monitor interval=30s # MyStonith, 是stonith_id;fence_virt, stonith_device_type; "pcmk_host_list=f1 op monitor interval=30s", stonith_device_options, pcmk_host_list是这个资源控制的机器列表.
# pcs stonith show [stonith_id] [--full] # `--full`显示所有配置的 stonith 选项
# pcs stonith update stonith_id [stonith_device_options] # 修改或添加选项
# pcs stonith delete stonith_id # 删除 fencing 设备
# pcs stonith fence node [--off] # 手动隔离某个节点. `--off`会使用 off API 调用 stonith，从而关闭节点，而不是重启节点.
# pcs stonith confirm node # 确定指定的节点目前是否已被关闭
# --- 配置集群资源
# pcs resource create VirtualIP ocf:heartbeat:IPaddr2 ip=192.168.0.120 cidr_netmask=24 op monitor interval=30s meta resource-stickiness=5O # 创建一个名为 VirtualIP，使用 ocf 标准, 由heartbeat 提供程序，以及类型 IPaddr2 的资源. 这个资源的浮动地址为 192.168.0.120. 为保证资源正常工作，可在资源定义中添加监控操作, 当前系统会每 30 秒检查一次，确定该资源是否运行. meta用于设置资源的元数据选项.
# pcs resource delete VirtualIP # 删除配置的资源
# pcs resource update VirtualIP ip=192.169.0.120 # 修改配置资源的参数
# pcs resource list  # 显示所有可用资源
# pcs resource standard  # 显示可用资源代理标准
# pcs resource providers # 显示可用资源代理提供程序列表
# pcs resource list string   # 显示根据指定字符串过滤的可用资源列表。可使用这个命令显示根据标准名称、提供程序或类型过滤的资源
# pcs resource describe standard:provider:type|type # 显示该资源设定的参数
# pcs resource defaults resource-stickiness=100 # 更改资源选项的默认值
# pcs resource defaults # 显示当前配置的默认值列表
# pcs resource op defaults [options] # 获取监控操作全局默认值
# pcs resource op defaults timeout=240s # 为所有监控操作将全局 timeout 值设定为 240s
# pcs resource meta dummy_resource failure-timeout=20s # 资源dummy_resource可在 20 秒内尝试在同一节点中重启
# pcs resource show --full # 显示所有配置的资源列表及为那些资源配置的参数
# pcs resource show dummy_resource 查看该资源的配置
# pcs resource group add shortcut IPaddr Email # 创建 shortcut 资源组，该资源组包含现有资源 IPaddr 和 Email
# pcs resource op remove VirtualIP stop interval=0s timeout=20s 删除停止超时操作
# pcs resource op add VirtualIP stop interval=0s timeout=40s # 添加停止超时操作
# pcs resource enable/disable resource_id # 启用和禁用集群资源
# pcs resource cleanup resource_id # 清除 resource_id 指定的资源
# --- 属性
# crm_attribute --type crm_config --name xxx --query # 查询指定的集群属性from cib
# crm_attribute --type crm_config --name xxx --update xxx # 设置指定的集群属性
```

#### crmsh
- cibadmin : 操作 CIB 的基础管理命令
- crm_attribute : 对 CIB 进行查询和修改
- crm_diff : 帮助生成和应用 CIB XML 补丁
- crm_verify : 用于验证 CIB 的一致性、检测其它错误和测试是否可以联机到正在运行的集群
- crm_resource : 负责与 CRM 进行交互, 可以启动、停止、删除或者迁移在集群节点上的资源
- crm_failcount : 查询当前资源错误统计

    错误统计是资源监视器的附加属性,它的值会根据资源监视到的故障而递增,它与资源错误粘性数值(migration-threshold)相乘得到结果为错误切换分值. 如果这个数值超过设置的大小,资源就会发上切换, 除非错误统计数被重置. **要想资源会切,需删除错误统计数**.

- crm_standby : 用来控制备机属性, 以决定资源是否可以运行在该节点上
- crm_mon : 配置和监视集群状态, 该命令输出节点数量、名称、UUID 以及状态


```bash
# crm_attribute --name maintenance-mode --query --type crm_config [--quiet] # 属性查询
scope=crm_config  name=maintenance-mode value=false
# crm_failcount --resource OS-1401dc36-node1-ha-flv --node OS-1401dc36-node1 --quiet # 管理记录每个资源的故障计数的计数器
# crm_resource --list-raw # 显示资源
# crm_resource --locate --resource OS-1401dc36-node1-ha-flv # 显示资源所在node
``` 

### 防​止​资​源​在​节​点​恢​复​后​移​动​
故障发生时，资源会迁移到正常节点上，但当故障节点恢复后，资源可能再次回到原来节点，这在有些情况下并非是最好的策略，因为资源的迁移是有停机时间的，资源在节点间每一次的来回流动都会造成那段时间内节点提供的服务无法被正常访问，特别是一些复杂的应用，如MySQL数据库，其停机时间会更长. 为了避免这种情况，可以根据需要，使用资源粘性策略: `rsc_defaults resourcce-stickiness=200`.


# keepalived
linux下轻量级ha解决方案, 主要通过虚拟路由冗余来实现ha. 与heartbeat相比, 没有其功能强大, 但在某些场景下更简单合适.
keepalived起初是为LVS设计, 专门监控(根据3~5层的交换机制)集群中各节点的状态, 通过剔除故障节点和加入正常节点来实现ha.
keepalived支持VRRP(Virtual Router Redundancy Protocol, 虚拟路由器冗余协议, 为了解决静态路由器出现的单点故障).

keepalived运行机制:
- 网络层

    使用ICMP探测节点状态
- 传输层

    利用tcp的端口连接和扫描来判断节点状态
- 应用层

    通常是自定义keepalived来处理.

keepalived配置: `/etc/Keepalived/Keepalived.conf`, 可分为:
- 全局配置(global configuration)
- vrrpd配置
- lvs配置

## vrrp
将两台及以上的物理路由器虚拟成一个虚拟路由器, 它有一个唯一标识vrid(与一组ip关联), 同一时刻仅有一台物理路由器对外提供服务.

# lb
## lvs(base kernel)
参考:
- [LVS原理篇：LVS简介、结构、四种模式、十种算法](https://blog.csdn.net/lcl_xiaowugui/article/details/81701949)

LVS的模型中有两个角色：
- 调度器:Director Server，又称为Dispatcher，Balancer

    调度器主要用于接受用户请求
- 真实主机:Real Server，简称为RS

    用于真正处理用户的请求

将所在角色的IP地址分为以下4种：
- Director Virtual IP:调度器用于与客户端通信的IP地址，简称为VIP
- Director IP:调度器用于与RealServer通信的IP地址，简称为DIP
- Real Server : 后端主机的用于与调度器通信的IP地址，简称为RIP
- CIP：Client IP，访问客户端的IP地址所

lvs的ip负载均衡是通过ipvs实现的, 具体机制分3种:
- NAT

    用户请求LVS到达director，director将请求的报文的目的IP改为RIP和将报文的目标端口也改为realserver的相应端口，最后将报文发送到realserver上，realserver将数据返回给director，director再将响应报文的源地址和源端口换成vip和相应端口, 最后把数据发送给用户

    优点:
    - 支持端口映射
    - RS可以使用任意操作系统
    - 节省公有IP地址
    
        RIP和DIP都应该使用同一网段私有地址，而且RS的网关要指向DIP. 使用nat另外一个好处就是后端的主机相对比较安全.

    缺点：
    - 效率低: 请求和响应报文都要经过Director转发;极高负载时，Director可能成为系统瓶颈
- FULLNAT


    FULLNAT模式也不需要DIP和RIP在同一网段
    FULLNAT和NAT相比的话：会保证RS的回包一定可到达LVS
    FULLNAT需要更新源IP，所以性能正常比NAT模式下降10%
- TUN

    基于ip隧道技术.
    用户请求LVS到达director，director通过IP-TUN技术将请求报文的包封装到一个新的IP包里面，目的IP为VIP(不变)，然后director将报文发送到realserver，realserver基于IP-TUN解析出来包的目的为VIP，检测网卡是否绑定了VIP，绑定了就处理这个包，如果在同一个网段，将请求直接返回给用户，否则通过网关返回给用户；如果没有绑定VIP就直接丢掉这个包

    优点：

    - RIP,VIP,DIP都应该使用公网地址，且RS网关不指向DIP
    - 只接受进站请求，请求报文经由Director调度，但是响应报文不需经由Director, 解决了LVS-NAT时的问题，减少负载

    缺点：
    - 不指向Director所以不支持端口映射
    - RS的OS必须支持隧道功能
    - 隧道技术会额外花费性能，增大开销, 同时运维麻烦, 一般不用.

- DR

    当Director接收到请求之后，通过调度方法选举出RealServer, 将目标地址的MAC地址改为RealServer的MAC地址. RealServer接受到转发而来的请求，发现目标地址是VIP, RealServer配置在lo接口上, 处理请求之后则使用lo接口上的VIP响应CIP.

    优点:
    - RIP可以使用私有地址，也可以使用公网地址, 只要求DIP和RIP的地址在同一个网段内.
    - 请求报文经由Director调度，但是响应报文不经由Director. 性能最高最好, 因此**最常用**.
    - RS可以使用大多数OS

    缺点：
    - 不支持端口映射
    - 不能跨局域网

四种模式的比较:
- 是否需要VIP和realserver在同一网段
    
    DR模式因为只修改包的MAC地址，需要通过ARP广播找到realserver，所以VIP和realserver必须在同一个网段，也就是说DR模式需要先确认这个IP是否只能挂在这个LVS下面；其他模式因为都会修改目的地址为realserver的IP地址，所以不需要在同一个网段内
- 是否需要在realserver上绑定VIP
    
    realserver在收到包之后会判断目的地址是否是自己的IP
    DR模式的目的地址没有修改，还是VIP，所以需要在realserver上绑定VIP
    IP TUN模式值是对包重新包装了一层，realserver解析后的包的IP仍然是VIP，所以也需要在realserver上绑定VIP
- 四种模式的性能比较
    
    DR模式、IP TUN模式都是在包进入的时候经过LVS，在包返回的时候直接返回给client；所以二者的性能比NAT高
    但TUN模式更加复杂，所以性能不如DR
    FULLNAT模式不仅更换目的IP还更换了源IP，所以性能比NAT下降10%
    性能比较：DR>TUN>NAT>FULLNAT

LVS的八种调度方法
- 静态方法:仅依据算法本身进行轮询调度

    - RR:Round Robin,轮调 : 均等地对待每一台服务器，不管服务器上的实际连接数和系统负载
    - WRR:Weighted RR，加权论调 : 加权，手动让能者多劳
    - SH:SourceIP Hash : 与目标地址散列调度算法类似，但它是根据源地址散列算法进行静态分配固定的服务器资源. 来自同一个IP地址的请求都将调度到同一个RealServer
    - DH:Destination Hash : 根据目标 IP 地址通过散列函数将目标 IP 与服务器建立映射关系，出现服务器不可用或负载过高的情况下，发往该目标 IP 的请求会固定发给该服务器. 不管IP，请求特定的东西，都定义到同一个RS上

- 动态方法:根据算法及RS的当前负载状态进行调度

    - LC:least connections(最小链接数) : 链接最少，也就是Overhead最小就调度给谁; 假如都一样，就根据配置的RS自上而下调度
    - WLC:Weighted Least Connection (加权最小连接数) : 这个是LVS的默认算法. 调度器可以自动问询真实服务器的负载情况，并动态调整权值
    - SED:Shortest Expection Delay(最小期望延迟) : WLC算法的改进. 不考虑非活动链接，谁的权重大，优先选择权重大的服务器来接收请求，但权重大的机器会比较忙
    - NQ:Never Queue : SED算法的改进
    - LBLC:Locality-Based Least-Connection,基于局部的的LC算法

        请求数据包的目标 IP 地址的一种调度算法，该算法先根据请求的目标 IP 地址寻找最近的该目标 IP 地址所有使用的服务器，如果这台服务器依然可用，并且有能力处理该请求，调度器会尽量选择相同的服务器，否则会继续选择其它可行的服务器
    - lblcr, 复杂的基于局部性最少的连接算法 : 记录的不是要给目标 IP 与一台服务器之间的连接记录，它会维护一个目标 IP 到一组服务器之间的映射关系，防止单点服务器负载过高

## haproxy
基于应用实现的软负载均衡, 基于tcp和http.

## FAQ
### pacemaker : Cannot use default route w/o netmask
该错误是`systemctl restart pacemaker`时`journalctl -f`截获的.

原因: 定义`ocf:heartbeat:IPaddr`时未指定`cidr_netmask`
