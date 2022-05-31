package ast

import (
	"fmt"
	"zimlit/graphene/token"
)

type ValueKind uint8

const (
	INT   = iota
	FLOAT = iota
)

func NewKind(t token.TokenKind) ValueKind {
	switch t {
	case token.INTK:
		return INT
	case token.FLOATK:
		return FLOAT
	}
	panic("unreachable")
}

func (k ValueKind) String() string {
	switch k {
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	default:
		return "INVALID"
	}
}

type Expr interface {
	expr()
	String() string
	Accept(v Visitor[any]) any
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (b Binary) expr() {}
func (b Binary) String() string {
	return fmt.Sprintf(
		"(%s %s %s)",
		b.Operator.Kind.String(),
		b.Left.String(),
		b.Right.String(),
	)
}
func (b Binary) Accept(v Visitor[any]) any {
	return v.visitBinary(b)
}

func NewBinary(left Expr, operator token.Token, right Expr) Binary {
	return Binary{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (u Unary) expr() {}
func (u Unary) String() string {
	return fmt.Sprintf(
		"(%s %s)",
		u.Operator.Kind.String(),
		u.Right.String(),
	)
}

func (u Unary) Accept(v Visitor[any]) any {
	return v.visitUnary(u)
}

func NewUnary(operator token.Token, right Expr) Unary {
	return Unary{
		Operator: operator,
		Right:    right,
	}
}

type Literal struct {
	Value string
	Kind  token.TokenKind
}

func (l Literal) expr() {}
func (l Literal) String() string {
	return fmt.Sprint(l.Value)
}

type Visitor[R any] interface {
	visitBinary(b Binary) R
	visitUnary(u Unary) R
	visitLiteral(l Literal) R
	visitGrouping(g Grouping) R
	visitVarDecl(v VarDecl) R
}

func (l Literal) Accept(v Visitor[any]) any {
	return v.visitLiteral(l)
}

func NewLiteral(value string, kind token.TokenKind) Literal {
	return Literal{value, kind}
}

type Grouping struct {
	Inner Expr
}

func (g Grouping) expr() {}
func (g Grouping) String() string {
	return g.Inner.String()
}

func (g Grouping) Accept(v Visitor[any]) any {
	return v.visitGrouping(g)
}

func NewGrouping(inner Expr) Grouping {
	return Grouping{inner}
}

type VarDecl struct {
	Name  string
	Kind  ValueKind
	Value Expr
}

func (v VarDecl) expr() {}
func (v VarDecl) String() string {
	return fmt.Sprintf("(let %s %s %s)", v.Name, v.Kind.String(), v.Value.String())
}

func (va VarDecl) Accept(v Visitor[any]) any {
	return v.visitVarDecl
}

func NewVarDecl(name string, kind ValueKind, value Expr) VarDecl {
	return VarDecl{
		Name:  name,
		Kind:  kind,
		Value: value,
	}
}
