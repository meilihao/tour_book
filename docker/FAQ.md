# FAQ

## prior storage driver "aufs" failed: driver not supported
docker 1.8.2,运行"sudo docker daemon",报以上错误.
linux kernel没有aufs驱动,试试其他docker storage-driver,比如"docker daemon -s overlay".

## 进入运行中的docker容器
[进入容器](https://yeasy.gitbooks.io/docker_practice/content/container/enter.html):
- `docker attach`,通过 docker attach 可以 attach 到**容器启动命令所用的终端**
使用 attach 命令有时候并不方便. 当多个窗口同时 attach 到同一个容器的时候，所有窗口都会同步显示. 当某个窗口因命令阻塞时,其他窗口也无法执行操作了.
且如果你用CTRL-c或`exit xxx`来退出，同时这个信号会kill该容器,但可通过 Ctrl+p 然后 Ctrl+q 组合键退出 attach 终端
- `docker exec`
- nsenter工具

attach 与 exec 主要区别如下:
1. attach 直接进入容器，不会启动新的进程;exec 则是在容器中打开**新的终端**，并且会启动**新的进程**
2. 如果想直接在终端中查看启动命令的输出，用 attach；其他情况使用 exec

## save/load和export/import 区别
export命令用于持久化容器(即容器快照);save命令用于持久化镜像(即镜像快照).
import/load用于导入镜像,但import用于操作export导出的容器,load用于操作save导出的镜像.

这两者的区别在于容器快照文件(`docker export`)将丢弃所有的历史记录和元数据信息（即仅保存容器当时的快照状态），而镜像存储文件(`docker save`)将保存完整记录，体积也要大。此外，从容器快照文件导入时可以重新指定标签等元数据信息.

> ![Image Layer（镜像层）](https://www.hi-linux.com/img/linux/docker-cmd-01.png)

## docker stop/kill区别
- `docker stop` : 先向容器发送SIGTERM信号,等待一段时间后(默认是10s),再发送SIGKILL信号终止容器
- `docker kill` : 直接发送SIGKILL信号来强行终止容器

## docker 退出
容器一旦运行完启动容器时的cmd命令后就会退出

## Dockerfile ADD/COPY
区别:
- ADD 的`<src>`可以为URL; COPY只能复制build context(即Dockerfile所在目录)中的文件
- ADD指令会将tar文件解压到指定位置，而COPY指令只做复制操作

**一般优先使用 COPY,因为它比 ADD 更透明**

> COPY的目标目录必须带`/`, 否则创建的文件会没有进入权限.

## Dockerfile RUN/ENTRYPOINT/CMD
RUN是只在build镜像时运行，固化在image中, 其先于CMD和ENTRYPOINT的. CMD会在每次启动容器的时候运行.

Dockerfile中可以有多个CMD, 但只有最后一个生效. **CMD可以被docker run指定的命令取代**.
Dockerfile中可以有多个ENTRYPOINT, 但只有最后一个生效. **CMD/docker run指定的命令会被当做参数传递给ENTRYPOINT**.

ENTRYPOINT是容器运行程序的入口.

ENTRYPOINT指令有两种格式，CMD指令有三种格式, RUN有两种格式：
```
ENTRYPOINT ["command", "arg1", "arg2"] // exec 模式的写法: docker直接执行[command], 推荐写法
ENTRYPOINT command arg1 arg2 // shell 模式的写法: docker执行`/bin/sh -c [command]`

CMD ["command", "arg1", "arg2"]
CMD command arg1 arg2
CMD ["arg1, "arg2"] # as default parameters to ENTRYPOINT, 必须配合exec 模式的ENTRYPOINT来使用

RUN ["command", "arg1", "arg2"]
RUN command
```

关于ENTRYPOINT和CMD的交互细节在[官方Dockerfile的Understand how CMD and ENTRYPOINT interact里](https://docs.docker.com/engine/reference/builder/).

永远使用Exec表示法: 组合使用ENTRYPOINT和CMD格式时确保你一定用的是Exec表示法. 如果用其中一个用的是Shell表示法, 或者一个是Shell表示法, 另一个是Exec表示法, 将永远得不到预期的效果:
```text
Dockerfile    Command

ENTRYPOINT /bin/ping -c 3
CMD localhost               /bin/sh -c '/bin/ping -c 3' /bin/sh -c localhost

ENTRYPOINT ["/bin/ping","-c","3"]
CMD localhost               /bin/ping -c 3 /bin/sh -c localhost

ENTRYPOINT /bin/ping -c 3
CMD ["localhost"]"          /bin/sh -c '/bin/ping -c 3' localhost

ENTRYPOINT ["/bin/ping","-c","3"]
CMD ["localhost"]            /bin/ping -c 3 localhost
```
从上面看出, 只有ENTRYPOINT和CMD都用Exec表示法, 才能得到预期的效果.

ps: shell 形式是以`/bin/sh -c [command]`启动，它不传递信号. 这意味着可执行文件将不是容器的PID 1，并且不会接收Unix信号，因此您的可执行文件将不会从docker stop <container>接收到SIGTERM

> 默认情况下，Docker 会提供一个隐含的 ENTRYPOINT:`/bin/sh -c`. 所以在不指定 ENTRYPOINT 时，实际上运行容器里的在完整进程是`/bin/sh -c ${CMD}`，即 CMD 的内容就是 ENTRYPOINT 的参数. 因此我们会统一称 Docker 容器的启动进程为 ENTRYPOINT，而不是 CMD. 推荐Dockerfile至少要指定一个CMD或者ENTRYPOINT命令.

## alpine无法运行golang程序
```sh
# /app/micro
/bin/sh: 19: /app/micro: not found # 明明存在/app/micro文件,且有执行权限
```

推测: alpine使用 musl libc取代了glibc,导致程序依赖库的动态链接库缺失(通过`file`,`ldd`,`readelf -d xxx`命令查看),因此**不推荐alpine镜像跑golang程序,推荐使用和目标服务器相同发行版的镜像作为base image**.同时go编译时[禁用cgo `CGO_ENABLED=0`](https://stackoverflow.com/questions/36279253/go-compiled-binary-wont-run-in-an-alpine-docker-container-on-ubuntu-host)可解决这个问题,或使用`go build -ldflags '-linkmode "external" -extldflags "-static"' server.go`进行静态编译cgo来解决(待定,测试后发现还是动态链接依赖).

> [也谈Go的可移植性](http://tonybai.com/2017/06/27/an-intro-about-go-portability/)

## Dockerfile的expose和docker run的-p
`-p`，是映射宿主端口和容器端口，即将容器的对应端口服务公开给外界访问，而 `EXPOSE`仅仅是声明容器打算使用什么端口而已，并不会自动在宿主进行端口映射.

## 去掉sudo
将用户加入docker用户组:
1. `sudo cat /etc/group | grep docker`
1. 如果不存在docker组，可以添加`sudo groupadd docker`
1. 添加当前用户到docker组，`sudo gpasswd -a $USER docker`/`sudo usermod -aG docker chen`
1. 重启docker服务,`sudo systemctl restart docker`,用户需要**注销/重启**系统使上一步的修改生效.
1. 如果权限不够，`sudo chmod a+rw /var/run/docker.sock`

## 安装docker
通过[Docker CE 镜像源站](https://yq.aliyun.com/articles/110806)安装.

或[官方方式](https://docs.docker.com/engine/installation/#server).

## go程序 in alpine容器 报: /usr/local/go/lib/time/zoneinfo.zip: no such file or directory
```sh
apk add --no-cache tzdata
```

## Dockerfile
Dockerfile 中的每个原语执行后都会生成一个对应的镜像层. 即使原语本身并没有明显地修改文件的操作（比如，ENV 原语），它对应的层也会存在. 只不过在外界看来，这个层是空的.

有时docker run后面不加命令是因为在 Dockerfile 中已经指定了 CMD, 否则就需要将进程的启动命令加在`docker run`后面.

## Layer/Image ID
镜像由一系列层组成, 每层都用64位的十六进制数来表示, 非常类似Git repo中的commit.
镜像最上层的layer ID就是该镜像的ID, 其默认存储在`/var/lib/docker`下.

官方推荐使用`dockerviz`工具来分析镜像.

## 多阶段构建
Docker image的多阶段构建中, 每个From语句开启一个构建阶段，并且可以通过`as`语法为此阶段构建命名(比如下面的builder).

```sh
//Dockerfile

FROM golang:alpine as builder

WORKDIR /go/src
COPY httpserver.go .

RUN go build -o httpd ./httpserver.go

From alpine:latest

WORKDIR /root/
COPY --from=builder /go/src/httpd . # 通过COPY命令在两个阶段构建产物之间传递数据
RUN chmod +x /root/httpd

ENTRYPOINT ["/root/httpd"]
```

## 日志位置
```
/var/lib/docker/containers/${containerid}/${containerid}.log-json.log
```

# error
## Error response from daemon: Driver overlay2 failed to remove root filesystem
完整错误:
```
Error response from daemon: Driver overlay2 failed to remove root filesystem 95ee7e853063ca485ef7ce82b17db977303280df34db4fac2f3fa0367ab50b2c: remove /var/lib/docker/overlay2/dd95ab1ff29c37f16450194f79b9876a7e34da2dfbb8ee609745f00e017cb91c/merged: device or resource busy
```

解决方法`umount /var/lib/docker/overlay2/容器id`:
```
sudo umount /var/lib/docker/overlay2/dd95ab1ff29c37f16450194f79b9876a7e34da2dfbb8ee609745f00e017cb91c
```

## Error response from daemon: driver failed programming external connectivity on endpoint seafile 
```sh
# docker start c2f56bcd3b32
Error response from daemon: driver failed programming external connectivity on endpoint seafile (c4e6c105b6c080b91895576eedcf2b94adfa34fe081d147e98d417fd20c7761f):  (iptables failed: iptables --wait -t nat -A DOCKER -p tcp -d 127.0.0.1 --dport 9090 -j DNAT --to-destination 172.17.0.2:80 ! -i docker0: iptables: No chain/target/match by that name.
 (exit status 1))
```

docker服务启动时定义的自定义链DOCKER由于某种原因被清掉,重启docker服务即可重新生成自定义链DOCKER.

## alpine 和 busybox 比较
单从image的size上来说，busybox更小. 不过busybox默认的libc实现是uClibc(1.2MB)，而我们通常运行环境使用的libc实现都是glibc，因此我们要么选择静态编译程序，要么使用busybox:glibc(4.46MB)镜像作为base image.

而 alpine image 是另外一种蝇量级 base image，它使用了比 glibc 更小更安全的 musl libc(5.53MB). 不过和 busybox image 相比，alpine image 体积还是略大. 除了因为 musl比uClibc 大一些之外，alpine还在镜像中添加了自己的**包管理系统apk**，开发者可以使用apk在基于alpine的镜像中添 加需要的包或工具. 因此，对于普通开发者而言，alpine image显然是更佳的选择, 不过alpine使用的libc实现为musl，与基于glibc上编译出来的应用程序不兼容, 通常因为找不到glibc的动态共享库文件而报`xxx "no such file or directory"`.

目前Docker官方已推荐使用Alpine作为基础镜像环境

对于Go应用来说，我们可以采用静态编译的程序，但一旦采用静态编译，也就意味着我们将失去一些libc提供的原生能力，比如, 在linux上就无法使用系统提供的DNS解析能力，只能使用Go自实现的DNS解析器, 不过对于使用并没有影响.

## the input device is not a TTY
jekins使用docker run构建项目时报:
```log
the input device is not a TTY
Build step 'Execute shell' marked build as failure
```

解决方法: `docker run`时使用了`-it`选项, 去掉即可.