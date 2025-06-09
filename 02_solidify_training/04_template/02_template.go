package main

import (
	"log"
	"os"
	"text/template"
)

type AddStruct struct {
	A int
	B int
}

/*模板传入函数*/
func main() {
	telFile := "D:\\golang projects\\go_training\\02_solidify_training\\04_template\\02_tpl.gohtml"
	tel := template.New(telFile)
	tel.Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})
	parseFiles, err := tel.ParseFiles(telFile)
	if err != nil {
		log.Fatalf("解析模板异常：%s\n", err)
	}
	err = parseFiles.Execute(os.Stdout, AddStruct{
		A: 10,
		B: 20,
	})
	if err != nil {
		log.Fatalf("执行模板异常：%s\n", err)
	}

}
