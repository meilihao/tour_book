# ulimit

## 描述

shell的内建命令, 可用于对一个特定的shell和shell开启的子进程进行资源限制.

可限制类型:
- 所创建的内核文件的大小
- 进程数据块的大小
- shell进程创建文件的大小
- 内存锁住的大小
- 常驻内存集的大小
- 打开文件描述符的数量
- 分配堆栈的最大大小
- cpu时间
- 单个用户的最大线程数
- shell进程所能使用的最大虚拟内存

它还支持硬资源限制(一旦设置不能增加)和软资源限制(设置后可用增加, 但不能超过硬资源设置).

`/etc/security/limits.conf`即是pam_limits.so的配置文件. 要使 limits.conf 文件配置生效，必须要确保 pam_limits.so 文件被加入到启动文件(/etc/pam.d/login)中, 在RHEL中可能还与`/etc/pam.d/system-auth`有关, 可见下面的扩展.

ulimit, limits.conf 和 pam_limits 的关系，大致是这样的：
1. 用户进行登录，触发 pam_limits
2. pam_limits 读取 limits.conf，相应地设定用户所获得的 shell 的 limits
3. 用户在 shell 中，可以通过 ulimit 命令，查看或者修改当前 shell 的 limits
4. 当用户在 shell 中执行程序时，该程序进程会继承 shell 的 limits 值. 于是，limits 在子进程中生效了.

// todo https://feichashao.com/ulimit_demo/
// http://cn.linux.vbird.org/linux_basic/0320bash.php#variable_ulimit

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

## 扩展
/etc/pam.d/system-auth文件由Red-Hat和类似系统用于将常见安全策略组合在一起.它通常包含在其他/etc/pam.d策略文件中.

当通过sshd通过ssh访问系统时,会查询/etc.pam.d/sshd策略文件.此文件包含/etc/pam.d/system-auth,因此对/etc/pam.d/system-auth的更改有效.

当通过/bin/login程序登录时,会查询文件/etc/pam.d/login,因此对它的任何更改只会影响/bin /login.