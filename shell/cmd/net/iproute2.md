# iproute2
- [新的网络管理工具 ip 替代 ifconfig 零压力](http://www.linuxstory.org/replacing-ifconfig-with-ip/)
- [ip命令以及与net-tools的映射](https://linux.cn/article-3144-1.html)
- [放弃 ifconfig，拥抱 ip 命令](https://linux.cn/article-13089-1.html)


## 组件
- ip : 用于管理路由表和网络接口
- tc : 用于流量控制管理
- ss : 用于转储套接字统计信息
- lnstat : 用于转储linux网络统计信息
- bridge : 用于管理网桥地址和设备
- nstat : 类似于netstat, 但比它提供更多的信息

    ```bash
    nstat -a
    nstat --json
    ```

    ```bash
    $ strace -e open nstat 2>&1 > /dev/null|grep /proc
    open("/proc/uptime", O_RDONLY)          = 4
    open("/proc/net/netstat", O_RDONLY)     = 4
    open("/proc/net/snmp6", O_RDONLY)       = 4
    open("/proc/net/snmp", O_RDONLY)        = 4

    $ strace -e open netstat -s 2>&1 > /dev/null|grep /proc
    open("/proc/net/snmp", O_RDONLY)        = 3
    open("/proc/net/netstat", O_RDONLY)     = 3
    ```

    参考:
    - [Linux network metrics: why you should use nstat instead of netstat](https://loicpefferkorn.net/2016/03/linux-network-metrics-why-you-should-use-nstat-instead-of-netstat/)
    - [Linux network statistics reference](https://loicpefferkorn.net/2018/09/linux-network-statistics-reference/)
