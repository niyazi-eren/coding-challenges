package compression_tool

import (
	"fmt"
	"testing"
)

func TestFileDoesNotExist(t *testing.T) {
	_, err := BuildFrequencyTable("doesNotExist")
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestFrequencyTable(t *testing.T) {
	table, _ := BuildFrequencyTable("book.txt")
	assertOccurrence(t, table, "X", 333)
	assertOccurrence(t, table, "t", 223000)
}

func TestGeneratePrefixCodeTable(t *testing.T) {
	prefixTable := GeneratePrefixCodeTable("test.txt")

	minCodeLen, maxCodeLen := 1000, 0
	for _, bits := range prefixTable {
		codeLen := len(bits)
		if minCodeLen > codeLen {
			minCodeLen = codeLen
		}
		if maxCodeLen < codeLen {
			maxCodeLen = codeLen
		}
	}
	codeE, _ := prefixTable["E"]
	codeK, _ := prefixTable["K"]
	codeZ, _ := prefixTable["Z"]
	assertTrue(t, len(codeE) == minCodeLen, fmt.Sprintf("size of prefix code %v != %v", codeE, minCodeLen))
	assertTrue(t, len(codeZ) == maxCodeLen, fmt.Sprintf("size of prefix code %v != %v", codeZ, maxCodeLen))
	assertTrue(t, len(codeK) == maxCodeLen, fmt.Sprintf("size of prefix code %v != %v", codeK, maxCodeLen))
}

func assertOccurrence(t *testing.T, table map[string]int, char string, expected int) {
	result, _ := table[char]
	if result != expected {
		t.Errorf("Error: expected %v, got: %v", expected, result)
	}
}

func assertTrue(t *testing.T, result bool, errorMsg string) {
	if !result {
		t.Errorf(errorMsg)
	}
}
