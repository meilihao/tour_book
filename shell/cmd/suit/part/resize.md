# resize

## FAQ
### [扩容](https://help.aliyun.com/document_detail/113316.html)
```bash
yum install cloud-utils-growpart / apt install cloud-guest-utils

qemu-img resize openEuler-image.qcow2 +20G
growpart /dev/vda 1 # 表示扩容系统盘的第一个分区，/dev/vda是系统盘，1是分区编号，/dev/vda和1之间需要空格分隔. 如果单盘有多个连续分区的情况, 那么扩容时, 只需要扩容最后一个分区即可

resize2fs /dev/vda1 # ext4扩容
xfs_growfs /media/vdc # xfs扩容
btrfs filesystem resize max /mountpoint # btrfs扩容

df -Th # 检查扩容后结果
```

可能遇到的错误:
- `growpart /dev/vda 1`报`unexpected output in sfdisk --version [sfdisk，来自 util-linux 2.23.2]` : 先运行`LANG=en_US.UTF-8`切换字符编码类型，然后再进行尝试