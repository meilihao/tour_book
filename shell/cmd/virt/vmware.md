# vmware

## workstation网络
ref:
- [vmware虚拟机网络配置详解](https://blog.51cto.com/u_15169172/2710721)
- [虚拟机 VMware Workstation 16 PRO 的网络配置](https://blog.csdn.net/weixin_41905135/article/details/123858658)
- [VMware Workstation Pro for Personal Use (For Linux)下载地址](https://support.broadcom.com/group/ecx/productdownloads?subfamily=VMware+Workstation+Pro)

  ```bash
  $ sudo apt install make pkexec                       // 需要构建工具
  $ sudo apt install linux-headers-$(uname -r)         // Ubuntu/Debian
  $ sudo dnf install "kernel-devel-$(uname -r)"        // AlmaLinux/Fedora
  ```

  **vmware workstation 17.5.2 无法在fedora 40 + kernel 6.9.4/ubuntu 24.04 + kernel 6.8 上构建所需的内核模块**, 兼任列表见[Supported host operating systems for Workstation Pro 16.x, 17.x and Workstation Player 16.x, 17.x](https://knowledge.broadcom.com/external/article/315653/supported-host-operating-systems-for-wor.html)

  vmware kernel module src位置见[VMware, Debian Kernel Upgrade](https://wiki.debian.org/VMware). patch在[VMware vmmon & vmnet 17.5.1 and Linux kernel 6.8.0 won't compile](https://unix.stackexchange.com/questions/773558/vmware-vmmon-vmnet-17-5-1-and-linux-kernel-6-8-0-wont-compile), 但当前只支持到17.5.1

vmware为提供了三种网络工作模式：Bridged（桥接模式）、NAT（网络地址转换模式）、Host-Only（仅主机模式）.

打开vmware虚拟机，可以在选项栏的“编辑”下的“虚拟网络编辑器”中看到VMnet0（桥接模式）、VMnet1（仅主机模式）、VMnet8（NAT模式）, 它们分别表示各自模式下的虚拟交换机. windows网络连接管理解决仅可见VMware Network Adapter VMnet1和VMware Network Adapter VMnet8.

|选择网络连接属性|意义|
|桥接模式（将虚拟机直接连接到外部网络）VMnet0|使用（连接）VMnet0虚拟交换机，此时虚拟机相当于网络上的一台独立计算机，与主机一样，拥有一个独立的IP地址|
|NAT模式（VMnet8）|使用（连接）VMnet8虚拟交换机，此时虚拟机可以通过主机单向访问网络上的其他工作站（包括Internet网络），其他工作站不能访问虚拟机|
|仅主机模式（VMnet1）|使用（连接）VMnet1虚拟交换机，此时虚拟交换机只能与虚拟机、主机互连，与网络上的其他工作站不能访问|

## vsphere网络
ref:
- [一张图了解VMware vSphere网络架构](https://blog.csdn.net/m0_60444349/article/details/135094818)
- [vSphere Distributed Switch 架构](https://docs.vmware.com/cn/VMware-vSphere/7.0/com.vmware.vsphere.networking.doc/GUID-B15C6A13-797E-4BCB-B9D9-5CBC5A60C3A6.html)

VMkernel适配器主要用作管理和vMotion, 支持创建多个VMkernel网卡将管理流量和vMotion流量分流.

标准交换机由ESXI主机创建，分布式交换机由 vCenter创建.

ESXI标准vSwitch支持自定义的VLAN ID，以实现网络隔离. 根据VLAN ID不同，可分为三种网络：
1. VLAN ID 0 阻止任何携带VLAN Tag的数据包=Access
1. VLAN ID 4095允许通过任何携带VLAN Tag的数据包=All
1. VLAN ID 1-4094仅允许携带指定VLAN ID tag的数据包=允许特定VLAN通过

## esxcli
ref:
- [VMwareCLI命令参考](https://xstarcd.github.io/wiki/Cloud/VMWareCLI.html)
- [ESXI中虚拟机常用命令](https://blog.csdn.net/qq_43642222/article/details/123333202)
- [使用 pktcap-uw 实用程序捕获和跟踪网络数据包](https://docs.vmware.com/cn/VMware-vSphere/7.0/com.vmware.vsphere.networking.doc/GUID-5CE50870-81A9-457E-BE56-C3FCEEF3D0D5.html)
- [使用 vSphere ESXi 5.5 中的 pktcap-uw 工具捕获网络流量](https://www.dell.com/support/kbdoc/zh-cn/000129233/%E6%8D%95%E8%8E%B7-%E7%BD%91%E7%BB%9C-%E6%B5%81%E9%87%8F-%E4%BD%BF%E7%94%A8-pktcap-uw-%E5%B7%A5%E5%85%B7-%E5%9C%A8-vsphere-esxi-5-5)

```sh
# esxcfg-nics -l # 查看网卡
# esxcli network nic get -n vmnic0 # 查看更加详细的网卡信息
# vmkchdev -l |grep vmnic0 # 显示网卡的VID，DID 等信息
# esxcli vm process list # 正在运行的vm
# vim-cmd vmsvc/getallvms # 获取vm id
# esxcli network vm list
World ID  Name                 Num Ports  Networks              
--------  -------------------  ---------  ----------------------
  193286  wzh-centos-25.155            1  VM Network            
  221168  wzh-DPMP-3.0-25.157          1  VM Network            
  184549  windows-25.158               1  VM Network            
  127634  DPMP3.0                      2  VM Network, VM Network
# esxcli network vm port list -w 221168
   Port ID: 33554458
   vSwitch: vSwitch0
   Portgroup: VM Network
   DVPort ID: 
   MAC Address: 00:50:56:ac:e2:16
   IP Address: 0.0.0.0
   Team Uplink: vmnic0
   Uplink Port ID: 33554436
   Active Filters:
# net-stats -l # 一步完成上面两个命令
# pktcap-uw -A # 支持的抓包位置
# pktcap-uw --switchport 33554458 -c 2500 -o vnic.pcap # --switchport是抓虚拟机网卡; 33554458是某虚拟机虚拟网卡编号; 抓2500个包; vnic.pcap是文件名. 只能抓从VM发出去的包
# pktcap-uw --switchport 33554458 -c 2500 --capture PortOutput -o vnic_in.pcap # 抓发往VM的包
# pktcap-uw --switchport 33554458 --capture PortInput -o PortInput1.pcap & pktcap-uw --switchport 33554458 --capture PortOutput -o PortOutput1.pcap & # 一步完成上面两个命令
# kill $(lsof |grep pktcap-uw |awk '{print $1}'| sort -u) # 结束抓包
```

> 抓包: 暂无可一次抓双向包的方法, 可开两个SSH2窗口同时分别抓入向和出向的包, 抓完之后可以把这2个.pcap文件用wireshark进行合并成一个.pcap文件：`wireshark-->file-->merge`.

![ESXi网络架构](/misc/img/shell/fddc724073a4b6e81ef165f9d991ecc6.png)
![ESXi抓包位置](https://blogs.vmware.com/vsphere/files/2018/12/pktcap-overview-1024x923.png): 常用位置+vmk网络

  > ![ESXi抓包位置offline](/misc/img/shell/pktcap.png)

  常用位置:
  - 虚拟机网卡(portgroup->vm): `--switchport <id> --capture VnicTx/VnicRx`
  - 分布式交换机出口(vds->portgroup): `--switchport <id> --capture PortInput/PortOutput`
  - 分布式交换机入口(kernel->vds): `--uplink <vmnicX> --capture PortInput/PortOutput`
  - kernel(物理网卡): `--uplink <vmnicX> --capture UplinkSndKernel/UplinkRcvKernel`

pktcap-uw其他参数:
- -c: 指定抓取个数
- --ip 117.132.191.126: 指定dst or src ip
- --ng: 写入注释信息, 需与`-o`合用
- --proto 0x01: 指定协议, 比如ICMP(0x01)
- --tcpport 443: 指定端口
- --uplink: 指定uplin(物理)网卡

## FAQ
### 支持kvm虚拟化
选择vm右键选择编辑->cpu->勾选`向客户机操作系统公开硬件辅助的虚拟化`和`I/O MMU(或硬件CPU和MMU)`

### 嵌套虚拟化
vmware仅支持两层嵌套:
1. 第一层由硬件虚拟化
2. 第二层由软件模拟

### vmware workstation pro 17.5.2 创建vm时无法选中`virtualization engine`
电脑本身支持虚拟化. 在编辑vm中选中即可

### 卸载
ref:
- [Uninstall Workstation Pro from a Linux Host](https://docs.vmware.com/en/VMware-Workstation-Pro/17/com.vmware.ws.using.doc/GUID-05E4C876-F32C-49D2-82B4-8C759691E7F5.html)

```bash
$ su root
# vmware-installer -u vmware-workstation
```

### [Open Virtualization Format (OVF) Tool](https://developer.broadcom.com/tools/open-virtualization-format-ovf-tool/latest)
`ovftool TrueNAS-SCALE-24.04.1.1.vmx  truenas.ovf`, ovf可被virtualbox导入.

### VMWare启用嵌套虚拟化后使用libguestfs时报"libguestfs: error: could not create appliance through libvirt...qemu-system-x86_64: error: failed to set MSR 0x48f to 0x7fffff00036dfb...qemu-system-x86_64: .../qemu-4.2.1/target/i386/kvm.c:2655: kvm_buf_set_msrs: Assertion `ret == cpu->kvm_msr_buf->nmsrs' failed. [code=1 int1=-1]"
ref:
- [解决 AMD CPU 使用 VMWare 在嵌套虚拟化中用 qemu 启动虚拟机提示 Assertion ret == cpu-kvm_msr_buf-nmsrs failed](https://blog.csdn.net/kunyus/article/details/106986126)
- [`[37/47] KVM: introduce module parameter for ignoring unknown MSRs accesses`](https://patchwork.kernel.org/project/kvm/patch/1250686963-8357-38-git-send-email-avi@redhat.com/)

  这个配置并没有真正解决问题，只是通过配置这个参数可以使 kvm 忽略相关错误
- [**Assertion `ret == cpu->kvm_msr_buf->nmsrs' failed**](https://bugs.launchpad.net/qemu/+bug/1661386)

  ESXi 6.0.0 Build 2494585, 有问题.

  `the problem is that with PMU disabled in VMWare config, it's not giving the right info to the guest to know it's disabled.`, 因此推测是vmware bug.
- [nested virt: kvm crash "kvm_put_msrs: Assertion `ret == cpu->kvm_msr_buf->nmsrs' failed."](https://github.com/kubernetes/minikube/issues/2968)

env:
- cpu: intel

解决:
1. 打开vmware `启用虚拟化CPU性能计数器`(已验证, 可行)或嵌套vm使用`-cpu host,pmu=off`(未验证) from [**Assertion `ret == cpu->kvm_msr_buf->nmsrs' failed**](https://bugs.launchpad.net/qemu/+bug/1661386)

  > `启用虚拟化CPU性能计数器`容易吃内存
1. 忽略error
```bash
# sudo tee /etc/modprobe.d/qemu-system-x86.conf << EOF
options kvm ignore_msrs=1
EOF
# reboot
```

ps:
1. `VMware ESXi, 6.5.0, 4564106`, 报错
1. `VMware ESXi, 7.0.3, 21930508`, 正常

### exsi开启ssh
方法:
1. 通过exsi ui

  1. 在左侧导航器中选中“主机”，然后鼠标右击，在“服务”中点击“启用安全Shell（SSH）”
  2. 开启成功后，在主机页面中会有ssh服务已经开启的提示。


1. 通过vSphere

  1. 登录vSphere vCenter后，选择要开启ssh服务的主机，点击“配置”，在配置页面左侧栏中选择“系统”下的“服务
  1. 选中SSH，再点击“启动”即可