package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func run(source []byte) {
	lexer := NewLexer(source)
	tokens := lexer.readTokens()

	fmt.Print(">>> ")
	for _, token := range tokens {
		fmt.Print(".", token)
	}
}

func runFile(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot read file", err)
	}
	run(bytes)

	if hadError {
		os.Exit(65) // TODO error constants
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		bytes, err := reader.ReadBytes('\n')
		fmt.Println(">>>", string(bytes))
		if err != nil {
			if err == io.EOF {
				fmt.Println("Bye!")
				break
			}
			fmt.Fprintln(os.Stderr, "cannot read input:", err)
		}
		run(bytes)
		hadError = false
	}
}

func main() {
	argv := os.Args[1:]
	argc := len(argv)
	if argc == 0 {
		runPrompt()
	} else if argc == 1 {
		runFile(argv[0])
	} else {
		fmt.Println("Usage: glox [script]")
	}
}
