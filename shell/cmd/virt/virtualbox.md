# virtualbox
ref:
- [VirtualBox 7种网络模式配置指南](https://segmentfault.com/a/1190000043810778)

## FAQ
### 退出VirtualBox VM独占的键盘和鼠标
`右边Alt` + `右边Ctrl`

默认设置不独占: 管理->全局设定->热键->取消勾选"自动独占键盘"

### 开启嵌套虚拟化
bios已开启虚拟化, 但virtualbox gui无法开启嵌套虚拟化

```bash
$ VBoxManage list vms
$ VBoxManage modifyvm truenas_test --nested-hw-virt on
```

### VBoxManage命令报`The vboxdrv kernel module is not loaded`
`sudo /sbin/vboxconfig`

### 设置磁盘image uuid(不是硬盘序列号, 而是vbox管理磁盘的唯一id)
`VBoxManage internalcommands sethduuid /home/chen/data/vm/virtualbox/TrueNAS-SCALE-24.04.1.1/TrueNAS-SCALE-24.04.1.1_1.vdi`

改完后会导致vm启动失败, 报`UUID {1fa8697f-299e-46b2-806e-cd6dcbbde74d} of the medium '/home/chen/data/vm/virtualbox/TrueNAS-SCALE-24.04.1.1/TrueNAS-SCALE-24.04.1.1_1.vdi' does not match the value {0cb712ab-6f6d-4cd8-92cc-a6e1cda569b6} stored in the media registry('/home/chen/.config/virtualbox/VirtualBox.xml').`.
解决方法: 管理->工具->虚拟介质管理
1. 选中它, 点`释放`
1. 选中它, 点`删除`
1. 点`注册`, 重新添加该盘

或使用`VBoxManage internalcommands sethduuid <路径\xxx.vdi> <UUID>`指定回uuid.

### 设置磁盘serial number
ref:
- [Configuring the Hard Disk Vendor Product Data (VPD)](https://www.virtualbox.org/manual/ch09.html#changevpd)

**未成功(无法确定setextradata的config path), 设置磁盘serial number场景建议使用virtmanager**

```bash
$ VBoxManage getextradata TrueNAS-SCALE-24.04.1.1 [enumerate]
$ VBoxManage showvminfo TrueNAS-SCALE-24.04.1.1 --details
$ VBoxManage setextradata TrueNAS-SCALE-24.04.1.1 VBoxInternal/Devices/lsilogic/0/Port1/Config/SerialNumber 8CTR22V3 # 设置后启动失败???
```

> `Type: LsiLogic, Instance: 0`(from TrueNAS-SCALE-24.04.1.1.vbox)=`/Devices/lsilogic/0`

> 磁盘serial来自其Vendor product data (VPD).

> ahci是sata硬盘，Port1代表第二个硬盘; ide是piix4ide, 其接口是Primarymaster/Primaryslave等

> setextradata在vm xml(from `VBoxManage showvminfo`中的`Config file`)中的`<ExtraData>`

> VBoxManage setextradata的value为空时即删除该设置项

### 设置`DmiSystemSerial`
`VBoxManage setextradata xxx VBoxInternal/Devices/pcbios/0/Config/DmiSystemSerial <id>`

### 安装扩展包
管理->工具->扩展包管理器, 选择安装即可