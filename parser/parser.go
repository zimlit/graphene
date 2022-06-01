package parser

import (
	"fmt"
	"strings"
	"zimlit/graphene/ast"
	"zimlit/graphene/token"

	"github.com/fatih/color"
)

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
	fmt.Fprintln(&str, m.lineStr)
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
	fmt.Fprintf(&str, "%s:%d:%d\n", u.fname, u.line, u.col)
	b(&str, "  |\n")
	b(&str, "%d | ", u.line)
	fmt.Fprintln(&str, u.lineStr)
	b(&str, "  |")
	if u.got == nil {
		for i := -1; i < u.col; i++ {
			fmt.Fprint(&str, " ")
		}
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
		if t == nil {
			return false, newUnexpectedTokenErr(p.peek(), types, p.lines[p.previous().Line-1], p.previous().Line, p.previous().Col, p.fname)
		}
		return false, newUnexpectedTokenErr(p.peek(), types, p.lines[p.peek().Line-1], p.peek().Line, p.peek().Col, p.fname)
	}
}

func (p *Parser) synchronize() {
	p.advance()

	for p.pos < len(p.tokens) {
		switch p.peek().Kind {
		case token.LET:
			return
		case token.IF:
			return
		case token.ELSE:
			return
		case token.ELSEIF:
			return
		case token.END:
			return
		}
		p.advance()
	}
}

func (p *Parser) Parse() ([]ast.Expr, error) {
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
		return nil, errs
	}

	return exprs, nil
}

func (p *Parser) expression() (ast.Expr, error) {
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
		c, err := p.consume(token.IDENT)
		if !c {
			return nil, err
		}
		name := p.previous()
		c, err = p.consume(token.COLON)
		if !c {
			return nil, err
		}
		c, err = p.consume(token.INTK, token.FLOATK)
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

	return p.comparison()
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.EQEQ, token.NEQ, token.LESS, token.LESSEQ, token.GREATER, token.GREATEREQ) {
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
