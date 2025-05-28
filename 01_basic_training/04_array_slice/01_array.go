package main

import (
	"fmt"
	"slices"
)

/*数组*/
/*引用类型赋值：slice、map、channel、指针共享底层数据，如果触发扩容或重新分配内存则不再共享*/
func main() {
	var arr [5]int = [...]int{1, 2, 3, 4, 5}
	for i := 0; i < len(arr); i++ {
		fmt.Printf("for循环value=%d\n", i)
	}
	for index, value := range arr {
		fmt.Printf("range循环：index=%d,value=%d\n", index, value)
	}
	fmt.Printf("查询第0个元素：%d\n", arr[0])
	arr[0] = 12
	fmt.Printf("修改第0个元素：%d\n", arr[0])
	/*使用切片进行修改前闭后开*/
	var arr1 [5]int
	copy(arr1[:], arr[:])
	fmt.Printf("arr地址：%p,arr1地址：%p\n", &arr, &arr1)
	slice := slices.Clone(arr[:])
	fmt.Printf("arr地址：%p,slice地址：%p", &arr, &slice)
	slice[0] = 100
	fmt.Printf("修改前：%d,修改后：%d\n", arr[0], slice[0])
	/*slice是引用数据类型虽然s1和s2地址不同但是由于数据的未超过底层数组的容量仍然使用相同的底层数组*/
	s1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9} // cap = 9
	s2 := s1[3:4]
	s2 = append(s2, 1)
	fmt.Println(s2)
	fmt.Println(s1)
	fmt.Printf("s1地址：%p,s2地址：%p\n", &s1, &s2)
}
