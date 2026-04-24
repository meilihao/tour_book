# faq
### tcp bbr
要求 : linux kernel >=4.9

```
# /etc/sysctl.conf添加
net.ipv4.tcp_congestion_control=bbr
net.core.default_qdisc=fq
```

```
# 使新配置生效
sudo sysctl -p

# 结果都有 bbr , 则证明你的内核已开启bbr
sudo sysctl net.ipv4.tcp_available_congestion_control
sudo sysctl net.ipv4.tcp_congestion_control

# 看到有 tcp_bbr 模块即说明bbr已启动
lsmod | grep bbr
```