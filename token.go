package main

// http://www.craftinginterpreters.com/scanning.html#lexeme-type

type TokenKind int

const (
	ILLEGAL TokenKind = iota
	EOF

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
)

type Token struct {
	kind    TokenKind
	lexeme  []byte      // string?
	literal interface{} // use an interface?
	line    int
}

func (t Token) String() string {
	// return t.kind + " " + t.lexeme + " " + t.literal
	return string(t.lexeme)
}
