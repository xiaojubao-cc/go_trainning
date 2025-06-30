package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

type User struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Posts []Post `gorm:"foreignKey:UserID"` // 1对多关联
}

type Post struct {
	ID       uint `gorm:"primaryKey"`
	Title    string
	UserID   uint      // 外键
	Comments []Comment `gorm:"foreignKey:PostID"` // 1对多关联
}

type Comment struct {
	ID      uint `gorm:"primaryKey"`
	Content string
	PostID  uint    // 外键
	Replies []Reply `gorm:"foreignKey:CommentID"` // 1对多关联
}

type Reply struct {
	ID        uint `gorm:"primaryKey"`
	Content   string
	CommentID uint // 外键
}

var associationDb *gorm.DB

func init() {
	associationDb = associationDbConnection()
}
func associationDbConnection() *gorm.DB {
	// 自定义日志配置
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // data source name
		DefaultStringSize:         256,                                                                          // default size for string fields
		DisableDatetimePrecision:  true,                                                                         // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                                         // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                         // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                                        // auto configure based on currently MySQL version
	}), &gorm.Config{
		//全局配置设置批量插入单个批次数量 /*这里是针对本次session设置批次*/db.Session(&gorm.Session{CreateBatchSize: 1000,})
		CreateBatchSize: 1000,
		Logger:          newLogger,
		/*跳过默认的事务*/
		SkipDefaultTransaction: true,
		/*禁止数据库外键约束在struct使用foreignKey标签时只会连表查询*/
		DisableForeignKeyConstraintWhenMigrating: true,
		/*命名规则*/
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatal("conn is fault")
		return nil
	}
	//db.AutoMigrate(&Member{}, &Company{})
	db.AutoMigrate(&User{}, &Post{}, &Comment{}, &Reply{})
	return db
}
func main() {
	// 构建完整嵌套结构
	newUser := []User{
		{Name: "Alice",
			Posts: []Post{
				{
					Title: "GORM Guide",
					Comments: []Comment{
						{
							Content: "Great post!",
							Replies: []Reply{
								{Content: "Thank you!"}, // 嵌套回复
							},
						},
						{
							Content: "Looking forward to more",
						},
					},
				},
			},
		},
		{
			Name: "Bob",
			Posts: []Post{
				{
					Title: "GORM Guide",
					Comments: []Comment{
						{
							Content: "Great post!",
						},
					},
				},
			},
		},
	}

	// 关键步骤：指定嵌套创建路径
	result := associationDb.Select("Posts", "Posts.Comments", "Posts.Comments.Replies").Create(&newUser)
	if result.Error != nil {
		panic(result.Error)
	}
}
