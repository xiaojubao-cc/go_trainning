package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	UserId   string `json:"id"`
	Username string `json:"username"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
}

func main() {
	person := &Person{
		UserId:   "1",
		Username: "tom",
		Age:      18,
		Address:  "beijing",
	}
	//序列化操作
	marshal, err := json.Marshal(person)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%++v\n", string(marshal))
	//json格式化输出
	marshalIndent, _ := json.MarshalIndent(person, "", "\t")
	fmt.Printf("%++v\n", string(marshalIndent))
	var person1 Person
	//反序列化操作
	jsonStr := `{"id":"2","username":"jerry","age":18,"address":"成都"}`
	err = json.Unmarshal([]byte(jsonStr), &person1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", person1)
}
