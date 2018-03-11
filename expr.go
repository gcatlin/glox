package main

type Expr interface {
	Accept(v Visitor) string
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
