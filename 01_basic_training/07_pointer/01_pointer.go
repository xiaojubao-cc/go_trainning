package main

import "fmt"

/*指针*/
func main() {
	var prt *int = new(int)
	*prt = 10
	fmt.Println(*prt)

	prt1 := &prt
	/*解引用(*)操作获取指针所指向的元素*/
	fmt.Println(**prt1)
}
