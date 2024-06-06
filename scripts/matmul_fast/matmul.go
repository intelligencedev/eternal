package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
)

const blockSize = 16
const numThreads = 8

// This implementation leverages Go’s concurrency model to parallelize the computation.
// Adjust the number of threads (numThreads) based on your CPU capabilities for optimal performance.
func blockedColumnParallelAtomicGemm(A, B, C []float64, N int) {
	var pos uint64
	var wg sync.WaitGroup

	for t := 0; t < numThreads; t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for colChunk := int(atomic.AddUint64(&pos, blockSize)) - blockSize; colChunk < N; colChunk = int(atomic.AddUint64(&pos, blockSize)) - blockSize {
				for row := 0; row < N; row++ {
					for tile := 0; tile < N; tile += blockSize {
						for tileRow := 0; tileRow < blockSize; tileRow++ {
							for idx := 0; idx < blockSize; idx++ {
								C[row*N+colChunk+idx] += A[row*N+tile+tileRow] * B[tile*N+tileRow*N+colChunk+idx]
							}
						}
					}
				}
			}
		}()
	}

	wg.Wait()
}

func main() {
	N := 1024 + 16 // Example size
	A := make([]float64, N*N)
	B := make([]float64, N*N)
	C := make([]float64, N*N)

	// Initialize matrices with random values
	for i := range A {
		A[i] = rand.Float64()
		B[i] = rand.Float64()
	}

	blockedColumnParallelAtomicGemm(A, B, C, N)

	fmt.Println("Matrix multiplication completed.")
}
