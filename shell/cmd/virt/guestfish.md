# guestfish
guestfish是一套虚拟机镜像管理的利器, 提供一系列对镜像管理的工具, 也提供对外的API. Guestfish 是libguestfs项目中的一个工具软件，提供修改虚机镜像内部配置的功能, 基于libguestfs.

guestfish主要包含以下工具：
- guestfish : 挂载镜像，并提供一个交互的shell

	guestfish --rw --add disk.img --mount /dev/vg_guest/lv_root --mount /dev/sda1:/boot edit /boot/grub/grub.conf # 编辑grub.conf	
- guestmount : 将镜像挂载到指定的目录

	guestmount -a t.img -m /dev/sda1 --ro /mnt # 以只读方式将镜像挂到/mnt	
- guestumount : 卸载镜像目录
- virt-alignment-scan : 镜像块对齐扫描, 检查镜像分区是否块对齐

	virt-alignment-scan -a xxx.qcow2	
- virt-builder : 快速镜像创建
		
	可轻松快速地构建供本地或云使用的各种虚拟机镜像
- virt-cat : 显示镜像中文件内容

	virt-cat -a t.qcow2 /root/anaconda-ks.cfg
- virt-copy-in : 拷贝文件到镜像内部

	virt-copy-in test.txt -a t.qcow2 /root	
- virt-copy-out : 拷贝镜像文件出来

	virt-copy-in -a t.qcow2 /root/test.txt .
- virt-customize : 定制虚拟机镜像

	- 设置 root 密码并禁用 cloud-init: `virt-customize -a MY-CLOUD-IMAGE.qcow2 --root-password password:SUPER-SECRET-PASSWORD --uninstall cloud-init`
	- `--inject-virtio-win`: 给windows注入virtio驱动
- virt-df :查看虚拟机镜像空间使用情况

	virt-df -a t.qcow2
- virt-diff : 不启动虚拟机的情况下，比较虚拟机内部两份文件差别
- virt-edit : 编辑虚拟机内部文件

	virt-edit -a t.qcow2 /root/anaconda-ks.cfg
- virt-filesystems : 显示镜像文件系统信息

    - [New tool: virt-filesystems](https://github.com/libguestfs/libguestfs/commit/fbc2555903be8c88ad9430d871cf0d27c8fded1e)
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

	virt-rescue --suggest -d fedora15 # 进入救援模式
- virt-resize : 虚拟机分区大小修改

	virt-resize --expand /dev/sda2 olddisk newdisk
	virt-resize --resize /dev/sda1=+200M --expand /dev/sda2 olddisk newdisk # 将boot分区增加200M, 剩下空间给/dev/sda2
	virt-resize --expand /dev/sda2 --LV-expand /dev/vg_guest/lv_root olddisk newdisk # lv扩展
	virt-resize --align-first always --expand /dev/vda1 kuai-no-vda kuai-no-vda-2-yes # 调整镜像并强制块对齐
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
- virt-win-reg : 修改windows硬盘镜像中的注册表

	- [16.10.3. 使用 virt-win-reg](https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/6/html/virtualization_administration_guide/sect-virt-win-reg-use)
	- [获取windows版本: `virt-win-reg Windows7 'HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion' ProductName`](https://libguestfs.org/virt-win-reg.1.html)
- libguestfs-test-tool :测试libguestfs
- hivex : 解压windows注册表文件
- hivexregedit :合并、并导出注册表文件内容
- hivexsh:注册表修改交互的shell
- hivexml:注册表转化为xml
- hivexget得到注册表键值
- guestfsd:guestfs服务

安装: `yum install libguestfs-tools libguestfs-winsupport`, 默认安装不支持windows.

# libguestfs-tools

> libguestfs-test-tool可检查libguestfs-tools环境是否正常, 正常输出`===== TEST FINISHED OK =====`

libguestfs是用于访问和修改虚拟机的磁盘镜像的一组工具集合. 它提供了访问和编辑客户机中的文件、 脚本化修改客户机中的信息、 监控磁盘使用和空闲的统计信
息、 P2V、 V2V、 创建客户机、 克隆客户机、 备份磁盘内容、 格式化磁盘、 调整磁盘大小等非常丰富的功能.

libguestfs除了支持KVM虚拟机， 它甚至支持VMware、Hyper-V等非开源的虚拟机. 同时， libguestfs还提供了一套C库以方便被链接到自己用
C/C++开发的管理程序之中. 它还有对其他很多流程编程语言（如： Python） 的绑定. 它的安装方法: `dnf install libguestfs-tools libguestfs-tools-c`

libguestfs-tools提供了很多工具， 可以分别对应不同的功能和使用场景， 如： virt-ls用于列出虚拟机中的文件, virt-copy-in用于往虚拟机中复制文件或目录， virt-copy-out用于从虚拟机往外复制文件或目录， virt-resize用于调整磁盘大小， virt-cat用于显示虚拟机中的一个文件的内容， virt-edit用于编辑虚拟机中的文件， virt-df用于查看虚拟机中文件系统空间使用情况, 等等.

libguestfs的一些工具用于Windows客户机镜像的操作时， 需要先安装libguestfswinsupport这个软件包； 当使用guestmount来挂载Windows镜像时, 还需要安装ntfs-3g软件包.

virt-inspector探测image信息.

## [guestfish shell](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/virtualization_deployment_and_administration_guide/sect-guest_virtual_machine_disk_access_with_offline_tools-the_guestfish_shell#doc-wrapper)
image权限需要`qemu:qemu`, 且qemu用户能访问到该文件

```bash
# --- 不支持lsblk
# guestfish --rw --add disk.img [-i] # -i: (--inspector) - Inspect the disks and mount the filesystems, 如果执行探测成功会自动执行run
><fs> run
><fs> list-filesystems
/dev/sda1: xfs
/dev/centos/root: xfs
/dev/centos/swap: swap
><fs> pvs
/dev/sda2
><fs> lvs
/dev/centos/root
/dev/centos/swap
><fs> xfs-repair /dev/centos/root [forcelogzero:true] # `forcelogzero:true`=`xfs_repair -L` // 有遇到: 1. 直接挂载没有问题, 强制修复后反而无法挂载; 2.第一次强制修复lvm+root分区失败(在本次重复多次修复均报错), 但退出再次执行guestfish命令重复修复后却能mount
><fs> mount /dev/centos/root /
```

ext4:
```bash
><fs> e2fsck /dev/centos/root [forceall:true] # correct:可能还是会报错
```

```bash
$ --- 支持lsblk
$ virt-rescue --format=raw  -a /dev/zvol123
...
Welcome to virt-rescue, the libguestfs rescue shell.

Note: The contents of / (root) are the rescue appliance.
You have to mount the guest’s partitions under /sysroot
before you can examine them.

><rescue> lsblk       
NAME MAJ:MIN RM SIZE RO TYPE MOUNTPOINT
sda    8:0    0   1G  0 disk 
sdb    8:16   0   4G  0 disk /
><rescue> mount /dev/sda /sysroot/
><rescue> ls /sysroot/
lost+found  test
```

支持的其他命令 from [guestfish - the guest filesystem shell](https://www.libguestfs.org/guestfish.1.html): ls, ll, cat, more, download, tar-out, ...


## 场景
整机保护(保护系统盘+若干数据盘)中修改fstab, grub, ip等.

## libguestfs原理
ref:
- [libguestfs详解](https://www.hanbaoying.com/2017/02/26/libguestfs.html)

原理:
1. 执行guestfish -a会动一个进程，也就是那个shell壳子，姑且称之为main program
2. 运行run的时候，会创建一个child process，在child process中，利用libvirt启动一个称为appliance的虚拟机

	appliance使用supermin和host的kernel制作而成
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

### 拷贝到disk
`guestfish add /dev/xxx : run : mount /dev/sda1 : copy-in /root/a.cfg /root : ls /root`

### xfs_repair: feature 'xfs' is not available in this\nbuild of libguestfs.
`yum install libguestfs-xfs`

### guestfish: error while loading shared libraries: libconfig.so.11: cannot open shared object file: No such file or directory
需要libconfig-1.7.x, 比如当前环境是libconfig-1.7.2

### guestfish: /usr/lib64/libselinux.so.1: no version information available (required by /usr/lib64/libguestfs.so.0)
`yum reinstall libselinux/apt install --reinstall libselinux1`

如果报错信息里还有`libguestfs: error: mount: mount exited with status 32: mount: /sysroot: wrong fs type, bad option, bad superblock on ...`(`mount /xxx /`), 那就是在guestfish中mount fs时因为要挂载的fs需要修复而引发的错误, 需要**逐个检查要挂载的fs**, 修复fs后再mount即可.