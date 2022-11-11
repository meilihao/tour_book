# nmcli
ref:
- [Managing IP Networking](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/networking_guide/part-managing_ip_networking#doc-wrapper)

RHEL 和 CentOS 系统默认使用 NetworkManager 来提供网络服务，这是一种动态管理网
络配置的守护进程，能够让网络设备保持连接状态. 其使用 nmcli 命令来管理 Network 
Manager 服务. 即nmcli 是一款基于命令行的网络配置工具.

**使用 nmcli 命令配置过的网络会话是永久生效的.**

nmtui是nmcli的terminal gui. nm-connection-editor是gui. [nmstatectl以yaml/json来配置网络](https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/9/html/configuring_and_managing_networking/proc_configuring-a-static-ethernet-connection-using-nmstatectl_configuring-an-ethernet-connection). control-center也可以, 但不支持与 nm-connection-editor 应用程序或 nmcli 实用程序一样多的配置选项.

> 在RHEL 8上，已经弃用Network.service，因此只能通过NetworkManager.service进行网络配置，包括动态IP和静态IP

> [cockpit](https://www.kclouder.cn/howtocockpit/)：redhat自带的基于web图形界面的"驾驶舱"工具，具有dashborad和基础管理功能, 系统管理员可以执行诸如存储管理、网络配置、检查日志、管理容器等任务.

## 配置文件
`/etc/NetworkManager/system-connections`

## example
```bash
# nmcli connection show # 查看网络信息或网络状态
# ### nmcli 支持网络会话功能, 便于切换网络, 比如公司和家里
# --- 设置静态地址
# nmcli connection add con-name company ifname eno16777736 autoconnect no type ethernet ip4 192.168.10.10/24 gw4 192.168.10.1 # autoconnect no 参数设置该网络会话默认不被自动激活，以及用 ip4 及 gw4 参数手动指定网络的 IP 地址
# nmcli connection add con-name house type ethernet ifname eno16777736 # 从外部 DHCP 服务器自动获得 IP 地址，因此不需要进行手动指定
# nmcli connection up house # 回家启用 house网络会话，网卡就能自动通过 DHCP 获取到 IP 地址
# nmcli -p dev # 查看设备状态
# nmcli connection modify Example-Connection ipv4.addresses 192.0.2.1/24 # 设置ipv4
# nmcli connection modify Example-Connection ipv6.addresses 2001:db8:1::1/64 # 设置ipv6
# nmcli connection modify Example-Connection ipv4.method manual # 设置为手动连接
# nmcli connection modify Example-Connection ipv6.method manual
# nmcli connection modify Example-Connection ipv4.gateway 192.0.2.254 # 设置默认网关
# nmcli connection modify Example-Connection ipv6.gateway 2001:db8:1::fffe
# nmcli connection modify Example-Connection ipv4.dns "192.0.2.200" # 设置dns, 可用`host xxx`来验证
# nmcli connection modify Example-Connection ipv6.dns "2001:db8:1::ffbb"
# nmcli connection modify Example-Connection ipv4.dns-search example.com # 设置搜索域
# nmcli connection modify Example-Connection ipv6.dns-search example.com
# --- 交互式编辑
# nmcli connection edit type ethernet con-name Example-Connection
nmcli> set connection.interface-name enp7s0 # 设置网络接口
nmcli> set ipv4.addresses 192.0.2.1/24 # 设置 IPv4 地址
nmcli> set ipv6.addresses 2001:db8:1::1/64 # 设置 IPv6 地址
nmcli> set ipv4.method manual # 将 IPv4 和 IPv6 连接方法设置为 manual
nmcli> set ipv6.method manual
nmcli> set ipv4.gateway 192.0.2.254 # 设置 IPv4 和 IPv6 默认网关
nmcli> set ipv6.gateway 2001:db8:1::fffe
nmcli> set ipv4.dns 192.0.2.200  # 设置 IPv4 和 IPv6 DNS 服务器地址. 要设置多个 DNS 服务器，以空格分隔并用引号括起来
nmcli> set ipv6.dns 2001:db8:1::ffbb
nmcli> set ipv4.dns-search example.com # 为 IPv4 和 IPv6 连接设置 DNS 搜索域
nmcli> set ipv6.dns-search example.com
nmcli> save persistent # 保存并激活连接
Saving the connection with 'autoconnect=yes'. That might result in an immediate activation of the connection.
Do you still want to save? (yes/no) [yes] yes
nmcli> quit # 退出
# nmcli device status # 显示设备和连接的状态
# nmcli connection show Example-Connection # 显示连接配置集的所有设置
# --- 设置动态地址
# nmcli connection add con-name Example-Connection ifname enp7s0 type ethernet
# nmcli connection modify Example-Connection ipv4.dhcp-hostname Example ipv6.dhcp-hostname Example # （可选）在使用 Example-Connection 配置文件时，更改 NetworkManager 发送给 DHCP 服务器的主机名
# nmcli connection modify Example-Connection ipv4.dhcp-client-id client-ID # （可选）在使用 Example-Connection 配置文件时，更改 NetworkManager 发送给 IPv4 DHCP 服务器的客户端 ID. 对于 IPv6 ，没有 dhcp-client-id 参数。要为 IPv6 创建一个标识符，请配置 dhclient 服务.
# --- 按接口名称使用单一连接配置文件配置多个以太网接口, [实现主机在具有动态 IP 地址分配的以太网之间漫游](https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/9/html/configuring_and_managing_networking/proc_configuring-multiple-ethernet-interfaces-using-a-single-connection-profile-by-interface-name_configuring-an-ethernet-connection), 也支持通过PCI ID来匹配`match.path "pci-0000:07:00.0 pci-0000:08:00.0"`
nmcli connection add con-name Example connection.multi-connect multiple match.interface-name enp* type ethernet
```

## dhcp
默认情况下，NetworkManager 使用其内部的 DHCP 客户端. 但是，如果使用内置客户端的 DHCP 客户端，也可以将 NetworkManager 配置为使用 dhclient. 请注意，RHEL 不提供 dhcpcd，因此 NetworkManager 无法使用这个客户端

```bash
# vim /etc/NetworkManager/conf.d/dhcp-client.conf
[main]
dhcp=dhclient
# dnf install dhcp-client
# systemctl restart NetworkManager
# cat /var/log/messages |grep dhclient # 验证
Apr 26 09:54:19 server NetworkManager[27748]: <info>  [1650959659.8483] dhcp-init: Using DHCP client 'dhclient'
# nmcli connection modify connection_name ipv4.dhcp-timeout 30 ipv6.dhcp-timeout 30 # 设置 ipv4.dhcp-timeout 和 ipv6.dhcp-timeout 属性
# nmcli connection modify connection_name ipv4.may-fail <value> # 配置如果网络管理器（NetworkManager）在超时前没有接收 IPv4 地址时的行为: yes, 尝试ipv6, 除非ipv6禁用或未配置, 一旦ipv6成功, 则不再尝试ipv4; no，连接会被停止. 如果启用了连接的 autoconnect, 会根据 autoconnect-retries 属性中设置的值尝试多次激活连接. 如果连接仍然无法获得 DHCP 地址，则自动激活会失败。请注意，5 分钟后，自动连接过程会再次启动，从 DHCP 服务器获取 IP 地址.
# nmcli connection modify connection_name ipv6.may-fail <value> # 配置如果网络管理器（NetworkManager）在超时前没有接收 IPv6 地址时的行为. 行为类似同上
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

## lib
- [ushiboy/nmcli](https://github.com/ushiboy/nmcli)
- [Wifx/gonetworkmanager](https://github.com/Wifx/gonetworkmanager)

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