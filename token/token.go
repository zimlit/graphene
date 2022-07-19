/*
	Copyright 2022 Devin Rockwell

	This file is part of Graphene.

	Graphene is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

	Graphene is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

	You should have received a copy of the GNU General Public License along with Graphene. If not, see <https://www.gnu.org/licenses/>.
*/

package token

import (
	"fmt"
	"strings"
)

type TokenKind uint8

const (
	INT = iota
	FLOAT
	PLUS
	MINUS
	STAR
	SLASH
	LPAREN
	RPAREN
	INTK
	FLOATK
	LET
	EQ
	EQEQ
	NEQ
	LESS
	GREATER
	LESSEQ
	GREATEREQ
	BANG
	IDENT
	COLON
	NIL
	IF
	ELSE
	ELSEIF
	END
	MUT
	WHILE
	STRING
	STRINGK
	FN
	COMMA
	RETURN
)

func (t TokenKind) String() string {
	switch t {
	case INT:
		return "integer literal"
	case FLOAT:
		return "float literal"
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case STAR:
		return "*"
	case SLASH:
		return "/"
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case INTK:
		return "int"
	case FLOATK:
		return "float"
	case LET:
		return "let"
	case EQ:
		return "="
	case EQEQ:
		return "=="
	case NEQ:
		return "!="
	case LESS:
		return "<"
	case GREATER:
		return ">"
	case LESSEQ:
		return "<="
	case GREATEREQ:
		return ">="
	case IDENT:
		return "identifier"
	case COLON:
		return ":"
	case NIL:
		return "nil"
	case IF:
		return "if"
	case ELSE:
		return "else"
	case ELSEIF:
		return "else if"
	case END:
		return "end"
	case MUT:
		return "mut"
	case STRING:
		return "string literal"
	case STRINGK:
		return "string"
	case FN:
		return "fn"
	case COMMA:
		return ","
	case RETURN:
		return "return"
	default:
		return "INVALID"
	}
}

type Token struct {
	Kind    TokenKind
	Literal string
	Line    int
	Col     int
}

func (t Token) String() string {
	var str strings.Builder

	fmt.Fprintln(&str, "Token {")
	fmt.Fprint(&str, "    kind:", t.Kind)
	fmt.Fprintln(&str, ",")
	fmt.Fprintf(&str, "    literal: \"%s\"", t.Literal)
	fmt.Fprintln(&str, ",")
	fmt.Fprint(&str, "   ", "line:", t.Line)
	fmt.Fprintln(&str, ",")
	fmt.Fprintln(&str, "   ", "col:", t.Col)
	fmt.Fprintln(&str, "}")

	return str.String()
}
