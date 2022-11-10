# cloudstack

## 概念
CloudStack部署架构:
- Zone：Zone 对应于现实中的一个数据中心，它是 CloudStack 中最大的一个单元。
- Pod：Pod 对应着一个机架。同一个 pod 中的机器在同一个子网(网段)中。
- Cluster：Cluster 是多个主机组成的一个集群。同一个 cluster 中的主机有相同的硬件，相同的 Hypervisor，和共用同样的存储。同一个 cluster 中的虚拟机，可以实现无中断服务地从一个主机迁移到另外一个上。
- Host：Host 就是运行虚拟机(VM)的主机。

	即从包含关系上来说，一个 zone 包含多个 pod，一个 pod 包含多个 cluster，一个 cluster 包含多个 host。概念上zone>pod>cluster>host.
- Primary storage：一级存储与 cluster 关联，它为该 cluster 中的主机的全部虚拟机提供磁盘卷。一个 cluster 至少有一个一级存储，且在部署时位置要临近主机以提供高性能。
- Secondary storage：二级存储与 zone 关联，它存储模板文件，ISO镜像和磁盘卷快照。

## source
- packaging

	构建rpm/deb

## 流程分析
ref:
- [Apache CloudStack API Documentation](https://cloudstack.apache.org/api/apidocs-4.17/)
- [CloudStack常用API](https://blog.csdn.net/dandanfengyun/article/details/107506781)
- [CloudStack High Availability源码分析](https://codeantenna.com/a/Keg1YuU7wg)

### 创建vm
ref:
- [CloudStack 创建VM 源码流程分析](https://blog.csdn.net/zt689/article/details/38119715)

1. `server/src/main/java/com/cloud/api/ApiServlet.java#processRequest(final HttpServletRequest req, final HttpServletResponse resp)`

	处理所有请求. 具体处理登录、登出请求其他请求进行转发, 进入下一步 by `apiServer.handleRequest`. 这里入参是[deployVirtualMachine](https://cloudstack.apache.org/api/apidocs-4.17/apis/deployVirtualMachine.html).

	req:
	```log
	http://localhost:8080/client/api
	?command=deployVirtualMachine
	&response=json
	&serviceofferingid=8dc2f63f-eab1-4fb8-b070-c73ca4aa93f8
	&templateid=9f3741df-8bd5-4b4b-8f24-36fcbc0f72b7
	&zoneid=16c8ac5d-84a7-4070-adec-5a0aaa17a686&name=CentOS7Test1
	&ipaddress=192.168.199.51
	&apiKey=zsElxq-OfzPyJ8LJcrv7V5vKO74GP6UM-3TuniJoGflIHMWDAfby2rGYypFnL7XROnzxuLBHNFKk4uQCvP1AMA
	```

	resp:
	```json
	{
	    "deployvirtualmachineresponse": {
	        "id": "f3269f7d-ea10-494d-bc1d-a460baf89ca7",
	        "jobid": "993f113e-231d-40ea-ad04-a5af01b4611c"
	    }
	}
	```

	浏览(listVirtualMachines, 这里与上面创建的vm不是同一个, 这里仅为了展示vm schema)返回的是:
	```json
	{
	    "listvirtualmachinesresponse": {
	        "count": 1,	
	        "virtualmachine": [
	            {
	                "id": "1a9a8606-921e-489b-9d27-6e22d4400b0b",
	                "name": "CentOS7Test",
	                "displayname": "CentOS7Test",
	                "account": "admin",
	                "userid": "d86f1bf0-ca2b-11ea-abc8-000c29972b82",
	                "username": "admin",
	                "domainid": "ab8efe85-ca2b-11ea-abc8-000c29972b82",
	                "domain": "ROOT",
	                "created": "2020-07-23T10:12:21+0800",
	                "state": "Running",
	                "haenable": false,
	                "groupid": "dcf670db-5001-459f-bf9a-6b8f1b7628fd",
	                "group": "test",
	                "zoneid": "2d372d51-fca5-409a-a98c-d09edb943138",
	                "zonename": "Zone1",
	                "hostid": "b9c1e868-81bf-4f2c-a28d-8876dafef3d1",
	                "hostname": "agent",
	                "templateid": "e3d7f9f9-6209-4f28-a16b-f48d544c44e6",
	                "templatename": "CentOS6MiniModel",
	                "templatedisplaytext": "CentOS6MinimalModel",
	                "passwordenabled": false,
	                "serviceofferingid": "bebcf7d6-b8e3-454c-b600-d2b2b3c50b4f",
	                "serviceofferingname": "Medium Instance",
	                "cpunumber": 1,
	                "cpuspeed": 1000,
	                "memory": 1024,
	                "cpuused": "6.12%",
	                "networkkbsread": 13,
	                "networkkbswrite": 0,
	                "diskkbsread": 45553,
	                "diskkbswrite": 1365,
	                "memorykbs": 1048576,
	                "memoryintfreekbs": 524288,
	                "memorytargetkbs": 524288,
	                "diskioread": 1931,
	                "diskiowrite": 258,
	                "guestosid": "c583de96-ca2b-11ea-abc8-000c29972b82",
	                "rootdeviceid": 0,
	                "rootdevicetype": "ROOT",
	                "securitygroup": [
	                    {
	                        "id": "d86f5621-ca2b-11ea-abc8-000c29972b82",
	                        "name": "default",
	                        "description": "Default Security Group",
	                        "account": "admin",
	                        "ingressrule": [],
	                        "egressrule": [],
	                        "tags": [],
	                        "virtualmachineids": []
	                    }
	                ],
	                "nic": [
	                    {
	                        "id": "945db421-96c5-4d40-9037-6064bebcf145",
	                        "networkid": "7f9883ab-d15d-45b1-a9c1-0aa9c1f7f230",
	                        "networkname": "defaultGuestNetwork",
	                        "netmask": "255.255.255.0",
	                        "gateway": "192.168.199.1",
	                        "ipaddress": "192.168.199.47",
	                        "broadcasturi": "vlan://untagged",
	                        "traffictype": "Guest",
	                        "type": "Shared",
	                        "isdefault": true,
	                        "macaddress": "1e:00:b7:00:00:1c",
	                        "secondaryip": [],
	                        "extradhcpoption": []
	                    }
	                ],
	                "hypervisor": "KVM",
	                "instancename": "i-2-8-VM",
	                "details": {
	                    "keyboard": "us",
	                    "memoryOvercommitRatio": "2.0",
	                    "Message.ReservedCapacityFreed.Flag": "false",
	                    "cpuOvercommitRatio": "2.0"
	                },
	                "affinitygroup": [],
	                "displayvm": true,
	                "isdynamicallyscalable": false,
	                "ostypeid": "c583de96-ca2b-11ea-abc8-000c29972b82",
	                "tags": []
	            }
	        ]
	    }
	}
	```

1. `server/src/main/java/com/cloud/api/ApiServer.java#handleRequest(final Map params, final String responseType, final StringBuilder auditTrailSb) throws ServerApiException`

	收到创建VM的请求，解析参数根据参数“command” 实例化处理该API的类，并将请求参数赋值到该类创建的对象


	对于创建VM，实例化的类为: `org.apache.cloudstack.api.command.user.vm.DeployVMCmd(api/src/main/java/org/apache/cloudstack/api/command/user/vm/DeployVMCmd.java)`, 因为它继承于BaseAsyncCreateCustomIdCmd, 所以需要首先调用DeployVMCmd的create函数，再调用execute函数.

	> 实例化类的对应关系: command = `@APICommand(name = "deployVirtualMachine",...)`
1. `org.apache.cloudstack.api.command.user.vm.DeployVMCmd.create()`

	检查zone、serviceoffering、templateId、diskoffering、zone是否启用了本地存储、获取ipv4,、ipv6地址，判断是基本zone，还是高级zone.
1. etc.

## FAQ
1. `Could not transfer metadata net.juniper.contrail:juniper-contrail-api:1.0-SNAPSHOT/maven-metadata.xml from/to maven-default-http-blocker (http://0.0.0.0/): transfer failed for http://0.0.0.0/net/juniper/contrail/juniper-contrail-api/1.0-SNAPSHOT/maven-metadata.xml`
Maven 3.8.1就禁止了所有HTTP协议的Maven仓库, 解决方法:

1. 添加blocked
```bash
<mirror>
    <id>maven-default-http-blocker</id>
    <mirrorOf>external:http:*</mirrorOf>
    <name>Pseudo repository to mirror external repositories initially using HTTP.</name>
    <url>http://0.0.0.0/</url>
    <blocked>true</blocked>
</mirror>
```
2. 使用https
3. 注释掉 $MAVEN_HOME/conf/settings.xml 中的拦截标签`<mirror><id>maven-default-http-blocker</id>...</mirror>`