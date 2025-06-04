package main

import (
	"fmt"
	"net/http"
)

// Handler /*路由*/
type Handler struct {
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("实现handler路由\n")
}
func main() {
	/*使用handler结构体*/
	//http.Handle("/", &Handler{})
	//http.ListenAndServe(":8080", nil)

	/*使用handleFunc*/
	http.HandleFunc("/handleFunc", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("实现handlerFunc路由\n")
	})
	http.ListenAndServe("localhost:8080", nil)
	fmt.Println("启动服务")
}
