# lsblk
列出所有的块设备

## 选项

- -a : 显示所有设备
- -b : 以bytes方式显示设备大小
- -d : 不显示 分区, slaves, holders
- -f/--fs : 显示文件系统信息
- -l : 以表格形式输出
- -n : 不输出表头(标题)
- -o : 输出指定的属性, 比如: name,maj:min,tran,PHY-SEC, LOG-SEC,...

	- ROTA : ROTA=1, HDD; ROTA=0, SSD
- O : 输出所有`-o`选项
- -P : 使用key="value"格式显示
- -r : 使用原始格式显示
- -S/--scsi : 仅输出scsi设备

	tran:
	- ata
	- sata
	- sas
	- fc

- -t : 显示拓扑结构信息
- --json : json输出

## example
```bash
$ lsblk -td # 显示磁盘的物理/逻辑扇区大小
```

## FAQ
### holders,slaves
sysfs 有 `/sys/block/$device/{holders,slaves}` 来表示哪些设备依赖于哪些设备, 比如lvm或软raid. 在lsof和fuser都找不到被占用问题时可尝试.