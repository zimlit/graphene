package parser

import (
	"fmt"
	"strings"
	"zimlit/graphene/ast"
	"zimlit/graphene/token"

	"github.com/fatih/color"
)

type ParseResult struct {
	Exprs ast.Exprs
	Err   ParseError
}

type ParseError []error

func (p ParseError) Error() string {
	var str strings.Builder

	for _, err := range p {
		fmt.Fprintln(&str, err.Error())
	}

	return str.String()
}

type MsgErr struct {
	msg     string
	line    int
	col     int
	lineStr string
	fname   string
}

func (m MsgErr) Error() string {
	var str strings.Builder
	r := color.New(color.FgHiRed, color.Bold).FprintfFunc()
	w := color.New(color.FgHiWhite, color.Bold).FprintfFunc()
	b := color.New(color.FgHiBlue, color.Bold).FprintfFunc()

	r(&str, "error")
	fmt.Fprint(&str, ": ")
	w(&str, "%s\n", m.msg)
	b(&str, " --> ")
	fmt.Fprintf(&str, "%s:%d:%d\n", m.fname, m.line, m.col)
	b(&str, "  |\n")
	b(&str, "%d | ", m.line)
	fmt.Fprintln(&str, m.lineStr)
	b(&str, "  |")
	for i := 0; i < m.col; i++ {
		fmt.Fprint(&str, " ")
	}
	r(&str, "^ %s\n", m.msg)

	return str.String()
}

func newMsgErr(msg string, line int, col int, lineStr string, fname string) MsgErr {
	return MsgErr{
		msg:     msg,
		line:    line,
		col:     col,
		lineStr: lineStr,
		fname:   fname,
	}
}

type UnexpectedTokenErr struct {
	got      *token.Token
	expected []token.TokenKind
	line     int
	col      int
	lineStr  string
	fname    string
}

func (u UnexpectedTokenErr) Error() string {
	var err strings.Builder
	var str strings.Builder
	r := color.New(color.FgHiRed, color.Bold).FprintfFunc()
	w := color.New(color.FgHiWhite, color.Bold).FprintfFunc()
	b := color.New(color.FgHiBlue, color.Bold).FprintfFunc()
	fmt.Fprint(&err, "Unexpected token expected ")

	for i, expected := range u.expected {
		switch expected {
		case token.INT:
			fmt.Fprint(&err, expected)
		case token.NIL:
			fmt.Fprint(&err, expected)
		case token.FLOAT:
			fmt.Fprint(&err, expected)
		default:
			fmt.Fprintf(&err, "\"%s\"", expected.String())
		}
		if i != len(u.expected)-1 {
			fmt.Fprint(&err, " or ")
		} else {
			fmt.Fprint(&err, " got ")
		}
	}
	if u.got == nil {
		fmt.Fprintf(&err, "EOF")
	} else {
		switch u.got.Kind {
		case token.INT:
			fmt.Fprint(&err, u.got.Kind)
		case token.NIL:
			fmt.Fprint(&err, u.got.Kind)
		case token.FLOAT:
			fmt.Fprint(&err, u.got.Kind)
		default:
			fmt.Fprintf(&err, "\"%s\"", u.got.Kind.String())
		}
	}

	r(&str, "error")
	fmt.Fprint(&str, ": ")
	w(&str, "%s\n", err.String())
	b(&str, " --> ")
	fmt.Fprintf(&str, "%s:%d:%d\n", u.fname, u.line, u.col)
	b(&str, "  |\n")
	b(&str, "%d | ", u.line)
	fmt.Fprint(&str, u.lineStr)
	if u.lineStr[len(u.lineStr)-1] != '\n' {
		fmt.Fprintf(&str, "\n")
	}
	b(&str, "  |")
	if u.got == nil {
		for i := -1; i < u.col; i++ {
			fmt.Fprint(&str, " ")
		}
	} else {
		for i := 0; i < u.col; i++ {
			fmt.Fprint(&str, " ")
		}
	}

	r(&str, "^ %s\n", err.String())

	return str.String()
}

func newUnexpectedTokenErr(got *token.Token, expected []token.TokenKind, lineStr string, line int, col int, fname string) UnexpectedTokenErr {
	return UnexpectedTokenErr{
		got:      got,
		expected: expected,
		lineStr:  lineStr,
		line:     line,
		col:      col,
		fname:    fname,
	}
}
