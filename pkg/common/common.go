package common

import (
	"bufio"
	"fmt"
	"os"
)

var keyboard = map[byte][2]int{
	'q': {0, 0}, 'w': {0, 1}, 'e': {0, 2}, 'r': {0, 3}, 't': {0, 4}, 'y': {0, 5}, 'u': {0, 6}, 'i': {0, 7}, 'o': {0, 8}, 'p': {0, 9},
	'a': {1, 0}, 's': {1, 1}, 'd': {1, 2}, 'f': {1, 3}, 'g': {1, 4}, 'h': {1, 5}, 'j': {1, 6}, 'k': {1, 7}, 'l': {1, 8},
	'z': {2, 0}, 'x': {2, 1}, 'c': {2, 2}, 'v': {2, 3}, 'b': {2, 4}, 'n': {2, 5}, 'm': {2, 6},
}

const Alphabet = "abcdefghijklmnopqrstuvwxyz"

// Define a function that reads the words from a file
func ReadWordsFromFile(filename string) []string {
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
func CalculateWeight(word string) int {
	weight := 0
	for i := 0; i < len(word)-1; i++ {
		dist := DistBytes(word[i], word[i+1])
		weight += dist
	}
	return weight
}

func DistBytes(a, b byte) int {
	pos1 := keyboard[a]
	pos2 := keyboard[b]
	return absDiffInt(pos1[0], pos2[0]) + absDiffInt(pos1[1], pos2[1])
}

func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}
