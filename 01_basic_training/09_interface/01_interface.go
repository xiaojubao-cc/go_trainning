package main

// Person /*接口*/
type Person interface {
	Say()
	Walk()
}

type PersonImpl struct {
	name string
	age  int
}

func (p *PersonImpl) Say() {
	println("PersonImpl Say")
}
func (p *PersonImpl) Walk() {
	println("PersonImpl Walk")
}
func main() {
	p := &PersonImpl{
		name: "张三",
		age:  18,
	}
	p.Say()
	p.Walk()
}
