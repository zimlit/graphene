/*
	Copyright 2022 Devin Rockwell

	This file is part of Graphene.

	Graphene is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

	Graphene is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

	You should have received a copy of the GNU General Public License along with Graphene. If not, see <https://www.gnu.org/licenses/>.
*/

package parser

import (
	"zimlit/graphene/ast"
	"zimlit/graphene/token"
)

func (p *Parser) advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

func (p *Parser) peek() *token.Token {
	if p.pos >= len(p.tokens) {
		return nil
	}
	return &p.tokens[p.pos]
}

func (p *Parser) previous() *token.Token {
	if p.pos > len(p.tokens) {
		return nil
	} else if len(p.tokens) == 0 {
		return nil
	} else if p.pos == 0 {
		return p.peek()
	}
	return &p.tokens[p.pos-1]
}

func (p *Parser) check(t token.TokenKind) bool {
	if p.pos >= len(p.tokens) {
		return false
	}

	return p.peek().Kind == t
}

func (p *Parser) match(types ...token.TokenKind) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(types ...token.TokenKind) (bool, error) {
	if p.match(types...) {
		return true, nil
	} else {
		t := p.peek()
		if p.previous().Line < t.Line {
			return false, newUnexpectedTokenErr(t, types, p.lines[p.previous().Line-1], p.previous().Line, p.previous().Col+1, p.fname)
		}
		if t == nil {
			return false, newUnexpectedTokenErr(t, types, p.lines[p.previous().Line-1], p.previous().Line, p.previous().Col, p.fname)
		}

		return false, newUnexpectedTokenErr(t, types, p.lines[t.Line-1], t.Line, t.Col, p.fname)
	}
}

func (p *Parser) synchronize() {
	p.advance()

	for p.pos < len(p.tokens) {
		switch p.peek().Kind {
		case token.LET:
			p.advance()
			return
		case token.IF:
			p.advance()
			return
		case token.ELSE:
			p.advance()
			return
		case token.ELSEIF:
			p.advance()
			return
		case token.END:
			p.advance()
			return
		}
		p.advance()
	}
}

func (p *Parser) kind() (ast.ValueKind, error) {
	if p.match(token.INTK) {
		return ast.INT, nil
	} else if p.match(token.FLOATK) {
		return ast.FLOAT, nil
	} else if p.match(token.STRINGK) {
		return ast.STRING, nil
	} else if p.match(token.FN) {
		_, err := p.consume(token.LPAREN)
		if err != nil {
			return nil, err
		}
		params := []ast.Param{}
		kind, err := p.kind()
		if err != nil {
			_, err = p.consume(token.RPAREN)
			if err != nil {
				return nil, err
			}
		}
		params = append(params, ast.NewParam("", kind))
		for p.match(token.COMMA) {
			kind, err := p.kind()
			if err != nil {
				return nil, err
			}
			params = append(params, ast.NewParam("", kind))
		}
		_, err = p.consume(token.RPAREN)
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.COLON)
		if err != nil {
			return nil, err
		}
		kind, err = p.kind()
		if err != nil {
			return nil, err
		}
		return ast.NewFnT(params, kind), nil
	}

	if p.peek() == nil {
		return nil, newUnexpectedTokenErr(p.peek(), []token.TokenKind{token.INTK, token.FLOATK, token.STRINGK, token.FN}, p.lines[p.previous().Line-1], p.previous().Line, p.previous().Col, p.fname)
	}
	return nil, newUnexpectedTokenErr(p.peek(), []token.TokenKind{token.INTK, token.FLOATK, token.STRINGK, token.FN}, p.lines[p.peek().Line-1], p.peek().Line, p.peek().Col, p.fname)
}
