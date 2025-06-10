package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"reflect"
	"strings"
)

// 定义封装的错误struct
type validationError struct {
	Namespace       string `json:"namespace"`
	Field           string `json:"field"`
	StructNamespace string `json:"structNamespace"`
	StructField     string `json:"structField"`
	Tag             string `json:"tag"`
	ActualTag       string `json:"actualTag"`
	//错误类型
	Kind string `json:"kind"`
	//错误类型
	Type string `json:"type"`
	//错误字段的值
	Value string `json:"value"`
	//错误字段tag设置的默认值
	Param   string `json:"param"`
	Message string `json:"message"`
}
type Gender uint

const (
	Male Gender = iota
	Female
	Intersex
)

func (gender Gender) String() string {
	terms := []string{"Male", "Female", "Intersex"}
	if gender < Male || gender > Intersex {
		return "unknown"
	}
	return terms[gender]
}

type Users struct {
	FirstName      string       `json:"fname"`
	LastName       string       `json:"lname"`
	Age            uint8        `validate:"gte=0,lte=130"`
	Email          string       `json:"e-mail" validate:"required,email"`
	FavouriteColor string       `validate:"hexcolor|rgb|rgba"`
	Addresses      []*Addresses `validate:"required,dive,required"` // a person can have a home and cottage...
	Gender         Gender       `json:"gender" validate:"required,gender_custom_validation"`
}

type Addresses struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required"`
}

// use a single instance of Validate, it caches struct info
var vali *validator.Validate

func main() {
	vali = validator.New()
	//缓存tag到validate
	vali.RegisterTagNameFunc(func(fl reflect.StructField) string {
		//对json标签进行两次分割
		name := strings.SplitN(fl.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	vali.RegisterStructValidation(UserStructLevelValidation, Users{})

	err := vali.RegisterValidation("gender_custom_validation", genderCustomValidation)
	if err != nil {
		log.Fatalf("RegisterValidation exception:%s", err)
		return
	}

	address := &Addresses{
		Street: "Eavesdown Docks",
		Planet: "Persphone",
		Phone:  "none",
		City:   "Unknown",
	}

	user := &Users{
		FirstName:      "",
		LastName:       "",
		Age:            45,
		Email:          "Badger.Smith@gmail",
		FavouriteColor: "#000",
		Addresses:      []*Addresses{address},
	}

	errs := vali.Struct(user)
	if errs != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(errs, &invalidValidationError) {
			log.Printf("invalidValidationError: %s", invalidValidationError)
			return
		}

		var validationErrors validator.ValidationErrors
		//匹配错误并赋值
		if errors.As(errs, &validationErrors) {
			for _, e := range validationErrors {
				/*频繁修改建议指针大结构体也需要使用指针,值类型属于复制副本*/
				validationErrors := &validationError{
					Namespace:       e.Namespace(),
					StructNamespace: e.StructNamespace(),
					Field:           e.Field(),
					StructField:     e.StructField(),
					Tag:             e.Tag(),
					ActualTag:       e.ActualTag(),
					Kind:            fmt.Sprintf("%v", e.Kind()),
					Type:            fmt.Sprintf("%v", e.Type()),
					Value:           fmt.Sprintf("%v", e.Value()),
					Param:           e.Param(),
					Message:         e.Error(),
				}
				indent, _ := json.MarshalIndent(validationErrors, "", " ")
				fmt.Println(string(indent))
			}
		}
		return
	}
}

func genderCustomValidation(fl validator.FieldLevel) bool {
	gender := fl.Field().Interface().(Gender)
	return gender.String() != "unknown"
}

// UserStructLevelValidation 自定义结构体验证
func UserStructLevelValidation(sl validator.StructLevel) {
	user := sl.Current().Interface().(Users)
	if len(user.FirstName) == 0 || len(user.LastName) == 0 {
		//自定义错误输出格式"Key: 'Users.fname' Error:Field validation for 'fname' failed on the 'fnameoflname' tag"
		sl.ReportError(user.FirstName, "fname", "FirstName", "fnameoflname", "")
		sl.ReportError(user.LastName, "lname", "LastName", "fnameoflname", "")
	}
}
