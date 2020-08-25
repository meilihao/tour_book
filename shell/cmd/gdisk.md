# gdisk
gdisk是支持gpt的分区工具.

> from gdisk_1.0.5-1_amd64

> cgdisk: 支持终端窗口功能的gdisk

## example
```bash
$ sudo  gdisk /dev/nbd0
Command (? for help): ? # help
Command (? for help): p # 打印分区表
Command (? for help): n # 新建分区, 此时默认已经创建了分区表
Partition number (1-128, default 1): 1 # 分区号
First sector (34-16777182, default = 2048) or {+-}size{KMGTP}: # 起始扇区
Last sector (2048-16777182, default = 16777182) or {+-}size{KMGTP}: +256M # 新分区的大小
Current type is 'Linux filesystem'
Hex code or GUID (L to show codes, Enter = 8300): # 默认即可之后再修改, 8300是用于格式化成ext4/xfs等linux文件系统用的
Command (? for help): p # 查看创建的分区, 及剩余的磁盘空间
Command (? for help): d # 删除分区
Partition number (1-2): 2 # 删除分区2
Command (? for help):  l # 查看支持的分区类型
Command (? for help): t 　　　  　 # 更改分区类型，这里输入l也可以查看分区的类型
Partition number (1-2): 2 　　　　 # 输入要更改的分区号
Current type is 'Linux filesystem'
Hex code or GUID (L to show codes, Enter = 8300): 8e00　　　　 # 输入分区类型的编号 
Command (? for help): c                              # 更改分区名称
Partition number (1-2): 2                            # 要更改的分区号
Enter name: pv1 LVM                                  # 更改后的名称
Command (? for help): w 　　　　　　　　 # 保存配置，如果不想保存可以输入q退出

Final checks complete. About to write GPT data. THIS WILL OVERWRITE EXISTING
PARTITIONS!!

Do you want to proceed? (Y/N): y # 询问是否相想继续，输入y继续
OK; writing new GUID partition table (GPT) to /dev/nbd0.
The operation has completed successfully.
```