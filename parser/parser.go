package parser

import (
	"errors"
	"zimlit/graphene/ast"
	"zimlit/graphene/token"
)

type Parser struct {
	tokens []token.Token
	pos    int
	lines  []string
}

func NewParser(tokens []token.Token, lines []string) Parser {
	return Parser{
		tokens: tokens,
		pos:    0,
		lines:  lines,
	}
}

func (p *Parser) advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.pos]
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.pos-1]
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

func (p *Parser) synchronize() {
	p.advance()

	for p.pos < len(p.tokens) {
		switch p.peek().Kind {

		}
		p.advance()
	}
}

func (p *Parser) Parse() (ast.Expr, error) {
	return p.expression()
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.term()
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
		expr = ast.NewBinary(expr, operator, right)
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
		expr = ast.NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(token.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return ast.NewUnary(operator, right), nil
	}

	return p.primary()
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(token.INT, token.FLOAT) {
		return ast.NewLiteral(p.previous().Literal), nil
	}

	if p.match(token.LPAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if !p.match(token.RPAREN) {
			return nil, errors.New("expect ')' after expression")
		}
		return ast.NewGrouping(expr), nil
	}

	return nil, errors.New("expect expression")
}
