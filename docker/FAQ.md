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