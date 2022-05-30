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

func (p *Parser) consume(message string, types ...token.TokenKind) (bool, error) {
	if p.match(types...) {
		return true, nil
	} else {
		return false, errors.New(message)
	}
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
	return p.varDecl()
}

func (p *Parser) varDecl() (ast.Expr, error) {
	if p.match(token.LET) {
		c, err := p.consume("expect variable name", token.IDENT)
		if !c {
			return nil, err
		}
		name := p.previous()
		c, err = p.consume("expect type anotation", token.COLON)
		if !c {
			return nil, err
		}
		c, err = p.consume("expect type name", token.INTK, token.FLOATK)
		if !c {
			return nil, err
		}
		kind := ast.NewKind(p.previous().Kind)
		var value ast.Expr = ast.NewLiteral("nil", token.NIL)
		if p.match(token.EQ) {
			value, err = p.expression()
			if err != nil {
				return nil, err
			}
		}
		return ast.NewVarDecl(name.Literal, kind, value), nil
	}

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
	if p.match(token.INT, token.FLOAT, token.NIL) {
		return ast.NewLiteral(p.previous().Literal, p.previous().Kind), nil

	}

	if p.match(token.LPAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		c, err := p.consume("expect closing paren", token.RPAREN)
		if !c {
			return nil, err
		}
		return ast.NewGrouping(expr), nil
	}

	return nil, errors.New("expect expression")
}
