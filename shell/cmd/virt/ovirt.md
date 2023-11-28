# ovirt
- [oVirt虚拟化关键概念、组件与技术原理](http://simiam.com/2018/07/11/oVirt%E8%99%9A%E6%8B%9F%E5%8C%96%E5%85%B3%E9%94%AE%E6%A6%82%E5%BF%B5%E3%80%81%E7%BB%84%E4%BB%B6%E4%B8%8E%E6%8A%80%E6%9C%AF%E5%8E%9F%E7%90%86/)
- [oVirt 架构学习](https://cloud.tencent.com/developer/article/1435899)
- [编译oVirt相关软件包](https://www.hikunpeng.com/document/detail/zh/kunpengcpfs/ecosystemEnable/oVirt/kunpengovirtoe_04_0004.html)
- [Next-Gen Backup & DR Solution for oVirt](https://resources.ovirt.org/site-files/2020/Next-Gen%20oVirt%20Backup%20%26%20DR%20Solution%20-%20vinchin%20-%20official%20presentation.pdf)

oVirt基于KVM，并整合使用了libvirt、gluster、patternfly、ansible等一系列优秀的开源软件，oVirt的定位是替代VMware vsphere，oVirt目前已经成为了企业 虚拟化  环境可选的解决方案，另外相比OpenStack的庞大和复杂，oVirt在企业私有云建设中具备部署和维护使用简单的优势.

oVirt项目中的不同组件主要包含三个部分：
- ovirt-engine: 用来进行管理虚拟机（创建、开关启停）、配置网络和存储等操作
- 一个或多个主机（节点）: 用来运行虚拟机。

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

## 部署
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
# vim /etc/hosts # 在所有需要访问ovirt-engine的机器上追加`192.168.88.152 egnine.myovirt.com`, 并ping test一下
# hostnamectl set-hostname egnine.myovirt.com # 给相应的ovirt role node配置正确的hostname
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
Host fully qualified DNS name of this server [egnine.myovirt.com]: # 建议配置正确, 如果使用默认的`localhost.localdomain`, 则部署成功后登入ovirt-engine会报`用于访问系统的 FQDN 不是一个有效的引擎 FQDN。您需要使用一个引擎 FQDN 或一个引擎备选的 FQDN 来访问系统。`
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
          Host FQDN                               : egnine.myovirt.com
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
              https://egnine.myovirt.com/ovirt-engine-auth/admin
Web access is enabled at:
              http://egnine.myovirt.com:80/ovirt-engine
              https://egnine.myovirt.com:443/ovirt-engine
...
[ INFO  ] Execution of setup completed successfully
```

访问`https://egnine.myovirt.com/ovirt-engine-auth/admin`, 用admin
访问`https://egnine.myovirt.com/ovirt-engine`, 用admin@ovirt

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
	路径: /images, 已创建(mkdir /images && chown -R vdsm:kvm /images/)
1. 存储->存储域, 新建disks

	数据中心: mydc
	存储类型: 主机本地
	主机: node151
	路径: /disks, 已创建(mkdir /images && chown -R vdsm:kvm /disks/)
1. 计算->虚拟机, 新建test
1. 用novnc访问 vm. novnc无法访问时见FAQ

Hosted Engine部署:
1. VM

	Engine VM FQDN: engine FQDN
	VM IP Address: engine ip
	Host FQDN: current virt node FQDN
1. Engine

	Admin Portal Password: xxx

未成功, ovirt node未知重启了.

## 备份
- [vacosta94/VirtBKP](https://github.com/vacosta94/VirtBKP)
- [allwaysoft/ovirtvmbackup](https://github.com/allwaysoft/ovirtvmbackup/blob/master/ovirtvmbackup.py)
- [oVirt/vdsm/incremental-backup.md](https://github.com/oVirt/vdsm/blob/master/doc/incremental-backup.md)
- [Incremental Backup](https://www.ovirt.org/develop/release-management/features/storage/incremental-backup.html)

官方sdk不支持差分备份, 需要自实现

## FAQ
### ovirt-engine FQDN使用了`localhost.localdomain`, 访问admin portal报`用于访问系统的 FQDN 不是一个有效的引擎 FQDN。您需要使用一个引擎 FQDN 或一个引擎备选的 FQDN 来访问系统`
ref:
- [使用 oVirt Engine Rename 工具重命名 Manager](https://access.redhat.com/documentation/zh-cn/red_hat_virtualization/4.3/html/administration_guide/chap-utilities#Renaming_the_Manager_with_the_Ovirt_Engine_Rename_Tool)
- [OVIRT取消FQDN访问限制](https://blog.csdn.net/h106140873/article/details/82179888)

	没采用, 还是改FQDN方便点

使用`/usr/share/ovirt-engine/setup/bin/ovirt-engine-rename`修改ovirt-engine的FQDN, 如果修改时报`Host name is not valid: egnine.myovirt.com did not resolve into an IP address`, ssh重新登入再ping一下egnine.myovirt.com, 成功则再次修改即可. 但修改后还是有问题见"ovirt-engine登入报`Internal Server Error`".

### 修改无效的FQDN后, ovirt-engine登入报`Internal Server Error`, 且`/var/log/httpd/ssl_error_log`报`oidc_authenticate_user: the URL hostname (localhost.localdomain) of the configured OIDCRedirectURI does not match the URL hostname of the URL being accessed (egnine.myovirt.com): the "state" and "session" cookies will not be shared between the two!, referer: https://egnine.myovirt.com/ovirt-engine/`

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