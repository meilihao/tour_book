# lspci
参考:
- [The PCI ID Repository](https://pci-ids.ucw.cz)

列出所有的pci设备即连接到 PCI 总线的所有设备, 比如主板, 声卡, 显卡, 网卡和usb接口设备等.

> lspci is from [pciutils](http://mj.ucw.cz/sw/pciutils/), 对应命令还有setpci(用于读写PCIe配置寄存器).

## 选项

- -d : 用于指定供应商、设备或类 ID 的所有设备, 比如`-d 14e4:165a`
- -D : 在输出信息中显示设备的domain
- -k : 显示内核加载了哪些驱动程序模块.

    - `Kernel driver in use` : 正在使用的内核驱动程序
    - `Kernel modules` : 列出可用于支持该设备的模块
- -n : `-nn`的精简输出

    `03:00.0 0200: 14e4:165a`解读:
    - 03:00.0 : bus number = 03; device number = 00; function = 0, 这3个编号会组合成一个16bits的pcie设备的识别码. 但在 Linux 是使用 Class ID + Vendor ID + Device ID 来代表设备, id与名称的对应关系见`/usr/share/hwdata/pci.ids`.

        - bus number, 8bits 2^8, 至多可可接 256 个bus(0 to ff)
        - device number, 5bits 2^5, 至多可接 32 个device(0 to 1f) 
        - function number, 3bits 2^3, 至多每种设备可有 8 种功能(0 to 7)
    - 0200 : class 0200 表示是 Network controller. class id(即`pci_dev->class`) from[PCIE Configuration Space – Class Code即pci code and id assignment specification](https://blog.ladsai.com/pci-configuration-space-class-code.html)
    - 14e4 : vendor ID 14e4 产商 Broadcom Corporation
    - 165a : device ID 1659 产品名称 NetXtreme BCM5721 Gigabit Ethernet PCI Express
- -nn ： 显示供应商和设备ID
- -m : 以机器可读格式输出信息, 还可用`-mm`, `-vmm`
- -Q : 使用 DNS 查询中央数据库, 需要联网

    update-pciids 可更新本地 PCI ID 数据库.
- -s : 总线号, lspci输出的第一列, 比如`87:00.0`, 指定仅输出该设备信息
- -v : 详细的pci设备信息
- -vvv : 更详细的pci设备信息
- -t : pcie tree, 常用`lspci -tvv`
    
    - `[0000:00]` : 树的根即是RC(Root Complex)
    - `+-1c.1-[02]` : 1c.1 是 PCIe 根集线器的插槽编号`1c`和功能编号`1`. 在这种情况下, 它包含一个 PCIe 桥接器, 这座桥后面的总线编号为是02,这里也可以是范围比如`02-3a`
- -xxx : 提供PCI 配置空间的十六进制转储

## example
```bash
$ sudo lspci -nn | grep -e VGA
01:00.0 VGA compatible controller [0300]: NVIDIA Corporation GK107 [GeForce GTX 650] [10de:0fc6] (rev a1) # 设备名称后的方括号内有用冒号分隔的数字，即供应商和设备 ID. 输出表明 Nvidia Corporation 制造的设备的供应商 ID 为 10de
$ sudo lspci -nn -d 10de: # 所有 Nvidia 设备
01:00.0 VGA compatible controller [0300]: NVIDIA Corporation GK107 [GeForce GTX 650] [10de:0fc6] (rev a1)
01:00.1 Audio device [0403]: NVIDIA Corporation GK107 HDMI Audio Controller [10de:0e1b] (rev a1)
# lspci -s 03:10.3 -Dn
0000:03:10.3 0200: 8086:1520 (rev 01)
# lspci -tvv # 获取pcie tree
# lspci -nvm # 用于以数字形式显示 PCI 器件供应商 ID 和器件 ID
```

"0000:03:10.3"表示设备在PCI/PCI-E总线中的具体位置，依次是设备的domain(0000) 、 bus(03) 、 slot(10) 、 function(3), 其中domain的值
一般为0(当机器有多个host bridge时, 其取值范围是0~0xffff), bus的取值范围是0~0xff, slot取值范围是0~0x1f, function取值范围是0~0x7，其中后面3个值一般简称为BDF(即bus:device:function). 在输出信息中, 设备的vendor ID是"8086"("8086"ID代表Intel Corporation), device ID是"1520"(代表i350 VF) .

## pci输出说明
1. `PME(D0+,D1-,D2-,D3hot+,D3cold+) power management event的状态, 在输出中`+`代表有被启用, `-`代表沒有被启用`.
1. LnkCap/LnkSta

    LnkCap : 系统能提供的最高频宽 PCI-Express 1.0 ( 2.5G ) Width x4=10G. 如果系统是提供 PCI-Express 2.0 那 x1 速度是 5G.
    LnkSta : 目前该PCI-E 装置跑的速度 PCI-Express 1.0 ( 2.5G ) Width x1=2.5G

    LnkSta 和 LnkCap 这两个速度有可能不一样, 系统所提供的是 PCI Express 是 2.0, 但装置还是使用 1.0 的

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

### lspci输出的pcie tree乱了
做底层驱动开发特别是PCIe驱动相关的开发, 经常看到的一个错误就是“树乱了”. 树乱了从表面上看, 即是`lspci -t`的内容发生了错乱, 跟正常情况下的树不一样了。本质上, “树乱了”由于系统中的某些PCIe设备出错或者异常导致的. `lspci -t`时, 系统会重新访问这颗树上的所有设备, 由于部分PCIe设备的异常, 访问失败(返回异常值或者全F), 从而影响pciutils构建出一颗完整的树, 出现部分树枝断裂(链路link down), 树叶错位(EP Bar空间异常)等现象.