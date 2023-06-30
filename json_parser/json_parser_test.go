package json_parser

import (
	"strconv"
	"strings"
	"testing"
)

func TestStepOne(t *testing.T) {
	fileNames := []string{"tests/step1/valid.json", "tests/step1/invalid.json"}
	for _, fileName := range fileNames {
		assertValidate(t, fileName)
	}
}

func TestStepTwo(t *testing.T) {
	fileNames := []string{"tests/step2/valid.json",
		"tests/step2/valid2.json",
		"tests/step2/invalid.json",
		"tests/step2/invalid2.json"}
	for _, fileName := range fileNames {
		assertValidate(t, fileName)
	}
}

func TestStepThree(t *testing.T) {
	fileNames := []string{"tests/step3/valid.json", "tests/step3/invalid.json"}
	for _, fileName := range fileNames {
		assertValidate(t, fileName)
	}
}

func TestStepFour(t *testing.T) {
	fileNames := []string{"tests/step4/valid.json", "tests/step4/valid2.json", "tests/step4/invalid.json"}
	for _, fileName := range fileNames {
		assertValidate(t, fileName)
	}
}

func TestStepFivePass(t *testing.T) {
	for i := 1; i <= 3; i++ {
		fileName := "tests/step5/pass" + strconv.Itoa(i) + ".json"
		assertValidate(t, fileName)
	}
}

func TestStepFiveFail(t *testing.T) {
	for i := 10; i <= 33; i++ {
		fileName := "tests/step5/fail" + strconv.Itoa(i) + ".json"
		assertValidate(t, fileName)
	}
}

func assertValidate(t *testing.T, fileName string) {
	var expected bool
	if strings.Contains(fileName, "fail") || strings.Contains(fileName, "invalid") {
		expected = false
	} else {
		expected = true
	}
	result := validate(fileName)
	if result != expected {
		t.Errorf("result for file %v is %v; expected %v", fileName, result, expected)
	}
}

func TestStepFour2(t *testing.T) {
	fileNames := []string{"tests/step4/valid2.json"}
	for _, fileName := range fileNames {
		assertValidate(t, fileName)
	}
}
