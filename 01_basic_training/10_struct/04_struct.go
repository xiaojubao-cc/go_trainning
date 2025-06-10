package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

// dive递归验证集合内元素

type diveValidate struct {
	Array []string          `validate:"required,gt=0,dive,required"`
	Map   map[string]string `validate:"required,gt=0,dive,keys,keyMax,endkeys,required,max=1000"`
}

var validated *validator.Validate

func main() {
	validated = validator.New()
	validated.RegisterAlias("keyMax", "max=10")
	var diveValidate diveValidate
	validating(diveValidate)
	diveValidate.Array = []string{""}
	diveValidate.Map = map[string]string{"key": "value"}
	validating(diveValidate)
}

func validating(v diveValidate) {
	fmt.Println("starting.....")
	err := validated.Struct(v)
	if err != nil {
		fmt.Printf("Err(s):\n%+v\n", err)
	}
	fmt.Println("ending.....")
}
