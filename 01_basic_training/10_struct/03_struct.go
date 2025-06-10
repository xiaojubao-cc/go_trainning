package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

/*标签中属性高定义：
标签                         说明                                         示例
required              字段存在并且非零值                               validate:"required"
omitempty             字段为空时跳过后续验证                            validate:"omitempty",min=3
len=N                 长度必须等于N包括(字符串/数组/切片/map)            	validate:"len=3"
min=N                 最小长度N包括(字符串/数组/切片/map)                validate:"min=3"
max=N                 最大长度N包括(字符串/数组/切片/map)                validate:"max=3"
eq=N                  判断值是否等于                                   validate:"eq=3"
ne=N                  判断值是否不等于                                 validate:"ne=2"
gt=N                  大于N                                          validate:"gt=3"
lt=N                  小于N                                         validate:"lt=3"
gte=N                 大于等于N                                      validate:"gte=3"
lte=N                 小于等于N                                      validate:"lte=3"
contains=text         字符串必须包含text                              validate:"contains=@"
excludes=text         字符串必须不包含text                             validate:"excludes=@"
alpha                 字符串必须全部是字母(a-Z)                         validate:"alpha"
alphanum              字符串必须全部是字母和数字(a-Z0-9)                 validate:"alphanum"
numeric               字符串必须全部是数字(0-9)                         validate:"numeric"
email                 字符串必须为邮箱格式                              validate:"email"
url                   字符串必须为URL格式                              validate:"url"
uuid                  字符串必须为UUID格式                             validate:"uuid"
lowercase             字符串必须全部是小写字母(a-z)                      validate:"lowercase"
uppercase             字符串必须全部是大写字母(a-z)                      validate:"uppercase"
dive                  递归验证(slice,map,数组)                         validate:"dive,required"
keys/endkeys          递归验证map结合dive使用endkeys后的规则作用于value   validate:"dive keys required endkeys"
eqfiled=field         必须等于指定字段的值                          validate:"eqfield=Name"
neqfield=field        必须不等于指定字段的值                         validate:"neqfield=Name"
gtfield=field         必须大于指定字段的值                          validate:"gtfield=Age"
oneof=val1  val2      必须是枚举val1和val2                         validate:"oneof=val1 val2"
datetime              必须符合指定时间格式(默认"2006-01-02")         validate:"datetime=2006-01-02 15:04"
timezone              必须合法时区                                 validate:"timezone"
*/

type DbBackedUser struct {
	Name sql.NullString `validate:"required"`
	Age  sql.NullInt64  `validate:"required"`
	ID   string         `validate:"required,length"`
}

var validate *validator.Validate

func main() {
	validate = validator.New()
	validate.RegisterCustomTypeFunc(ValidateValuer, sql.NullString{}, sql.NullInt64{}, sql.NullBool{}, sql.NullFloat64{})
	/*第三个参数是字段值为空是否跳过验证*/
	validate.RegisterValidation("length", customLengthValidator)
	user := DbBackedUser{Name: sql.NullString{String: "lisi", Valid: true}, Age: sql.NullInt64{Int64: 0, Valid: true}, ID: "12345678"}
	err := validate.Struct(user)
	if err != nil {
		fmt.Printf("Err(s):\n%+v\n", err)
	}
}

func customLengthValidator(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		if len(value) > 5 {
			return false
		}
	}
	return true
}

// ValidateValuer /*自定义函数预处理*/
func ValidateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		value, err := valuer.Value()
		if err == nil {
			return value
		}
	}
	return nil
}
