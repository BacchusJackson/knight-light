package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var globalCounter *Counter


func main() {
  log.SetOutput(os.Stdout)

	globalCounter = &Counter{
		CurrentCount: 0,
		LastUpdated:  time.Now(),
	}

	http.HandleFunc("/add", PostAdd)
	http.HandleFunc("/status", GetStatus)

  port, err := strconv.Atoi(os.Getenv("APP_PORT"))
  if err != nil {
    log.Fatalf("APP_PORT env variable must be set like '5000': %v\n", err)
  }

  log.Printf("Starting Counter API on port %d", port) 

  if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
    log.Fatal(err)
	}

}

type AddRequest struct {
	Count int `json:"count"`
}

type Counter struct {
	CurrentCount int       `json:"currentCount"`
	LastUpdated  time.Time `json:"lastUpdated"`
}

func GetStatus(w http.ResponseWriter, r *http.Request) {
  log.Printf("GET status request received from %s\n", r.RemoteAddr)
	writeStatus(w)
}

func PostAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Invalid Method:", r.Method)
		return
	}

	requestObj := &AddRequest{}
	if err := json.NewDecoder(r.Body).Decode(requestObj); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Add Request Decode Error:", err)
		return
	}

	log.Printf("Request Object: %+v\n", requestObj)

	globalCounter.CurrentCount += requestObj.Count
	globalCounter.LastUpdated = time.Now()

	writeStatus(w)
}

func writeStatus(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
  addHeaders(w)

	if err := json.NewEncoder(w).Encode(globalCounter); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Add Request Encode Error:", err)
		return
	}
}


func addHeaders(w http.ResponseWriter) {
  if os.Getenv("APP_ENV") == "local" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
  }
}


