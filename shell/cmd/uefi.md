# uefi
参考:
- [不需要硬件也可以开发UEFI](https://zhuanlan.zhihu.com/p/107360611)
- [BIOSandSecureBootAttacksUncovered_eko10.pdf](/misc/pdf/BIOSandSecureBootAttacksUncovered_eko10.pdf)
- [GRUB (简体中文)](https://wiki.archlinux.org/index.php/GRUB_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87)#Chainload_%E4%B8%80%E4%B8%AA_Arch_Linux_.efi_%E6%96%87%E4%BB%B6) 
- [Unified Extensible Firmware Interface (简体中文)](https://wiki.archlinux.org/index.php/Unified_Extensible_Firmware_Interface_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87))

- EDK2：    Intel发起的UEFI开发环境
- OVMF：    基于EDK2的开源虚拟机(qemu)UEFI固件
- seaBios： 基于传统BIOS的开源虚拟机(qemu)固件

## 真正的UEFI启动逻辑
参考:
- [UEFI启动背后的原理](https://i-m.dev/posts/20200808-135558.html)
- [关于Windows Boot Manager、Bootmgfw.efi、Bootx64.efi、bcdboot.exe 的详解](http://bbs.wuyou.net/forum.php?mod=viewthread&tid=303679)

```code
hdparmdef 从EFI文件启动(文件路径):
    根据文件路径加载EFI文件
    if 开启了SecureBoot功能:
        if not SecureBoot校验通过:
            启动失败报错退出 # 一般可以在屏幕上看到错误信息
    执行EFI # 一般来说也就从这里进入了操作系统或者Grub，所以不会返回

def 从存储设备启动(存储设备):
    # 可以是是GPT分区，也可以是MBR分区
    for 分区 in 存储设备:
        # 分区类型在MBR/GPT的分区表中标注，其本质是个普通FAT文件系统分区
        if 分区类型 == EFI系统分区:
            # 注意，FAT文件系统路径不区分大小写，且由于出自Windows，所以其路径习惯用反斜杠\而非斜杠/
            # 对于非x86_64类型的硬件，默认路径也有差异
            if 分区根目录下存在"\EFI\Boot\Bootx64.efi"文件:
                从EFI文件启动(分区索引 + 该文件在分区下路径)

def 从EFI启动项启动(启动项):
    if 该启动项指向的是一个EFI文件路径:
        从EFI文件启动(启动项所指文件路径)
    elif 启动项指向的是一个存储设备:
        从存储设备启动(启动项所指设备)

def 启动():
    根据主板NVRAM(非易失性随机访问存储器)中启动项配置，以及扫描的硬件（比如刚插入装机U盘等），生成启动顺序表
    if 用户无操控:
        for 启动项 in 启动顺序表:
            从EFI启动项启动(启动项) # 成功则返回，失败可能会返回，也可能是报错退出
    else:
        根据启动顺序表选择界面
        if 用户选择了一个启动项:
            # 这对应了手动从列表中选择一个启动设备
            从EFI启动项启动(启动项)
        elif 用户手动选择了一个具体的EFI文件:
            # 这对应了手动浏览EFI分区并选择某一个具体的文件
            从EFI文件启动(文件路径)
    这次启动失败了，死机？退出主界面？或者重新尝试？结果因机器而异


硬件扫描和准备
启动()
```

总结说来，UEFI并不是单纯地找到启动设备上的EFI分区里的某固定文件路径就启动了，而是会结合NVRAM中的配置来完成整个过程。NVRAM中的启动配置包括一系列启动项和响应的预设优先顺序.

## efibootmgr
UEFI规范定义了名为UEFI启动管理器的一项功能,它的定义如下:”UEFI启动管理器是一种固件策略引擎,可通过修改固件架构中定义的**全局NVRAM变量来进行配置**. 启动管理器将尝试按全局NVRAM变量定义的顺序依次加载UEFI驱动和UEFI应用程序(包括UEFI操作系统启动装载程序).”

简单来说,UEFI启动管理器可以管理UEFI的启动菜单,例如调整顺序,删除,添加.它最大的优点就是可以从上层系统修改UEFI的行为

Linux发行版包含UEFI启动管理器即efibootmgr, 可用`efibootmgr -v`查看启动项细节.

UEFI会在`${EF00,EFI System分区}`的`/EFI`文件夹中查找所有文件夹,并搜寻efi文件.

```
$ ls                
Boot  grub  Microsoft
```

如上,除了Boot文件夹,还有grub文件夹和Microsoft文件夹,这也是开机启动顺序所显示的名称. 而这两个文件夹则分别存放着grub和Windows Boot Manager, 这两个bootloader的efi文件:/EFI/grub/grubx64.efi和/EFI/Miscosoft/Boot/bootmgfw.efi

Boot是计算机默认引导文件所在的目录, 在主板主板NVRAM异常时起备份作用.

事实上, Boot/bootx64.efi是通用名,任何其他的引导文件都可以改成这个名称,放在/EFI/Boot目录下,从而成为计算机默认引导文件.

**注意**: 如果是在U盘上写入grub的话,在使用grub-install时需要添加--removable选项(将grubx64.efi文件改名为bootx64.efi);如果之前已经安装过grub需要重装时,则需要添加--recheck选项(删除原有grub相关的所有文件,再安装).

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

### 安装grub
1. 安装grub首先要确定有ESP分区,并将其挂载到一个目录中,例如/boot/efi.
1. 然后,使用grub-install命令生成grubx64.efi文件,并将该grub的模块放在/boot/efi/EFI/lfs下.

	`grub-install --target=x86_64-efi --efi-directory=/boot/efi --bootloader-id=lfs`
1. 再生成grub.cfg文件, 生成的行为实际上是综合了/etc/default/grub的选项和/etc/grub.d/里的脚本.

	`grub-mkconfig -o /boot/grub/grub.cfg`

	grub-mkconfig会执行/etc/grub.d/30_os-prober脚本文件,该文件会搜寻所有可加载的内核,并生成启动项.

### Grub是如何被加载的
```bash
# tree /boot/efi/EFI/ubuntu/
/boot/efi/EFI/ubuntu/
├── grub.cfg
└── grubx64.efi
# cat /boot/efi/EFI/ubuntu/grub.cfg 
search.fs_uuid 5786c7db-036e-42ff-a0e0-a676133b3dcf root 
set prefix=($root)'/grub'
configfile $prefix/grub.cfg
```
\EFI\ubuntu\grubx64.efi同目录下有一个\EFI\ubuntu\grub.cfg.

它的第一行的uuid是我Ubuntu系统所在boot分区的uuid，找到这个分区后，根据第二第三行的配置，Grub的EFIb部件会从中读取/boot/grub/grub.cfg文件作为grub启动配置（该文件其实就是执行grub-update的输出）.

然后，就是展示熟悉的Grub窗口，在然后根据选择的Grub启动项进入操作系统，不过这已经不是UEFI的范畴了.

## UEFI Shell
参考:
- [UEFI Shell命令操作总结](https://blog.csdn.net/kair_wu/article/details/48342093)
- [用efibootmgr管理UEFI启动项，添加丢失的启动项](https://blog.csdn.net/Pipcie/article/details/79971337)

Shell命令的通用选项:
- -b : 输出信息分屏显示

### 相关命令
- map : 显示设备映射的列表，即可用文件系统（fsN）和存储设备（blkN）的名称
- edit FS0:\EFI\refind\refind.conf : 类似nano的编辑器
- help : 获取help

### 使用UEFI Shell引导U盘/磁盘启动
一般UEFI启动分区是硬盘最前端的FAT分区, 因此首先要找到存放启动文件的分区，依次输入下列命令：
1. `fs0:`（ufeshell使用fs:x方式表示，序号从0开始, **冒号不能丢**）
1. ls（查看文件结构，确保其中有EFI目录，如果没有则重新选择其他fs驱动器）
1. cd EFI\Boot（进入引导目录，查看引导文件）
1. BOOTX64.efi（直接运行uefi可执行程序）

uefi存储标示:
- fsX : filesystem
- blkX : block device 或者 data storage device

## OVMF(Open Virtual Machine Firmware)
参考:
- [PCI passthrough via OVMF (简体中文)](https://wiki.archlinux.org/index.php/PCI_passthrough_via_OVMF_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87))
- [How to run OVMF](https://github.com/tianocore/tianocore.github.io/wiki/How-to-run-OVMF)
- [kraxel可下载已编译好的ovmf/已安装好os的qcow2 image](https://www.kraxel.org/repos/), 比如[edk2.git-ovmf-x64-0-20200515.1440.gcbccf99592.noarch.rpm](https://www.kraxel.org/repos/jenkins/edk2/edk2.git-ovmf-x64-0-20200515.1440.gcbccf99592.noarch.rpm), 使用rpm2cpio解压rpm即可, 自编译可参考他的[edk2.git.spec.template](https://git.kraxel.org/cgit/jenkins/edk2/tree/edk2.git.spec.template), 他的[官网](https://www.kraxel.org/repos/)中的section "Using the repo"有简短说明.
- [redhat ovmf whitepaper](http://people.redhat.com/~lersek/ovmf-whitepaper-c770f8c.txt)

> `sudo apt install ovmf` on deepin v20.

OVMF is an EDK II based project to enable UEFI support for Virtual Machines. OVMF contains sample UEFI firmware for QEMU and KVM.

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

	`-pflash OVMF-pure-efi.fd`=`-drive if=pflash,format=raw,readonly,file=/usr/share/ovmf/OVMF_CODE-pure-efi.fd \
	-drive if=pflash,format=raw,file=my_uefi_vars.bin`, my_uefi_vars.fd is copy from OVMF_VARS-pure-efi.fd.
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
参考:
- [edk2-ovmf](https://www.archlinux.org/packages/extra/any/edk2-ovmf/)右上角的[Source Files](https://github.com/archlinux/svntogit-packages/tree/packages/edk2/trunk)中的[构建脚本](https://github.com/archlinux/svntogit-packages/blob/packages/edk2/trunk/PKGBUILD).
- [gentoo edk2-ovmf-999999.ebuild](https://gitweb.gentoo.org/repo/gentoo.git/plain/sys-firmware/edk2-ovmf/edk2-ovmf-999999.ebuild)

```bash
$ sudo apt install acpica-tools nasm # acpica-tools即以前的iasl
$ git clone --depth 1 -b edk2-stable202005 git@github.com:tianocore/edk2.git
$ cd edk2-stable202005
$ git submodule update --init # 下载`.gitmodules`中的依赖, 以后用`git submodule update`更新submodule. 相关doc在edk2-stable202005/ReadMe.rst
$ . edksetup.sh
$ make -C BaseTools
$ --- with debug
$ OvmfPkg/build.sh -a X64 # with debug by doc edk2-stable202005/OvmfPkg/README
$ cp Build/OvmfX64/DEBUG_GCC5/FVOVMF.fd ~/test/uefi
$ --- no debug
$ OvmfPkg/build.sh -a X64 -b RELEASE -t GCC5 # no debug
$ cp Build/OvmfX64/RELEASE_GCC5/FV/OVMF.fd ~/test/uefi
$ --- test efi
$ cd ~/test/uefi
$ qemu-system-x86_64 -pflash ./OVMF.fd
```

## FAQ
### 检查os是否uefi启动
```
# efivar -l # by read efivarfs = /sys/firmware/efi/efivars
```

### grubx64.efi和shimx64.efi区别
- EFI\ubuntu\grubx64.efi : 
- EFI\ubuntu\shimx64.efi : 是安全启动选项


	启动shimx64.efi就像启动一样grubx64.efi. 在启用了安全启动的计算机上，启动shimx64.efi会间接启动GRUB，而无法直接启动grubx64.efi.
