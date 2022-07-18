/*
	Copyright 2022 Devin Rockwell

	This file is part of Graphene.

	Graphene is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

	Graphene is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

	You should have received a copy of the GNU General Public License along with Graphene. If not, see <https://www.gnu.org/licenses/>.
*/
package lexer

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type tmpLexErr struct {
	col   int
	line  int
	msg   string
	fname string
}

func (l *tmpLexErr) Error() string { return "" }

type LexErr struct {
	col     int
	line    int
	lineStr string
	msg     string
	fname   string
}

func (l *LexErr) Error() string {
	r := color.New(color.FgHiRed, color.Bold).FprintfFunc()
	w := color.New(color.FgHiWhite, color.Bold).FprintfFunc()
	b := color.New(color.FgHiBlue, color.Bold).FprintfFunc()
	var str strings.Builder
	r(&str, "error")
	fmt.Fprint(&str, ": ")
	w(&str, "%s\n", l.msg)
	b(&str, " --> ")
	fmt.Fprintf(&str, "%s:%d:%d\n", l.fname, l.line, l.col)
	b(&str, "  |\n")
	b(&str, "%d | ", l.line)
	if l.lineStr[0] == '\n' {
		fmt.Fprint(&str, l.lineStr[1:])
		l.col--
	} else {
		fmt.Fprint(&str, l.lineStr)
	}
	if l.lineStr[len(l.lineStr)-1] != '\n' {
		fmt.Fprint(&str, "\n")
	}

	b(&str, "  |")
	for i := 0; i < l.col; i++ {
		fmt.Fprint(&str, " ")
	}
	r(&str, "^ %s\n", l.msg)

	return str.String()
}

type LexErrs []LexErr

func (l *LexErrs) Error() string {
	var str strings.Builder
	for _, err := range *l {
		fmt.Fprintln(&str, err.Error())
	}

	return str.String()
}
