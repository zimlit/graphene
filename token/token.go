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
	IDENT
	COLON
	NIL
)

func (t TokenKind) String() string {
	switch t {
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case STAR:
		return "STAR"
	case SLASH:
		return "SLASH"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case INTK:
		return "INTK"
	case FLOATK:
		return "FLOATK"
	case LET:
		return "LET"
	case EQ:
		return "EQ"
	case IDENT:
		return "IDENT"
	case COLON:
		return "COLON"
	case NIL:
		return "NIL"
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
	fmt.Fprint(&str, "   ", "kind:", t.Kind)
	fmt.Fprintln(&str, ",")
	fmt.Fprint(&str, "   ", "literal:", t.Literal)
	fmt.Fprintln(&str, ",")
	fmt.Fprint(&str, "   ", "line:", t.Line)
	fmt.Fprintln(&str, ",")
	fmt.Fprintln(&str, "   ", "col:", t.Col)
	fmt.Fprintln(&str, "}")

	return str.String()
}
