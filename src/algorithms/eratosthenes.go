/*
	This file implements the Sieve of Eratosthenes algorithm for finding prime numbers.
	The implementation is optimized for distributed computing, allowing work to be
	divided into manageable chunks that can be processed independently by worker nodes.
	This implementation works efficiently for ranges up to 10^8.
*/

package algorithms

import (
	"fmt"
)

// ===== Min & max prime numbers for this algo =====
const MAX_PRIME_SOE = 100000000
const MIN_PRIME_SOE = 2

// ===== Main function to be called externally =====
func FindPrimesWithEratosthenes(start, end int) ([]int, error) {

	if err := validateRangeSOE(start, end); err != nil {
		return nil, err
	}

	isPrime := sieveOfEratosthenes(end)
	
	var primes []int
	for i := start; i <= end; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	
	return primes, nil
}

// ===== Helper functions =====
func sieveOfEratosthenes(end int) []bool {

	isPrime := make([]bool, end+1)
	for i := 2; i <= end; i++ {
		isPrime[i] = true
	}

	for p := 2; p*p < end; p++ {
		if isPrime[p] {
			for i := p * p; i <= end; i += p {
				isPrime[i] = false
			}
		}
	}

	return isPrime
}

func validateRangeSOE(start, end int) error {
	if start < MIN_PRIME_SOE {
		return fmt.Errorf("Start value (%d) is less than minimum prime number for SOE (%d)", start, MIN_PRIME_SOE)
	} else if end > MAX_PRIME_SOE {
		return fmt.Errorf("End value (%d) is greater than maximum prime number for SOE (%d)", start, MAX_PRIME_SOE)
	} else if start > end {
		return fmt.Errorf("Start (%d) must be less than or equal to end (%d)", start, end)
	}

	return nil
}

