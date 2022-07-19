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
)

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

type Return struct {
	Value Expr
}

func (r Return) String() string {
	return fmt.Sprintf("(return %s)", r.Value.String())
}

func (r Return) Accept(v Visitor[any]) any {
	return v.visitReturnExpr(r)
}

func NewReturn(value Expr) Return {
	return Return{
		Value: value,
	}
}
