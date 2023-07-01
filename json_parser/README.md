# Write Your Own JSON Parser

JSON parser implementation for https://codingchallenges.fyi/challenges/challenge-json-parser/

Building a JSON parser is an easy way to learn about parsing techniques which are useful for everything from parsing simple data formats through to building a fully featured compiler for a programming language.

Parsing is often broken up into two stages: lexical analysis and syntactic analysis. Lexical analysis is the process of dividing a sequence of characters into meaningful chunks, called tokens. Syntactic analysis (which is also sometimes referred to as parsing) is the process of analysing the list of tokens to match it to a formal grammar.

You can read far more about building lexers, parses and compilers in what is regarded as the definitive book on compilers: Compilers: Principles, Techniques, and Tools - widely known as the “Dragon Book” (because there’s an illustration of a dragon on the cover).
## Installation
To run the project, you need to have Go installed on your system. You can download and install the latest version of Go from the official Go website: https://golang.org/

## Usage
In the json_parser directory, run the command `go test -V` to run the  tests