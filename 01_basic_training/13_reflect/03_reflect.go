package main

import (
	"fmt"
	"reflect"
)

/*通过反射调用函数*/

func multiply(a, b int) (int, error) {
	return a * b, nil
}
func main() {
	/*通过类型反射获取函数的基本信息*/
	typeOf := reflect.TypeOf(multiply)
	fmt.Printf("获取函数反射的输入参数类型：%s\n", typeOf.Name())
	fmt.Printf("获取函数反射的输入参数个数：%d\n", typeOf.NumIn())
	fmt.Printf("获取函数反射的输出参数个数：%d\n", typeOf.NumOut())
	for i := 0; i < typeOf.NumIn(); i++ {
		fmt.Printf("获取函数反射的输入参数类型：%s\n", typeOf.In(i).Kind())
	}
	for i := 0; i < typeOf.NumOut(); i++ {
		fmt.Printf("获取函数反射的输出参数类型：%s\n", typeOf.Out(i).Kind())
	}

	/*通过值反射调用函数*/
	reflectValue := reflect.ValueOf(multiply)
	returnValues := reflectValue.Call([]reflect.Value{reflect.ValueOf(2), reflect.ValueOf(3)})
	for _, returnValue := range returnValues {
		fmt.Println(returnValue.Interface())
	}

}
