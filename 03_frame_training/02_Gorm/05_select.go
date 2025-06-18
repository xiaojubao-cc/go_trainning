package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/hints"
	"log"
	"os"
	"strings"
	"time"
)

/*
	1.使用内联条件:db.Find(&employee,Employee{Id:1})
	2.NOT，OR条件和Where类似
	3.选择指定字段输出：db.Select("Name", "Age").Find(&employees)
	4.使用Order:db.Order("id desc").Find(&employees)
	5.使用distinct:db.Distinct("name", "age").Order("name, age desc").Find(&results)
	6.连接和派生:query := db.Table("order").Select("MAX(order.finished_at) as latest").Joins("left join user user on order.user_id = user.id").Where("user.age > ?", 18).Group("order.user_id")
	db.Model(&Order{}).Joins("join (?) q on order.finished_at = q.latest", query).Scan(&results)
	// SELECT `order`.`user_id`,`order`.`finished_at` FROM `order` join (SELECT MAX(order.finished_at) as latest FROM `order` left join user user on order.user_id = user.id WHERE user.age > 18 GROUP BY `order`.`user_id`) q on order.finished_at = q.latest



*/
/*
	方法                 作用             是否需要结构体           是否支持关联操作           是否自动映射字段       是否触发钩子函数
Model(&Struct{})  绑定结构体与数据库表关系  ✅ 必须传入结构体指针      ✅ 支持关联操作            ✅ 自动映射字段           会
Table("table_name")   直接指定表名          ❌ 无需结构体          ❌ 不支持关联操作          ❌ 不自动映射字段         不会

使用场景                 推荐方法               说明
查询完整模型数据            Find           支持关联、分页、零值过滤
聚合/分组查询              Scan              需配合 AS 别名
动态字段处理               Scan        使用 map[string]interface{}
关联查询                  Find             支持 Preload、Joins
Raw SQL 查询          Raw().Scan()       最灵活，但需手动维护 SQL
*/

type Department struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name"`
	ManagerId int64      `json:"manager_id"` // 部门经理ID
	Phone     string     `json:"phone"`      // 联系电话
	Address   string     `json:"address"`    // 办公地址
	CreatedAt CustomTime `json:"created_at"`
}

type Employee struct {
	Id           *int64     `json:"id"`
	Name         string     `json:"name"` // 姓名
	Gender       int8       `json:"gender"`
	BirthDate    CustomTime `json:"birth_date"` // 出生日期
	Position     string     `json:"position"`   // 职位
	Salary       string     `json:"salary"`
	DepartmentId int8       `json:"department_id"` // 所属部门
	HireDate     CustomTime `json:"hire_date"`
	Email        string     `json:"email"`
	CreatedAt    CustomTime `json:"created_at"`
	UpdatedAt    CustomTime `json:"updated_at"`
}
type CustomTime time.Time

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%v\"", time.Time(ct).Format(time.DateTime))
	return []byte(formatted), nil
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	t1, err := time.Parse(time.DateTime, timeStr)
	*ct = CustomTime(t1)
	return err
}

// Value /*Value方法转化为数据库支持的类型*/
func (ct CustomTime) Value() (driver.Value, error) {
	tTime := time.Time(ct)
	return tTime.Format(time.DateTime), nil
}

// Scan /*Scan转化为自定义类型*/
func (ct *CustomTime) Scan(v interface{}) error {
	switch val := v.(type) {
	case time.Time:
		*ct = CustomTime(val)
	case []byte:
		t, err := time.Parse(time.DateTime, string(val))
		if err != nil {
			return err
		}
		*ct = CustomTime(t)
	case string:
		t, err := time.Parse(time.DateTime, val)
		if err != nil {
			return err
		}
		*ct = CustomTime(t)
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}

// TableName 实现自定义表名
func (Employee) TableName() string {
	return "employee"
}
func (Department) TableName() string {
	return "department"
}

var selectDb *gorm.DB

func init() {
	selectDb = selectBasicConnection()
}
func selectBasicConnection() *gorm.DB {
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
		/*所有模型字段都按其名称进行选择select * from table*/
		//QueryFields: true,
		//设置表名
		//NamingStrategy: schema.NamingStrategy{
		//	TablePrefix:   "",   //表名前缀
		//	SingularTable: true, //是否使用单数表名
		//},
	})

	if err != nil {
		log.Fatal("conn is fault")
		return nil
	}
	//自动创建表
	db.AutoMigrate(&Employee{}, &Department{})
	return db
}

// AfterFind /*钩子函数处理查询后的数据*/
func (employee *Employee) AfterFind(tx *gorm.DB) (err error) {
	//针对数据库查询的数据进行处理返回
	return nil
}

func main() {
	//CreateMultipleEmployee()
	SelectSingleEmployeeToStruct()
	//SelectSingleEmployeeToMap()
	/*如果主键是数字*/
	//SelectSingleOrMultiplyEmployeeByNumberPrimaryKey()
	/*如果主键是非数字*/
	//SelectSingleOrMultiplyEmployeeByNotNumberPrimaryKey()
	/*如果结构体有默认值*/
	//SelectSingleOrMultiplyEmployeeByStructDefaultValue()
	/*字符串作为查询条件*/
	//SelectSingleOrMultiplyEmployeeByWhereStringCondition()
	/*struct,map,slice作为条件查询*/
	//SelectSingleOrMultiplyEmployeeByWhereMapOrStructOrSliceCondition()
	/*limit和offset*/
	//SelectSingleOrMultiplyEmployeeByPage()
	/*group操作*/
	//SelectFieldGroupByTable()
	/*原生sql查询*/
	//SelectByRawAndScan()
	/*映射字段*/
	//SelectAllFields()
	/*锁定改行所有操作*/
	//GormForUpdateLock()
	/*锁定某行其他事务可以查询该行*/
	//GormForShareLock()
	/*获取不到锁则不等待*/
	//GormForWaitOption()
	/*跳过事务正在处理的行*/
	//GormForSkipLock()
	/*多列条件查询*/
	//GormMultipleColumns()
	/*命名参数*/
	//GormChristenParam()
	/*如果未找到匹配的记录，则初始化新实例*/
	//GormFirstOrInit()
	/*如果未找到匹配的记录，则创建*/
	//GormFirstOrCreate()
	/*指定索引*/
	//GormHintIndex()
	/*适合查询*/
	//GormIterate()
	/*适合查询出的数据进行处理*/
	//GormFindInBatches()
	/*查询某列数据返回切片*/
	//GormPluck()
	/*范围应用*/
	//GormScopes()
	/*计数*/
	//GormCount()
}

func CreateMultipleEmployee() {
	employees := make([]*Employee, 0)
	bir01, _ := time.Parse(time.DateTime, "20230921170000")
	bir02, _ := time.Parse(time.DateTime, "19940319170000")
	employees = append(employees, &Employee{
		Name:         "James",
		Gender:       1,
		BirthDate:    CustomTime(bir01),
		Position:     "Developer",
		Salary:       "5000",
		DepartmentId: 3,
		HireDate:     CustomTime(time.Now()),
		Email:        "james@gmail.com",
		CreatedAt:    CustomTime(time.Now()),
		UpdatedAt:    CustomTime(time.Now()),
	}, &Employee{
		Name:         "Elizabeth",
		Gender:       2,
		BirthDate:    CustomTime(bir02),
		Position:     "Developer",
		Salary:       "5000",
		DepartmentId: 4,
		HireDate:     CustomTime(time.Now()),
		Email:        "elizabeth@gmail.com",
		CreatedAt:    CustomTime(time.Now()),
		UpdatedAt:    CustomTime(time.Now()),
	})
	selectDb.Session(&gorm.Session{CreateBatchSize: 1}).CreateInBatches(employees, 1)

}

func CreateMultipleDepartment() {
	departments := make([]*Department, 0)
	departments = append(departments, &Department{
		Name:      "Development",
		ManagerId: 1,
		Phone:     "1234567890",
		Address:   "123 Main St",
	}, &Department{
		Name:      "Marketing",
		ManagerId: 2,
		Phone:     "1234567890",
		Address:   "123 Main St",
	})
	selectDb.Session(&gorm.Session{CreateBatchSize: 1})
	selectDb.CreateInBatches(departments, 2)
}

func SelectSingleEmployeeToStruct() {
	var employee Employee
	selectDb.Session(&gorm.Session{QueryFields: true}).Select("name").First(&employee, 1)
	marshal, _ := json.Marshal(&employee)
	fmt.Printf("SelectSingleEmployeeToStruct is %+v", string(marshal))
}

func SelectSingleEmployeeToMap() {
	var employeeMap map[string]interface{}
	//使用map进行映射必须使用Model(结构体指针)
	//selectDb.Model(&Employee{}).First(&employeeMap)
	//这里不能使用First
	selectDb.Table("employee").Take(&employeeMap)
	fmt.Printf("SelectSingleEmployeeToMap is %+v", employeeMap)
}

func SelectSingleOrMultiplyEmployeeByNumberPrimaryKey() {
	//var employee Employee
	//selectDb.First(&employee, 1)
	//marshal, _ := json.Marshal(&employee)
	//fmt.Printf("singleEmployee is %+v", string(marshal))

	var employees []Employee
	selectDb.Find(&employees, []int{1, 2, 3})
	for _, employee := range employees {
		marshal, _ := json.Marshal(&employee)
		fmt.Printf("SelectSingleOrMultiplyEmployeeByNumberPrimaryKey is %+v\n", string(marshal))
	}
}

func SelectSingleOrMultiplyEmployeeByNotNumberPrimaryKey() {
	/*sql查询会拼接为and进行查询*/
	//var employee Employee
	//selectDb.First(&employee, "id = ?", 1)
	//marshal, _ := json.Marshal(&employee)
	//fmt.Printf("SelectSingleOrMultiplyEmployeeByNotNumberPrimaryKey is %+v", string(marshal))

	var employees []Employee
	selectDb.Find(&employees, "id in ?", []string{"1", "2", "3"})
	for _, employee := range employees {
		marshal, _ := json.Marshal(&employee)
		fmt.Printf("SelectSingleOrMultiplyEmployeeByNotNumberPrimaryKey is %+v\n", string(marshal))
	}
}

func SelectSingleOrMultiplyEmployeeByStructDefaultValue() {
	//struct写法
	var employees []Employee
	selectDb.Find(&employees, &Employee{Name: "Jerry"})
	marshal, _ := json.Marshal(&employees)
	fmt.Printf("SelectSingleOrMultiplyEmployeeByStructDefaultValue is %+v\n", string(marshal))
}

func SelectSingleOrMultiplyEmployeeByWhereStringCondition() {
	var employees []Employee
	where := "id in ? and name like ?"
	selectDb.Where(where, []int{1, 2, 3}, "%Tom%").Find(&employees)
	for _, employee := range employees {
		marshal, _ := json.Marshal(&employee)
		fmt.Printf("SelectSingleOrMultiplyEmployeeByNotNumberPrimaryKey is %+v\n", string(marshal))
	}
}

func SelectSingleOrMultiplyEmployeeByWhereMapOrStructOrSliceCondition() {
	/*map,slice不需要传入指针,但是struct需要传入指针*/
	var employees []Employee
	//mapWhere := make(map[string]interface{})
	//mapWhere["name"] = "Tom"
	//mapWhere["id"] = []int{1, 2, 3}
	/*使用结构体是需要注意非零值才会被作为条件,过滤掉零值*/
	//selectDb.Where(mapWhere).Find(&employees)
	employee := Employee{Name: "Jerry"}
	/*指定查询条件中的struct特定字段,此时就可以拼接零值作为条件 SELECT * FROM `employee` WHERE `employee`.`name` = 'Jerry' AND `employee`.`salary` = ''*/
	selectDb.Where(&employee, "name", "salary").Find(&employees)
	//ids := make([]int, 0)
	//ids = append(ids, 1, 3, 5, 7, 9)
	//selectDb.Where(ids).Find(&employees)
	for _, employee := range employees {
		marshal, _ := json.Marshal(&employee)
		fmt.Printf("SelectSingleOrMultiplyEmployeeByWhereMapOrStructOrSliceCondition is %+v\n", string(marshal))
	}
}

func SelectSingleOrMultiplyEmployeeByPage() {
	var employees []Employee
	selectDb.Limit(2).Offset(2).Find(&employees)
	for _, employee := range employees {
		marshal, _ := json.Marshal(&employee)
		fmt.Printf("SelectSingleOrMultiplyEmployeeByWhereMapOrStructOrSliceCondition is %+v\n", string(marshal))
	}
}

type EmployeeDTO struct {
	DepartmentName string `gorm:"column:DepartmentName"`
	Salary         string `gorm:"column:Salary"`
}

func SelectFieldGroupByModel() {
	var employeeDTO []EmployeeDTO
	tx := selectDb.Model(&Employee{}).Select("department.name as DepartmentName,sum(employee.salary) AS Salary").Joins("left join department on department.id = employee.department_id")
	tx.Group("employee.department_id").Find(&employeeDTO)
	for _, employee := range employeeDTO {
		marshal, _ := json.Marshal(&employee)
		fmt.Printf("SelectFieldByGroup is %+v\n", string(marshal))
	}
}

func SelectFieldGroupByTable() {
	/*Table不能使用连表Joins但是可以使用子查询*/
	dataMap := make([]map[string]interface{}, 0)
	//selectDb.Table("employee").Select("department_id,SUM(salary)").Group("department_id").Scan(&dataMap)
	//for k, v := range dataMap {
	//	fmt.Printf("SelectFieldGroupByTable is %+v:%+v\n", k, v)
	//}
	rows, err := selectDb.Table("employee").Select("department_id,SUM(salary)").Group("department_id").Rows()
	if err != nil {
		fmt.Printf("SelectFieldGroupByTable is %+v\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		var departmentID int64
		var totalSalary string
		if err := rows.Scan(&departmentID, &totalSalary); err != nil {
			log.Printf("failed to scan row: %v", err)
			continue
		}
		dataMap = append(dataMap, map[string]interface{}{
			"department_id": departmentID,
			"total_salary":  totalSalary,
		})
	}
	for k, v := range dataMap {
		fmt.Printf("SelectFieldGroupByTable is %+v:%+v\n", k, v)
	}
}

func SelectByRawAndScan() {
	/*使用原生的sql*/
	var employees []Employee
	selectDb.Raw("select * from employee").Scan(&employees)
	for _, employee := range employees {
		marshal, _ := json.Marshal(&employee)
		fmt.Printf("SelectByRawAndScan is %+v\n", string(marshal))
	}
}

func SelectAllFields() {
	var maps = make(map[string]interface{})
	/*配置中添加QueryFields:true以所有模型字段都按其名称进行选择*/
	selectDb.Session(&gorm.Session{QueryFields: true}).Model(Employee{}).Where("id = ?", 1).Select("name", "gender").Scan(&maps)
	fmt.Printf("SelectAllFields is %+v\n", maps)
	for _, employee := range maps {
		marshal, _ := json.Marshal(&employee)
		fmt.Printf("SelectAllFields is %+v\n", string(marshal))
	}
}

func GormForUpdateLock() {
	/*
		1.update:锁定操作行 SELECT * FROM `employee` WHERE `employee`.`id` = 1 FOR UPDATE
	*/
	i := int64(1)
	var employee Employee = Employee{Id: &i}
	selectDb.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&employee)
	marshal, _ := json.Marshal(&employee)
	fmt.Printf("GormForUpdateLock is %+v\n", marshal)
}

func GormForShareLock() {
	/*允许其他事务针对操作行进行查询 SELECT * FROM `employee` WHERE `employee`.`id` = 1 FOR SHARE OF `employee`*/
	i := int64(1)
	var employee Employee = Employee{Id: &i}
	/*Table: clause.Table{Name: clause.CurrentTable}:当使用 JOIN 查询多个表时，若仅需锁定主表的数据行，避免对关联表加锁*/
	selectDb.Clauses(clause.Locking{Strength: "SHARE", Table: clause.Table{Name: clause.CurrentTable}}).Find(&employee)
}

func GormForWaitOption() {
	/*
			Options:NOWAIT:立即返回
		    SELECT * FROM `employee` WHERE `employee`.`id` = 1 FOR UPDATE NOWAIT
	*/
	i := int64(1)
	var employee Employee = Employee{Id: &i}
	selectDb.Clauses(clause.Locking{Strength: "UPDATE", Options: "NOWAIT"}).Find(&employee)
}

func GormForSkipLock() {
	/*
			SKIP LOCKED:跳过已被其他事务锁定的任何行,高并发有用
		    SELECT * FROM `employee` FOR UPDATE SKIP LOCKED
	*/
	var employee []Employee
	selectDb.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).Find(&employee)
}

func GormSubQuery() {
	/*子查询有性能问题*/

}

func GormMultipleColumns() {
	var employees []Employee
	selectDb.Where("(name,department_id,gender) IN ?", [][]interface{}{{"Tom", 1, 1}, {"Jerry", 2, 2}}).Find(&employees)
	for _, employee := range employees {
		marshal, _ := json.Marshal(&employee)
		fmt.Printf("GormMultipleColumns is %+v\n", string(marshal))
	}
}

func GormChristenParam() {
	var employees []Employee
	selectDb.Where("name = @name", sql.Named("name", "Tom")).Or("department_id in @Ids", sql.Named("Ids", []int{1, 2, 3})).Find(&employees)
	for _, employee := range employees {
		marshal, _ := json.Marshal(&employee)
		fmt.Printf("GormChristenParam is %+v\n", string(marshal))
	}
}

func GormFirstOrInit() {
	/*如果未找到匹配的记录，则初始化新实例*/
	var employee Employee
	//selectDb.FirstOrInit(&employee, Employee{Name: "Oliver"})
	//selectDb.FirstOrInit(&employee, map[string]interface{}{"name": "Oliver",})
	/*Attrs未找到记录才进行初始化数据 GormFirstOrInit is {"id":0,"name":"Oliver","gender":0,"birth_date":"0001-01-01 00:00:00","position":"","salary":"80000","department_id":0,"hire_date":"0001-01-01 00:00:00","email":"","created_at":"0001-01-01 00:00:00","updated_at":"0001-01-01 00:00:00"}*/
	//selectDb.Where(Employee{Name: "Oliver"}).Attrs(Employee{Salary: "80000"}).FirstOrInit(&employee)
	/*Assign无论是否找到数据都进行初始化,但是不会保存到数据库 GormFirstOrInit is {"id":1,"name":"Harden","gender":1,"birth_date":"2023-09-22 01:00:00","position":"Developer","salary":"5000","department_id":1,"hire_date":"2025-06-16 15:11:10","email":"tom@gmail.com","created_at":"2025-06-16 15:11:10","updated_at":"2025-06-16 15:11:10"}*/
	selectDb.Where(Employee{Name: "Tom"}).Assign(Employee{Name: "Harden"}).FirstOrInit(&employee)
	marshal, _ := json.Marshal(&employee)
	fmt.Printf("GormFirstOrInit is %+v\n", string(marshal))
}

func GormFirstOrCreate() {
	/*如果没有找到匹配的记录，则创建一个新记录*/
	var employee Employee
	//selectDb.FirstOrCreate(&employee, Employee{Name: "Oliver"})
	//selectDb.FirstOrCreate(&employee, map[string]interface{}{"name": "Oliver",})
	/*Attrs为新记录指定其他属性会保存到数据库*/
	selectDb.Where(Employee{Name: "Oliver"}).Attrs(Employee{Salary: "80000"}).FirstOrCreate(&employee)
	/*Assign为新纪录指定其他属性会保存到数据库*/
	selectDb.Where(Employee{Name: "Oliver"}).Attrs(Employee{Email: "oliver@gmail.com"}).FirstOrCreate(&employee)
}

func GormHintIndex() {
	var employees []Employee
	/*指定索引*/
	//selectDb.Clauses(hints.UseIndex("index_name")).Find(&employee)
	/*Join操作指定索引*/
	selectDb.Clauses(hints.ForceIndex("index_department_id")).Joins("left join department on department.id = employee.department_id").Find(&employees)
	for _, employee := range employees {
		marshal, _ := json.Marshal(&employee)
		fmt.Printf("GormHintIndex is %+v\n", string(marshal))
	}
}

func GormIterate() {
	/*使用流式查询需要关闭流*/
	var employees []Employee
	rows, _ := selectDb.Model(Employee{}).Rows()
	defer rows.Close()
	for rows.Next() {
		var employee Employee
		selectDb.ScanRows(rows, &employee)
		employees = append(employees, employee)
	}
	for _, employee := range employees {
		marshal, _ := json.Marshal(&employee)
		fmt.Printf("GormIterate is %+v\n", string(marshal))
	}
}

func GormFindInBatches() {
	var employees []Employee
	selectDb.FindInBatches(&employees, 2, func(tx *gorm.DB, batch int) error {
		for _, employee := range employees {
			marshal, _ := json.Marshal(&employee)
			fmt.Printf("GormFindInBatches is %+v\n", string(marshal))
		}
		return nil
	})
}

func GormPluck() {
	/*将查询的某一列数据组装为slice返回*/
	var names []string
	//selectDb.Model(Employee{}).Pluck("name", &names)
	selectDb.Table("employee").Distinct().Pluck("name", &names)
	for _, name := range names {
		fmt.Printf("GormPluck is %+v\n", name)
	}
}

// GormScopesWithName /*需要传参的scopes需要返回函数*/
func GormScopesWithName(names []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name in ?", names)
	}
}

// GormScopesWithSalary /*不需要传参的直接返回对象*/
func GormScopesWithSalary(db *gorm.DB) *gorm.DB {
	return db.Where("salary > ?", "3000")
}
func GormScopes() {
	/*将公共部分进行抽取*/
	var employees []Employee
	selectDb.Scopes(GormScopesWithSalary, GormScopesWithName([]string{"Tom", "Jerry"})).Find(&employees)
	for _, employee := range employees {
		marshal, _ := json.Marshal(&employee)
		fmt.Printf("GormScopes is %+v\n", string(marshal))
	}
}

func GormCount() {
	var count int64
	//selectDb.Model(Employee{}).Count(&count)
	/*使用distinct计数*/
	//selectDb.Model(Employee{}).Distinct("name").Count(&count)
	//selectDb.Model(Employee{}).Select("count(distinct(name))").Count(&count)
	/*使用分组计数*/
	selectDb.Table("employee").Group("name").Count(&count)
	fmt.Printf("employee total_count is %d", count)
}
