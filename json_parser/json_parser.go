package json_parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"unicode"
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

func isSpecialToken(token string) bool {
	values := []string{kvDelim, start, end, elementDelim}
	for _, val := range values {
		if token == val {
			return true
		}
	}
	return false
}

func isValidKey(token string) bool {
	if isSpecialToken(token) {
		return false
	}
	return token[0] == '"' && token[len(token)-1] == '"'
}

func isValidValue(token string) bool {
	if isSpecialToken(token) {
		return false
	}
	// case string
	if token[0] == '"' && token[len(token)-1] == '"' {
		return true
	}
	// case number
	if unicode.IsDigit(rune(token[0])) {
		_, err := strconv.Atoi(token)
		return err == nil
	}

	return isBoolean(token) || isNull(token)
}

func isBoolean(token string) bool {
	return token == "true" || token == "false"
}

func isNull(token string) bool {
	return token == "null"
}

func isValidKeyValuePair(tokens []string) bool {
	if len(tokens) == 2 {
		return true
	}
	for i := 1; i < len(tokens)-1; i++ {
		// validate key
		if (i%4 == 1 && !isValidKey(tokens[i])) ||

			// validate delimiter
			(i%4 == 2 && (tokens[i] != kvDelim)) ||

			// validate value
			(i%4 == 3 && !isValidValue(tokens[i])) ||

			// validate element delimiter
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
			token += stringDelim
			continue
			// append to string
		} else if parsingString && char != stringDelim {
			token += char
			// finished parsing string
		} else if char == stringDelim && parsingString {
			parsingString = false
			token += stringDelim
			tokens = append(tokens, token)
			token = ""
		} else if char == space || char == newLine {
			continue
		} else if !isSpecialToken(char) {
			token += char
		} else if isSpecialToken(char) {
			if token != "" {
				tokens = append(tokens, token)
				token = ""
			}
			tokens = append(tokens, char)
		}
	}
	return tokens
}
