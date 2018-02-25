package main

// http://www.craftinginterpreters.com/scanning.html#lexeme-type

type Literal interface{}

type TokenKind int

const (
	// EOF = -(iota+1)
	ILLEGAL TokenKind = iota

	// Single-character tokens
	LEFT_PAREN  // (
	RIGHT_PAREN // )
	LEFT_BRACE  // {
	RIGHT_BRACE // {
	COMMA       // ,
	DOT         // .
	MINUS       // -
	PLUS        // +
	SEMICOLON   // ;
	SLASH       // /
	STAR        // *

	// One or two character tokens
	BANG          // !
	BANG_EQUAL    // !=
	EQUAL         // =
	EQUAL_EQUAL   // ==
	GREATER       // >
	GREATER_EQUAL // >=
	LESS          // <
	LESS_EQUAL    // <=

	// Literals
	IDENTIFIER
	STRING
	NUMBER

	// Keywords
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	// Special
	COMMENT
	EOF
)

type Token struct {
	kind    TokenKind
	lexeme  []byte
	literal Literal
	line    int
	col     int
}

func (t Token) String() string {
	// return t.kind + " " + t.lexeme + " " + t.literal
	return string(t.lexeme)
}

var Keywords = map[string]TokenKind{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}
