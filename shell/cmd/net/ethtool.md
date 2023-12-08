# ethtool

配置硬件选项

## example
```bash
# 改变speed
ethtool -s eth0 speed 1000 duplex full

# 关闭GRO
ethtool -K eth0 gro off

# 开启网卡多队列
ethtool -L eth0 combined 4

# 开启vxlan offload
ethtool -K ens2f0 rx-checksum on
ethtool -K ens2f0 tx-udp_tnl-segmentation on

# 查询网卡统计
ethtool -S eth0
```