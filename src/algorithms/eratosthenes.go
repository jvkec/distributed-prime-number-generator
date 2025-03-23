/*
	This file implements the Sieve of Eratosthenes algorithm for finding prime numbers.
	The implementation is optimized for distributed computing, allowing work to be
	divided into manageable chunks that can be processed independently by worker nodes.
	This implementation works efficiently for ranges up to 10^8.
*/

package algorithms

import (
	"fmt"
	// "math"
	// "errors"
)

// ===== Min & max prime numbers for this algo =====
const MAX_PRIME = 100000000
const MIN_PRIME = 2

// ===== Main function to be called externally =====
func FindPrimesWithEratosthenes(start, end int) ([]int, error) {

	return []int{}, nil
}


// ===== Helper functions =====
func SieveOfEratosthenes(limit int) []bool {

	return make([]bool, limit+1)
}

func ValidateRange(start, end int) error {
	if start < MIN_PRIME {
		return fmt.Errorf("Start value (%d) is less than minimum prime number (%d)", start, MIN_PRIME)
	} else if end > MAX_PRIME {
		return fmt.Errorf("End value (%d) is greater than max. prime number (%d)", start, MIN_PRIME)
	} else if start > end {
		return fmt.Errorf("Start (%d) must be less than or equal to end (%d)", start, end)
	}

	return nil // nil -> no error
}

