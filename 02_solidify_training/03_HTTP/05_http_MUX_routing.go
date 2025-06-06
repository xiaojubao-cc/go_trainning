package main

import (
	"io"
	"net/http"
)

type CatHandler int

type DogHandler int

func (c CatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Cat")
}

func (d DogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Dog")
}
func main() {
	var catHandler CatHandler
	var dogHandler DogHandler

	serveMux := http.NewServeMux()
	/*/dog/和/dog区别后者是精确匹配,前者会匹配子路由*/
	serveMux.Handle("/cat/", catHandler)
	serveMux.Handle("/dog/", dogHandler)
	http.ListenAndServe(":8080", serveMux)
}
