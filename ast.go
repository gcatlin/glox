package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Expr interface {
	// String() string
	Accept(v Visitor) string
}

type IntLiteral int
type FloatLiteral float64
type StringLiteral string

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
	// Visit(expr Expr) (v Visitor)
	visitBinaryExpr(expr BinaryExpr) string
	visitGroupingExpr(expr GroupingExpr) string
	visitLiteralExpr(expr LiteralExpr) string
	visitUnaryExpr(expr UnaryExpr) string
}

type BinaryExpr struct {
	op  Token
	lhs Expr
	rhs Expr
}

type GroupingExpr struct {
	expr Expr
}

type LiteralExpr struct {
	value Literal
}

type UnaryExpr struct {
	op  Token
	rhs Expr
}

func (expr BinaryExpr) Accept(v Visitor) string {
	return v.visitBinaryExpr(expr)
}

func (expr GroupingExpr) Accept(v Visitor) string {
	return v.visitGroupingExpr(expr)
}

func (expr LiteralExpr) Accept(v Visitor) string {
	return v.visitLiteralExpr(expr)
}

func (expr UnaryExpr) Accept(v Visitor) string {
	return v.visitUnaryExpr(expr)
}

type AstPrinter struct {
}

func (p AstPrinter) Print(expr Expr) {
	fmt.Println(expr.Accept(p))
}

func (p AstPrinter) parenthesize(name []byte, exprs ...Expr) string {
	// var b strings.Builder
	// for expr, _ := range exprs {
	// fmt.Fprintf(&b, "%d...", i)
	// }
	// b.WriteString("ignition")
	// fmt.Println(b.String())

	var b strings.Builder
	b.Grow(1024)

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
