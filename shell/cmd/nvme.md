# nvme

# nvme-cli
ref:
- [Solid state drive/NVMe](https://wiki.archlinux.org/title/Solid_state_drive/NVMe)

nvme存储的命令行工具.

```bash
# apt install nvme-cli
# dnf install nvme-cli
```

## example
```bash
# nvme list # 列出你机器上所有的 NVMe 设备和命名空间
# nvme id-ctrl /dev/nvme0n1 [-H] # 获取更多关于该硬盘和它所支持的特性的信息
# nvme id-ns /dev/nvme0n1 -H # 查看nvme ns信息, 比如lba大小. `smartctl -c /dev/nvme0n1`也可查看lba大小
# nvme smart-log /dev/nvme0n1 # 了解硬盘的整体健康状况
# nvme error-log /dev/nvme0 # 查看nvme error log page
# nvme format [--lbaf=1] /dev/nvme0nX # 格式化一个 NVMe 驱动器. lbaf, 指定lba size, 通常使用4k可获得更高性能.
# nvme sanitize /dev/nvme0nX # 安全地擦除驱动器数据
# nvme fw-log /dev/nvme0 # 查看固件信息

# nvme discover -t rdma -a 192.168.80.100 -s 4420
# nvme connect -t rdma -n "nqn.2016-06.io.spdk:cnode1" -a 192.168.80.100 -s 4420
# nvme disconnect -n "nqn.2016-06.io.spdk:cnode1"
```

# nvmetcli
配置NVMe target工具

> 只在centos/fedora有该工具.