package main

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"
)

type Department struct {
	ID        int64      `json:"ID"`
	Name      string     `json:"name"`
	ManagerId int64      `json:"manager_id"` // 部门经理ID
	Phone     string     `json:"phone"`      // 联系电话
	Address   string     `json:"address"`    // 办公地址
	CreatedAt CustomTime `json:"created_at"`
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
		QueryFields: true,
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
	//db.AutoMigrate(&Employee{}, &Department{})
	return db
}
func main() {
	//GormSessionDryRun()
	GormToSQL()
}

// GormSessionDryRun /*获取完整的sql*/
func GormSessionDryRun() {
	/*用于生成sql*/
	statement := selectDb.Session(&gorm.Session{DryRun: true}).Table("employees").Where([]int{1, 2, 3}).Find(&Employee{}).Statement
	/*获得完整的执行sql*/
	sql := statement.Dialector.Explain(statement.SQL.String(), statement.Vars...)
	fmt.Println(sql)
}

// GormToSQL /*获取执行的sql*/
func GormToSQL() {
	var employees []Employee
	selectDb.FindInBatches(&employees, 2, func(tx *gorm.DB, batch int) error {
		sql := tx.ToSQL(func(in *gorm.DB) *gorm.DB {
			return in.Find(&employees)
		})
		fmt.Println(sql)
		return nil
	})
}

// GormTransaction 事务操作
func GormTransaction() {
	selectDb.Transaction(func(tx *gorm.DB) error {
		//嵌套事务可以回滚部分操作
		selectDb.Transaction(func(tx *gorm.DB) error {
			return nil
		})
		return nil
	})

	selectDb.Transaction(func(tx *gorm.DB) error {
		//所有操作在同一个事务
		return nil
	})
}

// GormConnection 无事务执行
func GormConnection() {
	selectDb.Connection(func(tx *gorm.DB) error {
		//可以执行查询修改等操作
		return nil
	})
}
