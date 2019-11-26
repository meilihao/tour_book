# orm使用
**推荐使用xorm**.

## gorm
### join:
```go
type CommentaryItem struct {
	Commentary `gorm:"embedded"`
	ScenicName string `json:"scenic_name"`
}

// gorm join 无法使用 alias
func CommentaryList(tdb *gorm.DB, pg *pager.Pager) []CommentaryItem {
	ls := make([]CommentaryItem, 0)

	if err := tdb.Table(_TableCommentary).
		Select("count(*) OVER() as var_count, commentaries.*,scenic.name as scenic_name").
		Joins("join scenic on commentaries.scenic_id = scenic.id").
		Limit(pg.Size).Offset(pg.Offset()).Find(&ls).Error; err != nil {
		log.Error(err)

		return ls[:0]
	}

	if len(ls) > 0 {
		pg.SetTotal(ls[0].VarCount)
	}

	return ls
}
```

### 位运算
```go
db.Model(n).Update("status", gorm.Expr(fmt.Sprintf("status | %d", 1))).Error
db.Model(n).Update("status", gorm.Expr(fmt.Sprintf("status & (~%d)", 1))).Error
```

### sqlite3 localtime
```go
/*
CREATE TABLE IF NOT EXISTS loginfo(
       uid INTEGER PRIMARY KEY AUTOINCREMENT,
       created DATE NULL
   );*/

func main() {
	db, err := gorm.Open("sqlite3", "sqlite3.db?_loc=Local")
	checkErr(err)

	db.LogMode(true)

	db.Exec("insert into loginfo values(2,datetime(?, 'unixepoch'));", time.Now().Unix())

	n := &logInfo{}
	checkErr(db.Table("loginfo").First(n).Error)
	fmt.Println(n)
}

type logInfo struct {
	Id      uint64 `gorm:"Column:uid" json: "uid"`
	Created string `gorm:"Column:created" json: "created"`
}
```

sqlite3使用NUMERIC保存DATE/datetime的数据. 驱动[github.com/mattn/go-sqlite3是将数值转换成utc再设置时区的](https://github.com/mattn/go-sqlite3/blob/master/sqlite3.go#L2036), 因此存入的应该是unixstamp.

> `datetime('now','utc')`存入再用`?_loc=Local`取出时间相差不是8小时???.