package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var MessageChan chan Message
var CancelChan chan struct{}
var UserChan chan User

func main() {

	MessageChan = make(chan Message)
	CancelChan = make(chan struct{})
	UserChan = make(chan User)

	go chatListener(MessageChan, CancelChan, UserChan)

	port := 3000
	app := fiber.New()

	app.Use("/ws", wsMiddleware)
	app.Get("/ws", websocket.New(handleWebSocket))

	app.Get("/health", healthCheck)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))

}

func wsMiddleware(c *fiber.Ctx) error {
	if c.Get("host") == "localhost:3000" {
		c.Locals("Host", "Localhost:3000")
		return c.Next()
	}
	return c.Status(403).SendString("Request origin not allowed")
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"status": "good"})
}

func handleWebSocket(c *websocket.Conn) {

	err := c.WriteMessage(websocket.TextMessage, []byte("Welcome to Knight Light Mono Chat! What is your username?"))
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

	if err := c.WriteMessage(websocket.TextMessage, []byte("Username saved. Happy chatting!")); err != nil {
		log.Println(err)
		return
	}
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
					if msg.user == user.username || user.alive == false {
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
			for _, user := range users {
				if user.username == newUser.username {
          log.Printf("User: %s already registered, setting to alive!\n", user.username)
					user.alive = true
				}
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
