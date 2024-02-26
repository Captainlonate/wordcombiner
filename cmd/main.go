package main

import (
	"captainlonate/wordcombiner/internal/conceptsdb"
	"captainlonate/wordcombiner/internal/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Required Environment Variables - Must be set in .env file
var requiredEnvVariables = []string{
	"OPENAI_API_KEY", "OPENAI_MODEL",
}

func main() {
	// Load environment variables from .env file (uses .env by default)
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Check if required environment variables exist and have non-empty values
	for _, key := range requiredEnvVariables {
		value := os.Getenv(key)
		if value == "" {
			fmt.Printf("Error: Environment variable %s is required but not set\n", key)
			return
		}
	}

	// Establish a connection to the redis database. If it fails, the program will exit.
	conceptsdb.InitializeRedisConnection()

	// Create a new router and start the web server
	router := routes.NewRouter()
	addr := fmt.Sprintf(":%d", 8080)
	fmt.Printf("Server listening on %s...\n", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
