package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const (
	ErrCompile = iota + 65
)

var hadError = false

func run(source []byte, filename string) {
	tokens := NewScanner(source, filename).ScanAll()
	for _, token := range tokens {
		if token.kind != EOF {
			fmt.Printf("%v\n", token)
		}
	}
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

		run(bytes, "repl")
		// if bytes[0] != '\n' {
		// 	run(bytes)
		// }
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
