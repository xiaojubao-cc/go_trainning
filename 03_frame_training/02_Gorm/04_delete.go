package main

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Roles struct {
	ID          uint           `gorm:"primaryKey;autoIncrement;comment:ID" json:"id"`
	Name        string         `gorm:"column:name;not null;comment:名称" json:"name"`
	Description string         `gorm:"not null;comment:描述" json:"description"`
	CreateAt    time.Time      `gorm:"column:create_at;comment:创建时间;default:CURRENT_TIMESTAMP" json:"create_at"`
	UpdateAt    time.Time      `gorm:"column:update_at;comment:更新时间;default:CURRENT_TIMESTAMP" json:"update_at"`
	Deleted     gorm.DeletedAt //这里可以替换为IsDel soft_delete.DeletedAt `gorm:"softDelete:flag"`
}

func (*Roles) TableName() string {
	return "role"
}
func (role *Roles) BeforeDelete(tx *gorm.DB) (err error) {
	if role.ID == 1 {
		log.Println("Illegal Deletion Admin ")
		return errors.New("Illegal Deletion Admin ")
	}
	return nil
}

var deleteDb *gorm.DB

func init() {
	deleteDb = deleteBasicConnection()
}
func deleteBasicConnection() *gorm.DB {
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
	db.AutoMigrate(&Roles{})
	return db
}

func main() {
	/*指定主键继续删除*/
	//GormDeleteByID()
	//SoftDelete()
	//var roles []Roles = []Roles{
	//	{ID: 1, Name: "Tom", Description: "boss"},
	//	{ID: 2, Name: "Jerry", Description: "dev"},
	//	{ID: 3, Name: "Jordan", Description: "test"},
	//}
	//deleteDb.Create(&roles)
}

func GormDeleteByID() {
	/*使用主键删除*/
	//var role Roles = Roles{ID: 1}
	//deleteDb.Delete(&role)
	//deleteDb.Delete(&Roles{}, 1)
	deleteDb.Model(Roles{}).Where("id = ?", 3).Delete(&Roles{})
}

func GormHandleGlobDelete() {
	/*使用原生sql执行全局删除*/
	//deleteDb.Exec("DELETE FROM role")
	/*修改session属性*/
	deleteDb.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Roles{})
}

func SoftDelete() {
	/*软删除:struct中新增字段Deleted,并且类型为gorm.DeleteAt;删除之前该字段是null,删除之后该字段是删除时间*/
	//deleteDb.Delete(&Roles{}, 1)
	/*查找软删除记录Unscoped()*/
	/*永久删除记录deleteDb.Unscoped().Delete(&Roles{ID:1})*/
	var roles []Roles
	deleteDb.Unscoped().Find(&roles)
	fmt.Println(roles)
}
