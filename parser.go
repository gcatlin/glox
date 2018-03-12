package main

type Err int

const (
	ParseError Err = iota
)

type Parser struct {
	current  int
	filename string
	source   []byte
	tokens   []Token
}

func (p *Parser) addition() Expr {
	expr := p.multiplication()
	for p.match(MINUS, PLUS) {
		op := p.previous()
		rhs := p.multiplication()
		expr = BinaryExpr{op: op, lhs: expr, rhs: rhs}
	}
	return expr
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(kind TokenKind) bool {
	return !p.isAtEnd() && p.peek().kind == kind
}

func (p *Parser) comparison() Expr {
	expr := p.addition()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		op := p.previous()
		rhs := p.addition()
		expr = BinaryExpr{op: op, lhs: expr, rhs: rhs}
	}
	return expr
}

func (p *Parser) consume(kind TokenKind, message string) Token {
	if p.check(kind) {
		return p.advance()
	}
	panic(p.err(p.peek(), message))
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		op := p.previous()
		rhs := p.comparison()
		expr = BinaryExpr{op: op, lhs: expr, rhs: rhs}
	}
	return expr
}

func (p *Parser) err(token Token, message string) Err {
	reportError(p.filename, token.line, token.col, len(token.lexeme), p.getLine(token),
		"[parser] "+message)
	return ParseError
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) getLine(token Token) string {
	// Find start of line
	start, end := 0, len(p.source)
	for i := token.line - 1; i > 0; i-- {
		for p.source[start] != '\n' {
			start++
		}
	}
	// Find end of line
	for j := start; j < end; j++ {
		if p.source[j] == '\n' {
			end = j
			break
		}
	}
	return string(p.source[start:end])
}

func (p *Parser) isAtEnd() bool {
	return p.peek().kind == EOF
}

func (p *Parser) match(kinds ...TokenKind) bool {
	for _, kind := range kinds {
		if p.check(kind) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) multiplication() Expr {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		op := p.previous()
		rhs := p.unary()
		expr = BinaryExpr{op: op, lhs: expr, rhs: rhs}
	}
	return expr
}

func (p *Parser) Parse() Expr {
	defer func() {
		// See https://github.com/golang/go/wiki/PanicAndRecover
		if r := recover(); r != ParseError {
			// panic(r)
		}
	}()
	return p.expression()
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return LiteralExpr{BoolLiteral(false)}
	}
	if p.match(TRUE) {
		return LiteralExpr{BoolLiteral(true)}
	}
	if p.match(NIL) {
		return LiteralExpr{nil}
	}
	if p.match(NUMBER, STRING) {
		return LiteralExpr{p.previous().literal}
	}
	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expected ')' after expression.")
		return GroupingExpr{expr}
	}

	panic(p.err(p.peek(), "Expected an expression."))
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().kind == SEMICOLON {
			return
		}

		switch p.peek().kind {
		case CLASS:
		case FN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}
		p.advance()
	}
}

func (p *Parser) unary() Expr {
	for p.match(BANG, MINUS) {
		op := p.previous()
		rhs := p.unary()
		return UnaryExpr{op: op, rhs: rhs}
	}
	return p.primary()
}
