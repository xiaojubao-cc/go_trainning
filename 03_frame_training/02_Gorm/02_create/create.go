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
		FullSaveAssociations:                     true,
	})

	if err != nil {
		log.Fatal("conn is fault")
		return nil
	}
	//db.AutoMigrate(&User{}, &CreditCard{})
	//db.AutoMigrate(&Employee{}, &Department{})
	//db.AutoMigrate(&Order{}, &OrderItem{}, &Product{})
	db.AutoMigrate(&Student{}, &Course{}, &Enrollment{})
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
	ID        uint `gorm:"primaryKey"`
	Quantity  int
	OrderID   uint
	ProductID uint
	Product   Product `gorm:"foreignKey:ProductID;references:ID"`
}

type Product struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Price float64
}

// Student /*多对多模型*/
type Student struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	CreatedAt time.Time
	Courses   []Course `gorm:"many2many:enrollments;"` // 多对多关系
}
type Course struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null;index"`
	Description string
	StartDate   time.Time
	Students    []Student `gorm:"many2many:enrollments;"` // 多对多关系
}

// Enrollment 选课关系模型（连接表，含额外字段）
type Enrollment struct {
	StudentID  uint      `gorm:"primaryKey"` // 复合主键
	CourseID   uint      `gorm:"primaryKey"` // 复合主键
	EnrolledAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Grade      float64
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
	//CreateHasManyModel()
	/*插入many2many模型*/
	//CreateMany2ManyModel()
	/*查询many2many*/
	QueryMany2ManyModel()
}

type StudentCourseGrade struct {
	StudentName string
	CourseTitle string
	Grade       float64
}

func QueryMany2ManyModel() {
	//queryStudentWithCourses(1)
	//queryCourseStudents("数学基础")
	queryAllGrades()
}

func queryAllGrades() {
	var studentCourseGrades []StudentCourseGrade = make([]StudentCourseGrade, 0)
	db.Table("courses").Select("courses.title as course_title,enrollments.grade,students.name as student_name").
		Joins("join enrollments on enrollments.course_id = courses.id").
		Joins("join students on students.id = enrollments.student_id").Scan(&studentCourseGrades)
	for _, studentCourseGrade := range studentCourseGrades {
		fmt.Printf("studentCourseGrade:%+v\n", studentCourseGrade)
	}
}

func queryCourseStudents(title string) {
	var courses []Course
	db.Table("courses").Select("courses.id,courses.title,courses.description,courses.start_date").
		Joins("left join enrollments on enrollments.course_id = courses.id").
		Joins("left join students on students.id = enrollments.student_id").
		Where("courses.title = ?", title).Scan(&courses)
	for _, course := range courses {
		fmt.Printf("Course:%+v\n", course)
	}
}
func queryStudentWithCourses(studentID uint) {
	var students []Student = make([]Student, 0)
	db.Preload("Courses").Find(&students, studentID)
	fmt.Printf("student:%+v", students)
}
func CreateMany2ManyModel() {
	// 创建课程
	courses := []Course{
		{Title: "数学基础", Description: "基础数学课程", StartDate: time.Now()},
		{Title: "编程入门", Description: "Python编程基础", StartDate: time.Now().AddDate(0, 1, 0)},
		{Title: "数据结构", Description: "算法与数据结构", StartDate: time.Now().AddDate(0, 2, 0)},
	}
	db.Create(&courses)

	// 创建学生
	students := []Student{
		{Name: "张三", Email: "zhangsan@example.com"},
		{Name: "李四", Email: "lisi@example.com"},
		{Name: "王五", Email: "wangwu@example.com"},
	}
	db.Create(&students)

	// 建立关联（学生选课）
	tx := db.Begin()

	// 方法1：使用Association添加关联
	mathCourse := courses[0]
	programmingCourse := courses[1]

	//Enrollment中的外键字段填充
	if err := tx.Model(&students[0]).Association("Courses").Append(&mathCourse, &programmingCourse); err != nil {
		tx.Rollback()
		panic("关联添加失败")
	}
	//Enrollment分数需要单独进行添加
	if err := tx.Model(&Enrollment{}).
		Where("student_id = ? AND course_id = ?", students[0].ID, mathCourse.ID).
		Update("grade", 92.5).Error; err != nil {
		tx.Rollback()
		panic("更新分数失败")
	}

	if err := tx.Model(&Enrollment{}).
		Where("student_id = ? AND course_id = ?", students[0].ID, programmingCourse.ID).
		Update("grade", 88.0).Error; err != nil {
		tx.Rollback()
		panic("更新分数失败")
	}
	tx.Commit()
}

func CreateHasManyModel() {
	//创建订单和订单明细
	//tx := db.Begin()
	//tx = tx.Debug()
	//var orders []Order = make([]Order, 0)
	//orders = append(orders, Order{
	//	Name:      "Order1",
	//	OrderDate: time.Now(),
	//}, Order{
	//	Name:      "Order2",
	//	OrderDate: time.Now(),
	//})
	//if err := tx.Create(&orders).Error; err != nil {
	//	tx.Rollback()
	//}
	////创建product
	//var products []Product = make([]Product, 0)
	//products = append(products, Product{
	//	Name:  "茶杯",
	//	Price: 100.00,
	//}, Product{
	//	Name:  "电脑",
	//	Price: 200.00,
	//})
	//if err := tx.Create(&products).Error; err != nil {
	//	tx.Rollback()
	//}
	//
	////创建订单明细
	//var orderItems []OrderItem = make([]OrderItem, 0)
	//orderItems = append(orderItems, OrderItem{
	//	OrderID:   orders[0].ID,
	//	Quantity:  1,
	//	ProductID: products[0].ID,
	//}, OrderItem{
	//	OrderID:   orders[1].ID,
	//	Quantity:  2,
	//	ProductID: products[1].ID,
	//})
	//
	//if err := tx.Create(&orderItems).Error; err != nil {
	//	tx.Rollback()
	//}
	//tx.Commit()
	//查询数据
	var Order []Order
	//db.Preload("OrderItems").Preload("OrderItems.Product").Find(&Order)
	db.Joins("JOIN order_items ON order_items.order_id = orders.id").
		Joins("JOIN products ON products.id = order_items.product_id").
		Where("products.price > ?", 100).
		//SELECT `orders`.`id`,`orders`.`name`,`orders`.`order_date` FROM `orders` JOIN order_items ON order_items.order_id = orders.id JOIN products ON products.id = order_items.product_id WHERE products.price > 100
		//这里为了填充嵌套字段OrderItems.Product和OrderItems数据
		Preload("OrderItems.Product").
		Find(&Order)
	for _, order := range Order {
		jsonData, _ := json.MarshalIndent(order, "", "  ")
		fmt.Println("Created order:", string(jsonData))
	}
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
