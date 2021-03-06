package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var hadError = false

func run(source []byte, filename string) {
	tokens := NewScanner(source, filename).ScanAll()
	if tokens[0].kind == EOF {
		return
	}

	parser := &Parser{current: 0, filename: filename, source: source, tokens: tokens}
	expr := parser.Parse()

	if hadError {
		return
	}

	var p AstPrinter
	p.Print(expr)
}

func runFile(path string) {
	bytes, _ := ioutil.ReadFile(path)
	run(bytes, path)
}

func repl() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(ANSI_BOLD + "glox> " + ANSI_RESET)
		bytes, err := reader.ReadBytes('\n')
		if err == io.EOF {
			fmt.Println(ANSI_RESET)
			break
		}

		run(bytes, "?")
		hadError = false
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		repl()
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		fmt.Println("Usage: glox [script]")
	}
}
