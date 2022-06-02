# qemu
参考:
- [Linux 虚拟化入门（一）Qemu，KVM，Virsh 概念指南](https://blog.frytea.com/archives/539/)
- [QEMU/KVM磁盘在线备份](https://cloud.tencent.com/developer/article/1468103)
- [OpenStack Virtual Machine Image Guide 之 Create images manually - 包括win+virtio](https://docs.openstack.org/image-guide/index.html)
- [openstack 构建自己的云主机镜像](https://www.programminghunter.com/article/88541281593/)
- [制作 OpenStack Linux 镜像 - 每天5分钟玩转 OpenStack](https://developer.aliyun.com/article/463345)
- [如何构建OpenStack镜像](https://www.jingh.top/2016/05/28/%E5%A6%82%E4%BD%95%E6%9E%84%E5%BB%BAOpenStack%E9%95%9C%E5%83%8F/)
- [DIB(disk image builder)是OpenStack TripleO项目的子项目，专门用于构建OpenStack镜像](https://www.jingh.top/2016/05/28/%E5%A6%82%E4%BD%95%E6%9E%84%E5%BB%BAOpenStack%E9%95%9C%E5%83%8F/)

Qemu是一个广泛使用的开源计算机仿真器和虚拟机.

**QEMU社区推荐使用功能丰富的`-device`参数来替代以前的一些参数（如`-usb`等）**.

操作:
```
$ sudo yum install qemu -y
$ sudo apt-get install qemu/qemu-system-x86 # ubuntu 20.04用qemu-system-x86
$ qemu- + <tab> 查看支持的arch
$ qemu-system-x86_64 -boot menu=on,splash-time=15000 # 查看seabios version
```

> qemu-x86_64: 仅仅模拟CPU, 即运行某种架构的程序; qemu-system-x86_64: 模拟整个PC即全系统模拟模式 (full-system emulation), 即运行某种架构系统的vm
> 在qemu新版本(比如2.8.1)中已经将qemu-kvm模块完全合并到qemu中去. 因此当需要使用专kvm特性时候，只需要qemu-system-x86_64 启动命令中增加`–enable-kvm`参数即可.

编译选项:
- xx-softmmu和xx-linux-user的区别和关系

    xxx-softmmu将编译qemu-system-xxx,这是一个用于xxx架构(系统仿真)的仿真机器.重置时,起点将是该架构的重置向量.而xxx-linux-user则编译qemu-xxx,它允许您在xxx架构中运行用户应用程序(用户模式仿真).这将寻找用户应用程序的主要功能,并从那里开始执行

    aarch64_be-linux-user是大端arm

编译:
> qemu 6.x.x编译可参考[这里](https://github.com/archlinux/svntogit-packages/blob/packages/qemu/trunk/PKGBUILD)

在QEMU源代码目录下， 运行完./configure命令后， 会生成一个config-host.mak， 在这个文件里可以看到使能了哪些选项（CONFIG_XX=y）.

```bash
$ sudo apt install -y ninja-build meson pkg-config libglib2.0-dev libpixman-1-dev libjemalloc-dev libzstd-dev liburing-dev # qemue 5.2.0开始构建系统切换到ninja+meson. libzstd-dev最低要求1.4.0, 而deepin v20是1.3.8, 可到[这里](https://packages.debian.org/buster-backports/libzstd-dev)下载libzstd1和libzstd-dev. liburing-dev 未在apt repo里, 可在[这里](https://packages.debian.org/sid/liburing-dev)下载liburing-dev和liburing1来安装
$ git ls-remote -t https://mirrors.tuna.tsinghua.edu.cn/git/qemu.git
$ git clone --depth=1 -b v5.2.0  https://mirrors.tuna.tsinghua.edu.cn/git/qemu.git && cd qemu
git submodule init
git submodule update --recursive
$ ./configure --target-list="x86_64-softmmu,x86_64-linux-user,aarch64-softmmu,aarch64-linux-user,aarch64_be-linux-user,riscv64-softmmu,riscv64-linux-user" \
              --prefix=/opt/qemu \
              --enable-kvm \
              --enable-sdl \
              --disable-xen \
              --enable-jemalloc \
              --enable-zstd \
              --enable-linux-io-uring \
              --mandir="\${prefix}/share/man"
    		#   --enable-opengl \
            #   --enable-gtk
$ make -j $(nproc) && sudo make install
$ qemu-system-x86_64 --version
# -- 避免某些旧应用(比如旧版libvirtd)使用到了旧命名qemu-kvm
$ ln -sf /usr/bin/qemu-system-x86_64 /usr/bin/qemu-kvm
$ ln -sf /usr/bin/qemu-system-x86_64 /usr/libexec/qemu-kvm
```

> 最简编译选项: `./configure --target-list="x86_64-softmmu" --enable-kvm`
> `--enable-sdl`是为了启用qemu gui, 便于查看early boot, 比如uefi/grub信息.
> qemu 5.2.0构建切换到ninja+meson后, 构建速度贼快.
> qemu 6.0.0需要libzstd>=1.4.0

生成的相关程序:
- ivshmem-client/server：这是一个 guest 和 host 共享内存的应用程序，遵循 C/S 的架构
- qemu-edid : [测试edid支持](https://www.kraxel.org/blog/2019/03/edid-support-for-qemu/)
- qemu-ga：这是一个不利用网络实现 guest 和 host 之间交互的应用程序（使用 virtio-serial），运行在 guest 中
- qemu-io：这是一个执行 Qemu I/O 操作的命令行工具
- qemu-system-x86_64：Qemu 的核心应用程序，虚拟机就由它创建的

    使用 qemu-system-x86 来启动 x86 架构的虚拟机:`qemu-system-x86_64 test-vm-1.qcow2`, 因为 test-vm-1.qcow2 中并未给虚拟机安装操作系统，所以会提示 “No bootable device”，无可启动设备.

    启动 VM 安装操作系统镜像: `qemu-system-x86_64 -m 2048 -enable-kvm test-vm-1.qcow2 -cdrom ./Centos-Desktop-x86_64-20-1.iso`:
    - -m 指定虚拟机内存大小，默认单位是 MB
    - -enable-kvm 使用 KVM 进行加速
    - -cdrom 添加 fedora 的安装镜像

    iso安装完成后重起虚拟机便会从硬盘 ( test-vm-1.qcow2 ) 启动.
- qemu-img：创建虚拟机镜像文件的工具(`apt install qemu-utils`)

    虚拟机镜像用来模拟虚拟机的硬盘，在启动虚拟机之前需要创建镜像文件. 镜像文件创建完成后，可使用 qemu-system-x86 来启动x86 架构的虚拟机.

    `qemu-img create -f qcow2 test-vm-1.qcow2 10G`:
    - -f 选项用于指定镜像的格式，qcow2 格式是 Qemu 最常用的镜像格式，采用来写时复制技术来优化性能
    - -o options : 功能选项, 可用`qemu-img create -f qcow2 -o?`查询qcow2支持的option

        - backing_file=xxx : 基于img xxx创建镜像, 即实现增量镜像的效果, 参考[KVM虚拟机镜像那点儿事，qcow2六大功能，内部快照和外部快照有啥区别？](https://sq.sf.163.com/blog/article/218146701477384192)

            backing file就是基于这个原理的用处，一个qcow2的image可以保存另一个disk image的改变，而不影响另一个image.
    - test-vm-1.qcow2 是镜像文件的名字
    - 10G是镜像文件大小

    > QED 的开发已被放弃, 建议使用qcow2.

    `qemu-img convert -p -c -f raw -O qcow2 vm500G.raw /path/new-vm500G.qcow2`, vm500G.raw其实也可以是/dev/sda:
    - -p : 显示转换进度

    `qemu-img check [-f fmt] filename`: 对磁盘镜像文件进行一致性检查， 查找镜像文件中的错误.
    - -f fmt : 指定文件的格式. 如果不指定格式, qemu-img会自动检测.

    `qemu-img commit [-f fmt] filename`: 提交filename文件中的更改到后端支持镜像文件（创建时通过backing_file指定的）中

    `qemu-img convert [-c] [-f fmt] [-O output_fmt] [-o options] filename[filename2[...]] output_filename`: 将fmt格式的filename镜像文件根据options选项转换为格式为output_fmt的, 名为output_filename的镜像文件. 比如将vmware vmdk转为qcow2. 一般来说， 输入文件格式fmt由qemu-img工具自动检测到， 而输出文件格式output_fmt根据自己需要来指定,  默认会被转换为raw文件格式（且默认使用稀疏文件的方式存储, 以节省存储空间）

    如果使用qcow2、 qcow等作为输出文件格式来转换raw格式的镜像文件（非稀疏文件格式）, 镜像转换还可以将镜像文件转化为更小的镜像，因为它可以将空的扇区删除, 使之在生成的输出文件中不存在.

    - -c : 表示对输出img进行压缩, 但仅qcow2和qcow支持压缩, 并且这种压缩是只读的， 如果压缩的扇区被重写， 则会被重写为未压缩的数据.

    `qemu-img info [-f fmt] filename`: 查看filename的信息. 如果是稀疏文件则还会显示预计分配大小和实际分配大小; 如果由snapshot则也会显示出来

    `qemu-img snapshot [-l|-a snapshot|-c snapshot|-d snapshot] filename`.
    - -l : 表示查询并列出镜像文件中的所有快照
    - -a snapshot : 让镜像文件使用某个快照
    - -c snapshot : 创建一个快照
    - -d : 删除一个快照

    `qemu-img rebase [-f fmt] [-t cache] [-p] [-u] -b backing_file [-F backing_fmt] filename`: 改变镜像文件的后端镜像文件,  只有qcow2和qed格式支持rebase命令

    这个命令可以工作于两种模式之下， 一种是安全模式（Safe Mode） ， 这是默认的模式， qemu-img会根据比较原来的后端镜像与现在的后端镜像的不同进行合理的处理； 另一种是非安全模式（Unsafe Mode） ， 是通过“-u”参数来指定的， 这种模式主要用于将后端镜像重命名或移动位置后对前端镜像文件的修复处理， 由用户去保证后端镜像的一致性
    - -b backing_file : 指定的文件作为后端镜像
    - -F backing_fmt : 后端镜像被转化为指定的backing_fmt后端镜像格式

    `qemu-img resize filename [+|-]size`

    改变镜像文件的大小， 使其不同于创建之时的大小. “+”和“-”分别表示增加和减少镜像文件的大小， size也支持K、 M、 G、 T等单位的使用. 缩小镜像的大小之前， 需要在客户机中保证其中的文件系统有空余空间， 否则数据会丢失。 另外， qcow2格式文件不支持缩小镜像的操作.

    在增加了镜像文件大小后， 也需启动客户机在其中应用“fdisk”“parted”等分区工具进行相应的操作, 才能真正让客户机使用到增加后的镜像空间。 不过使用resize命令时需要小心（做好备份） ， 如果失败， 可能会导致镜像文件无法正常使用， 而造成数据丢失.


    qemu-img支持的format可用`qemu-img -h|grep 'Supported formats'`查询:
    - raw
    
        原始的磁盘镜像格式， 也是qemu-img命令默认的文件格式. 这种格式的文件的优势在于它非常简单， 且非常容易移植到其他模拟器（emulator， QEMU也是一个emulator）上去使用. 如果客户机文件系统（如Linux的ext2/ext3/ext4、 Windows的NTFS） 支持“空洞”（hole） ， 那么镜像文件只有在被写有数据的扇区才会真正占用磁盘空间， 从而节省磁盘空间. qemu-img默认的raw格式的文件其实是稀疏文件（sparse file）.
        
        raw格式只有一个参数选项： preallocation, 它有3个值： off， falloc， full:
        - off就是禁止预分配空间， 即采用稀疏文件方式， 这是默认值
        - falloc是qemu-img创建镜像时候调用posix_fallocate()函数来预分配磁盘空间给镜像文件（但不往其中写入数据， 所以也能瞬时完成）.
        - full是除了实实在在地预分配空间以外， 还逐字节地写0， 所以很慢. 尽管一开始就实际占用磁盘空间的方式没有节省磁盘的效果， 不过这种方式在写入新的数据时不需要宿主机从现有磁盘中分配空间， 因此在第一次写入数据时， 这种方式的性能会比稀疏文件的方式更好一点

        > [稀疏文件](https://en.wikipedia.org/wiki/Sparse_file)是计算机系统块设备中能有效利用磁盘空间的文件类型， 它用元数据（metadata） 中的简要描述来标识哪些块是空的， 只有在空间被实际数据占用时， 才将数据实际写到磁盘中.
    - qcow2

        qcow2是QEMU目前推荐的镜像格式， 它是使用最广、 功能最多的格式。 它支持稀疏文件（即支持空洞） 以节省存储空间， 它支持可选的AES加密以提高镜像文件安全性， 支持基于zlib的压缩， 支持在一个镜像文件中有多个虚拟机快照

        在qemu-img命令中qcow2支持如下几个选项:
        - compat（兼容性水平， compatibility level） ， 可以等于0.10或者1.1， 表示适用于0.10版本以后的QEMU， 或者是1.1版本以后的QEMU
        - backing_file : 用于指定后端镜像文件
        - backing_fmt : 设置后端镜像的镜像格式
        - cluster_size : 设置镜像中簇的大小， 取值为512B～2MB， 默认值为64kB。 较小的簇可以节省镜像文件的空间， 而较大的簇可以带来更好的性能， 需要根据实际情况来平衡。一般采用默认值即可
        - preallocation， 设置镜像文件空间的预分配模式， 其值可为off、 falloc、 full、metadata. 前3种与raw格式的类似， metadata模式用于设置为镜像文件预分配metadata的磁盘空间， 所以这种方式生成的镜像文件稍大一点， 不过在其真正分配空间写入数据时效率更高. 生成镜像文件的大小依次是`off<metadata<falloc=full`， 性能上full最好， 其他3种依次递减.
        - encryption， 用于设置加密， 该选项将来会被废弃， 不推荐使用. 对于需要加密镜像的需求， 推荐使用Linux本身的Linux dm-crypt/LUKS系统
        - lazy_refcounts， 用于延迟引用计数（refcount） 的更新， 可以减少metadata的I/O操作， 以达到提高performance的效果.  适用于cache=writethrough这类不会自己组合metadata操作的情况. 它的缺点是一旦客户机意外崩溃， 下次启动时会隐含一次qemu-img check-rall的操作， 需要额外花费点时间. 它是当compact=1.1时才有的选项.
        - refcount_bits， 一个引用计数的比特宽度， 默认为16
    - qcow

        这是较旧的QEMU镜像格式， 现在已经很少使用了， 一般用于兼容比较老版本的QEMU, 推荐使用qcow2.
    - vdi
    
        兼容Oracle（Sun） VirtualBox1.1的镜像文件格式（Virtual Disk Image）
    - vmdk

        兼容VMware 4版本以上的镜像文件格式（Virtual Machine Disk Format）
    - vpc

        兼容Microsoft的Virtual PC的镜像文件格式（Virtual Hard Disk format）
    - vhdx

        兼容Microsoft Hyper-V的镜像文件格式
    - sheepdog

        heepdog项目是由日本NTT实验室发起的、 为QEMU/KVM做的一个开源的分布式存储系统， 为KVM虚拟化提供块存储.

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
- qemu-pr-helper : [Persistent reservation helper protocol](https://www.qemu.org/docs/master/interop/pr-helper.html)

如何更便捷地创建虚拟机呢? 答案就是libvirt. 这次就不再一个个指定虚拟机启动的参数，而是用 libvirt. 它管理 qemu 虚拟机，是基于 XML 文件，这样容易维护. xml中定义了kvm中domain的配置信息，可以使用virt-install来生成. 首先，需要安装 libvirt
```bash
# apt install libvirt-bin # 高版本上是libvirt-clients
# apt install virtinst
```

> kvm中vm的grub kernel启动参数添加`console=ttyS0`再执行update-grub并重启后, 就可以通过`virsh console ${vm_name}`，进入机器的控制台，可以不依赖于 SSH 和 IP 地址进行登录.

## qemu-system-x86_64
参考:
- [qemu-system-x86_64命令总结](http://blog.leanote.com/post/7wlnk13/%E5%88%9B%E5%BB%BAKVM%E8%99%9A%E6%8B%9F%E6%9C%BA)

### 选项
参考:
- [qemu命令行参数](https://blog.csdn.net/GerryLee93/article/details/106475710)

- boot :

    - once=d : 指定系统的启动顺序是首次光驱, guest reboot后根据默认order(启动顺序)启动
- cdrom : 分配给guest的光驱
- display

    - curses : 仅用于文本模式(text mode is used only with BIOS firmware), 当出现"640 x 480 Graphic mode"时表示guest已切换到图形模式.
    - sdl : qemu console gui
- -cpu <cpu>/help : help可获取qemu支持模拟的cpu
- -kernel bzImage : 使用linux bzImage, 但`CONFIG_PVH=y`时可直接使用vmlinux
- -m <N>G : 分配内存大小
- -M <machine>/help : 当前版本的Qemu工具支持的开发板列表
- -s : 设置gdbserver的监听端口, 等同于`-gdb tcp::1234`
- -S : 启动时cpu仅加电, 但不继续执行, 相当于将断点打在CPU加电后要执行的第一条指令处，也就是BIOS程序的第一条指令. 必须在qemu monitor输入`c`才能继续. 未使用`-monitor`时, 按`Ctrl+Alt+2`可进入qemu的monitor界面,`Ctrl+Alt+1`回到qemu
- -serial stdio : redirects the virtual serial port to the host's terminal input/output, 丢失early boot信息即加电到出现终端登入界面间的信息.
- smp <N> : 为对称多处理器结构分配N个vCPU
- -device : 要模拟的设备, help可获取qemu支持模拟的设备列表
- -monitor

    tcp – raw tcp sockets, **推荐**.
    telnet – the telnet protocol is used instead of raw tcp sockets. This is the preferred option over tcp as you can break out of the monitor using Ctrl-] then typing quit. You can’t break out of the monitor like this after connecting with the raw socket option
    10.1.77.82 – Listen on this host/IP only. You can use 127.0.0.1 if you want to only allow connections locally. If you want to listen on any ip address on the server, just leave this blank so you end up with two consecutive colons ie `::`.
    4444 – port number to listen on.
    server – listening in server mode
    nowait – qemu will wait for a client socket application to connect to the port before continuing unless this option is used. In most cases you’ll want to use the nowait option.

qemu默认会启动一个vnc server(port 5900), 可用vncviewer工具连上该server来查看guest.

> centos 7常使用tigervnc-server和tigervnc作为vnc server和vncviewer.

## qemu monitor
进入: 鼠标点击qemu窗口，然后ctrl+alt+2即可切换到控制台; ctrl+alt+1回到guest窗口.
滚屏: ctrl + 上/下
查看是否使用kvm: info kvm

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

## noVNC
ref:
- [NoVNC远程连接](https://www.jianshu.com/p/0f3b351a156c)

noVNC是一个 HTML5 VNC 客户端，采用 HTML 5 WebSockets, Canvas 和 JavaScript 实现，noVNC 被普遍用在各大云计算、虚拟机控制面板中，比如 OpenStack Dashboard 和 OpenNebula Sunstone 都用的是 noVNC.

noVNC采用WebSockets实现，但是目前大多数VNC服务器都不支持 WebSockets，所以noVNC是不能直接连接 VNC 服务器的，需要一个代理来做WebSockets和TCP sockets之间的转换. 这个代理在noVNC的目录里，叫做websockify.

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

### machine type选择
参考:
- [Platforms available in QEMU](https://wiki.qemu.org/Documentation/Platforms)

#### x86_64
参考:
- [Qemu X86架构的Machine Type](https://remimin.github.io/2019/07/09/qemu_machine_type/)

可通过`qemu-system-x86_64 --machine help`查看x86支持的所以machine type.

i440fx是1996年推出的架构, 已过时. q35是2009年推出的架构, 更现代.

#### riscv64
参考:
- [RISC-V QEMU Part 2: The RISC-V QEMU port is upstream](https://www.sifive.com/blog/risc-v-qemu-part-2-the-risc-v-qemu-port-is-upstream)

- Spike是官方的RISC-V模拟平台
- SiFive的E和U系列: E, 32-bit embedded cores; U, 64-bit application processors
- Virt是通用的虚拟化RISC-V平台，支持VirtIO设备

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
# egrep '^flags.*(vmx|svm)' /proc/cpuinfo # vmx is intel; svm is amd
# LC_ALL=C lscpu | grep Virtualization
# virt-host-validate # 该工具验证主机是否以合适的方式配置来运行 libvirt 管理程序驱动程序
```

硬件不支持, 检查bios/uefi是否关闭了虚拟化支持, 此时可用kvm-ok(ubuntu/debian from cpu-checker)检查.

arm检查:
```bash
# dmesg | grep -i kvm # 也发现过不输出`Hyp mode initialized successfully`的机器存在/dev/kvm
KVM [1]: Hyp mode initialized successfully
```

### Dmesg print "psmouse serio1: VMMouse at isa0060/serio1/input0 lost sync at byte 1" when using vmmouse
一旦鼠标移动到qemu gui上, 就会提示该信息. qemu 4.2正常, 5.1.0有该问题.

```bash
# rmmod psmouse # 卸载psmouse驱动, psmouse是触控板驱动.
```

### `qemu-system-x86_64 -kernel arch/x86_64/boot/bzImage -nographic` 退出
尝试按Ctrl + a，然后按c，以获得（qemu console）提示, 再一个简单的q即可退出.

推荐使用`qemu-system-x86_64 -kernel arch/x86_64/boot/bzImage -nographic -append "console=ttyS0"`运行, 退出方法同上.

### 强制退出QEMU虚拟机

先按ctrl + a 放开后，再按下 x, 这在运行启动虚拟机的命令后发生卡死现象时特别有用.

### 压缩qcow2
初始磁盘压缩方法:
```bash
//将默认raw格式的磁盘，简单压缩转换成qcow2格式
# cp --sparse=always vm500G.raw vm500G-new.raw   //--sparse=always稀疏拷贝，忽略全0数据
# qemu-img convert -c -f raw -O qcow2 vm500G.raw /path/new-vm500G.qcow2
 
//将默认qcow2格式的磁盘，导出为简单压缩后的qcow2格式
# qemu-img convert -p -c -O qcow2 vm500G.qcow2 new.img.qcow2
```

以上两种方法都能在一定程度上压缩减小导出后的镜像文件体积；但仅限于在虚拟机刚安装部署好，还没有进行过大量数据读写处理的情况下. 因为非空白(非全零)块无法压缩.

针对其他创建需要先用ncdu清理, 再vm内部写零操作:
```bash
# dd if=/dev/zero of=/null.dat   //创建一个全0的大文件，占满所有的剩余磁盘空间，需要很久时间
# rm -f /null.dat                //删除这个文件
```
最后关机并使用`初始磁盘压缩`方法即可.

### 获得更多输出
qemu窗口-视图-勾选"显示标签页"

### 查询qemu vm的默认选项
`virt-install --os-variant <os>`, 参数os是`osinfo-query os`输出的"Short ID".

### [windows virtio驱动](https://fedorapeople.org/groups/virt/virtio-win/direct-downloads/archive-virtio/virtio-win-0.1.96/)

### UEFI启动 花屏
UEFI启动需要使用qxl驱动, 即`-vga qxl`

### 桌面虚拟化
ref:
- [华为桌面云解决方案](https://blog.51cto.com/u_15162069/2768919)

SPICE (Simple Protocol for Independent Computing Environments) 是一个用于虚拟化环境中的通讯协定, 此协议透过网络来连结到虚拟化平台上之虚拟机器桌面. SPICE client的实现有Virt-viewer.

相较于通过浏览器连接(HTML5, 比如noVNC) 或其他VNC client, SPICE不仅支持虚拟机器音源输出且拥有较佳的影像显示.

### QEMU 支持直接引导Linux内核(vmlinuz，initrd，bzImage)，非常方便适用于内核调试
```bash
qemu-system-x86_64 -kernel bzImage -initrd initrd.img \
                  -append "root=/dev/sda1 init=/bin/bash" \
                  -hda rootfs.img
```
选项:
 - -kernel : 提供内核镜像,bzImage
 - -initrd : 提供initramfs
 - -append : 提供内核参数，指引rootfs所在分区，指引init命令路径
 - -hda : 提供rootfs根文件系统

区别于下列方式，其中rootfs.img包括了 grub + MBR + 多磁盘分区 + rootfs ,启动过程基本与传统PC启动过程无异.

qemu-system-x86_64 -m 512M -drive format=raw,file=rootfs.img

该种引导方式，内核模块和initramfs都包括在rootfs.img中的某个分区的文件系统中

### aarch64启动方式
[在 AArch64 架构上启动 VM 有两种方式](https://that.guru/blog/uefi-secure-boot-in-libvirt/)：
1. UEFI
2. 内核+initrd

一般除了调试没人会选择第二种, 因此它可等同于仅支持uefi.

libvirt 5.3 引入了对 QEMU 自 QEMU 2.9 提供的固件自动选择功能的支持. 此 QEMU 功能依赖于固件 JSON 文件(`/usr/share/qemu/firmware` from `edk2-ovmf`), 这些文件描述了每个固件文件的用途以及如何描述它.

`virsh domcapabilities --machine pc-q35-5.1 | xmllint --xpath '/domainCapabilities/os' -`看查看machine支持的uefi固件.

### hostdev与direct区别
hostdev为vm提供对 PCI 设备的直接访问, 缺点是无法模拟不同的设备类型. 因此, vm必须有一个可用的驱动程序，该驱动程序与管理程序提供的硬件类型相匹配.

direct并部署真的direct, 它依赖于管理程序配置的网络接口来提供与客户操作系统的连接. 它允许 KVM 本地模拟一些常见的网络接口类型，这些类型通常存在于大多数当前和旧版操作系统中.