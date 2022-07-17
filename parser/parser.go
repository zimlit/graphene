package parser

import (
	"fmt"
	"zimlit/graphene/ast"
	"zimlit/graphene/token"
)

type Parser struct {
	tokens []token.Token
	pos    int
	lines  []string
	fname  string
}

func NewParser(tokens []token.Token, lines []string, fname string) Parser {
	return Parser{
		tokens: tokens,
		pos:    0,
		lines:  lines,
		fname:  fname,
	}
}

func (p *Parser) Parse(c chan ParseResult) {
	exprs := []ast.Expr{}
	var errs ParseError = nil

	for p.pos < len(p.tokens) {
		expr, err := p.expression()
		if err != nil {
			errs = append(errs, err)
			p.synchronize()
		}
		exprs = append(exprs, expr)
	}

	if errs != nil {
		c <- ParseResult{nil, errs}
	}

	c <- ParseResult{exprs, nil}
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.whileExpr()
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(token.MINUS, token.BANG) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return ast.NewUnary(*operator, right), nil
	}

	return p.primary()
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(token.INT, token.FLOAT, token.NIL, token.IDENT) {
		return ast.NewLiteral(p.previous().Literal, p.previous().Kind), nil
	}
	if p.match(token.STRING) {
		return ast.NewLiteral(fmt.Sprintf("\"%s\"", p.previous().Literal), p.previous().Kind), nil
	}

	if p.match(token.LPAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		c, err := p.consume(token.RPAREN)
		if !c {
			return nil, err
		}
		return ast.NewGrouping(expr), nil
	}

	if p.peek() == nil {
		if p.previous() == nil {
			return nil, newMsgErr("Expected expression", 1, 1, "", p.fname)
		}
		return nil, newMsgErr("Expected expression", p.previous().Line, p.previous().Col+1, p.lines[p.previous().Line-1], p.fname)
	} else {
		return nil, newMsgErr("Expected expression", p.peek().Line, p.peek().Col, p.lines[p.peek().Line-1], p.fname)
	}

}
