# lsblk
列出所有的块设备

## 选项

- -a : 显示所有设备
- -b : 以bytes方式显示设备大小
- -d : 不显示 slaves 或 holders
- -f : 显示文件系统信息
- -l : 以表格形式输出
- -n : 不输出表头(标题)
- -o : 输出指定的属性, 比如: name,maj:min,tran

	- ROTA : ROTA=1, HDD; ROTA=0, SSD
- -P : 使用key="value"格式显示
- -r : 使用原始格式显示
- -S/--scsi : 仅输出scsi设备
- -t : 显示拓扑结构信息
- --json : json输出