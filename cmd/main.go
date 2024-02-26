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

func main() {
	// Load environment variables from .env file (uses .env by default)
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Check if required environment variables exist and have non-empty values
	requiredVariables := []string{"OPENAI_API_KEY", "OPENAI_MODEL"}
	for _, key := range requiredVariables {
		value := os.Getenv(key)
		if value == "" {
			fmt.Printf("Error: Environment variable %s is required but not set\n", key)
			return
		}
	}

	// Establish a connection to the redis database
	conceptsdb.InitializeRedisConnection()

	// Create a new router and start the web server
	router := routes.NewRouter()
	addr := fmt.Sprintf(":%d", 8080)
	fmt.Printf("Server listening on %s...\n", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
