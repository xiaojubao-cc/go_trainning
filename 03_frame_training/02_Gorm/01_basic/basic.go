package main

/*
1.字段级标签
标签                          作用                          示例
column:name              指定数据库列名             gorm:"column:name"
type:type                指定数据库类型             gorm:"type:varchar(100)"
size:size                指定字段大小               gorm:"size:100"
primaryKey               指定字段为主键             gorm:"primaryKey"
autoIncrement            指定字段为自增字段          gorm:"primaryKey;autoIncrement"
unique                   指定唯一索引               gorm:"unique"
index:name               指定索引名称               gorm:"index:index_name"
index:name,option        指定索引名称和选项          gorm:"index:index_name,class:(FULLTEXT, expression:(lower(name)))
uniqueIndex:name         创建唯一索引名称            gorm:"uniqueIndex:idx_name_unique"
not null                 非空约束                  gorm:"not null"
nullable                 允许空约束                 gorm:"nullable"
default:value            指定默认值                 gorm:"default:value"
comment:value            指定字段注释               gorm:"comment:value"
embeddable               嵌入结构体                 gorm:"embedded"
embeddedPrefix:prefix    嵌入结构体前缀              gorm:"embeddedPrefix:prefix"
->                       只读                      gorm:"->"
<-                       只写                      gorm:"<-"
<->                      读写                      gorm:"<->"
-                        忽略                      gorm:"-"
saveAssociations         保存关联记录               gorm:"saveAssociations"
serializer:json          JSON序列化                gorm:"serializer:json"
encrypt                  字段加密                  gorm:"encrypt:gcm"
2.关联标签
foreignKey:column        指定外键列名                gorm:"foreignKey:column"
references:column        指定关联列名                gorm:"references:column"
polymorphic:column       多态关联列名                gorm:"polymorphic:Commentable"
polymorphicValue:value   多态关联值                 gorm:"polymorphicValue:value"
many2many:table          多对多关联表名              gorm:"many2many:user_roles;"
joinForeignKey:column    多对多关联外键列名           gorm:"many2many:user_roles;joinForeignKey:UserID"
constraint:OnUpdate      指定外键更新策略            gorm:"constraint:OnUpdate:(CASCADE,SET NULL)"
constraint:OnDelete      指定外键删除策略            gorm:"constraint:OnDelete:(CASCADE,SET NULL) CASCADE:级联删除"
3.模型标签
table:table              指定表名                    gorm:"table:table_name"
singularTable            指定单数表名                 gorm:"singularTable:true"
disableDefaultConstraint 禁用默认外键约束              gorm:"disableDefaultConstraints"

4.struct中使用值类型还是指针类型
默认使用值类型：适用于大多数核心实体;简化代码逻辑;提高局部性，减少内存碎片
在以下情况使用指针类型：关联对象可能不存在;需要区分加载状态;大型结构体（减少复制开销）;需要修改原始对象;多态关联关系

5.joins和preload区别:
preload:将结果组装到嵌套结构体;关联层级较深(>2);数据量小于1000;不需要基于关联表字段过滤;不能映射到map;
joins:基于关联表字段过滤(字段或者条件或者关联表字段聚合);数据量大于1000;需要单次查询;
preload和joins结合使用:joins只能填充OrderItem的基础字段，preload填充OrderItem的关联字段Product

6.Model和Table区别
	方法                 作用             是否需要结构体           是否支持关联操作           是否自动映射字段       是否触发钩子函数
Model(&Struct{})  绑定结构体与数据库表关系  ✅ 必须传入结构体指针      ✅ 支持关联操作            ✅ 自动映射字段           会
Table("table_name")   直接指定表名          ❌ 无需结构体          ❌ 不支持关联操作          ❌ 不自动映射字段         不会

7.Find和Scan区别
	使用场景                 推荐方法               说明
	查询完整模型数据            Find           支持关联、分页、零值过滤
	聚合/分组查询              Scan              需配合 AS 别名
	动态字段处理               Scan        使用 map[string]interface{}
	关联查询                  Find             支持 Preload、Joins
	Raw SQL 查询          Raw().Scan()       最灵活，但需手动维护 SQL

8.原生sql操作
	1.原生查询:Raw()...Scan()
	2.原生执行:db.Exec()
	3.sql测试不会执行:db.Session(&gorm.Session{DryRun: true})
	4.用于生成sql:db.ToSQL(func)
	5.row和rows结合scan

9.事务
	事务只会捕获error，不会捕获panic,嵌套事务和回滚点

10.session
	/*DryRun: true,获取sql打印*/
/*PrepareStmt: true,缓存sql使用于高频相同查询*/
/*NewDB: true,创建新的连接*/
/*Initialized: true,*/
/*SkipHooks: true,跳过钩子函数*/
/*DisableNestedTransaction: true,禁用嵌套事务*/
/*AllowGlobalUpdate: true,允许全表更新*/
/*FullSaveAssociations: true,自动保存关联及引用*/
/*上下文*/
/*QueryFields: true,按字段查询*/
/*CreateBatchSize: 100,创建数据批次*/
