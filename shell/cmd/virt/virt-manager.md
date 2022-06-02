# virt-manager
基于libvirt的管理vm的gui工具.

virt-viewer经常用于替换传统的VNC客户端查看器， 因为后者通常不支持x509认证授权的SSL/TLS加密， 而virt-viewer是支持的.

virt-install命令行工具为虚拟客户机的安装提供了一个便捷易用的方式， 它也是用libvirt API来创建KVM、 Xen、 LXC等各种类型的客户机， 同时， 它也为
virt-manager的图形界面创建客户机提供了安装系统的API. 使用virt-install中的一些选项（--initrd-inject、 --extra-args等） 和Kickstart文件， 可以实现无人值守的自动化安装客户机系统.

virt-top是一个用于展示虚拟化客户机运行状态和资源使用率的工具.

## 安装与使用
```bash
# dnf install virt-manager
# virt-manager -c qemu+ssh://192.168.158.31/system
```

## 构建
前提:
1. qemu
1. libvirt
1. 源码中的`INSTALL.md`

```
# -- https://github.com/archlinux/svntogit-community/blob/packages/virt-manager/trunk/PKGBUILD
wget https://releases.pagure.org/virt-manager/virt-manager-2.2.1.tar.gz
tar -xvzf virt-manager-2.2.1.tar.gz
cd virt-manager-2.2.1
apt install intltool
# python3 setup.py configure --default-hvs qemu,lxc
python3 setup.py install [--force]
```

## 源码
入口在`virt-manager`, 入口是`virtmanager.runcli()`

vscode配置:
```bash
# cat ~/.vscode/setting
{
    "python.autoComplete.extraPaths": [
        "/usr/share/virt-manager",
        "/usr/lib/python3/dist-packages"
    ],
    "python.analysis.extraPaths": [
        "/usr/share/virt-manager",
        "/usr/lib/python3/dist-packages"
    ]
}
```

logger打印文件和行数, virt-manger使用了`virtinst`的`logger.py`, 修改方法即在`log = logging.getLogger("virtinst")`前加入:
```python
logging.basicConfig(format='%(asctime)s.%(msecs)03d [%(levelname)s] [%(pathname)s:%(lineno)d] %(message)s', datefmt='%Y-%m-%d %H:%M:%S')
```

获取日志`virt-manager --debug 2>&1 |tee -a vm.log`

通过日志, virt-manager创建vm最后一步的`Finish`按钮调用的是`createvm.py#_finish_clicked()`, 它通过`installer.set_install_defaults(guest)`+`installer.start_install(guest, meter=meter)`来创建guest xml.

具体生成xml逻辑在`virtinst/install/install.py#start_install()`:
1. 先通过`self._build_xml()`生成xml骨架, 比如没有具体的pci controller设备的vm xml.
1. 通过`self._create_guest()`将骨架xml提交给libvirtd, 由它填充细节并返回完整xml.

其他工具生成vm xml的方法类似上述的`start_install`, 先骨架后由libvirt补全:
- webvirtcloud

    调用入口在`instances/views.py`的`create_instance()`

    拼接骨架xml方法在`vrtManager/create.py#create_instance()`, 在其中最后是通过调用lib `libvirt`的`defineXML()`生成具体xml.
- truenas

    入口`middlewared/plugins/vm/vms.py#do_create()`在`await self.middleware.run_in_thread(self._add, vm_id)`, 调用链是`_add() -> _add_with_vm_data() -> VMSupervisor.__init__ -> update_domain()->__define_domain()->construct_xml()->get_domain_children()`, 具体拼接骨架xml和调用`libvirt`的`defineXML()`都在`__define_domain`里.

### 部分逻辑说明
- virtManager/vmmAddHardware.py : "添加硬件"按钮

## FAQ
参考:
- [How to compile virt-manager on Debian or Ubuntu](https://www.xmodulo.com/compile-virt-manager-debian-ubuntu.html)

### unable to execute 'intltool-update': No such file or directory
`apt install intltool`

### virt-manager运行报`Namespace LibvirtGlib not available`
参考gobject-introspection的[Namespaces are represented on disk by type libraries (.typelib files).](https://valadoc.org/gobject-introspection-1.0/GI.Repository.html), 应该是没有对应的`.typelib`文件, 它原本应由`apt install gir1.2-libvirt-glib-1.0`提供.

> GObject Introspection（简称 GI）用于产生与解析 C 程序库 API 元信息，以便于动态语言（或托管语言）绑定基于 C + GObject 的程序库, 具体可见[GObject Introspection 的作用与意义](http://garfileo.is-programmer.com/2012/2/20/gobject-introspection-introduction.32200.html).

> Typelibs将从环境变量GI_TYPELIB_PATH和`/usr/lib/girepository-1. 0/`中的路径加载.

安装libvirt-glib, 参考[libvirt-glib/trunk/PKGBUILD](https://github.com/archlinux/svntogit-community/blob/packages/libvirt-glib/trunk/PKGBUILD)
```bash
# 参考源码的INSTALL
apt install python3-gi libgirepository1.0-dev valac
wget https://libvirt.org/sources/glib/libvirt-glib-1.0.0.tar.gz
tar -xf libvirt-glib-1.0.0.tar.gz
cd libvirt-glib-1.0.0
./configure --enable-introspection --enable-vala
make
make install
cp /usr/local/lib/girepository-1.0/* /usr/lib/aarch64-linux/gnu/girepository-1.0 # 或设置变量GI_TYPELIB_PATH追加/usr/local/lib/girepository-1.0
```

### virt-manager运行报`No module named 'libvirt'`
未安装libvirt的python绑定: python3-libvirt. 参考[libvirt-python/trunk/PKGBUILD](https://github.com/archlinux/svntogit-community/blob/packages/libvirt-python/trunk/PKGBUILD), 安装即可.

```bash
wget https://libvirt.org/sources/python/libvirt-python-6.0.0.tar.gz
tar -xf libvirt-python-6.0.0.tar.gz
cd libvirt-python-6.0.0
python3 setup.py clean
python3 setup.py install --optimize=1
```

### virt-manager运行报`No module named 'libxml2'`
ubuntu 20.04: `apt install python3-libxml2`
ubuntu 16.04没有python3-libxml2, 用`pip3 install libxml2-python3`

### virt-manager运行报`pygobject3 3.22.0 or later is required.`
`pip3 install PyGobject==3.36.1` # version from Ubuntu 20.04

可能会遇到`No package 'cairo' found`, 解决方法: `apt install libcairo2-dev`, 再执行`pip3 install pycairo==1.16.2`

### virt-manager运行报`gtk 3.22.0 or later is required.`
它从`2.1.0`开始gtk必须是`3.24`及以上, 将virt-manager降级到`2.0.0`, gtk只要`3.14`(Ubuntu 16.04.6使用gtk 3.18).

### virt-manager运行报`cannot import name Vte, introspection typelib not found`
`apt install gir1.2-vte-2.91`

### virt-manager无法新建vm, `virt-manager --debug`报`cannot import name 'vmmDetails'`
经核对, vmmDetails明显存在于`/usr/share/virt-manager/virtManager/details.py`中, 估计是上次使用了2.2.1安装, 为了解决gtk报错使用`python3 setup.py install --force`降级安装了2.0.0, 因为历史文件干扰导致, 使用`rm -rf /usr/share/virt-manager`删除再用`python3 setup.py install --force`安装即可.

### virt-manager调试
`virt-manager --debug`

### virt-manager打开新建虚拟机界面报错"Error: No active connection to install on"
调试日志报"Autostart connect error: Unable to connect to libvirt qemu:///system."

原因未知.

> 当前环境是虚拟机, kvm-ok验证/dev/kvm不存在.

### virt-manager新建连接报错: `Cannot recv data: ssh_askpass: exec(/usr/bin/ssh-askpass): No such file or directory`
在virt-manager所在机器执行`apt install ssh-askpass`

### 动态迁移
ref:
- [<<KVM实战>> 4.3.3中的6.动态迁移]

在KVM虚拟环境中， 如果遇到宿主机负载过高或需要升级宿主机硬件等需求时， 可以选择将部分或全部客户机动态迁移到其他的宿主机上继续运行. 需要满足如下前提条件才能使动态迁移成功实施:
1. 源宿主机和目的宿主机使用共享存储, 如NFS、 iSCSI、 基于光纤通道的LUN、GFS2等， 而且它们挂载共享存储到本地的挂载路径需要完全一致，被迁移的客户机就是使用该共享存储上的镜像文件
1. 硬件平台和libvirt软件的版本要尽可能的一致， 如果软硬件平台差异较大， 可能会增加动态迁移失败的概率
1. 源宿主机和目的宿主机的网络通畅并且打开了对应的端口
1. 源宿主机和目的宿主机必须有相同的网络配置， 否则可能出现动态迁移之后客户机的网络不能正常工作的情况
1. 如果客户机使用了和源宿主机建立桥接的方式获得网络， 那么只能在同一个局域网（LAN） 中进行迁移， 否则客户机在迁移后， 其网络将无法正常工作

### virt-manager创建vm后连接到图形控制台时报`Error opening SPICE console: Namespace SpiceClientGtk not available`
`apt install gir1.2-spiceclientgtk-3.0`

### virt-install创建vm报`不能为架构 'x86_64' 找到任何 UEFI 二进制路径`
`apt install ovmf`
`yum install edk2-ovmf`

### 创建vm时如何选择firmware
在创建vm的`Create a new virtual machine`选择`Customize configuration before install`, 它允许在install前配置vm更多细节.

### 启动vm报`child reported (status=125): unable to open /mnt/nfs/xxx.iso: Read-only file system`
/mnt/nfs是nfs的挂载点.

libvirt的限制, 将iso拷贝到本地, 再使用其本地路径即可.

### 创建vm报`'spicevmc' is not a valid cha driver name`
是[qemu报错](https://gitlab.com/qemu-project/qemu/-/issues/488), qemu未构建spicevmc设备

### win10 vm启用virtio
vm配置:
- cdrom: sata
- disk/nic: virtio


**安装系统时加载virtio disk驱动可直接在virtio盘上安装os, 进入系统后可用iso里的`virtio-win-gt-x<arch>.exe(最新版仅支持win8/win server 2012及以上)`可一次性安装全部剩余virtio驱动**

> iso文件列表在`data/info.json`, 里面还能找到xp相关的virtio驱动.

> `virtio-win-gt-x<arch>.exe`从0.1.173-2开始出现且支持xp, 但安装其中的RHEV Agent时出现了失败, 选择不安装其中的`RHEV/SPICE Agent`(两个都要不选)其余的都能正常安装.

> virt-manager可在vm启动的情况下, 更换cdrom里的iso

virtio disk启用步骤:
1. 初始`windows安装程序`的磁盘列表为空, 选择下方的`加载驱动程序`->`浏览`
1. 选择第二个cdrom上的virtio-win-0.1.127.iso->选中`amd/w10`目录->选中要安装的驱动`Red Hat VirtIO SCSI controller (E:\amd64\w10\viostor.inf)`
1. 点`下一步`, 等待完成后, 磁盘列表会出现virtio磁盘.

> 设备管理器标志: 磁盘驱动器里有`Red Hat VirtIO SCSI Disk Device`, 存储控制器里出现`Red Hat VirtIO SCSI controller`

virtio-win-0.1.127.iso内容:
- amd64: scsi和scsi blk驱动的合集
- NetKVM : Virtio网络驱动
- viostor : Virtio块驱动
- vioscsi : Virtio SCSI驱动
- vioserial : virtio 串口驱动
- viorng : Virtio RNG驱动
- Balloon : Virtio 内存气球驱动
- qxl : 用于Windows 7及之前版本的QXL显卡驱动. (virtio-win-0.1.103-1和之后版本会创建)
- qxldod : 用于Windows 8及之后版本的QXL显卡驱动. (virtio-win-0.1.103-2和之后版本会创建)
- pvpanic : QEMU pvpanic 设备驱动 (virtio-win-0.1.103-2和之后版本会创建)
- guest-agent : QEMU Guest Agent 32bit 和 64bit 安装包
- qemupciserial : QEMU PCI 串口设备驱动
- qemufwcfg : qemu fwcfg设备驱动

virtio nic启用:
1. `设备管理器->其他设备->以太网控制器`, 右键选中`更新驱动程序`,选中`E:\NetKVM\w10\amd64`, 下一步即可
1. 驱动安装完成后, `设备管理器->网络适配器`里会出现`Red Hat VirtIO Ethernet Adapter`

其他virtio驱动对应:
1. `设备管理器->其他设备->PCI简单通讯控制器` -> `E:\vioserial\w10\amd64`

### virt-manager安装windows蓝屏
os: xp

蓝屏:
- 0X000000A5

    是ACPI的问题. 蓝屏原因是你电脑BIOS采用的ACPI规范版本太高了(qemu machine是q35), 而Windows XP系统内置的ACPI驱动最高仅能支持到ACPI 2.0规范, 将machine换为pc即可.
- 0X0000007B

    因为配置了sata/virtio盘, xp默认没有相关驱动, 需使用ide. 此时nic和显卡是virtio无影响.

    > xp刚启动按f6支持从**软盘**`--disk path=/tmp/virtio-win-0.1.96_x86.vfd,device=floppy`加载第三方驱动, 比如virtio-win-0.1.96_x86.vfd已包含virtio disk驱动, [从`0.1.196`开始vfd消失](https://fedorapeople.org/groups/virt/virtio-win/direct-downloads/archive-virtio/virtio-win-0.1.196-1/), 因此推荐先使用sata安装os后再安装virtio驱动并关机, 最后修改disk bus即可.

    xp启用virtio方法:
    1. 用ide安装os成功后关机
    2. 添加virtio盘, 开机并用iso安装相关设备驱动
    3. 在`计算机管理->磁盘管理`将新virtio盘初始化并新建分区成功表示virtio blk加载成功, 再关机
    4. 将原先ide盘配置成virtio盘即可(即修改`vm xml`后再`virsh define xxx.xml`).

    经验证: xp安装os成功后(已安装virio blk)的系统盘不能从ide转为virtio, 推测是os启动时没有加载virtio驱动(可能需要类似update-grub2的操作), 但新增virtio disk是可以的.

    > win查看系统信息`systeminfo`, 在存在virtio设备且未安装virtio驱动时该命令可能闪退(在xp上遇到),

    设备驱动映射(`设备管理器->其他设备`):
    1. 以太网控制器: e:\NetKVM\xp\86
    1. scsi控制器 : viostor(xp只支持virtio blk), e:\viostor\xp\86
    1. PCI简单通讯控制器 : E:\vioserial\xp\x86`

    **用`硬件更新向导`的`不要搜索`项(搜索即使指定正确位置也没有结果), 并强制安装驱动**

### kernel是否支持virtio
ref:
- [安装virtio驱动](https://help.aliyun.com/document_detail/62423.htm)
- [Linux 系统检查 Virtio 驱动](https://intl.cloud.tencent.com/zh/document/product/213/9929)

执行`grep -i virtio /boot/config-$(uname -r)`, 有三种情况:
- 如果返回结果中CONFIG_VIRTIO_BLK 参数和CONFIG_VIRTIO_NET 参数取值为 m

    ```bash
    # lsinitrd /boot/initramfs-$(uname -r).img | grep virtio # CentOS 6/CentOS 7/CentOS 8/RedHat 6/RedHat 7 操作系统
    # lsinitramfs /boot/initrd.img-$(uname -r) | grep virtio # Debian/Ubuntu 操作系统
    ```
- 如果在返回结果中CONFIG_VIRTIO_BLK 参数和CONFIG_VIRTIO_NET 参数取值为 y，表示kernel内置了 Virtio 驱动
- 如果在返回结果中没有CONFIG_VIRTIO_BLK 参数和CONFIG_VIRTIO_NET 参数的信息，表示该操作系统不支持virtio

### vm在线加盘后重启未生效
关机再开机即可生效.

### kylinv10 kvm iso安装无法进入安装界面而是停在uefi shell
出现这个UEFI SHELL, 说明系统在加载时出了问题，找不到efi启动文件.

创建vm时osvariant选择kylin v10, 而不是选中centos 7.
