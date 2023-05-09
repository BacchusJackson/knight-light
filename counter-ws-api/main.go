package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {

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
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
		}

		log.Printf("[%s] message: %s\n", c.Locals("Host"), msg)

		if err = c.WriteMessage(mt, []byte("Howdy, from Fiber")); err != nil {
			log.Println(err)
		}
	}
}
