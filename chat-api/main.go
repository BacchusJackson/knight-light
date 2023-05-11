package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var MessageChan chan Message
var CancelChan chan struct{}
var UserChan chan User

func main() {
	addr := flag.String("addr", ":5000", "http address and port")
	flag.Parse()

	MessageChan = make(chan Message)
	CancelChan = make(chan struct{})
	UserChan = make(chan User)

	go chatListener(MessageChan, CancelChan, UserChan)

	app := fiber.New()

	app.Use("/ws", wsMiddleware)
	app.Get("/ws", websocket.New(handleWebSocket))

	app.Get("/health", func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{"status": "Healthy"})
	})

	log.Fatal(app.Listen(*addr))

}

func wsMiddleware(c *fiber.Ctx) error {
	if c.Get("host") == "localhost:3000" {
		c.Locals("Host", "Localhost:3000")
		return c.Next()
	}
	return c.Status(403).SendString("Request origin not allowed")
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

	UserChan <- User{username: username, conn: c, alive: true}

	for {
		_, msgBytes, err := c.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}

		MessageChan <- NewMessage(username, msgBytes)
	}
}

// ChatListener will listen for new messages and send them to other user socket connections
func chatListener(msgChan chan Message, cancelChan chan struct{}, userChan chan User) {
	users := []User{}

	for {
		select {
		case msg := <-msgChan:
			go func() {
				for _, user := range users {
					if user.alive == false {
						continue
					}
					err := user.conn.WriteMessage(websocket.TextMessage, msg.Bytes())
					if err != nil {
						log.Printf("Message to user: %s failed: %v. User is no longer alive\n", user.username, err)
						user.alive = false
					}
				}

			}()
		case <-cancelChan:
			log.Printf("Chat Listener Cancelled")
			return
		case newUser := <-userChan:
			msg := "[server] Happy chatting!"

			for _, user := range users {
				if user.username == newUser.username {
					log.Printf("User: %s already registered, setting to alive!\n", user.username)
					user.alive = true
					msg = fmt.Sprintf("[server] Welcome back %s", user.username)
					newUser = user
					break
				}
			}

			if err := newUser.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				log.Println(err)
				return
			}

			users = append(users, newUser)
		}
	}
}

type User struct {
	username string
	conn     *websocket.Conn
	alive    bool
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
