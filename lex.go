package main

import "unicode/utf8"

type Lexer struct {
	current int
	start   int
	line    int

	source []byte
	tokens []Token
}

func NewLexer(source []byte) Lexer {
	return Lexer{
		current: 0,
		line:    1,
		start:   0,

		source: source,
		tokens: []Token{},
	}
}

func (l *Lexer) addToken(kind TokenKind) {
	text := l.source[l.start:l.current]
	l.tokens = append(l.tokens, Token{
		kind:   kind,
		lexeme: text,
		line:   l.line,
	})
}

func (l *Lexer) advance() rune {
	l.current++
	r, _ := utf8.DecodeRune(l.source[l.current:])
	return r
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) readToken() {
	c := l.advance()
	switch c {
	case '(':
		l.addToken(LEFT_PAREN)
	case ')':
		l.addToken(RIGHT_PAREN)
	case '{':
		l.addToken(LEFT_BRACE)
	case '}':
		l.addToken(RIGHT_BRACE)
	case ',':
		l.addToken(COMMA)
	case '.':
		l.addToken(DOT)
	case '-':
		l.addToken(MINUS)
	case '+':
		l.addToken(PLUS)
	case ';':
		l.addToken(SEMICOLON)
	case '*':
		l.addToken(STAR)
	}
}

func (l *Lexer) readTokens() []Token {
	for !l.isAtEnd() {
		l.start = l.current
		l.readToken()
	}

	l.tokens = append(l.tokens, Token{
		kind:    EOF,
		lexeme:  []byte{},
		line:    l.line,
		literal: nil,
	})
	return l.tokens
}
