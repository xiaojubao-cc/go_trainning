package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func validateMap(data, rule map[string]interface{}) {
	errs := validate.ValidateMap(data, rule)
	if errs != nil {
		for key, err := range errs {
			fmt.Printf("map key:%s,err:%s\n", key, err)
		}
	}
}

type validateNestMap struct {
	Map map[string]interface{} `validate:"nestMap"`
}

func main() {
	validate = validator.New()
	data := map[string]interface{}{"name": "zytel", "email": "zytel3301@gmail.com"}
	rule := map[string]interface{}{"name": "required,min=3,max=10,alpha", "email": "omitempty,required,email"}
	validateMap(data, rule)
	/*嵌套map验证*/
	data = map[string]interface{}{
		"name":  "Arshiya Kiani",
		"email": "zytel3301@gmail.com",
		"details": map[string]interface{}{
			"family_members": map[string]interface{}{
				"father_name": "Micheal",
				"mother_name": "Hannah",
			},
			"salary": "1000",
			"phones": []map[string]interface{}{
				{
					"number": "11-111-1111",
					"remark": "home",
				},
				{
					"number": "22-222-2222",
					"remark": "work",
				},
			},
		},
	}
	rule = map[string]interface{}{
		"name":  "min=4,max=32",
		"email": "required,email",
		"details": map[string]interface{}{
			"family_members": map[string]interface{}{
				"father_name": "required,min=4,max=32",
				"mother_name": "required,min=4,max=32",
			},
			"salary": "number",
			"phones": map[string]interface{}{
				"number": "required,min=4,max=32",
				"remark": "required,min=1,max=32",
			},
		},
	}
	validateMap(data, rule)
}
