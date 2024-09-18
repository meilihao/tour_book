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

## `ethtool -S`输出
ref:
- [`ethtool -S`](https://support.huawei.com/enterprise/zh/doc/EDOC1100280170/4697098e)
- [struct rtnl_link_stats64](https://docs.kernel.org/networking/statistics.html)

字段:
- rx_missed：表明网络接口丢失的数据包数量，通常与网络流量和接口的处理能力有关。
- rx_no_buffer_count：表明系统在接收数据包时没有足够的缓冲区，通常与系统内存或网络缓冲区配置有关
- rx_missed_errors: [在DMA传送完, 发送硬中断之前, 网卡的FIFO缓冲已经满了, 导致接收的数据要立即丢掉](https://lp007819.wordpress.com/2013/05/23/intel%E7%BD%91%E5%8D%A1%E7%8A%B6%E6%80%81%E7%BB%9F%E8%AE%A1%E7%96%91%E9%97%AE/)

    rx_missed_errors =E1000_MPC
    rx_fifo_errors=E1000_MPC+E1000_RQDPC
    rx_no_buffer_count=E1000_RNBC

    drop = rx_dropped(内核丢弃) + rx_missed_errors
    overrun = rx_fifo_error

## FAQ
### `/var/logs/messages 大量出现"tx hang... resetting adapter", "Detected Tx Unit Hang"`
ref:
- [Intel 82599网卡异常挂死原因](https://www.cnblogs.com/smith9527/p/10348953.html)

    通过检查PCI信息和网卡寄存器，发现寄存器值全为0xffffffff，可能原因包括PCIe接口接触不良导致的通信异常
- [ixgbe driver hang up | Detected Tx Unit Hang Tx Queue](https://forum.proxmox.com/threads/ixgbe-driver-hang-up-detected-tx-unit-hang-tx-queue.120328/)

    换成X710后正常
- [e1000e网卡驱动频繁报告“Detected Hardware Unit Hang”错误](https://www.aliencn.net/archives/412)

    可以看到rx-checksumming和tx-checksumming是on的，就是因为这个功能和当前系统不兼容导致的

这是网络接口ens192的传输故障导致的问题. "tx hang"表示传输队列挂起, "resetting"表示系统正在尝试重置网络接口以解决问题. 这可能是由于网络接口驱动程序或硬件故障引起的.
建议检查网络接口的驱动程序是否最新，并确保硬件连接正常. 如果问题持续存在，可能需要更换网络接口或与硬件供应商联系以获取支持.

```bash
# lspci -vvv -s 84:00.0
# ethtool -d eth0
```