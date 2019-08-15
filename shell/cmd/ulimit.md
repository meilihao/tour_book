# ulimit

## 描述

限制系统用户对shell资源的访问

## 语法格式

```
ulimit [OPTIONS] [LIMIT]
```

## 选项

-H : 设置硬资源限制.
-S : 设置软资源限制.
-a : 显示当前所有的资源限制.
- -c size :设置core文件的最大值.单位:blocks
- -d size :设置数据段的最大值.单位:kbytes
- -f size :设置创建文件的最大值.单位:blocks
- -l size :设置在内存中锁定进程的最大值.单位:kbytes
- -m size :设置可以使用的常驻内存的最大值.单位:kbytes
- -n size :设置内核可以同时打开的文件描述符的最大值.单位:n
- -p size :设置管道缓冲区的最大值.单位:kbytes
- -s size :设置堆栈的最大值.单位:kbytes
- -t size :设置CPU使用时间的最大上限.单位:seconds
- -v size :设置虚拟内存的最大值.单位:kbytes
- -u <程序数目> : 用户最多可开启的程序数目