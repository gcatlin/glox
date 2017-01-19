package main

// http://www.craftinginterpreters.com/scanning.html#lexeme-type

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var hadError = false

type Token string

type Lexer struct {
	source string
}

func (l *Lexer) scanTokens() []Token {
	tokens := make([]Token, len(l.source))
	for i, t := range strings.Split(l.source, "") {
		tokens[i] = Token(t)
	}
	return tokens
}

func error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Println("[line ", line, "] Error", where, ": ", message)
	hadError = true
}

func run(s string) {
	lexer := Lexer{source: s}
	tokens := lexer.scanTokens()

	fmt.Print(">>> ")
	for _, token := range tokens {
		fmt.Print(".", string(token))
	}
}

func runFile(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot read file", err)
	}
	run(string(bytes))

	if hadError {
		os.Exit(65)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Println("Bye!")
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot read input:", err)
		}

		run(str)
		hadError = false
	}
}

func main() {
	argv := os.Args[1:]
	argc := len(argv)
	if argc > 1 {
		fmt.Println("Usage: glox [script]")
	} else if argc == 1 {
		runFile(argv[0])
	} else {
		runPrompt()
	}
}
