package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

/*header输出格式如下*/
//2025/06/06 10:18:53 GET / HTTP/1.1
//2025/06/06 10:18:53 Host: localhost:8080
//2025/06/06 10:18:53 User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0
//2025/06/06 10:18:53 Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
//2025/06/06 10:18:53 Accept-Language: zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2
//2025/06/06 10:18:53 Accept-Encoding: gzip, deflate, br, zstd
//2025/06/06 10:18:53 Connection: keep-alive
//2025/06/06 10:18:53 Cookie: Hm_lvt_a1ff8825baa73c3a78eb96aa40325abc=1734505022,1735095595; Pycharm-8cbd502e=f6787538-7450-477f-9e62-8095d1b39b40
//2025/06/06 10:18:53 Upgrade-Insecure-Requests: 1
//2025/06/06 10:18:53 Sec-Fetch-Dest: document
//2025/06/06 10:18:53 Sec-Fetch-Mode: navigate
//2025/06/06 10:18:53 Sec-Fetch-Site: none
//2025/06/06 10:18:53 Sec-Fetch-User: ?1
//2025/06/06 10:18:53 Priority: u=0, i
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
		go handleConnServerGET(conn)
	}
}

func handleConnServerGET(conn net.Conn) {
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
	body := "hello golang"
	writer := bufio.NewWriter(conn)
	defer writer.Flush()
	// 写入响应头
	writer.WriteString("HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(writer, "Content-Length: %d\r\n", len(body))
	writer.WriteString("Content-Type: text/plain\r\n")
	writer.WriteString("\r\n")
	// 写入响应正文
	writer.WriteString(body)
}
