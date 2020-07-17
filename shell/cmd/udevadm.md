# udevadm
参考:
- [udevadm 中文手册](http://www.jinbuguo.com/systemd/udevadm.html)

udev 管理工具

## example
```bash
# udevadm monitor # 监视内核发出的设备事件(以"KERNEL"标记)， 以及udev在处理完udev规则之后发出的事件(以"UDEV"标记)，并在控制台上输出事件的设备路径(devpath)
# udevadm info --query=path --name=/dev/zd123 # 从udev数据库中提取设备信息
```