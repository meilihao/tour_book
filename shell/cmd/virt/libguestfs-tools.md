# libguestfs-tools
libguestfs是用于访问和修改虚拟机的磁盘镜像的一组工具集合. 它提供了访问和编辑客户机中的文件、 脚本化修改客户机中的信息、 监控磁盘使用和空闲的统计信
息、 P2V、 V2V、 创建客户机、 克隆客户机、 备份磁盘内容、 格式化磁盘、 调整磁盘大小等非常丰富的功能.

libguestfs除了支持KVM虚拟机， 它甚至支持VMware、Hyper-V等非开源的虚拟机. 同时， libguestfs还提供了一套C库以方便被链接到自己用
C/C++开发的管理程序之中. 它还有对其他很多流程编程语言（如： Python） 的绑定. 它的安装方法: `dnf install libguestfs-tools libguestfs-tools-c`

libguestfs-tools提供了很多工具， 可以分别对应不同的功能和使用场景， 如： virt-ls用于列出虚拟机中的文件, virt-copy-in用于往虚拟机中复制文件或目录， virt-copy-out用于从虚拟机往外复制文件或目录， virt-resize用于调整磁盘大小， virt-cat用于显示虚拟机中的一个文件的内容， virt-edit用于编辑虚拟机中的文件， virt-df用于查看虚拟机中文件系统空间使用情况, 等等.

libguestfs的一些工具用于Windows客户机镜像的操作时， 需要先安装libguestfswinsupport这个软件包； 当使用guestmount来挂载Windows镜像时, 还需要安装ntfs-3g软件包.

virt-inspector探测image信息.

## 场景
整机保护(保护系统盘+若干数据盘)中修改fstab, grub, ip等.

## libguestfs原理
ref:
- [libguestfs详解](https://www.hanbaoying.com/2017/02/26/libguestfs.html)

原理:
1. 执行guestfish -a会动一个进程，也就是那个shell壳子，姑且称之为main program
2. 运行run的时候，会创建一个child process，在child process中，利用libvirt启动一个称为appliance的虚拟机。
3. 在appliance中，运行了linux kernel和一系列用户空间的工具(LVM, ext2等)，以及一个后台进程guestfsd
4. main process中的libguestfs和这个guestfd通过RPC进行交互
5. 由child process的kernel来操作disk image

可用`export LIBGUESTFS_DEBUG=1`来查看详细的启动过程.