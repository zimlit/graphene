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
