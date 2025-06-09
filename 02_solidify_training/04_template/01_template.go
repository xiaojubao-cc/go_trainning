package main

import (
	"log"
	"os"
	"text/template"
)

/*模板传入数据*/
func main() {
	files, err := template.ParseFiles("D:\\golang projects\\go_training\\02_solidify_training\\04_template\\01_tpl.gohtml")
	if err != nil {
		log.Fatalf("解析模板异常：%s\n", err)
	}
	/*构造结构体进行多参数的传递*/
	param := struct {
		/*这里需要大写，否则模板无法获取*/
		Num int
		Arr []int
	}{
		Num: 10,
		Arr: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}
	err = files.Execute(os.Stdout, param)
	if err != nil {
		log.Fatalf("执行模板异常：%s\n", err)
	}
}
