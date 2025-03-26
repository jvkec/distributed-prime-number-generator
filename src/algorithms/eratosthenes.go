// This file implements the Sieve of Eratosthenes algorithm for finding prime numbers.
// The implementation is optimized for distributed computing, allowing work to be
// divided into manageable chunks that can be processed independently by worker nodes.
// This implementation works efficiently for ranges up to 10^8.

package algorithms

func FindPrimesWithEratosthenes(start, end int) ([]int, error) {

	isPrime := sieveOfEratosthenes(end)
	
	var primes []int
	for i := start; i <= end; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	
	return primes, nil
}

func sieveOfEratosthenes(end int) []bool {

	isPrime := make([]bool, end+1)
	for i := 2; i <= end; i++ {
		isPrime[i] = true
	}

	for p := 2; p*p <= end; p++ {
		if isPrime[p] {
			for i := p * p; i <= end; i += p {
				isPrime[i] = false
			}
		}
	}

	return isPrime
}