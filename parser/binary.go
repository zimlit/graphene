package parser

import (
	"zimlit/graphene/ast"
	"zimlit/graphene/token"
)

func (p *Parser) assignment() (ast.Expr, error) {
	if p.match(token.IDENT) {
		ident := p.previous()
		c, _ := p.consume(token.EQ)
		if !c {
			p.pos--
			return p.equality()
		}
		val, err := p.expression()
		if err != nil {
			return nil, err
		}

		return ast.NewAssignment(ident.Literal, val), nil
	}

	return p.equality()
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(token.EQEQ, token.NEQ) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinary(expr, *operator, right)
	}

	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.LESS, token.LESSEQ, token.GREATER, token.GREATEREQ) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinary(expr, *operator, right)
	}

	return expr, nil
}

func (p *Parser) term() (ast.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinary(expr, *operator, right)
	}

	return expr, nil
}

func (p *Parser) factor() (ast.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = ast.NewBinary(expr, *operator, right)
	}

	return expr, nil
}
