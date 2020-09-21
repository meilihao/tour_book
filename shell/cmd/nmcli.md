# nmcli
RHEL 和 CentOS 系统默认使用 NetworkManager 来提供网络服务，这是一种动态管理网
络配置的守护进程，能够让网络设备保持连接状态. 其使用 nmcli 命令来管理 Network 
Manager 服务. 即nmcli 是一款基于命令行的网络配置工具.

**使用 nmcli 命令配置过的网络会话是永久生效的.**

## example
```bash
# nmcli connection show # 查看网络信息或网络状态
# ### nmcli 支持网络会话功能, 便于切换网络, 比如公司和家里
# nmcli connection add con-name company ifname eno16777736 autoconnect no type ethernet ip4 192.168.10.10/24 gw4 192.168.10.1 # autoconnect no 参数设置该网络会话默认不被自动激活，以及用 ip4 及 gw4 参数手动指定网络的 IP 地址
# nmcli connection add con-name house type ethernet ifname eno16777736 # 从外部 DHCP 服务器自动获得 IP 地址，因此不需要进行手动指定
# nmcli connection up house # 回家启用 house网络会话，网卡就能自动通过 DHCP 获取到 IP 地址
```

## bonding
```bash
# vim /etc/sysconfig/network-scripts/ifcfg-eno16777736
TYPE=Ethernet 
BOOTPROTO=none 
ONBOOT=yes 
USERCTL=no 
DEVICE=eno16777736 
MASTER=bond0 
SLAVE=yes 
[root@linuxprobe ~]# vim /etc/sysconfig/network-scripts/ifcfg-eno33554968 
TYPE=Ethernet 
BOOTPROTO=none 
ONBOOT=yes 
USERCTL=no 
DEVICE=eno33554968 
MASTER=bond0 
SLAVE=yes
# vim /etc/sysconfig/network-scripts/ifcfg-bond0 
TYPE=Ethernet 
BOOTPROTO=none 
ONBOOT=yes 
USERCTL=no 
DEVICE=bond0 
IPADDR=192.168.10.10 
PREFIX=24 
DNS=192.168.10.1 
NM_CONTROLLED=no
#  vim /etc/modprobe.d/bond.conf # 时定义网卡以 mode6 模式进行绑定，且出现故障时自动切换的时间为 100 毫秒
alias bond0 bonding 
options bond0 miimon=100 mode=6
# systemctl restart network # 重启网络服务后网卡bonding操作即可成功. 正常情况下只有 bond0 网卡设备才会有 IP 地址等信息
```

常见的网卡bond驱动有三种模式—mode0、mode1和 mode6, 它们的使用的情景是:
- mode0（平衡负载模式）：平时两块网卡均工作，且自动备援，但需要在与服务器本地网卡相连的交换机设备上进行端口聚合来支持绑定技术
- mode1（自动备援模式）：平时只有一块网卡工作，在它故障后自动替换为另外的网卡
- mode6（平衡负载模式）：平时两块网卡均工作，且自动备援，无须交换机设备提供辅助支持