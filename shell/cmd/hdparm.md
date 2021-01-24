# hdparm

## 描述

用来检测、显示及设定IDE或SCSI硬盘的参数, 也是测试硬盘读性能的常用工具

## 格式

## 选项
	-C : 检测IDE硬盘的电源管理模式, 比如standby

## 例

    # hdparm /dev/sda # 获取硬盘的设置
    # hhdparm -g /dev/sda -w # 显示硬盘的柱面、磁头、扇区数
    # hdparm --read-sector 824769880 /dev/sdc # 读取指定sector
    # smartctl -l selftest /dev/sda # 硬盘坏道修复方法
