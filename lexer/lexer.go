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
	"unicode"
	"zimlit/graphene/token"
)

type Lexer struct {
	col      int
	pos      int
	line     int
	lineStr  string
	source   []rune
	fname    string
	keywords map[string]token.TokenKind
}

func NewLexer(source string, fname string) Lexer {
	l := Lexer{
		col:     1,
		pos:     0,
		line:    1,
		lineStr: "",
		source:  []rune(source),
		fname:   fname,
	}
	l.keywords = make(map[string]token.TokenKind)
	l.keywords["int"] = token.INTK
	l.keywords["float"] = token.FLOATK
	l.keywords["let"] = token.LET
	l.keywords["nil"] = token.NIL
	l.keywords["if"] = token.IF
	l.keywords["end"] = token.END
	l.keywords["else"] = token.ELSE
	l.keywords["mut"] = token.MUT
	l.keywords["while"] = token.WHILE
	l.keywords["string"] = token.STRINGK
	l.keywords["fn"] = token.FN
	l.keywords["return"] = token.RETURN

	return l
}

func (l *Lexer) string() (*token.Token, *tmpLexErr) {
	val := ""
	col := l.col
	l.advance()

Exit:
	for ; ; l.advance() {
		if l.pos >= len(l.source) {
			err := l.newTmpErr("Unclosed string")
			return nil, &err
		}
		switch l.peek() {
		case '"':
			break Exit
		case '\000':
			fallthrough
		case '\n':
			err := l.newTmpErr("Unclosed string")
			return nil, &err
		case '\\':
			l.advance()
			switch l.peek() {
			case 't':
				val += "\t"
			case 'n':
				val += "\n"
			case '"':
				val += "\""
			case '\\':
				val += "\\"
			case 'r':
				val += "\r"
			case 'v':
				val += "\v"
			default:
				err := l.newTmpErr("Invalid escape character")
				return nil, &err
			}
		default:
			val += string(l.peek())
		}
	}
	tok := l.newTokenAt(val, token.STRING, col)
	return &tok, nil
}

func (l *Lexer) num() (*token.Token, *tmpLexErr) {
	val := ""
	dot_count := 0
	col := l.col
	for ; l.pos < len(l.source); l.advance() {
		val += string(l.peek())
		if l.peek() == '.' {
			dot_count++
		}

		if !(unicode.IsDigit(l.peekNext()) || l.peekNext() == '.') {
			break
		}
	}

	var tok token.Token

	if dot_count > 1 {
		e := l.newTmpErrAt("to many dots in number literal", col)
		return nil, &e
	} else if dot_count == 1 {
		tok = l.newTokenAt(val, token.FLOAT, col)
	} else {
		tok = l.newTokenAt(val, token.INT, col)
	}

	return &tok, nil
}

func (l *Lexer) ident() token.Token {
	val := ""
	col := l.col

	for ; l.pos < len(l.source); l.advance() {
		val += string(l.peek())

		if !(unicode.IsLetter(l.peekNext()) || unicode.IsDigit(l.peekNext()) || l.peekNext() == '_') {
			break
		}
	}
	if val == "else" {
		if l.source[l.pos+2] == 'i' {
			if l.source[l.pos+3] == 'f' {
				l.advance()
				l.advance()
				l.advance()
				return l.newTokenAt(val, token.ELSEIF, col)
			}
		}
	}
	t := l.keywords[val]
	if t == 0 {
		return l.newTokenAt(val, token.IDENT, col)
	}
	return l.newTokenAt(val, t, col)
}

func (l *Lexer) Lex() ([]token.Token, []string, LexErrs) {
	toks := []token.Token{}
	var errs LexErrs = nil
	tmps := []tmpLexErr{}
	lines := []string{}
	for ; l.pos < len(l.source); l.advance() {
		switch l.peek() {
		case '+':
			toks = append(toks, l.newToken("+", token.PLUS))
		case '-':
			toks = append(toks, l.newToken("-", token.MINUS))
		case '*':
			toks = append(toks, l.newToken("*", token.STAR))
		case '/':
			toks = append(toks, l.newToken("/", token.SLASH))
		case '(':
			toks = append(toks, l.newToken("(", token.LPAREN))
		case ')':
			toks = append(toks, l.newToken(")", token.RPAREN))
		case ',':
			toks = append(toks, l.newToken(",", token.COMMA))
		case '=':
			if l.match('=') {
				toks = append(toks, l.newTokenAt("==", token.EQEQ, l.col-1))
			} else {
				toks = append(toks, l.newToken("=", token.EQ))
			}
		case '!':
			if l.match('=') {
				toks = append(toks, l.newTokenAt("!=", token.NEQ, l.col-1))
			} else {
				toks = append(toks, l.newToken("!", token.BANG))
			}
		case '<':
			if l.match('=') {
				toks = append(toks, l.newTokenAt("<=", token.LESSEQ, l.col-1))
			} else {
				toks = append(toks, l.newToken("<", token.LESS))
			}
		case '>':
			if l.match('=') {
				toks = append(toks, l.newTokenAt(">=", token.GREATEREQ, l.col-1))
			} else {
				toks = append(toks, l.newToken(">", token.GREATER))
			}
		case ':':
			toks = append(toks, l.newToken(":", token.COLON))
		case '"':
			t, err := l.string()
			if err != nil {
				tmps = append(tmps, *err)
			} else {
				toks = append(toks, *t)
			}
		case ' ':
		case '\t':
		case '\r':
		case '\v':
		case '\n':
			l.lineStr += "\n"
			for _, err := range tmps {
				errs = append(errs, l.newLexErr(err))
			}
			lines = append(lines, l.lineStr)
			l.line++
			l.lineStr = ""
			tmps = []tmpLexErr{}
			l.col = 1
		default:
			if unicode.IsDigit(l.peek()) {
				t, err := l.num()

				if err != nil {
					tmps = append(tmps, *err)
				} else {
					toks = append(toks, *t)
				}

			} else if unicode.IsLetter(l.peek()) || l.peek() == '_' {
				t := l.ident()
				toks = append(toks, t)
			} else {
				tmps = append(tmps, l.newTmpErr(fmt.Sprintf("Unexpected character '%s'", string(l.peek()))))

			}

		}
	}
	for _, err := range tmps {
		errs = append(errs, l.newLexErr(err))
	}
	if l.lineStr != "" {
		lines = append(lines, l.lineStr)
	}

	if errs != nil {
		return nil, nil, errs
	}
	return toks, lines, nil
}
