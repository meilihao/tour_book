# README

## 命令注入方法

共3种
1. 将命令追加到builder.sh的`docker run`中
1. 将构建脚本放在build中, 构建image时COPY进去, 再用`ENTRYPOINT`执行
1. 将上述两种方法结合使用, 可参考[docker redis](https://github.com/docker-library/redis/blob/master/6.0/alpine/docker-entrypoint.sh)