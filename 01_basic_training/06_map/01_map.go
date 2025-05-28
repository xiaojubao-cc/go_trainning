package main

import (
	"fmt"
	"sync"
)

/*字典*/
func main() {
	var map1 map[string]string = make(map[string]string, 16)
	map1["name"] = "golang"
	fmt.Printf("name value:%s", map1["name"])
	map1["age"] = "18"
	map1["sex"] = "male"
	/*判断某个key是否存在*/
	if value, ok := map1["name"]; ok {
		fmt.Printf("name value:%s", value)
	} else {
		fmt.Printf("name value not exist")
	}
	/*循环*/
	for key, value := range map1 {
		fmt.Printf("key:%s,value:%s\n", key, value)
	}
	/*线程安全的Map*/
	var syncMap sync.Map
	/*存储*/
	syncMap.Store("name", "golang")
	syncMap.Store("age", "18")
	/*取值*/ /*若不存在则存储*/
	if value, ok := syncMap.Load("name"); ok {
		fmt.Printf("name value:%s\n", value)
	} else {
		syncMap.LoadOrStore("name", "golang")
	}
	/*遍历*/
	syncMap.Range(func(key, value interface{}) bool {
		fmt.Printf("key:%s,value:%s\n", key, value)
		return true
	})

}
