# docker cli

## 常用操作
下载镜像 ： docker pull [registry_hostname/[group/]]namespace/name[:tag], 比如`docker pull ubuntu@sha256:c7ca486dce6710697a81e2d08796a865f48c6525a3f81f3aa936a4905c9846be`或`docker pull arm64/centos:latest`
给镜像打tag : docker tag SOURCE_IMAGE[:TAG] TARGET_IMAGE[:TAG] # 镜像的ID相同,仅创建别名而已
镜像/容器的详细信息 : docker inspect [OPTIONS] NAME|ID [NAME|ID...]
镜像查找 : docker search [OPTIONS] TERM
查看镜像列表：docker images [OPTIONS] [REPOSITORY[:TAG]]
删除镜像 : docker rmi [OPTIONS] IMAGE [IMAGE...] #IMAGE可以是标签或ID.使用标签:一个镜像有多标签时,仅删除指定标签而已,否则彻底删除该镜像;使用id,先尝试删除所有该镜像的标签,再彻底删除该镜像.彻底删除时,如果有该镜像创建的容器存在,其默认是无法删除的;建议先删除依赖该镜像的所有容器,再删除该镜像,不推荐使用`-f`来强制删除,防止出现标签为`<none>`的临时镜像.
存出镜像 : docker save [OPTIONS] IMAGE [IMAGE...]

	```bash
	# docker save myimage:latest > myimage_latest.tar
	# docker save myimage:latest | gzip > myimage_latest.tar.gz
	# docker export <dockernameortag> | gzip > myimage_latest.tar.gz # for running or paused docker, use export
	# docker load < myimage_latest.tar
	# gunzip -c myimage_latest.tar.gz | docker load
	# docker tag ${image_id} myimage:latest # 重建tag
    # docker export $(docker create <image id>) --output="latest.tar" # 将image导出成rootfs
	```
载入镜像 : docker load [OPTIONS]
上传镜像 : docker push [OPTIONS] NAME[:TAG]
用镜像创建容器并启动：sudo docker run
查看容器列表：docker ps [OPTIONS]
暂停/取消暂停容器：docker pause/unpause CONTAINER [CONTAINER...]
停止容器：docker stop [OPTIONS] CONTAINER [CONTAINER...] // 向容器进程发送SIGTERM信号
kill容器：docker kill [OPTIONS] CONTAINER [CONTAINER...] // 向容器进程发送SIGKILL信号
启动容器：docker start [OPTIONS] CONTAINER [CONTAINER...]
删除容器：docker rm [OPTIONS] CONTAINER [CONTAINER...]
端口映射信息 : docker port CONTAINER [PRIVATE_PORT[/PROTO]]
从容器拷文件 : docker cp 7229f381542a:/go/src/app .
向容器拷文件 : docker cp ./app 7229f381542a:/go/src/app
查看指定镜像的创建历史 : docker history
获取docker的实时事件  : docker events
容器运行日志 : docker logs [OPTIONS] CONTAINER
查看镜像、容器、数据卷所占用的空间: docker system df
查看容器的存储层变化: docker diff CONTAINER # 最上层(读写层)和其他层(只读层)的差异
查看当前映射的端口配置: docker port CONTAINER
docker环境信息: docker info
构建image: docker build [--no-cache] PATH // --no-cache: 不使用image缓存; PATH代表使用指定路径的 Dockerfile 进行构建
查看bridge信息: docker network inspect bridge
列出docker daemon的所有网络: docker network ls
为容器添加网卡(比如bridge) : `docker network connect ${network_name} CONTAINER`

### 创建镜像
1. 基于已有镜像的容器创建(不推荐) : `docker commit [OPTIONS] CONTAINER [REPOSITORY[:TAG]]`
`sudo docker commit -m "Added a new file" -a "chen" 0e65be4364a8 test`
1. 基于本地模板导入 : `docker import [OPTIONS] file|URL|- [REPOSITORY[:TAG]]`
`sudo cat ubuntu-16.04.tar.gz |docker import - ubuntu:16.04`
1. 基于Dockerfile创建,**推荐**
    构建失败时, 可`docker run ${上一步构建成功的中间镜像}`用于调试Dockerfile

### 创建容器及启动
1. docker run [OPTIONS] IMAGE [COMMAND] [ARG...]
    - -t : 为docker分配一个伪终端(pseudo-tty)并绑定到容器的标准输入上
    - -i : 让容器的标准输入保持打开
    - -d : 让容器在后台以守护态(daemonized)形式运行
    - -v : 本地目录必须使用绝对路径,但本地文件可以使用相对路径,**推荐使用目录**
    - --rm : 容器退出后删除该容器

    > `-it` : 将伪终端作为容器的输入
1. docker start [OPTIONS] CONTAINER [CONTAINER...]

### 批量操作
docker rmi -f $(docker images | grep "none" | awk '{print $3}')
docker rmi -f $(docker images -q)
docker  image   rm   $(docker  image  ls   -a  -q)
docker container   stop   $(docker  container  ls   -a  -q)
docker   container   rm  $(docker  container  ls   -a  -q)
