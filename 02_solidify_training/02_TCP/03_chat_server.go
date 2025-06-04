package main

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

func (cs *ChatServer) ChatInput(user *User, input chan *Message) {
	message := &Message{
		Username: user.Name,
		Text:     "Welcome to the chat!",
	}
	input <- message
}

func (cs *ChatServer) Run() {
	for {

	}
}
func main() {

}
