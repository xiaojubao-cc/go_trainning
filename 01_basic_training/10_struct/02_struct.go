/*组合*/
package main

import "fmt"

// Student
/*这里的组合类型Class是使用指针类型还是在值类型取决于Class数据是否需要共享，*Class类型的数据存储于堆中，Student只是存储指针*/
/*指针嵌入适合大型数据结构，值类型适合小型数据结构*/
type Student struct {
	name string
	age  int
	*Class
	Teacher
}

type Class struct {
	cname string
}

type Teacher struct {
	name string
	age  int
}

func main() {
	student1 := &Student{
		name:  "lihua",
		age:   18,
		Class: &Class{cname: "class one"},
		Teacher: Teacher{
			name: "zhangsan",
			age:  18,
		},
	}
	student2 := Student{
		name: "lihua",
		age:  18,
		Teacher: Teacher{
			name: "zhangsan",
			age:  18,
		},
	}
	fmt.Printf("Student:%+v\n", student1)
	fmt.Printf("Student:%+v\n", student2)
	/*使用指针作为属性可以通过语法糖解引用获取属性值(*student).name = student.name*/
	fmt.Printf("Student->name:%s,Student->age:%d,Student->Class->cname:%s,Student->Teacher->name:%s",
		student1.name, student1.age, student1.Class.cname, student1.Teacher.name)
}
