package ast

import (
	"fmt"
	"zimlit/graphene/token"
)

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
	return fmt.Sprintf("(%s %s)", l.Kind, l.Value)
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
