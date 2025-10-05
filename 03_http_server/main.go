package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Global port for the server
const port = ":8080"

// --- DATA STRUCTURES (Simulates the MCP Execute payload) ---

type ToolExecution struct {
	ToolName  string                 `json:"tool_name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// --- HANDLERS ---

// handleHealth (omitted for brevity)
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status": "ok", "message": "Server is healthy"}`)
}

// handleData (omitted for brevity)
func handleData(w http.ResponseWriter, r *http.Request) {
	type ServerInfo struct {
		Name string `json:"server_name"`
	}
	info := ServerInfo{Name: "MCP-Go-Server"}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// handleQuery (omitted for brevity)
func handleQuery(w http.ResponseWriter, r *http.Request) {
	toolName := r.URL.Query().Get("tool")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "success", "requested_tool": "%s"}`, toolName)
}

// handleSubmit handles the incoming JSON payload from the request body (POST method).
func handleSubmit(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Define a target struct variable named 'execution'
	var execution ToolExecution

	// 2. Decode the JSON body from r.Body into the 'execution' variable
	err := json.NewDecoder(r.Body).Decode(&execution)

	// 3. Check for decoding error and respond with a 400 Bad Request if invalid JSON is provided
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Bad Request: Invalid JSON or incorrect format.", http.StatusBadRequest)
		return
	}

	// IMPORTANT: r.Body is an io.ReadCloser. We should close it after reading the data.
	// Since json.NewDecoder().Decode() reads until EOF, we can defer the close now.
	defer r.Body.Close()

	// SUCCESS RESPONSE (Only executed if decoding was successful)
	w.Header().Set("Content-Type", "application/json")

	// We check for nil arguments to prevent a panic if the map is not initialized.
	argsCount := 0
	if execution.Arguments != nil {
		argsCount = len(execution.Arguments)
	}

	fmt.Fprintf(w, `{"status": "received", "tool": "%s", "args_count": %d}`,
		execution.ToolName, argsCount)
}

// --- MAIN FUNCTION AND ROUTING ---

func main() {
	fmt.Printf("--- Project 3: HTTP Server Fundamentals (net/http) ---\n")

	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/data", handleData)
	http.HandleFunc("/query", handleQuery)
	http.HandleFunc("/execute", handleSubmit) // New handler routed

	fmt.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
