// Worker implementation that handles prime number calculation tasks. Each worker
// connects to the coordinator, receives work assignments (number ranges), processes
// them using the specified algorithm, and returns results. Workers can operate
// independently across multiple machines.

package node

import (
	"bytes"
	"distributed-prime-number-generator/src/algorithms"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Worker struct {
	ID            string
	ServerURL     string
	Client        *http.Client
}

func NewWorker(serverURL string) *Worker {
	return &Worker{
		ID:        "",  // Will be assigned by the server upon registration
		ServerURL: serverURL,
		Client:    &http.Client{Timeout: 10 * time.Second},
	}
}

// Register registers this worker with the server
func (w *Worker) Register() error {
	resp, err := w.Client.Post(w.ServerURL+"/api/workers", "application/json", nil)
	if err != nil {
		return fmt.Errorf("registration failed: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("registration failed with status: %d", resp.StatusCode)
	}
	
	var result map[string]string
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}
	
	w.ID = result["workerId"]
	fmt.Printf("Worker registered with ID: %s\n", w.ID)
	return nil
}

// GetNextChunk requests the next available chunk from the server
func (w *Worker) GetNextChunk() (*WorkChunk, error) {
	url := fmt.Sprintf("%s/api/workers/%s/chunks", w.ServerURL, w.ID)
	resp, err := w.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get chunk: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get chunk failed with status: %d", resp.StatusCode)
	}
	
	var chunk WorkChunk
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	
	if err := json.Unmarshal(body, &chunk); err != nil {
		return nil, fmt.Errorf("failed to parse chunk: %v", err)
	}
	
	return &chunk, nil
}

// SubmitResult sends the calculation result back to the server
func (w *Worker) SubmitResult(result ChunkResult) error {
	url := fmt.Sprintf("%s/api/workers/%s/results", w.ServerURL, w.ID)
	
	jsonData, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %v", err)
	}
	
	// Post the result
	resp, err := w.Client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to submit result: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("submit result failed with status: %d", resp.StatusCode)
	}
	
	return nil
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
func (w *Worker) Run() error {
	fmt.Printf("Worker starting, connecting to %s\n", w.ServerURL)
	
	if err := w.Register(); err != nil {
		return err
	}
	
	for {
		chunk, err := w.GetNextChunk()
		if err != nil {
			return fmt.Errorf("error getting chunk: %v", err)
		}
		
		// No more chunks available
		if chunk == nil {
			fmt.Printf("Worker %s: no more chunks available\n", w.ID)
			break
		}
		
		result, err := w.ProcessChunk(chunk)
		if err != nil {
			return fmt.Errorf("error processing chunk: %v", err)
		}
		
		// Submit the result
		if err := w.SubmitResult(*result); err != nil {
			return fmt.Errorf("error submitting result: %v", err)
		}
	}
	
	fmt.Printf("Worker %s finished all assigned work\n", w.ID)
	return nil
}