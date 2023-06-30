package json_parser

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

const (
	EOF          = "EOF"
	stringDelim  = "\""
	kvDelim      = ":"
	objectStart  = "{"
	objectEnd    = "}"
	arrayStart   = "["
	arrayEnd     = "]"
	elementDelim = ","
	space        = " "
	newLine      = "\n"
	tab          = "\t"
	carr         = "\r"
)

func validate(fileName string) bool {
	var tokens = make([]string, 0)
	var valid bool
	file, _ := os.Open(fileName)
	reader := bufio.NewReader(file)
	ch, _, _ := reader.ReadRune()
	char := string(ch)
	if char == objectStart {
		tokens, valid = parseObject(reader)
		if !valid {
			return false
		}
	} else if char == arrayStart {
		_, err := parseArray(reader)
		if err != nil {
			return false
		} else {
			return true
		}
	}
	defer file.Close()
	return validateStructure(tokens)
}

func validateStructure(tokens []string) bool {
	return len(tokens) >= 2 &&
		((tokens[0] == objectStart && tokens[len(tokens)-1] == objectEnd) ||
			(tokens[0] == arrayStart && tokens[len(tokens)-1] == arrayEnd))
}

func parseObject(reader *bufio.Reader) ([]string, bool) {
	var tokens = make([]string, 1)
	tokens[0] = objectStart

	// edge case empty object
	ch, _ := reader.Peek(1)
	if string(ch) == objectEnd {
		tokens = append(tokens, objectEnd)
		return tokens, true
	}

	for {
		// expect key need to take care of first "
		_, err := expect([]string{stringDelim}, reader)
		if err != nil {
			return tokens, false
		}
		token, err := parseString(reader)
		if err != nil {
			return tokens, false
		}
		tokens = append(tokens, token)

		// expect colon
		_, err = expect([]string{kvDelim}, reader)
		if err != nil {
			return tokens, false
		}
		tokens = append(tokens, kvDelim)

		// expect value
		token, _ = parseValue(reader)
		tokens = append(tokens, token)

		// expect comma or object end
		token, err = expect([]string{elementDelim, objectEnd}, reader)
		if err != nil {
			return tokens, false
		}
		tokens = append(tokens, token)
		if token == objectEnd {
			return tokens, true
		}
	}
}

func parseArray(reader *bufio.Reader) (string, error) {
	token := arrayStart
	for {
		t, err := parseValue(reader)
		if err != nil {
			return token, err
		}
		token += t

		t, err = expect([]string{elementDelim, arrayEnd}, reader)
		token += t
		if t == arrayEnd {
			return token, nil
		}
	}
}

func parseValue(reader *bufio.Reader) (string, error) {
	for {
		ch, _, err := reader.ReadRune()
		char := string(ch)

		if err != nil {
			return "", err
		}

		switch char {
		case space, newLine, carr, tab:
			continue
		case stringDelim:
			token, err := parseString(reader)
			if err != nil {
				return "", err
			}
			return token, nil
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "-":
			token, err := parseNumber(char, reader)
			if err != nil {
				return "", err
			}
			return token, nil
		case "t", "n", "f":
			token := expectedToken(char)
			err := parseToken(token[1:], reader)
			if err != nil {
				return "", err
			}
			return token, nil
		case arrayStart:
			token, err := parseArray(reader)
			if err != nil {
				return "", err
			}
			return token, nil
		case objectStart:
			tokens, success := parseObject(reader)
			token := strings.Join(tokens, ", ")
			if !success {
				return "", errors.New("parseValue: error parsing inner object: " + token)
			}
			return token, nil

		default:
			return "", errors.New("parseValue: error unexpected character: " + char)
		}
	}
}

func expect(expected []string, reader *bufio.Reader) (string, error) {
	for {
		ch, _, _ := reader.ReadRune()
		char := string(ch)

		for _, exp := range expected {
			if char == exp {
				return char, nil
			}
		}

		switch char {
		case space, newLine, carr, tab:
			continue
		default:
			return "", errors.New("expect: error, unexpected token: " + char)
		}
	}
}

func expectedToken(char string) string {
	token := ""
	if char == "t" {
		token = "true"
	}
	if char == "n" {
		token = "null"
	}
	if char == "f" {
		token = "false"
	}
	return token
}

func parseToken(expectedToken string, reader *bufio.Reader) error {
	token := ""
	for i := 0; i < len(expectedToken); i++ {
		ch, _, _ := reader.ReadRune()
		token += string(ch)
	}
	if token != expectedToken {
		return errors.New("parseToken: error, unexpected token: " + token)
	}
	return nil
}

func parseString(reader *bufio.Reader) (string, error) {
	token := ""
	escaped := false
	for {
		ch, _, _ := reader.ReadRune()
		char := string(ch)
		if char == newLine || char == carr {
			return token, errors.New("parseString: error, forbidden character in string: " + char)
		}
		if char == space {
			continue
		}
		if char == tab {
			return token, errors.New("parseString: error, tab is forbidden in string")
		}
		switch {
		case escaped:
			escaped = false
			switch char {
			case "b", "f", "n", "r", "t", "/", "\"", "\\":
				token += "\\" + char
			case "u":
				// TODO handle unicode
				token += "\\" + char
			default:
				return token, errors.New("parseString: error, forbidden character after escape: " + char)
			}

		case char == "\\":
			escaped = true
		case char == stringDelim:
			return token, nil
		default:
			token += char
		}
	}
}

func parseNumber(first string, reader *bufio.Reader) (string, error) {
	token := first
	// case leading zeros:
	if first == "0" {
		bytes, _ := reader.Peek(1)
		char := ""
		if len(bytes) == 0 {
			return token, nil
		} else {
			char = string(bytes[0])
		}
		if char != "e" && char != "E" && char != "." && char != elementDelim && char != space {
			bytes, err := reader.Peek(5)
			token = ""
			if err == nil {
				token = string(bytes)
			}
			return token, errors.New("parseNumber: error leading zeros in token near: " + token)
		}
	}

	for {
		bytes, _ := reader.Peek(1)
		char := string(bytes[0])

		switch char {
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			ch, _, _ := reader.ReadRune()
			token += string(ch)
		case space, elementDelim, objectEnd, newLine, carr:
			return token, nil
		case ".":
			token, err := parseAfterDot(token, reader)
			if err != nil {
				return token, err
			}
			return token, nil
		case "E", "e":
			token, err := parseAfterExp(token, reader)
			if err != nil {
				return token, err
			}
			return token, nil
		default:
			return token, errors.New("ParseNumber: error, unexpected character: " + char)
		}
	}
}

func parseAfterDot(token string, reader *bufio.Reader) (string, error) {
	// scan and add first dot
	ch, _, _ := reader.ReadRune()
	token += string(ch)

	for {
		bytes, _ := reader.Peek(1)
		char := string(bytes[0])

		switch char {
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			ch, _, _ := reader.ReadRune()
			token += string(ch)
		case space, elementDelim, objectEnd, newLine, EOF:
			return token, nil
		case "e", "E":
			token, err := parseAfterExp(token, reader)
			if err != nil {
				return token, err
			}
			return token, nil
		default:
			return token, errors.New("AfterDot: error, unexpected character: " + char)
		}
	}
}

func parseAfterExp(token string, reader *bufio.Reader) (string, error) {
	// scan and add the exp sign
	ch, _, _ := reader.ReadRune()
	token += string(ch)
	hasSign := false

	for {
		bytes, _ := reader.Peek(1)
		char := string(bytes[0])

		switch char {
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			ch, _, _ := reader.ReadRune()
			token += string(ch)
		case space, elementDelim, objectEnd, newLine, EOF:
			return token, nil
		case "-", "+":
			if hasSign {
				return token, errors.New("parseAfterExp: error, token " + token + " has two signs")
			}
			hasSign = true
			ch, _, _ := reader.ReadRune()
			token += string(ch)
		default:
			return token, errors.New("parseAfterExp: error, unexpected character: " + char)
		}
	}
}
