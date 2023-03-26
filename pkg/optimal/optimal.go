package optimal

import (
	"fmt"
	"math"
	"password_gen/m/v2/pkg/common"
	"runtime"
	"sync"
)

type Result struct {
	minDist int
	path    []int
}

// Naive approach comparing each combinations in O(n^4)
func Find(dictPath string, startLength int, endLength int, numberOfWords int) {
	minDist := math.MaxInt32
	var minPath []int

	// Define a map of the keyboard keys and their positions
	keyboard := map[byte][2]int{
		'q': {0, 0}, 'w': {0, 1}, 'e': {0, 2}, 'r': {0, 3}, 't': {0, 4}, 'y': {0, 5}, 'u': {0, 6}, 'i': {0, 7}, 'o': {0, 8}, 'p': {0, 9},
		'a': {1, 0}, 's': {1, 1}, 'd': {1, 2}, 'f': {1, 3}, 'g': {1, 4}, 'h': {1, 5}, 'j': {1, 6}, 'k': {1, 7}, 'l': {1, 8},
		'z': {2, 0}, 'x': {2, 1}, 'c': {2, 2}, 'v': {2, 3}, 'b': {2, 4}, 'n': {2, 5}, 'm': {2, 6},
	}

	words := common.ReadWordsFromFile(dictPath)

	// Set the number of worker goroutines
	workerCount := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	wg.Add(workerCount)

	// Create a channel to communicate the results from the worker goroutines
	resultCh := make(chan Result)

	// Start the worker goroutines
	step := (len(words)) / workerCount
	for worker := 0; worker < workerCount; worker++ {
		start := worker * step
		end := (worker + 1) * step
		if end > len(words) {
			end = len(words)
		}
		fmt.Printf("starting %v to %v\n", start, end)
		go func(start, end int) {
			defer wg.Done()
			var localMinDist int = math.MaxInt32
			var localMinPath [4]int
			for i := start; i < end; i++ {
				for j := 0; j < len(words); j++ {
					for k := 0; k < len(words); k++ {
						for l := 0; l < len(words); l++ {
							combined := words[i] + words[j] + words[k] + words[l]
							if len(combined) < startLength || len(combined) > endLength {
								continue
							}
							dist := common.CalculateWeight(combined, keyboard)
							if dist < localMinDist && i != j && i != k && i != l && j != k && j != l && k != l {
								localMinDist = dist
								localMinPath = [4]int{i, j, k, l}
							}
						}
					}
				}
			}
			resultCh <- Result{localMinDist, localMinPath[:]}
		}(start, end)
	}

	// Start a goroutine to close the result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Process the results
	for result := range resultCh {
		if result.minDist < minDist {
			minDist = result.minDist
			minPath = result.path
		}
	}

	fmt.Printf("The four words are: %v, %v, %v, and %v\n", words[minPath[0]], words[minPath[1]], words[minPath[2]], words[minPath[3]])
}
