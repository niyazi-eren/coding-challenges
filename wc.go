package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	bufferSize = 2048
)

func main() {
	if len(os.Args) != 4 {
		usage()
	}

	cmd := os.Args[1]
	if cmd != "ccwc" {
		usage()
	}

	opt := os.Args[2]
	if opt != "-c" && opt != "-l" && opt != "-w" && opt != "-m" {
		usage()
	}

	fileName := os.Args[3]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	if opt == "-c" {
		nbytes, err := countBytes(file)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}
		fmt.Println(nbytes, fileName)
	}

	if opt == "-l" {
		nlines, err := countLines(file)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}
		fmt.Println(nlines, fileName)
	}

	if opt == "-w" {
		nwords, err := countWords(file)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}
		fmt.Println(nwords, fileName)
	}

	if opt == "-m" {
		nchars, err := countChars(file)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}
		fmt.Println(nchars, fileName)
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
	return totalBytes, nil
}

func usage() {
	panic("Usage: ccwc [-clwm] <file>")
}
