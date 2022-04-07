# ulimit

## 描述

shell的内建命令, 可用于对一个特定的shell和shell开启的子进程进行资源限制.

> 类似功能的还有cgroupv2(**更灵活**), 可用`cgcreate`创建.

## 语法格式

```
ulimit [OPTIONS] [LIMIT]
```

## 选项

-H : 设置硬资源限制.
-S : 设置软资源限制.
-a : 显示当前所有的资源限制, 默认是`-S`.
- -c size : 创建core文件的最大值.单位: blocks
- -d size : 程序数据段的最大值.单位: KB
- -e      : 非实时调度优先级
- -f size : 创建文件的最大值.单位:blocks
- -i      : 待处理信号的最大数量
- -l size : 最大能锁住的内存空间.单位: KB
- -m size : 可以使用的常驻内存的最大值.单位: KB
- -n size : 内核可以同时打开的文件描述符的最大值.单位:n
- -p size : 管道缓冲区的最大值.单位: KB
- -q      : POSIX message queue最大缓冲空间
- -r      : 实时调度优先级
- -s size : 堆栈的最大值.单位: KB
- -t size : CPU使用时间的最大上限.单位:seconds
- -u <程序数目> : 用户最多可开启的程序数目
- -v size : 虚拟内存的最大值.单位: KB
- -x      : 文件锁的最大数量

## 扩展
### ulimit配置方法
1. 对登录用户进行限制

    1. 在/etc/profile、/etc/bashrc、~/.bash_profile、~/.bashrc文件中写入ulimit命令
    2. 直接在控制台输入ulimit命令, 这是临时配置, 重启失效

1. 对应用程序进行限制

    对nginx进行限制, 编写启动脚本

    startup.sh:
    ```
    ulimit -s 512
    ulimit -n 65536
    ```

1. 对多个用户或用户组进行限制

    在/etc/security/limits.conf中输入<domain> <type> <item> <value>, 每一行一个限制. 比如:
    ```bash
    #<domain>    <type>    <item>    <value>
    couchdb      hard      nofile    65536
    couchdb      soft      nofile    65536
    ```

    > 也可在/etc/security/limits.d中添加.

    domain 表示用户或者组的名字，还可以使用 * 作为通配符. Type 可以有两个值，soft 和 hard. Item 则表示需要限定的资源，可以有很多候选值，如 stack，cpu等. 通过添加对应的一行描述，则可以产生相应的限制.

1. 对系统全局的进程进行限制

    /etc/sysctl.conf

    体现:
    ```bash
    # cat /proc/sys/fs/file-max
    # cat /proc/sys/fs/file-nr
    ```


### ulimit 配置文件: limits.conf

`/etc/security/limits.conf`是pam_limits.so的配置文件. 要使 limits.conf 文件配置生效，必须要确保 pam_limits.so 文件已加入到`/etc/pam.d/`下相应的策略文件中.

// todo https://feichashao.com/ulimit_demo/
// http://cn.linux.vbird.org/linux_basic/0320bash.php#variable_ulimit

ulimit, limits.conf 和 pam_limits 的关系，大致是这样的：
1. 用户进行登录，触发 pam_limits

	- 当通过sshd通过ssh访问系统时,会查询策略文件/etc/pam.d/sshd.
	- 当通过/bin/login程序登录时,会查询策略文件/etc/pam.d/login.
2. pam_limits 读取 limits.conf，相应地设定用户所获得的 shell 的 limits
3. 用户在 shell 中，可以通过 ulimit 命令，查看或者修改当前 shell 的 limits
4. 当用户在 shell 中执行程序时，该程序进程会继承 shell 的 limits 值. 于是，limits 在子进程中生效了.

> /etc/pam.d/system-auth文件由Red-Hat和类似系统用于将常见安全策略组合在一起.它通常包含在其他/etc/pam.d策略文件中, 比如策略文件sshd, login等.

### 硬资源限制和软资源限制
硬资源限制: 对资源的绝对限制, 在任何情况下都不允许用户超过这个限制, 除非进程有root权限
软资源限制: 指用户可以在一定时间范围内(默认时为一周,在/usr/include/sys/fs/ufs_quota.h文件中设置)超过软限制的额度,在硬限制的范围内继续申请资源,同时系统会在用户登录时给出警告信息和仍可继续申请资源剩余时间.

## FAQ
### 查看进行的limits
`cat /proc/<pid>/limits`