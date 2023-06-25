package json_parser

import (
	"testing"
)

func TestStepOneValid(t *testing.T) {
	result := isValidJsonFile("tests/step1/valid.json")
	expected := true
	if result != expected {
		t.Errorf("result is %v; expected %v", result, expected)
	}
}

func TestStepOneInvalid(t *testing.T) {
	result := isValidJsonFile("tests/step1/invalid.json")
	expected := false
	if result != expected {
		t.Errorf("result is %v; expected %v", result, expected)
	}
}

func TestStepTwoValid(t *testing.T) {
	result := isValidJsonFile("tests/step2/valid.json")
	expected := true
	if result != expected {
		t.Errorf("result is %v; expected %v", result, expected)
	}
}

func TestStepTwoValid2(t *testing.T) {
	result := isValidJsonFile("tests/step2/valid2.json")
	expected := true
	if result != expected {
		t.Errorf("result is %v; expected %v", result, expected)
	}
}

func TestStepTwoInvalid(t *testing.T) {
	result := isValidJsonFile("tests/step2/invalid.json")
	expected := false
	if result != expected {
		t.Errorf("result is %v; expected %v", result, expected)
	}
}

func TestStepTwoInvalid2(t *testing.T) {
	result := isValidJsonFile("tests/step2/invalid2.json")
	expected := false
	if result != expected {
		t.Errorf("result is %v; expected %v", result, expected)
	}
}

func TestStepThreeValid(t *testing.T) {
	result := isValidJsonFile("tests/step3/valid.json")
	expected := true
	if result != expected {
		t.Errorf("result is %v; expected %v", result, expected)
	}
}

func TestStepThreeInvalid(t *testing.T) {
	result := isValidJsonFile("tests/step3/invalid.json")
	expected := false
	if result != expected {
		t.Errorf("result is %v; expected %v", result, expected)
	}
}
