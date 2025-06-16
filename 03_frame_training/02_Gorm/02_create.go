package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

/*
	1.创建记录指定字段:db.Select("Name", "Age", "Birthday").Create(&User{})
	2.创建是忽略字段:db.Omit("Name", "Age", "Birthday").Create(&User{})
	3.批量插入:db.CreateInBatches([]*User{},1000) GORM 将生成一条 SQL 语句来插入所有数据并回填主键值，钩子方法也会被调用。当记录可以拆分为多个批次时，它将开始事务
	4.跳过钩子函数:db.Session(&gorm.Session{SkipHooks: true,}).Create(&User{});插入数据作为map创建时，不会调用钩子，不会保存关联，也不会回填主键值
	5.使用原生sql表达式gorm.Expr("NOW()")插入数据：db.Create(&User{Name: "Jinzhu", Age: 18, Birthday: gorm.Expr("NOW()")})

*/

const DATEFORMAT = "20060102150405"

var db *gorm.DB

func init() {
	db = basicConnection()
}
func basicConnection() *gorm.DB {
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
	})

	if err != nil {
		log.Fatal("conn is fault")
		return nil
	}
	return db
}

type User struct {
	Id         int
	Name       string
	Age        int
	Birthday   time.Time
	CreditCard CreditCard
}

type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
}

func main() {
	/*单条记录插入*/
	//CreateSingleUser()
	/*多条记录插入*/
	//CreateMultipleUser()
	/*通过map插入数据*/
	//CreateUserByMap()
	/*创建连表插入*/
	AssociationsCreate()
}

func AssociationsCreate() {
	/*跳过指定关联结构体:db.Omit("CreditCard").Create(user)*/
	/*跳过所有与user关联的结构体:db.Omit(clause.Associations).Create(user)*/
	user := &User{
		Name:     "Billi",
		Age:      18,
		Birthday: time.Now(),
		CreditCard: CreditCard{
			Number: "4111111111111111",
		},
	}
	db.Create(user)
}

func CreateUserByMap() {
	mapSlice := make([]map[string]interface{}, 0)
	mapSlice = append(mapSlice, map[string]interface{}{
		"Name":     "Tom",
		"Age":      18,
		"Birthday": time.Now(),
	}, map[string]interface{}{
		"Name":     "Jerry",
		"Age":      18,
		"Birthday": time.Now(),
	})
	/*取消默认批次*/
	db.Session(&gorm.Session{CreateBatchSize: 1})
	/*注意此处需要使用db.Model*/
	db.Model(&User{}).Create(mapSlice)
}

// BeforeCreate BeforeSave BeforeCreate AfterSave AfterCreate/*钩子函数*/
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("BeforeCreate")
	//set default value
	if user.Name == "" {
		user.Name = "default user"
	}
	return
}

func CreateSingleUser() {
	user := User{
		Name:     "汪建超",
		Age:      18,
		Birthday: time.Now(),
	}
	result := db.Create(&user)
	fmt.Printf("create user id is %d\n", user.Id)
	fmt.Printf("insert raw is %d\n", result.RowsAffected)
	fmt.Printf("insert error is %v\n", result.Error)
}

func CreateMultipleUser() {
	var users = make([]*User, 0)
	bir01, _ := time.Parse(DATEFORMAT, "20230921170000")
	bir02, _ := time.Parse(DATEFORMAT, "19940319170000")
	users = append(users, &User{Name: "汪梓榆", Age: 2, Birthday: bir01},
		&User{Name: "蔡娟", Age: 31, Birthday: bir02})
	result := db.Create(users)
	fmt.Printf("insert raw is %d\n", result.RowsAffected)
	fmt.Printf("insert error is %v\n", result.Error)
}
