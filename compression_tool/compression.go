package compression_tool

import (
	"bufio"
	"errors"
	"os"
)

var charFrequencyTable = make(map[string]int)

func BuildFrequencyTable(fileName string) (map[string]int, error) {
	file, err := os.Open(fileName)
	defer file.Close()

	if err != nil {
		return nil, errors.New("couldn't open file: " + fileName)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for i := 0; i < len(line); i++ {
			char := string(line[i])
			_, isPresent := charFrequencyTable[char]
			if isPresent {
				charFrequencyTable[char] = charFrequencyTable[char] + 1
			} else {
				charFrequencyTable[char] = 1
			}
		}
	}
	return charFrequencyTable, nil
}

func BuildHuffmanTree(fileName string) (FrequencyTree, error) {
	charFrequencyTable, err := BuildFrequencyTable(fileName)
	if err != nil {
		return FrequencyTree{}, err
	}

	ft := FrequencyTrees{}
	for ch, freq := range charFrequencyTable {
		ft = append(ft, FrequencyTree{letter: ch, frequency: freq})
	}
	return ft.BuildHuffmanTree()[0], nil
}
