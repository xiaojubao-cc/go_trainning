package main

import "fmt"

/*泛型方法*/
func sum[T int | float64](a, b T) T {
	return a + b
}

// GenericsSlice /*泛型切片*/
type GenericsSlice[T int | string] []T

// GenericsMap /*泛型字典*/
type GenericsMap[K string, V interface{}] map[K]V

// GenericsPipe 定义管道类型
type GenericsPipe[T int | string] chan T

// GenericsStruct /*泛型结构体*/
type GenericsStruct[T int | string, S int | string] struct {
	name string
	id   T
	list []S
}

// GenericsInterface /*泛型接口*/
type GenericsInterface[T int | string] interface {
	genericsWay() T
}

func (g *GenericsStruct[T, S]) genericsWay() T {
	return g.id
}

func main() {
	sumInt := sum[int](1, 2)
	fmt.Println(sumInt)
	sumFloat := sum[float64](1.1, 2.2)
	fmt.Printf("%.2f", sumFloat)

	slice := make(GenericsSlice[int], 5)
	for _, val := range append(slice, 1, 3, 5, 7, 9) {
		fmt.Println(val)
	}

	slices := make(GenericsSlice[string], 5)
	for _, val := range append(slices, "1", "3", "5", "7", "9") {
		fmt.Println(val)
	}

	genericsMap := GenericsMap[string, interface{}]{
		"name": "golang",
		"age":  18,
	}
	for key, val := range genericsMap {
		fmt.Println(key, val)
	}

	gs := &GenericsStruct[int, string]{
		name: "golang",
		id:   18,
		list: []string{"1", "2", "3", "4", "5"},
	}
	fmt.Printf("%+v\n", gs)
	fmt.Printf("实现泛型接口：%d", gs.genericsWay())

}
