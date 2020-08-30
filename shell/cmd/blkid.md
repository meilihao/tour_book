# blkid
获取磁盘上所采用文件系统信息, 比如文件系统类型,Label,UUID等.

## example
```bash
# blkid -p /dev/nvme0n1
/dev/nvme0n1: PTUUID="f8fbd0cd-27ad-4df0-bcc1-7c5e88bf6cad" PTTYPE="gpt"
# blkid
/dev/nvme0n1p1: UUID="643D-9978" TYPE="vfat" PARTUUID="8c465477-4444-4e2a-9306-6526f24cae36"
/dev/nvme0n1p2: UUID="5786c7db-036e-42ff-a0e0-a676133b3dcf" TYPE="ext4" PARTUUID="bd1fa4e9-95d0-4663-bd1a-f00fd3c35528"
/dev/nvme0n1p3: UUID="0cbd324e-7ef1-4152-8e62-a3997e8ad987" TYPE="ext4" PARTUUID="68b77c91-fe6d-4402-b7be-d160f2d80137"
/dev/nvme0n1p4: UUID="fcfd3de6-ae22-4e74-9f63-24ab6e81210e" TYPE="ext4" PARTUUID="2a2bb175-909a-4712-9a97-43e25e960785"
# gdisk -l /dev/nvme0n1
GPT fdisk (gdisk) version 1.0.3

Partition table scan:
  MBR: protective
  BSD: not present
  APM: not present
  GPT: present

Found valid GPT with protective MBR; using GPT.
Disk /dev/nvme0n1: 500118192 sectors, 238.5 GiB
Model: SAMSUNG MZVLV256HCHP-00000              
Sector size (logical/physical): 512/512 bytes
Disk identifier (GUID): F8FBD0CD-27AD-4DF0-BCC1-7C5E88BF6CAD
Partition table holds up to 128 entries
Main partition table begins at sector 2 and ends at sector 33
First usable sector is 34, last usable sector is 500118158
Partitions will be aligned on 2048-sector boundaries
Total free space is 2669 sectors (1.3 MiB)

Number  Start (sector)    End (sector)  Size       Code  Name
   1            2048         1050623   512.0 MiB   EF00  
   2         1050624         5244927   2.0 GiB     8300  
   3         5244928        88604671   39.7 GiB    8300  
   4        88604672       500117503   196.2 GiB   8300
```

## FAQ
### blk输出中的PTUUID, PARTUUID, UUID区别
- PTUUID : 磁盘分区表本身的标识, 通过gdisk也可看到, 等价于MBR分区磁盘上的磁盘签名.
- PARTUUID : partition-table-level UUID, 是GPT分区磁盘中的分区才有的功能. 因为不依赖分区的实际内容, 因此是使用某种未知的加密方法对分区进行加密时，该分区的唯一标识符.
- UUID : 是filesystem-level UUID, 从分区中的fs metadata中获取到, 仅在fs类型已知且可读的情况下才能获取到.

> 在MBR分区的磁盘上，分区表中没有实际的UUID. 因此，使用32位磁盘签名代替PTUUID，并通过在磁盘签名的末尾添加`-`和两位分区号来创建PARTUUID.