package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

var counterChan chan Counter
var connectionChan chan *Connection

var globalCounter Counter
var globalCounterMu sync.RWMutex

func main() {

	addr := flag.String("addr", ":5000", "http address and port")
	flag.Parse()

	log.SetOutput(os.Stdout)

	var err error

	if err != nil {
		panic(err)
	}

  counterChan = make(chan Counter)
  connectionChan = make(chan *Connection)

  go updater()

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/status", func(c *fiber.Ctx) error {

		return c.JSON(globalCounterCopy())
	})

	app.Get("/heath", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthly"})
	})

	app.Post("/add", func(c *fiber.Ctx) error {
		requestObj := &AddRequest{}

		if err := json.Unmarshal(c.Body(), requestObj); err != nil {
			log.Println(err)
			return err
		}

		log.Printf("Request Object: %+v\n", requestObj)

		counter := Counter{
			CurrentCount: requestObj.Count,
			LastUpdated:  time.Now(),
		}

		counterChan <- counter

		return c.JSON(&counter)
	})

	app.Get("/sse", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {

			killChan := make(chan struct{})

			conn := &Connection{
				id:       uuid.New(),
				writer:   w,
				killChan: killChan,
        alive: true,
			}
      fmt.Println("Connection Request:", conn.id.String())

			connectionChan <- conn

			// Wait until connection dies
			<-killChan

		}))

		return nil
	})

	log.Fatal(app.Listen(*addr))
}

type Connection struct {
	id       uuid.UUID
	writer   *bufio.Writer
	killChan chan struct{}
  alive bool
}

func globalCounterCopy() Counter {
	counter := Counter{}
	globalCounterMu.RLock()
	counter.CurrentCount = globalCounter.CurrentCount
	counter.LastUpdated = globalCounter.LastUpdated
	globalCounterMu.RUnlock()
	return counter
}

func updateGlobalCounter(n int, updated time.Time) {
	globalCounterMu.Lock()
	globalCounter.CurrentCount += n
	globalCounter.LastUpdated = updated
	globalCounterMu.Unlock()
}

func updater() {
  log.Println("Updater Started")
	connections := make([]*Connection, 0)

	for {
		select {
		case counter := <-counterChan:
			updateGlobalCounter(counter.CurrentCount, counter.LastUpdated)
			sendUpdates(connections)
		case newConn := <-connectionChan:
			connections = append(connections, newConn)
      log.Println("New Connection Registered", newConn.id.String())
		}
	}
}

func sendUpdates(connections []*Connection) {
  fmt.Println("Sending Updates!")

	for i := range connections {
    if !connections[i].alive {
      continue
    }

      if err := writeStatus(connections[i].writer); err != nil {
      log.Println("Failed to write to", connections[i].id)
      connections[i].alive = false
			connections[i].killChan <- struct{}{}

		}
	}
}

func writeStatus(w *bufio.Writer) error {
    log.Println("Writing Status...")
	  counter := globalCounterCopy()
		buf := new(bytes.Buffer)
		_ = json.NewEncoder(buf).Encode(counter)

		fmt.Fprintf(w, "data: %s\n\n", buf.String())
		return w.Flush()
}

// remove without consideration of order
func remove[T any](index int, items []T) {
  items = append(items, items[index])
  items = append(items[:index], items[index+1:len(items)-1]...)
}

type AddRequest struct {
	Count int `json:"count"`
}

type Counter struct {
	CurrentCount int       `json:"currentCount"`
	LastUpdated  time.Time `json:"lastUpdated"`
}
