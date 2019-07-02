# orm使用
gorm join:
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