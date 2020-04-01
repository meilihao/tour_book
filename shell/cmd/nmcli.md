# nmcli
参考:
- [linux中nmcli命令使用及网络配置](https://www.cnblogs.com/djlsunshine/p/9733182.html)
- [基于RHEL8/CentOS8的网络IP配置详解](https://zhuanlan.zhihu.com/p/56892392)

[nmcli命令与配置文件对应关系](/misc/img/shell/1482552-20180930195218267-1123602572.png)

> 在rhel8上默认不安装network.service(等同废弃)，因此只能通过NM进行网络配置，如果未开启NM，则无法使用网络.

> nmtui是nmcli的字符gui.

在nmcli中有2个命令最为常用：
- nmcli connection : 连接，可理解为配置文件，相当于ifcfg-ethX. 可以简写为nmcli c
- nmcli device : 设备，可理解为实际存在的网卡（包括物理网卡和虚拟网卡）.可以简写为nmcli d

在NM里，有2个维度：连接（connection）和设备（device），这是多对一的关系. 想给某个网卡配ip，首先NM要能纳管这个网卡, 设备里存在的网卡（即nmcli d可以看到的），就是NM纳管的. 接着，可以为一个设备配置多个连接（即nmcli c可以看到的），每个连接可以理解为一个ifcfg配置文件. 同一时刻，一个设备只能有一个连接活跃. 可以通过nmcli c up切换连接.

connection有2种状态：
- 活跃（带颜色字体）：表示当前该connection生效
- 非活跃（正常字体）：表示当前该connection不生效

device有4种常见状态：
- connected：已被NM纳管，并且当前有活跃的connection
- disconnected：已被NM纳管，但是当前没有活跃的connection
- unmanaged：未被NM纳管
- unavailable：不可用，NM无法纳管，通常出现于网卡link为down的时候（比如ip link set ethX down）

## example
```bash
# 查看ip（类似于ifconfig、ip addr）
nmcli

# 创建connection，配置静态ip（等同于配置ifcfg，其中BOOTPROTO=none，并ifup启动）
#  type ethernet：创建连接时候必须指定类型，类型有很多，可以通过`nmcli c add type -h`查看
# con-name ethX ifname ethX：第一个ethX表示连接（connection）的名字，这个名字可以任意定义，无需和网卡名相同；第二个ethX表示网卡名，这个ethX必须是在nmcli d里能看到的
# ipv4.addr '192.168.1.100/24,192.168.1.101/32'：配置2个ip地址，分别为192.168.1.100/24和192.168.1.101/32
# 如果这是为ethX创建的第一个连接，则自动生效；如果此时已有连接存在，则该连接不会自动生效，可以执行nmcli c up ethX来切换生效
nmcli c add type ethernet con-name ethX ifname ethX ipv4.addr 192.168.1.100/24 ipv4.gateway 192.168.1.1 ipv4.method manual

# 创建connection，配置动态ip（等同于配置ifcfg，其中BOOTPROTO=dhcp，并ifup启动）
nmcli c add type ethernet con-name ethX ifname ethX ipv4.method auto

# 修改ip（非交互式）
nmcli c modify ethX ipv4.addr '192.168.1.200/24'
nmcli c up ethX

# 修改ip（交互式）
nmcli c edit ethX
nmcli> goto ipv4.addresses
nmcli ipv4.addresses> change
Edit 'addresses' value: 192.168.1.200/24
Do you also want to set 'ipv4.method' to 'manual'? [yes]: yes
nmcli ipv4> save
nmcli ipv4> activate
nmcli ipv4> quit

# 启用connection（相当于ifup）
nmcli c up ethX

# 停止connection（相当于ifdown）
nmcli c down

# 删除connection（类似于ifdown并删除ifcfg）
nmcli c delete ethX

# 查看connection列表
# - 第一列是connection名字，简称con-name（注意con-name不是网卡名）
# - 第二列是connection的UUID
# - 最后一列才是网卡名（标准说法叫device名），可通过nmcil d查看device
nmcli c show

# 查看connection详细信息
nmcli c show ethX

# 重载所有ifcfg或route到connection（不会立即生效）
nmcli c reload

# 重载指定ifcfg或route到connection（不会立即生效）
nmcli c load /etc/sysconfig/network-scripts/ifcfg-ethX
nmcli c load /etc/sysconfig/network-scripts/route-ethX

# 立即生效(刷新)connection，有3种方法
nmcli c up ethX
nmcli d reapply ethX
nmcli d connect ethX # 由NM对指定网卡进行管理，同时刷新该网卡对应的活跃connection（如果之前有修改过connection配置）；如果有connection但是都处于非活跃状态，则自动选择一个connection并将其活跃；如果没有connection，则自动生成一个并将其活跃

# 查看device列表
nmcli d

# 查看所有device详细信息
nmcli d show

# 查看指定device的详细信息
nmcli d show ethX

# 激活网卡
nmcli d connect ethX

# 关闭无线网络（NM默认启用无线网络）
nmcli r all off

# 查看NM纳管状态
nmcli n

# 开启NM纳管
nmcli n on

# 关闭NM纳管（谨慎执行）
nmcli n off

# 监听事件
nmcli m

# 查看NM本身状态
nmcli

# 检测NM是否在线可用
nm-online

# 让NM**暂时(即重启恢复)**不管理指定网卡，此操作不会变更实际网卡的link状态，只会使对应的connection变成非活跃
nmcli d disconnect ethX
```

## 其他
ifcfg和NM connection的关联：虽然network.service被废弃了，但是redhat为了兼容传统的ifcfg，通过NM进行网络配置时候，会自动将connection同步到ifcfg配置文件中, 也可以通过nmcli c reload或者nmcli c load /etc/sysconfig/network-scripts/ifcfg-ethX的方式来让NM读取ifcfg配置文件到connection中. 因此ifcfg和connection是一对一的关系，connection和device是多对一的关系.