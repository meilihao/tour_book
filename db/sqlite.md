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