package main

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"
)

/*
原生sql和sql生成器:
1.原生查询:Raw()...Scan()
2.原生执行:db.Exec()
3.sql测试不会执行:db.Session(&gorm.Session{DryRun: true})
4.用于生成sql:db.ToSQL(func)
5.row和rows结合scan
*/
var nativeDb *gorm.DB

type CustomOriginalTime time.Time

// Value /*实现value方法匹配数据库字段类型*/
func (ct *CustomOriginalTime) Value() (driver.Value, error) {
	t := time.Time(*ct)
	return t, nil
}

// Scan /*实现Scan方法匹配自定义时间格式*/
func (ct *CustomOriginalTime) Scan(v interface{}) error {
	switch val := v.(type) {
	case time.Time:
		*ct = CustomOriginalTime(val)
	case []byte:
		parse, err := time.Parse(time.DateTime, string(val))
		if err != nil {
			return err
		}
		*ct = CustomOriginalTime(parse)
	case string:
		parse, err := time.Parse(time.DateTime, val)
		if err != nil {
			return err
		}
		*ct = CustomOriginalTime(parse)
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}
func (ct *CustomOriginalTime) MarshalJSON() ([]byte, error) {
	output := fmt.Sprintf("\"%s\"", time.Time(*ct).Format("2006-01-02 15:04:05"))
	return []byte(output), nil
}

func (ct *CustomOriginalTime) UnmarshalJSON(data []byte) error {
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	t1, err := time.Parse(time.DateTime, timeStr)
	*ct = CustomOriginalTime(t1)
	return err
}

type Permissions struct {
	ID           uint               `gorm:"primaryKey;autoIncrement"`
	Name         string             `gorm:"type:varchar(50)"`
	ResourceType string             `gorm:"type:varchar(100)"`
	Action       string             `gorm:"type:varchar(100)"`
	CreatedAt    CustomOriginalTime `gorm:"column:create_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt    CustomOriginalTime `gorm:"column:update_at;default:CURRENT_TIMESTAMP"`
}

func (Permissions) TableName() string {
	return "permissions"
}

func init() {
	nativeDb = nativeBasicConnection()
}
func nativeBasicConnection() *gorm.DB {
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
	db.AutoMigrate(&Permissions{})
	return db
}
func main() {
	/*展示执行SQL*/
	//GormToSQL()
	/*无事务下一个连接中执行多个sql*/
	//GormConnection()
	GormSessionNewDBs()
}

func GormToSQL() {
	//SELECT * FROM `permissions` ORDER BY `permissions`.`id` LIMIT 2
	var permissions []Permissions
	nativeDb.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(Permissions{}).FindInBatches(&permissions, 2, func(tx *gorm.DB, batch int) error {
			for _, permission := range permissions {
				marshal, _ := json.Marshal(&permission)
				fmt.Printf("ToSQL is %+v\n", string(marshal))
			}
			return nil
		})
	})
	//SELECT * FROM `permissions`
	nativeDb.Session(&gorm.Session{DryRun: true}).Model(Permissions{}).Rows()
}

func GormConnection() {
	nativeDb.Connection(func(tx *gorm.DB) error {
		/*在无事务的情况下执行多个语句*/
		var permissions []Permissions = []Permissions{{
			Name:         "user:list",
			ResourceType: "user",
			Action:       "get",
		}, {
			Name:         "user:create",
			ResourceType: "user",
			Action:       "post",
		}}
		tx.Create(&permissions)
		/*ID回填*/
		tx.Find(&permissions)
		for _, permission := range permissions {
			marshal, _ := json.Marshal(&permission)
			fmt.Printf("GormConnection is %+v\n", string(marshal))
		}
		return nil
	})
}

func GormClause() {
	/*
		clause 是一个强大的 SQL 构建工具包，用于生成复杂的 SQL 子句（如 WHERE, ORDER BY, LIMIT, JOIN, FOR UPDATE 等）
		clause 适用于需要精细控制 SQL 生成的场景，例如动态查询、多表关联、锁机制等。
	*/
	nativeDb.Clauses(clause.Locking{Strength: "UPDATE", Table: clause.Table{Name: clause.CurrentTable}})
}

func GormSessionNewDBs() {
	var dataMap []map[string]interface{}
	nativeDb.Session(&gorm.Session{DryRun: true}).Table("employee").Scan(&dataMap)
	for _, data := range dataMap {
		fmt.Printf("GormSessionNewDB is %+v\n", data)
	}
	nativeDb.Session(&gorm.Session{NewDB: true}).Table("permissions").Scan(&dataMap)
	for _, data := range dataMap {
		fmt.Printf("GormSessionNewDB is %+v\n", data)
	}
}
