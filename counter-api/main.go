package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

var globalCounter *Counter

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	globalCounter = &Counter{
		CurrentCount: 0,
		LastUpdated:  time.Now(),
	}

	http.HandleFunc("/add", PostAdd)
	http.HandleFunc("/status", GetStatus)

  if err := http.ListenAndServe(":" + port, nil); err != nil {
    panic(err)
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

	if err := json.NewEncoder(w).Encode(globalCounter); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Add Request Encode Error:", err)
		return
	}
}
