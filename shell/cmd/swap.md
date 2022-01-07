#  mkswap
将创建的分区格式为SWAP 分区.

## other
`cat /proc/swaps` : 查看swap状态
`swapon --show` : 查看swap状态
`swapon ${swap分区}`: 启用swap分区
`swapoff ${swap分区}`: 关闭swap分区

## FAQ
### swap在`/etc/fstab`的格式
```conf
/dev/sdb2 swap swap defaults 0 0
```

### 禁用swap
1. 注释`/etc/fstab`中swap项. 执行该操作并重启系统后发现 swap 分区依然被开启了
1. systemd禁用swap

```bash
# --- 需swap on的状态下执行
$ swapon --show
$ systemctl list-units | grep swap
dev-sda2.swap                                                                             loaded active active    Swap Partition
swap.target                                                                               loaded active active
$ systemctl cat dev-sda2.swap
$ sudo systemctl mask dev-sdXX.swap
```

### swap大小
ref:
- [RHEL推荐的SWAP空间的大小划分原则](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html/managing_storage_devices/getting-started-with-swap_managing-storage-devices)