package main

import (
	"fmt"
	"reflect"
)

type Interface interface {
	Say()
	Print()
}

type InterfaceImpl struct {
}

func (i *InterfaceImpl) Say() {
	fmt.Println("hello world")
}
func (i *InterfaceImpl) Print() {
	fmt.Println("hello world")
}

/*反射*/
/*elem()会剥去指针那一层，只会暴露原本的类型*/
func main() {
	str := "hello world"
	typeOf := reflect.TypeOf(str)
	fmt.Printf("str的类型是：%s\n", typeOf)
	fmt.Printf("获取变量的基本类型：%s\n", typeOf.Kind())
	/*elem指针，切片，数组，通道，映射表*/
	slice := make([]int, 0)
	fmt.Printf("获取引用类型的基本数据类型：%v\n", reflect.TypeOf(slice))
	fmt.Printf("获取引用类型的基本数据类型：%s\n", reflect.TypeOf(slice).Elem())

	var ptr = new([]int)
	fmt.Printf("获取引用类型的基本数据类型：%v\n", reflect.TypeOf(ptr))
	fmt.Printf("获取引用类型的基本数据类型：%s\n", reflect.TypeOf(ptr).Elem())
	fmt.Printf("获取引用类型的基本数据类型名称：%s\n", reflect.TypeOf(ptr).Elem().Kind())
	fmt.Printf("获取引用类型的字节长度：%d\n", reflect.TypeOf(ptr).Size())
	/*判断方法是否实现了某一个接口*/
	inter := new(Interface)
	interfaceImpl := new(InterfaceImpl)
	implements := reflect.TypeOf(interfaceImpl).Implements(reflect.TypeOf(inter).Elem())
	fmt.Printf("判断方法是否实现了某一个接口：%t\n", implements)
	/*判断一个类型是否可以被转换为另外一个类型*/
	convert := reflect.TypeOf(interfaceImpl).ConvertibleTo(reflect.TypeOf(inter).Elem())
	fmt.Printf("判断方法是否可以强转：%t\n", convert)
	fmt.Printf("%v\n", reflect.TypeOf(interfaceImpl).Elem())
	/*kind获取基本数据类型*/
	fmt.Printf("%v\n", reflect.TypeOf(interfaceImpl).Elem().Kind())
	fmt.Printf("%v\n", reflect.TypeOf(inter).Elem())
}
