# firewall-cmd

firewalld 服务是把配置好的防火墙策略交由内核层面的 nftables 包过滤框架来处理.

> firewall-config 是 firewalld 防火墙配置管理工具的 GUI（图形用户界面）版本.

## 描述

firewalld的字符界面管理工具.

区域就是 firewalld 预先准备了几套防火墙策略集合（策略模板），用户可
以根据生产场景的不同而选择合适的策略集合，从而实现防火墙策略之间的快速切换.

firewalld 中常用的区域名称及测了规则:
- trusted 允许所有的数据包
- home 拒绝流入的流量，除非与流出的流量相关；而如果流量与 ssh、mdns、ipp-client、amba-client 与 dhcpv6-client 服务相关，则允许流量
- internal 等同于 home 区域
- work 拒绝流入的流量，除非与流出的流量数相关；而如果流量与 ssh、ipp-client 与dhcpv6-client 服务相关，则允许流量
- public 拒绝流入的流量，除非与流出的流量相关；而如果流量与 ssh、dhcpv6-client 服务相关，则允许流量
- external 拒绝流入的流量，除非与流出的流量相关；而如果流量与 ssh 服务相关，则允许流量
- dmz 拒绝流入的流量，除非与流出的流量相关；而如果流量与 ssh 服务相关，则允许流量
- block 拒绝流入的流量，除非与流出的流量相关
- drop 拒绝流入的流量，除非与流出的流量相关

使用 firewalld 配置的防火墙策略默认为运行时（Runtime）模式，又称为当前生效模式，而且随着系统的重启会失效. 如果想让配置策略一直存在，就需要使用永`--permanent`参数.

流量转发命令格式为`firewall-cmd --permanent --zone=<区域> --add-forward-port=port=<源端口号>:proto=<协议>:toport=<目标端口号>:toaddr=<目标 IP 地址>`

## 选项:
- --zone=NAME          # 指定 zone
- --permanent          # 永久修改，--reload 后生效
- --timeout=seconds    # 持续效果，到期后自动移除，用于调试，不能与 --permanent 同时使用
- --get-default-zone 查询默认的区域名称
- --set-default-zone=<区域名称> 设置默认的区域，使其永久生效
- --get-zones 显示可用的区域
- --get-services 显示预先定义的服务
- --get-active-zones 显示当前正在使用的区域与网卡名称
- --add-source= 将源自此 IP 或子网的流量导向指定的区域
- --remove-source= 不再将源自此 IP 或子网的流量导向某个指定区域
- --add-interface=<网卡名称> 将源自该网卡的所有流量都导向某个指定区域
- --change-interface=<网卡名称> 将某个网卡与区域进行关联
- --list-all 显示当前区域的网卡配置参数、资源、端口以及服务等信息
- --list-all-zones 显示所有区域的网卡配置参数、资源、端口以及服务等信息
- --add-service=<服务名> 设置默认区域允许该服务的流量
- --add-port=<端口号/协议> 设置默认区域允许该端口的流量, 添加时使用了`--permanent`时, 那么移除时也需要添加它
- --remove-service=<服务名> 设置默认区域不再允许该服务的流量
- --remove-port=<端口号/协议> 设置默认区域不再允许该端口的流量
- --reload 让“永久生效”的配置规则立即生效，并覆盖当前的配置规则
- --panic-on 开启应急状况模式
- --panic-off 关闭应急状况模式
- --add-rich-rule : 富规则也叫复规则， 表示更细致、更详细的防火墙策略配置，它可以针对系统服务、端
口号、源地址和目标地址等诸多信息进行更有针对性的策略配置。它的优先级在所有的防火

## 参考
- [CentOS 7 下使用 Firewall](https://havee.me/linux/2015-01/using-firewalls-on-centos-7.html)

## 例

### 把 firewalld 服务中 eno16777728 网卡的默认区域修改为 external，并在系统重启后同样生效
```bash
# firewall-cmd --permanent --zone=external --change-interface=eno16777728 
success 
# firewall-cmd --get-zone-of-interface=eno16777728
public 
# firewall-cmd --permanent --get-zone-of-interface=eno16777728
external
```

### 添加端口
```
// 永久打开一个端口
firewall-cmd --permanent --zone=public --add-port=80/tcp --add-port=81/tcp
// 永久打开某项服务
firewall-cmd --permanent --zone=public --add-service=http
// 需重新加载防火墙
firewall-cmd --reload
```

### 检查是否生效
```
firewall-cmd --zone=public --query-port=80/tcp
```

### 列出所有的开放端口
```
firewall-cmd --list-all
```

### 删除端口
```
// 永久关闭一个端口
firewall-cmd --permanent --zone=public --remove-port=80/tcp
// 永久关闭某项服务
firewall-cmd --permanent --zone=public --remove-service=http
```

### 进行端口转发
```
firewall-cmd --permanent --add-forward-port=port=80:proto=tcp:toport=8080:toaddr=192.0.2.55
```

### 查看
```
# firewall-cmd --get-default-zone # 查看 firewalld 服务当前所使用的区域
# firewall-cmd --set-default-zone=public # 把 firewalld 服务的当前默认区域设置为 public
# firewall-cmd --get-zone-of-interface=eno16777728 # 查询 eno16777728 网卡在 firewalld 服务中的区域
// 查看已被激活的 Zone 信息
firewall-cmd --get-active-zones
// 查看指定接口的 Zone 信息
firewall-cmd --get-zone-of-interface=eth0
// 查看指定级别的接口
firewall-cmd --zone=public --list-interfaces
// 查看指定级别的所有信息,比如public
firewall-cmd --zone=public --list-all
// 查看所有级别被允许的信息
firewall-cmd --get-service
// 查看重启后所有 Zones 级别中被允许的服务，即永久放行的服务
firewall-cmd --get-service --permanent
# ###启动/关闭 firewalld 防火墙服务的应急状况模式，阻断一切网络连接（当远程控制服务器
时请慎用）
# firewall-cmd --panic-on
# firewall-cmd --panic-off
# ### 查询 public 区域是否允许请求 SSH 和 HTTPS 协议的流量
# firewall-cmd --zone=public --query-service=ssh 
# firewall-cmd --zone=public --query-service=https
# ### 把 firewalld 服务中请求 HTTPS 协议的流量设置为永久允许，并立即生效
# firewall-cmd --zone=public --add-service=https 
# firewall-cmd --permanent --zone=public --add-service=https 
# firewall-cmd --reload
# ### 把 firewalld 服务中请求 HTTP 协议的流量设置为永久拒绝，并立即生效
# firewall-cmd --permanent --zone=public --remove-service=http 
# firewall-cmd --reload
# ### 把在 firewalld 服务中访问 8080 和 8081 端口的流量策略设置为允许，但仅限当前生效
# firewall-cmd --zone=public --add-port=8080-8081/tcp 
# firewall-cmd --zone=public --list-ports 
8080-8081/tcp
# ### 把原本访问本机 888 端口的流量转发到 22 端口，要且求当前和长期均有效
# firewall-cmd --permanent --zone=public --add-forward-port= 
port=888:proto=tcp:toport=22:toaddr=192.168.10.10
# firewall-cmd --reload
# ### 拒绝192.168.10.0/24 网段的所有用户访问本机的 ssh 服务（22 端口）
# firewall-cmd --permanent --zone=public --add-rich-rule=" 
rule family="ipv4" source address="192.168.10.0/24" service name="ssh" reject" # 拒绝
192.168.10.0/24 网段的所有用户访问本机的 ssh 服务
墙策略中也是最高的
# firewall-cmd --reload
# firewall-cmd --complete-reload  # 中断所有连接的重新加载. 两者的区别就是`--reload`无需断开连接，就是firewalld特性之一动态添加规则，`--complete-reload`需要断开连接，类似重启服务
```

## rich rules
修改rule需要`--reload`, 这里图省事没加

```bash
# firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="192.168.0.200" port protocol="tcp" port="80" reject' # 限制IP为192.168.0.200的地址禁止访问80端口
# firewall-cmd --zone=public --list-rich-rules
# firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="192.168.0.200" port protocol="tcp" port="80" accept' # 解除上面被限制的192.168.0.200
# firewall-cmd --permanent --remove-rich-rule='rule family="ipv4" source address="192.168.0.200" port protocol="tcp" port="80" accept' # 同上, **推荐**
# firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="10.0.0.0/24" port protocol="tcp" port="80" reject' # 限制10.0.0.0-10.0.0.255整个段的IP, 禁止访问
```

### 运行、停止、禁用firewalld
```
# systemctl start firewalld
# systemctl status firewalld 或者 firewall-cmd --state
# systemctl disable firewalld
# systemctl stop firewalld
```

## other
TCP Wrappers 是 RHEL 7 系统中默认启用的一款流量监控程序，它能够根据来访主机的地址与本机的目标服务程序作出允许或拒绝的操作.

它是能允许或禁止 Linux 系统提供服务(比如ssh)的防火墙.