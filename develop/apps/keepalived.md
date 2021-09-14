# keepalived
它是集群管理中保证集群高可用的一个服务软件, 用来防止单点故障. VRRP协议是keepalived实现的基础.

keepalived可提供vrrp以及health-check功能，可以只用它提供双机浮动的vip（vrrp虚拟路由功能），这样可以简单实现一个双机热备高可用功能.

keepalived是以VRRP虚拟路由冗余协议为基础实现高可用的，可以认为是实现路由器高可用的协议，即将N台提供相同功能的路由器组成一个路由器组，这个组里面有一个master和多个backup，master上面有一个对外提供服务的vip（该路由器所在局域网内其他机器的默认路由为该vip），**master会发组播**，当backup收不到VRRP包时就认为master宕掉了，这时就需要根据VRRP的优先级来选举一个backup当master. 这样的话就可以保证路由器的高可用了.

在Keepalived 有两个角色：Master(一个)、Backup（多个），如果设置一个为Master，但Master挂了后再起来，必然再次业务又一次切换，这对于有状态服务是不可接受的. 解决方案就是所有机器都设置为Backup，且在优先级高的Backup配置中设置nopreemt(不抢占).

## 组件
keepalived也是模块化设计，不同模块复杂不同的功能，它主要有三个模块，分别是core、check和VRRP，其中：
- core模块：为keepalived的核心组件，负责主进程的启动、维护以及全局配置文件的加载和解析
- check：负责健康检查，包括常见的各种检查方式
- VRRP模块：是来实现VRRP协议的
- watch dog:监控check和vrrp进程的看管者，check负责检测器子进程的健康状态，当其检测到master上的服务不可用时则通告vrrp将其转移至backup服务器上

keepalived正常启动的时候，共启动3个进程：
一个是父进程，负责监控其子进程；一个是VRRP子进程，另外一个是checkers子进程. 两个子进程都被系统watchlog看管，两个子进程各自负责复杂自己的事.

Healthcheck子进程检查各自服务器的健康状况, 例如http,lvs. 如果healthchecks进程检查到master上服务不可用了，就会通知本机上的VRRP子进程，让他删除通告，并且去掉虚拟IP，转换为BACKUP状态.

## 配置
参考:
- [Keepalive详解](https://www.cnblogs.com/rexcheny/p/10778567.html)
- [keepalived 配置整理](https://weizhimiao.github.io/2017/02/11/keepalived%E9%85%8D%E7%BD%AE%E6%95%B4%E7%90%86/)

keepalived只有一个配置文件keepalived.conf，配置文件里面主要包括以下几个配置项，分别是global_defs、static_ipaddress、static_routes、VRRP_script、VRRP_instance和virtual_server.

总的来说，keepalived主要有三类配置区域：
1）全局配置(Global Configuration)

    全局配置又包括两个子配置：
    1. 全局定义(global definition)
    1. 静态路由配置(static ipaddress/routes)


    ```conf
    global_defs {
        notification_email { # 表示keepalived在发生诸如切换操作时需要发送email通知，邮件接收地址可以多个，每行一个
            admin@example.com
            guest@example.com
        }
         
        notification_email_from notify@example.com # 表示发送通知邮件时邮件源地址是谁
        smtp_server 127.0.0.1 # 表示发送email时使用的smtp服务器地址，这里可以用本地的sendmail来实现
        stmp_connect_timeout 30
        router_id node1 # 路由id,主备节点不能相同
        # VRRP的ipv4和ipv6的广播地址，配置了VIP的网卡向这个地址广播来宣告自己的配置信息，下面是默认值
        vrrp_mcast_group4 224.0.0.18
        vrrp_mcast_group6 ff02::12 # 默认只使用ipv4
    }
    ```
2）VRRPD配置

    VRRPD配置包括三个类：
    1. VRRP同步组(synchroization group)
    1. VRRP实例(VRRP Instance)

        ```conf
        vrrp_instance VI_1 {
            state MASTER               # 指定instance初始角色(MASTER 表示主节点，BACKUP 表示备份节点)，实际根据优先级决定. 与backup节点不一样
            interface eth0             # 表示发vrrp包的接口
            virtual_router_id 51       # VRID(0-255)，相同VRID为一个组，决定多播MAC地址. 主备节点需要设置为相同
            priority 100               # 优先级(1-255), 主节点的优先级需要设置比备份节点高. backup节点改为90.
            advert_int 1               # 设置主备之间的检查间隔，单位为秒
            authentication {
                auth_type PASS         # 认证方式，可以是pass或ha
                auth_pass 1111         # 认证密码. 同一个vrrp_instance下，MASTER和BACKUP的密码必须一致才能正常通信
            }
            # vip
            virtual_ipaddress {
                # 172.16.60.129/24         # 指定VIP，不指定网卡，默认为eth0,注意：不指定/prefix,默认为/32
                172.16.60.129/24 dev  eth0 # VIP地址
                # 192.168.1.33/24 brd 192.168.1.255 dev eno1 label eno1:1 # IP/掩码 dev 配置在哪个网卡的哪个别名上
            }
            notify_master "/etc/keepalived/keepalived.sh master" # 当前节点状态转为master时触发的脚本
            notify_backup "/etc/keepalived/keepalived.sh backup" # 当前节点状态转为backup时触发的脚本
            notify_fault "/etc/keepalived/keepalived.sh fault" # 当前节点keepalived出现故障转为"FAULT"状态时触发的脚本
            notify_stop "/etc/keepalived/keepalived.sh fault" # 当前节点keepalived停止时触发的脚本
            notify xxx # 表示只要状态切换都会调用的脚本，并且该脚本是在以上四个脚本执行之后再调用的
            # 追踪脚本，通常用于去执行vrrp_script中定义的脚本内容
            track_script {
                check_running
            }
            track_interface {  # 设置额外的监控，里面那个网卡出现问题都会切换. 通常不使用
                eth0
                eth1
            }
            #虚拟路由，当IP漂过来之后需要添加的路由信息, 可选.
            virtual_routes {
                172.16.0.0/12 via 10.210.214.1
                192.168.1.0/24 via 192.168.1.1 dev eth1
                default via 202.102.152.1
            }
            # debug # debug级别
            # 将多播vrrp通告调整为单播. keepalived在组播模式下所有的信息都会向224.0.0.18的组播地址发送, 产生众多的无用信息，并且会产生干扰和冲突，可以将组播的模式改为单拨. 这是一种安全的方法，避免局域网内有大量的keepalived造成虚拟路由id的冲突
            unicast_src_ip  172.19.1.14   # 本机ip, backup节点ip相反
            unicast_peer {              
                172.19.1.15      #对端ip
            }
            nopreempt                   # 定义工作模式为非抢占模式, 默认是抢占模式. **抢占模式时主节点故障恢复后, 就会重新抢回vip (根据配置里的优先级决定的).**. 首先nopreemt必须在state为BACKUP的节点上才生效（因为是BACKUP节点决定是否来成为MASTER的）. 推荐使用将所有节点的state都设置成BACKUP并且都加上nopreempt选项，这样就完成了关于autofailback功能，当想手动将某节点切换为MASTER时只需去掉该节点的nopreempt选项并且将priority改的比其他节点大，然后重新加载配置文件即可（等MASTER切过来之后再将配置文件改回去再reload一下）
            preempt_delay 300           # 抢占式模式下，节点上线后触发新选举操作的延迟时长, 避免节点还没进入工作状态就进行抢占导致小段时间内不可用. 这里的间隔时间要大于vrrp_script中定义的时长
        }
        ```

        ```bash
        vim /etc/keepalived/keepalived.sh
        #!/bin/bash
        # http://blog.mykernel.cn/2020/10/22/keepalived%E7%BC%96%E8%AF%91%E5%AE%89%E8%A3%85%E5%8F%8A%E4%BC%98%E5%8C%96/
        #author :Magedu
        #Description : an example of notify script
        contact='root@localhost'

        notify() {
            local mailsubject="$(hostname) to be $1, vip floating"
            local mailbody="$(date +'%F %T'): vrrp transition, $(hostname) changed to be $1"  #时间状态改变
            # 此步骤可以修改为调用python脚本完成微信报警
            echo "$mailbody" | mail -s "$mailsubject" $contact
        }


        case $1 in
        master)
            notify master
            systemctl start nginx.service
            exit 0
            ;;
        backup)
            notify backup
            systemctl restart nginx.service
            exit 0
            ;;
        fault)
            notify fault
            systemctl stop nginx.service
            exit 0
            ;;
        *)
            echo "Usage: $(basename $0) {master|backup|fault}"
            exit 1
            ;;
        esac
        ```

    1. VRRP脚本

        通过脚本来检测服务是否正常.

        ```conf
        vrrp_script <SCRIPT_NAME> {
           script <STRING>|<QUOTED-STRING> # path of the script to execute，需要运行的脚本，返回值为0表示正常; 其它值都会当成检测失败.
           interval <INTEGER>  # seconds between script invocations, default 1 second ，脚本运行时间，即隔多少秒去检测
           timeout <INTEGER>   # seconds after which script is considered to have failed，脚本运行的超时时间
           weight <INTEGER:-254..254>  # adjust priority by this weight, default 0
           rise <INTEGER>              # required number of successes for OK transition，配置几次检测成功才认为服务正常
           fall <INTEGER>              # required number of successes for KO transition，配置几次检测失败才认为服务异常
           user USERNAME [GROUPNAME]   # user/group names to run script under
                                       #   group default to group of user
           init_fail                   # assume script initially is in failed state，配置初始时失败状态
        }
        ```

        example:
        ```conf
        vrrp_script check_running {
           script "/usr/local/bin/check_running"
           interval 10
           weight 10
           timeout 2
           fall 3
        }
        ```

        ```bash
        # cat /usr/local/bin/check_running
        #!/bin/bash
        #author : panbuhei
        #nginx check script in keepalived

        NGINX_PIDNUM=`ps -ef | grep nginx | grep -v grep | wc -l`

        NGINX_PORTNUM=`ss -antpl | grep nginx | wc -l`

        if [ $NGINX_PIDNUM -eq 0 ];then
            exit 1
        elif [ $NGINX_PORTNUM -eq 0 ];then
            exit 1
        else
            exit 0
        fi
        # --- 其他 vrrp_script
        # cat /etc/keepalived/curl.sh
        #!/bin/bash
        curl -m 2 -I http://172.20.27.10:9000/haproxy-status &> /dev/null
        if [ $? -eq 0 ];then
            exit 0
        else
            exit 2
        fi
        # cat /etc/keepalived/ping.sh
        #!/bin/bash
        ping -c 2 172.20.0.1 &> /dev/null
        if [ $? -eq 0 ];then
            exit 0
        else
            exit 2
        fi
        ```
3）LVS配置

    如果没有配置LVS+keepalived，那么无需配置这段区域. 如果用的是nginx来代替LVS，也无需配置这里. 这里的LVS配置是专门为keepalived+LVS集成准备的.

## FAQ
### keepalived两个节点都出现了vip
env：kylinV10 (fork from centos 7.x)

根源: 节点间通信故障.

在backup节点执行`tcpdump -i enp1s0 vrrp -n`, 发现它一直在发送vrrp通告. 理论上备节点收到主节点的通告, 并发现其优先级高于自己就不会主动对外发送通告了, 转为一直接收主节点发送的vrrp通告, 即正常环境只有主节点会一直发送vrrp通告.

在主备节点分别执行:
```bash
setenforce 0                                                      # 临时关闭SELINUX
sed -i "s/^SELINUX=.*/SELINUX=disabled/g" /etc/selinux/config     # 永久关闭SELINUX
firewall-cmd --direct --permanent --add-rule ipv4 filter INPUT 0 --destination 224.0.0.18 --protocol vrrp -j ACCEPT
firewall-cmd --direct --permanent --add-rule ipv4 filter OUTPUT 0 --destination 224.0.0.18 --protocol vrrp -j ACCEPT
firewall-cmd --reload
```

如果还是不行(禁用selinux和firewalld后), 可用单播的`unicast_src_ip的unicast_peer`.

以前遇到过多播不通的情况即节点都在发送vrrp通告, 同网段的其他节点都收到了但它们自己就是不收, 此时`ethtool -S <net_dev>`都没有`multicast`属性但netdev上有`MULTICAST`, 已排除情况:
```bash
# 已排除selinux, 防火墙
# for i in /proc/sys/net/ipv4/conf/*/rp_filter ; do echo 0 > "$i";   done # 无效
# echo "0" > /proc/sys/net/ipv4/icmp_echo_ignore_broadcasts # 无效
# route add -net 224.0.0.0 netmask 240.0.0.0 enp1s0 # 无效
# 添加vrrp_mcast_group6也没生效
# nstat -a # 发现[UdpIgnoredMulti](https://elixir.bootlin.com/linux/v5.14/source/include/uapi/linux/snmp.h#L161)一直在增长, 通过`strace -e open nstat -a`可知, IgnoredMulti来源于`/proc/net/snmp`(参考[TCP 统计信息详解](https://github.com/moooofly/MarkSomethingDown/blob/master/Linux/TCP%20%E7%9B%B8%E5%85%B3%E7%BB%9F%E8%AE%A1%E4%BF%A1%E6%81%AF%E8%AF%A6%E8%A7%A3.md))
```

# VRRP
VRRP全称Virtual Router Redundancy Protocol，即虚拟路由冗余协议。对于VRRP，需要清楚知道的是：
1. VRRP是用来实现路由器冗余的协议。
1. VRRP协议是为了消除在静态缺省路由环境下路由器单点故障引起的网络失效而设计的主备模式的协议，使得发生故障而进行设计设备功能切换时可以不影响内外数据通信，不需要再修改内部网络的网络参数。
1. VRRP协议需要具有IP备份，优先路由选择，减少不必要的路由器通信等功能。
1. VRRP协议将两台或多台路由器设备虚拟成一个设备，对外提供虚拟路由器IP（一个或多个）。然而，在路由器组内部，如果实际拥有这个对外IP的路由器如果工作正常的话，就是master，或者是通过算法选举产生的，MASTER实现针对虚拟路由器IP的各种网络功能，如ARP请求，ICMP，以及数据的转发等，其他设备不具有该IP，状态是BACKUP。除了接收MASTER的VRRP状态通告信息外，不执行对外的网络功能，当主级失效时，BACKUP将接管原先MASTER的网络功能。
1. VRRP协议配置时，需要配置每个路由器的虚拟路由ID(VRID)和优先权值，**使用VRID将路由器进行分组，具有相同VRID值的路由器为同一个组**，VRID是一个0-255的整整数；**同一个组中的路由器通过使用优先权值来选举MASTER，优先权大者为MASTER**，优先权也是一个0-255的正整数。