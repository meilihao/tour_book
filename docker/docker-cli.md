## 常用操作
下载镜像 ： sudo docker pull
查看镜像列表：sudo docker images
删除镜像 : sudo docker rmi
用镜像创建容器：sudo docker creat
用镜像创建容器并启动：sudo docker run
查看容器列表：sudo docker ps
停止容器：sudo docker stop
启动容器：sudo docker start
删除容器：sudo docker rm

## 删除镜像
sudo docker rmi $(sudo docker images | awk '/^<none>/ { print $3 }')

