package main

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type User struct {
	FirstName      string     `validate:"required"`
	LastName       string     `validate:"required"`
	Age            uint8      `validate:"gte=0,lte=130"`
	Email          string     `validate:"required,email"`
	Gender         string     `validate:"oneof=male female prefer_not_to"`
	FavouriteColor string     `validate:"iscolor"`                // alias for 'hexcolor|rgb|rgba|hsl|hsla'
	Addresses      []*Address `validate:"required,dive,required"` // a person can have a home and cottage...
}

// Address houses a users address information
type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required,max=11"`
}

var v *validator.Validate

func validateStruct() {
	v = validator.New(validator.WithRequiredStructEnabled())
	address := Address{
		Street: "Eavesdown Docks",
		City:   "Halifax",
		Planet: "Tatooine",
		Phone:  "131571456015",
	}
	user := User{
		FirstName:      "Badger",
		LastName:       "Smith",
		Age:            135,
		Email:          "<EMAIL>",
		Gender:         "male",
		FavouriteColor: "#000-",
		Addresses:      []*Address{&address},
	}
	err := v.Struct(user)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			fmt.Println(err)
			return
		}

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, e := range validationErrors {
				//获取错误字段的命名空间User.Addresses[0].Street
				fmt.Println(e.Namespace())
				//获取错误字段的名称Street
				fmt.Println(e.Field())
				//获取错误字段完整路径User.Addresses[0].Street
				fmt.Println(e.StructNamespace())
				//获取错误字段的属性名称Street
				fmt.Println(e.StructField())
				//获取错误字段的标签
				fmt.Println(e.Tag())
				//获取错误字段的标签
				fmt.Println(e.ActualTag())
				//获取错误字段的类型
				fmt.Println(e.Kind())
				//获取错误字段的reflect.Type
				fmt.Println(e.Type())
				//获取错误字段的值
				fmt.Println(e.Value())
				//获取错误字段的参数值
				fmt.Println(e.Param())
				fmt.Println()
			}
		}
	}
	return
}

// 验证变量是否符合指定规则
func validateVar() {
	email := "1171264943@qq.com"
	errs := v.Var(email, "required,email")
	if errs != nil {
		fmt.Println(errs)
		return
	}
}
func main() {
	validateStruct()
	validateVar()

}
