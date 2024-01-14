# vlan

## example
```bash
配置VLAN
# 安装并加载内核模块
apt-get install vlan
modprobe 8021q

# 添加vlan
vconfig add eth0 100
ifconfig eth0.100 192.168.100.2 netmask 255.255.255.0

# 删除vlan
vconfig rem eth0.100
```