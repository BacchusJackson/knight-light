package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
)

var MessageChan chan Message
var CancelChan chan struct{}

var users *Users

func main() {
	addr := flag.String("addr", ":5000", "http address and port")
	flag.Parse()

	MessageChan = make(chan Message)

	users = &Users{users: make([]*User, 0)}

	go messageWorker()

	app := fiber.New()

	app.Use(logger.New())
	// Allow Localhost request origin connection
	app.Use(cors.New())

	app.Get("/ws", websocket.New(handleWebSocket))

	app.Get("/health", func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{"status": "Healthy"})
	})

	log.Fatal(app.Listen(*addr))

}

func handleWebSocket(c *websocket.Conn) {

	err := c.WriteMessage(websocket.TextMessage, []byte("[server] Welcome to Knight Light Mono Chat! What is your username?"))
	if err != nil {
		log.Println(err)
		return
	}
	_, msgBytes, err := c.ReadMessage()

	if err != nil {
		log.Println(err)
		return
	}

	msg := string(msgBytes)
	runes := []rune(msg)
	if len(runes) > 20 {
		runes = runes[:20]
	}
	username := string(runes)

	user := users.Register(username, c)

	MessageChan <- NewMessage("server", []byte(fmt.Sprintf("User %s connected\n", user.username)))

	for {
		log.Println(user.username, ": Starting Message Loop")
		_, msgBytes, err := c.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}

		msg := NewMessage(user.username, msgBytes)
		log.Println(msg)

		MessageChan <- NewMessage(username, msgBytes)
	}

}

func messageWorker() {
	for {
		receivedMsg := <-MessageChan
		users.SendAll(receivedMsg)
	}
}

type User struct {
	username string
	conn     *websocket.Conn
	alive    bool
}

func (u *User) Send(msg Message) {
	err := u.conn.WriteMessage(websocket.TextMessage, msg.Bytes())
	if err != nil {
    u.alive = false
		log.Printf("User: %+v error: %s\n", u, err)
	}
}

type Users struct {
	users []*User
	mu    sync.Mutex
}

func (u *Users) Register(username string, conn *websocket.Conn) *User {

	u.mu.Lock()
	defer u.mu.Unlock()
	log.Printf("Registering username: %s\n", username)

	for i := range u.users {
		if username != u.users[i].username {
			continue
		}

		u.users[i].alive = true
		log.Printf("Existing user: %s\n", username)
		return u.users[i]
	}

	newUser := &User{username: username, conn: conn, alive: true}
	u.users = append(u.users, newUser)

	log.Printf("Registered username: %s\n", username)
	return newUser
}

func (u *Users) SendAll(msg Message) {
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, user := range u.users {
		user.Send(msg)
	}
}

type Message struct {
	user string
	body string
}

func (m Message) String() string {
	return fmt.Sprintf("[%s]: %s", m.user, m.body)
}

func (m Message) Bytes() []byte {
	return []byte(m.String())
}

func NewMessage(username string, bodyBytes []byte) Message {
	return Message{
		user: username,
		body: string(bodyBytes),
	}
}
