# windows


## 注册表
- 网络配置: `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces\<InterfaceGUID>`
- initiator iqn, 见[/shell/cmd/suit/target.md]

## FAQ
### 安装驱动
[`pnputil -i -a *.inf`, 需管理员权限](https://help.aliyun.com/document_detail/217543.html#section-1kb-hov-812)

### 查看dll依赖
- [lucasg/Dependencies](https://zhuanlan.zhihu.com/p/395557318)

	需要.net framework>=4.6.2
- `git for windows`里的ldd

### 磁盘签名
windows为了区分计算机系统上的存储设备, 每个存储设备都标有一个名为`磁盘签名`的唯一编号, 用于标识. 唯一磁盘标识符存储在分区表信息中.

在XP和Vista之类的旧版本中，签名冲突通常会被忽视，因为Windows系统会自动替换报告重复签名的磁盘签名. 但是对于Windows 7，Windows 8和Windows 10，磁盘签名冲突的处理方式不同, 当两个存储设备具有相同的磁盘签名时，创建磁盘签名冲突的辅助驱动器将关闭，并且在修复冲突之前无法安装使用.

签名位置:
- mbr: Windows磁盘标签占用引导程序后的4个字节，其地址在偏移1B8H~1BBH处，是Windows系统对硬盘初始化时写入的一个磁盘标签
- gpt: GPT磁盘的GUID

```msdoc
> diskpart
> list disk
> select disk <N> # 选择磁盘
> list partition # 查看分区
> uniqueid disk # 查看磁盘签名
> uniqueid disk ID = 1456ACBD # 修改签名: MBR磁盘的十六进制格式或GPT磁盘的GUID
```

### 修复ntfs
`ntfsfix /dev/sdx`

### 查看端口占用的进场
步骤:
- `netstat –ano|findstr <端口号>`, 找到进程pid
- `tasklist|findstr <PID号>`, 找到进程

### 获取cpu序列号
`wmic CPU get ProcessorID`

### 是否支持virtio网卡
1. 列出系统中所有网络适配器的信息，包括驱动程序名称（DriverName）和制造商（Manufacturer）等

	- PowerShell: `Get-NetAdapter | Select-Object Name, InterfaceDescription, DriverName, DriverVersion, Manufacturer`
	- cmd: `wmic nic get AdapterType, Name, Manufacturer, NetConnectionID`
1. 如果 VirtIO 网卡已正确安装并且正在使用，应该能在输出中找到相关信息，比如驱动程序名称可能包含 "virtio" 或者制造商是 "Red Hat"

