// Entry point for the coordinator and API server. Initializes the coordinator service,
// sets up the REST API endpoints, and starts listening for worker connections and
// client requests. Can be run on a dedicated server to manage the distributed system.

package main

import (
	"distributed-prime-number-generator/src/api"
	"distributed-prime-number-generator/src/node"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Parse command-line flags
	port := flag.Int("port", 8080, "API server port")
	flag.Parse()

	fmt.Println("=====================================================")
	fmt.Println("  Distributed Prime Number Generator - Server")
	fmt.Println("=====================================================")
	
	coordinator := node.NewCoordinator()
	fmt.Println("Coordinator initialized")
	
	// Create and start the API server
	server := api.NewServer(coordinator, *port)
	fmt.Printf("Starting API server on port %d...\n", *port)
	
	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	fmt.Println("Server is running. Press Ctrl+C to shutdown.")
	
	// Wait for termination signal
	<-sigChan
	fmt.Println("\nShutting down server...")
}