# env

## 描述

查看所有与此终端进程相关的环境变量

## 例

    # env
    # cat /proc/$PID/environ | tr '\0' '\n' # 查看某个进程的环境变量, `/proc/$PID/environ`由null字符（\0）分隔
    # env https_proxy=192.168.0.111:3142 http_proxy=192.168.0.111:3142 env|grep -i proxy # 一次设置多个env用空格分隔
