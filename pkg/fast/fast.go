package fast

import (
	"fmt"
	"math"
	"password_gen/m/v2/pkg/common"
	"runtime"
	"sync"
)

// Find not guaranted the best password, but one of the best with O(N^2) complexity
func Find(dictPath string) {
	// Define a map for storing best next word to word
	bestNextWords := make(map[int]int)
	bestWordsMutex := sync.RWMutex{}

	// Define a map of the keyboard keys and their positions
	keyboard := map[byte][2]int{
		'q': {0, 0}, 'w': {0, 1}, 'e': {0, 2}, 'r': {0, 3}, 't': {0, 4}, 'y': {0, 5}, 'u': {0, 6}, 'i': {0, 7}, 'o': {0, 8}, 'p': {0, 9},
		'a': {1, 0}, 's': {1, 1}, 'd': {1, 2}, 'f': {1, 3}, 'g': {1, 4}, 'h': {1, 5}, 'j': {1, 6}, 'k': {1, 7}, 'l': {1, 8},
		'z': {2, 0}, 'x': {2, 1}, 'c': {2, 2}, 'v': {2, 3}, 'b': {2, 4}, 'n': {2, 5}, 'm': {2, 6},
	}
	keyboardMutex := sync.RWMutex{}

	// Set the number of worker goroutines
	workerCount := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	wg.Add(workerCount)

	//read words
	words := common.ReadWordsFromFile(dictPath)
	weights := make([]int, len(words))
	for i := 0; i < len(words); i++ {
		weights[i] = -1
	}

	// Create a channel to communicate the results from the worker goroutines
	resultCh := make(chan [5]int)

	// Start the worker goroutines
	// Start building chains of optimal words by weight
	step := (len(words)) / workerCount
	for worker := 0; worker < workerCount; worker++ {
		start := worker * step
		end := (worker + 1) * step
		if end > len(words) {
			end = len(words)
		}
		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				minDist := math.MaxInt32
				minIdx := -1
				word1 := words[i]
				if weights[i] < 0 {
					weights[i] = common.CalculateWeight(word1, &keyboard, &keyboardMutex)
				}
				for j, word2 := range words {
					if i != j {
						if weights[j] < 0 {
							weights[j] = common.CalculateWeight(word2, &keyboard, &keyboardMutex)
						}
						dist := weights[i] + weights[j] + common.CalculateWeight(word1[len(word1)-1:]+word2[:1], &keyboard, &keyboardMutex)
						if dist < minDist {
							minDist = dist
							minIdx = j
						}
					}
				}
				bestWordsMutex.Lock()
				bestNextWords[i] = minIdx
				bestWordsMutex.Unlock()
			}
		}(start, end)
	}

	wg.Wait()

	// Start analyzing chains of optimal words
	minDist := math.MaxInt32
	var minPath [4]int

	wg.Add(workerCount)
	for worker := 0; worker < workerCount; worker++ {
		start := worker * step
		end := (worker + 1) * step
		if end > len(words) {
			end = len(words)
		}
		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				currentPath := make([]int, 4)
				currentDist := weights[i]
				chainLength := 0
				current := i
				for j := 0; j < 4; j++ {
					currentPath[j] = current
					chainLength += len(words[current])
					if j > 0 {
						currentDist += weights[current]
						currentDist += common.CalculateWeight(words[currentPath[j-1]][len(words[currentPath[j-1]])-1:]+words[current][:1], &keyboard, &keyboardMutex)
					}
					bestWordsMutex.RLock()
					current = bestNextWords[current]
					bestWordsMutex.RUnlock()
				}

				if currentDist < minDist && chainLength >= 20 && chainLength <= 24 {
					resultCh <- [5]int{currentDist, currentPath[0], currentPath[1], currentPath[2], currentPath[3]}
				}
			}
		}(start, end)
	}

	// Start a goroutine to close the result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Process the results
	for result := range resultCh {
		if result[0] < minDist {
			minDist = result[0]
			minPath = [4]int{result[1], result[2], result[3], result[4]}
		}
	}

	fmt.Printf("The four words are: %v, %v, %v, and %v\n", words[minPath[0]], words[minPath[1]], words[minPath[2]], words[minPath[3]])
}
