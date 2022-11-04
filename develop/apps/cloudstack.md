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
### 创建vm
1. `server/src/main/java/com/cloud/api/ApiServlet.java#processRequest(final HttpServletRequest req, final HttpServletResponse resp)`

	处理所有请求. 具体处理登录、登出请求其他请求进行转发, 进入下一步 by `apiServer.handleRequest`

1. `server/src/main/java/com/cloud/api/ApiServer.java#handleRequest(final Map params, final String responseType, final StringBuilder auditTrailSb) throws ServerApiException`

	收到创建VM的请求，解析参数根据参数“command” 实例化处理该API的类，并将请求参数赋值到该类创建的对象


	对于创建VM，实例化的类为: `org.apache.cloudstack.api.command.user.vm.DeployVMCmd(api/src/main/java/org/apache/cloudstack/api/command/user/vm/DeployVMCmd.java)`, 因为它继承于BaseAsyncCreateCustomIdCmd, 所以需要首先调用DeployVMCmd的create函数，再调用execute函数.

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