package main

var hadError bool = false

func error(line int, message string) {
	reportError(line, "", message)
}

func reportError(line int, where string, message string) {
	report(line, where, message)
	hadError = true
}
