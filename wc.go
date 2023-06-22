package main

import (
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
	if opt != "-c" {
		usage()
	}

	fileName := os.Args[3]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	nbytes, err := countBytes(file)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	fmt.Println(nbytes, fileName)
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
	panic("Usage: ccwc -c <file>")
}
