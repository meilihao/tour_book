## 常用操作
下载镜像 ： docker pull [REGISTRY_HOSTNAME/]NAME[:TAG]
添加镜像tag : docker tag SOURCE_IMAGE[:TAG] TARGET_IMAGE[:TAG] # 镜像的ID相同,仅创建别名而已
镜像/容器的详细信息 : docker inspect [OPTIONS] NAME|ID [NAME|ID...]
镜像查找 : docker search [OPTIONS] TERM
查看镜像列表：docker images [OPTIONS] [REPOSITORY[:TAG]]
删除镜像 : docker rmi [OPTIONS] IMAGE [IMAGE...] #IMAGE可以是标签或ID.使用标签:一个镜像有多标签时,仅删除指定标签而已,否则彻底删除该镜像;使用id,先尝试删除所有该镜像的标签,再彻底删除该镜像.彻底删除时,如果有该镜像创建的容器存在,其默认是无法删除的;建议先删除依赖该镜像的所有容器,再删除该镜像,不推荐使用`-f`来强制删除,防止出现标签为`<none>`的临时镜像.
存出镜像 : docker save [OPTIONS] IMAGE [IMAGE...]
载入镜像 : docker load [OPTIONS]
上传镜像 : docker push [OPTIONS] NAME[:TAG]
用镜像创建容器并启动：sudo docker run
查看容器列表：docker ps [OPTIONS]
停止容器：docker stop [OPTIONS] CONTAINER [CONTAINER...]
启动容器：docker start [OPTIONS] CONTAINER [CONTAINER...]
删除容器：docker rm [OPTIONS] CONTAINER [CONTAINER...]
端口映射信息 : docker port CONTAINER [PRIVATE_PORT[/PROTO]]
容器运行日志 : docker logs [OPTIONS] CONTAINER

### 创建镜像

1. 基于已有镜像的容器创建(不推荐) : `docker commit [OPTIONS] CONTAINER [REPOSITORY[:TAG]]`
`sudo docker commit -m "Added a new file" -a "chen" 0e65be4364a8 test`
1. 基于本地模板导入 : `docker import [OPTIONS] file|URL|- [REPOSITORY[:TAG]]`
`sudo cat ubuntu-16.04.tar.gz |docker import - ubuntu:16.04`
1. 基于Dockerfile创建

### 创建容器及启动
1. docker create [OPTIONS] IMAGE [COMMAND] [ARG...] + docker start [OPTIONS] CONTAINER [CONTAINER...]
1. docker run [OPTIONS] IMAGE [COMMAND] [ARG...]
- -t : 让docker分配一个伪终端(pseudo-tty)并绑定到容器的标准输入上
- -i : 当前shell的标准输入绑定到容器的标准输出上
- -d : 让容器在后台以守护态(daemonized)形式运行
- -v : 本地目录必须使用绝对路径,但本地文件可以使用相对路径,**推荐使用目录**

## 删除镜像
sudo docker rmi $(sudo docker images | awk '/^<none>/ { print $3 }')

