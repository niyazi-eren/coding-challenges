package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	bufferSize = 2048
)

func main() {
	numArgs := len(os.Args)
	if numArgs != 3 && numArgs != 4 {
		usage()
		return
	}

	cmd := os.Args[1]
	if cmd != "ccwc" {
		usage()
		return
	}

	fileName := parseFileName(os.Args)
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}

	opt := parseOptions(os.Args)
	switch opt {
	case "c":
		nbytes, _ := countBytes(file)
		fmt.Println(nbytes, fileName)
	case "l":
		nlines, _ := countLines(file)
		fmt.Println(nlines, fileName)
	case "w":
		nwords, _ := countWords(file)
		fmt.Println(nwords, fileName)
	case "m":
		nchars, _ := countChars(file)
		fmt.Println(nchars, fileName)
	case "":
		nlines, _ := countLines(file)
		nbytes, _ := countBytes(file)
		nwords, _ := countWords(file)
		fmt.Println(nlines, nbytes, nwords, fileName)
	default:
		usage()
	}
	defer file.Close()
}

func parseOptions(args []string) string {
	if len(args) == 3 {
		return ""
	} else {
		opt := args[2]
		// get chars after the -
		return strings.Split(opt, "-")[1]
	}
}

func parseFileName(args []string) string {
	if len(args) == 3 {
		return args[2]
	} else {
		return args[3]
	}
}

func countLines(file *os.File) (int, error) {
	scanner := bufio.NewScanner(file)
	numLines := 0
	for scanner.Scan() {
		scanner.Text()
		numLines++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}
	file.Seek(0, 0)
	return numLines, nil
}

func countChars(file *os.File) (int, error) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanBytes)

	numChars := 0
	for scanner.Scan() {
		scanner.Text()
		numChars++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}
	file.Seek(0, 0)
	return numChars, nil
}

func countWords(file *os.File) (int, error) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	numWords := 0
	for scanner.Scan() {
		scanner.Text()
		numWords++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}
	file.Seek(0, 0)
	return numWords, nil
}

func countBytes(file *os.File) (int, error) {
	buffer := make([]byte, bufferSize)
	totalBytes := 0

	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}
		totalBytes += bytesRead
	}
	file.Seek(0, 0)
	return totalBytes, nil
}

func usage() {
	panic("Usage: ccwc [-clwm] <file>")
}
