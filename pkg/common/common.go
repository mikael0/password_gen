package common

import (
	"bufio"
	"fmt"
	"os"
)

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
func CalculateWeight(word string, keyboard map[byte][2]int) int {
	weight := 0
	for i := 0; i < len(word)-1; i++ {
		pos1 := keyboard[word[i]]
		pos2 := keyboard[word[i+1]]
		dist := absDiffInt(pos1[0], pos2[0]) + absDiffInt(pos1[1], pos2[1])
		weight += dist
	}
	return weight
}

func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}
