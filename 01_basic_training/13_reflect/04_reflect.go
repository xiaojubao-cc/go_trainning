package main

import (
	"fmt"
	"reflect"
)

/*结构体反射*/

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	money   int
}

func (p Person) Talk(msg string) string {
	return msg
}
func main() {

	/*通过索引访问结构体*/
	elem := reflect.TypeOf(new(Person)).Elem()
	fmt.Println(elem.NumField())
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fmt.Println(field.Name, field.Tag)
	}

	/*通过名称访问字段*/
	if name, ok := elem.FieldByName("Name"); ok {
		fmt.Println(name.Name, name.Tag)
		fmt.Println(name.Tag.Get("json"))
		fmt.Println(name.Tag.Lookup("json"))
	}
	/*通过字段访问Tag*/

	/*通过反射修改字段的值*/
	var p = new(Person)
	name := reflect.ValueOf(p).Elem().FieldByName("Name")
	if name.CanSet() {
		name.SetString("Tom")
	}
	fmt.Printf("%+v\n", p)
	/*修改私有属性的值*/
	money := reflect.ValueOf(p).Elem().FieldByName("money")
	/*判断结构体中是否有money字段*/
	if (money != reflect.Value{}) {
		/*构造私有字段的指针映射*/
		newPointer := reflect.NewAt(money.Type(), money.Addr().UnsafePointer())
		newPointer.Elem().SetInt(1000)
	}
	fmt.Printf("%+v\n", p)

	/*获取结构体中的方法*/
	p1 := reflect.TypeOf(new(Person)).Elem()
	for i := 0; i < p1.NumMethod(); i++ {
		method := p1.Method(i)
		for j := 0; j < method.Func.Type().NumIn(); j++ {
			fmt.Println(method.Func.Type().In(j))
		}
		for j := 0; j < method.Func.Type().NumOut(); j++ {
			fmt.Println(method.Func.Type().Out(j))
		}
	}
	/*调用struct中的方法*/
	p2 := reflect.ValueOf(new(Person)).Elem()
	talk := p2.MethodByName("Talk")
	if (talk != reflect.Value{}) {
		call := talk.Call([]reflect.Value{reflect.ValueOf("hello golang reflect")})
		for _, v := range call {
			fmt.Println(v.Interface())
		}
	}

	/*判断两个对象是否相等*/
	a := make([]int, 100)
	b := make([]int, 100)
	fmt.Println(reflect.DeepEqual(a, b))
}
