package compression_tool

import (
	"bufio"
	"errors"
	"os"
)

type PrefixTable = map[string]string

func BuildFrequencyTable(fileName string) (map[string]int, error) {
	file, err := os.Open(fileName)
	defer file.Close()

	if err != nil {
		return nil, errors.New("couldn't open file: " + fileName)
	}

	charFrequencyTable := make(map[string]int)

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

func GeneratePrefixCodeTable(fileName string) PrefixTable {
	huffmanTree, _ := BuildHuffmanTree(fileName)
	prefixTable := make(map[string]string)
	TraverseHuffmanTree(huffmanTree, "", &prefixTable)
	return prefixTable
}

func TraverseHuffmanTree(huffTree FrequencyTree, code string, prefixTable *map[string]string) {
	if huffTree.left == nil && huffTree.right == nil {
		(*prefixTable)[huffTree.letter] = code
	}

	if huffTree.right != nil {
		TraverseHuffmanTree(*huffTree.right, code+"1", prefixTable)
	}

	if huffTree.left != nil {
		TraverseHuffmanTree(*huffTree.left, code+"0", prefixTable)
	}
}
