// Simple REST API that allows clients to request prime number calculations.
// Provides endpoints for submitting calculation requests, checking status,
// and retrieving results. Communicates with the coordinator to initiate distributed
// processing and monitor progress.

package api

import (
	"distributed-prime-number-generator/src/node"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	Coordinator *node.Coordinator
	Port        int
}

func NewServer(coordinator *node.Coordinator, port int) *Server {
	return &Server{
		Coordinator: coordinator,
		Port:        port,
	}
}

func (s *Server) Start() error {
	// Register endpoints
	http.HandleFunc("/api/jobs", s.handleJobs)
	http.HandleFunc("/api/jobs/", s.handleJobById)
	http.HandleFunc("/api/workers", s.handleWorkers)
	http.HandleFunc("/api/workers/", s.handleWorkerById)
	
	// Start the server
	addr := fmt.Sprintf(":%d", s.Port)
	fmt.Printf("API server listening on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}

type CreateJobRequest struct {
	Start     int `json:"start"`
	End       int `json:"end"`
	Rounds	  int `json:"rounds"`
	ChunkSize int `json:"chunkSize"`
}

type JobResponse struct {
	JobID string `json:"jobId"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (s *Server) handleJobs(w http.ResponseWriter, r *http.Request) {
	// Only support POST method for now
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req CreateJobRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	
	if req.Start < 2 {
		sendErrorResponse(w, "Start must be at least 2", http.StatusBadRequest)
		return
	}
	
	if req.End <= req.Start {
		sendErrorResponse(w, "End must be greater than start", http.StatusBadRequest)
		return
	}
	
	if req.ChunkSize <= 0 {
		// Default chunk size
		req.ChunkSize = 10000
	}
	
	jobID, err := s.Coordinator.CreateJob(req.Start, req.End, req.Rounds, req.ChunkSize)
	if err != nil {
		sendErrorResponse(w, fmt.Sprintf("Failed to create job: %v", err), http.StatusInternalServerError)
		return
	}
	
	response := JobResponse{JobID: jobID}
	sendJSONResponse(w, response, http.StatusCreated)
}

func (s *Server) handleJobById(w http.ResponseWriter, r *http.Request) {
    // Only support GET method for now
    if r.Method != http.MethodGet {
        sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Extract job ID from URL
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 3 {
        sendErrorResponse(w, "Invalid job ID", http.StatusBadRequest)
        return
    }
    
    jobID := parts[len(parts)-1]
    
    // TODO in future: filter results by job ID
    fmt.Printf("Getting results for job: %s\n", jobID)
    
    results, err := s.Coordinator.GetJobResults(jobID)
    if err != nil {
        sendErrorResponse(w, fmt.Sprintf("Error: %v", err), http.StatusNotFound)
        return
    }
    
    sendJSONResponse(w, results, http.StatusOK)
}

func (s *Server) handleWorkers(w http.ResponseWriter, r *http.Request) {
	// Only support POST method for worker registration
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	workerID := fmt.Sprintf("worker-%d", time.Now().UnixNano())
	
	s.Coordinator.RegisterWorker(workerID)

	response := map[string]string{"workerId": workerID}
	sendJSONResponse(w, response, http.StatusCreated)
}

func (s *Server) handleWorkerById(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 3 {
		sendErrorResponse(w, "Invalid worker ID", http.StatusBadRequest)
		return
	}
	
	var workerID string
    for i, part := range parts {
        if part == "workers" && i+1 < len(parts) {
            workerID = parts[i+1]
            break
        }
    }

	_, exists := s.Coordinator.Workers[workerID]
    if !exists {
        sendErrorResponse(w, fmt.Sprintf("Worker not found: %s", workerID), http.StatusBadRequest)
        return
    }

	if strings.HasSuffix(r.URL.Path, "/chunks") {
		s.handleGetNextChunk(w, r, workerID)
	} else if strings.Contains(r.URL.Path, "/results") {
		s.handleSubmitResults(w, r, workerID)
	} else {
		sendErrorResponse(w, "Method not allowed or invalid endpoint", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleGetNextChunk(w http.ResponseWriter, r *http.Request, workerID string) {

	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	chunk, err := s.Coordinator.GetNextChunk(workerID)
	if err != nil {
		sendErrorResponse(w, fmt.Sprintf("Error: %v", err), http.StatusBadRequest)
		return
	}
	
	if chunk == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	
	sendJSONResponse(w, chunk, http.StatusOK)
}

func (s *Server) handleSubmitResults(w http.ResponseWriter, r *http.Request, workerID string) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var result node.ChunkResult
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&result); err != nil {
		sendErrorResponse(w, "Invalid result format", http.StatusBadRequest)
		return
	}
	
	err := s.Coordinator.SubmitResult(result)
	if err != nil {
		sendErrorResponse(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}

// Helper to send JSON responses
func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Printf("Error encoding JSON response: %v\n", err)
	}
}

// Helper to send error responses
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := ErrorResponse{Error: message}
	sendJSONResponse(w, response, statusCode)
}