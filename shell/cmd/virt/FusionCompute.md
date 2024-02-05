# FusionCompute
ref:
- [](https://forum.huawei.com/enterprise/en/fusioncompute-architecture/thread/667270024320139264-667213860102352896)

	FusionCompute由计算节点代理（CNA, Computing Node Agent）和虚拟资源管理（VRM, Virtualization Resource Management）两部分组成.

	VRM是FusionCompute的管理中枢, 负责资源的分配、调度、维护、监控、资源管理管理, 是基于Linux操作系统的但运行在CNA上的VM.

	CNA类似于KVM中的QEMU+KVM模块. CNA提供虚拟化功能, 部署在集群中，将集群中的计算资源、存储资源和网络资源虚拟化为资源池供用户使用. CNA也基于Linux操作系统.

	CNA的组成:
	1. UVP (Unified Virtualization Platform)
		- UVP实现底层硬件的虚拟化
	1. VNA （Virtualization Node Agent ）
		- VNA实现对接VRM，向上提供管理接口
- [**华为虚拟化平台使用与管理**](https://juejin.cn/post/7184252320920797244)

	FusionCompute基于KVM
- [FusionCompute零星笔记](https://zhangfupeng.com/1860.html)

	FusionCompute采用的是硬件辅助全虚拟化.
- [[云计算]HCIE-Cloud 云计算运营 - 华为云计算解决方案](https://cnblogs.com/Skybiubiu/p/14618384.html)
- [HCIA-FusionCompute虚拟化产品](https://www.cnblogs.com/tjane/p/16845648.html)

	默认账号
- [**HCIA-Cloud Computing**](https://space.bilibili.com/630399494/channel/collectiondetail?sid=548989)
- [华为云计算HCIA V5.0](https://space.bilibili.com/412127397/channel/collectiondetail?sid=1479257)
- [FusionCompute8.0部署环境](https://space.bilibili.com/341704369)

## VRM/CNA
- db: VRM使用gaussdb; CNA未使用db

## vm
vm cdrom底层可用nbd, 避免拷贝iso. 对比过vrm和cna的fs(`df -h`)变化, 应该不是全量拷贝到某个节点在做nbd, 而是先拷贝部分, 等到vm读cdrom时, 再通过nbd+websocket读取所需部分. 因为关掉操作光驱的弹窗, 再用dd取代nbd设备报`Input/output error`, 缺点: 安装过程慢.

### disk
配置模式:
- 普通延迟置零: raw, disk size初始是创建时的大小

	cbt第一次备份大小是virtual size
- 普通: raw, disk size初始是创建时的大小

	cbt第一次备份大小是virtual size
- 精简: raw, disk size初始是4K

	cbt第一次备份大小是virtual size

## 备份
ref:
- [云祺容灾备份系统「华为FusionCompute备份」实操演示](https://www.bilibili.com/video/BV1cm4y1K7yb)
- [FusionComputeGolangSDK](https://github.com/KubeOperator/FusionComputeGolangSDK)
- [针对FusionComputer场景数据保护，AnyBackup有何妙招](https://www.aishu.cn/cn/blog/60)

flow:
- [华为FusionCompute](https://storware.gitbook.io/backup-and-recovery/protecting-virtual-machines/virtual-machines/huawei-fusion-compute)

	**type 'normal' when the VM is without VMTools or type 'CBTbackup' if the VM has VMTools installed**
- [eBackup有哪几种备份组网方式，各备份组网方式主要的应用场景及备份流程？](https://blog.51cto.com/u_15069486/4057651)
- [华为虚拟化备份方案](https://www.hcieonline.com/03-%E5%A4%87%E4%BB%BD/)
- [eSDK Cloud V100R005C00 开发指南(FusionSphere Backup&Restore SDK)](https://download.huawei.com/edownload/e/download.do?actionFlag=download&nid=EDOC1000067596&partNo=h001&mid=SUPE_DOC&_t=1501982801000)

底层是普通延迟置零, 首次创建cbt备份后, disk.source变为qcow2 compat 1.1, backingStore.source应该是原先的磁盘.

### `BCManager和FusionCompute Scoket接口`使用
1. 打CBT快照前, 不能存在prepareBackupResource
1. 第一次cbt备份, getCbtDiffBitmap返回的bitNum和bit1Num均为0
1. version找不到指定值, 用0
1. sequence: 单调递增即可, server返回同样的sequence
1. 读远端磁盘文件的返回: header+ployload(包含了file data)
1. crc用: ChecksumIEEE
1. 一个backupresource只能prepareBackupResource一次, 即不支持并发

	如果getBackupResource返回的targetFile即是cbt全备时需要下载的文件(xml中的backingStore.source), 同时其格式是raw(by qemu-img即页面上的"普通延迟置零"), 那么服务器上的该文件和下载到的该文件的md5相同
1. 如果数据全部读取完成, server断会自行关闭socket, 不用调用`关闭远端磁盘文件`. 想重新读取需要重新deleteBackupResource+延迟(删除BackupResource需要时间, 可用getBackupResource检查是否已删除)+prepareBackupResource
1. prepareBackupResource使用lanssl时, 只需要用tls.Dial()连接server即可

	可能的错误:
	- `x509: cannot validate certificate for xxx because it doesn't contain any IP SANs`: 因为server使用了自签名证书, 使用`tls.Config{InsecureSkipVerify: true}`即可
1. prepareBackupResource限制8个socket, 可以根据resp error(10300412/10300413)判断
1. 创建vm by VRM CreateVm:

 空闲cpu: CpuQuantity
 空闲内存: MemQuantityMB, 没有vrm管理页面没有显示可用的空闲内存, 参考CpuQuantity, 找到MemQuantityMB

## FAQ
### vm和host网段不通
vm使用端口组的VLAN应为0, 表示不使用VLAN标签, 再用nmtui配置ip即可

### vm iso启动黑屏
单纯慢, 多等待

### vm
- 挂载Tools: libvirt vm xml添加cdrom+vmtools-linux.iso相关配置

	ref:
	- [FusionCompute 安装，linux下安装vmtools报错Unsupported linux distribution](https://blog.csdn.net/csdnxiaohua/article/details/128832029)

	```bash
	# mkdir cdrom
	# mount /dev/sr0 cdrom
	# cp cdrom/vmtools-3.0.0.024.tar.bz2 .
	# tar -xf vmtools-3.0.0.024.tar.bz2 # 需要lbzip2, centos 8.1 官方qcow2没有该lib, `yum install bzip2`
	# cd vmtools
	# ./install
	# reboot
	```

	验证是否安装成功:
	1. `systemctl status vm-agent`
	1. 查看vm详情中的`Tools`状态

	centos 7.9成功, centos 8.1失败.
### ssh
ref:
- [华为FusionCompute-VRM密码重置](https://blog.csdn.net/sj349781478/article/details/122662166)

VRM/CNA不能直接用root登入, 需用其他账号比如gandalf, 再`su root`+`root密码`切换到root

### 端口
VRM:
	- web: 8443
	- https api: 7443

### API
- 查询CBT

	- 10430043: 创建CBT快照成功后, 该接口返回`"cbtflag": true`
- createVmSnapshot:

	- 10300026: 安装Tools. 安装Tools后, 在vm关机下也可执行创建CBT快照
- 查询虚拟机卷CBT差量位图

	- 10430041: 需要指定snapUuid

### 添加磁盘报`虚拟机磁盘存在未完成任务, 请任务结束后重试`
该vm存在"准备备份资源", 该vm的"任务和事件"里不存在该任务, 但`系统管理-任务与日志-任务中心`有该任务, 先结束即可

### vrm
- pg_hba.conf在/opt/gaussdb/data下. 登入db: `psql -U galax vrm`+`SingleLOUD!1`(不行的话试试ssh gandalf/root密码)
- 代码是java, 在/opt/galax/vrmportal/tomcat/webapps/ROOT/WEB-INF下