# scsi
[sg3-utils](http://sg.danny.cz/sg/)是一个工具包，提供了与SCSI设备通信的命令工具.

sg = SCSI Generic
ses = SCSI Enclosure Service

ses标准:
- [SCSI Enclosure Services - 3 (SES-3) ](https://www.t10.org/members/w_ses3.htm)
- [SCSI Enclosure Services - 4 (SES-4) ](https://www.t10.org/members/w_ses4.htm)

# sg_map
Linux实现了一个通用的SCSI设备驱动，如果一个设备支持SCSI协议，那么当它插入后，SCSI设备驱动将自动识别它，并创建一个相关联的设备文件，通常为/dev/sg0、/dev/sg1等（一切设备皆文件).

确定哪个 SCSI 设备表示某个 sg 接口,即sg设备到scsi设备的映射. SCSI 磁盘的标签为 /dev/sda、/dev/sdb 和 /dev/sdc 等.

通过这些通用的驱动器接口(/dev/sgXXX)，就可以将 SCSI 命令直接发送到 SCSI 设备，而不需要经过在 SCSI 磁盘上创建（并装载到某个目录）的文件系统.

## 选项

- -i : 输出标准的查询+vendor(供应商)+产品名及修订字符串
- -x : 显示sg设备的pcie地址: <host_number> <bus> <scsi_id> <lun> <scsi_type>, 前4个数字是scsi设备映射到的SCSI地址, 与`/sys/class/enclosure`下的目录对应


### example
```bash
# sg_map -i
# sg_map | awk '{if($2==""){print $1}}' # 显示扩展柜和主控, 可直接用`sg_map -x`的scsi_type过滤出扩展柜(scsi_type=13=0xd)
# ls /dev/sg* -l |awk '{if($4=="root"){print $10}}' # 显示扩展柜和主控
# lsscsi -g|grep "enclosu" # 显示扩展柜
```

# sg_ses - send controls and fetch status from a SCSI EnclosureServices (SES) device
## 选项
- --index   : Element index (not slot index)
- --set     : set   a status of specify element
- --clear   : clear a status of specify element
- --dev-slot-num=SN : 槽位号
- -r : 输出原始信息, 即二进制信息

### example
```bash
# sg_ses -p 0x0 /dev/sg7 # 查看sg设备支持的pages
# sg_ses -p 0x2 /dev/sg7 # 根据`sg_ses -p 0x0`返回的结果, 查看指定的page, 这里的0x2表示`Enclosure status/control (SES) [0x2]`
# sg_ses -p 2 -I 27 /dev/sg7 # 查看指定element的enclosure status
# sg_ses -p 0xa /dev/sg7 # 获取扩展柜中设备的SAS address, 槽位号
# sg_ses -p 0xa /dev/sg7 |grep -E 'slot|Element' |sed 'N;s/\n//' |awk '{print $3,$15}' # 获取element_index与slot_number的对应关系
# sg_ses -p 0xa /dev/sg7 |grep -E 'slot|Element' |sed 'N;s/\n//' |awk '{print $15,$3}' |sort -n # 获取slot_number与element_index的对应关系
# sg_ses -p 0xa /dev/sg7 |grep -E 'slot number|  SAS address' |sed 'N;s/\n//' |awk '{print $12,$15}' |sort -n # 槽位对应的SAS address, 0x0000000000000000或0x0表示没有盘
# sg_ses -ee # 查看允许设置的状态
# --- disk fault: (Red LED light on)
# --- `--set/--clear/--get`对应的格式是`<start_byte>:<start_bit>[:<number_of_bits>]`, <number_of_bits>未提供时默认是1 
# sg_ses --index=27 --set=fault  /dev/sg5
# sg_ses --index=27 --set=3:5:1  /dev/sg5 # 同上
# sg_ses --index 27 --clear=3:5 /dev/sg5 清除状态
# ---
```

## sginfo
```bash
# sginfo - # 获取sg设备的信息
```

## sg_inq
```
# sg_inq /dev/sda # 获取磁盘的概要信息
# sg_inq -p 0x0 /dev/sda # 获取磁盘支持的pages
# sg_inq -p 0x83 /dev/sda # 获取磁盘的设备标识
```

## 扩展
### Control LED
Insert harddisk, Blue LED light on, it is controlled by hardware.
Read/write harddisk, Blue LED flash, it is controlled by hardware also.

```bash
# sg_ses -ee
```
