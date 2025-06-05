package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type User struct {
	Name   string
	Output chan Message
}

type Message struct {
	Username string
	Text     string
}

type ChatServer struct {
	Users map[string]User
	Join  chan User
	Leave chan User
	Input chan Message
}

func (cs *ChatServer) Run() {
	for {
		select {
		case user := <-cs.Join:
			cs.Users[user.Name] = user
			go func() {
				message := Message{
					Username: "SYSTEM",
					Text:     fmt.Sprintf("Welcome to %s the chat!", user.Name),
				}
				cs.Input <- message
			}()
		case user := <-cs.Leave:
			delete(cs.Users, user.Name)
			go func() {
				message := Message{
					Username: "SYSTEM",
					Text:     fmt.Sprintf("Bye %s the chat!", user.Name),
				}
				cs.Input <- message
			}()
		case message := <-cs.Input:
			for _, user := range cs.Users {
				/*使用select属于非阻塞写入数据*/
				select {
				case user.Output <- message:
				default:
				}
			}
		}
	}
}

/*客户端监听操作*/
func handleConn(chatServer *ChatServer, conn net.Conn) {
	defer conn.Close()
	io.WriteString(conn, "Enter your name: ")
	/*客户端监听输入*/
	scanner := bufio.NewScanner(conn)
	/*阻塞监听客户端的输入*/
	scanner.Scan()
	user := User{
		Name:   scanner.Text(),
		Output: make(chan Message, 2),
	}
	chatServer.Join <- user
	defer func() {
		chatServer.Leave <- user
	}()
	go func() {
		/*持续监听客户端的输入*/
		for scanner.Scan() {
			text := scanner.Text()
			chatServer.Input <- Message{Username: user.Name, Text: text}
		}
	}()
	/*消息群发*/
	for message := range user.Output {
		if user.Name != message.Username {
			_, err := io.WriteString(conn, message.Username+": "+message.Text+"\n")
			if err != nil {
				break
			}
		}
	}
}
func main() {
	listener, err := net.Listen("tcp", ":8080")
	defer listener.Close()
	if err != nil {
		fmt.Errorf(err.Error())
	}
	chatServer := &ChatServer{
		Users: make(map[string]User),
		Join:  make(chan User),
		Leave: make(chan User),
		Input: make(chan Message),
	}
	/*启动服务端*/
	go chatServer.Run()
	for {
		/*接受客户端请求*/
		conn, err := listener.Accept()
		if err != nil {
			fmt.Errorf("accept exception")
		}
		/*为每个客户端连接启用一个协程*/
		go handleConn(chatServer, conn)
	}
}
