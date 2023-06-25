package json_parser

import (
	"bufio"
	"fmt"
	"os"
)

const (
	stringDelim  = "\""
	kvDelim      = ":"
	start        = "{"
	end          = "}"
	arrayStart   = "["
	arrayEnd     = "]"
	elementDelim = ","
	space        = " "
	newLine      = "\n"
)

func isValidJsonFile(fileName string) bool {
	tokens := tokenize(fileName)
	return len(tokens) >= 2 &&
		tokens[0] == start &&
		tokens[len(tokens)-1] == end &&
		isValidKeyValuePair(tokens)
}

func isForbiddenToken(token string) bool {
	values := []string{kvDelim, start, end, arrayStart, arrayEnd, elementDelim}
	for _, val := range values {
		if token == val {
			return true
		}
	}
	return false
}

func isValidKeyValuePair(tokens []string) bool {
	if !(len(tokens) > 2) {
		return true
	}
	for i := 1; i < len(tokens)-1; i++ {
		if ((i%4 == 1 || i%4 == 3) && isForbiddenToken(tokens[i])) ||
			(i%4 == 2 && (tokens[i] != kvDelim)) ||
			(i%4 == 0 && tokens[i] != elementDelim) ||
			// edge case where last token is a comma
			(i%4 == 0 && (i+4 > len(tokens)-1) && tokens[i] == elementDelim) {
			return false
		}
	}
	return true
}

func tokenize(fileName string) []string {
	file, _ := os.Open(fileName)
	defer file.Close()

	var tokens = make([]string, 0, 100)
	parsingString := false
	reader := bufio.NewReader(file)

	token := ""
	for {
		ch, _, err := reader.ReadRune()
		char := string(ch)

		if err != nil && err.Error() != "EOF" {
			fmt.Println("Error reading file:", err)
		}
		if err != nil && err.Error() == "EOF" {
			break
		}

		// handle key value pairs for string type
		if char == stringDelim && !parsingString {
			parsingString = true
			continue
		} else if char == stringDelim && parsingString {
			parsingString = false
			tokens = append(tokens, token)
			token = ""
		} else if parsingString && char != stringDelim {
			token += char
		} else if char != space && char != newLine {
			tokens = append(tokens, char)
		}
	}
	return tokens
}
