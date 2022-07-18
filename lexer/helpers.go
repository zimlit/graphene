/*
	Copyright 2022 Devin Rockwell

	This file is part of Graphene.

	Graphene is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

	Graphene is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

	You should have received a copy of the GNU General Public License along with Graphene. If not, see <https://www.gnu.org/licenses/>.
*/

package lexer

import "zimlit/graphene/token"

func (l *Lexer) newTmpErr(msg string) tmpLexErr {
	return tmpLexErr{
		col:   l.col,
		line:  l.line,
		msg:   msg,
		fname: l.fname,
	}
}

func (l *Lexer) newTmpErrAt(msg string, col int) tmpLexErr {
	return tmpLexErr{
		col:   col,
		line:  l.line,
		msg:   msg,
		fname: l.fname,
	}
}

func (l *Lexer) newLexErr(tmp tmpLexErr) LexErr {
	return LexErr{
		col:     tmp.col,
		line:    tmp.line,
		msg:     tmp.msg,
		lineStr: l.lineStr,
		fname:   tmp.fname,
	}
}

func (l *Lexer) newToken(literal string, kind token.TokenKind) token.Token {
	return token.Token{
		Kind:    kind,
		Literal: literal,
		Line:    l.line,
		Col:     l.col,
	}
}

func (l *Lexer) newTokenAt(literal string, kind token.TokenKind, col int) token.Token {
	return token.Token{
		Kind:    kind,
		Literal: literal,
		Line:    l.line,
		Col:     col,
	}
}

func (l *Lexer) peek() rune {
	if l.pos < len(l.source) {
		return l.source[l.pos]
	}
	return '\000'
}

func (l *Lexer) peekNext() rune {
	if l.pos < len(l.source)-1 {
		return l.source[l.pos+1]
	}
	return '\000'
}

func (l *Lexer) advance() {
	if l.pos < len(l.source) {
		l.lineStr += string(l.source[l.pos])

	}
	l.pos++
	l.col++
}

func (l *Lexer) match(c rune) bool {
	if l.peekNext() == c {
		l.advance()
		return true
	}

	return false
}
