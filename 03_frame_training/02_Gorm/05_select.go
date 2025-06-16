package main

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	方法                 作用             是否需要结构体           是否支持关联操作             是否自动映射字段
Model(&Struct{})  绑定结构体与数据库表关系  ✅ 必须传入结构体指针      ✅ 支持关联操作             ✅ 自动映射字段
Table("table_name")   直接指定表名          ❌ 无需结构体          ❌ 不支持关联操作           ❌ 不自动映射字段

使用场景                 推荐方法               说明
查询完整模型数据            Find           支持关联、分页、零值过滤
聚合/分组查询              Scan              需配合 AS 别名
动态字段处理               Scan        使用 map[string]interface{}
关联查询                  Find             支持 Preload、Joins
Raw SQL 查询          Raw().Scan()       最灵活，但需手动维护 SQL
*/

type Department struct {
	Id        int64      `json:"id" gorm:"table:department"`
	Name      string     `json:"name"`
	ManagerId int64      `json:"manager_id"` // 部门经理ID
	Phone     string     `json:"phone"`      // 联系电话
	Address   string     `json:"address"`    // 办公地址
	CreatedAt CustomTime `json:"created_at"`
}

type Employee struct {
	Id           int64      `json:"id" gorm:"table:employee"`
	Name         string     `json:"name"` // 姓名
	Gender       int8       `json:"gender"`
	BirthDate    CustomTime `json:"birth_date"` // 出生日期
	Position     string     `json:"position"`   // 职位
	Salary       string     `json:"salary"`
	DepartmentId int64      `json:"department_id"` // 所属部门
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

func (employee *Employee) AfterFind(tx *gorm.DB) (err error) {

	return nil
}

func main() {
	//CreateMultipleEmployee()
	//SelectSingleEmployeeToStruct()
	//SelectSingleEmployeeToMap()
	//如果主键是数字
	//SelectSingleOrMultiplyEmployeeByNumberPrimaryKey()
	//如果主键是非数字
	//SelectSingleOrMultiplyEmployeeByNotNumberPrimaryKey()
	//如果结构体有默认值
	//SelectSingleOrMultiplyEmployeeByStructDefaultValue()
	//SelectSingleOrMultiplyEmployeeByWhereStringCondition()
	//SelectSingleOrMultiplyEmployeeByWhereMapOrStructOrSliceCondition()
	//SelectSingleOrMultiplyEmployeeByPage()
	//SelectFieldGroupByTable()
	//SelectByRawAndScan()
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
	selectDb.First(&employee, 1)
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
	/*Table不能使用连表Joins*/
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
