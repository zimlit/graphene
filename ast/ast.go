package ast

import (
	"fmt"
	"zimlit/graphene/token"
)

type Expr interface {
	expr()
	String() string
}

type Binary struct {
	left     Expr
	operator token.Token
	right    Expr
}

func (b Binary) expr() {}
func (b Binary) String() string {
	return fmt.Sprintf(
		"(%s %s %s)",
		b.operator.Kind.String(),
		b.left.String(),
		b.right.String(),
	)
}

func NewBinary(left Expr, operator token.Token, right Expr) Binary {
	return Binary{
		left:     left,
		operator: operator,
		right:    right,
	}
}

type Unary struct {
	operator token.Token
	right    Expr
}

func (u Unary) expr() {}
func (u Unary) String() string {
	return fmt.Sprintf(
		"(%s %s)",
		u.operator.Kind.String(),
		u.right.String(),
	)
}

func NewUnary(operator token.Token, right Expr) Unary {
	return Unary{
		operator: operator,
		right:    right,
	}
}

type Literal struct {
	value string
}

func (l Literal) expr() {}
func (l Literal) String() string {
	return l.value
}

func NewLiteral(value string) Literal {
	return Literal{value}
}

type Grouping struct {
	inner Expr
}

func (g Grouping) expr() {}
func (g Grouping) String() string {
	return g.inner.String()
}

func NewGrouping(inner Expr) Grouping {
	return Grouping{inner}
}
