# dd

## 描述

文件复制和转换工具

## 选项
- if : 输入来源
    - `/dev/zero` : 提供无限个空字符（NULL/0x00）
    - `/dev/urandom` : 提供不为空字符的随机字节数据流
- of : 输出到指定位置
- bs : 单次读取和写入的字节数, 默认是512
- count : 操作次数, 因此文件大小 = bs * count
- skip=n : 用于`if`, 表示从开头跳过n个块开始读取.
- seek=BLOCKS : 用于`of`, 跳过输出文件中指定大小的块数来写，但是并不实际写入, 因此速度很快. 同时因为不实际写入磁盘，所以在容量只有10G的硬盘上创建100G的此类文件也是可以的
- status=progress : 显示进度
- iflag/oflag : dd做读写测试时，要加两个参数 iflag=nocache 和 oflag=direct 参数. 没有的话dd有时会显示从内存中传输数据的结果，速度会不准确.
- conv
    - notrunc: 不截断输出文件

## example
```sh
$ dd if=/etc/fstab of=/etc/fstab.bak
3+1 records in # bs默认是512, 3个完整的512块+1个未满512的块
3+1 records out
$ dd if=/dev/sda2 of=/opt/sda2.bak # 备份整个sda2分区
$ dd if=/opt/sda2.bak of=/dev/sda2 # 恢复分区, 恢复前需卸载分区
$ dd if=/dev/zero of=test bs=1M count=1000 # 在当前目录下会生成一个1000M的test文件，文件内容为全0, 但是这样为实际写入硬盘，文件产生速度取决于硬盘读写速度
$ seq 1000000 | xargs -i dd if=/dev/zero of={}.dat bs=1024 count=1 # 随机生成1百万个1K的文件
$ dd if=/dev/cdrom of=centos.iso # 将光驱设备中的光盘制作成 iso 格式的镜像文件
$ dd bs=8k count=4k if=/dev/zero of=test.log conv=fdatasync/fsync # fdatasync/fsync区别是conv=fsync会把文件的“数据”和“metadata”都写入磁盘, 而fdatasync仅数据落盘, 两者时间相差不大. 单纯磁盘性能测试推荐用fdatasync. dd默认启用写缓存(先把数据写到os的“写缓存”，就算完成了写操作, 再由os周期性地调用sync函数，把“写缓存”中的数据刷入磁盘. 因此“写缓存”的存在，会测试出一个超级快的错误性能值. from [正确使用 dd 测试磁盘读写速度](https://cloud.tencent.com/developer/article/1114720)
$ dd if=/dev/urandom of=xxx bs=512 count=1 skip=2 conv=notrunc # 在1024B后替换512的内容
```

读取mbr:
```bash
# dd if=/dev/sda of=mbr.hex bs=512 count=1
# hexdump -C mbr.hex
```

## FAQ
### dd大文件导致内核cache占用很大
dd 60G的zvol到zfs fs, 发现/proc/meminfo的Buffers变得很高(16G, 内存共63G, kernel cache占28G), kill dd进程后Buffers立马变为200M

解决(未验证):
1. dd bs=4M iflag=nocache oflag=sync
1. dd iflag=direct oflag=direct
1. 调整虚拟内存管理参数

    - vm.dirty_ratio：表示系统内存中可以占用多少比例的数据在同步到磁盘之前保持在内存中。如果内存占用过高，dd 可能会导致写入数据堆积在内存中，增加缓存压力
    - vm.dirty_background_ratio：表示当内存中有多少比例的脏数据时，后台的写入进程会开始把脏数据写回磁盘

    比如:
    ```bash
    sysctl vm.dirty_ratio=10
    sysctl vm.dirty_background_ratio=5
    ```