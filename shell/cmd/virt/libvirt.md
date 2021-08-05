# libvirt
## 构建
参考:
- [Centos7.6 下编译安装 Libvirt 7.5](https://blog.frytea.com/archives/546/)
- [KVM安装及使用指南](https://bbs.huaweicloud.com/forum/forum.php?mod=viewthread&tid=113876)

前提:
1. qemu

参考源码中的README
```
# -- https://github.com/archlinux/svntogit-community/blob/packages/libvirt/trunk/PKGBUILD
rm /usr/lib/aarch64-linux-gun/libvirt*
wget https://libvirt.org/sources/libvirt-6.0.0.tar.xz
tar -xf libvirt-6.0.0.tar.xz
cd libvirt-6.0.0
mkdir build && cd build
apt install gnutls-bin ebtables
apt install libgnutls-dev libpciaccess-dev libnl-3-dev libnl-route-3-dev libdevmapper-dev libyajl-dev
../configure --prefix=/usr --sysconfdir=/etc --localstatedir=/var --with-qemu # `../autogen.sh --system --with-qemu` # 根据当前环境自动选择编译选项 # configure没有识别到自编译的qemu, 因此需要追加`--with-qemu`; `-–system`参数会尽可能保证安装目录与原生发行版系统的一致性
make && make install
ldconfig
systemctl start libvirtd
virsh version # 验证virsh
reboot # 否则virt-manger连接可能报错
```

libvirt-7.6.0编译:
```bash
$ apt install meson ninja # meson至少0.54
$ pip3 install rst2html5
$ mkdir build
$ meson build --prefix=/usr # meson build -Dsystem=true
$ ninja -C build
$ sudo ninja -C build install
```

> libvirt 7.2.0开始要求`glib 2.56`

> libvirt 6.3.0编译html docs时会报错

## FAQ
### `failed to connect to the hypervisor` & `failed to connect socket to '/var/run/libvirt/libvirt-sock': No such file or directory`
原因: libvirt服务未启动，找不到libvirt-sock.

解决方法: `systemctl start libvirtd`

### `Cannot read CA certificate '/etc/pki/CA/cacert.pem': No such file or directory`
当连接指定主机名(`qemu://system`或`qemu://session`, 使用**两个正斜杠**)时, QEMU 传输默认为 TLS, 这会要求证书.

> 使用三个正斜杠连接到的是本地主机, 例如`qemu:///system`, 不用tls.

解决方法: 见[TLSCreateCACert](https://wiki.libvirt.org/page/TLSCreateCACert)

ca.info:
```conf
cn = libvirt.org
ca
cert_signing_key
expiration_days = 7000
```

### `Cannot read certificate '/etc/pki/libvirt/servercert.pem': No such file or directory`
解决方法: 见[TLSCreateServerCerts](https://wiki.libvirt.org/page/TLSCreateServerCerts)

server.info:
```conf
organization = libvirt.org
cn = xxx
tls_www_server
encryption_key
signing_key
expiration_days = 7000
```

### `systemctl start libvirtd`日志报`direct firewall backed requested, but /sbin/ebtables is not avaliable: No such file or directory`
`apt install ebtables`

### `journalctl`报`failed to connect socket to '/var/run/libvirt/virtqemud-sock': No such file or directory`
参考:
- [libvirt.spec.in](https://github.com/libvirt/libvirt/blob/master/libvirt.spec.in)

virt-manager访问本地qemu虚拟机时用到, 应该是 virtqemud 服务没起来导致的, 还可能是编译libvirt时没有追加`--with-qemu`导致没有编译出该服务.

### libvirtd调试
参考:
- [libvirt ：debug and logging](https://blog.csdn.net/ggg_xxx/article/details/80060672)

方法:
1. `LIBVIRT_DEBUG=1 libvirtd`
1. `vim /etc/libvirt/libvirtd.conf`将`log_level`设为2

> 1时log不友好. 

### `YAJL 2 is required to build QEMU driver`
`apt install libyajl-dev`

### `XDR is required for remote driver `
`apt install libtirpc-dev`

### `Cannot check QEMU binary /usr/libexec/qemu-kvm: No such file or directory`
virt-manger创建vm时提示"警告:KVM不可用", 此时libvirtd就报该错. [根据qemu文档: The KVM project used to maintain a fork of QEMU called qemu-kvm. All feature differences have been merged into QEMU upstream and the development of the fork suspended.](https://wiki.qemu.org/Features/KVM), 即qemu-kvm现在已被弃用, 其中的所有代码以合并入 qemu-system-x86_64.

当前未找到libvirtd将qemu-kvm切换到qemu-system-x86_64的方法, 只是将`/usr/libexec/qemu-kvm`指向了qemu-system-x86_64, 但启用kvm需要添加显示参数`--enable-kvm`

此时libvirtd还会报`this function is not supported by the connect driver: cannot detect host CPU mode for aarch64 architecture`(虚拟机和物理机都会, 且物理机的kvm-ok是提示`KVM acceleration can be used`), 这个问题确定是libvirt接口不兼容引起的, 可参考[鲲鹏服务器安装kvm虚拟机，cannot detect host CPU model for aarch64 architect](https://bbs.huaweicloud.com/forum/thread-75455-1-1.html)和[libvirt启动时出现cannot detect host CPU model for architecture](https://bbs.huaweicloud.com/forum/thread-59280-1-1.html)

### virt-manger创建vm时提示`Falied to setup UEFI: 不能为架构'aarch64'找到任何UEFI二进制路径\nInstall options are limited`
`virt-manager --debug`报`Did not find any UEFI binary path for ...`

```bash
apt install qemu-efi # curl -o /etc/yum.repos.d/firmware.repo https://www.kraxel.org/repos/firmware.repo && yum install -y edk2.git-aarch64 
cat >> /etc/libvirt/qemu.conf <<'EOF'
nvram = [
    "/usr/share/AAVMF/AAVMF_CODE.fd:/usr/share/AAVMF/AAVMF_VARS.fd" # for ubuntu
    "/usr/share/edk2.git/aarch64/QEMU_EFI-pflash.raw:/usr/share/edk2.git/aarch64/vars-template-pflash.raw" # for centos
]
EOF
systemctl restart libvirtd
```

ps:
centos KVM对UEFI的支持包如下:
- 32位的arm处理器需要edk2.git-arm
- 64位的arm处理器需要edk2.git-aarch64
- x86-64处理器需要edk2.git-ovmf-x64

### libvirtd.service自动退出
Ubuntu16.04.6+飞腾主板+libvirt 6.0.0, systemd里没有报错日志, 也没有crash.

### `'HWCAP_CPUID' undeclared`
内核版本太低, 比如4.4, HWCAP_CPUID没有定义. libvirt从6.4.0开始引入它.

## virtsh
virsh 属于 libvirt 工具， libvirt 是目前使用最为广泛的对 KVM 虚拟机进行管理的工具和 API, 它还可管理 VMware, VirtualBox, Hyper-V等.

Libvirt 分服务端和客户端, Libvirtd 是一个 daemon 进程, 是服务端, 可以被本地的 virsh 调用, 也可以被远程的 virsh 调用, virsh 相当于客户端.

### 常用命令
> 可参考`man virsh`

如下命令启动虚拟机： `virsh create <name of virtual machine>`
启动虚拟机： `virsh start <name>`
列出所有虚拟机 (不管是否运行)： `virsh list --all`, `--all`包括没运行的vm, 则只输出运行中的vm
正常关闭 guest ： `virsh shutdown <virtual machine (name | id | uuid)>`
强制关闭 guest ： `virsh destroy <virtual machine (name | id | uuid)>`
挂起vm: `virsh suspend <name>`
恢复被挂起的vm: `virsh resumed <name>`
开机自启动vm: `virsh autostart <name>`
连接vm: `virsh console <name>`
保存虚拟机快照到文件： `virsh save <virtual machine (name | id | uuid)> <filename>`
从快照恢复虚拟机： `virsh restore <filename>`
查看虚拟机配置文件： `virsh dumpxml <virtual machine (name | id | uuid)`
删除vm的配置文件: `virsh undifine <name>`
根据配置文件定义vm: `virsh define <file.xml>`
列出全部 virsh 可用命令： `virsh help`

    ```conf
    # virsh help
    command：

     Domain Management (help keyword 'domain'):
        attach-device                  从一个XML文件附加装置
        attach-disk                    附加磁盘设备
        attach-interface               获得网络界面
        autostart                      自动开始一个域
        blkdeviotune                   设定或者查询块设备 I/O 调节参数。
        blkiotune                      获取或者数值 blkio 参数
        blockcommit                    启动块提交操作。
        blockcopy                      启动块复制操作。
        blockjob                       管理活跃块操作
        blockpull                      使用其后端映像填充磁盘。
        blockresize                    创新定义域块设备大小
        change-media                   更改 CD 介质或者软盘驱动器
        console                        连接到客户会话
        cpu-stats                      显示域 cpu 统计数据
        create                         从一个 XML 文件创建一个域
        define                         从一个 XML 文件定义（但不开始）一个域
        desc                           显示或者设定域描述或者标题
        destroy                        销毁（停止）域
        detach-device                  从一个 XML 文件分离设备
        detach-device-alias            detach device from an alias
        detach-disk                    分离磁盘设备
        detach-interface               分离网络界面
        domdisplay                     域显示连接 URI
        domfsfreeze                    Freeze domain's mounted filesystems.
        domfsthaw                      Thaw domain's mounted filesystems.
        domfsinfo                      Get information of domain's mounted filesystems.
        domfstrim                      在域挂载的文件系统中调用 fstrim。
        domhostname                    输出域主机名
        domid                          把一个域名或 UUID 转换为域 id
        domif-setlink                  设定虚拟接口的链接状态
        domiftune                      获取/设定虚拟接口参数
        domjobabort                    忽略活跃域任务
        domjobinfo                     域任务信息
        domname                        将域 id 或 UUID 转换为域名
        domrename                      rename a domain
        dompmsuspend                   使用电源管理功能挂起域
        dompmwakeup                    从 pmsuspended 状态唤醒域
        domuuid                        把一个域名或 id 转换为域 UUID
        domxml-from-native             将原始配置转换为域 XML
        domxml-to-native               将域 XML 转换为原始配置
        dump                           把一个域的内核 dump 到一个文件中以方便分析
        dumpxml                        XML 中的域信息
        edit                           编辑某个域的 XML 配置
        event                          Domain Events
        inject-nmi                     在虚拟机中输入 NMI
        iothreadinfo                   view domain IOThreads
        iothreadpin                    control domain IOThread affinity
        iothreadadd                    add an IOThread to the guest domain
        iothreaddel                    delete an IOThread from the guest domain
        send-key                       向虚拟机发送序列号
        send-process-signal            向进程发送信号
        lxc-enter-namespace            LXC 虚拟机进入名称空间
        managedsave                    管理域状态的保存
        managedsave-remove             删除域的管理保存
        managedsave-edit               edit XML for a domain's managed save state file
        managedsave-dumpxml            Domain information of managed save state file in XML
        managedsave-define             redefine the XML for a domain's managed save state file
        memtune                        获取或者数值内存参数
        perf                           Get or set perf event
        metadata                       show or set domain's custom XML metadata
        migrate                        将域迁移到另一个主机中
        migrate-setmaxdowntime         设定最大可耐受故障时间
        migrate-getmaxdowntime         get maximum tolerable downtime
        migrate-compcache              获取/设定压缩缓存大小
        migrate-setspeed               设定迁移带宽的最大值
        migrate-getspeed               获取最长迁移带宽
        migrate-postcopy               Switch running migration from pre-copy to post-copy
        numatune                       获取或者数值 numa 参数
        qemu-attach                    QEMU 附加
        qemu-monitor-command           QEMU 监控程序命令
        qemu-monitor-event             QEMU Monitor Events
        qemu-agent-command             QEMU 虚拟机代理命令
        reboot                         重新启动一个域
        reset                          重新设定域
        restore                        从一个存在一个文件中的状态恢复一个域
        resume                         重新恢复一个域
        save                           把一个域的状态保存到一个文件
        save-image-define              为域的保存状态文件重新定义 XML
        save-image-dumpxml             在 XML 中保存状态域信息
        save-image-edit                为域保存状态文件编辑 XML
        schedinfo                      显示/设置日程安排变量
        screenshot                     提取当前域控制台快照并保存到文件中
        set-lifecycle-action           change lifecycle actions
        set-user-password              set the user password inside the domain
        setmaxmem                      改变最大内存限制值
        setmem                         改变内存的分配
        setvcpus                       改变虚拟 CPU 的号
        shutdown                       关闭一个域
        start                          开始一个（以前定义的）非活跃的域
        suspend                        挂起一个域
        ttyconsole                     tty 控制台
        undefine                       取消定义一个域
        update-device                  从 XML 文件中关系设备
        vcpucount                      域 vcpu 计数
        vcpuinfo                       详细的域 vcpu 信息
        vcpupin                        控制或者查询域 vcpu 亲和性
        emulatorpin                    控制火车查询域模拟器亲和性
        vncdisplay                     vnc 显示
        guestvcpus                     query or modify state of vcpu in the guest (via agent)
        setvcpu                        attach/detach vcpu or groups of threads
        domblkthreshold                set the threshold for block-threshold event for a given block device or it's backing chain element

     Domain Monitoring (help keyword 'monitor'):
        domblkerror                    在块设备中显示错误
        domblkinfo                     域块设备大小信息
        domblklist                     列出所有域块
        domblkstat                     获得域设备块状态
        domcontrol                     域控制接口状态
        domif-getlink                  获取虚拟接口链接状态
        domifaddr                      Get network interfaces' addresses for a running domain
        domiflist                      列出所有域虚拟接口
        domifstat                      获得域网络接口状态
        dominfo                        域信息
        dommemstat                     获取域的内存统计
        domstate                       域状态
        domstats                       get statistics about one or multiple domains
        domtime                        domain time
        list                           列出域

     Host and Hypervisor (help keyword 'host'):
        allocpages                     Manipulate pages pool size
        capabilities                   性能
        cpu-baseline                   计算基线 CPU
        cpu-compare                    使用 XML 文件中描述的 CPU 与主机 CPU 进行对比
        cpu-models                     CPU models
        domcapabilities                domain capabilities
        freecell                       NUMA可用内存
        freepages                      NUMA free pages
        hostname                       打印管理程序主机名
        hypervisor-cpu-baseline        compute baseline CPU usable by a specific hypervisor
        hypervisor-cpu-compare         compare a CPU with the CPU created by a hypervisor on the host
        maxvcpus                       连接 vcpu 最大值
        node-memory-tune               获取或者设定节点内存参数
        nodecpumap                     节点 cpu 映射
        nodecpustats                   输出节点的 cpu 状统计数据。
        nodeinfo                       节点信息
        nodememstats                   输出节点的内存状统计数据。
        nodesuspend                    在给定时间段挂起主机节点
        sysinfo                        输出 hypervisor sysinfo
        uri                            打印管理程序典型的URI
        version                        显示版本

     Interface (help keyword 'interface'):
        iface-begin                    生成当前接口设置快照，可在今后用于提交 (iface-commit) 或者恢复 (iface-rollback)
        iface-bridge                   生成桥接设备并为其附加一个现有网络设备
        iface-commit                   提交 iface-begin 后的更改并释放恢复点
        iface-define                   define an inactive persistent physical host interface or modify an existing persistent one from an XML file
        iface-destroy                  删除物理主机接口（启用它请执行 "if-down"）
        iface-dumpxml                  XML 中的接口信息
        iface-edit                     为物理主机界面编辑 XML 配置
        iface-list                     物理主机接口列表
        iface-mac                      将接口名称转换为接口 MAC 地址
        iface-name                     将接口 MAC 地址转换为接口名称
        iface-rollback                 恢复到之前保存的使用 iface-begin 生成的更改
        iface-start                    启动物理主机接口（启用它请执行 "if-up"）
        iface-unbridge                 分离其辅助设备后取消定义桥接设备
        iface-undefine                 取消定义物理主机接口（从配置中删除）

     Network Filter (help keyword 'filter'):
        nwfilter-define                使用 XML 文件定义或者更新网络过滤器
        nwfilter-dumpxml               XML 中的网络过滤器信息
        nwfilter-edit                  为网络过滤器编辑 XML 配置
        nwfilter-list                  列出网络过滤器
        nwfilter-undefine              取消定义网络过滤器
        nwfilter-binding-create        create a network filter binding from an XML file
        nwfilter-binding-delete        delete a network filter binding
        nwfilter-binding-dumpxml       XML 中的网络过滤器信息
        nwfilter-binding-list          list network filter bindings

     Networking (help keyword 'network'):
        net-autostart                  自动开始网络
        net-create                     从一个 XML 文件创建一个网络
        net-define                     define an inactive persistent virtual network or modify an existing persistent one from an XML file
        net-destroy                    销毁（停止）网络
        net-dhcp-leases                print lease info for a given network
        net-dumpxml                    XML 中的网络信息
        net-edit                       为网络编辑 XML 配置
        net-event                      Network Events
        net-info                       网络信息
        net-list                       列出网络
        net-name                       把一个网络UUID 转换为网络名
        net-start                      开始一个(以前定义的)不活跃的网络
        net-undefine                   undefine a persistent network
        net-update                     更新现有网络配置的部分
        net-uuid                       把一个网络名转换为网络UUID

     Node Device (help keyword 'nodedev'):
        nodedev-create                 根据节点中的 XML 文件定义生成设备
        nodedev-destroy                销毁（停止）节点中的设备
        nodedev-detach                 将节点设备与其设备驱动程序分离
        nodedev-dumpxml                XML 中的节点设备详情
        nodedev-list                   这台主机中中的枚举设备
        nodedev-reattach               重新将节点设备附加到他的设备驱动程序中
        nodedev-reset                  重置节点设备
        nodedev-event                  Node Device Events

     Secret (help keyword 'secret'):
        secret-define                  定义或者修改 XML 中的 secret
        secret-dumpxml                 XML 中的 secret 属性
        secret-event                   Secret Events
        secret-get-value               secret 值输出
        secret-list                    列出 secret
        secret-set-value               设定 secret 值
        secret-undefine                取消定义 secret

     Snapshot (help keyword 'snapshot'):
        snapshot-create                使用 XML 生成快照
        snapshot-create-as             使用一组参数生成快照
        snapshot-current               获取或者设定当前快照
        snapshot-delete                删除域快照
        snapshot-dumpxml               为域快照转储 XML
        snapshot-edit                  编辑快照 XML
        snapshot-info                  快照信息
        snapshot-list                  为域列出快照
        snapshot-parent                获取快照的上级快照名称
        snapshot-revert                将域转换为快照

     Storage Pool (help keyword 'pool'):
        find-storage-pool-sources-as   找到潜在存储池源
        find-storage-pool-sources      发现潜在存储池源
        pool-autostart                 自动启动某个池
        pool-build                     建立池
        pool-create-as                 从一组变量中创建一个池
        pool-create                    从一个 XML 文件中创建一个池
        pool-define-as                 在一组变量中定义池
        pool-define                    define an inactive persistent storage pool or modify an existing persistent one from an XML file
        pool-delete                    删除池
        pool-destroy                   销毁（删除）池
        pool-dumpxml                   XML 中的池信息
        pool-edit                      为存储池编辑 XML 配置
        pool-info                      存储池信息
        pool-list                      列出池
        pool-name                      将池 UUID 转换为池名称
        pool-refresh                   刷新池
        pool-start                     启动一个（以前定义的）非活跃的池
        pool-undefine                  取消定义一个不活跃的池
        pool-uuid                      把一个池名称转换为池 UUID
        pool-event                     Storage Pool Events

     Storage Volume (help keyword 'volume'):
        vol-clone                      克隆卷。
        vol-create-as                  从一组变量中创建卷
        vol-create                     从一个 XML 文件创建一个卷
        vol-create-from                生成卷，使用另一个卷作为输入。
        vol-delete                     删除卷
        vol-download                   将卷内容下载到文件中
        vol-dumpxml                    XML 中的卷信息
        vol-info                       存储卷信息
        vol-key                        为给定密钥或者路径返回卷密钥
        vol-list                       列出卷
        vol-name                       为给定密钥或者路径返回卷名
        vol-path                       为给定密钥或者路径返回卷路径
        vol-pool                       为给定密钥或者路径返回存储池
        vol-resize                     创新定义卷大小
        vol-upload                     将文件内容上传到卷中
        vol-wipe                       擦除卷

     Virsh itself (help keyword 'virsh'):
        cd                             更改当前目录
        echo                           echo 参数
        exit                           退出这个非交互式终端
        help                           打印帮助
        pwd                            输出当前目录
        quit                           退出这个非交互式终端
        connect                        连接（重新连接）到 hypervisor
    ```

创建vm:
```bash
qemu-img create -f qcow2 centos_kvm1.qcow2 16G
virt-install \
--virt-type=kvm \
--name=centos-kvm \
--hvm \
--vcpus=1 \
--memory=1024 \
--cdrom=/srv/kvm/CentOS-7-x86_64-Minimal-1810.iso \
--disk path=/srv/kvm/centos_kvm1.qcow2,size=16,format=qcow2 \
--graphics vnc,password=kvm,listen=::,port=5911 \
--network bridge=virbr0 \
--autostart \
--force
```

安装成功后使用任意一个可以访问 KVM 宿主机的带有桌面的设备上的 VNC viewer 进入 `<vm宿主机ip>:5911`, 输入密码 `kvm` 就可以进入虚拟机, 然后继续安装了.

install 常用参数说明展开目录:
```conf
–name指定虚拟机名称
–memory分配内存大小.
–vcpus分配CPU核心数，最大与实体机CPU核心数相同
–disk指定虚拟机镜像，size指定分配大小单位为G.
–network网络类型，此处用的是默认，一般用的应该是bridge桥接.
–accelerate加速
–cdrom指定安装镜像iso
–vnc启用VNC远程管理，一般安装系统都要启用.
–vncport指定VNC监控端口，默认端口为5900，端口不能重复.
–vnclisten指定VNC绑定IP，默认绑定127.0.0.1，这里改为0.0.0.0
–os-type=linux,windows
–os-variant=rhel6

--name      指定虚拟机名称
--ram       虚拟机内存大小，以 MB 为单位
--vcpus     分配CPU核心数，最大与实体机CPU核心数相同
–-vnc       启用VNC远程管理，一般安装系统都要启用.
–-vncport   指定VNC监控端口，默认端口为5900，端口不能重复.
–-vnclisten  指定VNC绑定IP，默认绑定127.0.0.1，这里改为0.0.0.0
--network   虚拟机网络配置
  # 其中子选项，bridge=br0 指定桥接网卡的名称.

–os-type=linux,windows
–os-variant=rhel7.2

--disk 指定虚拟机的磁盘存储位置
  # size，初始磁盘大小，以 GB 为单位.

--location 指定安装介质路径，如光盘镜像的文件路径.
--graphics 图形化显示配置
  # 全新安装虚拟机过程中可能会有很多交互操作，比如设置语言，初始化 root 密码等等.
  # graphics 选项的作用就是配置图形化的交互方式，可以使用 vnc（一种远程桌面软件）进行链接.
  # 我们这列使用命令行的方式安装，所以这里要设置为 none，但要通过 --extra-args 选项指定终端信息，
  # 这样才能将安装过程中的交互信息输出到当前控制台.
--extra-args 根据不同的安装方式设置不同的额外选项
```