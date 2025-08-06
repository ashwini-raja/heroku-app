package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/heroku/x/dynoid"
)

func main() {
	// Get port from environment variable (Heroku sets PORT)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Define routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get dyno ID using the heroku/x/dynoid package
		// Using "applink.staging.herokudev.com" as the audience
		dynoID, err := dynoid.ReadLocal("applink")
		if err != nil {
			log.Printf("Error getting dyno ID: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Create HTTP client
		client := &http.Client{}

		// Create request to the specified URL
		req, err := http.NewRequest("GET", "https://applink.staging.herokudev.com/up", nil)
		if err != nil {
			log.Printf("Error creating request: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Add authorization header with dyno ID as bearer token
		if dynoID != "" {
			req.Header.Set("Authorization", "Bearer "+dynoID)
			req.Header.Set("Content-Type", "application/json")
		}

		// Make the request
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error making request: %v", err)
			http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response: %v", err)
			http.Error(w, "Failed to read response", http.StatusInternalServerError)
			return
		}

		// Set response headers
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(resp.StatusCode)

		// Write response body
		w.Write(body)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
