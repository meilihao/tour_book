# lspci
参考:
- [The PCI ID Repository](https://pci-ids.ucw.cz)

列出所有的pci设备即连接到 PCI 总线的所有设备, 比如主板, 声卡, 显卡, 网卡和usb接口设备等.

## 选项

- -d : 用于指定供应商、设备或类 ID 的所有设备
- -D : 在输出信息中显示设备的domain
- -k : 显示内核加载了哪些驱动程序模块.

    - `Kernel driver in use` : 正在使用的内核驱动程序
    - `Kernel modules` : 列出可用于支持该设备的模块
- -nn ： 显示供应商和设备ID
- -Q : 使用 DNS 查询中央数据库, 需要联网

    update-pciids 可更新本地 PCI ID 数据库.
- -s : 总线号, lspci输出的第一列, 比如`87:00.0`, 指定仅输出该设备信息
- -v : 详细的pci设备信息
- -vvv : 更详细的pci设备信息

## example
```bash
$ sudo lspci -nn | grep -e VGA
01:00.0 VGA compatible controller [0300]: NVIDIA Corporation GK107 [GeForce GTX 650] [10de:0fc6] (rev a1) # 设备名称后的方括号内有用冒号分隔的数字，即供应商和设备 ID. 输出表明 Nvidia Corporation 制造的设备的供应商 ID 为 10de
$ sudo lspci -nn -d 10de: # 所有 Nvidia 设备
01:00.0 VGA compatible controller [0300]: NVIDIA Corporation GK107 [GeForce GTX 650] [10de:0fc6] (rev a1)
01:00.1 Audio device [0403]: NVIDIA Corporation GK107 HDMI Audio Controller [10de:0e1b] (rev a1)
# lspci -s 03:10.3 -Dn
0000:03:10.3 0200: 8086:1520 (rev 01)
```

"0000:03:10.3"表示设备在PCI/PCI-E总线中的具体位置，依次是设备的domain(0000) 、 bus(03) 、 slot(10) 、 function(3), 其中domain的值
一般为0(当机器有多个host bridge时, 其取值范围是0~0xffff), bus的取值范围是0~0xff, slot取值范围是0~0x1f, function取值范围是0~0x7，其中后面3个值一般简称为BDF(即bus:device:function). 在输出信息中, 设备的vendor ID是"8086"("8086"ID代表Intel Corporation), device ID是"1520"(代表i350 VF) .

## FAQ
### 查找pci controller的驱动
```bash
# --- 直接通过lspci查找
$ lspci -nk/lspci -v
# --- 根据总线号查找
$ sudo lspci
...
02:00.0 Ethernet controller: Realtek Semiconductor Co., Ltd. RTL8111/8168B PCI Express Gigabit Ethernet controller (rev 01)
$ find /sys | grep drivers.*02:00
/sys/bus/pci/drivers/r8169/0000:02:00.0
$ lspci -vvv -s 87:00.0
```