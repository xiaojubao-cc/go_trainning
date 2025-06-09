package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)

type Person struct {
	FirstName string
	LastName  string
}

func main() {
	tel, err := template.ParseFiles("D:\\golang projects\\go_training\\02_solidify_training\\04_template\\03_tpl.gohtml")
	if err != nil {
		log.Fatalf("parse tempalte err：%s\n", err)
	}
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		firstName := req.FormValue("first")
		lastName := req.FormValue("last")
		fmt.Printf("firstName is %s, lastName is %s\n", firstName, lastName)
		/*这里必须使用请求过滤*/
		if req.Method == "POST" {
			file, header, err := req.FormFile("file")
			if err != nil {
				log.Fatalf("get form file err：%s\n", err)
			}
			defer file.Close()
			fmt.Printf("file header is %+v", header)
			io.Copy(os.Stdout, file)
		}
		tel.Execute(resp, Person{
			FirstName: firstName,
			LastName:  lastName,
		})
	})

	http.ListenAndServe(":8080", nil)
}
