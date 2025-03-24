// This file contains the Miller-Rabin primality test implementation. This probabilistic
// algorithm is more efficient for testing larger numbers. It provides configurable
// accuracy through multiple rounds of testing and is suitable for ranges containing
// very large numbers.

package algorithms

import (
	// "fmt"
	"math/big"
	"math/rand"
	"time"
)

func FindPrimesWithMillerRabin(start, end int, rounds int) ([]int, error) {
	
	if rounds <= 0 {
		rounds = 5
	}

	rand.Seed(time.Now().UnixNano())

	var primes []int
	for num := start; num <= end; num++ {
		if isMillerRabinPrime(num, rounds) {
			primes = append(primes, num)
		}
	}

	return primes, nil
}

func isMillerRabinPrime(n int, rounds int) bool {

	if n%2 == 0 {
		return false
	}

	nBig := big.NewInt(int64(n))
	r, d := decompose(n-1)

	nMinusOne := big.NewInt(int64(n - 1))
	one := big.NewInt(1)
	two := big.NewInt(2)

	// Primality test
	for i := 0; i < rounds; i++ {
		a := randomBigInt(2, n-2)
		aBig := big.NewInt(int64(a))

		x := new(big.Int).Exp(aBig, d, nBig)

		if x.Cmp(one) == 0 || x.Cmp(nMinusOne) == 0 {
			continue
		}

		if !checkComposite(x, r, nBig, one, two, nMinusOne) {
			return false
		}
	}

	return true
}

// decompose expresses n as 2^r * d where d is odd
func decompose(n int) (r int, d *big.Int) {
	d = big.NewInt(int64(n))
	r = 0

	// Count how many times d is divisible by 2
	two := big.NewInt(2)
	zero := big.NewInt(0)

	for new(big.Int).Mod(d, two).Cmp(zero) == 0 {
		d.Div(d, two)
		r++
	}

	return r, d
}

// checkComposite performs the r-1 iterations for the Miller-Rabin test
func checkComposite(x *big.Int, r int, nBig, one, two, nMinusOne *big.Int) bool {
	for j := 0; j < r-1; j++ {
		// x = x^2 % n
		x.Exp(x, two, nBig)

		// If x == 1, the number is composite
		if x.Cmp(one) == 0 {
			return false
		}

		// If x == n-1, the number might be prime
		if x.Cmp(nMinusOne) == 0 {
			return true
		}
	}

	return false
}

func randomBigInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}
