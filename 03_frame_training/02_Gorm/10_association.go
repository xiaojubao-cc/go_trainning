package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

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
	db.AutoMigrate(&Member{}, &Company{})
	return db
}

type Member struct {
	gorm.Model
	Name      string
	CompanyID int
	Company   *Company `gorm:"foreignKey:ID;references:CompanyID"`
}

type Company struct {
	ID   int `gorm:"primaryKey"`
	Code string
	Name string
}

type People struct {
	ID           uint
	Name         string
	IdentityCard *IdentityCard `gorm:"foreignKey:PeopleID;references:ID"`
}

type IdentityCard struct {
	gorm.Model
	Number   string
	PeopleID uint
}

func main() {
	/*成员属于公司,公司作为成员属性foreignKey,reference,primaryKey标签的使用*/
	GormAssociationBelongsTo()
	/*has one一对一模型person has one CreditCard,CreditCard对象是person的属性,personId作为CreditCard属性*/
	GormAssociationHasOne()
}

func GormAssociationHasOne() {
	var peoples []People
	peoples = append(peoples, People{
		ID:   1,
		Name: "Tom",
		IdentityCard: &IdentityCard{
			Number:   "123456789",
			PeopleID: 1,
		},
	}, People{
		ID:   2,
		Name: "Jerry",
		IdentityCard: &IdentityCard{
			Number:   "123456788",
			PeopleID: 2,
		},
	})
}

func GormAssociationBelongsTo() {
	//var member []Member = []Member{
	//	{
	//		Name:      "Tom",
	//		CompanyID: 1,
	//		Company: Company{
	//			ID:   1,
	//			Code: "001",
	//			Name: "huawei",
	//		},
	//	},
	//	{
	//		Name:      "Jerry",
	//		CompanyID: 2,
	//		Company: Company{
	//			ID:   2,
	//			Code: "002",
	//			Name: "xiaomi",
	//		},
	//	},
	//
	//	{
	//		Name:      "Mike",
	//		CompanyID: 3,
	//		Company: Company{
	//			ID:   3,
	//			Code: "003",
	//			Name: "apple",
	//		},
	//	},
	//}
	//associationDb.Create(&member)
	//associationDb.Omit(clause.Associations).Create(&member)
	var member []Member
	associationDb.Preload("Company").Find(&member)
	for _, member := range member {
		marshal, _ := json.Marshal(&member)
		fmt.Printf("SelectAllFields is %+v\n", string(marshal))
	}
}
