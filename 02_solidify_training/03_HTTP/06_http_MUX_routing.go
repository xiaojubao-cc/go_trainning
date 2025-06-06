package main

import (
	"io"
	"net/http"
)

func HandleDog(res http.ResponseWriter, r *http.Request) {
	io.WriteString(res, "doggy doggy doggy")
}

func HandleCat(res http.ResponseWriter, r *http.Request) {
	io.WriteString(res, "catty catty catty")
}
func main() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/cat/", HandleCat)
	serveMux.HandleFunc("/dog/", HandleDog)
	http.ListenAndServe("127.0.0.1:8080", serveMux)
}
