# vinchin
ref
- [虚拟/物理机备份中深度有效数据提取应用原理](https://blog.csdn.net/qq_42934452/article/details/133942908)

    原理: 识别出文件系统层的有效数据位图

## FAQ
### [db账号](https://github.com/Chocapikk/CVE-2024-22899-to-22903-ExploitChain/blob/main/exploit.py)
ref:
- [VINCHIN BACKUP & RECOVERY V4.0 Quick Installation Guide Install on XenServer Virtual Server](https://www.vinchin.com/en/res/pdf/Quick_Installation_Guide_for_Citrix_XenServer%20_Vinchin_Backup_&_Recovery_v4.0.pdf)

config dir(/etc/vinchin):
- database_script/vinchin_db.sql => vinchin_db
- pt_server.conf.xml : 被编码了
- /etc/backup_system_common_server.conf.xml: 同上被编码了


```
self.db_user = "vinchin"
self.db_password = "yunqi123456"
self.db_name = "vinchin_db"
```

### log
1. 存储设备里:`<storage_uuid>/task_log`
1. `/var/log/vinchin`

### 虚拟机备份
> 增量和差异是互斥的

vm disk使用qcow2保存, 没用backing_chains

> 增量/差分场景都可基于backing_chains进行.

#### 瞬间恢复
创建disk.img通过NFS(`/mnt/kvinfs`)导出, 在fc host的数据存储里创建NAS存储, 再使用该共享img启动vm