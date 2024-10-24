# lsscsi
列出SCSI设备（或主机）及它们的属性。如果H:C:T:L给出参数，然后它作为一个过滤器，只匹配它被列出的设备。这里H是指Host，C是指Channel，T是指Id，L是指Lun.

**lsscsi会按照总线顺序枚举设备, 因此存在父子关系的设备肯定是连续出现的, 比如mediumx(磁带柜)和tape(磁带驱动器)**.

## 选项
- -g : 显示SCSI通用设备文件名称, 默认不显示`/dev/sg<N>`
- -k : 显示内核名称而不是设备节点名
- -d : 显示设备节点的主要号码和次要号码
- -H : 列出当前连接到系统的SCSI主机而不是SCSI设备
- -l : 显示每一个SCSI设备（主机）的附加信息
- -c : 相对于执行cat /proc/scsi/scsi命令的输出, 不推荐使用
- -p : 显示额外的数据完整性（保护）的信息
- -t : 显示传输信息
- -L : 以“属性名=值”的方式显示附加信息

        [`-t -L`联用可通过`transport`](http://sg.danny.cz/scsi/lsscsi.html#mozTocId961094)很方便的看出是什么设备:
        - fc
        - sata
        - pcie : nvme
        - iSCSI : iscsi
        - sas
        - spi
        - usb
- -v : 当信息找到时输出目录名
- -y<路径> :  假设sysfs挂载在指定路径而不是默认的`/sys`

## example
```bash
$ lsscsi # 列出SCSI设备及它们的属性, 机柜类型是enclosu
[1:0:0:0]    disk    ATA      Hoodisk SSD      61.3  /dev/sda 
$ lsscsi 2:0:0:0 # 显示匹配“2:0:0:0”的SCSI设备
$ lsscsi -d # 显示SCSI设备节点的主要号码和次要号码
$ lsscsi -L # 以“属性名=值”的方式输出SCSI设备的附加信息
$ lsscsi -v # 显示SCSI设备属性时也显示目录名
$ lsscsi -t # 显示SCSI设备的传输信息
$ lsscsi -c # 以相当于执行cat /proc/scsi/scsi命令的输出方式显示SCSI设备
$ lsscsi -t -L # `[1:0:0:0]`是SCSI设备id([H(SCSI adapter number, 比如hba):C(channel number即bus):T(target):L(LUN ID)])即[SCSI Addressing](https://www.tldp.org/HOWTO/SCSI-2.4-HOWTO/scsiaddr.html).
```
