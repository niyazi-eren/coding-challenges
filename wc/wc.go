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

	file, err := openFile(os.Args)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}

	opt := parseOptions(os.Args)
	switch opt {
	case "c":
		nbytes, _ := countBytes(file)
		fmt.Println(nbytes, file.Name())
	case "l":
		nlines, _ := countLines(file)
		fmt.Println(nlines, file.Name())
	case "w":
		nwords, _ := countWords(file)
		fmt.Println(nwords, file.Name())
	case "m":
		nchars, _ := countChars(file)
		fmt.Println(nchars, file.Name())
	case "":
		nlines, _ := countLines(file)
		nbytes, _ := countBytes(file)
		nwords, _ := countWords(file)
		fmt.Println(nlines, nbytes, nwords, file.Name())
	default:
		usage()
	}
	defer cleanup(file)
}

func cleanup(file *os.File) {
	file.Close()
	if file.Name() == "tmp.txt" {
		err := os.Remove(file.Name())
		if err != nil {
			fmt.Println("Error deleting file:", err)
			return
		}
	}
}

func parseOptions(args []string) string {
	if len(args) == 3 && args[2][0] == '-' {
		opt := args[2]
		return strings.Split(opt, "-")[1]
	} else if len(args) == 4 {
		opt := args[2]
		return strings.Split(opt, "-")[1]
	} else {
		return ""
	}
}

func openFile(args []string) (*os.File, error) {
	if len(args) == 3 && args[2][0] != '-' {
		fileName := args[2]
		return os.Open(fileName)
	} else if len(args) == 4 {
		fileName := args[3]
		return os.Open(fileName)
	} else {
		// input file case
		scanner := bufio.NewScanner(os.Stdin)
		tmpFile, _ := os.OpenFile("tmp.txt", os.O_APPEND|os.O_CREATE, 0644)

		for scanner.Scan() {
			input := scanner.Text()
			tmpFile.WriteString(input + "\n")
		}

		tmpFile.Seek(0, 0)
		return tmpFile, nil
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
