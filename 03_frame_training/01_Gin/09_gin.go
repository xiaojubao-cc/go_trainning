package main

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"reflect"
	"sync"
)

// 自定义验证器
/*这里是惰性加载*/
type defaultValidator struct {
	validate *validator.Validate
	once     sync.Once
}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyInit()
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyInit()
	return v.validate
}

/*Go 语言中用于静态检查接口实现的惯用写法，其核心目的是确保 defaultValidator 类型完全实现了 binding.StructValidator 接口的所有方法*/
var _ binding.StructValidator = &defaultValidator{}

func (v *defaultValidator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
func main() {
	binding.Validator = new(defaultValidator)
}
