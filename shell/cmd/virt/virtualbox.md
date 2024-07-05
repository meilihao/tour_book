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
