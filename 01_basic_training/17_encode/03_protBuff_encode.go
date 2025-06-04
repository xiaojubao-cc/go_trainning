package main

import (
	"fmt"
)

type Student struct {
	Id   int
	Name string
	Age  int
	Class
}

type Class struct {
	Name string
	Teacher
}
type Teacher struct {
	Name string
}

func main() {
	student := &Student{
		Id:   1,
		Name: "张三",
		Age:  18,
		Class: Class{
			Name: "1班",
			Teacher: Teacher{
				Name: "李老师",
			},
		},
	}
	fmt.Println(student)
}
