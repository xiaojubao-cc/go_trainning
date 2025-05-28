package main

import (
	"fmt"
)

//Log 结构体嵌入在 Customer 结构体中。
/**
    指针嵌入：
	内存结构：Customer 结构体包含一个指向 Log 结构体的指针。
	存储方式：Log 的实例独立分配在内存中（通常在堆上），Customer 仅存储其地址。
    指针类型嵌入 (*Log)： 适合大型结构体或需要共享数据，但需处理指针初始化和 nil 安全问题。
*/
type Log struct {
	msg string
}

type Customer struct {
	Name string
	log  *Log
}

func main() {
	// c := new(Customer)
	// c.Name = "Barak Obama"
	// c.log = new(Log)
	// c.log.msg = "1 - Yes we can!"
	// shorter:
	//必须显式分配内存&Log{"1 - Yes we can!"}
	c := &Customer{"Barak Obama", &Log{"1 - Yes we can!"}}
	// fmt.Println(c)   // &{Barak Obama 1 - Yes we can!}
	//直接调用指针c.Log.Add()
	//c.Log().Add("2 - After me the world will be a better place!")
	c.log.Add("3 - After me the world will be a better place!")
	//fmt.Println(c.log)
	fmt.Println(c.Log())
}

func (l *Log) Add(s string) {
	l.msg += "\n" + s
}

func (l *Log) String() string {
	return l.msg
}

func (c *Customer) Log() *Log {
	return c.log
}
