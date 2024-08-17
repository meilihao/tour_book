# nvidia + deepin 20.1

> 深度显卡驱动管理器需在安装nvidia-detect后才能探测到nvidia显卡

1. 探测是否存在nvidia显卡
```bash
apt install nvidia-detect
nvidia-detect # 其实直接使用`lspci | egrep 'VGA|3D'`也可
```

1. 禁用NVIDIA开源驱动nouveau
```bash
sudo bash -c "echo blacklist nouveau > /etc/modprobe.d/blacklist-nvidia-nouveau.conf"
sudo bash -c "echo options nouveau modeset=0 >> /etc/modprobe.d/blacklist-nvidia-nouveau.conf"
sudo update-initramfs -u # for ubuntu / sudo dracut --force # for fedora, rhel, centos
sudo reboot
```

1. 重启后查看nouveau是否已被禁用
```bash
lsmod |grep -i nouveau
```

1. 清理旧nvidia驱动和大黄蜂方案
```bash
sudo apt --purge remove nvidia* bumblebee* # nvidia-detect不算
```

1. 安装nvidia
前提:
```bash
sudo  dnf install kernel-devel libglvnd
```

change to non GUI mode: Ctrl+Alt+F2 (works on DeepinOS) or Ctrl+Alt+F1(works on Other Distro)

**经测试小米笔记本GeForce 940MX使用nvidia官方驱动成功(需结合下方的修改`lightdm.conf+xorg.conf+.xinitrc`), 而deepin 20.1 apt驱动失败**

nvidia官方驱动:
```bash
chmod +x NVIDIA-Linux-x86_64-455.45.01.run # NVIDIA-Linux-x86_64-455.45.01.run包含了nvidia-smi
sudo systemctl stop lightdm.service
sudo init 3  #切换运行级别3来运行驱动安装程序（不切换可能安装失败）
sudo ./NVIDIA-Linux-x86_64-455.45.01.run # 
安装过程中，开始两个选项选YES，第三个选项选中间那个带over的，最后一步问是否重启系统，选NO
sudo reboot
```

> 安装里的选项里，只有yes和no的选yes

deepin 20.1 apt驱动:
```bash
sudo systemctl stop lightdm.service
sudo init 3
sudo apt install nvidia-driver nvidia-settings nvidia-smi xserver-xorg-video-nvidia
sudo reboot
```

1. 查看nvidia驱动是否生效
```bash
sudo nvidia-smi # 看是否使用了显存即可
```

---
此时如果nvidia驱动不能正常工作, 那么进行如下操作:
1. sudo vim /etc/lightdm/display_setup.sh
```bash
#!/usr/bin/env bash
xrandr --setprovideroutputsource modesetting NVIDIA-0
xrandr --auto
xrandr --dpi 96
```
1. sudo chmod +x /etc/lightdm/display_setup.sh
1. sudo vim /etc/lightdm/lightdm.conf
```bash
display-setup-script=/etc/lightdm/display_setup.sh # 找到 display-setup-script这一行，去掉前面的注释，将display_setup.sh路径赋给它
```

1. 重启再次执行nvidia-smi

---
如果NVIDIA的显卡并没有工作，显存一点都没占用，显示`No running Processes found`, 则还要做下面的工作:
1. `lspci | egrep 'VGA|3D'`获取到显卡的BusID, `00:02:00 --> 0:2:0  01:00:00 --> 1:0:0`
1. sudo vim /etc/X11/xorg.conf, 没有xorg.conf时要新建一个
```conf
Section "Module"
    Load "modesetting"
EndSection

Section "Device"
    Identifier "nvidia"
    Driver "nvidia"
    BusID "PCI:1:0:0"
    Option "AllowEmptyInitialConfiguration"
EndSection
```
1. vim ~/.xinitrc
```bash
xrandr --setprovideroutputsource modesetting NVIDIA-0
xrandr --auto
xrandr --dpi 96
```
1. 重启再次执行nvidia-smi

---
所有尝试失败后建议回滚上述已操作过的步骤并清除nvidia驱动, 回退到nouveau或继续使用集显驱动.

## FAQ
### `sudo ./NVIDIA-Linux-x86_64-455.45.01.run`报`The NVIDIA probe routine was not called for 1 device(s)`
通过`locale nvidia`清理**所有**旧安装的nvidia文件后重新安装即可.

### `apt install nvidia-driver-525`重启黑屏, 桌面只有鼠标的箭头
os: ubuntu 22.04

```bash
# cat /etc/X11/default-display-manager
lightdm
# systemctl start lightdm # 此时就能正常进入桌面, 且之后重启也不会再黑屏
# nvidia-smi # 查看是否正确
```

后面一次重启还是出现黑屏, 查看lightdm已启动, 重启lightdm后能进入桌面但下方的任务栏消失. 重启后选择kernel 5.19(原先是6.1.28)一切正常.
