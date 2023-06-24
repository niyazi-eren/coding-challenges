package json_parser

import (
	"testing"
)

func TestStepOneValid(t *testing.T) {
	result := Validate("tests/step1/valid.json")
	expected := true
	if result != expected {
		t.Errorf("result is %v; expected %v", result, expected)
	}
}

func TestStepOneInvalid(t *testing.T) {
	result := Validate("tests/step1/invalid.json")
	expected := false
	if result != expected {
		t.Errorf("result is %v; expected %v", result, expected)
	}
}
