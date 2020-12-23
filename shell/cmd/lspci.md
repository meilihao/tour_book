# lspci
列出所有的pci设备, 比如主板, 声卡, 显卡, 网卡和usb接口设备等.

## 选项

- -v : 更详细的pci设备信息

## FAQ
### 查找pci controller的驱动
```bash
$ sudo lspci
...
02:00.0 Ethernet controller: Realtek Semiconductor Co., Ltd. RTL8111/8168B PCI Express Gigabit Ethernet controller (rev 01)
$ find /sys | grep drivers.*02:00
/sys/bus/pci/drivers/r8169/0000:02:00.0
# # --- 直接通过lspci查找
$ lspci -nk/lspci -v
```