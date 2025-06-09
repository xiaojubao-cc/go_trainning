package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

/*访问路径127.0.0.1:8080重定向到https://127.0.0.1:10443/*/

const (
	FORWARDURL = "https://127.0.0.1:10443"
)

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}
func forwardHttp(resp http.ResponseWriter, req *http.Request) {
	joinPath, err := url.JoinPath(FORWARDURL, strings.TrimPrefix(req.URL.Path, "/"))
	if err != nil {
		log.Fatalf("join path err：%s\n", err)
	}
	http.Redirect(resp, req, joinPath, http.StatusMovedPermanently)
}
func main() {
	http.HandleFunc("/", handler)
	go http.ListenAndServe(":8080", http.HandlerFunc(forwardHttp))
	err := http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
	if err != nil {
		fmt.Println("ListenAndServeTLS error:", err)
	}
}
