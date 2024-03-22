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
- [华为FusionCompute详解（二）FusionCompute总体介绍以及规划部署](https://blog.csdn.net/qq_46254436/article/details/105810195)
- [华为FusionCompute详解（一）FusionSphere虚拟化套件介绍](https://blog.csdn.net/qq_46254436/article/details/105807057)
- [**华为FusionCompute 8.x部署文档**](https://www.bilibili.com/read/cv19164631/)

## VRM/CNA
- db: VRM使用gaussdb; CNA未使用db

## vm
vm cdrom底层可用nbd, 避免拷贝iso. 对比过vrm和cna的fs(`df -h`)变化, 应该不是全量拷贝到某个节点在做nbd, 而是先拷贝部分, 等到vm读cdrom时, 再通过nbd+websocket读取所需部分. 因为关掉操作光驱的弹窗, 再用dd取代nbd设备报`Input/output error`, 缺点: 安装过程慢.

### disk
配置模式:
- 普通延迟置零: raw, disk size初始是创建时的大小

	根据磁盘容量为磁盘分配空间，创建时不会擦除物理设备上保留的任何数据，但后续从虚拟机首次执行写操作时会按需要将其置零. 创建速度比“普通”模式快；IO性能介于“普通”和“精简”两种模式之间.

	cbt第一次备份大小是virtual size
- 普通: raw, disk size初始是创建时的大小

	根据磁盘容量为磁盘分配空间，在创建过程中会将物理设备上保留的数据置零. 这种格式的磁盘性能要优于其他两种磁盘格式，但创建这种格式的磁盘所需的时间可能会比创建其他类型的磁盘长，且预留空间和实际占用空间相等，建议系统盘使用该模式.

	cbt第一次备份大小是virtual size
- 精简: raw, disk size初始是4K

	该模式下，系统首次仅分配磁盘容量配置值的部分容量，后续根据使用情况，逐步进行分配，直到分配总量达到磁盘容量配置值为止.

	cbt第一次备份大小是virtual size

注意:
1. "iscsi"挂入CNA的盘不能作为vm系统盘

## 备份
ref:
- [云祺容灾备份系统「华为FusionCompute备份」实操演示](https://www.bilibili.com/video/BV1cm4y1K7yb)
- [sdk: go-fusion-compute](https://github.com/LawyZheng/go-fusion-compute)
- [sdk: FusionComputeGolangSDK](https://github.com/KubeOperator/FusionComputeGolangSDK)
- [sdk: eSDK_FC_SDK_Java](https://github.com/jacklongway/eSDK_FC_SDK_Java)
- [针对FusionComputer场景数据保护，AnyBackup有何妙招](https://www.aishu.cn/cn/blog/60)
- [VMware 虚拟化编程(10) — VMware 数据块修改跟踪技术 CBT](https://www.cnblogs.com/jmilkfan-fanguiju/p/10589802.html)

flow:
- [华为FusionCompute](https://storware.gitbook.io/backup-and-recovery/protecting-virtual-machines/virtual-machines/huawei-fusion-compute)

	**type 'normal' when the VM is without VMTools or type 'CBTbackup' if the VM has VMTools installed**
- [eBackup有哪几种备份组网方式，各备份组网方式主要的应用场景及备份流程？](https://blog.51cto.com/u_15069486/4057651)
- [华为虚拟化备份方案](https://www.hcieonline.com/03-%E5%A4%87%E4%BB%BD/)
- [eSDK Cloud V100R005C00 开发指南(FusionSphere Backup&Restore SDK)](https://download.huawei.com/edownload/e/download.do?actionFlag=download&nid=EDOC1000067596&partNo=h001&mid=SUPE_DOC&_t=1501982801000)
- [华为发布全闪备份一体机旗舰新品，并宣布备份软件开源](https://www.163.com/dy/article/IMDEST6Q0511B8LM.html)

	开源时间线:
	- [open-eBackup：共建数据备份新生态](https://we.yesky.com/blog/309938)

底层是普通延迟置零, 首次创建cbt备份后, disk.source(应该就是disk的增量)变为qcow2 compat 1.1, backingStore.source应该是原先的磁盘.

### CBT
ref:
- [**华为云计算学习：备份之CBT技术**](https://blog.csdn.net/yangshihuz/article/details/104600311)
- [什么是CBT](https://itnobita.com/archives/yun-ji-suan-bi-ji)

CBT需要内存位图和CBT文件:
- 内存位图

	以1bit大小的0和1记录虚拟机的每个块的变更信息进行登记，0表示未改变，1表示改变.

	每次备份后会被重置重置为0

- CBT文件

	用来记录不同备份点上的数据块的CBT版本号

	> CBT版本号: 用于标识备份次数相对应的序号(大小4字节)，用来标记可从第几次备份中找到数据

	每次全备后可重置版本号.

备份过程:
1. 第一次: 

	内存位图mb0全为0; CBT文件cbt0全为0

	备份: 全备+cb0
2. 第2次

	步骤:
	1. mb1=copy(mb0), 将mb0全置0变为mb2
	2. 根据mb1更新cbt0相应位置为1, 得到cbt1

	备份: 数据变化量+cbt1
3. ...

cbt大小预估: `(volGB * 2<<30)/(sectorOfBlock*512) * 4B`

> 观察到7/1/2G卷的cbt文件都是8.1M???, 应该是预留了空间存储其他数据.

> 7G盘全备getCbtDiffBitmap返回的bitNum=3584和bit1Num=2991, bitNum!=bit1Num, 推测qcow也有bitmap, 可以辅助判断哪些数据块有数据. 

### `BCManager和FusionCompute Scoket接口`使用
1. 打CBT快照前, 不能存在prepareBackupResource
1. 第一次cbt备份, getCbtDiffBitmap返回的bitNum和bit1Num均为0

	第二次cbt备份(上一次cbt快照未删除), getCbtDiffBitmap返回的bitNum和bit1Num均为0.

	原因: BlockNum不能为0

	> GetCbtDiffBitmapReq不要使用`, omitempty`, 容易报400; 也没必要使用SnapUuid, 实际测试SnapUuid没有效果
	> 备份时, prepareBackupResource使用的Snap必须存在
	> CNA重启好后CBT快照的ChgID会从1重新开始

	vinchin(by `截取VNA请求`)每次BlockNum的step是1024, 具体传参:
	- 全备: ChgID为空即"0"
	- 增量: ChgID为cbtbackup snap的pre ChgID

		备份内容是pre ChgID ~ cur ChgID的变化数据, 因此如果ChgID为cur ChgID, 返回的Bit1Num是0.

	有效即根据返回的bit1Num计算.

	GetCbtDiffBitmapReq可以先仅传Type和VolCBTCreateTime

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

	并发备份vm的两个盘, 它们备份资源返回的hostPort都是35001, 小盘备份完成时CNA会关闭35001导致大盘备份失败, 参考vinchin推测FusionCompute本身设计就是这样.

	vinchin备份时(多vm):
	- 快照模式=串行: 逐台vm逐个disk备份再删除snap都是串行操作.

		```go
		for _,vm:=range vms{
			createSnap
			scanVmDisk
			for _,d:=range vm.Disks{
				tranDisk
			}
			deleteSnap
			执行保留策略
		}
		```
	- 快照模式=并行: 创建快照并行, 逐台vm逐个disk备份再删除snap即都是串行操作

		```go
		wg:=sync.WaitGroup{}

		for _,vm:=range vms{
			wg.Add(1)

			go func(){
				defer wg.Done()

				createSnap
			}
		}

		wg.Wait()

		for _,vm:=range vms{
			scanVmDisk
			for _,d:=range vm.Disks{
				tranDisk
			}
			deleteSnap
			执行保留策略
		}
		```

	并发还原是没有问题的.

1. 创建vm by VRM CreateVm:

 空闲cpu: CpuQuantity
 空闲内存: MemQuantityMB, 没有vrm管理页面没有显示可用的空闲内存, 参考CpuQuantity, 找到MemQuantityMB
 vmName/portGroupName可以重复
 网卡的mac需要唯一, 空表示自动生成, 否则填具体mac

## FAQ
### [X-Auth-UserType](https://zhuanlan.zhihu.com/p/648587865)

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
- createVmSnapshot with CBTbackup:

	- 10300026: 安装Tools. 安装Tools后, 在vm关机下也可执行创建CBT快照
- 查询虚拟机卷CBT差量位图

	- 10430041: 需要指定snapUuid

### 添加磁盘报`虚拟机磁盘存在未完成任务, 请任务结束后重试`
该vm存在"准备备份资源", 该vm的"任务和事件"里不存在该任务, 但`系统管理-任务与日志-任务中心`有该任务, 先结束即可

### vrm
- pg_hba.conf在/opt/gaussdb/data下. 登入db: `psql -U galax vrm`+`SingleLOUD!1`(不行的话试试ssh gandalf/root密码)
- 代码是java, 在/opt/galax/vrmportal/tomcat/webapps/ROOT/WEB-INF下

### novnc
前端:
1. browser -> vrm返回vrm proxy url(wss)
	
	返回的端口是6084, 但wss链接时实际使用的是6083

	看vrm上websockify进程有个`--proxy-port 6084`可能与此有关, 官方websockify没有这个参数
2. browser打开vrm proxy url
3. vrm proxy(novnc) -> CNA nginx 5903 -> CNA qemu-kvm 5903(localhost)

	qemu-kvm vnc的授权方式是sasl, 参考[VNC 身份验证](https://documentation.suse.com/zh-cn/sles/15-SP2/html/SLES-all/cha-libvirt-connect.html#sec-libvirt-connect-auth-vnc), 可得实际由/etc/sasl2/qemu.conf配置里的saslauthd验证身份.

	再根据`saslauthd -m /run/saslauthd -a pam`的参数`pam`和vm xml的`graphics.authz+authz`推断应该是使用os用户鉴权.

	> `testsaslauthd -u hehe -p 123  [-f /run/saslauthd/mux]      //验证sasl工作是否正常`

	> novnc授权文件在`/opt/galax/vrmportal/tomcat/webapps/OmsPortal/novnc/utils/websockify`下, 但实际发现其授权里的vnc_username+vnc_port和CNA上的vm xml的authz.identity对不上???

### 截取VNA请求
```bash
# vim /usr/local/nginx/conf/nginx.conf
修改`/var/log/galaxenginlog/vna-nginx/nginx-access.log`对于的format, `log_format format_main escape=json '... $request_body'`, 参考[nginx记录post数据](https://cloud.tencent.com/developer/article/1501467)
# /usr/local/nginx/sbin/nginx -s reload
```