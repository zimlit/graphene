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

func (p *Parser) returnExpr() (ast.Expr, error) {
	if p.match(token.RETURN) {
		value, err := p.expression()
		if err != nil {
			return nil, err
		}
		return ast.NewReturn(value), nil
	}

	return p.whileExpr()
}

func (p *Parser) whileExpr() (ast.Expr, error) {
	if p.match(token.WHILE) {
		cond, err := p.expression()
		if err != nil {
			return nil, err
		}

		var body []ast.Expr
		for {
			b, err := p.expression()
			if p.peek() == nil {
				return nil, newUnexpectedTokenErr(nil, []token.TokenKind{token.END}, p.lines[p.previous().Line-1], p.previous().Line, p.previous().Col, p.fname)
			}
			if err != nil {
				if p.peek().Kind == token.END {
					break
				}

				return nil, err
			}
			body = append(body, b)
		}

		c, err := p.consume(token.END)
		if !c {
			return nil, err
		}

		return ast.NewWhileExpr(cond, body), nil
	}

	return p.ifExpr()
}

func (p *Parser) ifExpr() (ast.Expr, error) {
	if p.match(token.IF) {
		cond, err := p.expression()
		if err != nil {
			return nil, err
		}
		var body []ast.Expr
		for {
			b, err := p.expression()
			if err != nil {
				if p.peek() == nil {
					return nil, newUnexpectedTokenErr(nil, []token.TokenKind{token.END}, p.lines[p.previous().Line-1], p.previous().Line, p.previous().Col, p.fname)
				}
				if p.peek().Kind == token.ELSE || p.peek().Kind == token.ELSEIF || p.peek().Kind == token.END {
					break
				}

				return nil, err
			}
			body = append(body, b)
		}

		var else_ifs []ast.IfExpr
		for p.match(token.ELSEIF) {
			econd, err := p.expression()
			if err != nil {
				return nil, err
			}

			var ebody []ast.Expr
			for {
				e, err := p.expression()
				if p.peek() == nil {
					return nil, newUnexpectedTokenErr(nil, []token.TokenKind{token.END}, p.lines[p.previous().Line-1], p.previous().Line, p.previous().Col, p.fname)
				}
				if err != nil {
					if p.peek().Kind == token.ELSE || p.peek().Kind == token.ELSEIF || p.peek().Kind == token.END {
						break
					}
					return nil, err
				}
				ebody = append(ebody, e)
			}

			else_if := ast.NewIfExpr(econd, ebody, nil, nil)

			else_ifs = append(else_ifs, else_if)
		}

		var el []ast.Expr
		if p.match(token.ELSE) {
			for {
				e, err := p.expression()
				if p.peek() == nil {
					return nil, newUnexpectedTokenErr(nil, []token.TokenKind{token.END}, p.lines[p.previous().Line-1], p.previous().Line, p.previous().Col, p.fname)
				}
				if err != nil {
					if p.peek().Kind == token.ELSE || p.peek().Kind == token.ELSEIF || p.peek().Kind == token.END {
						break
					}
					return nil, err
				}
				el = append(el, e)
			}
		}
		c, err := p.consume(token.END)
		if !c {
			return nil, err
		}

		return ast.NewIfExpr(cond, body, else_ifs, el), nil

	}

	return p.varDecl()
}

func (p *Parser) varDecl() (ast.Expr, error) {
	if p.match(token.LET) {
		is_mut := p.match(token.MUT)
		c, err := p.consume(token.IDENT)
		if !c {
			return nil, err
		}
		name := p.previous()
		c, err = p.consume(token.COLON)
		if !c {
			return nil, err
		}
		kind, err := p.kind()
		if err != nil {
			return nil, err
		}
		var value ast.Expr = ast.NewLiteral("nil", token.NIL)
		if p.match(token.EQ) {
			value, err = p.expression()
			if err != nil {
				return nil, err
			}
		}
		return ast.NewVarDecl(name.Literal, kind, value, is_mut), nil
	}

	return p.fn()
}

func (p *Parser) fn() (ast.Expr, error) {
	if p.match(token.FN) {
		name := ""
		if p.match(token.IDENT) {
			name = p.previous().Literal
		}
		_, err := p.consume(token.LPAREN)
		if err != nil {
			return nil, err
		}
		params := []ast.Param{}
		if p.match(token.IDENT) {
			name := p.previous().Literal
			_, err := p.consume(token.COLON)
			if err != nil {
				return nil, err
			}
			kind, err := p.kind()
			if err != nil {
				return nil, err
			}
			params = append(params, ast.NewParam(name, kind))
			for p.match(token.COMMA) {
				_, err := p.consume(token.IDENT)
				if err != nil {
					return nil, err
				}
				name = p.previous().Literal
				_, err = p.consume(token.COLON)
				if err != nil {
					return nil, err
				}
				kind, err := p.kind()
				if err != nil {
					return nil, err
				}
				params = append(params, ast.NewParam(name, kind))
			}
		}
		_, err = p.consume(token.RPAREN)
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.COLON)
		if err != nil {
			return nil, err
		}
		kind, err := p.kind()
		if err != nil {
			return nil, err
		}

		var body []ast.Expr
		for {
			b, err := p.expression()
			if p.peek() == nil {
				return nil, newUnexpectedTokenErr(nil, []token.TokenKind{token.END}, p.lines[p.previous().Line-1], p.previous().Line, p.previous().Col, p.fname)
			}
			if err != nil {
				if p.peek().Kind == token.END {
					break
				}

				return nil, err
			}
			body = append(body, b)
		}

		c, err := p.consume(token.END)
		if !c {
			return nil, err
		}
		f := ast.NewFn(params, body, kind)
		if name != "" {
			return ast.NewVarDecl(name, ast.NewFnT(params, kind), f, false), nil
		}
		return f, nil
	}

	return p.assignment()
}
