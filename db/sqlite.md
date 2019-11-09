# sqlite

## sqlite3命令
```
.mode column # 按row显示数据
.mode line # 每个属性一行
.header ON|OFF # 是否显示表头
```

export:
```sh
$ sqlite3 xxx.db3
> .output xxx.sql # 导出文件名
>.dump [table_name...] # 默认导出全库
> .q
```
　　

import:
```sh
$ sqlite3 xxx.db3
> .read xxx.sql
> .q
```

## FAQ
### database is locked
sqlite3只支持一写多读.

sqlite同一时间只能进行一个写操作，当同时有两个写操作的时候,后执行的只能先等待,如果等待时间超过5秒,就会产生这种错误. 同样一个文件正在写入,重复打开数据库操作更容易导致这种问题的发生.