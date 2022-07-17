/*
	Copyright 2022 Devin Rockwell

	This file is part of Graphene.

	Graphene is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

	Graphene is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

	You should have received a copy of the GNU General Public License along with Graphene. If not, see <https://www.gnu.org/licenses/>.
*/

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
	if m.lineStr[0] == '\n' {
		fmt.Fprint(&str, m.lineStr[1:])
		m.col--
	} else {
		fmt.Fprint(&str, m.lineStr)
	}
	if m.lineStr[len(m.lineStr)-1] != '\n' {
		fmt.Fprint(&str, "\n")
	}
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
		u.line++
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
	fmt.Fprintf(&str, "%s:%d:", u.fname, u.line)
	if u.got == nil {
		fmt.Fprintln(&str, "0")
	} else {
		fmt.Fprintln(&str, u.col)
	}
	b(&str, "  |\n")
	b(&str, "%d | ", u.line)
	if u.got != nil {
		if u.lineStr[0] == '\n' {
			fmt.Fprint(&str, u.lineStr[1:])
			u.col--
		} else {
			fmt.Fprint(&str, u.lineStr)
		}
	} else {
		fmt.Fprintln(&str, "")
	}

	b(&str, "  |")
	if u.got == nil {
		fmt.Fprint(&str, " ")
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
