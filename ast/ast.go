package ast

import (
	"fmt"
	"strings"
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
	visitIfExpr(i IfExpr) R
	visitAssignment(a Assignment) R
	visitWhileExpr(w WhileExpr) R
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
	Name   string
	Kind   ValueKind
	is_mut bool
	Value  Expr
}

func (v VarDecl) expr() {}
func (v VarDecl) String() string {
	if !v.is_mut {
		return fmt.Sprintf("(let %s %s %s)", v.Name, v.Kind.String(), v.Value.String())
	} else {
		return fmt.Sprintf("(let mut %s %s %s)", v.Name, v.Kind.String(), v.Value.String())
	}
}

func (va VarDecl) Accept(v Visitor[any]) any {
	return v.visitVarDecl(va)
}

func NewVarDecl(name string, kind ValueKind, value Expr, is_mut bool) VarDecl {
	return VarDecl{
		Name:   name,
		Kind:   kind,
		is_mut: is_mut,
		Value:  value,
	}
}

type IfExpr struct {
	Condition Expr
	Body      []Expr
	Else_ifs  []IfExpr
	Else      []Expr
}

func (i IfExpr) expr() {}
func (i IfExpr) String() string {
	var str strings.Builder

	fmt.Fprintf(&str, "(if %s (", i.Condition.String())
	for _, b := range i.Body {
		fmt.Fprint(&str, b)
	}
	if i.Else_ifs == nil {
		fmt.Fprint(&str, ")")
	} else {
		fmt.Fprint(&str, ") ")
	}

	for _, e := range i.Else_ifs {
		fmt.Fprintf(&str, "%s ", e.String())
	}
	if i.Else != nil {
		fmt.Fprint(&str, "(")
	}
	for _, e := range i.Else {
		fmt.Fprintf(&str, "%s", e.String())
	}
	if i.Else != nil {
		fmt.Fprint(&str, ")")
	}
	fmt.Fprint(&str, ")")

	return str.String()
}

func (i IfExpr) Accept(v Visitor[any]) any {
	return v.visitIfExpr(i)
}

func NewIfExpr(condition Expr, body []Expr, else_ifs []IfExpr, el []Expr) IfExpr {
	return IfExpr{
		Condition: condition,
		Body:      body,
		Else_ifs:  else_ifs,
		Else:      el,
	}
}

type Assignment struct {
	name  string
	value Expr
}

func (a Assignment) expr() {}
func (a Assignment) String() string {
	return fmt.Sprintf("(= %s %s)", a.name, a.value.String())
}

func (a Assignment) Accept(v Visitor[any]) any {
	return v.visitAssignment(a)
}

func NewAssignment(name string, value Expr) Assignment {
	return Assignment{
		name:  name,
		value: value,
	}
}

type WhileExpr struct {
	cond Expr
	body []Expr
}

func (w WhileExpr) expr() {}
func (w WhileExpr) String() string {
	var str strings.Builder
	fmt.Fprintf(&str, "(while %s (", w.cond.String())
	for _, e := range w.body {
		fmt.Fprint(&str, e)
	}
	fmt.Fprintf(&str, "))")

	return str.String()
}

func (w WhileExpr) Accept(v Visitor[any]) any {
	return v.visitWhileExpr(w)
}

func NewWhileExpr(cond Expr, body []Expr) WhileExpr {
	return WhileExpr{
		cond: cond,
		body: body,
	}
}
