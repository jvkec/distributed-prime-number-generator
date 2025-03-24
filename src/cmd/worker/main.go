// Entry point for worker nodes. Initializes the worker process, connects to the
// coordinator, and begins accepting and processing work assignments. Multiple instances
// can be run across different machines to form the distributed computing cluster.

package main

import (
	"distributed-prime-number-generator/src/node"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	serverURL := flag.String("server", "http://localhost:8080", "URL of the coordinator server")
	flag.Parse()

	fmt.Println("=====================================================")
	fmt.Println("  Distributed Prime Number Generator - Worker")
	fmt.Println("=====================================================")
	
	worker := node.NewWorker(*serverURL)
	
	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// Start worker in a goroutine
	errChan := make(chan error, 1)
	go func() {
		err := worker.Run()
		if err != nil {
			errChan <- err
		}
	}()
	
	fmt.Printf("Worker started, connecting to %s\n", *serverURL)
	fmt.Println("Press Ctrl+C to shutdown.")
	
	select {
	case err := <-errChan:
		if err != nil {
			log.Fatalf("Worker error: %v", err)
		}
		fmt.Println("Worker completed all work successfully.")
	case <-sigChan:
		fmt.Println("\nShutting down worker...")
	}
}