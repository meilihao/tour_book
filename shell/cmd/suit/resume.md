# resume
ref:
- [详解在 Ubuntu 中引导到救援模式或紧急模式](https://linux.cn/article-14709-1.html)

进入resume方法:
1. 在GRUB 菜单中，选择第一项，并按下 e 按键来编辑它, 在其启动参数`linux /boot/vmlinuz-5.15.0-30-generic ... ro quiet splash`(如果存在`$vt_handoff`就删除它)后追加` systemd.unit=rescue.target`, 再按`Ctrl + x`或按`F10`来引导到救援模式

	也可选择`recovery mode`引导项, 通常发行版均提供了该引导项

	Ubuntu 22.04还提供了`Recovery menu`, 选择`root`项即可, 修复完成后执行`exit`还是会回到该menu, 此时选择`resume`开始正常引导即可

	> `rescue.target`也可替换为`systemd.unit=emergency.target`, 或使用`systemctl emergency`/`systemctl rescue`来切换
1. 执行`mount -n -o remount,rw /`后可读写根文件系统
1. 执行`systemctl default`或`exit`来继续引导到正常模式, 或使用`systemctl reboot`重启系统

## FAQ
### 单用户/emergency/rescue区别
单用户模式是一个runlevel, emergency.target/rescue.target均可进入.

紧急模式加载了带有只读根文件系统文件系统的最小环境，没有启用任何网络或其他服务; 但救援模式尝试挂载所有本地文件系统并尝试启动一些重要的服务，包括网络.