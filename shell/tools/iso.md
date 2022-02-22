## iso

## centos 7.6 iso安装原理
ref:
- [Centos系统安装启动原理 - 包含"可启动盘原理"](https://keenjin.github.io/2019/10/centos_system_install_internal/)
- [CentOS7全自动安装光盘制作详解](https://cloud.tencent.com/developer/article/1452074)
- [WORKING WITH ISO IMAGES](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/anaconda_customization_guide/sect-iso-images)
- [制作自动化安装的CENTOS7 ISO](http://linux.laoqinren.net/linux/build-auto-install-centos7-iso/)
- [systemd时代的开机启动流程(UEFI+systemd) - **推荐**](https://www.junmajinlong.com/linux/systemd/systemd_bootup/)
- [CentOS 7 完整的引导流程及常见的系统故障处理](https://blog.csdn.net/weixin_44983653/article/details/97792586)

iso目录:
- isolinux  目录主要存放光盘启动时的加载kernel的bootloader isolinux.bin以及其配置文件isolinux.cfg，kernel文件vmlinuz，虚拟文件系统initrd.img
- images   目录包括了PXE启动时必要的启动映像文件
- Packages 目录存放需要安装到系统上的rpm软件包及信息
- Repodata 目录存放rpm包依赖信息
- LiveOS   目录存放了一个重要只读文件系统squashfs.img镜像，安装程序anaconda就放在镜像中
- EFI      目录存放的EFI引导模式下所需的引导文件，如BOOTX64.EFI、grubx64.efi

Linux CentOS ISO镜像安装系统的Boot过程和系统安装好后的boot过程在主体流程是一致的，只是用的bootloader有差异。目前ISO镜像安装盘使用的isolinux作为其bootloader,而正常的CentOS系统启动现阶段主流的bootloader一般是Grub2.

启动过程:
1. 系统固件初始化
1. 加载bootloader

	根据 El-Torito 规范，BIOS 会读取 ISO 的准确地址进行判断，ISO 是否可以进行引导启动（在第 71 字节地址保存引导目录）。如果判断可引导，那么就加载对应的引导程序（bootloader）.

	BIOS会读取MBR中存放的bootloader，即图3所示的isolinux目录下的isolinux.cat文件，isolinux.cat会协助进一步加载isolinux.bin。Isolinux.bin会根据isolinux.cfg配置文件内容显示相应的boot菜单选择进行kernel的加载.

	**UEFI待补充**
1. 加载kernel与initrd

	isolinux.bin完成加载后，会根据用户在boot menu中的选择进行对应kernel的加载。Kenel自身是一个压缩文件，但是在kernel文件头部嵌入有解压的代码.

	同时bootloader会加载initrd虚拟文件系统, 一般使用的文件名为initrd.img. 这个文件能够通过 Boot Loader 来加载到内存中，然后这个文件会被解压缩且在内存作为rootfs, 它能够提供一个可执行的程序，通过该程序来加载开机过程中所最需要的核心模块， 通常这些模块就是USB, RAID, LVM, SCSI等文件系统与磁盘接口的驱动程序.

	等到载入完成后，会帮助Kernel重新呼叫systemd来开始后续的正常开机流程.
1. 系统初始化

	Kernel所需的基础运行环境到此已经准备就绪，Kernel会调用虚拟根文件系统中的init(systemd)程序进行系统初始化.

	Systemd调用起来默认执行default.target, 它一般指向是initrd.target. initrd.target是需要读入一堆例如 basic.target, sysinit.target等等的硬件侦测、核心功能启用的流程，然后开始让系统顺利运作。最终才又卸除 initramfs 的小型文件系统，然后去挂载系统实际的根目录.
1. 挂载squashfs.img

	Squashfs.img镜像为一个只读压缩文件系统, squashfs.img会挂载作为CentOS安装过程的根文件系统，形成一个特殊的安装环境.
1. 启动anaconda

	Anaconda启动后就会进行系统相关RPM包的安装。Anaconda安装程序的启动是有systemd执行anaconda.target的结果


系统安装的时候是按照iso root目录下的ks.cfg文件的内容进行安装.