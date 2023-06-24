package json_parser

import (
	"bufio"
	"fmt"
	"os"
)

func Validate(fileName string) bool {
	var tokens = make([]string, 0, 100)
	file, err := os.Open(fileName)
	if err != nil {
		return false
	}

	reader := bufio.NewReader(file)
	for {
		char, _, err := reader.ReadRune()
		if err != nil && err.Error() != "EOF" {
			fmt.Println("Error reading file:", err)
			return false
		}
		if err != nil && err.Error() == "EOF" {
			break
		}
		token := string(char)
		tokens = append(tokens, token)
	}

	defer file.Close()

	return len(tokens) >= 2 && tokens[0] == "{" && tokens[len(tokens)-1] == "}"

}
