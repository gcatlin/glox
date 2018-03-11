package main

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

const NUL = 0

// const NUL = '\000'

type Scanner struct {
	current   int
	start     int
	line      int
	col       int
	sourceLen int
	filename  string
	source    []byte
	tokens    []Token
}

func NewScanner(source []byte, filename string) *Scanner {
	return &Scanner{
		current:   0,
		start:     0,
		line:      1,
		col:       0,
		filename:  filename,
		sourceLen: len(source),
		source:    source,
		tokens:    make([]Token, 0, 256),
	}
}

func (s *Scanner) addToken(kind TokenKind) {
	s.addTokenLiteral(kind, nil)
}

func (s *Scanner) addTokenLiteral(kind TokenKind, literal Literal) {
	lexeme := s.source[s.start:s.current]
	token := Token{kind: kind, lexeme: lexeme, line: s.line, literal: literal}
	// s.info(&token)
	s.tokens = append(s.tokens, token)
}

func (s *Scanner) addTokenFor(ch rune, matched TokenKind, unmatched TokenKind) {
	if s.match(ch) {
		s.addToken(matched)
	} else {
		s.addToken(unmatched)
	}
}

func (s *Scanner) getLine(line int) string {
	pos, start, end := 0, 0, s.sourceLen-1
	for line--; line != 0; line-- {
		for s.source[pos] != '\n' && pos < end {
			pos++
		}
	}
	for start = pos; pos < end; pos++ {
		if s.source[pos] == '\n' {
			break
		}
	}
	return string(s.source[start:end])
}

func (s *Scanner) info(tok *Token) {
	len := len(tok.lexeme)
	reportInfo(s.filename, s.line, s.col-len, len, s.getLine(s.line), tok.kind.String())
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= s.sourceLen
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() || s.readRune(s.current) != expected {
		return false
	}

	s.current++
	s.col++
	return true
}

func (s *Scanner) Next() rune {
	ch := s.readRune(s.current)
	s.current++
	if ch == '\n' {
		s.line++
		s.col = 0
	} else {
		s.col++
	}

	return ch
}

func (s *Scanner) Peek() rune {
	if s.isAtEnd() {
		return NUL
	}
	return s.readRune(s.current)
}

func (s *Scanner) PeekNext() rune {
	if s.current+1 >= s.sourceLen {
		return NUL
	}
	return s.readRune(s.current + 1)
}

func (s *Scanner) readRune(offset int) rune {
	ch, _ := utf8.DecodeRune(s.source[offset:])
	return ch
}

func (s *Scanner) Scan() {
	ch := s.Next()
	switch ch {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		s.addTokenFor('=', BANG_EQUAL, BANG)
	case '=':
		s.addTokenFor('=', EQUAL_EQUAL, EQUAL)
	case '<':
		s.addTokenFor('=', LESS_EQUAL, LESS)
	case '>':
		s.addTokenFor('=', GREATER_EQUAL, GREATER)
	case '/':
		if s.match('/') {
			s.scanComment()
		} else {
			s.addToken(SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		break
	case '"':
		s.scanString()
	default:
		if isDigit(ch) {
			s.scanNumber()
		} else if isAlpha(ch) {
			s.scanIdentifier()
		} else {
			reportError(s.filename, s.line, s.col-1, 1, s.getLine(s.line),
				"Unexpected character: '"+string(ch)+"'")
			// exit
		}
	}
}

func (s *Scanner) ScanAll() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.Scan()
	}

	s.tokens = append(s.tokens, Token{kind: EOF, line: s.line})
	return s.tokens
}

func (s *Scanner) scanComment() {
	s.scanUntil('\n')
	s.addToken(COMMENT)
}

func (s *Scanner) scanIdentifier() {
	for isAlphaOrDigit(s.Peek()) {
		s.Next()
	}

	lexeme := string(s.source[s.start:s.current])
	kind, ok := Keywords[lexeme]
	if !ok {
		kind = IDENTIFIER
	}
	s.addToken(kind)
}

func (s *Scanner) scanNumber() {
	for isDigit(s.Peek()) {
		s.Next()
	}

	// Look for a fractional part
	if s.Peek() == '.' && isDigit(s.PeekNext()) {
		// Consume the "."
		s.Next()

		for isDigit(s.Peek()) {
			s.Next()
		}
	}

	// TODO error handling
	float, _ := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)
	s.addTokenLiteral(NUMBER, FloatLiteral(float))
}

func (s *Scanner) scanString() {
	s.scanUntil('"')
	if s.isAtEnd() {
		reportError(s.filename, s.line, s.col, 1, s.getLine(s.line),
			"Unterminated string.")
		return
	}

	// Consume the closing double-quote and return string excluding quotes
	s.Next()
	str := s.source[s.start+1 : s.current-1]
	s.addTokenLiteral(STRING, StringLiteral(str))
}

func (s *Scanner) scanUntil(until rune) {
	for s.Peek() != until && !s.isAtEnd() {
		s.Next()
	}
}

func isAlpha(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isAlphaOrDigit(ch rune) bool {
	return isAlpha(ch) || isDigit(ch)
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}
