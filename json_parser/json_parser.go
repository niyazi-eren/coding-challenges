package json_parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
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
	content := getContent(fileName)
	tokens := tokenize(content)
	return len(tokens) >= 2 &&
		tokens[0] == start &&
		tokens[len(tokens)-1] == end &&
		isValidKeyValuePair(tokens)
}

func getContent(fileName string) string {
	content := ""
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return ""
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content += scanner.Text()
	}
	return content
}

func tokenize(content string) []string {
	tokens := make([]string, 0, 200)
	var parsingString, parsingObject, parsingArray bool
	var token string

	for _, ch := range content {
		char := string(ch)

		switch {
		case parsingObject:
			if char == end {
				parsingObject = false
				token += char
				tokens = append(tokens, token)
				token = ""
			} else if char == space || char == newLine {
				continue
			} else {
				token += char
			}
		case parsingArray:
			if char == arrayEnd {
				parsingArray = false
				token += char
				tokens = append(tokens, token)
				token = ""
			} else {
				token += char
			}
		case char == stringDelim && !parsingString:
			parsingString = true
			token += stringDelim
			continue
		case parsingString && char != stringDelim:
			token += char
		case char == stringDelim && parsingString:
			parsingString = false
			token += stringDelim
			tokens = append(tokens, token)
			token = ""
		case char == space || char == newLine:
			continue
		case isSpecialToken(char):
			if char == start && len(tokens) == 0 {
				tokens = append(tokens, char)
			} else if char == start {
				parsingObject = true
				token += char
			} else if char == arrayStart {
				parsingArray = true
				token += char
			} else if token != "" {
				tokens = append(tokens, token, char)
				token = ""
			} else {
				tokens = append(tokens, char)
			}
		default:
			token += char
		}
	}
	return tokens
}

func isSpecialToken(token string) bool {
	values := []string{kvDelim, start, end, elementDelim, arrayStart, arrayEnd}
	for _, val := range values {
		if token == val {
			return true
		}
	}
	return false
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
			// check if is an inner object
			(i%4 == 3 && (string(tokens[i][0]) == start) && !isValidKeyValuePair(tokenize(tokens[i]))) ||
			// validate value otherwise
			(i%4 == 3 && (string(tokens[i][0]) != start) && !isValidValue(tokens[i])) ||
			// validate element delimiter
			(i%4 == 0 && tokens[i] != elementDelim) ||
			// edge case where last token is a comma
			(i%4 == 0 && (i+4 > len(tokens)-1) && tokens[i] == elementDelim) {
			return false
		}
	}
	return true
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

	if unicode.IsDigit(rune(token[0])) || (rune(token[0]) == '-') {
		return isValidNumber(token)
	}

	// case array
	if string(token[0]) == arrayStart && string(token[len(token)-1]) == arrayEnd {
		// split array by values
		values := strings.Split(token[1:len(token)-1], elementDelim)
		for _, value := range values {
			// trim spaces
			value = strings.TrimSpace(value)

			if value != "" && !(isValidValue(value)) {
				return false
			}
		}
		return true
	}

	return isBoolean(token) || isNull(token)
}

func isValidNumber(token string) bool {
	// case with leading zeros
	if len(token) > 2 && token[0] == '0' && !(token[1] == 'e' || token[1] == 'E' || token[1] == '-' || token[1] == '.') {
		return false
	}
	exclude, _ := regexp.MatchString("^\\d+[+-]\\d+$", token)
	if exclude {
		return false
	}

	match, _ := regexp.MatchString("([-+]?\\d+(\\.\\d+)?([eE][-+]?\\d+)?|\\d+\n)", token)
	return match
}

func isBoolean(token string) bool {
	return token == "true" || token == "false"
}

func isNull(token string) bool {
	return token == "null"
}
