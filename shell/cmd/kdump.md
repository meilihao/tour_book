# kdump

## cmds
```bash
# kdump-config show # 查看kdump配置
```

## FAQ
### grub追加`crashkernel=512M-:384M`
kernel 4.4(ubuntu 14.04)不支持该形式, 直接使用`crashkernel=384M`即可.

### ubuntu 14.04 /boot/grub/grub.cfg添加了`crashkernel=384M`, reboot后未在`/proc/cmdline`中出现
参考:
- [Kernel Crash Dump](https://ubuntu.com/server/docs/kernel-crash-dump)

需要调用`dpkg-reconfigure kexec-tools`和`dpkg-reconfigure kdump-tools`配置这两工具, 配置时都选`yes`即可.