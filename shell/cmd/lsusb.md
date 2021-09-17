# lsusb
列出所有的use设备即连接到 usb 总线的所有设备

## 选项

- -D : 显示设备详情， 包括供应商和设备ID

    update-usbids 可更新本地 PCI ID 数据库.

## example
```bash
lsusb -D /dev/bus/usb/001/002
```