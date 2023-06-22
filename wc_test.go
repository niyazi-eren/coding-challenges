package main

import (
	"os"
	"testing"
)

var testFile, _ = os.Open("test.txt")

func TestStepOne(t *testing.T) {
	result, _ := countBytes(testFile)
	expected := 341836
	if result != expected {
		t.Errorf("result is %d bytes; expected %d", result, expected)
	}
}

func TestStepTwo(t *testing.T) {
	result, _ := countLines(testFile)
	expected := 7137
	if result != expected {
		t.Errorf("result is %d lines; expected %d", result, expected)
	}
}