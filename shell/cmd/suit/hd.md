# hd
## 获取有关硬件方面的信息

**DMI (Desktop Management Interface)**就是帮助收集电脑系统信息的管理系统，DMI信息的收集必须在严格遵照SMBIOS规范的前提下进行。

[Dmidecode命令详解](http://www.ha97.com/4120.html)

## FAQ
### mpt3sas_cm0: log_info(0x3112010a): originator(PL), code(0x12), sub_code(0x010a)
参考:
- [mpt2sas0: log_info(0x31120303) 问题分析与解决](https://blog.csdn.net/weixin_44648216/article/details/104070284)
- [hotplug_sata_drive.md](https://github.com/huataihuang/cloud-atlas-draft/blob/master/os/linux/kernel/storage/hotplug_sata_drive.md)
- [mpt2sas故障处理](https://huataihuang.gitbooks.io/cloud-atlas/content/storage/das/mpt2sas/troubleshooting/mpt2sas_offline_fail_disk.html)
- [使用命令行工具对LSI阵列卡进行高效管理](https://blog.51cto.com/1130739/1771506)

> mpt3sas is the driver for the SATA host bus adapter.

log_info是一个U32长度的变量, 后面是对它的解释, 定义在[`IOC LOGINFO defines`](https://elixir.bootlin.com/linux/v5.10.10/source/drivers/message/fusion/lsi/mpi_log_sas.h#L20)

log_info逻辑在[_base_sas_log_info](https://elixir.bootlin.com/linux/v5.10.10/source/drivers/scsi/mpt3sas/mpt3sas_base.c#L1230).

originator在:
```c
// https://elixir.bootlin.com/linux/v5.10.10/source/drivers/scsi/mpt3sas/mpt3sas_base.c#L1257
	switch (sas_loginfo.dw.originator) {
	case 0:
		originator_str = "IOP";
		break;
	case 1:
		originator_str = "PL";
		break;
	case 2:
		if (!ioc->hide_ir_msg)
			originator_str = "IR";
		else
			originator_str = "WarpDrive";
		break;
	}
```

[originator](https://elixir.bootlin.com/linux/v5.10.10/source/drivers/message/fusion/lsi/mpi_log_sas.h#L31):
0x0 对应 IOP 意思为IO Processor
0x1 对应 PL 意思为 Protocol Layer
0x2 对应 IR 意思为 Intergrated RAID

`code(0x12)的`定义是[PL_LOGINFO_CODE_ABORT](https://elixir.bootlin.com/linux/v5.10.10/source/drivers/message/fusion/lsi/mpi_log_sas.h#L137)

`sub_code(0x010a)`的定义是[PL_LOGINFO_SUB_CODE_<XXX>](https://elixir.bootlin.com/linux/v5.10.10/source/drivers/message/fusion/lsi/mpi_log_sas.h#L144)

sub_code(0x010a)=`PL_LOGINFO_SUB_CODE_OPEN_FAILURE|PL_LOGINFO_SUB_CODE_OPEN_FAIL_BREAK`

other example:
```log
0x31110d01 == (
MPI_IOCLOGINFO_TYPE_SAS |
IOC_LOGINFO_ORIGINATOR_PL |
PL_LOGINFO_CODE_RESET |
PL_LOGINFO_SUB_CODE_SATA_LINK_DOWN |
PL_LOGINFO_SUB_CODE_OPEN_FAIL_NO_DEST_TIME_OUT
)
```

> /var/log/messages日志中不断出现hard reset，表明存储卡出现了异常即服务器RAID卡需要替换维修.

> 存储hba卡通常是来自LSI公司（LSI Corporation）,  一般地，支持RAID 5的卡，称其为阵列卡，都可以使用LSI官方提供的MegaCli, SAS2IRCU, SAS3IRCU等工具来管理; 而不支持RAID 5的卡，称其为SAS卡，使用lsiutil工具来管理. HP的服务器则使用其特有的hpacucli工具来管理.
