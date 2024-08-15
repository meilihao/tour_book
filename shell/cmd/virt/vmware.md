# vmware

## 网络
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

  `the problem is that with PMU disabled in VMWare config, it's not giving the right info to the guest to know it's disabled.`, 因此推测是vmware bug, 打开vmware `启用虚拟化CPU性能计数器`或嵌套vm使用`-cpu host,pmu=off`可能有用(未验证).

- [nested virt: kvm crash "kvm_put_msrs: Assertion `ret == cpu->kvm_msr_buf->nmsrs' failed."](https://github.com/kubernetes/minikube/issues/2968)

env:
- cpu: intel

解决:
```bash
# sudo tee /etc/modprobe.d/qemu-system-x86.conf << EOF
options kvm ignore_msrs=1
EOF
# reboot
```

ps:
1. `VMware ESXi, 7.0.3, 21930508`, 正常