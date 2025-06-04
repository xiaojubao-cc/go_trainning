package main

import (
	"fmt"
	"io"
	"net/http"
)

/*client*/
func main() {
	response, err := http.Get("http://localhost:8080")
	if err != nil {
		fmt.Errorf("服务器异常：%s", err)
	}
	defer response.Body.Close()
	content, _ := io.ReadAll(response.Body)
	fmt.Println(string(content))

}
