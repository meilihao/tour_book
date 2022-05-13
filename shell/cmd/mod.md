# mod
mod工具在module-init-tools.

# depmod
用于分析可载入模块的相依性, 以供modprobe使用.

## 选项
- -a : 分析所有可用的模块
- -v : 输出详细信息

depmod会遍历文件/lib/modules/`uname -r`/modules.dep解析模块依赖关系, 该文件由depmod -a 命令建立的，保存了内核模块的依赖关系.

# lsmod
显示已加载内核模块的状态, 信息来自`/proc/modules`.

输出信息:
- Module  : 模块名称
- Size  : 模块大小
- Used by : 依赖其他模块的个数 + 被其他模块依赖的列表

# modinfo
显示内核模块的信息

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

- vermagic

    加载模块时, 将vermagic检查值中的字符串是否匹配, 如果它们不匹配, 将得到一个错误, 内核将拒绝加载该模块. modprobe可以通过使用--force标志来克服这一点. 当然, 这些检查是为了保护kernel, 因此使用此选项很危险.
- depends : 依赖的mod

### example
```bash
# modinfo -F filename qla2xxx # 检查光纤驱动模块是否存在
# modinfo  first_time.ko # 查看模块信息
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