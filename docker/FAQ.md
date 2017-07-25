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
export命令用于持久化容器(不是镜像);save命令用于持久化镜像(不是容器).
import/load用于导入镜像.

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

## alpine无法运行golang程序
```sh
# /app/micro
/bin/sh: 19: /app/micro: not found # 明明存在/app/micro文件,且有执行权限
```

推测: alpine使用 musl libc取代了glibc,导致程序依赖库的缺失,因此**不推荐alpine镜像跑golang程序,推荐使用和目标服务器相同发行版的镜像作为base image**

## Dockerfile的expose和docker run的-p
`-p`，是映射宿主端口和容器端口，即将容器的对应端口服务公开给外界访问，而 `EXPOSE`仅仅是声明容器打算使用什么端口而已，并不会自动在宿主进行端口映射.