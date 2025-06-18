package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

/*
	tx.Statement.Changed("name")注意事项:
	1.使用db.Model(&role{ID:1}).UpdateColumn("name", "Tom"),这里&role{ID:1}只传入了ID会从数据库查询数据和updates
      后面的数据进行比较;如果传入了name例如:&role{ID:1,Name:"Tom"}那么直接进行比较;
	2.updates方法使用结构体是是否触发tx.Statement.Changed("name")校验取决与struct的标签是否设置了`gorm:"column:name"`
    3.updateColumn和updateColumns方法不会调用钩子函数和自动更新时间
*/

type Role struct {
	ID          uint      `gorm:"primaryKey;autoIncrement;comment:ID" json:"id"`
	Name        string    `gorm:"column:name;not null;comment:名称" json:"name"`
	Description string    `gorm:"not null;comment:描述" json:"description"`
	CreateAt    time.Time `gorm:"column:create_at;comment:创建时间;default:CURRENT_TIMESTAMP" json:"create_at"`
	UpdateAt    time.Time `gorm:"column:update_at;comment:更新时间;default:CURRENT_TIMESTAMP" json:"update_at"`
}

func (*Role) TableName() string {
	return "role"
}

var saveDb *gorm.DB

func init() {
	saveDb = saveBasicConnection()
}
func saveBasicConnection() *gorm.DB {
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
	db.AutoMigrate(&Role{})
	return db
}

func main() {
	/*保存或者更新*/
	//GormSave()
	/*更新struct或者map非零字段的值*/
	GormDesignateFieldsUpdate()
	/*批量更新*/
	//GormBatchUpdate()
	/*全局更新*/
	//GormHandleGlobUpdate()
	/*通过表达式更新数据*/
	//GormUpdateByExpression()
	/*通过子查询进行更新*/
	//GormUpdateBySubQuery()
	/*返回修改之后的数据*/
	//GormUpdateRawReturnData()
}

func GormSave() {
	/*save操作如果携带主键则进行更新操作,否则进行新增操作*/
	var role Role = Role{
		Name:        "Jerry",
		Description: "jerry is exploitation",
	}
	saveDb.Select("name", "description").Save(&role)
}
func (role *Role) AfterUpdate(tx *gorm.DB) (err error) {
	fmt.Println("AfterUpdate")
	return nil
}

// BeforeUpdate /*检查字段是否已更改*/
func (role *Role) BeforeUpdate(tx *gorm.DB) (err error) {
	/*这里似乎只能在updates(map[string]interface{})*/
	//if tx.Statement.Changed("Name") {
	//	fmt.Println("name is not allow to update")
	//	return errors.New("name is not allow to update")
	//}
	//if tx.Statement.Changed("Description") {
	//	return errors.New("description is not allow to update")
	//}
	//if tx.Statement.Changed("CreateAt") {
	//	fmt.Println("create_at is not allow to update")
	//}
	return nil
}

// BeforeSave /*在数据落库前执行的操作*/
func (role *Role) BeforeSave(tx *gorm.DB) (err error) {
	//比如数据处理时间更新操作等
	//tx.Statement.SetColumn("name", "Tom", true)
	return nil
}
func GormDesignateFieldsUpdate() {
	//var role Role = Role{ID: 1}
	/*指定struct*/
	//saveDb.Model(&Role{ID: 1, Name: "Jerry"}).Updates(map[string]interface{}{"name": "Jerry"})
	//saveDb.Model(&Role{ID: 1, Name: "Jerry"}).Update("name", "Tom")
	/*指定map*/
	//saveDb.Model(&role).Updates(map[string]interface{}{"create_at": time.Now().Add(1 * time.Minute), "name": "Jerry"})
	/*指定字段不会走钩子函数,使用table也不会触发钩子函数*/
	//saveDb.Model(&role).UpdateColumn("name", "oliva")
	//saveDb.Model(&role).UpdateColumns(Role{Name: "", Description: ""})
	saveDb.Model(Role{}).Where(" id = ?", 1).UpdateColumns(map[string]interface{}{"name": "Jerry", "description": "Jerry is exploitation"})
}

func GormBatchUpdate() {
	/*批量更新*/
	statement := saveDb.Session(&gorm.Session{DryRun: true}).Model(&Role{}).Select("name", "create_at").Updates(Role{Name: "Jerry", CreateAt: time.Now().Add(1 * time.Hour)}).Statement
	fmt.Println(statement.SQL.String())
	fmt.Println(statement.Vars)
	//saveDb.Model(&Role{}).Where("id in @Ids", sql.Named("Ids", []int{1, 2, 3, 4})).Updates(map[string]interface{}{"name": "", "create_at": time.Now().Add(1 * time.Hour)})
}

func GormHandleGlobUpdate() {
	/*阻止全局更新,gorm没有条件不会更新但是Exec方法除外*/
	//saveDb.Model(Role{}).Exec("UPDATE ROLE SET NAME = ?", "Tom")
	/*这里报错Model方法需要使用指针类型的对象,但是换成其他表又可以使用值类型*/
	saveDb.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&Role{}).Update("name", "Jerry")
}

func GormUpdateByExpression() {
	/*通过表达式更新*/
	saveDb.Model(&Role{}).Where("id = ?", 1).Update("name", gorm.Expr("CONCAT('my name is ', ?)", "jerry"))
}

func GormUpdateBySubQuery() {
	/*通过子查询更新*/
	saveDb.Model(&Role{}).Update("name", saveDb.Model(Role{}).Select("name").Where("id = ?", 1))
}

func GormUpdateRawReturnData() {
	var roles []Role
	/*返回所有的字段Mysql类似不支持*/
	saveDb.Model(&roles).Clauses(clause.Returning{}).Where("id = ?", 2).Update("name", "Tom")
	fmt.Printf("GormUpdateRawReturnData is %+v\n", roles)
}
