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

## 备份
ref:
- [云祺容灾备份系统「华为FusionCompute备份」实操演示](https://www.bilibili.com/video/BV1cm4y1K7yb)
- [FusionComputeGolangSDK](https://github.com/KubeOperator/FusionComputeGolangSDK)
- [针对FusionComputer场景数据保护，AnyBackup有何妙招](https://www.aishu.cn/cn/blog/60)

flow:
- [华为FusionCompute](https://storware.gitbook.io/backup-and-recovery/protecting-virtual-machines/virtual-machines/huawei-fusion-compute)
- [eBackup有哪几种备份组网方式，各备份组网方式主要的应用场景及备份流程？](https://blog.51cto.com/u_15069486/4057651)
- [华为虚拟化备份方案](https://www.hcieonline.com/03-%E5%A4%87%E4%BB%BD/)
- [eSDK Cloud V100R005C00 开发指南(FusionSphere Backup&Restore SDK)](https://download.huawei.com/edownload/e/download.do?actionFlag=download&nid=EDOC1000067596&partNo=h001&mid=SUPE_DOC&_t=1501982801000)