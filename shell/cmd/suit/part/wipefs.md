# wipefs
wipefs可以擦除指定设备上的文件系统，raid或分区表签名（魔术字符串），从而使libblkid看不见签名.

wipefs不会擦除文件系统本身或设备中的任何其他数据; 
当不带任何选项使用时，wipefs会列出所有可见的文件系统及其基本签名的偏移量.

当擦除分区表签名以将更改通知内核时，wipefs调用BLKRRPART ioctl.

## example
```bash
# wipefs -af /dev/sdd # Erase any file system, partition table, or RAID signatures
dd if=/dev/zero of=/dev/sdd bs=512 count=1 conv=notrunc # 清除mbr
sfdisk --delete /dev/sdd # 清除分区表
```