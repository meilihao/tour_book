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
根据连接配置文件的目的，将其保存在以下目录中：
- `/etc/NetworkManager/system-connections` ：用户创建的持久配置文件的通用位置，也可以对其进行编辑. NetworkManager 将它们自动复制到 `/etc/NetworkManager/system-connections/`
- `/run/NetworkManager/system-connections` ：用于在重启系统时自动删除的临时配置文件
- `/usr/lib/NetworkManager/system-connections` ：用于预先部署的不可变的配置文件. 当使用 NetworkManager API 编辑此类配置文件时，NetworkManager 会将此配置文件复制到持久性存储或临时存储中

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

## wifi
```bash
# nmcli radio wifi on # 使用wifi前需要启用wifi radio
# nmcli device wifi list # 罗列可用的wifi. ssid=`--`表示, 该wifi不会广播ssid
# nmcli device wifi connect Office --ask # 连接wifi. 该方法是互方式输入密码
Password: wifi-password
# nmcli device wifi connect Office wifi-password # 连接wifi, 非交互式输入密码
# ---  使用wifi时手动指定ip
# nmcli connection modify Office ipv4.method manual ipv4.addresses 192.0.2.1/24 ipv4.gateway 192.0.2.254 ipv4.dns 192.0.2.200 ipv4.dns-search example.com # 配置 IPv4
# nmcli connection modify Office ipv6.method manual ipv6.addresses 2001:db8:1::1/6 # 配置 IPv6
# nmcli connection up Office # 配置ip后激活即可
# nmcli connection show --active # 验证
```

## 配置默认网关
默认网关是在没有其他路由与数据包的目的地匹配时转发网络数据包的路由器. 在本地网络中，默认网关通常是与距离互联网有一跳的主机.

通常可会在创建连接时设置默认网关.

注意:
1. 只有在很少情况下（比如使用多路径 TCP 时），在主机上需要多个默认网关。在大多数情况下，只配置一个默认网关，来避免意外的路由行为或异步路由问题.
1. 要将流量路由到不同的互联网提供商，请使用基于策略的路由，而不是多个默认网关.

```bash
# --- 非交互式设置
# nmcli connection modify example ipv4.gateway "192.0.2.1"
# nmcli connection modify example ipv6.gateway "2001:db8:1::1"
# nmcli connection up example # 重启网络连接以使更改生效
# --- 交互式设置
# nmcli connection edit example
nmcli> set ipv4.gateway 192.0.2.1
nmcli> set ipv6.gateway 2001:db8:1::1
nmcli> print
...
ipv4.gateway:                           192.0.2.1
...
ipv6.gateway:                           2001:db8:1::1
...
nmcli> save persistent # 保存配置
nmcli> activate example # 重启网络连接以使更改生效
nmcli> quit
# --- 验证
# ip -4 route
default via 192.0.2.1 dev example proto static metric 100
# ip -6 route
default via 2001:db8:1::1 dev example proto static metric 100 pref medium
# nmcli -f GENERAL.CONNECTION,IP4.GATEWAY,IP6.GATEWAY device show enp1s0
GENERAL.CONNECTION:      Corporate-LAN
IP4.GATEWAY:             192.0.2.1
IP6.GATEWAY:             2001:db8:1::1
# nmcli connection modify connection_name ipv4.never-default yes ipv6.never-default yes # 如果连接使用动态 IP 配置，则需配置 NetworkManager 不使用该连接作为 IPv4 和 IPv6 连接的默认路由. 将 ipv4.never-default 和 ipv6.never-default 设为 yes，会自动从连接配置文件中删除相应协议默认网关的 IP 地址
```

## 配置静态路由
路由可确保可以在相互连接的网络间发送和接收流量。在较大环境中，管理员通常配置服务以便路由器可以动态地了解其他路由器。在较小的环境中，管理员通常会配置静态路由，以确保流量可以从一个网络到下一个网络访问.

```bash
# nmcli connection modify example +ipv4.routes "198.51.100.0/24 192.0.2.10, 203.0.113.0/24 192.0.2.10" # 要在一个步骤中设置多个路由，使用逗号分隔单个路由传递给该命令。例如，要将路由添加到 198.51.100.0/24 和 203.0.113.0/24 网络，它们都通过 192.0.2.10 网关路由
# nmcli connection modify example +ipv6.routes "2001:db8:2::/64 2001:db8:1::10"
# nmcli connection up example
```

## [配置基于策略的路由以定义其他路由](https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/9/html/configuring_and_managing_networking/configuring-policy-based-routing-to-define-alternative-routes_configuring-and-managing-networking)

默认情况下，RHEL 中的内核决定使用路由表根据目标地址转发网络数据包。基于策略的路由允许您配置复杂的路由场景.

```bash
# nmcli connection add type ethernet con-name Provider-A ifname enp7s0 ipv4.method manual ipv4.addresses 198.51.100.1/30 ipv4.gateway 198.51.100.2 ipv4.dns 198.51.100.200 connection.zone external # connection.zone firewalld_zone表示将网络接口分配给定义的 firewalld 区域. 请注意，firewalld 会为分配给 external 区域的接口自动启用伪装
# nmcli connection add type ethernet con-name Provider-B ifname enp1s0 ipv4.method manual ipv4.addresses 192.0.2.1/30 ipv4.routes "0.0.0.0/0 192.0.2.2 table=5000" connection.zone external # 此命令使用 ipv4.routes 参数而不是 ipv4.gateway 来设置默认网关。这需要将这个连接的默认网关分配给不同于默认的路由表(5000)。当连接被激活时，NetworkManager 会自动创建这个新的路由表
# nmcli connection add type ethernet con-name Internal-Workstations ifname enp8s0 ipv4.method manual ipv4.addresses 10.0.0.1/24 ipv4.routes "10.0.0.0/24 table=5000" ipv4.routing-rules "priority 5 from 10.0.0.0/24 table 5000" connection.zone trusted #  此命令使用 ipv4.routes 参数将静态路由添加到 ID 为 5000 的路由表中。10.0.0.0/24 子网的这个静态路由使用到供应商 B 的本地网络接口的 IP 地址(192.0.2.1)来作为下一跳。另外，命令使用 ipv4.routing-rules 参数来添加优先级为 5 的路由规则，该规则将来自 10.0.0.0/24 子网的流量路由到表 5000。低的值具有更高的优先级
# nmcli connection add type ethernet con-name Servers ifname enp9s0 ipv4.method manual ipv4.addresses 203.0.113.1/24 connection.zone trusted
# ip route list table 5000 # 验证
```

## bonding
**推荐使用team**

可在不同类型的设备中创建绑定，例如：
- 物理和虚拟以太网设备
- 网络桥接
- 网络团队（team）
- VLAN 设备

active-backup、balance-tlb 和 balance-alb 模式不需要网络交换机的任何具体配置。然而，其他绑定模式需要配置交换机来聚合链接. 例如，对于模式 0、2 和 3，Cisco 交换机需要 EtherChannel ，但对于模式 4，需要链接聚合控制协议(LACP)和 EtherChannel.

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
# cat /proc/net/bonding/bond0 # 显示绑定状态
```

常见的网卡bond驱动有三种模式—mode0、mode1和 mode6, 它们的使用的情景是:
- mode0（平衡负载模式）：平时两块网卡均工作，且自动备援，但需要在与服务器本地网卡相连的交换机设备上进行端口聚合来支持绑定技术
- mode1（自动备援模式）：平时只有一块网卡工作，在它故障后自动替换为另外的网卡
- mode6（平衡负载模式）：平时两块网卡均工作，且自动备援，无须交换机设备提供辅助支持

## 配置 VLAN
VLAN 是物理网络中的一个逻辑网络。当 VLAN 接口通过接口时，VLAN 接口标签带有 VLAN ID 的数据包，并删除返回的数据包的标签.

[前提](https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/9/html/configuring_and_managing_networking/configuring-vlan-tagging_configuring-and-managing-networking):
1. 用作虚拟 VLAN 接口的父接口支持 VLAN
1. 主机连接的switch支持vlan
1. 其他

  1. 绑定的端口是up
  1. 这个绑定没有使用 fail_over_mac=follow 选项进行配置。VLAN 虚拟设备无法更改其 MAC 地址以匹配父设备的新 MAC 地址。在这种情况下，流量仍会与不正确的源 MAC 地址一同发送
  1. 这个绑定通常不会预期从 DHCP 服务器或 IPv6 自动配置获取 IP 地址。在创建绑定时通过设置 ipv4.method=disable 和 ipv6.method=ignore 选项来确保它。否则，如果 DHCP 或 IPv6 自动配置在一段时间后失败，接口可能会关闭

```bash
# nmcli connection add type vlan con-name vlan10 ifname vlan10 vlan.parent enp1s0 vlan.id 10 # 创建一个使用 enp1s0 作为其父接口，使用 VLAN ID 10 标记数据包，名为 vlan10 的 VLAN 接口. VLAN 必须在范围 0 到 4094 之间
# nmcli connection modify vlan10 ethernet.mtu 2000 # 默认情况下，VLAN 连接会继承上级接口的最大传输单元（MTU）
# --- 配置 VLAN 设备的 IP 设置. 如果要使用这个 VLAN 设备作为其它设备的端口，请跳过这一步
# nmcli connection modify vlan10 ipv4.addresses '192.0.2.1/24' # 配置ipv4
# nmcli connection modify vlan10 ipv4.gateway '192.0.2.254'
# nmcli connection modify vlan10 ipv4.dns '192.0.2.253'
# nmcli connection modify vlan10 ipv4.method manual
# nmcli connection modify vlan10 ipv6.addresses '2001:db8:1::1/32' # 配置ipv6
# nmcli connection modify vlan10 ipv6.gateway '2001:db8:1::fffe'
# nmcli connection modify vlan10 ipv6.dns '2001:db8:1::fffd'
# nmcli connection modify vlan10 ipv6.method manual
# nmcli connection up vlan10 # 激活
# ip -d addr show vlan10 # 验证
```

## vxlan
```bash
# nmcli connection add con-name Example ifname enp1s0 type ethernet
# nmcli connection modify Example ipv4.addresses 198.51.100.2/24 ipv4.method manual ipv4.gateway 198.51.100.254 ipv4.dns 198.51.100.200 ipv4.dns-search example.com # 配置 IPv4. 如果网络使用 DHCP，请跳过这一步
# nmcli connection up Example # 激活 Example 连接
# --- 要使虚拟可扩展局域网(VXLAN)对虚拟机(VM)不可见, 需在主机上创建一个网桥，并将 VXLAN 附加给网桥
# nmcli connection add type bridge con-name br0 ifname br0 ipv4.method disabled ipv6.method disabled # 此命令在网桥设备上设置没有 IPv4 和 IPv6 地址，因为此网桥在 2 层工作
# nmcli connection add type vxlan slave-type bridge con-name br0-vxlan10 ifname vxlan10 id 10 local 198.51.100.2 remote 203.0.113.1 master br0 #  ID 10 ：设置 VXLAN 标识符; local 198.51.100.2 ：设置传出数据包的源 IP 地址; remote 203.0.113.1 ：当目的地链路层地址在 VXLAN 设备转发数据库中未知时，设置要在传出数据包中使用的单播或多播IP地址; Master br0 ：将要创建的此 VXLAN 连接设为 br0 连接中的一个端口; ipv4.method disabled 和 ipv6.method disabled: 在网桥上禁用 IPv4 和 IPv6. 默认情况下，NetworkManager 使用 8472 作为目的地端口。如果目的地端口不同，还要将 destination-port <port_number> 选项传给命令.
# nmcli connection up br0
# firewall-cmd --permanent --add-port=8472/udp
# firewall-cmd --reload
# bridge fdb show dev vxlan10 # 显示转发表, 验证配置 
# --- 在带有现有网桥的 libvirt 中创建一个虚拟网络
# vim ~/vxlan10-bridge.xml
<network>
 <name>vxlan10-bridge</name>
 <forward mode="bridge" />
 <bridge name="br0" />
</network>
# virsh net-define ~/vxlan10-bridge.xml # 使用 ~/vxlan10-bridge.xml 文件来在 libvirt 中创建一个新的虚拟网络
# virsh net-start vxlan10-bridge # 启动 vxlan10-bridge 虚拟网络
# virsh net-autostart vxlan10-bridge # 将 vxlan10-bridge 虚拟网络配置为在 libvirtd 服务启动时自动启动
# virsh net-list # 显示虚拟网络列表
# --- libvirt使用vxlan
# virt-install ... --network network:vxlan10-bridge # 创建新的虚拟机，并将其配置为使用 vxlan10-bridge 网络
# virt-xml VM_name --edit --network network=vxlan10-bridge # 更改现有虚拟机的网络设置
```

## bridge
网络桥接是一个链路层设备，它可根据 MAC 地址列表转发网络间的流量。网桥通过侦听网络流量并了解连接到每个网络的主机来构建 MAC 地址表.

当配置网桥时，网桥被称为 controller，其使用的设备为 ports.

可在不同类型的设备中创建桥接，例如：
- 物理和虚拟以太网设备
- bond
- team
- VLAN 设备

注意: 由于 IEEE 802.11 标准指定在 Wi-Fi 中使用 3 个地址帧以便有效地使用随机时间，因此无法通过 Ad-Hoc 或者 Infrastructure 模式中的 Wi-Fi 网络配置网桥

```bash
# nmcli connection add type bridge con-name bridge0 ifname bridge0 # 创建bridge
# nmcli connection modify bond0 master bridge0 # 将接口分配给网桥
# --- 配置网桥的 IP 设置。如果要使用这个网桥作为其它设备的端口，请跳过这一步
# nmcli connection modify bridge0 ipv4.addresses '192.0.2.1/24' # 配置ipv4
# nmcli connection modify bridge0 ipv4.gateway '192.0.2.254'
# nmcli connection modify bridge0 ipv4.dns '192.0.2.253'
# nmcli connection modify bridge0 ipv4.dns-search 'example.com'
# nmcli connection modify bridge0 ipv4.method manual
# nmcli connection modify bridge0 ipv6.addresses '2001:db8:1::1/64' # 配置ipv6
# nmcli connection modify bridge0 ipv6.gateway '2001:db8:1::fffe'
# nmcli connection modify bridge0 ipv6.dns '2001:db8:1::fffd'
# nmcli connection modify bridge0 ipv6.dns-search 'example.com'
# nmcli connection modify bridge0 ipv6.method manual
# nmcli connection modify bridge0 bridge.priority '16384' # 配置网桥的其他属性. 例如将 bridge0 的生成树协议(STP)优先级设为 16384
# nmcli connection modify bridge0 connection.autoconnect-slaves 1 # 启用网桥连接的 connection.autoconnect-slaves. 当激活连接的任何端口时，NetworkManager 也会激活网桥，但不会激活它的其它端口, 配置connection.autoconnect-slaves后自动启用所有端口
# nmcli connection up bridge0
# --- 验证
# ip link show master bridge0 # 显示作为特定网桥端口的以太网设备的链接状态
# bridge link show # 显示作为任意网桥设备端口的以太网设备状态
```

## team
可以在不同类型的设备中创建team, 例如：
- 物理和虚拟以太网设备
- 网络绑定
- 网络桥接
- VLAN 设备

可用的 runner 如下：
- broadcast ：转换所有端口上的数据。
- roundrobin ：依次转换所有端口上的数据。
- activebackup ：转换一个端口上的数据，而其他端口上的数据则作为备份保留。
- loadbalance：转换所有具有活跃的 Tx 负载均衡和基于 Berkeley 数据包过滤器(BPF)的 Tx 端口选择器的端口上的数据。
- random ：转换随机选择的端口上的数据。
- lacp ：实现 802.3ad 链路聚合控制协议(LACP)。 

teamd 服务使用链路监视器来监控从属设备的状态。可用的 link-watchers 如下：
- ethtool ：libteam 库使用 ethtool 工具来监视链接状态的变化。这是默认的 link-watcher, **推荐**
- arp_ping: libteam 库使用 arp_ping 工具来监控使用地址解析协议(ARP)的远端硬件地址是否存在。
- nsna_ping: 在 IPv6 连接上，libteam 库使用来自 IPv6 邻居发现协议的邻居广告和邻居请求功能来监控邻居接口的存在。 

注意: 每个 runner 都可以使用任何链接监视器，但 lacp 除外。此 runner 只能使用 ethtool 链接监视器

```bash
# dnf install teamd NetworkManager-team
# nmcli connection show team-team0 # 显示 team-team0
# teamdctl team0 config dump actual > /tmp/team0.json # 将 team0 设备的配置导出到 JSON 文件中
# --- 删除 team-team0 连接配置集以及相关端口的配置集
# nmcli connection delete team-team0
# nmcli connection delete team-team0-port1
# nmcli connection delete team-team0-port2
# nmcli connection add type team con-name team0 ifname team0 team.runner activebackup # 创建team
# nmcli connection modify team0 team.link-watchers "name=ethtool delay-up=2500, name=arp_ping source-host=192.0.2.1 target-host=192.0.2.2" # 设置链接监视器. 不同的连接监视器以逗号分隔.
# --- 为team分配端口接口, 下面二选一
# nmcli connection add type ethernet slave-type team con-name team0-port1 ifname enp7s0 master team0 # 没有配置要分配给team的接口，则为其创建新的连接配置集
# nmcli connection modify enp7s0 master team0 # 将现有的连接配置文件分配给team
# --- 配置团队的 IP 设置. 如果要使用这个team作为其它设备的端口，请跳过这一步
# nmcli connection modify team0 ipv4.addresses '192.0.2.1/24' # 配置 IPv4
# nmcli connection modify team0 ipv4.gateway '192.0.2.254'
# nmcli connection modify team0 ipv4.dns '192.0.2.253'
# nmcli connection modify team0 ipv4.dns-search 'example.com'
# nmcli connection modify team0 ipv4.method manual
# nmcli connection modify team0 ipv6.addresses '2001:db8:1::1/64' # 配置 IPv6
# nmcli connection modify team0 ipv6.gateway '2001:db8:1::fffe'
# nmcli connection modify team0 ipv6.dns '2001:db8:1::fffd'
# nmcli connection modify team0 ipv6.dns-search 'example.com'
# nmcli connection modify team0 ipv6.method manual
# nmcli connection up team0
# teamdctl team0 state # 验证
```

## 更改主机名
```bash
# nmcli general hostname # 显示当前主机名
# nmcli general hostname new-hostname.example.com # NetworkManager 自动重启 systemd-hostnamed 以激活新名称
```

## [端口镜像](https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/9/html/configuring_and_managing_networking/assembly_port-mirroring_configuring-and-managing-networking)
可以使用端口镜像复制从一个网络设备传输到另一个网络设备的入站和出站网络流量。管理员使用端口镜像来监控网络流量，并收集网络流量，用于：
- 调试网络问题并调整网络流
- 检查并分析网络流量，来对网络问题进行故障排除
- 检测入侵

## 配置 NetworkManager 以忽略某些设备
默认情况下，NetworkManager 管理除 lo （环回）设备以外的所有设备。但是，可以将某些设备设置为 非受管设备 来配置网络管理器(NetworkManager)忽略这些设备.

```bash
# vim /etc/NetworkManager/conf.d/99-unmanaged-devices.conf # 永久设置
unmanaged-devices=interface-name:interface_1;interface-name:interface_2;...
# systemctl reload NetworkManager
# nmcli device status # 设备旁边的 unmanaged 状态表示 NetworkManager 没有管理该设备
# nmcli device set enp1s0 managed no # 临时设置
```

## 配置网络设备以接受来自所有 MAC 地址的流量
网络设备通常会拦截和读取编程的控制器接收的数据包, 可以在虚拟交换机或端口组层面上，将网络设备配置为接受来自所有 MAC 地址的流量。

可以使用这个网络模式来：
- 诊断网络连接问题
- 出于安全原因监控网络活动
- 拦截传输中的私有数据或网络中的入侵

```bash
# --- iproute2
# ip link set enp1s0 promisc on # 启用accept-all-mac-addresses
# ip link show enp1s0 # 设备描述中的 PROMISC 标志表示启用了该模式
# --- nmcli
# nmcli connection modify enp1s0 ethernet.accept-all-mac-addresses yes
# nmcli connection up enp1s0
# nmcli connection show enp1s0 # 802-3-ethernet.accept-all-mac-addresses: true 表示该模式已启用
```

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