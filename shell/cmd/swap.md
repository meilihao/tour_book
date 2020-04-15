#  mkswap
将创建的分区格式为SWAP 分区.

## other
` swapon ${swap分区}`: 启用swap分区
` swapoff ${swap分区}`: 关闭swap分区

## FAQ
### swap在`/etc/fstab`的格式
```conf
/dev/sdb2 swap swap defaults 0 0
```