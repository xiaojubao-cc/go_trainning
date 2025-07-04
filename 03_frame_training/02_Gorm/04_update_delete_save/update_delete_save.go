package main

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/soft_delete"
	"log"
	"os"
	"strings"
	"time"
)

type Department struct {
	ID        int64                 `json:"ID"`
	Name      string                `json:"name"`
	ManagerId int64                 `json:"manager_id"` // 部门经理ID
	Phone     string                `json:"phone"`      // 联系电话
	Address   string                `json:"address"`    // 办公地址
	CreatedAt CustomTime            `json:"created_at"`
	IsDelete  soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
}

type Employee struct {
	ID           *int64     `json:"ID"`
	Name         string     `json:"name"` // 姓名
	Gender       int8       `json:"gender"`
	BirthDate    CustomTime `json:"birth_date"` // 出生日期
	Position     string     `json:"position"`   // 职位
	Salary       string     `json:"salary"`
	DepartmentId int8       `json:"department_id"` // 所属部门
	Department   Department `gorm:"foreignKey:DepartmentId;references:ID"`
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
		CreateBatchSize:                          1000,
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
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
	return db
}

func (employee *Employee) BeforeUpdate(tx *gorm.DB) (err error) {
	//检测数据库字段是否更新,如果执行更新前查询过该字段会被保存在Statement中Dest中
	/*db.Model(&User{ID: 1, Name: "jinzhu"}).Updates(map[string]interface{"name": "jinzhu2"})  Changed("Name") => true*/
	if tx.Statement.Changed("Name") {
		fmt.Printf("Name field is updated")
	}
	fmt.Println("execute BeforeUpdate")
	return nil
}
func main() {
	/*更新操作 save和update操作struct是否有ID*/
	//GormUpdateBySingleColumn()

	//GormUpdateByMultipleColumn()

	/*指定字段更新*/
	//GormUpdateDesignateColumn()

	/*获取更新的记录数*/
	//GormUpdateReturnRows()

	/*使用表达式更新*/
	//GormUpdateByGormExpression()

	/*使用主键删除*/
	GormDeleteByPrimaryKey()

	/*软删除使用Unscoped来物理删除和查询已经逻辑删除的数据*/
	GormDeleteBySoftDelete()
}

func GormDeleteBySoftDelete() {
	selectDb.Unscoped().Find(&Employee{}, 1)
	selectDb.Unscoped().Delete(&Employee{Salary: "50000"})
}

func GormDeleteByPrimaryKey() {
	statement := selectDb.Session(&gorm.Session{DryRun: true}).Delete(Employee{}, 1).Statement
	explain := statement.Dialector.Explain(statement.SQL.String(), statement.Vars...)
	fmt.Printf("delete sql is %s", explain)
	/*mysql不支持这种操作*/
	//var department Department
	//selectDb.Clauses(clause.Returning{Columns: []clause.Column{{Name: "name"}, {Name: "gender"}}}).Delete(&department, 3)
	//fmt.Printf("%+v", department)
}

func GormUpdateByGormExpression() {
	selectDb.Model(Employee{}).Where([]int{1}).Update("salary", gorm.Expr("salary + ?", 1000))
}

func GormUpdateReturnRows() {
	tx := selectDb.Model(Employee{}).Where([]int{1, 2}).Update("salary", "9000")
	fmt.Printf("update rows is %d", tx.RowsAffected)
}

func GormUpdateDesignateColumn() {
	/*配合这select可以更新struct的零值字段*/
	selectDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
		//return selectDb.Table("employees").Session(&gorm.Session{AllowGlobalUpdate: true}).Select("salary").Updates(map[string]interface{}{"department_id": 2, "salary": 8000})
		/*UPDATE `employees` SET `salary`='',`department_id`=2*/
		//return selectDb.Model(Employee{}).Select("salary", "department_id").Updates(Employee{DepartmentId: 2})
		return selectDb.Table("employees").Session(&gorm.Session{AllowGlobalUpdate: true}).Select("salary", "department_id").Updates(Employee{DepartmentId: 2})
	})
}

func GormUpdateByMultipleColumn() {
	/*Table不会调用钩子函数;updates方法要调用钩子函数*/
	//selectDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
	//	return selectDb.Table("employees").Where("id = ?", 1).Updates(map[string]interface{}{"department_id": 2, "salary": 8000})
	//	return selectDb.Table("employees").Where("id = ?", 1).UpdateColumns(map[string]interface{}{"department_id": 2, "salary": 8000})
	//})
	selectDb.Model(Employee{}).Where("id = ?", 1).Updates(map[string]interface{}{"department_id": 2, "salary": 8000})
}

func GormUpdateBySingleColumn() {
	/*Table不会调用钩子函数;update方法要调用钩子函数*/
	//selectDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
	//	return selectDb.Table("employees").Where("id = ?", 1).Update("salary", 5000)
	//	//return selectDb.Table("employees").Where("id = ?", 1).UpdateColumn("salary", 5000)
	//})
	selectDb.Model(Employee{}).Where("id = ?", 1).Update("salary", 1000)
}
