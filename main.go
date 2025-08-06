package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/heroku/x/dynoid"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func init() {
	// Get Redis URL from environment variable
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Println("REDIS_URL not set, Redis functionality will be disabled")
		return
	}

	// Parse Redis URL and create client
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("Error parsing REDIS_URL: %v", err)
		return
	}

	// Configure TLS for secure Redis connections (like Heroku Redis)
	if opts.TLSConfig != nil {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: true, // Skip certificate verification for Heroku Redis
		}
	}

	redisClient = redis.NewClient(opts)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Error connecting to Redis: %v", err)
		redisClient = nil
	} else {
		log.Println("Successfully connected to Redis")
	}
}

func main() {
	// Get port from environment variable (Heroku sets PORT)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Define routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get dyno ID using the heroku/x/dynoid package
		// Using "applink" as the audience
		token, err := dynoid.ReadLocal("applink")
		if err != nil {
			log.Printf("Error getting dyno ID (running locally?): %v", err)
			token = "local-dev" // Use a default value for local development
		}

		// Redis operations
		var redisInfo string
		if redisClient != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Increment a counter
			counter, err := redisClient.Incr(ctx, "request_counter").Result()
			if err != nil {
				log.Printf("Redis error: %v", err)
				redisInfo = "Redis error occurred"
			} else {
				// Store dyno ID with timestamp
				timestamp := time.Now().Format(time.RFC3339)
				key := fmt.Sprintf("dyno:%s:last_request", token)
				err = redisClient.Set(ctx, key, timestamp, 24*time.Hour).Err()
				if err != nil {
					log.Printf("Redis set error: %v", err)
				}

				redisInfo = fmt.Sprintf("Request #%d from dyno %s at %s", counter, token, timestamp)
			}
		} else {
			redisInfo = "Redis not available"
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
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("User-Agent", "heroku-app/1.0")

		// Add authorization header with dyno ID as bearer token
		if token != "" {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
			//			req.Header.Set("Authorization", "Bearer "+dynoID)
			if redisClient != nil {
				redisClient.Set(context.Background(), "dynoID", token, 0)
			}
		} else {
			if redisClient != nil {
				redisClient.Set(context.Background(), "dynoID", "unknown", 0)
			}
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

		// Create response with Redis info
		response := fmt.Sprintf("External API Response: %s\n\nRedis Info: %s", string(body), redisInfo)

		// Set response headers
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		// Write response body
		w.Write([]byte(response))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		status := "OK"
		if redisClient != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			_, err := redisClient.Ping(ctx).Result()
			if err != nil {
				status = "Redis connection failed"
			}
		}
		fmt.Fprintf(w, status)
	})

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
