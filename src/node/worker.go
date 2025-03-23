// Worker implementation that handles prime number calculation tasks. Each worker
// connects to the coordinator, receives work assignments (number ranges), processes
// them using the specified algorithm, and returns results. Workers can operate
// independently across multiple machines.

package node

import (
	"distributed-prime-number-generator/src/algorithms"
	"fmt"
	"time"
)

// Worker represents a worker node in the distributed system
type Worker struct {
	ID            string
	CoordinatorID string
}

// NewWorker creates a new worker instance
func NewWorker(id string, coordinatorID string) *Worker {
	return &Worker{
		ID:            id,
		CoordinatorID: coordinatorID,
	}
}

// ProcessChunk handles the calculation of primes in a given chunk
func (w *Worker) ProcessChunk(chunk *WorkChunk) (*ChunkResult, error) {
	startTime := time.Now()
	
	var primes []int
	var err error

	if chunk.Algorithm == SOE {
		fmt.Printf("Worker %s processing chunk %s with Sieve of Eratosthenes\n", w.ID, chunk.ID)
		primes, err = algorithms.FindPrimesWithEratosthenes(chunk.Start, chunk.End)
	} else {
		fmt.Printf("Worker %s processing chunk %s with Miller-Rabin\n", w.ID, chunk.ID)
		primes, err = algorithms.FindPrimesWithMillerRabin(chunk.Start, chunk.End, chunk.Rounds)
	}
	
	if err != nil {
		return nil, fmt.Errorf("error processing chunk: %v", err)
	}
	
	runtime := time.Since(startTime)
	
	result := &ChunkResult{
		ChunkID: chunk.ID,
		Primes:  primes,
		Runtime: runtime,
	}
	
	fmt.Printf("Worker %s finished chunk %s (found %d primes in %v)\n", 
		w.ID, chunk.ID, len(primes), runtime)
	
	return result, nil
}

// Run starts the worker's processing loop
func (w *Worker) Run(coordinator *Coordinator) error {
	fmt.Printf("Worker %s starting\n", w.ID)
	
	coordinator.RegisterWorker(w.ID)

	for {
		chunk, err := coordinator.GetNextChunk(w.ID)
		if err != nil {
			return fmt.Errorf("error getting chunk: %v", err)
		}
		
		if chunk == nil {
			fmt.Printf("Worker %s: no more chunks available\n", w.ID)
			break
		}
		
		result, err := w.ProcessChunk(chunk)
		if err != nil {
			return fmt.Errorf("error processing chunk: %v", err)
		}
		
		err = coordinator.SubmitResult(*result)
		if err != nil {
			return fmt.Errorf("error submitting result: %v", err)
		}
	}
	
	fmt.Printf("Worker %s finished all assigned work\n", w.ID)
	return nil
}