package main

import (
	"fmt"
	"net/http"
)

/*获取请求路径携带的参数*/
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		key := "q"
		query := r.URL.Query()
		value := query.Get(key)
		fmt.Printf("获取的参数为：%s\n", value)
	})
	http.ListenAndServe(":8080", nil)
}
