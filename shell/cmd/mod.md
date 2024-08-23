# mod
mod工具在module-init-tools.

# depmod
用于分析可载入模块的相依性, 以供modprobe使用.

## 选项
- -a : 分析所有可用的模块
- -v : 输出详细信息

depmod会遍历文件/lib/modules/`uname -r`/modules.dep解析模块依赖关系, 该文件在编译kernel时由depmod命令建立的，保存了内核模块的依赖关系.

# lsmod
显示已加载内核模块的状态, 信息来自`/proc/modules`.

输出信息:
- Module  : 模块名称
- Size  : 模块大小
- Used by : 依赖其他模块的个数 + 被其他模块依赖的列表

# modinfo
显示内核模块的信息, 比如配置参数

> 查看mod配置也可用: /sys/module/<mod_name>/parameters

字段说明:
- firmware : driver支持的固件
- alias:

    比如`pci:v00008086d00005A84sv*sd*bc03sc*i*`, 可将其分成以下字符部分：
    - v00008086：v代表供应商ID, 它标识硬件制造商. 该清单由PCI特别兴趣小组维护. `0x8086`代表"英特尔公司"
    - d00005A84：d代表由制造商选择的设备ID. 此ID通常与供应商ID配对, 以形成硬件设备的唯一32位标识符
    - `sv*，sd*`：子系统供应商版本和子系统设备版本用于设备的进一步标识（`*`指示它将匹配任何内容）
    - bc03：基类. 它定义了它是哪种设备: IDE interface, Ethernet controller, USB Controller, .... bc03表示Display controller. lspci将数字映射到设备类.
    - `sc*`：基类的子类
    - `i*`：界面

    depmod分析模块并提前其中的设备表以获得ID, 并写入/lib/modules/<version>/modules.alias

- vermagic

    加载模块时, 将vermagic检查值中的字符串是否匹配, 如果它们不匹配, 将得到一个错误, 内核将拒绝加载该模块. modprobe可以通过使用--force标志来克服这一点. 当然, 这些检查是为了保护kernel, 因此使用此选项很危险.
- depends : 依赖的mod

### example
```bash
# modinfo -F filename qla2xxx # 检查光纤驱动模块是否存在
# modinfo first_time.ko # 查看模块信息
```

# modprobe命令
用于对Linux内核中添加或移除模块

加载指定的模块(会自动处理依赖)到内核时, 若在载入过程中发生错误，在modprobe会卸载整组的模块. 即modprobe会自动处理依赖.

modprobe会查看模块 目录`/lib/modules/$(uname -r)`里面的所有模块和文件，除了可选的/etc/modprobe.conf配置文件和/etc/modprobe.d目录外.
modprobe需要一个最新的modules.dep(`/lib/modules/$(uname -r)/modules.dep`)文件，可以用depmod来生成. 该文件列出了每一个模块需要的其他模块，modprobe使用这个去自动添加或删除模块的依赖.

> modprobe就是调用insmod和rmmod来实现的.

### 选项
- -a : 加载一组匹配的模块
- -c : 输出所有模块的配置信息
- -C : 重载默认配置文件(/etc/modprobe.conf或/etc/modprobe.d)
- -D : 打印模块依赖
- -n : 不实际执行. 可以和-v选项一起使用，调试非常有用
- -r : 选项后指定模块时为卸载指定模块(会清理依赖). 与rmmod功能相同.
- -v : 详细信息
- -q : 不提示任何错误信息

### example
```bash
# modprobe qla2xxx # 载入光纤驱动模块
# modprobe -r igb  # 删除igb模块
# modprobe igb  max_vfs=7 # 模块选项方法1
# echo "options igb max_vfs=7"  >>/etc/modprobe.d/igb.conf # 模块选项方法2
# modprobe igb
```

# insmod
插入内核模块

### example
```bash
# insmod /lib/modules/`uname -r`/kernel/zfs/zfs.ko
```

> 安装已加载的mod会报错: "insmod: ERROR: could not insert module xxx.ko: File exists"

# rmmod
移除内核模块

### example
```bash
# rmmod uvcvideo
```

# ksyms
显示内核符号和模块符号表的信息. 信息来自`/proc/kallsyms`

# dracut
ref:
- [dracut.bootup](https://man7.org/linux/man-pages/man7/dracut.bootup.7.html)
- [Dracut on shutdown](https://github.com/redhat-plumbers/dracut-fedora/blob/main/man/dracut.asc)

Dracut 是一个用于构建 initramfs cpio 档案的工具.

在`/etc/dracut.conf.d`配置, 配置文件格式是`add_drivers+=" xxx xxx "`(**两边需有空格**)

```bash
# dracut --list-modules # 列出系统上所有可用的 dracut 模块. 所有 dracut 模块都位于 /usr/lib/dracut/modules.d 目录中. 在此目录中，所有模块都表示为子目录, 并包含一系列脚本. 每个模块都提供特定的功能
# dracut --kver 5.14.14-300.fc35.x86_64 # 为特定内核版本构建 initramfs
# dracut --regenerate-all --force # 为所有现有内核构建或重新构建 initramfs. 如果特定内核的 initramfs 已经存在, 则需要`--force`
# dracut -H --force # 通常使用 dracut 生成 initramfs 时, 会创建通用主机配置, 即包含了启动通用机器所需的所有内容, 以确保最大可能的兼容性. 如果只想将特定机器实际需要的内容放入 initramfs 中, 可以使用 -H 选项（--hostonly 的缩写）.
# lsinitrd /boot/initramfs-5.14.0-130.el9.x86_64.img # 列出 initramfs 中包含的文件, 该脚本实际就是使用了dracut
# dracut --include /custom-content.conf /etc/custom-content.conf --force # 使用 --include <sourcePath> <targetPath> 在 initramfs 中包含额外文件
# dracut --install "/custom-content.conf /custom-content0.conf" --force # --install 可用于在 initramfs 中包含文件. 与 --include 的主要区别在于: 文件安装在 initramfs 中, 与它们在系统中的位置相同.
# man dracut.conf
```

> mkinitrd（make initial ramdisk）是一个兼容包装器, 它调用dracut来生成initramfs.

# update-initramfs
```bash
# update-initramfs
```

## FAQ
### modprobe和insmod区别
insmod不能处理依赖, 而modprobe可以.

### 黑名单
在 /etc/modprobe.d/ 中创建 `.conf` 文件，使用 **blacklist 关键字屏蔽不需要的模块. blacklist仅屏蔽自动装入, 而不禁止手动操作**.

注意: blacklist 命令会屏蔽一个模块，所以它不会自动装入，但是如果其它非屏蔽模块依赖该模块或手动加载该模块，系统依然会装入它. 要避免这个行为，可以让 modprobe 使用自定义的 install 命令，直接返回导入失败：
```conf
$ vim /etc/modprobe.d/blacklist.conf
...
install MODULE /bin/false
...
```

这样就可以"屏蔽"模块及所有依赖它的模块.

同样可以通过内核命令行(位于 GRUB2 或 Syslinux)禁用模块：
```conf
modprobe.blacklist=modname1,modname2,modname3 # 当某个模块导致系统无法启动时，可以使用此方法禁用模块
```

如果出现模块在启动时未加载，而且启动日志中(journalctl -b) 显示模块被屏蔽，但是 /etc/modprobe.d/ 中未找到屏蔽设置，请检查 /usr/lib/modprobe.d/ 目录.

### 开机自动加载
当内核检测到新设备时，它将运行modprobe并向其传递一个标识该设备的名称. 大多数设备通过供应商和型号的注册号进行标识，例如PCI或USB标识符. modprobe查询模块别名表以查找包含该特定设备的驱动程序的文件的名称. 类似的原理适用于非硬件设备的驱动程序，例如文件系统和密码算法.

一旦modprobe确定了哪个模块文件（.ko）包含请求的驱动程序，它将模块文件加载到内核中：模块代码被动态加载到内核中. 如果模块加载成功，它将出现在lsmod列表中.

当内核检测到新的可热插拔硬件时，例如USB外设连接时，模块会自动加载. 操作系统还对枚举在启动过程中早期在系统上存在的所有硬件进行了检查，以便为启动时存在的外围设备加载驱动程序.

也可以使用modprobeor insmod命令手动请求加载模块。大多数发行版都包含一个启动脚本，用于加载中列出的模块/etc/modules。加载模块的另一种方式是，如果它们是模块的依赖项：如果模块A依赖于模块B，则modprobe A在加载A之前先加载B。

加载模块后，即使使用该驱动程序的所有设备都已断开连接，模块也将保持加载状态，直到明确卸载为止。很久以前，有一种机制可以自动卸载未使用的模块，但是，如果我没记错的话，udev出现在现场时，它已被删除。我怀疑自动模块卸载不是一个常见功能，因为可能需要自动卸载的系统大多数是台式计算机，它们无论如何都具有大量内存（按驱动程序代码的大小）。

### module开机自加载配置
参考:
- [Kernel module (简体中文)](https://wiki.archlinux.org/title/Kernel_module_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87))

方法有两种:
1. `etc/modules-load.d` by systemd

    ```bash
    cat << EOF > /etc/modules-load.d/nbd.conf
    nbd
    EOF
    ```
2. `/etc/modules`

    `echo "nbd" >> /etc/modules`


配置加载模块时的参数:
```bash
cat << EOF > /etc/modprobe.d/nbd.conf
options nbd nbds_max=512
EOF
```

### 卸载dkms module
```bash
# dkms status
rtl8812AU, 4.3.14, 4.4.0-45-generic, x86_64: installed
rtl8812AU, 4.3.14, 4.4.0-47-generic, x86_64: installed
# dkms remove rtl8812AU/4.3.14 --all # 指定module version
# dkms uninstall -k 4.4.0-45-generic rtl8812AU # 移除指定kernel的
```

src残留:
1. `/usr/src`

### 查看initrd
lsinitrd xxx.img

### dracut添加驱动
ref:
- [安装原生的KVM驱动](https://support.huaweicloud.com/usermanual-ims/ims_01_0326.html)
- [Installing Native KVM Drivers](https://support.huaweicloud.com/intl/en-us/usermanual-ims/ims_01_0326.html)

    **推荐使用修改`/etc/dracut.conf`的方法**

```bash
# check for virtio drivers
lsinitrd /boot/initramfs-$version.img | grep virtio

# if not found, add them
cd /boot
dracut -f initramfs-$version.img --add-drivers "virtio_blk virtio_scsi virtio_net virtio virtio_pci virtio_ring virtio-rng virtio_console virtio_balloon virtio_gpu virtio_input"

# prove it worked
lsinitrd /boot/initramfs-$version.img | grep virtio
```

drivers:
- qla2xxx
- sg : scsi
- mpt3sas: sas
- ahci: sata

可通过`cat /lib/modules/`uname -r`/modules.dep |grep virtio`或`cat /boot/config-$version | grep -i virtio`查找virtio驱动, 因为部分virtio驱动可能不存在(比如virtio_ring), 还有部分virtio驱动已直接构建在内核中

### dracut/update-initramfs
dracut/update-initramfs都是管理initramfs的工具.

dracut 主要用于 Fedora、CentOS/RHEL 等, update-initramfs 用于 Debian、Ubuntu 等发行版.