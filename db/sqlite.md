# sqlite

## sqlite3命令
```
.mode column # 按row显示数据
.mode line # 每个属性一行, 推荐
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

显示create table语句:
```
$ sqlite3 xxx.db3
> .schema <table>
```

## FAQ
### database is locked
sqlite3只支持一写多读.

sqlite同一时间只能进行一个写操作，当同时有两个写操作的时候,后执行的只能先等待,如果等待时间超过5秒,就会产生这种错误. 同样一个文件正在写入,重复打开数据库操作更容易导致这种问题的发生.

### datetime, date类型
sqlite3默认没有datetime和date类型, 但[SQLite支持列的亲和类型](https://www.runoob.com/sqlite/sqlite-data-types.html), 比如定义时是datetime, 实际保持时是亲和类型NUMERIC, 但无法通过`col > ${unixstamp}`之类的数值比较进行操作.

操作:
```sql
INSERT INTO loginfo(username, created) values(?,datetime(?, 'unixepoch')) // "xxx", time.Now().Unix()
SELECT * FROM loginfo where created > datetime(?, 'unixepoch'); // time.Now().Unix()
```

### json支持
```bash
go build --tags "sqlite_json1"  github.com/mattn/go-sqlite3
# t.go

type Product struct {
	gorm.Model
	Code   string
	Price  uint
	Author Author `gorm:"type:JSON;serializer:json"`
}

type Author struct {
	Name  string
	Email string
}
```

> 看到过sqlite3用`blocks TEXT NOT NULL`存储json内容

### `REPLACE INTO`
`INSERT OR REPLACE INTO mtime_cache (path, mtime, mtime_nsec, size, blocks) VALUES (?,?,?,?,?)`