# Distributed Prime Number Generator

A scalable system that distributes prime number calculations across multiple worker nodes, using different algorithms optimized for various number ranges.

## Why I built this

One of my first projects in Go's and I wanted to do something related to my interests in distributed systems. I thought Go's concurrency was a good use case for generating a heck ton of prime numbers!

## Overview

This distributed system allows you to find prime numbers within specified ranges by distributing the workload across multiple worker nodes. It automatically selects the most appropriate algorithm based on the range:

- **Sieve of Eratosthenes**: Efficient for smaller ranges (up to 10^8)
- **Miller-Rabin Primality Test**: Used for larger ranges (10^8 to 10^12)

## Features

- **Work Distribution**: Divides calculation ranges into manageable chunks
- **Multiple Algorithms**: Selects the optimal algorithm based on number size
- **REST API**: Submit jobs and retrieve results via HTTP endpoints
- **Stateless Workers**: Add or remove workers dynamically as needed

## Getting Started

### Prerequisites

- Go 1.16 or higher
- Network connectivity between server and worker nodes

### Running the Server

```bash
go run cmd/server/main.go -port 8080
```

### Running Worker Nodes

You can run multiple workers on different machines:

```bash
go run cmd/worker/main.go -server http://server-ip:8080
```

### Creating a Job

Use the API to create a prime calculation job:

```bash
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -d '{"start": 2, "end": 1000000, "rounds": 10, "chunkSize": 10000}'
```

Parameters:
- `start`: Beginning of the range to search for primes
- `end`: End of the range
- `rounds`: Number of rounds for Miller-Rabin test (5-40 recommended)
- `chunkSize`: Size of each work unit (affects distribution granularity)

### Retrieving Results

```bash
curl http://localhost:8080/api/jobs/job-id
```

Replace `job-id` with the ID returned when creating the job.

## Miller-Rabin Round Recommendations

For the Miller-Rabin primality test, the number of rounds affects accuracy:

- **5-7 rounds**: Good balance for most applications
- **10-15 rounds**: Higher confidence for cryptographic applications
- **20-40 rounds**: Very high confidence for critical applications

## Usage Examples

### Finding primes in a small range (uses Sieve of Eratosthenes)

```bash
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -d '{"start": 2, "end": 10000, "rounds": 5, "chunkSize": 1000}'
```

### Finding primes in a large range (uses Miller-Rabin)

```bash
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -d '{"start": 100000000, "end": 100001000, "rounds": 10, "chunkSize": 100}'
```
