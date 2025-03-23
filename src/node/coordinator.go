// Coordinator service that manages the distribution of work among connected worker
// nodes. It divides the requested prime number range into chunks, assigns these
// chunks to available workers, tracks progress, and collects results. Implements
// basic load balancing and worker failure handling.

package node

import (
	"fmt"
	"sync"
	"time"
)

type AlgorithmType string

const (
    SOE AlgorithmType = "Sieve of Eratosthenes"
    MRPT  AlgorithmType = "Miller Rabin Primality Test"
    TRANSITION_THRESHOLD = 100000000
)

type WorkChunk struct {
    ID        string
    Start     int
    End       int
    Algorithm AlgorithmType
}

type ChunkResult struct {
	ChunkID string
	Primes  []int
	Runtime time.Duration
}

type WorkerInfo struct {
	ID            string
	LastHeartbeat time.Time
	ActiveChunks  []string
	CompletedJobs int
}

type Coordinator struct {
	workers       map[string]*WorkerInfo
	chunks        map[string]*WorkChunk
	results       map[string]*ChunkResult
	pendingChunks []string
	mutex         sync.Mutex
}

// "Constructor"; keep track of workers, chunks, and results in hashmaps
func NewCoordinator() *Coordinator {
	return &Coordinator{
		workers:       make(map[string]*WorkerInfo),
		chunks:        make(map[string]*WorkChunk),
		results:       make(map[string]*ChunkResult),
		pendingChunks: []string{},
	}
}

// RegisterWorker adds a new worker to the pool
func (c *Coordinator) RegisterWorker(workerID string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.workers[workerID] = &WorkerInfo{
		ID:            workerID,
		LastHeartbeat: time.Now(),
		ActiveChunks:  []string{},
		CompletedJobs: 0,
	}

	fmt.Printf("Worker registered: %s\n", workerID)
}

// Divides a range into chunks and prepares them for processing
func (c *Coordinator) CreateJob(start, end, chunkSize int) (string, error) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    jobID := fmt.Sprintf("job-%d", time.Now().UnixNano())
    
    for chunkStart := start; chunkStart <= end; chunkStart += chunkSize {
        chunkEnd := chunkStart + chunkSize - 1
        if chunkEnd > end {
            chunkEnd = end
        }
        
        algorithm := SOE
        if chunkStart >= TRANSITION_THRESHOLD || chunkEnd >= TRANSITION_THRESHOLD {
            algorithm = MRPT
        }
        
        chunkID := fmt.Sprintf("%s-chunk-%d-%d", jobID, chunkStart, chunkEnd)
        chunk := &WorkChunk{
            ID:        chunkID,
            Start:     chunkStart,
            End:       chunkEnd,
            Algorithm: algorithm,
        }
        
        c.chunks[chunkID] = chunk
        c.pendingChunks = append(c.pendingChunks, chunkID)
        
        fmt.Printf("Created chunk: %s (%d to %d) using %s\n", 
            chunkID, chunkStart, chunkEnd, algorithm)
    }
    
    return jobID, nil
}

// GetNextChunk assigns the next available chunk to a worker
func (c *Coordinator) GetNextChunk(workerID string) (*WorkChunk, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	worker, exists := c.workers[workerID]
	if !exists {
		return nil, fmt.Errorf("unknown worker (%s)", workerID)
	}
	
	worker.LastHeartbeat = time.Now()
	
	if len(c.pendingChunks) == 0 {
		return nil, nil
	}
	
	chunkID := c.pendingChunks[0]
	c.pendingChunks = c.pendingChunks[1:]
	
	worker.ActiveChunks = append(worker.ActiveChunks, chunkID)
	
	fmt.Printf("Assigned chunk %s to worker %s\n", chunkID, workerID)
	
	return c.chunks[chunkID], nil
}

// SubmitResult stores the result of a processed chunk
func (c *Coordinator) SubmitResult(result ChunkResult) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.results[result.ChunkID] = &result
	
	for workerID, worker := range c.workers {
		for i, chunkID := range worker.ActiveChunks {
			if chunkID == result.ChunkID {
				worker.ActiveChunks = append(worker.ActiveChunks[:i], worker.ActiveChunks[i+1:]...)
				worker.CompletedJobs++
				
				fmt.Printf("Worker %s completed chunk %s (found %d primes in %v)\n", 
					workerID, result.ChunkID, len(result.Primes), result.Runtime)
				break
			}
		}
	}
	
	return nil
}

// GetResults combines all results for completed chunks
func (c *Coordinator) GetResults() []int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var allPrimes []int
	for _, result := range c.results {
		allPrimes = append(allPrimes, result.Primes...)
	}
	
	return allPrimes
}