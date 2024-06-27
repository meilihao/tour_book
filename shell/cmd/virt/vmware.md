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

  **vmware workstation 17.5.2 无法在 kernel 6.9.4 上构建所需的内核模块**, 兼任列表见[Supported host operating systems for Workstation Pro 16.x, 17.x and Workstation Player 16.x, 17.x](https://knowledge.broadcom.com/external/article/315653/supported-host-operating-systems-for-wor.html)

  vmware kernel module src位置见[VMware, Debian Kernel Upgrade](https://wiki.debian.org/VMware). patch在[VMware vmmon & vmnet 17.5.1 and Linux kernel 6.8.0 won't compile](https://unix.stackexchange.com/questions/773558/vmware-vmmon-vmnet-17-5-1-and-linux-kernel-6-8-0-wont-compile), 但当前只支持到17.5.1

vmware为提供了三种网络工作模式：Bridged（桥接模式）、NAT（网络地址转换模式）、Host-Only（仅主机模式）.

打开vmware虚拟机，可以在选项栏的“编辑”下的“虚拟网络编辑器”中看到VMnet0（桥接模式）、VMnet1（仅主机模式）、VMnet8（NAT模式）, 它们分别表示各自模式下的虚拟交换机. windows网络连接管理解决仅可见VMware Network Adapter VMnet1和VMware Network Adapter VMnet8.

|选择网络连接属性|意义|
|桥接模式（将虚拟机直接连接到外部网络）VMnet0|使用（连接）VMnet0虚拟交换机，此时虚拟机相当于网络上的一台独立计算机，与主机一样，拥有一个独立的IP地址|
|NAT模式（VMnet8）|使用（连接）VMnet8虚拟交换机，此时虚拟机可以通过主机单向访问网络上的其他工作站（包括Internet网络），其他工作站不能访问虚拟机|
|仅主机模式（VMnet1）|使用（连接）VMnet1虚拟交换机，此时虚拟交换机只能与虚拟机、主机互连，与网络上的其他工作站不能访问|

## FAQ
### 支持kvm虚拟化
选择vm右键选择编辑->cpu->勾选`向客户机操作系统公开硬件辅助的虚拟化`和`I/O MMU`

### 嵌套虚拟化
vmware仅支持两层嵌套:
1. 第一层由硬件虚拟化
2. 第二层由软件模拟

### vmware workstation pro 17.5.2 创建vm时无法选中`virtualization engine`
电脑本身支持虚拟化. 在编辑vm中选中即可
