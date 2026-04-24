# bluetooth
## tools
```bash
$ dnf install bluez
$ bluetoothctl # 交互式命令
# scan on # 开始扫描
```

## FAQ
### linux蓝牙无法连接(开关已打开)
```
$ sudo journalctl -f 
3月 21 11:44:56 chen-pc bluetoothd[1881]: Failed to set mode: Blocked through rfkill (0x12)
...
$ rfkill list                                                                                                                                         11:47:23
0: phy0: Wireless LAN
	Soft blocked: no
	Hard blocked: no
1: hci0: Bluetooth
	Soft blocked: yes # 命令deepin控制台的蓝牙开关已打开, 却还是显示blocked.
	Hard blocked: no
$ rfkill unblock 1
```

bluetoothctl命令可查看蓝牙状态, 比如`scan on`监听蓝牙设备的变化.
blueman是管理蓝牙的gui工具.