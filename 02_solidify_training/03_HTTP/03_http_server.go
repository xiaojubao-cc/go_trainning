package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic:", err)
		}
	}()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("listen exception")
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("accept exception")
		}
		go handleConnServerPost(conn)
	}
}

func handleConnServerPost(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic:", err)
		}
	}()
	defer conn.Close()
	i := 0
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
		if i == 0 {
			method := strings.Fields(ln)[0]
			log.Printf("METHOLD: %s", method)
		} else {
			if len(ln) == 0 {
				break
			}
		}
		i++
	}
	//response返回响应使用带缓冲的性能更好
	// 构建响应
	body := `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title></title>
</head>
<body>
	<form method="POST">
		<input type="text" name="key" value="">
		<input type="submit">
	</form>
</body>
</html>
	`
	writer := bufio.NewWriter(conn)
	defer writer.Flush()
	// 写入响应头
	writer.WriteString("HTTP/1.1 302 OK\r\n")
	fmt.Fprintf(writer, "Content-Length: %d\r\n", len(body))
	/*跳转路径Location和状态码有关(302,303)*/
	writer.WriteString("Location: http://www.baidu.com\r\n")
	/*非文本*/
	//writer.WriteString("Content-Type: text/plain\r\n")
	writer.WriteString("\r\n")
	// 写入响应正文
	writer.WriteString(body)
}
