package main

import (
	"fmt"
	"io"
	"net"
)

/*使用命令访问curl --raw telnet://localhost:8080*/
func main() {
	listener, err := net.Listen("tcp", ":8080")
	defer listener.Close()
	if err != nil {
		panic(err)
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Errorf("listen exception")
		}
	}()
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Errorf("accept exception")
		}
		io.WriteString(conn, "hello world")
		conn.Close()
	}
}
