package main

import (
	"fmt"
	"reflect"
)

/*值类型*/
func main() {
	str := "hello world"
	fmt.Printf("%T\n", reflect.ValueOf(str))
	fmt.Println(reflect.ValueOf(str).Type())
	fmt.Println(reflect.ValueOf(str).Kind())

	var prt = new(string)
	*prt = "hello world"
	fmt.Printf("%v\n", reflect.ValueOf(prt).Elem().Kind())
	fmt.Printf("获取原本值：%s\n", reflect.ValueOf(prt).Elem().Interface())

}
