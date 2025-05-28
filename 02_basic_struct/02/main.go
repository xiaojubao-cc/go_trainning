package main

import "fmt"

// Log 结构体嵌入在 Customer 结构体中。
/**
指针嵌入方式：
	内存结构：Customer 结构体直接包含 Log 结构体的所有字段。
	存储方式：Log 的实例直接内联在 Customer 的内存空间中。
    值类型嵌入 (Log)： 适合小型结构体，无需担心 nil 问题，但可能涉及内存复制。
*/
type Log struct {
	msg string
}

type Customer struct {
	Name string
	Log
}

func main() {
	//需要显示初始化值Log{"1 - Yes we can!"}
	c := &Customer{"Barak Obama", Log{"1 - Yes we can!"}}
	//调用方法自动转换为(&c.Log).Add()
	//c.Add("2 - After me the world will be a better place!")
	c.Add("2 - After me the world will be a better place!")
	fmt.Println(c)
}

func (l *Log) Add(s string) {
	l.msg += "\n" + s
}

func (c *Customer) String() string {
	return c.Name + "\nLog:" + fmt.Sprintln(c.Log)
}

func (l *Log) String() string {
	return l.msg
}

/* Output:
Barak Obama
Log:{1 - Yes we can!
2 - After me the world will be a better place!}
*/
