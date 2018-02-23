package main

import "fmt"

func report(line int, where string, message string) {
	fmt.Println("[line ", line, "] Error", where, ": ", message)
}
