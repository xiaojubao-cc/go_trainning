package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

var syncMap sync.Map

func handleRedisOrder(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fields := strings.Fields(text)
		if len(fields) < 1 {
			io.WriteString(conn, "INVALID COMMAND\n")
			continue
		}
		switch fields[0] {
		case "GET":
			if len(fields) != 2 {
				io.WriteString(conn, "INVALID COMMAND\n")
				continue
			}
			value, ok := syncMap.Load(fields[1])
			if ok {
				builder := strings.Builder{}
				builder.WriteString("data:")
				builder.WriteString(value.(string))
				builder.WriteString("\n")
				io.WriteString(conn, builder.String())
			} else {
				io.WriteString(conn, "NOT FOUND\n")
			}

		case "SET":
			if len(fields) < 3 {
				io.WriteString(conn, "INVALID COMMAND\n")
				continue
			}
			syncMap.Store(fields[1], fields[2])
			io.WriteString(conn, "OK")
		case "DEL":
			if len(fields) != 2 {
				io.WriteString(conn, "INVALID COMMAND\n")
				continue
			}
			syncMap.Delete(fields[1])
			io.WriteString(conn, "OK")
		default:
			io.WriteString(conn, "INVALID COMMAND "+fields[0]+"\n")
		}
	}
}

/*监听客户端的输入*/
func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln(r)
		}
	}()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("accept exception")
		}
		go handleRedisOrder(conn)
	}
}
