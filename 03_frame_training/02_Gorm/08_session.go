package main

import (
	"context"
	"gorm.io/gorm"
	"time"
)

func main() {
	nativeDb.Session(&gorm.Session{
		DryRun: true,
		/*缓存sql使用于高频相同查询*/
		PrepareStmt: true,
		/*创建新的连接*/
		NewDB: true,
		/**/
		Initialized: true,
		/*跳过钩子函数*/
		SkipHooks: true,
		/*禁用嵌套事务*/
		DisableNestedTransaction: true,
		/*允许全表更新*/
		AllowGlobalUpdate: true,
		/*自动保存关联及引用*/
		FullSaveAssociations: true,
		/*上下文*/

		/*按字段查询*/
		QueryFields: true,
		/*创建数据批次*/
		CreateBatchSize: 100,
	})
}

func GormSessionDryRun() {
	/*用于生成sql*/
	statement := nativeDb.Session(&gorm.Session{DryRun: true}).Table("employee").Statement
	/*获得完整的执行sql*/
	nativeDb.Dialector.Explain(statement.SQL.String(), statement.Vars...)
}

func GormSessionNewDB() {
	/*
			场景              推荐配置                说明
		生成 SQL 预览        NewDB: true     避免 DryRun 影响后续真实查询。
		单元测试             NewDB: true    为每个测试用例创建独立实例，防止数据污染。
		事务隔离             NewDB: true     在嵌套事务中管理独立的回滚/提交逻辑
		临时修改日志级别      NewDB: false    短时调整日志输出，无需创建新连接池。
		高并发请求           NewDB: false    复用连接池资源，避免频繁创建/销毁连接。
	*/
	nativeDb.Session(&gorm.Session{NewDB: true})
}

func GormSessionWithContext() {
	/*超时上下文取消*/
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*100)
	tx := nativeDb.Session(&gorm.Session{Context: ctx})
	tx.Table("employee").Find(&Employee{})
	/*清理session配置*/
	tx.Session(&gorm.Session{})
	tx.Table("employee").Where("id = ?", 1).Find(&Employee{})
}

func GormSessionQueryFields() {
	/*
	   未开启查询的sql:select * from table
	   配合struct的column标签进行字段映射
	   type User struct {
	         ID   uint   `gorm:"column:id"`
	         Name string `gorm:"column:name"`
	         Role string `gorm:"-"` // 不会被查询
	     }
	*/
	nativeDb.Session(&gorm.Session{QueryFields: true}).Table("employee").Find(&Employee{})
}
