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
