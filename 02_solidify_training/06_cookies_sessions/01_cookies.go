package main

import (
	"fmt"
	"log"
	"net/http"
)

func cookiesOperate(resp http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("my-cookie")
	if err != nil {
		log.Fatalf("get cookies is fault")
	}
	fmt.Printf("fetch cookies is %s", cookie.Value)
	http.SetCookie(resp, &http.Cookie{
		Name:  "my-cookie",
		Value: "some value",
	})
}
func main() {
	http.HandleFunc("/", cookiesOperate)
	http.ListenAndServe(":8080", nil)
}
