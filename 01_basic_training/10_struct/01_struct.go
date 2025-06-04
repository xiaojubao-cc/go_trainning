package main

import "fmt"

// Person /*结构体*/
/*结构嵌套优先使用匿名字段类似于继承,嵌套实体类的字段会在json序列化时自动展开*/
type Person struct {
	Name     string
	Age      int
	Address  string
	Salary   float64
	Birthday string
}

func (p *Person) getName() string {
	return p.Name
}

// PersonOptions /*声明一个类型*/
type PersonOptions func(p *Person)

// WithName /*赋值名字*/
func WithName(name string) PersonOptions {
	return func(p *Person) {
		p.Name = name
	}
}
func WithAge(age int) PersonOptions {
	return func(p *Person) {
		p.Age = age
	}
}

func WithAddress(address string) PersonOptions {
	return func(p *Person) {
		p.Address = address
	}
}

func WithSalary(salary float64) PersonOptions {
	return func(p *Person) {
		p.Salary = salary
	}
}

func WithBirthday(birthday string) PersonOptions {
	return func(p *Person) {
		p.Birthday = birthday
	}
}

// NewPerson /*构造函数*/
func NewPerson(opt ...PersonOptions) *Person {
	person := &Person{}
	for _, option := range opt {
		option(person)
	}
	/*这里还可以针对传入的值进行校验*/
	if person.Age < 0 {
		person.Age = 0
	}
	return person
}
func main() {
	var person = NewPerson(
		/*这里相当于是已经调用了函数返回了PersonOptions类型函数*/
		WithName("Jerry"),
		WithAge(18),
		WithAddress("China"),
		WithSalary(5000),
		WithBirthday("1990-01-01"),
	)
	fmt.Printf(person.getName())
}
