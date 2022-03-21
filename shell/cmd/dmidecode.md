# dmidecode
dmidecode遵循SMBIOS/DMI标准(Desktop Management Interface,DMI)，输出硬件信息, 比如BIOS、系统、主板、处理器、内存、缓存等等.

> 由于DMI信息可以人为修改，因此里面的信息不一定是系统准确的信息.

SMBIOS规范定义了以下DMI类型:
Type   Information
────────────────────────────────────────────
   0   BIOS
   1   System
   2   Baseboard
   3   Chassis
   4   Processor
   5   Memory Controller
   6   Memory Module
   7   Cache
   8   Port Connector
   9   System Slots
  10   On Board Devices
  11   OEM Strings
  12   System Configuration Options
  13   BIOS Language
  14   Group Associations
  15   System Event Log
  16   Physical Memory Array
  17   Memory Device
  18   32-bit Memory Error
  19   Memory Array Mapped Address
  20   Memory Device Mapped Address
  21   Built-in Pointing Device
  22   Portable Battery
  23   System Reset
  24   Hardware Security
  25   System Power Controls
  26   Voltage Probe
  27   Cooling Device
  28   Temperature Probe
  29   Electrical Current Probe
  30   Out-of-band Remote Access
  31   Boot Integrity Services
  32   System Boot
  33   64-bit Memory Error
  34   Management Device
  35   Management Device Component
  36   Management Device Threshold Data
  37   Memory Channel
  38   IPMI Device
  39   Power Supply
  40   Additional Information
  41   Onboard Devices Extended Information
  42   Management Controller Host Interface

输出时可以使用关键字来代替数字形式的类型, 需要添加--type参数:
Keyword     Types
──────────────────────────────
bios            0, 13
system       1, 12, 15, 23, 32
baseboard  2, 10, 41
chassis       3
processor   4
memory      5, 6, 16, 17
cache         7
connector   8
slot             9

## 输出字段说明
- Handle: 标识符号


## example
```bash
dmidecode # 获取全部信息
dmidecode -t 1  # 查看服务器信息
dmidecode |grep 'Product Name' # 查看服务器型号 
dmidecode |grep 'Serial Number' # 查看主板的序列号 
dmidecode -t 2  # 查看主板信息
dmidecode -s system-serial-number # 查看系统序列号 
dmidecode -t memory # 查看内存信息 
dmidecode -t 11 # 查看OEM信息 
dmidecode -t 17 # 查看内存条数
dmidecode -t 16 # 查询内存信息
dmidecode -t 4  # 查看CPU信息
dmidecode -t 0  # 查看bios

cat /proc/scsi/scsi # 查看服务器硬盘信息
```