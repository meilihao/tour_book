# ovirt
- [**oVirt虚拟化关键概念、组件与技术原理**](http://simiam.com/2018/07/11/oVirt%E8%99%9A%E6%8B%9F%E5%8C%96%E5%85%B3%E9%94%AE%E6%A6%82%E5%BF%B5%E3%80%81%E7%BB%84%E4%BB%B6%E4%B8%8E%E6%8A%80%E6%9C%AF%E5%8E%9F%E7%90%86/)
- [oVirt 架构学习](https://cloud.tencent.com/developer/article/1435899)
- [编译oVirt相关软件包](https://www.hikunpeng.com/document/detail/zh/kunpengcpfs/ecosystemEnable/oVirt/kunpengovirtoe_04_0004.html)
- [Next-Gen Backup & DR Solution for oVirt](https://resources.ovirt.org/site-files/2020/Next-Gen%20oVirt%20Backup%20%26%20DR%20Solution%20-%20vinchin%20-%20official%20presentation.pdf)
- [架构 Background](https://www.ovirt.org/develop/release-management/features/virt/enhance-import-export-with-ova.html)
- [oVirt导入导出ova格式虚机](https://www.cnovirt.com/archives/938)

oVirt基于KVM，并整合使用了libvirt、gluster、patternfly、ansible等一系列优秀的开源软件，oVirt的定位是替代VMware vsphere，oVirt目前已经成为了企业 虚拟化  环境可选的解决方案，另外相比OpenStack的庞大和复杂，oVirt在企业私有云建设中具备部署和维护使用简单的优势.

oVirt项目中的不同组件主要包含三个部分：
- ovirt-engine: 用来进行管理虚拟机（创建、开关启停）、配置网络和存储等操作
- 一个或多个主机（节点, 比如定制的ovirt node）: 用来运行虚拟机

	主机节点是安装有VDSM和libvirt组件的Linux发行版，也包含一些用来实现网络虚拟化和其它系统服务的组件
- 一个或多个存储节点: 用来存放虚拟机镜像和ISO镜像

	存储节点可以使用块存储或文件存储，可以利用主机节点自身的存储做存储节点（local on host模式），或者使用外部的存储，例如通过 NFS  访问，或者是IP-SAN/FC-SAN。还有一种就是超融合架构，通过Gluster将主机节点自身的磁盘组成池来使用，同时能够实现高可用和冗余。

组件:
1. Engine(ovirt-engine)-对物理服务器进行统一管理，允许用户创建和部署新的虚拟机
2. Admin Portal-engine的webadmin页面，主要用于系统的管理员登陆，里面可以进行平台内所有的操作
3. User Portal-主要用于普通用户登陆的页面，里面只能进行一些虚拟机的基本操作
4. REST API-可以通过应用调用它提供的API进行一系列的虚拟化操作，它对外提供命令行和python SDK
5. CLI/SDK-engine对外提供命令行接口和SDK来和它进行操作
6. Database-使用的是Postgres数据库，主要用于提供对平台数据的存储，支持备份和恢复功能
7. Host agent(VDSM)-oVirt engine通过和它进行通信在服务器上对虚拟机进行操作，作为一个中间层的作用，向上和ovirt engine通信，向下连接虚拟化软件qemu-kvm，主要对网络/计算/存储资源的调度和管理
8. Guest Agent-它被安装在虚拟机内部，主要功能是给ovirt-engine提供这台虚拟机的各种资源的使用信息，和ovirt-engine通过虚拟串口通信
9. AD/IPA-目录服务系统，ovirt-engine通过他们获取到平台内域用户的登陆权限
10. DWH(Data Warehouse)-数据仓库通过ETL架构进行数据的采集，并将数据保存在history DB(数据库-‘ovirt_engine_history’)内
11. Reports Engine-使用Jasper Reports架构对history DB内的数据进行页面展示
12. SPICE client-用户可以通过spice客户端访问虚拟机

概念:
1. 数据中心是ovirt虚拟化环境中的根容器（最高一级的逻辑项），它包括以下三个子容器（子项）

	1. 集群: 用来保存与集群相关的信息

		集群由一个至多个计算节点（主机）组成，这些主机具有相互兼容处理器内核。一个集群组成了一个虚拟机的迁移域，虚拟机可以被实时迁移到同一集群中的其它主机上。一个数据中心可以包括多个集群，一个集群可以包括多个主机。

	1. 存储: 用来保存存储类型、存储域的信息，以及存储域间的连接信息。存储域在数据中心一级上定义，并可以被数据中心的所有集群使用（即可以被集群中的所有主机挂载）
	1. 网络: 用来保存与数据中心中的逻辑网络相关的信息，如网络地址、VLAN标签等信息。逻辑网络在数据中心一级上定义，并可以在集群一级上使用
1. ovirt网络: 用来处理虚拟机间网络连接的逻辑网络是通过计算节点上的基于软件的网桥实现的。在默认情况 下，ovirt-engine在安装过程中会创建一个名为“ovirtmgmt管理网络”的逻辑网络。此外系统管理员还可以添加专用的存储逻辑网络和专用的显示逻辑网络
1. 存储域就是一系列具有公共存储接口的镜像的集合，存储域中包括了虚拟机模板、快照、数据镜像、ISO文件以及存储域本身的元数据。一个存储域可以由块设备（块存储）组成，也可以由文件系统（文件存储）组成

	分类:
	1. 数据（Data）存储域：保存ovirt虚拟化环境中的所有虚拟机的磁盘镜像。这些磁盘镜像包括安装的操作系统，或由虚拟机产生或保存的数据。数据存储域支持NFS、iSCSI、FCP、GlusterFS或POSIX兼容的存储系统
	1. 导出（Export）存储域：它可以在不同数据中心间转移磁盘镜像和虚拟机模板提供一个中间存储，并可以用来保存虚拟机的备份。导出存储域支持NFS存储。一个导出域可以被多个不同的数据中心访问，但它同时只能被一个数据中心使用
	1. ISO存储域：用来存储ISO文件（也称为镜像，它是物理的CD或DVD的代表）

	导出（Export）存储域和ISO存储域将被数据（Data）存储域取代.

	自动恢复（激活）机制:
	1. 主机根据其所在数据中心中的存储域元数据信息来监测存储域，当该数据中心中的所有主机都报告某个存储域无法访问时，这个存储域会被认定为“不活跃”.
	1. ovirt-engine监测到某个存储域不活跃时并不会断开与它的连接，而是会认为这可能是一个临时的网络故障导致的，engine会每隔5分钟尝试重新激活任何不活跃的存储域
1. 磁盘镜像存储分配策略

	1. 预分配存储（Preallocated Storage）: 虚拟磁盘镜像所需要的所有存储空间在虚拟机创建前就需要被完全分配

		因为在进行写操作时不需要进行磁盘空间分配的动作，所以预分配存储策略有更好的写性能. 但是，预分配存储的大小不能被扩展，这就失去了一些灵活性。另外，它也会降低ovirt-engine进行存储“over-commitment”的能力。预分配存储策略适用于需要大量I/O操作（特别是写操作比较频繁时），并对存储速率有较高要求的虚拟机，一般情况下，作为应用服务器的虚拟机推荐使用预分配存储策略。

	1. 稀疏分配存储（Sparsely Allocated Storage）, 即存储精简配置策略. 在创建虚拟机的时候，为虚拟磁盘镜像设定一个存储空间上限，而磁盘镜像在开始时并不使用任何 存储域中的存储空间。当虚拟机需要向磁盘中写数据时，磁盘会从存储域中获得一定的存储空间（默认为1G），当磁盘数据量达到所设置的磁盘空间上限时将不会再为虚拟磁盘增加容量
1. 元数据版本

	当前使用的是V3版本: 适用于NFS、GlusterFS、POSIX、iSCSI和FC存储域.

	ovirt-engine在启动时会将每个数据中心中的存储域配置信息下发给各个主机上的VDSM实例，VDSM分根据接收到的存储域配置信息将相关存储域挂载至主机
1. SPM

	能修改存储域元数据的主机就是SPM主机, 机制是`一人写，多人读`. 它通过SPM主机选举机制产生

## 部署
ref:
- [oVirt4.4本地存储架构部署教程（v4.4.3）（Engine独立部署）](https://www.cnovirt.com/archives/2807)
- [oVirt4.4本地存储架构部署教程（v4.4.3）（HostedEngine方式）](https://www.cnovirt.com/archives/2851)

> [官方repo](https://resources.ovirt.org/pub)里centos 7最多支持到ovirt 4.3

> dnf -y install https://resources.ovirt.org/pub/yum-repo/ovirt-release44.rpm

使用oVirt Node iso安装好后可使用cockpit https://ip:9090 by root 进入管理界面

> hosted-engine: 管理端高可用模式, 即将engine作为一台虚拟机漂在所有的node节点上， 如果当engine所在的node节点损坏, 那么会自动迁移到另一个node节点上

> oVirt 引擎包括一个数据仓库，用于收集有关主机、虚拟机和存储的监控数据

[Installing oVirt as a standalone Manager with local databases](https://www.ovirt.org/documentation/installing_ovirt_as_a_standalone_manager_with_local_databases/index.html), 该方式添加host报错:
```bash
# 停止firewalld/selinux
# dnf install -y centos-release-ovirt45 # 官方repo最高就ovirt44, [官方文档里支持45](https://www.ovirt.org/download/)
# dnf install -y ovirt-openvswitch # 见[Installing on RHEL or derivatives](https://www.ovirt.org/download/install_on_rhel.html)
# --- [Enabling the oVirt Engine Repositories](https://www.ovirt.org/documentation/installing_ovirt_as_a_standalone_manager_with_local_databases/index.html), Common procedure valid for both 4.4 and 4.5 on **Enterprise Linux 8 only**
# dnf module -y enable javapackages-tools
# dnf module -y enable pki-deps
# dnf module -y enable postgresql:12
# dnf module -y enable mod_auth_openidc:2.3
# dnf module -y enable nodejs:14
# dnf distro-sync --nobest # 不加`--nobest`可能更安全
# dnf upgrade --nobest
# --- 正式开始
# vim /etc/hosts # 在所有需要访问ovirt-engine的机器上追加`192.168.88.152 engine.myovirt.com`, 并ping test一下
# hostnamectl set-hostname engine.myovirt.com # 给相应的ovirt role node配置正确的hostname
# dnf install -y python3.11-pip # 见ovirt-engine添加host报错
# python3.11 -m pip install netaddr
# dnf install ovirt-engine
# engine-setup # 下面是配置option
Configure Cinderlib integration (Currently in tech preview) (Yes, No) [No]:
Configure Engine on this host (Yes, No) [Yes]:
Configure ovirt-provider-ovn (Yes, No) [Yes]:
Configure WebSocket Proxy on this host (Yes, No) [Yes]:
Configure Data Warehouse on this host (Yes, No) [Yes]
Configure Keycloak on this host (Yes, No) [Yes]:
Configure VM Console Proxy on this host (Yes, No) [Yes]:
Configure Grafana on this host (Yes, No) [Yes]:
Host fully qualified DNS name of this server [engine.myovirt.com]: # 建议配置正确, 如果使用默认的`localhost.localdomain`, 则部署成功后登入ovirt-engine会报`用于访问系统的 FQDN 不是一个有效的引擎 FQDN。您需要使用一个引擎 FQDN 或一个引擎备选的 FQDN 来访问系统。`
Do you want Setup to configure the firewall? (Yes, No) [Yes]:
Firewall manager to configure (firewalld): firewalld
Where is the DWH database located? (Local, Remote) [Local]:
Would you like Setup to automatically configure postgresql and create DWH database, or prefer to perform that manually? (Automatic, Manual) [Automatic]:
Where is the Keycloak database located? (Local, Remote) [Local]:
Would you like Setup to automatically configure postgresql and create Keycloak database, or prefer to perform that manually? (Automatic, Manual) [Automatic]:
Where is the Engine database located? (Local, Remote) [Local]:
Would you like Setup to automatically configure postgresql and create Engine database, or prefer to perform that manually? (Automatic, Manual) [Automatic]:
Engine admin password: password
Use weak password? (Yes, No) [No]: Yes
Application mode (Virt, Gluster, Both) [Both]:
Use Engine admin password as initial keycloak admin password (Yes, No) [Yes]:
Default SAN wipe after delete (Yes, No) [No]:
Organization name for certificate [myovirt.com]:
Do you wish to set the application as the default page of the web server? (Yes, No) [Yes]:
Do you wish Setup to configure that, or prefer to perform that manually? (Automatic, Manual) [Automatic]:
(1, 2)[1]:
Use Engine admin password as initial Grafana admin password (Yes, No) [Yes]:
...
--== CONFIGURATION PREVIEW ==--
         
          Application mode                        : both
          Default SAN wipe after delete           : False
          Host FQDN                               : engine.myovirt.com
          Firewall manager                        : firewalld
          Update Firewall                         : True
          Set up Cinderlib integration            : False
          Configure local Engine database         : True
          Set application as default page         : True
          Configure Apache SSL                    : True
          Keycloak installation                   : True
          Engine database host                    : localhost
          Engine database port                    : 5432
          Engine database secured connection      : False
          Engine database host name validation    : False
          Engine database name                    : engine
          Engine database user name               : engine
          Engine installation                     : True
          PKI organization                        : myovirt.com
          Set up ovirt-provider-ovn               : True
          DWH installation                        : True
          DWH database host                       : localhost
          DWH database port                       : 5432
          DWH database secured connection         : False
          DWH database host name validation       : False
          DWH database name                       : ovirt_engine_history
          Configure local DWH database            : True
          Grafana integration                     : True
          Grafana database user name              : ovirt_engine_history_grafana
          Keycloak database host                  : localhost
          Keycloak database port                  : 5432
          Keycloak database secured connection    : False
          Keycloak database host name validation  : False
          Keycloak database name                  : ovirt_engine_keycloak
          Keycloak database user name             : ovirt_engine_keycloak
          Configure local Keycloak database       : True
          Configure VMConsole Proxy               : True
          Configure WebSocket Proxy               : True
         
          Please confirm installation settings (OK, Cancel) [OK]:OK
# --- engine-setup开始自动部署, 部署完成会看到`--== SUMMARY ==-...[ INFO  ] Execution of setup completed successfully`
...
Web access for Keycloak Administration Console is enabled at:
              https://engine.myovirt.com/ovirt-engine-auth/admin
Web access is enabled at:
              http://engine.myovirt.com:80/ovirt-engine
              https://engine.myovirt.com:443/ovirt-engine
...
[ INFO  ] Execution of setup completed successfully
```

访问`https://engine.myovirt.com/ovirt-engine-auth/admin`, 用admin
访问`https://engine.myovirt.com/ovirt-engine`, 用admin@ovirt
访问前将ovirt ca证书导入浏览器(见FAQ), 否则很多web操作都可能出问题.

admin portal登入后显示数据中心`Default`未初始化, 开始配置:
1. 计算->数据中心, 新建数据中心"mydc"

	存储类型选本地,  即kvm主机节点的虚拟磁盘保存在本地
1. 计算->集群, 配置集群mycluster

	数据中心选mydc
	CPU架构选x86_64
	防火墙类型选firewalld: host > centos7选firewall
1. 计算->主机, 配置主机node151

	主机集群选mycluster
	ip: 192.168.88.151
	密码: xxx

	不配置"电源管理"

	确定后, 看到主机状态是"Installing", 查看该主机的"事件"tab页, 可看到具体安装进度, 如果遇到错误, 还可根据event里的部署log检查具体错误.

	解决错误后, 选择主机列表右上角的"安装"->"重新安装"即可

	可能的错误:
	1. [`The ipaddr filter requires python's netaddr be installed on the ansible controller`](https://github.com/oVirt/ovirt-ansible-collection/issues/695)

		```
		# 在ovirt-engine端执行(ovirt-engine好像使用py3.11, 当时`pip3 install netaddr`提示已安装netaddr)
		# dnf install python3.11-pip
		# python3.11 -m pip install netaddr
		```
	1. `Fail host deploy if firewall type is iptables for hosts other than CentOS 7/ RHEL 7`

		修改集群防火墙类型为firewalld, 如果集群中存在host, 此时要先删除host再修改集群
	1. [`Task If output plugin is elasticsearch, validate host address is set failed to execute`/`'elasticsearch_host' is undefined`](https://github.com/oVirt/ovirt-engine/issues/895)


		考虑到DWH是收集监控信息, 在重新部署engine-setup并禁用DWH和Grafana(与DWH关联)后, 重新添加host还是报该错误

		> 禁用DWH和Grafana后, admin portal的dashboard将没有统计内容

		网上检索, 发现[ovirt-engine-metrics/roles/ovirt_initial_validations/defaults/main.yml](https://github.com/oVirt/ovirt-engine-metrics/blob/master/roles/ovirt_initial_validations/defaults/main.yml), output plugin正是elasticsearch.

		解决: 见`https://github.com/oVirt/ovirt-engine-metrics/pull/35/files`, 按PR更新`roles/ovirt_initial_validations/tasks/check_logging_collectors.yml`即可
1. 存储->存储域, 新建images

	数据中心: mydc
	存储类型: 主机本地
	主机: node151
	路径: /images, 已创建(mkdir /images && chown -R vdsm:kvm /images/), 挂载路径需要唯一

	注意: 从v4.4开始, ovirt计划废弃iso域, 而是将iso上传到数据域, 具体见[oVirt上传ISO镜像的方法](https://www.cnovirt.com/archives/1190), 该文章也提供了上传iso到iso域的方法
1. 存储->存储域, 新建disks

	数据中心: mydc
	存储类型: 主机本地
	主机: node151
	路径: /disks, 已创建(mkdir /images && chown -R vdsm:kvm /disks/)
1. 计算->虚拟机, 新建test

	创建磁盘是不启用"启用增量备份", 那么其格式是raw
1. 用novnc访问 vm. novnc无法访问时见FAQ

Hosted Engine部署:
1. 可以使用[cnovirt-node-ng-installer iso](https://www.cnovirt.com/%e5%ae%89%e8%a3%85%e5%8c%85%e4%b8%8b%e8%bd%bd), 它已集成了engine appliance rpm包. 因为使用官方最新版ovirt-node-ng-installer-4.5.4-2022120615.el8.iso部署过程中node重启了导致部署失败
1. 配置hosts和hostname

	```bash
	echo "192.168.88.151 node.myovirt.com" >> /etc/hosts
	hostnamectl set-hostname node.myovirt.com
	```
1. 准备nfs, 见[oVirt4.4本地存储架构部署教程（v4.4.3）（HostedEngine方式）](https://www.cnovirt.com/archives/2851)

	```bash
	# mkdir /data/images/nfs
	# chown vdsm:kvm /data/images/nfs
	# vi /etc/exports
	```
1. VM

	Engine VM FQDN: engine FQDN
	VM IP Address: engine ip
	Host FQDN: current virt node FQDN
1. Engine

	Admin Portal Password: xxx
1. 配置nfs: 192.168.88.151:/data/images/nfs
1. 导入ca证书并访问`https://engine.myovirt.com/ovirt-engine`, 用admin@internal

## [api](https://www.ovirt.org/documentation/)
- [python-ovirt-engine-sdk4](https://github.com/oVirt/python-ovirt-engine-sdk4)

	使用增量备份时需要在vm disks上开启"启用增量备份", 否则备份时没有生成checkpoint

	examples:
	```
	# ./list_vms.py # 获取vm list
	test: 95c839b9-04c6-4f9e-b113-83e0aefcddb1
	test3: 4bb770db-3302-4f7f-846b-f4e12705985c
	test_restore: 24f25298-a7ea-4a59-be34-bb78c276dc4d
	# ./list_vm_disks.py # 获取vm disks by name use "test"
	name: CentOS-8-GenericCloud-8.1.1911
	id: 198668cf-7f01-467a-badc-85396034c1de
	status: ok
	provisioned_size: 10737418240
	# ./list_disk_snapshots.py -c engine1 198668cf-7f01-467a-badc-85396034c1de
	[
	  {
	    "actual_size": 768610304,
	    "format": "cow",
	    "id": "c0586bf2-e409-4e4e-98e1-f62ef967da26",
	    "parent": null,
	    "status": "ok"
	  }
	]
	# ./list_vm_snapshots.py # 获取vms的snaps
	test:Active VM:CentOS-8-GenericCloud-8.1.1911:images
	test3:Active VM:test3_Disk1:disks
	test_restore:Active VM:CentOS-8-GenericCloud-8.1.1911:disks
	# --- 全量备份
	# ./backup_vm.py -c engine1 --debug full --backup-dir bak_dir 95c839b9-04c6-4f9e-b113-83e0aefcddb1
	[   0.0 ] Starting full backup for VM '95c839b9-04c6-4f9e-b113-83e0aefcddb1'
	[   1.3 ] Waiting until backup '1f085fec-8ea5-42ee-bb25-fc8aad52188f' is ready
	[  23.6 ] Created checkpoint '0ec49a4c-06e5-4c08-ad4f-20ecce88b78a'
	[  23.6 ] Downloading full backup for disk '198668cf-7f01-467a-badc-85396034c1de'
	[  23.6 ] Creating backup file 'bak_dir/20231130115919.0ec49a4c-06e5-4c08-ad4f-20ecce88b78a.198668cf-7f01-467a-badc-85396034c1de.full.qcow2'
	[  25.2 ] Image transfer 'dc28ada5-9562-47ce-ba69-4622d666605c' is ready
	[ 100% ] 10.00 GiB, 51.10 s, 200.41 MiB/s
	[  76.3 ] Finalizing image transfer
	[  81.5 ] Download completed successfully
	[  81.5 ] Finalizing backup
	[ 136.5 ] Full backup '1f085fec-8ea5-42ee-bb25-fc8aad52188f' completed successfully
	# --- 增量备份
	# ./backup_vm.py -c engine1 --debug incremental --backup-dir bak_dir --from-checkpoint-uuid 0ec49a4c-06e5-4c08-ad4f-20ecce88b78a 95c839b9-04c6-4f9e-b113-83e0aefcddb1
	[   0.0 ] Starting incremental backup for VM '95c839b9-04c6-4f9e-b113-83e0aefcddb1'
	[   1.1 ] Waiting until backup 'd31fc6d1-17c7-4bfc-b64d-644974a69323' is ready
	[  13.2 ] Created checkpoint '2045a517-836d-43ba-a1ca-3c8f81fab921'
	[  13.3 ] Downloading incremental backup for disk '198668cf-7f01-467a-badc-85396034c1de'
	[  13.3 ] Creating backup file 'bak_dir/20231130120504.2045a517-836d-43ba-a1ca-3c8f81fab921.198668cf-7f01-467a-badc-85396034c1de.incremental.qcow2'
	[  13.3 ] Using backing file '20231130115919.0ec49a4c-06e5-4c08-ad4f-20ecce88b78a.198668cf-7f01-467a-badc-85396034c1de.full.qcow2'
	[  14.7 ] Image transfer '5329d8c7-cf11-43bb-bc7a-7cc9cef0fca6' is ready
	[ 100% ] 10.00 GiB, 0.59 s, 16.84 GiB/s                                        
	[  15.3 ] Finalizing image transfer
	[  17.3 ] Download completed successfully
	[  17.3 ] Finalizing backup
	[  69.3 ] Incremental backup 'd31fc6d1-17c7-4bfc-b64d-644974a69323' completed successfully
	# ./backup_vm.py -c engine1 --debug stop 95c839b9-04c6-4f9e-b113-83e0aefcddb1 4cff2a36-ddbb-4737-bbf3-5f903f331686 # `vm_id bakcup_id`, 取消备份, 可解锁备份失败时锁定的disks
	# --- 分步全备
	# ./backup_vm.py -c engine1 --debug start 4bb770db-3302-4f7f-846b-f4e12705985c
	[   0.0 ] Starting full backup for VM '4bb770db-3302-4f7f-846b-f4e12705985c'
	[   1.3 ] Waiting until backup '9b098edf-4f23-4291-84de-38ffae8d59e5' is ready
	[  23.7 ] Created checkpoint 'b91276e2-06ab-461e-bd36-efdec14473ed'
	[  23.8 ] Backup '9b098edf-4f23-4291-84de-38ffae8d59e5' is ready
	#./backup_vm.py -c engine1 --debug download --backup-uuid 9b098edf-4f23-4291-84de-38ffae8d59e5 4bb770db-3302-4f7f-846b-f4e12705985c
	[   0.0 ] Downloading VM '4bb770db-3302-4f7f-846b-f4e12705985c' disks
	[   0.5 ] Downloading full backup for disk '63dbc882-70fc-48cf-8401-9c41b8cd0e06'
	[   0.5 ] Creating backup file './20231129165654.b91276e2-06ab-461e-bd36-efdec14473ed.63dbc882-70fc-48cf-8401-9c41b8cd0e06.full.qcow2'
	[   2.1 ] Image transfer 'd4d6ebc0-b802-4103-babb-6768c171680c' is ready
	[ 100% ] 15.00 GiB, 0.31 s, 48.19 GiB/s
	[   2.4 ] Finalizing image transfer
	[   4.5 ] Download completed successfully
	[   4.6 ] Finished downloading disks
	# ./backup_vm.py -c engine1 --debug stop 4bb770db-3302-4f7f-846b-f4e12705985c 9b098edf-4f23-4291-84de-38ffae8d59e5
	[   0.0 ] Finalizing backup '9b098edf-4f23-4291-84de-38ffae8d59e5'
	[  45.3 ] Backup '9b098edf-4f23-4291-84de-38ffae8d59e5' completed successfully
	#  --- 分步增备
	# ./backup_vm.py -c engine1 --debug start --from-checkpoint-uuid b91276e2-06ab-461e-bd36-efdec14473ed 4bb770db-3302-4f7f-846b-f4e12705985c
	[   0.0 ] Starting incremental backup since checkpoint 'b91276e2-06ab-461e-bd36-efdec14473ed' for VM '4bb770db-3302-4f7f-846b-f4e12705985c'
	[   1.1 ] Waiting until backup 'bd474e84-f44a-4fcd-8cce-327c39f78e6d' is ready
	[  12.3 ] Created checkpoint '088a83ea-e897-4bb6-a7ca-cf1f75c6f7bd'
	[  12.4 ] Backup 'bd474e84-f44a-4fcd-8cce-327c39f78e6d' is ready
	# ./backup_vm.py -c engine1 --debug download --backup-uuid bd474e84-f44a-4fcd-8cce-327c39f78e6d --incremental 4bb770db-3302-4f7f-846b-f4e12705985c
	[   0.0 ] Downloading VM '4bb770db-3302-4f7f-846b-f4e12705985c' disks
	[   0.4 ] Downloading incremental backup for disk '63dbc882-70fc-48cf-8401-9c41b8cd0e06'
	[   0.4 ] Creating backup file './20231129170320.088a83ea-e897-4bb6-a7ca-cf1f75c6f7bd.63dbc882-70fc-48cf-8401-9c41b8cd0e06.incremental.qcow2'
	[   0.4 ] Using backing file '20231129165654.b91276e2-06ab-461e-bd36-efdec14473ed.63dbc882-70fc-48cf-8401-9c41b8cd0e06.full.qcow2'
	[   1.8 ] Image transfer 'ea1f86d7-5f08-40d0-9076-550a0d678bd3' is ready
	[ 100% ] 15.00 GiB, 0.27 s, 56.06 GiB/s
	[   2.1 ] Finalizing image transfer
	[   5.1 ] Download completed successfully
	[   5.2 ] Finished downloading disks
	# ./backup_vm.py -c engine1 --debug stop 4bb770db-3302-4f7f-846b-f4e12705985c bd474e84-f44a-4fcd-8cce-327c39f78e6d
	[   0.0 ] Finalizing backup 'bd474e84-f44a-4fcd-8cce-327c39f78e6d'
	[  56.4 ] Backup 'bd474e84-f44a-4fcd-8cce-327c39f78e6d' completed successfully
	# ./upload_disk.py -c engine1 --disk-spare --sd-name data 20231129170320.088a83ea-e897-4bb6-a7ca-cf1f75c6f7bd.63dbc882-70fc-48cf-8401-9c41b8cd0e06.incremental.qcow2
	```

	> upload_disk.py基于ovirt_imageio.client.upload, 该函数默认backing_chain=True即还原时自动处理qcow2 backing chain.

	> `-c engine1`中的engine1是~/.config/ovirt.conf(from examples/ovirt.conf)的section

	> backup_vm.py的start/download/stop子命令是将备份过程分成了3步处理

	> backup_vm.py增量备份可能失败, 见FAQ.

	> backup_vm.py开始备份后会在`/var/run/vdsm/backup`创建同uuid的目录(backup_vm.py stop后会被清理掉), /var/lib/vdsm/storage/transient_disks也有一些文件(结束后也会被清理掉), 记录见`select * from vm_backups`, `select * from vm_backup_disk_map`, `select * from vm_checkpoints`, `select * from vm_checkpoints_disks_map`

	其他:
	```bash
	# ./list_storage_domains.py # 仅支持按名称查询
	```

增量备份分析:
- cmd_incremental

## 备份
- [**Incremental Backup in oVirt**](https://resources.ovirt.org/site-files/2020/Back_to_the_future-incremental_backup_in_oVirt.pdf)或[Back_to_the_future-incremental_backup_in_oVirt.pdf](/misc/pdf/Back_to_the_future-incremental_backup_in_oVirt.pdf)

	Changed block tracking

	Will be in tech preview in oVirt 4.4 - requires libvirt 6.0.z and qemu 4.2

	优点:
	1. Speed up incremental backup by copying only blocks that changed since the last backup
	1. Speed up full backup by copying only the data extents and skipping zero extents
	1. No need to create and delete a snapshot
	1. Access raw guest data in backup and restore regardless of the underlying disk format and snapshots
	1. imageio client library can upload/download raw/qcow2 images (including the backing files)

	还原增量镜像:
	```bash
	$ qemu-img rebase -u incr-backup-2020-01-14.qcow2 -b incr-backup-2020-01-13.qcow2 -F qcow2
	$ qemu-img rebase -u incr-backup-2020-01-13.qcow2 -b full-backup-2020-01-12.qcow2 -F qcow2
	$ --- incr-backup-2020-01-14.qcow2 就是可用的镜像文件了
	````

	该pdf包含备份过程中的关键打点

- [3.2.5. 使用增加备份和恢复 API 备份和恢复虚拟机](https://access.redhat.com/documentation/zh-cn/red_hat_virtualization/4.4/html/administration_guide/chap-backups_and_migration#backing_up_and_restoring_virtual_machines_using_the_incremental_backup_and_restore_api)
- [vacosta94/VirtBKP](https://github.com/vacosta94/VirtBKP)
- [allwaysoft/ovirtvmbackup](https://github.com/allwaysoft/ovirtvmbackup/blob/master/ovirtvmbackup.py)
- [oVirt/vdsm/incremental-backup.md](https://github.com/oVirt/vdsm/blob/master/doc/incremental-backup.md)
- [Incremental Backup](https://www.ovirt.org/develop/incremental-backup-guide/incremental-backup-guide.html)
- [Incremental Backup](https://www.ovirt.org/develop/release-management/features/storage/incremental-backup.html)
- [ovirt增量备份](https://blog.csdn.net/allway2/article/details/102979449)
- [LWN: QEMU中的数据变动跟踪与差分备份](https://blog.csdn.net/Linux_Everything/article/details/110848435)
- [华为云计算学习：备份之CBT技术](https://www.vinchin.com/blog/vinchin-technique-share-details.html?id=8374)
- [ovirtsdk4.services](https://ovirt.github.io/python-ovirt-engine-sdk4/master/services.m.html#ovirtsdk4.services.VmService.snapshots_service)
- [Hybrid Backup](https://www.ovirt.org/media/Hybrid-backup-v8.pdf)
- [Python SDK 指南](https://access.redhat.com/documentation/zh-cn/red_hat_virtualization/4.4/html-single/python_sdk_guide/index#chap-Python_Examples)

增量备份要求:
1. qcow2 v3即compat 1.1, libvirt 增量备份是基于CBT备份. raw镜像仅支持全备, 因此raw和qcow2 v3混合增量备份时, raw盘还是全备
2. 需要指定backing_file, 否则备份生成的文件, `qemu-img info --backing-chain`查询不到backing-chain, 从而导致还原出的disk image是错误的

官方sdk不支持差分备份, 需要自实现

> `python-ovirt-engine-sdk4-4.6.1/examples/vm_backup.py`基于snapshots_service()仅支持全量备份; 而同目录的`backup_vm.py`基于backup_service支持增量备份

备份还原细节 from bareos ovirt plugin:
1. len(snaps_service.list()) > 1 时不能备份

### 混合备份
```
[Hybrid backup is enabled by default in in oVirt 4.5](https://www.ovirt.org/media/Hybrid-backup-v8.pdf). To disable it globally:
# engine-config -s UseHybridBackup=false
# systemctl restart ovirt-engine
```

## tools
- engine-config: 查看engine配置

## FAQ
### ovirt-engine FQDN使用了`localhost.localdomain`, 访问admin portal报`用于访问系统的 FQDN 不是一个有效的引擎 FQDN。您需要使用一个引擎 FQDN 或一个引擎备选的 FQDN 来访问系统`
ref:
- [使用 oVirt Engine Rename 工具重命名 Manager](https://access.redhat.com/documentation/zh-cn/red_hat_virtualization/4.3/html/administration_guide/chap-utilities#Renaming_the_Manager_with_the_Ovirt_Engine_Rename_Tool)
- [OVIRT取消FQDN访问限制](https://blog.csdn.net/h106140873/article/details/82179888)

	没采用, 还是改FQDN方便点

使用`/usr/share/ovirt-engine/setup/bin/ovirt-engine-rename`修改ovirt-engine的FQDN, 如果修改时报`Host name is not valid: engine.myovirt.com did not resolve into an IP address`, ssh重新登入再ping一下engine.myovirt.com, 成功则再次修改即可. 但修改后还是有问题见"ovirt-engine登入报`Internal Server Error`".

### 修改无效的FQDN后, ovirt-engine登入报`Internal Server Error`, 且`/var/log/httpd/ssl_error_log`报`oidc_authenticate_user: the URL hostname (localhost.localdomain) of the configured OIDCRedirectURI does not match the URL hostname of the URL being accessed (engine.myovirt.com): the "state" and "session" cookies will not be shared between the two!, referer: https://engine.myovirt.com/ovirt-engine/`

用`grep -r "localhost.localdomain" /etc`发现ovirt-engine-rename后很多配置还是使用了修改前的FQDN, 果断重装

### 创建self-hostd时Engine VM FQDN报"localhost is not a valid address"
检查/etc/hosts, 是否已添加相关的engine, node FQDN, 并将相关的ovirt node的hostname设置为相应的FQDN

### ovirt node执行`virsh list`需要账号
```bash
[root@ovirt3 ~]# find / -name libvirtconnection.py 
/usr/lib/python3.6/site-packages/vdsm/common/libvirtconnection.py
[root@ovirt3 ~]# egrep SASL_USERNAME /usr/lib/python3.6/site-packages/vdsm/common/libvirtconnection.py
SASL_USERNAME = "vdsm@ovirt"
[root@ovirt3 ~]# find / -name libvirt_password          
/etc/pki/vdsm/keys/libvirt_password
[root@ovirt3 ~]# cat /etc/pki/vdsm/keys/libvirt_password
shibboleth
```

或使用`virsh -c qemu:///system?authfile=/etc/ovirt-hosted-engine/virsh_auth.conf`

`virsh -c qemu+ssh://root@192.168.88.151/system list --all`需要root密码

### ovirt vm 控制台的原生客户端
`apt install virt-viewer`

### ovirt novnc 访问vm时浏览器报`Something went wrong, connection is closed`, 且ovirt-engine的系统日志报`ovirt-websocket-proxy.py[36643]: ovirt-websocket-proxy[36643] INFO msg:630 handler exception: [SSL: SSLV3_ALERT_CERTIFICATE_UNKNOWN] sslv3 alert certificate unknown (_ssl.c:897)`
从admin portal下载ca cert, 导入浏览器并重启, 再访问其novnc即可

- chrome

	1. 打开设置, 找到"管理证书"
	1. 选中"授权机构", 并导入, 导入时选择信任所有授权项
- firefox
	1. 找到"证书管理器"
	1. 选中"证书颁发机构", 并导入, 导入时选择信任所有授权项

其他ca证书导致的问题:
1. 上传iso报, 报`被系统暂停`. 此时可通过`上传`按钮下拉里的`取消`进行取消上传

### CentOS Linux and Stream cloud images
ref:
- [【Cloud】修改CentOS官方 云镜像的ROOT密码](https://developer.aliyun.com/article/799104)

选"带有BIOS的Q35芯片组", 因为这些qcow2不支持uefi

启动过程中可能卡在"probing edd (edd=off to disable)... ok", 继续等待即可

root密码是未知的, 需要自行修改镜像密码:
```bash
# yum install libguestfs-tools
# virt-customize -a disk/CentOS-7-x86_64-GenericCloud-1511.qcow2 --root-password password:123456
# virt-customize -a disk/CentOS-7-x86_64-GenericCloud-1511.qcow2 --root-password random # 随机密码, 修改过程会打印该密码
# virt-customize -a disk/CentOS-7-x86_64-GenericCloud-1511.qcow2 --password tao:password:taoyuhang # 为其他用户设置密码, 且只覆盖密码, 并不能创建用户
```

### api调用报`ovirtsdk4.AuthError: Error during SSO authentication access_denied : Cannot authenticate user No valid profile found in credentials`
当在api中使用admin账号时应是"admin@ovirt@internalsso"

### python sdk报`error:0A000126:SSL routines::unexpected eof while reading`
sdk.Connection的url写错域名

### 解锁disk
```
Unlock all disks
/usr/share/ovirt-engine/setup/dbutils/unlock_entity.sh -t all

Unlock specific VM disk
/usr/share/ovirt-engine/setup/dbutils/unlock_entity.sh -t vm UUID_OF_VM
```

可以解锁损坏的clone disk操作, 但无法解锁使用ovirt sdk增量备份时锁住的disk, 此时可通过`./backup_vm.py stop`解决, 忘记backup uuid时可从`/var/log/ovirt-engine/engine.log`检索`backup`来获取, 比如找到的`Change VM '95c839b9-04c6-4f9e-b113-83e0aefcddb1' backup '304dbd1b-1074-4d3d-8c9c-d09b1e117cdf' phase from 'STARTING' to 'READY'`


### api log
- /var/log/httpd/ssl_access_log
- /var/log/httpd/ssl_request_log
- /var/log/httpd/ovirt-requests-log
- /var/log/ovirt-engine/engine.log
- `/var/log/vdsm/vdsm.log(.\d+.xz)?`

vdsm loglevel修改: [README.logging](https://github.com/oVirt/vdsm/blob/master/README.logging)

### sdk 增量备份 engine.log 报`Bitmap does not exist`
ref:
- [incremental backup failed: Bitmap does not exist](https://github.com/oVirt/ovirt-engine/issues/896)
- [hybrid backup VM after a snapshot fails - start_nbd_server error=Bitmap does not exist](https://bugzilla.redhat.com/show_bug.cgi?id=2068104)
- [Dirty Bitmaps and Incremental Backup](https://qemu-project.gitlab.io/qemu/interop/bitmaps.html)

错误起点应在vdsm.log.

解决方法: qcow2 版本要求v3即compat 1.1

ps: 在ovirt 4.4/4.5 存储域中新建的disk image(已启用"启用增量备份"), compat已是1.1, 本次出问题的镜像是centos cloud image, 其compat是0.10

### sdk api调用报`SSL certificate problem: unable to get local issuer certificate`
engine_url和ca.pem中的域名信息不匹配

### 打开控制台by novnc报`设置 VM ticket 失败`
使用SPICE访问

### sdk 全备 engine.log 报`Cannot store dirty bitmaps in qcow2 v2 files`
ref:
- [Backup fails with "Cannot store dirty bitmaps in qcow2 v2 files"](https://github.com/abbbi/virtnbdbackup#backup-fails-with-cannot-store-dirty-bitmaps-in-qcow2-v2-files)
- [Features/Qcow3](https://wiki.qemu.org/Features/Qcow3)

原因: qcow2 版本要求v3

查看方法:
1. file xxx.qcow2
2. `qemu-img info`: qcow2 v3的compat应是1.1

修改qcow2 compat: `qemu-img amend -f qcow2 -o compat=1.1 /tmp/test.qcow2`

### 查找vm 对应的磁盘
在存储域下, `grep <disk id>`查找

### backup_vm.py增量备份报`...qemu-img: warning: Could not verify backing image. This may become an error in future versions. Image is not in qcow2 format...RuntimeError("Timeout waiting for qemu-nbd socket")`
情况1:
1. 增量备份的目标目录存在其他文件, 文件名刚好符合backup_vm.py#find_backing_file的匹配导致使用了错误的backing file
2. 增量备份时要验证所依赖的backing file, 该文件不能删除, 即不能删除本次备份所有直接依赖的backing file, 间接依赖的backing file可以删除

	因为client.download调用了qemu-img create创建新qcow2, 其参数需要该backing file.

### engine db 配置
`/etc/ovirt-engine/engine.conf.d/10-setup-database.conf`

```
# su postgres
# psql -s engine
```

### upload_disk.py报`Cannot add Virtual Disk. Disk configuration (COW Preallocated backup-None) is incompatible with the storage domain type.`
上传的disk image时精简的, 需要`--disk-sparse`参数

### python lxml解析disks
namespaces是xpath()时使用的ns

```python3
ovf = lxml.etree.parse("vm.ovf")

namespaces = {
    'ovf': 'http://schemas.dmtf.org/ovf/envelope/1',
    'xsi': 'http://www.w3.org/2001/XMLSchema-instance'
}

disk_elements = ovf.xpath(
    '//Section[@xsi:type="ovf:DiskSection_Type"]/Disk',
    namespaces=namespaces
)

for disk_element in disk_elements:
    # Get disk properties:
    props = {}
    for key, value in disk_element.items():
        key = key.replace('{%s}' % namespaces['ovf'], '')
        props[key] = value
    print(props)
```

或使用绝对路径: `/ovf:Envelope/Section[@xsi:type="ovf:DiskSection_Type"]/Disk`

### export ova报`Operation not permitted abortedcode=100`
ref:
- [[Libguestfs] [PATCH v2 2/3] v2v: ovf: Create OVF more aligned with the standard](https://listman.redhat.com/archives/libguestfs/2018-February/018293.html)

仅允许export正在运行的vm

### export ova xml与vms_service().list()获取的xml格式有差异
ref:
- [ISO/IEC 17203:2017 - Information technology — Open Virtualization Format (OVF) specification](https://webstore.iec.ch/preview/info_isoiec17203%7Bed2.0%7Den.pdf)
- [OVF 与 OVA 区别与转换](https://blog.k8s.li/ovf-to-ova.html)
- [虚拟机包 OVF和OVA的区别](https://blog.51cto.com/wolfgang/1125864)

Centos8.ova是包含vm xml和vm disks的tar包.

```bash
# tar -xvf Centos8.ova
vm.ovf
e0ea67f1-577b-44e1-8c58-15121398beb3
```

e0ea67f1-577b-44e1-8c58-15121398beb3是`/ovf:Envelope/References/File[@ovf:href="e0ea67f1-577b-44e1-8c58-15121398beb3"]`与`/ovf:Envelope/DiskSection/Disk[@ovf:fileRef="e0ea67f1-577b-44e1-8c58-15121398beb3"]`相对应

export ova时disk xml path是`/ovf:Envelope/DiskSection/Disk`且`Disk[@ovf:fileRef]`仅一层, 使用`connection.system_service().vms_service().list(id=xxx, all_content=True).initialization.configuration.data.encode('utf-8')`时path是`/ovf:Envelope/Section[@xsi:type="ovf:DiskSection_Type"]/Disk`且`Disk[@ovf:fileRef]`有2层, 第二层的内容与export ova的`Disk[@ovf:fileRef]`相同. 通过`upload_ova_as_vm_or_template.py#upload_ova_as_vm_or_template.py`看到, xml差异应是types.ConfigurationType的OVA和OVF不同导致的, 可通过`vms_service().list()`时添加`ovf_as_ova=True`来解决, 此时xml内容与export ova的一致

ovf标准使用的是DiskSection, 即export ova时所采用的xml格式.

### ovirt node 4.4.10备份已关机的vm失败
[add_bitmap failed: Operation not permitted](https://github.com/oVirt/ovirt-engine/issues/902)

ovirt-engine-4.5.4-1.el8.noarch + ovirt-node-ng-installer-4.5.4-2022120615.el8.iso 上全/增量均备份正常

原因: vdsmd.service权限是vdsm:kvm, 而[path](https://github.com/oVirt/vdsm/blob/v4.40.100.2/lib/vdsm/storage/outOfProcess.py#L152C24-L152C41)是root:root, 应该是之前调试时误操作导致的.

### api调用报`Error during SSO authentication access_denied : Cannot authenticate user ‘admin@N/A’: No valid profile found in credentials..`
账号错误, 使用了`admin`, 正确应是`admin@internal`