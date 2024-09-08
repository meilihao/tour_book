# bridge
ref:
- [bridge-utils(brctl) -> iproute2(ip/bridge)](https://qiita.com/s1061123/items/54c9b4c001877135c4ff)

## brctl
```bash
brctl show # bridge id: 显示桥接设备的桥接 ID。桥接 ID 是一个唯一标识符，用于生成树协议（STP）中识别和管理桥接设备。桥接 ID 是由桥接设备的 MAC 地址和优先级组成的; STP enabled: 显示生成树协议（STP）是否启用。STP 是一种网络协议，用于防止网络环路和广播风暴;interfaces: 列出连接到桥接设备的所有接口. 这些接口可以是物理网络接口、虚拟网络接口或其他网络设备
brctl showmacs br0 # 查看桥学到的所有 MAC 地址. `port no`: 该mac所在的端口号. 端口号对应于桥接设备的物理接口, 需结合brctl showstp br0确定; `is loal`: yes表示这个 MAC 地址是桥接设备本地接口的地址。即这个 MAC 地址不是通过桥接设备的学习功能学到的，而是桥接设备接口自身的地址. 非 local（如 forward, learn）: 表示这些 MAC 地址是桥接设备学习到的，并且这些地址是通过桥接设备的端口转发流量的; `ageing timer`, MAC 地址在桥接设备的 MAC 地址表中存在的时间（通常以秒为单位）. 这个字段表示自 MAC 地址被学习以来的时间, 如果时间过长，可能意味着这个 MAC 地址长时间没有活动，也可能被从表中删除.
brctl stp br0 off # 多个以太网桥可以工作在一起组成一个更大的网络，利用 802.1d 协议在两个网络之间寻找最短路径，STP 的作用是防止以太网桥之间形成回路，如果确定只有一个网桥，则可以关闭STP
bridge fdb show # 一个用于显示 Linux 桥接设备 (bridge) 的 MAC 地址表的命令. `mac`: 网络设备的mac; `dev`: 与该 MAC 地址关联的网络设备接口; vlan: 显示与该 MAC 地址关联的 VLAN ID;master: 显示当前 MAC 地址的状态，通常是 master。表示这个 MAC 地址在桥接设备的转发数据库中作为主要或活动地址. self: 表示该 MAC 地址属于本地设备自身的接口。即MAC 地址是桥接设备本身的地址. permanent: 表示该 MAC 地址是永久的(相对于学习来的mac)，不会被老化或自动删除。permanent MAC 地址通常是通过静态配置添加到桥接设备的，不会因时间过长而从 MAC 地址表中被清除.
```