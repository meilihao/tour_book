# vinchin

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
不知Vinhin技术细节, 但备份到qcow2上应该挺合适, 尤其是增量/差分都可基于backing_chains进行.