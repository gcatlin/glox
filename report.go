package main

import (
	"fmt"
	"os"
)

const (
	FILENAME_STYLE = ANSI_RESET
	LINE_STYLE     = ANSI_RESET
	LINE_NUM_STYLE = ANSI_RESET + ANSI_FG_BLUE
	COL_NUM_STYLE  = ANSI_RESET + ANSI_FG_CYAN
	MESSAGE_STYLE  = ANSI_RESET + ANSI_BOLD
	INFO_STYLE     = ANSI_RESET + ANSI_FG_GREEN + ANSI_BOLD
	ERROR_STYLE    = ANSI_RESET + ANSI_FG_RED + ANSI_BOLD
)

type LogLevel int

const (
	Info LogLevel = iota
	Error
)

type LogConfig struct {
	level string
	style string
	line  string
}

var LogLevelConfig = map[LogLevel]LogConfig{
	Info:  {level: "info", style: INFO_STYLE, line: "-----------------------------"},
	Error: {level: "error", style: ERROR_STYLE, line: "^^^^^^^^^^^^^^^^^^^^^^^^^^^^^"},
}

//
// https://blog.rust-lang.org/2016/08/10/Shape-of-errors-to-come.html
//
// error[E0499]: cannot borrow `foo.bar1` as mutable more than once at a time
//   --> src/test/compile-fail/borrowck/borrowck-borrow-for-owned-ptr.rs:29:22
//    |
// 28 |      let bar1 = &mut foo.bar1;
//    |                      -------- first mutable borrow occurs here
// 29 |      let _bar2 = &mut foo.bar1;
//    |                       ^^^^^^^^ second mutable borrow occurs here
// 30 |      *bar1;
// 31 |  }
//    |  - first borrow ends here
//
func report(level LogLevel, filename string, line, col, len int, srcLine, message string) {
	config := LogLevelConfig[level]
	padding := countDigits(line)

	// Message
	fmt.Fprintf(os.Stderr, "%s%s", config.style, config.level)
	fmt.Fprintf(os.Stderr, MESSAGE_STYLE+": %s\n", message)
	fmt.Fprintf(os.Stderr, LINE_NUM_STYLE+" %*s--> ", padding, "")
	fmt.Fprintf(os.Stderr, FILENAME_STYLE+"%s", filename)
	fmt.Fprintf(os.Stderr, LINE_NUM_STYLE+":%d"+COL_NUM_STYLE+":%d\n", line, col+1)
	fmt.Fprintf(os.Stderr, LINE_NUM_STYLE+" %*s | \n", padding, "")
	fmt.Fprintf(os.Stderr, LINE_NUM_STYLE+" %d | ", line)

	// Code
	fmt.Fprintf(os.Stderr, LINE_STYLE+"%s", srcLine[:col])
	fmt.Fprintf(os.Stderr, "%s%s", config.style, srcLine[col:col+len])
	fmt.Fprintf(os.Stderr, LINE_STYLE+"%s\n", srcLine[col+len:])

	// Annotation
	fmt.Fprintf(os.Stderr, LINE_NUM_STYLE+" %*s | ", padding, "")
	fmt.Fprintf(os.Stderr, "%s%*s%.*s", config.style, col, "", len, config.line)
	fmt.Fprintf(os.Stderr, " %s%s"+ANSI_RESET+"\n", config.style, message)
}

func reportInfo(filename string, line, col, len int, srcLine, message string) {
	report(Info, filename, line, col, len, srcLine, message)
}

func reportError(filename string, line, col, len int, srcLine, message string) {
	report(Error, filename, line, col, len, srcLine, message)
	hadError = true
}

func countDigits(i int) int {
	switch {
	case i < 10:
		return 1
	case i < 100:
		return 2
	case i < 1000:
		return 3
	case i < 10000:
		return 4
	case i < 100000:
		return 5
	case i < 1000000:
		return 6
	case i < 10000000:
		return 7
	case i < 100000000:
		return 8
	default:
		return 9
	}
}
