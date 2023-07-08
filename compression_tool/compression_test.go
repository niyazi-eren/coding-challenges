package compression_tool

import (
	"testing"
)

func TestFileDoesNotExist(t *testing.T) {
	fileName := "doesNotExist"
	_, err := BuildFrequencyTable(fileName)

	if err == nil {
		t.Errorf("expected error")
	}
}

func TestFrequencyTable(t *testing.T) {
	fileName := "test.txt"
	table, _ := BuildFrequencyTable(fileName)
	assertOccurrence(t, table, "X", 333)
	assertOccurrence(t, table, "t", 223000)
}

func assertOccurrence(t *testing.T, table map[string]int, char string, expected int) {
	result, _ := table[char]
	if result != expected {
		t.Errorf("Error: expected %v, got: %v", expected, result)
	}
}
