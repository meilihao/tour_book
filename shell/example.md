# example
1. Check if run in rescue mode
```bash
if [ "`grep 'rescue' /proc/cmdline 2> /dev/null || true`" != "" ] # 舍弃stderr且[`|| true`](https://unix.stackexchange.com/questions/325705/why-is-pattern-command-true-useful/325727)保证命令的exit code是0, 因此不会中断脚本
```

1. trap
参考:
- [为shell布置陷阱：trap捕捉信号方法论](https://www.cnblogs.com/f-ck-need-u/p/7454174.html)

```bash
trap on_exit EXIT # Assign exit handler. bash 提供的一个叫做 EXIT 的伪信号(exit为0)，trap 它时当脚本因为任何原因退出时，相应的命令或函数就会执行. `on_exit`函数用`-`替换时表示不执行exit handler
```

1.  other
```sh
cat /proc/cmdline # 获取内核启动参数
# BOOT_IMAGE=/boot/vmlinuz-5.3.0-40-generic root=UUID=df199c17-0b8f-4337-835b-788bce1ecd7c ro quiet splash vt.handoff=7
blkid -U  df199c17-0b8f-4337-835b-788bce1ecd7c # 获取分区by uuid
# /dev/md0
if [ -d "/sys/block/`basename /dev/md0`/md" ] # /dev/md0 is md
disk_size_sectors="`blockdev --getsz /dev/md0`" # 获取/dev/md0的sector数
dis_size_gb=$((${disk_size_sectors}*512/1024/1024/1024)) # 获取/dev/md0的大小, 单位gb.
dd if=/dev/md0 of=/tmp/mbr.dat bs=512 count=1 # 备份mbr
dd if=/dev/md0 of=/tmp/mbr.dat bs=512 count=1 # 还原mbr
df -T # 只可以查看已经挂载的分区和文件系统类型
file -s /dev/sda4 # 分区文件系统类型
lsblk -f # 查看分区挂载点和文件系统
parted -l #可以查看未挂载的文件系统类型，以及哪些分区尚未格式化

 if [ "`dmidecode -s system-product-name 2>&1 | grep 'SSG-2028R-DE2CR24L' 2>&1`" != "" ] # 匹配硬件型号
elif [ "`dmidecode -s system-product-name 2>&1 | grep 'VMware Virtual Platform' 2>&1`" != "" ]  #匹配 VMware Virtual Platform
elif [ "`dmidecode -s system-product-name 2>&1 | grep 'VirtualBox' 2>&1`" != "" ]  # VirtualBox
elif [ "`dmidecode -s system-product-name 2>&1 | grep 'Bochs' 2>&1`" != "" ]  # QEMU/KVM Guest OS

ls -d /sys/block/sd*/device/enclosure_device:* 2>&1 # 获取共享存储列表, disk上线可能需要时间, 可查看dmesg日志确定, 因此需要等待一段时间再执行该命令
/sys/block/sde/device/enclosure_device:Slot01 # 机箱(这里是共享存储整列)的槽位
lsblk -d -o name,rota # 区别hdd和ssd靠rota; 区分ssd和nvme是name的前缀.

for devname in `ls -d /sys/block/sd*/device/enclosure_device:* | cut -d/ -f4` __AVOID_EMPTY__ # __AVOID_EMPTY__表示给前面的list追加了一项
        do
            echo $devname
           if [ "`sg_inq -e -p 0xb1 /dev/${devname} 2> /dev/null | grep 'Non-rotating medium'`" = "" ]  # 判断不是SSD, 因为hdd磁片是可旋转的, 它靠/sys/block/${dev}/queue/rotational的返回值, 1表示可旋转是hdd, 这种方法有个问题，那就是/sys/block/下面不只有硬盘，还可能有别的块设备，它们会干扰你的判断. 其他方法: `fdisk -l`列出磁盘详情, 可在输出结果中寻找一些HDD特有的关键字，比如：”heads”（磁头），”track”（磁道）和”cylinders”（柱面）.
        done

grep -v "^appliance_type=" /root/appliance.conf > /root/appliance.conf.tmp 将除appliance_type开头外的内容移入/root/appliance.conf.tmp

lsblk |grep "^vd" > /dev/null 2>&1
if "$?" = "0" ] # run in kvm, 也可用virt-what解决.

# /etc/localtime是用来描述本机时间(对应date命令)，而 /etc/timezone是用来描述本机所属的时区
cat /etc/timezone 2>/dev/null # 获取时区
timedatectl set-timezone 'Asia/Shanghai'  2>/dev/null # 设置时区

if [ "`grep 'inet static' /etc/network/interfaces || true`" !="" ] # 已设置network
ls -d /sys/class/net/eth* # 查看是否存在ethN之类的网卡
grep PCI_SLOT_NAME /sys/class/net/${netdev_name}/device/uevent 2 >/dev/null |cut -c15- # 获取设备的PCI_SLOT_NAME
lspci -mm -s ${PCI_SLOT_NAME} |sed -e 's/" "/\n/g' |sed -e 's/"/\n/g' |head -1 # 获取slot
lspci -mm -s ${PCI_SLOT_NAME} |sed -e 's/" "/\n/g' |sed -e 's/"/\n/g' |head -3|tail -1 # 获取vendor(供应商)
lspci -mm -s ${PCI_SLOT_NAME} |sed -e 's/" "/\n/g' |sed -e 's/"/\n/g' |head -4|tail -1 # 获取设备类型
ls /sys/class/net/*/address # 已分配过ip的netdev
skip_netdev="docker0 inter0 lo"
# dialog --inputbox text height width
dialog  --title "xxx" --inputbox "please input your name:" 10 30 2> /tmp/name.txt #这里的`2>`表示将stderr重定向到/tmp/name.txt
if [[ ${hostname} =~^[A-Za-z][A-Za-z0-9-]*$ ]] # 验证hostname by RFC(952,1123)
ifdown ${iface} # 网卡下线
ifup ${iface} # 网卡上线
cat /etc/iscsi/initiatorname.iscsi # it saved a unique InitiatorName
service open-iscsi restart

{
    read a
    read b
} < /tmp/xxx # {}表示语句块, 即读取/tmp/xxx后执行{}内的命令.

ipmitool raw 0x30 0x70 0x20 2>&1|tr -d ' ' # 获取双控节点的节点标志, 需要硬件支持
/lib/udev/scsi_id --page=0x83 --whitelisted /dev/sda # 查看/dev/sda的设备信息
#> 1ATA     Hoodisk SSD                             KATMC9A11220768
sg_ses --page=7 /dev/bsg/1:0:0:0
#>  ATA       Hoodisk SSD       61.3
#>    disk device (not an enclosure)
#>sg_ses failed: Illegal request, Invalid opcode
echo /sys/class/enclosure/*/*/device/block/* | tr " " "\n" # 获取disks in SAS enclosure
```