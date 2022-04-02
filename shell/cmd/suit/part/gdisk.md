# gdisk
参考:
- [gPT fdisk](https://wiki.archlinux.org/index.php/GPT_fdisk_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87))
- [Partitioning (简体中文)](https://wiki.archlinux.org/index.php/Partitioning_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87))
- [C语言读取GPT分区信息](https://blog.csdn.net/qq_37734256/article/details/88384750)
- [使用 sgdisk 管理分区](https://linux.cn/article-10771-1.html)

gdisk是支持gpt的分区工具.

> from gdisk_1.0.5-1_amd64

> cgdisk: 支持终端窗口功能的gdisk

## example
```bash
$ sudo  gdisk /dev/nbd0
Command (? for help): ? # help
Command (? for help): p # 打印分区表
Command (? for help): n # 新建分区, 此时默认已经创建了分区表
Partition number (1-128, default 1): 1 # 分区号
First sector (34-16777182, default = 2048) or {+-}size{KMGTP}: # 起始扇区
Last sector (2048-16777182, default = 16777182) or {+-}size{KMGTP}: +256M # 新分区的大小
Current type is 'Linux filesystem'
Hex code or GUID (L to show codes, Enter = 8300): # 默认即可之后再修改, 8300是用于格式化成ext4/xfs等linux文件系统用的
Command (? for help): p # 查看创建的分区, 及剩余的磁盘空间
Command (? for help): d # 删除分区
Partition number (1-2): 2 # 删除分区2
Command (? for help):  l # 查看支持的分区类型
Command (? for help): t 　　　  　 # 更改分区类型，这里输入l也可以查看分区的类型
Partition number (1-2): 2 　　　　 # 输入要更改的分区号
Current type is 'Linux filesystem'
Hex code or GUID (L to show codes, Enter = 8300): 8e00　　　　 # 输入分区类型的编号 
Command (? for help): c                              # 更改分区名称
Partition number (1-2): 2                            # 要更改的分区号
Enter name: pv1 LVM                                  # 更改后的名称
Command (? for help): w 　　　　　　　　 # 保存配置，如果不想保存可以输入q退出

Final checks complete. About to write GPT data. THIS WILL OVERWRITE EXISTING
PARTITIONS!!

Do you want to proceed? (Y/N): y # 询问是否相想继续，输入y继续
OK; writing new GUID partition table (GPT) to /dev/nbd0.
The operation has completed successfully.
```

## gdisk分区代号
- EF02  BIOS boot partition

    gnu grub用它来引导基于legacy bios但启动设备上却包含GPT格式分区表时的操作系统. 这种结构有时候被称为BIOS/GPT启动.

    此时使用`sudo parted -l /dev/nvme0n1`查看时, Partition Table: gpt, Disk Flags不为空, 是pmbr_boot. 

- EF00  EFI System
- 8300  Linux filesystem

## script分区
参考:
- [Creating GPT partitions easily on the command line](https://suntong.github.io/blogs/2015/12/26/creating-gpt-partitions-easily-on-the-command-line/)

> sgdisk支持linux+unix, gdisk支持linux.

```bash
# sed -e 's/\s*\([\+0-9a-zA-Z]*\).*/\1/' << EOF | gdisk /dev/nbd0
  o # new gpt
  Y # Proceed
  n # new partition
  1 # Partition number
    # default - start at beginning of disk
  +1M # 1 MB BIOS boot partition
  ef02  # new partition
  n
  2

  +256M
  ef00 # EFI System
  n
  3

  +2048M
  8300
  n
  4
  
  
  
  w # write GPT data
  Y # want to proceed
EOF
# gdisk /dev/nbd0 # 上面命令的分区效果
GPT fdisk (gdisk) version 1.0.3

Partition table scan:
  MBR: protective
  BSD: not present
  APM: not present
  GPT: present

Found valid GPT with protective MBR; using GPT.

Command (? for help): p
Disk /dev/nbd0: 16777216 sectors, 8.0 GiB
Sector size (logical/physical): 512/512 bytes
Disk identifier (GUID): D43F5226-18A1-4A4B-B6B4-F77E5BAF6961
Partition table holds up to 128 entries
Main partition table begins at sector 2 and ends at sector 33
First usable sector is 34, last usable sector is 16777182
Partitions will be aligned on 2048-sector boundaries
Total free space is 2014 sectors (1007.0 KiB)

Number  Start (sector)    End (sector)  Size       Code  Name
   1            2048            4095   1024.0 KiB  EF02  BIOS boot partition
   2            4096          528383   256.0 MiB   EF00  EFI System
   3          528384         4722687   2.0 GiB     8300  Linux filesystem
   4         4722688        16777182   5.7 GiB     8300  Linux filesystem
# --- format /dev/sdb as GPT, GUID Partition Table
# sgdisk -Z /dev/nbd0
# sgdisk -n 0:0:+1M -t 0:ef02 -c 0:"bios_boot" /dev/nbd0
# sgdisk -n 0:0:+256M -t 0:ef00 -c 0:"efi" /dev/nbd0
# sgdisk -n 0:0:+2G -t 0:8300 -c 0:"linux_boot" /dev/nbd0
# sgdisk -n 0:0:+1G -t 0:0700 -c 0:"windows" /dev/nbd0
# sgdisk -n 0:0:+1G -t 0:8200 -c 0:"linux_swap" /dev/nbd0
# sgdisk -n 0:0:+1G -t 0:8300 -c 0:"os1" /dev/nbd0
# sgdisk -n 0:0:0 -t 0:8300 -c 0:"data" /dev/nbd0
# sgdisk -p /dev/nbd0
# sgdisk -p /dev/nbd0
Disk /dev/nbd0: 16777216 sectors, 8.0 GiB
Sector size (logical/physical): 512/512 bytes
Disk identifier (GUID): 46C2F0D3-8206-45A8-98D3-BDB09981A754
Partition table holds up to 128 entries
Main partition table begins at sector 2 and ends at sector 33
First usable sector is 34, last usable sector is 16777182
Partitions will be aligned on 2048-sector boundaries
Total free space is 2014 sectors (1007.0 KiB)

Number  Start (sector)    End (sector)  Size       Code  Name
   1            2048            4095   1024.0 KiB  EF02  bios_boot
   2            4096          528383   256.0 MiB   EF00  efi
   3          528384         4722687   2.0 GiB     8300  linux_boot
   4         4722688         6819839   1024.0 MiB  0700  windows
   5         6819840         8916991   1024.0 MiB  8200  linux_swap
   6         8916992        11014143   1024.0 MiB  8300  os1
   7        11014144        16777182   2.7 GiB     8300  data
```

sgdisk命令解析：
1. `sgdisk -n 1:2048:4095 -c 1:"BIOS Boot Partition" -t 1:ef02 /dev/nbd0`
1. `sgdisk -n 0:0:+1M -t 0:ef02 -c 0:"bios_boot" /dev/nbd0`

- `1:2048:4095`: 1, 分区编号(从1开始); 2048, 起始扇区; 4095, 结束扇区
- `0:0:+1M`: 0, 分区编号sgdisk自动计算; 0, 起始扇区sgdisk自动计算; +1M, 该分区大小
- `-t 0:ef02`: 0, 前面自动生成的分区号, ef02是分区类型
- `-c 0:"bios_boot"`: 0, 前面自动生成的分区号, bios_boot是该分区的name

### 获取pt uuid/part uuid/fs uuid
`sudo blkid /dev/nbd0`

fs uuid: `ls -l /dev/disk/by-uuid/`
pt uuid: `gdik -l /dev/disk`