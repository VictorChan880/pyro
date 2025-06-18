package main

import (
	"bufio"
	"fmt"
	"os"
)

var hasError bool

func main() {
	// Create a reader to read input from the user
	args := os.Args
	hasError = false
	for _, arg := range args {
		fmt.Println(arg)
	}
	if len(args) > 2 {
		fmt.Println("Usage: goPyro [script]")
		return
	} else if len(args) == 2 {
		runFile(args[1])
	} else {
		runPrompt()
	}

}
func runFile(fileName string) {
	fileBytes, err := os.ReadFile(fileName)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	run(string(fileBytes))
}
func runPrompt() {
	input := bufio.NewReader(os.Stdin)
	for {
		hasError = false
		line, err := input.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading line:", err)
			return
		}
		if line == "" {
			break
		}
		run(line)

	}
}

func run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.scanTokens()
	// fmt.Println(tokens)
	parser := Parser{tokens, 0}
	statements, _ := parser.parse()
	// fmt.Println(statements)
	interpreter := Interpreter{
		Environment: NewEnvironment(),
	}
	interpreter.interpret(statements)

	if hasError {
		os.Exit(65)
	}

}
