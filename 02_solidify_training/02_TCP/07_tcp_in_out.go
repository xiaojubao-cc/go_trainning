package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

var data sync.Map

type Command struct {
	Fields []string
	Result chan string
}

/*服务端处理*/
func handleRedisServer(commands chan Command) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("handleRedisServer panic:", err)
		}
	}()
	//取出信道中的数据进行处理
	for command := range commands {
		fields := command.Fields
		if len(fields) == 0 {
			command.Result <- "command error"
		}
		switch fields[0] {
		case "GET":
			if len(fields) != 2 {
				command.Result <- "command error"
			}
			value, ok := data.Load(fields[1])
			if ok {
				command.Result <- value.(string)
			} else {
				command.Result <- "key not found"
			}

		case "SET":
			if len(fields) != 3 {
				command.Result <- "command error"
			}
			data.Store(fields[1], fields[2])
			command.Result <- "OK"

		case "DEL":
			data.Delete(fields[1])
			command.Result <- "OK"

		default:
			command.Result <- "command error"
		}
	}
}

/*客户端处理*/
func handleRedisClient(commands chan Command, conn net.Conn) {
	defer conn.Close()
	/*将客户端的输入全部放入自己的信道*/
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		ch := make(chan string)
		fields := strings.Fields(text)
		commands <- Command{fields, ch}
		io.WriteString(conn, <-ch)
	}

}
func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("listen exception")
	}
	defer listener.Close()
	commands := make(chan Command)
	go handleRedisServer(commands)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("accept exception")
		}
		go handleRedisClient(commands, conn)
	}
}
