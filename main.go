package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Global backends list
var backends []string
var mu sync.RWMutex // Mutex to handle concurrent read/write

// Struct to hold the response or error from each backend
type result struct {
	response *http.Response
	err      error
	backend  string
}

// Function to load backends from environment variable at startup
func loadBackendsFromEnv() {
	backendEnv := os.Getenv("BACKEND_SERVERS")
	if backendEnv == "" {
		log.Fatal("Environment variable BACKEND_SERVERS is not set")
	}

	// Split the environment variable by commas
	backends = strings.Split(backendEnv, ",")
	log.Printf("Loaded backends from environment: %v", backends)
}

// API endpoint to update the backend list at runtime
func updateBackendsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse new backends from request body
	newBackends := r.URL.Query().Get("backends")
	if newBackends == "" {
		http.Error(w, "No backends provided", http.StatusBadRequest)
		return
	}

	// Update backends list
	mu.Lock()
	backends = strings.Split(newBackends, ",")
	mu.Unlock()

	log.Printf("Updated backends at runtime: %v", backends)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Backends updated successfully")
}

// Function to forward requests to a single backend
func forwardRequest(backend string, r *http.Request, resultChan chan<- result) {
	req, err := http.NewRequest(r.Method, backend+r.RequestURI, r.Body)
	if err != nil {
		resultChan <- result{nil, err, backend}
		return
	}

	for k, v := range r.Header {
		req.Header[k] = v
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	resultChan <- result{resp, err, backend}
}

// Proxy handler to forward requests to backends
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.RequestURI)

	resultChan := make(chan result, len(backends))
	var wg sync.WaitGroup

	mu.RLock() // Read lock while accessing backends
	for _, backend := range backends {
		wg.Add(1)
		go func(backend string) {
			defer wg.Done()
			forwardRequest(backend, r, resultChan)
		}(backend)
	}
	mu.RUnlock() // Release the read lock

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for res := range resultChan {
		if res.err == nil && res.response.StatusCode >= 200 && res.response.StatusCode < 300 {
			log.Printf("Success from backend %s: %s", res.backend, res.response.Status)
			for k, v := range res.response.Header {
				w.Header()[k] = v
			}
			w.WriteHeader(res.response.StatusCode)
			io.Copy(w, res.response.Body)
			res.response.Body.Close()
			return
		} else if res.err != nil {
			log.Printf("Error from backend %s: %v", res.backend, res.err)
		} else {
			log.Printf("Non-success status from backend %s: %s", res.backend, res.response.Status)
		}
	}

	http.Error(w, "All backend requests failed", http.StatusBadGateway)
}

func main() {
	loadBackendsFromEnv() // Load backends at startup

	// Explicitly handle update-backends route first
	http.HandleFunc("/-update-backends", updateBackendsHandler)

	// Catch-all proxy handler for other routes
	http.HandleFunc("/", proxyHandler)

	log.Println("Starting proxy server on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
