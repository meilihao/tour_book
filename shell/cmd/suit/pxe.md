# pxe
PXE（ Preboot eXecute Environment，预启动执行环境）是由 Intel 公司开发的技术，能够让计算机通过网络来启动操作系统（前提是计算机上安装的网卡支持 PXE 技术），主要用于
在无人值守安装系统中引导客户端主机安装 Linux 操作系统.

# kickstart
ref:
- [Kickstart](https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/7/html/installation_guide/chap-kickstart-installations)

Kickstart 是一种无人值守的安装方式，其工作原理是预先把原本需要运维人员手工填写的参数保存成一个 ks.cfg 文件，当安装过程中需要填写参数时则自动匹配 Kickstart 生成的文件.

# 无人值守系统
组件:
- dhcpd 分配网卡信息及指引获取驱动文件
- tftp-server 提供驱动及引导文件的传输
- SYSLinux 提供驱动及引导文件

	SYSLinux 是一个用于提供引导加载的服务程序。与其说 SYSLinux 是一个服务程序， 不如说是一个包含了很多引导文件的文件夹。在安装好 SYSLinux 服务程序后， /usr/share/syslinux
目录中会出现很多引导文件
- vsftpd 提供完整系统镜像的传输
- Kickstart 提供安装过程中选项的应答配置

	其实在 root 管理员的家目录中有一个名为anaconda-ks.cfg 的文件, 它就是Kickstart 应答文件.

env:
192.168.10.10: 无人值守系统

```bash
# --- 配置dns
# dnf install -y dhcp-server
# vim /etc/dhcp/dhcpd.conf
allow booting;
allow bootp;
ddns-update-style none;
ignore client-updates;
subnet 192.168.10.0 netmask 255.255.255.0 {
	option subnet-mask 255.255.255.0;
	option domain-name-servers 192.168.10.10;
	range dynamic-bootp 192.168.10.100 192.168.10.200; # 允许了 BOOTP 引导程序协议，旨在让局域网内暂时没有操作系统的主机也能获取静态 IP 地址
	default-lease-time 21600;
	max-lease-time 43200;
	next-server 192.168.10.10;
	filename "pxelinux.0"; # 在配置文件的最下面加载了引导驱动文件 pxelinux.0（这个文件会在下面的步骤中创建），其目的是让客户端主机获取到 IP 地址后主动获取引导驱动文件，自行进入下一步的安装过程
}
# systemctl restart dhcpd
# systemctl status dhcpd
# --- 配置tftp-server
# dnf install -y tftp-server xinetd
# vim /etc/xinetd.d/tftp
service tftp
{
	socket_type = dgram
	protocol = udp
	wait = yes
	user = root
	server = /usr/sbin/in.tftpd
	server_args = -s /var/lib/tftpboot
	disable = no
	per_source = 11
	cps = 100 2
	flags = IPv4
}
# systemctl restart xinetd
# --- 配置syslinux
# dnf install -y syslinux
# cd /var/lib/tftpboot
# cp /usr/share/syslinux/pxelinux.0 . # 把 SYSLinux 提供的引导文件（ 也就是前文提到的文件 pxelinux.0） 复制到 TFTP 服务程序的默认目录中，这样客户端主机就能够顺利地获取到引导文件了。另外在RHEL 8 系统光盘镜像中也有一些需要调取的引导文件
# cp /media/cdrom/images/pxeboot/* . # cdrom已挂载需要安装的iso
# cp /media/cdrom/isolinux/* .
# mkdir pxelinux.cfg
# cp /media/cdrom/isolinux/isolinux.cfg pxelinux.cfg/default # default是开机时的选项菜单, 默认的开机菜单中有 3 个选项： 安装系统、 对安装介质进行检验、 排错模式
# vim pxelinux.cfg/default
default linux # 默认执行linux选项
...
label linux
	menu label ^Install Red Hat Enterprise Linux 8.0.0
	kernel vmlinuz
	append initrd=initrd.img inst.stage2=ftp://192.168.10.10 ks=ftp://192.168.10.10/pub/ks.cfg quiet # 将默认的光盘镜像安装方式修改成 FTP 文件传输方式，并指定好光盘镜像的获取网址以及 Kickstart 应答文件的获取路径
# --- 配置vsftpd
# dnf install -y vsftpd
# vim /etc/vsftpd/vsftpd.conf
...
anonymous_enable=YES
...
# systemctl restart vsftpd
# cp -r /media/cdrom/* /var/ftp # 拷贝光盘内容到ftp
# setsebool -P ftpd_connect_all_unreserved=on
# --- 配置Kickstart
# cp ~/anaconda-ks.cfg /var/ftp/pub/ks.cfg
# chmod +r /var/ftp/pub/ks.cfg
# /var/ftp/pub/ks.cfg
#version=RHEL8
ignoredisk --only-use=sda
autopart --type=lvm
# Partition clearing information
clearpart --none --initlabel
# Use graphical install
graphical
repo --name="AppStream" --baseurl=ftp://192.168.10.10/AppStream # 软件仓库改为由 FTP 服务器提供的网络路径
# Use CDROM installation media
url --url=ftp://192.168.10.10/BaseOS # 由 CDROM 改为网络安装源
# --- 以上是表示安装硬盘的名称为 sda 及使用 LVM 技术
# Keyboard layouts
keyboard --vckeymap=us --xlayouts='us' # 指定键盘布局
# System language
lang en_US.UTF-8

# Network information
network --bootproto=dhcp --device=ens160 --onboot=on --ipv6=auto --activate # 让网卡默认处于 DHCP 模式，否则在几十、上百台主机同时被创建出来后，会因为 IP 地址相互冲突而导致后续无法管理
network --hostname=linuxprobe.com
# Root password
rootpw --iscrypted $6$EzIFyouUyBvWRIXv$y3bW3JZ2vD4c8bwVyKt7J90gyjULALTMLrnZ
ZmvVujA75EpCCn50rlYm64MHAInbMAXAgn2Bmlgou/pYjUZzL1
# X Window System configuration information
xconfig --startxonboot
# Run the Setup Agent on first boot
firstboot --enable
# System services
services --disabled="chronyd"
# System timezone
timezone Asia/Shanghai --isUtc --nontp # 定义了系统默认时区为“上海”
user --name=linuxprobe --password=$6$a5fEjghDXGPvEoQc$HQqzvBlGVyhsJjgKFDTpiCEavS.inAwNTLZm/I5R5ALLKaMdtxZoKgb4/EaDyiPSSNNHGqrEkRnfJWap56m./. --iscrypted --gecos="linuxprobe" # 创建了一个普通用户，密码值可复制/etc/shadow文件中的加密密文, 它由系统自动创建

%packages
@^graphical-server-environment # 要安装的软件来源. graphical-server-environment 即带有图形化界面的服务器环境，它对应的是安装界面中的 Server With GUI 选项

%end

%addon com_redhat_kdump --disable --reserve-mb='auto'

%end

%anaconda
pwpolicy root --minlen=6 --minquality=1 --notstrict --nochanges --notempty
pwpolicy user --minlen=6 --minquality=1 --notstrict --nochanges --emptyok
pwpolicy luks --minlen=6 --minquality=1 --notstrict --nochanges --notempty
%end
```