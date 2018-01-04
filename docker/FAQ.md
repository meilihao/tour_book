## prior storage driver "aufs" failed: driver not supported
docker 1.8.2,运行"sudo docker daemon",报以上错误.
linux kernel没有aufs驱动,试试其他docker storage-driver,比如"docker daemon -s overlay".

## 进入运行中的docker容器
[进入容器](https://yeasy.gitbooks.io/docker_practice/content/container/enter.html):
- `docker attach`
使用 attach 命令有时候并不方便。当多个窗口同时 attach 到同一个容器的时候，所有窗口都会同步显示。当某个窗口因命令阻塞时,其他窗口也无法执行操作了
- `docker exec`
- nsenter工具

## save/load和export/import 区别
export命令用于持久化容器(即容器快照);save命令用于持久化镜像(即镜像快照).
import/load用于导入镜像,但import用于操作export导出的容器,load用于操作save导出的镜像.

这两者的区别在于容器快照文件(`docker export`)将丢弃所有的历史记录和元数据信息（即仅保存容器当时的快照状态），而镜像存储文件(`docker save`)将保存完整记录，体积也要大。此外，从容器快照文件导入时可以重新指定标签等元数据信息.

## docker stop/kill区别
- `docker stop` : 先向容器发送SIGTERM信号,等待一段时间后(默认是10s),再发送SIGKILL信号终止容器
- `docker kill` : 直接发送SIGKILL信号来强行终止容器

## docker 退出
容器一旦运行完启动容器时的cmd命令后就会退出

## Dockerfile ADD/COPY
区别:
- ADD 的`<src>`可以为URL
- ADD指令会将tar文件解压到指定位置，而COPY指令只做复制操作

**一般优先使用 COPY,因为它比 ADD 更透明**
## Dockerfile RUN/ENTRYPOINT/CMD
ENTRYPOINT指令有两种格式，CMD指令有三种格式：
```
ENTRYPOINT ["程序名", "参数1", "参数2"]
ENTRYPOINT 命令 参数1 参数2

CMD ["程序名", "参数1", "参数2"]
CMD 命令 参数1 参数2
CMD 参数1 参数2 # as default parameters to ENTRYPOINT
```
每个Dockerfile只能有一条CMD命令,有多条时,只有最后一条会被执行;用户启动容器时指定的命令会覆盖CMD指定的命令.
每个Dockerfile只能有一条ENTRYPOINT命令,有多条时,只有最后一条会被执行;用户启动容器时指定的命令会作为参数传递给ENTRYPOINT命令.

ENTRYPOINT是容器运行程序的入口.

RUN是在build成镜像时就运行的，先于CMD和ENTRYPOINT的，CMD会在每次启动容器的时候运行，而RUN只在创建镜像时执行一次，固化在image中.

ENTRYPOINT和CMD的不同点在于执行docker run时参数传递方式，**CMD指定的命令可以被docker run传递的命令覆盖**.
ENTRYPOINT会把容器名后面的所有内容都当成参数传递给其指定的命令(**不会对命令覆盖**)

关于ENTRYPOINT和CMD的交互，用一个官方表格可以说明：
<table>
<thead>
<tr>
<th align="left"></th>
<th align="left"><strong>No ENTRYPOINT</strong></th>
<th align="left"><strong>ENTRYPOINT exec_entry p1_entry</strong></th>
<th align="left"><strong>ENTRYPOINT ["exec_entry", "p1_entry"]</strong></th>
</tr>
</thead>
<tbody>
<tr>
<td align="left"><strong>No CMD</strong></td>
<td align="left">error, not allowed</td>
<td align="left">/bin/sh -c exec_entry p1_entry</td>
<td align="left">exec_entry p1_entry</td>
</tr>
<tr>
<td align="left"><strong>CMD ["exec_cmd", "p1_cmd"]</strong></td>
<td align="left">exec_cmd p1_cmd</td>
<td align="left">/bin/sh -c exec_entry p1_entry</td>
<td align="left">exec_entry p1_entry exec_cmd p1_cmd</td>
</tr>
<tr>
<td align="left"><strong>CMD ["p1_cmd", "p2_cmd"]</strong></td>
<td align="left">p1_cmd p2_cmd</td>
<td align="left">/bin/sh -c exec_entry p1_entry</td>
<td align="left">exec_entry p1_entry p1_cmd p2_cmd</td>
</tr>
<tr>
<td align="left"><strong>CMD exec_cmd p1_cmd</strong></td>
<td align="left">CMD exec_cmd p1_cmd</td>
<td align="left">/bin/sh -c exec_entry p1_entry</td>
<td align="left">exec_entry p1_entry /bin/sh -c exec_cmd p1_cmd</td>
</tr></tbody></table>

ps: shell 形式防止使用任何CMD或运行命令行参数，但是缺点是您的ENTRYPOINT将作/bin/sh -c的子命令启动，它不传递信号。这意味着可执行文件将不是容器的PID 1，并且不会接收Unix信号，因此您的可执行文件将不会从docker stop <container>接收到SIGTERM

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
1. 重启docker服务,`sudo systemctl restart docker`,用户需要**重新登录**系统使上一步的修改生效.
1. 如果权限不够，`sudo chmod a+rw /var/run/docker.sock`

## 安装docker
通过[Docker CE 镜像源站](https://yq.aliyun.com/articles/110806)安装.

## go程序 in alpine容器 报: /usr/local/go/lib/time/zoneinfo.zip: no such file or directory
```sh
apk add --no-cache tzdata
```

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