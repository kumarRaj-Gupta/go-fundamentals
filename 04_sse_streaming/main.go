package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const port = ":8081"

// handleStream handles a request by continuously streaming events.
func handleStream(w http.ResponseWriter, r *http.Request) {
	log.Println("Client connected for streaming...")

	// 1. --- SET MANDATORY SSE HEADERS ---
	// Your code here: Set Content-Type, Cache-Control, and Connection headers.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// 2. --- GET THE FLUSHER INTERFACE ---
	// Your code here: Perform type assertion and check for success.
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Fatal("The ResponseWriter doesn't support Flusher")
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}
	// If the connection doesn't support flushing, we can't stream.
	// if !ok { ... log error ... return }

	// --- CORE STREAMING LOGIC (Will be completed in Task 2) ---

	eventID := 0
	for i := 0; i < 3; i++ {
		eventID++

		time.Sleep(time.Second)
		// The SSE format requires a "data: " prefix followed by a newline, and then a double newline (\n\n) to end the event.
		// For simplicity here, we only use the data field.
		data := fmt.Sprintf("Tool Execution ID: %d at %s", eventID, time.Now().Format("03:04:05 PM"))

		// Write the event data to the ResponseWriter
		fmt.Fprintf(w, "data: %s\n\n", data)

		// CRITICAL: Force the data packet to be sent immediately.
		flusher.Flush()
	}

	log.Println("Stream finished and connection closed.")
}

// main sets up the server and routing.
func main() {
	fmt.Printf("--- Project 4: SSE Streaming (net/http & http.Flusher) ---\n")

	// Route all requests to /events to our streaming handler
	http.HandleFunc("/events", handleStream)

	fmt.Printf("Starting SSE server on port %s...\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
