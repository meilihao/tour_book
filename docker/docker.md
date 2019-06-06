# docker
## 概念
### 容器
轻量级的应用打包分发技术,是为了解决应用运行环境的一致性而开发, 组成部分:
- 应用本身
- 应用依赖的运行环境rootfs(即容器镜像)

### OCI, Open Container Initiative
定义容器规范的组织, 现有规范:
- runtime spec : 管理容器和容器镜像的软件, runC是该规范的实现
- image format spec : 定义容器镜像的标准

### 编排引擎
管理集群资源的生命周期并提供服务, 包括容器管理/调度, 集群定义, 服务发现. 当前主流是Kubernetes.

### 架构
[图解 Docker 架构(细)](https://www.hi-linux.com/posts/13732.html)
![Docker 架构(粗)](images/arch.jpg)

docker是典型的C/S架构:
- docker cli : 向docker daemon发送请求
- docker daemon : 根据docker cli的请求创建,运行,监控容器; 构建,存储镜像. 默认仅监听localhost client的请求.

![一张图掌握 Docker 命令](https://static.oschina.net/uploads/img/201702/09111906_odFS.png)

## 镜像
镜像的实现使用了Union File System 也叫 UnionFS，其最主要的功能是将多个不同位置的目录联合挂载（union mount）到同一个目录下.

所有对容器的改动, 包括添加, 修改, 删除都只发生在容器层(即可读写层), 细节:
1. 添加 : 发生在容器层
1. 读取 : Docker会从上往下依次在各层查找, 找到后载入内存
1. 修改 : 不在容器层时，Docker会从上往下依次在各层查找, 找到后复制到容器层, 再修改
    使用了CopyOnWrite机制: 仅当需要修改时才复制数据, 优势是共享数据，减少物理空间占用
1. 删除 : 只读层存在删除对象，使用 whiteout(对文件)/opaque(对目录) 机制: 通过在容器层建立对应的 whiteout/opaque，来遮挡下层分支中的所有路径名相同的文件/目录

当前默认GraphDriver是overlay2. overlay2仅有两层, 性能上比使用多层的aufs(不推荐)有优势.

### [Dockerfile](https://docs.docker.com/engine/reference/builder/)
`FROM scratch`表示从零开始构建,不推荐: 因为没有常用的调试工具. 常用alpine.

## docker运行
容器的生命周期依赖启动时执行的命令, 只要该命令不结束, 容器就会一直运行.

容器的启动命令可通过`docker ps -a`中的`COMMAND`列查看.

docker run的restart策略:
- no : 容器退出时不要自动重启, 是默认值
- always : 无论容器因何原因退出(包括正常退出)都立即重启
- on-failure[:max-retries] : 只在容器以非0状态码退出时重启. 可选的max-retries表示尝试重启容器的次数
- unless-stopped : 不管退出状态码是什么始终重启容器, 但是不考虑在Docker daemon启动时就已经停止了的容器

![容器生命周期](http://seo-1255598498.file.myqcloud.com/full/0c3da8ddb467b2c716c4570a4760bb4cb95fdb3a.jpg)

### 资源限制
资源限制可使用`progrium/stress`镜像验证.

docker run的资源限制参数:
- -m : 内存限制
- --memory-swap : 内存+swap的限制
- -c : 使用cpu的权重(仅cpu资源紧张的情况下有效), 具体分配到的cpu取决于该值占所有容器cpu share总和的比例.
- --blkio-weight : io的权重值, 也`-c`类似