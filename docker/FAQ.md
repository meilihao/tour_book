## prior storage driver "aufs" failed: driver not supported
docker 1.8.2,运行"sudo docker daemon",报以上错误.
linux kernel没有aufs驱动,试试其他docker storage-driver,比如"docker daemon -s overlay".

## 进入运行中的docker容器
[进入容器](http://dockerpool.com/static/books/docker_practice/container/enter.html),推荐使用**nsenter**.
