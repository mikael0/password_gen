package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"
	"sync"
)

// Define a function that reads the words from a file
func readWordsFromFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	words := make([]string, 0)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
	return words
}

// Define a function that calculates the weight of a word
func calculateWeight(word string, keyboard *map[byte][2]int, keyboardMutex *sync.RWMutex) int {
	weight := 0
	for i := 0; i < len(word)-1; i++ {
		(*keyboardMutex).RLock()
		pos1 := (*keyboard)[word[i]]
		pos2 := (*keyboard)[word[i+1]]
		(*keyboardMutex).RUnlock()
		dist := int(math.Abs(float64(pos1[0]-pos2[0])) + math.Abs(float64(pos1[1]-pos2[1])))
		weight += dist
	}
	return weight
}

//Naive approach comparing each combinations in O(n^4)
func findOptimal(dictPath string) {
	minDist := math.MaxInt32
	var minPath [4]int

	// Define a map of the keyboard keys and their positions
	keyboard := map[byte][2]int{
		'q': {0, 0}, 'w': {0, 1}, 'e': {0, 2}, 'r': {0, 3}, 't': {0, 4}, 'y': {0, 5}, 'u': {0, 6}, 'i': {0, 7}, 'o': {0, 8}, 'p': {0, 9},
		'a': {1, 0}, 's': {1, 1}, 'd': {1, 2}, 'f': {1, 3}, 'g': {1, 4}, 'h': {1, 5}, 'j': {1, 6}, 'k': {1, 7}, 'l': {1, 8},
		'z': {2, 0}, 'x': {2, 1}, 'c': {2, 2}, 'v': {2, 3}, 'b': {2, 4}, 'n': {2, 5}, 'm': {2, 6},
	}
	keyboardMutex := sync.RWMutex{}

	words := readWordsFromFile(dictPath)

	// Set the number of worker goroutines
	workerCount := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	wg.Add(workerCount)

	// Create a channel to communicate the results from the worker goroutines
	resultCh := make(chan [5]int)

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
							if len(combined) < 20 || len(combined) > 24 {
								continue
							}
							dist := calculateWeight(combined, &keyboard, &keyboardMutex)
							if dist < localMinDist && i != j && i != k && i != l && j != k && j != l && k != l {
								localMinDist = dist
								localMinPath = [4]int{i, j, k, l}
							}
						}
					}
				}
			}
			resultCh <- [5]int{localMinDist, localMinPath[0], localMinPath[1], localMinPath[2], localMinPath[3]}
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

//Naive approach comparing each combinations in O(n^4)
func findMaybeOptimalFast(dictPath string) {
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
	words := readWordsFromFile(dictPath)
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
		fmt.Printf("starting %v to %v\n", start, end)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				minDist := math.MaxInt32
				minIdx := -1
				word1 := words[i]
				if weights[i] < 0 {
					weights[i] = calculateWeight(word1, &keyboard, &keyboardMutex)
				}
				for j, word2 := range words {
					if i != j {
						if weights[j] < 0 {
							weights[j] = calculateWeight(word2, &keyboard, &keyboardMutex)
						}
						dist := weights[i] + weights[j] + calculateWeight(word1[len(word1)-1:]+word2[:1], &keyboard, &keyboardMutex)
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
		fmt.Printf("starting %v to %v\n", start, end)
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
						currentDist += calculateWeight(words[currentPath[j-1]][len(words[currentPath[j-1]])-1:]+words[current][:1], &keyboard, &keyboardMutex)
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

func main() {
	dictPath := flag.String("dict", "", "Path to the dictionary file")
	mode := flag.String("mode", "fast", "Mode of operation (fast or optimal)")

	flag.Parse()

	if *dictPath == "" {
		fmt.Println("Error: --dict is required")
		flag.Usage()
		os.Exit(1)
	}

	if *mode != "fast" && *mode != "optimal" {
		fmt.Println("Error: --mode must be 'fast' or 'optimal'")
		flag.Usage()
		os.Exit(1)
	}

	switch *mode {
		case "fast": 
			findMaybeOptimalFast(*dictPath)
		case "optimal": 
			findOptimal(*dictPath)
	}
}
