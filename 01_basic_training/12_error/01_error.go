package main

import (
	"errors"
	"fmt"
)

// ERROR /*error:正常的流程出错，需要处理，直接忽略掉不处理程序也不会崩溃*/
var (
	ERROR = errors.New("square root of negative number")
)

// NorMathError /*实现ERROR接口*/
type NorMathError struct {
	code int
	msg  string
}

func (n *NorMathError) Error() string {
	return n.msg
}

func main() {
	norMathError := &NorMathError{code: 1001, msg: "square root of negative number"}
	fmt.Printf("自定义异常：%++v\n", norMathError)
	printErr := fmt.Errorf("打印异常信息：square root of negative number")
	fmt.Printf("%++v", printErr)
}
