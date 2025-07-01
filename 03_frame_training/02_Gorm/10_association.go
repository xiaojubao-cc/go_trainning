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

/*
特性  	       值类型 (Order)	           指针类型 (Product)
内存占用	        直接存储数据	               存储内存地址
空值表示       零值结构体 (所有字段为零值)	nil (明确表示不存在)
修改行为	       方法内修改不影响原始值	    方法内修改会影响原始对象
数据库关联	   适合作为"拥有者"实体	    适合作为"被引用"实体
JSON 序列化	     总是输出所有字段	          可省略 nil 字段
适用场景	      核心实体、总是存在的对象	  可选关联、可能不存在的对象
*/
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
	//db.AutoMigrate(&People{}, &IdentityCard{})
	//db.AutoMigrate(&Order{}, &OrderItem{}, &Product{})
	db.AutoMigrate(&Human{}, &Language{})
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

type Order struct {
	ID         uint
	OrderItems []*OrderItem `gorm:"foreignKey:OrderID;references:ID"`
}

func (Order) TableName() string {
	return "`order`"
}

type OrderItem struct {
	ID      uint
	OrderID uint
	Product *Product `gorm:"foreignKey:OrderItemID;references:ID"`
}

type Product struct {
	ID          uint
	Name        string
	OrderItemID uint
}

type Human struct {
	ID        uint
	Name      string
	Languages []*Language `gorm:"many2many:human_languages;"`
}

type Language struct {
	ID     uint
	Name   string
	Humans []*Human `gorm:"many2many:human_languages;"`
}

func main() {
	/*成员属于公司,公司作为成员属性foreignKey,reference,primaryKey标签的使用*/
	//GormAssociationBelongsTo()
	/*has one一对一模型person has one CreditCard,CreditCard对象是person的属性,personId作为CreditCard属性*/
	//GormAssociationHasOne()
	/*一对多*/
	//GormAssociationHasMany()
}

func GormAssociationHasMany() {
	//var OrderItems []*OrderItem
	//OrderItems = append(OrderItems, &OrderItem{
	//	ID:      1,
	//	OrderID: 1,
	//	Product: &Product{
	//		ID:          1,
	//		Name:        "product1",
	//		OrderItemID: 1,
	//	},
	//}, &OrderItem{
	//	ID:      2,
	//	OrderID: 1,
	//	Product: &Product{
	//		ID:   2,
	//		Name: "product2",
	//	},
	//})
	//var orders []Order
	//orders = append(orders, Order{
	//	ID:         1,
	//	OrderItems: OrderItems,
	//})
	//associationDb.Create(&orders)
	var orders []Order
	associationDb.Preload("OrderItems"). // 加载 OrderItems
						Preload("OrderItems.Product"). // 加载 OrderItems 的 Product
						Find(&orders)
	for _, order := range orders {
		marshal, _ := json.Marshal(&order)
		fmt.Printf("GormAssociationHasMany is %+v\n", string(marshal))
	}
}

func GormAssociationHasOne() {
	var peoples []People
	//peoples = append(peoples, People{
	//	ID:   1,
	//	Name: "Tom",
	//	IdentityCard: &IdentityCard{
	//		Number:   "123456789",
	//		PeopleID: 1,
	//	},
	//}, People{
	//	ID:   2,
	//	Name: "Jerry",
	//	IdentityCard: &IdentityCard{
	//		Number:   "123456788",
	//		PeopleID: 2,
	//	},
	//})
	//associationDb.Create(&peoples)
	//associationDb.Preload("IdentityCard").Where(&People{ID: 1}).Find(&peoples)
	associationDb.Find(&peoples)
	for _, people := range peoples {
		marshal, _ := json.Marshal(&people)
		fmt.Printf("GormAssociationHasOne is %+v\n", string(marshal))
	}
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
