package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type eventStream struct {
	clients map[chan []Request]bool
	mutex   sync.Mutex
}

var eventMapData []Request

func main() {
	stream := &eventStream{
		clients: make(map[chan []Request]bool),
	}

	http.HandleFunc("/events", stream.handleEvents)
	http.HandleFunc("/send-event", stream.handleSendEvent)

	port := ":8080"
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func (s *eventStream) handleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	if len(eventMapData) < 1 {
		eventMapData = append(eventMapData, Request{
			Text: "Pesanan Sedang Dikemas",
			Time: time.Now(),
		})
	}
	jsonData, _ := json.Marshal(eventMapData)
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.(http.Flusher).Flush()
}

func (s *eventStream) handleSendEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	request := Request{
		Text: req.Text,
		Time: time.Now(),
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	eventMapData = append(eventMapData, Request{
		Text: request.Text,
		Time: request.Time,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Event sent successfully"))
}

type Request struct {
	Text string    `json:"text"`
	Time time.Time `json:"time"`
}
