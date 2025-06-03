# gorm
ref:
- [Go 语言数据库开发利器：Gorm 框架指南](https://zhuanlan.zhihu.com/p/1908102185616646147)

Gorm 是一个用 Go 语言编写的全功能 ORM 库.

与其他 Go 语言的 ORM 框架相比，Gorm 具有以下优势：

1. 与 SQLx 对比

    SQLx 更接近原生 SQL，需要开发者编写更多 SQL 语句
    Gorm 提供更高级的抽象，减少 SQL 编写，更适合快速开发
    Gorm 的自动迁移和关联关系管理是 SQLx 所不具备的

2. 与 XORM 对比

    XORM 功能丰富但 API 相对复杂
    Gorm 的 API 设计更简洁直观，学习曲线更平缓
    Gorm 的社区活跃度和文档质量高且支持中、英、法等11种语言

3. 与 Ent 对比

    Ent 是 Facebook 开发的图数据库 ORM，专注于图数据模型
    Gorm 更通用，适用于传统关系型数据库
    Gorm 的 API 更符合 Go 语言习惯，而 Ent 的 API 风格更接近 GraphQL


## log
ref:
- [Monitoring Gin and GORM with OpenTelemetry](https://dev.to/vmihailenco/monitoring-gin-and-gorm-with-opentelemetry-53o0)

## preload
Preload是gorm特有的方法，用于自动加载关联字段（belongTo、many2many、hasOne、hasMany）中的数据, 即允许使用 Preload通过多个SQL中来直接加载关系.
Joins 属于SQL语法中的一部分，使用First、Find等方法时，会生成对应的SQL语句。一个Jions方法只支持一对一(hasOne, belongTo)关系和一对多关系（hasMany），多对多关系需要两个Joins方法（因为要多关联一个中间表）.

GORM 的 Preload() 是用来预加载关联数据的，比如有一对多或多对一关系。它会：
1. 先发出主查询（Find 主表）
2. 接着根据主表结果生成子查询来加载关联数据

这些 SQL 是分开执行的，但sql日志顺序是反的

```go
package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Child struct {
	gorm.Model
	ChildName string `gorm:"column:child_name;unique"`
	Toys      []Toy  `gorm:"foreignKey:ChildId"` //关联字段 hasMany
}

func (Child) TableName() string {
	return "child"
}

func (o *Child) AfterFind(tx *gorm.DB) (err error) {
	fmt.Println("--- find Child")
	return
}

type Toy struct {
	gorm.Model
	Name    string `gorm:"column:name"`
	ChildId uint   `gorm:"column:child_id"`
	Child   *Child `gorm:"foreignKey:ChildId"` //关联字段 belongTo
}

func (Toy) TableName() string {
	return "toy"
}

func (o *Toy) AfterFind(tx *gorm.DB) (err error) {
	fmt.Println("--- find Toy")
	return
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		//QueryFields: true, // 显示所有字段
	})
	if err != nil {
		panic("failed to connect database:" + err.Error())
	}

	// InsertData(db)
	// Preload(db)
	Join(db)
}

func Join(db *gorm.DB) {
	var toys11 []Toy
	// 注意where中的表别名一定要加双引号
	/*
	   SELECT "toy"."id","toy"."created_at","toy"."updated_at","toy"."deleted_at","toy"."child_id","toy"."name","Child"."id" AS "Child__id","Child"."created_at" AS "Child__created_at","Child"."updated_at" AS "Child__updated_at","Child"."deleted_at" AS "Child__deleted_at","Child"."child_name" AS "Child__child_name" FROM "toy" LEFT JOIN "child" "Child" ON "toy"."child_id" = "Child"."id" AND "Child"."deleted_at" IS NULL WHERE "Child".child_name = '刘涛' AND "toy"."deleted_at" IS NULL
	*/
	db.Joins("Child").Where("\"Child\".child_name = ?", "刘涛").Find(&toys11) // Find()对应的是主表

	spew.Dump(toys11)

	//也可以写成这样，这里是不需要双引号的
	/*
	  SELECT `toy`.`id`,`toy`.`created_at`,`toy`.`updated_at`,`toy`.`deleted_at`,`toy`.`name`,`toy`.`child_id`,`Child`.`id` AS `Child__id`,`Child`.`created_at` AS `Child__created_at`,`Child`.`updated_at` AS `Child__updated_at`,`Child`.`deleted_at` AS `Child__deleted_at`,`Child`.`child_name` AS `Child__child_name` FROM `toy` LEFT JOIN `child` `Child` ON `toy`.`child_id` = `Child`.`id` AND `Child`.`deleted_at` IS NULL WHERE `Child`.`child_name` = "刘涛" AND `toy`.`deleted_at` IS NULL
	*/
	var toys12 []Toy
	db.Debug().Joins("Child").Clauses(clause.Eq{
		Column: "Child.child_name",
		Value:  "刘涛",
	}).Find(&toys12)

	spew.Dump(toys12)

	//查询拥有‘汽车’玩具的孩子们
	/*
		SELECT `child`.`id`,`child`.`created_at`,`child`.`updated_at`,`child`.`deleted_at`,`child`.`child_name`,`Toys`.`id` AS `Toys__id`,`Toys`.`created_at` AS `Toys__created_at`,`Toys`.`updated_at` AS `Toys__updated_at`,`Toys`.`deleted_at` AS `Toys__deleted_at`,`Toys`.`name` AS `Toys__name`,`Toys`.`child_id` AS `Toys__child_id` FROM `child` LEFT JOIN `toy` `Toys` ON `child`.`id` = `Toys`.`child_id` AND `Toys`.`deleted_at` IS NULL WHERE `Toys`.`name` LIKE "%汽车%" AND `child`.`deleted_at` IS NULL
	*/
	children := []Child{}
	db.Joins("Toys").Clauses(clause.Like{
		Column: "Toys.name",
		Value:  "%" + "汽车" + "%",
	}).Find(&children)

	spew.Dump(children)

	// --- 上面方式生成的sql均存在: 右表的`deleted_at` IS NULL, 比如"LEFT JOIN `child` `Child` ON `toy`.`child_id` = `Child`.`id` AND `Child`.`deleted_at` IS NULL", 而下面没有该条件

	// 使Joins支持一对多关系的查询
	// SELECT `t1`.`id`,`t1`.`created_at`,`t1`.`updated_at`,`t1`.`deleted_at`,`t1`.`child_name` FROM child AS t1 LEFT JOIN toy AS t2 on t1.id =t2.child_id  WHERE t2.name = "纸飞机" AND `t1`.`deleted_at` IS NULL
	var children2 []Child
	db.Table(fmt.Sprintf("%v AS t1", Child{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %v AS t2 on t1.id =t2.child_id ", Toy{}.TableName())).
		Where("t2.name = ?", "纸飞机").Find(&children2)

	spew.Dump(children2)

	/*
		SELECT `t1`.`id`,`t1`.`created_at`,`t1`.`updated_at`,`t1`.`deleted_at`,`t1`.`child_name` FROM child AS t1 LEFT JOIN toy AS t2 on t1.id =t2.child_id  WHERE t2.name = "纸飞机" AND `t1`.`deleted_at` IS NULL
		SELECT * FROM `toy` WHERE `toy`.`child_id` = 1 AND `toy`.`deleted_at` IS NULL
	*/
	var children3 []Child
	db.Preload("Toys").
		Table(fmt.Sprintf("%v AS t1", Child{}.TableName())).
		Joins(fmt.Sprintf("LEFT JOIN %v AS t2 on t1.id =t2.child_id ", Toy{}.TableName())).
		Where("t2.name = ?", "纸飞机").Find(&children3)

	spew.Dump(children3)
}

func Preload(db *gorm.DB) {
	children := []Child{}
	/*
		SELECT * FROM `child` WHERE `child`.`deleted_at` IS NULL
		SELECT * FROM `child` WHERE `child`.`deleted_at` IS NULL
	*/
	//db.Debug().Preload("Toys", db.Where(&Toy{Name: "纸飞机"})).Find(&children)
	// 同上
	db.Preload("Toys", "name = ?", "纸飞机").Find(&children)
	spew.Dump(children)

	// 每个人只展示一个玩具
	// 下面代码是错误的, 不支持对 Preload 设置 LIMIT, 因此王斌没有关联到Toys
	/*
		SELECT * FROM `child` WHERE `child`.`deleted_at` IS NULL
		SELECT * FROM `toy` WHERE `toy`.`child_id` IN (1,2) AND `toy`.`deleted_at` IS NULL LIMIT 1
	*/
	var children2 []Child
	limit1 := func(db *gorm.DB) *gorm.DB { return db.Limit(1) }
	db.Preload("Toys", limit1).Find(&children2)
	spew.Dump(children2)

	// 正确方法
	var children21 []Child
	db.Find(&children21)
	// 对每个孩子查一个玩具（可用 JOIN / 或者批量一次性查询）
	for i := range children21 {
		var toy Toy
		db.Where("child_id = ?", children21[i].ID).Limit(1).Find(&toy)
		children21[i].Toys = []Toy{toy}
	}
	spew.Dump(children21)
}

func InsertData(db *gorm.DB) {
	var children = []Child{
		{
			ChildName: "刘涛",
			Toys: []Toy{
				{Name: "纸飞机"},
				{Name: "小火车"},
			},
		},
		{
			ChildName: "王斌",
			Toys: []Toy{
				{Name: "玩具兵"},
			},
		},
	}

	db.AutoMigrate(&Toy{})   // CREATE TABLE `child` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`child_name` text,CONSTRAINT `uni_child_child_name` UNIQUE (`child_name`))
	db.AutoMigrate(&Child{}) // CREATE TABLE `toy` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`name` text,`child_id` integer,CONSTRAINT `fk_child_toys` FOREIGN KEY (`child_id`) REFERENCES `child`(`id`))

	// toy的child_id怎么提前生成了??? 应该是log打印的问题
	// INSERT INTO `toy` (`created_at`,`updated_at`,`deleted_at`,`name`,`child_id`) VALUES ("2025-06-03 21:59:08.548","2025-06-03 21:59:08.548",NULL,"纸飞机",1),("2025-06-03 21:59:08.548","2025-06-03 21:59:08.548",NULL,"小火车",1),("2025-06-03 21:59:08.548","2025-06-03 21:59:08.548",NULL,"玩具兵",2) ON CONFLICT (`id`) DO UPDATE SET `child_id`=`excluded`.`child_id` RETURNING `id`
	// INSERT INTO `child` (`created_at`,`updated_at`,`deleted_at`,`child_name`) VALUES ("2025-06-03 21:59:08.547","2025-06-03 21:59:08.547",NULL,"刘涛"),("2025-06-03 21:59:08.547","2025-06-03 21:59:08.547",NULL,"王斌") RETURNING `id`
	if err := db.Create(&children).Error; err != nil {
		panic(err)
	}
}
```