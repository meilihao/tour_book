# guestfish
guestfish是一套虚拟机镜像管理的利器, 提供一系列对镜像管理的工具, 也提供对外的API.

guestfish主要包含以下工具：
- guestfish : 挂载镜像，并提供一个交互的shell

	guestfish --rw --add disk.img --mount /dev/vg_guest/lv_root --mount /dev/sda1:/boot edit /boot/grub/grub.conf # 编辑grub.conf	
- guestmount : 将镜像挂载到指定的目录

	guestmount -a t.img -m /dev/sda1 --ro /mnt # 以只读方式将镜像挂到/mnt	
- guestumount : 卸载镜像目录
- virt-alignment-scan : 镜像块对齐扫描, 检查镜像分区是否块对齐
- virt-builder : 快速镜像创建
- virt-cat : 显示镜像中文件内容

	virt-cat -a t.qcow2 /root/anaconda-ks.cfg
- virt-copy-in : 拷贝文件到镜像内部

	virt-copy-in test.txt -a t.qcow2 /root	
- virt-copy-out : 拷贝镜像文件出来

	virt-copy-in -a t.qcow2 /root/test.txt .
- virt-customize : 定制虚拟机镜像
- virt-df :查看虚拟机镜像空间使用情况

	virt-df -a t.qcow2
- virt-diff : 不启动虚拟机的情况下，比较虚拟机内部两份文件差别
- virt-edit : 编辑虚拟机内部文件

	virt-edit -a t.qcow2 /root/anaconda-ks.cfg
- virt-filesystems : 显示镜像文件系统信息
- virt-format : 格式化镜像内部磁盘
- virt-inspector : 镜像信息测试

	显示os版本, kernel, driver, mountpoint等

	virt-inspector [--query] xxx.qcow2, `--query`输出便于解析

	virt-inspector2支持xml输出, 系统信息比virt-inspector详细
- virt-list-filesystems : 列出镜像文件系统

	列出vm镜像的fs, 分区, 块设备, lvm信息
- virt-list-partitions : 列出镜像分区信息
- virt-log : 显示镜像日志
- virt-ls : 列出镜像文件

	virt-ls -a t.qcow2 /root
- virt-make-fs : 镜像中创建文件系统
- virt-p2v :物理机转虚拟机
- virt-p2v-make-disk : 创建物理机转虚拟机ISO光盘
- virt-p2v-make-kickstart : 创建物理机转虚拟机kickstart文件
- virt-rescue:进去虚拟机救援模式
- virt-resize : 虚拟机分区大小修改

	virt-resize --expand /dev/sda2 olddisk newdisk
	virt-resize --resize /dev/sda1=+200M --expand /dev/sda2 olddisk newdisk # 将boot分区增加200M, 剩下空间给/dev/sda2
	virt-resize --expand /dev/sda2 --LV-expand /dev/vg_guest/lv_root olddisk newdisk # lv扩展
- virt-convert : 转换vm镜像格式

	virt-convert -i raw -o qcow2 old.img new.qcow2
- virt-sparsify : 镜像稀疏空洞消除

	virt-sparsify -x test.qcow2 --convert qcow2 test2.qcow2	
- virt-sysprep :镜像初始化
- virt-tar : 文件打包并传入传出镜像

	virt-tar -x domname /home home.tar # 将vm home目录复制出来并打包
	virt-tar -u domname xxx.tar /tmp # 上传压缩文件到vm并解压
- virt-tar-in : 文件打包并传入镜像

	virt-tar-in -a t.qcow2 data.tar /root	
- virt-tar-out : 文件打包并传出镜像

	virt-tar-out -a t.qcow2 /root data.tar
- virt-v2v :其他格式虚拟机镜像转KVM镜像
- virt-win-reg : 注册表导入镜像
- libguestfs-test-tool :测试libguestfs
- hivex : 解压windows注册表文件
- hivexregedit :合并、并导出注册表文件内容
- hivexsh:注册表修改交互的shell
- hivexml:注册表转化为xml
- hivexget得到注册表键值
- guestfsd:guestfs服务

安装: `yum install libguestfs-tools libguestfs-winsupport`, 默认安装不支持windows.

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

## FAQ
### windows 2008 server自检慢
方法1:

```bash
yum install ntfs-3g
losetup -f # 获取可用的loop设备
/dev/loop0
losetup /dev/loop0 /dev/vmVG/ptyyb-webzb-57_vda # 挂载镜像
kpartx -av /dev/loop0 # 使用kpartx将镜像分区映射: 如果只有一个分区, 默认挂载在/dev/mapper/loop0p1; 有boot分区时, C盘在/dev/mapper/loop0p2
ntfsfix -b -d /dev/mapper/loop0p1 # `-b -d`清除ntfs的检查标志信息
kpartx -dv /dev/loop0 # 分离镜像映射
losetup -d /dev/loop0
```

2. guestfish