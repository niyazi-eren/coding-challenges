# Write Your Own wc Tool

WC tool implementation for https://codingchallenges.fyi/challenges/challenge-wc

The Unix command line tools are a great metaphor for good software engineering and they follow the Unix Philosophies of:

Writing simple parts connected by clean interfaces - each tool does just one thing and provides a simple CLI that handles text input from either files or file streams.
Design programs to be connected to other programs - each tool can be easily connected to other tools to create incredibly powerful compositions.
Following these philosophies has made the simple unix command line tools some of the most widely used software engineering tools - allowing us to create very complex text data processing pipelines from simple command line tools. There’s even a Coursera course on Linux and Bash for Data Engineering.

You can read more about the Unix Philosophy in the excellent book The Art of Unix Programming.
## Installation
To run the project, you need to have Go installed on your system. You can download and install the latest version of Go from the official Go website: https://golang.org/
Once you have Go installed, you can clone the wc repository or download the source code as a ZIP file. To clone the repository, open a terminal and run the following command:
`git clone https://github.com/niyazi-eren/coding-challenges.git`
## Usage
In the json_parser directory, run the command `go run wc.go ccwc [-clwm] test.txt`
- ccwc: The module name to use for the word count operation.
- [-clwm]: The flags to specify the counting options. Choose one or more of the following:
    - -c: Count the number of bytes in the input file. 
    - -l: Count the number of lines in the input file. 
    - -w: Count the number of words in the input file. 
    - -m: Count the number of characters in the input file. 
    - test.txt: The input file to perform the word count operation on. Replace test.txt with the actual path to your desired file.
## Test
In the json_parser directory, run the command `go test -V` to run the tests