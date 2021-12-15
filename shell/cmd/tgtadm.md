# tgtadm

tgtadm常用于管理三类对象：
    - target:创建new，删除，查看
    - lun：创建，查看，删除
    - account：创建用户，绑定，解绑定，删除，查看

## 格式

    tgtadm --lld [driver] --op [operation] --mode [mode] [OPTION]...

选项:
- -L --lld : 表示驱动程序driver, 通常是iscsi
- -m --mode : 指定操作的对象，mode为target，logicalunit等
- -o --op [operation]:对指定的对象所要做的操作，operation有delete,new,bind,show,unbind等

OPTION常用选项:
- -t --tid :用来指定Target的ID
- -T --targetname ：指定Target名称

    Target的名称格式：`iqn.xxxx-yy.reversedoamin.STRING[:substring]`, 比如`iqn.2015-11.com.a.web:server1`, 其中：
    - iqn为iqn前缀
    - xxxx为年份
    - yy为月份
    - reversedomain为所在域名的反写
    - STRING为字符串
    - substring为子字符串

- -l --lun :指定lun的号码
- -b --backing-store ：关联到指定lun上的后端存储设备，此例为分区
- -I --initiator-address ：指定可以访问Target的IP地址

## example
```bash
# --- 添加一个新的 target 且其ID为 [id]， 名字为 [name]: `--lld [driver] --op new --mode target --tid=[id] --targetname [name]``
tgtadm --lld iscsi --op new --mode target --tid 1 -T iqn.2013-05.com.magedu:tsan.disk1

# --- 显示所有或某个特定的target: `--lld [driver] --op show --mode target [--tid=[id]]`
tgtadm --lld iscsi --op show --mode target
tgtadm --lld iscsi --op show --mode target --tid 1 # 显示刚创建的target
tgtadm --lld iscsi --op show --mode system # 查看tgtadm支持的参数

# --- 向某ID为[id]的设备上添加一个新的LUN，其号码为[lun]，且此设备提供给initiator使用。[path]是某“块设备”的路径，此块设备也可以是raid或lvm设备。lun0已经被系统预留: `--lld [driver] --op new --mode=logicalunit --tid=[id] --lun=[lun] --backing-store [path]`
tgtadm --lld iscsi --op new --mode logicalunit --tid 1 --lun 1 -b /dev/sda1 # 创建LUN，号码为1

# --- 删除ID为[id]的target: `--lld [driver] --op delete --mode target --tid=[id]`
tgtadm --lld iscsi --op delete --mode target --tid 1
# --- 删除target [id]中的LUN [lun]：`-lld [driver] --op delete --mode=logicalunit --tid=[id] --lun=[lun]`
sudo tgtadm --lld iscsi --op delete --mode logicalunit --tid 1 --lun 1

# --- 开放给192.168.0.0/24网络中的主机访问(其中的-I相当于--initiator-address)：
tgtadm --lld iscsi --op bind --mode target --tid 1 -I 192.168.85.0/24


# ---定义某target的基于主机的访问控制列表，其中[address]表示允许访问此target的initiator客户端的列表: `--lld [driver] --op bind --mode=target --tid=[id] --initiator-address=[address]`
# --- Create a new account
tgtadm --lld iscsi --op new --mode account --user administrator --password 123456
tgtadm --lld iscsi --op show --mode account

# --- Assign this account to a target:
tgtadm --lld iscsi --op bind --mode account --tid 1 --user administrator
tgtadm --lld iscsi --op show --mode target

# --- 解除target [id]的访问控制列表中[address]的访问控制权限：`--lld [driver] --op unbind --mode=target --tid=[id] --initiator-address=[address]`
# ---Set up an outgoing account. First, you need to create a new account like the previous example
tgtadm --lld iscsi --op new --mode account --user abc --password 123456
tgtadm --lld iscsi --op show --mode account
tgtadm --lld iscsi --op bind --mode account --tid 1 --user abc --outgoing
tgtadm --lld iscsi --op show --mode target
```