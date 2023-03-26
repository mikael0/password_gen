package fast

import (
	"container/heap"
	"fmt"
	"math"
	"password_gen/m/v2/pkg/common"
	"runtime"
	"sync"

	"golang.org/x/exp/slices"
)

type Result struct {
	minDist int
	path    []int
}

// Find not guaranted the best password, but one of the best with O(N) complexity
func Find(dictPath string, startLength int, endLength int, numberOfWords int) {
	// Define a map of the keyboard keys and their positions
	keyboard := map[byte][2]int{
		'q': {0, 0}, 'w': {0, 1}, 'e': {0, 2}, 'r': {0, 3}, 't': {0, 4}, 'y': {0, 5}, 'u': {0, 6}, 'i': {0, 7}, 'o': {0, 8}, 'p': {0, 9},
		'a': {1, 0}, 's': {1, 1}, 'd': {1, 2}, 'f': {1, 3}, 'g': {1, 4}, 'h': {1, 5}, 'j': {1, 6}, 'k': {1, 7}, 'l': {1, 8},
		'z': {2, 0}, 'x': {2, 1}, 'c': {2, 2}, 'v': {2, 3}, 'b': {2, 4}, 'n': {2, 5}, 'm': {2, 6},
	}

	// Set the number of worker goroutines
	workerCount := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	wg.Add(workerCount)

	//read words
	words := common.ReadWordsFromFile(dictPath)
	weightsMutex := sync.RWMutex{}
	weights := make([]int, len(words))
	for i := 0; i < len(words); i++ {
		weights[i] = -1
	}

	weightsChannel := make(chan WordWithWeight)

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
				word := words[i]
				weight := common.CalculateWeight(word, keyboard)
				wordWithWeight := WordWithWeight{word, i, weight}

				weightsMutex.Lock()
				weights[i] = weight
				weightsMutex.Unlock()

				weightsChannel <- wordWithWeight
			}
		}(start, end)
	}

	go func() {
		wg.Wait()
		close(weightsChannel)
	}()

	wordMap := make(map[byte]*MinHeap)
	for wordWithWeight := range weightsChannel {
		firstLetter := wordWithWeight.word[0]
		if _, exists := wordMap[firstLetter]; !exists {
			minHeap := &MinHeap{}
			heap.Init(minHeap)
			heap.Push(minHeap, wordWithWeight)
			wordMap[firstLetter] = minHeap
		} else {
			heap.Push(wordMap[firstLetter], wordWithWeight)
		}
	}

	resultCh := make(chan Result)
	wordMapMutex := sync.RWMutex{}
	wg.Add(workerCount)
	for worker := 0; worker < workerCount; worker++ {
		start := worker * step
		end := (worker + 1) * step
		if end > len(words) {
			end = len(words)
		}
		go func(start, end, worker int) {
			defer wg.Done()
			minDist := math.MaxInt32
			var minPath []int
			for i := start; i < end; i++ {
				currentPath := make([]int, 4)
				chainLength := 0
				current := i

				weightsMutex.Lock()
				currentDist := weights[i]
				weightsMutex.Unlock()

				for j := 0; j < numberOfWords; j++ {
					currentPath[j] = current
					currentWord := words[current]
					currentWordLength := len(currentWord)
					lastLetter := currentWord[currentWordLength-1]
					chainLength += currentWordLength

					currentHeap := new(MinHeap)

					wordMapMutex.Lock()
					*currentHeap = *wordMap[lastLetter]
					wordMapMutex.Unlock()

					var currentItem WordWithWeight
					for ok := true; ok; ok = slices.Contains(currentPath, currentItem.idx) {
						if len(*currentHeap) <= 0 {
							fmt.Printf("failed chain - no candidates left\n")
							return
						}
						currentItem = heap.Pop(currentHeap).(WordWithWeight)
					}
					currentDist += currentItem.weight

					current = currentItem.idx
				}
				if currentDist < minDist && chainLength >= startLength && chainLength <= endLength {
					// fmt.Printf("[%v]Candidate words are: %v, %v, %v, and %v with weight %v (previous %v)\n", worker, words[currentPath[0]], words[currentPath[1]], words[currentPath[2]], words[currentPath[3]], currentDist, minDist)
					minDist = currentDist
					minPath = currentPath
				}
			}
			resultCh <- Result{minDist, minPath}
		}(start, end, worker)
	}

	// Start a goroutine to close the result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	minDist := math.MaxInt32
	var minPath []int
	// Process the results
	for result := range resultCh {
		if result.minDist < minDist {
			minDist = result.minDist
			minPath = result.path
		}
	}

	fmt.Printf("The four words are: %v, %v, %v, and %v with weight %v\n", words[minPath[0]], words[minPath[1]], words[minPath[2]], words[minPath[3]], minDist)
}
