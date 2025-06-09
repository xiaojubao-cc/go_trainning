package main

import (
	"io"
	"net/http"
	"os"
)

func fetchFile(res http.ResponseWriter, req *http.Request) {
	file, err := os.Open("D:\\go_projects\\go_trainning\\02_solidify_training\\03_HTTP\\toby.jpg")
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		defer file.Close()
		io.CopyBuffer(res, file, make([]byte, 1024))
		//使用零拷贝技术、智能头处理、并发控制优化
		//http.ServeFile(res, req, "D:\\go_projects\\go_trainning\\02_solidify_training\\03_HTTP\\toby.jpg")
	}
}

func main() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", fetchFile)
	http.ListenAndServe(":8080", serveMux)
}
