package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	/*使用dial访问8080*/
	dialConn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer dialConn.Close()
	response, err := io.ReadAll(dialConn)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(response))
}
