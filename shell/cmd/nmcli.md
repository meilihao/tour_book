# nmcli
RHEL 和 CentOS 系统默认使用 NetworkManager 来提供网络服务，这是一种动态管理网
络配置的守护进程，能够让网络设备保持连接状态. 其使用 nmcli 命令来管理 Network 
Manager 服务. 即nmcli 是一款基于命令行的网络配置工具.

**使用 nmcli 命令配置过的网络会话是永久生效的.**

nmtui是nmcli的gui.

> 在RHEL 8上，已经弃用Network.service，因此只能通过NetworkManager.service进行网络配置，包括动态IP和静态IP

> [cockpit](https://www.kclouder.cn/howtocockpit/)：redhat自带的基于web图形界面的"驾驶舱"工具，具有dashborad和基础管理功能, 系统管理员可以执行诸如存储管理、网络配置、检查日志、管理容器等任务.

## 配置文件
`/etc/NetworkManager/system-connections`

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

## FAQ
### 管理工具选择
在CentOS8中默认使用NetworkManager守护进程来监控和管理网络设置

> 在RHEL 7上，同时支持Network.service和NetworkManager.service（简称NM），默认情况下，这两个服务都开启.

Ubuntu 17.10和更高Ubuntu版本使用Netplan作为默认网络管理工具, Netplan支持两个渲染器：NetworkManager和Systemd-networked. 且Ubuntu 20.04 LTS 在网络管理上相比较于18.04有很大的不同, 其不再使用networking服务(`/etc/network/interfaces`)进行管理网络了.

Ubuntu Server 20.04使用nmcli配置网络连接:
```bash
sudo apt install network-manager
sudo vim /etc/NetworkManager/NetworkManager.conf # 编辑NetworkManager.conf文件并启用network-manager
managed=false -> managed=true
sudo systemctl restart NetworkManager.service
sudo systemctl enable NetworkManager.service
sudo vim /etc/netplan/00-installer-config.yaml # 编辑/etc/netplan/*.yaml(01-network-manager-all.yml/50-cloud-init.yaml等不同环境命名下的netplan配置文件)，使用下面的内容替换它
network:
  version: 2
  renderer: NetworkManager
sudo netplan apply
```

netplan2NM.sh:
```bash
#!/usr/bin/env bash

# netplan2NM.sh
# Ubuntu server 20.04  Change from netplan to NetworkManager for all interfaces

echo 'Changing netplan to NetowrkManager on all interfaces'
# backup existing yaml file
cd /etc/netplan
cp 01-netcfg.yaml 01-netcfg.yaml.bak

# re-write the yaml file
cat << EOF > /etc/netplan/01-netcfg.yaml
# This file describes the network interfaces available on your system
# For more information, see netplan(5).
network:
  version: 2
  renderer: NetworkManager
EOF

# setup netplan for NM
# netplan generate
netplan apply
# make sure NM is running
systemctl enable NetworkManager.service
systemctl restart NetworkManager.service

echo 'Done!'
```