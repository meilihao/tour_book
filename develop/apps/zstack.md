# zstack 4.1.0
参考:
- [zstack - 任杰的博客](https://blog.csdn.net/Snail_Ren?type=blog)
- [zstack工作流程分析](https://github.com/SnailJie/ZstackWithComments)
- [zstack二次开发](https://blog.csdn.net/Snail_Ren/article/details/70318783)

zstack是一款数据中心的资源管理软件.

## 构建
env: centos 7.6

```bash
$ sudo apt install openjdk-8-jdk
$ export JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64
$ sudo apt install maven
$ sudo vim /etc/maven/settings.xml # 在`<mirrors>`配置节点加入如下`<mirror>`内容
  <mirrors>
    <mirror>
      <id>aliyun</id>
      <mirrorOf>*</mirrorOf>
      <name>Aliyun Mirror</name>
      <url>https://maven.aliyun.com/repository/public</url>
    </mirror>
  </mirrors>
$ mkdir ~/test/zstack-test && cd ~/test/zstack-test
$ wget https://downloads.apache.org/tomcat/tomcat-8/v8.5.65/bin/apache-tomcat-8.5.65.tar.gz
$ tar -xf apache-tomcat-8.5.65.tar.gz
$ export CATALINA_HOME=~/test/zstack-test/apache-tomcat-8.5.65
$ vim $CATALINA_HOME/conf/tomcat-users.xml # 添加tomcat用户
<role rolename="admin-gui"/>
<role rolename="manager-gui"/>
<role rolename="manager-jmx"/>
<role rolename="manager-script"/>
<role rolename="manager-status"/>
<user username="tomcat" password="tomcat" roles="admin-gui,manager-gui,manager-script,manager-jmx,manager-status"/>
$ git clone --depth 1 git@github.com:zstackio/zstack.git
$ git clone --depth 1 git@github.com:zstackio/zstack-utility.git
$ git clone --depth 1 git@github.com:zstackio/zstack-vyos.git
$ git clone --depth 1 https://github.com/zstackio/zstack-dashboard
# --- 编译zstack-dashboard
cd ~/test/zstack-test/zstack-dashboard
sudo apt install python-pip
python2 setup.py sdist # 打包结果在dist/zstack_dashboard-0.7.tar.gz
cp zstack-dashboard/dist/zstack_dashboard-*.tar.gz ~/test/zstack-test/apache-tomcat-8.5.65/webapps/zstack/WEB-INF/classes/tools/ # 拷贝到tomcat里
$ /etc/init.d/zstack-dashboard stop # 关闭dashboard
$ zstack-ctl install_ui # 安装dashboard
$ /etc/init.d/zstack-dashboard start # 启动dashboard
# --- 编译zstack
$ cd ~/test/zstack-test/zstack/build
$ ./deploydb.sh root password localhost 3306
$ vim ../conf/zstack.properties # 更新db参数, 在文件部署到tomcat时路径是webapps/zstack/WEB-INF/classes/zstack.properties
DB.user=root
DB.password=password
$ ./deploy.sh
# --- 下面命令仅说明, 无法用于实际
$ cd ~/test/zstack-test/zstack-utility/zstackbuild
$ vim build.properties # 因为多个组件源码的缺失, 不能通过zstackbuild来打包
$ vim build.xml
$ ant -Dzstack_build_root=/root/zstack all-in-one # 编译结束后, 会在 zstack-utility/zstackbuild/target 目录下产生 zstack-installer*.bin, 即为安装包
```

## 架构
参考:
- [how-to-learn-zstack-code](https://jingtyu.gitbooks.io/how-to-learn-zstack-code)

- zstack使用Java编写, 是ZStack的核心, 负责IaaS各种资源管理调度和消息通讯
- zstack-utility目前主要使用Python编写, 包含ZStack的各种终端代理和其他工具. 这些终端代理负责接收来自ZStack核心的消息并执行对应的操作, 例如和Libvirt通讯来管理VM的生命周期、各种存储（例如Ceph, iSCSI, SFTP）的管理、 虚拟路由器里管理VM的IP地址等等. 除了终端代理工具外, 这个软件仓库还包含了ZStack其他的工具, 例如ZStack的编辑打包工具、 ZStack安装程序、ZStack命令行工具、ZStack管控工具等等
- zstack-dashboard使用JavaScript编写, 是ZStack的图形界面, **已不再更新**

zstack目录:
- build	用于Java部分的编译、打包、部署等
- compute	有关计算资源的操作, 各类工作流. 有放置、集群、主机、虚拟机、区域的相关操作. 如虚拟机模块：提供具体的虚拟机在创建、管理过程中, 需要做哪些步骤, 垃圾如何回收等
- conf	配置文件及SQL文件的放置;Spring Service配置存放;持久化文件配置
- configuration	模板管理. 创建管理虚拟机模板、云盘模板、虚拟路由模板. 同时可以配置全局属性
- console	疑似控制模块, 用于控制消息的管理console angent
- core	核心模块. 实现系统的核心功能——包括数据库、消息总线、工作流实现等等
- header	消息以及Entity的定义
- identity	账户管理. 用户的登入登出操作、sessionToken的获取等等
- image	镜像的管理. 处理镜像的消息, 如镜像的添加、查询等等
- network	网络管理, 包含二层网络、三层网络的创建、管理等操作
- plugin	顾名思义. 其中不少组件都以插件化开发, 提供较高的灵活性
- portal	消息总分发入口, 包括ManageMent Node的管理、根据消息分发. 在接受到消息体后, 对消息中的信息进行验证, 调用对应消息所需的消息预处理. ManagementNodeManagerImpl中實現了管理節點的啓動流程
- search	查询操作处理模块. 与数据库相关的操作
- sdk	测试库使用的SDK
- simulator	对于测试库支持的又一模块, 主要用户simulator agent的行为.  此模块比较独立, 其他模块没有调用这个模块的代码
- tag	Tag系统, 负责Tag的生成和管理
- testlib	测试库
- test	测试模块
- tool	工具包. 目前仅仅支持了doc生成
- utils	代码中使用的工具类
- 其他	功能实现模块

zstack-utility目录
- agentcli	也是一个提供CLI服务的模块. 根据描述, 应该是用于agent测试用的, 只识别写在文件中的命令, 应该是用于执行测试命令脚本
- apibinding	为zstackcli服务, 对cli命令相关API转换提供服务, 生成请求格式, 往Tomcat发出REST请求
- appliancevm	开启HttpServer, 监听7759端口. 主要负责防火墙的规则设置
- buildsystem	zstack系统的安装代码, 负责系统的安装流程, 安装后的服务启动
- cephbackupstorage	
- cephprimarystorage	
- consoleproxy	开启HttpServer, 监听7758端口. 守护进程, 包含VNC连接. 在ZStack运行的过程中会一直开启consoleproxy守护进程
- fusionstorbackupstorage	
- fusionstorprimarystorage	
- imagestorebackupstorage	
- installation	
- iscsifilesystemagent	
- kvmagen kvmagent,用于部署到子节点上进行实际的kvm操作的代理. 在这个文件夹下面包含很多plugin, 在虚拟机创建的时候会依次载入这些plugin
- puppets	
- setting	
- sftpbackupstorage	
- virtualrouter	
- zstackbuild	ZStack build目录. 进入此目录后进行All-in-One包编译
- zstackcli	zstackcli命令管理工具
- zstackctl	zstackctl命令管理工具
- zstacklib	提供运行中所需要的一些工具. 守护进程、Http服务等

zstack-dashboard目录:
- zstack_dashboard    web所有的文件
- zstacl_dashboard/static 静态资源，包含页面的html文件以及js、image等文件
- zstack_dashboard/web.py 后台文件，对消息的收发
- zstack_dashboard.sh 部署脚本

安装目录:
- /usr/local/zstack

    ZStack的主要安装路径(All-In-One安装路径)，为ManagementNode安装保存路径，即此路径会安装在ManagementNode上。此路径包含ZStack的Tomcat服务、Ansible服务、安装源码包、垃圾收集等等。

    ManagementNode会不断轮询Agent，来收集Agent信息。当发现有服务down后，会自动重新部署安装Agent节点，替换代码

- /var/lib/zstack

    zstack的Agent运行时环境，即此路径会安装在host节点. 在virtualenv中包含console代理、virtualenv、KVM Agent等等.

    KVM Agent使用cherrypy来提供http服务. 即在Agent端提供Http server功能，接收来自ManagementNode的请求，如该agent的状态收集请求、控制请求等

### 总览
zstack其实质上，是一个单进程多线程程序.

在ZStack内部，利用的是消息队列进行的交互
在zstack agent和zstack之间，是利用RESTful进行的交互

源码阅读见: [hello_zstack](https://gitee.com/chenhao/hello_zstack)

### 异步架构
参考:
- [ZStack--可拓展性的秘密武器1：异步架构](https://mp.weixin.qq.com/s?__biz=MzI0NTc4MzE4Mw==&mid=100000033&idx=1&sn=1a0389051572055b7bb9a23f2acb5e21)

Iaas的核心应该做的是管控层, 而不是数据层. 故ZStack仅仅也是做出一些“决策”而已——在设计系统的时候, 应不考虑在这些决策的执行上消耗大量的资源. 在面对大量请求或者“决策”的时候, 如果使用多线程来处理阻塞式IO模型时会遇到一些问题：
- 阻塞模型的吞吐量受到线程池大小的限制
- 创建并使用许多线程会耗费额外的时间用于上下文切换, 影响系统性能

而非阻塞、异步的消息驱动系统可以只运行少量的线程, 并且不阻塞这些线程, 只在需要计算资源时才使用它们. 这大大提高了系统的响应速度, 并且能够更高效地利用系统资源. 

故, ZStack采用了异步架构, 分别由三个部分组成：
- 异步消息
- 异步方法
- 异步HTTP 请求

如果在系统中的一部分采用异步设计, 是不行的. 这样还是会因为同步而没法享受异步带来的“福利”. 故此整个系统都得采用异步架构.

相对的, 开发者们在编写异步代码的时候得格外小心. 

在系统设计中, 异步调用可以减少系统在IO上出现瓶颈的可能性.

### 无状态服务
- [ZStack—可拓展性秘密武器2：无状态的服务](https://mp.weixin.qq.com/s?__biz=MzI0NTc4MzE4Mw==&mid=100000033&idx=2&sn=18c269f0108b0d39d46cd89e48fdf5a1)

在ZStack中, 每一个服务都是独立存在的. 为了方便的管理更多的物理机, ZStack推荐采用集群部署MN. 但这样就会遇到一个问题, 不同MN下面有着不同的几个服务存在, 在这里我们设其为X个服务. 在10个MN部署的情况下, 可能就是10X个服务. 那么在一个资源需要操作时, 我需要发送向对应的MN. 那么如何找到那个MN呢？最直观的想法就是在各个MN中保存相应的“服务表”, 这即是一种状态. 那么在分布式系统中, 采用有状态的服务绝对不是一个好的选择, 它会严重影响系统的扩展性. ZStack巧妙的采用了一致性哈希算法+MQ解决了这个问题. 

这在系统设计中实为是一种使用一致性hash技术的负载均衡


### 无锁架构
参考:
- [ZStack--可拓展性秘密武器3：无锁架构](https://mp.weixin.qq.com/s?__biz=MzI0NTc4MzE4Mw==&mid=100000033&idx=3&sn=9006182bc8c8992e4771a48061d4481b)

解决并发的问题不一定要用显式的锁, 也可以对同一资源做操作的任务做成队列使其串行执行.

### 项目模块化
在Intellij中打开ZStack的代码, 会发现大多数目录底下都会有一个pom.xml文件, ZStack采用了模块化项目. 模块化的好处在工程实践中不言而喻的, 比如：
- 可以在不影响整个系统的情况下替换某个模块
- 开发者只要专心的在自己的模块中工作即可
- 减少系统耦合度, 提高内聚, 减少资源循环依赖, 增强系统框架设计

### 通信 by mq
在ZStack中, 每个功能实现模块都会被称为服务——一个独立的服务. 各个服务之间的通信由MQ来承担. 这就像是传统的CSE, C和E是不耦合的, 通过S来交互. 同样的, 一个服务需要向另一个服务发起调用, 只需往消息总线发送消息, 并指定这个服务ID（Service ID）即可. 如果某个服务的代码需要大量重构或者做成微服务, 只要提供相同的服务并注册到MQ上就可以了. 这就是事件驱动架构（Event Driven Architecture）的一种典型实现.

> CSE：Controller、Service、Entity. 注：称作Domain或者Model都是不专业的. Domain是一个领域对象, 往往在做传统Java软件web开发中, 这些Domain都是贫血模型, 是没有行为的, 或是没有足够的领域模型的行为的, 所以, 以这个理论来讲, 这些Domain都应该是一个普通的entity对象, 并非领域对象, 所以请把包名改为:com.xxx.entity.

### 框架
#### 使用Spring
在代码中, 每当我们New出一个对象时, 这个模块便对这个对象产生了依赖. 当我们需要测试的时候就不得不去Mock它. 当依赖的对象or Field 有成千上万个的时候, 这就是一场灾难了. 代码变得愈发不可测, 坑就越多, 开发者在扩展or维护项目的时候就会愈发的乏力. 这就像是我们之前提到的MQ, 服务1->MQ->服务2, 由于中间隔了一个MQ, 于是服务1和服务2没有必然的关系. 同样的, 从对象1->调用->对象2到对象1->调用->Spring提供的IOC容器->对象2, 这样使对象与对象之间也没有了直接调用关系, 对象1只要知道它要调用的对象实现了其需要的Interface就是可以调用的. 

除了Autowired的正确使用姿势. 在ZStack中, 还有一类很有意思的代码, 一般称之为xxxExtensionPoint. 其本质就是定义一个接口, 然后其实现类作为Bean通过XML注册到IOC中. 在需要使用的时候, 通过Spring获取到所有实现该接口的对象, 调用其函数. 这样就会使代码变得非常的灵活. 

例如, ZStack分为多个版本——开源版、企业版、混合云版等. 如果一个服务在不同版本中的处理逻辑需要稍许不同, 那么就可以在开源版的代码中注册一个接口, 在另一个版本的服务中实现该接口. 这样也不会影响到开源版的原有逻辑. 从模块上看我们代码的是松耦合并且无法直接调用的, 但是在内存中, 却是可以调用得到的. 

## FAQ
### tomcat管理
```bash
$ cd ~/test/zstack-test/apache-tomcat-8.5.65/bin
$ ./startup.sh # 访问localhost:8080即可
$ ./shutdown.sh
```

> tomcat log在$CATALINA_HOME下, `tail -f management-server.log`即可

tomcat user role(5种):
- admin-gui
- manager-gui - allows access to the HTML GUI and the status pages
- manager-script - allows access to the text interface and the status pages
- manager-jmx - allows access to the JMX proxy and the status pages
- manager-status - allows access to the status pages only

tomcat在`webapps/manager/META-INF/context.xml`限制仅允许本机访问管理页面, 不限制可修改allow属性为`allow="^.*$"`.

### mvn打包报`No compiler is provided in this environment. Perhaps you are running on a JRE rather than a JDK`
原因: mvn打包需要jdk的环境, 而它自己没有找到jdk的配置目录

解决方法:`export JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64`

> JRE和JDK的区别: JRE是Java Runtime Envrionment, 是用来运行Java环境的, 并不能用来开发; JDK是Java Development Kit, 是Java的开发组件, 既能运行Java程序又能进行开发

> JRE中带不带headless的区别: 带headless的是用来运行不包含GUI的java程序的, 不带headless的可以运行带GUI的java程序

### java.lang.NoClassDefFoundError: javax/xml/bind/JAXBException
JAXB API是java EE 的API, 因此在java SE 9.0 中不再包含这个 Jar 包. java 9 中引入了模块的概念, 默认情况下, Java SE中将不再包含java EE 的Jar包, 而在 java 6/7 / 8 时关于这个API 都是捆绑在一起的.

解决方法:`export JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64`