package optimal

import (
	"fmt"
	"math"
	"password_gen/m/v2/pkg/common"
	"sort"
)

type WordEntry struct {
	Weight    int
	word      string
	firstChar byte
}

type Engine struct {
	dictionary         []WordEntry
	optimalWeight      int
	optimalWordIndices []int
}

func (engine *Engine) find(charDistances [][]int, size, minLength, maxLength, minWordLength int) {
	if len(engine.dictionary) < size {
		return
	}

	maxWordLength := maxLength - (size-1)*minWordLength
	engine.optimalWordIndices = make([]int, size)

	visited := make([]int, size)
	for i := 0; i < size; i++ {
		visited[i] = -1
	}

	weightLengthData := make([]struct {
		Weight int
		length int
	}, size)

	for index, word := range engine.dictionary {
		curWeight := word.Weight
		curLength := len(word.word)
		if curWeight >= engine.optimalWeight {
			break
		}
		if curLength > maxWordLength {
			continue
		}

		weightLengthData[0] = struct {
			Weight int
			length int
		}{curWeight, curLength}
		wordIndices := []int{index}
		for len(wordIndices) > 0 {
			lastWord := engine.dictionary[wordIndices[len(wordIndices)-1]].word
			lastChar := lastWord[len(lastWord)-1] - 'a'
			visitedLimit := visited[len(wordIndices)]
			foundNext := false
			for i, nextWord := range engine.dictionary {
				newWeight := weightLengthData[len(wordIndices)-1].Weight + nextWord.Weight
				if newWeight >= engine.optimalWeight {
					break
				}

				if i <= visitedLimit || len(nextWord.word) > maxWordLength {
					continue
				}
				newWeight += charDistances[lastChar][nextWord.firstChar]
				if newWeight >= engine.optimalWeight {
					continue
				}
				newLength := weightLengthData[len(wordIndices)-1].length + len(nextWord.word)
				if newLength > maxLength || contains(wordIndices, i) {
					continue
				}

				visited[len(wordIndices)] = i
				if len(wordIndices) == size-1 {
					if newLength >= minLength && newLength <= maxLength {
						engine.optimalWeight = newWeight
						copy(engine.optimalWordIndices, append(wordIndices, i))
					}
					continue
				}

				weightLengthData[len(wordIndices)] = struct {
					Weight int
					length int
				}{newWeight, newLength}
				wordIndices = append(wordIndices, i)
				foundNext = true
				break
			}
			if !foundNext {
				wordIndices = wordIndices[:len(wordIndices)-1]
				for i := len(wordIndices) + 1; i < size; i++ {
					visited[i] = -1
				}
			}
		}
	}
}

func contains(arr []int, elem int) bool {
	for _, element := range arr {
		if elem == element {
			return true
		}
	}

	return false
}

func Find(path string, min, max, length int) {
	dict := common.ReadWordsFromFile(path)

	minWordLength := 0
	sortedDict := make([]WordEntry, len(dict))

	for i := range dict {
		word := dict[i]
		sortedDict[i] = WordEntry{common.CalculateWeight(word), word, word[0] - 'a'}
	}

	sort.Slice(sortedDict, func(i, j int) bool {
		return sortedDict[i].Weight < sortedDict[j].Weight
	})

	alphabetLength := len(common.Alphabet)
	charDistances := make([][]int, alphabetLength)
	for i := 0; i < alphabetLength; i++ {
		charDistances[i] = make([]int, alphabetLength)
		for j := 0; j < alphabetLength; j++ {
			charDistances[i][j] = common.DistBytes(common.Alphabet[i], common.Alphabet[j])
		}
	}

	engine := &Engine{dictionary: sortedDict, optimalWeight: math.MaxInt32}
	engine.find(charDistances, length, min, max, minWordLength)

	totalWeight := 0
	totalLength := 0
	var lastChar byte
	fmt.Printf("Words are \n")
	for _, i := range engine.optimalWordIndices {
		wordEntry := engine.dictionary[i]
		totalWeight += wordEntry.Weight
		currWord := wordEntry.word
		totalLength += len(currWord)

		fmt.Printf("%s (%v)\n", currWord, wordEntry.Weight)

		if lastChar != 0 {
			totalWeight += charDistances[lastChar-'a'][currWord[0]-'a']
		}
		lastChar = currWord[len(currWord)-1]
	}
	fmt.Printf("With weight %v and length %v \n", totalWeight, totalLength)
}
