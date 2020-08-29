# uefi
参考:
- [不需要硬件也可以开发UEFI](https://zhuanlan.zhihu.com/p/107360611)
- [BIOSandSecureBootAttacksUncovered_eko10.pdf](/misc/pdf/BIOSandSecureBootAttacksUncovered_eko10.pdf)
- [GRUB (简体中文)](https://wiki.archlinux.org/index.php/GRUB_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87)#Chainload_%E4%B8%80%E4%B8%AA_Arch_Linux_.efi_%E6%96%87%E4%BB%B6) 
- [Unified Extensible Firmware Interface (简体中文)](https://wiki.archlinux.org/index.php/Unified_Extensible_Firmware_Interface_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87))

- EDK2：    Intel发起的UEFI开发环境
- OVMF：    基于EDK2的开源虚拟机(qemu)UEFI固件
- seaBios： 基于传统BIOS的开源虚拟机(qemu)固件


## efibootmgr
UEFI 规范定义了名为 UEFI 启动管理器的一项功能, Linux发行版包含名为efibootmgr 的工具，可用于更改 UEFI 启动管理器的配置, 可用`efibootmgr -v`查看启动项细节.

### example
```bash
# efivar # 可检测是否支持efi by `apt install efivar`
# efibootmgr -v
BootCurrent: 0000 # 目前从“启动菜单”的哪个项上启动
Timeout: 0 seconds # 如果固件的 UEFI 启动管理器显示了类似启动菜单的界面，那么这一行表示继续启动默认项之前的超时
BootOrder: 0000,0003,0002,2001,2002,2003 # UEFI 固件将按照BootOrder 中列出的顺序，尝试从“启动菜单”中的每个“项”进行启动, 其余输出显示了实际的启动项
Boot0000* deepin	HD(1,GPT,8c465477-4444-4e2a-9306-6526f24cae36,0x800,0x100000)/File(\EFI\deepin\shimx64.efi)
Boot0002* Linpus lite	HD(1,GPT,8c465477-4444-4e2a-9306-6526f24cae36,0x800,0x100000)/File(\EFI\Boot\grubx64.efi)RC
Boot0003* ubuntu	HD(1,GPT,8c465477-4444-4e2a-9306-6526f24cae36,0x800,0x100000)/File(\EFI\ubuntu\grubx64.efi)RC
Boot2001* EFI USB Device	RC
Boot2002* EFI DVD/CDROM	RC
Boot2003* EFI Network	RC
# qemu-system-x86_64 -bios "/usr/share/ovmf/OVMF.fd" # 查看uefi
```

## OVMF
参考:
- [PCI passthrough via OVMF (简体中文)](https://wiki.archlinux.org/index.php/PCI_passthrough_via_OVMF_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87))
- [How to run OVMF](https://github.com/tianocore/tianocore.github.io/wiki/How-to-run-OVMF)
- [kraxel可下载已编译好的ovmf/已安装好os的qcow2 image](https://www.kraxel.org/repos/), 比如[edk2.git-ovmf-x64-0-20200515.1440.gcbccf99592.noarch.rpm](https://www.kraxel.org/repos/jenkins/edk2/edk2.git-ovmf-x64-0-20200515.1440.gcbccf99592.noarch.rpm), 使用rpm2cpio解压rpm即可
- [redhat ovmf whitepaper](http://people.redhat.com/~lersek/ovmf-whitepaper-c770f8c.txt)

edk2.git-ovmf-x64-0-20200515.1440.gcbccf99592.noarch.rpm解压说明:
```bash
# → ll
总用量 15M
-rw-r--r-- 1 chen chen 1.9M 8月  29 16:19 OVMF_CODE-need-smm.fd
-rw-r--r-- 1 chen chen 1.9M 8月  29 16:19 OVMF_CODE-pure-efi.fd
-rw-r--r-- 1 chen chen 1.9M 8月  29 16:19 OVMF_CODE-with-csm.fd
-rw-r--r-- 1 chen chen 128K 8月  29 16:19 OVMF_VARS-need-smm.fd
-rw-r--r-- 1 chen chen 128K 8月  29 16:19 OVMF_VARS-pure-efi.fd
-rw-r--r-- 1 chen chen 128K 8月  29 16:19 OVMF_VARS-with-csm.fd
-rw-r--r-- 1 chen chen 2.0M 8月  29 16:19 OVMF-need-smm.fd
-rw-r--r-- 1 chen chen 2.0M 8月  29 16:27 OVMF-pure-efi.fd
-rw-r--r-- 1 chen chen 2.0M 8月  29 16:19 OVMF-with-csm.fd
-rw-r--r-- 1 chen chen 2.5M 8月  29 16:19 UefiShell.iso
```

smm指系统管理模式（System Management mode）是Intel在80386SL之后引入x86体系结构的一种CPU的执行模式, SMM模式对操作系统透明，换句话说，操作系统根本不知道系统何时进入SMM模式，也无法感知SMM模式曾经执行过.

支持安全启动的UEFI, 就需要smm的支持, qemu + smm看[这里](https://github.com/tianocore/tianocore.github.io/wiki/Testing-SMM-with-QEMU,-KVM-and-libvirt).

csm指兼容支持模块(Compatibility Support Module), 是UEFI的一个特殊模块，对于不支持UEFI的系统提供兼容性支持.

OVMF-pure-efi = OVMF_CODE-pure-efi + OVMF_VARS-pure-efi, 即程序与配置分离. 通常qemu使用`OVMF-pure-efi.fd`, 而libvirt使用`nvram = [ "/usr/share/edk2.git/ovmf-x64/OVMF_CODE-pure-efi.fd:/usr/share/edk2.git/ovmf-x64/OVMF_VARS-pure-efi.fd"]`. 因为针对多个虚拟机场景，对于loader来说，属性是readonly，可以共用，对于变量文件，要复制到相应的目录，避免冲突.

pure-efi推测是仅efi的意思, 由此递推:
need-smm = pure-efi + smm
with-csm = pure-efi + csm

[qemu的设置uefi/bios的参数](https://github.com/tianocore/edk2/blob/master/OvmfPkg/README):
- -pflash : 支持仿真flash, 可以将uefi vars保存到该flash中, **推荐**
- -bios : 仅能模拟部分uefi vars, 且重启后会非易失性(即保存的)uefi vars会丢失
- -L : 用于指定bios.bin的目录, 而bios.bin由OVMF.fd重命名或软连接而来, 此外与`-bios`相同

UefiShell.iso:
由于OVMF并未附带安装任何SecureBoot密钥，要实现真正的secure boot, 我们需要执行一个EFI程序来导入一系列证书, 来模仿MS认证的UEFI机器.
 UefiShell.iso就是这样一个不错的iso工具, 它里面内置了UEFI shell和EFI程序. 通过`qemu-system-x86_64 -M q35 -pflash OVMF-need-smm.fd -cdrom /usr/share/edk2/ovmf/UefiShell.iso -m 400`, 且没有`-kernel, -initrd`选项. 如果一切顺利, 就可在shell中添加key了.

```shell
Shell> fs0:
FS0:\> EnrollDefaultKeys.efi
FS0:\> reset
```

vm重新启动后会在dmesg中看到"Secure boot enabled"了.

> 不是使用OVMF-need-smm.fd时, EnrollDefaultKeys.efi会报错: "error: GetVariable("SetupMode", 8BE4DF61-93CA-11D2AA0D-00E098032B8C): Not Found". 即使换成了"OVMF-need-smm.fd"还是会卡在qemu的启动界面上(ovmf之前), 原因可看[这里](https://askbot.fedoraproject.org/en/question/86384/wiki-using-uefi-with-qemu/)和[SMM support](https://github.com/tianocore/edk2/blob/master/OvmfPkg/README), 但对照检查后还是不行, 只能使用OVMF-pure-efi.fd了.

### 编译
参考[edk2-ovmf](https://www.archlinux.org/packages/extra/any/edk2-ovmf/)右上角的[Source Files](https://github.com/archlinux/svntogit-packages/tree/packages/edk2/trunk)中的[构建脚本](https://github.com/archlinux/svntogit-packages/blob/packages/edk2/trunk/PKGBUILD).

```bash
$ git clone --depth 1 -b edk2-stable202005 git@github.com:tianocore/edk2.git
$ git submodule update --init --recursive # 下载`.gitmodules`中的依赖, 以后用`git submodule update --remote`更新submodule
```