package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		file, _, err := req.FormFile("q")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		io.Copy(os.Stdout, file)
		res.Header().Set("Content-Type", "text/html")
		io.WriteString(res, `
			<form method="POST" enctype="multipart/form-data">
			  <input type="file" name="q">
			  <input type="submit">
			</form>`)
	})
	http.ListenAndServe(":8080", nil)
}
