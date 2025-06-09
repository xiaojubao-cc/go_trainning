package main

import (
	"fmt"
	"html"
	"io"
	"net/http"
)

/*获取表单的数据*/
func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		name := "q"
		formValue := req.FormValue(name)
		fmt.Printf("获取的参数为：%s\n", html.UnescapeString(formValue))
		res.Header().Set("Content-Type", "text/html")
		io.WriteString(res, `<form method="POST">
		 <input type="text" name="q">
		 <input type="submit">
		</form>`)
	})
	http.ListenAndServe(":8080", nil)
}
