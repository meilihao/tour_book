# iptables

iptables会把配置好的防火墙策略交由内核层面的 netfilter 网络过滤器来处理, 已被nftables取代.

iptables 服务把用于处理或过滤流量的策略条目称之为规则，多条规则可以组成一个规则链，而规则链则依据数据包处理位置的不同进行分类，具体如下：
1. 在进行路由选择前处理数据包（PREROUTING）
1. 处理流入的数据包（INPUT）
1. 处理流出的数据包（OUTPUT）
1. 处理转发的数据包（FORWARD）
1. 在进行路由选择后处理数据包（POSTROUTING）

使用最多的就是INPUT 规则链.

iptables 服务的动作分别是:
- ACCEPT（允许流量通过）
- REJECT（拒绝流量通过）
- LOG（记录日志信息）
- DROP（拒绝流量通过）

REJECT 和 DROP 的不同点: 就 DROP 来说，它是直接将流量丢弃而且不响应；REJECT 则会在拒绝流量后再回复一条提示信息，从而让流量发送方清晰地看到数据被拒绝的响应信息.

> 规则链的默认拒绝动作只能是 DROP，而不能是 REJECT

> 使用 iptables 命令配置的防火墙规则默认会在系统下一次重启时失效，如果想让配置的防火墙策略永久生效，还要执行保存命令`service iptables save`

iptables 命令可以根据流量的源地址、目的地址、传输协议、服务类型等信息进行匹配，一旦匹配成功，iptables 就会根据策略规则所预设的动作来处理这些流量.

## 参数
- -P : 设置默认策略
- -F : 清空规则链
- -L : 查看规则链
- -A : 在规则链的末尾加入新规则
- -I num : 在规则链的头部加入新规则
- -D num : 删除某一条规则
- -s : 匹配来源地址 IP/MASK，加叹号“!”表示除这个 IP 外
- -d : 匹配目标地址
- -i : 网卡名称 匹配从这块网卡流入的数据
- -o : 网卡名称 匹配从这块网卡流出的数据
- -p : 匹配协议，如 TCP、UDP、ICMP 
- --dport num : 匹配目标端口号
- --sport num : 匹配来源端口号

## example
```bash
# iptables -L # 查看已有的防火墙规则链
# iptables -F # 清空已有的防火墙规则链
# iptables -P INPUT DROP # 把 INPUT 规则链的默认策略设置为拒绝
# iptables -I INPUT -p icmp -j ACCEPT # 允许 ICMP 流量进入
# iptables -D INPUT 1 # 删除 INPUT 规则链中刚刚加入的一条策略
# iptables -P INPUT ACCEPT # 把默认策略设置为允许
# ### 只允许指定网段的主机访问本机的 22 端口，拒绝来自其他所有主机的流量
# iptables -I INPUT -s 192.168.10.0/24 -p tcp --dport 22 -j ACCEPT 
# iptables -A INPUT -p tcp --dport 22 -j REJECT
# ### 拒绝所有人访问本机 12345 端口的策略规则
# iptables -I INPUT -p tcp --dport 12345 -j REJECT
# iptables -I INPUT -p udp --dport 12345 -j REJECT
# ###
# iptables -I INPUT -p tcp -s 192.168.10.5 --dport 80 -j REJECT # 加拒绝 192.168.10.5 主机访问本机 80 端口（Web 服务）
# ### 拒绝所有主机访问本机 1000～1024 端口
# iptables -A INPUT -p tcp --dport 1000:1024 -j REJECT 
# iptables -A INPUT -p udp --dport 1000:1024 -j REJECT
# ### 
```