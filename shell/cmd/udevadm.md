# udevadm
参考:
- [udevadm 中文手册](http://www.jinbuguo.com/systemd/udevadm.html)

udev 管理工具

## example
```bash
# udevadm monitor # 监视内核发出的设备事件(以"KERNEL"标记)， 以及udev在处理完udev规则之后发出的事件(以"UDEV"标记)，并在控制台上输出事件的设备路径(devpath)
# udevadm info --query=path --name=/dev/zd123 # 从udev数据库中提取设备信息
# udevadm info /sys/class/net/enp2s0 | grep ID_PATH= # 获取设备路径. PCI ID 是连接到系统的设备的唯一标识符.
# udevadm test /sys/class/net/eth0 # 测试应用到eth0上的rule是否能生效
```

`udevadm info`获取mediumx(磁带柜)其`/dev/sgX`设备比`/dev/schY`设备输出信息更多, 但tape和disk却恰好相反, 是`/dev/stZ`或`/dev/sdZ`输出更多信息.

## FAQ
### udev rules优先级
/etc/udev/rules.d中的任何规则都将优先于/lib/udev/rules.d中的规则

通过自定义规则可在`/dev/disk/by-xxx`自定义路径下创建相应设备的软链接

### `/etc/udev/rules.d/99-vmware-scsi-udev.rules:8 Invalid value "/bin/sh -c 'echo 180 >/sys$DEVPATH/timeout'" for RUN (char 27: invalid substitution type), ignoring, but please fix it.`
```bash
vim /etc/udev/rules.d/99-vmware-scsi-udev.rules
..., RUN+="/bin/sh -c 'echo 180 >/sys$$DEVPATH/device/timeout'"
```

> DEVPATH from `udevadm info /dev/sda`, `/sys$$DEVPATH`即device path
