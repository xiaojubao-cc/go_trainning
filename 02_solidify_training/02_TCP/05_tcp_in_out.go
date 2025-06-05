package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

/*将网络接收的消息打印到终端*/
func main() {
	listener, err := net.Listen("tcp", ":8080")
	defer listener.Close()
	if err != nil {
		fmt.Errorf(err.Error())
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Errorf(err.Error())
		}
		/*这里必须使用协程因为io.copy是阻塞性质的方法*/
		go io.Copy(os.Stdout, conn)
		/*将终端的数据回写客户端*/
		go io.Copy(conn, os.Stdin)
	}
}
