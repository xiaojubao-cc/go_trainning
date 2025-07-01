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

type User struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	Email      string `gorm:"unique"`
	CreatedAt  time.Time
	CreditCard CreditCard `gorm:"foreignKey:UserID;references:ID"` // HasOne 关系(外键在关联表中)
	Profile    Profile    `gorm:"foreignKey:UserID;references:ID"` // HasOne 关系
	CompanyID  uint
	Company    Company    `gorm:"foreignKey:CompanyID;references:ID"` // BelongsTo 关系(外键在当前表中)
	Orders     []Order    `gorm:"foreignKey:UserID;references:ID"`    // HasMany 关系
	Subscribed []*Product `gorm:"many2many:user_subscriptions"`       // Many2Many 关系
}

type CreditCard struct {
	ID     uint   `gorm:"primaryKey"`
	Number string `gorm:"unique"`
	Expiry string
	UserID uint // 外键
}

type Profile struct {
	ID      uint `gorm:"primaryKey"`
	Bio     string
	Website string
	UserID  uint // 外键
}

type Company struct {
	ID      uint `gorm:"primaryKey"`
	Name    string
	Address string
}

type Order struct {
	ID         uint `gorm:"primaryKey"`
	OrderDate  time.Time
	UserID     uint        // 外键
	OrderItems []OrderItem `gorm:"foreignKey:OrderID"` // HasMany 关系
}

type OrderItem struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   uint // 外键
	ProductID uint // 外键
	Quantity  int
	Product   Product `gorm:"foreignKey:ProductID"` // BelongsTo 关系
}

type Product struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Price       float64
	Description string
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
	db.AutoMigrate(&User{}, &CreditCard{}, &Profile{}, &Company{}, &Order{}, &OrderItem{}, &Product{})
	return db
}
func main() {
	//创建product
	tx := associationDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("事务回滚:", r)
		}
	}()

	// 创建产品
	products := []Product{
		{Name: "MacBook Pro", Price: 1999.99, Description: "高性能笔记本电脑"},
		{Name: "iPhone 14", Price: 999.99, Description: "旗舰智能手机"},
		{Name: "AirPods Pro", Price: 249.99, Description: "降噪无线耳机"},
		{Name: "iPad Pro", Price: 799.99, Description: "专业平板电脑"},
	}
	if err := tx.Create(&products).Error; err != nil {
		tx.Rollback()
		log.Println("创建产品失败:", err)
		return
	}

	//创建公司
	company := Company{Name: "Apple Inc.", Address: "One Apple Park Way, Cupertino, CA 95014"}
	if err := tx.Create(&company).Error; err != nil {
		tx.Rollback()
		log.Println("创建公司失败:", err)
		return
	}

	//创建用户
	user1 := User{
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		CreditCard: CreditCard{
			Number: "4111111111111111",
			Expiry: "12/25",
		},
		Profile: Profile{
			Bio:     "技术爱好者，Apple产品粉丝",
			Website: "https://johndoe.dev",
		},
		Company: company,
		Orders: []Order{
			{
				OrderDate: time.Now(),
				OrderItems: []OrderItem{
					{ProductID: products[0].ID, Quantity: 1},
					{ProductID: products[2].ID, Quantity: 2},
				},
			},
			{
				OrderDate: time.Now().Add(-24 * time.Hour),
				OrderItems: []OrderItem{
					{ProductID: products[1].ID, Quantity: 1},
				},
			},
		},
		Subscribed: []*Product{&products[0], &products[3]},
	}
	// 使用Select精确控制关联创建
	if err := tx.Select(
		"CreditCard",
		"Profile",
		"Company",
		"Orders",
		"Orders.OrderItems",
		"Subscribed",
	).Create(&user1).Error; err != nil {
		tx.Rollback()
		log.Fatal("创建用户失败:", err)
	}

	// 创建另一个用户
	user2 := User{
		Name:       "Jane Smith",
		Email:      "jane@example.com",
		Company:    Company{Name: "Smith Enterprises"},
		Subscribed: []*Product{&products[1], &products[3]},
	}
	if err := tx.Select("Company", "Subscribed").Create(&user2).Error; err != nil {
		tx.Rollback()
		log.Fatal("创建用户2失败:", err)
	}

	tx.Commit()
	fmt.Println("示例数据创建完成")
	// 查询用户及其信用卡
	var userWithCard User
	associationDb.Preload("CreditCard").First(&userWithCard, "name = ?", "John Doe")
	printJSON("用户及信用卡:", userWithCard)

	// 查询用户及其档案
	var userWithProfile User
	associationDb.Preload("Profile").First(&userWithProfile, "email = ?", "john@example.com")
	printJSON("用户及档案:", userWithProfile)

	// 3.2 BelongsTo 关系查询
	fmt.Println("\n=== BelongsTo 关系查询 ===")

	// 查询用户所属公司
	var userWithCompany User
	associationDb.Preload("Company").First(&userWithCompany, "name = ?", "John Doe")
	printJSON("用户所属公司:", userWithCompany)

	// 查询订单项及其产品
	var orderItemWithProduct OrderItem
	associationDb.Preload("Product").First(&orderItemWithProduct)
	printJSON("订单项及产品:", orderItemWithProduct)

	// 3.3 HasMany 关系查询
	fmt.Println("\n=== HasMany 关系查询 ===")

	// 查询用户的所有订单
	var userWithOrders User
	associationDb.Preload("Orders").First(&userWithOrders, "name = ?", "John Doe")
	printJSON("用户的所有订单:", userWithOrders)

	// 查询订单的所有订单项
	var orderWithItems Order
	associationDb.Preload("OrderItems").First(&orderWithItems)
	printJSON("订单的所有订单项:", orderWithItems)

	// 深度预加载: 用户 + 订单 + 订单项 + 产品
	var userFull User
	associationDb.Preload("Orders.OrderItems.Product").First(&userFull, "name = ?", "John Doe")
	printJSON("用户完整订单信息:", userFull)

	// 3.4 Many2Many 关系查询
	fmt.Println("\n=== Many2Many 关系查询 ===")

	// 查询用户订阅的产品
	var userWithSubscriptions User
	associationDb.Preload("Subscribed").First(&userWithSubscriptions, "name = ?", "John Doe")
	printJSON("用户订阅的产品:", userWithSubscriptions)

	// 查询订阅某个产品的用户
	var product Product
	associationDb.Preload("Users").First(&product, "name = ?", "MacBook Pro")
	// 注意: 需要反向声明关系（在Product中添加 Users []User `gorm:"many2many:user_subscriptions"`）
	printJSON("订阅MacBook Pro的用户:", product)

	// 条件预加载
	var jane User
	associationDb.Preload("Subscribed", "price > ?", 500).First(&jane, "name = ?", "Jane Smith")
	printJSON("Jane订阅的高价产品:", jane)

	// 3.5 关联方法查询
	fmt.Println("\n=== 关联方法查询 ===")

	// 查找用户的所有订单
	var user User
	associationDb.First(&user, "name = ?", "John Doe")

	var orders []Order
	associationDb.Model(&user).Association("Orders").Find(&orders)
	printJSON("John的订单:", orders)

	// 统计用户订单数量
	orderCount := associationDb.Model(&user).Association("Orders").Count()
	fmt.Printf("John的订单数量: %d\n", orderCount)

	// 添加新订单
	newOrder := Order{
		OrderDate: time.Now(),
		OrderItems: []OrderItem{
			{ProductID: 4, Quantity: 1}, // iPad Pro
		},
	}
	associationDb.Model(&user).Association("Orders").Append(&newOrder)
}

// 辅助函数：打印JSON格式数据
func printJSON(prefix string, data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	fmt.Printf("%s:\n%s\n", prefix, jsonData)
}
