/*
	Copyright 2022 Devin Rockwell

    This file is part of Graphene.

    Graphene is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

    Graphene is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

    You should have received a copy of the GNU General Public License along with Graphene. If not, see <https://www.gnu.org/licenses/>.
*/

package ast

import (
	"fmt"
	"strings"
	"zimlit/graphene/token"
)

type ValueKind interface {
	String() string
	vkind()
}

type Const uint8

const (
	INT Const = iota
	FLOAT
	STRING
)

func (c Const) vkind() {}
func (c Const) String() string {
	switch c {
	case INT:
		return "int"
	case FLOAT:
		return "float"
	case STRING:
		return "string"
	}
	panic("unreachable")
}

type Fn struct {
	Params []Param
	Rtype  ValueKind
}

func (f Fn) vkind() {}
func (f Fn) String() string {
	var str strings.Builder
	fmt.Fprintf(&str, "(fn %s (", f.Rtype.String())
	for i, e := range f.Params {
		fmt.Fprintf(&str, "%s", e.Kind.String())
		if i+1 != len(f.Params) {
			fmt.Fprint(&str, " ")
		}
	}
	fmt.Fprintf(&str, "))")
	return str.String()
}
func NewFnT(params []Param, rtype ValueKind) Fn {
	return Fn{
		Params: params,
		Rtype:  rtype,
	}
}

type Exprs []Expr

type Expr interface {
	String() string
	Accept(v Visitor[any]) any
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

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
	visitFnExpr(f FnExpr) R
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

func (i IfExpr) String() string {
	var str strings.Builder

	fmt.Fprintf(&str, "(if %s (", i.Condition.String())
	for j, b := range i.Body {
		fmt.Fprint(&str, b)
		if j+1 != len(i.Body) {
			fmt.Fprint(&str, " ")
		}
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
	Name  string
	Value Expr
}

func (a Assignment) String() string {
	return fmt.Sprintf("(= %s %s)", a.Name, a.Value.String())
}

func (a Assignment) Accept(v Visitor[any]) any {
	return v.visitAssignment(a)
}

func NewAssignment(name string, value Expr) Assignment {
	return Assignment{
		Name:  name,
		Value: value,
	}
}

type WhileExpr struct {
	Cond Expr
	Body []Expr
}

func (w WhileExpr) String() string {
	var str strings.Builder
	fmt.Fprintf(&str, "(while %s (", w.Cond.String())
	for i, e := range w.Body {
		fmt.Fprint(&str, e)
		if i+1 != len(w.Body) {
			fmt.Fprint(&str, " ")
		}
	}
	fmt.Fprintf(&str, "))")

	return str.String()
}

func (w WhileExpr) Accept(v Visitor[any]) any {
	return v.visitWhileExpr(w)
}

func NewWhileExpr(cond Expr, body []Expr) WhileExpr {
	return WhileExpr{
		Cond: cond,
		Body: body,
	}
}

type Param struct {
	Name string
	Kind ValueKind
}

func NewParam(name string, kind ValueKind) Param {
	return Param{
		Name: name,
		Kind: kind,
	}
}

type FnExpr struct {
	Params []Param
	Body   []Expr
	Rtype  ValueKind
}

func (f FnExpr) String() string {
	var str strings.Builder
	fmt.Fprintf(&str, "(fn %s (", f.Rtype.String())
	for i, e := range f.Params {
		fmt.Fprintf(&str, "%s: %s", e.Name, e.Kind.String())
		if i+1 != len(f.Params) {
			fmt.Fprint(&str, " ")
		}
	}
	fmt.Fprint(&str, ") (")
	for i, e := range f.Body {
		fmt.Fprint(&str, e)
		if i+1 != len(f.Body) {
			fmt.Fprint(&str, " ")
		}
	}
	fmt.Fprintf(&str, "))")
	return str.String()
}

func (f FnExpr) Accept(v Visitor[any]) any {
	return v.visitFnExpr(f)
}

func NewFn(params []Param, body []Expr, rtype ValueKind) FnExpr {
	return FnExpr{
		Params: params,
		Body:   body,
		Rtype:  rtype,
	}
}
