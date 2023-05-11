package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var globalCounter *Counter

func main() {

	addr := flag.String("addr", ":5000", "http address and port")
	flag.Parse()

	log.SetOutput(os.Stdout)

	var err error

	if err != nil {
		panic(err)
	}

	globalCounter = &Counter{
		CurrentCount: 0,
		LastUpdated:  time.Now(),
	}

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/status", func(c *fiber.Ctx) error {
		return c.JSON(globalCounter)
	})

	app.Get("/heath", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthly"})
	})

	app.Post("/add", func(c *fiber.Ctx) error {
		requestObj := new(AddRequest)

		if err := c.BodyParser(requestObj); err != nil {
			return err
		}

		globalCounter.CurrentCount += requestObj.Count
		globalCounter.LastUpdated = time.Now()

		return c.JSON(globalCounter)
	})

	log.Fatal(app.Listen(*addr))
}

type AddRequest struct {
	Count int `json:"count"`
}

type Counter struct {
	CurrentCount int       `json:"currentCount"`
	LastUpdated  time.Time `json:"lastUpdated"`
}
