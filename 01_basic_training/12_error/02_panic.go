package main

import "fmt"

/**panic：很严重的问题，程序应该在处理完问题后立即退出*/
/*
	recover 必须在 defer 函数中调用：
	recover 只有在 defer 函数中才能捕获到 panic。
	defer 确保在函数返回前执行 recover。
	panic 会传播：
	如果没有 recover 捕获，panic 会沿着调用栈向上传播，最终导致程序崩溃。
	可以在多个层级设置 defer 来捕获 panic。
	recover 只能捕获当前协程的 panic：
	recover 无法捕获其他协程的 panic。
	每个协程需要独立处理自己的 panic。
	避免滥用 panic：
	panic 应该用于处理程序无法恢复的严重错误。
	对于可恢复的错误，应使用 error 返回值。
*/
func division(a, b int) float32 {
	if b == 0 {
		panic("除数不能为0")
	}
	return float32(a / b)
}
func customizeRecover() {
	if r := recover(); r != nil {
		fmt.Printf("recovered from %s", r)
	}
}
func main() {
	defer customizeRecover()
	division(10, 0)
}
