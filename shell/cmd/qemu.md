# qemu
Qemu是一个广泛使用的开源计算机仿真器和虚拟机.

操作:
```
$ sudo yum install qemu -y
$ sudo apt-get install qemu
$ qemu- + <tab> 查看支持的arch
$ qemu-system-x86_64 -boot menu=on,splash-time=15000 # 查看seabios version
```

> qemu-x86_64: 仅仅模拟CPU; qemu-system-x86_64: 模拟整个PC
> 在qemu新版本(比如2.8.1)中已经将qemu-kvm模块完全合并到qemu中去. 因此当需要使用专kvm特性时候，只需要qemu-system-x86_64 启动命令中增加`–enable-kvm`参数即可.

编译选项:
- xx-softmmu和xx-linux-user的区别和关系

    xxx-softmmu将编译qemu-system-xxx,这是一个用于xxx架构(系统仿真)的仿真机器.重置时,起点将是该架构的重置向量.而xxx-linux-user则编译qemu-xxx,它允许您在xxx架构中运行用户应用程序(用户模式仿真).这将寻找用户应用程序的主要功能,并从那里开始执行

    aarch64_be-linux-user是大端arm

编译:
```bash
$ git ls-remote -t https://mirrors.tuna.tsinghua.edu.cn/git/qemu.git
$ git clone --depth=1 -b v5.1.0  https://mirrors.tuna.tsinghua.edu.cn/git/qemu.git && cd qemu
git submodule init
git submodule update --recursive
$ ./configure --target-list="x86_64-softmmu,x86_64-linux-user,aarch64-softmmu,aarch64-linux-user,aarch64_be-linux-user,riscv64-softmmu,riscv64-linux-user" \
			  --enable-kvm \
			  --enable-sdl \
              --mandir="\${prefix}/share/man"
    		#   --enable-opengl \
            #   --enable-gtk
$ make && sudo make install
$ qemu-system-x86_64 --version
```

> `--enable-sdl`是为了启用qemu console gui, 便于查看early boot, 比如uefi/grub信息.
> qemu 5.1.0源码编译一次后重新make会卡住, 即使提前make clean过也不行, 只能用一份新源码来编译.

生成的相关程序:
- ivshmem-client/server：这是一个 guest 和 host 共享内存的应用程序，遵循 C/S 的架构
- qemu-ga：这是一个不利用网络实现 guest 和 host 之间交互的应用程序（使用 virtio-serial），运行在 guest 中
- qemu-io：这是一个执行 Qemu I/O 操作的命令行工具
- qemu-system-x86_64：Qemu 的核心应用程序，虚拟机就由它创建的

    使用 qemu-system-x86 来启动 x86 架构的虚拟机:`qemu-system-x86_64 test-vm-1.qcow2`, 因为 test-vm-1.qcow2 中并未给虚拟机安装操作系统，所以会提示 “No bootable device”，无可启动设备.

    启动 VM 安装操作系统镜像: `qemu-system-x86_64 -m 2048 -enable-kvm test-vm-1.qcow2 -cdrom ./Centos-Desktop-x86_64-20-1.iso`:
    - -m 指定虚拟机内存大小，默认单位是 MB
    - -enable-kvm 使用 KVM 进行加速
    - -cdrom 添加 fedora 的安装镜像

    iso安装完成后重起虚拟机便会从硬盘 ( test-vm-1.qcow2 ) 启动.
- qemu-img：创建虚拟机镜像文件的工具

    虚拟机镜像用来模拟虚拟机的硬盘，在启动虚拟机之前需要创建镜像文件. 镜像文件创建完成后，可使用 qemu-system-x86 来启动x86 架构的虚拟机.

    `qemu-img create -f qcow2 test-vm-1.qcow2 10G`:
    - -f 选项用于指定镜像的格式，qcow2 格式是 Qemu 最常用的镜像格式，采用来写时复制技术来优化性能
    - test-vm-1.qcow2 是镜像文件的名字
    - 10G是镜像文件大小

- qemu-nbd：磁盘挂载工具

  ```bash
  # qemu-img create -f <fmt> <image filename> <size of disk>
  qemu-img create -f qcow2 lfs.img 8G
  sudo modprobe -v nbd
  sudo qemu-nbd -c /dev/nbd0 lfs.img
  sudo gdisk /dev/nbd0
  # see result
  sudo  gdisk -l /dev/nbd0
  qemu-nbd --disconnect /dev/nbd0
  ```
  > 通过`modinfo nbd`获知默认一个nbd设备的最大分区数是16, 可通过max_part修改.

如何更便捷地创建虚拟机呢? 答案就是libvirt. 这次就不再一个个指定虚拟机启动的参数，而是用 libvirt. 它管理 qemu 虚拟机，是基于 XML 文件，这样容易维护. xml中定义了kvm中domain的配置信息，可以使用virt-install来生成. 首先，需要安装 libvirt
```bash
# apt install libvirt-bin
# apt install virtinst
```

> kvm中vm的grub kernel启动参数添加`console=ttyS0`再执行update-grub并重启后, 就可以通过`virsh console ${vm_name}`，进入机器的控制台，可以不依赖于 SSH 和 IP 地址进行登录.

## qemu-system-x86_64
参考:
- [qemu-system-x86_64命令总结](http://blog.leanote.com/post/7wlnk13/%E5%88%9B%E5%BB%BAKVM%E8%99%9A%E6%8B%9F%E6%9C%BA)

### 选项
参考:
- [qemu命令行参数](https://blog.csdn.net/GerryLee93/article/details/106475710)

- display

    - curses : 仅用于文本模式(text mode is used only with BIOS firmware), 当出现"640 x 480 Graphic mode"时表示guest已切换到图形模式.
    - sdl : qemu console gui
- -cpu <cpu>/help : help可获取qemu支持模拟的cpu
- -M <machine>/help : 当前版本的Qemu工具支持的开发板列表
- -s : 设置gdbserver的监听端口, 等同于`-gdb tcp::1234`
- -S : 启动时cpu仅加电, 但不继续执行, 相当于将断点打在CPU加电后要执行的第一条指令处，也就是BIOS程序的第一条指令. 必须在qemu monitor输入`c`才能继续. 未使用`-monitor`时, 按`Ctrl+Alt+2`可进入qemu的monitor界面,`Ctrl+Alt+1`回到qemu
- -serial stdio : redirects the virtual serial port to the host's terminal input/output, 丢失early boot信息即加电到出现终端登入界面间的信息.
- -monitor

    tcp – raw tcp sockets, **推荐**.
    telnet – the telnet protocol is used instead of raw tcp sockets. This is the preferred option over tcp as you can break out of the monitor using Ctrl-] then typing quit. You can’t break out of the monitor like this after connecting with the raw socket option
    10.1.77.82 – Listen on this host/IP only. You can use 127.0.0.1 if you want to only allow connections locally. If you want to listen on any ip address on the server, just leave this blank so you end up with two consecutive colons ie `::`.
    4444 – port number to listen on.
    server – listening in server mode
    nowait – qemu will wait for a client socket application to connect to the port before continuing unless this option is used. In most cases you’ll want to use the nowait option.

## qemu monitor
滚屏: ctrl + 上/下

## 操作
### 模拟cpu加电
```
# qemu-system-x86_64  -S -monitor tcp::4444,server,nowait # qemu起来的窗口太小, `info registers`展示不完全, 因此使用qemu monitor来解决.
# nc localhost 4444 # 另一个terminal
QEMU 2.8.1 monitor - type 'help' for more information
(qemu) info registers
info registers
EAX=00000000 EBX=00000000 ECX=00000000 EDX=00000663
ESI=00000000 EDI=00000000 EBP=00000000 ESP=00000000
EIP=0000fff0 EFL=00000002 [-------] CPL=0 II=0 A20=1 SMM=0 HLT=0
ES =0000 00000000 0000ffff 00009300
CS =f000 ffff0000 0000ffff 00009b00
SS =0000 00000000 0000ffff 00009300
DS =0000 00000000 0000ffff 00009300
FS =0000 00000000 0000ffff 00009300
GS =0000 00000000 0000ffff 00009300
LDT=0000 00000000 0000ffff 00008200
TR =0000 00000000 0000ffff 00008b00
GDT=     00000000 0000ffff
IDT=     00000000 0000ffff
CR0=60000010 CR2=00000000 CR3=00000000 CR4=00000000
DR0=0000000000000000 DR1=0000000000000000 DR2=0000000000000000 DR3=0000000000000000 
DR6=00000000ffff0ff0 DR7=0000000000000400
EFER=0000000000000000
FCW=037f FSW=0000 [ST=0] FTW=00 MXCSR=00001f80
FPR0=0000000000000000 0000 FPR1=0000000000000000 0000
FPR2=0000000000000000 0000 FPR3=0000000000000000 0000
FPR4=0000000000000000 0000 FPR5=0000000000000000 0000
FPR6=0000000000000000 0000 FPR7=0000000000000000 0000
XMM00=00000000000000000000000000000000 XMM01=00000000000000000000000000000000
XMM02=00000000000000000000000000000000 XMM03=00000000000000000000000000000000
XMM04=00000000000000000000000000000000 XMM05=00000000000000000000000000000000
XMM06=00000000000000000000000000000000 XMM07=00000000000000000000000000000000
```

## qcow2
### 修改qcow2 image的方法
1. libguestfs-tools
```
$ sudo apt-get install libguestfs-tools
$ guestmount  -a  x.qcow2 -i  --rw  {mount_dir} # 挂载qcow2
$ sudo umount {mount_dir}
$ guestfish --rw -a centos6.5-minimal.qcow2 # 进入qcow2直接修改即可, 与系统进入修复模式类似.
><fs> run
><fs> list-filesystems # 查找文件系统
/dev/sda1: ext4
><fs> mount /dev/sda1 / # 挂载文件系统
><fs> touch /etc/rc.local
><fs> edit /etc/rc.local
><fs> chmod 0755 /etc/sysconfig/modules/8021q.modules
><fs> exit
```

1. qemu-nbd

## FAQ
### qemu编译依赖
```
# sudo apt install libglib2.0-dev # RROR: glib-2.48 gthread-2.0 is required to compile QEMU
# sudo apt install libpixman-1-dev # Please install the pixman devel package
# sudo apt install flex bison
# sudo apt install libsdl2-dev # Install SDL2-devel, otherwise "VNC server running on 127.0.0.1:5900 : 缺SDL"
# 删除qemu源码, 重新解压编译 # No rule to make target 'x86_64-softmmu/config-devices.mak', needed by 'config-all-devices.mak'
```

### qemu使用32bit寄存器 on x86_64
QEMU prints the CPU state in the 32 bit format if the CPU is
currently in 32-bit mode, and in 64 bit format if it is currently
in 64-bit mode. So it simply depends what the CPU happens to be
doing at any given time.

可能是seabios最高支持到32bit的原因???.

### 如何将qcow2打内容克隆到磁盘
`qemu-img dd -f qcow2 -O raw bs=4M if=/vm-images/image.qcow2 of=/dev/sdd1`支持将qcow2 dd到磁盘

### x86 machine type选择
参考:
- [Qemu X86架构的Machine Type](https://remimin.github.io/2019/07/09/qemu_machine_type/)

可通过`qemu-system-x86_64 --machine help`查看x86支持的所以machine type.

i440fx是1996年推出的架构, 已过时. q35是2009年推出的架构, 更现代.

### [qemu mirror](https://mirrors.tuna.tsinghua.edu.cn/help/qemu.git/)

### Could not access KVM kernel module: No such file or directory
`sudo modprobe kvm-intel`

### could not insert 'kvm_intel': Operation not supported
```
# sudo dmesg | grep kvm
[783653.034467] kvm: no hardware support
```
或通过`/proc/cpuinfo`查询:
```
# egrep '^flags.*(vmx|svm)' /proc/cpuinfo
# LC_ALL=C lscpu | grep Virtualization
```

硬件不支持, 检查bios/uefi是否关闭了虚拟化支持.