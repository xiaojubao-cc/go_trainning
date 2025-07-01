package main

import (
	"encoding/json"
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
		/*跳过默认的事务*/
		SkipDefaultTransaction: true,
		/*自动创建数据库外键约束*/
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Fatal("conn is fault")
		return nil
	}
	//db.AutoMigrate(&User{}, &CreditCard{})
	//db.AutoMigrate(&Employee{}, &Department{})
	db.AutoMigrate(&Order{}, &OrderItem{})
	return db
}

// User CreditCard
/*
	hasOne模型 user拥有creditCard,外键在creditCard,creditCard属于子模型
	父实体"拥有"子实体作为其扩展如
	需要级联删除关系
    主要从父实体访问子实体信息
	子实体是父实体的可选补充信息
*/
type User struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	Name       string
	Age        int
	Birthday   time.Time
	CreditCard CreditCard `gorm:"foreignKey:UserID;references:ID"`
}

type CreditCard struct {
	ID     uint `gorm:"primaryKey;autoIncrement"`
	Number string
	UserID uint
}

// Employee Department
/*
	belongTo模型 Employee属于Department,外键在Employee,Employee属于子模型
	子实体不能独立于父实体存在
	需要从子实体访问父实体信息 如：从评论(comment)访问文章(post)
	外键自然存在于子实体中
*/
type Employee struct {
	ID           uint
	Name         string
	Salary       float64
	DepartmentID uint
	Department   Department
}

type Department struct {
	ID   uint
	Name string
}

// Order OrderItem /*hasMany模型*/
type Order struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	OrderDate  time.Time
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;references:ID"`
}

type OrderItem struct {
	ID       uint `primaryKey`
	Quantity int
	OrderID  uint
}

func main() {
	/*单条记录插入*/
	//CreateSingleHasOneModel()
	/*多条记录插入*/
	//CreateMultipleHasOneModel()
	/*通过map插入数据*/
	//CreateUserByMap()
	/*创建连表插入*/
	//AssociationsCreate()
	/*插入belongTo模型*/
	//CreateBelongToModel()
	/*插入hasMany模型*/
	CreateHasManyModel()
}

func CreateHasManyModel() {
	//创建订单和订单明细
	tx := db.Begin()
	var order = Order{
		Name:      "Order1",
		OrderDate: time.Now(),
	}
	if err := tx.Create(&order); err != nil {
		tx.Rollback()
		return
	}
	var orderItems []OrderItem = make([]OrderItem, 0)
	orderItems = append(orderItems, OrderItem{
		OrderID:  order.ID,
		Quantity: 1,
	}, OrderItem{
		OrderID:  order.ID,
		Quantity: 2,
	})
	if err := tx.Model(&order).Association("OrderItems").Append(orderItems); err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	db.Preload("OrderItems").Find(&order)
	fmt.Printf("%+v", order)
}

func CreateBelongToModel() {
	var employees []Employee = make([]Employee, 0)
	tx := db.Begin()
	//先创建部门
	var departments = make([]Department, 0)
	departments = append(departments, Department{
		Name: "IT",
	}, Department{
		Name: "Manager",
	})
	if err := db.CreateInBatches(&departments, 2); err != nil {
		tx.Rollback()
	}
	employees = append(employees, Employee{
		Name:   "Jordan",
		Salary: 1000,
		Department: Department{
			ID:   departments[0].ID,
			Name: "IT",
		},
	}, Employee{
		Name:   "Jerry",
		Salary: 1000,
		Department: Department{
			ID:   departments[1].ID,
			Name: "Manager",
		},
	}, Employee{
		Name:   "Tom",
		Salary: 1000,
		Department: Department{
			ID:   departments[0].ID,
			Name: "IT",
		},
	}, Employee{
		Name:   "Mike",
		Salary: 1000,
		Department: Department{
			ID:   departments[1].ID,
			Name: "Manager",
		},
	}, Employee{
		Name:   "Lucy",
		Salary: 1000,
		Department: Department{
			ID:   departments[1].ID,
			Name: "Manager",
		},
	},
	)
	if err := db.CreateInBatches(employees, 2); err != nil {
		tx.Rollback()
	}
	tx.Commit()
	//belongTo模型查询、插入 先父后子
	db.Preload("Department").Find(&employees)
	for _, employee := range employees {
		marshal, _ := json.Marshal(&employee)
		fmt.Printf("CreateBelongToModel is %+v\n", string(marshal))
	}
}

func AssociationsCreate() {
	/*跳过指定关联结构体:db.Omit("CreditCard").Create(user)*/
	/*跳过所有与user关联的结构体:db.Omit(clause.Associations).Create(user)*/
	var user = &User{
		Name:     "Jordan",
		Age:      18,
		Birthday: time.Now(),
		CreditCard: CreditCard{
			Number: "513821199504206415",
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

func CreateSingleHasOneModel() {
	//hasOne模型单条数据插入
	var user = &User{
		Name:     "Tom",
		Age:      18,
		Birthday: time.Now(),
		CreditCard: CreditCard{
			Number: "513821199504206416",
		},
	}
	db.Create(user)
	//hasOne模型的查询先父后子
	db.Preload("CreditCard").Find(user)
	fmt.Printf("user is %+v\n", user)
}

func CreateMultipleHasOneModel() {
	//hasOne模型多条数据插入
	var users = make([]User, 0)
	users = append(users, User{
		Name:     "Jerry",
		Age:      18,
		Birthday: time.Now(),
		CreditCard: CreditCard{
			Number: "513821199504206417",
		},
	}, User{
		Name:     "Mike",
		Age:      18,
		Birthday: time.Now(),
		CreditCard: CreditCard{
			Number: "513821199504206418",
		},
	}, User{
		Name:     "Lucy",
		Age:      18,
		Birthday: time.Now(),
		CreditCard: CreditCard{
			Number: "513821199504206419",
		},
	}, User{
		Name:     "Lily",
		Age:      18,
		Birthday: time.Now(),
		CreditCard: CreditCard{
			Number: "513821199504206420",
		},
	})
	db.Session(&gorm.Session{CreateBatchSize: 2}).CreateInBatches(&users, 2)
	db.Preload("CreditCard").Find(&users)
}
