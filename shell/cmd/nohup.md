# nohup

## 描述

可以忽略挂断信号，继续执行命令,避免终端关闭时，其关联进程也自动被关闭.

nohup原理: 终端关闭后会给此终端下的每一个进程发送SIGHUP信号，而使用nohup运行的进程则会忽略这个信号，因此终端关闭后进程也不会退出.

## 语法格式

```
nohup COMMAND [ARG]...
```

## 例

	# 2>&1是将标准错误（2）重定向到标准输出（&1），标准输出（&1）再被重定向输入到myout.file文件中
    # nohup command > myout.file 2>&1 &

参考:

[Linux 技巧：让进程在后台可靠运行的几种方法](https://www.ibm.com/developerworks/cn/linux/l-cn-nohup/)