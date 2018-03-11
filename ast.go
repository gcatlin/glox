package main

import (
	"fmt"
	"strconv"
	"strings"
)

type BoolLiteral bool
type IntLiteral int
type FloatLiteral float64
type StringLiteral string

func (b BoolLiteral) String() string {
	if b {
		return "true"
	}
	return "false"
}

func (f FloatLiteral) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

func (i IntLiteral) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (s StringLiteral) String() string {
	return string(s)
}

type Visitor interface {
	visitBinaryExpr(expr BinaryExpr) string
	visitGroupingExpr(expr GroupingExpr) string
	visitLiteralExpr(expr LiteralExpr) string
	visitUnaryExpr(expr UnaryExpr) string
}

type AstPrinter struct {
}

func (p AstPrinter) Print(expr Expr) {
	fmt.Println(expr.Accept(p))
}

func (p AstPrinter) parenthesize(name []byte, exprs ...Expr) string {
	var b strings.Builder

	b.WriteByte('(')
	b.Write(name)
	for _, expr := range exprs {
		b.WriteByte(' ')
		b.WriteString(expr.Accept(p))
	}
	b.WriteByte(')')

	return b.String()
}

func (p AstPrinter) visitBinaryExpr(expr BinaryExpr) string {
	return p.parenthesize(expr.op.lexeme, expr.lhs, expr.rhs)
}

func (p AstPrinter) visitGroupingExpr(expr GroupingExpr) string {
	return p.parenthesize([]byte("group"), expr.expr)
}

func (p AstPrinter) visitLiteralExpr(expr LiteralExpr) string {
	if expr.value == nil {
		return "nil"
	}
	return expr.value.String()
}

func (p AstPrinter) visitUnaryExpr(expr UnaryExpr) string {
	return p.parenthesize(expr.op.lexeme, expr.rhs)
}
